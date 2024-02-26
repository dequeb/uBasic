// Package semcheck implements a static semantic analysis checker for µC.
package semcheck

import (
	"fmt"
	"strconv"
	"strings"
	"uBasic/ast"
	"uBasic/ast/astutil"
	"uBasic/errors"
	"uBasic/token"
	"uBasic/types"
)

// NoNestedFunctions disables the checking for nested functions
var NoNestedFunctions = false

// Check performs static semantic analysis on the given file.
func Check(file *ast.File) error {

	f := func(node ast.Node) error {
		switch node := node.(type) {
		case ast.FuncOrSub:
			if err := checkFunctionReturns(node); err != nil {
				return err
			}
			// check parameter definition
			// Optional at the end of the list
			// ParamArray as last without optional before
			if err := checkParamDef(node); err != nil {
				return err
			}
		case *ast.ExitStmt:
			// Check invalid exit statement
			if err := checkInvalidExitStmt(node); err != nil {
				return err
			}
		case *ast.CallOrIndexExpr:
			// check if function call returns a value
			if err := checkReturnValueNotAssigned(node); err != nil {
				return err
			}
		case *ast.BinaryExpr:
			// check if value is assigned to a constant
			if err := checkConstAssignment(node); err != nil {
				return err
			}
			// check type of operands
			// if err := checkBinaryExpr(node); err != nil {
			// 	return err
			// }
		// case *ast.UnaryExpr:
		// check type of operand
		// if err := checkUnaryExpr(node); err != nil {
		// 	return err
		// }
		case *ast.SpecialStmt:
			// validate that redim is done on a dynamic array
			if strings.ToLower(node.Keyword1.Literal) == "redim" {
				if err := checkRedim(node); err != nil {
					return err
				}
			}
			// validate erase statement
			if strings.ToLower(node.Keyword1.Literal) == "erase" {
				if err := checkErase(node); err != nil {
					return err
				}
			}
		// case *ast.BasicLit:
		// if date constant is validated by the parser
		case *ast.ForStmt:
			// check if next iterator correspond to for statement
			if err := checkNextIterator(node); err != nil {
				return err
			}
		case *ast.IfStmt:
			// check if there is no assignment in condition
			if err := checkIfStmt(node); err != nil {
				return err
			}
		}
		return nil

	}
	return astutil.Walk(file, f)
}

// checkInvalidExitStmt reports an error if the given function or sub declaration contains an Exit Statement of the wrong tyype .
func checkInvalidExitStmt(exitStmt *ast.ExitStmt) error {
	var f func(node ast.Node) error
	f = func(node ast.Node) error {
		stmtList, ok := node.(*ast.StatementList)
		if ok {
			return f(stmtList.GetParent())
		}

		stmt, ok := node.(ast.HasBody)
		if !ok {
			return errors.Newf(exitStmt.Token().Position, "exit statement of type %q cannot be used in this context", exitStmt.ExitType)
		}

		switch stmt.(type) {
		case *ast.FuncDecl:
			if exitStmt.ExitType.Kind == token.KwFunction {
				return nil
			}
		case *ast.SubDecl:
			if exitStmt.ExitType.Kind == token.KwSub {
				return nil
			}
		case *ast.ForStmt:
			if exitStmt.ExitType.Kind == token.KwFor {
				return nil
			}
		case *ast.WhileStmt, *ast.DoUntilStmt, *ast.DoWhileStmt, *ast.UntilStmt:
			if exitStmt.ExitType.Kind == token.KwDo {
				return nil
			}
		}

		return f(stmt.GetParent())
	}
	return f(exitStmt.GetParent())

	// var f func(node ast.Node) error
	// f = func(node ast.Node) error {
	// 	stmt, ok := node.(ast.HasBody)
	// 	if !ok {
	// 		return nil
	// 	}
	// 	for _, stmt := range stmt.GetBody() {
	// 		if exitStmt, ok := stmt.(*ast.ExitStmt); ok {
	// 			switch strings.ToLower(exitStmt.ExitType) {
	// 			case "function":
	// 				if !isFunction {
	// 					return errors.Newf(exitStmt.Token().Position, "function cannot exit with Exit Sub statement")
	// 				}
	// 			case "sub":
	// 				if isFunction {
	// 					return errors.Newf(exitStmt.Token().Position, "subroutine cannot exit with Exit Function statement")
	// 				}
	// 			}
	// 		} else if subStmt, ok := stmt.(ast.HasBody); ok {
	// 			if err := f(subStmt); err != nil {
	// 				return err
	// 			}
	// 		}
	// 	}
	// 	return nil
	// }
	// return astutil.Walk(fn, f)
}

