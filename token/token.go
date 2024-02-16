// Package token defines constants representing the lexical tokens of the µBASIC
// programming language.
package token

import (
	"fmt"
	"strings"
)

// a position is a line and column number in the input text
type Position struct {
	fmt.Stringer
	Line     int
	Column   int
	Absolute int
}

func (p *Position) Copy() Position {
	return Position{Line: p.Line, Column: p.Column, Absolute: p.Absolute}
}

// String returns the string representation of the token.
func (pos *Position) String() string {
	return fmt.Sprintf("line: %d, col: %d", pos.Line, pos.Column)
}

// A Token represents a lexical token of the µBAISC programming language.
type Token struct {
	// The token type.
	Kind
	// The string value of the token.
	Literal string
	// Start position in the input string.
	Position Position
}

func (tok Token) String() string {
	return fmt.Sprintf(`token.Token{Kind: token.%v, Val: %q, Pos: %v}`, tok.Kind.String(), tok.Literal, tok.Position)
}

// Equals reports whether tok and other are the same token.
func (tok Token) Equals(other *Token) bool {
	if other == nil {
		return false
	}
	return tok.Kind == other.Kind && tok.Literal == other.Literal && tok.Position.Equals(&other.Position)
}

// Equals reports whether pos and other are the same position.
func (pos Position) Equals(other *Position) bool {
	if other == nil {
		return false
	}
	return pos.Line == other.Line && pos.Column == other.Column && pos.Absolute == other.Absolute
}

//go:generate stringer -type Kind
//go:generate gorename -from kind_string.go::i -to kind
//go:generate gorename -from kind_string.go::String -to GoString

// Kind is the set of lexical token types of the µC programming language.
type Kind uint16

// Token types.
const (
	// Special tokens.
	EOF     Kind = iota // End of file
	Illegal             // Token value holds an error message (e.g. unterminated string)
	Comment             // /* block comment */ or // line comment

	literalStart

	// Identifiers and basic literals.
	Identifier  // main (also includes type names)
	LongLit     // 123
	DoubleLit   // 123.45
	StringLit   // "asdf;lkj"
	BooleanLit  // true, false
	DateLit     // 2019-01-01 00:00:00
	CurrencyLit // 55.00$
	literalEnd

	operatorStart

	// Operators and delimiters.
	Add        // +
	Concat     // &
	Minus      // -
	Mul        // *
	Div        // /
	Exponent   // ^
	Eq         // ==
	Assign     // =
	Neq        // !=
	Lt         // <
	Le         // <=
	Gt         // >
	Ge         // >=
	Colon      // :
	Semicolon  // ;
	Lparen     // (
	Rparen     // )
	Comma      // ,
	Dot        // .
	Underscore // _
	EOL        // End of line

	operatorEnd

	keywordStart

	// Keywords.
	IntDiv // Div
	Mod    // Mod
	And    // And
	Or     // Or
	Not    // Not
	KwAs
	KwConst
	KwDim
	KwDo
	KwLoop
	KwElse
	KwEnd
	KwEnum
	KwErase
	KwExit
	KwFunction
	KwIf
	KwElseIf
	KwIn
	KwRedim
	KwSub
	KwThen
	KwUntil
	KwWhile
	KwSelect
	KwCase
	KwStop
	KwOn
	KwError
	KwResume
	KwGoto
	KwLong
	KwInteger
	KwSingle
	KwDouble
	KwString
	KwBoolean
	KwDate
	KwVariant
	KwNothing
	KwTrue
	KwFalse
	KwCurrency
	KwFor
	KwNext
	KwStep
	KwTo
	KwEach
	KwByVal
	KwByRef
	KwOptional
	KwParamArray
	KwPreserve
	KwCall
	KwLet

	keywordEnd
)

