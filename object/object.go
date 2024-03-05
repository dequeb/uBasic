package object

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
	"uBasic/ast"
	"uBasic/source"
	"uBasic/token"
	"uBasic/types"

	"github.com/mewkiz/pkg/term"
)

type ObjectType string

var UserColor = false

// any is a type that can hold any value.
type Object interface {
	fmt.Stringer
	Type() ObjectType
	Copy() Object
	Position() token.Position
	IsConstant() bool
	Equals(Object) bool
	GetValue() any
}

const (
	BOOLEAN_OBJ  ObjectType = "Boolean"
	LONG_OBJ     ObjectType = "Long"
	INTEGER_OBJ  ObjectType = "Integer"
	SINGLE_OBJ   ObjectType = "Single"
	DOUBLE_OBJ   ObjectType = "Double"
	CURRENCY_OBJ ObjectType = "Currency"
	STRING_OBJ   ObjectType = "String"
	ARRAY_OBJ    ObjectType = "Array"
	VARIANT_OBJ  ObjectType = "Variant"
	DATE_OBJ     ObjectType = "Date"
	NOTHING_OBJ  ObjectType = "Nothing"
	USERDEF_OBJ  ObjectType = "UserDefined"
	ERROR_OBJ    ObjectType = "Error"
	RETURN_OBJ   ObjectType = "Return value"
	EXIT_OBJ     ObjectType = "Exit"
	FUNCTION_OBJ ObjectType = "Function"
	SUB_OBJ      ObjectType = "Sub"
	CLASS_OBJ    ObjectType = "Class"
	RESUME_OBJ   ObjectType = "Resume"
)

var (
	NOTHING = &Nothing{}
)

// ----------------------- Environment -----------------------

type Environment struct {
	fmt.Stringer
	store  map[string]Object
	From   ast.Node
	Parent *Environment
}

func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]Object)}
}

func (e *Environment) Get(name string) (Object, bool) {
	name = strings.ToLower(name)
	obj, ok := e.store[name]
	if !ok && e.Parent != nil {
		obj, ok = e.Parent.Get(name)
	}
	return obj, ok
}

func (e *Environment) CallStack() string {
	str := ""
	if e.From != nil && e.From.Token() != nil {
		pos := e.From.Token().Position
		str += strconv.Itoa(pos.Line) + "\n"
	}
	if e.Parent != nil {
		str += e.Parent.CallStack()
	}
	return str
}

func (e *Environment) Set(name string, val Object) {
	name = strings.ToLower(name)
	e.store[name] = val
}

func (e *Environment) Delete(name string) {
	name = strings.ToLower(name)
	delete(e.store, name)
}

func (e *Environment) Copy() *Environment {
	env := NewEnvironment()
	for k, v := range e.store {
		env.Set(k, v)
	}
	return env
}

func (e *Environment) Extend() *Environment {

	env := NewEnvironment()
	env.Parent = e
	return env
}

func (e *Environment) Merge(env *Environment) {
	for k, v := range env.store {
		e.Set(k, v)
	}
}

func (e *Environment) String() string {
	buf := strings.Builder{}

	if e.Parent != nil && e.Parent.Parent != nil {
		buf.WriteString(e.Parent.String())
		buf.WriteString("━━━━━\n")
	}

	// must make sur the order remains constant to avoid flicker in the display
	keys := make([]string, len(e.store))
	i := uint64(0)
	for k := range e.store {
		keys[i] = k
		i++
	}
	// sort keys
	sort.Strings(keys)

	var content string
	for _, k := range keys {
		v := e.store[k]
		if v != nil {
			// display all except functions
			if v.Type() != FUNCTION_OBJ && v.Type() != SUB_OBJ {

				if UserColor {
					content = term.RedBold(fmt.Sprint(v.String()))
				} else {
					content = fmt.Sprint(v.String())
				}
				buf.WriteString(fmt.Sprintf("%s:%s\n", k, content))
			}
		}
	}
	return buf.String()
}

// ----------------------- Nothing -----------------------

type Nothing struct{}

func (s *Nothing) String() string {
	return ""
}

func (s *Nothing) Type() ObjectType {
	return NOTHING_OBJ
}

func (s *Nothing) Copy() Object {
	return NOTHING
}

func (s *Nothing) IsConstant() bool {
	return true
}

func (s *Nothing) Position() token.Position {
	return token.Position{Line: 0, Column: 0, Absolute: -1}
}

func (s *Nothing) Equals(other Object) bool {
	_, ok := other.(*Nothing)
	return ok
}

