package irgen

import (
	"fmt"
	"testing"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
)

func TestIntGEN(*testing.T) {
	// Create a new LLVM IR module.
	m := ir.NewModule()
	// Create a new global variable of type i32 and name it "ten".
	ten := constant.NewInt(types.I64, 10)
	tenGlobal := m.NewGlobalDef("ten", ten)
	// Create a new global variable of type [15]i8 and name it "str".
	hello := constant.NewCharArrayFromString("Hello, %ld!\n\x00")
	str := m.NewGlobalDef("str", hello)
	// Add external function declaration of printf.
	printf := m.NewFunc("printf", types.I32, ir.NewParam("", types.NewPointer(types.I8)))
	printf.Sig.Variadic = true

	// Create a new function main which returns an i32.
	main := m.NewFunc("main", types.I32)
	entry := main.NewBlock("")
	// Allocate memory for ten pointer.
	tenPtr := entry.NewAlloca(ten.Typ)
	entry.NewStore(ten, tenPtr)
	// store the pointer to ten in tenPtr
	// Cast *[15]i8 to *i8.

	zero := constant.NewInt(types.I64, 0)
	gep := constant.NewGetElementPtr(hello.Typ, str, zero, zero)
	gep10 := constant.NewGetElementPtr(ten.Typ, tenGlobal, zero)
	tmp1 := entry.NewLoad(types.I64, gep10)

	// printf tests
	entry.NewCall(printf, gep, tmp1)

	// Return 0 from main.
	entry.NewRet(constant.NewInt(types.I32, 0))
	fmt.Println(m)
}

func TestStrAllocGEN(*testing.T) {
	// Create a new LLVM IR module.
	m := ir.NewModule()
	// Convenience types and values.
	i32 := types.I32
	i8 := types.I8
	i8ptr := types.NewPointer(i8)

	// create link to stdlib.h
	// add string functions
	strcpy := m.NewFunc("strcpy", i8ptr, ir.NewParam("dst", i8ptr), ir.NewParam("src", i8ptr))
	puts := m.NewFunc("puts", i32, ir.NewParam("s", i8ptr))

	// memory management
	malloc := m.NewFunc("malloc", i8ptr, ir.NewParam("size", types.I32))
	free := m.NewFunc("free", types.Void, ir.NewParam("ptr", i8ptr))

	// Create a global variable of type string.
	str := m.NewGlobal("str", i8ptr)
	str.Init = constant.NewNull(i8ptr)

	// test constant
	constantStr0 := constant.NewCharArrayFromString("Hello, World!\n\x00")
	constantStr1 := m.NewGlobalDef(".str0", constantStr0)

	// ---------------------------------------------------------
	// main ()
	// ---------------------------------------------------------
	// Create a new function main which returns an i32.
	main := m.NewFunc("main", types.I32)
	entry := main.NewBlock("")

	length := int64(len(constantStr0.X))
	// allocate heap memory for global strings
	str2 := entry.NewCall(malloc, constant.NewInt(types.I32, length))
	entry.NewStore(str2, str)

	// copy string to allocated memory
	str3 := entry.NewLoad(i8ptr, str)
	entry.NewCall(strcpy, str3, constantStr1)

	// capture the pointer to the string
	str10 := entry.NewLoad(i8ptr, str)
	// call puts
	entry.NewCall(puts, str10)

	// free memory
	entry.NewCall(free, str10)

	// Return 0 from main.
	entry.NewRet(constant.NewInt(types.I32, 0))
	fmt.Println(m)
}
func TestBoolean1(*testing.T) {
	// Create a new LLVM IR module.
	m := ir.NewModule()

	// global constant
	trueStr0 := constant.NewCharArrayFromString("True\x00")
	trueStr1 := m.NewGlobalDef("true", trueStr0)
	falseStr0 := constant.NewCharArrayFromString("False\x00")
	falseStr1 := m.NewGlobalDef("false", falseStr0)

	// print the result
	printf := m.NewFunc("puts", types.I32, ir.NewParam("str", types.I8Ptr))

	// ---- [ conversion function ] --------------------------------------------------

	// create a function to convert from I1 to string
	param := ir.NewParam("value", types.I1)
	booleanToString := m.NewFunc("_fromCharToStringBoolean_", types.I8Ptr, param)

	entry := booleanToString.NewBlock("")
	// create a switch to check the value of the boolean

	// create a new block for the true case
	trueBlock := booleanToString.NewBlock("")
	// create a new block for the false case
	falseBlock := booleanToString.NewBlock("")
	// create a new block to return result
	endBlock := booleanToString.NewBlock("")

	// create a switch to check the value of the boolean
	cmp := entry.NewICmp(enum.IPredEQ, booleanToString.Params[0], constant.NewInt(types.I1, 1))
	entry.NewCondBr(cmp, trueBlock, falseBlock)

	// true block - load true value

	gep1 := constant.NewGetElementPtr(trueStr0.Typ, trueStr1, constant.NewInt(types.I64, 0))
	//v1 := trueBlock.NewLoad(types.I8, gep1)
	trueBlock.NewBr(endBlock)

	// false block - load false value

	gep2 := constant.NewGetElementPtr(falseStr0.Typ, falseStr1, constant.NewInt(types.I64, 0))
	//v2 := falseBlock.NewLoad(types.I8, gep2)
	falseBlock.NewBr(endBlock)

	// end block - return value
	v3 := endBlock.NewPhi(ir.NewIncoming(gep1, trueBlock), ir.NewIncoming(gep2, falseBlock))
	endBlock.NewRet(v3)
	// ---- [ main ] --------------------------------------------------

	main := m.NewFunc("main", types.I32)
	entry = main.NewBlock("")
	// call the function with 1
	tmp1 := entry.NewCall(booleanToString, constant.NewInt(types.I1, 1))
	// tmp2 := entry.NewLoad(types.I8Ptr, tmp1)
	entry.NewCall(printf, tmp1)

	// call the function with 1
	tmp3 := entry.NewCall(booleanToString, constant.NewInt(types.I1, 0))
	// tmp4 := entry.NewLoad(types.I8Ptr, tmp3)
	entry.NewCall(printf, tmp3)

	// Return 0 from main.
	entry.NewRet(constant.NewInt(types.I32, 0))
	fmt.Println(m)
}