// if function returns values in a naive approach: we search for at least one return statement
func checkFunctionReturns(fn ast.FuncOrSub) error {
	// check if function has a return type
	fnType, err := fn.Type()
	if err != nil {
		return err
	}
	if fnType == nil {
		return nil
	}
	// check if function is a subroutine
	if _, ok := fn.(*ast.SubDecl); ok {
		return nil
	}

	var f func(stmt ast.HasBody) error
	found := false
	f = func(stmt ast.HasBody) error {
		for _, stmtList := range stmt.GetBody() {
			for _, stmt := range stmtList.Statements {

				// search for expression statements
				if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
					// check if expression is an assignment to function name
					if binaryExpr, ok := exprStmt.Expression.(*ast.BinaryExpr); ok {
						if ident, ok := binaryExpr.Left.(*ast.Identifier); ok {
							if ident.Name == fn.Name().Name {
								found = true
								return nil
							}
						}
					}
				} else if subStmt, ok := stmt.(ast.HasBody); ok {
					if err := f(subStmt); err != nil {
						return err
					} else if found {
						return nil
					}
				}
			}
		}
		return nil
	}
	if err := f(fn); err != nil {
		return err
	} else if !found {
		return errors.Newf(fn.Token().Position, "function does not return a value")
	}
	return nil
}

// checkParamDef reports an error if the given function or sub declaration contains an invalid parameter definition.
func checkParamDef(fn ast.FuncOrSub) error {
	isOptional := false
	isParamArray := false

	for _, param := range fn.GetParams() {
		newOptional := param.Optional
		// verify that default function parameters are all at the end
		if newOptional && !isOptional {
			isOptional = true
		} else if !newOptional && (isOptional || isParamArray) {
			return errors.Newf(param.Token().Position, "optional parameters must be at the end of the parameter list")
		} else if param.ParamArray {
			if isParamArray {
				return errors.Newf(param.Token().Position, "parameter array must appears once and be the last parameter")
			}
			isParamArray = true
			if isOptional {
				return errors.Newf(param.Token().Position, "parameter array cannot follow optional parameters")
			}
			// must be an array
			paramType, err := param.Type()
			if err != nil {
				return err
			}
			if _, ok := paramType.(*types.Array); !ok {
				return errors.Newf(param.Token().Position, "parameter array must be an array")
			}
			// check if paramArray variable is an array
			// if _, ok := param.Type().(*ast.ArrayType); !ok {
			// 	return errors.Newf(param.Token().Position, "parameter array must be an array")
			// }
		}
	}
	return nil
}

// checkNextIterator reports an error if the given for loop contains an invalid next iterator.
func checkNextIterator(forStmt *ast.ForStmt) error {
	// check if next iterator correspond to for statement
	next := forStmt.Next
	if next == nil {
		return nil
	}
	nextIterator := next.Name
	forIterator, ok := forStmt.ForExpression.(*ast.ForNextExpr)
	if !ok && nextIterator != "" {
		return errors.Newf(forStmt.Token().Position, "invalid next identifier %q for statement", nextIterator)
	} else if ok && nextIterator != forIterator.Variable.Name {
		return errors.Newf(forStmt.Token().Position, "next iterator %q does not correspond to for statement", nextIterator)
	}
	return nil
}

// checkReturnValueNotAssigned reports an error if the given function call does not use the return value.
func checkReturnValueNotAssigned(call *ast.CallOrIndexExpr) error {

	callType, err := call.Identifier.Decl.Type()
	if err != nil {
		return err
	}
	if callType == nil {
		return nil
	}
	// check if function is a subroutine
	if _, ok := call.Identifier.Decl.(*ast.SubDecl); ok {
		return nil
	}

	// check if function call result is used
	if _, ok := call.GetParent().(*ast.ExprStmt); ok {
		return errors.Newf(call.Token().Position, "function call does not use return value")
	}
	return nil
}

// checkConstAssignment reports an error if the given binary expression assigns a value to a constant.
func checkConstAssignment(binaryExpr *ast.BinaryExpr) error {
	if binaryExpr.OpKind == token.Assign {
		if ident, ok := binaryExpr.Left.(*ast.Identifier); ok {
			if ident.Decl != nil {
				if _, ok := ident.Decl.(*ast.ConstDeclItem); ok {
					return errors.Newf(binaryExpr.Token().Position, "cannot assign value to constant %q", ident.Name)
				}
			}
		}
	}
	return nil
}