func (s *Nothing) GetValue() any {
	return nil
}

// ----------------------- Boolean -----------------------

type Boolean struct {
	Value bool
	Const bool
	Pos   token.Position
}

func (b *Boolean) String() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

func NewBooleanByBool(value bool, pos token.Position) *Boolean {
	return &Boolean{Value: value, Pos: pos}
}

func NewBoolean(value string, pos token.Position) (*Boolean, error) {
	val, err := strconv.ParseBool(value)
	if err != nil {
		return nil, err
	}
	return NewBooleanByBool(val, pos), nil
}

func (b *Boolean) Copy() Object {
	return NewBooleanByBool(b.Value, b.Pos)
}

func (b *Boolean) IsConstant() bool {
	return b.Const
}

func (s *Boolean) Position() token.Position {
	return token.Position{Line: 0, Column: 0, Absolute: -1}
}

func (b *Boolean) Equals(other Object) bool {
	otherBool, ok := other.(*Boolean)
	if !ok {
		return false
	}
	return b.Value == otherBool.Value
}

func (b *Boolean) GetValue() any {
	return b.Value
}

// ----------------------- Integer -----------------------

type Integer struct {
	Value int32
	Const bool
	Pos   token.Position
}

func (l *Integer) String() string {
	return fmt.Sprintf("%d", l.Value)
}

func (l *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

func (l *Integer) IsConstant() bool {
	return l.Const
}

func NewInteger(value string, pos token.Position) (*Integer, error) {
	val, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return nil, err
	}

	return &Integer{Value: int32(val), Pos: pos}, nil
}

func NewIntegerByInt(value int32, pos token.Position) *Integer {
	return &Integer{Value: value, Pos: pos}
}

func (l *Integer) Copy() Object {
	return &Integer{Value: l.Value, Pos: l.Pos}
}

func (s *Integer) Position() token.Position {
	return token.Position{Line: 0, Column: 0, Absolute: -1}
}

func (l *Integer) Equals(other Object) bool {
	otherInt, ok := other.(*Integer)
	if !ok {
		return false
	}
	return l.Value == otherInt.Value
}

func (l *Integer) GetValue() any {
	return l.Value
}

// ----------------------- Long -----------------------

type Long struct {
	Value int64
	Const bool
	Pos   token.Position
}

func (l *Long) String() string {
	return fmt.Sprintf("%d", l.Value)
}

func (l *Long) Type() ObjectType {
	return LONG_OBJ
}

func (l *Long) IsConstant() bool {
	return l.Const
}

func NewLong(value string, pos token.Position) (*Long, error) {
	val, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return nil, err
	}

	return &Long{Value: val, Pos: pos}, nil
}

func NewLongByInt(value int64, pos token.Position) *Long {
	return &Long{Value: value, Pos: pos}
}
func (l *Long) Copy() Object {
	return &Long{Value: l.Value, Pos: l.Pos}
}

func (l *Long) Position() token.Position {
	return l.Pos
}

func (l *Long) Equals(other Object) bool {
	otherLong, ok := other.(*Long)
	if !ok {
		return false
	}
	return l.Value == otherLong.Value
}

func (l *Long) GetValue() any {
	return l.Value
}

// ----------------------- Single -----------------------

type Single struct {
	Value float32
	Const bool
	Pos   token.Position
}

func (d *Single) String() string {
	return fmt.Sprintf("%f", d.Value)
}

func (d *Single) Type() ObjectType {
	return SINGLE_OBJ
}

func (d *Single) IsConstant() bool {
	return d.Const
}

func NewSingle(value string, pos token.Position) (*Single, error) {
	val, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return nil, err
	}
	return &Single{Value: (float32)(val), Pos: pos}, nil
}

func NewSingleByFloat(value float32, pos token.Position) *Single {
	return &Single{Value: value, Pos: pos}
}

func (d *Single) Copy() Object {
	return &Single{Value: d.Value, Pos: d.Pos}
}

func (s *Single) Position() token.Position {
	return s.Pos
}

func (d *Single) Equals(other Object) bool {
	otherSingle, ok := other.(*Single)
	if !ok {
		return false
	}
	return d.Value == otherSingle.Value
}

func (l *Single) GetValue() any {
	return l.Value
}

// ----------------------- Double -----------------------

type Double struct {
	Value float64
	Const bool
	Pos   token.Position
}

func (d *Double) String() string {
	return fmt.Sprintf("%f", d.Value)
}

func (d *Double) Type() ObjectType {
	return DOUBLE_OBJ
}

