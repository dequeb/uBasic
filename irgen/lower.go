package irgen

// TODO: Add convenience functions for creating instruction in emit.go, to
// remove if err != nil { panic("foo") } from the irgen code.

import (
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"uBasic/ast"
	"uBasic/ast/astutil"
	"uBasic/errors"
	"uBasic/eval"
	"uBasic/object"
	"uBasic/sem"
	"uBasic/token"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	irtypes "github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func GenToFile(file *ast.File, info *sem.Info, filename string) error {
	// replace extension with .ll
	filename = strings.TrimSuffix(filename, ".bas")
	if !strings.HasSuffix(filename, ".ll") {
		filename += ".ll"
	}
	// create file
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	fmt.Fprint(f, gen(file, info, filename).String())

	// run compiler
	cmd := exec.Command("lli", filename)
	if err = cmd.Run(); err != nil {
		return err
	}

	return nil
}

// Gen generates LLVM IR based on the syntax tree of the given file.
func Gen(file *ast.File, info *sem.Info, filename string) *ir.Module {
	return gen(file, info, filename)
}

// === [ File scope ] ==========================================================

// gen generates LLVM IR based on the syntax tree of the given file.
func gen(file *ast.File, info *sem.Info, filename string) *ir.Module {
	m := NewModule(info)
	m.SourceFilename = filename
	m.genExternals()
	m.genErrorHandler()

	// process function and sub declarations
	for _, statementList := range file.Body {
		for _, statement := range statementList.Statements {
			switch statement := statement.(type) {
			case *ast.DimDecl:
				for _, vars := range statement.Vars {
					m.globalVarDecl(vars)
				}
			case *ast.ConstDecl:
				for _, cnst := range statement.Consts {
					m.newGlobalConstant(&cnst)
				}
			case *ast.FuncDecl:
				m.funcDecl(statement)
			case *ast.SubDecl:
				m.subDecl(statement)
			}
		}
	}
	// local variables of main function were declared globally
	m.SkipLocalVariables = true

	// process main function
	params := []*ir.Param{} // no command line arguments
	f := NewFunc("main", irtypes.I32, params...)

	// Generate function body.
	dbg.Printf("create function definition: main")
	mainParams := []ast.ParamItem{}
	m.funcBody(f, mainParams, file.Body, nil)
	return m.Module
}

// --- [ Function declaration ] ------------------------------------------------

// funcDecl lowers the given function declaration to LLVM IR, emitting code to
// m.
func (m *Module) funcDecl(n *ast.FuncDecl) {
	// Generate function signature.
	ident := n.Name()
	name := ident.String()
	uBasicType, err := ast.NewType(n.FuncType)
	if err != nil {
		panic(fmt.Sprintf("unable to create type; %v", err))
	}
	typ := toIrType(uBasicType)
	sig, ok := typ.(*irtypes.FuncType)
	if !ok {
		panic(fmt.Sprintf("invalid function type; expected *irtypes.FuncType, got %T", typ))
	}
	var params []*ir.Param
	for i, p := range n.FuncType.Params {
		paramType := sig.Params[i]
		param := ir.NewParam(p.Name().String(), paramType)
		params = append(params, param)
	}
	f := NewFunc(name, sig.RetType, params...)

	// Emit Variadic attribute: check if last parameter is ParamArray
	if len(n.FuncType.Params) > 0 {
		if n.FuncType.Params[len(n.FuncType.Params)-1].ParamArray {
			f.Sig.Variadic = true
		}
	}

	if !astutil.IsDef(n) {
		dbg.Printf("create function declaration: %v", n)
		// Emit function declaration.
		m.emitFunc(f)
		return
	}
	m.setIdentValue(ident, f.Func)

	// Generate function body.
	dbg.Printf("create function definition: %v", n)
	m.funcBody(f, n.FuncType.Params, n.Body, n.FuncName)
}

// subDecl lowers the given sub declaration to LLVM IR, emitting code to m.
func (m *Module) subDecl(n *ast.SubDecl) {
	// Generate function signature.
	ident := n.Name()
	name := ident.String()
	typ := irtypes.Void
	sig := irtypes.FuncType{RetType: typ}
	var params []*ir.Param
	for i, p := range n.SubType.Params {
		paramType := sig.Params[i]
		param := ir.NewParam(p.Name().String(), paramType)
		params = append(params, param)
	}
	f := NewFunc(name, sig.RetType, params...)
	// Emit Variadic attribute: check if last parameter is ParamArray
	if len(n.SubType.Params) > 0 {
		if n.SubType.Params[len(n.SubType.Params)-1].ParamArray {
			f.Sig.Variadic = true
		}
	}

	if !astutil.IsDef(n) {
		dbg.Printf("create subroutine declaration: %v", n)
		// Emit function declaration.
		m.emitFunc(f)
		return
	}
	m.setIdentValue(ident, f.Func)

	// Generate function body.
	dbg.Printf("create subroutine definition: %v", n)
	m.funcBody(f, n.SubType.Params, n.Body, nil)
}

