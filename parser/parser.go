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
}

func New(r *lexer.Lexer) *Parser {
	p := &Parser{r: r}
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
	default:
		return nil, unexpectedToken{exp: "statement", t: p.currentToken}
	}
}

func (p *Parser) parseVarStatement() (*ast.VarStatement, error) {
	varStat := &ast.VarStatement{}

	p.nextToken()
	if p.currentToken.Type != token.Ident {
		return nil, unexpectedToken{exp: "var statement", t: p.currentToken}
	}
	varStat.Name = &ast.Identifier{Value: p.currentToken.Literal}

	p.nextToken()
	if p.currentToken.Type != token.Assign {
		return nil, unexpectedToken{exp: "=", t: p.currentToken}
	}

	// skip expression for now
	for p.currentToken.Type != token.Semicolon {
		p.nextToken()
	}

	return varStat, nil
}
