package object

import "strconv"

var (
	True  = &boolean{Value: true}
	False = &boolean{Value: false}
)

const ObjBoolean = "BOOLEAN"

type boolean struct {
	Value bool
}

func (b *boolean) Type() ObjectType {
	return ObjBoolean
}

func (b *boolean) String() string {
	return strconv.FormatBool(b.Value)
}
