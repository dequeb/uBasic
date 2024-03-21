package eval

import (
	"fmt"
	"os"
	"testing"
	"uBasic/errors"
	"uBasic/lexer"
	"uBasic/object"
	"uBasic/parser"
	"uBasic/sem"
	"uBasic/source"
)

func TestEvalBooleanExpression(t *testing.T) {
	test := []struct {
		input    string
		expected bool
	}{
		{"Print True", true},
		{"Print False", false},
		{"Print Not True ", false},
		{"Print Not False ", true},
		{"Print Not (1 == 2)  ", true},
		{"Print True == True", true},
		{"Print True == False", false},
		{"Print True <>True", false},
		{"Print True <> False", true},
		{"Print 1.05$ == 1.05", true},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}
	return true
}

func TestEvalLongExpression(t *testing.T) {
	test := []struct {
		input    string
		expected int64
	}{
		{"Print 5", 5},
		{"Print 10", 10},
		{"Print -5", -5},
		{"Print -10", -10},
		{"Print 5 + 5 + 5 + 5 - 10", 10},
		{"Print 2 * 2 * 2 * 2 * 2", 32},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testLongObject(t, evaluated, tt.expected)
	}
}

func TestLongEvalDim(t *testing.T) {
	test := []struct {
		input    string
		expected int64
	}{
		{"Dim a As Long:Let a = 5:Print a", 5},
		{"Dim a As Long\nLet a = 5\nPrint a\nLet a = 10\nPrint a", 10},
		{"Dim a As Long\nLet a = 5*7\nPrint a\nDim b As Long\nLet b = 10\nPrint a+b", 45},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testLongObject(t, evaluated, tt.expected)
	}
}
func testLongObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Long)
	if !ok {
		t.Errorf("object is not Long. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}

func TestEvalDoubleLongExpression(t *testing.T) {
	test := []struct {
		input    string
		expected float64
	}{
		{"Print 5.5 + 5$", 10.5},
		{"Print 10.5 + 5", 15.5},
		{"Print -5.5 + 5", -0.5},
		{"Print -10.5 + 5", -5.5},
		{"Print 5.5 + 5.5 + 5.5 + 5 - 10.5", 11},
		{"Print 2.5 * 2.5$ * 2.5 * 2.5 * 2.5", 97.65625},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testDoubleObject(t, evaluated, tt.expected)
	}
}
func TestDoubleEvalDim(t *testing.T) {
	test := []struct {
		input    string
		expected float64
	}{
		{"Dim a As Double:Let a = 5:Print a", 5.0},
		{"Dim a As Double:Let a = 5:Print a:Let a = a+5.5:Print a", 10.5},
		{"Dim a As Double:Let a = 5*7:Print a\nDim b As Double\nLet b = 10 * a\nPrint a+b", 385.0},
		{"Dim a As Double:Let a = 5 EXP 2:Print a", 25.0},
		{"Dim a As Double:Let a = 5.0/2:Print a", 2.5},
		{"Dim a As Double:Let a = 5 MOD 2:Print a", 1.0},
		{"Dim a As Double:Let a = 5.0 DIV 2:Print a", 2.0},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testDoubleObject(t, evaluated, tt.expected)
	}
}

func TestEvalDoubleExpression(t *testing.T) {
	test := []struct {
		input    string
		expected float64
	}{
		{"Print 5.5", 5.5},
		{"Print 10.5", 10.5},
		{"Print -5.5", -5.5},
		{"Print -10.5", -10.5},
		{"Print 5.5 + 5.5 + 5.5 + 5.5 - 10.5", 11.5},
		{"Print 2.5 * 2.5 * 2.5 * 2.5 * 2.5", 97.65625},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testDoubleObject(t, evaluated, tt.expected)
	}
}

func testDoubleObject(t *testing.T, obj object.Object, expected float64) bool {
	result, ok := obj.(*object.Double)
	if !ok {
		t.Errorf("object is not Double. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%f, want=%f", result.Value, expected)
		return false
	}
	return true
}

func TestEvalStringExpression(t *testing.T) {
	test := []struct {
		input    string
		expected string
	}{
		{`Print "Hello World"`, `Hello World`},
		{`Print "Hello" & " " & "World"`, `Hello World`},
		{`Dim s as string: let s= "Hello" & " " & "World": print s`, `Hello World`},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}

func TestEvalNothingExpression(t *testing.T) {
	test := []struct {
		input string
	}{
		{`Print `},
		{`Print ;`},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		_, ok := evaluated.(*object.Nothing)
		if !ok {
			t.Errorf("object is not Nothing. got=%T (%+v)", evaluated, evaluated)
		}
	}
}
func TestEvalDateExpression(t *testing.T) {
	test := []struct {
		input    string
		expected string
	}{
		{`Print #2019-01-02 00:00:00#`, "2019-01-02"},
		{`Print #00:00:00#`, "0000-01-01"},
		{`Print #2019-01-28 10:11:12#`, "2019-01-28 10:11:12"},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testDateObject(t, evaluated, tt.expected)
	}
}

func testDateObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.Date)
	if !ok {
		t.Errorf("object is not Date. got=%T (%+v)", obj, obj)
		return false
	}
	if result.String() != expected {
		t.Errorf("object has wrong value. got=%s, want=%s", result.String(), expected)
		return false
	}
	return true
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	file := p.ParseFile()
	src := &source.Source{Input: input, Name: "testEval"}
	if file == nil {
		for _, err := range p.Errors() {
			e := err.(*errors.Error)
			e.Source = src
			fmt.Println(err)
		}
	} else {
		info, err := sem.Check(file)
		if err != nil {
			e := err.(*errors.Error)
			e.Source = src
			fmt.Println(err)
		} else {
			env := Define(info, os.Stdin, os.Stdout)
			obj := Run(file, env, nil)
			error, ok := obj.(*object.Error)
			if ok {
				error.Stack.Source = src
				fmt.Println(obj)
				return nil
			}
			return obj
		}
	}
	return nil
}

func TestFunctionObject(t *testing.T) {
	input := `Function add(x As Long, y As Long) As Long
		Let add= x + y
	End Function`
	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}
	if len(fn.Parameters.Params) != 2 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
	}
	if fn.Parameters.Params[0].String() != "x As Long" {
		t.Fatalf("parameter is not Long. got=%q", fn.Parameters.Params[0].String())
	}
	if fn.Parameters.Params[1].String() != "y As Long" {
		t.Fatalf("parameter is not Long. got=%q", fn.Parameters.Params[1].String())
	}
	if fn.Parameters.Result.String() != "Long" {
		t.Fatalf("result is not Long. got=%q", fn.Parameters.Result.String())
	}
	if fn.Body[0].String() != "Let add = x + y\n" {
		t.Fatalf("body is not Let add= x + y. got=%q", fn.Body[0].String())
	}
}

