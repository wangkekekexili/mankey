package lexer

import (
	"github.com/wangkekekexili/mankey/token"
)

type Lexer struct {
	input string
	pos   int // last checked position
}

func New(input string) *Lexer {
	return &Lexer{
		input: input,
		pos:   -1,
	}
}

func (r *Lexer) currentChar() (byte, bool) {
	if r.pos >= 0 && r.pos < len(r.input) {
		return r.input[r.pos], true
	} else {
		return 0, false
	}
}

func (r *Lexer) mustCurrentChar() byte {
	b, ok := r.currentChar()
	if !ok {
		panic("failed to get current char")
	}
	return b
}

func (r *Lexer) nextChar() (byte, bool) {
	next := r.pos + 1
	if next >= len(r.input) {
		return 0, false
	}
	r.pos = next
	return r.input[r.pos], true
}

func (r *Lexer) peekNextChar() (byte, bool) {
	next := r.pos + 1
	if next >= len(r.input) {
		return 0, false
	}
	return r.input[next], true
}

func (r *Lexer) advance() {
	if r.pos < len(r.input) {
		r.pos++
	}
}

func (r *Lexer) mustCurrentIdentifier() string {
	start := r.pos
	for {
		b, ok := r.peekNextChar()
		if !ok {
			return r.input[start : r.pos+1]
		}
		if isLetter(b) {
			r.advance()
			continue
		} else {
			return r.input[start : r.pos+1]
		}
	}
}

func (r *Lexer) mustCurrentNumber() string {
	start := r.pos
	for {
		b, ok := r.peekNextChar()
		if !ok {
			return r.input[start : r.pos+1]
		}
		if isDigit(b) {
			r.advance()
			continue
		} else {
			return r.input[start : r.pos+1]
		}
	}
}

func (r *Lexer) mustCurrentString() string {
	// start points to the starting quote.
	start := r.pos

	for {
		r.advance()
		ch, ok := r.currentChar()
		if !ok {
			panic("unfinishing string")
		}
		if ch == '"' {
			break
		}
	}

	return r.input[start+1 : r.pos]
}

func (r *Lexer) skipWhitespace() {
	for {
		b, ok := r.peekNextChar()
		if !ok {
			return
		}
		if b == ' ' || b == '\t' || b == '\n' {
			r.advance()
			continue
		} else {
			break
		}
	}
}

func (r *Lexer) NextToken() *token.Token {
	r.skipWhitespace()
	b, ok := r.nextChar()
	if !ok {
		return token.New(token.EOF, "")
	}
	switch b {
	case '"':
		return token.New(token.String, r.mustCurrentString())
	case '=':
		n, ok := r.peekNextChar()
		if ok && n == '=' {
			r.advance()
			return token.New(token.Equal, "==")
		} else {
			return token.New(token.Assign, "=")
		}
	case '!':
		n, ok := r.peekNextChar()
		if ok && n == '=' {
			r.advance()
			return token.New(token.NotEqual, "!=")
		} else {
			return token.New(token.Not, "!")
		}
	case '>':
		n, ok := r.peekNextChar()
		if ok && n == '=' {
			r.advance()
			return token.New(token.Gte, ">=")
		} else {
			return token.New(token.Gt, ">")
		}
	case '<':
		n, ok := r.peekNextChar()
		if ok && n == '=' {
			r.advance()
			return token.New(token.Lte, "<=")
		} else {
			return token.New(token.Lt, "<")
		}
	case '+':
		return token.New(token.Add, "+")
	case '-':
		return token.New(token.Minus, "-")
	case '/':
		return token.New(token.Divide, "/")
	case '*':
		return token.New(token.Multiply, "*")
	case ',':
		return token.New(token.Comma, ",")
	case ';':
		return token.New(token.Semicolon, ";")
	case '(':
		return token.New(token.LParen, "(")
	case ')':
		return token.New(token.RParen, ")")
	case '[':
		return token.New(token.LBracket, "[")
	case ']':
		return token.New(token.RBracket, "]")
	case '{':
		return token.New(token.LBrace, "{")
	case '}':
		return token.New(token.RBrace, "}")
	default: // other non single character operation
		if isLetter(b) {
			ident := r.mustCurrentIdentifier()
			identType := token.LookupIdent(ident)
			return token.New(identType, ident)
		} else if isDigit(b) {
			n := r.mustCurrentNumber()
			return token.New(token.Number, n)
		}
	}
	return token.New(token.Illegal, "")
}

func isLetter(b byte) bool {
	return b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z' || b == '_'
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}
