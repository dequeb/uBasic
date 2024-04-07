package parser

import (
	"testing"
	"uBasic/ast"
	"uBasic/lexer"
	"uBasic/token"
)

func TestDimStatement(t *testing.T) {
	input := `Dim x as Integer: Dim w as integer
	Dim y as Single, v() as UneEnum
	Dim z as Double, u(1,2) as Double, t(1,2,3) as Double
	Dim a as String, s(4) as string
	Dim b as Boolean, r as string
	Dim c as Date
	Dim d as Variant`

	l := lexer.New(input)
	p := New(l)

	file := p.ParseFile()
	if file == nil {
		for _, e := range p.Errors() {
			t.Log(e.Error())
		}
		t.Fatalf("ParseFile() returned nil")
	}
	if len(file.Body) != 7 {
		t.Fatalf("file.Statements does not contain 7 statements. got=%d", len(file.Body))
	}

	tests := []struct {
		expectedIdentifier string
		expectedType       string
	}{
		{"x", "Integer"},
		{"y", "Single"},
		{"z", "Double"},
		{"a", "String"},
		{"b", "Boolean"},
		{"c", "Date"},
		{"d", "Variant"},
	}

	for i, tt := range tests {
		stmt := file.Body[i].Statements[0] // only validate first statement of each line

		if !testDimStatement(t, stmt, tt.expectedIdentifier, tt.expectedType) {
			return
		}
	}
	t.Log(file.String())
}

func testDimStatement(t *testing.T, s ast.Node, name string, typ string) bool {
	if s.Token().Literal != "Dim" {
		t.Errorf("s.TokenLiteral not 'Dim'. got=%q", s.Token().Literal)
		return false
	}

	dimStmt, ok := s.(*ast.DimDecl)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if dimStmt.Vars[0].Name().Name != name {
		t.Errorf("dimStmt.Vars[0].String()  not '%s'. got=%s", name, dimStmt.Vars[0].String())
		return false
	}

	dimType, err := dimStmt.Vars[0].Type()
	if err != nil {
		t.Errorf("dimStmt.Vars[0].Type()  not '%s'. got=%s", typ, err)
		return false
	}
	if dimType.String() != typ {
		t.Errorf("letStmt.Type.Value not '%s'. got=%s", typ, dimType.String())
		return false
	}

	return true
}

func TestConstStatement(t *testing.T) {
	input := `Const a as Integer = 2: const b as Single = 3.14
	Const c as Double = 3.14:const  d as String = "hello"
	Const e as Boolean = True: const f as Date = #2019/12/25#: const  g as Variant = 3.14`

	l := lexer.New(input)
	p := New(l)

	file := p.ParseFile()
	if file == nil {
		t.Fatalf("ParseFile() returned nil")
	}
	if len(file.Body) != 3 {
		t.Fatalf("file.Statements does not contain 3 statements. got=%d", len(file.Body))
	}

	tests := []struct {
		expectedIdentifier string
		expectedType       string
		expectedValue      string
	}{
		{"a", "Integer", "2"},
		{"b", "Single", "3.14"},
		{"c", "Double", "3.14"},
		{"d", "String", "\"hello\""},
		{"e", "Boolean", "True"},
		{"f", "Date", "#2019/12/25#"},
		{"g", "Variant", "3.14"},
	}

	var test int
	for _, stmtList := range file.Body {
		for _, stmt := range stmtList.Statements {
			if !testConstStatement(t, stmt, tests[test].expectedIdentifier, tests[test].expectedType, tests[test].expectedValue) {
				return
			}
			test++
		}
	}
	t.Log(file.String())
}

