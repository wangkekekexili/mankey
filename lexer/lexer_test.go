package lexer

import (
	"testing"

	"github.com/wangkekekexili/mankey/token"
)

func TestReadChar(t *testing.T) {
	lexer := New("abc")
	for _, exp := range []byte{'a', 'b', 'c', 0} {
		got, ok := lexer.nextChar()
		if exp == 0 {
			if ok {
				t.Fatalf("expect no more char, got %v", got)
			}
		} else {
			if !ok {
				t.Fatalf("expect to get char %v", exp)
			}
			if exp != got {
				t.Fatalf("expect to get char %v; got %v", exp, got)
			}
		}
	}
}

func TestNextToken(t *testing.T) {
	tests := []struct {
		input     string
		expTokens []*token.Token
	}{
		{
			input:     "",
			expTokens: nil,
		},
		{
			input: "{(+=-)}",
			expTokens: []*token.Token{
				token.New(token.LBrace, "{"),
				token.New(token.LParen, "("),
				token.New(token.Add, "+"),
				token.New(token.Assign, "="),
				token.New(token.Minus, "-"),
				token.New(token.RParen, ")"),
				token.New(token.RBrace, "}"),
			},
		},
		{
			input: "var   age =   10;  ",
			expTokens: []*token.Token{
				token.New(token.Var, "var"),
				token.New(token.Ident, "age"),
				token.New(token.Assign, "="),
				token.New(token.Number, "10"),
				token.New(token.Semicolon, ";"),
			},
		},
		{
			input: `var name = "ke";`,
			expTokens: []*token.Token{
				token.New(token.Var, "var"),
				token.New(token.Ident, "name"),
				token.New(token.Assign, "="),
				token.New(token.String, "ke"),
				token.New(token.Semicolon, ";"),
			},
		},
		{
			input: `var five = 5;
var ten = 10;
!-/*5;
5 < 10 > 5;
   `,
			expTokens: []*token.Token{
				token.New(token.Var, "var"),
				token.New(token.Ident, "five"),
				token.New(token.Assign, "="),
				token.New(token.Number, "5"),
				token.New(token.Semicolon, ";"),

				token.New(token.Var, "var"),
				token.New(token.Ident, "ten"),
				token.New(token.Assign, "="),
				token.New(token.Number, "10"),
				token.New(token.Semicolon, ";"),

				token.New(token.Not, "!"),
				token.New(token.Minus, "-"),
				token.New(token.Divide, "/"),
				token.New(token.Multiply, "*"),
				token.New(token.Number, "5"),
				token.New(token.Semicolon, ";"),

				token.New(token.Number, "5"),
				token.New(token.Lt, "<"),
				token.New(token.Number, "10"),
				token.New(token.Gt, ">"),
				token.New(token.Number, "5"),
				token.New(token.Semicolon, ";"),
			},
		},
		{
			input: `
if (5 < 10) {
    return true;
} else {
    return false;
}
`,
			expTokens: []*token.Token{
				token.New(token.If, "if"),
				token.New(token.LParen, "("),
				token.New(token.Number, "5"),
				token.New(token.Lt, "<"),
				token.New(token.Number, "10"),
				token.New(token.RParen, ")"),
				token.New(token.LBrace, "{"),

				token.New(token.Return, "return"),
				token.New(token.True, "true"),
				token.New(token.Semicolon, ";"),

				token.New(token.RBrace, "}"),
				token.New(token.Else, "else"),
				token.New(token.LBrace, "{"),

				token.New(token.Return, "return"),
				token.New(token.False, "false"),
				token.New(token.Semicolon, ";"),
				token.New(token.RBrace, "}"),
			},
		},
	}
	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			lexer := New(test.input)
			for _, exp := range test.expTokens {
				got := lexer.NextToken()
				if !got.Equals(exp) {
					t.Fatalf("got %v; want %v", got, exp)
				}
			}
			last := lexer.NextToken()
			if !last.Equals(token.New(token.EOF, "")) {
				t.Fatalf("expect no token left; got %v", last)
			}
		})
	}
}
