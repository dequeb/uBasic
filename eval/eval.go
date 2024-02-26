package eval

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
	"uBasic/ast"
	"uBasic/eval/rtlib"
	"uBasic/object"
	"uBasic/sem"
	"uBasic/token"
)

const returnName = "result!"

func Define(info *sem.Info) *object.Environment {
	env := object.NewEnvironment()
	scopes := info.Scopes
	for _, s := range scopes {
		for k, v := range s.Outer.Decls {
			switch v := v.(type) {
			case *ast.ScalarDecl:
				env.Set(k, evalScalarDecl(v, env))
			case *ast.ArrayDecl:
				env.Set(k, evalArrayDecl(v, env))
			case *ast.ConstDeclItem:
				env.Set(k, evalConstDeclItem(v, env))
			case *ast.EnumDecl:
				env.Set(k, evalEnumDecl(v, env))
			case *ast.TypeDef:
				env.Set(k, evalTypeDef(v, env))
			case *ast.FuncDecl:
				env.Set(k, evalFuncDecl(v, env))
			case *ast.SubDecl:
				env.Set(k, evalSubDecl(v, env))
			case *ast.ClassDecl:
				env.Set(k, evalClassDecl(v, env))
			default:
				fmt.Println("unknown declaration: ", v.String())
			}
		}
	}
	return env.Extend()
}

var callBack func(ast.Node) bool
var currentLine int

func Run(file *ast.File, env *object.Environment, f func(ast.Node) bool) object.Object {
	callBack = f
	return Eval(nil, file, env)
}

func Eval(class *object.Class, node ast.Node, env *object.Environment) object.Object {
	if callBack != nil {
		// only call back on node with token (for line number)
		if node.Token() != nil && node.Token().Position.Line != currentLine {
			currentLine = node.Token().Position.Line
			if !callBack(node) {
				return &object.Exit{}
			}
		}
	}

	switch node := node.(type) {
	case *ast.BasicLit:
		return evalBasicLit(node, env)
	case *ast.File:
		return evalFile(node, env)
	case *ast.StatementList:
		return evalStatementList(node, env)
	case *ast.FuncDecl:
		return evalFuncDecl(node, env)
	case *ast.SubDecl:
		return evalSubDecl(node, env)
	case *ast.Identifier:
		return evalIdentifier(class, node, env)
	case *ast.ArrayDecl:
		return evalArrayDecl(node, env)
	case *ast.ScalarDecl:
		return evalScalarDecl(node, env)
	case *ast.DimDecl:
		return evalDimDecl(node, env)
	case *ast.ConstDecl:
		return evalConstDecl(node, env)
	case *ast.ConstDeclItem:
		return evalConstDeclItem(node, env)
	case *ast.EnumDecl:
		return evalEnumDecl(node, env)
	case *ast.TypeDef:
		return evalTypeDef(node, env)
	case *ast.ExprStmt:
		return evalExprStmt(node, env)
	case *ast.EmptyStmt:
		return evalEmptyStmt(node, env)
	case *ast.IfStmt:
		return evalIfStmt(node, env)
	case *ast.WhileStmt:
		return evalWhileStmt(node, env)
	case *ast.UntilStmt:
		return evalUntilStmt(node, env)
	case *ast.DoWhileStmt:
		return evalDoWhileStmt(node, env)
	case *ast.DoUntilStmt:
		return evalDoUntilStmt(node, env)
	case *ast.ForStmt:
		return evalForStmt(node, env)
	case *ast.BinaryExpr:
		return evalBinaryExpr(node, env)
	case *ast.UnaryExpr:
		return evalUnaryExpr(node, env)
	case *ast.CallOrIndexExpr:
		return evalCallOrIndexExpr(class, node, env)
	case *ast.SpecialStmt:
		return evalSpecialStmt(node, env)
	case *ast.ExitStmt:
		return evalExitStmt(node, env)
	case *ast.ParenExpr:
		return evalParenExpr(node, env)
	case *ast.CallSubStmt:
		return evalCallSubStmt(class, node, env)
	case *ast.CallSelectorExpr:
		return evalCallSelectorExpr(node, env)
	case *ast.ForNextExpr:
		return evalFinishedForNextExpr(node, env)
	case *ast.ForEachExpr:
		return evalFinishedForEachExpr(node, env)
	case *ast.SelectStmt:
		return evalSelectStmt(node, env)
	default:
		return object.NewError(node.Token().Position, "unknown node type: "+node.String())
	}
}

func evalBasicLit(node *ast.BasicLit, env *object.Environment) object.Object {
	switch node.Kind {
	case token.LongLit:
		str := node.Value.(string)
		val, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return object.NewError(node.Token().Position, err.Error())
		}
		return object.NewLongByInt(val, node.Token().Position)
	case token.DoubleLit:
		str := node.Value.(string)
		val, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return object.NewError(node.Token().Position, err.Error())
		}
		return object.NewDoubleByFloat(val, node.Token().Position)
	case token.StringLit:
		return object.NewString(node.Value.(string), node.Token().Position)
	case token.DateLit:
		dt, err := object.NewDate(node.Value.(string), node.Token().Position)
		if err != nil {
			return object.NewError(node.Token().Position, err.Error())
		}
		return dt
	case token.CurrencyLit:
		cur, err := object.NewCurrency(node.Value.(string), node.Token().Position)
		if err != nil {
			return object.NewError(node.Token().Position, err.Error())
		}
		return cur
	case token.BooleanLit:
		return object.NewBooleanByBool(node.Value.(bool), node.Token().Position)
	case token.KwTrue:
		return object.NewBooleanByBool(true, node.Token().Position)
	case token.KwFalse:
		return object.NewBooleanByBool(false, node.Token().Position)
	case token.KwNothing:
		return object.NOTHING
	}
	return object.NewError(node.Token().Position, "unknown basic literal: "+node.String())
}

func evalFile(node *ast.File, env *object.Environment) object.Object {
	return evalBody(node.StatementLists, env)
}

func evalBody(node []ast.StatementList, env *object.Environment) object.Object {
	var result object.Object

	for _, stmtList := range node {
		result = Eval(nil, &stmtList, env)
		if result != nil {
			switch result := result.(type) {
			case *object.Error:
				result.Stack.Push(stmtList.Token().Position)
				return result
			case *object.ReturnValue:
				return result.Value
			case *object.Exit:
				return result
			}
		}
	}
	return result
}

func evalStatementList(node *ast.StatementList, env *object.Environment) object.Object {
	var result object.Object
	for _, stmt := range node.Statements {
		result = Eval(nil, stmt, env)
		if result != nil {
			switch result := result.(type) {
			case *object.Error, *object.ReturnValue, *object.Exit:
				return result
			}
		}
	}
	return result
}

func evalFuncDecl(node *ast.FuncDecl, env *object.Environment) object.Object {
	params := node.FuncType
	body := node.Body
	fn := &object.Function{Definition: node, Parameters: params, Body: body, Env: env}
	// add the function to the environment
	env.Set(node.FuncName.Name, fn)
	return fn
}

func evalSubDecl(node *ast.SubDecl, env *object.Environment) object.Object {
	params := node.SubType
	body := node.Body
	sub := &object.Sub{Definition: node, Parameters: params, Body: body, Env: env}
	// add the subroutine to the environment
	env.Set(node.SubName.Name, sub)
	return sub
}

// create a new instance of the class in environment
func evalClassDecl(node *ast.ClassDecl, env *object.Environment) object.Object {
	className := node.ClassName
	class := object.NewClass(className.Name, node.Token().Position)
	// add class members to class object
	for name, member := range node.Members {
		switch member := member.(type) {
		case *ast.ScalarDecl:
			memberObj := evalScalarDecl(member, env)
			class.Members[name] = memberObj
		case *ast.ArrayDecl:
			memberObj := evalArrayDecl(member, env)
			class.Members[name] = memberObj
		case *ast.FuncDecl:
			memberObj := evalFuncDecl(member, env)
			class.Members[name] = memberObj
		case *ast.SubDecl:
			memberObj := evalSubDecl(member, env)
			class.Members[name] = memberObj
		case *ast.ClassDecl:
			memberObj := evalClassDecl(member, env)
			class.Members[name] = memberObj
		default:
			return object.NewError(node.Token().Position, "unknown class member: "+member.String())
		}
	}
	// add the class to the environment
	env.Set(className.Name, class)
	return class
}