func TestBoolean2(*testing.T) {
	// Create a new LLVM IR module.
	m := ir.NewModule()

	// global constant
	trueStr0 := constant.NewCharArrayFromString("True\x00")
	trueStr1 := m.NewGlobalDef("true", trueStr0)
	falseStr0 := constant.NewCharArrayFromString("False\x00")
	falseStr1 := m.NewGlobalDef("false", falseStr0)

	// print the result
	printf := m.NewFunc("puts", types.I32, ir.NewParam("str", types.I8Ptr))

	// ---- [ conversion function ] --------------------------------------------------

	// create a function to convert from I1 to string
	param := ir.NewParam("value", types.I1)
	booleanToString := m.NewFunc("_fromCharToStringBoolean_", types.I8Ptr, param)

	entry := booleanToString.NewBlock("")

	// create a switch to check the value of the boolean
	cmp := entry.NewICmp(enum.IPredEQ, booleanToString.Params[0], constant.NewInt(types.I1, 1))
	res := entry.NewSelect(cmp,
		constant.NewGetElementPtr(trueStr0.Typ, trueStr1, constant.NewInt(types.I64, 0)),
		constant.NewGetElementPtr(falseStr0.Typ, falseStr1, constant.NewInt(types.I64, 0)))

	entry.NewRet(res)
	// ---- [ main ] --------------------------------------------------

	main := m.NewFunc("main", types.I32)
	entry = main.NewBlock("")
	// call the function with 1
	tmp1 := entry.NewCall(booleanToString, constant.NewInt(types.I1, 1))
	// tmp2 := entry.NewLoad(types.I8Ptr, tmp1)
	entry.NewCall(printf, tmp1)

	// call the function with 1
	tmp3 := entry.NewCall(booleanToString, constant.NewInt(types.I1, 0))
	// tmp4 := entry.NewLoad(types.I8Ptr, tmp3)
	entry.NewCall(printf, tmp3)

	// Return 0 from main.
	entry.NewRet(constant.NewInt(types.I32, 0))
	fmt.Println(m)
}