// funcBody lowers the given function declaration to LLVM IR, emitting code to
// m.
func (m *Module) funcBody(f *Function, params []ast.ParamItem, body []ast.StatementList, resultIdentifier *ast.Identifier) {
	// Initialize function body.
	f.startBody()

	// main calls _main to intercept errors
	if f.Name() == "main" {
		// main noParams
		noParams := []*ir.Param{}
		innerF := NewFunc(".main", irtypes.Void, noParams...)

		// main function
		exception := f.NewBlock("exception")
		normalCall := f.NewBlock("normalCall")
		end := f.NewBlock("end")

		// entry:
		tmp2 := f.currentBlock.NewCall(m.LookupFunction("setjmp"), m.LookupGlobal(JumpBuffer))
		cmp := f.currentBlock.NewICmp(enum.IPredEQ, tmp2, constant.NewInt(types.I32, 0))
		f.currentBlock.NewCondBr(cmp, normalCall.Block, exception.Block)

		// normalCall:
		f.Blocks = append(f.Blocks, f.currentBlock.Block)
		f.currentBlock = normalCall
		normalCall.NewCall(innerF)
		normalCall.NewBr(end.Block)

		// exception:
		// em2 = exception.NewLoad(types.I8Ptr, errorMessage)
		f.Blocks = append(f.Blocks, f.currentBlock.Block)
		f.currentBlock = exception
		en := f.currentBlock.NewLoad(types.I32, m.LookupGlobal(ErrorNumber))
		f.currentBlock.NewCall(m.LookupFunction("printf"), m.LookupGlobal(ErrorMessage))
		f.currentBlock.NewRet(en)

		// end:
		f.Blocks = append(f.Blocks, f.currentBlock.Block)
		f.currentBlock = end

		if err := m.endBody(f, resultIdentifier); err != nil {
			panic(fmt.Sprintf("unable to finalize function body; %v", err))
		}
		// Emit function definition.
		m.emitFunc(f)

		// call inner function.
		astParams := []ast.ParamItem{}
		m.funcBody(innerF, astParams, body, nil)

		return
	} else {
		// Emit local variable declarations for function parameters.
		for i, param := range f.Params {
			p := m.funcParam(f, param)
			dbg.Printf("create function parameter: %v", params[i])
			ident := params[i].Name()
			got := f.genUnique(ident)
			if ident.Name != got {
				panic(fmt.Sprintf("unable to generate identical function parameter name; expected %q, got %q", ident, got))
			}
			f.setIdentValue(ident, p)
		}
		// Emit local variable declaration for result identifier.
		if resultIdentifier != nil {
			result := f.currentBlock.NewAlloca(f.Sig.RetType)
			f.currentBlock.NewStore(constZero(f.Sig.RetType), result)
			dbg.Printf("create result identifier: %v", resultIdentifier)
			f.setIdentValue(resultIdentifier, result)
		}
		// Generate function body.
		dbg.Printf("create function definition: " + f.Name())
		m.BodyStmt(f, body)

		// Finalize function body.
		if err := m.endBody(f, resultIdentifier); err != nil {
			panic(fmt.Sprintf("unable to finalize function body; %v", err))
		}

		// Emit function definition.
		m.emitFunc(f)
	}
}

// funcParam lowers the given function parameter to LLVM IR, emitting code to f.
func (m *Module) funcParam(f *Function, param *ir.Param) value.Value {
	// Input:
	//    void f(int a) {
	//    }
	// Output:
	//    %1 = alloca i32
	//    store i32 %a, i32* %1
	addr := f.currentBlock.NewAlloca(param.Type())
	f.currentBlock.NewStore(param, addr)
	return addr
}

// --- [ Global variable declaration ] -----------------------------------------

// globalVarDecl lowers the given global variable declaration to LLVM IR,
// emitting code to m.
func (m *Module) globalVarDecl(n ast.VarDecl) {

	switch n := n.(type) {
	case *ast.ScalarDecl:
		m.globalScalarDecl(n)
	case *ast.ArrayDecl:
		m.globalArrayDecl(n)
	default:
		panic(fmt.Sprintf("support for global variable declaration %T not yet implemented", n))
	}
}

func (m *Module) globalScalarDecl(n *ast.ScalarDecl) {
	// Input:
	//    int x;
	// Output:
	//    @x = global i32 0
	ident := n.Name()
	dbg.Printf("create global variable: %v", n)
	typ0, err := n.Type()
	if err != nil {
		panic(fmt.Sprintf("unable to create type; %v", err))
	}
	typ := toIrType(typ0)
	var val constant.Constant
	if intType, ok := typ.(*irtypes.IntType); ok {
		val = constant.NewInt(intType, 0)
	} else if floatType, ok := typ.(*irtypes.FloatType); ok {
		val = constant.NewFloat(floatType, 0)
	} else if ptrType, ok := typ.(*irtypes.PointerType); ok {
		val = constant.NewNull(ptrType)
	} else {
		val = constant.NewZeroInitializer(typ)
	}
	global := ir.NewGlobalDef(ident.Name, val)
	m.setIdentValue(ident, global)
	// Emit global variable definition.
	m.emitGlobal(global)
}