func evalIdentifier(root object.Object, node *ast.Identifier, env *object.Environment) object.Object {
	if !object.IsNil(root) {
		// check if the identifier is a class
		class, ok := root.(*object.Class)
		if ok {
			// check if the identifier is a class member
			member, ok := class.Members[strings.ToLower(node.Name)]
			if !ok {
				return object.NewError(node.Token().Position, "undefined class member: "+node.Name)
			}
			return member
		}
		// check if root is an enum
		enum, ok := root.(*object.UserDefined)
		if ok {
			// check if the identifier is an enum member
			for _, v := range enum.Decl.Values {
				if v.Name == node.Name {
					return object.NewUserDefined(node.Name, enum.Decl, node.Token().Position)
				}
			}
		}
		return object.NewError(node.Token().Position, "undefined enum member: "+node.Name)
	}
	val, ok := env.Get(node.Name)
	if !ok {

		if !ok {
			return object.NewError(node.Token().Position, "undefined variable: "+node.Name)
		}
	}
	return val
}

func evalArrayDecl(node *ast.ArrayDecl, env *object.Environment) object.Object {
	dimension := make([]uint32, len(node.VarType.Dimensions))
	for i, dim := range node.VarType.Dimensions {
		val := Eval(nil, dim, env)
		if val != nil {
			switch val := val.(type) {
			case *object.Error:
				return val
			case *object.Integer:
				dimension[i] = uint32(val.Value)
			case *object.Long:
				dimension[i] = uint32(val.Value)
			}
		}
	}
	typ := strings.ToLower(node.VarType.Type.(*ast.Identifier).Name)
	switch typ {
	case "long", "_long":
		env.Set(node.VarName.Name, &object.Array{SubType: object.LONG_OBJ, Dimensions: dimension, Pos: node.Token().Position})
	case "integer", "_integer":
		env.Set(node.VarName.Name, &object.Array{SubType: object.INTEGER_OBJ, Dimensions: dimension, Pos: node.Token().Position})
	case "single", "_single":
		env.Set(node.VarName.Name, &object.Array{SubType: object.SINGLE_OBJ, Dimensions: dimension, Pos: node.Token().Position})
	case "double", "_double":
		env.Set(node.VarName.Name, &object.Array{SubType: object.DOUBLE_OBJ, Dimensions: dimension, Pos: node.Token().Position})
	case "string", "_string":
		env.Set(node.VarName.Name, &object.Array{SubType: object.STRING_OBJ, Dimensions: dimension, Pos: node.Token().Position})
	case "date", "_date":
		env.Set(node.VarName.Name, &object.Array{SubType: object.DATE_OBJ, Dimensions: dimension, Pos: node.Token().Position})
	case "boolean", "_boolean":
		env.Set(node.VarName.Name, &object.Array{SubType: object.BOOLEAN_OBJ, Dimensions: dimension, Pos: node.Token().Position})
	case "currency", "_currency":
		env.Set(node.VarName.Name, &object.Array{SubType: object.CURRENCY_OBJ, Dimensions: dimension, Pos: node.Token().Position})
	case "variant", "_variant":
		env.Set(node.VarName.Name, &object.Array{SubType: object.VARIANT_OBJ, Dimensions: dimension, Pos: node.Token().Position})
	default:
		return object.NewError(node.Token().Position, "unknown type: "+typ)
	}
	return nil
}

func evalScalarDecl(node *ast.ScalarDecl, env *object.Environment) object.Object {
	typ, ok := node.VarType.(*ast.Identifier)
	if ok {
		name := strings.ToLower(typ.Name)
		switch name {
		// see comment under UserDefinedType
		case "long", "long$":
			return object.NewLongByInt(0, node.Token().Position)
		case "integer", "integer$":
			return object.NewIntegerByInt(0, node.Token().Position)
		case "single", "single$":
			return object.NewSingleByFloat(0, node.Token().Position)
		case "double", "double$":
			return object.NewDoubleByFloat(0, node.Token().Position)
		case "string", "string$":
			return object.NewString("", node.Token().Position)
		case "date", "date$":
			return object.NewDateByTime(time.Time{}, node.Token().Position)
		case "boolean", "boolean$":
			return object.NewBooleanByBool(false, node.Token().Position)
		case "currency", "currency$":
			return object.NewCurrencyByFloat(0, node.Token().Position)
		case "variant", "variant$":
			return object.NewVariantByObject(object.NOTHING, node.Token().Position)
		}
	}
	return object.NewError(node.Token().Position, "unknown type: "+typ.Name)
}

func evalDimDecl(node *ast.DimDecl, env *object.Environment) object.Object {
	for _, v := range node.Vars {
		val := Eval(nil, v, env)
		if val != nil {
			switch val := val.(type) {
			case *object.Error:
				return val
			}
			identifier := v.Name()
			env.Set(identifier.Name, val)
		}
	}
	return nil
}
func evalConstDecl(node *ast.ConstDecl, env *object.Environment) object.Object {
	for _, v := range node.Consts {
		val := Eval(nil, &v, env)
		if val != nil {
			switch val := val.(type) {
			case *object.Error:
				return val
			}
			identifier := v.Name()
			env.Set(identifier.Name, val)
		}
	}
	return nil
}

func evalConstDeclItem(node *ast.ConstDeclItem, env *object.Environment) object.Object {
	value := Eval(nil, node.ConstValue, env)

	userDefinedConst, ok := node.ConstType.(*ast.UserDefinedType)
	if ok {
		decl, ok := userDefinedConst.Identifier.Decl.(*ast.EnumDecl)
		if ok {
			// check if value in values of the enum
			enum := decl.Identifier.Name
			enumDecl, ok := env.Get(enum)
			if !ok {
				return object.NewError(node.Token().Position, "undefined enum: "+enum)
			}
			enumObj, ok := enumDecl.(*object.UserDefined)
			if !ok {
				return object.NewError(node.Token().Position, "not an enum: "+enum)
			}
			enumValues := enumObj.Decl.Values
			found := false
			for _, v := range enumValues {
				if v.Name == value.String() {
					found = true
					break
				}
			}
			if !found {
				return object.NewError(node.Token().Position, "invalid value for enum: "+value.String())
			}
			//  we have an enum constant
			obj := object.NewUserDefined(node.ConstName.Name, decl, node.Token().Position)
			obj.Const = true
			obj.Value = value.String()
			return obj
		}
	}
	typ, ok := node.ConstType.(*ast.Identifier)
	if ok {
		val := node.ConstValue.String()
		switch typ.Name {
		case "Long$":
			obj, err := object.NewLong(val, node.Token().Position)
			if err != nil {
				return object.NewError(node.Token().Position, "invalid value for Long: "+val)
			}
			obj.Const = true
			return obj
		case "Integer$":
			obj, err := object.NewInteger(val, node.Token().Position)
			if err != nil {
				return object.NewError(node.Token().Position, "invalid value for Integer: "+val)
			}
			obj.Const = true
			return obj
		case "Single$":
			obj, err := object.NewSingle(val, node.Token().Position)
			if err != nil {
				return object.NewError(node.Token().Position, "invalid value for Single: "+val)
			}
			obj.Const = true
			return obj
		case "Double$":
			obj, err := object.NewDouble(val, node.Token().Position)
			if err != nil {
				return object.NewError(node.Token().Position, "invalid value for Double: "+val)
			}
			obj.Const = true
			return obj
		case "String$":
			obj := object.NewString(val, node.Token().Position)
			obj.Const = true
			return obj
		case "Date$":
			obj, err := object.NewDate(val, node.Token().Position)
			if err != nil {
				return object.NewError(node.Token().Position, "invalid value for Date: "+val)
			}
			obj.Const = true
			return obj
		case "Boolean$":
			val = strings.ToLower(val)
			if val == "true" {
				return object.NewBooleanByBool(true, node.Token().Position)
			} else if val == "false" {
				return object.NewBooleanByBool(false, node.Token().Position)
			} else {
				return object.NewError(node.Token().Position, "invalid value for Boolean: "+val)
			}
		case "Currency$":
			obj, err := object.NewCurrency(val, node.Token().Position)
			if err != nil {
				return object.NewError(node.Token().Position, "invalid value for Currency: "+val)
			}
			obj.Const = true
			return obj
		case "Variant$":
			obj, err := object.NewVariant(val, node.Token().Position)
			if err != nil {
				return object.NewError(node.Token().Position, "invalid value for Variant: "+val)
			}
			obj.Const = true
			return obj
		default:
			return object.NewError(node.Token().Position, "unknown type: "+typ.Name)
		}
	}
	return object.NewError(node.Token().Position, "unknown type: "+node.ConstType.String())
}

