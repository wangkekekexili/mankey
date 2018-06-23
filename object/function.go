package object

import (
	"strings"

	"github.com/wangkekekexili/mankey/ast"
)

const ObjFunction = "FUNCTION"

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType {
	return ObjFunction
}

func (f *Function) String() string {
	var paraStrs []string
	for _, para := range f.Parameters {
		paraStrs = append(paraStrs, para.String())
	}

	var b strings.Builder
	b.WriteString("func(")
	b.WriteString(strings.Join(paraStrs, ","))
	b.WriteString(")")
	b.WriteString(f.Body.String())
	return b.String()
}