// globalArrayDecl lowers the given global array declaration to LLVM IR, emitting code to m.
func (m *Module) globalArrayDecl(n *ast.ArrayDecl) {
	// Input:
	//    int x[3];
	// Output:
	//    @x = global [3 x i32] zeroinitializer
	ident := n.Name()
	dbg.Printf("create global array: %v", n)
	typ0, err := n.Type()
	if err != nil {
		panic(fmt.Sprintf("unable to create type; %v", err))
	}
	// get array type
	arrayTyp := toIrType(typ0)
	dimensions := n.VarType.Dimensions
	if len(dimensions) == 0 {
		//m.glocalDynamicArrayDecl(n)
		panic("dynamic arrays not yet implemented")
	}

	// for multi-dimensional arrays, create an array of single-dimension
	// we will multiply the dimensions to get the total size of the array

	// calculate the total size of the array
	env := object.NewEnvironment()
	size := int64(1)
	for i := len(dimensions) - 1; i >= 0; i-- {
		result := eval.Eval(nil, dimensions[i], env)
		var dimension int64
		switch result.(type) {
		case *object.Long:
			dimension = result.GetValue().(int64)
			size *= dimension
		default:
			panic("unknown expression: " + result.String())
		}

		// allocate length of array constant
		constName := fmt.Sprintf(".%s_%d", ident.Name, i)
		cnst := m.NewGlobalDef(constName, constant.NewInt(irtypes.I32, dimension))
		cnst.Immutable = true

	}

	// Create a new global variable of type [15]i8 and name it "str".
	array0 := constant.NewArray(&irtypes.ArrayType{Len: uint64(size), ElemType: arrayTyp})
	init := constant.NewZeroInitializer(array0.Typ)
	global := ir.NewGlobalDef(ident.Name, init)

	m.setIdentValue(ident, global)
	// Emit global variable definition.
	m.emitGlobal(global)

}

// --- [ Global constant declaration ] -----------------------------------------

func (m Module) newGlobalConstant(node *ast.ConstDeclItem) {
	// all constants are evaluated at compile time
	Object := eval.Eval(nil, node.ConstValue, m.env)
	m.env.Set(node.ConstName.Name, Object)

	value := Object.String()
	if Object.Type() == object.ERROR_OBJ {
		panic("error evaluating constant : " + Object.String())
	}

	typ := strings.ToLower(node.ConstType.Token().Literal)
	switch typ {
	case "integer":
		val, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			panic("error converting string to int")
		}
		cnst := m.NewGlobalDef(node.ConstName.Name, constant.NewInt(irtypes.I32, val))
		cnst.Immutable = true
		m.setIdentValue(node.ConstName, cnst)
	case "long":
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			panic("error converting string to long")
		}
		cnst := m.NewGlobalDef(node.ConstName.Name, constant.NewInt(irtypes.I64, val))
		cnst.Immutable = true
		m.setIdentValue(node.ConstName, cnst)
	case "single", "currency":
		val, err := strconv.ParseFloat(value, 32)
		if err != nil {
			panic("error converting string to float")
		}
		cnst := m.NewGlobalDef(node.ConstName.Name, constant.NewFloat(irtypes.Float, val))
		cnst.Immutable = true
		m.setIdentValue(node.ConstName, cnst)
	case "double":
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			panic("error converting string to double")
		}
		cnst := m.NewGlobalDef(node.ConstName.Name, constant.NewFloat(irtypes.Double, val))
		cnst.Immutable = true
		m.setIdentValue(node.ConstName, cnst)
	case "boolean":
		cnst := m.NewGlobalDef(node.ConstName.Name, constant.NewBool(strings.EqualFold(value, "true")))
		cnst.Immutable = true
		m.setIdentValue(node.ConstName, cnst)
	case "string":
		m.newGlobalStringConstant(value, node.ConstName.Name)
	case "date":
		floatDate := FromDateStringToFloat(value)
		cnst := m.NewGlobalDef(node.ConstName.Name, constant.NewFloat(irtypes.Double, floatDate))
		cnst.Immutable = true
		m.setIdentValue(node.ConstName, cnst)
	default:
		panic("unknown type")
	}

}

var globalCounter int

func (m *Module) newGlobalStringConstant(val, name string) *ir.Global {
	basicText := m.cleanString(val)
	if name == "" {
		globalCounter++
		name = fmt.Sprintf("_%d", globalCounter)
	}
	text := constant.NewCharArrayFromString(basicText + "\x00")
	value := m.NewGlobalDef(name, text)
	value.Immutable = true
	return value
}

func (m *Module) cleanString(s string) string {
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}

	s = strings.Replace(s, "\"\"", "\"", -1)
	s = strings.Replace(s, "\x00", "", -1)
	return s
}

// === [ Function scope ] ======================================================

// --- [ Local variable definition ] -------------------------------------------

// localVarDef lowers the given local variable definition to LLVM IR, emitting
// code to f.
func (m *Module) localVarDef(f *Function, n ast.VarDecl) value.Value {
	// Input:
	//    void f() {
	//       int a;           // <-- relevant line
	//    }
	// Output:
	//    %a = alloca i32
	ident := n.Name()
	dbg.Printf("create local variable: %v", n)
	typ0, err := n.Type()
	if err != nil {
		panic(fmt.Sprintf("unable to create type; %v", err))
	}
	typ := toIrType(typ0)
	allocaInst := f.currentBlock.NewAlloca(typ)
	// Emit local variable definition.
	return f.emitLocal(ident, allocaInst)
}

