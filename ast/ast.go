package ast

type Statement interface{}
type Expression interface{}

type Program struct {
	Statements []Statement
}

type VarStatement struct {
	Name  *Identifier
	Value Expression
}

type Identifier struct {
	Value string
}
