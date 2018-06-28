package evaluator

import (
	"github.com/wangkekekexili/mankey/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				panic("len accept one argument")
			}
			switch e := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(e.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(e.Elements))}
			default:
				panic("unexpected object for len")
			}
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) <= 1 {
				panic("push accepts at least 2 arguments")
			}
			arr, ok := args[0].(*object.Array)
			if !ok {
				panic("the first argument for push must be an array")
			}
			newArr := &object.Array{
				Elements: make([]object.Object, 0, len(arr.Elements)+len(args)-1),
			}
			for _, e := range arr.Elements {
				newArr.Elements = append(newArr.Elements, e)
			}
			newArr.Elements = append(newArr.Elements, args[1:]...)
			return newArr
		},
	},
}