// constDecl lowers the given constant declaration to LLVM IR, emitting code to
// f.
func (m *Module) localConstDecl(f *Function, cnst *ast.ConstDeclItem) {
	switch cnst.ConstValue.(type) {
	case *ast.BasicLit:
		switch cnst.ConstValue.(*ast.BasicLit).Kind {
		case token.LongLit:
			// %a = alloca i32
			a := f.currentBlock.NewAlloca(irtypes.I64)
			a.SetName(cnst.ConstName.Name)
			// get value
			env := object.NewEnvironment()
			Object := eval.Eval(nil, cnst.ConstValue, env)
			var Value int64
			switch Object.(type) {
			case *object.Long:
				Value = Object.GetValue().(int64)
			default:
				panic("unknown expression: " + Object.String())
			}

			// store i32 32, i32* %
			f.currentBlock.NewStore(constant.NewInt(irtypes.I64, Value), a)
		}
	}

}

// --- [ Statements ] ----------------------------------------------------------

// stmt lowers the given statement to LLVM IR, emitting code to f.
func (m *Module) stmt(f *Function, stmt ast.Statement) {
	switch stmt := stmt.(type) {
	// case *ast.BlockStmt:
	// 	m.blockStmt(f, stmt)
	// 	return
	case *ast.EmptyStmt, *ast.Comment, *ast.SubDecl, *ast.FuncDecl:
		// nothing to do.
		// for function and sub declarations, it was done during global scope
		// @see `gen` function
	case *ast.ExprStmt:
		m.exprStmt(f, stmt)
	case *ast.IfStmt:
		m.ifStmt(f, stmt)
	case *ast.WhileStmt:
		m.whileStmt(f, stmt)
	case *ast.CallSubStmt:
		m.callSubStmt(f, stmt)
	case *ast.SpecialStmt:
		m.specialStmt(f, stmt)
	case *ast.DimDecl:
		if !m.SkipLocalVariables {
			for _, vars := range stmt.Vars {
				m.localVarDef(f, vars)
			}
		}
	case *ast.ConstDecl:
		if !m.SkipLocalVariables {
			for _, cnst := range stmt.Consts {
				m.localConstDecl(f, &cnst)
			}
		}
	default:
		panic(fmt.Sprintf("support for %T not yet implemented", stmt))
	}
}

// BodyStmt lowers the given body statement to LLVM IR, emitting code to f.
func (m *Module) BodyStmt(f *Function, body []ast.StatementList) {
	for _, statementList := range body {
		m.StatementList(f, &statementList)
	}
}

// StatementList lowers the given block statement to LLVM IR, emitting code to f.
func (m *Module) StatementList(f *Function, stmtList *ast.StatementList) {
	for _, statement := range stmtList.Statements {
		m.stmt(f, statement)
	}
}

// exprStmt lowers the given expression statement to LLVM IR, emitting code to
// f.
func (m *Module) exprStmt(f *Function, stmt *ast.ExprStmt) {
	m.expr(f, stmt.Expression)
}

// callSubStmt lowers the given call statement to LLVM IR, emitting code to f.
func (m *Module) callSubStmt(f *Function, stmt *ast.CallSubStmt) {
	m.expr(f, stmt.Definition)
}

// ifStmt lowers the given if statement to LLVM IR, emitting code to f.
func (m *Module) ifStmt(f *Function, stmt *ast.IfStmt) {
	cond := m.cond(f, stmt.Condition)
	trueBranch := f.NewBlock("")
	end := f.NewBlock("")
	falseBranch := end
	if stmt.Else != nil {
		falseBranch = f.NewBlock("")
	}
	termCondBr := ir.NewCondBr(cond, trueBranch.Block, falseBranch.Block)
	f.currentBlock.SetTerm(termCondBr)
	f.currentBlock = trueBranch
	m.BodyStmt(f, stmt.Body)
	// Emit jump if body doesn't end with return statement (i.e. the current
	// basic block is none nil).
	if f.currentBlock != nil {
		termBr := ir.NewBr(end.Block)
		f.currentBlock.SetTerm(termBr)
	}
	if stmt.Else != nil {
		f.currentBlock = falseBranch
		m.BodyStmt(f, stmt.Else)
		// Emit jump if body doesn't end with return statement (i.e. the current
		// basic block is none nil).
		if f.currentBlock != nil {
			termBr := ir.NewBr(end.Block)
			f.currentBlock.SetTerm(termBr)
		}
	}
	f.currentBlock = end
}

// // exitStmt lowers the given return statement to LLVM IR, emitting code to f.
// func (m *Module) exitStmt(f *Function, stmt *ast.ExitStmt) {
// 	// Input:
// 	//    int f() {
// 	//       return 42;       // <-- relevant line
// 	//    }
// 	// Output:
// 	//    ret i32 42
// 	if stmt.Result == nil {
// 		termRet := ir.NewRet(nil)
// 		f.curBlock.SetTerm(termRet)
// 		f.curBlock = nil
// 		return
// 	}
// 	result := m.expr(f, stmt.Result)
// 	// Implicit conversion.
// 	resultType := f.Sig.RetType
// 	result = m.convert(f, result, resultType)
// 	termRet := ir.NewRet(result)
// 	f.curBlock.SetTerm(termRet)
// 	f.curBlock = nil
// }

