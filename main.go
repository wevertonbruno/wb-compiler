package main

import (
	"github.com/wevertonbruno/wb-compiler/analyzers/lexer"
	"github.com/wevertonbruno/wb-compiler/analyzers/parser"
	"github.com/wevertonbruno/wb-compiler/analyzers/reader"
)

func main() {
	_reader := reader.NewFile("test_code.wb")
	_lexer := lexer.NewLexer(_reader)
	_parser := parser.NewParser(_lexer)
	_parser.Parse()
}
