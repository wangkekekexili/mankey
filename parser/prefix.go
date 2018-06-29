package parser

import (
	"strconv"

	"github.com/wangkekekexili/mankey/ast"
	"github.com/wangkekekexili/mankey/token"
)

func (p *Parser) parseIdentifier() (ast.Expression, error) {
	ident := &ast.Identifier{Value: p.currentToken.Literal}
	if p.peekToken.Type != token.LParen {
		return ident, nil
	}
	p.nextToken()
	arguments, err := p.parseExpressionList(token.RParen)
	if err != nil {
		return nil, err
	}
	return &ast.CallExpression{
		Function:  ident,
		Arguments: arguments,
	}, nil
}

func (p *Parser) parseBoolean() (ast.Expression, error) {
	v, err := strconv.ParseBool(p.currentToken.Literal)
	if err != nil {
		return nil, err
	}
	return &ast.Boolean{Value: v}, nil
}

func (p *Parser) parseInteger() (ast.Expression, error) {
	v, err := strconv.ParseInt(p.currentToken.Literal, 10, 64)
	if err != nil {
		return nil, err
	}
	return &ast.Integer{Value: v}, nil
}

func (p *Parser) parseString() (ast.Expression, error) {
	return &ast.String{Value: p.currentToken.Literal}, nil
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

func (p *Parser) parseGroupedExpression() (ast.Expression, error) {
	p.nextToken()
	expr, err := p.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}
	if p.peekToken.Type != token.RParen {
		return nil, errUnexpectedToken{t: p.peekToken, exp: ")"}
	}
	p.nextToken()
	return expr, nil
}

func (p *Parser) parseIfExpression() (ast.Expression, error) {
	ifExpression := &ast.IfExpression{}

	p.nextToken()
	if p.currentToken.Type != token.LParen {
		return nil, errUnexpectedToken{t: p.currentToken, exp: "("}
	}
	p.nextToken()
	expr, err := p.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}
	ifExpression.Condition = expr
	p.nextToken()
	if p.currentToken.Type != token.RParen {
		return nil, errUnexpectedToken{t: p.currentToken, exp: ")"}
	}

	p.nextToken()
	if p.currentToken.Type != token.LBrace {
		return nil, errUnexpectedToken{t: p.currentToken, exp: "{"}
	}
	block, err := p.parseBlockStatement()
	if err != nil {
		return nil, err
	}
	ifExpression.Consequence = block

	if p.peekToken.Type == token.Else {
		p.nextToken()
		p.nextToken()
		if p.currentToken.Type != token.LBrace {
			return nil, errUnexpectedToken{t: p.currentToken, exp: "{"}
		}
		block, err := p.parseBlockStatement()
		if err != nil {
			return nil, err
		}
		ifExpression.Alternative = block
	}

	return ifExpression, nil
}

func (p *Parser) parseFunction() (ast.Expression, error) {
	function := &ast.Function{}

	p.nextToken()
	if p.currentToken.Type != token.LParen {
		return nil, errUnexpectedToken{t: p.currentToken, exp: "("}
	}
	list, err := p.parseParameterList()
	if err != nil {
		return nil, err
	}
	function.Parameters = list

	p.nextToken()
	if p.currentToken.Type != token.LBrace {
		return nil, errUnexpectedToken{t: p.currentToken, exp: "{"}
	}
	function.Body, err = p.parseBlockStatement()
	if err != nil {
		return nil, err
	}

	return function, nil
}

func (p *Parser) parseArray() (ast.Expression, error) {
	arr := &ast.Array{}

	list, err := p.parseExpressionList(token.RBracket)
	if err != nil {
		return nil, err
	}
	arr.Elements = list

	return arr, nil
}

func (p *Parser) parseHash() (ast.Expression, error) {
	hash := &ast.Hash{Value: make(map[ast.Expression]ast.Expression)}

	p.nextToken()
	if p.currentToken.Type == token.RBrace {
		return hash, nil
	}

	for p.currentToken.Type != token.EOF {
		key, err := p.parseExpression(Lowest)
		if err != nil {
			return nil, err
		}

		p.nextToken()
		if p.currentToken.Type != token.Colon {
			return nil, errUnexpectedToken{t: p.currentToken, exp: ":"}
		}
		p.nextToken()
		value, err := p.parseExpression(Lowest)
		if err != nil {
			return nil, err
		}

		hash.Value[key] = value

		p.nextToken()
		if p.currentToken.Type == token.RBrace {
			break
		}
		if p.currentToken.Type != token.Comma {
			return nil, errUnexpectedToken{t: p.currentToken, exp: ","}
		}
		p.nextToken()
	}
	return hash, nil
}
