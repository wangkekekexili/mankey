package object

import "strings"

const ObjArray = "ARRAY"

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType {
	return ObjArray
}

func (a *Array) String() string {
	var strs []string
	for _, element := range a.Elements {
		strs = append(strs, element.String())
	}
	return "[" + strings.Join(strs, ",") + "]"
}
