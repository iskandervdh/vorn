package ast

import (
	"testing"

	"github.com/iskandervdh/vorn/token"
)

func TestExpressionStatement(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&ExpressionStatement{
				Token: token.Token{Type: token.IDENT, Literal: "myVar", Line: 1, Column: 1},
				Expression: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
			},
			&ExpressionStatement{
				Token:      token.Token{Type: token.IDENT, Literal: "anotherVar"},
				Expression: nil,
			},
		},
	}

	if program.String() != "myVar" {
		t.Errorf("program.String() wrong. got '%q'", program.String())
	}

	if program.Statements[0].String() != "myVar" {
		t.Errorf("program.Statements[0].String() wrong. got '%q'", program.Statements[0].String())
	}

	if program.Statements[0].TokenLiteral() != "myVar" {
		t.Errorf("program.Statements[0].TokenLiteral() wrong. got '%q'", program.Statements[0].TokenLiteral())
	}

	if program.Statements[1].String() != "" {
		t.Errorf("program.Statements[1].String() wrong. got '%q'", program.Statements[1].String())
	}

	if program.Statements[0].Line() != 1 {
		t.Errorf("program.Statements[0].Line() wrong. got '%q'", program.Statements[0].Line())
	}

	if program.Statements[0].Column() != 1 {
		t.Errorf("program.Statements[0].Column() wrong. got '%q'", program.Statements[0].Column())
	}

	program.Statements[1].statementNode()
}

func TestVariableStatement(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&VariableStatement{
				Token: token.Token{Type: token.LET, Literal: "let", Line: 2, Column: 10},
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

	if program.Statements[0].String() != "let myVar = anotherVar;" {
		t.Errorf("program.Statements[0].String() wrong. got '%q'", program.Statements[0].String())
	}

	if program.Statements[0].TokenLiteral() != "let" {
		t.Errorf("program.Statements[0].TokenLiteral() wrong. got '%q'", program.Statements[0].TokenLiteral())
	}

	if program.Statements[0].Line() != 2 {
		t.Errorf("program.Statements[0].Line() wrong. got '%q'", program.Statements[0].Line())
	}

	if program.Statements[0].Column() != 10 {
		t.Errorf("program.Statements[0].Column() wrong. got '%q'", program.Statements[0].Column())
	}

	program.Statements[0].statementNode()
}
