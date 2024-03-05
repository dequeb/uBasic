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
	10.998$
	""
	"abc"
	"abc""def"`
	tests1 := []struct {
		expectedType    token.Kind
		expectedLiteral string
	}{
		{token.DateLit, "#1999/12/31 23:59:59#"},
		{token.EOL, "\n"},
		{token.DateLit, "#1999-12-31#"},
		{token.EOL, "\n"},
		{token.DateLit, "#23:59:59#"},
		{token.EOL, "\n"},
		{token.Illegal, "#"},
		{token.Identifier, "abc"},
		{token.Illegal, "#"},
		{token.EOL, "\n"},
		{token.CurrencyLit, "10$"},
		{token.Illegal, "#1333"},
		{token.EOL, "\n"},
		{token.CurrencyLit, "10.998$"},
		{token.EOL, "\n"},
		{token.StringLit, "\"\""},
		{token.EOL, "\n"},
		{token.StringLit, "\"abc\""},
		{token.EOL, "\n"},
		{token.StringLit, "\"abc\"\"def\""},
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
function a (x as single, y as integer) As Double
a = x + y	' comment
End Function 
' another comment
result = _
a(five, ten)`
	tests := []struct {
		expectedType    token.Kind
		expectedLiteral string
		expectedLine    int
		expectedColumn  int
	}{
		{token.Identifier, "five", 1, 1},
		{token.Eq, "=", 1, 6},
		{token.DoubleLit, "5.0", 1, 8},
		{token.EOL, "\n", 1, 11},
		{token.Identifier, "ten", 2, 1},
		{token.Eq, "=", 2, 5},
		{token.LongLit, "10", 2, 7},
		{token.EOL, "\n", 2, 9},
		{token.KwFunction, "function", 3, 1},
		{token.Identifier, "a", 3, 10},
		{token.Lparen, "(", 3, 12},
		{token.Identifier, "x", 3, 13},
		{token.KwAs, "as", 3, 15},
		{token.KwSingle, "single", 3, 18},
		{token.Comma, ",", 3, 24},
		{token.Identifier, "y", 3, 26},
		{token.KwAs, "as", 3, 28}, // in lowercase
		{token.KwInteger, "integer", 3, 31},
		{token.Rparen, ")", 3, 38},
		{token.KwAs, "As", 3, 40},
		{token.KwDouble, "Double", 3, 43},
		{token.EOL, "\n", 3, 49},
		{token.Identifier, "a", 4, 1},
		{token.Eq, "=", 4, 3},
		{token.Identifier, "x", 4, 5},
		{token.Add, "+", 4, 7},
		{token.Identifier, "y", 4, 9},
		{token.EOL, "\n", 4, 20},
		{token.KwEnd, "End", 5, 1},
		{token.KwFunction, "Function", 5, 5},
		{token.EOL, "\n", 5, 14},
		{token.EOL, "\n", 6, 18},
		{token.Identifier, "result", 7, 1},
		{token.Eq, "=", 7, 8},
		{token.Underscore, "_", 7, 10},
		{token.EOL, "\n", 7, 11},
		{token.Identifier, "a", 8, 1},
		{token.Lparen, "(", 8, 2},
		{token.Identifier, "five", 8, 3},
		{token.Comma, ",", 8, 7},
		{token.Identifier, "ten", 8, 9},
		{token.Rparen, ")", 8, 12},
		{token.EOF, "", 8, 13},
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

		if tok.Position.Line != tt.expectedLine || tok.Position.Column != tt.expectedColumn {
			t.Fatalf("tests[%d]=%s - tok.Position wrong. expected=%d,%d, got=%d,%d", i, tok.Literal, tt.expectedLine, tt.expectedColumn, tok.Position.Line, tok.Position.Column)
		}
	}
}
