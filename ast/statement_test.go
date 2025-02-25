package ast

import (
	"testing"

	"github.com/iskandervdh/vorn/token"
)

func TestExpressionStatement(t *testing.T) {
	statement := &ExpressionStatement{
		Token: token.Token{
			Type:    token.IDENT,
			Literal: "foo",
			Line:    1,
			Column:  1,
		},
		Expression: &Identifier{
			Token: token.Token{
				Type:    token.IDENT,
				Literal: "foo",
			},
			Value: "foo",
		},
	}

	if statement.String() != "foo" {
		t.Errorf("statement.String() wrong. got=%q", statement.String())
	}

	if statement.TokenLiteral() != "foo" {
		t.Errorf("statement.TokenLiteral() wrong. got=%q", statement.TokenLiteral())
	}

	if statement.Expression.String() != "foo" {
		t.Errorf("statement.Expression.String() wrong. got=%q", statement.Expression.String())
	}

	if statement.Expression.TokenLiteral() != "foo" {
		t.Errorf("statement.Expression.TokenLiteral() wrong. got=%q", statement.Expression.TokenLiteral())
	}

	if statement.Line() != 1 {
		t.Errorf("statement.Line() wrong. got=%d", statement.Line())
	}

	if statement.Column() != 1 {
		t.Errorf("statement.Column() wrong. got=%d", statement.Column())
	}

	statement.Expression = nil

	if statement.String() != "" {
		t.Errorf("statement.String() wrong. got=%q", statement.String())
	}

	statement.statementNode()
}

func TestVariableStatement(t *testing.T) {
	statement := &VariableStatement{
		Token: token.Token{
			Type:    token.LET,
			Literal: "let",
			Line:    1,
			Column:  1,
		},
		Name: &Identifier{
			Token: token.Token{
				Type:    token.IDENT,
				Literal: "foo",
			},
			Value: "foo",
		},
		Value: &Identifier{
			Token: token.Token{
				Type:    token.IDENT,
				Literal: "bar",
			},
			Value: "bar",
		},
	}

	if statement.String() != "let foo = bar;" {
		t.Errorf("statement.String() wrong. got=%q", statement.String())
	}

	if statement.TokenLiteral() != "let" {
		t.Errorf("statement.TokenLiteral() wrong. got=%q", statement.TokenLiteral())
	}

	if statement.Name.String() != "foo" {
		t.Errorf("statement.Name.String() wrong. got=%q", statement.Name.String())
	}

	if statement.Name.TokenLiteral() != "foo" {
		t.Errorf("statement.Name.TokenLiteral() wrong. got=%q", statement.Name.TokenLiteral())
	}

	if statement.Value.String() != "bar" {
		t.Errorf("statement.Value.String() wrong. got=%q", statement.Value.String())
	}

	if statement.Value.TokenLiteral() != "bar" {
		t.Errorf("statement.Value.TokenLiteral() wrong. got=%q", statement.Value.TokenLiteral())
	}

	if statement.Line() != 1 {
		t.Errorf("statement.Line() wrong. got=%d", statement.Line())
	}

	if statement.Column() != 1 {
		t.Errorf("statement.Column() wrong. got=%d", statement.Column())
	}

	statement.statementNode()
}

func TestReturnStatement(t *testing.T) {
	statement := &ReturnStatement{
		Token: token.Token{
			Type:    token.RETURN,
			Literal: "return",
			Line:    1,
			Column:  1,
		},
		ReturnValue: &Identifier{
			Token: token.Token{
				Type:    token.IDENT,
				Literal: "foo",
			},
			Value: "foo",
		},
	}

	if statement.String() != "return foo;" {
		t.Errorf("statement.String() wrong. got=%q", statement.String())
	}

	if statement.TokenLiteral() != "return" {
		t.Errorf("statement.TokenLiteral() wrong. got=%q", statement.TokenLiteral())
	}

	if statement.ReturnValue.String() != "foo" {
		t.Errorf("statement.ReturnValue.String() wrong. got=%q", statement.ReturnValue.String())
	}

	if statement.ReturnValue.TokenLiteral() != "foo" {
		t.Errorf("statement.ReturnValue.TokenLiteral() wrong. got=%q", statement.ReturnValue.TokenLiteral())
	}

	if statement.Line() != 1 {
		t.Errorf("statement.Line() wrong. got=%d", statement.Line())
	}

	if statement.Column() != 1 {
		t.Errorf("statement.Column() wrong. got=%d", statement.Column())
	}

	statement.statementNode()
}

