package main

import (
	"os"

	"github.com/wangkekekexili/mankey/repl"
)

func main() {
	repl.Do(os.Stdin, os.Stdout)
}
