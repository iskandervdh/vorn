package ast

import (
	"testing"

	"github.com/iskandervdh/vorn/token"
)

func TestIdentifier(t *testing.T) {
	expression := &Identifier{
		Value: "foo",
		Token: token.Token{
			Type:    token.IDENT,
			Literal: "foo",
			Line:    1,
			Column:  1,
		},
	}

	if expression.String() != "foo" {
		t.Errorf("Identifier.String() = %s; want foo", expression.String())
	}

	if expression.TokenLiteral() != "foo" {
		t.Errorf("Identifier.TokenLiteral() = %s; want foo", expression.TokenLiteral())
	}

	if expression.Line() != 1 {
		t.Errorf("Identifier.Line() = %d; want 1", expression.Line())
	}

	if expression.Column() != 1 {
		t.Errorf("Identifier.Column() = %d; want 1", expression.Column())
	}

	expression.expressionNode()
}

func TestPrefixExpression(t *testing.T) {
	expression := &PrefixExpression{
		Operator: "!",
		Right: &Identifier{
			Value: "foo",
			Token: token.Token{
				Type:    token.IDENT,
				Literal: "foo",
				Line:    1,
				Column:  1,
			},
		},
		Token: token.Token{
			Type:    token.EXCLAMATION,
			Literal: "!",
			Line:    1,
			Column:  1,
		},
	}

	if expression.String() != "(!foo)" {
		t.Errorf("PrefixExpression.String() = %s; want (!foo)", expression.String())
	}

	if expression.TokenLiteral() != "!" {
		t.Errorf("PrefixExpression.TokenLiteral() = %s; want !", expression.TokenLiteral())
	}

	if expression.Line() != 1 {
		t.Errorf("PrefixExpression.Line() = %d; want 1", expression.Line())
	}

	if expression.Column() != 1 {
		t.Errorf("PrefixExpression.Column() = %d; want 1", expression.Column())
	}

	expression.expressionNode()
}

func TestInfixExpression(t *testing.T) {
	expression := &InfixExpression{
		Left: &Identifier{
			Value: "foo",
			Token: token.Token{
				Type:    token.IDENT,
				Literal: "foo",
				Line:    1,
				Column:  1,
			},
		},
		Operator: "+",
		Right: &Identifier{
			Value: "bar",
			Token: token.Token{
				Type:    token.IDENT,
				Literal: "bar",
				Line:    1,
				Column:  7,
			},
		},
		Token: token.Token{
			Type:    token.PLUS,
			Literal: "+",
			Line:    1,
			Column:  5,
		},
	}

	if expression.String() != "(foo + bar)" {
		t.Errorf("InfixExpression.String() = %s; want (foo + bar)", expression.String())
	}

	if expression.TokenLiteral() != "+" {
		t.Errorf("InfixExpression.TokenLiteral() = %s; want +", expression.TokenLiteral())
	}

	if expression.Line() != 1 {
		t.Errorf("InfixExpression.Line() = %d; want 1", expression.Line())
	}

	if expression.Column() != 5 {
		t.Errorf("InfixExpression.Column() = %d; want 1", expression.Column())
	}

	expression.expressionNode()
}

func TestIfElseExpression(t *testing.T) {
	program := NewProgram()

	expression := &IfExpression{
		Condition: &Identifier{
			Value: "foo",
			Token: token.Token{
				Type:    token.IDENT,
				Literal: "foo",
				Line:    1,
				Column:  1,
			},
		},
		Consequence: &BlockStatement{
			Statements: []Statement{
				&ExpressionStatement{
					Expression: &Identifier{
						Value: "bar",
						Token: token.Token{
							Type:    token.IDENT,
							Literal: "bar",
							Line:    1,
							Column:  1,
						},
					},
				},
			},
			Parent: program,
			Token: token.Token{
				Type:    token.LBRACE,
				Literal: "{",
				Line:    1,
				Column:  1,
			},
		},
		Alternative: &BlockStatement{
			Statements: []Statement{
				&ExpressionStatement{
					Expression: &Identifier{
						Value: "baz",
						Token: token.Token{
							Type:    token.IDENT,
							Literal: "baz",
							Line:    1,
							Column:  1,
						},
					},
				},
			},
			Parent: program,
			Token: token.Token{
				Type:    token.LBRACE,
				Literal: "{",
				Line:    1,
				Column:  1,
			},
		},
		Token: token.Token{
			Type:    token.IF,
			Literal: "if",
			Line:    1,
			Column:  1,
		},
	}

	expected := `if (foo) {
  bar
} else {
  baz
}`

	if expression.String() != expected {
		t.Errorf("IfElseExpression.String() = %s, want %s", expression.String(), expected)
	}

	if expression.TokenLiteral() != "if" {
		t.Errorf("IfElseExpression.TokenLiteral() = %s; want if", expression.TokenLiteral())
	}

	if expression.Line() != 1 {
		t.Errorf("IfElseExpression.Line() = %d; want 1", expression.Line())
	}

	if expression.Column() != 1 {
		t.Errorf("IfElseExpression.Column() = %d; want 1", expression.Column())
	}

	expression.expressionNode()
}

