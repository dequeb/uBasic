// Package types declares the data types of the µBASIC programming language.
package types

import (
	"fmt"
	"strings"
)

// A Type represents a type of µBASIC, and has one of the following underlying
// types.
//
//	*Basic
//	*Array
//	*Func
type Type interface {
	// Equal reports whether t and u are of equal type.
	Equal(u Type) bool
	IsVariant() bool
	fmt.Stringer
}

// A Numerical type is numerical if so specified by IsNumerical.
type Numerical interface {
	// IsNumerical reports whether the given type is numerical.
	IsNumerical() bool
}

// Types.
type (
	// A Basic represents a basic type.
	//
	// Examples.
	//
	//    char
	//    int
	Basic struct {
		// Kind of basic type.
		Kind BasicKind
	}
	// An Array represents an array type.
	//
	// Examples.
	//
	//    Dim a(10) As Long
	//    Dim b(5, 2) As String
	Array struct {
		// Element type.
		Type Type
		// Array length.
		Dimensions []int
	}

	// A Func represents a function signature.
	//
	// Examples.
	//
	//    int(void)
	//    int(int a, int b)
	Func struct {
		// Return type.
		Result Type
		// Function parameter types; or nil if void parameter.
		Params []*Field
	}

	UserDefined struct {
		Name string
	}
)

//go:generate stringer -type BasicKind
//go:generate gorename -from basickind_string.go::i -to kind

// BasicKind describes the kind of basic type.
type BasicKind int

// Basic type.
const (
	Invalid BasicKind = iota // invalid type

	String
	Long
	Double
	Boolean
	DateTime
	Variant
	Single
	Integer
	Currency
	Nothing // no type
)

// A Field represents a field declaration in a struct type, or a parameter
// declaration in a function signature.
//
// Examples.
//
//	char
//	int a
type Field struct {
	// Field type.
	Type Type
	// Field name; or empty.
	Name string
	// What is the default value for this field?
	DefaultValue string
	// Is field by ref
	ByRef bool
	// is field optional
	Optional bool
	// is field a parameter array
	ParamArray bool
}

func (field *Field) String() string {
	if len(field.Name) > 0 {
		return fmt.Sprintf("%v %v", field.Type, field.Name)
	}
	return field.Type.String()
}

// Equal reports whether t and u are of equal type.
func (t *Basic) Equal(u Type) bool {
	if u == nil {
		return false
	}
	if u, ok := u.(*Basic); ok {
		return t.Kind == u.Kind
	}
	if u, ok := u.(*Func); ok {
		return u.Equal(t)
	}
	return false
}

func (t *UserDefined) Equal(u Type) bool {
	if u == nil {
		return false
	}
	if u, ok := u.(*UserDefined); ok {
		return t.Name == u.Name
	}
	if u, ok := u.(*Func); ok {
		return u.Equal(t)
	}
	return false
}

// Equal reports whether t and u are of equal type.
func (t *Array) Equal(u Type) bool {
	if u == nil {
		return false
	}
	if u, ok := u.(*Array); ok {
		if !Equal(t.Type, u.Type) {
			return false
		}
		if len(u.Dimensions) != len(t.Dimensions) {
			return false
		}
		for i := range t.Dimensions {
			if t.Dimensions[i] != u.Dimensions[i] {
				return false
			}
		}
		return true
	}
	return false
}

// Equal reports whether t and u are of equal type.
func (t *Func) Equal(u Type) bool {
	if u == nil {
		return false
	}
	// cannot assign functions in µBASIC
	// but we can assign function name as return value
	return t.Result.Equal(u)
}

// Equal reports whether t and u are of equal type.
func Equal(t, u Type) bool {
	if t == nil && u == nil {
		return true
	}
	if u == nil || t == nil {
		return false
	}
	return t.Equal(u)
}

// IsInteger reports whether the given type is an integer (i.e. "int" or
// "char").
func IsInteger(t Type) bool {
	if t == nil {
		return false
	}

	if t, ok := t.(*Basic); ok {
		switch t.Kind {
		case Long:
			return true
		}
	}
	return false
}

// IsNumerical reports whether the given type is numerical.
func (t *Basic) IsNumerical() bool {
	if t == nil {
		return false
	}
	switch t.Kind {
	case Long, Double, Variant, Single, Integer, Currency:
		return true
	case Boolean, String, Nothing, DateTime:
		return false
	default:
		panic(fmt.Sprintf("types.Basic.IsNumerical: unknown basic type (%d)", int(t.Kind)))
	}
}

func (t *Basic) String() string {
	names := map[BasicKind]string{
		String:   "String",
		Long:     "Long",
		Double:   "Double",
		Boolean:  "Boolean",
		DateTime: "DateTime",
		Variant:  "Variant",
		Single:   "Single",
		Integer:  "Integer",
		Currency: "Currency",
		Nothing:  "Nothing",
	}
	if s, ok := names[t.Kind]; ok {
		return s
	}
	return fmt.Sprintf("unknown basic type (%d)", int(t.Kind))
}

func (t *UserDefined) String() string {
	return t.Name
}

func (t *Array) String() string {
	buf := strings.Builder{}
	buf.WriteString("(")
	for i, dim := range t.Dimensions {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(fmt.Sprintf("%d", dim))
	}
	buf.WriteString(") As ")
	buf.WriteString(t.Type.String())
	return buf.String()
}

func (t *Func) String() string {
	buf := strings.Builder{}
	buf.WriteString("(")
	for _, param := range t.Params {
		buf.WriteString(param.String())
	}
	if t.Result != nil {
		buf.WriteString(") As ")
		buf.WriteString(t.Result.String())
	}
	return buf.String()
}

func (t *Basic) IsVariant() bool {
	return t.Kind == Variant
}
func (t *Array) IsVariant() bool {
	return t.Type.IsVariant()
}
func (t *Func) IsVariant() bool {
	return t.Result.IsVariant()
}
func (t *UserDefined) IsVariant() bool {
	return false
}

// Verify that the µBasic types implement the Type interface.
var (
	_ Type = &Basic{}
	_ Type = &Array{}
	_ Type = &Func{}
	_ Type = &UserDefined{}
)
