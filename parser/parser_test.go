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
			t.Fatalf("expect var statement; got %v", stat)
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
			t.Fatalf("expect var statement; got %v", stat)
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
		t.Fatal("expect to get an expression statement")
	}
	identifier, ok := expressionStat.Value.(*ast.Identifier)
	if !ok {
		t.Fatal("expect to get an identifier")
	}
	if identifier.Value != "apple" {
		t.Fatalf("got identifer value %v; want %v", identifier.Value, "apple")
	}
}
