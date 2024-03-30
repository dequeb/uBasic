// Package typecheck implements type-checking of parse trees.
package typecheck

import (
	"uBasic/ast"
	"uBasic/ast/astutil"
	"uBasic/errors"
	"uBasic/token"
	"uBasic/types"
)

// Check type-checks the given file, and store a mapping from expression nodes
// to types in exprTypes.
func Check(file *ast.File, exprTypes map[ast.Expression]types.Type) error {
	// Deduce the types of expressions.
	if err := deduce(file, exprTypes); err != nil {
		return err
	}

	// Type-check file.
	if err := check(file, exprTypes); err != nil {
		return err
	}

	return nil
}

// check type-checks the given file.
func check(file *ast.File, exprTypes map[ast.Expression]types.Type) error {
	// funcsAndSubs is a stack of function and subroutine declarations, where the top-most entry
	// represents the currently active function.
	var funcsAndSubs []types.SubOrFunc

	// check type-checks the given node.
	check := func(n ast.Node) error {
		switch n := n.(type) {
		case *ast.UnaryExpr:
			// Check that the operand type is compatible with the operator.
			if n.OpKind == token.Not { // logical negation
				typ := &types.Basic{Kind: types.Boolean}
				if !isCompatible(exprTypes[n.Right], typ) {
					return errors.Newf(n.OpToken.Position, "incompatible type for operator %q; expected %q, got %q", n.OpKind, typ, exprTypes[n.Right])
				}
			} else if n.OpKind == token.Minus { // arithmetic negation
				typ := &types.Basic{Kind: types.Long}
				if !isCompatible(exprTypes[n.Right], typ) {
					return errors.Newf(n.OpToken.Position, "incompatible type for operator %q; expected %q, got %q", n.OpKind, typ, exprTypes[n.Right])
				}
			}
		case *ast.BinaryExpr:
			// Check that the operand types are compatible with the operator.
			typX := exprTypes[n.Left]
			typY := exprTypes[n.Right]
			switch n.OpKind {
			case token.Add, token.Minus, token.Mul, token.Div, token.IntDiv, token.Mod, token.Exponent: // arithmetic operators
				typ := &types.Basic{Kind: types.Long}
				if !isCompatible(typX, typ) || !isCompatible(typY, typ) {
					return errors.Newf(n.OpToken.Position, "incompatible type for operator %q; expected %q, got %q and %q", n.OpKind, typ, typX, typY)
				}
			case token.Eq, token.Neq, token.Lt, token.Le, token.Gt, token.Ge: // relational operators
				if !isCompatible(typX, typY) && !isFunctionCompatible(n.Left, n.Right) {
					return errors.Newf(n.OpToken.Position, "incompatible type for operator %q; between %q and %q", n.OpKind, typX, typY)
				}
			case token.And, token.Or: // logical operators
				typ := &types.Basic{Kind: types.Boolean}
				if !isCompatible(typX, typ) || !isCompatible(typY, typ) {
					return errors.Newf(n.OpToken.Position, "incompatible type for operator %q; expected %q, got %q and %q", n.OpKind, typ, typX, typY)
				}
			case token.Assign: // assignment operator
				// is assignable was verified in deduce
				// if !isAssignable(n.X) {
				// 	return errors.Newf(n.OpToken.Position, "cannot assign to %q of type %q", n.X, typX)
				// }
				if t, ok := typX.(*types.Func); ok {
					typX = t.Result
				}

				if !isCompatible(typX, typY) {
					return errors.Newf(n.OpToken.Position, "incompatible type for operator %q; between %q and %q", n.OpKind, typX, typY)
				}
			case token.Concat: // concatenation operator
				typ := &types.Basic{Kind: types.String}
				if !isCompatible(typX, typ) || !isCompatible(typY, typ) {
					return errors.Newf(n.OpToken.Position, "incompatible type for operator %q; expected %q, got %q and %q", n.OpKind, typ, typX, typY)
				}
			default:
				return errors.Newf(n.OpToken.Position, "support for operator %q not yet implemented", n.OpKind)
			}
		case *ast.BasicLit:
			// assign type to basic literal received from lexer
			switch n.Kind {
			case token.StringLit:
				exprTypes[n] = &types.Basic{Kind: types.String}
			case token.LongLit:
				exprTypes[n] = &types.Basic{Kind: types.Long}
			case token.DoubleLit:
				exprTypes[n] = &types.Basic{Kind: types.Double}
			case token.BooleanLit, token.KwTrue, token.KwFalse:
				exprTypes[n] = &types.Basic{Kind: types.Boolean}
			case token.DateLit:
				exprTypes[n] = &types.Basic{Kind: types.Date}
			case token.CurrencyLit:
				exprTypes[n] = &types.Basic{Kind: types.Currency}
			case token.KwNothing:
				exprTypes[n] = &types.Basic{Kind: types.Variant}
			default:
				return errors.Newf(n.Token().Position, "support for basic literal type %v not yet implemented", n.Kind)
			}
		case *ast.ConstDecl:
			// Check that the constant type is compatible with the expression type.
			for _, cons := range n.Consts {
				if consType, err := cons.Type(); err != nil {
					return err
				} else if !isCompatible(exprTypes[cons.ConstValue], consType) {
					return errors.Newf(cons.Token().Position, "incompatible type for constant %q; expected %q, got %q", cons.ConstName, consType, exprTypes[cons.ConstValue])
				}
			}
		case *ast.ScalarDecl:
			// Check that the variable type is compatible with the expression type.
			if n.VarValue != nil {
				if nType, err := n.Type(); err != nil {
					return err
				} else if !isCompatible(exprTypes[n.VarValue], nType) {
					return errors.Newf(n.Token().Position, "incompatible type for variable %q; expected %q, got %q", n.VarName, nType, exprTypes[n.VarValue])
				}
			}
		case ast.FuncOrSub:
			if astutil.IsDef(n) {
				// push function declaration.
				var subType *types.Sub
				nType, err := n.Type()
				if err != nil {
					return err
				}
				funcType, ok := nType.(*types.Func)
				if !ok {
					subType, ok = nType.(*types.Sub)
					if !ok {
						return errors.Newf(n.Token().Position, "cannot declare non-function %s of type %s", n.Name(), nType)
					}
				}

				if funcType != nil {
					funcType.Params = make([]*types.Field, len(n.GetParams()))
					for i, param := range n.GetParams() {
						funcType.Params[i], err = ast.NewField(&param)
						if err != nil {
							return err
						}
						if param.DefaultValue != nil {
							funcType.Params[i].DefaultValue = param.DefaultValue.String()
						}
					}
					nTyp, err := n.Type()
					if err != nil {
						return err
					}
					funcType.Result = nTyp
					funcsAndSubs = append(funcsAndSubs, funcType)
				} else if subType != nil {
					subType.Params = make([]*types.Field, len(n.GetParams()))
					for i, param := range n.GetParams() {
						subType.Params[i], err = ast.NewField(&param)
						if err != nil {
							return err
						}
						if param.DefaultValue != nil {
							subType.Params[i].DefaultValue = param.DefaultValue.String()
						}
					}
					funcsAndSubs = append(funcsAndSubs, subType)
				}
			}
		case *ast.CallOrIndexExpr:
			fType, err := n.Identifier.Decl.Type()
			if err != nil {
				return err
			}
			funcType, ok := fType.(*types.Func)
			if !ok {
				subType, ok := fType.(*types.Sub)
				if ok {
					// verify that parameters are compabible with arguments
					if err := VerifyParameters(n, subType, exprTypes); err != nil {
						return err
					}
				} else {
					array, ok := fType.(*types.Array)
					if !ok {
						return errors.Newf(n.Lparen.Position, "cannot call non-function or non-array %q of type %q", n.Identifier, funcType)
					}
					// verify number of dimensions of array
					if len(n.Args) != len(array.Dimensions) {
						if len(array.Dimensions) == 0 && len(n.Args) == 1 {
							// dynamic array
							return nil
						}
						return errors.Newf(n.Lparen.Position, "calling %q with wrong number of dimensions; expected %d, got %d", n.Identifier, len(array.Dimensions), len(n.Args))
					}
				}
			} else {
				// verify that parameters are compabible with arguments
				if err := VerifyParameters(n, funcType, exprTypes); err != nil {
					return err
				}
			}
		case *ast.Identifier:
			if n.Decl == nil {
				return errors.Newf(n.Token().Position, "undeclared identifier %q", n.Name)
			}
		case *ast.File, *ast.EmptyStmt, *ast.ParamItem:
			// nothing to do
			//default:
			// fmt.Printf("not type-checked: %T\n", n)
		}
		return nil
	}

	// after reverts to the outer function after traversing function definitions.
	after := func(n ast.Node) error {
		switch n := n.(type) {
		case *ast.FuncDecl, *ast.SubDecl:
			if astutil.IsDef(n) {
				// pop function declaration.
				funcsAndSubs = funcsAndSubs[:len(funcsAndSubs)-1]
			}
		}
		return nil
	}

	// Walk the AST of the given file to perform type-checking.
	if err := astutil.WalkBeforeAfter(file, check, after); err != nil {
		return err
	}

	return nil
}

