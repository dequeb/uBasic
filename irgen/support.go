package irgen

import (
	"fmt"
	"strings"
	"time"

	"uBasic/ast"
	"uBasic/object"
	uBasictypes "uBasic/types"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// compileSpecialStmt compiles a special statement.
func (m *Module) specialStmt(f *Function, node *ast.SpecialStmt) {
	switch strings.ToLower(node.Keyword1.Literal) {
	case "print", "debug.print", "msgbox":
		m.printStmt(f, node)
	default:
		panic("unknown special statement")
	}
}

func FromDateStringToFloat(date string) float64 {
	targetDate, err := object.FromStringToTime(date)
	if err != nil {
		panic(err)
	}
	return FromDateToFloat(targetDate)
}

func FromFloatToDateString(f float64) string {
	targetDate := FromFloatToDate(f)
	return object.FromTimeToString(targetDate)
}

func FromDateToFloat(date time.Time) float64 {
	ref := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
	return date.Sub(ref).Seconds()
}

func FromFloatToDate(f float64) time.Time {
	ref := time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)
	return ref.Add(time.Duration(f) * time.Second)
}

// compilePrintStmt compiles a print statement.
func (m *Module) printStmt(f *Function, node *ast.SpecialStmt) {
	zero := constant.NewInt(types.I32, 0)
	printf := m.LookupFunction("printf")
	puts := m.LookupFunction("puts")
	for _, arg := range node.Args {
		switch arg := arg.(type) {
		case *ast.BasicLit:
			switch argValue := arg.Value.(type) {
			case string:
				str0 := m.cleanString(arg.Value.(string))
				argType, _ := arg.Type()
				if argType.(*uBasictypes.Basic).Kind == uBasictypes.Date {
					// remove # from the date
					str0 = strings.Trim(str0, "#")
				}
				str1 := constant.NewCharArrayFromString(str0 + "\x00")
				str2 := f.currentBlock.NewAlloca(str1.Typ)
				f.currentBlock.NewStore(str1, str2)
				f.currentBlock.NewCall(printf, str2)
			case bool:
				var value *ir.Global
				if argValue {
					value = m.LookupGlobal("true")
				} else {
					value = m.LookupGlobal("false")
				}
				typ := value.Typ.ElemType
				gep := constant.NewGetElementPtr(typ, value, zero, zero)
				tmp1 := f.currentBlock.NewLoad(typ, gep)
				f.currentBlock.NewCall(puts, tmp1)
			default:
				panic("unknown basic literal")
			}
		case *ast.Identifier:
			// search identifier
			var value1 value.Value
			var ok bool
			var tmp1 value.Value
			var gep value.Value
			declPos := arg.Decl.Name().Token().Position.Absolute
			// is it local or global
			if value1, ok = f.idents[declPos]; ok {
				var variableValue value.Value
				var typ types.Type
				typ = value1.Type()
				// if param passed by reference, get the value pointed to
				if paramItem, ok := arg.Decl.(*ast.ParamItem); ok {
					if !paramItem.ByVal {
						// get the type of the identifier
						variableValue = f.currentBlock.NewLoad(typ, value1)
						typ = variableValue.Type().(*types.PointerType).ElemType
					}
				}
				variableValue = f.currentBlock.NewLoad(typ, value1)

				m.printValue(f, arg, variableValue)
			} else {
				value1 = m.LookupGlobal(arg.Name)
				globalValue := value1.(*ir.Global)
				if globalValue != nil {
					typ := globalValue.Typ.ElemType
					// is it a pointer?
					argType, _ := arg.Decl.Type()
					if argType.(*uBasictypes.Basic).Kind == uBasictypes.String {
						switch typ.(type) {
						case *types.PointerType:
							gep = constant.NewGetElementPtr(typ, globalValue, zero)
							tmp1 = f.currentBlock.NewLoad(typ, gep)
						default:
							gep = globalValue
							tmp1 = gep
						}
					} else {
						gep = globalValue
						tmp1 = f.currentBlock.NewLoad(typ, gep)
					}
					m.printValue(f, arg, tmp1)
				} else {
					panic("unknown identifier")
				}
			}
		case *ast.CallOrIndexExpr:
			val := m.expr(f, arg)
			m.printValue(f, arg, val)
		default:
			val := m.expr(f, arg)
			m.printValue(f, arg, val)
			// // free the memory
			// switch val := val.(type) {
			// case *ir.InstCall:
			// 	_, ok := val.Type().(*types.PointerType)
			// 	if ok {
			// 		f.currentBlock.NewCall(m.LookupFunction("free"), val)
			// 	}
			// }
		}
	}
	if node.Semicolon == nil {
		value := m.LookupGlobal("vbEmpty")
		typ := value.Typ.ElemType
		gep := constant.NewGetElementPtr(typ, value, zero, zero)
		f.currentBlock.NewCall(puts, gep)
	}
}

