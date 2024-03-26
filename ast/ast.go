package ast

import (
	"fmt"
	"strings"
	"uBasic/token"
	"uBasic/types"

	"github.com/iancoleman/strcase"
)

// A File represents a µBASIC source file.
type File struct {
	// Top-level declarations.
	Body []StatementList
	// File name.
	Name string
	// parent node
	Parent Node // always nil
}

// dictionary of symbols
type Dict map[string]int

type IdList map[int]string

// A Node represents a node within the abstract syntax tree, and has one of the
// following underlying types.
//
//	*File
//	Decl
//	Stmt
//	Expr
type Node interface {
	fmt.Stringer
	// Token returns the first token of the Node
	Token() *token.Token
	// parent node
	GetParent() Node
	// set parent node
	SetParent(Node)
	// comparaison of two nodes. Returns nil if nodes are equal
	EqualsNode(Node) Node
}

// A Decl node represents a declaration, and has one of the following underlying
// types.
//
//			*FuncDecl
//			*SubDecl
//			*VarDecl
//			*EnumDecl
//			*ConstDeclItem
//			*ArrayDecl
//		 *DimDecl
//		 *ScalarDecl
//	  *LabelDecl
//
// Pseudo-code representation of a declaration.
//
//	type ident [= value]
type Decl interface {
	Node
	// Type returns the type of the declared identifier.
	Type() (types.Type, error)
	// Name returns the name of the declared identifier.
	Name() *Identifier
	// Value returns the initializing value of the defined identifier; or nil if
	// declaration or tentative definition.
	//
	// Underlying type for type definitions.
	//    Type
	Value() Node
	// isDecl ensures that only declaration nodes can be assigned to the Decl
	// interface.
	isDecl()
}

// a VarDecl node represents a variable declaration.
//
// Examples.
//
//	Dim x As Integer
//	Dim x As Integer = 42
//	Dim x(10) As Integer
//	Dim x(10, 20) As Integer
//
// It represents the following types
//
//	*ScalarDecl
//	*ArrayDecl
type (

	// A TypeDef node represents a type definition.
	//
	// Examples.
	//
	//    typedef int foo
	TypeDef struct {
		// Position of `typedef` keyword.
		Typedef *token.Token
		// Underlying type of type definition.
		DeclType Type
		// Type name.
		TypeName *Identifier
		// Underlying type of type definition.
		Val types.Type
		// parent node
		Parent Node
	}

	VarDecl interface {
		Node
		// Type returns the type of the declared identifier.
		Type() (types.Type, error)
		// Name returns the name of the declared identifier.
		Name() *Identifier
		// isVarDecl ensures that only declaration nodes can be assigned to the VarDecl
		// interface.
		isVarDecl()
	}

	FuncOrSub interface {
		Node
		// Name returns the name of the declared sub of function.
		Name() *Identifier
		// GetBody returns the body of the function or sub
		GetBody() []StatementList
		// Type returns the type of the declared sub or function.
		Type() (types.Type, error)
		// list of parameters
		GetParams() []ParamItem
	}

	HasBody interface {
		Node
		// GetBody returns the body of the function or sub
		GetBody() []StatementList
	}
)

// Declaration nodes.
type (
	// A FuncDecl node represents a function declaration.
	//
	// Examples.
	//
	//    Function add(a As Integer, b As Integer) As Integer
	FuncDecl struct {
		// start position of the function keyword
		FunctionKw *token.Token
		// Function signature.
		FuncType *FuncType
		// Function name.
		FuncName *Identifier
		// Function body
		Body []StatementList
		// parent node
		Parent Node
	}

	// comment is kept to allow syntax highlighting
	Comment struct {
		// Position of the comment
		Text *token.Token
		// parent node
		Parent Node
	}

	// A SubDecl node represents a Subroutine declaration.
	//
	// Examples.
	//
	//    Sub add(a As Integer, b As Integer)
	SubDecl struct {
		// start position of the sub keyword
		SubKw *token.Token
		// Subroutine signature.
		SubType *SubType
		// Sub name.
		SubName *Identifier
		// Sub body
		Body []StatementList
		// parent node
		Parent Node
	}

	// A ParamItem node represents a variable declaration.
	//
	// Examples.
	//
	// X as Integer
	// Optional X as Integer = 42
	// ByVal X as Integer
	// Optional ByVal X as Integer = 42
	// ParamArray X
	ParamItem struct {
		// Variable name.
		VarName *Identifier
		// is optional
		Optional bool
		// byVal - passed by value
		ByVal bool
		// paramArray
		ParamArray bool
		// is an array
		IsArray bool
		// default value if optional
		DefaultValue Expression
		// Variable type.
		VarType Type
		// parent node
		Parent Node
	}

	// A ConstDeclItem node represents a constant declaration.
	//
	// Examples.
	//
	// X as Integer = 42
	ConstDeclItem struct {
		// Variable name.
		ConstName *Identifier
		// Variable value expression; or nil if variable declaration (i.e. not
		// variable definition).
		ConstValue Expression
		// Variable type.
		ConstType Type
		// parent node
		Parent Node
	}

	// A DimDecl node represents a variable declaration.
	//
	// Examples.
	//
	//    Dim x As Integer
	//    Dim x As Integer = 42
	//    Dim x As Integer(10)
	//    Dim x As Integer(10, 20)
	DimDecl struct {
		// Position of `dim` keyword.
		DimKw *token.Token
		// Variable declarations
		Vars []VarDecl
		// parent node
		Parent Node
	}

	// a ConstDecl node represents a constant declaration.
	//
	// Examples.
	//
	// Const X as Integer = 42
	ConstDecl struct {
		// Position of `const` keyword.
		ConstKw *token.Token
		// Const declarations
		Consts []ConstDeclItem
		// parent node
		Parent Node
	}

	// ArrayDecl node represents an array declaration.
	//
	// Examples.
	//
	//    Dim x As Long(10)
	//    Dim x As Long(10, 20)
	ArrayDecl struct {
		// Variable name.
		VarName *Identifier
		// Arrray
		VarType *ArrayType
		// parent node
		Parent Node
	}

	// ScalarDecl node represents an array declaration.
	//
	// Examples.
	//
	//    Dim x As Double
	//    Dim x As Long = 0
	ScalarDecl struct {
		// Variable name.
		VarName *Identifier
		// Variable type.
		VarType Type
		// Initial Value
		VarValue Expression
		// parent node
		Parent Node
	}

	// An EnumDecl node represents an enum declaration.
	//
	// Examples.
	//
	//    Enum Animal
	//        Cat
	//        Dog
	//    End Enum
	EnumDecl struct {
		// Position of `enum` keyword.
		EnumKw *token.Token
		// Enum name.
		Identifier *Identifier
		// Enum values.
		Values []Identifier
		// parent node
		Parent Node
	}

	// an ClassDecl node represents a class declaration
	//
	// Examples.
	//
	//    Class Animal
	//        Public Sub New()
	//        End Sub
	//    End Class
	ClassDecl struct {
		// Position of `class` keyword.
		ClassKw *token.Token
		// Class name.
		ClassName *Identifier
		// Class body
		Members map[string]Decl
		// parent node
		Parent Node
	}

	// A JumpLabelDecl node represents a label declaration.
	//
	// Examples.
	//
	//    Label:
	JumpLabelDecl struct {
		// Label name.
		Label *Identifier
		// parent node
		Parent Node
	}
)

