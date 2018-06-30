package object

import "strconv"

var (
	True  = &Boolean{Value: true}
	False = &Boolean{Value: false}
)

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

func (b *Boolean) HashKey() HashKey {
	if b.Value {
		return ObjBoolean + "_true"
	} else {
		return ObjBoolean + "_false"
	}
}