// Print1Value prints a value to the console
func (m *Module) printValue(f *Function, arg ast.Expression, val value.Value) {
	printf := m.LookupFunction("printf")

	astType, _ := arg.Type()
	if funcType, ok := astType.(*uBasictypes.Func); ok {
		// get the function type
		astType = funcType.Result
	} else if arrayType, ok := astType.(*uBasictypes.Array); ok {
		// get the array type
		astType = arrayType.Type
	}
	var format string
	switch astType.(*uBasictypes.Basic).Kind {
	case uBasictypes.Integer:
		format = "%d"
	case uBasictypes.Long:
		format = "%ld"
	case uBasictypes.Single:
		format = "%f"
	case uBasictypes.Currency:
		format = "%f"
	case uBasictypes.Double:
		format = "%lf"
	case uBasictypes.String:
		format = "%s"
	case uBasictypes.Boolean:
		truestr0 := m.LookupGlobal("true")
		truestr1 := truestr0.Init
		falsestr0 := m.LookupGlobal("false")
		falsestr1 := falsestr0.Init
		// create a switch to check the value of the boolean
		cmp := f.currentBlock.NewICmp(enum.IPredEQ, val, constant.NewInt(types.I1, 1))
		val = f.currentBlock.NewSelect(cmp,
			constant.NewGetElementPtr(truestr1.Type(), truestr0, constant.NewInt(types.I64, 0)),
			constant.NewGetElementPtr(falsestr1.Type(), falsestr0, constant.NewInt(types.I64, 0)))
		format = "%s"
	case uBasictypes.Date:
		// format = "%s"
		format = "%f" // TODO: fix when we have a better way to handle dates
	case uBasictypes.Nothing:
		format = "%s"
	default:
		panic("unknown type")
	}
	format1 := constant.NewCharArrayFromString(format + "\x00")
	format2 := f.currentBlock.NewAlloca(format1.Typ)
	f.currentBlock.NewStore(format1, format2)
	f.currentBlock.NewCall(printf, format2, val)
}

// toIrType converts the given uBasic type to the corresponding LLVM IR type.
func toIrType(n uBasictypes.Type) types.Type {
	var t types.Type
	switch uBasicType := n.(type) {
	case *uBasictypes.Basic:
		switch uBasicType.Kind {
		case uBasictypes.Boolean:
			t = types.NewInt(1)
		case uBasictypes.Integer:
			t = types.NewInt(32)
		case uBasictypes.Long:
			t = types.NewInt(64)
		case uBasictypes.Single, uBasictypes.Currency, uBasictypes.Date:
			t = types.Float
		case uBasictypes.Double:
			t = types.Double
		case uBasictypes.Nothing:
			t = types.Void
		case uBasictypes.String:
			t = types.NewPointer(types.I8)
		}
	case *uBasictypes.Array:
		elem := toIrType(uBasicType.Type)
		var length = 0
		for _, dim := range uBasicType.Dimensions {
			length *= dim
		}
		if length == 0 {
			t = types.NewPointer(elem) // dynamic array
		} else {
			t = types.NewArray(uint64(length), elem) // static array
		}
	case *uBasictypes.Func:
		var params []types.Type
		variadic := false
		for _, p := range uBasicType.Params {
			pt := toIrType(p.Type)
			dbg.Printf("converting type %#v to %#v", p.Type, pt)
			params = append(params, pt)
		}
		result := toIrType(uBasicType.Result)
		typ := types.NewFunc(result, params...)
		typ.Variadic = variadic
		t = typ
	case *uBasictypes.Sub:
		var params []types.Type
		variadic := false
		for _, p := range uBasicType.Params {
			pt := toIrType(p.Type)
			dbg.Printf("converting type %#v to %#v", p.Type, pt)
			params = append(params, pt)
		}
		result := types.Void
		typ := types.NewFunc(result, params...)
		typ.Variadic = variadic
		t = typ
	default:
		panic(fmt.Sprintf("support for translating type %T not yet implemented.", uBasicType))
	}
	if t == nil {
		panic(fmt.Errorf("conversion failed: %#v", n))
	}
	return t
}

