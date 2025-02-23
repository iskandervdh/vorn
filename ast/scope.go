package ast

// Scope is an interface that represents a scope in the AST.
// A scope is a block of code that has its own set of variables.
// A scope can be nested inside another scope.
// All scopes are nested inside the program scope.
type Scope interface {
	GetParentScope() Scope
	GetScopeStatements() []Statement
}
