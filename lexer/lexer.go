package lexer

import (
	"time"
	"uBasic/token"
)

type Lexer struct {
	input          string
	lcPosition     token.Position
	lcReadPosition token.Position
	position       int
	readPosition   int
	ch             byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.lcReadPosition.Line = 1
	l.lcReadPosition.Column = 1
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.lcPosition.Line = l.lcReadPosition.Line
	l.lcPosition.Column = l.lcReadPosition.Column
	l.readPosition++
	l.lcReadPosition.Column++
}

func (l *Lexer) Rewind() {
	l.readPosition--
	l.lcReadPosition.Column--
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case token.Add.String()[0]:
		tok = newToken(token.Add, l.ch, l.lcReadPosition)
	case token.Sub.String()[0]:
		tok = newToken(token.Sub, l.ch, l.lcReadPosition)
	case token.Mul.String()[0]:
		tok = newToken(token.Mul, l.ch, l.lcReadPosition)
	case token.Div.String()[0]:
		tok = newToken(token.Div, l.ch, l.lcReadPosition)
	case token.Lparen.String()[0]:
		tok = newToken(token.Lparen, l.ch, l.lcReadPosition)
	case token.Rparen.String()[0]:
		tok = newToken(token.Rparen, l.ch, l.lcReadPosition)
	case token.Comma.String()[0]:
		tok = newToken(token.Comma, l.ch, l.lcReadPosition)
	case token.Underscore.String()[0]:
		tok = newToken(token.Underscore, l.ch, l.lcReadPosition)
	case token.Concat.String()[0]:
		tok = newToken(token.Concat, l.ch, l.lcReadPosition)
	case token.Eq.String()[0]:
		tok = newToken(token.Eq, l.ch, l.lcReadPosition)
	case token.Lt.String()[0]:
		if l.peekChar() == '=' {
			tok.Literal = "<="
			tok.Kind = token.Le
			tok.Position = l.lcPosition
			l.readChar()
			return tok
		}
		tok = newToken(token.Lt, l.ch, l.lcReadPosition)
		if l.peekChar() == '>' {
			tok.Literal = "<>"
			tok.Kind = token.Ne
			tok.Position = l.lcPosition
			l.readChar()
			return tok
		}
		tok = newToken(token.Lt, l.ch, l.lcReadPosition)
	case token.Gt.String()[0]:
		if l.peekChar() == '=' {
			tok.Literal = ">="
			tok.Kind = token.Ge
			tok.Position = l.lcPosition
			l.readChar()
			return tok
		}
		tok = newToken(token.Gt, l.ch, l.lcReadPosition)
	case '\n':
		tok = newToken(token.EOL, l.ch, l.lcReadPosition)
		l.lcReadPosition.Line++
		l.lcReadPosition.Column = 1
	case 0:
		tok.Literal = ""
		tok.Kind = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Kind = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Kind = token.LongLit
			tok.Literal = l.readNumber()
			// is it a long, a double or a currency?
			if l.ch == '$' {
				l.readChar()
				tok.Literal += "$"
				tok.Kind = token.CurrencyLit
			} else if l.ch == '.' {
				l.readChar()
				tok.Kind = token.DoubleLit
				tok.Literal += "." + l.readNumber()
				if l.ch == '$' {
					l.readChar()
					tok.Kind = token.CurrencyLit
					tok.Literal += "$"
				}
			}
			return tok
		} else if l.ch == '#' {
			l.readChar() // skip the first #
			tok.Kind = token.DateTimeLit
			tok.Literal = "#" + l.readDateTime() + "#"
			if isValidDateTime(tok.Literal) {
				l.readChar() // skip the last #
				return tok
			} else {
				l.Rewind()
				tok.Literal = tok.Literal[:len(tok.Literal)-1] // remove the # we put in advance
				tok.Kind = token.Illegal
			}
		} else {
			tok = newToken(token.Illegal, l.ch, l.lcReadPosition)
		}
	}
	l.readChar()

	return tok
}

func newToken(kind token.Kind, ch byte, pos token.Position) token.Token {
	return token.Token{Kind: kind, Literal: string(ch), Position: token.Position{Line: pos.Line, Column: pos.Column}}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

func (l *Lexer) readDateTime() string {
	position := l.position
	for isDateTimeChar(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isDateTimeChar(ch byte) bool {
	return '0' <= ch && ch <= '9' || ch == '-' || ch == ':' || ch == '/' || ch == ' '
}

func isValidDateTime(s string) bool {
	s = s[1 : len(s)-1]
	// List of date formats
	dateFormats := []string{
		"2006-01-02",
		"20061/01/02",
		"2006-01-02 15:04:05",
		"2006/01/02 15:04:05",
		"15:04:05",
	}

	for _, format := range dateFormats {
		_, err := time.Parse(format, s)
		if err == nil {
			return true
		}
	}
	return false
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.readChar()
	}
}
