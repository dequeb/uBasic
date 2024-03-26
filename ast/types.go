package ast

import (
	"strconv"
	"uBasic/errors"
	"uBasic/token"
	"uBasic/types"
)

// NewType returns a new type equivalent to the given type node.
func NewType(n Node) (types.Type, error) {
	switch n := n.(type) {
	case *ClassDecl:
		var err error
		typ := &types.Class{Name: n.ClassName.Name}
		typ.Members = make(map[string]types.Type)
		for _, member := range n.Members {
			switch member := member.(type) {
			case *SubDecl:
				typ.Members[member.SubName.Name], err = NewType(member.SubType)
				if err != nil {
					return nil, err
				}
			case *FuncDecl:
				typ.Members[member.FuncName.Name], err = NewType(member.FuncType)
				if err != nil {
					return nil, err
				}
			case *ScalarDecl:
				typ.Members[member.VarName.Name], err = NewType(member.VarType)
				if err != nil {
					return nil, err
				}
			case *ArrayDecl:
				typ.Members[member.VarName.Name], err = NewType(member.VarType)
				if err != nil {
					return nil, err
				}
			default:
				return nil, errors.Newf(n.Token().Position, "invalid declaration %T", member)
			}
		}
		return typ, nil
	case *BasicLit:
		switch n.Kind {
		case token.StringLit:
			return &types.Basic{Kind: types.String}, nil
		case token.LongLit:
			return &types.Basic{Kind: types.Long}, nil
		case token.DoubleLit:
			return &types.Basic{Kind: types.Double}, nil
		case token.BooleanLit, token.KwTrue, token.KwFalse:
			return &types.Basic{Kind: types.Boolean}, nil
		case token.CurrencyLit:
			return &types.Basic{Kind: types.Currency}, nil
		case token.DateLit:
			return &types.Basic{Kind: types.Date}, nil
		case token.KwNothing:
			return &types.Basic{Kind: types.Nothing}, nil
		default:
			return nil, errors.Newf(n.Token().Position, "invalid basic literal %v", n.Kind)
		}
	// case *ArrayType:
	// 	return &types.Array{Elem: newType(n.Elem), Len: n.Len}
	case *FuncType:
		params := make([]*types.Field, len(n.Params))
		for i := range n.Params {
			var err error
			params[i], err = NewField(&n.Params[i])
			if err != nil {
				return nil, err
			}
		}
		if n.Result == nil {
			return nil, nil
		}
		if typ, err := NewType(n.Result); err != nil {
			return nil, err
		} else {
			return &types.Func{Result: typ, Params: params}, nil
		}
	case *SubType:
		params := make([]*types.Field, len(n.Params))
		for i := range n.Params {
			var err error
			params[i], err = NewField(&n.Params[i])
			if err != nil {
				return nil, err
			}
		}
		return &types.Sub{Params: params}, nil
	case *Identifier:

		if n.Decl == nil {
			return NewBasic(n), nil
		}
		return n.Decl.Type()
	case *ParamItem:
		if n.IsArray {
			typ, err := NewType(n.VarType)
			if err != nil {
				return nil, err
			}
			array := &types.Array{Type: typ}
			return array, nil
		} else {
			return NewType(n.VarType)
		}
	case *ArrayType:
		typ, err := NewType(n.Type)
		if err != nil {
			return nil, err
		}
		array := &types.Array{Type: typ}
		for _, dim := range n.Dimensions {
			if dimLit, ok := dim.(*BasicLit); ok {
				if dimension, err := strconv.Atoi(dimLit.String()); err != nil {
					return nil, errors.Newf(n.Token().Position, "invalid array dimension %v", dimLit.String())
				} else if dimension < 1 {
					return nil, errors.Newf(n.Token().Position, "invalid array dimension %v", dimLit.String())
				} else {
					array.Dimensions = append(array.Dimensions, dimension)
				}
			} else if dimVar, ok := dim.(*Identifier); ok {
				// identifier must be of type Long
				if dimVar.Decl == nil {
					return nil, errors.Newf(n.Token().Position, "undeclared identifier %v", dimVar.Name)
				}
				typ, err := dimVar.Decl.Type()
				if err != nil {
					return nil, err
				}
				if basicType, ok := typ.(*types.Basic); ok && basicType.Kind == types.Long {
					array.Dimensions = append(array.Dimensions, -1)
				} else {
					return nil, errors.Newf(n.Token().Position, "invalid array dimension %v", dimVar.Name)
				}
			} else {
				typ, err := dim.Type()
				if err != nil {
					return nil, err
				}
				if basicType, ok := typ.(*types.Basic); ok && basicType.Kind == types.Long {
					array.Dimensions = append(array.Dimensions, -1)
				} else {
					return nil, errors.Newf(n.Token().Position, "invalid array dimension %v", dim.String())
				}
				array.Dimensions = append(array.Dimensions, -1)
			}
		}
		return array, nil
	case *BinaryExpr:
		// Check that the operand types are compatible with the operator.
		typX, err := NewType(n.Left)
		if err != nil {
			return nil, err
		}
		typY, err := NewType(n.Right)
		if err != nil {
			return nil, err
		}
		switch n.OpKind {
		case token.IntDiv:
			return &types.Basic{Kind: types.Long}, nil
		case token.Add, token.Minus, token.Mul, token.Div, token.Mod, token.Exponent: // arithmetic operators
			return HigherPrecision(typX, typY)
		case token.Eq, token.Neq, token.Lt, token.Le, token.Gt, token.Ge: // relational operators
			return &types.Basic{Kind: types.Boolean}, nil
		case token.And, token.Or: // logical operators
			return &types.Basic{Kind: types.Boolean}, nil
		case token.Concat: // concatenation operator
			return &types.Basic{Kind: types.String}, nil
		default:
			return nil, errors.Newf(n.Token().Position, "support for operator %q not yet implemented", n.OpKind)
		}
	case *UnaryExpr:
		return NewType(n.Right)
	case *CallOrIndexExpr:
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
			return nil, nil
		default:
			return nil, errors.Newf(n.Token().Position, "cannot call non-function or non-array %q of type %q", n.Identifier, typ)
		}
	case *UserDefinedType:
		return NewUserDefinedType(n.Identifier.Name), nil
	case *ParenExpr:
		return NewType(n.Expr)
	default:
		return nil, errors.Newf(n.Token().Position, "support for type %T not yet implemented", n)
	}
}

