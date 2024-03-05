package rtlib

// r-t lib is a runtime library for the µBASIC interpreter. It provides
// functionality for evaluating expressions, statements, and other µBASIC related
// tasks.

// to transform it into a µBASIC runtime library:
// see https://github.com/charlierobin/creating-dylibs-with-go/tree/master
// add import "C" to the file
// add //export <function_name> to the function
// generate the .h file with go tool cgo -godefs <file.go>
// compile the .dylib with go build -buildmode=c-shared -o <file>.dylib <file>.go

import (
	"fmt"
	"io"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"uBasic/ast"
	"uBasic/object"
	"uBasic/sem"
	"uBasic/token"
)

type Terminal struct {
	In  io.Reader
	Out io.Writer
}

var terminal *Terminal

func Init(term *Terminal) {
	terminal = term
}

// EvalSpecialStatement evaluates a special statement.
func EvalSpecialStatement(node *ast.SpecialStmt, params []object.Object, env *object.Environment) object.Object {
	kind := node.Keyword1.Kind
	var value object.Object
	switch kind {
	case token.Identifier:
		keyword := strings.ToLower(node.Keyword1.Literal)
		switch keyword {
		case "print", "debug.print", "msgbox":
			buf := strings.Builder{}
			value = object.NOTHING

			if len(params) > 0 {
				for _, value = range params {
					switch value := value.(type) {
					case *object.String:
						buf.WriteString(value.Value)
					default:
						str := fmt.Sprint(value)
						buf.WriteString(str)
					}
				}
			}
			if node.Semicolon == nil {
				buf.WriteString("\n")
			}
			// fmt.Print(buf.String())
			terminal.Out.Write([]byte(buf.String()))
			return value // return last print value or Nothing for automatic testing
		case "input":
			if len(params) == 0 {
				return object.NewError(node.Keyword1.Position, "missing argument")
			}
			if len(params) > 2 {
				return object.NewError(node.Keyword1.Position, "too many arguments")
			}
			prompt := ""
			if len(params) == 2 {
				if params[1].Type() != object.STRING_OBJ {
					return object.NewError(params[1].Position(), "argument must be a string")
				}
				prompt = params[1].(*object.String).Value
			}
			if params[0].Type() != object.STRING_OBJ {
				return object.NewError(params[0].Position(), "argument must be a string")
			}
			// print prompt
			terminal.Out.Write([]byte(prompt))
			var input string
			fmt.Fscanln(terminal.In, &input)
			// update parameter 2 with the input
			if len(params) == 2 {
				params[1].(*object.String).Value = input
			}
			return object.NewString(input, params[0].Position())
		}
	}
	return object.NewError(node.Keyword1.Position, "unknown special statement: "+node.Keyword1.Literal)
}

// EvalBody evaluates the body of a function.
func EvalBody(fn object.Object, params []object.Object, env *object.Environment) object.Object {
	// get function name
	var name string
	switch fn := fn.(type) {
	case *object.Function:
		name = fn.Definition.FuncName.Name
	case *object.Sub:
		name = fn.Definition.SubName.Name
	}
	name = strings.ToLower(name)
	switch name {
	case "chr":
		return evalChr(params, env)
	case "instr":
		return evalInstr(params, env)
	case "lcase":
		return evalLCase(params, env)
	case "left":
		return evalLeft(params, env)
	case "len":
		return evalLen(params, env)
	case "ltrim":
		return evalLTrim(params, env)
	case "mid":
		return evalMid(params, env)
	case "right":
		return evalRight(params, env)
	case "rtrim":
		return evalRTrim(params, env)
	case "space":
		return evalSpace(params, env)
	case "strcomp":
		return evalStrComp(params, env)
	case "strng":
		return evalStrng(params, env)
	case "strreverse":
		return evalStrReverse(params, env)
	case "trim":
		return evalTrim(params, env)
	case "ucase":
		return evalUCase(params, env)

		// --------------------------------
		// ------- date/time functions ----
		// --------------------------------

	case "dte":
		return evalDte(params, env)
	case "dateadd":
		return evalDateAdd(params, env)
	case "datediff":
		return evalDateDiff(params, env)
	case "datepart":
		return evalDatePart(params, env)
	case "dateserial":
		return evalDateSerial(params, env)
	case "datevalue":
		return evalDateValue(params, env)
	case "day":
		return evalDay(params, env)
	case "hour":
		return evalHour(params, env)
	case "minute":
		return evalMinute(params, env)
	case "month":
		return evalMonth(params, env)
	case "now":
		return evalNow(params, env)
	case "second":
		return evalSecond(params, env)
	case "time":
		return evalTime(params, env)
	case "timer":
		return evalTimer(params, env)
	case "timeserial":
		return evalTimeSerial(params, env)
	case "timevalue":
		return evalTimeValue(params, env)
	case "weekday":
		return evalWeekday(params, env)
	case "year":
		return evalYear(params, env)

		// --------------------------------
		// ------- conversion functions ---
		// --------------------------------
	case "cbool":
		return evalCBool(params, env)
	case "cdate":
		return evalCDate(params, env)
	case "cdbl":
		return evalCDbl(params, env)
	case "clng":
		return evalCLng(params, env)
	case "cstr":
		return evalCStr(params, env)
	case "cvar":
		return evalCVar(params, env)
	case "asc":
		return evalAsc(params, env)
	case "format":
		return evalFormat(params, env)
	case "hex":
		return evalHex(params, env)
	case "oct":
		return evalOct(params, env)

		// --------------------------------
		// ------- mathematical functions -
		// --------------------------------

	case "abs":
		return evalAbs(params, env)
	case "atn":
		return evalAtn(params, env)
	case "cos":
		return evalCos(params, env)
	case "expn":
		return evalExpn(params, env)
	case "fix":
		return evalFix(params, env)
	case "int":
		return evalInt(params, env)
	case "log":
		return evalLog(params, env)
	case "rnd":
		return evalRnd(params, env)
	case "sgn":
		return evalSgn(params, env)
	case "sin":
		return evalSin(params, env)
	case "sqr":
		return evalSqr(params, env)
	case "tan":
		return evalTan(params, env)

		// --------------------------------
		// ------- array functions -------
		// --------------------------------
	case "lbound":
		return evalLBound(params, env)
	case "ubound":
		return evalUBound(params, env)

		// --------------------------------
		// ------- input output -----------
		// --------------------------------
	case "input":
		return evalInput(params, env)

	default:
		pos := token.Position{Line: 0, Column: 0, Absolute: -1}
		return object.NewError(pos, "unknown function: "+name)
	}
}

func EvalClassProperty(class *object.Class, member object.Object, env *object.Environment) object.Object {
	switch class.Name {
	case "Application":
		return evalApplicationProperties(member, env)
	default:
		return object.NewError(class.Position(), "unknown class: "+class.Name)
	}
}

func EvalClassMethod(class *object.Class, method *object.Function, env *object.Environment) object.Object {
	switch class.Name {
	case "Application":
		return evalApplicationMethod(method, env)
	default:
		return object.NewError(class.Position(), "unknown class: "+class.Name)
	}
}

// evalApplicationMethod evaluates the Application class.
func evalApplicationMethod(fn *object.Function, env *object.Environment) object.Object {
	switch fn.Definition.FuncName.Name {
	case "getOS":
		return object.NewString("macOS", fn.Position())
	default:
		return object.NewError(fn.Position(), "unknown class method: "+fn.Definition.FuncName.Name)
	}
}

