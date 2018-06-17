package parser

import (
	"fmt"

	"github.com/wangkekekexili/mankey/token"
)

type errUnexpectedToken struct {
	exp string
	t   *token.Token
}

func (e errUnexpectedToken) Error() string {
	return fmt.Sprintf("expect %v; got token %v", e.exp, e.t)
}

type errNoPrefixParseFunction struct {
	t *token.Token
}

func (e errNoPrefixParseFunction) Error() string {
	return fmt.Sprintf("no prefix parse function for %v", e.t)
}

type errNoInfixParseFunction struct {
	t *token.Token
}

func (e errNoInfixParseFunction) Error() string {
	return fmt.Sprintf("no infix parse function for %v", e.t)
}