func testConstStatement(t *testing.T, s ast.Node, name string, typ string, value string) bool {
	if s.Token().Kind != token.KwConst {
		t.Errorf("s.TokenLiteral not 'Const'. got=%q", s.Token().Literal)
		return false
	}

	Stmt, ok := s.(*ast.ConstDecl)
	if !ok {
		t.Errorf("s not *ast.ConstDecl. got=%T", s)
		return false
	}

	if Stmt.Consts[0].Name().Name != name {
		t.Errorf("dimStmt.Vars[0].String()  not '%s'. got=%s", name, Stmt.Consts[0].String())
		return false
	}

	constType, err := Stmt.Consts[0].Type()
	if err != nil {
		t.Errorf("Stmt.Consts[0].Type()  not '%s'. got=%s", typ, err)
		return false
	}
	if constType.String() != typ {
		t.Errorf("letStmt.Type.Value not '%s'. got=%s", typ, constType.String())
		return false
	}
	testValue := Stmt.Consts[0].Value().String()
	if testValue != value {
		t.Errorf("Stmt.Consts[0].Value.String() not '%s'. got=%s", value, testValue)
		return false
	}

	return true
}

func TestEnumStatement(t *testing.T) {
	input := `Enum numbers
		one 
		two
		three
	End Enum`

	l := lexer.New(input)
	p := New(l)

	file := p.ParseFile()
	if file == nil {
		for _, e := range p.Errors() {
			t.Log(e.Error())
		}
		t.Fatalf("ParseFile() returned nil")
	}
	if len(file.Body) != 1 {
		t.Fatalf("file.Statements does not contain 1 statements. got=%d", len(file.Body))
	}

	tests := []struct {
		expectedIdentifier string
		expectedValues     []string
	}{
		{"numbers", []string{"one", "two", "three"}},
	}

	var test int
	for _, stmtList := range file.Body {
		for _, stmt := range stmtList.Statements {
			if !testEnumStatement(t, stmt, tests[test].expectedIdentifier, tests[test].expectedValues) {
				return
			}
			test++
		}
	}
	t.Log(file.String())
}

func testEnumStatement(t *testing.T, s ast.Node, name string, values []string) bool {
	if s.Token().Kind != token.KwEnum {
		t.Errorf("s.TokenLiteral not 'Const'. got=%q", s.Token().Literal)
		return false
	}

	Stmt, ok := s.(*ast.EnumDecl)
	if !ok {
		t.Errorf("s not *ast.EnumDecl. got=%T", s)
		return false
	}

	if Stmt.Name().Name != name {
		t.Errorf("Enum Stmt.Name().Name   not '%s'. got=%s", name, Stmt.Name().Name)
		return false
	}

	if len(Stmt.Values) != len(values) {
		t.Errorf("Enum.Values length  not '%s'. got=%s", values, Stmt.Values)
		return false
	}

	for i, v := range Stmt.Values {
		if v.Name != values[i] {
			t.Errorf("dimStmt.Vars[0].String()  not '%s'. got=%s", values[i], v.String())
			return false
		}
	}

	return true
}

func TestFunctionStatement(t *testing.T) {
	input := `function test() as Integer
		Dim x as Integer: Dim w as integer
		' test = 1
		end function
	function test2(byval a as integer, byRef b as long, paramArray c() as currency) as String
		Dim y as Single, v() as UneEnum
		' test2 = "hello"
	end function
	function test3(byref optional c as double = 2.2) as Double
		const a as Integer = 2
		' test3 = 3.14
		exit function
	end function`

	l := lexer.New(input)
	p := New(l)

	file := p.ParseFile()
	if file == nil {
		for _, e := range p.Errors() {
			t.Log(e.Error())
		}
		t.Fatalf("ParseFile() returned nil")
	}
	if len(file.Body) != 3 {
		t.Fatalf("file.Statements does not contain 3 statements. got=%d", len(file.Body))
	}

	tests := []struct {
		expectedIdentifier string
		expectedParams     []string
		expectedType       string
	}{
		{"test", []string{}, "Integer"},
		{"test2", []string{"a", "b", "c"}, "String"},
		{"test3", []string{"c"}, "Double"},
	}

	var test int
	for _, stmtList := range file.Body {
		for _, stmt := range stmtList.Statements {
			// only validate first statement of each line
			if !testFunctionStatement(t, stmt, tests[test].expectedIdentifier, tests[test].expectedParams, tests[test].expectedType) {
				return
			}
			test++
		}
	}
	t.Log(file.String())
}

