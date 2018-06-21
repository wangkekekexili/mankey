package evaluator

import (
	"fmt"

	"github.com/wangkekekexili/mankey/ast"
	"github.com/wangkekekexili/mankey/object"
)

func Eval(node ast.Node) (object.Object, error) {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Value)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node)
	case *ast.Integer:
		return &object.Integer{Value: node.Value}, nil
	case *ast.Boolean:
		return evalBoolean(node.Value), nil
	default:
		return nil, nil
	}
}

func evalStatements(stats []ast.Statement) (object.Object, error) {
	if len(stats) == 0 {
		return object.Null, nil
	}
	var result object.Object
	var err error
	for _, stat := range stats {
		result, err = Eval(stat)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func evalPrefixExpression(n *ast.PrefixExpression) (object.Object, error) {
	value, err := Eval(n.Value)
	if err != nil {
		return nil, err
	}
	switch n.Op {
	case "!":
		boolean, ok := value.(*object.Boolean)
		if !ok {
			return nil, fmt.Errorf("'!' only works on boolean value")
		}
		return evalBoolean(!boolean.Value), nil
	case "-":
		integer, ok := value.(*object.Integer)
		if !ok {
			return nil, fmt.Errorf("'-' only works on integer value")
		}
		return &object.Integer{Value: -integer.Value}, nil
	default:
		return nil, fmt.Errorf("unknown prefix operator: %v", n.Op)
	}
}

func evalBoolean(v bool) object.Object {
	if v {
		return object.True
	} else {
		return object.False
	}
}
