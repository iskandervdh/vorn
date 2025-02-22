package token

import (
	"testing"
)

func TestTokenNew(t *testing.T) {
	token := New(EOF, 0, 1, 1)

	if token.Type != EOF {
		t.Errorf("Expected token type to be %s, got %s", EOF, token.Type)
	}

	if token.Line != 1 {
		t.Errorf("Expected token line to be 1, got %d", token.Line)
	}

	if token.Column != 1 {
		t.Errorf("Expected token column to be 1, got %d", token.Column)
	}

	if token.Literal != "\x00" {
		t.Errorf("Expected token literal to be empty, got %q", token.Literal)
	}

	token = New(IDENT, 't', 1, 1)

	if token.Type != IDENT {
		t.Errorf("Expected token type to be %s, got %s", IDENT, token.Type)
	}

	if token.Line != 1 {
		t.Errorf("Expected token line to be 1, got %d", token.Line)
	}

	if token.Column != 1 {
		t.Errorf("Expected token column to be 1, got %d", token.Column)
	}

	if token.Literal != "t" {
		t.Errorf("Expected token literal to be 't', got %q", token.Literal)
	}
}

func TestLookupIdent(t *testing.T) {
	ident := LookupIdent("func")

	if ident != FUNCTION {
		t.Errorf("Expected ident to be %s, got %s", FUNCTION, ident)
	}

	ident = LookupIdent("const")

	if ident != CONST {
		t.Errorf("Expected ident to be %s, got %s", CONST, ident)
	}

	ident = LookupIdent("not_a_keyword")

	if ident != IDENT {
		t.Errorf("Expected ident to be %s, got %s", IDENT, ident)
	}
}
