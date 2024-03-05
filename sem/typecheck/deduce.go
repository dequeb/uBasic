package typecheck

import (
	"fmt"
	"uBasic/token"
	"uBasic/types"

	"uBasic/ast"
	"uBasic/ast/astutil"
	"uBasic/errors"
	//"token"
)

// deduce performs type deduction of expressions, and store the result in
// exprTypes.
func deduce(file *ast.File, exprTypes map[ast.Expression]types.Type) error {
	// deduce performs type deduction of the given expression.
	deduce := func(n ast.Node) error {
		if expr, ok := n.(ast.Expression); ok {
			typ, err := TypeOf(expr)
			if err != nil {
				return err
			}
			exprTypes[expr] = typ
		}
		return nil
	}

	// Walk the AST of the given file to deduce the types of expression nodes.
	if err := astutil.Walk(file, deduce); err != nil {
		return err
	}

	return nil
}

// TypeOf returns the type of the given expression.
func TypeOf(n ast.Expression) (types.Type, error) {
	switch n := n.(type) {
	case *ast.BasicLit:
		switch n.Kind {
		case token.StringLit:
			return &types.Basic{Kind: types.String}, nil
		case token.LongLit:
			return &types.Basic{Kind: types.Long}, nil
		case token.DoubleLit:
			return &types.Basic{Kind: types.Double}, nil
		case token.BooleanLit, token.KwTrue, token.KwFalse:
			return &types.Basic{Kind: types.Boolean}, nil
		case token.DateLit:
			return &types.Basic{Kind: types.Date}, nil
		case token.KwNothing:
			return &types.Basic{Kind: types.Variant}, nil
		case token.CurrencyLit:
			return &types.Basic{Kind: types.Currency}, nil
		default:
			panic(fmt.Sprintf("support for basic literal type %v not yet implemented", n.Kind))
		}
	case *ast.BinaryExpr:
		xType, err := TypeOf(n.Left)
		if err != nil {
			return nil, err
		}
		yType, err := TypeOf(n.Right)
		if err != nil {
			return nil, err
		}
		if n.OpKind == token.Assign {
			if !isAssignable(n.Left) {
				return nil, errors.Newf(n.OpToken.Position, "cannot assign to %q of type %q", n.Left, xType)
			}
			return xType, nil
		} else {
			if n.OpKind == token.Gt || n.OpKind == token.Ge || n.OpKind == token.Lt || n.OpKind == token.Le || n.OpKind == token.Eq || n.OpKind == token.Neq {
				if !isCompatible(xType, yType) && !isFunctionCompatible(n.Left, n.Right) {
					return nil, errors.Newf(n.Token().Position, "invalid operation: %v (type mismatch between %q and %q)", n, xType, yType)
				}
				return &types.Basic{Kind: types.Boolean}, nil
			}
		}
		if !isCompatible(xType, yType) {
			return nil, errors.Newf(n.Token().Position, "invalid operation: %v (type mismatch between %q and %q)", n, xType, yType)
		}
		return ast.HigherPrecision(xType, yType)
	case *ast.CallOrIndexExpr:
		typ, err := n.Identifier.Decl.Type()
		if err != nil {
			return nil, err
		}

		switch subType := typ.(type) {
		case *types.Func:
			return subType.Result, nil
		case *types.Array:
			return subType.Type, nil
		case *types.Sub:
			// search a call statement in parent nodes
			for parent := n.GetParent(); parent != nil; parent = parent.GetParent() {
				if _, ok := parent.(*ast.CallSubStmt); ok {
					return subType, nil
				}
			}
			return nil, errors.Newf(n.Lparen.Position, "cannot call non-function or non-array %q of type %q", n.Identifier, typ)
		default:
			return nil, errors.Newf(n.Lparen.Position, "cannot call non-function or non-array %q of type %q", n.Identifier, typ)
		}

	case *ast.Identifier:
		if n.Decl == nil {
			panic("not yet implemented")
		}
		return n.Decl.Type()
	case *ast.CallSelectorExpr:
		var selectorDecl ast.Decl
		switch selector := n.Selector.(type) {
		case *ast.Identifier:
			selectorDecl = selector.Decl
		case *ast.CallOrIndexExpr:
			identifier := selector.Identifier
			selectorDecl = identifier.Decl
		default:
			panic(fmt.Sprintf("support for type %T not yet implemented.", n.Root))
		}
		if selectorDecl == nil {
			panic("not yet implemented")
		}
		return selectorDecl.Type()
	case *ast.ParenExpr:
		return TypeOf(n.Expr)
	case *ast.UnaryExpr:
		return TypeOf(n.Right)
	default:
		panic(fmt.Sprintf("support for type %T not yet implemented.", n))
	}
}

// isAssignable reports whether the given expression is assignable (i.e. a valid
// lvalue).
func isAssignable(x ast.Expression) bool {
	switch x := x.(type) {
	case *ast.BasicLit:
		return false
	case *ast.BinaryExpr:
		return false
	case *ast.CallOrIndexExpr:
		xType, err := x.Identifier.Decl.Type()
		if err != nil {
			return false
		}
		if _, ok := xType.(*types.Array); ok {
			return true
		}
		return false
	case *ast.Identifier:
		xType, err := x.Decl.Type()
		if err != nil {
			return false
		}

		if xType == nil {
			return false // this is a Subroutine
		}
		switch typ := xType.(type) {
		case *types.Basic:
			return true
		case *types.Array:
			return false
		case *types.Func:
			// if we are in function scope, then we can assign to a function
			// identifier.
			for parent := x.GetParent(); parent != nil; parent = parent.GetParent() {
				if p, ok := parent.(*ast.FuncDecl); ok {
					if p.FuncName.Name == x.Name {
						return true
					}
				}
			}
			return false
		case *types.UserDefined:
			return true
		default:
			panic(fmt.Sprintf("support for declaration type %T not yet implemented", typ))
		}
	case *ast.ParenExpr:
		return isAssignable(x.Expr)
	case *ast.UnaryExpr:
		return false
	case *ast.CallSelectorExpr:
		return false

	default:
		panic(fmt.Sprintf("support for expression type %T not yet implemented", x))
	}
}
