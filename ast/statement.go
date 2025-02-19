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

func (bs *BlockStatement) GetParentScope() Scope {
	return bs.Parent
}

func (bs *BlockStatement) GetScopeStatements() []Statement {
	return bs.Statements
}

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
	out.WriteString(fs.Condition.String())
	out.WriteString(" ")
	out.WriteString(fs.Body.String())

	return out.String()
}

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