func evalEnumDecl(node *ast.EnumDecl, env *object.Environment) object.Object {
	// // define the enum values as constants
	obj := object.NewClass(node.Identifier.Name, node.Token().Position)
	obj.Members = make(map[string]object.Object)
	for _, v := range node.Values {
		member := object.NewUserDefined(v.Name, node, v.Token().Position)
		member.Const = true
		obj.Members[v.Name] = member
	}
	env.Set(node.Identifier.Name, obj)
	return nil
}

func evalTypeDef(node *ast.TypeDef, env *object.Environment) object.Object {
	// skip basic types
	// the "_" at the begining had to be added to break conflict between Date basic type and Date function

	switch strings.ToLower(node.TypeName.Name) {
	case "long$":
		// nothing to do
	case "integer$":
		// nothing to do
	case "single$":
		// nothing to do
	case "double$":
		// nothing to do
	case "string$":
		// nothing to do
	case "datetime$":
		// nothing to do
	case "boolean$":
		// nothing to do
	case "currency$":
		// nothing to do
	case "variant$":
		// nothing to do
	case "date$":
		// nothing to do
	default:
		return object.NewError(node.Token().Position, "unknown type: "+node.TypeName.Name)
	}
	return nil
}

func evalExprStmt(node *ast.ExprStmt, env *object.Environment) object.Object {
	return Eval(nil, node.Expression, env)
}

func evalEmptyStmt(node *ast.EmptyStmt, env *object.Environment) object.Object {
	return nil
}

func evalIfStmt(node *ast.IfStmt, env *object.Environment) object.Object {
	condition := Eval(nil, node.Condition, env)
	cond, ok := condition.(*object.Boolean)
	if !ok {
		return object.NewError(node.Token().Position, "condition is not a Boolean: "+node.Condition.String())
	}
	if cond.Value {
		return evalBody(node.Body, env)
	} else if node.ElseIf != nil {
		for _, stmt := range node.ElseIf {
			condition := Eval(nil, stmt.Condition, env)
			cond, ok := condition.(*object.Boolean)
			if !ok {
				return object.NewError(node.Token().Position, "condition is not a Boolean: "+stmt.Condition.String())
			}
			if cond.Value {
				return evalBody(stmt.Body, env)
			}
		}
	}
	if node.Else != nil {
		return evalBody(node.Else, env)
	}
	return nil
}

func evalWhileStmt(node *ast.WhileStmt, env *object.Environment) object.Object {
	for {
		condition := Eval(nil, node.Condition, env)
		cond, ok := condition.(*object.Boolean)
		if !ok {
			return object.NewError(node.Token().Position, "condition is not a Boolean: "+node.Condition.String())
		}
		if !cond.Value {
			break
		}
		result := evalBody(node.Body, env)
		if result != nil {
			switch result := result.(type) {
			case *object.Error:
				return result
			case *object.ReturnValue:
				return result
			case *object.Exit:
				if result.Kind == token.KwDo {
					return object.NOTHING
				}
			}
		}
	}
	return nil
}

func evalUntilStmt(node *ast.UntilStmt, env *object.Environment) object.Object {
	for {
		condition := Eval(nil, node.Condition, env)
		cond, ok := condition.(*object.Boolean)
		if !ok {
			return object.NewError(node.Token().Position, "condition is not a Boolean: "+node.Condition.String())
		}
		if cond.Value {
			break
		}
		result := evalBody(node.Body, env)
		if result != nil {
			switch result := result.(type) {
			case *object.Error:
				return result
			case *object.ReturnValue:
				return result
			case *object.Exit:
				if result.Kind == token.KwDo {
					return object.NOTHING
				}
			}
		}
	}
	return nil
}

func evalDoWhileStmt(node *ast.DoWhileStmt, env *object.Environment) object.Object {
	for {
		result := evalBody(node.Body, env)
		if result != nil {
			switch result := result.(type) {
			case *object.Error:
				return result
			case *object.ReturnValue:
				return result
			case *object.Exit:
				if result.Kind == token.KwDo {
					return object.NOTHING
				}
			}
		}
		condition := Eval(nil, node.Condition, env)
		cond, ok := condition.(*object.Boolean)
		if !ok {
			return object.NewError(node.Token().Position, "condition is not a Boolean: "+node.Condition.String())
		}
		if !cond.Value {
			break
		}

	}
	return nil
}

func evalDoUntilStmt(node *ast.DoUntilStmt, env *object.Environment) object.Object {
	for {
		result := evalBody(node.Body, env)
		if result != nil {
			switch result := result.(type) {
			case *object.Error:
				return result
			case *object.ReturnValue:
				return result
			case *object.Exit:
				if result.Kind == token.KwDo {
					return object.NOTHING
				}
			}
		}
		condition := Eval(nil, node.Condition, env)
		cond, ok := condition.(*object.Boolean)
		if !ok {
			return object.NewError(node.Token().Position, "condition is not a Boolean: "+node.Condition.String())
		}
		if cond.Value {
			break
		}

	}
	return nil
}

func evalForStmt(node *ast.ForStmt, env *object.Environment) object.Object {
	// initialize the loop variable
	forExpression := node.ForExpression
	var iterator object.Object
	switch forExpression := forExpression.(type) {
	case *ast.ForNextExpr:
		iterator = initForNextExpr(forExpression, env)
	case *ast.ForEachExpr:
		iterator = initForEachExpr(forExpression, env)
	}
	if iterator == nil {
		return object.NewError(node.Token().Position, "invalid loop variable: "+forExpression.String())
	} else {
		switch iterator := iterator.(type) {
		case *object.Error:
			return iterator
		case *object.Nothing:
			return iterator
		}
	}

	for {
		// execute the loop
		result := evalBody(node.Body, env)
		if result != nil {
			switch result := result.(type) {
			case *object.Error:
				return result
			case *object.ReturnValue:
				return result
			case *object.Exit:
				if result.Kind == token.KwFor {
					return object.NOTHING
				}
			}

			// evaluate the loop condition
			condition := Eval(nil, node.ForExpression, env)
			isFinished, ok := condition.(*object.Boolean)
			if !ok {
				return object.NewError(node.Token().Position, "condition is not a Boolean: "+node.ForExpression.String())
			}
			if isFinished.Value {
				break
			}
		}
		// execute increment
		var increment object.Object
		switch forExpression := forExpression.(type) {
		case *ast.ForNextExpr:
			increment = incrementForNextExpr(forExpression, env)
		case *ast.ForEachExpr:
			increment = incrementForEachExpr(forExpression, env)
		}
		if increment != nil {
			switch increment := increment.(type) {
			case *object.Error:
				return increment
			}
		}

	}
	return nil
}
func evalFinishedForNextExpr(node *ast.ForNextExpr, env *object.Environment) object.Object {
	// create comparaison node and evaluate it
	var stepValue any
	stepValue = 1
	step := node.Step
	if step != nil {
		val := Eval(nil, step, env)
		stepValue = val.GetValue()
	}
	var operation *ast.BinaryExpr
	var numStepValue float64
	switch stepValue := stepValue.(type) {
	case int32:
		numStepValue = float64(stepValue)
	case int64:
		numStepValue = float64(stepValue)
	case float32:
		numStepValue = float64(stepValue)
	case float64:
		numStepValue = stepValue
	case int:
		numStepValue = float64(stepValue)
	default:
		return object.NewError(node.Token().Position, "invalid step value: "+step.String())
	}

	if numStepValue > 0 {
		operation = &ast.BinaryExpr{Left: node.Variable, OpKind: token.Gt, OpToken: node.Token(), Right: node.To}
	} else {
		operation = &ast.BinaryExpr{Left: node.Variable, OpKind: token.Lt, OpToken: node.Token(), Right: node.To}
	}
	return Eval(nil, operation, env)
}

