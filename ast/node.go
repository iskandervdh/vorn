package ast

type Node interface {
	TokenLiteral() string
	String() string
	Line() int
	Column() int
}