func TestEvalFunctionExpression(t *testing.T) {
	test := []struct {
		input    string
		expected int64
	}{
		{`Function add(byval x As Long, byval y As Long) As Long
			Let add= x + y
		End Function
		Print add(5, 5)`, 10},
		{`Function e(byval x As Long, byval y As Long) As Long
			Let e= x EXP y
		End Function
		Print e(5, 2)`, 25},
		{`Function mul(byval x As Long, byval y As Long) As Long
			Let mul= x * y
		End Function
		Print mul(5, 5) + mul(5, 5)`, 50},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testLongObject(t, evaluated, tt.expected)
	}
}

func TestEvalSubExpression(t *testing.T) {
	test := []struct {
		input    string
		expected int64
	}{
		{`sub add(byval x As Long, byval y As Long, result as long) 
			Let result = x + y
			Let x = 20
		End sub
		Dim a As Long
		dim b As Long
		let b = 5
		call add(b, 5, a)
		Print a+b
		'print add(b, 5, a)`, 15},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testLongObject(t, evaluated, tt.expected)
	}
}

func TestEvalArray(t *testing.T) {
	test := []struct {
		input    string
		expected int64
	}{
		{`Dim a(5) As Long
		Let a(0) = 5
		Print a(0)`, 5},
		{`Dim a(5) As Long
		Let a(0) = 5
		Let a(1) = 10
		Print a(0) + a(1)`, 15},
		{`Dim a(6) As Long
		Let a(0) = 5
		Let a(1) = -10
		Let a(2) = 15
		Let a(3) = 20
		Let a(4) = 25
		Let a(5) = 30
		Print a(0) + a(1) + a(2) + a(3) + a(4) + a(5)`, 85},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testLongObject(t, evaluated, tt.expected)
	}
}

