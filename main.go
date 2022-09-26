package main

import (
	"fmt"
	"github.com/wevertonbruno/wb-compiler/analyzers"
	"github.com/wevertonbruno/wb-compiler/analyzers/token"
)

func main() {
	lexer := analyzers.NewLexer("test_code.wb", true)
	t, err := lexer.GetToken()
	if err != nil {
		panic(err)
	}

	for t.Kind != token.EOF {
		fmt.Println(t)
		t, err = lexer.GetToken()
		if err != nil {
			panic(err)
		}
	}
}