// A ForExpr node represents a for loop expression. It has one of the following
// underlying types.
//
//	*ForNextExpr
//	*ForEachExpr
//
// Pseudo-code representation of a declaration.
//
//	for ident = start to end [step step]
//	for each ident in collection
type ForExpr interface {
	Node
	// isForExpr ensures that only for loop expression nodes can be assigned to
	// the ForExpr interface.
	isForExpr()
}

// For loop expression nodes.
type (
	// A ForNextExpr node represents a for loop expression.
	//
	// Examples.
	//
	//    for i = 0 to 10 step 2
	ForNextExpr struct {
		// Loop variable.
		Variable *Identifier
		// From value.
		From Expression
		// To value.
		To Expression
		// Step value.
		Step Expression
		// parent node
		Parent Node
	}

	// A ForEachExpr node represents a for each loop expression.
	//
	// Examples.
	//
	//    for each i in collection
	ForEachExpr struct {
		// Loop variable.
		Variable *Identifier
		// Collection.
		Collection Expression
		// parent node
		Parent Node
	}
)

// A StatementList represents a list of statements.
type StatementList struct {
	// List of statements.
	Statements []Statement
	// parent node
	Parent Node
}

// A Statement node represents a statement, and has one of the following underlying
// types.
//
//		*EmptyStmt
//		*ExprStmt
//		*IfStmt
//		*SelectStmt
//		*ExitStmt
//		*SpecialStmt
//	 	*DoWhileStmt
//		*DoUntilStmt
//		*WhileStmt
//	    *UntilStmt
//		*ForStmt
//	    *DimDecl
//		*ConstDecl
//		*EnumDecl
//	    *LabelDecl
//	    *CallSubStmt
type Statement interface {
	Node
	// isStmt ensures that only statement nodes can be assigned to the Stmt
	// interface.
	isStmt()
}

// Statement nodes.
type (
	// An ExitStmt node represents one of the following:
	//
	//    Exit Function
	//    Exit Sub
	//    Exit For
	//    Exit Do
	ExitStmt struct {
		// Exit keyword.
		ExitKw *token.Token
		// Exit type
		ExitType *token.Token
		// parent node
		Parent Node
	}

	// A SpecialStmt node represents a call statement
	// to one of the following:
	//
	//    Debug.Print "hello world"
	//    MsgBox "hello world", ubOkOnly
	// 	  Redim Preserve x(10)
	SpecialStmt struct {
		Keyword1 *token.Token
		// keyword2 (used to store "Preserve" in Redim Preserve)
		Keyword2 string
		// Function arguments.
		Args []Expression
		// semicolon
		Semicolon *token.Token // optional, to keep the print statement on the same line
		// parent node
		Parent Node
	}

	// An EmptyStmt node represents an empty statement (i.e. an empty line).
	//
	// Examples.
	//
	//
	EmptyStmt struct {
		// Position of semicolon end-of-line.
		EOL *token.Token
		// parent node
		Parent Node
	}

	// An ExprStmt node represents a stand-alone expression in a statement list.
	//
	// Examples.
	//
	//    42
	//    f()
	ExprStmt struct {
		// Stand-alone expression.
		Expression Expression
		// parent node
		Parent Node
	}

	// An IfStmt node represents an if statement.
	//
	// Examples.
	//
	//    If isPrime(x) Then
	//        Print x
	//    Else
	//        Print "not prime"
	//    End If
	IfStmt struct {
		// Position of `if` keyword.
		IfKw *token.Token
		// Condition.
		Condition Expression
		// True branch.
		Body []StatementList
		// elsif branches
		ElseIf []ElseIfStmt
		// False branch; or nil if 1-way conditional.
		Else []StatementList
		// parent node
		Parent Node
	}

	// An SelectStmt node represents a select statement.
	//
	// Examples.
	//
	//    Select Case x
	//    Case 1
	//        Print "one"
	//    Case 2
	//        Print "two"
	//    Case Else
	//        Print "other"
	//    End Select
	SelectStmt struct {
		// Position of `select` keyword.
		SelectKw *token.Token
		// Condition.
		Condition Expression
		// Case branches.
		Body []CaseStmt
		// parent node
		Parent Node
	}

	// A CaseStmt node represents a case expression.
	//
	// Examples.
	//
	//    Case 1
	//        Print "one"
	//    Case 2
	//        Print "two"
	//    Case Else
	//        Print "other"
	CaseStmt struct {
		// Position of `case` keyword.
		CaseKw *token.Token
		// Condition.
		Condition Expression
		// Case body.
		Body []StatementList
		// parent node
		Parent Node
	}

	// An ElseIfStmt node represents an else if statement.
	//
	// Examples.
	//
	//    ElseIf isPrime(x) Then
	//        Print x
	ElseIfStmt struct {
		// Position of `else` keyword.
		ElseIfKw *token.Token
		// Condition.
		Condition Expression
		// True branch.
		Body []StatementList
		// parent node
		Parent Node
	}

	// A WhileStmt node represents a Do While statement.
	//
	// Examples.
	//
	//    "Do" "While" Expr "\n" Stmt "Loop"
	WhileStmt struct {
		// Position of `do` keyword.
		DoKw *token.Token
		// Condition.
		Condition Expression
		// Loop body.
		Body []StatementList
		// parent node
		Parent Node
	}

	// A UntilStmt node represents a Do Until statement.
	//
	// Examples.
	//
	//    "Do" "Until" Expr "\n" Stmt "Loop"
	UntilStmt struct {
		// Position of `do` keyword.
		DoKw *token.Token
		// Condition.
		Condition Expression
		// Loop body.
		Body []StatementList
		// parent node
		Parent Node
	}

	// A DoWhileStmt node represents a Do While statement.
	//
	// Examples.
	//
	//    "Do" Stmt "Loop" "While" Expr "\n"
	DoWhileStmt struct {
		// Position of `do` keyword.
		DoKw *token.Token
		// Condition.
		Condition Expression
		// Loop body.
		Body []StatementList
		// parent node
		Parent Node
	}

	// A DoUntilStmt node represents a Do Until statement.
	//
	// Examples.
	//
	//    "Do" Stmt "Loop" "Until" Expr "\n"
	DoUntilStmt struct {
		// Position of `do` keyword.
		DoKw *token.Token
		// Condition.
		Condition Expression
		// Loop body.
		Body []StatementList
		// parent node
		Parent Node
	}

	// A ForStmt node represents a for loop statement.
	//
	// Examples.
	//
	//    "For" Expr "To" Expr ["Step" Expr] "\n" ClosedStmt "Next"
	ForStmt struct {
		// Position of `for` keyword.
		ForKw *token.Token
		// Loop variable.
		ForExpression ForExpr
		// Loop body.
		Body []StatementList
		// Next var
		Next *Identifier
		// parent node
		Parent Node
	}
)

// An Expression node represents an expression, and has one of the following
// underlying types.
//
//	*BasicLit
//	*BinaryExpr
//	*Identifier
//	*IndexExpr
//	*ParenExpr
//	*UnaryExpr
type Expression interface {
	Node
	// isExpr ensures that only expression nodes can be assigned to the Expr
	// interface.
	isExpr()
	// Type returns the type of the expression.
	Type() (types.Type, error)
}