func TestEvalVariantLong(t *testing.T) {
	test := []struct {
		input    string
		expected int64
	}{
		{`Dim a As Variant
		Let a = 5
		Print a`, 5},
		{`Dim a As Variant
		Let a = 5
		Let a = a + 5
		Print a`, 10},
		{`Dim a As Variant
		Let a = 5
		Let a = a + 5
		Let a = a + 5
		Print a`, 15},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testVariantLongObject(t, evaluated, tt.expected)
	}
}

func testVariantLongObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Variant)
	if !ok {
		t.Errorf("object is not Variant. got=%T (%+v)", obj, obj)
		return false
	}
	number, ok := result.Value.(*object.Long)
	if !ok {
		t.Errorf("object is not Long. got=%T (%+v)", result.Value, result.Value)
		return false
	}
	if number.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}

func TestEvalVariantString(t *testing.T) {
	test := []struct {
		input    string
		expected string
	}{
		{`Dim a As Variant
		Let a = "Hello"
		Print a`, "Hello"},
		{`Dim a As Variant
		Let a = "Hello"
		Let a = a & " World"
		Print a`, "Hello World"},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testVariantStringObject(t, evaluated, tt.expected)
	}
}

func testVariantStringObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.Variant)
	if !ok {
		t.Errorf("object is not Variant. got=%T (%+v)", obj, obj)
		return false
	}
	number, ok := result.Value.(*object.String)
	if !ok {
		t.Errorf("object is not String. got=%T (%+v)", result.Value, result.Value)
		return false
	}
	if number.Value != expected {
		t.Errorf("object has wrong value. got=%s, want=%s", result.Value, expected)
		return false
	}
	return true
}

func TestEnum(t *testing.T) {
	test := []struct {
		input    string
		expected string
	}{
		{`Enum Color
		Red
		Green
		Blue
		End Enum
		Print Color.Red`, "Red"},
		{`Enum Color
		Red
		Green
		Blue
		End Enum
		Print Color.Green`, "Green"},
		{`Enum Color
		Red
		Green
		Blue
		End Enum
		Print Color.Blue`, "Blue"},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testEnumObject(t, evaluated, tt.expected)
	}
}

func testEnumObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.UserDefined)
	if !ok {
		t.Errorf("object is not User Defined. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%s, want=%s", result.Value, expected)
		return false
	}
	return true
}

func TestCallSelectorExpr(t *testing.T) {
	test := []struct {
		input    string
		expected string
	}{
		{`Print Application.Name`, "uBasic"},
		{`Print Application.User`, "user1"},
		{`call Application.GetOS()`, "macOS"},
		{`Print Application.GetOS()`, "macOS"},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}

func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.String)
	if !ok {
		t.Errorf("object is not String. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%s, want=%s", result.Value, expected)
		return false
	}
	return true
}

