package ide

import (
	"fmt"
	"image/color"
	"strings"
	"uBasic/ast"
	"uBasic/token"

	// "gioui.org/font/gofont"
	// "gioui.org/widget/material"
	"gioui.org/unit"
	"gioui.org/x/richtext"
)

const identSize = 4

var (
	typ               = color.NRGBA{G: 170, A: 255}
	keyword           = color.NRGBA{B: 170, A: 255}
	carret            = color.NRGBA{R: 170, A: 255}
	lineNum           = color.NRGBA{A: 255}
	label             = color.NRGBA{R: 0, G: 170, B: 170, A: 255}
	constants         = color.NRGBA{G: 210, A: 255}
	statement         = color.NRGBA{R: 170, B: 170, A: 255}
	identifier        = color.NRGBA{B: 255, R: 150, G: 150, A: 255}
	operator          = color.NRGBA{R: 170, G: 170, A: 255}
	comment           = color.NRGBA{G: 85, A: 255}
	stringCol         = color.NRGBA{R: 170, G: 85, B: 0, A: 255}
	keyword2          = color.NRGBA{R: 70, G: 70, B: 70, A: 255}
	other             = color.NRGBA{R: 100, B: 100, G: 100, A: 255}
	unknown           = color.NRGBA{R: 255, A: 255}
	lineNumber        = 0
	identation        = 0
	currentLineNumber = 0
	size              = unit.Sp(12)
	font              = fonts[2].Font
	fontComment       = fonts[1].Font
)

func ColorText(node ast.Node, currentLine int) []richtext.SpanStyle {
	currentLineNumber = currentLine
	lineNumber = 0
	identation = 0
	return colorText(node)
}

func colorText(node ast.Node) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	if node == nil {
		return styles
	}

	switch node := node.(type) {
	case *ast.CaseStmt:
		styles = append(styles, colorCaseStmt(node)...)
	case *ast.SelectStmt:
		styles = append(styles, colorSelectStmt(node)...)
	case *ast.CallSubStmt:
		styles = append(styles, colorCallSubStmt(node)...)
	case *ast.CallSelectorExpr:
		styles = append(styles, colorCallSelectorExpr(node)...)
	case *ast.WhileStmt:
		styles = append(styles, colorWhileStmt(node)...)
	case *ast.ForStmt:
		styles = append(styles, colorForStmt(node)...)
	case *ast.ForNextExpr:
		styles = append(styles, colorForNextExpr(node)...)
	case *ast.ForEachExpr:
		styles = append(styles, colorForEachExpr(node)...)
	case *ast.UntilStmt:
		styles = append(styles, colorUntilStmt(node)...)
	case *ast.DoWhileStmt:
		styles = append(styles, colorDoWhileStmt(node)...)
	case *ast.DoUntilStmt:
		styles = append(styles, colorDoUntilStmt(node)...)
	case *ast.ExitStmt:
		styles = append(styles, colorExitStmt(node)...)
	case *ast.ParenExpr:
		styles = append(styles, colorParenExpr(node)...)
	case *ast.UnaryExpr:
		styles = append(styles, colorUnaryExpr(node)...)
	case *ast.BinaryExpr:
		styles = append(styles, colorBinaryExpr(node)...)
	case *ast.ExprStmt:
		styles = append(styles, colorExprStmt(node)...)
	case *ast.ConstDecl:
		styles = append(styles, colorConstDecl(node)...)
	case *ast.ConstDeclItem:
		styles = append(styles, colorConstDeclItem(node)...)
	case *ast.ElseIfStmt:
		styles = append(styles, colorElseIfStmt(node)...)
	case *ast.IfStmt:
		styles = append(styles, colorIfStmt(node)...)
	case *ast.EnumDecl:
		styles = append(styles, colorEnumDecl(node)...)
	case *ast.CallOrIndexExpr:
		styles = append(styles, colorCallOrIndexExpr(node)...)
	case *ast.Comment:
		styles = append(styles, colorComment(node))
	case *ast.BasicLit:
		styles = append(styles, colorBasicLit(node)...)
	case *ast.DimDecl:
		styles = append(styles, colorDimDecl(node)...)
	case *ast.ScalarDecl:
		styles = append(styles, colorScalarDecl(node)...)
	case *ast.ArrayDecl:
		styles = append(styles, colorArrayDecl(node)...)
	case *ast.ArrayType:
		styles = append(styles, colorArrayType(node)...)
	case *ast.Identifier:
		styles = append(styles, colorIdentifier(node)...)
	case *ast.FuncDecl:
		styles = append(styles, colorFuncDecl(node)...)
	case *ast.SubDecl:
		styles = append(styles, colorSubDecl(node)...)
	case *ast.FuncType:
		styles = append(styles, colorFuncType(node)...)
	case *ast.SubType:
		styles = append(styles, colorSubType(node)...)
	case *ast.SpecialStmt:
		styles = append(styles, colorSpecialStmt(node)...)
	case *ast.File:
		styles = append(styles, colorFile(node)...)
	case *ast.StatementList:
		styles = append(styles, colorStatementList(node)...)
	case *ast.EmptyStmt:
		styles = append(styles, colorLineHeader())
	case *ast.JumpLabelDecl:
		styles = append(styles, colorJumpLabelDecl(node)...)
	case *ast.UserDefinedType:
		styles = append(styles, colorUserDefinedType(node)...)
	case *ast.ParamItem:
		styles = append(styles, colorParamItem(node)...)
	case *ast.ClassDecl:
		styles = append(styles, colorClassDecl(node)...)
	case *ast.JumpStmt:
		styles = append(styles, colorJumpStmt(node)...)

	default:
		// styles = append(styles, richtext.SpanStyle{Color: unknown, Content: node.String(), Size: size, Font: font})
		panic(fmt.Sprintf("unhandled node type %T", node))
	}
	return styles
}

