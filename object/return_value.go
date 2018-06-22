package object

const ObjReturnValue = "RETURN_VALUE"

type ReturnValue struct {
	Value Object
}

func (b *ReturnValue) Type() ObjectType {
	return ObjReturnValue
}

func (b *ReturnValue) String() string {
	return b.Value.String()
}
