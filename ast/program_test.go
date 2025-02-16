package ast

import (
	"testing"

	"github.com/iskandervdh/vorn/token"
)

func TestEmptyProgram(t *testing.T) {
	program := NewProgram()

	if program == nil {
		t.Fatalf("NewProgram() returned nil")
	}

	if len(program.Statements) != 0 {
		t.Fatalf("program.Statements has wrong length. got %d", len(program.Statements))
	}
}

func TestSimpleProgram(t *testing.T) {
	program := NewProgram()

	if program.TokenLiteral() != "" {
		t.Fatalf("expected empty string, got '%s'", program.TokenLiteral())
	}

	program.Statements = append(program.Statements, &VariableStatement{
		Token: token.Token{Type: token.LET, Literal: "let"},
		Name: &Identifier{
			Token: token.Token{Type: token.IDENT, Literal: "x"},
			Value: "x",
		},
	})
	program.Statements = append(program.Statements, &VariableStatement{
		Token: token.Token{Type: token.CONST, Literal: "const"},
		Name: &Identifier{
			Token: token.Token{Type: token.IDENT, Literal: "y"},
			Value: "y",
		},
	})

	if len(program.Statements) != 2 {
		t.Fatalf("program.Statements has wrong length. got %d", len(program.Statements))
	}

	if program.Statements[0] == nil {
		t.Fatalf("program.Statements[0] is nil")
	}

	if program.TokenLiteral() != "let" {
		t.Fatalf("program.TokenLiteral() wrong. got '%s'", program.TokenLiteral())
	}

	if !program.Statements[0].(*VariableStatement).IsLet() {
		t.Fatalf("program.Statements[0] is not a let statement")
	}

	if !program.Statements[1].(*VariableStatement).IsConst() {
		t.Fatalf("program.Statements[1] is not a const statement")
	}
}