func colorBody(node []ast.StatementList) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	for _, stmt := range node {
		styles = append(styles, colorText(&stmt)...)
	}
	return styles
}

// colorJumpStmt returns the color of a jump statement.
func colorJumpStmt(n *ast.JumpStmt) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	if n.OnError {
		styles = append(styles, richtext.SpanStyle{Color: statement, Content: "On error ", Size: size, Font: font})
	}
	styles = append(styles, richtext.SpanStyle{Color: statement, Content: n.JumpKw.Literal, Size: size, Font: font})
	styles = append(styles, richtext.SpanStyle{Color: statement, Content: " ", Size: size, Font: font})

	if n.NextKw != nil {
		styles = append(styles, richtext.SpanStyle{Color: statement, Content: "Next ", Size: size, Font: font})
	}
	if n.Label != nil {
		styles = append(styles, colorText(n.Label)...)
	}
	if n.Number != nil {
		styles = append(styles, richtext.SpanStyle{Color: identifier, Content: n.Number.Literal, Size: size, Font: font})
	}
	return styles
}

// colorClassDecl returns the color of a class declaration
func colorClassDecl(n *ast.ClassDecl) []richtext.SpanStyle {
	panic("not implemented")
}

// colorParamItem returns the color of a parameter item.
func colorParamItem(n *ast.ParamItem) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	if n.Optional {
		styles = append(styles, richtext.SpanStyle{Color: statement, Content: "Optional ", Size: size, Font: font})
	}
	if n.ByVal {
		styles = append(styles, richtext.SpanStyle{Color: statement, Content: "ByVal ", Size: size, Font: font})
	}
	if n.ParamArray {
		styles = append(styles, richtext.SpanStyle{Color: statement, Content: "ParamArray ", Size: size, Font: font})
	}
	styles = append(styles, colorText(n.VarName)...)
	if n.IsArray {
		styles = append(styles, richtext.SpanStyle{Color: operator, Content: "()", Size: size, Font: font})
	}

	if n.VarType != nil {
		styles = append(styles, richtext.SpanStyle{Color: other, Content: " As ", Size: size, Font: font})
		styles = append(styles, colorText(n.VarType)...)
	}
	if n.DefaultValue != nil {
		styles = append(styles, richtext.SpanStyle{Color: other, Content: " = ", Size: size, Font: font})
		styles = append(styles, colorText(n.DefaultValue)...)
	}
	return styles
}

// colorUserDefinedType returns the color of a user defined type.
func colorUserDefinedType(n *ast.UserDefinedType) []richtext.SpanStyle {
	return colorIdentifier(n.Identifier)
}

