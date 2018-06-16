package ast

import (
	"fmt"
	"strconv"
	"strings"
)

type Node interface {
	String() string
}

type Statement interface {
	Node
}

type Expression interface {
	Node
}

type Identifier struct {
	Value string
}

func (i *Identifier) String() string {
	return i.Value
}

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	s := strings.Builder{}
	for _, stat := range p.Statements {
		s.WriteString(stat.String())
	}
	return s.String()
}

type VarStatement struct {
	Name  *Identifier
	Value Expression
}

func (s *VarStatement) String() string {
	return fmt.Sprintf("var %v = %v;", s.Name, s.Value)
}

type ReturnStatement struct {
	Value Expression
}

func (s *ReturnStatement) String() string {
	return fmt.Sprintf("return %v;", s.Value)
}

type ExpressionStatement struct {
	Value Expression
}

func (s *ExpressionStatement) String() string {
	return s.Value.String()
}

type IntegerLiteral struct {
	Value int64
}

func (s *IntegerLiteral) String() string {
	return strconv.FormatInt(s.Value, 10)
}
