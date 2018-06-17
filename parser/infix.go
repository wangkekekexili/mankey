package parser

import "github.com/wangkekekexili/mankey/ast"

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