// evalApplicationProperties evaluates the Application class.
func evalApplicationProperties(member object.Object, env *object.Environment) object.Object {
	name := member.(*object.String).Value
	switch name {
	// TODO: implement all properties
	case "Name":
		return object.NewString("uBasic", member.Position())
	case "User":
		return object.NewString("user1", member.Position())
	default:
		return object.NewError(member.Position(), "unknown member: "+member.String())
	}
}

// evalChr evaluates the chr function.
func evalChr(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}
	param := params[0]
	if param.Type() == object.VARIANT_OBJ {
		param = param.(*object.Variant).Value
	}
	if param.Type() != object.LONG_OBJ {
		return object.NewError(params[0].Position(), "argument must be a long")
	}
	value := int32(params[0].(*object.Long).Value)
	return object.NewString(string(value), params[0].Position())
}

// evalAsc evaluates the asc function.
func evalAsc(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}
	if params[0].Type() != object.STRING_OBJ {
		return object.NewError(params[0].Position(), "argument must be a string")
	}
	value := params[0].(*object.String).Value
	if len(value) == 0 {
		return object.NewError(params[0].Position(), "argument must not be empty")
	}
	return object.NewIntegerByInt(int32(value[0]), params[0].Position())
}

// evalMid evaluates the mid function.
func evalMid(params []object.Object, env *object.Environment) object.Object {
	if len(params) < 2 || len(params) > 3 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=2 or 3", len(params)))
	}
	param0 := params[0]
	if param0.Type() == object.VARIANT_OBJ {
		param0 = param0.(*object.Variant).Value
	}
	if param0.Type() != object.STRING_OBJ {
		return object.NewError(params[0].Position(), "first argument must be a string")
	}
	param1 := params[1]
	if param1.Type() == object.VARIANT_OBJ {
		param1 = param1.(*object.Variant).Value
	}
	if param1.Type() != object.LONG_OBJ {
		return object.NewError(params[1].Position(), "second argument must be a longr")
	}
	param2 := params[2]
	if param2.Type() == object.VARIANT_OBJ {
		param2 = param2.(*object.Variant).Value
	}
	if param2.Type() != object.LONG_OBJ {
		return object.NewError(params[2].Position(), "third argument must be a long")
	}
	length := int32(params[2].(*object.Long).Value)
	value := params[0].(*object.String).Value
	start := int32(params[1].(*object.Long).Value)

	// if len is negative, start from the end
	if start < 0 {
		start = int32(len(value)) + start + 1
	}
	if length < 0 {
		// evaluate with 2 arguments
		return object.NewString(value[start-1:], params[0].Position())
	}
	if length+start-1 > int32(len(value)) {
		length = int32(len(value)) - start + 1
	}
	return object.NewString(value[start-1:start+length-1], params[0].Position())
}

// evalLeft evaluates the left function.
func evalLeft(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 2 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=2", len(params)))
	}
	param0 := params[0]
	if param0.Type() == object.VARIANT_OBJ {
		param0 = param0.(*object.Variant).Value
	}
	if param0.Type() != object.STRING_OBJ {
		return object.NewError(params[0].Position(), "first argument must be a string")
	}
	param1 := params[1]
	if param1.Type() == object.VARIANT_OBJ {
		param1 = param1.(*object.Variant).Value
	}
	if param1.Type() != object.LONG_OBJ {
		return object.NewError(params[1].Position(), "second argument must be a long")
	}
	value := param0.(*object.String).Value
	length := int(param1.(*object.Long).Value)
	return object.NewString(value[:length], params[0].Position())
}

// evalRight evaluates the right function.
func evalRight(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 2 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=2", len(params)))
	}
	param0 := params[0]
	if param0.Type() == object.VARIANT_OBJ {
		param0 = param0.(*object.Variant).Value
	}
	if param0.Type() != object.STRING_OBJ {
		return object.NewError(params[0].Position(), "first argument must be a string")
	}
	param1 := params[1]
	if param1.Type() == object.VARIANT_OBJ {
		param1 = param1.(*object.Variant).Value
	}
	if param1.Type() != object.LONG_OBJ {
		return object.NewError(params[1].Position(), "second argument must be a long")
	}
	value := param0.(*object.String).Value
	length := param1.(*object.Long).Value
	return object.NewString(value[len(value)-int(length):], params[0].Position())
}

// evalLen evaluates the len function.
func evalLen(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}
	param := params[0]
	if param.Type() == object.VARIANT_OBJ {
		param = param.(*object.Variant).Value
	}

	if param.Type() != object.STRING_OBJ {
		return object.NewError(params[0].Position(), "argument must be a string")
	}
	value := param.(*object.String).Value
	return object.NewIntegerByInt(int32(len(value)), params[0].Position())
}

// evalLCase evaluates the lcase function.
func evalLCase(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}
	param := params[0]
	if param.Type() == object.VARIANT_OBJ {
		param = param.(*object.Variant).Value
	}

	if param.Type() != object.STRING_OBJ {
		return object.NewError(param.Position(), "argument must be a string")
	}
	value := param.(*object.String).Value
	return object.NewString(strings.ToLower(value), param.Position())
}

// evalLTrim evaluates the ltrim function.
func evalLTrim(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}
	param0 := params[0]
	if param0.Type() == object.VARIANT_OBJ {
		param0 = param0.(*object.Variant).Value
	}
	if param0.Type() != object.STRING_OBJ {
		return object.NewError(params[0].Position(), "argument must be a string")
	}
	value := param0.(*object.String).Value
	return object.NewString(strings.TrimLeft(value, " "), params[0].Position())
}

// evalRTrim evaluates the rtrim function.
func evalRTrim(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}
	param0 := params[0]
	if param0.Type() == object.VARIANT_OBJ {
		param0 = param0.(*object.Variant).Value
	}
	if param0.Type() != object.STRING_OBJ {
		return object.NewError(params[0].Position(), "argument must be a string")
	}
	value := param0.(*object.String).Value
	return object.NewString(strings.TrimRight(value, " "), params[0].Position())
}

// evalSpace evaluates the space function.
func evalSpace(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}
	if params[0].Type() != object.LONG_OBJ {
		return object.NewError(params[0].Position(), "argument must be a long")
	}
	value := params[0].(*object.Long).Value
	return object.NewString(strings.Repeat(" ", int(value)), params[0].Position())
}

// evalStrComp evaluates the strcomp function.
func evalStrComp(params []object.Object, env *object.Environment) object.Object {
	if len(params) < 2 || len(params) > 3 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=2 or 3", len(params)))
	}
	param0 := params[0]
	if param0.Type() == object.VARIANT_OBJ {
		param0 = param0.(*object.Variant).Value
	}
	if param0.Type() != object.STRING_OBJ {
		return object.NewError(params[0].Position(), "first argument must be a string")
	}
	param1 := params[1]
	if param1.Type() == object.VARIANT_OBJ {
		param1 = param1.(*object.Variant).Value
	}
	if param1.Type() != object.STRING_OBJ {
		return object.NewError(params[1].Position(), "second argument must be a string")
	}
	value1 := param0.(*object.String).Value
	value2 := param1.(*object.String).Value
	var option int64
	var param2 object.Object
	if len(params) == 3 {
		param2 = params[2]
		if param2.Type() == object.VARIANT_OBJ {
			param2 = param2.(*object.Variant).Value
		}
		if param2.Type() != object.LONG_OBJ {
			return object.NewError(params[2].Position(), "third argument must be a long")
		}
		option = param2.(*object.Long).Value
	}
	if len(params) == 2 || option == 0 { // compare case-sensitive
		result := strings.Compare(value1, value2)
		return object.NewIntegerByInt(int32(result), params[0].Position())
	}
	if option < 0 || option > 2 {
		return object.NewError(params[2].Position(), "third argument must be between 0 and 1")
	}
	// compare case-sensitive
	value1 = strings.ToLower(value1)
	value2 = strings.ToLower(value2)
	result := strings.Compare(value1, value2)
	return object.NewIntegerByInt(int32(result), params[0].Position())
}

