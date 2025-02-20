package ast

import (
	"bytes"
	"strings"

	"github.com/iskandervdh/vorn/token"
)

type Statement interface {
	Node
	statementNode()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}

	return ""
}
func (es *ExpressionStatement) Line() int   { return es.Token.Line }
func (es *ExpressionStatement) Column() int { return es.Token.Column }

type VariableStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (vs *VariableStatement) statementNode()       {}
func (vs *VariableStatement) TokenLiteral() string { return vs.Token.Literal }
func (vs *VariableStatement) IsLet() bool          { return vs.Token.Type == token.LET }
func (vs *VariableStatement) IsConst() bool        { return vs.Token.Type == token.CONST }
func (vs *VariableStatement) String() string {
	var out bytes.Buffer

	out.WriteString(vs.TokenLiteral() + " ")
	out.WriteString(vs.Name.String())
	out.WriteString(" = ")

	if vs.Value != nil {
		out.WriteString(vs.Value.String())
	}

	out.WriteString(";")

	return out.String()
}
func (vs *VariableStatement) Line() int   { return vs.Token.Line }
func (vs *VariableStatement) Column() int { return vs.Token.Column }

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}
func (rs *ReturnStatement) Line() int   { return rs.Token.Line }
func (rs *ReturnStatement) Column() int { return rs.Token.Column }

type BlockStatement struct {
	Parent     Scope
	Token      token.Token // the { token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}
func (bs *BlockStatement) Line() int   { return bs.Token.Line }
func (bs *BlockStatement) Column() int { return bs.Token.Column }

func (bs *BlockStatement) GetParentScope() Scope {
	return bs.Parent
}

func (bs *BlockStatement) GetScopeStatements() []Statement {
	return bs.Statements
}

type WhileStatement struct {
	Token     token.Token // The 'while' token
	Condition Expression

	Consequence *BlockStatement
}

func (ws *WhileStatement) statementNode()       {}
func (ws *WhileStatement) TokenLiteral() string { return ws.Token.Literal }
func (ws *WhileStatement) String() string {
	var out bytes.Buffer

	out.WriteString("while")
	out.WriteString(ws.Condition.String())
	out.WriteString(" ")
	out.WriteString(ws.Consequence.String())

	return out.String()
}
func (ws *WhileStatement) Line() int   { return ws.Token.Line }
func (ws *WhileStatement) Column() int { return ws.Token.Column }

type ForStatement struct {
	Parent     Scope
	Token      token.Token
	Statements []Statement

	Init      Statement
	Condition Expression
	Update    Expression
	Body      *BlockStatement
}

func (fs *ForStatement) statementNode()       {}
func (fs *ForStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *ForStatement) String() string {
	var out bytes.Buffer

	out.WriteString("for")
	out.WriteString(fs.Init.String())
	out.WriteString("; ")
	out.WriteString(fs.Condition.String())
	out.WriteString("; ")
	out.WriteString(fs.Update.String())
	out.WriteString(" ")
	out.WriteString(fs.Body.String())

	return out.String()
}
func (fs *ForStatement) Line() int   { return fs.Token.Line }
func (fs *ForStatement) Column() int { return fs.Token.Column }

func (fs *ForStatement) GetParentScope() Scope {
	return fs.Parent
}

func (fs *ForStatement) GetScopeStatements() []Statement {
	return fs.Statements
}

type FunctionStatement struct {
	Token      token.Token // the 'func' token
	Name       *Identifier
	Parameters []*Identifier

	Body *BlockStatement
}

func (fs *FunctionStatement) statementNode()       {}
func (fs *FunctionStatement) TokenLiteral() string { return fs.Token.Literal }
func (fs *FunctionStatement) String() string {
	var out bytes.Buffer

	params := []string{}

	for _, p := range fs.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(fs.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(fs.Name.String())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")

	out.WriteString(fs.Body.String())

	return out.String()
}
func (fs *FunctionStatement) Line() int   { return fs.Token.Line }
func (fs *FunctionStatement) Column() int { return fs.Token.Column }