// whileStmt lowers the given while statement to LLVM IR, emitting code to f.
func (m *Module) whileStmt(f *Function, stmt *ast.WhileStmt) {
	condBranch := f.NewBlock("")
	termBr := ir.NewBr(condBranch.Block)
	f.currentBlock.SetTerm(termBr)
	f.currentBlock = condBranch
	cond := m.cond(f, stmt.Condition)
	bodyBranch := f.NewBlock("")
	endBranch := f.NewBlock("")
	termCondBr := ir.NewCondBr(cond, bodyBranch.Block, endBranch.Block)
	f.currentBlock.SetTerm(termCondBr)
	f.currentBlock = bodyBranch
	m.BodyStmt(f, stmt.Body)
	// Emit jump if body doesn't end with return statement (i.e. the current
	// basic block is none nil).
	if f.currentBlock != nil {
		termBr := ir.NewBr(condBranch.Block)
		f.currentBlock.SetTerm(termBr)
	}
	f.currentBlock = endBranch
}

// --- [ Expressions ] ----------------------------------------------------------

// cond lowers the given condition expression to LLVM IR, emitting code to f.
func (m *Module) cond(f *Function, expr ast.Expression) value.Value {
	cond := m.expr(f, expr)
	if cond.Type().Equal(irtypes.I1) {
		return cond
	}
	// Create boolean expression if cond is not already of boolean type.
	//
	//    cond != 0
	// zero is the integer constant 0.
	zero := constZero(cond.Type())
	return f.currentBlock.NewICmp(enum.IPredNE, cond, zero)
}

// expr lowers the given expression to LLVM IR, emitting code to f.
func (m *Module) expr(f *Function, expr ast.Expression) value.Value {
	switch expr := expr.(type) {
	case *ast.BasicLit:
		return m.basicLit(expr)
	case *ast.BinaryExpr:
		return m.binaryExpr(f, expr)
	case *ast.CallOrIndexExpr:
		return m.callOrIndexExpr(f, expr)
	case *ast.Identifier:
		return m.identUse(f, expr)
	case *ast.ParenExpr:
		return m.expr(f, expr.Expr)
	case *ast.UnaryExpr:
		return m.unaryExpr(f, expr)
	default:
		panic(fmt.Sprintf("support for type %T not yet implemented", expr))
	}
}

// basicLit lowers the given basic literal to LLVM IR, emitting code to f.
func (m *Module) basicLit(n *ast.BasicLit) value.Value {
	typ := m.typeOf(n)
	switch n.Kind {
	case token.CurrencyLit:
		floatType, ok := typ.(*irtypes.FloatType)
		floatType.Kind = irtypes.FloatKindFloat
		if !ok {
			panic(fmt.Errorf("invalid currency literal type; expected *irtypes.FloatType, got %T", typ))
		}
		// convert currency to float
		val := n.Value.(string)
		// remove $ from the currency
		val = strings.TrimSuffix(val, "$")
		c, err := constant.NewFloatFromString(floatType, val)
		if err != nil {
			panic(fmt.Errorf("unable to parse float literal %q; %v", val, err))
		}
		return c

	case token.LongLit:
		intType, ok := typ.(*irtypes.IntType)
		intType.BitSize = 64
		if !ok {
			panic(fmt.Errorf("invalid integer literal type; expected *irtypes.IntType, got %T", typ))
		}
		c, err := constant.NewIntFromString(intType, n.Value.(string))
		if err != nil {
			panic(fmt.Errorf("unable to parse integer literal %q; %v", n.Value, err))
		}
		return c
	case token.DoubleLit:
		floatType, ok := typ.(*irtypes.FloatType)
		floatType.Kind = irtypes.FloatKindDouble
		if !ok {
			panic(fmt.Errorf("invalid float literal type; expected *irtypes.FloatType, got %T", typ))
		}
		c, err := constant.NewFloatFromString(floatType, n.Value.(string))
		if err != nil {
			panic(errors.Newf(n.Token().Position, "unable to parse float literal %q; %v", n.Value, err))
		}
		return c
	case token.StringLit:
		return m.newGlobalStringConstant(n.Value.(string), "")
	case token.KwTrue:
		return constant.NewBool(true)
	case token.KwFalse:
		return constant.NewBool(false)
	case token.KwNothing:
		return constant.NewNull(irtypes.I1Ptr)
	case token.DateLit:
		floatType, ok := typ.(*irtypes.FloatType)
		floatType.Kind = irtypes.FloatKindDouble
		if !ok {
			panic(fmt.Errorf("invalid integer literal type; expected *irtypes.IntType, got %T", typ))
		}
		// convert date to float
		numSecs := FromDateStringToFloat(n.Value.(string))

		c := &constant.Float{
			Typ: floatType,
			X:   big.NewFloat(numSecs),
		}
		return c

	default:
		panic(fmt.Sprintf("support for basic literal kind %v not yet implemented", n.Kind))
	}
}