func (d *Double) IsConstant() bool {
	return d.Const
}

func NewDouble(value string, pos token.Position) (*Double, error) {
	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil, err
	}
	return &Double{Value: val, Pos: pos}, nil
}

func NewDoubleByFloat(value float64, pos token.Position) *Double {
	return &Double{Value: value, Pos: pos}
}

func (d *Double) Copy() Object {
	return &Double{Value: d.Value, Pos: d.Pos}
}

func (s *Double) Position() token.Position {
	return s.Pos
}

func (d *Double) Equals(other Object) bool {
	otherDouble, ok := other.(*Double)
	if !ok {
		return false
	}
	return d.Value == otherDouble.Value
}

func (l *Double) GetValue() any {
	return l.Value
}

// ----------------------- Currency -----------------------

type Currency struct {
	Value float64
	Const bool
	Pos   token.Position
}

func (d *Currency) String() string {
	return fmt.Sprintf("%f", d.Value)
}

func (d *Currency) Type() ObjectType {
	return CURRENCY_OBJ
}

func (d *Currency) IsConstant() bool {
	return d.Const
}

func NewCurrencyByFloat(value float64, pos token.Position) *Currency {
	return &Currency{Value: value, Pos: pos}
}

func NewCurrency(value string, pos token.Position) (*Currency, error) {
	// remove last currency symbol if present
	value = strings.Trim(value, "$")
	val, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return nil, err
	}

	return &Currency{Value: val, Pos: pos}, nil
}

func (d *Currency) Copy() Object {
	return &Currency{Value: d.Value, Pos: d.Pos}
}

func (s *Currency) Position() token.Position {
	return s.Pos
}

func (d *Currency) Equals(other Object) bool {
	otherCurrency, ok := other.(*Currency)
	if !ok {
		return false
	}
	return d.Value == otherCurrency.Value
}

func (l *Currency) GetValue() any {
	return l.Value
}

// ----------------------- String -----------------------

type String struct {
	Value string
	Const bool
	Pos   token.Position
}

func (s *String) String() string {
	return s.Value
}

func (s *String) Type() ObjectType {
	return STRING_OBJ
}

func (s *String) IsConstant() bool {
	return s.Const
}

func NewString(value string, pos token.Position) *String {
	value = strings.Trim(value, "\"")
	value = strings.ReplaceAll(value, `""`, `"`)
	return &String{Value: value, Pos: pos}
}

func (s *String) Copy() Object {
	return &String{Value: s.Value, Pos: s.Pos}
}

func (s *String) Position() token.Position {
	return s.Pos
}

func (s *String) Equals(other Object) bool {
	otherString, ok := other.(*String)
	if !ok {
		return false
	}
	return s.Value == otherString.Value
}

func (l *String) GetValue() any {
	return l.Value
}

// ----------------------- Array -----------------------

type Array struct {
	Values     map[uint32]Object
	Const      bool
	SubType    ObjectType
	Dimensions []uint32
	Pos        token.Position
}

func (a *Array) Type() ObjectType {
	return ARRAY_OBJ
}

func (a *Array) String() string {
	buf := strings.Builder{}
	if a.Values != nil {
		buf.WriteString("(")
		for i := 0; i < int(a.Dimensions[0]); i++ {
			v := a.Values[uint32(i)]
			if v == nil {
				v = NewEmptyByType(a.SubType, a.Pos)
			}
			if i != 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(v.String())
		}
		buf.WriteString(")")
	} else {
		buf.WriteString("()")
	}
	return buf.String()
}

func (a *Array) IsConstant() bool {
	return a.Const
}

func NewArray(subType ObjectType, dimensions []uint32, pos token.Position) *Array {
	return &Array{SubType: subType, Dimensions: dimensions, Pos: pos}
}

func (a *Array) Copy() Object {
	arr := &Array{Pos: a.Pos, SubType: a.SubType, Dimensions: a.Dimensions}
	if a.Values != nil {
		arr.Values = make(map[uint32]Object)
		for k, v := range a.Values {
			arr.Values[k] = v.Copy()
		}
	}
	return arr
}

// Set sets the value at the given index
func (a *Array) Set(index []uint32, Value Object) Object {
	// set default value
	if Value == nil {
		Value = NewEmptyByType(a.SubType, a.Pos)
	}
	// initialize empty array
	if a.Values == nil {
		a.Values = make(map[uint32]Object)
	}

	i := index[0] // boundaries checked by convertIndex()
	// last array dimension
	if len(index) == 1 {
		_, ok := a.Values[i]
		if !ok {
			a.Values[i] = Value.Copy() // always use a copy to avoid changing the original value
		}
		return a.Values[i]
	}
	// sub array
	var array *Array
	if a.Values[i] == nil {
		array = NewArray(a.SubType, index[1:], a.Pos)
		a.Values[i] = array
	} else {
		array = a.Values[i].(*Array)
	}
	return array.Set(index[1:], Value)
}

