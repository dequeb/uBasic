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
				typ := value1.Type()
				variableValue := f.currentBlock.NewLoad(typ, value1)
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
		default:
			val := m.expr(f, arg)
			m.printValue(f, arg, val)
			// free the memory
			switch val := val.(type) {
			case *ir.InstCall:
				_, ok := val.Type().(*types.PointerType)
				if ok {
					f.currentBlock.NewCall(m.LookupFunction("free"), val)
				}
			}
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
	var err error
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
	if err != nil {
		panic(err)
	}
	if t == nil {
		panic(fmt.Errorf("conversion failed: %#v", n))
	}
	return t
}

func (m *Module) genExternals() {
	// Convenience types and values.
	i32 := types.I32
	i8 := types.I8
	i8ptr := types.NewPointer(i8)
	//zero := constant.NewInt(i32, 0)

	// garbage collection -----------------
	// %struct.GarbageCollector = type { ptr, i8, ptr, i64 }
	// @gc = external global %struct.GarbageCollector, align 8
	// m.NewGlobal("gc", types.NewStruct(types.NewPointer(i8), types.I8, types.NewPointer(i8), types.I64))

	// declare void @gc_start(ptr noundef, ptr noundef) #1
	// declare i64 @gc_stop(ptr noundef) #1
	// declare ptr @gc_malloc(ptr noundef, i64 noundef) #1
	// m.NewFunc("gc_start", types.Void, ir.NewParam("ptr", types.NewPointer(types.Void)), ir.NewParam("ptr", types.NewPointer(types.Void)))
	// m.NewFunc("gc_stop", types.I64, ir.NewParam("ptr", types.NewPointer(types.Void)))
	// m.NewFunc("gc_malloc", types.NewPointer(types.Void), ir.NewParam("ptr", types.NewPointer(types.Void)), ir.NewParam("size", types.I64))
	m.NewFunc("malloc", types.NewPointer(types.I8), ir.NewParam("size", types.I64))
	m.NewFunc("free", types.Void, ir.NewParam("ptr", types.NewPointer(types.I8)))

	// string functions -----------------
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
}
