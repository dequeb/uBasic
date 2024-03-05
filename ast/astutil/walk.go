package astutil

import (
	"fmt"

	"uBasic/ast"
)

// Walk traverses the given parse tree, calling f(n) for each node n in the
// tree, in a bottom-up traversal.
func Walk(node ast.Node, f func(ast.Node) error) error {
	nop := func(n ast.Node) error { return nil }
	return WalkBeforeAfter(node, nop, f)
}

// WalkBeforeAfter traverses the given parse tree, calling before(n) before
// traversing the node's children, and after(n) afterwards, in a bottom-up
// traversal.
func WalkBeforeAfter(node ast.Node, before, after func(ast.Node) error) error {

	// if nill, do nothing
	if node == nil {
		return nil
	}

	switch n := node.(type) {
	// Source file.
	case *ast.File:
		if n != nil {
			return walkFile(n, before, after)
		}

	// Declarations.
	case *ast.FuncDecl:
		if n != nil {
			return walkFuncDecl(n, before, after)
		}
	case *ast.SubDecl:
		if n != nil {
			return walkSubDecl(n, before, after)
		}
	case *ast.ParamItem:
		if n != nil {
			return walkVarDecl(n, before, after)
		}

	case *ast.ArrayDecl:
		if n != nil {
			return walkArrayDecl(n, before, after)
		}
	case *ast.ArrayType:
		if n != nil {
			return walkArrayType(n, before, after)
		}
	case *ast.ScalarDecl:
		if n != nil {
			return walkScalarDecl(n, before, after)
		}
	case *ast.DimDecl:
		if n != nil {
			return walkDimDecl(n, before, after)
		}
	case *ast.EnumDecl:
		if n != nil {
			return walkEnumDecl(n, before, after)
		}
	case *ast.ConstDecl:
		if n != nil {
			return walkConstDecl(n, before, after)
		}
	case *ast.ConstDeclItem:
		if n != nil {
			return walkConstItemDecl(n, before, after)
		}
	case *ast.JumpLabelDecl:
		if n != nil {
			return walkJumpLabelDecl(n, before, after)
		}

	// Statements.
	case *ast.JumpStmt:
		if n != nil {
			return walkJumpStmt(n, before, after)
		}
	case *ast.EmptyStmt:
		if n != nil {
			return walkEmptyStmt(n, before, after)
		}
	case *ast.ExprStmt:
		if n != nil {
			return walkExprStmt(n, before, after)
		}
	case *ast.IfStmt:
		if n != nil {
			return walkIfStmt(n, before, after)
		}
	case *ast.WhileStmt:
		if n != nil {
			return walkWhileStmt(n, before, after)
		}
	case *ast.UntilStmt:
		if n != nil {
			return walkUntilStmt(n, before, after)
		}
	case *ast.DoWhileStmt:
		if n != nil {
			return walkWDohileStmt(n, before, after)
		}
	case *ast.DoUntilStmt:
		if n != nil {
			return walkDoUntilStmt(n, before, after)
		}
	case *ast.ForStmt:
		if n != nil {
			return walkForStmt(n, before, after)
		}
	case *ast.SpecialStmt:
		if n != nil {
			return walkSpecialStmt(n, before, after)
		}
	case *ast.ExitStmt:
		if n != nil {
			return walkExitStmt(n, before, after)
		}
	case *ast.StatementList:
		if n != nil {
			return WalkStatementList(n, before, after)
		}
	case *ast.SelectStmt:
		if n != nil {
			return walkSelectStmt(n, before, after)
		}
	case *ast.CaseStmt:
		if n != nil {
			return walkCaseExpr(n, before, after)
		}
	case *ast.CallSubStmt:
		if n != nil {
			return walkCallSubStmt(n, before, after)
		}

	// Expressions.
	case *ast.ForEachExpr:
		if n != nil {
			return walkForEachExpr(n, before, after)
		}
	case *ast.ForNextExpr:
		if n != nil {
			return walkForNextExpr(n, before, after)
		}
	case *ast.BasicLit:
		if n != nil {
			return walkBasicLit(n, before, after)
		}
	case *ast.BinaryExpr:
		if n != nil {
			return walkBinaryExpr(n, before, after)
		}
	case *ast.CallOrIndexExpr:
		if n != nil {
			return walkCallOrIndexExpr(n, before, after)
		}
	case *ast.Identifier:
		if n != nil {
			return walkIdent(n, before, after)
		}
	case *ast.ParenExpr:
		if n != nil {
			return walkParenExpr(n, before, after)
		}
	case *ast.UnaryExpr:
		if n != nil {
			return walkUnaryExpr(n, before, after)
		}
	case *ast.CallSelectorExpr:
		if n != nil {
			return walkCallSelectorExpr(n, before, after)
		}
	// Types.
	case *ast.FuncType:
		if n != nil {
			return walkFuncType(n, before, after)
		}
	case *ast.SubType:
		if n != nil {
			return walkSubType(n, before, after)
		}
	case *ast.UserDefinedType:
		if n != nil {
			return walkUserDefinedType(n, before, after)
		}
	case *ast.Comment:
		if n != nil {
			return walkComment(n, before, after)
		}

	default:
		panic(fmt.Sprintf("support for walking node of type %T not yet implemented", node))
	}

	return nil
}