func TestBreakExpression(t *testing.T) {
	expression := &BreakExpression{
		Token: token.Token{
			Type:    token.BREAK,
			Literal: "break",
			Line:    1,
			Column:  1,
		},
	}

	if expression.String() != "break" {
		t.Errorf("BreakExpression.String() = %s; want break", expression.String())
	}

	if expression.TokenLiteral() != "break" {
		t.Errorf("BreakExpression.TokenLiteral() = %s; want break", expression.TokenLiteral())
	}

	if expression.Line() != 1 {
		t.Errorf("BreakExpression.Line() = %d; want 1", expression.Line())
	}

	if expression.Column() != 1 {
		t.Errorf("BreakExpression.Column() = %d; want 1", expression.Column())
	}

	expression.expressionNode()
}

func TestContinueExpression(t *testing.T) {
	expression := &ContinueExpression{
		Token: token.Token{
			Type:    token.CONTINUE,
			Literal: "continue",
			Line:    1,
			Column:  1,
		},
	}

	if expression.String() != "continue" {
		t.Errorf("ContinueExpression.String() = %s; want continue", expression.String())
	}

	if expression.TokenLiteral() != "continue" {
		t.Errorf("ContinueExpression.TokenLiteral() = %s; want continue", expression.TokenLiteral())
	}

	if expression.Line() != 1 {
		t.Errorf("ContinueExpression.Line() = %d; want 1", expression.Line())
	}

	if expression.Column() != 1 {
		t.Errorf("ContinueExpression.Column() = %d; want 1", expression.Column())
	}

	expression.expressionNode()
}

func TestCallExpression(t *testing.T) {
	expression := &CallExpression{
		Function: &Identifier{
			Value: "foo",
			Token: token.Token{
				Type:    token.IDENT,
				Literal: "foo",
				Line:    1,
				Column:  1,
			},
		},
		Arguments: []Expression{
			&Identifier{
				Value: "bar",
				Token: token.Token{
					Type:    token.IDENT,
					Literal: "bar",
					Line:    1,
					Column:  1,
				},
			},
		},
		Token: token.Token{
			Type:    token.LPAREN,
			Literal: "(",
			Line:    1,
			Column:  1,
		},
	}

	if expression.String() != "foo(bar)" {
		t.Errorf("CallExpression.String() = %s; want foo(bar)", expression.String())
	}

	if expression.TokenLiteral() != "(" {
		t.Errorf("CallExpression.TokenLiteral() = %s; want (", expression.TokenLiteral())
	}

	if expression.Line() != 1 {
		t.Errorf("CallExpression.Line() = %d; want 1", expression.Line())
	}

	if expression.Column() != 1 {
		t.Errorf("CallExpression.Column() = %d; want 1", expression.Column())
	}

	expression.expressionNode()
}

func TestIndexExpression(t *testing.T) {
	expression := &IndexExpression{
		Left: &Identifier{
			Value: "foo",
			Token: token.Token{
				Type:    token.IDENT,
				Literal: "foo",
				Line:    1,
				Column:  1,
			},
		},
		Index: &Identifier{
			Value: "bar",
			Token: token.Token{
				Type:    token.IDENT,
				Literal: "bar",
				Line:    1,
				Column:  1,
			},
		},
		Token: token.Token{
			Type:    token.LBRACKET,
			Literal: "[",
			Line:    1,
			Column:  1,
		},
	}

	if expression.String() != "(foo[bar])" {
		t.Errorf("IndexExpression.String() = %s; want (foo[bar])", expression.String())
	}

	if expression.TokenLiteral() != "[" {
		t.Errorf("IndexExpression.TokenLiteral() = %s; want [", expression.TokenLiteral())
	}

	if expression.Line() != 1 {
		t.Errorf("IndexExpression.Line() = %d; want 1", expression.Line())
	}

	if expression.Column() != 1 {
		t.Errorf("IndexExpression.Column() = %d; want 1", expression.Column())
	}

	expression.expressionNode()
}

