package parser

import (
	"fmt"

	"github.com/wangkekekexili/mankey/token"
)

type unexpectedToken struct {
	exp string
	t   *token.Token
}

func (t unexpectedToken) Error() string {
	return fmt.Sprintf("expect %v; got token %v", t.exp, t.t)
}
