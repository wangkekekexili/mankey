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
	return Eval(program)
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

func TestEvalInteger(t *testing.T) {
	tests := []struct {
		code   string
		expInt int64
	}{
		{"5", 5},
		{"42", 42},
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