func testFunctionStatement(t *testing.T, s ast.Node, name string, params []string, typ string) bool {
	if s.Token().Kind != token.KwFunction {
		t.Errorf("s.TokenLiteral not 'Function'. got=%q", s.Token().Literal)
		return false
	}

	Stmt, ok := s.(*ast.FuncDecl)
	if !ok {
		t.Errorf("s not *ast.FunctionDecl. got=%T", s)
		return false
	}

	if Stmt.FuncName.Name != name {
		t.Errorf("Function Stmt.FuncName.Name   not '%s'. got=%s", name, Stmt.FuncName.Name)
		return false
	}

	funcType := Stmt.FuncType

	if len(funcType.Params) != len(params) {
		t.Errorf("Function.Params length  not '%d'. got=%d", len(params), len(funcType.Params))
		return false
	}

	for i, p := range funcType.Params {
		if p.Name().Name != params[i] {
			t.Errorf("Function.Params[%d].Name().Name  not '%s'. got=%s", i, params[i], p.Name().Name)
			return false
		}
	}
	if funcType.Result.String() != typ {
		t.Errorf("Stmt.Type().Value not '%s'. got=%s", typ, funcType.Result.String())
		return false
	}

	return true
}

func TestSubStatement(t *testing.T) {
	input := `Sub test() 
	Dim x as Integer: Dim w as integer
	Let test = 1
	end Sub
Sub test2(byval a as integer, byRef b as long, paramArray c() as currency) 
	Dim y as Single, v() as UneEnum
	Let test2 = "hello"
end Sub
Sub test3(byref optional c as double = 2.2) 
	const a as Integer = 2
	Let test3 = 3.14
	exit Sub
end sub`

	l := lexer.New(input)
	p := New(l)

	file := p.ParseFile()
	if file == nil {
		for _, e := range p.Errors() {
			t.Log(e.Error())
		}
		t.Fatalf("ParseFile() returned nil")
	}
	if len(file.Body) != 3 {
		t.Fatalf("file.Statements does not contain 3 statements. got=%d", len(file.Body))
	}

	tests := []struct {
		expectedIdentifier string
		expectedParams     []string
	}{
		{"test", []string{}},
		{"test2", []string{"a", "b", "c"}},
		{"test3", []string{"c"}},
	}

	var test int
	for _, stmtList := range file.Body {
		for _, stmt := range stmtList.Statements {
			// only validate first statement of each line
			if !testSubStatement(t, stmt, tests[test].expectedIdentifier, tests[test].expectedParams) {
				return
			}
			test++
		}
	}
	t.Log(file.String())
}

func testSubStatement(t *testing.T, s ast.Node, name string, params []string) bool {
	if s.Token().Kind != token.KwSub {
		t.Errorf("s.TokenLiteral not 'Sub'. got=%q", s.Token().Literal)
		return false
	}

	Stmt, ok := s.(*ast.SubDecl)
	if !ok {
		t.Errorf("s not *ast.SubDecl. got=%T", s)
		return false
	}

	if Stmt.SubName.Name != name {
		t.Errorf("Function Stmt.SubName.Name   not '%s'. got=%s", name, Stmt.SubName.Name)
		return false
	}

	subType := Stmt.SubType

	if len(subType.Params) != len(params) {
		t.Errorf("Sub.Params length  not '%d'. got=%d", len(params), len(subType.Params))
		return false
	}

	for i, p := range subType.Params {
		if p.Name().Name != params[i] {
			t.Errorf("Function.Params[%d].Name().Name  not '%s'. got=%s", i, params[i], p.Name().Name)
			return false
		}
	}

	return true
}