const ErrorHandler = ".ErrorHandler"
const ErrorMessage = ".ErrorMessage"
const ErrorNumber = ".ErrorNumber"
const JumpBuffer = ".JumpBuffer"
const throwException = ".throwException"

func (m *Module) genErrorHandler() {
	// throw exception function
	throwException := m.NewFunc(throwException, types.Void)
	entry0 := throwException.NewBlock("")
	entry0.NewCall(m.LookupFunction("longjmp"), m.LookupGlobal(JumpBuffer), constant.NewInt(types.I32, 1))
	entry0.NewUnreachable()

}

func (m *Module) genExternals() {
	m.genGC()
	// Convenience types and values.
	i32 := types.I32
	i8 := types.I8
	i8ptr := types.NewPointer(i8)

	// io functions -----------------
	printf := m.NewFunc("printf", i32, ir.NewParam("format", i8ptr))
	printf.Sig.Variadic = true
	m.NewFunc("puts", i32, ir.NewParam("s", i8ptr))
	scanf := m.NewFunc("scanf", i32, ir.NewParam("format", i8ptr))
	scanf.Sig.Variadic = true
	// string manipulation -----------------
	m.NewFunc("strcpy", i8ptr, ir.NewParam("dst", i8ptr), ir.NewParam("src", i8ptr))
	m.NewFunc("strcat", i8ptr, ir.NewParam("dst", i8ptr), ir.NewParam("src", i8ptr))
	m.NewFunc("sscanf", i32, ir.NewParam("str", i8ptr), ir.NewParam("format", i8ptr), ir.NewParam("dst", i8ptr))
	m.NewFunc("strlen", i32, ir.NewParam("str", i8ptr))

	sprintf := m.NewFunc("sprintf", i32, ir.NewParam("str", i8ptr), ir.NewParam("format", i8ptr))
	sprintf.Sig.Variadic = true

	// constants
	m.newGlobalStringConstant("", "vbEmpty")
	m.newGlobalStringConstant("\x0D", "vbCR")
	m.newGlobalStringConstant("\x0A", "vbLF")
	m.newGlobalStringConstant("\x0D\x0A", "vbCrLf")
	m.newGlobalStringConstant("\x09", "vbTab")
	m.newGlobalStringConstant("True", "true")
	m.newGlobalStringConstant("False", "false")

	// exception handling global variables
	m.NewFunc("exit", types.Void, ir.NewParam("status", types.I32))
	jump_bufferType := types.NewArray(48, types.I32)
	jump_buffer := m.NewGlobal(JumpBuffer, jump_bufferType)
	jump_buffer.Init = constant.NewZeroInitializer(jump_bufferType)
	errorNumber := m.NewGlobal(ErrorNumber, types.I32)
	errorNumber.Init = constant.NewInt(types.I8, 0)
	errorMessageType := types.NewArray(256, types.I8)
	errorMessage := m.NewGlobal(ErrorMessage, errorMessageType)
	errorMessage.Init = constant.NewZeroInitializer(errorMessageType)

	// setjmp and longjmp
	m.NewFunc("setjmp", types.I32, ir.NewParam("", types.NewPointer(types.I32)))
	m.NewFunc("longjmp", types.Void, ir.NewParam("", types.NewPointer(types.I32)), ir.NewParam("", types.I32))

	// exception constant
	divisionByZero := constant.NewCharArrayFromString("Division by zero\n\x00")
	dv0 := m.NewGlobalDef(".divisionByZero", divisionByZero)
	dv0.Linkage = enum.LinkagePrivate
	arrayOutOfBounds := constant.NewCharArrayFromString("Array index out of bounds\n\x00")
	aob := m.NewGlobalDef(".arrayIndexOutOfBounds", arrayOutOfBounds)
	aob.Linkage = enum.LinkagePrivate
}

