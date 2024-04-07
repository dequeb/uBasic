// Package astutil implements utility functions for handling parse trees.
package astutil

import "uBasic/ast"

// IsDef reports whether the given declaration is a definition.
func IsDef(decl ast.Node) bool {
	switch decl.(type) {
	case *ast.ScalarDecl, *ast.ArrayDecl, *ast.ConstDeclItem, *ast.FuncDecl, *ast.SubDecl, *ast.EnumDecl, *ast.ClassDecl, *ast.ParamItem:
		return true
	default:
		return false
	}
}
