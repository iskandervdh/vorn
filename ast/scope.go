package ast

type Scope interface {
	GetParentScope() Scope
	GetScopeStatements() []Statement
}