// ----- error numbers -----

const (
	// ErrorNumberDivisionByZero is the error number for division by zero.
	NoError int64 = iota
	ErrorDivisionByZero
	ErrorIndexOutOfBounds
)

// checkIfDivisionByZero checks if the given expression is zero, if so generate an error message and return.
func (m *Module) checkIfDivisionByZero(f *Function, val value.Value) {
	// branching blocks
	trueBranch := f.NewBlock("")
	end := f.NewBlock("")

	switch val.Type().(type) {
	case *types.FloatType:
		// check if the value is zero
		zero := constant.NewFloat(types.Float, 0)
		cond := f.currentBlock.NewFCmp(enum.FPredOEQ, val, zero)
		f.currentBlock.NewCondBr(cond, trueBranch.Block, end.Block)
	case *types.IntType:
		// check if the value is zero
		zero := constant.NewInt(types.I32, 0)
		cond := f.currentBlock.NewICmp(enum.IPredEQ, val, zero)
		f.currentBlock.NewCondBr(cond, trueBranch.Block, end.Block)
	}

	// trueBranch:
	f.Blocks = append(f.Blocks, f.currentBlock.Block)
	f.currentBlock = trueBranch
	errorNum := constant.NewInt(types.I32, ErrorDivisionByZero)
	f.currentBlock.NewStore(errorNum, m.LookupGlobal(ErrorNumber))
	f.currentBlock.NewCall(m.LookupFunction("strcpy"), m.LookupGlobal(ErrorMessage), m.LookupGlobal(".divisionByZero"))
	f.currentBlock.NewCall(m.LookupFunction(throwException))
	f.currentBlock.NewUnreachable()
	// end:
	f.Blocks = append(f.Blocks, f.currentBlock.Block)
	f.currentBlock = end
}

// checkArrayBounds checks if the given index is out of bounds, if so generate an error message and return.
func (m *Module) checkArrayBounds(f *Function, index value.Value, length value.Value) {

	// TODO add compile time chech for array bounds

	zero := constZero(index.Type())
	trueBranch1 := f.NewBlock("")
	end1 := f.NewBlock("")
	end2 := f.NewBlock("")

	// inconming block:
	cond := f.currentBlock.NewICmp(enum.IPredULT, index, zero)
	f.currentBlock.NewCondBr(cond, trueBranch1.Block, end1.Block)

	// trueBranch1:
	f.Blocks = append(f.Blocks, f.currentBlock.Block)
	f.currentBlock = trueBranch1
	errorNumber := constant.NewInt(types.I32, ErrorIndexOutOfBounds)
	f.currentBlock.NewStore(errorNumber, m.LookupGlobal(ErrorNumber))
	f.currentBlock.NewCall(m.LookupFunction("strcpy"), m.LookupGlobal(ErrorMessage), m.LookupGlobal(".arrayIndexOutOfBounds"))
	f.currentBlock.NewCall(m.LookupFunction(throwException))
	f.currentBlock.NewUnreachable()

	// end1:
	f.Blocks = append(f.Blocks, f.currentBlock.Block)
	f.currentBlock = end1
	cond = f.currentBlock.NewICmp(enum.IPredUGE, index, length)
	f.currentBlock.NewCondBr(cond, trueBranch1.Block, end2.Block)

	// end2:
	f.Blocks = append(f.Blocks, f.currentBlock.Block)
	f.currentBlock = end2
}
