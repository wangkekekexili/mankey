package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/wangkekekexili/mankey/lexer"
	"github.com/wangkekekexili/mankey/token"
)

func Do(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		x := lexer.New(scanner.Text())
		for {
			t := x.NextToken()
			if t.Equals(token.New(token.EOF, "")) {
				break
			}
			fmt.Fprintln(w, t)
		}
	}
}