func TestIfStatement(t *testing.T) {
	input := `If true Then
		Dim x As Integer
	ElseIf false Then
		Dim y As Single
	Else
		Dim z As Double
	End If`

	l := lexer.New(input)
	p := New(l)

	file := p.ParseFile()
	if file == nil {
		for _, e := range p.Errors() {
			t.Log(e.Error())
		}
		t.Fatalf("ParseFile() returned nil")
	}
	if len(file.Body) != 1 {
		t.Fatalf("file.Statements does not contain 1 statements. got=%d", len(file.Body))
	}

	tests := []struct {
		expectedExpression      string
		expectedStatements      []string
		expectedElsifExpression string
		expectedElsifStatements []string
		expectedElseStatements  []string
	}{
		{"true", []string{"Dim x As Integer\n"}, "false", []string{"Dim y As Single\n"}, []string{"Dim z As Double"}},
	}

	var test int
	for _, stmtList := range file.Body {
		for _, stmt := range stmtList.Statements {
			// only validate first statement of each line
			if !testIfStatement(t, stmt, tests[test].expectedExpression, tests[test].expectedStatements, tests[test].expectedElsifExpression, tests[test].expectedElsifStatements, tests[test].expectedElseStatements) {
				return
			}
			test++
		}
	}
	t.Log(file.String())
}

func testIfStatement(t *testing.T, s ast.Node, expression string, statements []string, elsifExpression string, elsifStatements []string, elseStatements []string) bool {
	if s.Token().Kind != token.KwIf {
		t.Errorf("s.TokenLiteral not 'If'. got=%q", s.Token().Literal)
		return false
	}

	Stmt, ok := s.(*ast.IfStmt)
	if !ok {
		t.Errorf("s not *ast.IfStmt. got=%T", s)
		return false
	}

	if Stmt.Condition.String() != expression {
		t.Errorf("If Stmt.Condition.String()   not '%s'. got=%s", expression, Stmt.Condition.String())
		return false
	}

	if len(Stmt.Body) != len(statements) {
		t.Errorf("If.Body length  not '%d'. got=%d", len(statements), len(Stmt.Body))
		return false
	}

	for i, p := range Stmt.Body {
		if p.String() != statements[i] {
			t.Errorf("If.Statements[%d]  not '%s'. got=%s", i, statements[i], p.String())
			return false
		}
	}

	if len(Stmt.ElseIf) > 0 {
		if Stmt.ElseIf[0].Condition.String() != elsifExpression {
			t.Errorf("If.Elsif[0].Condition.String()   not '%s'. got=%s", elsifExpression, Stmt.ElseIf[0].Condition.String())
			return false
		}

		if len(Stmt.ElseIf[0].Body) != len(elsifStatements) {
			t.Errorf("If.Elsif[0].Statements length  not '%d'. got=%d", len(elsifStatements), len(Stmt.ElseIf[0].Body))
			return false
		}

		for i, p := range Stmt.ElseIf[0].Body {
			if p.String() != elsifStatements[i] {
				t.Errorf("If.Elsif[0].Statements[%d]  not '%s'. got=%s", i, elsifStatements[i], p.String())
				return false
			}
		}
	}

	if len(Stmt.Else) > 0 {
		if len(Stmt.Else[0].Statements) != len(elseStatements) {
			t.Errorf("If.Else[0].Statements length  not '%d'. got=%d", len(elseStatements), len(Stmt.Else[0].Statements))
			return false
		}

		for i, p := range Stmt.Else[0].Statements {
			if p.String() != elseStatements[i] {
				t.Errorf("If.Else[0].Statements[%d]  not '%s'. got=%s", i, elseStatements[i], p.String())
				return false
			}
		}
	}
	return true
}

func TestForStatement(t *testing.T) {
	input := `For i = 1 To 10 step 2
		Dim x As Integer
		Next
		for each j in arr
			Dim y As Single
		Next j`

	l := lexer.New(input)
	p := New(l)

	file := p.ParseFile()
	if file == nil {
		for _, e := range p.Errors() {
			t.Log(e.Error())
		}
		t.Fatalf("ParseFile() returned nil")
	}
	if len(file.Body) != 2 {
		t.Fatalf("file.Statements does not contain 2 statements. got=%d", len(file.Body))
	}

	tests := []struct {
		expectedExpression string
		expectedStatements []string
	}{
		{"i = 1 To 10 Step 2", []string{"Dim x As Integer\n"}},
		{"Each j In arr", []string{"Dim y As Single\n"}},
	}

	var test int
	for _, stmtList := range file.Body {
		for _, stmt := range stmtList.Statements {
			// only validate first statement of each line
			if !testForStatement(t, stmt, tests[test].expectedExpression, tests[test].expectedStatements) {
				return
			}
			test++
		}
	}
	t.Log(file.String())
}

