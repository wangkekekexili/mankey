package evaluator

import (
	"errors"
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
	case *ast.InfixExpression:
		return evalInfixExpression(node)
	case *ast.IfExpression:
		return evalIfExpression(node)
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

func evalInfixExpression(n *ast.InfixExpression) (object.Object, error) {
	left, err := Eval(n.Left)
	if err != nil {
		return nil, err
	}
	right, err := Eval(n.Right)
	if err != nil {
		return nil, err
	}
	switch {
	case left.Type() == object.ObjInteger && right.Type() == object.ObjInteger:
		return evalIntegerInfixExpression(n.Op, left.(*object.Integer).Value, right.(*object.Integer).Value)
	case left.Type() == object.ObjBoolean && right.Type() == object.ObjBoolean:
		return evalBooleanInfixExpression(n.Op, left.(*object.Boolean).Value, right.(*object.Boolean).Value)
	default:
		return nil, fmt.Errorf("unsupported operator %v for operands %v and %v", n.Op, left, right)
	}
}

func evalIfExpression(ifExpression *ast.IfExpression) (object.Object, error) {
	cond, err := Eval(ifExpression.Condition)
	if err != nil {
		return nil, err
	}
	condBool, ok := cond.(*object.Boolean)
	if !ok {
		return nil, fmt.Errorf("non-boolean value for the if expression")
	}
	if condBool.Value {
		return evalStatements(ifExpression.Consequence.Statements)
	} else {
		if ifExpression.Alternative == nil {
			return object.Null, nil
		} else {
			return evalStatements(ifExpression.Alternative.Statements)
		}
	}
}

func evalIntegerInfixExpression(op ast.Operator, left, right int64) (object.Object, error) {
	switch op {
	case "+":
		return &object.Integer{Value: left + right}, nil
	case "-":
		return &object.Integer{Value: left - right}, nil
	case "*":
		return &object.Integer{Value: left * right}, nil
	case "/":
		if right == 0 {
			return nil, errors.New("divide by zero")
		}
		return &object.Integer{Value: left / right}, nil
	case ">":
		return &object.Boolean{Value: left > right}, nil
	case ">=":
		return &object.Boolean{Value: left >= right}, nil
	case "<":
		return &object.Boolean{Value: left < right}, nil
	case "<=":
		return &object.Boolean{Value: left <= right}, nil
	case "==":
		return &object.Boolean{Value: left == right}, nil
	case "!=":
		return &object.Boolean{Value: left != right}, nil
	default:
		return nil, fmt.Errorf("unexpected operator %v for integer operands", op)
	}
}

func evalBooleanInfixExpression(op ast.Operator, left, right bool) (object.Object, error) {
	switch op {
	case "==":
		return &object.Boolean{Value: left == right}, nil
	case "!=":
		return &object.Boolean{Value: left != right}, nil
	default:
		return nil, fmt.Errorf("unexpected operator %v for boolean operands", op)
	}
}

func evalBoolean(v bool) object.Object {
	if v {
		return object.True
	} else {
		return object.False
	}
}
