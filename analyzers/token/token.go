package token

import (
	"fmt"
	"github.com/wevertonbruno/wb-compiler/analyzers/reader"
)

type Token struct {
	Kind     Kind            `json:"kind"`
	Spelling string          `json:"spelling"`
	Position reader.Position `json:"position"`
}

func NewToken(kind Kind, position reader.Position) Token {
	return Token{
		Kind:     kind,
		Spelling: spellMapping[kind],
		Position: position,
	}
}

func NewTokenString(kind Kind, spell string, position reader.Position) Token {
	if kind == IDENTIFIER {
		if val, exists := reservedKeywords[spell]; exists {
			kind = val
		}
	}
	return Token{
		Kind:     kind,
		Spelling: spell,
		Position: position,
	}
}

func (t Token) String() string {
	return fmt.Sprintf("{Kind: %v, Spelling: %v, Position: %v}", spellMapping[t.Kind], t.Spelling, t.Position)
}

var (
	spellMapping = map[Kind]string{
		IDENTIFIER: "<identifier>",
		NEWLINE:    "<new line>",
		NUMBER:     "<number>",
		CHAR:       "<char>",
		STRING:     "<string>",

		IF:       "if",
		ELSE:     "else",
		THEN:     "then",
		FUNCTION: "func",
		WHILE:    "while",
		VAR:      "var",

		OP_PLUS:   "<plus>",
		OP_MINUS:  "<minus>",
		OP_MULTI:  "<mult>",
		OP_DIVIDE: "<div>",
		OP_EQ:     "<equals>",
		OP_NOTEQ:  "<not equals>",
		OP_LT:     "<less than>",
		OP_LTE:    "<less than or equals>",
		OP_GT:     "<grather than>",
		OP_GTE:    "<grather than or equals>",

		ASSIGN:    "=",
		SEMICOLON: ";",
		L_BRACKET: "(",
		R_BRACKET: ")",
		L_BRACE:   "{",
		R_BRACE:   "}",
		ERROR:     "<error>",
		EOF:       "<eof>",
	}

	reservedKeywords = map[string]Kind{
		spellMapping[IF]:       IF,
		spellMapping[ELSE]:     ELSE,
		spellMapping[THEN]:     THEN,
		spellMapping[FUNCTION]: FUNCTION,
		spellMapping[WHILE]:    WHILE,
		spellMapping[VAR]:      VAR,
	}
)