func testForStatement(t *testing.T, s ast.Node, expression string, statements []string) bool {
	if s.Token().Kind != token.KwFor {
		t.Errorf("s.TokenLiteral not 'For'. got=%q", s.Token().Literal)
		return false
	}

	Stmt, ok := s.(*ast.ForStmt)
	if !ok {
		t.Errorf("s not *ast.ForStmt. got=%T", s)
		return false
	}

	if Stmt.ForExpression.String() != expression {
		t.Errorf("For Stmt.Condition.String()   not '%s'. got=%s", expression, Stmt.ForExpression.String())
		return false
	}

	if len(Stmt.Body) != len(statements) {
		t.Errorf("For.Body length  not '%d'. got=%d", len(statements), len(Stmt.Body))
		return false
	}

	for i, p := range Stmt.Body {
		if p.String() != statements[i] {
			t.Errorf("For.Statements[%d]  not '%s'. got=%s", i, statements[i], p.String())
			return false
		}
	}

	return true
}

func TestSpecialStatement(t *testing.T) {
	input := `Print 2
	redim x(1)
	redim preserve x(2)
	erase x
	Debug.Print "hello", """", "World";
	goto begin
	Stop`

	l := lexer.New(input)
	p := New(l)

	file := p.ParseFile()
	if file == nil {
		for _, e := range p.Errors() {
			t.Log(e.Error())
		}
		t.Fatalf("ParseFile() returned nil")
	}
	if len(file.Body) != 7 {
		t.Fatalf("file.Statements does not contain 7 statements. got=%d", len(file.Body))
	}

	tests := []struct {
		expectedKeyword1 string
		expectedKeyword2 string
		expectedArgs     []string
	}{
		{"Print", "", []string{"2"}},
		{"redim", "", []string{"x(1)"}},
		{"redim", "preserve", []string{"x(2)"}},
		{"erase", "", []string{"x"}},
		{"Debug.Print", "", []string{"\"hello\"", "\"\"\"\"", "\"World\""}},
		{"goto", "", []string{"begin"}},
		{"Stop", "", []string{}},
	}

	for i, tt := range tests {
		stmt := file.Body[i].Statements[0] // only validate first statement of each line

		if !testSpecialStatement(t, stmt, tt.expectedKeyword1, tt.expectedKeyword2, tt.expectedArgs) {
			return
		}
	}
	t.Log(file.String())
}

func testSpecialStatement(t *testing.T, s ast.Node, keyword1 string, keyword2 string, args []string) bool {
	if s.Token().Literal != keyword1 {
		t.Errorf("s.TokenLiteral not '%s'. got=%q", keyword1, s.Token().Literal)
		return false
	}

	Stmt, ok := s.(*ast.SpecialStmt)
	if !ok {
		t.Errorf("s not *ast.SpecialStmt. got=%T", s)
		return false
	}

	if Stmt.Keyword1.Literal != keyword1 {
		t.Errorf("SpecialStmt.Keyword1   not '%s'. got=%s", keyword1, Stmt.Keyword1)
		return false
	}

	if Stmt.Keyword2 != keyword2 {
		t.Errorf("SpecialStmt.Keyword2   not '%s'. got=%s", keyword2, Stmt.Keyword2)
		return false
	}

	if len(Stmt.Args) != len(args) {
		t.Errorf("SpecialStmt.Args length  not '%d'. got=%d", len(args), len(Stmt.Args))
		return false
	}

	for i, p := range Stmt.Args {
		if p.String() != args[i] {
			t.Errorf("SpecialStmt.Args[%d]  not '%s'. got=%s", i, args[i], p.String())
			return false
		}
	}

	return true
}