func (a *Array) Redimension(preserved bool, dimension uint32) Object {
	if !preserved {
		a.Values = nil
		a.Dimensions = []uint32{dimension}
		return a
	}
	// shrink or expand the array?
	if len(a.Dimensions) == 1 {
		// shrink
		if dimension < a.Dimensions[0] {
			for i := dimension; i < a.Dimensions[0]; i++ {
				delete(a.Values, i)
			}
			a.Dimensions[0] = dimension
		}
		// expand
		if dimension > a.Dimensions[0] {
			a.Dimensions[0] = dimension
		}
		return a
	} else {
		// cannot resize a multi-dimensional array
		return NewError(a.Pos, "Cannot resize a multi-dimensional array")
	}
}

// Get gets the value at the given index
// the value must be bound to the array
// to make sure that we can receive the updated value
func (a *Array) Get(index []uint32, defaultValue Object) Object {
	// set default value
	if defaultValue == nil {
		defaultValue = NewEmptyByType(a.SubType, a.Pos)
	}
	// initialize empty array
	if a.Values == nil {
		a.Values = make(map[uint32]Object)
	}

	i := index[0] // boundaries checked by convertIndex()
	// last array dimension
	if len(index) == 1 {
		_, ok := a.Values[i]
		if !ok {
			a.Values[i] = defaultValue
		}
		return a.Values[i]
	}
	// sub array
	var array *Array
	if a.Values[i] == nil {
		array = NewArray(a.SubType, index[1:], a.Pos)
		a.Values[i] = array
	} else {
		array = a.Values[i].(*Array)
	}
	return array.Get(index[1:], defaultValue)
}

func (a *Array) converIndex(index []Object, value Object) ([]uint32, Object) {
	// convert index to uint32
	indexes := make([]uint32, len(index))
	for i, v := range index {
		if v.Type() != INTEGER_OBJ {
			if v.Type() != LONG_OBJ {
				return indexes, NewError(v.Position(), "index must be an integer or long")
			}
			if v.(*Long).Value < 0 {
				return indexes, NewError(v.Position(), "index must be positive")
			}
			if v.(*Long).Value > int64(a.Dimensions[i]) {
				return indexes, NewError(a.Pos, "index out of range")
			}

			indexes[i] = uint32(v.(*Long).Value)
			continue
		}
		if v.(*Integer).Value < 0 {
			return indexes, NewError(v.Position(), "index must be positive")
		}
		if uint32(v.(*Integer).Value) >= a.Dimensions[i] {
			return indexes, NewError(a.Pos, "index out of range")
		}

		indexes[i] = uint32(v.(*Integer).Value)
	}
	return indexes, value

}

func (a *Array) SetValueAt(index []Object, value Object) Object {
	indexes, value := a.converIndex(index, value)
	if value != nil && value.Type() == ERROR_OBJ {
		return value
	}
	return a.Set(indexes, value)
}

func (a *Array) GetValueAt(index []Object) Object {
	indexes, value := a.converIndex(index, nil)
	if value != nil && value.Type() == ERROR_OBJ {
		return value
	}
	return a.Get(indexes, nil)
}

func (a *Array) Position() token.Position {
	return a.Pos
}

func (a *Array) Equals(other Object) bool {
	otherArray, ok := other.(*Array)
	if !ok {
		return false
	}
	if len(a.Dimensions) != len(otherArray.Dimensions) {
		return false
	}
	for i, d := range a.Dimensions {
		if d != otherArray.Dimensions[i] {
			return false
		}
	}

	if a.Values == nil && otherArray.Values == nil {
		for i, v := range a.Values {
			if !v.Equals(otherArray.Values[i]) {
				return false
			}
		}
	}
	return true
}

func (l *Array) GetValue() any {
	return l.Values
}

// ----------------------- Variant -----------------------

type Variant struct {
	Value Object
	Const bool
	Pos   token.Position
}

func (v *Variant) Type() ObjectType {
	return VARIANT_OBJ
}

func (v *Variant) String() string {
	return v.Value.String()
}

func (v *Variant) IsConstant() bool {
	return v.Const
}
func NewVariantByObject(value Object, pos token.Position) *Variant {
	return &Variant{Value: value, Pos: pos}
}

