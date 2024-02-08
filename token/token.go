// Package token defines constants representing the lexical tokens of the µBASIC
// programming language.
package token

import (
	"fmt"

	"github.com/iancoleman/strcase"
)

// a position is a line and column number in the input text
type Position struct {
	Line   int
	Column int
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
	Ident       // main (also includes type names)
	LongLit     // 123
	DoubleLit   // 123.45
	StringLit   // "asdf;lkj"
	BooleanLit  // true, false
	DateTimeLit // 2019-01-01 00:00:00
	NothingLit  // nothing
	CurrencyLit // 55.00$
	literalEnd

	operatorStart

	// Operators and delimiters.
	Add        // +
	Concat     // &
	Sub        // -
	Mul        // *
	Div        // /
	Exponent   // ^
	Eq         // =
	Ne         // !=
	Lt         // <
	Le         // <=
	Gt         // >
	Ge         // >=
	Lparen     // (
	Rparen     // )
	Comma      // ,
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
	KwElse
	KwEnd
	KwErase
	KwExit
	KwFunction
	KwIf
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
	kwGoto
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
	keywordEnd
)

func (kind Kind) String() string {
	names := map[Kind]string{
		EOF:         "EOF",
		Illegal:     "error",
		Comment:     "comment",
		Ident:       "identifier",
		LongLit:     "integer literal",
		DoubleLit:   "floating-point literal",
		StringLit:   "string literal",
		DateTimeLit: "date-time literal",
		NothingLit:  "Nothing constant",
		BooleanLit:  "boolean literal",
		CurrencyLit: "currency literal",
		Add:         "+",
		Concat:      "&",
		Sub:         "-",
		Mul:         "*",
		Div:         "/",
		IntDiv:      "Div",
		Mod:         "Mod",
		Exponent:    "Exp",
		Eq:          "=",
		Ne:          "<>",
		Lt:          "<",
		Le:          "<=",
		Gt:          ">",
		Ge:          ">=",
		And:         "And",
		Or:          "Or",
		Not:         "Not ",
		Lparen:      "(",
		Rparen:      ")",
		Comma:       ",",
		Underscore:  "_",
		KwAs:        "As",
		KwConst:     "Const",
		KwDim:       "Dim",
		KwDo:        "Do",
		KwElse:      "Else",
		KwEnd:       "End",
		KwErase:     "Erase",
		KwExit:      "Exit",
		KwFunction:  "Function",
		KwIf:        "If",
		KwIn:        "In",
		KwRedim:     "Redim",
		KwSub:       "Sub",
		KwThen:      "Then",
		KwUntil:     "Until",
		KwWhile:     "While",
		KwSelect:    "Select",
		KwCase:      "Case",
		KwStop:      "Stop",
		KwOn:        "On",
		KwError:     "Error",
		KwResume:    "Resume",
		kwGoto:      "Goto",
		KwLong:      "Long",
		KwInteger:   "Integer",
		KwSingle:    "Single",
		KwDouble:    "Double",
		KwString:    "String",
		KwBoolean:   "Boolean",
		KwDate:      "Date",
		KwVariant:   "Variant",
		KwNothing:   "Nothing",
		KwTrue:      "True",
		KwFalse:     "False",
		KwCurrency:  "Currency",
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
	ident = strcase.ToCamel(ident)
	if kind, ok := Keywords[ident]; ok {
		return kind
	}
	return Ident
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
	"As":       KwAs,
	"Const":    KwConst,
	"Dim":      KwDim,
	"Do":       KwDo,
	"Else":     KwElse,
	"End":      KwEnd,
	"Erase":    KwErase,
	"Exit":     KwExit,
	"Function": KwFunction,
	"If":       KwIf,
	"In":       KwIn,
	"Redim":    KwRedim,
	"Sub":      KwSub,
	"Then":     KwThen,
	"Until":    KwUntil,
	"While":    KwWhile,
	"Select":   KwSelect,
	"Case":     KwCase,
	"Stop":     KwStop,
	"On":       KwOn,
	"Error":    KwError,
	"Resume":   KwResume,
	"Goto":     kwGoto,
	"Long":     KwLong,
	"Integer":  KwInteger,
	"Single":   KwSingle,
	"Double":   KwDouble,
	"String":   KwString,
	"Boolean":  KwBoolean,
	"Date":     KwDate,
	"Variant":  KwVariant,
	"Nothing":  KwNothing,
	"True":     KwTrue,
	"False":    KwFalse,
	"Currency": KwCurrency,
}