func initForNextExpr(node *ast.ForNextExpr, env *object.Environment) object.Object {
	// check if begin equals end
	compare := &ast.BinaryExpr{Left: node.From, OpKind: token.Eq, OpToken: node.Token(), Right: node.To}
	result := Eval(nil, compare, env)
	if result != nil {
		switch result := result.(type) {
		case *object.Error:
			return result
		case *object.Boolean:
			if result.Value {
				return object.NOTHING
			}
		}
	}

	// create assignment node and evaluate it
	operation := &ast.BinaryExpr{Left: node.Variable, OpKind: token.Assign, OpToken: node.Token(), Right: node.From}
	return Eval(nil, operation, env)
}

func incrementForNextExpr(node *ast.ForNextExpr, env *object.Environment) object.Object {
	// create increment node and evaluate it
	var operation *ast.BinaryExpr
	if node.Step == nil {
		operation = &ast.BinaryExpr{Left: node.Variable, OpKind: token.Add, OpToken: node.Token(), Right: &ast.BasicLit{Kind: token.LongLit, Value: "1", ValPos: node.Token()}}
	} else {
		operation = &ast.BinaryExpr{Left: node.Variable, OpKind: token.Add, OpToken: node.Token(), Right: node.Step}
	}
	assignment := &ast.BinaryExpr{Left: node.Variable, OpKind: token.Assign, OpToken: node.Token(), Right: operation}
	return Eval(nil, assignment, env)
}

func evalFinishedForEachExpr(node *ast.ForEachExpr, env *object.Environment) object.Object {
	// get iterator
	iterator, ok := env.Get(node.Variable.Name + "%")
	if !ok {
		return object.NewError(node.Token().Position, "undefined variable for each iterator: "+node.Variable.Name+"%")
	} else if iterator == nil {
		return object.NewError(node.Token().Position, "undefined variable for each iterator: "+node.Variable.Name+"%")
	}
	// get collection length
	var length int32
	val := Eval(nil, node.Collection, env)
	if val != nil {
		switch val := val.(type) {
		case *object.Error:
			return val
		case *object.Array:
			length = int32(len(val.Values))
		case *object.Variant:
			variant := val.Value
			switch variant := variant.(type) {
			case *object.Array:
				length = int32(len(variant.Values))
			}
		default:
			return object.NewError(node.Token().Position, "invalid type for loop variable: "+val.String())
		}
	}
	// compare the iterator to the collection length
	index := iterator.(*object.Integer)
	if index.Value < length {
		return object.NewBooleanByBool(false, node.Token().Position)
	}
	env.Delete(node.Variable.Name + "%")
	return object.NewBooleanByBool(true, node.Token().Position)
}

func initForEachExpr(node *ast.ForEachExpr, env *object.Environment) object.Object {
	// initialize the loop variable
	identifier := node.Variable
	val := Eval(nil, node.Collection, env)
	if val != nil {
		switch val := val.(type) {
		case *object.Error:
			return val
		case *object.Array:
			env.Set(identifier.Name, val.Values[0])
		case *object.Variant:
			variant := val.Value
			switch variant := variant.(type) {
			case *object.Array:
				env.Set(identifier.Name, variant.Values[0])
			}
		default:
			return object.NewError(node.Token().Position, "invalid type for loop variable: "+val.String())
		}
	}

	// we store an iterator variable identifier% in the environment
	env.Set(identifier.Name+"%", object.NewIntegerByInt(0, node.Token().Position))
	return val
}

func incrementForEachExpr(node *ast.ForEachExpr, env *object.Environment) object.Object {
	val := Eval(nil, node.Collection, env)

	// get array
	var array *object.Array
	switch val := val.(type) {
	case *object.Error:
		return val
	case *object.Array:
		array = val
	case *object.Variant:
		variant := val.Value
		switch variant := variant.(type) {
		case *object.Array:
			array = variant
		}
	default:
		return object.NewError(node.Token().Position, "invalid type for loop variable: "+val.String())
	}
	// find current index
	i, ok := env.Get(node.Variable.Name + "%")
	if !ok {
		return object.NewError(node.Token().Position, "undefined variable for each iterator: "+node.Variable.Name+"%")
	}

	// increment index
	index := i.(*object.Integer)
	index.Value++
	if index.Value < int32(len(array.Values)) {
		env.Set(node.Variable.Name, array.Values[uint32(index.Value)])
		return array.Values[uint32(index.Value)]
	}
	return nil
}

func evalSelectStmt(node *ast.SelectStmt, env *object.Environment) object.Object {
	// evaluate the select expression
	expression := Eval(nil, node.Condition, env)
	if expression == nil {
		return object.NewError(node.Token().Position, "invalid select expression: "+node.Condition.String())
	} else if expression.Type() == object.ERROR_OBJ {
		return expression
	}

	// evaluate the case statements
	for _, caseStmt := range node.Body {
		// case else branch
		if caseStmt.Condition == nil {
			return evalBody(caseStmt.Body, env)
		}
		// evaluate the case statement condition
		condition := Eval(nil, caseStmt.Condition, env)
		if condition == nil {
			return object.NewError(node.Token().Position, "invalid case expression: "+caseStmt.Condition.String())
		} else if condition.Type() == object.ERROR_OBJ {
			return condition
		}
		// compare the case statement to the select expression
		if expression.Equals(condition) {
			// execute the case statement
			return evalBody(caseStmt.Body, env)
		}
	}

	return nil
}

func evalCallSelectorExpr(node *ast.CallSelectorExpr, env *object.Environment) object.Object {
	// root is either an built-in class or an enum

	root := Eval(nil, node.Root, env)
	if root.Type() == object.ERROR_OBJ {
		return root
	}
	if root.Type() == object.CLASS_OBJ {
		class := root.(*object.Class)
		switch selector := node.Selector.(type) {
		case *ast.Identifier:
			// check if object is a user defined class
			value := class.Members[selector.Name]
			if value != nil {
				return value
			}
			// try in the built-in class
			selectorIdentifier := object.NewString(selector.Name, selector.Token().Position)
			return rtlib.EvalClassProperty(class, selectorIdentifier, env)
		case *ast.CallOrIndexExpr:
			// this will call ftlib.EvalClassMethod()
			return evalCallOrIndexExpr(class, selector, env)
		}
	}
	return object.NewError(node.Token().Position, "unknown call selector expression: "+node.String())

}

