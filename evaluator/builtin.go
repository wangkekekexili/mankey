package evaluator

import "github.com/wangkekekexili/mankey/object"

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
}
