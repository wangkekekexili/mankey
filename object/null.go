package object

var Null = &null{}

const ObjNull = "NULL"

type null struct{}

func (n *null) Type() ObjectType {
	return ObjNull
}

func (n *null) String() string {
	return ObjNull
}