func evalBinaryExpr(node *ast.BinaryExpr, env *object.Environment) object.Object {
	// check if one of the operands is a Boolean
	// go for lazy evaluation

	left := Eval(nil, node.Left, env)
	var leftBool *object.Boolean
	if left.Type() == object.ERROR_OBJ {
		return left
	}

	// check if left is a function definition
	// on an assignment. If so, we'll evaluate the result instead
	// of the function definition
	if left.Type() == object.FUNCTION_OBJ {
		if node.OpKind == token.Assign {
			resultNode := &ast.Identifier{Name: returnName, Tok: node.Token()}
			left = Eval(nil, resultNode, env)
			if left.Type() == object.ERROR_OBJ {
				return left
			}
		}
	}

	// check if left is a variant
	if left.Type() == object.VARIANT_OBJ {
		if left.(*object.Variant).Value.Type() == object.BOOLEAN_OBJ {
			leftBool = left.(*object.Variant).Value.(*object.Boolean)
		}
	} else {
		leftBool, _ = left.(*object.Boolean)
	}
	if leftBool != nil && node.OpKind != token.Assign {
		return evalBooleanBinaryExpr(node, env, leftBool)
	}

	// if not assignment, convert left variant to its basic type
	if node.OpKind != token.Assign && left.Type() == object.VARIANT_OBJ {
		left = left.(*object.Variant).Value
	}

	// evaluate right operand and convert variant to its basic type
	right := Eval(nil, node.Right, env)
	if right.Type() == object.ERROR_OBJ {
		return right
	}
	if right.Type() == object.VARIANT_OBJ {
		right = right.(*object.Variant).Value
	}
	// process enum types first
	if left.Type() == object.USERDEF_OBJ {
		leftEnum := left.(*object.UserDefined)
		switch right := right.(type) {
		case *object.UserDefined:
			if leftEnum.Decl.Identifier.Name != right.Decl.Identifier.Name {
				return object.NewError(node.Token().Position, "mismatched enum types: "+leftEnum.Decl.Identifier.Name+" and "+right.Decl.Identifier.Name)
			}
			return evalEnumBinaryExpr(node, env, leftEnum, right)
		default:
			return object.NewError(node.Token().Position, "right operand is not an enum or string: "+node.Right.String())
		}
	}

	// get precision of the operands
	leftBasicType := object.GetBasicType(left.Type())
	rightBasicType := object.GetBasicType(right.Type())
	// get the result precision
	leftPrecision := ast.GetPrecisionOrder(leftBasicType)
	rightPrecision := ast.GetPrecisionOrder(rightBasicType)

	resultPrecision := rightPrecision
	if node.OpKind != token.Assign {
		if rightPrecision < leftPrecision {
			resultPrecision = leftPrecision
		}
	}
	// convert left and right to the same type
	switch resultPrecision {
	case ast.LongPrecision:
		// convert  right to Long
		rightLong, err := object.NewLong(right.String(), node.Token().Position)
		if err != nil {
			return object.NewError(node.Token().Position, "right operand is not a number: "+node.Right.String())
		}
		return evalLongBinaryExpr(node, env, left, rightLong.Value)
	case ast.IntegerPrecision:
		// convert  right to Integer
		rightInt, err := object.NewInteger(right.String(), node.Token().Position)
		if err != nil {
			return object.NewError(node.Token().Position, "right operand is not a number: "+node.Right.String())
		}
		return evalLongBinaryExpr(node, env, left, int64(rightInt.Value))
	case ast.SinglePrecision:
		// convert  right to Single
		rightSingle, err := object.NewSingle(right.String(), node.Token().Position)
		if err != nil {
			return object.NewError(node.Token().Position, "right operand is not a number: "+node.Right.String())
		}
		return evalDoubleBinaryExpr(node, env, left, float64(rightSingle.Value))
	case ast.DoublePrecision:
		// convert  right to Double
		rightDouble, err := object.NewDouble(right.String(), node.Token().Position)
		if err != nil {
			return object.NewError(node.Token().Position, "right operand is not a number: "+node.Right.String())
		}
		return evalDoubleBinaryExpr(node, env, left, rightDouble.Value)
	case ast.CurrencyPrecision:
		// convert right to Currency
		rightCurrency, err := object.NewCurrency(right.String(), node.Token().Position)
		if err != nil {
			return object.NewError(node.Token().Position, "right operand is not a number: "+node.Right.String())
		}
		return evalDoubleBinaryExpr(node, env, left, rightCurrency.Value)
	case ast.DatePrecision:
		// convert right to Date
		rightDate, err := object.NewDate(right.String(), node.Token().Position)
		if err != nil {
			return object.NewError(node.Token().Position, "right operand is not a date: "+node.Right.String())
		}
		return evalDateBinaryExpr(node, env, left, rightDate.Value)
	case ast.StringPrecision:
		// convert right to String
		rightString := object.NewString(right.String(), node.Token().Position)
		return evalStringBinaryExpr(node, env, left, rightString.Value)
	case ast.BooleanPrecision:
		return evalBooleanBinaryExpr(node, env, left)
	default:
		return object.NewError(node.Token().Position, "unknown binary expression: "+node.String())
	}
}

func evalLongBinaryExpr(node *ast.BinaryExpr, env *object.Environment, leftObject object.Object, rightValue int64) object.Object {
	// TODO: implement value overflow

	// convert left to Long if it is not an assignment
	var leftValue int64
	if node.OpKind != token.Assign {
		switch leftObject := leftObject.(type) {
		case *object.Long:
			leftValue = leftObject.Value
		case *object.Integer:
			leftValue = int64(leftObject.Value)
		case *object.Variant:
			variant := leftObject.Value
			switch variant := variant.(type) {
			case *object.Long:
				leftValue = variant.Value
			case *object.Integer:
				leftValue = int64(variant.Value)
			default:
				return object.NewError(node.Token().Position, "left operand is not a number: "+variant.String())
			}
		default:
			return object.NewError(node.Token().Position, "left operand is not a number: "+leftObject.String())
		}
	}

	switch node.OpKind {
	case token.Add:
		return object.NewLongByInt(leftValue+rightValue, node.Token().Position)
	case token.Minus:
		return object.NewLongByInt(leftValue-rightValue, node.Token().Position)
	case token.Mul:
		return object.NewLongByInt(leftValue*rightValue, node.Token().Position)
	case token.Div:
		if rightValue == 0 {
			return object.NewError(node.Token().Position, "division by zero")
		}
		return object.NewLongByInt(leftValue/rightValue, node.Token().Position)
	case token.Mod:
		return object.NewLongByInt(leftValue%rightValue, node.Token().Position)
	case token.IntDiv:
		if rightValue == 0 {
			return object.NewError(node.Token().Position, "division by zero")
		}
		return object.NewLongByInt(leftValue/rightValue, node.Token().Position)
	case token.Exponent:
		return object.NewLongByInt(int64(math.Pow(float64(leftValue), float64(rightValue))), node.Token().Position)
	case token.Eq:
		return object.NewBooleanByBool(leftValue == rightValue, node.Token().Position)
	case token.Neq:
		return object.NewBooleanByBool(leftValue != rightValue, node.Token().Position)
	case token.Lt:
		return object.NewBooleanByBool(leftValue < rightValue, node.Token().Position)
	case token.Le:
		return object.NewBooleanByBool(leftValue <= rightValue, node.Token().Position)
	case token.Gt:
		return object.NewBooleanByBool(leftValue > rightValue, node.Token().Position)
	case token.Ge:
		return object.NewBooleanByBool(leftValue >= rightValue, node.Token().Position)
	case token.Assign:
		if leftObject.IsConstant() {
			return object.NewError(node.Token().Position, "constant cannot be assigned: "+leftObject.String())
		}

		switch l := leftObject.(type) {
		case *object.Long:
			l.Value = rightValue
			return l
		case *object.Integer:
			l.Value = int32(rightValue)
			return l
		case *object.Single:
			l.Value = (float32)(rightValue)
			return l
		case *object.Double:
			l.Value = (float64)(rightValue)
			return l
		case *object.Currency:
			l.Value = (float64)(rightValue)
			return l
		case *object.Variant:
			l.Value = object.NewLongByInt(rightValue, node.Token().Position)
		}
	default:
		return object.NewError(node.Token().Position, "unknown binary expression: "+node.String())
	}
	return nil
}

