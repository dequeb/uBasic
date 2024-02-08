package lexer

import (
	"testing"
	"uBasic/token"
)

func TestNextToken1(t *testing.T) {
	input1 := `=+(),_`
	tests1 := []struct {
		expectedType    token.Kind
		expectedLiteral string
	}{
		{token.Eq, "="},
		{token.Add, "+"},
		{token.Lparen, "("},
		{token.Rparen, ")"},
		{token.Comma, ","},
		{token.Underscore, "_"},
		{token.EOF, ""},
	}

	l := New(input1)
	for i, tt := range tests1 {
		tok := l.NextToken()
		if tok.Kind != tt.expectedType {
			t.Fatalf("tests[%d] - tok.Kind wrong. expected=%q, got=%qat position %d,%d", i, tt.expectedLiteral, tok.Literal, tok.Position.Line, tok.Position.Column)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - tok.Literal wrong. expected=%q, got=%qat position %d,%d", i, tt.expectedLiteral, tok.Literal, tok.Position.Line, tok.Position.Column)
		}
	}
}

func TestNextToken2(t *testing.T) {
	input1 := `#1999/12/31 23:59:59#
	#1999-12-31#
	#23:59:59#
	#abc#
	10$#1333
	10.998$`
	tests1 := []struct {
		expectedType    token.Kind
		expectedLiteral string
	}{
		{token.DateTimeLit, "#1999/12/31 23:59:59#"},
		{token.EOL, "\n"},
		{token.DateTimeLit, "#1999-12-31#"},
		{token.EOL, "\n"},
		{token.DateTimeLit, "#23:59:59#"},
		{token.EOL, "\n"},
		{token.Illegal, "#"},
		{token.Ident, "abc"},
		{token.Illegal, "#"},
		{token.EOL, "\n"},
		{token.CurrencyLit, "10$"},
		{token.Illegal, "#1333"},
		{token.EOL, "\n"},
		{token.CurrencyLit, "10.998$"},
		{token.EOF, ""},
	}
	l := New(input1)
	for i, tt := range tests1 {
		tok := l.NextToken()
		if tok.Kind != tt.expectedType {
			t.Fatalf("tests[%d] - tok.Kind wrong. expected=%q, got=%qat position %d,%d", i, tt.expectedLiteral, tok.Literal, tok.Position.Line, tok.Position.Column)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - tok.Literal wrong. expected=%q, got=%qat position %d,%d", i, tt.expectedLiteral, tok.Literal, tok.Position.Line, tok.Position.Column)
		}
	}
}

func TestNextToken(t *testing.T) {
	input := `five = 5.0
	ten = 10
	Function a (x As Single, y as Integer) As Double
		a = x + y
	End Function 
	result = _
	a(five, ten)`
	tests := []struct {
		expectedType    token.Kind
		expectedLiteral string
	}{
		{token.Ident, "five"},
		{token.Eq, "="},
		{token.DoubleLit, "5.0"},
		{token.EOL, "\n"},
		{token.Ident, "ten"},
		{token.Eq, "="},
		{token.LongLit, "10"},
		{token.EOL, "\n"},
		{token.KwFunction, "Function"},
		{token.Ident, "a"},
		{token.Lparen, "("},
		{token.Ident, "x"},
		{token.KwAs, "As"},
		{token.KwSingle, "Single"},
		{token.Comma, ","},
		{token.Ident, "y"},
		{token.KwAs, "as"}, // in lowercase
		{token.KwInteger, "Integer"},
		{token.Rparen, ")"},
		{token.KwAs, "As"},
		{token.KwDouble, "Double"},
		{token.EOL, "\n"},
		{token.Ident, "a"},
		{token.Eq, "="},
		{token.Ident, "x"},
		{token.Add, "+"},
		{token.Ident, "y"},
		{token.EOL, "\n"},
		{token.KwEnd, "End"},
		{token.KwFunction, "Function"},
		{token.EOL, "\n"},
		{token.Ident, "result"},
		{token.Eq, "="},
		{token.Underscore, "_"},
		{token.EOL, "\n"},
		{token.Ident, "a"},
		{token.Lparen, "("},
		{token.Ident, "five"},
		{token.Comma, ","},
		{token.Ident, "ten"},
		{token.Rparen, ")"},
		{token.EOF, ""},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Kind != tt.expectedType {
			t.Fatalf("tests[%d] - tok.Kind wrong. expected=%q, got=%qat position %d,%d", i, tt.expectedLiteral, tok.Literal, tok.Position.Line, tok.Position.Column)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - tok.Literal wrong. expected=%q, got=%q at position %d,%d", i, tt.expectedLiteral, tok.Literal, tok.Position.Line, tok.Position.Column)
		}
	}
}