// Expression nodes.
type (
	// A BasicLit node represents a basic literal.
	//
	// Examples.
	//
	//    42
	//    45.9
	//    "hello world"
	//    #2019-01-01 00:00:00#
	//    True
	//	  Nothing
	BasicLit struct {
		// Position of basic literal.
		ValTok *token.Token
		// Basic literal type, one of the following.
		//
		//  token.IntLit
		//	token.FloatLit
		//	token.StringLit
		//	token.DateTimeLit
		//	token.BooleanLit
		//	token.NothingLit
		//  token.SubLit		pseudo-literal for sub return type
		Kind token.Kind
		// Basic literal value; e.g. 123, "hello world"
		Value any
		// parent node
		Parent Node
	}

	// An BinaryExpr node represents a binary expression; X op Y.
	//
	// Examples.
	//
	//    x + y
	//    x = 42
	BinaryExpr struct {
		// First operand.
		Left Expression
		// Position of binary operator.
		OpToken *token.Token
		// Operator, one of the following.
		//    token.Add      // +
		//    token.Sub      // -
		//    token.Mul      // *
		//    token.Div      // /
		//    token.Lt       // <
		//    token.Gt       // >
		//    token.Le       // <=
		//    token.Ge       // >=
		//    token.Ne       // <>
		//    token.Eq       // ==
		//    token.Land     // &&
		//    token.Assign   // =
		//    token.Concat   // &
		//    token.Eqv      // =
		OpKind token.Kind
		// Second operand.
		Right Expression
		// parent node
		Parent Node
	}

	// A CallOrIndexExpr node represents a call expression
	// or an array index expression.
	//
	// Examples.
	//
	//    foo()
	//    bar(42)
	CallOrIndexExpr struct {
		// Function name.
		Identifier *Identifier
		// Position of left-parenthesis `(`.
		Lparen *token.Token
		// Function arguments.
		Args []Expression
		// Position of right-parenthesis `)`.
		Rparen token.Position
		// parent node
		Parent Node
	}

	// A callSubStmt node represents a call to a subroutine.
	//
	// Examples.
	//
	//    foo
	//    bar 42
	CallSubStmt struct {
		// call token
		CallKw *token.Token
		// subroutine name.
		Definition Expression
		// parent node
		Parent Node
	}

	// A JumpStmt node represents a jump statement.
	//
	// Examples.
	//
	//    Goto Label
	//    On error resumn Label
	JumpStmt struct {
		// on error indicator
		OnError bool
		// jump token
		JumpKw *token.Token
		// label name.
		Label *Identifier
		// next kw
		NextKw *token.Token
		// numeric label
		Number *token.Token
		// parent node
		Parent Node
	}

	// A CallSelector node represents an object method call expression.
	//
	// Examples.
	//
	//    foo.bar
	//    foo().bar
	CallSelectorExpr struct {
		// root name.
		Root Expression
		// Position of the dot `.`.
		Dot token.Position
		// Function arguments.
		Selector Expression
		// parent node
		Parent Node
	}

	// An Identifier node represents an identifier.
	//
	// Examples.
	//
	//    x
	//    int
	Identifier struct {
		// Position of identifier.
		Tok *token.Token
		// Identifier name.
		Name string
		// Corresponding function, variable or type definition. The declaration
		// mapping is added during the semantic analysis phase, based on the
		// lexical scope of the identifier.
		Decl Decl
		// parent node
		Parent Node
	}

	// A ParenExpr node represents a parenthesised expression.
	ParenExpr struct {
		// Position of left-parenthesis `(`.
		Lparen *token.Token
		// Parenthesised expression.
		Expr Expression
		// Position of right-parenthesis `)`.
		Rparen token.Position
		// parent node
		Parent Node
	}

	// An UnaryExpr node represents an unary expression; op X.
	//
	// Examples.
	//
	//    -42
	//    !(x == 3 || x == 10)
	UnaryExpr struct {
		// Position of unary operator.
		OpToken *token.Token
		// Operator, one of the following.
		//    token.Sub   // -
		//    token.Not   // !
		OpKind token.Kind
		// Operand.
		Right Expression
		// parent node
		Parent Node
	}
)

// A Type node represents a type of µC, and has one of the following underlying
// types.
//
//		*ArrayType
//		*FuncType
//	 *SubType
//		*Ident
type Type interface {
	Node
	// isType ensures that only type nodes can be assigned to the Type interface.
	isType()
}

// Type nodes.
type (
	// A user-defined type.
	UserDefinedType struct {
		// identifier
		Identifier *Identifier
		// Parent node
		Parent Node
	}

	// A FuncType node represents a function signature.
	//
	// Examples.
	//
	//    (first as Integer, second as Integer) as Integer
	FuncType struct {
		// Position of left-parenthesis `(`.
		Lparen *token.Token
		// Function parameters.
		Params []ParamItem
		// Position of right-parenthesis `)`.
		Rparen token.Position
		// Return type.
		Result Type
		// parent node
		Parent Node
	}

	// A SubType node represents a Subroutine signature.
	//
	// Examples.
	//
	//    (first as Integer, second as Integer) as Integer
	SubType struct {
		// Position of left-parenthesis `(`.
		Lparen *token.Token
		// Function parameters.
		Params []ParamItem
		// Position of right-parenthesis `)`.
		Rparen token.Position
		// parent node
		Parent Node
	}

	// An ArrayType node represents an array type.
	//
	// Examples.
	//
	//    Integer(10)
	//    Integer(10, 20)
	ArrayType struct {
		// Position of left-parenthesis `(`.
		Lparen *token.Token
		// Array dimensions.
		Dimensions []Expression
		// Position of right-parenthesis `)`.
		Rparen *token.Position
		// Element type.
		Type Type
		// parent node
		Parent Node
	}
)

func (n *TypeDef) String() string {
	return fmt.Sprintf("typedef %v %v", n.DeclType, n.TypeName)
}

func (n *BasicLit) String() string {
	return fmt.Sprint(n.Value)
}

func (n *BinaryExpr) String() string {
	return fmt.Sprintf("%v %v %v", n.Left, n.OpKind, n.Right)
}

func (n *ForNextExpr) String() string {
	buf := strings.Builder{}
	buf.WriteString(n.Variable.String())
	buf.WriteString(" = ")
	buf.WriteString(n.From.String())
	buf.WriteString(" To ")
	buf.WriteString(n.To.String())
	if n.Step != nil {
		buf.WriteString(" Step ")
		buf.WriteString(n.Step.String())
	}
	return buf.String()
}

func (n *ForEachExpr) String() string {
	buf := strings.Builder{}
	buf.WriteString("Each ")
	buf.WriteString(n.Variable.String())
	buf.WriteString(" In ")
	buf.WriteString(n.Collection.String())
	return buf.String()
}

func (n *ExitStmt) String() string {
	buf := strings.Builder{}
	buf.WriteString("Exit ")
	buf.WriteString(n.ExitType.Literal)
	return buf.String()
}

func (n *SpecialStmt) String() string {
	buf := strings.Builder{}
	buf.WriteString(strcase.ToCamel(n.Keyword1.Literal))
	if len(n.Keyword2) > 0 {
		buf.WriteString(" ")
		buf.WriteString(strcase.ToCamel(n.Keyword2))
	}
	buf.WriteString(" ")
	for i, arg := range n.Args {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(arg.String())
	}
	if n.Semicolon != nil {
		buf.WriteString(";")
	}
	return buf.String()
}