// NewField returns a new field type equivalent to the given field node.
func NewField(decl *ParamItem) (*types.Field, error) {

	typDecl, err := decl.Type()
	if err != nil {
		return nil, err
	}
	typ := &types.Field{Type: typDecl}
	if decl.VarName != nil {
		typ.Name = decl.VarName.Name
	}
	typ.ByVal = decl.ByVal
	typ.Optional = decl.Optional
	typ.ParamArray = decl.ParamArray
	if decl.DefaultValue != nil {
		typ.DefaultValue = decl.DefaultValue.String()
	}
	return typ, nil
}

// universePos specifies a pseudo-position used for identifiers declared in the
// universe scope.
var universePos = token.Position{Line: 0, Column: 0, Absolute: -1}
var universalToken = token.Token{Kind: token.Identifier, Literal: "", Position: universePos}

// NewBasic returns a new basic type equivalent to the given identifier.
func NewBasic(ident *Identifier) types.Type {
	// TODO: Check if we may come up with a cleaner solution. At least, this
	// works for now.
	switch ident.Name {
	case "String":
		stringIdent := &Identifier{Tok: &universalToken, Name: "String"}
		stringType := &types.Basic{Kind: types.String}
		stringDecl := &TypeDef{DeclType: stringIdent, TypeName: stringIdent, Val: stringType}
		stringIdent.Decl = stringDecl
		ident.Decl = stringDecl
		return stringType
	case "Long":
		longIdent := &Identifier{Tok: &universalToken, Name: "Long"}
		longType := &types.Basic{Kind: types.Long}
		longDecl := &TypeDef{DeclType: longIdent, TypeName: longIdent, Val: longType}
		longIdent.Decl = longDecl
		ident.Decl = longDecl
		return longType
	case "Double":
		doubleIdent := &Identifier{Tok: &universalToken, Name: "Double"}
		doubleType := &types.Basic{Kind: types.Double}
		doubleDecl := &TypeDef{DeclType: doubleIdent, TypeName: doubleIdent, Val: doubleType}
		doubleIdent.Decl = doubleDecl
		ident.Decl = doubleDecl
		return doubleType
	case "Boolean":
		booleanIdent := &Identifier{Tok: &universalToken, Name: "Boolean"}
		booleanType := &types.Basic{Kind: types.Boolean}
		booleanDecl := &TypeDef{DeclType: booleanIdent, TypeName: booleanIdent, Val: booleanType}
		booleanIdent.Decl = booleanDecl
		ident.Decl = booleanDecl
		return booleanType
	case "Variant":
		variantIdent := &Identifier{Tok: &universalToken, Name: "Variant"}
		variantType := &types.Basic{Kind: types.Variant}
		variantDecl := &TypeDef{DeclType: variantIdent, TypeName: variantIdent, Val: variantType}
		variantIdent.Decl = variantDecl
		ident.Decl = variantDecl
		return variantType
	case "Currency":
		currencyIdent := &Identifier{Tok: &universalToken, Name: "Currency"}
		currencyType := &types.Basic{Kind: types.Currency}
		currencyDecl := &TypeDef{DeclType: currencyIdent, TypeName: currencyIdent, Val: currencyType}
		currencyIdent.Decl = currencyDecl
		ident.Decl = currencyDecl
		return currencyType
	case "DateTime":
		dateTimeIdent := &Identifier{Tok: &universalToken, Name: "DateTime"}
		dateTimeType := &types.Basic{Kind: types.Date}
		dateTimeDecl := &TypeDef{DeclType: dateTimeIdent, TypeName: dateTimeIdent, Val: dateTimeType}
		dateTimeIdent.Decl = dateTimeDecl
		ident.Decl = dateTimeDecl
		return dateTimeType
	case "Nothing":
		nothingIdent := &Identifier{Tok: &universalToken, Name: "Nothing"}
		nothingType := &types.Basic{Kind: types.Nothing}
		nothingDecl := &TypeDef{DeclType: nothingIdent, TypeName: nothingIdent, Val: nothingType}
		nothingIdent.Decl = nothingDecl
		ident.Decl = nothingDecl
		return nothingType
	case "Single":
		singleIdent := &Identifier{Tok: &universalToken, Name: "Single"}
		singleType := &types.Basic{Kind: types.Single}
		singleDecl := &TypeDef{DeclType: singleIdent, TypeName: singleIdent, Val: singleType}
		singleIdent.Decl = singleDecl
		ident.Decl = singleDecl
		return singleType
	case "Integer":
		integerIdent := &Identifier{Tok: &universalToken, Name: "Integer"}
		integerType := &types.Basic{Kind: types.Integer}
		integerDecl := &TypeDef{DeclType: integerIdent, TypeName: integerIdent, Val: integerType}
		integerIdent.Decl = integerDecl
		ident.Decl = integerDecl
		return integerType
	default:
		return NewUserDefinedType(ident.Name)
	}
}