// colorJumpLabelDecl returns the color of a jump label declaration.
func colorJumpLabelDecl(n *ast.JumpLabelDecl) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	styles = append(styles, richtext.SpanStyle{Color: label, Content: n.Label.Name, Size: size, Font: font})
	styles = append(styles, richtext.SpanStyle{Color: statement, Content: ":", Size: size, Font: font})
	return styles
}

// colorCaseStmt returns the color of a case statement.
func colorCaseStmt(n *ast.CaseStmt) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	styles = append(styles, colorLineHeader())
	styles = append(styles, richtext.SpanStyle{Color: statement, Content: "Case ", Size: size, Font: font})
	if n.Condition != nil {
		styles = append(styles, colorText(n.Condition)...)
	} else {
		styles = append(styles, richtext.SpanStyle{Color: statement, Content: "Else", Size: size, Font: font})
	}
	styles = append(styles, colorEOL())
	identation += identSize
	for _, stmt := range n.Body {
		styles = append(styles, colorText(&stmt)...)
	}
	identation -= identSize
	return styles
}

// colorSelectStmt returns the color of a select statement.
func colorSelectStmt(n *ast.SelectStmt) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	styles = append(styles, richtext.SpanStyle{Color: statement, Content: "Select Case ", Size: size, Font: font})
	styles = append(styles, colorText(n.Condition)...)
	styles = append(styles, colorEOL())
	identation += identSize
	for _, stmt := range n.Body {
		styles = append(styles, colorText(&stmt)...)
	}
	identation -= identSize
	styles = append(styles, colorLineHeader())
	styles = append(styles, richtext.SpanStyle{Color: statement, Content: "End Select", Size: size, Font: font})
	return styles
}

// colorCallSubStmt returns the color of a call sub statement.
func colorCallSubStmt(n *ast.CallSubStmt) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	styles = append(styles, richtext.SpanStyle{Color: statement, Content: "call ", Size: size, Font: font})
	styles = append(styles, colorText(n.Definition)...)
	return styles
}

// colorCallSelectorExpr returns the color of a call selector expression.
func colorCallSelectorExpr(n *ast.CallSelectorExpr) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	styles = append(styles, colorText(n.Root)...)
	styles = append(styles, richtext.SpanStyle{Color: operator, Content: ".", Size: size, Font: font})
	styles = append(styles, colorText(n.Selector)...)
	return styles
}

// colorForStmt returns the color of a for statement.
func colorForStmt(n *ast.ForStmt) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	styles = append(styles, richtext.SpanStyle{Color: statement, Content: "For ", Size: size, Font: font})
	styles = append(styles, colorText(n.ForExpression)...)
	styles = append(styles, colorEOL())
	identation += identSize
	styles = append(styles, colorBody(n.Body)...)
	identation -= identSize
	styles = append(styles, colorLineHeader())
	styles = append(styles, richtext.SpanStyle{Color: statement, Content: " Next", Size: size, Font: font})
	return styles
}

// colorForNextExpr returns the color of a for next expression.
func colorForNextExpr(n *ast.ForNextExpr) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	styles = append(styles, colorText(n.Variable)...)
	styles = append(styles, richtext.SpanStyle{Color: statement, Content: " = ", Size: size, Font: font})
	styles = append(styles, colorText(n.From)...)
	styles = append(styles, richtext.SpanStyle{Color: statement, Content: " To ", Size: size, Font: font})
	styles = append(styles, colorText(n.To)...)
	if n.Step != nil {
		styles = append(styles, richtext.SpanStyle{Color: statement, Content: " Step ", Size: size, Font: font})
		styles = append(styles, colorText(n.Step)...)
	}
	return styles
}

// colorForEachExpr returns the color of a for each expression.
func colorForEachExpr(n *ast.ForEachExpr) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	styles = append(styles, richtext.SpanStyle{Color: statement, Content: "Each ", Size: size, Font: font})
	styles = append(styles, colorText(n.Variable)...)
	styles = append(styles, richtext.SpanStyle{Color: statement, Content: " In ", Size: size, Font: font})
	styles = append(styles, colorText(n.Collection)...)
	return styles
}

// colorWhileStmt returns the color of a while statement.
func colorWhileStmt(n *ast.WhileStmt) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	styles = append(styles, richtext.SpanStyle{Color: statement, Content: "Do While ", Size: size, Font: font})
	styles = append(styles, colorText(n.Condition)...)
	styles = append(styles, colorEOL())
	identation += identSize
	styles = append(styles, colorBody(n.Body)...)
	identation -= identSize
	styles = append(styles, colorLineHeader())
	styles = append(styles, richtext.SpanStyle{Color: statement, Content: "Loop", Size: size, Font: font})
	return styles
}

