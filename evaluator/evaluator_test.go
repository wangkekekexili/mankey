package evaluator

import (
	"fmt"
	"testing"

	"github.com/wangkekekexili/mankey/lexer"
	"github.com/wangkekekexili/mankey/object"
	"github.com/wangkekekexili/mankey/parser"
)

func eval(code string) (object.Object, error) {
	program, err := parser.New(lexer.New(code)).ParseProgram()
	if err != nil {
		return nil, err
	}
	return Eval(program, object.NewEnvironment())
}

func assertIntegerObject(o object.Object, v int64) error {
	integer, ok := o.(*object.Integer)
	if !ok {
		return fmt.Errorf("expected to get an integer object; got %T", o)
	}
	if integer.Value != v {
		return fmt.Errorf("got integer value %v; want %v", integer.Value, v)
	}
	return nil
}

func assertBoolObject(o object.Object, v bool) error {
	boolean, ok := o.(*object.Boolean)
	if !ok {
		return fmt.Errorf("expected to get a boolean object; got %T", o)
	}
	if boolean.Value != v {
		return fmt.Errorf("got boolean value %v; want %v", boolean.Value, v)
	}
	return nil
}

func assertNullIntBool(o object.Object, v interface{}) error {
	if v == nil {
		if o == object.Null {
			return nil
		} else {
			return fmt.Errorf("expected to get null; got %v", o)
		}
	}
	switch v := v.(type) {
	case int:
		return assertIntegerObject(o, int64(v))
	case bool:
		return assertBoolObject(o, v)
	default:
		return fmt.Errorf("bug in test: unexpected type %T", v)
	}
}

func TestEvalInteger(t *testing.T) {
	tests := []struct {
		code   string
		expInt int64
	}{
		{"5", 5},
		{"-5", -5},
		{"--5", 5},
		{"---5", -5},
		{"42", 42},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}
	for _, test := range tests {
		o, err := eval(test.code)
		if err != nil {
			t.Fatal(err)
		}
		err = assertIntegerObject(o, test.expInt)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestEvalBoolean(t *testing.T) {
	tests := []struct {
		code    string
		expBool bool
	}{
		{"true", true},
		{"false", false},
		{"!false", true},
		{"!!false", false},
		{"!true", false},
		{"!!true", true},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}
	for _, test := range tests {
		o, err := eval(test.code)
		if err != nil {
			t.Fatal(err)
		}
		err = assertBoolObject(o, test.expBool)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestEvalIfElseExpression(t *testing.T) {
	tests := []struct {
		code string
		exp  interface{}
	}{
		{"if (true) { 10 }", 10},
		{"if (false) { 10; }", nil},
		{"if (1 < 2) { 10; }", 10},
		{"if (1 > 2) { 10 }", nil},
		{"if (1 > 2) { true } else { false }", false},
		{"if (1 < 2) { true } else { false }", true},
	}
	for _, test := range tests {
		o, err := eval(test.code)
		if err != nil {
			t.Fatal(err)
		}
		err = assertNullIntBool(o, test.exp)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		code   string
		expInt int64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{
			code: `
if (true) {
  if (true) {
    return 42;
  }
}
return 10;
`,
			expInt: 42,
		},
	}
	for _, test := range tests {
		o, err := eval(test.code)
		if err != nil {
			t.Fatal(err)
		}
		err = assertIntegerObject(o, test.expInt)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestEvalVarStatement(t *testing.T) {
	tests := []struct {
		code   string
		expInt int64
	}{
		{"var n = 42;", 42},
		{"var n = 42; n+1", 43},
		{"var n = 42; var m = n - 2; m", 40},
	}
	for _, test := range tests {
		o, err := eval(test.code)
		if err != nil {
			t.Fatal(err)
		}
		err = assertIntegerObject(o, test.expInt)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestEvalFunction(t *testing.T) {
	o, err := eval("func(a) {return a;}")
	if err != nil {
		t.Fatal(err)
	}
	fn, ok := o.(*object.Function)
	if !ok {
		t.Fatalf("expected to get a function object; got %T", o)
	}
	if len(fn.Parameters) != 1 {
		t.Fatalf("expected to get 1 parameter; got %v", fn.Parameters)
	}
	if len(fn.Body.Statements) != 1 {
		t.Fatalf("expected to get 1 statement in the body; got %v", fn.Body)
	}
}

func TestCallExpression(t *testing.T) {
	tests := []struct {
		code   string
		expInt int64
	}{
		{"var fn = func(x) {return x;};fn(42)", 42},
		//{"var double = func(x) { x * 2; }; double(5);", 10},
		//{"var add = func(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
	}
	for _, test := range tests {
		o, err := eval(test.code)
		if err != nil {
			t.Fatal(err)
		}
		err = assertIntegerObject(o, test.expInt)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestError(t *testing.T) {
	codes := []string{
		"!10",
		"-true",
		"true + false",
		"if (1) {1}",
		"foobar",
	}
	for _, code := range codes {
		_, err := eval(code)
		if err == nil {
			t.Fatal("error expected")
		}
	}
}

func TestClosures(t *testing.T) {
	code := `
var newAdder = func(x) {
	func(y) { x + y };
};
var addTwo = newAdder(2);
addTwo(9);`
	o, err := eval(code)
	if err != nil {
		t.Fatal(err)
	}
	assertIntegerObject(o, 11)
}
