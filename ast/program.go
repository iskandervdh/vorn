package ast

import "bytes"

type Program struct {
	Statements []Statement
}

func NewProgram() *Program {
	return &Program{
		Statements: []Statement{},
	}
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

func (p *Program) Line() int {
	if len(p.Statements) > 0 {
		return p.Statements[0].Line()
	} else {
		return 0
	}
}

func (p *Program) Column() int {
	if len(p.Statements) > 0 {
		return p.Statements[0].Column()
	} else {
		return 0
	}
}

func (p *Program) GetParentScope() Scope {
	return nil
}

func (p *Program) GetScopeStatements() []Statement {
	return p.Statements
}