func TestChainingExpression(t *testing.T) {
	expression := &ChainingExpression{
		Left: &Identifier{
			Value: "foo",
			Token: token.Token{
				Type:    token.IDENT,
				Literal: "foo",
				Line:    1,
				Column:  1,
			},
		},
		Right: &Identifier{
			Value: "bar",
			Token: token.Token{
				Type:    token.IDENT,
				Literal: "bar",
				Line:    1,
				Column:  1,
			},
		},
		Token: token.Token{
			Type:    token.DOT,
			Literal: ".",
			Line:    1,
			Column:  1,
		},
	}

	if expression.String() != "(foo.bar)" {
		t.Errorf("ChainingExpression.String() = %s; want (foo.bar)", expression.String())
	}

	if expression.TokenLiteral() != "." {
		t.Errorf("ChainingExpression.TokenLiteral() = %s; want .", expression.TokenLiteral())
	}

	if expression.Line() != 1 {
		t.Errorf("ChainingExpression.Line() = %d; want 1", expression.Line())
	}

	if expression.Column() != 1 {
		t.Errorf("ChainingExpression.Column() = %d; want 1", expression.Column())
	}

	expression.expressionNode()
}

func TestReassignmentExpression(t *testing.T) {
	expression := &ReassignmentExpression{
		Name: &Identifier{
			Value: "foo",
			Token: token.Token{
				Type:    token.IDENT,
				Literal: "foo",
				Line:    1,
				Column:  1,
			},
		},
		Value: &Identifier{
			Value: "bar",
			Token: token.Token{
				Type:    token.IDENT,
				Literal: "bar",
				Line:    1,
				Column:  1,
			},
		},
		Token: token.Token{
			Type:    token.ASSIGN,
			Literal: "=",
			Line:    1,
			Column:  1,
		},
	}

	if expression.String() != "foo = bar" {
		t.Errorf("ReassignmentExpression.String() = %s; want foo = bar", expression.String())
	}

	if expression.TokenLiteral() != "=" {
		t.Errorf("ReassignmentExpression.TokenLiteral() = %s; want =", expression.TokenLiteral())
	}

	if expression.Line() != 1 {
		t.Errorf("ReassignmentExpression.Line() = %d; want 1", expression.Line())
	}

	if expression.Column() != 1 {
		t.Errorf("ReassignmentExpression.Column() = %d; want 1", expression.Column())
	}

	expression.expressionNode()
}

func TestIncrementDecrementExpression(t *testing.T) {
	expression := &IncrementDecrementExpression{
		Identifier: &Identifier{
			Value: "foo",
			Token: token.Token{
				Type:    token.IDENT,
				Literal: "foo",
				Line:    1,
				Column:  1,
			},
		},
		Token: token.Token{
			Type:    token.INCREMENT,
			Literal: "++",
			Line:    1,
			Column:  1,
		},
	}

	if expression.String() != "foo++" {
		t.Errorf("IncrementDecrementExpression.String() = %s; want foo++", expression.String())
	}

	if expression.TokenLiteral() != "++" {
		t.Errorf("IncrementDecrementExpression.TokenLiteral() = %s; want ++", expression.TokenLiteral())
	}

	if expression.Line() != 1 {
		t.Errorf("IncrementDecrementExpression.Line() = %d; want 1", expression.Line())
	}

	if expression.Column() != 1 {
		t.Errorf("IncrementDecrementExpression.Column() = %d; want 1", expression.Column())
	}

	expression.expressionNode()
}

func TestIncrementDecrementExpressionAfter(t *testing.T) {
	expression := &IncrementDecrementExpression{
		Identifier: &Identifier{
			Value: "foo",
			Token: token.Token{
				Type:    token.IDENT,
				Literal: "foo",
				Line:    1,
				Column:  1,
			},
		},
		Token: token.Token{
			Type:    token.DECREMENT,
			Literal: "--",
			Line:    1,
			Column:  1,
		},
		Before: true,
	}

	if expression.String() != "--foo" {
		t.Errorf("IncrementDecrementExpression.String() = %s; want --foo", expression.String())
	}

	if expression.TokenLiteral() != "--" {
		t.Errorf("IncrementDecrementExpression.TokenLiteral() = %s; want --", expression.TokenLiteral())
	}

	if expression.Line() != 1 {
		t.Errorf("IncrementDecrementExpression.Line() = %d; want 1", expression.Line())
	}

	if expression.Column() != 1 {
		t.Errorf("IncrementDecrementExpression.Column() = %d; want 1", expression.Column())
	}

	expression.expressionNode()
}