func evalDoubleBinaryExpr(node *ast.BinaryExpr, env *object.Environment, leftObject object.Object, rightValue float64) object.Object {
	// TODO: implement value overflow
	var leftValue float64
	if leftObject.Type() == object.VARIANT_OBJ {
		variant := leftObject.(*object.Variant).Value
		if variant.Type() == object.DOUBLE_OBJ {
			leftValue = variant.(*object.Double).Value
		} else if variant.Type() == object.SINGLE_OBJ {
			leftValue = float64(variant.(*object.Single).Value)
		} else if variant.Type() == object.CURRENCY_OBJ {
			leftValue = variant.(*object.Currency).Value
		} else if variant.Type() == object.LONG_OBJ {
			leftValue = float64(variant.(*object.Long).Value)
		} else if variant.Type() == object.INTEGER_OBJ {
			leftValue = float64(variant.(*object.Integer).Value)
		} else if node.OpKind != token.Assign {
			return object.NewError(node.Token().Position, "left operand is not a number: "+variant.String())
		}
	} else if leftObject.Type() == object.DOUBLE_OBJ {
		leftValue = leftObject.(*object.Double).Value
	} else if leftObject.Type() == object.SINGLE_OBJ {
		leftValue = float64(leftObject.(*object.Single).Value)
	} else if leftObject.Type() == object.CURRENCY_OBJ {
		leftValue = leftObject.(*object.Currency).Value
	} else if leftObject.Type() == object.LONG_OBJ {
		leftValue = float64(leftObject.(*object.Long).Value)
	} else if leftObject.Type() == object.INTEGER_OBJ {
		leftValue = float64(leftObject.(*object.Integer).Value)
	} else {
		return object.NewError(node.Token().Position, "left operand is not a number: "+leftObject.String())
	}

	switch node.OpKind {
	case token.Add:
		return &object.Double{Value: leftValue + rightValue}
	case token.Minus:
		return &object.Double{Value: leftValue - rightValue}
	case token.Mul:
		return &object.Double{Value: leftValue * rightValue}
	case token.Div:
		if rightValue == 0 {
			return object.NewError(node.Token().Position, "division by zero")
		}
		return &object.Double{Value: leftValue / rightValue}
	case token.IntDiv:
		if rightValue == 0 {
			return object.NewError(node.Token().Position, "division by zero")
		}
		return object.NewLongByInt((int64)(leftValue/rightValue), node.Token().Position)
	case token.Exponent:
		return &object.Double{Value: math.Pow(leftValue, rightValue)}
	case token.Eq:
		return object.NewBooleanByBool(leftValue == rightValue, node.Token().Position)
	case token.Neq:
		return object.NewBooleanByBool(leftValue != rightValue, node.Token().Position)
	case token.Lt:
		return object.NewBooleanByBool(leftValue < rightValue, node.Token().Position)
	case token.Le:
		return object.NewBooleanByBool(leftValue <= rightValue, node.Token().Position)
	case token.Gt:
		return object.NewBooleanByBool(leftValue > rightValue, node.Token().Position)
	case token.Ge:
		return object.NewBooleanByBool(leftValue >= rightValue, node.Token().Position)
	case token.Assign:
		if leftObject.IsConstant() {
			return object.NewError(node.Token().Position, "constant cannot be assigned: "+leftObject.String())
		}

		switch left := leftObject.(type) {
		case *object.Double:
			left.Value = rightValue
			return left
		case *object.Single:
			left.Value = (float32)(rightValue)
			return left
		case *object.Currency:
			left.Value = (float64)(rightValue)
			return left
		case *object.Variant:
			left.Value = object.NewDoubleByFloat(rightValue, node.Token().Position)
			return left
		default:
			return object.NewError(node.Token().Position, "numberic conversion error: "+left.String())
		}
	default:
		return object.NewError(node.Token().Position, "unknown binary expression: "+node.String())
	}
}

func evalBooleanBinaryExpr(node *ast.BinaryExpr, env *object.Environment, leftObject object.Object) object.Object {
	var leftValue bool
	if leftObject.Type() == object.VARIANT_OBJ {
		variant := leftObject.(*object.Variant)
		if variant.Value.Type() == object.BOOLEAN_OBJ {
			leftValue = variant.Value.(*object.Boolean).Value
		} else if node.OpKind != token.Assign {
			return object.NewError(node.Token().Position, "left operand is not a Boolean: "+leftObject.String())
		}
	} else if leftObject.Type() == object.BOOLEAN_OBJ {
		left, ok := leftObject.(*object.Boolean)
		if !ok {
			return object.NewError(node.Token().Position, "left operand is not a Boolean: "+leftObject.String())
		}
		leftValue = left.Value
	} else {
		return object.NewError(node.Token().Position, "left operand is not a Boolean: "+leftObject.String())
	}

	switch node.OpKind {
	case token.And:
		if leftValue {
			right := Eval(nil, node.Right, env)
			if isError(right) {
				return right
			}
			if right.Type() == object.VARIANT_OBJ {
				right = right.(*object.Variant).Value
			}
			rightBool, ok := right.(*object.Boolean)
			if !ok {
				return object.NewError(node.Token().Position, "right operand is not a Boolean: "+node.Right.String())
			}
			if rightBool.Value {
				return object.NewBooleanByBool(true, node.Token().Position)
			}
		}
		return object.NewBooleanByBool(false, node.Token().Position)
	case token.Or:
		if leftValue {
			return object.NewBooleanByBool(true, node.Token().Position)
		}

		right := Eval(nil, node.Right, env)
		if isError(right) {
			return right
		}
		if right.Type() == object.VARIANT_OBJ {
			right = right.(*object.Variant).Value
		}
		rightBool, ok := right.(*object.Boolean)
		if !ok {
			return object.NewError(node.Token().Position, "right operand is not a Boolean: "+node.Right.String())
		}
		return object.NewBooleanByBool(rightBool.Value, node.Token().Position)
	case token.Eq:
		right := Eval(nil, node.Right, env)
		if isError(right) {
			return right
		}
		if right.Type() == object.VARIANT_OBJ {
			right = right.(*object.Variant).Value
		}
		rightBool, ok := right.(*object.Boolean)
		if !ok {
			return object.NewError(node.Token().Position, "right operand is not a Boolean: "+node.Right.String())
		}
		return object.NewBooleanByBool(leftValue == rightBool.Value, node.Token().Position)
	case token.Neq:
		right := Eval(nil, node.Right, env)
		if isError(right) {
			return right
		}
		if right.Type() == object.VARIANT_OBJ {
			right = right.(*object.Variant).Value
		}

		rightBool, ok := right.(*object.Boolean)
		if !ok {
			return object.NewError(node.Token().Position, "right operand is not a Boolean: "+node.Right.String())
		}
		return object.NewBooleanByBool(leftValue != rightBool.Value, node.Token().Position)
	case token.Assign:

		// check if left operand is a constant
		if leftObject.IsConstant() {
			return object.NewError(node.Token().Position, "cannot assign a value to a constant: "+leftObject.String())
		}

		right := Eval(nil, node.Right, env)
		if isError(right) {
			return right
		}
		// validate the type
		if right.Type() == object.VARIANT_OBJ {
			variant := right.(*object.Variant)
			if variant.Value.Type() != object.BOOLEAN_OBJ {
				return object.NewError(node.Token().Position, "right operand is not a Boolean: "+node.Right.String())
			}
			right = variant.Value
		}

		rightValue, ok := right.(*object.Boolean)
		if !ok {
			return object.NewError(node.Token().Position, "right operand is not a Boolean: "+node.Right.String())
		} else {
			switch left := leftObject.(type) {
			case *object.Boolean:
				left.Value = rightValue.Value
			case *object.Variant:
				// we don't need to test prior type to assing
				left.Value = object.NewBooleanByBool(rightValue.Value, node.Token().Position)
			}
		}
		return leftObject
	default:
		return object.NewError(node.Token().Position, "unknown binary expression: "+node.String())
	}
}

func evalUnaryExpr(node *ast.UnaryExpr, env *object.Environment) object.Object {
	switch node.OpKind {
	case token.Not:
		right := Eval(nil, node.Right, env)
		if isError(right) {
			return right
		}
		if right.Type() == object.VARIANT_OBJ {
			right = right.(*object.Variant).Value
		}
		switch right := right.(type) {
		case *object.Boolean:
			return object.NewBooleanByBool(!right.Value, node.Token().Position)
		default:
			return object.NewError(node.Token().Position, "right operand is not a Boolean: "+node.Right.String())
		}
	case token.Minus:
		// create multiply by -1 node and evaluate it
		minusOne := &ast.BasicLit{Kind: token.LongLit, ValPos: node.OpToken, Value: "-1"}
		multiply := &ast.BinaryExpr{Left: minusOne, OpKind: token.Mul, OpToken: node.OpToken, Right: node.Right}
		assignment := &ast.BinaryExpr{Left: node.Right, OpKind: token.Assign, OpToken: node.OpToken, Right: multiply}
		return Eval(nil, assignment, env)
	}
	return object.NewError(node.Token().Position, "unknown unary expression: "+node.String())
}