// evalStrng evaluates the string function.
func evalStrng(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 2 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=2", len(params)))
	}
	param0 := params[0]
	if param0.Type() == object.VARIANT_OBJ {
		param0 = param0.(*object.Variant).Value
	}
	if param0.Type() != object.LONG_OBJ {
		return object.NewError(params[0].Position(), "first argument must be a long")
	}
	param1 := params[1]
	if param1.Type() == object.VARIANT_OBJ {
		param1 = param1.(*object.Variant).Value
	}
	if param1.Type() != object.STRING_OBJ {
		return object.NewError(params[1].Position(), "second argument must be a string")
	}
	value := param0.(*object.Long).Value
	return object.NewString(strings.Repeat(param1.(*object.String).Value, int(value)), params[0].Position())
}

// evalStrReverse evaluates the strreverse function.
func evalStrReverse(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}
	param0 := params[0]
	if param0.Type() == object.VARIANT_OBJ {
		param0 = param0.(*object.Variant).Value
	}
	if param0.Type() != object.STRING_OBJ {
		return object.NewError(params[0].Position(), "argument must be a string")
	}
	value := param0.(*object.String).Value
	runes := []rune(value)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return object.NewString(string(runes), params[0].Position())
}

// evalTrim evaluates the trim function.
func evalTrim(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}
	param0 := params[0]
	if param0.Type() == object.VARIANT_OBJ {
		param0 = param0.(*object.Variant).Value
	}
	if param0.Type() != object.STRING_OBJ {
		return object.NewError(params[0].Position(), "argument must be a string")
	}
	value := param0.(*object.String).Value
	return object.NewString(strings.TrimSpace(value), params[0].Position())
}

// evalUCase evaluates the ucase function.
func evalUCase(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}
	param := params[0]
	if param.Type() == object.VARIANT_OBJ {
		param = param.(*object.Variant).Value
	}

	if param.Type() != object.STRING_OBJ {
		return object.NewError(param.Position(), "argument must be a string")
	}
	value := param.(*object.String).Value
	return object.NewString(strings.ToUpper(value), param.Position())
}

// evalDte evaluates the date function.
func evalDte(params []object.Object, env *object.Environment) object.Object {
	t := time.Now()
	today := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, t.Nanosecond(), t.Location())
	return object.NewDateByTime(today, sem.UniversePos)
}

// evalDateAdd evaluates the dateadd function.
func evalDateAdd(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 3 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=3", len(params)))
	}
	param0 := params[0]
	if param0.Type() == object.VARIANT_OBJ {
		param0 = param0.(*object.Variant).Value
	}
	if param0.Type() != object.STRING_OBJ {
		return object.NewError(params[0].Position(), "first argument must be a string")
	}
	param1 := params[1]
	if param1.Type() == object.VARIANT_OBJ {
		param1 = param1.(*object.Variant).Value
	}
	if param1.Type() != object.LONG_OBJ {
		return object.NewError(params[1].Position(), "second argument must be a long")
	}

	param2 := params[2]
	if param2.Type() == object.VARIANT_OBJ {
		param2 = param2.(*object.Variant).Value
	}
	var date time.Time
	switch param2 := param2.(type) {
	case *object.String:
		value := param2.Value
		dateObj, err := object.NewDate(value, params[0].Position()) // convert to date object
		if err != nil {
			return object.NewError(params[0].Position(), err.Error())
		}
		date = dateObj.Value
	case *object.Date:
		date = param2.Value
	default:
		return object.NewError(params[2].Position(), "third argument must be a date or a string")
	}
	number := params[1].(*object.Long).Value
	unit := strings.ToLower(param0.(*object.String).Value)

	switch unit {
	case "yyyy":
		date = date.AddDate(int(number), 0, 0)
	case "q":
		date = date.AddDate(0, int(number)*3, 0)
	case "m":
		date = date.AddDate(0, int(number), 0)
	case "y":
		date = date.AddDate(0, 0, int(number))
	case "d":
		date = date.AddDate(0, 0, int(number))
	case "w":
		date = date.AddDate(0, 0, int(number*7))
	default:
		return object.NewError(params[2].Position(), "unknown unit: "+unit)
	}
	return object.NewDateByTime(date, param0.Position())
}

// evalDateDiff evaluates the datediff function.
func evalDateDiff(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 3 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=3", len(params)))
	}
	param0 := params[0]
	if param0.Type() == object.VARIANT_OBJ {
		param0 = param0.(*object.Variant).Value
	}
	if param0.Type() != object.STRING_OBJ {
		return object.NewError(params[0].Position(), "first argument must be a string")
	}
	param1 := params[1]
	if param1.Type() == object.VARIANT_OBJ {
		param1 = param1.(*object.Variant).Value
	}
	var date1 time.Time
	switch param1 := param1.(type) {
	case *object.String:
		value := param1.Value
		dateObj, err := object.NewDate(value, params[0].Position()) // convert to date object
		if err != nil {
			return object.NewError(params[0].Position(), err.Error())
		}
		date1 = dateObj.Value
	case *object.Date:
		date1 = param1.Value
	default:
		return object.NewError(params[2].Position(), "third argument must be a date or a string")
	}

	param2 := params[2]
	if param2.Type() == object.VARIANT_OBJ {
		param2 = param2.(*object.Variant).Value
	}
	var date2 time.Time
	switch param2 := param2.(type) {
	case *object.String:
		value := param2.Value
		dateObj, err := object.NewDate(value, params[0].Position()) // convert to date object
		if err != nil {
			return object.NewError(params[0].Position(), err.Error())
		}
		date2 = dateObj.Value
	case *object.Date:
		date2 = param2.Value
	default:
		return object.NewError(params[2].Position(), "third argument must be a date or a string")
	}
	var number int64
	unit := strings.ToLower(param0.(*object.String).Value)

	switch unit {
	case "yyyy":
		number = int64(date2.Year() - date1.Year())
	case "q":
		number = int64(date2.Month()/3 - date1.Month()/3)
	case "m":
		number = int64(date2.Month() - date1.Month())
	case "y":
		number = int64(date2.YearDay() - date1.YearDay())
	case "d":
		number = int64(date2.YearDay() - date1.YearDay())
	case "w":
		number = int64((date2.YearDay() - date1.YearDay()) / 7)
	case "ww":
		number = int64((date2.YearDay() - date1.YearDay()) / 7)
	case "h":
		number = int64(date2.Hour() - date1.Hour())
	case "n":
		number = int64(date2.Minute() - date1.Minute())
	case "s":
		number = int64(date2.Second() - date1.Second())

	default:
		return object.NewError(params[2].Position(), "unknown unit: "+unit)
	}
	return object.NewLongByInt(number, param0.Position())
}

