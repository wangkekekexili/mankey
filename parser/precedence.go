package parser

type precedence int

const (
	Lowest precedence = iota + 1
	Equal
	LteGte
	Add
	Multi
	Prefix
	Call
)