func TestIfGEN(*testing.T) {
	// Create a new LLVM IR module.
	m := ir.NewModule()
	// Create a new global variable of type i32 and name it "ten".
	ten := constant.NewInt(types.I64, 10)
	tenGlobal := m.NewGlobal("ten", ten.Typ)
	tenGlobal.Init = ten

	// Create a new global variable of type [15]i8 and name it "str".
	hello := constant.NewCharArrayFromString("Value: %ld!\n\x00")
	str := m.NewGlobalDef("str", hello)
	// Add external function declaration of printf.
	printf := m.NewFunc("printf", types.I32, ir.NewParam("", types.NewPointer(types.I8)))
	printf.Sig.Variadic = true

	// --------------------------------------------------------
	// main ()
	// --------------------------------------------------------
	// Create a new function main which returns an i32.
	main := m.NewFunc("main", types.I32)
	entry := main.NewBlock("")

	// create a new block for the true case
	ifTrue := main.NewBlock("")
	// create a new block for the false case
	ifFalse := main.NewBlock("")
	// create a new block to return result
	ifEnd := main.NewBlock("")

	// if ten == 1 then
	//     ten = ten - 1
	// else
	//     ten = ten + 1
	// end if

	// load the value of ten
	tmp1 := entry.NewLoad(types.I64, tenGlobal)
	// create a switch to check the value of the boolean
	cmp := entry.NewICmp(enum.IPredEQ, tmp1, constant.NewInt(types.I64, 1))
	entry.NewCondBr(cmp, ifTrue, ifFalse)

	// True - substract 1 from ten
	tmp2 := ifTrue.NewSub(tmp1, constant.NewInt(types.I64, 1))
	ifTrue.NewBr(ifEnd)

	// False - add 1 to ten
	tmp4 := ifFalse.NewAdd(tmp1, constant.NewInt(types.I64, 1))
	ifFalse.NewBr(ifEnd)

	// End if - return value
	tmp5 := ifEnd.NewPhi(ir.NewIncoming(tmp2, ifTrue), ir.NewIncoming(tmp4, ifFalse))
	ifEnd.NewStore(tmp5, tenGlobal)

	// get pointers for printf
	zero := constant.NewInt(types.I64, 0)
	gep := constant.NewGetElementPtr(hello.Typ, str, zero, zero)
	gep10 := constant.NewGetElementPtr(ten.Typ, tenGlobal, zero)
	tmp10 := ifEnd.NewLoad(types.I64, gep10)

	// printf tests
	ifEnd.NewCall(printf, gep, tmp10)

	// Return 0 from main.
	ifEnd.NewRet(constant.NewInt(types.I32, 0))
	fmt.Println(m)
}

func TestMultFloat(*testing.T) {
	// Create a new LLVM IR module.
	m := ir.NewModule()
	// --------------------------------------------------------
	// main ()
	// --------------------------------------------------------
	main := m.NewFunc("main", types.I32)
	entry := main.NewBlock("")

	tmp0 := entry.NewAlloca(types.Double)
	entry.NewStore(constant.NewFloat(types.Double, 10.9), tmp0)
	tmp1 := entry.NewAlloca(types.Double)
	entry.NewStore(constant.NewFloat(types.Double, 0.98), tmp1)
	entry.NewMul(entry.NewLoad(types.Double, tmp0), entry.NewLoad(types.Double, tmp1))

	entry.NewRet(constant.NewInt(types.I32, 0))
	fmt.Println(m)
}