func NewVariant(value string, pos token.Position) (*Variant, error) {
	if value == "" {
		return &Variant{Value: NOTHING, Pos: pos}, nil
	}
	if strings.EqualFold(value, "true") {
		return &Variant{Value: NewBooleanByBool(true, pos), Pos: pos}, nil
	}
	if strings.EqualFold(value, "false") {
		return &Variant{Value: NewBooleanByBool(false, pos), Pos: pos}, nil
	}
	if value[0] == '"' {
		return &Variant{Value: NewString(value, pos), Pos: pos}, nil
	}
	if value[0] == '#' {
		d, err := NewDate(value, pos)
		if err != nil {
			return nil, err
		}
		return &Variant{Value: d, Pos: pos}, nil
	}
	if value[len(value)-1] == '$' {
		c, err := NewCurrency(value, pos)
		if err != nil {
			return nil, err
		}
		return &Variant{Value: c, Pos: pos}, nil
	}
	if strings.Contains(value, ".") {
		d, err := NewDouble(value, pos)
		if err != nil {
			return nil, err
		}
		return &Variant{Value: d, Pos: pos}, nil
	}
	i, err := NewInteger(value, pos)
	if err != nil {
		return nil, err
	}
	return &Variant{Value: i, Pos: pos}, nil
}

func (v *Variant) Copy() Object {
	return &Variant{Value: v.Value.Copy(), Pos: v.Pos}
}

func (v *Variant) Position() token.Position {
	return v.Pos
}

func (v *Variant) Equals(other Object) bool {
	otherVariant, ok := other.(*Variant)
	if !ok {
		return false
	}
	return v.Value.Equals(otherVariant.Value)
}

func (l *Variant) GetValue() any {
	if l.Value == nil {
		return nil
	}
	return l.Value.GetValue()
}

// ----------------------- Class -----------------------
type Class struct {
	Name    string
	Pos     token.Position
	Members map[string]Object
}

func (c *Class) Type() ObjectType {
	return CLASS_OBJ
}

func (c *Class) String() string {
	return c.Name
}

func NewClass(name string, pos token.Position) *Class {
	return &Class{Name: name, Pos: pos, Members: make(map[string]Object)}
}

func (c *Class) IsConstant() bool {
	return true
}

func (c *Class) Copy() Object {
	return c // object is immutable
}

func (c *Class) Position() token.Position {
	return c.Pos
}

func (c *Class) Equals(other Object) bool {
	otherClass, ok := other.(*Class)
	if !ok {
		return false
	}
	if c.Name != otherClass.Name {
		return false
	}
	if len(c.Members) != len(otherClass.Members) {
		return false
	}
	for k, v := range c.Members {
		otherV, ok := otherClass.Members[k]
		if !ok {
			return false
		}
		if !v.Equals(otherV) {
			return false
		}
	}
	return true
}

func (l *Class) GetValue() any {
	return l.Members
}

// ----------------------- Date -----------------------
type Date struct {
	Value time.Time
	Const bool
	Decl  ast.Decl
	Pos   token.Position
}

const YYYYMMDD_HHMMSS = "2006-01-02 15:04:05" // default date-time format
const YYYYMMDD = "2006-01-02"                 // default date
const HHMMSS = "15:04:05"                     // default time

func NewDate(value string, pos token.Position) (*Date, error) {
	var dateFormats = []string{
		YYYYMMDD,
		"2006/01/02",
		YYYYMMDD_HHMMSS,
		"2006/01/02 15:04:05",
		HHMMSS,
	}
	var err error
	var dt time.Time
	s := strings.Trim(value, "#")
	for _, format := range dateFormats {
		dt, err = time.Parse(format, s)
		if err == nil {
			return &Date{Value: dt, Pos: pos}, nil
		}
	}
	return nil, err
}

func NewDateByTime(value time.Time, pos token.Position) *Date {
	return &Date{Value: value, Pos: pos}
}

func (d *Date) Type() ObjectType {
	return DATE_OBJ
}

func (d *Date) String() string {
	if d.Value.IsZero() {
		return ""
	} else if d.Value.Hour() == 0 && d.Value.Minute() == 0 && d.Value.Second() == 0 {
		return d.Value.Format(YYYYMMDD)
	} else if d.Value.Year() == 0 && d.Value.Month() == 1 && d.Value.Day() == 1 {
		return d.Value.Format(HHMMSS)
	}
	return d.Value.Format(YYYYMMDD_HHMMSS)
}