// evalDatePart evaluates the datepart function.
func evalDatePart(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 2 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=2", len(params)))
	}
	param0 := params[0]
	if param0.Type() == object.VARIANT_OBJ {
		param0 = param0.(*object.Variant).Value
	}
	if param0.Type() != object.STRING_OBJ {
		return object.NewError(params[0].Position(), "first argument must be a string")
	}
	param1 := params[1]
	if param1.Type() == object.VARIANT_OBJ {
		param1 = param1.(*object.Variant).Value
	}
	var date time.Time
	switch param1 := param1.(type) {
	case *object.String:
		value := param1.Value
		dateObj, err := object.NewDate(value, params[0].Position()) // convert to date object
		if err != nil {
			return object.NewError(params[0].Position(), err.Error())
		}
		date = dateObj.Value
	case *object.Date:
		date = param1.Value
	default:
		return object.NewError(params[1].Position(), "third argument must be a date or a string")
	}

	unit := strings.ToLower(param0.(*object.String).Value)
	switch unit {
	case "yyyy":
		return object.NewIntegerByInt(int32(date.Year()), param0.Position())
	case "q":
		return object.NewIntegerByInt(int32(date.Month()/3), param0.Position())
	case "m":
		return object.NewIntegerByInt(int32(date.Month()), param0.Position())
	case "y":
		return object.NewIntegerByInt(int32(date.YearDay()), param0.Position())
	case "d":
		return object.NewIntegerByInt(int32(date.YearDay()), param0.Position())
	case "w":
		return object.NewIntegerByInt(int32(date.Weekday()+1), param0.Position())
	case "ww":
		return object.NewIntegerByInt(int32(date.YearDay()/7+1), param0.Position())
	case "h":
		return object.NewIntegerByInt(int32(date.Hour()), param0.Position())
	case "n":
		return object.NewIntegerByInt(int32(date.Minute()), param0.Position())
	case "s":
		return object.NewIntegerByInt(int32(date.Second()), param0.Position())
	default:
		return object.NewError(params[0].Position(), "unknown unit: "+unit)
	}

}

// evalDateSerial evaluates the dateserial function.
func evalDateSerial(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 3 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=3", len(params)))
	}
	param0 := params[0]
	if param0.Type() == object.VARIANT_OBJ {
		param0 = param0.(*object.Variant).Value
	}
	if param0.Type() != object.LONG_OBJ {
		return object.NewError(params[0].Position(), "first argument must be a long")
	}
	param1 := params[1]
	if param1.Type() == object.VARIANT_OBJ {
		param1 = param1.(*object.Variant).Value
	}
	if param1.Type() != object.LONG_OBJ {
		return object.NewError(params[1].Position(), "second argument must be a long")
	}
	param2 := params[2]
	if param2.Type() == object.VARIANT_OBJ {
		param2 = param2.(*object.Variant).Value
	}
	if param2.Type() != object.LONG_OBJ {
		return object.NewError(params[2].Position(), "third argument must be a long")
	}
	year := params[0].(*object.Long).Value
	month := params[1].(*object.Long).Value
	day := params[2].(*object.Long).Value

	date := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.Local)
	return object.NewDateByTime(date, param0.Position())
}

// evalDateValue evaluates the datevalue function.
func evalDateValue(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}
	param0 := params[0]
	if param0.Type() == object.VARIANT_OBJ {
		param0 = param0.(*object.Variant).Value
	}
	var date time.Time
	switch param0 := param0.(type) {
	case *object.String:
		value := param0.Value
		dateObj, err := object.NewDate(value, params[0].Position()) // convert to date object
		if err != nil {
			return object.NewError(params[0].Position(), err.Error())
		}
		date = dateObj.Value
	case *object.Date:
		date = param0.Value
	default:
		return object.NewError(params[2].Position(), "first argument must be a date or a string")
	}
	return object.NewDateByTime(date, param0.Position())
}

// evalDay evaluates the day function.
func evalDay(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}
	param0 := params[0]
	if param0.Type() == object.VARIANT_OBJ {
		param0 = param0.(*object.Variant).Value
	}
	var date time.Time
	switch param0 := param0.(type) {
	case *object.String:
		value := param0.Value
		dateObj, err := object.NewDate(value, params[0].Position()) // convert to date object
		if err != nil {
			return object.NewError(params[0].Position(), err.Error())
		}
		date = dateObj.Value
	case *object.Date:
		date = param0.Value
	default:
		return object.NewError(params[2].Position(), "first argument must be a date or a string")
	}

	return object.NewIntegerByInt(int32(date.Day()), param0.Position())
}

// evalHour evaluates the hour function.
func evalHour(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}
	param0 := params[0]
	if param0.Type() == object.VARIANT_OBJ {
		param0 = param0.(*object.Variant).Value
	}
	var date time.Time
	switch param0 := param0.(type) {
	case *object.String:
		value := param0.Value
		dateObj, err := object.NewDate(value, params[0].Position()) // convert to date object
		if err != nil {
			return object.NewError(params[0].Position(), err.Error())
		}
		date = dateObj.Value
	case *object.Date:
		date = param0.Value
	default:
		return object.NewError(params[2].Position(), "first argument must be a date or a string")
	}
	return object.NewIntegerByInt(int32(date.Hour()), param0.Position())
}

// evalMinute evaluates the minute function.
func evalMinute(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}
	param0 := params[0]
	if param0.Type() == object.VARIANT_OBJ {
		param0 = param0.(*object.Variant).Value
	}
	var date time.Time
	switch param0 := param0.(type) {
	case *object.String:
		value := param0.Value
		dateObj, err := object.NewDate(value, params[0].Position()) // convert to date object
		if err != nil {
			return object.NewError(params[0].Position(), err.Error())
		}
		date = dateObj.Value
	case *object.Date:
		date = param0.Value
	default:
		return object.NewError(params[2].Position(), "first argument must be a date or a string")
	}
	return object.NewIntegerByInt(int32(date.Minute()), param0.Position())
}

// evalMonth evaluates the month function.
func evalMonth(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}
	param0 := params[0]
	if param0.Type() == object.VARIANT_OBJ {
		param0 = param0.(*object.Variant).Value
	}
	var date time.Time
	switch param0 := param0.(type) {
	case *object.String:
		value := param0.Value
		dateObj, err := object.NewDate(value, params[0].Position()) // convert to date object
		if err != nil {
			return object.NewError(params[0].Position(), err.Error())
		}
		date = dateObj.Value
	case *object.Date:
		date = param0.Value
	default:
		return object.NewError(params[2].Position(), "first argument must be a date or a string")
	}
	return object.NewIntegerByInt(int32(date.Month()), param0.Position())
}

// evalNow evaluates the now function.
func evalNow(params []object.Object, env *object.Environment) object.Object {
	return object.NewDateByTime(time.Now(), sem.UniversePos)
}

// evalSecond evaluates the second function.
func evalSecond(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}
	param0 := params[0]
	if param0.Type() == object.VARIANT_OBJ {
		param0 = param0.(*object.Variant).Value
	}
	var date time.Time
	switch param0 := param0.(type) {
	case *object.String:
		value := param0.Value
		dateObj, err := object.NewDate(value, params[0].Position()) // convert to date object
		if err != nil {
			return object.NewError(params[0].Position(), err.Error())
		}
		date = dateObj.Value
	case *object.Date:
		date = param0.Value
	default:
		return object.NewError(params[2].Position(), "first argument must be a date or a string")
	}
	return object.NewIntegerByInt(int32(date.Second()), param0.Position())
}

