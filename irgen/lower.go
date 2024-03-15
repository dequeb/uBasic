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
	fmt.Fprint(f, Gen(file, info).String())

	// run compiler
	cmd := exec.Command("llc", "-filetype=obj", filename)
	if err = cmd.Run(); err != nil {
		return err
	}

	// run linker
	cmd = exec.Command(("clang"), "-o", strings.TrimSuffix(filename, ".ll"), strings.TrimSuffix(filename, ".ll")+".o")
	if err = cmd.Run(); err != nil {
		return err
	}

	return nil
}

// Gen generates LLVM IR based on the syntax tree of the given file.
func Gen(file *ast.File, info *sem.Info) *ir.Module {
	return gen(file, info)
}

// === [ File scope ] ==========================================================

// gen generates LLVM IR based on the syntax tree of the given file.
func gen(file *ast.File, info *sem.Info) *ir.Module {
	m := NewModule(info)
	m.genExternals()
	// m.genInternals()

	// process function and sub declarations
	for _, statementList := range file.StatementLists {
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
	// main params
	I8PtrPtr := types.NewPointer(types.I8Ptr)
	params := []*ir.Param{ir.NewParam("argc", types.I32), ir.NewParam("argv", I8PtrPtr)}
	f := NewFunc("main", irtypes.I32, params...)
	dbg.Printf("create function declaration: %v", "main")

	// Generate function body.
	dbg.Printf("create function definition: main")
	mainParams := []ast.ParamItem{}
	m.funcBody(f, mainParams, file.StatementLists)
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
		panic(fmt.Sprintf("invalid function type; expected *types.FuncType, got %T", typ))
	}
	var params []*ir.Param
	for i, p := range n.FuncType.Params {
		paramType := sig.Params[i]
		param := ir.NewParam(p.Name().String(), paramType)
		params = append(params, param)
	}
	f := NewFunc(name, sig.RetType, params...)
	if !astutil.IsDef(n) {
		dbg.Printf("create function declaration: %v", n)
		// Emit function declaration.
		m.emitFunc(f)
		return
	}
	m.setIdentValue(ident, f.Func)

	// Generate function body.
	dbg.Printf("create function definition: %v", n)
	m.funcBody(f, n.FuncType.Params, n.Body)
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
	if !astutil.IsDef(n) {
		dbg.Printf("create subroutine declaration: %v", n)
		// Emit function declaration.
		m.emitFunc(f)
		return
	}
	m.setIdentValue(ident, f.Func)

	// Generate function body.
	dbg.Printf("create subroutine definition: %v", n)
	m.funcBody(f, n.SubType.Params, n.Body)
}

