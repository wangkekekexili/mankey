package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/wangkekekexili/mankey/lexer"
	"github.com/wangkekekexili/mankey/parser"
)

func Do(r io.Reader, w io.Writer) {

	fmt.Fprintf(w, ">> ")
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "exit" {
			return
		}
		p, err := parser.New(lexer.New(text)).ParseProgram()
		if err != nil {
			fmt.Fprintln(w, err)
			fmt.Fprintf(w, ">> ")
			continue
		}
		fmt.Fprintln(w, p)
		fmt.Fprintf(w, ">> ")
	}
}