// === [ Source file ] ===

// walkFile walks the parse tree of the given source file in depth first order.
func walkFile(file *ast.File, before, after func(ast.Node) error) error {
	if err := before(file); err != nil {
		return err
	}
	for _, stmtList := range file.StatementLists {
		if err := WalkBeforeAfter(&stmtList, before, after); err != nil {
			return err
		}
	}
	if err := after(file); err != nil {
		return err
	}
	return nil
}

// === [ Top-level declarations ] ===

// walkFuncDecl walks the parse tree of the given function declaration in depth
// first order.
func walkFuncDecl(decl *ast.FuncDecl, before, after func(ast.Node) error) error {
	if err := before(decl); err != nil {
		return err
	}
	if err := WalkBeforeAfter(decl.FuncName, before, after); err != nil {
		return err
	}
	if err := WalkBeforeAfter(decl.FuncType, before, after); err != nil {
		return err
	}
	for _, stmtList := range decl.Body {
		if err := WalkBeforeAfter(&stmtList, before, after); err != nil {
			return err
		}
	}
	if err := after(decl); err != nil {
		return err
	}
	return nil
}

// walkSubDecl walks the parse tree of the given sub declaration in depth first
// order.
func walkSubDecl(decl *ast.SubDecl, before, after func(ast.Node) error) error {
	if err := before(decl); err != nil {
		return err
	}
	if err := WalkBeforeAfter(decl.SubName, before, after); err != nil {
		return err
	}
	if err := WalkBeforeAfter(decl.SubType, before, after); err != nil {
		return err
	}
	for _, stmtList := range decl.Body {
		if err := WalkBeforeAfter(&stmtList, before, after); err != nil {
			return err
		}
	}
	if err := after(decl); err != nil {
		return err
	}
	return nil
}

// walkVarDecl walks the parse tree of the given variable declaration in depth
// first order.
func walkVarDecl(decl *ast.ParamItem, before, after func(ast.Node) error) error {
	if err := before(decl); err != nil {
		return err
	}
	if err := WalkBeforeAfter(decl.VarType, before, after); err != nil {
		return err
	}
	if err := WalkBeforeAfter(decl.VarName, before, after); err != nil {
		return err
	}
	if err := WalkBeforeAfter(decl.DefaultValue, before, after); err != nil {
		return err
	}
	if err := after(decl); err != nil {
		return err
	}
	return nil
}

// walkArrayDecl walks the parse tree of the given array declaration in depth
// first order.
func walkArrayDecl(decl *ast.ArrayDecl, before, after func(ast.Node) error) error {
	if err := before(decl); err != nil {
		return err
	}
	if err := WalkBeforeAfter(decl.VarType, before, after); err != nil {
		return err
	}
	if err := WalkBeforeAfter(decl.VarName, before, after); err != nil {
		return err
	}
	if err := after(decl); err != nil {
		return err
	}
	return nil
}