func TestArrays(t *testing.T) {
	// Create a new LLVM IR module.
	m := ir.NewModule()
	// Create a new global variable of type [15]i8 and name it "str".
	array0 := constant.NewArray(&types.ArrayType{Len: 2, ElemType: types.I32})
	init := constant.NewZeroInitializer(array0.Typ)
	array := m.NewGlobalDef("array", init)

	// Add external function declaration of printf.
	printf := m.NewFunc("printf", types.I32, ir.NewParam("", types.NewPointer(types.I8)))
	printf.Sig.Variadic = true

	// Create a new global variable of type [15]i8 and name it "str".
	format := constant.NewCharArrayFromString("==> %d\n\x00")
	formatstr := m.NewGlobalDef("str", format)
	// --------------------------------------------------------
	// main ()
	// --------------------------------------------------------
	// Create a new function main which returns an i32.
	main := m.NewFunc("main", types.I32)
	entry := main.NewBlock("")

	// initialize array
	zero := constant.NewInt(types.I64, 0)
	one := constant.NewInt(types.I64, 1)
	gep0 := constant.NewGetElementPtr(array0.Typ, array, zero, zero)
	entry.NewStore(constant.NewInt(types.I32, 10), gep0)
	gep0 = constant.NewGetElementPtr(array0.Typ, array, zero, one)
	entry.NewStore(constant.NewInt(types.I32, 20), gep0)

	// get pointers for printf
	gep := constant.NewGetElementPtr(format.Typ, formatstr, zero, zero)

	// printf array[0]
	gep1 := constant.NewGetElementPtr(array0.Typ, array, zero, zero)
	tmp1 := entry.NewLoad(types.I32, gep1)
	entry.NewCall(printf, gep, tmp1)

	// printf array[1]
	gep2 := constant.NewGetElementPtr(array0.Typ, array, zero, one)
	tmp2 := entry.NewLoad(types.I32, gep2)
	entry.NewCall(printf, gep, tmp2)

	// printf array[0] + array[1]
	gep3 := constant.NewGetElementPtr(array0.Typ, array, zero, zero)
	tmp3 := entry.NewLoad(types.I32, gep3)
	gep4 := constant.NewGetElementPtr(array0.Typ, array, zero, one)
	tmp4 := entry.NewLoad(types.I32, gep4)
	tmp5 := entry.NewAdd(tmp3, tmp4)
	entry.NewCall(printf, gep, tmp5)

	// Return add from main.
	entry.NewRet(tmp5)
	fmt.Println(m)
}

func TestDynamicArrays(t *testing.T) {
	// Create a new LLVM IR module.
	m := ir.NewModule()
	// Create a new global variable of type [15]i8 and name it "str".
	array0 := constant.NewArray(&types.ArrayType{Len: 2, ElemType: types.I32})
	init := constant.NewZeroInitializer(array0.Typ)
	array := m.NewGlobalDef("array", init)

	// Add external function declaration of printf.
	printf := m.NewFunc("printf", types.I32, ir.NewParam("", types.NewPointer(types.I8)))
	printf.Sig.Variadic = true

	// Create a new global variable of type [15]i8 and name it "str".
	format := constant.NewCharArrayFromString("==> %d\n\x00")
	formatstr := m.NewGlobalDef("str", format)
	// --------------------------------------------------------
	// main ()
	// --------------------------------------------------------
	// Create a new function main which returns an i32.
	main := m.NewFunc("main", types.I32)
	entry := main.NewBlock("")

	// initialize array
	zero := constant.NewInt(types.I64, 0)
	one := constant.NewInt(types.I64, 1)
	gep0 := constant.NewGetElementPtr(array0.Typ, array, zero, zero)
	entry.NewStore(constant.NewInt(types.I32, 10), gep0)
	gep0 = constant.NewGetElementPtr(array0.Typ, array, zero, one)
	entry.NewStore(constant.NewInt(types.I32, 20), gep0)

	// get pointers for printf
	gep := constant.NewGetElementPtr(format.Typ, formatstr, zero, zero)

	// printf array[0]
	gep1 := constant.NewGetElementPtr(array0.Typ, array, zero, zero)
	tmp1 := entry.NewLoad(types.I32, gep1)
	entry.NewCall(printf, gep, tmp1)

	// printf array[1]
	gep2 := constant.NewGetElementPtr(array0.Typ, array, zero, one)
	tmp2 := entry.NewLoad(types.I32, gep2)
	entry.NewCall(printf, gep, tmp2)

	// printf array[0] + array[1]
	gep3 := constant.NewGetElementPtr(array0.Typ, array, zero, zero)
	tmp3 := entry.NewLoad(types.I32, gep3)
	gep4 := constant.NewGetElementPtr(array0.Typ, array, zero, one)
	tmp4 := entry.NewLoad(types.I32, gep4)
	tmp5 := entry.NewAdd(tmp3, tmp4)
	entry.NewCall(printf, gep, tmp5)

	// Return add from main.
	entry.NewRet(tmp5)
	fmt.Println(m)
}

