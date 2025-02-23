package ast

import (
	"testing"

	"github.com/iskandervdh/vorn/token"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&VariableStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "someVar"},
					Value: "someVar",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}
	if program.String() != "let someVar = anotherVar;" {
		t.Errorf("program.String() wrong. got '%q'", program.String())
	}
}
