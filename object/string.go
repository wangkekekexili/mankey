package object

import "fmt"

const ObjString = "String"

type String struct {
	Value string
}

func (s *String) Type() ObjectType {
	return ObjString
}

func (s *String) String() string {
	return s.Value
}

func (s *String) HashKey() HashKey {
	return HashKey(fmt.Sprintf("%v_%v", ObjString, s.Value))
}