func TestException(t *testing.T) {
	// Create a new LLVM IR module.
	m := ir.NewModule()

	// exception handling global variables
	jump_bufferType := types.NewArray(48, types.I32)
	jump_buffer := m.NewGlobal("jump_buffer", jump_bufferType)
	jump_buffer.Init = constant.NewZeroInitializer(jump_bufferType)
	errorNumber := m.NewGlobal("errorNumber", types.I32)
	errorNumber.Init = constant.NewInt(types.I32, 0)
	emType := types.NewArray(256, types.I8)
	errorMessage := m.NewGlobal("errorMessage", emType)
	errorMessage.Init = constant.NewZeroInitializer(emType)

	// constants
	constantStr0 := constant.NewCharArrayFromString("Division by zero: %d\n\x00")
	constantStr1 := m.NewGlobalDef(".str0", constantStr0)

	// exception handling functions
	// setjmp and longjmp
	setjmp := m.NewFunc("setjmp", types.I32, ir.NewParam("", types.NewPointer(types.I32)))
	longjmp := m.NewFunc("longjmp", types.Void, ir.NewParam("", types.NewPointer(types.I32)), ir.NewParam("", types.I32))
	// printf
	printf := m.NewFunc("printf", types.I32, ir.NewParam("", types.NewPointer(types.I8)))
	printf.Sig.Variadic = true
	// strcpy
	strcpy := m.NewFunc("strcpy", types.I8Ptr, ir.NewParam("dst", types.I8Ptr), ir.NewParam("src", types.I8Ptr))

	// throw exception function
	throwException := m.NewFunc("throwException", types.Void)
	entry0 := throwException.NewBlock("")
	entry0.NewCall(longjmp, jump_buffer, constant.NewInt(types.I32, 1))
	entry0.NewUnreachable()

	// function function_that_might_throw_exception
	function_that_might_throw_exception := m.NewFunc("function_that_might_throw_exception", types.Void)
	entry1 := function_that_might_throw_exception.NewBlock("")
	// set error message
	entry1.NewCall(strcpy, errorMessage, constantStr1)
	// set error number
	entry1.NewStore(constant.NewInt(types.I32, 17), errorNumber)
	entry1.NewCall(throwException)
	entry1.NewRet(nil)

	// main function
	main := m.NewFunc("main", types.I32)
	entry := main.NewBlock("")
	exception := main.NewBlock("exception")
	normalCall := main.NewBlock("normalCall")
	end := main.NewBlock("end")

	// entry:
	tmp2 := entry.NewCall(setjmp, jump_buffer)
	cmp := entry.NewICmp(enum.IPredEQ, tmp2, constant.NewInt(types.I32, 0))
	entry.NewCondBr(cmp, normalCall, exception)

	// normalCall:
	normalCall.NewCall(function_that_might_throw_exception)
	normalCall.NewBr(end)

	// exception:
	// em2 = exception.NewLoad(types.I8Ptr, errorMessage)
	en := exception.NewLoad(types.I32, errorNumber)
	exception.NewCall(printf, errorMessage, en)
	exception.NewRet(en)

	// end:
	end.NewRet(constant.NewInt(types.I32, 0))
	fmt.Println(m)
}

