package object

import "strconv"

const ObjBoolean = "BOOLEAN"

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return ObjBoolean
}

func (b *Boolean) String() string {
	return strconv.FormatBool(b.Value)
}