// funcBody lowers the given function declaration to LLVM IR, emitting code to
// m.
func (m *Module) funcBody(f *Function, params []ast.ParamItem, body []ast.StatementList) {
	// Initialize function body.
	f.startBody()

	if f.Name() == "main" {
		param := f.Params[0]
		tmp0 := f.currentBlock.NewAlloca(param.Type())
		f.currentBlock.NewStore(param, tmp0)
		dbg.Printf("create function parameter: %v", param)

		// initialize garbage collector
		gc_start := m.LookupFunction(".gc_start")
		gc := m.LookupGlobal(".gc")
		f.currentBlock.NewCall(gc_start, gc, tmp0)
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
	}

	// Generate function body.
	m.BodyStmt(f, body)

	// Finalize function body.
	if err := m.endBody(f); err != nil {
		panic(fmt.Sprintf("unable to finalize function body; %v", err))
	}

	// Emit function definition.
	m.emitFunc(f)
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

// --- [ Global constant declaration ] -----------------------------------------

func (m Module) newGlobalConstant(node *ast.ConstDeclItem) {
	env := object.NewEnvironment()
	Object := eval.Eval(nil, node.ConstValue, env)
	value := Object.String()

	typ := strings.ToLower(node.ConstType.Token().Literal)
	switch typ {
	case "integer":
		val, err := strconv.ParseInt(value, 10, 32)
		if err != nil {
			panic("error converting string to int")
		}
		cnst := m.NewGlobalDef(node.ConstName.Name, constant.NewInt(types.I32, val))
		cnst.Immutable = true
		m.setIdentValue(node.ConstName, cnst)
	case "long":
		val, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			panic("error converting string to long")
		}
		cnst := m.NewGlobalDef(node.ConstName.Name, constant.NewInt(types.I64, val))
		cnst.Immutable = true
		m.setIdentValue(node.ConstName, cnst)
	case "single", "currency":
		val, err := strconv.ParseFloat(value, 32)
		if err != nil {
			panic("error converting string to float")
		}
		cnst := m.NewGlobalDef(node.ConstName.Name, constant.NewFloat(types.Float, val))
		cnst.Immutable = true
		m.setIdentValue(node.ConstName, cnst)
	case "double":
		val, err := strconv.ParseFloat(value, 64)
		if err != nil {
			panic("error converting string to double")
		}
		cnst := m.NewGlobalDef(node.ConstName.Name, constant.NewFloat(types.Double, val))
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
		cnst := m.NewGlobalDef(node.ConstName.Name, constant.NewFloat(types.Double, floatDate))
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
	basicText := strings.Trim(s, "\"")
	basicText = strings.Replace(basicText, "\"\"", "\"", -1)
	basicText = strings.Replace(basicText, "\x00", "", -1)
	return basicText
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
			a := f.currentBlock.NewAlloca(types.I64)
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
			f.currentBlock.NewStore(constant.NewInt(types.I64, Value), a)
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
			panic(fmt.Errorf("invalid currency literal type; expected *types.FloatType, got %T", typ))
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
			panic(fmt.Errorf("invalid integer literal type; expected *types.IntType, got %T", typ))
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
			panic(fmt.Errorf("invalid float literal type; expected *types.FloatType, got %T", typ))
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
			panic(fmt.Errorf("invalid integer literal type; expected *types.IntType, got %T", typ))
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
		return f.currentBlock.NewAdd(x, y)

	// -
	case token.Minus:
		x, y := m.expr(f, n.Left), m.expr(f, n.Right)
		x, y = m.implicitConversion(f, x, y)
		return f.currentBlock.NewSub(x, y)

	// *
	case token.Mul:
		x, y := m.expr(f, n.Left), m.expr(f, n.Right)
		x, y = m.implicitConversion(f, x, y)
		return f.currentBlock.NewMul(x, y)

	// /
	case token.Div:
		x, y := m.expr(f, n.Left), m.expr(f, n.Right)
		x, y = m.implicitConversion(f, x, y)
		return f.currentBlock.NewSDiv(x, y)

	// <
	case token.Lt:
		x, y := m.expr(f, n.Left), m.expr(f, n.Right)
		x, y = m.implicitConversion(f, x, y)
		return f.currentBlock.NewICmp(enum.IPredSLT, x, y)

	// >
	case token.Gt:
		x, y := m.expr(f, n.Left), m.expr(f, n.Right)
		x, y = m.implicitConversion(f, x, y)
		return f.currentBlock.NewICmp(enum.IPredSGT, x, y)

	// <=
	case token.Le:
		x, y := m.expr(f, n.Left), m.expr(f, n.Right)
		x, y = m.implicitConversion(f, x, y)
		return f.currentBlock.NewICmp(enum.IPredSLE, x, y)

	// >=
	case token.Ge:
		x, y := m.expr(f, n.Left), m.expr(f, n.Right)
		x, y = m.implicitConversion(f, x, y)
		return f.currentBlock.NewICmp(enum.IPredSGE, x, y)

	// <>
	case token.Neq:
		x, y := m.expr(f, n.Left), m.expr(f, n.Right)
		x, y = m.implicitConversion(f, x, y)
		return f.currentBlock.NewICmp(enum.IPredNE, x, y)

	// ==
	case token.Eq:
		x, y := m.expr(f, n.Left), m.expr(f, n.Right)
		x, y = m.implicitConversion(f, x, y)
		return f.currentBlock.NewICmp(enum.IPredEQ, x, y)

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

	// =
	case token.Assign:
		y := m.expr(f, n.Right)
		switch expr := n.Left.(type) {
		case *ast.Identifier:
			m.identDef(f, expr, y)
		// case *ast.CallOrIndexExpr:
		// 	m.indexExprDef(f, expr, y)
		default:
			panic(fmt.Sprintf("support for assignment to type %T not yet implemented", expr))
		}
		return y

	default:
		panic(fmt.Sprintf("support for binary operator %v not yet implemented", n.OpToken))
	}
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
		panic(fmt.Sprintf("invalid function type; expected *types.FuncType, got %T", typ))
	}
	params := sig.Params
	result := sig.RetType
	// TODO: Validate result against function return type.
	_ = result
	var args []value.Value
	// TODO: Add support for variadic arguments.
	for i, arg := range callOrIndexExpr.Args {
		expr := m.expr(f, arg)
		expr = m.convert(f, expr, params[i])
		args = append(args, expr)
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

// identDef lowers the given identifier definition to LLVM IR, emitting code to
// f.
func (m *Module) identDef(f *Function, ident *ast.Identifier, v value.Value) {
	// i8 := types.I8
	// i8ptr := types.NewPointer(i8)

	addr := m.ident(f, ident)
	addrType, ok := addr.Type().(*irtypes.PointerType)
	if !ok {
		panic(fmt.Sprintf("invalid pointer type; expected *types.PointerType, got %T", addr.Type()))
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
		length := f.currentBlock.NewCall(m.LookupFunction(".strlen"), v)

		// allocate heap memory for global strings
		gc := m.LookupGlobal(".gc")
		memoryBlock := f.currentBlock.NewCall(m.LookupFunction(".gc_malloc"), gc, length)
		f.currentBlock.NewStore(memoryBlock, dest)
		f.currentBlock.NewCall(m.LookupFunction(".strcpy"), memoryBlock, v)
	} else {
		f.currentBlock.NewStore(v, addr)
	}
}

// // indexExpr lowers the given index expression to LLVM IR, emitting code to f.
// func (m *Module) indexExpr(f *Func, n *ast.IndexExpr) value.Value {
// 	index := m.expr(f, n.Index)
// 	// Extend the index to a 64-bit integer.
// 	if !irtypes.Equal(index.Type(), irtypes.I64) {
// 		index = m.convert(f, index, irtypes.I64)
// 	}
// 	typ := m.typeOf(n.Name)
// 	array := m.valueFromIdent(f, n.Name)

// 	// Dereference pointer pointer.
// 	elem := typ
// 	addr := array
// 	zero := constZero(irtypes.I64)
// 	indices := []value.Value{zero, index}
// 	if typ, ok := typ.(*irtypes.PointerType); ok {
// 		elem = typ.ElemType

// 		// Emit load instruction.
// 		// TODO: Validate typ against array.Elem().
// 		arrayElemType := array.Type().(*irtypes.PointerType).ElemType
// 		addr = f.curBlock.NewLoad(arrayElemType, array)
// 		indices = []value.Value{index}
// 	}

// 	// Emit getelementptr instruction.
// 	addrElemType := addr.Type().(*irtypes.PointerType).ElemType
// 	if m.isGlobal(n.Name) {
// 		var is []constant.Constant
// 		for _, index := range indices {
// 			i, ok := index.(constant.Constant)
// 			if !ok {
// 				break
// 			}
// 			is = append(is, i)
// 		}
// 		if len(is) == len(indices) {
// 			// In accordance with Clang, emit getelementptr constant expressions
// 			// for global variables.
// 			// TODO: Validate elem against addr.
// 			_ = elem
// 			if addr, ok := addr.(constant.Constant); ok {
// 				return constant.NewGetElementPtr(addrElemType, addr, is...)
// 			}
// 			panic(fmt.Sprintf("invalid constant address type; expected constant.Constant, got %T", addr))
// 		}
// 	}
// 	// TODO: Validate elem against array.Elem().
// 	return f.curBlock.NewGetElementPtr(addrElemType, addr, indices...)
// }

// // indexExprUse lowers the given index expression usage to LLVM IR, emitting
// // code to f.
// func (m *Module) indexExprUse(f *Func, n *ast.IndexExpr) value.Value {
// 	v := m.indexExpr(f, n)
// 	typ := m.typeOf(n)
// 	if isRef(typ) {
// 		return v
// 	}
// 	// TODO: Validate typ against v.Elem().
// 	elemType := v.Type().(*irtypes.PointerType).ElemType
// 	return f.curBlock.NewLoad(elemType, v)
// }

// // indexExprDef lowers the given identifier expression definition to LLVM IR,
// // emitting code to f.
// func (m *Module) indexExprDef(f *Func, n *ast.IndexExpr, v value.Value) {
// 	addr := m.indexExpr(f, n)
// 	addrType, ok := addr.Type().(*irtypes.PointerType)
// 	if !ok {
// 		panic(fmt.Sprintf("invalid pointer type; expected *types.PointerType, got %T", addr.Type()))
// 	}
// 	v = m.convert(f, v, addrType.ElemType)
// 	f.curBlock.NewStore(v, addr)
// }

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
		// TODO: Replace `(x != 0) ^ 1` with `x == 0`. Using the former for now to
		// simplify test cases, as they are generated by Clang.

		// Input:
		//    int g() {
		//       int y;
		//       not y;              // <-- relevant line
		//    }
		// Output:
		//    %2 = icmp ne i32 %1, 0
		//    %3 = xor i1 %2, true
		cond := m.cond(f, n.Right)
		one := constOne(cond.Type())
		notCond := f.currentBlock.NewXor(cond, one)
		return f.currentBlock.NewZExt(notCond, m.typeOf(n.Right))
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
	//TODO: implement arrays
	// case *ast.ArrayDecl:
	// 	 return m.indexExpr(f, n)
	default:
		panic(fmt.Sprintf("support for type %T not yet implemented", n))
	}
}