// evalTime evaluates the time function.
func evalTime(params []object.Object, env *object.Environment) object.Object {
	t := time.Now()
	now := time.Date(0, 1, 1, t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
	return object.NewDateByTime(now, sem.UniversePos)
}

// evalTimer evaluates the timer function.
func evalTimer(params []object.Object, env *object.Environment) object.Object {
	// calculate the number of seconds since midnight
	t := time.Now()
	seconds := t.Hour()*3600 + t.Minute()*60 + t.Second()
	return object.NewDoubleByFloat(float64(seconds), sem.UniversePos)
}

// evalTimeSerial evaluates the timeserial function.
func evalTimeSerial(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 3 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=3", len(params)))
	}
	param0 := params[0]
	if param0.Type() == object.VARIANT_OBJ {
		param0 = param0.(*object.Variant).Value
	}
	if param0.Type() != object.LONG_OBJ {
		return object.NewError(params[0].Position(), "first argument must be a long")
	}
	param1 := params[1]
	if param1.Type() == object.VARIANT_OBJ {
		param1 = param1.(*object.Variant).Value
	}
	if param1.Type() != object.LONG_OBJ {
		return object.NewError(params[1].Position(), "second argument must be a long")
	}
	param2 := params[2]
	if param2.Type() == object.VARIANT_OBJ {
		param2 = param2.(*object.Variant).Value
	}
	if param2.Type() != object.LONG_OBJ {
		return object.NewError(params[2].Position(), "third argument must be a long")
	}
	hours := params[0].(*object.Long).Value
	minutes := params[1].(*object.Long).Value
	seconds := params[2].(*object.Long).Value

	date := time.Date(0, 1, 1, int(hours), int(minutes), int(seconds), 0, time.Local)
	return object.NewDateByTime(date, param0.Position())

}

// evalTimeValue evaluates the timevalue function.
func evalTimeValue(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}
	param0 := params[0]
	if param0.Type() == object.VARIANT_OBJ {
		param0 = param0.(*object.Variant).Value
	}
	var date time.Time
	switch param0 := param0.(type) {
	case *object.String:
		value := param0.Value
		dateObj, err := object.NewDate(value, params[0].Position()) // convert to date object
		if err != nil {
			return object.NewError(params[0].Position(), err.Error())
		}
		date = dateObj.Value
	case *object.Date:
		date = param0.Value
	default:
		return object.NewError(params[2].Position(), "first argument must be a date or a string")
	}
	return object.NewDateByTime(date, param0.Position())
}

// evalWeekday evaluates the weekday function.
func evalWeekday(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}
	param0 := params[0]
	if param0.Type() == object.VARIANT_OBJ {
		param0 = param0.(*object.Variant).Value
	}
	var date time.Time
	switch param0 := param0.(type) {
	case *object.String:
		value := param0.Value
		dateObj, err := object.NewDate(value, params[0].Position()) // convert to date object
		if err != nil {
			return object.NewError(params[0].Position(), err.Error())
		}
		date = dateObj.Value
	case *object.Date:
		date = param0.Value
	default:
		return object.NewError(params[2].Position(), "first argument must be a date or a string")
	}
	return object.NewIntegerByInt(int32(date.Weekday()+1), param0.Position())
}

// evalYear evaluates the year function.
func evalYear(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}
	param0 := params[0]
	if param0.Type() == object.VARIANT_OBJ {
		param0 = param0.(*object.Variant).Value
	}
	var date time.Time
	switch param0 := param0.(type) {
	case *object.String:
		value := param0.Value
		dateObj, err := object.NewDate(value, params[0].Position()) // convert to date object
		if err != nil {
			return object.NewError(params[0].Position(), err.Error())
		}
		date = dateObj.Value
	case *object.Date:
		date = param0.Value
	default:
		return object.NewError(params[2].Position(), "first argument must be a date or a string")
	}
	return object.NewIntegerByInt(int32(date.Year()), param0.Position())
}

// evalCBool evaluates the cbool function.
func evalCBool(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}
	param := params[0]
	if param.Type() == object.VARIANT_OBJ {
		param = param.(*object.Variant).Value
	}
	switch param := param.(type) {
	case *object.Boolean:
		return param
	case *object.Integer:
		if param.Value != 0 {
			return object.NewBooleanByBool(true, param.Position())
		}
		return object.NewBooleanByBool(false, param.Position())
	case *object.Double:
		if param.Value != 0 {
			return object.NewBooleanByBool(true, param.Position())
		}
		return object.NewBooleanByBool(false, param.Position())
	case *object.Single:
		if param.Value != 0 {
			return object.NewBooleanByBool(true, param.Position())
		}
		return object.NewBooleanByBool(false, param.Position())
	case *object.Long:
		if param.Value != 0 {
			return object.NewBooleanByBool(true, param.Position())
		}
		return object.NewBooleanByBool(false, param.Position())
	case *object.Currency:
		if param.Value != 0 {
			return object.NewBooleanByBool(true, param.Position())
		}
		return object.NewBooleanByBool(false, param.Position())
	case *object.Date:
		if !param.Value.IsZero() {
			return object.NewBooleanByBool(true, param.Position())
		}
		return object.NewBooleanByBool(false, param.Position())
	case *object.String:
		obj, err := object.NewBoolean(param.Value, param.Position())
		if err != nil {
			return object.NewBooleanByBool(false, param.Position())
		}
		return obj
	}
	return object.NewBooleanByBool(false, param.Position())
}

// evalCDate evaluates the cdate function.
func evalCDate(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}
	param := params[0]
	if param.Type() == object.VARIANT_OBJ {
		param = param.(*object.Variant).Value
	}
	switch param := param.(type) {
	case *object.Date:
		return param
	case *object.String:
		date, err := object.NewDate(param.Value, param.Position())
		if err != nil {
			return object.NewDateByTime(time.Time{}, param.Position())
		}
		return date
	}
	return object.NewDateByTime(time.Time{}, params[0].Position())
}

// evalCDbl evaluates the cdbl function.
func evalCDbl(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}
	param := params[0]
	if param.Type() == object.VARIANT_OBJ {
		param = param.(*object.Variant).Value
	}
	switch param := param.(type) {
	case *object.Double:
		return param
	case *object.Integer:
		return object.NewDoubleByFloat(float64(param.Value), param.Position())
	case *object.Single:
		return object.NewDoubleByFloat(float64(param.Value), param.Position())
	case *object.Long:
		return object.NewDoubleByFloat(float64(param.Value), param.Position())
	case *object.Currency:
		return object.NewDoubleByFloat(float64(param.Value), param.Position())
	case *object.Date:
		return object.NewDoubleByFloat(float64(param.Value.Unix()), param.Position())
	case *object.String:
		value, err := strconv.ParseFloat(param.Value, 64)
		if err != nil {
			return object.NewDoubleByFloat(0, param.Position())
		}
		return object.NewDoubleByFloat(value, param.Position())
	}
	return object.NewDoubleByFloat(0, params[0].Position())
}