// VerifyParameters verifies that the given call arguments are compatible with
// the given function parameters.
func VerifyParameters(call *ast.CallOrIndexExpr, funcOrSubType types.SubOrFunc, exprTypes map[ast.Expression]types.Type) error {
	lastParamIsParamArray := false
	if len(funcOrSubType.GetParams()) > 0 {
		lastParamIsParamArray = funcOrSubType.GetParams()[len(funcOrSubType.GetParams())-1].ParamArray
	}

	// there are 3 possibles cases:
	// 1. same number of parameters and arguments
	// 2. more parameters than arguments (optional parameters)
	// 3. more arguments than parameters (paramArray argument)
	// lets start with the first case
	if len(call.Args) == len(funcOrSubType.GetParams()) && !lastParamIsParamArray {
		// verify that required parameters are compatible with arguments
		for i, arg := range call.Args {
			if !isCompatibleArg(funcOrSubType.GetParams()[i].Type, exprTypes[arg]) {
				return errors.Newf(arg.Token().Position, "incompatible type for argument %d; expected %q, got %q", i+1, funcOrSubType.GetParams()[i].Type, exprTypes[arg])
			}
		}
		// second case
	} else if len(call.Args) < len(funcOrSubType.GetParams()) {
		// verify that required parameters are compatible with arguments
		for i, arg := range call.Args {
			if !isCompatibleArg(funcOrSubType.GetParams()[i].Type, exprTypes[arg]) {
				return errors.Newf(arg.Token().Position, "incompatible type for argument %d; expected %q, got %q", i+1, funcOrSubType.GetParams()[i].Type, exprTypes[arg])
			}
		}
		// verify that optional parameters are compatible with arguments
		for i := len(call.Args); i < len(funcOrSubType.GetParams()); i++ {
			if !funcOrSubType.GetParams()[i].Optional && !funcOrSubType.GetParams()[i].ParamArray {
				return errors.Newf(call.Lparen.Position, "missing argument for parameter %v", funcOrSubType.GetParams()[i].Name)
			}
		}
		// third case
	} else {
		// no arguments or paramArray
		if len(funcOrSubType.GetParams()) == 0 {
			return errors.Newf(call.Lparen.Position, "too many arguments for function %q", call.Identifier)
		}

		// verify that required parameters are compatible with arguments
		for i := 0; i < len(funcOrSubType.GetParams())-2; i++ { // do not include paramArray
			arg := call.Args[i]
			if !isCompatibleArg(funcOrSubType.GetParams()[i].Type, exprTypes[arg]) {
				return errors.Newf(arg.Token().Position, "incompatible type for argument %d; expected %q, got %q", i+1, funcOrSubType.GetParams()[i].Type, exprTypes[arg])
			}
		}

		// verify if last parameter is a paramArray
		if !funcOrSubType.GetParams()[len(funcOrSubType.GetParams())-1].ParamArray {
			return errors.Newf(call.Lparen.Position, "too many arguments for function %q", call.Identifier)
		}

		// verify that other parameters are compatible with paramArray type
		paramArray := funcOrSubType.GetParams()[len(funcOrSubType.GetParams())-1]
		elementType := paramArray.Type.(*types.Array).Type
		for i := len(funcOrSubType.GetParams()) - 1; i < len(call.Args); i++ {
			arg := call.Args[i]
			if !isCompatibleArg(elementType, exprTypes[arg]) {
				return errors.Newf(call.Lparen.Position, "incompatible type for parameter array; expected %q, got %q", elementType, exprTypes[call.Args[len(call.Args)-1]])
			}
		}
	}
	return nil
}

