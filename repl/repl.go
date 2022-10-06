package repl

import (
	"bufio"
	"fmt"
	"github.com/wevertonbruno/wb-compiler/analyzers/lexer"
	"github.com/wevertonbruno/wb-compiler/analyzers/parser"
	"github.com/wevertonbruno/wb-compiler/analyzers/reader"
	"io"
)

const (
	PROMPT  = "-> "
	WELCOME = "Welcome to wb-lang!"
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	fmt.Println(WELCOME)
	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := reader.NewInput(scanner.Text())
		lex := lexer.NewLexer(line)
		_parser := parser.NewParser(lex)
		program := _parser.Parse()
		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}
