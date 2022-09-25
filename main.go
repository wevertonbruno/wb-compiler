package main

import (
	"compiler/analyzers"
	"compiler/analyzers/token"
	"fmt"
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
