package parser

import (
	"github.com/wangkekekexili/mankey/ast"
	"github.com/wangkekekexili/mankey/token"
)

func (p *Parser) parseInfixExpression(left ast.Expression) (ast.Expression, error) {
	infixExpression := &ast.InfixExpression{Left: left, Op: ast.Operator(p.currentToken.Literal)}
	d := p.currentPrecedence()
	p.nextToken()
	right, err := p.parseExpression(d)
	if err != nil {
		return nil, err
	}
	infixExpression.Right = right
	return infixExpression, nil
}

func (p *Parser) parseIndexExpression(left ast.Expression) (ast.Expression, error) {
	indexExpression := &ast.IndexExpression{Left: left}

	p.nextToken()
	index, err := p.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}
	indexExpression.Index = index

	if p.peekToken.Type != token.RBracket {
		return nil, errUnexpectedToken{t: p.peekToken, exp: "]"}
	}
	p.nextToken()
	return indexExpression, nil
}
