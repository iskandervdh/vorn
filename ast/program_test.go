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

func TestProgramLocation(t *testing.T) {
	program := NewProgram()

	if program.Line() != 0 {
		t.Fatalf("program.Line() wrong. got %d", program.Line())
	}

	if program.Column() != 0 {
		t.Fatalf("program.Column() wrong. got %d", program.Column())
	}

	program.Statements = []Statement{&VariableStatement{
		Token: token.Token{
			Type:    token.LET,
			Literal: "let",
			Line:    1,
			Column:  1,
		},
		Name: &Identifier{
			Token: token.Token{
				Type:    token.IDENT,
				Literal: "x",
				Line:    1,
				Column:  5,
			},
			Value: "x",
		},
		Value: &IntegerLiteral{
			Token: token.Token{
				Type:    token.INT,
				Literal: "5",
				Line:    1,
				Column:  7,
			},
			Value: 5,
		},
	}}

	if program.Line() != 1 {
		t.Fatalf("program.Line() wrong. got %d", program.Line())
	}

	if program.Column() != 1 {
		t.Fatalf("program.Column() wrong. got %d", program.Column())
	}

	if program.Line() != program.Statements[0].Line() {
		t.Fatalf("program.Line() wrong. got %d", program.Line())
	}

	if program.Column() != program.Statements[0].Column() {
		t.Fatalf("program.Column() wrong. got %d", program.Column())
	}
}

func TestProgramScope(t *testing.T) {
	program := NewProgram()

	if program.GetParentScope() != nil {
		t.Fatalf("program.GetParentScope() is not nil")
	}

	program.Statements = []Statement{&VariableStatement{
		Token: token.Token{
			Type:    token.LET,
			Literal: "let",
		},
		Name: &Identifier{
			Token: token.Token{
				Type:    token.IDENT,
				Literal: "x",
			},
			Value: "x",
		},
	}}

	if len(program.GetScopeStatements()) != len(program.Statements) {
		t.Fatalf("program.GetScopeStatements() wrong length. got %d, want %d", len(program.GetScopeStatements()), len(program.Statements))
	}

	for i, statement := range program.GetScopeStatements() {
		if statement != program.Statements[i] {
			t.Fatalf("program.GetScopeStatements() wrong statement at index %d. got %v, want %v", i, statement, program.Statements[i])
		}
	}
}
