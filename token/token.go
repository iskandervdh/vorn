package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string

	Line   int
	Column int
}

const (
	// Special tokens
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
	COMMENT = "COMMENT"

	// Identifiers + literals
	IDENT  = "IDENT"
	INT    = "INT"
	FLOAT  = "FLOAT"
	STRING = "STRING"

	// Operators
	ASSIGN      = "="
	PLUS        = "+"
	MINUS       = "-"
	EXCLAMATION = "!"
	ASTERISK    = "*"
	SLASH       = "/"
	PERCENT     = "%"

	// Comparison operators
	LT  = "<"
	GT  = ">"
	LTE = "<="
	GTE = ">="

	EQ     = "=="
	NOT_EQ = "!="

	// Logical operators
	OR  = "||"
	AND = "&&"

	// Delimiters
	DOT       = "."
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	FUNCTION = "FUNCTION"
	CONST    = "CONST"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	NULL     = "NULL"
	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	WHILE    = "WHILE"
	FOR      = "FOR"
	BREAK    = "BREAK"
	CONTINUE = "CONTINUE"
)

/*
Keywords are stored in a map where the key is the keyword and the value is the TokenType.

When we encounter a keyword in the source code, we can look it up in the map to determine its TokenType.
*/
var keywords = map[string]TokenType{
	"func":     FUNCTION,
	"const":    CONST,
	"let":      LET,
	"true":     TRUE,
	"false":    FALSE,
	"null":     NULL,
	"if":       IF,
	"else":     ELSE,
	"return":   RETURN,
	"while":    WHILE,
	"for":      FOR,
	"break":    BREAK,
	"continue": CONTINUE,
}

/*
Create a new token with the given TokenType, character, line and column.

Returns the new token.
*/
func New(tokenType TokenType, ch byte, line int, column int) Token {
	return Token{Type: tokenType, Literal: string(ch), Line: line, Column: column}
}

/*
Lookup the TokenType of a given identifier if it is a keyword.

Returns the TokenType of the identifier. If the identifier is not a keyword, return IDENT.
*/
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}