func walkArrayType(decl *ast.ArrayType, before, after func(ast.Node) error) error {
	if err := before(decl); err != nil {
		return err
	}
	for _, dim := range decl.Dimensions {
		if err := WalkBeforeAfter(dim, before, after); err != nil {
			return err
		}
	}
	if err := WalkBeforeAfter(decl.Type, before, after); err != nil {
		return err
	}
	if err := after(decl); err != nil {
		return err
	}
	return nil
}

// walkScalarDecl walks the parse tree of the given scalar declaration in depth
// first order.
func walkScalarDecl(decl *ast.ScalarDecl, before, after func(ast.Node) error) error {
	if err := before(decl); err != nil {
		return err
	}
	if err := WalkBeforeAfter(decl.VarType, before, after); err != nil {
		return err
	}
	if err := WalkBeforeAfter(decl.VarName, before, after); err != nil {
		return err
	}
	if err := WalkBeforeAfter(decl.VarValue, before, after); err != nil {
		return err
	}
	if err := after(decl); err != nil {
		return err
	}
	return nil
}

// walkDimDecl walks the parse tree of the given dim declaration in depth first
// order.
func walkDimDecl(decl *ast.DimDecl, before, after func(ast.Node) error) error {
	if err := before(decl); err != nil {
		return err
	}
	for _, variable := range decl.Vars {
		if err := WalkBeforeAfter(variable, before, after); err != nil {
			return err
		}
	}
	if err := after(decl); err != nil {
		return err
	}
	return nil
}

// walkEnumDecl walks the parse tree of the given enum declaration in depth first
// order.
func walkEnumDecl(decl *ast.EnumDecl, before, after func(ast.Node) error) error {
	if err := before(decl); err != nil {
		return err
	}
	if err := WalkBeforeAfter(decl.Identifier, before, after); err != nil {
		return err
	}
	for _, val := range decl.Values {
		if err := WalkBeforeAfter(&val, before, after); err != nil {
			return err
		}
	}
	if err := after(decl); err != nil {
		return err
	}
	return nil
}

// walkConstDecl walks the parse tree of the given constant declaration in depth
// first order.
func walkConstDecl(decl *ast.ConstDecl, before, after func(ast.Node) error) error {
	if err := before(decl); err != nil {
		return err
	}
	for _, constant := range decl.Consts {
		if err := WalkBeforeAfter(&constant, before, after); err != nil {
			return err
		}
	}
	if err := after(decl); err != nil {
		return err
	}
	return nil
}

// walkConstItemDecl walks the parse tree of the given constant declaration item in depth
// first order.
func walkConstItemDecl(decl *ast.ConstDeclItem, before, after func(ast.Node) error) error {
	if err := before(decl); err != nil {
		return err
	}
	if err := WalkBeforeAfter(decl.ConstName, before, after); err != nil {
		return err
	}
	if err := WalkBeforeAfter(decl.ConstType, before, after); err != nil {
		return err
	}
	if err := WalkBeforeAfter(decl.ConstValue, before, after); err != nil {
		return err
	}
	if err := after(decl); err != nil {
		return err
	}
	return nil
}

// === [ Statements ] ===

// walkEmptyStmt walks the parse tree of the given empty statement in depth
// first order.
func walkEmptyStmt(stmt *ast.EmptyStmt, before, after func(ast.Node) error) error {
	if err := before(stmt); err != nil {
		return err
	}
	if err := after(stmt); err != nil {
		return err
	}
	return nil
}

// walkExprStmt walks the parse tree of the given expression statement in depth
// first order.
func walkExprStmt(stmt *ast.ExprStmt, before, after func(ast.Node) error) error {
	if err := before(stmt); err != nil {
		return err
	}
	if err := WalkBeforeAfter(stmt.Expression, before, after); err != nil {
		return err
	}
	if err := after(stmt); err != nil {
		return err
	}
	return nil
}

// walkIfStmt walks the parse tree of the given if statement in depth first
// order.
func walkIfStmt(stmt *ast.IfStmt, before, after func(ast.Node) error) error {
	if err := before(stmt); err != nil {
		return err
	}
	if err := WalkBeforeAfter(stmt.Condition, before, after); err != nil {
		return err
	}
	for _, stmtList := range stmt.Body {
		if err := WalkBeforeAfter(&stmtList, before, after); err != nil {
			return err
		}
	}
	for _, stmtList := range stmt.Else {
		if err := WalkBeforeAfter(&stmtList, before, after); err != nil {
			return err
		}
	}
	if err := after(stmt); err != nil {
		return err
	}
	return nil
}

