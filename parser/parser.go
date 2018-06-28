package parser

import (
	"github.com/wangkekekexili/mankey/ast"
	"github.com/wangkekekexili/mankey/lexer"
	"github.com/wangkekekexili/mankey/token"
)

type Parser struct {
	r *lexer.Lexer

	currentToken *token.Token
	peekToken    *token.Token

	prefixParseFnMap map[token.TokenType]prefixParseFn
	infixParseFnMap  map[token.TokenType]infixParseFn
}

func New(r *lexer.Lexer) *Parser {
	p := &Parser{
		r: r,
	}
	p.prefixParseFnMap = map[token.TokenType]prefixParseFn{
		token.Ident:    p.parseIdentifier,
		token.Number:   p.parseInteger,
		token.String:   p.parseString,
		token.Minus:    p.parsePrefixExpression,
		token.Not:      p.parsePrefixExpression,
		token.True:     p.parseBoolean,
		token.False:    p.parseBoolean,
		token.LParen:   p.parseGroupedExpression,
		token.LBracket: p.parseArray,
		token.If:       p.parseIfExpression,
		token.Func:     p.parseFunction,
	}
	p.infixParseFnMap = map[token.TokenType]infixParseFn{
		token.Equal:    p.parseInfixExpression,
		token.NotEqual: p.parseInfixExpression,
		token.Lt:       p.parseInfixExpression,
		token.Lte:      p.parseInfixExpression,
		token.Gt:       p.parseInfixExpression,
		token.Gte:      p.parseInfixExpression,
		token.Add:      p.parseInfixExpression,
		token.Minus:    p.parseInfixExpression,
		token.Multiply: p.parseInfixExpression,
		token.Divide:   p.parseInfixExpression,
		token.LBracket: p.parseIndexExpression,
	}

	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.r.NextToken()
}

func (p *Parser) ParseProgram() (*ast.Program, error) {
	program := &ast.Program{}
	for p.currentToken.Type != token.EOF {
		stat, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		program.Statements = append(program.Statements, stat)
		p.nextToken()
	}
	return program, nil
}

func (p *Parser) parseBlockStatement() (*ast.BlockStatement, error) {
	block := &ast.BlockStatement{}

	p.nextToken()
	for p.currentToken.Type != token.RBrace && p.currentToken.Type != token.EOF {
		stat, err := p.parseStatement()
		if err != nil {
			return nil, err
		}
		block.Statements = append(block.Statements, stat)
		p.nextToken()
	}
	if p.currentToken.Type != token.RBrace {
		return nil, errUnexpectedToken{t: p.currentToken, exp: "}"}
	}
	return block, nil
}

func (p *Parser) parseStatement() (ast.Statement, error) {
	switch p.currentToken.Type {
	case token.Var:
		return p.parseVarStatement()
	case token.Return:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseVarStatement() (*ast.VarStatement, error) {
	varStat := &ast.VarStatement{}

	p.nextToken()
	if p.currentToken.Type != token.Ident {
		return nil, errUnexpectedToken{exp: "var statement", t: p.currentToken}

	}
	varStat.Name = &ast.Identifier{Value: p.currentToken.Literal}

	p.nextToken()
	if p.currentToken.Type != token.Assign {
		return nil, errUnexpectedToken{exp: "=", t: p.currentToken}
	}

	p.nextToken()
	expr, err := p.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}
	varStat.Value = expr

	if p.peekToken.Type == token.Semicolon {
		p.nextToken()
	}

	return varStat, nil
}

func (p *Parser) parseReturnStatement() (*ast.ReturnStatement, error) {
	returnStatement := &ast.ReturnStatement{}

	p.nextToken()
	expr, err := p.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}
	returnStatement.Value = expr

	if p.peekToken.Type == token.Semicolon {
		p.nextToken()
	}

	return returnStatement, nil
}

func (p *Parser) parseExpressionStatement() (*ast.ExpressionStatement, error) {
	expressionStatement := &ast.ExpressionStatement{}
	var err error
	expressionStatement.Value, err = p.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}
	if p.peekToken.Type == token.Semicolon {
		p.nextToken()
	}
	return expressionStatement, nil
}

func (p *Parser) parseExpression(d precedence) (ast.Expression, error) {
	prefixFn, ok := p.prefixParseFnMap[p.currentToken.Type]
	if !ok {
		return nil, errNoPrefixParseFunction{t: p.currentToken}
	}
	expr, err := prefixFn()
	if err != nil {
		return nil, err
	}
	for p.peekToken.Type != token.Semicolon && d < p.peekPrecedence() {
		p.nextToken()
		infixFn, ok := p.infixParseFnMap[p.currentToken.Type]
		if !ok {
			return nil, errNoInfixParseFunction{t: p.currentToken}
		}
		expr, err = infixFn(expr)
		if err != nil {
			return nil, err
		}
	}
	return expr, nil
}

func (p *Parser) parseParameterList() ([]*ast.Identifier, error) {
	if p.peekToken.Type == token.RParen {
		p.nextToken()
		return nil, nil
	}

	var list []*ast.Identifier

	p.nextToken()
	if p.currentToken.Type != token.Ident {
		return nil, errUnexpectedToken{t: p.currentToken, exp: "identifier"}
	}
	list = append(list, &ast.Identifier{Value: p.currentToken.Literal})

	for p.peekToken.Type == token.Comma {
		p.nextToken()
		p.nextToken()
		if p.currentToken.Type != token.Ident {
			return nil, errUnexpectedToken{t: p.currentToken, exp: "identifier"}
		}
		list = append(list, &ast.Identifier{Value: p.currentToken.Literal})
	}

	if p.peekToken.Type != token.RParen {
		return nil, errUnexpectedToken{t: p.peekToken, exp: ")"}
	}
	p.nextToken()

	return list, nil
}

func (p *Parser) parseExpressionList(end token.TokenType) ([]ast.Expression, error) {
	if p.peekToken.Type == end {
		p.nextToken()
		return nil, nil
	}

	var list []ast.Expression

	p.nextToken()
	expr, err := p.parseExpression(Lowest)
	if err != nil {
		return nil, err
	}
	list = append(list, expr)

	for p.peekToken.Type == token.Comma {
		p.nextToken()
		p.nextToken()
		expr, err = p.parseExpression(Lowest)
		if err != nil {
			return nil, err
		}
		list = append(list, expr)
	}

	if p.peekToken.Type != end {
		return nil, errUnexpectedToken{t: p.peekToken, exp: string(end)}
	}
	p.nextToken()

	return list, nil
}