// colorUntilStmt returns the color of an until statement.
func colorUntilStmt(n *ast.UntilStmt) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	styles = append(styles, richtext.SpanStyle{Color: statement, Content: "Do Until ", Size: size, Font: font})
	styles = append(styles, colorText(n.Condition)...)
	styles = append(styles, colorEOL())
	identation += identSize
	styles = append(styles, colorBody(n.Body)...)
	identation -= identSize
	styles = append(styles, colorLineHeader())
	styles = append(styles, richtext.SpanStyle{Color: statement, Content: " Loop", Size: size, Font: font})
	return styles
}

// colorDoWhileStmt returns the color of a do while statement.
func colorDoWhileStmt(n *ast.DoWhileStmt) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	styles = append(styles, richtext.SpanStyle{Color: statement, Content: "Do ", Size: size, Font: font})
	styles = append(styles, colorEOL())
	identation += identSize
	styles = append(styles, colorBody(n.Body)...)
	identation -= identSize
	styles = append(styles, colorLineHeader())
	styles = append(styles, richtext.SpanStyle{Color: statement, Content: " Loop while ", Size: size, Font: font})
	styles = append(styles, colorText(n.Condition)...)
	return styles
}

// colorDoUntilStmt returns the color of a do until statement.
func colorDoUntilStmt(n *ast.DoUntilStmt) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	styles = append(styles, richtext.SpanStyle{Color: statement, Content: "Do", Size: size, Font: font})
	styles = append(styles, colorEOL())
	identation += identSize
	styles = append(styles, colorBody(n.Body)...)
	identation -= identSize
	styles = append(styles, colorLineHeader())
	styles = append(styles, richtext.SpanStyle{Color: statement, Content: " Loop Until ", Size: size, Font: font})
	styles = append(styles, colorText(n.Condition)...)
	return styles
}

// colorExitStmt returns the color of an exit statement.
func colorExitStmt(n *ast.ExitStmt) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	styles = append(styles, richtext.SpanStyle{Color: statement, Content: "Exit ", Size: size, Font: font})
	styles = append(styles, richtext.SpanStyle{Color: statement, Content: n.ExitType.Literal, Size: size, Font: font})
	return styles
}

// colorParenExpr returns the color of a parenthesized expression.
func colorParenExpr(n *ast.ParenExpr) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	styles = append(styles, richtext.SpanStyle{Color: operator, Content: "(", Size: size, Font: font})
	styles = append(styles, colorText(n.Expr)...)
	styles = append(styles, richtext.SpanStyle{Color: operator, Content: ")", Size: size, Font: font})
	return styles
}

// colorUnaryExpr returns the color of a unary expression.
func colorUnaryExpr(n *ast.UnaryExpr) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	styles = append(styles, richtext.SpanStyle{Color: operator, Content: n.OpKind.String(), Size: size, Font: font})
	styles = append(styles, colorText(n.Right)...)
	return styles
}

// colorBinaryExpr returns the color of a binary expression.
func colorBinaryExpr(n *ast.BinaryExpr) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	styles = append(styles, colorText(n.Left)...)
	styles = append(styles, richtext.SpanStyle{Color: operator, Content: " " + n.OpKind.String() + " ", Size: size, Font: font})
	styles = append(styles, colorText(n.Right)...)
	return styles
}

// colorExprStmt returns the color of an expression statement.
func colorExprStmt(n *ast.ExprStmt) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	styles = append(styles, richtext.SpanStyle{Color: keyword, Content: "Let ", Size: size, Font: font})
	styles = append(styles, colorText(n.Expression)...)
	return styles
}

// colorConstDecl returns the color of a constant declaration.
func colorConstDecl(n *ast.ConstDecl) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	styles = append(styles, richtext.SpanStyle{Color: keyword, Content: "Const ", Size: size, Font: font})
	for i, constDecl := range n.Consts {
		if i != 0 {
			styles = append(styles, richtext.SpanStyle{Color: operator, Content: ", ", Size: size, Font: font})
		}
		styles = append(styles, colorText(&constDecl)...)
	}
	return styles
}

