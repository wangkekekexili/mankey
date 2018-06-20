package evaluator

import (
	"github.com/wangkekekexili/mankey/ast"
	"github.com/wangkekekexili/mankey/object"
)

func Eval(node ast.Node) (object.Object, error) {
	switch node := node.(type) {
	case *ast.Program:
		if len(node.Statements) == 0 {
			return object.Null, nil
		}
		var result object.Object
		var err error
		for _, s := range node.Statements {
			result, err = Eval(s)
			if err != nil {
				return nil, err
			}
		}
		return result, nil
	case *ast.ExpressionStatement:
		return Eval(node.Value)
	case *ast.Integer:
		return &object.Integer{Value: node.Value}, nil
	case *ast.Boolean:
		if node.Value {
			return object.True, nil
		} else {
			return object.False, nil
		}
	default:
		return nil, nil
	}
}