func (d *Date) IsConstant() bool {
	return d.Const
}

func (d *Date) Copy() Object {
	return &Date{Value: d.Value}
}

func (d *Date) Position() token.Position {
	return d.Pos
}

func (d *Date) Equals(other Object) bool {
	otherDate, ok := other.(*Date)
	if !ok {
		return false
	}
	return d.Value.Equal(otherDate.Value)
}

func (l *Date) GetValue() any {
	return l.Value
}

// ----------------------- UserDefined -----------------------
type UserDefined struct {
	Value string
	Const bool
	Decl  *ast.EnumDecl
	Pos   token.Position
}

func (u *UserDefined) Type() ObjectType {
	return USERDEF_OBJ
}

func (u *UserDefined) String() string {
	return u.Value
}

func (u *UserDefined) IsConstant() bool {
	return u.Const
}

func (u *UserDefined) Copy() Object {
	return &UserDefined{Value: u.Value, Decl: u.Decl, Pos: u.Pos}
}

func (u *UserDefined) Position() token.Position {
	return u.Pos
}

func NewUserDefined(value string, decl *ast.EnumDecl, pos token.Position) *UserDefined {
	return &UserDefined{Value: value, Decl: decl, Pos: pos}
}

func (u *UserDefined) Equals(other Object) bool {
	otherUserDefined, ok := other.(*UserDefined)
	if !ok {
		return false
	}
	return u.Value == otherUserDefined.Value
}

func (l *UserDefined) GetValue() any {
	return l.Value
}

// ----------------------- Error -----------------------

type StackPosition struct {
	Positions []token.Position
	Source    *source.Source
}

func (s *StackPosition) String() string {
	buf := strings.Builder{}

	for i, p := range s.Positions {
		line := ""
		if s.Source != nil {
			line = ": " + strings.Trim(s.Source.Line(s.Positions[i]), "\n\r\t ")
		}
		buf.WriteString(p.String() + line + "\n")
	}
	return buf.String()
}

func (s *StackPosition) Push(pos token.Position) {
	s.Positions = append(s.Positions, pos)
}

func NewStack() *StackPosition {
	positions := make([]token.Position, 0)
	return &StackPosition{Positions: positions}
}

type Error struct {
	Value string
	Pos   token.Position
	Stack *StackPosition
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}

func (e *Error) String() string {
	line := ""
	if e.Stack.Source != nil {
		line = ": " + strings.Trim(e.Stack.Source.Line(e.Pos), "\n\r\t ")
	}
	if UserColor {
		return e.Pos.String() + ": " + line + term.Color(" ━━≻ ", term.Bold) + term.RedBold(e.Value) + "\n" + e.Stack.String()
	} else {
		return e.Pos.String() + ": " + line + " ━━≻ " + e.Value + "\n" + e.Stack.String()
	}
}

func (e *Error) Push(pos token.Position) {
	e.Stack.Push(pos)
}

func (e *Error) IsConstant() bool {
	return true
}

func (e *Error) Error() string {
	return e.String()
}

func NewError(pos token.Position, text string) *Error {
	err := &Error{
		Pos:   pos,
		Value: text,
	}
	err.Stack = NewStack()
	return err
}

func (e *Error) Copy() Object {
	return nil
}

func (e *Error) Position() token.Position {
	return e.Pos
}

func (e *Error) Equals(other Object) bool {
	otherError, ok := other.(*Error)
	if !ok {
		return false
	}
	return e.Value == otherError.Value
}

func (l *Error) GetValue() any {
	return l.Value
}

// ----------------------- return value -----------------------

type ReturnValue struct {
	Value Object
	Pos   token.Position
}

func NewReturnValue(value Object, pos token.Position) *ReturnValue {
	return &ReturnValue{Value: value, Pos: pos}
}

func (rv *ReturnValue) Type() ObjectType {
	return RETURN_OBJ
}

func (rv *ReturnValue) IsConstant() bool {
	return true
}

func (rv *ReturnValue) String() string {
	return rv.Value.String()
}

func (rv *ReturnValue) Copy() Object {
	return &ReturnValue{Value: rv.Value.Copy()}
}

func (rv *ReturnValue) Position() token.Position {
	return rv.Pos
}

func (rv *ReturnValue) Equals(other Object) bool {
	otherReturnValue, ok := other.(*ReturnValue)
	if !ok {
		return false
	}
	return rv.Value.Equals(otherReturnValue.Value)
}

func (l *ReturnValue) GetValue() any {
	if l.Value == nil {
		return nil
	}
	return l.Value.GetValue()
}

