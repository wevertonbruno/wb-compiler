package main

import (
	"github.com/wevertonbruno/wb-compiler/repl"
	"os"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
