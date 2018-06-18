package parser

import (
	"strconv"

	"github.com/wangkekekexili/mankey/ast"
)

func (p *Parser) parseIdentifier() (ast.Expression, error) {
	return &ast.Identifier{Value: p.currentToken.Literal}, nil
}

func (p *Parser) parseBoolean() (ast.Expression, error) {
	v, err := strconv.ParseBool(p.currentToken.Literal)
	if err != nil {
		return nil, err
	}
	return &ast.Boolean{Value: v}, nil
}

func (p *Parser) parseIntegerLiteral() (ast.Expression, error) {
	v, err := strconv.ParseInt(p.currentToken.Literal, 10, 64)
	if err != nil {
		return nil, err
	}
	return &ast.IntegerLiteral{Value: v}, nil
}

func (p *Parser) parsePrefixExpression() (ast.Expression, error) {
	prefixExpression := &ast.PrefixExpression{Op: ast.Operator(p.currentToken.Literal)}
	p.nextToken()
	expr, err := p.parseExpression(Prefix)
	if err != nil {
		return nil, err
	}
	prefixExpression.Value = expr
	return prefixExpression, nil
}