func TestEnumCompare(t *testing.T) {
	test := []struct {
		input    string
		expected bool
	}{
		{`Enum Color
		Red
		Green
		Blue
		End Enum
		Print Color.Red == Color.Red`, true},
		{`Enum Color
		Red
		Green
		Blue
		End Enum
		Print Color.Red == Color.Green`, false},
		{`Enum Color
		Red
		Green
		Blue
		End Enum
		Print Color.Red <> Color.Green`, true},
		{`Enum Color
		Red
		Green
		Blue
		End Enum: Print Color.Red <> Color.Red`, false},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestForNextLoop(t *testing.T) {
	test := []struct {
		input    string
		expected int64
	}{
		{`Dim a As Long
		For a = 1 To 5 step 1
			Print a
		Next a
		Print a`, 6},
		{`Dim a As Long
		For a = 1 To 5
			Print a
		Next
		Print a`, 6},
		{`Dim a As Long
		For a = 1+0 To 10 div 2 Step 2
			Print a
		Next
		Print a`, 7},
		{`Dim a As Long
		For a = 5 To 1 Step -1
			Print a
		Next
		Print a`, 0},
		{`Dim a As Long
		For a = 5 To 5 
			Print a
		Next
		Print a`, 0},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testLongObject(t, evaluated, tt.expected)
	}
}

func TestForEachLoop(t *testing.T) {
	test := []struct {
		input    string
		expected int64
	}{
		{`Dim a(5) As Long
		Let a(0) = 5
		Let a(1) = 10
		Let a(2) = 15
		Let a(3) = 20
		Let a(4) = 25
		Let a(5) = 30
		Dim b As Long
		For Each b In a
			Print b
		Next
		Print b`, 30},
		{`Dim a as Variant
		dim b() as long
		for each a in b
			print a
		next
		Print 0`, 0},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testLongObject(t, evaluated, tt.expected)
	}
}

// test arrays
func TestArray(t *testing.T) {
	test := []struct {
		input    string
		expected int64
	}{
		{`Dim a(5) As Long
		Let a(0) = 5
		Print a(0)`, 5},
		{`Dim a(5) As Long
		Let a(0) = 5
		Let a(1) = a(0)
		Let a(0) = 10
		Print a(0) + a(1)`, 15},
		{`Dim a(6, 6) As Long
		Let a(0, 0) = 5
		Let a(1,1) = -10
		Let a(2,2) = 15
		Let a(3,3) = 20
		Let a(4,4) = 25
		Let a(5,5) = 30
		Print a(0,0) + a(1,1) + a(2,2) + a(3,3) + a(4,4) + a(5,5)`, 85},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testLongObject(t, evaluated, tt.expected)
	}
}

// test variant
func TestVariant(t *testing.T) {
	test := []struct {
		input    string
		expected int64
	}{
		{`Dim a As Variant
		Let a = FALSE
		print a
		Let a = 5
		Print a`, 5},
		{`Dim a As Variant
		Let a = "allo"
		Print a
		Let a = 5$
		Print a
		Let a = 5
		Let a = a + 5
		Print a`, 10},
		{`Dim a As Variant
		Let a = 5
		Let a = a + 5
		Let a = a + 5
		Print a`, 15},
		{`Dim a As Variant, b as variant
		Let a = TRUE
		Let b = FALSE
		print a OR b
		Let a = 5
		Let a = 5 EXP 2
		Print a`, 25},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testVariantLongObject(t, evaluated, tt.expected)
	}
}

func TestSelectCase(t *testing.T) {
	test := []struct {
		input    string
		expected int64
	}{
		{`Dim a As Long
		Let a = 5
		Select Case a
			Case 1
				Print 11
			Case 2
				Print 22
			Case 3
				Print 33
			Case 4
				Print 44
			Case 5
				Print 55
			Case Else
				Print 66
		End Select`, 55},
		{`Dim a As Long
		Let a = 6
		Select Case a
			Case 1
				Print 11
			Case 2
				Print 22
			Case 3
				Print 33
			Case 4
				Print 44
			Case 5
				Print 55
			Case Else
				Print 66
		End Select`, 66},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testLongObject(t, evaluated, tt.expected)
	}
}

// interrupting for loop
func TestInterruptForLoop(t *testing.T) {
	test := []struct {
		input    string
		expected int64
	}{
		{`Dim a As Long
		For a = 1 To 5
			Print a
			If a == 3 Then
				Exit For
			end if
		Next
		Print a`, 3},
		{`Dim a As Long
		For a = 1 To 5
			Print a
			If a == 3 Then
				Exit For
			Else
				Print "hello"
			End If
		Next
		Print a`, 3},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testLongObject(t, evaluated, tt.expected)
	}
}

// test interrupting while loop
func TestInterruptWhileLoop(t *testing.T) {
	test := []struct {
		input    string
		expected int64
	}{
		{`Dim a As Long
		Let a = 1
		do While a < 5
			Print a
			If a == 3 Then
				Exit do
			End If
			Let a = a + 1
		Loop
		Print a`, 3},
		{`Dim a As Long
		Let a = 1
		do until a == 5
			Print a
			If a == 3 Then
				Exit do
			Else
				Print "hello"
			End If
			Let a = a + 1
		loop
		Print a`, 3},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testLongObject(t, evaluated, tt.expected)
	}
}

// first interactive test
func TestInteractive(t *testing.T) {
	test := []struct {
		input    string
		expected string
	}{
		{`
Dim name as String
Input "What is your name? ", name
Print "Hello " & name
`, "Hello Michel"},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}

// test some string functions
func TestStringFunctions(t *testing.T) {
	test := []struct {
		input    string
		expected string
	}{
		{`Print CStr(Len("Hello"))`, "5"},
		{`Print CStr(Len("Hello " & "World"))`, "11"},
		{`Print LCase("Hello")`, "hello"},
		{`Print LCase("Hello " & "World")`, "hello world"},
		{`Print UCase("Hello")`, "HELLO"},
		{`Print UCase("Hello " & "World")`, "HELLO WORLD"},
		{`Print Left("Hello", 2)`, "He"},
		{`Print Left("Hello " & "World", 5)`, "Hello"},
		{`Print Right("Hello", 2)`, "lo"},
		{`Print Right("Hello " & "World", 5)`, "World"},
		{`Print Mid("Hello", 2, 2)`, "el"},
		{`Print Mid("Hello " & "World", 7, 5)`, "World"},
		{`Print Mid("Hello " & "World", 7)`, "World"},
		{`Print Mid("Hello " & "World", -3)`, "rld"},
		{`Print Mid("Hello " & "World", 6, 10)`, " World"},
		{`Print Mid("Hello " & "World", 6, 15)`, " World"},
		{`Print Mid("Hello " & "World", 5, 0)`, ""},
		{`Print Mid("Hello " & "World", 5, -1)`, "o World"},
		{`Print Mid("Hello " & "World", -5, 3)`, "Wor"},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}

// test chr, instr, ltrim, rtrim, space, strcomp, string, strreverse, trim, ucase
func TestStringFunctions2(t *testing.T) {
	test := []struct {
		input    string
		expected string
	}{
		{`Print Chr(65)`, "A"},
		{`Print Chr(65) & Chr(66) & Chr(67)`, "ABC"},
		{`Print Cstr(InStr(1, "Hello", "l"))`, "3"},
		{`Print Cstr(InStr(1, "Hello", "L", vbTextCompare))`, "3"},
		{`Print Cstr(InStr(4, "Hello", "l"))`, "4"},
		{`Print Cstr(InStr(5, "Hello", "l"))`, "0"},
		{`Print LTrim("  Hello")`, "Hello"},
		{`Print RTrim("Hello  ")`, "Hello"},
		{`Print Space(5)`, "     "},
		{`Print Cstr(StrComp("Hello", "hello"))`, "-1"},
		{`Print Cstr(StrComp("Hello", "hello", vbTextCompare))`, "0"},
		{`Print Strng(5, "a")`, "aaaaa"},
		{`Print StrReverse("Hello")`, "olleH"},
		{`Print Trim("  Hello  ")`, "Hello"},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}
func TestDateFunctions(t *testing.T) {
	test := []struct {
		input    string
		expected string
	}{
		// {`Print CStr(Dte())`, "2024-02-22 00:00:00"},		/// ever changing
		// {`Print CStr(Now())`, "0000-11-30 13:48:50"},		/// ever changing
		{`Print CStr(DateAdd("d", 5, #2000/12/20#))`, "2000-12-25"},
		{`Print CStr(DateAdd("m", 1, "2000-11-20"))`, "2000-12-20"},
		{`Print CStr(DateDiff("d", #2019-01-01#, #2019-01-07#))`, "6"},
		{`Print CStr(DatePart("d", #2019-01-07#))`, "7"},
		{`Print CStr(DatePart("w", #2024-02-04#))`, "1"},
		{`Print CStr(DatePart("ww", #2024-02-04#))`, "6"},
		{`Print CStr(DateSerial(2019, 1, 7))`, "2019-01-07"},
		{`Print CStr(DateValue("2019-01-07"))`, "2019-01-07"},
		{`Print CStr(Day(#2019-01-07#))`, "7"},
		{`Print CStr(Hour(#2019-01-07 10:11:12#))`, "10"},
		{`Print CStr(Minute(#2019-01-07 10:11:12#))`, "11"},
		{`Print CStr(Month(#2019-01-07#))`, "1"},
		{`Print CStr(Second(#2019-01-07 10:11:12#))`, "12"},
		// {`Print CStr(Time())`, "00:00:00"},					/// ever changing
		// {`Print CStr(Timer())`, "0"},						/// ever changing
		{`Print CStr(TimeSerial(10, 11, 12))`, "10:11:12"},
		{`Print CStr(TimeValue("10:11:12"))`, "10:11:12"},
		{`Print CStr(Weekday(#2024-02-04#))`, "1"},
		{`Print CStr(Year(#2019-01-07#))`, "2019"},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}

func TestConversionFunctions(t *testing.T) {
	test := []struct {
		input    string
		expected string
	}{
		{`Print CStr(CBool("TRUE"))`, "true"},
		{`Print CStr(CBool(0))`, "false"},
		{`Print CStr(CDate("2019-01-07"))`, "2019-01-07"},
		{`Print CStr(CDate(""))`, ""},
		{`Print CStr(CDbl(5))`, "5.000000"},
		{`Print CStr(CLng(5.5))`, "5"},
		{`Print CStr(CStr(5))`, "5"},
		{`Print CStr(CVar(5))`, "5"},
		{`Print CStr(Asc("A"))`, "65"},
		// not compatible with the original function - must use go's fmt package
		// see: https://pkg.go.dev/fmt
		{`Print Format(5, "%03d")`, "005"}, // instead of "000" format
		{`Print CStr(Hex(255))`, "FF"},
		{`Print CStr(Oct(255))`, "377"},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}

func TestMathFunctions(t *testing.T) {
	test := []struct {
		input    string
		expected string
	}{
		{`Print CStr(Abs(-5))`, "5"},
		{`Print CStr(Atn(5))`, "1.373401"},
		{`Print CStr(Cos(5))`, "0.283662"},
		{`Print CStr(Expn(5))`, "148.413159"},
		{`Print CStr(Fix(5.5))`, "5.000000"},
		{`Print CStr(Int(5.5))`, "5.000000"},
		{`Print CStr(Log(5))`, "1.609438"},
		// {`Print CStr(Rnd())`, "0.000000"},			/// ever changing
		{`Print CStr(Sgn(-5))`, "-1"},
		{`Print CStr(Sin(5))`, "-0.958924"},
		{`Print CStr(Sqr(5))`, "2.236068"},
		{`Print CStr(Tan(5))`, "-3.380515"},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}

func TestLongArray(t *testing.T) {
	test := []struct {
		input    string
		expected int64
	}{
		{`Dim a() as long
		Redim a(10 + 5 -1)
		Print ubound(a)`, 13},
		{`dim a(5,3) as long
		Let a(4,2) = 5
		print a(ubound(a,1), ubound(a,2))`, 5},
		{`dim a() as long
		redim a(1)
		let a(0) = 19
		redim preserve a(10)
		print a(0)`, 19},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testLongObject(t, evaluated, tt.expected)
	}
}

func TestFibo(t *testing.T) {
	test := []struct {
		input    string
		expected int64
	}{
		{`' Fibonacci Sequence
		' uBASIC Math Project
		' ------------------------ 
		' The array Fibo holds the Fibonacci numbers
		Dim Fibo(11) As Long
		Let Fibo(0) = 0
		Let Fibo(1) = 1
		Print "Fibonacci Sequence"
		Print "-------------------"
		Print "0"
		
		Dim N As Integer
		For N = 1 To 10
			Let Fibo(N+1) = Fibo(N) + Fibo(N-1)
			Print Fibo(N)
		Next N
		Print Fibo(10)`, 55},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testLongObject(t, evaluated, tt.expected)
	}
}

func TestBooleanConstants(t *testing.T) {
	test := []struct {
		input    string
		expected string
	}{
		{`Print CStr(True)`, "true"},
		{`Print CStr(False)`, "false"},
		// {`Const bool as Boolean = True
		// let bool=False
		// Print bool`, "true"},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testStringObject(t, evaluated, tt.expected)
	}
}

func Test8Queens(t *testing.T) {
	test := []struct {
		input string
	}{
		{`' 8 Queens
		' uBASIC Math Project
		' ------------------------ 
		' The array Queens holds the column position of the Queens
		Dim Queens(9) As Long
		Dim W As Integer
		Let W= 1
		Let Queens(W) = 0
		
		Function IsSafe(W As Integer, ByRef Queens() As Double) As Boolean
			Dim i As Integer
			dim Result as Boolean
			Let i = 1
			Let result = True
			Do While (i < W) And result
				Let result = Queens(i) <> Queens(W) And (Abs(Queens(i) - Queens(W)) <> W - i)
				Let i = i + 1
			Loop
			Let IsSafe = result
		End Function

		Do While W > 0
			Let Queens(W) = Queens(W) + 1
			Do While Queens(W) <= 8 And Not IsSafe(W, Queens)
				Let Queens(W) = Queens(W) + 1
			Loop
			If Queens(W) <= 8 Then
				If W == 8 Then
					Print Queens
					exit do
				Else
					Let W = W + 1
					Let Queens(W) = 0
				End If
			Else
				Let W = W - 1
			End If
		Loop`},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testNothingObject(t, evaluated)
	}
}

func testNothingObject(t *testing.T, obj object.Object) bool {
	_, ok := obj.(*object.Nothing)
	if !ok {
		t.Errorf("object is not Nothing. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestBublleSort(t *testing.T) {
	test := []struct {
		input string
	}{
		{`' uBASIC Math Project
		' ------------------------ 
		' The array A holds the numbers to be sorted

		Sub BubbleSort(byref MyArray() As Variant)
		'Sorts a one-dimensional VBA array from smallest to largest
		'using the bubble sort algorithm.
		Dim i As Long, j As Long
		Dim Temp As Variant
		
		For i = LBound(MyArray, 1) To UBound(MyArray, 1) - 1
			For j = i + 1 To UBound(MyArray, 1)
				If MyArray(i) > MyArray(j) Then
					Let Temp = MyArray(j)
					Let MyArray(j) = MyArray(i)
					Let MyArray(i) = Temp
				End If
			Next j
		Next i
	  End Sub
	  
	  Sub Main()
		Dim MyArray(26) As Variant
		Dim i As Long
		
		'Fill the array with a permutation of the characters a-z
		For i = 0 To 25
			Let MyArray(i) = Chr(96 + rnd(i+1))
		Next i
		
		'Print the original array
		For i = 0 To 25
			Debug.Print MyArray(i);
		Next i
		Debug.Print
	  
		
		'Sort the array
		call BubbleSort(MyArray)
		
		'Print the sorted array
		For i = 0 To 25
			Debug.Print MyArray(i);
		Next i
	  End Sub
	  
	  call Main()
	  `},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testNil(t, evaluated)
	}
}

func testNil(t *testing.T, obj object.Object) bool {
	if obj != nil {
		t.Errorf("object is not nil. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}

func TestQuickSort(t *testing.T) { // not working
	test := []struct {
		input string
	}{
		{`' uBASIC Math Project
		' ------------------------ 
		' The array A holds the numbers to be sorted

		Sub QuickSort(ByRef A() As Variant, ByVal Low As Long, ByVal High As Long)
		' Sorts a one-dimensional VBA array from smallest to largest
		' using the quick sort algorithm.
		Dim i As Long, j As Long
		Dim Temp As Variant
		Dim Pivot As Variant
		
		Let i = Low
		Let j = High
		Let Pivot = A((Low + High) DIV 2)
		
		Do
			Do While A(i) < Pivot
				Let i = i + 1
			Loop
			Do While A(j) > Pivot
				Let j = j - 1
			Loop
			If i <= j Then
				Let Temp = A(i)
				Let A(i) = A(j)
				Let A(j) = Temp
				Let i = i + 1
				Let j = j - 1
			End If
		Loop While i <= j
		
		If Low < j Then
			call QuickSort(A, Low, j)
		End If
		If i < High Then
			call QuickSort(A, i, High)
		End If
	  End Sub
	  
	  Sub Main()
		Dim A(26) As Variant
		Dim i As Long
		
		'Fill the array with a permutation of the characters a-z
		For i = 0 To 25
			Let A(i) = Chr(96 + rnd(i+1))
		Next i
		
		'Print the original array
		For i = 0 To 25
			Debug.Print A(i);
		Next i
		Debug.Print
	  
		
		'Sort the array
		call QuickSort(A, LBound(A, 1), UBound(A, 1))
		
		'Print the sorted array
		For i = 0 To 25
			Debug.Print A(i);
		Next i
	  End Sub
	  
	  call Main()
	  `},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input) // not working
		testNil(t, evaluated)
	}
}

func TestPrimes(t *testing.T) {
	test := []struct {
		input string
	}{
		{`' Primes
		' uBASIC Math Project
		' ------------------------ 
		' The array Primes holds the prime numbers
		Dim Primes(100) As Long
		Dim N As Long
		Dim P As Long
		Dim IsPrime As Boolean
		Let Primes(1) = 2
		Let Primes(2) = 3
		Let N = 2
		Let P = 5
		
		Do While N < 100
			Let IsPrime = True
			Let P = P + 2
			Let N = N + 1
			Let Primes(N) = P
			dim i as long
			For i = 1 To N
				Let IsPrime = IsPrime And (P Mod Primes(i) <> 0)
			Next i
			If IsPrime Then
				Let N = N + 1
				Let Primes(N) = P
			End If
		Loop
		Print Primes
		`},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testNil(t, evaluated)
	}
}

func TestPrimes2(t *testing.T) {
	test := []struct {
		input string
	}{
		{`' Primes
		' uBASIC Math Project
		' ------------------------ 
Dim Primes() As Long

Function IsPrime(Number As Long) As Boolean
  Dim I As Long
  For I = LBound(Primes) To UBound(Primes)
      If (Number Mod Primes(I)== 0) Then 
        Let IsPrime = False
        Exit Function
      End If
      If (Primes(I) >= Sqr(Number)) Then 
        Exit For
      End If
  Next
  Let IsPrime = True
End Function


 Sub BuildPrimes(Max As Long)
  If (Max < 3) Then 
    Exit Sub
  End If

  Dim I As Long
  For I = 2 To Max
    If (IsPrime(I)) Then
      ReDim Preserve Primes(UBound(Primes) + 1)
      Let Primes(UBound(Primes)) = I
    End If
  Next
End Sub

Sub Initialize()
  ReDim Primes(1)
  Let Primes(0) = 2
End Sub
call Initialize()
call BuildPrimes(100)
dim i as long
for i = 0 to ubound(Primes)
  print Primes(i), "-";
next`},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testNothingObject(t, evaluated)
	}
}

func TestPrime2(t *testing.T) {
	test := []struct {
		input string
	}{
		{`' option base 0
		' option explicit
		Dim Primes() As Long

		Function IsPrime(Number As Long) As Boolean
		  Dim I As Long
		  For I = LBound(Primes) To UBound(Primes)
			  If (Number Mod Primes(I)== 0) Then 
				Let IsPrime = False
				Exit Function
			  End If
			  If (Primes(I) >= Sqr(Number)) Then 
				Exit For
			  End If
		  Next
		  Let IsPrime = True
		End Function
		
		
		 Sub BuildPrimes(Max As Long)
		  If (Max < 3) Then 
			Exit Sub
		  End If
		
		  Dim I As Long
		  For I = 2 To Max
			If (IsPrime(I)) Then
			  ReDim Preserve Primes(UBound(Primes) + 2)	'option base 0 -> ubound == 0
			  Let Primes(UBound(Primes)) = I
			End If
		  Next
		End Sub
		
		Sub Initialize()
		  ReDim Primes(1)
		  Let Primes(0) = 2
		End Sub
		call Initialize()
		call BuildPrimes(100)
		
		dim i as Long
		print lbound(Primes),  ubound(Primes)
		for i = LBound(Primes) to UBound(Primes)
		  Debug.Print Primes(i)
		next i
		
		`},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		testNil(t, evaluated)
	}
}

func TestExits(t *testing.T) {
	test := []struct {
		input    string
		expected int64
	}{
		{`Dim a As Long
		For a = 1 To 5
			Print a
			If a == 3 Then
				Exit For
			end if
		Next
		Print a`, 3},
		{`Dim a As Long
		For a = 1 To 5
			Print a
			If a == 3 Then
				Exit For
			Else
				Print "hello"
			End If
		Next
		Print a`, 3},
		{`Dim a As Long
		Let a = 1
		do While a < 5
			Print a
			If a == 3 Then
				Exit do
			End If
			Let a = a + 1
		Loop
		Print a`, 3},
		{`Dim a As Long
		Let a = 1
		do until a == 5
			Print a
			If a == 3 Then
				Exit do
			Else
				Print "hello"
			End If
			Let a = a + 1
		loop
		Print a`, 3},
		{`Dim a As Long
		' text exit sub
		Let a = 1
		sub test()
			Let a = 5
			Exit Sub
			let a = 10
			Print "hello"
		end sub
		call test()
		Print a`, 5},
		{`' test exit function
		function test() as long
			let test = 5
			Exit function
			let test = 10
		end function
		print test()`, 5},
	}
	for _, tt := range test {
		evaluated := testEval(tt.input)
		fmt.Println()
		testLongObject(t, evaluated, tt.expected)
	}
}
