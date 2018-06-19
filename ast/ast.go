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

type BlockStatement struct {
	Statements []Statement
}

func (b *BlockStatement) String() string {
	s := strings.Builder{}
	s.WriteByte('{')
	for _, stat := range b.Statements {
		s.WriteString(stat.String())
	}
	s.WriteByte('}')
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

type Boolean struct {
	Value bool
}

func (b *Boolean) String() string {
	return strconv.FormatBool(b.Value)
}

type Integer struct {
	Value int64
}

func (s *Integer) String() string {
	return strconv.FormatInt(s.Value, 10)
}

type Operator string

type PrefixExpression struct {
	Op    Operator
	Value Expression
}

func (p *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%v)", p.Op, p.Value)
}

type InfixExpression struct {
	Left  Expression
	Op    Operator
	Right Expression
}

func (i *InfixExpression) String() string {
	return fmt.Sprintf("(%v%s%v)", i.Left, i.Op, i.Right)
}

type IfExpression struct {
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (i *IfExpression) String() string {
	s := fmt.Sprintf("if (%v) %v", i.Condition, i.Consequence)
	if i.Alternative != nil {
		s += fmt.Sprintf(" else %v", i.Alternative)
	}
	return s
}

type Function struct {
	Parameters []*Identifier
	Body       *BlockStatement
}

func (f *Function) String() string {
	paramStrs := make([]string, 0, len(f.Parameters))
	for _, para := range f.Parameters {
		paramStrs = append(paramStrs, para.String())
	}
	return fmt.Sprintf("func (%v) %v", strings.Join(paramStrs, ", "), f.Body)
}