func (kind Kind) String() string {
	names := map[Kind]string{
		EOF:          "EOF",
		Illegal:      "error",
		Comment:      "comment",
		Identifier:   "identifier",
		LongLit:      "integer literal",
		DoubleLit:    "floating-point literal",
		StringLit:    "string literal",
		DateLit:      "date-time literal",
		BooleanLit:   "boolean literal",
		CurrencyLit:  "currency literal",
		Add:          "+",
		Concat:       "&",
		Minus:        "-",
		Mul:          "*",
		Div:          "/",
		IntDiv:       "Div",
		Mod:          "Mod",
		Exponent:     "Exp",
		Eq:           "==",
		Assign:       "=",
		Neq:          "<>",
		Lt:           "<",
		Le:           "<=",
		Gt:           ">",
		Ge:           ">=",
		Colon:        ":",
		Semicolon:    ";",
		And:          "And",
		Or:           "Or",
		Not:          "Not ",
		Lparen:       "(",
		Rparen:       ")",
		Comma:        ",",
		Dot:          ".",
		Underscore:   "_",
		KwAs:         "As",
		KwConst:      "Const",
		KwDim:        "Dim",
		KwDo:         "Do",
		KwLoop:       "Loop",
		KwElse:       "Else",
		KwEnd:        "End",
		KwEnum:       "Enum",
		KwErase:      "Erase",
		KwExit:       "Exit",
		KwFunction:   "Function",
		KwIf:         "If",
		KwElseIf:     "ElseIf",
		KwIn:         "In",
		KwRedim:      "Redim",
		KwSub:        "Sub",
		KwThen:       "Then",
		KwUntil:      "Until",
		KwWhile:      "While",
		KwSelect:     "Select",
		KwCase:       "Case",
		KwStop:       "Stop",
		KwOn:         "On",
		KwError:      "Error",
		KwResume:     "Resume",
		KwGoto:       "Goto",
		KwLong:       "Long",
		KwInteger:    "Integer",
		KwSingle:     "Single",
		KwDouble:     "Double",
		KwString:     "String",
		KwBoolean:    "Boolean",
		KwDate:       "Date",
		KwVariant:    "Variant",
		KwNothing:    "Nothing",
		KwTrue:       "True",
		KwFalse:      "False",
		KwCurrency:   "Currency",
		KwFor:        "For",
		KwNext:       "Next",
		KwStep:       "Step",
		KwTo:         "To",
		KwEach:       "Each",
		KwByVal:      "ByVal",
		KwByRef:      "ByRef",
		KwOptional:   "Optional",
		KwParamArray: "ParamArray",
		KwPreserve:   "Preserve",
		KwCall:       "Call",
		KwLet:        "Let",
	}
	return names[kind]
}

// IsKeyword reports whether kind is a keyword.
func (kind Kind) IsKeyword() bool {
	return keywordStart < kind && kind < keywordEnd
}

func IsKeyword(str string) bool {
	_, ok := Keywords[str]
	return ok
}

func LookupIdent(ident string) Kind {
	ident = strings.ToLower(ident)
	if kind, ok := Keywords[ident]; ok {
		return kind
	}
	return Identifier
}

// IsLiteral reports whether kind is an identifier or a basic literal.
func (kind Kind) IsLiteral() bool {
	return literalStart < kind && kind < literalEnd
}

// IsOperator reports whether kind is an operator or a delimiter.
func (kind Kind) IsOperator() bool {
	return operatorStart < kind && kind < operatorEnd
}

// Keywords is the set of valid keywords in the µBASIC programming language.
var Keywords = map[string]Kind{
	"as":         KwAs,
	"const":      KwConst,
	"dim":        KwDim,
	"do":         KwDo,
	"loop":       KwLoop,
	"else":       KwElse,
	"end":        KwEnd,
	"enum":       KwEnum,
	"erase":      KwErase,
	"exit":       KwExit,
	"function":   KwFunction,
	"if":         KwIf,
	"elseif":     KwElseIf,
	"in":         KwIn,
	"redim":      KwRedim,
	"sub":        KwSub,
	"then":       KwThen,
	"until":      KwUntil,
	"while":      KwWhile,
	"select":     KwSelect,
	"case":       KwCase,
	"stop":       KwStop,
	"on":         KwOn,
	"error":      KwError,
	"resume":     KwResume,
	"goto":       KwGoto,
	"long":       KwLong,
	"integer":    KwInteger,
	"single":     KwSingle,
	"double":     KwDouble,
	"string":     KwString,
	"boolean":    KwBoolean,
	"date":       KwDate,
	"variant":    KwVariant,
	"nothing":    KwNothing,
	"true":       KwTrue,
	"false":      KwFalse,
	"currency":   KwCurrency,
	"to":         KwTo,
	"step":       KwStep,
	"each":       KwEach,
	"for":        KwFor,
	"next":       KwNext,
	"byval":      KwByVal,
	"byref":      KwByRef,
	"optional":   KwOptional,
	"paramarray": KwParamArray,
	"preserve":   KwPreserve,
	"call":       KwCall,
	"let":        KwLet,
	"and":        And,
	"or":         Or,
	"not":        Not,
	"div":        IntDiv,
	"mod":        Mod,
	"exp":        Exponent,
}