func (n *CallSubStmt) String() string {
	buf := strings.Builder{}
	buf.WriteString("Call ")
	buf.WriteString(n.Definition.String())
	return buf.String()
}

func (n *JumpStmt) String() string {
	buf := strings.Builder{}
	if n.OnError {
		buf.WriteString("On Error ")
	}
	buf.WriteString(n.JumpKw.Literal)
	if n.NextKw != nil {
		buf.WriteString(" ")
		buf.WriteString(n.NextKw.Literal)
	}
	if n.Label != nil {
		buf.WriteString(" ")
		buf.WriteString(n.Label.String())
	}
	if n.Number != nil {
		buf.WriteString(" ")
		buf.WriteString(n.Number.Literal)
	}
	return buf.String()
}

func (n *CallOrIndexExpr) String() string {
	buf := strings.Builder{}
	buf.WriteString(n.Identifier.String())
	buf.WriteString("(")
	for i, arg := range n.Args {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(arg.String())
	}
	buf.WriteString(")")
	return buf.String()
}

func (n *CallSelectorExpr) String() string {
	buf := strings.Builder{}
	buf.WriteString(n.Root.String())
	buf.WriteString(".")
	buf.WriteString(n.Selector.String())
	return buf.String()
}

func (n *EmptyStmt) String() string {
	return ""
}

func (n *ExprStmt) String() string {
	return fmt.Sprintf("Let %v", n.Expression)
}

func (n *File) String() string {
	buf := strings.Builder{}
	for _, node := range n.Body {
		buf.WriteString(node.String())
	}
	return buf.String()
}

func (n *StatementList) String() string {
	buf := strings.Builder{}
	for i, node := range n.Statements {
		if i != 0 {
			buf.WriteString(" : ")
		}
		buf.WriteString(node.String())
	}
	buf.WriteString("\n")
	return buf.String()
}

func (n *UserDefinedType) String() string {
	return n.Identifier.Name
}

func (n *FuncDecl) String() string {
	buf := strings.Builder{}
	buf.WriteString("Function ")
	buf.WriteString(n.FuncName.String())
	buf.WriteString("(")
	for i, param := range n.FuncType.Params {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(param.String())
	}
	buf.WriteString(") As ")
	buf.WriteString(n.FuncType.Result.String())
	buf.WriteString("\n")
	for _, stmt := range n.Body {
		buf.WriteString(stmt.String())
	}
	buf.WriteString("End Function")
	return buf.String()
}

func (n *SubDecl) String() string {
	buf := strings.Builder{}
	buf.WriteString("Sub ")
	buf.WriteString(n.SubName.Name)
	buf.WriteString("(")
	for i, param := range n.SubType.Params {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(param.String())
	}
	buf.WriteString("\n")
	for _, stmt := range n.Body {
		buf.WriteString(stmt.String())
	}
	buf.WriteString("End Sub")
	return buf.String()
}

func (n *Comment) String() string {
	return n.Text.Literal
}

func (n *FuncType) String() string {
	buf := strings.Builder{}
	buf.WriteString("(")
	for i, param := range n.Params {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(param.VarType.String())
	}
	buf.WriteString(")")
	buf.WriteString(n.Result.String())
	return buf.String()
}

func (n *SubType) String() string {
	buf := strings.Builder{}
	buf.WriteString("(")
	for i, param := range n.Params {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(param.VarType.String())
	}
	buf.WriteString(")")
	return buf.String()
}

func (n *Identifier) String() string {
	return n.Name
}

func (n *IfStmt) String() string {
	buf := strings.Builder{}
	buf.WriteString("If ")
	buf.WriteString(n.Condition.String())
	buf.WriteString(" Then\n")
	for _, stmt := range n.Body {
		buf.WriteString(stmt.String())
	}
	for _, stmt := range n.ElseIf {
		buf.WriteString(stmt.String())
	}

	if n.Else != nil {
		buf.WriteString("Else\n")
		for _, stmt := range n.Else {
			buf.WriteString(stmt.String())
		}
	}
	buf.WriteString("End If")
	return buf.String()
}

func (n *SelectStmt) String() string {
	buf := strings.Builder{}
	buf.WriteString("Select Case ")
	buf.WriteString(n.Condition.String())
	buf.WriteString("\n")
	for _, stmt := range n.Body {
		buf.WriteString(stmt.String())
	}
	buf.WriteString("End Select")
	return buf.String()
}

func (n *CaseStmt) String() string {
	buf := strings.Builder{}
	if n.Condition == nil {
		buf.WriteString("Case Else\n")
	} else {
		buf.WriteString("Case ")
		buf.WriteString(n.Condition.String())
		buf.WriteString("\n")
	}
	for _, stmt := range n.Body {
		buf.WriteString(stmt.String())
	}
	return buf.String()
}

func (n *ElseIfStmt) String() string {
	buf := strings.Builder{}
	buf.WriteString("ElseIf ")
	buf.WriteString(n.Condition.String())
	buf.WriteString(" Then\n")
	for _, stmt := range n.Body {
		buf.WriteString(stmt.String())
	}
	return buf.String()
}
func (n *ParenExpr) String() string {
	return fmt.Sprintf("(%v)", n.Expr)
}

func (n *UnaryExpr) String() string {
	return fmt.Sprintf("%v%v", n.OpKind, n.Right)
}

func (n *EnumDecl) String() string {
	buf := strings.Builder{}
	buf.WriteString("Enum ")
	buf.WriteString(n.Identifier.String())
	buf.WriteString("\n")
	for _, value := range n.Values {
		buf.WriteString(value.String())
		buf.WriteString("\n")
	}
	buf.WriteString("End Enum")
	return buf.String()
}

func (n *JumpLabelDecl) String() string {
	return n.Label.String() + ":"
}

func (n *ConstDecl) String() string {
	buf := strings.Builder{}
	buf.WriteString("Const ")
	for i, constDecl := range n.Consts {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(constDecl.String())
	}
	return buf.String()
}

func (n *ArrayDecl) String() string {
	buf := strings.Builder{}
	buf.WriteString(n.VarName.String())
	buf.WriteString(n.VarType.String())
	return buf.String()
}

func (n *ArrayType) String() string {
	buf := strings.Builder{}
	buf.WriteString("(")
	for i, dim := range n.Dimensions {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(dim.String())
	}
	buf.WriteString(")")
	if n.Type != nil {
		buf.WriteString(" As ")
		buf.WriteString(n.Type.String())
	}
	return buf.String()
}

func (n *ScalarDecl) String() string {
	buf := strings.Builder{}
	buf.WriteString(n.VarName.String())
	if n.VarType != nil {
		buf.WriteString(" As ")
		buf.WriteString(n.VarType.String())
	}
	if n.VarValue != nil {
		buf.WriteString(" = ")
		buf.WriteString(n.VarValue.String())
	}
	return buf.String()
}

func (n *ConstDeclItem) String() string {
	buf := strings.Builder{}
	buf.WriteString(n.ConstName.String())
	if n.ConstType != nil {
		buf.WriteString(" As ")
		buf.WriteString(n.ConstType.String())
	}
	buf.WriteString(" = ")
	buf.WriteString(n.ConstValue.String())
	return buf.String()
}

