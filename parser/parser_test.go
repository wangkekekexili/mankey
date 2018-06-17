package parser

import (
	"testing"

	"github.com/wangkekekexili/mankey/ast"
	"github.com/wangkekekexili/mankey/lexer"
)

func TestString(t *testing.T) {
	program := &ast.Program{
		Statements: []ast.Statement{
			&ast.VarStatement{
				Name: &ast.Identifier{
					Value: "x",
				},
				Value: &ast.Identifier{
					Value: "y",
				},
			},
			&ast.ReturnStatement{
				Value: &ast.Identifier{
					Value: "42",
				},
			},
		},
	}
	expStr := "var x = y;return 42;"
	if program.String() != expStr {
		t.Fatalf("String(): %v; want %v", program.String(), expStr)
	}
}

func TestVarStatement(t *testing.T) {
	code := `
var hello = 1;
var world = 2;
var add = a+b;
`
	expVarIdents := []string{"hello", "world", "add"}
	program, err := New(lexer.New(code)).ParseProgram()
	if err != nil {
		t.Fatal(err)
	}
	if len(program.Statements) != len(expVarIdents) {
		t.Fatalf("got %v statements; want %v", len(program.Statements), len(expVarIdents))
	}
	for i, stat := range program.Statements {
		varStat, ok := stat.(*ast.VarStatement)
		if !ok {
			t.Fatalf("expect var statement; got %T", stat)
		}
		if varStat.Name.Value != expVarIdents[i] {
			t.Fatalf("got identifier name %v; want %v", varStat.Name.Value, expVarIdents[i])
		}
	}
}

func TestReturnStatement(t *testing.T) {
	code := `
return 1;
return add(1, 2);
`
	program, err := New(lexer.New(code)).ParseProgram()
	if err != nil {
		t.Fatal(err)
	}
	for _, stat := range program.Statements {
		_, ok := stat.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("expect var statement; got %T", stat)
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	code := `apple;`
	program, err := New(lexer.New(code)).ParseProgram()
	if err != nil {
		t.Fatal(err)
	}
	if len(program.Statements) != 1 {
		t.Fatalf("expect 1 statement; got %v", program.Statements)
	}

	expressionStat, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expect to get an expression statement; got %T", program.Statements[0])
	}
	identifier, ok := expressionStat.Value.(*ast.Identifier)
	if !ok {
		t.Fatalf("expect to get an identifier; got %T", expressionStat.Value)
	}
	if identifier.Value != "apple" {
		t.Fatalf("got identifer value %v; want %v", identifier.Value, "apple")
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	code := `42;`
	program, err := New(lexer.New(code)).ParseProgram()
	if err != nil {
		t.Fatal(err)
	}
	if len(program.Statements) != 1 {
		t.Fatalf("expect 1 statement; got %v", program.Statements)
	}

	expressionStat, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expect to get an expression statement; got %T", program.Statements[0])
	}
	integerLiteral, ok := expressionStat.Value.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expect to get an identifier; got %T", expressionStat.Value)
	}
	if integerLiteral.Value != 42 {
		t.Fatalf("got integer value %v; want %v", integerLiteral.Value, 42)
	}
}

func TestPrefixExpression(t *testing.T) {
	tests := []struct {
		code  string
		expOp ast.Operator
	}{
		{"-5;", "-"},
		{"!10", "!"},
	}
	for _, test := range tests {
		program, err := New(lexer.New(test.code)).ParseProgram()
		if err != nil {
			t.Fatal(err)
		}
		if len(program.Statements) != 1 {
			t.Fatalf("expect 1 statement; got %v", program.Statements)
		}

		expressionStat, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("expect to get an expression statement; got %T", program.Statements[0])
		}
		prefixExpression, ok := expressionStat.Value.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("expect to get an prefix expression; got %T", expressionStat.Value)
		}
		if prefixExpression.Op != test.expOp {
			t.Fatalf("got operator %v; want %v", prefixExpression.Op, test.expOp)
		}
	}
}

func TestInfixExpression(t *testing.T) {
	tests := []struct {
		code  string
		expOp ast.Operator
	}{
		{"5*5;", "*"},
		{"6/6", "/"},
		{"42 == 42", "=="},
	}
	for _, test := range tests {
		program, err := New(lexer.New(test.code)).ParseProgram()
		if err != nil {
			t.Fatal(err)
		}
		if len(program.Statements) != 1 {
			t.Fatalf("expect 1 statement; got %v", program.Statements)
		}

		expressionStat, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("expect to get an expression statement; got %T", program.Statements[0])
		}
		infixExpression, ok := expressionStat.Value.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("expect to get an infix expression; got %T", expressionStat.Value)
		}
		if infixExpression.Op != test.expOp {
			t.Fatalf("got operator %v; want %v", infixExpression.Op, test.expOp)
		}
	}
}