// walkWhileStmt walks the parse tree of the given while statement in depth
// first order.
func walkWhileStmt(stmt *ast.WhileStmt, before, after func(ast.Node) error) error {
	if err := before(stmt); err != nil {
		return err
	}
	if err := WalkBeforeAfter(stmt.Condition, before, after); err != nil {
		return err
	}
	for _, stmtList := range stmt.Body {
		if err := WalkBeforeAfter(&stmtList, before, after); err != nil {
			return err
		}
	}
	if err := after(stmt); err != nil {
		return err
	}
	return nil
}

// walkUntilStmt walks the parse tree of the given until statement in depth
// first order.
func walkUntilStmt(stmt *ast.UntilStmt, before, after func(ast.Node) error) error {
	if err := before(stmt); err != nil {
		return err
	}
	if err := WalkBeforeAfter(stmt.Condition, before, after); err != nil {
		return err
	}
	for _, stmtList := range stmt.Body {
		if err := WalkBeforeAfter(&stmtList, before, after); err != nil {
			return err
		}
	}
	if err := after(stmt); err != nil {
		return err
	}
	return nil
}

// walkWDohileStmt walks the parse tree of the given do-while statement in depth
// first order.
func walkWDohileStmt(stmt *ast.DoWhileStmt, before, after func(ast.Node) error) error {
	if err := before(stmt); err != nil {
		return err
	}
	for _, stmtList := range stmt.Body {
		if err := WalkBeforeAfter(&stmtList, before, after); err != nil {
			return err
		}
	}
	if err := WalkBeforeAfter(stmt.Condition, before, after); err != nil {
		return err
	}
	if err := after(stmt); err != nil {
		return err
	}
	return nil
}

// walkDoUntilStmt walks the parse tree of the given do-until statement in depth
// first order.
func walkDoUntilStmt(stmt *ast.DoUntilStmt, before, after func(ast.Node) error) error {
	if err := before(stmt); err != nil {
		return err
	}
	for _, stmtList := range stmt.Body {
		if err := WalkBeforeAfter(&stmtList, before, after); err != nil {
			return err
		}
	}
	if err := WalkBeforeAfter(stmt.Condition, before, after); err != nil {
		return err
	}
	if err := after(stmt); err != nil {
		return err
	}
	return nil
}

// walkForStmt walks the parse tree of the given for statement in depth first
// order.
func walkForStmt(stmt *ast.ForStmt, before, after func(ast.Node) error) error {
	if err := before(stmt); err != nil {
		return err
	}
	if err := WalkBeforeAfter(stmt.ForExpression, before, after); err != nil {
		return err
	}
	for _, stmtList := range stmt.Body {
		if err := WalkBeforeAfter(&stmtList, before, after); err != nil {
			return err
		}
	}
	if err := after(stmt); err != nil {
		return err
	}
	return nil
}

// walkSpecialStmt walks the parse tree of the given special statement in depth
// first order.
func walkSpecialStmt(stmt *ast.SpecialStmt, before, after func(ast.Node) error) error {
	if err := before(stmt); err != nil {
		return err
	}
	for _, arg := range stmt.Args {
		if err := WalkBeforeAfter(arg, before, after); err != nil {
			return err
		}
	}
	if err := after(stmt); err != nil {
		return err
	}
	return nil
}

// walkExitStmt walks the parse tree of the given exit statement in depth first
// order.
func walkExitStmt(stmt *ast.ExitStmt, before, after func(ast.Node) error) error {
	if err := before(stmt); err != nil {
		return err
	}
	if err := after(stmt); err != nil {
		return err
	}
	return nil
}

// WalkStatementList walks the parse tree of the given Stateement List in depth first
func WalkStatementList(stmt *ast.StatementList, before, after func(ast.Node) error) error {
	if err := before(stmt); err != nil {
		return err
	}
	for _, stmt := range stmt.Statements {
		if err := WalkBeforeAfter(stmt, before, after); err != nil {
			return err
		}
	}
	if err := after(stmt); err != nil {
		return err
	}
	return nil
}