func (n *DimDecl) String() string {
	buf := strings.Builder{}
	buf.WriteString("Dim ")
	for i, varDecl := range n.Vars {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(varDecl.String())
	}
	return buf.String()
}

func (n *ParamItem) String() string {
	buf := strings.Builder{}
	if n.Optional {
		buf.WriteString("Optional ")
	}
	if n.ByVal {
		buf.WriteString("ByVal ")
	} else {
		buf.WriteString("ByRef ")
	}
	if n.ParamArray {
		buf.WriteString("ParamArray ")
	}
	buf.WriteString(n.VarName.String())
	if n.IsArray {
		buf.WriteString("()")
	}
	if n.VarType != nil {
		buf.WriteString(" As ")
		buf.WriteString(n.VarType.String())
	}
	if n.DefaultValue != nil {
		buf.WriteString(" = ")
		buf.WriteString(n.DefaultValue.String())
	}
	return buf.String()
}

func (n *WhileStmt) String() string {
	buf := strings.Builder{}
	buf.WriteString("Do While ")
	buf.WriteString(n.Condition.String())
	buf.WriteString("\n")
	for _, stmt := range n.Body {
		buf.WriteString(stmt.String())
	}
	buf.WriteString("Loop")
	return buf.String()
}

func (n *UntilStmt) String() string {
	buf := strings.Builder{}
	buf.WriteString("Do Until ")
	buf.WriteString(n.Condition.String())
	buf.WriteString("\n")
	for _, stmt := range n.Body {
		buf.WriteString(stmt.String())
	}
	buf.WriteString("Loop")
	return buf.String()
}

func (n *DoWhileStmt) String() string {
	buf := strings.Builder{}
	buf.WriteString("Do\n")
	for _, stmt := range n.Body {
		buf.WriteString(stmt.String())
	}
	buf.WriteString("Loop While ")
	buf.WriteString(n.Condition.String())
	return buf.String()
}

func (n *DoUntilStmt) String() string {
	buf := strings.Builder{}
	buf.WriteString("Do\n")
	for _, stmt := range n.Body {
		buf.WriteString(stmt.String())
	}
	buf.WriteString("Loop Until ")
	buf.WriteString(n.Condition.String())
	return buf.String()
}

func (n *ForStmt) String() string {
	buf := strings.Builder{}
	buf.WriteString("For ")
	buf.WriteString(n.ForExpression.String())
	buf.WriteString("\n")
	for _, stmt := range n.Body {
		buf.WriteString(stmt.String())
	}
	buf.WriteString("Next")
	if n.Next != nil {
		buf.WriteString(" ")
		buf.WriteString(n.Next.String())
	}
	return buf.String()
}

func (n *ClassDecl) String() string {
	buf := strings.Builder{}
	buf.WriteString("Class ")
	buf.WriteString(n.ClassName.String())
	buf.WriteString("\n")
	for _, member := range n.Members {
		buf.WriteString(member.String())
		buf.WriteString("\n")
	}
	buf.WriteString("End Class")
	return buf.String()
}

// Token returns the first token of the Node
func (n *StatementList) Token() *token.Token {
	if len(n.Statements) > 0 {
		return n.Statements[0].Token()
	}
	return nil
}

// Token returns the first token of the Node
func (n *TypeDef) Token() *token.Token {
	return n.Typedef
}

// Token returns the first token of the Node
func (n *ConstDecl) Token() *token.Token {
	return n.ConstKw
}

// // Token returns the first token of the Node
func (n *ArrayType) Token() *token.Token {
	return n.Lparen
}

// Token returns the first token of the Node
func (n *BasicLit) Token() *token.Token {
	return n.ValTok
}

// Token returns the first token of the Node
func (n *BinaryExpr) Token() *token.Token {
	return n.Left.Token()
}

// Token returns the first token of the Node
func (n *ForNextExpr) Token() *token.Token {
	return n.Variable.Token()
}

// Token returns the first token of the Node
func (n *ForEachExpr) Token() *token.Token {
	return n.Variable.Token()
}

// Token returns the first token of the Node
func (n *ExitStmt) Token() *token.Token {
	return n.ExitKw
}

// Token returns the first token of the Node
func (n *SpecialStmt) Token() *token.Token {
	return n.Keyword1
}

// Token returns the first token of the Node
func (n *CallSubStmt) Token() *token.Token {
	return n.CallKw
}

// Token returns the first token of the Node
func (n *JumpStmt) Token() *token.Token {
	return n.JumpKw
}

// Token returns the first token of the Node
func (n *CallOrIndexExpr) Token() *token.Token {
	return n.Identifier.Token()
}

// Token returns the first token of the Node
func (n *CallSelectorExpr) Token() *token.Token {
	return n.Root.Token()
}

// Token returns the first token of the Node
func (n *EmptyStmt) Token() *token.Token {
	return n.EOL
}

// Token returns the first token of the Node
func (n *ExprStmt) Token() *token.Token {
	return n.Expression.Token()
}

// Token returns the first token of the Node
func (n *File) Token() *token.Token {
	if len(n.Body) > 0 {
		return n.Body[0].Token()
	}
	return nil
}

// Token returns the first token of the Node
func (n *UserDefinedType) Token() *token.Token {
	return n.Identifier.Token()
}

// Token returns the first token of the Node
func (n *FuncDecl) Token() *token.Token {
	return n.FunctionKw
}

// Token returns the first token of the Node
func (n *SubDecl) Token() *token.Token {
	return n.SubKw
}

// Token returns the first token of the Node
func (n *Comment) Token() *token.Token {
	return n.Text
}

// Token returns the first token of the Node
func (n *FuncType) Token() *token.Token {
	return n.Lparen
}

// Token returns the first token of the Node
func (n *SubType) Token() *token.Token {
	return n.Lparen
}

// Token returns the first token of the Node
func (n *Identifier) Token() *token.Token {
	return n.Tok
}

// Token returns the first token of the Node
func (n *IfStmt) Token() *token.Token {
	return n.IfKw
}

// Token returns the first token of the Node
func (n *SelectStmt) Token() *token.Token {
	return n.SelectKw
}

// Token returns the first token of the Node
func (n *CaseStmt) Token() *token.Token {
	return n.CaseKw
}

// Token returns the first token of the Node
func (n *ElseIfStmt) Token() *token.Token {
	return n.ElseIfKw
}

// Token returns the first token of the Node
func (n *ParenExpr) Token() *token.Token {
	return n.Lparen
}

// Token returns the first token of the Node
func (n *UnaryExpr) Token() *token.Token {
	return n.OpToken
}

// Token returns the first token of the Node
func (n *EnumDecl) Token() *token.Token {
	return n.EnumKw
}

// Token returns the first token of the Node
func (n *JumpLabelDecl) Token() *token.Token {
	return n.Label.Token()
}

// Token returns the first token of the Node
func (n *ScalarDecl) Token() *token.Token {
	return n.VarName.Token()
}

// Token returns the first token of the Node
func (n *ArrayDecl) Token() *token.Token {
	return n.VarName.Token()
}

// Token returns the first token of the Node
func (n *ConstDeclItem) Token() *token.Token {
	return n.ConstName.Token()
}

// Token returns the first token of the Node
func (n *DimDecl) Token() *token.Token {
	return n.DimKw
}

// Token returns the first token of the Node
func (n *ParamItem) Token() *token.Token {
	return n.VarType.Token()
}

// Token returns the first token of the Node
func (n *WhileStmt) Token() *token.Token {
	return n.DoKw
}