// ----------------------- exit -----------------------

type Exit struct {
	Kind token.Kind
	Pos  token.Position
}

func NewExit(kind token.Kind, pos token.Position) *Exit {
	return &Exit{Kind: kind, Pos: pos}
}

func (ex *Exit) Type() ObjectType {
	return EXIT_OBJ
}

func (ex *Exit) IsConstant() bool {
	return true
}

func (ex *Exit) String() string {
	return ex.Kind.String()
}

func (ex *Exit) Copy() Object {
	return &Exit{Kind: ex.Kind, Pos: ex.Pos}
}

func (ex *Exit) Position() token.Position {
	return ex.Pos
}

func (ex *Exit) Equals(other Object) bool {
	otherExit, ok := other.(*Exit)
	if !ok {
		return false
	}
	return ex.Kind == otherExit.Kind
}

func (l *Exit) GetValue() any {
	return l.Kind
}

// ----------------------- function -----------------------

type Function struct {
	Definition *ast.FuncDecl
	Parameters *ast.FuncType
	Body       []ast.StatementList
	Env        *Environment
	Pos        token.Position
}

func NewFunction(def *ast.FuncDecl, params *ast.FuncType, body []ast.StatementList, env *Environment, pos token.Position) *Function {
	return &Function{
		Definition: def,
		Parameters: params,
		Body:       body,
		Env:        env,
		Pos:        pos,
	}
}

func (f *Function) Type() ObjectType {
	return FUNCTION_OBJ
}

func (f *Function) String() string {
	if f.Definition != nil {
		return f.Definition.String()
	}
	return "Function"
}

func (f *Function) IsConstant() bool {
	return true
}

func (f *Function) Copy() Object {
	return &Function{
		Definition: f.Definition,
		Parameters: f.Parameters,
		Body:       f.Body,
		Env:        f.Env.Copy(),
		Pos:        f.Pos,
	}
}

func (f *Function) Position() token.Position {
	return f.Pos
}

func (f *Function) Equals(other Object) bool {
	otherFunction, ok := other.(*Function)
	if !ok {
		return false
	}
	return f.Definition.String() == otherFunction.Definition.String()
}

func (l *Function) GetValue() any {
	return nil
}

// ----------------------- sub -----------------------

type Sub struct {
	Definition *ast.SubDecl
	Parameters *ast.SubType
	Body       []ast.StatementList
	Env        *Environment
	Pos        token.Position
}

func NewSub(def *ast.SubDecl, params *ast.SubType, body []ast.StatementList, env *Environment, pos token.Position) *Sub {
	return &Sub{
		Definition: def,
		Parameters: params,
		Body:       body,
		Env:        env,
		Pos:        pos,
	}
}

func (s *Sub) Type() ObjectType {
	return SUB_OBJ
}

func (s *Sub) String() string {
	return s.Definition.String()
}

func (s *Sub) IsConstant() bool {
	return true
}

func (s *Sub) Copy() Object {
	return &Sub{
		Definition: s.Definition,
		Parameters: s.Parameters,
		Body:       s.Body,
		Env:        s.Env.Copy(),
		Pos:        s.Pos,
	}
}

func (s *Sub) Position() token.Position {
	return s.Pos
}

func (s *Sub) Equals(other Object) bool {
	otherSub, ok := other.(*Sub)
	if !ok {
		return false
	}
	return s.Definition.String() == otherSub.Definition.String()
}

func (l *Sub) GetValue() any {
	return nil
}

// ----------------------- Resume -----------------------
// in error handling, the resume object indicates the next statement to execute
type Resume struct {
	Label string
	Pos   token.Position
}

func NewResume(label string, pos token.Position) *Resume {
	return &Resume{Label: label, Pos: pos}
}

func (r *Resume) Type() ObjectType {
	return RESUME_OBJ
}

func (r *Resume) String() string {
	return r.Label
}

func (r *Resume) IsConstant() bool {
	return true
}

func (r *Resume) Copy() Object {
	return &Resume{Label: r.Label, Pos: r.Pos}
}

func (r *Resume) Position() token.Position {
	return r.Pos
}

func (r *Resume) Equals(other Object) bool {
	otherResume, ok := other.(*Resume)
	if !ok {
		return false
	}
	return r.Label == otherResume.Label
}

func (l *Resume) GetValue() any {
	return l.Label
}

