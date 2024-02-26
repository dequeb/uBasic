package sem

import (
	"uBasic/ast"
	"uBasic/ast/astutil"
	"uBasic/errors"
)

// setParent initialize the parent node of the given node.
func SetParents(file *ast.File) error {
	if err := setParent(file, nil); err != nil {
		return err
	}

	// check if all the nodes were given a parent node.
	return checkParents(file)

}

// check if all the nodes were gvin a parent node.
func checkParents(file *ast.File) error {
	// check if all the nodes were given a parent node.
	f := func(node ast.Node) error {
		// if node has no parent, return error
		if node.GetParent() == nil {
			return errors.Newf(node.Token().Position, "parent node is unknown for %v", node)
		}
		// if node has a parent, return nil
		return nil
	}
	for _, stmtList := range file.StatementLists {
		if err := astutil.Walk(&stmtList, f); err != nil {
			return err
		}
	}
	return nil
}

// Set parent node recursively.
func setParent(node ast.Node, parent ast.Node) error {
	if node == nil {
		return errors.New(UniversePos, "invalid empty node")
	}

	// set parent node
	switch n := (node).(type) {
	case *ast.StatementList:
		n.SetParent(parent)
		for i := range n.Statements {
			setParent(n.Statements[i], n) // enure to pass original object, not copy
		}
	case *ast.SelectStmt:
		n.SetParent(parent)
		setParent(n.Condition, n)
		for i := range n.Body {
			setParent(&n.Body[i], n)
		}
	case *ast.CaseStmt:
		n.SetParent(parent)
		setParent(n.Condition, n)
		for i := range n.Body {
			setParent(&n.Body[i], n)
		}
	case *ast.ForStmt:
		n.SetParent(parent)
		setParent(n.ForExpression, n)
		for i := range n.Body {
			setParent(&n.Body[i], n)
		}
	case *ast.ArrayDecl:
		n.SetParent(parent)
		setParent(n.VarName, n)
		setParent(n.VarType, n)
	case *ast.ArrayType:
		n.SetParent(parent)
		for _, dim := range n.Dimensions {
			setParent(dim, n)
		}
		setParent(n.Type, n)
	case *ast.BasicLit:
		n.SetParent(parent)
	case *ast.BinaryExpr:
		n.SetParent(parent)
		setParent(n.Left, n)
		setParent(n.Right, n)
	case *ast.ParenExpr:
		n.SetParent(parent)
		setParent(n.Expr, n)
	case *ast.UnaryExpr:
		n.SetParent(parent)
		setParent(n.Right, n)
	case *ast.ExitStmt:
		n.SetParent(parent)
	case *ast.SpecialStmt:
		n.SetParent(parent)
		for _, arg := range n.Args {
			setParent(arg, n)
		}
	case *ast.CallOrIndexExpr:
		n.SetParent(parent)
		setParent(n.Identifier, n)
		for _, arg := range n.Args {
			setParent(arg, n)
		}
	case *ast.CallSelectorExpr:
		n.SetParent(parent)
		setParent(n.Root, n)
		setParent(n.Selector, n)
	case *ast.EmptyStmt:
		n.SetParent(parent)
	case *ast.ExprStmt:
		n.SetParent(parent)
		setParent(n.Expression, n)
	case *ast.File:
		n.SetParent(parent)
		for i := range n.StatementLists {
			setParent(&n.StatementLists[i], n)
		}
	case *ast.FuncDecl:
		n.SetParent(parent)
		setParent(n.FuncName, n)
		setParent(n.FuncType, n)
		for i := range n.Body {
			setParent(&n.Body[i], n)
		}
		for _, param := range n.GetParams() {
			setParent(&param, n)
		}
	case *ast.SubDecl:
		n.SetParent(parent)
		setParent(n.SubName, n)
		setParent(n.SubType, n)
		for i := range n.Body {
			setParent(&n.Body[i], n)
		}
		for _, param := range n.GetParams() {
			setParent(&param, n)
		}
	case *ast.FuncType:
		n.SetParent(parent)
		for i := range n.Params {
			setParent(&n.Params[i], n)
		}
		setParent(n.Result, n)
	case *ast.SubType:
		n.SetParent(parent)
		for i := range n.Params {
			setParent(&n.Params[i], n)
		}
	case *ast.Identifier:
		n.SetParent(parent)
	case *ast.IfStmt:
		n.SetParent(parent)
		setParent(n.Condition, n)
		for i := range n.Body {
			setParent(&n.Body[i], n)
		}
		for i := range n.Else {
			setParent(&n.Else[i], n)
		}
	case *ast.TypeDef:
		n.SetParent(parent)
		setParent(n.TypeName, n)
		setParent(n.DeclType, n)
	case *ast.DimDecl:
		n.SetParent(parent)
		for i := range n.Vars {
			setParent(n.Vars[i], n)
		}
	case *ast.ConstDecl:
		n.SetParent(parent)
		for i := range n.Consts {
			setParent(&n.Consts[i], n)
		}
	case *ast.ConstDeclItem:
		n.SetParent(parent)
		setParent(n.ConstName, n)
		setParent(n.ConstValue, n)
		setParent(n.ConstType, n)
	case *ast.ScalarDecl:
		n.SetParent(parent)
		setParent(n.VarName, n)
		setParent(n.VarType, n)
		setParent(n.VarValue, n)
	case *ast.WhileStmt:
		n.SetParent(parent)
		setParent(n.Condition, n)
		for i := range n.Body {
			setParent(&n.Body[i], n)
		}
	case *ast.UntilStmt:
		n.SetParent(parent)
		setParent(n.Condition, n)
		for i := range n.Body {
			setParent(&n.Body[i], n)
		}
	case *ast.DoWhileStmt:
		n.SetParent(parent)
		setParent(n.Condition, n)
		for i := range n.Body {
			setParent(&n.Body[i], n)
		}
	case *ast.DoUntilStmt:
		n.SetParent(parent)
		setParent(n.Condition, n)
		for i := range n.Body {
			setParent(&n.Body[i], n)
		}
	case *ast.EnumDecl:
		n.SetParent(parent)
		setParent(n.Identifier, n)
		for i := 0; i < len(n.Values); i++ {
			setParent(&n.Values[i], n)
		}
	case *ast.ParamItem:
		n.SetParent(parent)
		setParent(n.VarName, n)
		setParent(n.VarType, n)
		setParent(n.DefaultValue, n)
	case *ast.UserDefinedType:
		n.SetParent(parent)
		setParent(n.Identifier, n)
	case *ast.ForEachExpr:
		n.SetParent(parent)
		setParent(n.Collection, n)
		setParent(n.Variable, n)
	case *ast.ForNextExpr:
		n.SetParent(parent)
		setParent(n.Variable, n)
		setParent(n.From, n)
		setParent(n.Step, n)
		setParent(n.To, n)
	case *ast.CallSubStmt:
		n.SetParent(parent)
		setParent(n.Definition, n)
	default:
		return errors.Newf(node.Token().Position, "unknown node type %T", node)
	}
	return nil
}
