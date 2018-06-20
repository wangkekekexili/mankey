package object

import "strconv"

const ObjInteger = "INTEGER"

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType {
	return ObjInteger
}

func (i *Integer) Inspect() string {
	return strconv.FormatInt(i.Value, 64)
}
