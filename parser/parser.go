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
		token.Ident:  p.parseIdentifier,
		token.Number: p.parseIntegerLiteral,
		token.Minus:  p.parsePrefixExpression,
		token.Not:    p.parsePrefixExpression,
		token.True:   p.parseBoolean,
		token.False:  p.parseBoolean,
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

	// skip expression for now
	for p.currentToken.Type != token.Semicolon {
		p.nextToken()
	}

	return varStat, nil
}

func (p *Parser) parseReturnStatement() (*ast.ReturnStatement, error) {
	returnStatement := &ast.ReturnStatement{}

	p.nextToken()
	for p.currentToken.Type != token.Semicolon {
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
