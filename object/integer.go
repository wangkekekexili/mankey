package object

import (
	"fmt"
	"strconv"
)

const ObjInteger = "INTEGER"

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType {
	return ObjInteger
}

func (i *Integer) String() string {
	return strconv.FormatInt(i.Value, 10)
}

func (i *Integer) HashKey() HashKey {
	return HashKey(fmt.Sprintf("%v_%v", ObjInteger, i.Value))
}
