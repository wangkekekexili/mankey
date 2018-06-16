package parser

import (
	"strconv"

	"github.com/wangkekekexili/mankey/ast"
)

func (p *Parser) parseIdentifier() (ast.Expression, error) {
	return &ast.Identifier{Value: p.currentToken.Literal}, nil
}

func (p *Parser) parseIntegerLiteral() (ast.Expression, error) {
	v, err := strconv.ParseInt(p.currentToken.Literal, 10, 64)
	if err != nil {
		return nil, err
	}
	return &ast.IntegerLiteral{Value: v}, nil
}
