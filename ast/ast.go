package ast

type Statement interface{}
type Expression interface{}

type Identifier struct {
	Value string
}

type Program struct {
	Statements []Statement
}

type VarStatement struct {
	Name  *Identifier
	Value Expression
}

type ReturnStatement struct {
	Value Expression
}