// walkSelectStmt walks the parse tree of the given select statement in depth
// first order.
func walkSelectStmt(stmt *ast.SelectStmt, before, after func(ast.Node) error) error {
	if err := before(stmt); err != nil {
		return err
	}
	for _, stmtList := range stmt.Body {
		WalkBeforeAfter(&stmtList, before, after)
	}
	err := WalkBeforeAfter(stmt.Condition, before, after)
	if err != nil {
		return err
	}
	if err := after(stmt); err != nil {
		return err
	}
	return nil
}

// walkCaseExpr walks the parse tree of the given case expression in depth
// first order.
func walkCaseExpr(expr *ast.CaseStmt, before, after func(ast.Node) error) error {
	if err := before(expr); err != nil {
		return err
	}
	for _, stmtList := range expr.Body {
		for _, stmt := range stmtList.Statements {
			if err := WalkBeforeAfter(stmt, before, after); err != nil {
				return err
			}
		}
	}
	if err := WalkBeforeAfter(expr.Condition, before, after); err != nil {
		return err
	}
	if err := after(expr); err != nil {
		return err
	}
	return nil
}

// walkCallSubStmt walks the parse tree of the given case expression in depth
// first order.
func walkCallSubStmt(stmt *ast.CallSubStmt, before, after func(ast.Node) error) error {
	if err := before(stmt); err != nil {
		return err
	}
	if err := WalkBeforeAfter(stmt.Definition, before, after); err != nil {
		return err
	}
	if err := after(stmt); err != nil {
		return err
	}
	return nil
}

// === [ Expressions ] ===

// walkBasicLit walks the parse tree of the given basic literal expression in
// depth first order.
func walkBasicLit(lit *ast.BasicLit, before, after func(ast.Node) error) error {
	if err := before(lit); err != nil {
		return err
	}
	if err := after(lit); err != nil {
		return err
	}
	return nil
}

// walkBinaryExpr walks the parse tree of the given binary expression in depth
// first order.
func walkBinaryExpr(expr *ast.BinaryExpr, before, after func(ast.Node) error) error {
	if err := before(expr); err != nil {
		return err
	}
	if err := WalkBeforeAfter(expr.Left, before, after); err != nil {
		return err
	}
	if err := WalkBeforeAfter(expr.Right, before, after); err != nil {
		return err
	}
	if err := after(expr); err != nil {
		return err
	}
	return nil
}

// walkCallExpr walks the parse tree of the given call expression in depth first
// order.
func walkCallOrIndexExpr(call *ast.CallOrIndexExpr, before, after func(ast.Node) error) error {
	if err := before(call); err != nil {
		return err
	}
	if err := WalkBeforeAfter(call.Identifier, before, after); err != nil {
		return err
	}
	for _, arg := range call.Args {
		if err := WalkBeforeAfter(arg, before, after); err != nil {
			return err
		}
	}
	if err := after(call); err != nil {
		return err
	}
	return nil
}

// walkForEachExpr walks the parse tree of the given for-each expression in depth
// first order.
func walkForEachExpr(expr *ast.ForEachExpr, before, after func(ast.Node) error) error {
	if err := before(expr); err != nil {
		return err
	}
	if err := WalkBeforeAfter(expr.Variable, before, after); err != nil {
		return err
	}
	if err := WalkBeforeAfter(expr.Collection, before, after); err != nil {
		return err
	}
	if err := after(expr); err != nil {
		return err
	}
	return nil
}

// walkForNextExpr walks the parse tree of the given for-next expression in depth
// first order.
func walkForNextExpr(expr *ast.ForNextExpr, before, after func(ast.Node) error) error {
	if err := before(expr); err != nil {
		return err
	}
	if err := WalkBeforeAfter(expr.Variable, before, after); err != nil {
		return err
	}
	if err := WalkBeforeAfter(expr.From, before, after); err != nil {
		return err
	}
	if err := WalkBeforeAfter(expr.To, before, after); err != nil {
		return err
	}
	if err := WalkBeforeAfter(expr.Step, before, after); err != nil {
		return err
	}

	if err := after(expr); err != nil {
		return err
	}
	return nil
}

