package lexer

import (
	"time"
	"uBasic/token"
)

type Lexer struct {
	input          string
	lcPosition     token.Position
	lcReadPosition token.Position
	ch             byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.lcReadPosition.Line = 1
	l.lcReadPosition.Column = 0
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.lcReadPosition.Absolute >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.lcReadPosition.Absolute]
	}
	l.lcPosition.Absolute = l.lcReadPosition.Absolute
	l.lcReadPosition.Absolute++
	l.lcReadPosition.Column++
}

func (l *Lexer) Rewind() {
	l.lcReadPosition.Absolute--
	l.lcReadPosition.Column--
}

func (l *Lexer) NextToken() *token.Token {
	var tok token.Token

	l.skipWhitespace()

	l.lcPosition.Line = l.lcReadPosition.Line
	l.lcPosition.Column = l.lcReadPosition.Column

	switch l.ch {
	case token.Dot.String()[0]:
		tok = newToken(token.Dot, l.ch, l.lcPosition)
	case token.Colon.String()[0]:
		tok = newToken(token.Colon, l.ch, l.lcPosition)
	case token.Semicolon.String()[0]:
		tok = newToken(token.Semicolon, l.ch, l.lcPosition)
	case token.Add.String()[0]:
		tok = newToken(token.Add, l.ch, l.lcPosition)
	case token.Minus.String()[0]:
		tok = newToken(token.Minus, l.ch, l.lcPosition)
	case token.Mul.String()[0]:
		tok = newToken(token.Mul, l.ch, l.lcPosition)
	case token.Div.String()[0]:
		tok = newToken(token.Div, l.ch, l.lcPosition)
	case token.Lparen.String()[0]:
		tok = newToken(token.Lparen, l.ch, l.lcPosition)
	case token.Rparen.String()[0]:
		tok = newToken(token.Rparen, l.ch, l.lcPosition)
	case token.Comma.String()[0]:
		tok = newToken(token.Comma, l.ch, l.lcPosition)
	case token.Underscore.String()[0]:
		tok = newToken(token.Underscore, l.ch, l.lcPosition)
	case token.Concat.String()[0]:
		tok = newToken(token.Concat, l.ch, l.lcPosition)
	case token.Eq.String()[0]:
		if l.peekChar() == '=' {
			tok.Literal = "=="
			tok.Kind = token.Eq
			tok.Position = l.lcPosition.Copy()
			l.readChar()
			l.readChar()
			return &tok
		}
		tok = newToken(token.Assign, l.ch, l.lcPosition)
	case token.Lt.String()[0]:
		if l.peekChar() == '=' {
			tok.Literal = "<="
			tok.Kind = token.Le
			tok.Position = l.lcPosition.Copy()
			l.readChar()
			l.readChar()
			return &tok
		}
		tok = newToken(token.Lt, l.ch, l.lcPosition)
		if l.peekChar() == '>' {
			tok.Literal = "<>"
			tok.Kind = token.Neq
			tok.Position = l.lcPosition.Copy()
			l.readChar()
			l.readChar()
			return &tok
		}
		tok = newToken(token.Lt, l.ch, l.lcPosition)
	case token.Gt.String()[0]:
		if l.peekChar() == '=' {
			tok.Literal = ">="
			tok.Kind = token.Ge
			tok.Position = l.lcPosition.Copy()
			l.readChar()
			l.readChar()
			return &tok
		}
		tok = newToken(token.Gt, l.ch, l.lcPosition)
	case '\n', '\r':
		tok = newToken(token.EOL, l.ch, l.lcPosition)
		l.lcReadPosition.Line++
		l.lcReadPosition.Column = 0
	case 0:
		tok.Literal = ""
		tok.Kind = token.EOF
		tok.Position = l.lcPosition.Copy()
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Kind = token.LookupIdent(tok.Literal)
			tok.Position = l.lcPosition.Copy()
			return &tok
		} else if isDigit(l.ch) {
			tok.Kind = token.LongLit
			tok.Literal = l.readNumber()
			tok.Position = l.lcPosition.Copy()
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
			return &tok
		} else if l.ch == '#' {
			// read a date time
			l.readChar() // skip the first #
			tok.Kind = token.DateLit
			tok.Literal = "#" + l.readDateTime() + "#"
			tok.Position = l.lcPosition.Copy()
			if isValidDateTime(tok.Literal) {
				l.readChar() // skip the last #
				return &tok
			} else {
				if l.ch == '#' {
					l.readChar() // skip the last #
				} else {
					// rewind to the last character
					l.Rewind()
					tok.Literal = tok.Literal[:len(tok.Literal)-1]
				}

				tok.Kind = token.Illegal
				tok.Position = l.lcPosition.Copy()
			}
		} else if l.ch == '"' {
			// read a string
			l.readChar() // skip the first "
			tok.Kind = token.StringLit
			tok.Literal = "\"" + l.readString() + "\""
			tok.Position = l.lcPosition.Copy()
			if l.ch == '"' {
				l.readChar() // skip the last "
			}
			return &tok
		} else if l.ch == '\'' {
			// comment: skip the rest of the line
			for l.ch != '\n' && l.ch != '\r' && l.ch != 0 {
				l.readChar()
			}
			// ajust the position
			// l.lcReadPosition.Line++
			// l.lcReadPosition.Column = 0
			return l.NextToken()
		} else {
			tok = newToken(token.Illegal, l.ch, l.lcPosition)
		}
	}
	l.readChar()

	return &tok
}

func newToken(kind token.Kind, ch byte, pos token.Position) token.Token {
	return token.Token{Kind: kind, Literal: string(ch), Position: token.Position{Line: pos.Line, Column: pos.Column, Absolute: pos.Absolute}}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func isIdentifer(ch byte) bool {
	return isLetter(ch) || ch == '_' || isDigit(ch)
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) readIdentifier() string {
	position := l.lcPosition.Absolute
	for isIdentifer(l.ch) {
		l.readChar()
	}
	return l.input[position:l.lcPosition.Absolute]
}

func (l *Lexer) readNumber() string {
	position := l.lcPosition.Absolute
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.lcPosition.Absolute]
}

func (l *Lexer) peekChar() byte {
	if l.lcReadPosition.Absolute >= len(l.input) {
		return 0
	}
	return l.input[l.lcReadPosition.Absolute]
}

func (l *Lexer) readDateTime() string {
	position := l.lcPosition.Absolute
	for isDateTimeChar(l.ch) {
		l.readChar()
	}
	return l.input[position:l.lcPosition.Absolute]
}

func (l *Lexer) readString() string {
	position := l.lcPosition.Absolute
	for l.isStringChar(l.ch) && l.ch != 0 {
		l.readChar()
	}
	length := l.lcPosition.Absolute
	if l.ch == 0 {
		length--
	}
	return l.input[position:l.lcPosition.Absolute]
}

func (l *Lexer) isStringChar(ch byte) bool {
	res := ch != '"' || l.peekChar() == '"'
	if ch == '"' && l.peekChar() == '"' {
		l.readChar() // skip the next "
	}
	return res
}

func isDateTimeChar(ch byte) bool {
	return '0' <= ch && ch <= '9' || ch == '-' || ch == ':' || ch == '/' || ch == ' '
}

func isValidDateTime(s string) bool {
	s = s[1 : len(s)-1]
	// List of date formats
	dateFormats := []string{
		"2006-01-02",
		"2006/01/02",
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
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\f' {
		l.readChar()
	}
}
