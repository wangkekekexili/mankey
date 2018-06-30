package object

import (
	"fmt"
	"strings"
)

const ObjHash = "HASH"

type HashPair struct {
	K, V Object
}

type Hash struct {
	Hash map[HashKey]*HashPair
}

func (h *Hash) Type() ObjectType {
	return ObjHash
}

func (h *Hash) String() string {
	var strs []string
	for _, v := range h.Hash {
		strs = append(strs, fmt.Sprintf("%v: %v", v.K, v.V))
	}
	return "{" + strings.Join(strs, ",") + "}"
}