func evalStringBinaryExpr(node *ast.BinaryExpr, env *object.Environment, leftObject object.Object, right string) object.Object {
	var leftValue string
	if leftObject.Type() == object.VARIANT_OBJ {
		variant := leftObject.(*object.Variant)
		if variant.Value.Type() == object.STRING_OBJ {
			leftValue = variant.Value.(*object.String).Value
		} else if node.OpKind != token.Assign {
			return object.NewError(node.Token().Position, "left operand is not a String: "+leftObject.String())
		}
	} else if leftObject.Type() == object.STRING_OBJ {
		leftValue = leftObject.(*object.String).Value
	} else {
		return object.NewError(node.Token().Position, "left operand is not a String: "+leftObject.String())
	}

	switch node.OpKind {
	case token.Concat:
		value := leftValue + right
		return object.NewString(value, node.Token().Position)
	case token.Eq:
		return object.NewBooleanByBool(leftValue == right, node.Token().Position)
	case token.Neq:
		return object.NewBooleanByBool(leftValue != right, node.Token().Position)
	case token.Lt:
		return object.NewBooleanByBool(leftValue < right, node.Token().Position)
	case token.Le:
		return object.NewBooleanByBool(leftValue <= right, node.Token().Position)
	case token.Gt:
		return object.NewBooleanByBool(leftValue > right, node.Token().Position)
	case token.Ge:
		return object.NewBooleanByBool(leftValue >= right, node.Token().Position)
	case token.Assign:
		if leftObject.IsConstant() {
			return object.NewError(node.Token().Position, "cannot assign a value to a constant: "+leftObject.String())
		}
		switch left := leftObject.(type) {
		case *object.String:
			left.Value = right
			return left
		case *object.Variant:
			left.Value = object.NewString(right, node.Token().Position)
			return left
		}
	default:
		return object.NewError(node.Token().Position, "unknown binary expression: "+node.String())
	}
	return object.NewError(node.Token().Position, "unknown binary expression: "+node.String())
}

func evalDateBinaryExpr(node *ast.BinaryExpr, env *object.Environment, leftObject object.Object, right time.Time) object.Object {
	var leftValue time.Time
	if leftObject.Type() == object.VARIANT_OBJ {
		variant := leftObject.(*object.Variant)
		if variant.Value.Type() == object.DATE_OBJ {
			leftValue = variant.Value.(*object.Date).Value
		} else if node.OpKind != token.Assign {
			return object.NewError(node.Token().Position, "left operand is not a Date: "+leftObject.String())
		}
	} else if leftObject.Type() == object.DATE_OBJ {
		leftValue = leftObject.(*object.Date).Value
	} else {
		return object.NewError(node.Token().Position, "left operand is not a Date: "+leftObject.String())
	}

	switch node.OpKind {
	case token.Eq:
		return object.NewBooleanByBool(time.Time.Equal(leftValue, right), node.Token().Position)
	case token.Neq:
		return object.NewBooleanByBool(!time.Time.Equal(leftValue, right), node.Token().Position)
	case token.Lt:
		return object.NewBooleanByBool(leftValue.Before(right), node.Token().Position)
	case token.Le:
		return object.NewBooleanByBool(leftValue.Before(right) || time.Time.Equal(leftValue, right), node.Token().Position)
	case token.Gt:
		return object.NewBooleanByBool(leftValue.After(right), node.Token().Position)
	case token.Ge:
		return object.NewBooleanByBool(leftValue.After(right) || time.Time.Equal(leftValue, right), node.Token().Position)
	case token.Assign:
		// verify if left is a constant
		if leftObject.IsConstant() {
			return object.NewError(node.Token().Position, "cannot assign a value to a constant: "+leftObject.String())
		}
		switch left := leftObject.(type) {
		case *object.Date:
			left.Value = right
			return left
		case *object.Variant:
			left.Value = object.NewDateByTime(right, node.Token().Position)
			return left
		}
	default:
		return object.NewError(node.Token().Position, "unknown binary expression: "+node.String())
	}
	return object.NewError(node.Token().Position, "unknown binary expression: "+node.String())
}

func evalEnumBinaryExpr(node *ast.BinaryExpr, env *object.Environment, leftObject object.Object, right object.Object) object.Object {
	var leftValue string
	var leftType string
	var rightValue string
	var rightType string

	if leftObject.Type() == object.VARIANT_OBJ {
		variant := leftObject.(*object.Variant)
		if variant.Value.Type() == object.USERDEF_OBJ {
			leftValue = variant.Value.(*object.UserDefined).Value
			leftType = variant.Value.(*object.UserDefined).Decl.Identifier.Name
		} else if node.OpKind != token.Assign {
			return object.NewError(node.Token().Position, "left operand is not an enum: "+leftObject.String())
		}
	} else if leftObject.Type() == object.USERDEF_OBJ {
		leftValue = leftObject.(*object.UserDefined).Value
		leftType = leftObject.(*object.UserDefined).Decl.Identifier.Name
	} else {
		return object.NewError(node.Token().Position, "left operand is not an enum: "+leftObject.String())
	}
	// get the right value and type
	if right.Type() == object.VARIANT_OBJ {
		variant := right.(*object.Variant)
		if variant.Value.Type() == object.USERDEF_OBJ {
			rightValue = variant.Value.(*object.UserDefined).Value
			rightType = variant.Value.(*object.UserDefined).Decl.Identifier.Name
		} else {
			return object.NewError(node.Token().Position, "right operand is not an enum: "+right.String())
		}
	} else if right.Type() == object.USERDEF_OBJ {
		rightValue = right.(*object.UserDefined).Value
		rightType = right.(*object.UserDefined).Decl.Identifier.Name
	} else {
		return object.NewError(node.Token().Position, "right operand is not an enum: "+right.String())
	}

	// validate types
	if leftType != rightType {
		return object.NewError(node.Token().Position, "mismatched enum types: "+leftType+" and "+rightType)
	}

	switch node.OpKind {
	case token.Eq:
		return object.NewBooleanByBool(leftValue == rightValue, node.Token().Position)
	case token.Neq:
		return object.NewBooleanByBool(leftValue != rightValue, node.Token().Position)
	case token.Assign:
		// check if leftObject is a constant
		if leftObject.IsConstant() {
			return object.NewError(node.Token().Position, "cannot assign a value to a constant: "+leftObject.String())
		}
		switch left := leftObject.(type) {
		case *object.UserDefined:
			left.Value = rightValue
			return left
		case *object.Variant:
			left.Value = object.NewUserDefined(rightValue, leftObject.(*object.UserDefined).Decl, node.Token().Position)
			left.Value.(*object.UserDefined).Value = rightValue
			return left
		}
	default:
		return object.NewError(node.Token().Position, "unknown binary expression: "+node.String())
	}
	return object.NewError(node.Token().Position, "unknown binary expression: "+node.String())
}

func evalCallOrIndexExpr(class *object.Class, node *ast.CallOrIndexExpr, env *object.Environment) object.Object {
	target := node.Identifier.Decl
	switch target := target.(type) {
	case *ast.ParamItem:
		array := Eval(class, node.Identifier, env)
		if isError(array) {
			return array
		}
		// get index values
		indexes := evalExpressions(node.Args, env)
		if len(indexes) == 1 && isError(indexes[0]) {
			return indexes[0]
		}
		return evalArrayIndex(array, indexes, node.Args)

	case *ast.SubDecl:
		sub := Eval(class, node.Identifier, env)
		if isError(sub) {
			return sub
		}
		// get parameters values
		values := evalExpressions(node.Args, env)
		if len(values) == 1 && isError(values[0]) {
			return values[0]
		}
		return applyFunction(class, sub, values, node.Args)
	case *ast.FuncDecl:
		function := Eval(class, node.Identifier, env)
		if isError(function) {
			return function
		}
		// get parameters values
		values := evalExpressions(node.Args, env)
		if len(values) == 1 && isError(values[0]) {
			return values[0]
		}
		return applyFunction(class, function, values, node.Args)
	case *ast.ArrayDecl:
		array := Eval(class, node.Identifier, env)
		if isError(array) {
			return array
		}
		// get index values
		indexes := evalExpressions(node.Args, env)
		if len(indexes) == 1 && isError(indexes[0]) {
			return indexes[0]
		}
		return evalArrayIndex(array, indexes, node.Args)
	default:
		return object.NewError(node.Token().Position, "unknown declaration: "+target.String())
	}
}

func evalCallSubStmt(class *object.Class, node *ast.CallSubStmt, env *object.Environment) object.Object {
	sub := Eval(class, node.Definition, env)
	// don't evaluate the body if the sub is an error
	// we'll return it anyway
	if sub.Type() == object.EXIT_OBJ {
		return nil
	}
	return sub
}

