package lexer

import (
	"github.com/iskandervdh/vorn/token"
)

type Lexer struct {
	input        string
	position     int  // current character position
	readPosition int  // current reading position in input (after current character)
	char         byte // current character under examination

	line   int // current line number
	column int // current column number
}

func New(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 1,
	}

	// Read a character to initialize l.char
	l.readChar()

	return l
}

func (l *Lexer) NextToken() token.Token {
	var t token.Token

	l.skipWhitespace()

	switch l.char {
	case '=':
		if l.peekChar() == '=' {
			ch := l.char
			l.readChar()
			literal := string(ch) + string(l.char)

			t = token.Token{
				Type:    token.EQ,
				Literal: literal,
				Line:    l.line,
				Column:  l.column,
			}
		} else {
			t = token.New(token.ASSIGN, l.char, l.line, l.column)
		}
	case '+':
		t = token.New(token.PLUS, l.char, l.line, l.column)
	case '-':
		t = token.New(token.MINUS, l.char, l.line, l.column)
	case '!':
		if l.peekChar() == '=' {
			char := l.char
			l.readChar()
			literal := string(char) + string(l.char)

			t = token.Token{
				Type:    token.NOT_EQ,
				Literal: literal,
				Line:    l.line,
				Column:  l.column,
			}
		} else {
			t = token.New(token.EXCLAMATION, l.char, l.line, l.column)
		}
	case '/':
		t = token.New(token.SLASH, l.char, l.line, l.column)
	case '*':
		t = token.New(token.ASTERISK, l.char, l.line, l.column)
	case '<':
		if l.peekChar() == '=' {
			char := l.char
			l.readChar()
			literal := string(char) + string(l.char)

			t = token.Token{
				Type:    token.LTE,
				Literal: literal,
				Line:    l.line,
				Column:  l.column,
			}
		} else {
			t = token.New(token.LT, l.char, l.line, l.column)
		}
	case '>':
		if l.peekChar() == '=' {
			char := l.char
			l.readChar()
			literal := string(char) + string(l.char)

			t = token.Token{
				Type:    token.GTE,
				Literal: literal,
				Line:    l.line,
				Column:  l.column,
			}
		} else {
			t = token.New(token.GT, l.char, l.line, l.column)
		}
	case ';':
		t = token.New(token.SEMICOLON, l.char, l.line, l.column)
	case ':':
		t = token.New(token.COLON, l.char, l.line, l.column)
	case ',':
		t = token.New(token.COMMA, l.char, l.line, l.column)
	case '{':
		t = token.New(token.LBRACE, l.char, l.line, l.column)
	case '}':
		t = token.New(token.RBRACE, l.char, l.line, l.column)
	case '(':
		t = token.New(token.LPAREN, l.char, l.line, l.column)
	case ')':
		t = token.New(token.RPAREN, l.char, l.line, l.column)
	case '"':
		t.Line = l.line
		t.Column = l.column
		t.Type = token.STRING
		t.Literal = l.readString()
	case '[':
		t = token.New(token.LBRACKET, l.char, l.line, l.column)
	case ']':
		t = token.New(token.RBRACKET, l.char, l.line, l.column)
	case 0:
		t.Line = l.line
		t.Column = l.column
		t.Literal = ""
		t.Type = token.EOF
	default:
		if isLetter(l.char) {
			t.Line = l.line
			t.Column = l.column
			t.Literal = l.readIdentifier()
			t.Type = token.LookupIdent(t.Literal)

			return t
		} else if isDigit(l.char) {
			t.Line = l.line
			t.Column = l.column
			t.Literal, t.Type = l.readNumber()

			return t
		} else {
			t = token.New(token.ILLEGAL, l.char, l.line, l.column)
		}
	}

	l.readChar()

	return t
}

func (l *Lexer) skipWhitespace() {
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r' {
		if l.char == '\n' {
			l.line += 1
			l.column = 1
		}

		l.readChar()
	}
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition += 1
	l.column += 1
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

func (l *Lexer) readIdentifier() string {
	position := l.position

	for isLetter(l.char) || isDigit(l.char) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readNumber() (string, token.TokenType) {
	position := l.position
	tokenType := token.TokenType(token.INT)

	for isDigit(l.char) || l.char == '.' {
		if l.char == '.' {
			tokenType = token.FLOAT
		}

		l.readChar()
	}

	return l.input[position:l.position], tokenType
}

func (l *Lexer) readString() string {
	position := l.position + 1

	for {
		l.readChar()

		if l.char == '"' || l.char == 0 {
			break
		}
	}

	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