// binaryExpr lowers the given binary expression to LLVM IR, emitting code to f.
func (m *Module) binaryExpr(f *Function, n *ast.BinaryExpr) value.Value {
	switch n.OpToken.Kind {
	// +
	case token.Add:
		x, y := m.expr(f, n.Left), m.expr(f, n.Right)
		x, y = m.implicitConversion(f, x, y)
		switch x.Type().String() {
		case "i32", "i64":
			return f.currentBlock.NewAdd(x, y)
		case "float", "double":
			return f.currentBlock.NewFAdd(x, y)
		}
	// -
	case token.Minus:
		x, y := m.expr(f, n.Left), m.expr(f, n.Right)
		x, y = m.implicitConversion(f, x, y)
		switch x.Type().String() {
		case "i32", "i64":
			return f.currentBlock.NewSub(x, y)
		case "float", "double":
			return f.currentBlock.NewFSub(x, y)
		}

	// *
	case token.Mul:
		x, y := m.expr(f, n.Left), m.expr(f, n.Right)
		x, y = m.implicitConversion(f, x, y)
		switch x.Type().String() {
		case "i32", "i64":
			return f.currentBlock.NewMul(x, y)
		case "float", "double":
			return f.currentBlock.NewFMul(x, y)
		}

	// /
	case token.Div:
		x, y := m.expr(f, n.Left), m.expr(f, n.Right)
		x, y = m.implicitConversion(f, x, y)
		m.checkIfDivisionByZero(f, y)
		switch x.Type().String() {
		case "i32", "i64":
			return f.currentBlock.NewSDiv(x, y)
		case "float", "double":
			return f.currentBlock.NewFDiv(x, y)
		}

	// <
	case token.Lt:
		x, y := m.expr(f, n.Left), m.expr(f, n.Right)
		x, y = m.implicitConversion(f, x, y)
		switch x.Type().String() {
		case "i32", "i64":
			return f.currentBlock.NewICmp(enum.IPredSLT, x, y)
		case "float", "double":
			return f.currentBlock.NewFCmp(enum.FPredOLT, x, y)
		}

		// >
	case token.Gt:
		x, y := m.expr(f, n.Left), m.expr(f, n.Right)
		x, y = m.implicitConversion(f, x, y)
		switch x.Type().String() {
		case "i32", "i64":
			return f.currentBlock.NewICmp(enum.IPredSGT, x, y)
		case "float", "double":
			return f.currentBlock.NewFCmp(enum.FPredOGT, x, y)
		}

	// <=
	case token.Le:
		x, y := m.expr(f, n.Left), m.expr(f, n.Right)
		x, y = m.implicitConversion(f, x, y)
		switch x.Type().String() {
		case "i32", "i64":
			return f.currentBlock.NewICmp(enum.IPredSLE, x, y)
		case "float", "double":
			return f.currentBlock.NewFCmp(enum.FPredOLE, x, y)
		}

	// >=
	case token.Ge:
		x, y := m.expr(f, n.Left), m.expr(f, n.Right)
		x, y = m.implicitConversion(f, x, y)
		switch x.Type().String() {
		case "i32", "i64":
			return f.currentBlock.NewICmp(enum.IPredSGE, x, y)
		case "float", "double":
			return f.currentBlock.NewFCmp(enum.FPredOGE, x, y)
		}

	// <>
	case token.Neq:
		x, y := m.expr(f, n.Left), m.expr(f, n.Right)
		x, y = m.implicitConversion(f, x, y)
		switch x.Type().String() {
		case "i32", "i64":
			return f.currentBlock.NewICmp(enum.IPredNE, x, y)
		case "float", "double":
			return f.currentBlock.NewFCmp(enum.FPredONE, x, y)
		}
	// ==
	case token.Eq:
		x, y := m.expr(f, n.Left), m.expr(f, n.Right)
		x, y = m.implicitConversion(f, x, y)
		switch x.Type().String() {
		case "i32", "i64":
			return f.currentBlock.NewICmp(enum.IPredEQ, x, y)
		case "float", "double":
			return f.currentBlock.NewFCmp(enum.FPredOEQ, x, y)
		}

	// And
	case token.And:
		x := m.cond(f, n.Left)

		start := f.currentBlock
		trueBranch := f.NewBlock("")
		end := f.NewBlock("")
		termCondBr := ir.NewCondBr(x, trueBranch.Block, end.Block)
		f.currentBlock.SetTerm(termCondBr)
		f.currentBlock = trueBranch

		y := m.cond(f, n.Right)
		termBr := ir.NewBr(end.Block)
		trueBranch.SetTerm(termBr)
		f.currentBlock = end

		var incs []*ir.Incoming
		zero := constZero(irtypes.I1)
		inc := ir.NewIncoming(zero, start.Block)
		incs = append(incs, inc)
		inc = ir.NewIncoming(y, trueBranch.Block)
		incs = append(incs, inc)
		return f.currentBlock.NewPhi(incs...)
	// or
	case token.Or:
		x := m.cond(f, n.Left)

		start := f.currentBlock
		falseBranch := f.NewBlock("")
		end := f.NewBlock("")
		termCondBr := ir.NewCondBr(x, end.Block, falseBranch.Block)
		f.currentBlock.SetTerm(termCondBr)
		f.currentBlock = falseBranch

		y := m.cond(f, n.Right)
		termBr := ir.NewBr(end.Block)
		falseBranch.SetTerm(termBr)
		f.currentBlock = end

		var incs []*ir.Incoming
		one := constOne(irtypes.I1)
		inc := ir.NewIncoming(one, start.Block)
		incs = append(incs, inc)
		inc = ir.NewIncoming(y, falseBranch.Block)
		incs = append(incs, inc)
		return f.currentBlock.NewPhi(incs...)
	// =
	case token.Assign:
		right := m.expr(f, n.Right)
		switch left := n.Left.(type) {
		case *ast.Identifier:
			m.identDef(f, left, right)
			// case *ast.CallOrIndexExpr:
			// 	m.indexExprDef(f, expr, y)
		default:
			panic(fmt.Sprintf("support for assignment to type %T not yet implemented", left))
		}
		return right
	case token.Concat:
		// -----------------------------------------------------------
		// gc := m.LookupGlobal("gc")
		// -----------------------------------------------------------

		// allocate heap memory for intermediate strings
		// calculate length of string
		x := m.expr(f, n.Left)
		y := m.expr(f, n.Right)
		lengthX := f.currentBlock.NewCall(m.LookupFunction("strlen"), x)
		lengthY := f.currentBlock.NewCall(m.LookupFunction("strlen"), y)
		length := f.currentBlock.NewAdd(lengthX, lengthY)
		// -----------------------------------------------------------
		// memoryBlock := f.currentBlock.NewCall(m.LookupFunction("gc_malloc"), gc, length)
		memoryBlock := f.currentBlock.NewCall(m.LookupFunction("malloc"), length)
		// -----------------------------------------------------------

		f.currentBlock.NewCall(m.LookupFunction("strcpy"), memoryBlock, x)
		f.currentBlock.NewCall(m.LookupFunction("strcat"), memoryBlock, y)
		return memoryBlock
	default:
		panic(fmt.Sprintf("support for binary operator %v not yet implemented", n.OpToken))
	}
	panic("unreachable")
}