// ----------------------- validation of types -----------------------
// check that variables implement the Variable interface
var (
	_ Object = &Boolean{}
	_ Object = &Integer{}
	_ Object = &Long{}
	_ Object = &Single{}
	_ Object = &Double{}
	_ Object = &Currency{}
	_ Object = &String{}
	_ Object = &Array{}
	_ Object = &Variant{}
	_ Object = &Date{}
	_ Object = &Error{}
	_ Object = &ReturnValue{}
	_ Object = &Nothing{}
	_ Object = &Function{}
	_ Object = &Sub{}
	_ Object = &Exit{}
	_ Object = &UserDefined{}
	_ Object = &Class{}
	_ Object = &Resume{}
)

// IsNil checks if an object is nil
func IsNil(obj Object) bool {
	switch obj := obj.(type) {
	case *Nothing:
		return obj == nil
	case *Error:
		return obj == nil
	case *ReturnValue:
		return obj == nil
	case *Function:
		return obj == nil
	case *Sub:
		return obj == nil
	case *Exit:
		return obj == nil
	case *UserDefined:
		return obj == nil
	case *Class:
		return obj == nil
	case *Resume:
		return obj == nil
	default:
		return false
	}
}

func GetBasicType(objType ObjectType) types.BasicKind {
	switch objType {
	case BOOLEAN_OBJ:
		return types.Boolean
	case LONG_OBJ:
		return types.Long
	case INTEGER_OBJ:
		return types.Integer
	case SINGLE_OBJ:
		return types.Single
	case DOUBLE_OBJ:
		return types.Double
	case CURRENCY_OBJ:
		return types.Currency
	case STRING_OBJ:
		return types.String
	case DATE_OBJ:
		return types.Date
	case NOTHING_OBJ:
		return types.Nothing
	case ERROR_OBJ:
		return types.Nothing
	case RETURN_OBJ:
		return types.Nothing
	case FUNCTION_OBJ:
		return types.Nothing
	case SUB_OBJ:
		return types.Nothing
	case RESUME_OBJ:
		return types.Nothing
	case USERDEF_OBJ:
		return types.Enum
	default:
		return types.Nothing
	}
}

func NewEmptyByKind(typ token.Kind, pos token.Position) Object {
	switch typ {
	case token.KwBoolean:
		return NewBooleanByBool(false, pos)
	case token.KwLong:
		return NewLongByInt(0, pos)
	case token.KwInteger:
		return NewIntegerByInt(0, pos)
	case token.KwSingle:
		return NewSingleByFloat(0.0, pos)
	case token.KwDouble:
		return NewDoubleByFloat(0.0, pos)
	case token.KwCurrency:
		return NewCurrencyByFloat(0.0, pos)
	case token.KwString:
		return NewString("", pos)
	case token.KwDate:
		return NewDateByTime(time.Time{}, pos)
	case token.KwNothing:
		return NOTHING
	case token.KwResume:
		return NewResume("", pos)
	default:
		return NOTHING
	}
}

func KindToType(kind token.Kind) ObjectType {
	switch kind {
	case token.KwBoolean:
		return BOOLEAN_OBJ
	case token.KwLong:
		return LONG_OBJ
	case token.KwInteger:
		return INTEGER_OBJ
	case token.KwSingle:
		return SINGLE_OBJ
	case token.KwDouble:
		return DOUBLE_OBJ
	case token.KwCurrency:
		return CURRENCY_OBJ
	case token.KwString:
		return STRING_OBJ
	case token.KwDate:
		return DATE_OBJ
	case token.KwNothing:
		return NOTHING_OBJ
	case token.KwResume:
		return RESUME_OBJ
	default:
		return NOTHING_OBJ
	}
}
func NewEmptyByType(typ ObjectType, pos token.Position) Object {
	switch typ {
	case BOOLEAN_OBJ:
		return NewBooleanByBool(false, pos)
	case LONG_OBJ:
		return NewLongByInt(0, pos)
	case INTEGER_OBJ:
		return NewIntegerByInt(0, pos)
	case SINGLE_OBJ:
		return NewSingleByFloat(0.0, pos)
	case DOUBLE_OBJ:
		return NewDoubleByFloat(0.0, pos)
	case CURRENCY_OBJ:
		return NewCurrencyByFloat(0.0, pos)
	case STRING_OBJ:
		return NewString("", pos)
	case DATE_OBJ:
		return NewDateByTime(time.Time{}, pos)
	case VARIANT_OBJ:
		return NewVariantByObject(NOTHING, pos)
	case NOTHING_OBJ:
		return NOTHING
	case RESUME_OBJ:
		return NewResume("", pos)
	default:
		return nil
	}
}
