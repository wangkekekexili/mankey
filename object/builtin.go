package object

type BuiltinFunction func(...Object) Object

const ObjBuiltin = "BUILTIN"

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType {
	return ObjBuiltin
}

func (b *Builtin) String() string {
	return ObjBuiltin
}