func TestBlockStatement(t *testing.T) {
	program := NewProgram()

	statement := &BlockStatement{
		Token: token.Token{
			Type:    token.LBRACE,
			Literal: "{",
			Line:    1,
			Column:  1,
		},
		Statements: []Statement{
			&ExpressionStatement{
				Token: token.Token{
					Type:    token.IDENT,
					Literal: "foo",
				},
				Expression: &Identifier{
					Token: token.Token{
						Type:    token.IDENT,
						Literal: "foo",
					},
					Value: "foo",
				},
			},
		},
		Parent: program,
	}

	if statement.String() != "{\n  foo\n}" {
		t.Errorf("statement.String() wrong. got=%q", statement.String())
	}

	if statement.TokenLiteral() != "{" {
		t.Errorf("statement.TokenLiteral() wrong. got=%q", statement.TokenLiteral())
	}

	if statement.Line() != 1 {
		t.Errorf("statement.Line() wrong. got=%d", statement.Line())
	}

	if statement.Column() != 1 {
		t.Errorf("statement.Column() wrong. got=%d", statement.Column())
	}

	if statement.GetParentScope() != program {
		t.Errorf("statement.GetParentScope() wrong. got=%+v", statement.GetParentScope())
	}

	if len(statement.GetScopeStatements()) != len(statement.Statements) {
		t.Errorf("statement.GetScopeStatements() wrong. got=%+v", statement.GetScopeStatements())
	} else {
		for i, s := range statement.GetScopeStatements() {
			if s != statement.Statements[i] {
				t.Errorf("statement.GetScopeStatements() wrong. got=%+v", statement.GetScopeStatements())
				break
			}
		}
	}

	statement.statementNode()
}

func TestWhileStatement(t *testing.T) {
	program := NewProgram()

	statement := &WhileStatement{
		Token: token.Token{
			Type:    token.WHILE,
			Literal: "while",
			Line:    1,
			Column:  1,
		},
		Condition: &Identifier{
			Token: token.Token{
				Type:    token.IDENT,
				Literal: "foo",
			},
			Value: "foo",
		},
		Consequence: &BlockStatement{
			Token: token.Token{
				Type:    token.LBRACE,
				Literal: "{",
			},
			Statements: []Statement{
				&ExpressionStatement{
					Token: token.Token{
						Type:    token.IDENT,
						Literal: "foo",
					},
					Expression: &Identifier{
						Token: token.Token{
							Type:    token.IDENT,
							Literal: "foo",
						},
						Value: "foo",
					},
				},
			},
			Parent: program,
		},
	}

	if statement.String() != "while (foo) {\n  foo\n}" {
		t.Errorf("statement.String() wrong. got=%q", statement.String())
	}

	if statement.TokenLiteral() != "while" {
		t.Errorf("statement.TokenLiteral() wrong. got=%q", statement.TokenLiteral())
	}

	if statement.Condition.String() != "foo" {
		t.Errorf("statement.Condition.String() wrong. got=%q", statement.Condition.String())
	}

	if statement.Condition.TokenLiteral() != "foo" {
		t.Errorf("statement.Condition.TokenLiteral() wrong. got=%q", statement.Condition.TokenLiteral())
	}

	if statement.Line() != 1 {
		t.Errorf("statement.Line() wrong. got=%d", statement.Line())
	}

	if statement.Column() != 1 {
		t.Errorf("statement.Column() wrong. got=%d", statement.Column())
	}

	statement.statementNode()
}