// Token returns the first token of the Node
func (n *UntilStmt) Token() *token.Token {
	return n.DoKw
}

// Token returns the first token of the Node
func (n *DoWhileStmt) Token() *token.Token {
	return n.DoKw
}

// Token returns the first token of the Node
func (n *DoUntilStmt) Token() *token.Token {
	return n.DoKw
}

// Token returns the first token of the Node
func (n *ForStmt) Token() *token.Token {
	return n.ForKw
}

// Token returns the first token of the Node
func (n *ClassDecl) Token() *token.Token {
	return n.ClassKw
}

// GetParent returns the parent node of the current node.
func (n *StatementList) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *TypeDef) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *ConstDecl) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *BasicLit) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *BinaryExpr) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *ForNextExpr) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *ForEachExpr) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *ExitStmt) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *SpecialStmt) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *CallSubStmt) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *JumpStmt) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *CallOrIndexExpr) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *CallSelectorExpr) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *EmptyStmt) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *ExprStmt) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *File) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *FuncDecl) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *SubDecl) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *Comment) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *FuncType) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *SubType) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *Identifier) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *IfStmt) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *SelectStmt) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *CaseStmt) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *ElseIfStmt) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *ParenExpr) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *UnaryExpr) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *EnumDecl) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *JumpLabelDecl) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *ScalarDecl) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *ArrayDecl) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *ArrayType) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *ConstDeclItem) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *DimDecl) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *ParamItem) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *WhileStmt) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *UntilStmt) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *DoWhileStmt) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *DoUntilStmt) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *ForStmt) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *UserDefinedType) GetParent() Node {
	return n.Parent
}

// GetParent returns the parent node of the current node.
func (n *ClassDecl) GetParent() Node {
	return n.Parent
}

// SetParent sets the parent node of the current node.
func (n *StatementList) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *TypeDef) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *ConstDecl) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *BasicLit) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *BinaryExpr) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *ForNextExpr) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *ForEachExpr) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *ExitStmt) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *SpecialStmt) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *CallSubStmt) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *JumpStmt) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *CallOrIndexExpr) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *CallSelectorExpr) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *EmptyStmt) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *ExprStmt) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *File) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *FuncDecl) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *SubDecl) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *Comment) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *FuncType) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *SubType) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *Identifier) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *IfStmt) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *SelectStmt) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *CaseStmt) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *ElseIfStmt) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *ParenExpr) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *UnaryExpr) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *EnumDecl) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *JumpLabelDecl) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *ScalarDecl) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *ArrayDecl) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *ArrayType) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *ConstDeclItem) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *DimDecl) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *ParamItem) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *WhileStmt) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *UntilStmt) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *DoWhileStmt) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *DoUntilStmt) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *ForStmt) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *UserDefinedType) SetParent(node Node) {
	n.Parent = node
}

// SetParent sets the parent node of the current node.
func (n *ClassDecl) SetParent(node Node) {
	n.Parent = node
}

// Verify that all nodes implement the Node interface.
var (
	_ Node = &TypeDef{}
	_ Node = &StatementList{}
	_ Node = &ArrayType{}
	_ Node = &BasicLit{}
	_ Node = &BinaryExpr{}
	_ Node = &ExitStmt{}
	_ Node = &SpecialStmt{}
	_ Node = &CallSubStmt{}
	_ Node = &CallOrIndexExpr{}
	_ Node = &CallSelectorExpr{}
	_ Node = &EmptyStmt{}
	_ Node = &ExprStmt{}
	_ Node = &File{}
	_ Node = &FuncDecl{}
	_ Node = &SubDecl{}
	_ Node = &Comment{}
	_ Node = &FuncType{}
	_ Node = &SubType{}
	_ Node = &Identifier{}
	_ Node = &IfStmt{}
	_ Node = &SelectStmt{}
	_ Node = &CaseStmt{}
	_ Node = &ElseIfStmt{}
	_ Node = &ParenExpr{}
	_ Node = &UnaryExpr{}
	_ Node = &DimDecl{}
	_ Node = &ConstDeclItem{}
	_ Node = &ArrayDecl{}
	_ Node = &ScalarDecl{}
	_ Node = &WhileStmt{}
	_ Node = &UntilStmt{}
	_ Node = &DoWhileStmt{}
	_ Node = &DoUntilStmt{}
	_ Node = &ForStmt{}
	_ Node = &EnumDecl{}
	_ Node = &JumpLabelDecl{}
	_ Node = &UserDefinedType{}
	_ Node = &ParamItem{}
	_ Node = &ForEachExpr{}
	_ Node = &ForNextExpr{}
	_ Node = &ClassDecl{}
	_ Node = &CaseStmt{}
	_ Node = &JumpStmt{}
)

// Verify that all for loop expression nodes implement the ForExpr interface.
var (
	_ ForExpr = &ForNextExpr{}
	_ ForExpr = &ForEachExpr{}
)

// Type returns the type of the declared identifier.
func (n *TypeDef) Type() (types.Type, error) {
	if n.Val != nil {
		return n.Val, nil
	}
	var err error
	if n.Val, err = NewType(n.DeclType); err != nil {
		return nil, err
	}
	return n.Val, nil
}

// Type returns the type of the declared identifier.
func (n *ClassDecl) Type() (types.Type, error) {
	return NewType(n)
}

// Type returns the type of the declared identifier.
func (n *FuncDecl) Type() (types.Type, error) {
	return NewType(n.FuncType)
}

// Type returns the type of the declared identifier.
func (n *SubDecl) Type() (types.Type, error) {
	return NewType(n.SubType) // Will return Nothing
}

// Type returns the type of the declared identifier.
func (n *ScalarDecl) Type() (types.Type, error) {
	return NewType(n.VarType)
}

// Type returns the type of the declared identifier.
func (n *ArrayDecl) Type() (types.Type, error) {
	return NewType(n.VarType)
}

// Type returns the type of the declared identifier.
func (n *DimDecl) Type() (types.Type, error) {
	return nil, nil
}

// Type returns the type of the declared identifier.
func (n *ParamItem) Type() (types.Type, error) {
	// check if it is an array
	if n.IsArray {
		return NewType(n)
	}
	if n.VarType != nil {
		return NewType(n.VarType)
	} else {
		return n.VarName.Type()
	}

}

// Type returns the type of the declared constant.
func (n *ConstDeclItem) Type() (types.Type, error) {
	return NewType(n.ConstType)
}

// Type returns the type of the declared identifier.
func (n *EnumDecl) Type() (types.Type, error) {
	typ := &types.UserDefined{Name: n.Identifier.Name}
	return typ, nil
}

// Type returns the type of the declared identifier.
func (n *JumpLabelDecl) Type() (types.Type, error) {
	return nil, nil
}

// Type returns the type of the declared identifier.
func (n *Identifier) Type() (types.Type, error) {
	if n.Decl != nil {
		return n.Decl.Type()
	}
	return nil, nil
}

// Type returns the type of the declared identifier.
func (n *UserDefinedType) Type() (types.Type, error) {
	return NewType(n.Identifier)
}

// Type returns the type of the declared identifier.
func (n *UnaryExpr) Type() (types.Type, error) {
	return NewType(n.Right)
}

// Type returns the type of the declared identifier.
func (n *ParenExpr) Type() (types.Type, error) {
	return NewType(n.Expr)
}

