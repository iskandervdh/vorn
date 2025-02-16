package ast

import "github.com/iskandervdh/vorn/token"

type Node interface {
	TokenLiteral() string
	String() string
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string { return i.Value }
