package object

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