// evalCLng evaluates the clng function.
func evalCLng(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}
	param := params[0]
	if param.Type() == object.VARIANT_OBJ {
		param = param.(*object.Variant).Value
	}
	switch param := param.(type) {
	case *object.Long:
		return param
	case *object.Integer:
		return object.NewLongByInt(int64(param.Value), param.Position())
	case *object.Double:
		return object.NewLongByInt(int64(param.Value), param.Position())
	case *object.Single:
		return object.NewLongByInt(int64(param.Value), param.Position())
	case *object.Currency:
		return object.NewLongByInt(int64(param.Value), param.Position())
	case *object.Date:
		return object.NewLongByInt(param.Value.Unix(), param.Position())
	case *object.String:
		value, err := strconv.ParseInt(param.Value, 10, 64)
		if err != nil {
			return object.NewLongByInt(0, param.Position())
		}
		return object.NewLongByInt(value, param.Position())
	}
	return object.NewLongByInt(0, params[0].Position())
}

// ** evalCStr evaluates the cstr function.
func evalCStr(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}
	param := params[0]
	if param.Type() == object.VARIANT_OBJ {
		param = param.(*object.Variant).Value
	}

	switch param.Type() {
	case object.INTEGER_OBJ:
		return object.NewString(param.(*object.Integer).String(), param.Position())
	case object.BOOLEAN_OBJ:
		return object.NewString(param.(*object.Boolean).String(), param.Position())
	case object.DOUBLE_OBJ:
		return object.NewString(param.(*object.Double).String(), param.Position())
	case object.SINGLE_OBJ:
		return object.NewString(param.(*object.Single).String(), param.Position())
	case object.LONG_OBJ:
		return object.NewString(param.(*object.Long).String(), param.Position())
	case object.CURRENCY_OBJ:
		return object.NewString(param.(*object.Currency).String(), param.Position())
	case object.DATE_OBJ:
		return object.NewString(param.(*object.Date).String(), param.Position())
	case object.STRING_OBJ:
		return param
	default:
		return object.NewString("", param.Position())
	}
}

// evalCVar evaluates the cvar function.
func evalCVar(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}

	return object.NewVariantByObject(params[0], params[0].Position())
}

// evalFormat evaluates the format function.
// not compatible with the original function - must use go's fmt package
// see: https://pkg.go.dev/fmt
func evalFormat(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 2 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=2", len(params)))
	}
	// get value
	param0 := params[0]
	if param0.Type() == object.VARIANT_OBJ {
		param0 = param0.(*object.Variant).Value
	}

	// get format
	param1 := params[1]
	if param1.Type() == object.VARIANT_OBJ {
		param1 = param1.(*object.Variant).Value
	}
	if param1.Type() != object.STRING_OBJ {
		return object.NewError(params[1].Position(), "second argument must be a string")
	}
	format := param1.(*object.String).Value

	// format the value
	// TODO: implement the format function
	switch param0 := param0.(type) {
	case *object.Integer:
		return object.NewString(fmt.Sprintf(format, param0.Value), param0.Position())
	case *object.Long:
		return object.NewString(fmt.Sprintf(format, param0.Value), param0.Position())
	case *object.Double:
		return object.NewString(fmt.Sprintf(format, param0.Value), param0.Position())
	case *object.Single:
		return object.NewString(fmt.Sprintf(format, param0.Value), param0.Position())
	case *object.Currency:
		return object.NewString(fmt.Sprintf(format, param0.Value), param0.Position())
	case *object.Date:
		return object.NewString(fmt.Sprintf(format, param0.Value), param0.Position())
	case *object.String:
		return object.NewString(fmt.Sprintf(format, param0.Value), param0.Position())
	}
	return object.NewString("", params[0].Position())
}

// evalHex evaluates the hex function.
func evalHex(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}
	param := params[0]
	if param.Type() == object.VARIANT_OBJ {
		param = param.(*object.Variant).Value
	}
	switch param := param.(type) {
	case *object.Integer:
		return object.NewString(fmt.Sprintf("%X", param.Value), param.Position())
	case *object.Long:
		return object.NewString(fmt.Sprintf("%X", param.Value), param.Position())
	case *object.String:
		value, err := strconv.ParseInt(param.Value, 10, 64)
		if err != nil {
			return object.NewString("", param.Position())
		}
		return object.NewString(fmt.Sprintf("%X", value), param.Position())
	}
	return object.NewString("", params[0].Position())
}

// evalOct evaluates the oct function.
func evalOct(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}
	param := params[0]
	if param.Type() == object.VARIANT_OBJ {
		param = param.(*object.Variant).Value
	}
	switch param := param.(type) {
	case *object.Integer:
		return object.NewString(fmt.Sprintf("%o", param.Value), param.Position())
	case *object.Long:
		return object.NewString(fmt.Sprintf("%o", param.Value), param.Position())
	case *object.String:
		value, err := strconv.ParseInt(param.Value, 10, 64)
		if err != nil {
			return object.NewString("", param.Position())
		}
		return object.NewString(fmt.Sprintf("%o", value), param.Position())
	}
	return object.NewString("", params[0].Position())
}

// evalAbs evaluates the abs function.
func evalAbs(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}
	param := params[0]
	if param.Type() == object.VARIANT_OBJ {
		param = param.(*object.Variant).Value
	}
	switch param := param.(type) {
	case *object.Integer:
		return object.NewIntegerByInt(int32(math.Abs(float64(param.Value))), param.Position())
	case *object.Long:
		return object.NewLongByInt(int64(math.Abs(float64(param.Value))), param.Position())
	case *object.Double:
		return object.NewDoubleByFloat(math.Abs(param.Value), param.Position())
	case *object.Single:
		return object.NewSingleByFloat(float32(math.Abs(float64(param.Value))), param.Position())
	case *object.Currency:
		return object.NewCurrencyByFloat(math.Abs(param.Value), param.Position())
	default:
		return object.NewError(params[0].Position(), "argument must be a number")
	}
}

// evalAtn evaluates the atn function.
func evalAtn(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}
	param := params[0]
	if param.Type() == object.VARIANT_OBJ {
		param = param.(*object.Variant).Value
	}
	switch param := param.(type) {
	case *object.Integer:
		return object.NewDoubleByFloat(math.Atan(float64(param.Value)), param.Position())
	case *object.Long:
		return object.NewDoubleByFloat(math.Atan(float64(param.Value)), param.Position())
	case *object.Double:
		return object.NewDoubleByFloat(math.Atan(param.Value), param.Position())
	case *object.Single:
		return object.NewSingleByFloat(float32(math.Atan(float64(param.Value))), param.Position())
	case *object.Currency:
		return object.NewDoubleByFloat(math.Atan(param.Value), param.Position())
	default:
		return object.NewError(params[0].Position(), "argument must be a number")
	}
}

// evalCos evaluates the cos function.
func evalCos(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}
	param := params[0]
	if param.Type() == object.VARIANT_OBJ {
		param = param.(*object.Variant).Value
	}
	switch param := param.(type) {
	case *object.Integer:
		return object.NewDoubleByFloat(math.Cos(float64(param.Value)), param.Position())
	case *object.Long:
		return object.NewDoubleByFloat(math.Cos(float64(param.Value)), param.Position())
	case *object.Double:
		return object.NewDoubleByFloat(math.Cos(param.Value), param.Position())
	case *object.Single:
		return object.NewSingleByFloat(float32(math.Cos(float64(param.Value))), param.Position())
	case *object.Currency:
		return object.NewDoubleByFloat(math.Cos(param.Value), param.Position())
	default:
		return object.NewError(params[0].Position(), "argument must be a number")
	}
}