// checkRedim reports an error if the given redim statement is done on a non-dynamic array.
func checkRedim(redim *ast.SpecialStmt) error {
	// check if redim is done on a dynamic array
	if strings.ToLower(redim.Keyword1.Literal) == "redim" {
		// expect only one parameter on Redim statement
		if len(redim.Args) != 1 {
			return errors.Newf(redim.Token().Position, "redim expects one argument")
		}
		indexExpr, ok := redim.Args[0].(*ast.CallOrIndexExpr)
		if !ok {
			return errors.Newf(redim.Token().Position, "redim expects an array")
		}
		// check the number of parameters of redim statement
		arrayArgs := indexExpr.Args
		if len(arrayArgs) != 1 {
			return errors.Newf(redim.Token().Position, "redim expects one argument")
		}
		dimension := arrayArgs[0]
		if lit, ok := dimension.(*ast.BasicLit); ok {
			// ensure new dimension is valid (>0)
			if lit.Kind != token.LongLit {
				return errors.Newf(redim.Token().Position, "redim expects a valid number of dimensions")
			}
			// ensure new dimension is valid (>0)
			dimInt, err := strconv.Atoi(fmt.Sprint(lit.Value))
			if err != nil {
				return errors.Newf(redim.Token().Position, "redim expects a valid number of dimensions, not a string")
			}
			if dimInt <= 0 {
				return errors.Newf(redim.Token().Position, "redim expects a valid number of dimensions, greater than 0")
			}
		} else if indexExpr, ok := dimension.(ast.Expression); ok {
			// ensure expresion is long type
			typ, err := indexExpr.Type()
			if err != nil {
				return err
			}
			if typ, ok := typ.(*types.Basic); ok {
				if typ.Kind != types.Long {
					return errors.Newf(redim.Token().Position, "redim expects a valid number of dimensions of type long")
				}
			} else {
				return errors.Newf(redim.Token().Position, "redim expects a valid number of dimensions of type long")
			}
		} else {
			return errors.Newf(redim.Token().Position, "redim expects a valid number of dimensions")
		}
		// validate declaration
		decl := indexExpr.Identifier.Decl
		arrayDecl, ok := decl.(*ast.ArrayDecl)
		if !ok {
			return errors.Newf(redim.Token().Position, "redim expects an array")
		}
		// validate that redim is done on a dynamic array
		dimensions := arrayDecl.VarType.Dimensions
		if len(dimensions) == 0 { // dynamic array
			return nil
		}
		return errors.Newf(redim.Token().Position, "redim expects a dynamic array")
	}
	// all other statements are ignored
	return nil
}

// checkErase reports an error if the given erase statement is done on a non-dynamic array.
func checkErase(erase *ast.SpecialStmt) error {
	// check if erase is done on a dynamic array
	if strings.ToLower(erase.Keyword1.Literal) == "erase" {
		// expect only one parameter on erase statement
		if len(erase.Args) != 1 {
			return errors.Newf(erase.Token().Position, "erase expects one argument")
		}
		_, ok := erase.Args[0].(*ast.Identifier)
		if !ok {
			return errors.Newf(erase.Token().Position, "erase expects an identifier")
		}
		// get declaration
		ident := erase.Args[0].(*ast.Identifier)
		arrayDecl, ok := ident.Decl.(*ast.ArrayDecl)
		if !ok {
			return errors.Newf(erase.Token().Position, "erase expects an array")
		}
		dimensions := arrayDecl.VarType.Dimensions
		if len(dimensions) == 0 { // dynamic array
			return nil
		}
		return errors.Newf(erase.Token().Position, "erase expects a dynamic array")
	}
	// all other statements are ignored
	return nil
}

// checkIfStmt reports an error if the given if statement contains an assignment in the condition.
func checkIfStmt(ifStmt *ast.IfStmt) error {
	// check if there is no assignment in condition
	condition := ifStmt.Condition
	if binaryExpr, ok := condition.(*ast.BinaryExpr); ok {
		if binaryExpr.OpKind == token.Assign {
			return errors.Newf(binaryExpr.Token().Position, "assignment in condition")
		}
	}
	return nil
}