func NewUserDefinedType(name string) types.Type {
	ident := &Identifier{Tok: &universalToken, Name: name}
	userType := &types.UserDefined{Name: name}
	decl := &TypeDef{DeclType: ident, TypeName: ident, Val: userType}
	ident.Decl = decl
	return userType
}

const (
	VariantPrecision = iota
	BooleanPrecision
	IntegerPrecision
	LongPrecision
	CurrencyPrecision
	SinglePrecision
	DoublePrecision
	DatePrecision
	StringPrecision
	EnumPrecision
	ErrorPrecision // error flag
)

// higherPrecision returns the type of higher precision.
func HigherPrecision(dest, source types.Type) (types.Type, error) {
	if d, ok := dest.(*types.Basic); ok {
		if s, ok := source.(*types.Basic); ok {
			// same type
			if s.Kind == d.Kind {
				return &types.Basic{Kind: s.Kind}, nil
			}
			// order of precision
			// 0 = integer, long, currency, single, double, String, Variant
			precisionDest := GetPrecisionOrder(d.Kind)
			precisionSource := GetPrecisionOrder(s.Kind)
			precisionResult := precisionDest
			if precisionSource > precisionResult {
				precisionResult = precisionSource
			}

			switch precisionResult {
			case IntegerPrecision:
				return &types.Basic{Kind: types.Integer}, nil
			case LongPrecision:
				return &types.Basic{Kind: types.Long}, nil
			case CurrencyPrecision:
				return &types.Basic{Kind: types.Currency}, nil
			case SinglePrecision:
				return &types.Basic{Kind: types.Single}, nil
			case DoublePrecision:
				return &types.Basic{Kind: types.Double}, nil
			case StringPrecision:
				return &types.Basic{Kind: types.String}, nil
			case VariantPrecision:
				return &types.Basic{Kind: types.Variant}, nil
			}

		} else if d.Kind == types.Variant {
			// Variant can absorb any type.
			return &types.Basic{Kind: types.Variant}, nil
		}
	} else {
		if d, ok := dest.(*types.UserDefined); ok {
			if s, ok := source.(*types.UserDefined); ok {
				// Check for same types.
				if d.Name == s.Name {
					return d, nil
				}
			}
		}
	}
	return nil, errors.Newf(universePos, "Types %T and %T are not compatible.", dest, source)
}

func GetPrecisionOrder(typ types.BasicKind) int {
	switch typ {
	case types.Boolean:
		return BooleanPrecision
	case types.Integer:
		return IntegerPrecision
	case types.Long:
		return LongPrecision
	case types.Single:
		return SinglePrecision
	case types.Currency:
		return CurrencyPrecision
	case types.Double:
		return DoublePrecision
	case types.Date:
		return DatePrecision
	case types.Variant:
		return VariantPrecision
	case types.String:
		return StringPrecision
	case types.Enum:
		return EnumPrecision
	}
	return ErrorPrecision

}
