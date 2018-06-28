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

func TestStringLiteralExpression(t *testing.T) {
	code := `"hello world";`
	expressionStat, err := assertOneExpressionStatement(code)
	if err != nil {
		t.Fatal(err)
	}
	str, ok := expressionStat.Value.(*ast.String)
	if !ok {
		t.Fatalf("expected to get a string; got %T", expressionStat.Value)
	}
	if str.Value != "hello world" {
		t.Fatalf("expected to get string '%v'; got %v", "hello world", str.Value)
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
		{
			expr: "a + add(b * c) + d",
			expExpression: &ast.InfixExpression{
				Left: &ast.InfixExpression{
					Left: &ast.Identifier{Value: "a"},
					Op:   "+",
					Right: &ast.CallExpression{
						Function: &ast.Identifier{Value: "add"},
						Arguments: []ast.Expression{
							&ast.InfixExpression{
								Left:  &ast.Identifier{Value: "b"},
								Op:    "*",
								Right: &ast.Identifier{Value: "c"},
							},
						},
					},
				},
				Op:    "+",
				Right: &ast.Identifier{Value: "d"},
			},
		},
		{
			expr: "a * b[2*1]",
			expExpression: &ast.InfixExpression{
				Left: &ast.Identifier{Value: "a"},
				Op:   "*",
				Right: &ast.IndexExpression{
					Left: &ast.Identifier{Value: "b"},
					Index: &ast.InfixExpression{
						Left:  &ast.Integer{Value: 2},
						Op:    "*",
						Right: &ast.Integer{Value: 1},
					},
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

func TestFunction(t *testing.T) {
	expressionStat, err := assertOneExpressionStatement("func (x, y) {return x+y;}")
	if err != nil {
		t.Fatal(err)
	}
	function, ok := expressionStat.Value.(*ast.Function)
	if !ok {
		t.Fatalf("expect to get a function; got %T", expressionStat.Value)
	}
	if len(function.Parameters) != 2 {
		t.Fatalf("expect to get 2 parameters; got %v", function.Parameters)
	}
	if function.Parameters[0].Value != "x" {
		t.Fatalf("expect first parameter to be x; got %v", function.Parameters[0].Value)
	}
	if function.Parameters[1].Value != "y" {
		t.Fatalf("expect second parameter to be y; got %v", function.Parameters[1].Value)
	}
	if len(function.Body.Statements) != 1 {
		t.Fatalf("expect 1 statement in the block; got %v", function.Body.Statements)
	}
}

func TestFunctionParameters(t *testing.T) {
	tests := []struct {
		code          string
		expParameters []string
	}{
		{"func() {return 0;}", nil},
		{"func(a) {a;}", []string{"a"}},
		{"func(x, y) {x+y;}", []string{"x", "y"}},
	}
	for _, test := range tests {
		var expIdentifiers []*ast.Identifier
		for _, para := range test.expParameters {
			expIdentifiers = append(expIdentifiers, &ast.Identifier{Value: para})
		}

		expressionStat, err := assertOneExpressionStatement(test.code)
		if err != nil {
			t.Fatal(err)
		}
		function, ok := expressionStat.Value.(*ast.Function)
		if !ok {
			t.Fatalf("expect to get a function; got %T", expressionStat.Value)
		}
		if !reflect.DeepEqual(function.Parameters, expIdentifiers) {
			t.Fatalf("got parameters %v; want %v", function.Parameters, expIdentifiers)
		}
	}
}

func TestCallExpression(t *testing.T) {
	expressionStat, err := assertOneExpressionStatement("add(a, 1+2)")
	if err != nil {
		t.Fatal(err)
	}
	callExpression, ok := expressionStat.Value.(*ast.CallExpression)
	if !ok {
		t.Fatalf("expected to get a call expression; got %T", expressionStat.Value)
	}
	err = assertIdentifier(callExpression.Function, "add")
	if err != nil {
		t.Fatal(err)
	}
	if len(callExpression.Arguments) != 2 {
		t.Fatalf("expected to get 2 arguments; got %v", callExpression.Arguments)
	}
}

func TestCallArguments(t *testing.T) {
	tests := []struct {
		code         string
		expArguments []string
	}{
		{"print()", nil},
		{"abs(x)", []string{"x"}},
		{"minus(b, d)", []string{"b", "d"}},
	}
	for _, test := range tests {
		var expIdentifiers []ast.Expression
		for _, para := range test.expArguments {
			expIdentifiers = append(expIdentifiers, &ast.Identifier{Value: para})
		}

		expressionStat, err := assertOneExpressionStatement(test.code)
		if err != nil {
			t.Fatal(err)
		}
		call, ok := expressionStat.Value.(*ast.CallExpression)
		if !ok {
			t.Fatalf("expect to get a call expression; got %T", expressionStat.Value)
		}
		if !reflect.DeepEqual(call.Arguments, expIdentifiers) {
			t.Fatalf("got arguments %v; want %v", call.Arguments, expIdentifiers)
		}
	}
}

func TestArray(t *testing.T) {
	code := `[1, "hello", 42];`
	expressionStat, err := assertOneExpressionStatement(code)
	if err != nil {
		t.Fatal(err)
	}
	arr, ok := expressionStat.Value.(*ast.Array)
	if !ok {
		t.Fatalf("expected to get an array; got %T", expressionStat.Value)
	}
	if len(arr.Elements) != 3 {
		t.Fatalf("expected to get 3 elements in the array; got %v", arr.Elements)
	}
}

func TestIndex(t *testing.T) {
	code := `arr[3];`
	expressionStat, err := assertOneExpressionStatement(code)
	if err != nil {
		t.Fatal(err)
	}
	index, ok := expressionStat.Value.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("expected to get an index expression; got %T", expressionStat.Value)
	}
	err = assertIdentifier(index.Left, "arr")
	if err != nil {
		t.Fatal(err)
	}
	err = assertIntegerLiteral(index.Index, 3)
	if err != nil {
		t.Fatal(err)
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