// evalExpn evaluates the Exp function.
func evalExpn(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}
	param := params[0]
	if param.Type() == object.VARIANT_OBJ {
		param = param.(*object.Variant).Value
	}
	switch param := param.(type) {
	case *object.Integer:
		return object.NewDoubleByFloat(math.Exp(float64(param.Value)), param.Position())
	case *object.Long:
		return object.NewDoubleByFloat(math.Exp(float64(param.Value)), param.Position())
	case *object.Double:
		return object.NewDoubleByFloat(math.Exp(param.Value), param.Position())
	case *object.Single:
		return object.NewSingleByFloat(float32(math.Exp(float64(param.Value))), param.Position())
	case *object.Currency:
		return object.NewDoubleByFloat(math.Exp(param.Value), param.Position())
	default:
		return object.NewError(params[0].Position(), "argument must be a number")
	}
}

// evalFix evaluates the fix function.
func evalFix(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}
	param := params[0]
	if param.Type() == object.VARIANT_OBJ {
		param = param.(*object.Variant).Value
	}
	switch param := param.(type) {
	case *object.Integer:
		return object.NewIntegerByInt(int32(param.Value), param.Position())
	case *object.Long:
		return object.NewLongByInt(int64(param.Value), param.Position())
	case *object.Double:
		return object.NewDoubleByFloat(math.Floor(param.Value), param.Position())
	case *object.Single:
		return object.NewSingleByFloat(float32(math.Floor(float64(param.Value))), param.Position())
	case *object.Currency:
		return object.NewCurrencyByFloat(math.Floor(param.Value), param.Position())
	default:
		return object.NewError(params[0].Position(), "argument must be a number")
	}
}

// evalInt evaluates the int function.
func evalInt(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}
	param := params[0]
	if param.Type() == object.VARIANT_OBJ {
		param = param.(*object.Variant).Value
	}
	switch param := param.(type) {
	case *object.Integer:
		return object.NewIntegerByInt(int32(param.Value), param.Position())
	case *object.Long:
		return object.NewLongByInt(int64(param.Value), param.Position())
	case *object.Double:
		return object.NewDoubleByFloat(math.Floor(param.Value), param.Position())
	case *object.Single:
		return object.NewSingleByFloat(float32(math.Floor(float64(param.Value))), param.Position())
	case *object.Currency:
		return object.NewCurrencyByFloat(math.Floor(param.Value), param.Position())
	default:
		return object.NewError(params[0].Position(), "argument must be a number")
	}
}

// evalLog evaluates the log function.
func evalLog(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 && len(params) != 2 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1 or 2", len(params)))
	}
	param0 := params[0]
	if param0.Type() == object.VARIANT_OBJ {
		param0 = param0.(*object.Variant).Value
	}
	switch param0 := param0.(type) {
	case *object.Integer:
		if len(params) == 1 {
			return object.NewDoubleByFloat(math.Log(float64(param0.Value)), param0.Position())
		}
		param1 := params[1]
		if param1.Type() == object.VARIANT_OBJ {
			param1 = param1.(*object.Variant).Value
		}
		switch param1 := param1.(type) {
		case *object.Integer:
			return object.NewDoubleByFloat(math.Log(float64(param0.Value))/math.Log(float64(param1.Value)), param0.Position())
		case *object.Long:
			return object.NewDoubleByFloat(math.Log(float64(param0.Value))/math.Log(float64(param1.Value)), param0.Position())
		case *object.Double:
			return object.NewDoubleByFloat(math.Log(float64(param0.Value))/param1.Value, param0.Position())
		case *object.Single:
			return object.NewSingleByFloat(float32(math.Log(float64(param0.Value))/float64(param1.Value)), param0.Position())
		case *object.Currency:
			return object.NewDoubleByFloat(math.Log(float64(param0.Value))/param1.Value, param0.Position())
		default:
			return object.NewError(params[1].Position(), "second argument must be a number")
		}
	case *object.Long:
		if len(params) == 1 {
			return object.NewDoubleByFloat(math.Log(float64(param0.Value)), param0.Position())
		}
		param1 := params[1]
		if param1.Type() == object.VARIANT_OBJ {
			param1 = param1.(*object.Variant).Value
		}
		switch param1 := param1.(type) {
		case *object.Integer:
			return object.NewDoubleByFloat(math.Log(float64(param0.Value))/math.Log(float64(param1.Value)), param0.Position())
		case *object.Long:
			return object.NewDoubleByFloat(math.Log(float64(param0.Value))/math.Log(math.Float64frombits(uint64(param1.Value))), param0.Position())
		}
	}
	return object.NewError(params[0].Position(), "argument must be a number")
}

// evalRnd evaluates the rnd function.
func evalRnd(params []object.Object, env *object.Environment) object.Object {
	if len(params) > 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=0 or 1", len(params)))
	}
	if len(params) == 0 {
		return object.NewDoubleByFloat(rand.Float64(), sem.UniversePos)
	}
	param := params[0]
	if param.Type() == object.VARIANT_OBJ {
		param = param.(*object.Variant).Value
	}
	// the argument must be greater than 0
	if param.Type() != object.LONG_OBJ {
		return object.NewError(params[0].Position(), "argument must be a number")
	}
	if param.(*object.Long).Value <= 0 {
		return object.NewError(params[0].Position(), "argument must be greater than 0")
	}

	switch param := param.(type) {
	case *object.Integer:
		return object.NewIntegerByInt(int32(rand.Intn(int(param.Value))), param.Position())
	case *object.Long:
		return object.NewLongByInt(int64(rand.Int63n(param.Value)), param.Position())
	}
	return object.NewError(params[0].Position(), "argument must be a number")
}

// evalSgn evaluates the sgn function.
func evalSgn(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}
	param := params[0]
	if param.Type() == object.VARIANT_OBJ {
		param = param.(*object.Variant).Value
	}
	switch param := param.(type) {
	case *object.Integer:
		if param.Value > 0 {
			return object.NewIntegerByInt(1, param.Position())
		}
		if param.Value < 0 {
			return object.NewIntegerByInt(-1, param.Position())
		}
		return object.NewIntegerByInt(0, param.Position())
	case *object.Long:
		if param.Value > 0 {
			return object.NewLongByInt(1, param.Position())
		}
		if param.Value < 0 {
			return object.NewLongByInt(-1, param.Position())
		}
		return object.NewLongByInt(0, param.Position())
	case *object.Double:
		if param.Value > 0 {
			return object.NewDoubleByFloat(1, param.Position())
		}
		if param.Value < 0 {
			return object.NewDoubleByFloat(-1, param.Position())
		}
		return object.NewDoubleByFloat(0, param.Position())
	case *object.Single:
		if param.Value > 0 {
			return object.NewSingleByFloat(1, param.Position())
		}
		if param.Value < 0 {
			return object.NewSingleByFloat(-1, param.Position())
		}
		return object.NewSingleByFloat(0, param.Position())
	case *object.Currency:
		if param.Value > 0 {
			return object.NewCurrencyByFloat(1, param.Position())
		}
		if param.Value < 0 {
			return object.NewCurrencyByFloat(-1, param.Position())
		}
		return object.NewCurrencyByFloat(0, param.Position())
	}
	return object.NewError(params[0].Position(), "argument must be a number")
}

