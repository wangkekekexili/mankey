package token

import "fmt"

type TokenType string

const (
	Illegal = "ILLEGAL"
	EOF     = "EOF"

	Ident  = "IDENTIFIER"
	Number = "NUMBER"
	String = "STRING"

	Assign   = "="
	Add      = "+"
	Minus    = "-"
	Divide   = "/"
	Multiply = "*"

	Equal    = "=="
	Not      = "!"
	NotEqual = "!="

	Lt  = "<"
	Lte = "<="
	Gt  = ">"
	Gte = ">="

	Comma     = ","
	Semicolon = ";"

	LParen   = "("
	RParen   = ")"
	LBracket = "["
	RBracket = "]"
	LBrace   = "{"
	RBrace   = "}"

	Func   = "func"
	Var    = "var"
	True   = "true"
	False  = "false"
	If     = "if"
	Else   = "else"
	Return = "return"
)

var keywords = map[string]TokenType{
	"func":   Func,
	"var":    Var,
	"true":   True,
	"false":  False,
	"if":     If,
	"else":   Else,
	"return": Return,
}

func LookupIdent(ident string) TokenType {
	if typ, ok := keywords[ident]; ok {
		return typ
	} else {
		return Ident
	}
}

type Token struct {
	Type    TokenType
	Literal string
}

func New(typ TokenType, literal string) *Token {
	return &Token{
		Type:    typ,
		Literal: literal,
	}
}

func (t *Token) Equals(k *Token) bool {
	if t.Type != k.Type {
		return false
	}
	return t.Literal == k.Literal
}

func (t *Token) String() string {
	return fmt.Sprintf("[type=%v;literal=%v]", t.Type, t.Literal)
}