// walkIdent walks the parse tree of the given identifier expression in depth
// first order.
func walkIdent(ident *ast.Identifier, before, after func(ast.Node) error) error {
	if err := before(ident); err != nil {
		return err
	}
	if err := after(ident); err != nil {
		return err
	}
	return nil
}

// walkParenExpr walks the parse tree of the given parenthesized expression in
// depth first order.
func walkParenExpr(expr *ast.ParenExpr, before, after func(ast.Node) error) error {
	if err := before(expr); err != nil {
		return err
	}
	if err := WalkBeforeAfter(expr.Expr, before, after); err != nil {
		return err
	}
	if err := after(expr); err != nil {
		return err
	}
	return nil
}

// walkUnaryExpr walks the parse tree of the given unary expression in depth
// first order.
func walkUnaryExpr(expr *ast.UnaryExpr, before, after func(ast.Node) error) error {
	if err := before(expr); err != nil {
		return err
	}
	if err := WalkBeforeAfter(expr.Right, before, after); err != nil {
		return err
	}
	if err := after(expr); err != nil {
		return err
	}
	return nil
}

// walkCallSelectorExpr walks the parse tree of the given call selector expression
// in depth first order.
func walkCallSelectorExpr(expr *ast.CallSelectorExpr, before, after func(ast.Node) error) error {
	if err := before(expr); err != nil {
		return err
	}
	if err := WalkBeforeAfter(expr.Root, before, after); err != nil {
		return err
	}
	if err := WalkBeforeAfter(expr.Selector, before, after); err != nil {
		return err
	}
	if err := after(expr); err != nil {
		return err
	}
	return nil
}

// === [ Types ] ===

// walkFuncType walks the parse tree of the given function signature in depth
// first order.
func walkFuncType(fn *ast.FuncType, before, after func(ast.Node) error) error {
	if err := before(fn); err != nil {
		return err
	}
	if err := WalkBeforeAfter(fn.Result, before, after); err != nil {
		return err
	}
	for i := range fn.Params {
		if err := WalkBeforeAfter(&fn.Params[i], before, after); err != nil {
			return err
		}
	}
	if err := after(fn); err != nil {
		return err
	}
	return nil
}

// walkSubType walks the parse tree of the given subroutine signature in depth
// first order.
func walkSubType(fn *ast.SubType, before, after func(ast.Node) error) error {
	if err := before(fn); err != nil {
		return err
	}
	for i := range fn.Params {
		if err := WalkBeforeAfter(&fn.Params[i], before, after); err != nil {
			return err
		}
	}
	if err := after(fn); err != nil {
		return err
	}
	return nil
}

// walkUserDefinedType walks the parse tree of the given user-defined type in
// depth first order.
func walkUserDefinedType(typ *ast.UserDefinedType, before, after func(ast.Node) error) error {
	if err := before(typ); err != nil {
		return err
	}
	if err := WalkBeforeAfter(typ.Identifier, before, after); err != nil {
		return err
	}
	if err := after(typ); err != nil {
		return err
	}
	return nil
}

// === [ Jumps / Error Handling ] ===
// walkJumpStmt walks the parse tree of the given jump statement in depth first
// order.
func walkJumpLabelDecl(typ *ast.JumpLabelDecl, before, after func(ast.Node) error) error {
	if err := before(typ); err != nil {
		return err
	}
	if err := WalkBeforeAfter(typ.Label, before, after); err != nil {
		return err
	}
	if err := after(typ); err != nil {
		return err
	}
	return nil
}

// walkJumpStmt walks the parse tree of the given jump statement in depth first
// order.
func walkJumpStmt(stmt *ast.JumpStmt, before, after func(ast.Node) error) error {
	if err := before(stmt); err != nil {
		return err
	}
	if err := WalkBeforeAfter(stmt.Label, before, after); err != nil {
		return err
	}
	if err := after(stmt); err != nil {
		return err
	}
	return nil
}

// walkComment walks the parse tree of the given comment in depth first
// order.
func walkComment(stmt *ast.Comment, before, after func(ast.Node) error) error {
	if err := before(stmt); err != nil {
		return err
	}
	if err := after(stmt); err != nil {
		return err
	}
	return nil
}