// Type returns the type of the declared identifier.
func (n *BinaryExpr) Type() (types.Type, error) {
	return NewType(n)
}

// Type returns the type of the declared identifier.
func (n *CallOrIndexExpr) Type() (types.Type, error) {
	return NewType(n.Identifier)
}

// Type returns the type of the declared identifier.
func (n *CallSelectorExpr) Type() (types.Type, error) {
	return NewType(n.Selector)
}

// Type returns the type of the declared identifier.
func (n *BasicLit) Type() (types.Type, error) {
	return NewType(n)
}

// Name returns the name of the declared identifier.
func (n *TypeDef) Name() *Identifier {
	return n.TypeName
}

// Name returns the name of the declared identifier.
func (n *FuncDecl) Name() *Identifier {
	return n.FuncName
}

// Name returns the name of the declared identifier.
func (n *SubDecl) Name() *Identifier {
	return n.SubName
}

// Name returns the name of the declared identifier.
func (n *ScalarDecl) Name() *Identifier {
	return n.VarName
}

// Name returns the name of the declared identifier.
func (n *ArrayDecl) Name() *Identifier {
	return n.VarName
}

// Name returns the name of the declared identifier.
func (n *ConstDeclItem) Name() *Identifier {
	return n.ConstName
}

// Name returns the name of the declared identifier.
func (n *EnumDecl) Name() *Identifier {
	return n.Identifier
}

// Name returns the name of the declared identifier.
func (n *JumpLabelDecl) Name() *Identifier {
	return n.Label
}

// Name returns the name of the declared identifier.
func (n *DimDecl) Name() *Identifier {
	return nil
}

// Name returns the name of the declared identifier.
func (n *ParamItem) Name() *Identifier {
	return n.VarName
}

// Name returns the name of the declared identifier.
func (n *ClassDecl) Name() *Identifier {
	return n.ClassName
}

// Value returns the initializing value of the defined identifier; or nil if
// declaration or tentative definition.
//
// Underlying type for type definitions.
//
//	Type
func (n *TypeDef) Value() Node {
	// ref: https://golang.org/doc/faq#nil_error
	if n.DeclType != nil {
		return n.DeclType
	}
	return nil
}

// Value returns the initializing value of the defined identifier; or nil if
// declaration or tentative definition.
//
// Underlying type for function declarations.
func (n *FuncDecl) Value() Node {
	return nil
}

// Value returns the initializing value of the defined identifier; or nil if
// declaration or tentative definition.
//
// Underlying type for sub declarations.
func (n *SubDecl) Value() Node {
	return nil
}

func (n *DimDecl) Value() Node {
	return nil
}

// Value returns the initializing value of the defined identifier; or nil if
// declaration or tentative definition.
//
// Underlying type for variable declarations.
//
//	Expr
func (n *ParamItem) Value() Node {
	// ref: https://golang.org/doc/faq#nil_error
	if n.DefaultValue != nil {
		return n.DefaultValue
	}
	return nil
}

// Value returns the initializing value of the scalar declaration
func (n *ScalarDecl) Value() Node {
	return n.VarValue
}

// Value returns the initializing value of the defined array
//
// As it is not desirable to initialize an array with a value, this method
// always returns nil.
func (n *ArrayDecl) Value() Node {
	return nil
}

// Value returns the initializing value of the defined enum
//
// As it is not desirable to initialize an enum with a value, this method
// always returns nil.
func (n *EnumDecl) Value() Node {
	return nil
}

// Value returns the initializing value of a label
//
// As it is not desirable to initialize a label with a value, this method
// always returns nil.
func (n *JumpLabelDecl) Value() Node {
	return nil
}

// Value returns the initializing value of the defined constant
//
// Underlying type for variable declarations.
//
//	Expr
func (n *ConstDeclItem) Value() Node {
	// ref: https://golang.org/doc/faq#nil_error
	if n.ConstValue != nil {
		return n.ConstValue
	}
	return nil
}

// Value returns the initializing value of the defined identifier; or nil if
// declaration or tentative definition.
//
// Underlying type for type definitions.
//
//	Type
func (n *ClassDecl) Value() Node {
	return nil
}

// isDecl ensures that only declaration nodes can be assigned to the Decl
// interface.
func (n *TypeDef) isDecl()       {}
func (n *FuncDecl) isDecl()      {}
func (n *SubDecl) isDecl()       {}
func (n *DimDecl) isDecl()       {}
func (n *ConstDeclItem) isDecl() {}
func (n *ScalarDecl) isDecl()    {}
func (n *ArrayDecl) isDecl()     {}
func (n *EnumDecl) isDecl()      {}
func (n *JumpLabelDecl) isDecl() {}
func (n *ParamItem) isDecl()     {}
func (n *ClassDecl) isDecl()     {}

// Verify that the declaration nodes implement the Decl interface.
var (
	_ Decl = &TypeDef{}
	_ Decl = &FuncDecl{}
	_ Decl = &SubDecl{}
	_ Decl = &DimDecl{}
	_ Decl = &ConstDeclItem{}
	_ Decl = &ArrayDecl{}
	_ Decl = &ScalarDecl{}
	_ Decl = &EnumDecl{}
	_ Decl = &JumpLabelDecl{}
	_ Decl = &ParamItem{}
	_ Decl = &ClassDecl{}
)

// isStmt ensures that only statement nodes can be assigned to the Stmt
// interface.
func (n *ExitStmt) isStmt()      {}
func (n *SpecialStmt) isStmt()   {}
func (n *CallSubStmt) isStmt()   {}
func (n *EmptyStmt) isStmt()     {}
func (n *ExprStmt) isStmt()      {}
func (n *IfStmt) isStmt()        {}
func (n *SelectStmt) isStmt()    {}
func (n *WhileStmt) isStmt()     {}
func (n *UntilStmt) isStmt()     {}
func (n *DoWhileStmt) isStmt()   {}
func (n *DoUntilStmt) isStmt()   {}
func (n *ForStmt) isStmt()       {}
func (n *DimDecl) isStmt()       {}
func (n *ConstDecl) isStmt()     {}
func (n *EnumDecl) isStmt()      {}
func (n *JumpLabelDecl) isStmt() {}
func (n *FuncDecl) isStmt()      {}
func (n *SubDecl) isStmt()       {}
func (n *JumpStmt) isStmt()      {}
func (n *Comment) isStmt()       {}

// Verify that the statement nodes implement the Stmt interface.
var (
	_ Statement = &ExitStmt{}
	_ Statement = &SpecialStmt{}
	_ Statement = &CallSubStmt{}
	_ Statement = &EmptyStmt{}
	_ Statement = &ExprStmt{}
	_ Statement = &IfStmt{}
	_ Statement = &SelectStmt{}
	_ Statement = &WhileStmt{}
	_ Statement = &UntilStmt{}
	_ Statement = &DoWhileStmt{}
	_ Statement = &DoUntilStmt{}
	_ Statement = &ForStmt{}
	_ Statement = &DimDecl{}
	_ Statement = &ConstDecl{}
	_ Statement = &EnumDecl{}
	_ Statement = &JumpLabelDecl{}
	_ Statement = &FuncDecl{}
	_ Statement = &SubDecl{}
	_ Statement = &JumpStmt{}
	_ Statement = &Comment{}
)

