package parser

import (
	"fmt"
	"reflect"
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
	expressionStat, err := assertOneExpressionStatement(code)
	if err != nil {
		t.Fatal(err)
	}
	err = assertIdentifier(expressionStat.Value, "apple")
	if err != nil {
		t.Fatal(err)
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	code := `42;`
	expressionStat, err := assertOneExpressionStatement(code)
	if err != nil {
		t.Fatal(err)
	}
	err = assertIntegerLiteral(expressionStat.Value, 42)
	if err != nil {
		t.Fatal(err)
	}
}

func TestPrefixExpression(t *testing.T) {
	tests := []struct {
		code     string
		expOp    ast.Operator
		expValue int
	}{
		{"-5;", "-", 5},
		{"!10", "!", 10},
	}
	for _, test := range tests {
		expressionStat, err := assertOneExpressionStatement(test.code)
		if err != nil {
			t.Fatal(err)
		}
		prefixExpression, ok := expressionStat.Value.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("expect to get an prefix expression; got %T", expressionStat.Value)
		}
		if prefixExpression.Op != test.expOp {
			t.Fatalf("got operator %v; want %v", prefixExpression.Op, test.expOp)
		}
		assertIntegerLiteral(prefixExpression.Value, test.expValue)
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		code    string
		expBool bool
	}{
		{"true;", true},
		{"false;", false},
	}
	for _, test := range tests {
		expressionStat, err := assertOneExpressionStatement(test.code)
		if err != nil {
			t.Fatal(err)
		}
		err = assertBoolean(expressionStat.Value, test.expBool)
		if err != nil {
			t.Fatal(err)
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
		expressionStat, err := assertOneExpressionStatement(test.code)
		if err != nil {
			t.Fatal(err)
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

func TestParseOperatorPrecedence(t *testing.T) {
	tests := []struct {
		expr          string
		expExpression ast.Expression
	}{
		{
			expr:          "42",
			expExpression: &ast.Integer{Value: 42},
		},
		{
			expr: "-42",
			expExpression: &ast.PrefixExpression{
				Op:    "-",
				Value: &ast.Integer{Value: 42},
			},
		},
		{
			expr: "5+6",
			expExpression: &ast.InfixExpression{
				Left:  &ast.Integer{Value: 5},
				Op:    "+",
				Right: &ast.Integer{Value: 6},
			},
		},
		{
			expr: "5+6",
			expExpression: &ast.InfixExpression{
				Left:  &ast.Integer{Value: 5},
				Op:    "+",
				Right: &ast.Integer{Value: 6},
			},
		},
		{
			expr: "a+b*c",
			expExpression: &ast.InfixExpression{
				Left: &ast.Identifier{Value: "a"},
				Op:   "+",
				Right: &ast.InfixExpression{
					Left:  &ast.Identifier{Value: "b"},
					Op:    "*",
					Right: &ast.Identifier{Value: "c"},
				},
			},
		},
		{
			expr: "1*3 == 3",
			expExpression: &ast.InfixExpression{
				Left: &ast.InfixExpression{
					Left:  &ast.Integer{Value: 1},
					Op:    "*",
					Right: &ast.Integer{Value: 3},
				},
				Op:    "==",
				Right: &ast.Integer{Value: 3},
			},
		},
		{
			expr: "(2+2)*3",
			expExpression: &ast.InfixExpression{
				Left: &ast.InfixExpression{
					Left:  &ast.Integer{Value: 2},
					Op:    "+",
					Right: &ast.Integer{Value: 2},
				},
				Op:    "*",
				Right: &ast.Integer{Value: 3},
			},
		},
		{
			expr: "!(true == false)",
			expExpression: &ast.PrefixExpression{
				Op: "!",
				Value: &ast.InfixExpression{
					Left:  &ast.Boolean{Value: true},
					Op:    "==",
					Right: &ast.Boolean{Value: false},
				},
			},
		},
	}
	for _, test := range tests {
		gotProgram, err := New(lexer.New(test.expr + ";")).ParseProgram()
		if err != nil {
			t.Fatal(err)
		}
		expProgram := &ast.Program{
			Statements: []ast.Statement{
				&ast.ExpressionStatement{
					Value: test.expExpression,
				},
			},
		}
		if !reflect.DeepEqual(expProgram, gotProgram) {
			t.Fatalf("expected to get %v; got %v", expProgram, gotProgram)
		}
	}
}

func TestIfExpression(t *testing.T) {
	expressionStat, err := assertOneExpressionStatement("if (a < b) { 42 }")
	if err != nil {
		t.Fatal(err)
	}
	ifExpression, ok := expressionStat.Value.(*ast.IfExpression)
	if !ok {
		t.Fatalf("expect to get an if expression; got %T", expressionStat.Value)
	}
	_, ok = ifExpression.Condition.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("expect to get an infix expression; got %T", ifExpression.Condition)
	}
	if ifExpression.Consequence == nil {
		t.Fatal("expect non-nil consequence")
	}
	if ifExpression.Alternative != nil {
		t.Fatal("expect nil alternative")
	}
}

func TestIfElseExpression(t *testing.T) {
	expressionStat, err := assertOneExpressionStatement("if (a < b) { 42 } else { a; b;}")
	if err != nil {
		t.Fatal(err)
	}
	ifExpression, ok := expressionStat.Value.(*ast.IfExpression)
	if !ok {
		t.Fatalf("expect to get an if expression; got %T", expressionStat.Value)
	}
	_, ok = ifExpression.Condition.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("expect to get an infix expression; got %T", ifExpression.Condition)
	}
	if ifExpression.Consequence == nil {
		t.Fatal("expect non-nil consequence")
	}
	if ifExpression.Alternative == nil {
		t.Fatal("expect non-nil alternative")
	}
	if len(ifExpression.Alternative.Statements) != 2 {
		t.Fatalf("expect to get 2 statements in the else block; got %v", ifExpression.Alternative)
	}
}

func assertOneExpressionStatement(code string) (*ast.ExpressionStatement, error) {
	p, err := New(lexer.New(code)).ParseProgram()
	if err != nil {
		return nil, err
	}
	if p == nil || len(p.Statements) != 1 {
		return nil, fmt.Errorf("expect to get 1 statement; got %v", p)
	}
	expressionStatement, ok := p.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		return nil, fmt.Errorf("expect to get an expression statement; got %T", p.Statements[0])
	}
	return expressionStatement, nil
}

func assertIntegerLiteral(e ast.Expression, vi interface{}) error {
	var v int64
	switch vi := vi.(type) {
	case int:
		v = int64(vi)
	case int64:
		v = int64(vi)
	default:
		return fmt.Errorf("assertIntegerLiteral accepts an integer literal as the second parameter")
	}
	intLiteral, ok := e.(*ast.Integer)
	if !ok {
		return fmt.Errorf("expect to get interger literal; got %T", e)
	}
	if intLiteral.Value != v {
		return fmt.Errorf("expect to get %v; got %v", v, intLiteral.Value)
	}
	return nil
}

func assertIdentifier(e ast.Expression, v string) error {
	identifier, ok := e.(*ast.Identifier)
	if !ok {
		return fmt.Errorf("expect to get an identifier; got %T", e)
	}
	if identifier.Value != v {
		return fmt.Errorf("expect to get %v; got %v", v, identifier.Value)
	}
	return nil
}

func assertBoolean(e ast.Expression, v bool) error {
	boolean, ok := e.(*ast.Boolean)
	if !ok {
		return fmt.Errorf("expect to get a boolean; got %T", e)
	}
	if boolean.Value != v {
		return fmt.Errorf("expect to get %v; got %v", v, boolean.Value)
	}
	return nil
}