// callExpr lowers the given identifier to LLVM IR, emitting code to f.
func (m *Module) callExpr(f *Function, callOrIndexExpr *ast.CallOrIndexExpr) value.Value {
	typ0, err := callOrIndexExpr.Identifier.Decl.Type()
	if err != nil {
		panic(fmt.Sprintf("unable to create type; %v", err))
	}
	typ := toIrType(typ0)
	sig, ok := typ.(*irtypes.FuncType)
	if !ok {
		panic(fmt.Sprintf("invalid function type; expected *irtypes.FuncType, got %T", typ))
	}
	params := sig.Params
	result := sig.RetType
	_ = result
	var args []value.Value
	for i, arg := range callOrIndexExpr.Args {
		expr := m.expr(f, arg)
		expr = m.convert(f, expr, params[i])
		args = append(args, expr)
	}
	if len(callOrIndexExpr.Args) < len(params) {
		// add optional parameters
		var funcDecl *ast.FuncDecl
		var subDecl *ast.SubDecl
		var ok bool
		var calleeParams []ast.ParamItem
		if funcDecl, ok = callOrIndexExpr.Identifier.Decl.(*ast.FuncDecl); ok {
			calleeParams = funcDecl.FuncType.Params
		} else if subDecl, ok = callOrIndexExpr.Identifier.Decl.(*ast.SubDecl); ok {
			calleeParams = subDecl.SubType.Params
		}
		for i := len(callOrIndexExpr.Args); i < len(calleeParams); i++ {
			expr := m.expr(f, calleeParams[i].DefaultValue)
			expr = m.convert(f, expr, params[i])
			args = append(args, expr)
		}
	}

	v := m.valueFromIdent(f, callOrIndexExpr.Identifier)
	callee, ok := v.(*ir.Func)
	if !ok {
		panic(fmt.Sprintf("invalid callee type; expected *ir.Func, got %T", v))
	}

	return f.currentBlock.NewCall(callee, args...)
}

// ident lowers the given identifier to LLVM IR, emitting code to f.
func (m *Module) ident(f *Function, ident *ast.Identifier) value.Value {
	// if variable is a local variable of the name of the function
	if ident.Name == f.Name() {
		pos := ident.Decl.Name().Tok.Position.Absolute
		return f.idents[pos]
	}

	switch typ := m.typeOf(ident).(type) {
	case *irtypes.ArrayType:
		array := m.valueFromIdent(f, ident)
		arrayElemType := array.Type().(*irtypes.PointerType).ElemType
		zero := constZero(irtypes.I64)
		indices := []value.Value{zero, zero}

		// Emit getelementptr instruction.
		if m.isGlobal(ident) {
			var is []constant.Constant
			for _, index := range indices {
				i, ok := index.(constant.Constant)
				if !ok {
					break
				}
				is = append(is, i)
			}
			if len(is) == len(indices) {
				// In accordance with Clang, emit getelementptr constant expressions
				// for global variables.
				// TODO: Validate typ against array.
				_ = typ
				if array, ok := array.(constant.Constant); ok {
					return constant.NewGetElementPtr(arrayElemType, array, is...)
				}
				panic(fmt.Sprintf("invalid constant array type; expected constant.Constant, got %T", array))
			}
		}
		return f.currentBlock.NewGetElementPtr(arrayElemType, array, indices...)
	case *irtypes.PointerType:
		// Emit load instruction.
		// TODO: Validate typ against srcAddr.Elem().
		src := m.valueFromIdent(f, ident)
		srcElemType := src.Type().(*irtypes.PointerType).ElemType
		return f.currentBlock.NewLoad(srcElemType, src)
	default:
		return m.valueFromIdent(f, ident)
	}
}

// identUse lowers the given identifier usage to LLVM IR, emitting code to f.
func (m *Module) identUse(f *Function, ident *ast.Identifier) value.Value {
	v := m.ident(f, ident)
	typ := m.typeOf(ident)
	if isRef(typ) {
		return v
	}
	// TODO: Validate typ against v.Elem()
	elemType := v.Type().(*irtypes.PointerType).ElemType
	return f.currentBlock.NewLoad(elemType, v)
}