func TestParamArray(t *testing.T) {
	// Create a new LLVM IR module.
	m := ir.NewModule()

	// Add external function declaration of printf.
	printf := m.NewFunc("printf", types.I32, ir.NewParam("", types.NewPointer(types.I8)))
	printf.Sig.Variadic = true

	// Create a new global variable of type [15]i8 and name it "str".
	format := constant.NewCharArrayFromString("==> %d \n\x00")
	formatstr := m.NewGlobalDef("str", format)

	// --------------------------------------------------------
	// Create a new function varia which returns an i32
	// that adds the value of an array.
	varia := m.NewFunc("varia", types.I32, ir.NewParam(".values_0)", types.I32), ir.NewParam("values", types.I32Ptr))

	// define branching labels
	entry := varia.NewBlock("entry")
	loopCond := varia.NewBlock("loop.cond")
	loopBody := varia.NewBlock("loop.body")
	loopEnd := varia.NewBlock("loop.end")

	// allocate variables
	total := entry.NewAlloca(types.I32)
	total.SetName("total")
	count := entry.NewAlloca(types.I32)
	count.SetName("count")

	i := entry.NewAlloca(types.I32)
	i.SetName("i")

	// read parameters
	entry.NewStore(varia.Params[0], count)

	// initialize i
	entry.NewStore(constant.NewInt(types.I32, 0), i)

	// initialize total
	entry.NewStore(constant.NewInt(types.I32, 0), total)
	entry.NewBr(loopCond)

	// loop condition
	iVal := loopCond.NewLoad(types.I32, i)
	countVal := loopCond.NewLoad(types.I32, count)
	cond := loopCond.NewICmp(enum.IPredULT, iVal, countVal)
	loopCond.NewCondBr(cond, loopBody, loopEnd)

	// loop body
	gep := loopBody.NewGetElementPtr(types.I32, varia.Params[1], iVal)
	tmp := loopBody.NewLoad(types.I32, gep)
	totalVal := loopBody.NewLoad(types.I32, total)
	tmp1 := loopBody.NewAdd(tmp, totalVal)
	loopBody.NewStore(tmp1, total)
	iVal = loopBody.NewLoad(types.I32, i)
	inc := loopBody.NewAdd(iVal, constant.NewInt(types.I32, 1))
	loopBody.NewStore(inc, i)
	loopBody.NewBr(loopCond)

	// loop end
	loopEnd.NewRet(loopEnd.NewLoad(types.I32, total))

	// --------------------------------------------------------
	// main ()
	// --------------------------------------------------------
	// Create a new function main which returns an i32.
	main := m.NewFunc("main", types.I32)
	entryMain := main.NewBlock("")
	countMain := constant.NewInt(types.I32, 2)
	valueType := types.ArrayType{Len: 2, ElemType: types.I32}
	values := entryMain.NewAlloca(&valueType)

	// initialize array
	zero := constant.NewInt(types.I64, 0)
	gep1 := entryMain.NewGetElementPtr(types.I32, values, zero)
	entryMain.NewStore(constant.NewInt(types.I32, 10), gep1)

	one := constant.NewInt(types.I64, 1)
	gep2 := entryMain.NewGetElementPtr(types.I32, values, one)
	entryMain.NewStore(constant.NewInt(types.I32, 20), gep2)

	// print value
	value1 := entryMain.NewLoad(types.I32, gep1)
	entryMain.NewCall(printf, formatstr, value1)
	value2 := entryMain.NewLoad(types.I32, gep2)
	entryMain.NewCall(printf, formatstr, value2)

	// call varia
	tmpVaria := entryMain.NewCall(varia, countMain, values)

	// get pointers for printf
	entryMain.NewCall(printf, formatstr, tmpVaria)

	// Return 0 from main.
	entryMain.NewRet(constant.NewInt(types.I32, 0))
	fmt.Println(m)
}