func TestForStatement(t *testing.T) {
	program := NewProgram()

	statement := &ForStatement{
		Token: token.Token{
			Type:    token.FOR,
			Literal: "for",
			Line:    1,
			Column:  1,
		},
		Init: &ExpressionStatement{
			Token: token.Token{
				Type:    token.IDENT,
				Literal: "foo",
			},
			Expression: &Identifier{
				Token: token.Token{
					Type:    token.IDENT,
					Literal: "foo",
				},
				Value: "foo",
			},
		},
		Condition: &Identifier{
			Token: token.Token{
				Type:    token.IDENT,
				Literal: "foo",
			},
			Value: "foo",
		},
		Update: &ReassignmentExpression{
			Token: token.Token{
				Type:    token.IDENT,
				Literal: "foo",
			},
			Name: &Identifier{
				Token: token.Token{
					Type:    token.IDENT,
					Literal: "foo",
				},
				Value: "foo",
			},
			Value: &Identifier{
				Token: token.Token{
					Type:    token.IDENT,
					Literal: "foo",
				},
				Value: "foo + 1",
			},
		},
		Body: &BlockStatement{
			Token: token.Token{
				Type:    token.LBRACE,
				Literal: "{",
			},
			Statements: []Statement{
				&ExpressionStatement{
					Token: token.Token{
						Type:    token.IDENT,
						Literal: "foo",
					},
					Expression: &Identifier{
						Token: token.Token{
							Type:    token.IDENT,
							Literal: "foo",
						},
						Value: "foo",
					},
				},
			},
			Parent: program,
		},
		Parent: program,
	}

	if statement.String() != "for (foo; foo; foo = foo + 1) {\n  foo\n}" {
		t.Errorf("statement.String() wrong. got=%q", statement.String())
	}

	if statement.TokenLiteral() != "for" {
		t.Errorf("statement.TokenLiteral() wrong. got=%q", statement.TokenLiteral())
	}

	if statement.Init.String() != "foo" {
		t.Errorf("statement.Init.String() wrong. got=%q", statement.Init.String())
	}

	if statement.Init.TokenLiteral() != "foo" {
		t.Errorf("statement.Init.TokenLiteral() wrong. got=%q", statement.Init.TokenLiteral())
	}

	if statement.Condition.String() != "foo" {
		t.Errorf("statement.Condition.String() wrong. got=%q", statement.Condition.String())
	}

	if statement.Condition.TokenLiteral() != "foo" {
		t.Errorf("statement.Condition.TokenLiteral() wrong. got=%q", statement.Condition.TokenLiteral())
	}

	if statement.Update.String() != "foo = foo + 1" {
		t.Errorf("statement.Update.String() wrong. got=%q", statement.Update.String())
	}

	if statement.Update.TokenLiteral() != "foo" {
		t.Errorf("statement.Update.TokenLiteral() wrong. got=%q", statement.Update.TokenLiteral())
	}

	if statement.Line() != 1 {
		t.Errorf("statement.Line() wrong. got=%d", statement.Line())
	}

	if statement.Column() != 1 {
		t.Errorf("statement.Column() wrong. got=%d", statement.Column())
	}

	if statement.GetParentScope() != program {
		t.Errorf("statement.GetParentScope() wrong. got=%+v", statement.GetParentScope())
	}

	if len(statement.GetScopeStatements()) != len(statement.Statements) {
		t.Errorf("statement.GetScopeStatements() wrong. got=%+v", statement.GetScopeStatements())
	} else {
		for i, s := range statement.GetScopeStatements() {
			if s != statement.Body.Statements[i] {
				t.Errorf("statement.GetScopeStatements() wrong. got=%+v", statement.GetScopeStatements())
				break
			}
		}
	}

	statement.statementNode()
}

func TestFunctionStatement(t *testing.T) {
	program := NewProgram()

	statement := &FunctionStatement{
		Token: token.Token{
			Type:    token.FUNCTION,
			Literal: "function",
			Line:    1,
			Column:  1,
		},
		Name: &Identifier{
			Token: token.Token{
				Type:    token.IDENT,
				Literal: "foo",
			},
			Value: "foo",
		},
		Arguments: []*Identifier{
			{
				Token: token.Token{
					Type:    token.IDENT,
					Literal: "bar",
				},
				Value: "bar",
			},
		},
		Body: &BlockStatement{
			Token: token.Token{
				Type:    token.LBRACE,
				Literal: "{",
			},
			Statements: []Statement{
				&ExpressionStatement{
					Token: token.Token{
						Type:    token.IDENT,
						Literal: "foo",
					},
					Expression: &Identifier{
						Token: token.Token{
							Type:    token.IDENT,
							Literal: "foo",
						},
						Value: "foo",
					},
				},
			},
			Parent: program,
		},
	}

	if statement.String() != "function foo(bar) {\n  foo\n}" {
		t.Errorf("statement.String() wrong. got=%q", statement.String())
	}

	if statement.TokenLiteral() != "function" {
		t.Errorf("statement.TokenLiteral() wrong. got=%q", statement.TokenLiteral())
	}

	if statement.Name.String() != "foo" {
		t.Errorf("statement.Name.String() wrong. got=%q", statement.Name.String())
	}

	if statement.Name.TokenLiteral() != "foo" {
		t.Errorf("statement.Name.TokenLiteral() wrong. got=%q", statement.Name.TokenLiteral())
	}

	if statement.Line() != 1 {
		t.Errorf("statement.Line() wrong. got=%d", statement.Line())
	}

	if statement.Column() != 1 {
		t.Errorf("statement.Column() wrong. got=%d", statement.Column())
	}

	statement.statementNode()
}