// colorConstDeclItem returns the color of a constant declaration item.
func colorConstDeclItem(n *ast.ConstDeclItem) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	styles = append(styles, colorText(n.ConstName)...)
	if n.ConstType != nil {
		styles = append(styles, richtext.SpanStyle{Color: other, Content: " As ", Size: size, Font: font})
		styles = append(styles, colorText(n.ConstType)...)
	}
	styles = append(styles, richtext.SpanStyle{Color: other, Content: " = ", Size: size, Font: font})
	styles = append(styles, colorText(n.ConstValue)...)
	return styles
}

// colorElseIfStmt returns the color of an else if statement.
func colorElseIfStmt(n *ast.ElseIfStmt) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	styles = append(styles, colorLineHeader())
	styles = append(styles, richtext.SpanStyle{Color: keyword, Content: "ElseIf ", Size: size, Font: font})
	styles = append(styles, colorText(n.Condition)...)
	styles = append(styles, richtext.SpanStyle{Color: keyword, Content: " Then", Size: size, Font: font})
	styles = append(styles, colorEOL())
	identation += identSize
	styles = append(styles, colorBody(n.Body)...)
	identation -= identSize
	return styles
}

// colorIfStmt returns the color of an if statement.
func colorIfStmt(n *ast.IfStmt) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	styles = append(styles, richtext.SpanStyle{Color: statement, Content: "If ", Size: size, Font: font})
	styles = append(styles, colorText(n.Condition)...)
	styles = append(styles, richtext.SpanStyle{Color: statement, Content: " Then", Size: size, Font: font})
	styles = append(styles, colorEOL())
	identation += identSize
	styles = append(styles, colorBody(n.Body)...)
	identation -= identSize
	for _, stmt := range n.ElseIf {
		styles = append(styles, colorText(&stmt)...)
	}
	if n.Else != nil {
		styles = append(styles, colorLineHeader())
		styles = append(styles, richtext.SpanStyle{Color: statement, Content: "Else", Size: size, Font: font})
		styles = append(styles, colorEOL())
		identation += identSize
		styles = append(styles, colorBody(n.Else)...)
		identation -= identSize
	}
	styles = append(styles, colorLineHeader())
	styles = append(styles, richtext.SpanStyle{Color: statement, Content: "End If", Size: size, Font: font})
	return styles
}

// colorEnumDecl returns the color of an enum declaration.
func colorEnumDecl(n *ast.EnumDecl) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	styles = append(styles, richtext.SpanStyle{Color: keyword, Content: "Enum ", Size: size, Font: font})
	styles = append(styles, colorText(n.Identifier)...)
	styles = append(styles, colorEOL())
	identation += identSize
	for _, value := range n.Values {
		styles = append(styles, colorLineHeader())
		styles = append(styles, colorText(&value)...)
		styles = append(styles, colorEOL())
	}
	identation -= identSize
	styles = append(styles, colorLineHeader())
	styles = append(styles, richtext.SpanStyle{Color: keyword, Content: "End Enum", Size: size, Font: font})
	return styles
}

// colorCallOrIndexExpr returns the color of a call or index expression.
func colorCallOrIndexExpr(n *ast.CallOrIndexExpr) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	styles = append(styles, colorText(n.Identifier)...)
	styles = append(styles, richtext.SpanStyle{Color: operator, Content: "(", Size: size, Font: font})
	for i, arg := range n.Args {
		if i != 0 {
			styles = append(styles, richtext.SpanStyle{Color: operator, Content: ", ", Size: size, Font: font})
		}
		styles = append(styles, colorText(arg)...)
	}
	styles = append(styles, richtext.SpanStyle{Color: operator, Content: ")", Size: size, Font: font})
	return styles
}

// colorComment returns the color of a comment.
func colorComment(n *ast.Comment) richtext.SpanStyle {
	return richtext.SpanStyle{Color: comment, Content: n.Text.Literal, Size: size, Font: fontComment}
}

// colorBasicLit returns the color of a basic literal.
func colorBasicLit(n *ast.BasicLit) []richtext.SpanStyle {
	var color color.NRGBA
	switch n.Kind {
	case token.StringLit:
		color = stringCol
	case token.CurrencyLit, token.DateLit, token.DoubleLit, token.LongLit:
		color = constants
	case token.KwTrue, token.KwFalse:
		color = identifier
	}
	return []richtext.SpanStyle{{Color: color, Content: fmt.Sprint(n.Value), Size: size, Font: font}}
}

