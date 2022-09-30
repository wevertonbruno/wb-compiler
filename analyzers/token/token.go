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

func (t Token) Match(k Kind) bool {
	return t.Kind == k
}

func (t Token) String() string {
	return fmt.Sprintf("{Kind: %v, Spelling: %v, Position: %v}", spellMapping[t.Kind], t.Spelling, t.Position)
}

var (
	reservedKeywords = map[string]Kind{
		spellMapping[IF]:       IF,
		spellMapping[ELSE]:     ELSE,
		spellMapping[THEN]:     THEN,
		spellMapping[FUNCTION]: FUNCTION,
		spellMapping[WHILE]:    WHILE,
		spellMapping[VAR]:      VAR,
		spellMapping[PRINT]:    PRINT,
		spellMapping[RETURN]:   RETURN,

		spellMapping[INTEGER]: INTEGER,
		spellMapping[DECIMAL]: DECIMAL,
		spellMapping[STRING]:  STRING,
		spellMapping[CHAR]:    CHAR,
		spellMapping[BOOLEAN]: BOOLEAN,
		spellMapping[TRUE]:    TRUE,
		spellMapping[FALSE]:   FALSE,
	}
)
