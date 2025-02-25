package ast

import (
	"testing"

	"github.com/iskandervdh/vorn/token"
)

func TestIntegerLiteral(t *testing.T) {
	literal := &IntegerLiteral{
		Token: token.Token{
			Type:    token.INT,
			Literal: "5",
			Line:    1,
			Column:  1,
		},
		Value: 5,
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral() wrong. got=%s", literal.TokenLiteral())
	}

	if literal.String() != "5" {
		t.Errorf("literal.String() wrong. got=%s", literal.String())
	}

	if literal.Line() != 1 {
		t.Errorf("literal.Line() wrong. got=%d", literal.Line())
	}

	if literal.Column() != 1 {
		t.Errorf("literal.Column() wrong. got=%d", literal.Column())
	}

	literal.expressionNode()
}

func TestBooleanLiteral(t *testing.T) {
	literal := &BooleanLiteral{
		Token: token.Token{
			Type:    token.TRUE,
			Literal: "true",
			Line:    1,
			Column:  1,
		},
		Value: true,
	}

	if literal.TokenLiteral() != "true" {
		t.Errorf("literal.TokenLiteral() wrong. got=%s", literal.TokenLiteral())
	}

	if literal.String() != "true" {
		t.Errorf("literal.String() wrong. got=%s", literal.String())
	}

	if literal.Line() != 1 {
		t.Errorf("literal.Line() wrong. got=%d", literal.Line())
	}

	if literal.Column() != 1 {
		t.Errorf("literal.Column() wrong. got=%d", literal.Column())
	}

	literal.expressionNode()
}

func TestFloatLiteral(t *testing.T) {
	literal := &FloatLiteral{
		Token: token.Token{
			Type:    token.FLOAT,
			Literal: "5.5",
			Line:    1,
			Column:  1,
		},
		Value: 5.5,
	}

	if literal.TokenLiteral() != "5.5" {
		t.Errorf("literal.TokenLiteral() wrong. got=%s", literal.TokenLiteral())
	}

	if literal.String() != "5.5" {
		t.Errorf("literal.String() wrong. got=%s", literal.String())
	}

	if literal.Line() != 1 {
		t.Errorf("literal.Line() wrong. got=%d", literal.Line())
	}

	if literal.Column() != 1 {
		t.Errorf("literal.Column() wrong. got=%d", literal.Column())
	}

	literal.expressionNode()
}

func TestNullLiteral(t *testing.T) {
	literal := &NullLiteral{
		Token: token.Token{
			Type:    token.NULL,
			Literal: "null",
			Line:    1,
			Column:  1,
		},
	}

	if literal.TokenLiteral() != "null" {
		t.Errorf("literal.TokenLiteral() wrong. got=%s", literal.TokenLiteral())
	}

	if literal.String() != "null" {
		t.Errorf("literal.String() wrong. got=%s", literal.String())
	}

	if literal.Line() != 1 {
		t.Errorf("literal.Line() wrong. got=%d", literal.Line())
	}

	if literal.Column() != 1 {
		t.Errorf("literal.Column() wrong. got=%d", literal.Column())
	}

	literal.expressionNode()
}

func TestFunctionLiteral(t *testing.T) {
	literal := &FunctionLiteral{
		Token: token.Token{
			Type:    token.FUNCTION,
			Literal: "func",
			Line:    1,
			Column:  1,
		},
		Arguments: []*Identifier{
			{
				Token: token.Token{
					Type:    token.IDENT,
					Literal: "a",
					Line:    1,
					Column:  6,
				},
				Value: "a",
			},
		},
		Body: &BlockStatement{
			Token: token.Token{
				Type:    token.LBRACE,
				Literal: "{",
				Line:    1,
				Column:  8,
			},
			Statements: []Statement{},
		},
	}

	if literal.TokenLiteral() != "func" {
		t.Errorf("literal.TokenLiteral() wrong. got=%s", literal.TokenLiteral())
	}

	if literal.String() != `func(a) {
}` {
		t.Errorf("literal.String() wrong. got=%s", literal.String())
	}

	if literal.Line() != 1 {
		t.Errorf("literal.Line() wrong. got=%d", literal.Line())
	}

	if literal.Column() != 1 {
		t.Errorf("literal.Column() wrong. got=%d", literal.Column())
	}

	literal.expressionNode()
}

func TestStringLiteral(t *testing.T) {
	literal := &StringLiteral{
		Token: token.Token{
			Type:    token.STRING,
			Literal: "hello",
			Line:    1,
			Column:  1,
		},
		Value: "hello",
	}

	if literal.TokenLiteral() != "hello" {
		t.Errorf("literal.TokenLiteral() wrong. got=%s", literal.TokenLiteral())
	}

	if literal.String() != "hello" {
		t.Errorf("literal.String() wrong. got=%s", literal.String())
	}

	if literal.Line() != 1 {
		t.Errorf("literal.Line() wrong. got=%d", literal.Line())
	}

	if literal.Column() != 1 {
		t.Errorf("literal.Column() wrong. got=%d", literal.Column())
	}

	literal.expressionNode()
}

func TestArrayLiteral(t *testing.T) {
	literal := &ArrayLiteral{
		Token: token.Token{
			Type:    token.LBRACKET,
			Literal: "[",
			Line:    1,
			Column:  1,
		},
		Elements: []Expression{
			&IntegerLiteral{
				Token: token.Token{
					Type:    token.INT,
					Literal: "5",
					Line:    1,
					Column:  2,
				},
				Value: 5,
			},
		},
	}

	if literal.TokenLiteral() != "[" {
		t.Errorf("literal.TokenLiteral() wrong. got=%s", literal.TokenLiteral())
	}

	if literal.String() != "[5]" {
		t.Errorf("literal.String() wrong. got=%s", literal.String())
	}

	if literal.Line() != 1 {
		t.Errorf("literal.Line() wrong. got=%d", literal.Line())
	}

	if literal.Column() != 1 {
		t.Errorf("literal.Column() wrong. got=%d", literal.Column())
	}

	literal.expressionNode()
}

func TestHashLiteral(t *testing.T) {
	literal := &HashLiteral{
		Token: token.Token{
			Type:    token.LBRACE,
			Literal: "{",
			Line:    1,
			Column:  1,
		},
		Pairs: map[Expression]Expression{
			&StringLiteral{
				Token: token.Token{
					Type:    token.STRING,
					Literal: "a",
					Line:    1,
					Column:  2,
				},
				Value: "a",
			}: &IntegerLiteral{
				Token: token.Token{
					Type:    token.INT,
					Literal: "5",
					Line:    1,
					Column:  5,
				},
				Value: 5,
			},
		},
	}

	if literal.TokenLiteral() != "{" {
		t.Errorf("literal.TokenLiteral() wrong. got=%s", literal.TokenLiteral())
	}

	if literal.String() != `{a: 5}` {
		t.Errorf("literal.String() wrong. got=%s", literal.String())
	}

	if literal.Line() != 1 {
		t.Errorf("literal.Line() wrong. got=%d", literal.Line())
	}

	if literal.Column() != 1 {
		t.Errorf("literal.Column() wrong. got=%d", literal.Column())
	}

	literal.expressionNode()
}
