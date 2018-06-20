package object

const ObjNull = "NULL"

type Null struct{}

func (n *Null) Type() ObjectType {
	return ObjNull
}

func (n *Null) String() string {
	return ObjNull
}
