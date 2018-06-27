package evaluator

import (
	"errors"
	"fmt"

	"github.com/wangkekekexili/mankey/ast"
	"github.com/wangkekekexili/mankey/object"
)

func Eval(node ast.Node, env *object.Environment) (object.Object, error) {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.VarStatement:
		return evalVarStatement(node, env)
	case *ast.ReturnStatement:
		return evalReturnStatement(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Value, env)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node, env)
	case *ast.InfixExpression:
		return evalInfixExpression(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.Function:
		return evalFunction(node, env), nil
	case *ast.CallExpression:
		return evalCallExpression(node, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.Integer:
		return &object.Integer{Value: node.Value}, nil
	case *ast.Boolean:
		return evalBoolean(node.Value), nil
	case *ast.String:
		return &object.String{Value: node.Value}, nil
	default:
		return nil, fmt.Errorf("cannot evaluate %T", node)
	}
}

func evalProgram(node *ast.Program, env *object.Environment) (object.Object, error) {
	if len(node.Statements) == 0 {
		return object.Null, nil
	}
	var result object.Object
	var err error
	for _, stat := range node.Statements {
		result, err = Eval(stat, env)
		if err != nil {
			return nil, err
		}
		returnValue, ok := result.(*object.ReturnValue)
		if ok {
			return returnValue.Value, nil
		}
	}
	return result, nil
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) (object.Object, error) {
	if len(block.Statements) == 0 {
		return object.Null, nil
	}
	var result object.Object
	var err error
	for _, stat := range block.Statements {
		result, err = Eval(stat, env)
		if err != nil {
			return nil, err
		}
		if result.Type() == object.ObjReturnValue {
			return result, nil
		}
	}
	return result, nil
}

func evalVarStatement(node *ast.VarStatement, env *object.Environment) (object.Object, error) {
	o, err := Eval(node.Value, env)
	if err != nil {
		return nil, err
	}
	env.Set(node.Name.Value, o)
	return o, nil
}

func evalReturnStatement(node *ast.ReturnStatement, env *object.Environment) (object.Object, error) {
	o, err := Eval(node.Value, env)
	if err != nil {
		return nil, err
	}
	return &object.ReturnValue{Value: o}, nil
}

func evalPrefixExpression(n *ast.PrefixExpression, env *object.Environment) (object.Object, error) {
	value, err := Eval(n.Value, env)
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

func evalInfixExpression(n *ast.InfixExpression, env *object.Environment) (object.Object, error) {
	left, err := Eval(n.Left, env)
	if err != nil {
		return nil, err
	}
	right, err := Eval(n.Right, env)
	if err != nil {
		return nil, err
	}
	switch {
	case left.Type() == object.ObjInteger && right.Type() == object.ObjInteger:
		return evalIntegerInfixExpression(n.Op, left.(*object.Integer).Value, right.(*object.Integer).Value)
	case left.Type() == object.ObjBoolean && right.Type() == object.ObjBoolean:
		return evalBooleanInfixExpression(n.Op, left.(*object.Boolean).Value, right.(*object.Boolean).Value)
	case left.Type() == object.ObjString && right.Type() == object.ObjString:
		return evalStringInfixExpression(n.Op, left.(*object.String).Value, right.(*object.String).Value)
	default:
		return nil, fmt.Errorf("unsupported operator %v for operands %v and %v", n.Op, left, right)
	}
}

func evalIfExpression(ifExpression *ast.IfExpression, env *object.Environment) (object.Object, error) {
	cond, err := Eval(ifExpression.Condition, env)
	if err != nil {
		return nil, err
	}
	condBool, ok := cond.(*object.Boolean)
	if !ok {
		return nil, fmt.Errorf("non-boolean value for the if expression")
	}
	if condBool.Value {
		return evalBlockStatement(ifExpression.Consequence, env)
	} else {
		if ifExpression.Alternative == nil {
			return object.Null, nil
		} else {
			return evalBlockStatement(ifExpression.Alternative, env)
		}
	}
}

func evalFunction(fn *ast.Function, env *object.Environment) object.Object {
	return &object.Function{
		Parameters: fn.Parameters,
		Body:       fn.Body,
		Env:        env,
	}
}

func evalCallExpression(call *ast.CallExpression, env *object.Environment) (object.Object, error) {
	functionObj, err := Eval(call.Function, env)
	if err != nil {
		return nil, err
	}
	exprs, err := evalExpressions(call.Arguments, env)
	if err != nil {
		return nil, err
	}

	switch functionObj := functionObj.(type) {
	case *object.Function:
		if len(functionObj.Parameters) != len(call.Arguments) {
			return nil, fmt.Errorf("function expects %v parameter; %v provided", len(functionObj.Parameters), len(call.Arguments))
		}
		enclosedEnv := object.NewEnclosedEnvironment(functionObj.Env)
		for i := range functionObj.Parameters {
			enclosedEnv.Set(functionObj.Parameters[i].Value, exprs[i])
		}
		return unwrapReturnObject(evalBlockStatement(functionObj.Body, enclosedEnv))
	case *object.Builtin:
		return functionObj.Fn(exprs...), nil
	default:
		return nil, fmt.Errorf("unknown type of function %T", functionObj)
	}
}

func unwrapReturnObject(o object.Object, err error) (object.Object, error) {
	if err != nil {
		return nil, err
	}
	returnObj, ok := o.(*object.ReturnValue)
	if ok {
		o = returnObj.Value
	}
	return o, nil
}

func evalExpressions(exprs []ast.Expression, env *object.Environment) ([]object.Object, error) {
	var result []object.Object
	for _, expr := range exprs {
		o, err := Eval(expr, env)
		if err != nil {
			return nil, err
		}
		result = append(result, o)
	}
	return result, nil
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

func evalStringInfixExpression(op ast.Operator, left, right string) (object.Object, error) {
	if op == "+" {
		return &object.String{Value: left + right}, nil
	} else {
		return nil, fmt.Errorf("unexpected operator %v for string operands", op)
	}
}

func evalBoolean(v bool) object.Object {
	if v {
		return object.True
	} else {
		return object.False
	}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) (object.Object, error) {
	o, ok := env.Get(node.Value)
	if ok {
		return o, nil
	}
	o, ok = builtins[node.Value]
	if ok {
		return o, nil
	}
	return nil, fmt.Errorf("undefined identifier %v", node.Value)
}