// colorStatementList returns the color of a statement list.
func colorStatementList(n *ast.StatementList) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	styles = append(styles, colorLineHeader())
	for i, stmt := range n.Statements {
		if i != 0 && stmt.Token().Kind != token.Comment {
			styles = append(styles, richtext.SpanStyle{Color: other, Content: " : ", Size: size, Font: font})
		} else if stmt.Token().Kind == token.Comment {
			styles = append(styles, richtext.SpanStyle{Color: other, Content: " ", Size: size, Font: font})
		}
		styles = append(styles, colorText(stmt)...)
	}
	styles = append(styles, colorEOL())
	return styles
}

func colorLineHeader() richtext.SpanStyle {
	lineNumber++
	if lineNumber == currentLineNumber {
		lineHeader := fmt.Sprintf("%s% 3d", ">>", lineNumber) + " " + strings.Repeat(" ", identation)
		return richtext.SpanStyle{Color: carret, Content: lineHeader, Size: size, Font: font}
	}
	lineHeader := fmt.Sprintf("% 6d", lineNumber) + " " + strings.Repeat(" ", identation)
	return richtext.SpanStyle{Color: lineNum, Content: lineHeader, Size: size, Font: font}
}

func colorEOL() richtext.SpanStyle {
	return richtext.SpanStyle{Color: keyword, Content: "\n", Size: size, Font: font}
}

// colorFile returns the color of a file.
func colorFile(n *ast.File) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	for _, node := range n.Body {
		styles = append(styles, colorText(&node)...)
	}
	return styles
}

// colorDimDecl returns the color of a dim declaration.
func colorDimDecl(n *ast.DimDecl) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	styles = append(styles, richtext.SpanStyle{Color: keyword, Content: "Dim ", Size: size, Font: font})
	for i, varDecl := range n.Vars {
		if i != 0 {
			styles = append(styles, richtext.SpanStyle{Color: operator, Content: ", ", Size: size, Font: font})
		}
		styles = append(styles, colorText(varDecl)...)
	}
	return styles
}

// colorScalarDecl returns the color of a scalar declaration.
func colorScalarDecl(n *ast.ScalarDecl) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	styles = append(styles, colorText(n.VarName)...)
	if n.VarType != nil {
		styles = append(styles, richtext.SpanStyle{Color: other, Content: " As ", Size: size, Font: font})
		styles = append(styles, colorText(n.VarType)...)
	}
	if n.VarValue != nil {
		styles = append(styles, richtext.SpanStyle{Color: other, Content: " = ", Size: size, Font: font})
		styles = append(styles, colorText(n.VarValue)...)
	}
	return styles
}

// colorArrayDecl returns the color of an array declaration.
func colorArrayDecl(n *ast.ArrayDecl) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	styles = append(styles, colorText(n.VarName)...)
	styles = append(styles, colorText(n.VarType)...)
	return styles
}

// colorArrayType returns the color of an array type.
func colorArrayType(n *ast.ArrayType) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	styles = append(styles, richtext.SpanStyle{Color: operator, Content: "(", Size: size, Font: font})
	for i, dim := range n.Dimensions {
		if i != 0 {
			styles = append(styles, richtext.SpanStyle{Color: operator, Content: ", ", Size: size, Font: font})
		}
		styles = append(styles, colorText(dim)...)
	}
	styles = append(styles, richtext.SpanStyle{Color: operator, Content: ")", Size: size, Font: font})
	if n.Type != nil {
		styles = append(styles, richtext.SpanStyle{Color: other, Content: " As ", Size: size, Font: font})
		styles = append(styles, colorText(n.Type)...)
	}
	return styles
}

// colorIdentifier returns the color of an identifier.
func colorIdentifier(n *ast.Identifier) []richtext.SpanStyle {
	var color color.NRGBA
	switch strings.ToLower(n.Name) {
	case "long", "integer", "single", "double", "currency", "date", "string", "boolean", "variant":
		color = typ
	default:
		color = identifier
	}
	return []richtext.SpanStyle{{Color: color, Content: n.Name, Size: size, Font: font}}
}