// isExpr ensures that only expression nodes can be assigned to the Expr
// interface.
func (n *BasicLit) isExpr()         {}
func (n *BinaryExpr) isExpr()       {}
func (n *CallOrIndexExpr) isExpr()  {}
func (n *CallSelectorExpr) isExpr() {}
func (n *Identifier) isExpr()       {}
func (n *ParenExpr) isExpr()        {}
func (n *UnaryExpr) isExpr()        {}

// Verify that the expression nodes implement the Expr interface.
var (
	_ Expression = &BasicLit{}
	_ Expression = &BinaryExpr{}
	_ Expression = &CallOrIndexExpr{}
	_ Expression = &CallSelectorExpr{}
	_ Expression = &Identifier{}
	_ Expression = &ParenExpr{}
	_ Expression = &UnaryExpr{}
)

// this should theoricly useless but conversion from interface{} to []Expr is not possible
func ConvertToExprs(value any) ([]Expression, error) {
	exprs := make([]Expression, 0)
	tryExprs, ok := value.([]Expression)
	if ok {
		return tryExprs, nil
	}
	tryExpr, ok := value.(Expression)
	if ok {
		exprs = append(exprs, tryExpr)
		return exprs, nil
	}

	// TODO: might not be necessary
	for _, val := range value.([]any) {
		expression, err := ConvertToExpr(val)
		if err != nil {
			return nil, err
		}
		exprs = append(exprs, expression)
	}
	return exprs, nil
}

// this should theoricly useless but conversion from interface{} to Expr is not possible
func ConvertToExpr(value any) (Expression, error) {
	var expr Expression
	switch value := value.(type) {
	case *BasicLit:
		expr = value
	case *BinaryExpr:
		expr = value
	case *CallOrIndexExpr:
		expr = value
	case *CallSelectorExpr:
		expr = value
	case *Identifier:
		expr = value
	case *ParenExpr:
		expr = value
	case *UnaryExpr:
		expr = value
	default:
		return nil, fmt.Errorf("cannot convert %T to Expr", value)
	}
	return expr, nil
}

// isForExpr ensures that only for loop expression nodes can be assigned to the ForExpr interface.
func (n *ForNextExpr) isForExpr() {}
func (n *ForEachExpr) isForExpr() {}

// isType ensures that only type nodes can be assigned to the Type interface.
func (n *Identifier) isType()      {}
func (n *FuncType) isType()        {}
func (n *SubType) isType()         {}
func (n *UserDefinedType) isType() {}
func (n *ArrayType) isType()       {}

// Verify that the type nodes implement the Type interface.
var (
	_ Type = &Identifier{}
	_ Type = &ArrayType{}
	_ Type = &UserDefinedType{}
	_ Type = &FuncType{}
	_ Type = &SubType{}
)

// isVarDecl ensures that only variable declaration nodes can be assigned to the VarDecl interface.
func (n *ScalarDecl) isVarDecl() {}
func (n *ArrayDecl) isVarDecl()  {}

// Verify that the variable declaration nodes implement the VarDecl interface.
var (
	_ VarDecl = &ScalarDecl{}
	_ VarDecl = &ArrayDecl{}
)

func TokenToIdent(tok any) (*Identifier, error) {
	token, ok := tok.(token.Token)
	if !ok {
		return nil, fmt.Errorf("cannot convert %T to *token.Token", token)
	}
	return &Identifier{Tok: &token, Name: string(token.Literal)}, nil
}

// GetBody returns the body of the statement.
func (n *FuncDecl) GetBody() []StatementList {
	return n.Body
}

// GetBody returns the body of the statement.
func (n *SubDecl) GetBody() []StatementList {
	return n.Body
}

// GetBody returns the body of the statement.
func (n *IfStmt) GetBody() []StatementList {
	return n.Body
}

// GetBody returns the body of the statement.
func (n ElseIfStmt) GetBody() []StatementList {
	return n.Body
}

// GetBody returns the body of the statement.
func (n *WhileStmt) GetBody() []StatementList {
	return n.Body
}

// GetBody returns the body of the statement.
func (n *UntilStmt) GetBody() []StatementList {
	return n.Body
}

// GetBody returns the body of the statement.
func (n *DoWhileStmt) GetBody() []StatementList {
	return n.Body
}

// GetBody returns the body of the statement.
func (n *DoUntilStmt) GetBody() []StatementList {
	return n.Body
}

// GetBody returns the body of the statement.
func (n *ForStmt) GetBody() []StatementList {
	return n.Body
}

// GetParams returns the parameters of the function or sub.
func (n *FuncDecl) GetParams() []ParamItem {
	return n.FuncType.Params
}

// GetParams returns the parameters of the function or sub.
func (n *SubDecl) GetParams() []ParamItem {
	return n.SubType.Params
}

// Verify that the declaration implements function or sub
var (
	_ FuncOrSub = &FuncDecl{}
	_ FuncOrSub = &SubDecl{}
)

// Verify that the declaration implements HasBody
var (
	_ HasBody = &FuncDecl{}
	_ HasBody = &SubDecl{}
	_ HasBody = &IfStmt{}
	_ HasBody = &ElseIfStmt{}
	_ HasBody = &WhileStmt{}
	_ HasBody = &UntilStmt{}
	_ HasBody = &DoWhileStmt{}
	_ HasBody = &DoUntilStmt{}
	_ HasBody = &ForStmt{}
)

func IsNil(n Node) bool {
	switch n := n.(type) {
	case nil:
		return true
	case *TypeDef:
		return n == nil
	case *FuncDecl:
		return n == nil
	case *SubDecl:
		return n == nil
	case *DimDecl:
		return n == nil
	case *ConstDecl:
		return n == nil
	case *ConstDeclItem:
		return n == nil
	case *ScalarDecl:
		return n == nil
	case *ArrayDecl:
		return n == nil
	case *EnumDecl:
		return n == nil
	case *JumpLabelDecl:
		return n == nil
	case *JumpStmt:
		return n == nil
	case *ParamItem:
		return n == nil
	case *StatementList:
		return n == nil
	case *BasicLit:
		return n == nil
	case *BinaryExpr:
		return n == nil
	case *CallOrIndexExpr:
		return n == nil
	case *CallSelectorExpr:
		return n == nil
	case *Identifier:
		return n == nil
	case *ParenExpr:
		return n == nil
	case *UnaryExpr:
		return n == nil
	case *ExitStmt:
		return n == nil
	case *SpecialStmt:
		return n == nil
	case *CallSubStmt:
		return n == nil
	case *EmptyStmt:
		return n == nil
	case *ExprStmt:
		return n == nil
	case *File:
		return n == nil
	case *FuncType:
		return n == nil
	case *SubType:
		return n == nil
	case *IfStmt:
		return n == nil
	case *SelectStmt:
		return n == nil
	case *CaseStmt:
		return n == nil
	case *ElseIfStmt:
		return n == nil
	case *UserDefinedType:
		return n == nil
	case *ArrayType:
		return n == nil
	case *WhileStmt:
		return n == nil
	case *UntilStmt:
		return n == nil
	case *DoWhileStmt:
		return n == nil
	case *DoUntilStmt:
		return n == nil
	case *ForStmt:
		return n == nil
	case *ForEachExpr:
		return n == nil
	case *ForNextExpr:
		return n == nil
	case *Comment:
		return n == nil

	default:
		panic("Unknown type")
	}
}