// identDef lowers the given identifier definition to LLVM IR, emitting code to f.
func (m *Module) identDef(f *Function, ident *ast.Identifier, v value.Value) {
	// i8 := irtypes.I8
	// i8ptr := irtypes.NewPointer(i8)

	addr := m.ident(f, ident)
	addrType, ok := addr.Type().(*irtypes.PointerType)
	if !ok {
		panic(fmt.Sprintf("invalid pointer type; expected *irtypes.PointerType, got %T", addr.Type()))
	}
	v = m.convert(f, v, addrType.ElemType)
	// if string we need to allocate memory for it
	t, _ := ident.Decl.Type()
	if t.String() == "String" {
		// find base memory address
		dest := m.valueFromIdent(f, ident)
		// dest contains the address of the destination variable
		// addr contains the loaded address of the old destination variable
		// newAddr contains the loaded address of the new destination variable
		// v contains the loaded address of the source variable

		// calculate length of string
		length := f.currentBlock.NewCall(m.LookupFunction("strlen"), v)

		// -----------------------------------------------------------
		// allocate heap memory for global strings
		// gc := m.LookupGlobal("gc")
		// memoryBlock := f.currentBlock.NewCall(m.LookupFunction("gc_malloc"), gc, length)
		memoryBlock := f.currentBlock.NewCall(m.LookupFunction("malloc"), length)
		// -----------------------------------------------------------

		f.currentBlock.NewStore(memoryBlock, dest)
		f.currentBlock.NewCall(m.LookupFunction("strcpy"), memoryBlock, v)
	} else {
		f.currentBlock.NewStore(v, addr)
	}
}

// unaryExpr lowers the given unary expression to LLVM IR, emitting code to f.
func (m *Module) unaryExpr(f *Function, n *ast.UnaryExpr) value.Value {
	switch n.OpToken.Kind {
	// -expr
	case token.Minus:
		// Input:
		//    void f() {
		//       int x;
		//       -x;              // <-- relevant line
		//    }
		// Output:
		//    %2 = sub i32 0, %1
		expr := m.expr(f, n.Right)
		zero := constZero(expr.Type())
		return f.currentBlock.NewSub(zero, expr)
	// Not expr
	case token.Not:
		cond := m.cond(f, n.Right)
		one := constOne(cond.Type())
		notCond := f.currentBlock.NewXor(cond, one)
		return notCond
		// return f.currentBlock.NewZExt(notCond, m.typeOf(n.Right))
	default:
		panic(fmt.Sprintf("support for unary operator %v not yet implemented", n.OpToken))
	}
}

// callOrIndexExpr lowers the given call or index expression to LLVM IR, emitting
// code to f.
func (m *Module) callOrIndexExpr(f *Function, n *ast.CallOrIndexExpr) value.Value {
	switch n.Identifier.Decl.(type) {
	case *ast.FuncDecl, *ast.SubDecl:
		return m.callExpr(f, n)
	case *ast.ArrayDecl:
		return m.indexExpr(f, n)
	default:
		panic(fmt.Sprintf("support for type %T not yet implemented", n))
	}
}

// indexExpr lowers the given index expression to LLVM IR, emitting code to f.
func (m *Module) indexExpr(f *Function, n *ast.CallOrIndexExpr) value.Value {

	// evaluate dimensions at comopile time
	astDimensions := n.Identifier.Decl.(*ast.ArrayDecl).VarType.Dimensions
	dimensions := make([]int64, len(astDimensions))
	for i, dim := range astDimensions {
		result := eval.Eval(nil, dim, object.NewEnvironment())
		switch resultObj := result.(type) {
		case *object.Long:
			dimensions[i] = resultObj.Value
		default:
			panic("error evaluating array dimensions")
		}
	}
	indices := make([]value.Value, len(n.Args))
	var compoundIndex value.Value

	// verify that the number of indices is equal to the number of dimensions of the array
	// or one less if the array is a dynamic array

	if len(n.Args) == len(dimensions) || (len(n.Args) == 1 && len(dimensions) == 0) {
		for i, index := range n.Args {
			indices[i] = m.expr(f, index)
			// load array dimension
			dimVarName := fmt.Sprintf("%s_dim_%d", n.Identifier.Name, i)
			dimVar := m.LookupGlobal(dimVarName)
			dim := f.currentBlock.NewLoad(irtypes.I64, dimVar)
			// compare index with dimension
			m.checkArrayBounds(f, indices[i], dim)

			// calculate compound index
			if i == 0 {
				compoundIndex = indices[i]
			} else {
				compoundIndex = f.currentBlock.NewMul(compoundIndex, dim)
				compoundIndex = f.currentBlock.NewAdd(compoundIndex, indices[i])
			}
		}
	} else {
		panic("invalid number of indices") // should have been caught by the parser
	}
	// get array address
	array := m.valueFromIdent(f, n.Identifier)
	// get array element type
	arrayElemType := array.Type().(*irtypes.PointerType).ElemType
	// calculate address of array element
	// Emit getelementptr instruction.
	return f.currentBlock.NewGetElementPtr(arrayElemType, array, compoundIndex)
}