// isCompatibleArg reports whether the given call argument and function or array
// parameter types are compatible.
func isCompatibleArg(arg, param types.Type) bool {
	if isCompatible(arg, param) {
		return true
	}
	if arg, ok := arg.(*types.Array); ok {
		if param, ok := param.(*types.Array); ok {
			if len(arg.Dimensions) == 0 {
				// dynamic array of len(dimension) = 1
				if len(param.Dimensions) == 1 {
					return isCompatible(arg.Type, param.Type)
				}
			}
		}
	}
	return false
}

// isCompatible reports whether t and u are of compatible types.
func isCompatible(t, u types.Type) bool {
	if t, ok := t.(*types.ByRef); ok {
		return isCompatible(t.Type, u)
	}
	if u, ok := u.(*types.ByRef); ok {
		return isCompatible(t, u.Type)
	}

	if types.Equal(t, u) {
		return true
	}
	// if t is a variant, we can assign all types to it
	if t.IsVariant() {
		return true
	}
	if t, ok := t.(types.Numerical); ok {
		if u, ok := u.(types.Numerical); ok {
			return t.IsNumerical() && u.IsNumerical()
		}
	}
	return false
}

// isFunctionCompatible reports whether t and u are of compatible types within function scope
func isFunctionCompatible(left, right ast.Expression) bool {
	leftType, err := TypeOf(left)
	if err != nil {
		return false
	}
	rightType, err := TypeOf(right)
	if err != nil {
		return false
	}
	functType, ok := leftType.(*types.Func)
	if ok {
		if isCompatible(functType.Result, rightType) {
			// are we in the function scope?
			for parent := left.GetParent(); parent != nil; parent = parent.GetParent() {
				if p, ok := parent.(*ast.FuncDecl); ok {
					if p.FuncName.Name == left.(*ast.Identifier).Name {
						return true
					}
				}
			}
		}
	}
	functType, ok = rightType.(*types.Func)
	if ok {
		if isCompatible(leftType, functType.Result) {
			// are we in the function scope?
			for parent := right.GetParent(); parent != nil; parent = parent.GetParent() {
				if p, ok := parent.(*ast.FuncDecl); ok {
					if p.FuncName.Name == right.(*ast.Identifier).Name {
						return true
					}
				}
			}
		}
	}

	return false
}
