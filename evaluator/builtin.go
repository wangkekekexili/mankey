package evaluator

import "github.com/wangkekekexili/mankey/object"

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				panic("len accept one argument")
			}
			str, ok := args[0].(*object.String)
			if !ok {
				panic("string expected for len")
			}
			return &object.Integer{Value: int64(len(str.Value))}
		},
	},
}