func TestSelectStatement(t *testing.T) {
	input := `Select Case x
	Case 1
		Print "one"
	Case 2
		Print "two"
	Case 3
		Print "three"
	Case Else
		Print "other"
	End Select`

	l := lexer.New(input)
	p := New(l)

	file := p.ParseFile()
	if file == nil {
		for _, e := range p.Errors() {
			t.Log(e.Error())
		}
		t.Fatalf("ParseFile() returned nil")
	}
	if len(file.Body) != 1 {
		t.Fatalf("file.Statements does not contain 1 statements. got=%d", len(file.Body))
	}

	tests := []struct {
		expectedCondition string
		expectedCases     []string
	}{
		{"x", []string{"1", "2", "3", ""}},
	}

	for i, tt := range tests {
		stmt := file.Body[i].Statements[0] // only validate first statement of each line

		if !testSelectStatement(t, stmt, tt.expectedCondition, tt.expectedCases) {
			return
		}
	}
	t.Log(file.String())
}

func testSelectStatement(t *testing.T, s ast.Node, condition string, cases []string) bool {
	if s.Token().Kind != token.KwSelect {
		t.Errorf("s.TokenLiteral not 'Select'. got=%q", s.Token().Literal)
		return false
	}

	Stmt, ok := s.(*ast.SelectStmt)
	if !ok {
		t.Errorf("s not *ast.SelectStmt. got=%T", s)
		return false
	}

	if Stmt.Condition.String() != condition {
		t.Errorf("SelectStmt.Condition.String()   not '%s'. got=%s", condition, Stmt.Condition.String())
		return false
	}

	if len(Stmt.Body) != len(cases) {
		t.Errorf("SelectStmt.Cases length  not '%d'. got=%d", len(cases), len(Stmt.Body))
		return false
	}

	for i, p := range Stmt.Body {
		if p.Condition != nil {
			if p.Condition.String() != cases[i] {
				t.Errorf("SelectStmt.Cases[%d]  not '%s'. got=%s", i, cases[i], p.Condition.String())
				return false
			}
		} else {
			if cases[i] != "" {
				t.Errorf("SelectStmt.Cases[%d]  not '%s'. got=%s", i, cases[i], "")
				return false
			}
		}
	}
	return true
}

func TestDoStatement(t *testing.T) {
	input := `Do
		Dim x As Integer
	Loop While x 
	Do While x 
		Print x
	Loop
	Do Until x 
		Print x
	Loop
	Do
		Print x
	Loop Until x 
	`

	l := lexer.New(input)
	p := New(l)

	file := p.ParseFile()
	if file == nil {
		for _, e := range p.Errors() {
			t.Log(e.Error())
		}
		t.Fatalf("ParseFile() returned nil")
	}
	if len(file.Body) != 4 {
		t.Fatalf("file.Statements does not contain 4 statements. got=%d", len(file.Body))
	}

	tests := []struct {
		expectedType      string
		expectedCondition string
	}{
		{"LoopWhile", "x"},
		{"DoWhile", "x"},
		{"DoUntil", "x"},
		{"LoopUntil", "x"},
	}

	for i, tt := range tests {
		stmt := file.Body[i].Statements[0] // only validate first statement of each line

		if !testDoStatement(t, stmt, tt.expectedType, tt.expectedCondition) {
			return
		}
	}
	t.Log(file.String())
}

func testDoStatement(t *testing.T, s ast.Node, doType string, condition string) bool {
	if s.Token().Kind != token.KwDo {
		t.Errorf("s.TokenLiteral not 'Do'. got=%q", s.Token().Literal)
		return false
	}
	var conditionStr string

	switch s := s.(type) {
	case *ast.WhileStmt:
		if doType != "DoWhile" {
			t.Errorf("s not *ast.WhileStmt. got=%T", s)
			return false
		}
		conditionStr = s.Condition.String()
	case *ast.UntilStmt:
		if doType != "DoUntil" {
			t.Errorf("s not *ast.UntilStmt. got=%T", s)
			return false
		}
		conditionStr = s.Condition.String()
	case *ast.DoWhileStmt:
		if doType != "LoopWhile" {
			t.Errorf("s not *ast.DoWhileStmt. got=%T", s)
			return false
		}
		conditionStr = s.Condition.String()
	case *ast.DoUntilStmt:
		if doType != "LoopUntil" {
			t.Errorf("s not *ast.DoUntilStmt. got=%T", s)
			return false
		}
		conditionStr = s.Condition.String()
	default:
		t.Errorf("s not *ast.DoStmt. got=%T", s)
		return false
	}

	if conditionStr != condition {
		t.Errorf("DoStmt.Condition.String()   not '%s'. got=%s", condition, conditionStr)
		return false
	}

	return true
}