// colorFuncDecl returns the color of a function declaration.
func colorFuncDecl(n *ast.FuncDecl) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	styles = append(styles, richtext.SpanStyle{Color: keyword, Content: "Function ", Size: size, Font: font})
	styles = append(styles, colorText(n.FuncName)...)
	styles = append(styles, richtext.SpanStyle{Color: operator, Content: "(", Size: size, Font: font})
	for i, param := range n.FuncType.Params {
		if i != 0 {
			styles = append(styles, richtext.SpanStyle{Color: operator, Content: ", ", Size: size, Font: font})
		}
		styles = append(styles, colorText(&param)...)
	}
	styles = append(styles, richtext.SpanStyle{Color: operator, Content: ") As ", Size: size, Font: font})
	styles = append(styles, colorText(n.FuncType.Result)...)
	styles = append(styles, richtext.SpanStyle{Color: keyword, Content: "\n", Size: size, Font: font})
	identation += identSize

	for _, stmt := range n.Body {
		styles = append(styles, colorText(&stmt)...)
	}
	identation -= identSize
	styles = append(styles, colorLineHeader())
	styles = append(styles, richtext.SpanStyle{Color: keyword, Content: "End Function", Size: size, Font: font})
	return styles
}

// colorSubDecl returns the color of a sub declaration.
func colorSubDecl(n *ast.SubDecl) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	styles = append(styles, richtext.SpanStyle{Color: keyword, Content: "Sub ", Size: size, Font: font})
	styles = append(styles, colorText(n.SubName)...)
	styles = append(styles, richtext.SpanStyle{Color: operator, Content: "(", Size: size, Font: font})
	for i, param := range n.SubType.Params {
		if i != 0 {
			styles = append(styles, richtext.SpanStyle{Color: operator, Content: ", ", Size: size, Font: font})
		}
		styles = append(styles, colorText(&param)...)
	}
	styles = append(styles, richtext.SpanStyle{Color: operator, Content: ")\n", Size: size, Font: font})
	identation += identSize
	for _, stmt := range n.Body {
		styles = append(styles, colorText(&stmt)...)
	}
	identation -= identSize
	styles = append(styles, colorLineHeader())
	styles = append(styles, richtext.SpanStyle{Color: keyword, Content: "End Sub", Size: size, Font: font})
	return styles
}

// colorFuncType returns the color of a function type.
func colorFuncType(n *ast.FuncType) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	styles = append(styles, richtext.SpanStyle{Color: operator, Content: "(", Size: size, Font: font})
	for i, param := range n.Params {
		if i != 0 {
			styles = append(styles, richtext.SpanStyle{Color: operator, Content: ", ", Size: size, Font: font})
		}
		styles = append(styles, colorText(&param)...)
	}
	styles = append(styles, richtext.SpanStyle{Color: operator, Content: ")", Size: size, Font: font})
	styles = append(styles, colorText(n.Result)...)
	return styles
}

// colorSubType returns the color of a sub type.
func colorSubType(n *ast.SubType) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	styles = append(styles, richtext.SpanStyle{Color: operator, Content: "(", Size: size, Font: font})
	for i, param := range n.Params {
		if i != 0 {
			styles = append(styles, richtext.SpanStyle{Color: operator, Content: ", ", Size: size, Font: font})
		}
		styles = append(styles, colorText(&param)...)
	}
	styles = append(styles, richtext.SpanStyle{Color: operator, Content: ")", Size: size, Font: font})
	return styles
}

// colorSpecialStmt returns the color of a special statement.
func colorSpecialStmt(n *ast.SpecialStmt) []richtext.SpanStyle {
	var styles []richtext.SpanStyle
	styles = append(styles, richtext.SpanStyle{Color: keyword, Content: n.Keyword1.Literal + " ", Size: size, Font: font})
	if len(n.Keyword2) > 0 {
		styles = append(styles, richtext.SpanStyle{Color: keyword2, Content: n.Keyword2 + " ", Size: size, Font: font})
	}
	for i, arg := range n.Args {
		if i != 0 {
			styles = append(styles, richtext.SpanStyle{Color: keyword2, Content: ", ", Size: size, Font: font})
		}
		styles = append(styles, colorText(arg)...)
	}
	if n.Semicolon != nil {
		styles = append(styles, richtext.SpanStyle{Color: keyword2, Content: "; ", Size: size, Font: font})
	}
	return styles
}
