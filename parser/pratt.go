package parser

import (
	"github.com/wangkekekexili/mankey/ast"
)

type (
	prefixParseFn func() (ast.Expression, error)
	infixParseFn  func(ast.Expression) (ast.Expression, error)
)
