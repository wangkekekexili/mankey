package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/wangkekekexili/mankey/lexer"
	"github.com/wangkekekexili/mankey/parser"
)

func Do(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		p, err := parser.New(lexer.New(scanner.Text())).ParseProgram()
		if err != nil {
			fmt.Fprintln(w, err)
			continue
		}
		fmt.Fprintln(w, p)
	}
}