func TestSubCallStatement(t *testing.T) {
	input := `Call HouseCalc(99800, 43100 )
call MsgBox("This house is affordable.")`

	l := lexer.New(input)
	p := New(l)

	file := p.ParseFile()
	if file == nil {
		for _, e := range p.Errors() {
			t.Log(e.Error())
		}
		t.Fatalf("ParseFile() returned nil")
	}
	if len(file.Body) != 2 {
		t.Fatalf("file.Statements does not contain 2 statements. got=%d", len(file.Body))
	}

	tests := []struct {
		ExpectedSubName string
		expectedParams  []string
	}{
		{"HouseCalc", []string{"99800", "43100"}},
		{"MsgBox", []string{"\"This house is affordable.\""}},
	}

	for i, tt := range tests {
		stmt := file.Body[i].Statements[0] // only validate first statement of each line

		if !testSubCallStatement(t, stmt, tt.ExpectedSubName, tt.expectedParams) {
			return
		}
	}
	t.Log(file.String())
}

func testSubCallStatement(t *testing.T, s ast.Node, subName string, params []string) bool {

	if s.Token().Kind != token.KwCall {
		t.Errorf("s.TokenLiteral not 'Call'. got=%q", s.Token().Literal)
		return false
	}

	Stmt, ok := s.(*ast.CallSubStmt)
	if !ok {
		t.Errorf("s not *ast.SubCall. got=%T", s)
		return false
	}

	definition := Stmt.Definition.(*ast.CallOrIndexExpr)
	if definition.Identifier.Name != subName {
		t.Errorf("SubCall.Name   not '%s'. got=%s", subName, definition.Identifier.Name)
		return false
	}

	if len(definition.Args) != len(params) {
		t.Errorf("SubCall.Args length  not '%d'. got=%d", len(params), len(definition.Args))
		return false
	}

	for i, p := range definition.Args {
		if p.String() != params[i] {
			t.Errorf("SubCall.Args[%d]  not '%s'. got=%s", i, params[i], p.String())
			return false
		}
	}

	return true
}

func TestExpressionStatement(t *testing.T) {
	input := `Let x = 1
	Let y = 2 * 3 + 4
	Let z(x+4) = 5 == 6
	Let a = b <> c
	Let d(1,2) = (Sin(3) + Cos(4)) * 3.1415926
	Let e = f(1,2,3)
	Let root.value = -array(44)
	Let A = 77 & "hello" & "world" & "!"
	`

	l := lexer.New(input)
	p := New(l)

	file := p.ParseFile()
	if file == nil {
		for _, e := range p.Errors() {
			t.Log(e.Error())
		}
		t.Fatalf("ParseFile() returned nil")
	}
	if len(file.Body) != 8 {
		t.Fatalf("file.Statements does not contain 8 statements. got=%d", len(file.Body))
	}

	tests := []struct {
		ExpectedLeft  string
		expectedRight string
	}{
		{"x", "1"},
		{"y", "2 * 3 + 4"},
		{"z(x + 4)", "5 == 6"},
		{"a", "b <> c"},
		{"d(1, 2)", "(Sin(3) + Cos(4)) * 3.1415926"},
		{"e", "f(1, 2, 3)"},
		{"root.value", "-array(44)"},
		{"A", "77 & \"hello\" & \"world\" & \"!\""},
	}

	for i, tt := range tests {
		stmt := file.Body[i].Statements[0] // only validate first statement of each line

		if !testExpressionStatement(t, stmt, tt.ExpectedLeft, tt.expectedRight) {
			return
		}
	}
	t.Log(file.String())
}

