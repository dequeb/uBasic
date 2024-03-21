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
type SubOrFunc interface {
	isSubOrFunc()
	GetParams() []*Field
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
	Func struct {
		// Return type.
		Result Type
		// Function parameter types;
		Params []*Field
	}

	// A Sub represents a subroutine signature.
	//
	Sub struct {
		// Subroutine parameter types;
		Params []*Field
	}
	UserDefined struct {
		Name string
	}

	// a Class represents a class type
	Class struct {
		// Name of the class
		Name string
		// Members of the class
		Members map[string]Type
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
	Date
	Variant
	Single
	Integer
	Currency
	Enum
	Nothing // no type
)

// A Field represents a parameter
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

func (field *Field) Equal(u Field) bool {
	return field.Type.Equal(u.Type) && field.Name == u.Name &&
		field.DefaultValue == u.DefaultValue && field.ByRef == u.ByRef &&
		field.Optional == u.Optional && field.ParamArray == u.ParamArray
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
	// check if parameters are equal
	if u, ok := u.(*Func); ok {
		if len(t.Params) != len(u.Params) {
			return false
		}
		for i := range t.Params {
			if !t.Params[i].Equal(*u.Params[i]) {
				return false
			}
		}
		return t.Result.Equal(u)
	}
	return false
}

// Equal reports whether t and u are of equal type.
func (t *Sub) Equal(u Type) bool {
	if u == nil {
		return false
	}
	// check if parameters are equal
	if u, ok := u.(*Sub); ok {
		if len(t.Params) != len(u.Params) {
			return false
		}
		for i := range t.Params {
			if !t.Params[i].Equal(*u.Params[i]) {
				return false
			}
		}
		return true
	}
	return false
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

// Eqaul reports whether t and u are of equal type.
func (t *Class) Equal(u Type) bool {
	if u == nil {
		return false
	}
	if t == nil {
		return false
	}
	if u, ok := u.(*Class); ok {
		if t.Name != u.Name {
			return false
		}
		if len(t.Members) != len(u.Members) {
			return false
		}
		for i := range t.Members {
			if !t.Members[i].Equal(u.Members[i]) {
				return false
			}
		}
		return true
	}
	return false
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
	case Boolean, String, Nothing, Date:
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
		Date:     "DateTime",
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
	for i, param := range t.Params {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(param.String())
	}
	if t.Result != nil {
		buf.WriteString(") As ")
		buf.WriteString(t.Result.String())
	}
	return buf.String()
}

func (t *Sub) String() string {
	buf := strings.Builder{}
	buf.WriteString("(")
	for i, param := range t.Params {
		if i != 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(param.String())
	}
	return buf.String()
}

func (t *Class) String() string {
	buf := strings.Builder{}
	buf.WriteString("Class ")
	buf.WriteString(t.Name)
	buf.WriteString(" {")
	for _, field := range t.Members {
		buf.WriteString(field.String())
	}
	buf.WriteString("}")
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
func (t *Sub) IsVariant() bool {
	return false
}

func (t *UserDefined) IsVariant() bool {
	return false
}

func (t *Class) IsVariant() bool {
	return false
}

// Verify that the µBasic types implement the Type interface.
var (
	_ Type = &Basic{}
	_ Type = &Array{}
	_ Type = &Func{}
	_ Type = &Sub{}
	_ Type = &Class{}
	_ Type = &UserDefined{}
)

// GetBasicType returns the basic type for the given type.
func GetBasicType(t Type) BasicKind {
	if t, ok := t.(*Basic); ok {
		return t.Kind
	}
	if t, ok := t.(*Array); ok {
		return GetBasicType(t.Type)
	}
	if t, ok := t.(*Func); ok {
		return GetBasicType(t.Result)
	}
	return Invalid
}

func (s *Sub) isSubOrFunc()  {}
func (f *Func) isSubOrFunc() {}

func (s *Sub) GetParams() []*Field {
	return s.Params
}

func (f *Func) GetParams() []*Field {
	return f.Params
}

// type validation
var (
	_ SubOrFunc = &Sub{}
	_ SubOrFunc = &Func{}
)
