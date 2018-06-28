package parser

import "github.com/wangkekekexili/mankey/token"

type precedence int

const (
	Lowest precedence = iota + 1
	Equal
	LteGte
	Add
	Multi
	Prefix
	Call
	Index
)

var precedences = map[token.TokenType]precedence{
	token.Equal:    Equal,
	token.NotEqual: Equal,
	token.Lt:       LteGte,
	token.Lte:      LteGte,
	token.Gt:       LteGte,
	token.Gte:      LteGte,
	token.Add:      Add,
	token.Minus:    Add,
	token.Multiply: Multi,
	token.Divide:   Multi,
	token.LBracket: Index,
}

func (p *Parser) currentPrecedence() precedence {
	pre, ok := precedences[p.currentToken.Type]
	if ok {
		return pre
	}
	return Lowest
}

func (p *Parser) peekPrecedence() precedence {
	pre, ok := precedences[p.peekToken.Type]
	if ok {
		return pre
	}
	return Lowest
}