func evalArrayIndex(array object.Object, indexes []object.Object, args []ast.Expression) object.Object {
	switch array := array.(type) {
	case *object.Array:
		// check if number of indexes is equal to the number of dimensions
		if len(indexes) != len(array.Dimensions) {
			return object.NewError(args[0].Token().Position, "mismatched number of indexes: "+(string)(len(indexes))+" and "+(string)(len(array.Dimensions)))
		}
		return array.GetValueAt(indexes)
	default:
		return object.NewError(token.Position{Line: 0, Column: 0, Absolute: 0}, "not an array: "+(string)(array.Type()))
	}
}

func applyFunction(class *object.Class, fn object.Object, values []object.Object, args []ast.Expression) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv, err := extendFunctionEnv(fn, &values)
		if err != nil {
			return object.NewError(args[0].Token().Position, err.Error())
		}
		// do we have a built-in library call?
		var val object.Object
		if fn.Body == nil {
			if class != nil {
				val = rtlib.EvalClassMethod(class, fn, extendedEnv)
			} else {
				val = rtlib.EvalBody(fn, values, extendedEnv)
			}
		} else {
			val = evalBody(fn.Body, extendedEnv)
		}
		if isError(val) {
			return val
		}

		// if normal exit, wrap result in ReturnValue
		switch val := val.(type) {
		case *object.ReturnValue:
			// do nothing
		default:
			val = object.NewReturnValue(val, fn.Position())
		}

		return unwrapReturnValue(val)
	case *object.Sub:
		extendedEnv, err := extendFunctionEnv(fn, &values)
		if err != nil {
			return object.NewError(args[0].Token().Position, err.Error())
		}
		// do we have a built-in library call?
		var val object.Object
		if fn.Body == nil {
			val = rtlib.EvalBody(fn, values, extendedEnv)
		} else {
			val = evalBody(fn.Body, extendedEnv)
		}
		if isError(val) {
			return val
		}

		return unwrapReturnValue(val)
	default:
		return object.NewError(fn.Position(), "not a function: "+(string)(fn.Type()))
	}
}

func extendFunctionEnv(fn object.Object, args *[]object.Object) (*object.Environment, error) {
	var extendedEnv *object.Environment
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv = fn.Env.Extend()
		err := assignValueToParameters(fn.Parameters.Params, args, extendedEnv)
		if err != nil {
			return nil, err
		}
		extendedEnv.Set(returnName, object.NewEmptyByKind(fn.Definition.FuncType.Result.Token().Kind, fn.Position()))
	case *object.Sub:
		extendedEnv = fn.Env.Extend()
		err := assignValueToParameters(fn.Parameters.Params, args, extendedEnv)
		if err != nil {
			return nil, err
		}

	}
	return extendedEnv, nil
}

func assignValueToParameters(params []ast.ParamItem, values *[]object.Object, env *object.Environment) error {
	// three cases:
	// 1. exact number of parameters and values
	// 2. number of parameters is less than the number of values
	// 3. number of parameters is greater than the number of values

	// 1. exact number of parameters and values
	if len(params) == len(*values) {
		for i, param := range params {
			if param.ByVal {
				// create a copy of the value
				env.Set(param.VarName.Name, (*values)[i].Copy())
			} else {
				env.Set(param.VarName.Name, (*values)[i])
			}
		}
		return nil
	}
	// 2. number of parameters is less than the number of values
	if len(params) < len(*values) {
		for i, param := range params {
			if i < len(params) {
				if param.ByVal {
					// create a copy of the value
					env.Set(param.VarName.Name, (*values)[i].Copy())
				} else {
					env.Set(param.VarName.Name, (*values)[i])
				}
			} else {
				// check if last parameter is an array
				if param.ParamArray {
					// calculate the number of elements
					dimensions := []uint32{uint32(len(*values) - len(params))}
					// create an array
					arr := object.NewArray(object.KindToType(param.VarType.Token().Kind), dimensions, param.Token().Position)
					env.Set(param.VarName.Name, arr)
					// assign the rest of the values to the array
					for j := i; j < len(*values); j++ {
						// get the index
						index := j - i
						// build index array
						indexes := make([]uint32, 1)
						indexes[0] = uint32(index)
						// assign the value to the array
						arr.Set(indexes, (*values)[j])
					}
					return nil
				}
			}
		}
		return nil
	}
	// 3. number of parameters is greater than the number of values
	for i, param := range params {
		if i < len(*values) {
			if param.ByVal {
				// create a copy of the value
				env.Set(param.VarName.Name, (*values)[i].Copy())
			} else {
				env.Set(param.VarName.Name, (*values)[i])
			}
		} else if param.Optional {
			// set default value
			value := Eval(nil, param.DefaultValue, env)
			env.Set(param.VarName.Name, value)
			// add parameter to values
			*values = append(*values, value)
		}
	}
	return nil
}

func evalExpressions(args []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object
	for _, arg := range args {
		val := Eval(nil, arg, env)
		if isError(val) {
			return []object.Object{val}
		}
		result = append(result, val)
	}
	return result
}

func unwrapReturnValue(obj object.Object) object.Object {
	if obj != nil {
		switch obj := obj.(type) {
		case *object.ReturnValue:
			return obj.Value
		}
	}
	return obj
}

func evalSpecialStmt(node *ast.SpecialStmt, env *object.Environment) object.Object {
	kind := node.Keyword1.Kind
	var val object.Object
	val = object.NOTHING // default value
	switch kind {
	case token.KwLet:
		// evaluate the expression
		val = Eval(nil, node.Args[0], env)
	case token.KwRedim:
		// only evaluate the array index expression
		// expect a node of type CallOrIndexExpr
		arg, ok := node.Args[0].(*ast.CallOrIndexExpr)
		if !ok {
			return object.NewError(node.Keyword1.Position, "argument must be an array")
		}
		// read the array dimensions
		dim := arg.Args[0]
		// evaluate the dimension
		value := Eval(nil, dim, env)
		params := make([]object.Object, 0)
		params = append(params, value)

		// check if the array is preserved
		preserved := false
		// check node if preserve is set
		if strings.EqualFold(node.Keyword2, "preserve") {
			preserved = true
		}
		// check if the first parameter. It contains the dimension to redim
		param0 := params[0]
		if params[0].Type() == object.VARIANT_OBJ {
			param0 = params[0].(*object.Variant).Value
		}
		paramLong, ok := param0.(*object.Long)
		if !ok {
			return object.NewError(params[0].Position(), "dimension must be a long")
		}
		// array name is in the node
		arrayExpr, ok := node.Args[0].(*ast.CallOrIndexExpr)
		if !ok {
			return object.NewError(node.Keyword1.Position, "invalid array name")
		}
		name := arrayExpr.Identifier.Name
		// find the array in the environment
		obj, ok := env.Get(name)
		if !ok {
			return object.NewError(node.Keyword1.Position, "array not found: "+name)
		}
		arr, ok := obj.(*object.Array)
		if !ok {
			return object.NewError(node.Keyword1.Position, "not an array: "+name)
		}
		arr.Redimension(preserved, uint32(paramLong.Value))
		return arr
	default:
		params := make([]object.Object, 0)
		for _, expr := range node.Args {
			val = Eval(nil, expr, env)
			params = append(params, val)
		}
		val = rtlib.EvalSpecialStatement(node, params, env)
	}
	return val
}

func evalExitStmt(node *ast.ExitStmt, env *object.Environment) object.Object {
	switch node.ExitType.Kind {
	case token.KwSub, token.KwDo, token.KwFor:
		return object.NewExit(node.ExitType.Kind, node.Token().Position)
	case token.KwFunction:
		value, ok := env.Get(returnName)
		if ok {
			return object.NewReturnValue(value, node.Token().Position)
		}
		return object.NewError(node.Token().Position, "undefined function return value: ")
	}
	return object.NewReturnValue(object.NOTHING, node.Token().Position)
}

func evalParenExpr(node *ast.ParenExpr, env *object.Environment) object.Object {
	return Eval(nil, node.Expr, env)
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