func testExpressionStatement(t *testing.T, s ast.Node, left string, right string) bool {

	Stmt, ok := s.(*ast.ExprStmt)
	if !ok {
		t.Errorf("s not *ast.ExprStmt. got=%T", s)
		return false
	}

	expr, ok := Stmt.Expression.(*ast.BinaryExpr)
	if !ok {
		t.Errorf("s not *ast.BinaryExpr. got=%T", s)
		return false
	}

	if expr.Left.String() != left {
		t.Errorf("ExprStmt.Left   not '%s'. got=%s", left, expr.Left.String())
		return false
	}

	if expr.Right.String() != right {
		t.Errorf("ExprStmt.Right   not '%s'. got=%s", right, expr.Right.String())
		return false
	}

	if expr.OpKind != token.Assign {
		t.Errorf("ExprStmt.Op   not '%s'. got=%s", "=", expr.OpToken.Literal)
		return false
	}

	return true
}

func TestErrorHandlingStatement(t *testing.T) {
	input := `On Error Resume Next
	On Error _
		GoTo 0
	On Error GoTo ErrorHandler
	Resume Next
	Resume ErrorHandler
	`

	l := lexer.New(input)
	p := New(l)

	file := p.ParseFile()
	if file == nil {
		for _, e := range p.Errors() {
			t.Log(e.Error())
		}
		t.Fatalf("ParseFile() returned nil")
	}
	if len(file.Body) != 5 {
		t.Fatalf("file.Statements does not contain 5 statements. got=%d", len(file.Body))
	}

	tests := []struct {
		exprectedOnError bool
		ExpectedKeyword  string
		ExpectedLabel    string
	}{
		{true, "Resume", "Next"},
		{true, "GoTo", "0"},
		{true, "GoTo", "ErrorHandler"},
		{false, "Resume", "Next"},
		{false, "Resume", "ErrorHandler"},
	}

	for i, tt := range tests {
		stmt := file.Body[i].Statements[0] // only validate first statement of each line

		if !testErrorHandlingStatement(t, stmt, tt.exprectedOnError, tt.ExpectedKeyword, tt.ExpectedLabel) {
			return
		}
	}
}

func testErrorHandlingStatement(t *testing.T, s ast.Node, onError bool, keyword string, label string) bool {

	Stmt, ok := s.(*ast.JumpStmt)
	if !ok {
		t.Errorf("s not *ast.ErrorHandlingStmt. got=%T", s)
		return false
	}

	if Stmt.OnError != onError {
		t.Errorf("ErrorHandlingStmt.OnError   not '%t'. got=%t", onError, Stmt.OnError)
		return false
	}

	if Stmt.JumpKw.Literal != keyword {
		t.Errorf("ErrorHandlingStmt.Keyword   not '%s'. got=%s", keyword, Stmt.JumpKw.Literal)
		return false
	}

	var labelStr string
	if Stmt.Label != nil {
		labelStr = Stmt.Label.Name
	} else if Stmt.NextKw != nil {
		labelStr = Stmt.NextKw.Literal
	} else if Stmt.Number != nil {
		labelStr = Stmt.Number.Literal
	}
	if labelStr != label {
		t.Errorf("ErrorHandlingStmt.Label   not '%s'. got=%s ", label, labelStr)
		return false
	}

	return true
}

// test labels
func TestLabelStatement(t *testing.T) {
	input := `label:
	`

	l := lexer.New(input)
	p := New(l)

	file := p.ParseFile()
	if file == nil {
		for _, e := range p.Errors() {
			t.Log(e.Error())
		}
		t.Fatalf("ParseFile() returned nil")
	}
	if len(file.Body) != 1 {
		t.Fatalf("file.Statements does not contain 1 statement. got=%d", len(file.Body))
	}

	tests := []struct {
		expectedLabel string
	}{
		{"label"},
	}

	for i, tt := range tests {
		stmt := file.Body[i].Statements[0] // only validate first statement of each line

		if !testLabelStatement(t, stmt, tt.expectedLabel) {
			return
		}
	}
}

func testLabelStatement(t *testing.T, s ast.Node, label string) bool {

	Stmt, ok := s.(*ast.JumpLabelDecl)
	if !ok {
		t.Errorf("s not *ast.LabelStmt. got=%T", s)
		return false
	}

	if Stmt.Label.Name != label {
		t.Errorf("LabelStmt.Label   not '%s'. got=%s", label, Stmt.Label.Name)
		return false
	}

	return true
}
