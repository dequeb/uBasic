package irgen

import (
	"fmt"
	"strconv"

	"uBasic/ast"
	"uBasic/eval"
	"uBasic/object"

	"github.com/llir/llvm/ir/constant"
	irtypes "github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

// implicitConversion implicitly converts the value of the smallest type to the
// largest type of x and y, emitting code to f. The new values of x and y are
// returned.
func (m *Module) implicitConversion(f *Function, x, y value.Value) (value.Value, value.Value) {
	// Implicit conversion.
	switch {
	case isLarger(x.Type(), y.Type()) || firstOnlyIsFloat(x.Type(), y.Type()):
		y = m.convert(f, y, x.Type())
	case isLarger(y.Type(), x.Type()) || firstOnlyIsFloat(y.Type(), x.Type()):
		x = m.convert(f, x, y.Type())
	}
	return x, y
}

// convert converts the given value to the specified type, emitting code to f.
// No conversion is made, if v is already of the correct type.
func (m *Module) convert(f *Function, v value.Value, to irtypes.Type) value.Value {
	// Early return if v is already of the correct type.
	from := v.Type()
	if irtypes.Equal(from, to) {
		return v
	}
	fromType, ok := from.(*irtypes.IntType)
	if ok {
		return m.convertIntType(f, v, fromType, to)
	} else if fromType, ok := from.(*irtypes.FloatType); ok {
		return m.convertFloatType(f, v, fromType, to)
	} else if _, ok := from.(*irtypes.PointerType); ok {
		return m.convertPointerType(f, v, to)
	} else {
		panic(fmt.Sprintf("support for converting to type %T not yet implemented", to))
	}
}

// convertIntType converts the given integer value to the specified type, emitting code to f.
func (m *Module) convertIntType(f *Function, v value.Value, fromType *irtypes.IntType, to irtypes.Type) value.Value {
	toIntType, ok := to.(*irtypes.IntType)
	if ok {
		toSize := toIntType.BitSize
		fromSize := fromType.BitSize
		if toSize > fromSize {
			return f.currentBlock.NewSExt(v, toIntType)
		} else if toSize < fromSize {
			return f.currentBlock.NewTrunc(v, toIntType)
		} else {
			return v
		}
	} else if toFloatType, ok := to.(*irtypes.FloatType); ok {
		return f.currentBlock.NewSIToFP(v, toFloatType)
	} else if pt, ok := to.(*irtypes.PointerType); ok {

		return m.convert(f, v, pt.ElemType)
	} else {
		panic(fmt.Sprintf("support for converting to type %T not yet implemented", to))

	}
}

// convertFloatType converts the given integer value to the specified type, emitting code to f.
func (m *Module) convertFloatType(f *Function, v value.Value, fromType *irtypes.FloatType, to irtypes.Type) value.Value {
	toFloadType, ok := to.(*irtypes.FloatType)
	if ok {
		toSize := 64
		switch toFloadType.Kind {
		case irtypes.FloatKindFloat:
			toSize = 32
		case irtypes.FloatKindDouble:
			toSize = 64
		}
		fromSize := 64
		switch fromType.Kind {
		case irtypes.FloatKindFloat:
			fromSize = 32
		case irtypes.FloatKindDouble:
			fromSize = 64
		}
		if toSize > fromSize {
			return f.currentBlock.NewFPExt(v, toFloadType)
		} else if toSize < fromSize {
			return f.currentBlock.NewFPTrunc(v, toFloadType)
		} else {
			return v
		}
	} else if toIntType, ok := to.(*irtypes.IntType); ok {
		return f.currentBlock.NewFPToSI(v, toIntType)
	} else if pt, ok := to.(*irtypes.PointerType); ok {
		return m.convert(f, v, pt.ElemType)
	} else {
		panic(fmt.Sprintf("support for converting to type %T not yet implemented", to))
	}

}

// convertPointerType converts the given integer value to the specified type, emitting code to f.
func (m *Module) convertPointerType(f *Function, v value.Value, to irtypes.Type) value.Value {

	var isPointer bool
	var isPointerTo bool
	_, isPointer = v.Type().(*irtypes.PointerType)
	_, isPointerTo = to.(*irtypes.PointerType)

	if isPointer && isPointerTo {
		return v
	}
	return m.convert(f, f.currentBlock.NewLoad(to, v), to)
}

// isLarger reports whether t has higher precision than u.
func isLarger(t, u irtypes.Type) bool {
	// A Sizer is a type with a size in number of bits.
	type Sizer interface {
		// Size returns the size of t in number of bits.
		Size() int
	}
	var tSize int
	var uSize int

	if ts, ok := t.(Sizer); ok {
		tSize = ts.Size()
	} else if ti, ok := t.(*irtypes.IntType); ok {
		tSize = int(ti.BitSize)
	} else if t, ok := t.(*irtypes.FloatType); ok {
		switch t.Kind {
		case irtypes.FloatKindFloat:
			tSize = 32
		case irtypes.FloatKindDouble:
			tSize = 64
		}
	}
	if us, ok := u.(Sizer); ok {
		uSize = us.Size()
	} else if ui, ok := u.(*irtypes.IntType); ok {
		uSize = int(ui.BitSize)
	} else if u, ok := u.(*irtypes.FloatType); ok {
		switch u.Kind {
		case irtypes.FloatKindFloat:
			uSize = 32
		case irtypes.FloatKindDouble:
			uSize = 64
		}
	}
	return tSize > uSize

}

// firstOnlyIsFloat reports whether the first type is a floating-point type and
// the second type is not a floating-point type.
func firstOnlyIsFloat(t, u irtypes.Type) bool {
	_, tIsFloat := t.(*irtypes.FloatType)
	_, uIsFloat := u.(*irtypes.FloatType)
	return tIsFloat && !uIsFloat
}

// isRef reports whether the given type is a reference type; e.g. pointer or
// array.
func isRef(typ irtypes.Type) bool {
	switch typ.(type) {
	case *irtypes.ArrayType:
		return true
	case *irtypes.PointerType:
		return true
	default:
		return false
	}
}

// constZero returns the integer constant zero of the given type.
func constZero(typ irtypes.Type) constant.Constant {
	intType, ok := typ.(*irtypes.IntType)
	if !ok {
		panic(fmt.Errorf("invalid integer literal type; expected *types.IntType, got %T", typ))
	}
	return constant.NewInt(intType, 0)
}

// constOne returns the integer constant one of the given type.
func constOne(typ irtypes.Type) constant.Constant {
	intType, ok := typ.(*irtypes.IntType)
	if !ok {
		panic(fmt.Errorf("invalid integer literal type; expected *types.IntType, got %T", typ))
	}
	return constant.NewInt(intType, 1)
}

// genUnique generates a unique local variable name based on the given identifier.
func (f *Function) genUnique(ident *ast.Identifier) string {
	name := ident.Name
	if !f.exists[name] {
		f.exists[name] = true
		return name
	}
	for i := 1; ; i++ {
		name := fmt.Sprintf("%s%d", ident.Name, i)
		if !f.exists[name] {
			f.exists[name] = true
			return name
		}
	}
}

// isGlobal reports whether the given identifier is a global definition.
func (m *Module) isGlobal(ident *ast.Identifier) bool {
	pos := ident.Decl.Name().Tok.Position.Absolute
	_, exists := m.idents[pos]
	return exists
}

// valueFromIdent returns the LLVM IR value associated with the given
// identifier. Only search for global values if f is nil.
func (m *Module) valueFromIdent(f *Function, ident *ast.Identifier) value.Value {
	pos := ident.Decl.Name().Tok.Position.Absolute
	if v, ok := m.idents[pos]; ok {
		return v
	}
	if f != nil {
		if v, ok := f.idents[pos]; ok {
			return v
		}
	}
	panic(fmt.Sprintf("unable to locate value associated with identifier %q (declared at source code position %d)", ident, pos))
}

// valueFromIdent returns the LLVM IR value associated with the given
// identifier. Only search for global values if f is nil.
func (m *Module) arrayDimensionValueFromIdent(f *Function, ident *ast.Identifier, dimension int) value.Value {
	pos := ident.Decl.Name().Tok.Position.Absolute + dimension + 1 // +1 to avoid conflict with the array itself
	if v, ok := m.idents[pos]; ok {
		return v
	}
	if f != nil {
		if v, ok := f.idents[pos]; ok {
			return v
		}
	}
	panic(fmt.Sprintf("unable to locate value associated with identifier %q (declared at source code position %d)", ident, pos))
}

// setIdentValue maps the given global identifier to the associated value.
func (m *Module) setIdentValue(ident *ast.Identifier, v value.Value) {
	pos := ident.Decl.Name().Tok.Position.Absolute
	if old, ok := m.idents[pos]; ok {
		panic(fmt.Sprintf("unable to map identifier %q to value %v; already mapped to value %v", ident, v, old))
	}
	m.idents[pos] = v
}

// setIdentValue maps the given local identifier to the associated value.
func (f *Function) setIdentValue(ident *ast.Identifier, v value.Value) {
	pos := ident.Decl.Name().Tok.Position.Absolute
	if old, ok := f.idents[pos]; ok {
		panic(fmt.Sprintf("unable to map identifier %q to value %v; already mapped to value %v", ident, v, old))
	}
	f.idents[pos] = v
}

// setIdentValue maps the given global identifier to the associated value.
func (m *Module) setArrayDimensionIdentValue(f *Function, ident *ast.Identifier, dimension int, v value.Value) {
	if v == nil {
		panic(fmt.Sprintf("unable to map identifier %q to nil value", ident))
	}
	pos := ident.Decl.Name().Tok.Position.Absolute + dimension + 1 // +1 to avoid conflict with the array itself
	if f != nil {
		if old, ok := f.idents[pos]; ok {
			panic(fmt.Sprintf("unable to map identifier %q to value %v; already mapped to value %v", ident, v, old))
		}
		f.idents[pos] = v
	} else {
		if old, ok := m.idents[pos]; ok {
			panic(fmt.Sprintf("unable to map identifier %q to value %v; already mapped to value %v", ident, v, old))
		}
		m.idents[pos] = v
	}
}

// typeOf returns the LLVM IR type of the given expression.
func (m *Module) typeOf(expr ast.Expression) irtypes.Type {
	if typ, ok := m.info.Types[expr]; ok {
		return toIrType(typ)
	}
	panic(fmt.Sprintf("unable to locate type for expression %v", expr))
}

// pointerToValue returns the LLVM IR pointer value of the given expression.
func (m *Module) pointerToValue(f *Function, x value.Value) value.Value {
	if ptrType, ok := x.Type().(*irtypes.PointerType); ok {
		if _, ok := ptrType.ElemType.(*irtypes.ArrayType); !ok {
			x = f.currentBlock.NewLoad(ptrType.ElemType, x)
			return m.pointerToValue(f, x)
		}
	}
	return x
}

func (m *Module) constantAstToValues(node *ast.ConstDeclItem) (valuestr string, valueInt int64, valueFloat float64, valueBool bool) {
	// all constants are evaluated at compile time
	Object := eval.Eval(nil, node.ConstValue, m.env)
	m.env.Set(node.ConstName.Name, Object)

	switch obj := Object.(type) {
	case *object.Integer:
		valueInt = int64(obj.Value)
		valuestr = strconv.FormatInt(valueInt, 10)
		valueFloat = float64(valueInt)
		valueBool = valueInt != 0
	case *object.Long:
		valueInt = obj.Value
		valuestr = strconv.FormatInt(valueInt, 10)
		valueFloat = float64(valueInt)
		valueBool = valueInt != 0
	case *object.Single:
		valueFloat = float64(obj.Value)
		valuestr = strconv.FormatFloat(valueFloat, 'f', -1, 32)
		valueInt = int64(valueFloat)
		valueBool = valueInt != 0
	case *object.Currency:
		valueFloat = obj.Value
		valuestr = strconv.FormatFloat(valueFloat, 'f', -1, 32)
		valueInt = int64(valueFloat)
		valueBool = valueInt != 0
	case *object.Double:
		valueFloat = obj.Value
		valuestr = strconv.FormatFloat(valueFloat, 'f', -1, 64)
		valueInt = int64(valueFloat)
		valueBool = valueInt != 0
	case *object.Boolean:
		valueBool = obj.Value
		valuestr = strconv.FormatBool(valueBool)
		if valueBool {
			valueInt = 1
			valueFloat = 1.0
		} else {
			valueInt = 0
			valueFloat = 0.0
		}
	case *object.String:
		valuestr = obj.Value
		valueInt = 0
		valueFloat = 0.0
		valueBool = false
	case *object.Date:
		valueFloat = FromDateToFloat(obj.Value)
		valuestr = obj.String()
		valueInt = int64(valueFloat)
		valueBool = valueInt != 0
	case *object.Error:
		panic(obj)
	default:
		panic(fmt.Sprintf("unknown type %T", obj))
	}
	return
}

// find the declaration scope of the array
func (m *Module) getDeclarationScope(ident *ast.Identifier) string {
	var decl ast.Node
	decl = ident.Decl
	for decl.GetParent() != nil {
		decl = decl.GetParent()
		functOrSub, ok := decl.(ast.FuncOrSub)
		if ok {
			return functOrSub.Name().Name
		}
	}
	return ""
}