// evalSin evaluates the sin function.
func evalSin(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}
	param := params[0]
	if param.Type() == object.VARIANT_OBJ {
		param = param.(*object.Variant).Value
	}
	switch param := param.(type) {
	case *object.Integer:
		return object.NewDoubleByFloat(math.Sin(float64(param.Value)), param.Position())
	case *object.Long:
		return object.NewDoubleByFloat(math.Sin(float64(param.Value)), param.Position())
	case *object.Double:
		return object.NewDoubleByFloat(math.Sin(param.Value), param.Position())
	case *object.Single:
		return object.NewSingleByFloat(float32(math.Sin(float64(param.Value))), param.Position())
	case *object.Currency:
		return object.NewDoubleByFloat(math.Sin(param.Value), param.Position())
	default:
		return object.NewError(params[0].Position(), "argument must be a number")
	}
}

// evalSqr evaluates the sqr function.
func evalSqr(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}
	param := params[0]
	if param.Type() == object.VARIANT_OBJ {
		param = param.(*object.Variant).Value
	}
	switch param := param.(type) {
	case *object.Integer:
		return object.NewDoubleByFloat(math.Sqrt(float64(param.Value)), param.Position())
	case *object.Long:
		return object.NewDoubleByFloat(math.Sqrt(float64(param.Value)), param.Position())
	case *object.Double:
		return object.NewDoubleByFloat(math.Sqrt(param.Value), param.Position())
	case *object.Single:
		return object.NewSingleByFloat(float32(math.Sqrt(float64(param.Value))), param.Position())
	case *object.Currency:
		return object.NewDoubleByFloat(math.Sqrt(param.Value), param.Position())
	default:
		return object.NewError(params[0].Position(), "argument must be a number")
	}
}

// evalTan evaluates the tan function.
func evalTan(params []object.Object, env *object.Environment) object.Object {
	if len(params) != 1 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1", len(params)))
	}
	param := params[0]
	if param.Type() == object.VARIANT_OBJ {
		param = param.(*object.Variant).Value
	}
	switch param := param.(type) {
	case *object.Integer:
		return object.NewDoubleByFloat(math.Tan(float64(param.Value)), param.Position())
	case *object.Long:
		return object.NewDoubleByFloat(math.Tan(float64(param.Value)), param.Position())
	case *object.Double:
		return object.NewDoubleByFloat(math.Tan(param.Value), param.Position())
	case *object.Single:
		return object.NewSingleByFloat(float32(math.Tan(float64(param.Value))), param.Position())
	case *object.Currency:
		return object.NewDoubleByFloat(math.Tan(param.Value), param.Position())
	default:
		return object.NewError(params[0].Position(), "argument must be a number")
	}
}

// evalLBound evaluates the lbound function.
func evalLBound(params []object.Object, env *object.Environment) object.Object {
	// verify the number of arguments
	if len(params) > 2 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1 or 2", len(params)))
	}
	// verify the type of the first argument
	param0 := params[0]
	if param0.Type() == object.VARIANT_OBJ {
		param0 = param0.(*object.Variant).Value
	}
	if param0.Type() != object.ARRAY_OBJ {
		return object.NewError(params[0].Position(), "first argument must be an array")
	}
	// second argument is optional and not used

	// all arrays in uBasic are 0-based
	return object.NewIntegerByInt(0, params[0].Position())
}

// evalUBound evaluates the ubound function.
func evalUBound(params []object.Object, env *object.Environment) object.Object {
	// verify the number of arguments
	if len(params) > 2 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1 or 2", len(params)))
	}
	// verify the type of the first argument
	param0 := params[0]
	if param0.Type() == object.VARIANT_OBJ {
		param0 = param0.(*object.Variant).Value
	}
	if param0.Type() != object.ARRAY_OBJ {
		return object.NewError(params[0].Position(), "first argument must be an array")
	}
	// the second argument must be a number if present
	var param1 object.Object
	if len(params) == 2 {
		param1 = params[1]
		if param1.Type() == object.VARIANT_OBJ {
			param1 = param1.(*object.Variant).Value
		}
		if param1.Type() != object.LONG_OBJ {
			return object.NewError(params[1].Position(), "second argument must be a number")
		}
	}
	// get the array
	array := param0.(*object.Array)
	// get the dimension
	dimension := 0
	if param1 != nil {
		dimension = int(param1.(*object.Long).Value)
	}
	// verify the dimension
	if dimension < 1 || dimension > len(array.Dimensions) {
		return object.NewError(params[0].Position(), "dimension out of range")
	}
	// get the upper bound
	upperBound := array.Dimensions[dimension-1] - 1
	return object.NewLongByInt(int64(upperBound), params[0].Position())

}

// evalInstr evaluates the instr function.
func evalInstr(params []object.Object, env *object.Environment) object.Object {
	param0 := params[0]
	if param0.Type() == object.VARIANT_OBJ {
		param0 = param0.(*object.Variant).Value
	}
	if param0.Type() != object.LONG_OBJ {
		return object.NewError(params[0].Position(), "first argument must be a long")
	}

	param1 := params[1]
	if param1.Type() == object.VARIANT_OBJ {
		param1 = param1.(*object.Variant).Value
	}
	if param1.Type() != object.STRING_OBJ {
		return object.NewError(params[1].Position(), "second argument must be a string")
	}

	param2 := params[2]
	if param2.Type() == object.VARIANT_OBJ {
		param2 = param2.(*object.Variant).Value
	}
	if param2.Type() != object.STRING_OBJ {
		return object.NewError(params[1].Position(), "third argument must be a string")
	}
	start := param0.(*object.Long).Value
	value := param1.(*object.String).Value
	substr := param2.(*object.String).Value

	compare := 0
	if len(params) == 4 {
		param3 := params[3]
		if param3.Type() == object.VARIANT_OBJ {
			param3 = param3.(*object.Variant).Value
		}
		if param3.Type() != object.LONG_OBJ {
			return object.NewError(params[3].Position(), "fourth argument must be a long")
		}
		compare = int(param3.(*object.Long).Value)
	}

	if start == 1 {
		// compare case-sensitive
		if compare == 0 {
			return object.NewIntegerByInt(int32(strings.Index(value, substr)+1), params[0].Position())
		}
		// compare case-insensitive
		return object.NewIntegerByInt(int32(strings.Index(strings.ToLower(value), strings.ToLower(substr))+1), params[0].Position())
	}
	// compare case-sensitive
	if compare == 0 {

		position := strings.Index(value[start-1:], substr)
		if position == -1 {
			return object.NewIntegerByInt(0, params[0].Position())
		}
		return object.NewIntegerByInt(int32(position+int(start)), params[0].Position())
	}
	// compare case-insensitive
	position := strings.Index(strings.ToLower(value[start-1:]), strings.ToLower(substr))
	if position == -1 {
		return object.NewIntegerByInt(0, params[0].Position())
	}
	return object.NewIntegerByInt(int32(position+int(start)), params[0].Position())

}

// evalInput evaluates the input function.
func evalInput(params []object.Object, env *object.Environment) object.Object {
	if len(params) < 1 || len(params) > 2 {
		return object.NewError(params[0].Position(), fmt.Sprintf("wrong number of arguments. got=%d, want=1 or 2", len(params)))
	}
	if params[0].Type() != object.STRING_OBJ {
		return object.NewError(params[0].Position(), "first argument must be a string")
	}
	fmt.Print(params[0].(*object.String).Value)
	if len(params) == 2 {
		if params[1].Type() != object.STRING_OBJ {
			return object.NewError(params[1].Position(), "second argument must be a string")
		}
		fmt.Print(params[1].(*object.String).Value)
	}
	var input string
	fmt.Scanln(&input)
	return object.NewString(input, params[0].Position())
}
