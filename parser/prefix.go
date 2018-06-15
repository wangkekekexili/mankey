package parser

import "github.com/wangkekekexili/mankey/ast"

func (p *Parser) parseIdentifier() (ast.Expression, error) {
	return &ast.Identifier{
		Value: p.currentToken.Literal,
	}, nil
}
