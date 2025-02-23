package ast

import (
	"bytes"
	"strings"

	"github.com/iskandervdh/vorn/token"
)

type Expression interface {
	Node
	expressionNode()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }
func (i *Identifier) Line() int            { return i.Token.Line }
func (i *Identifier) Column() int          { return i.Token.Column }

type PrefixExpression struct {
	Token    token.Token // The prefix token, e.g. !
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}
func (pe *PrefixExpression) Line() int   { return pe.Token.Line }
func (pe *PrefixExpression) Column() int { return pe.Token.Column }

type InfixExpression struct {
	Token    token.Token // The operator token, e.g. +
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}
func (ie *InfixExpression) Line() int   { return ie.Token.Line }
func (ie *InfixExpression) Column() int { return ie.Token.Column }

type IfExpression struct {
	Token     token.Token // The 'if' token
	Condition Expression

	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if (")
	out.WriteString(ie.Condition.String())
	out.WriteString(") ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}
func (ie *IfExpression) Line() int   { return ie.Token.Line }
func (ie *IfExpression) Column() int { return ie.Token.Column }

type BreakExpression struct {
	Token token.Token // The 'break' token
}

func (be *BreakExpression) expressionNode()      {}
func (be *BreakExpression) TokenLiteral() string { return be.Token.Literal }
func (be *BreakExpression) String() string       { return be.Token.Literal }
func (be *BreakExpression) Line() int            { return be.Token.Line }
func (be *BreakExpression) Column() int          { return be.Token.Column }

type ContinueExpression struct {
	Token token.Token // The 'continue' token
}

func (ce *ContinueExpression) expressionNode()      {}
func (ce *ContinueExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *ContinueExpression) String() string       { return ce.Token.Literal }
func (ce *ContinueExpression) Line() int            { return ce.Token.Line }
func (ce *ContinueExpression) Column() int          { return ce.Token.Column }

type CallExpression struct {
	Token     token.Token // The '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}

	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}
func (ce *CallExpression) Line() int   { return ce.Token.Line }
func (ce *CallExpression) Column() int { return ce.Token.Column }

type IndexExpression struct {
	Token token.Token // The [ token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}
func (ie *IndexExpression) Line() int   { return ie.Token.Line }
func (ie *IndexExpression) Column() int { return ie.Token.Column }

type ReassignmentExpression struct {
	Token token.Token // The = token
	Name  *Identifier
	Value Expression
}

type ChainingExpression struct {
	Token token.Token // The '.' token
	Left  Expression
	Right Expression
}

func (ce *ChainingExpression) expressionNode()      {}
func (ce *ChainingExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *ChainingExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ce.Left.String())
	out.WriteString(".")
	out.WriteString(ce.Right.String())
	out.WriteString(")")

	return out.String()
}
func (ce *ChainingExpression) Line() int   { return ce.Token.Line }
func (ce *ChainingExpression) Column() int { return ce.Token.Column }

func (rs *ReassignmentExpression) expressionNode()      {}
func (rs *ReassignmentExpression) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReassignmentExpression) String() string {
	var out bytes.Buffer

	out.WriteString(rs.Name.String())
	out.WriteString(" = ")
	out.WriteString(rs.Value.String())

	return out.String()
}
func (rs *ReassignmentExpression) Line() int   { return rs.Token.Line }
func (rs *ReassignmentExpression) Column() int { return rs.Token.Column }
