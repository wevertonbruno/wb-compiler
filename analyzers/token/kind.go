package token

type Kind byte

func (k Kind) Name() string {
	return spellMapping[k]
}

const (
	IDENTIFIER Kind = iota
	NEWLINE
	INTLIT
	DECIMALLIT
	CHARLIT
	STRINGLIT
	BOOLEANLIT

	IF
	ELSE
	THEN
	FUNCTION
	WHILE
	VAR
	PRINT
	RETURN
	INTEGER
	DECIMAL
	STRING
	CHAR
	BOOLEAN
	TRUE
	FALSE

	OP_PLUS
	OP_MULTI
	OP_MINUS
	OP_DIVIDE
	OP_EQ
	OP_NOTEQ
	OP_LT
	OP_LTE
	OP_GT
	OP_GTE
	NOT

	ASSIGN
	SEMICOLON
	COLON
	L_BRACKET
	R_BRACKET
	L_BRACE
	R_BRACE

	ERROR
	EOF
)

var (
	spellMapping = map[Kind]string{
		IDENTIFIER: "<identifier>",
		NEWLINE:    "<new line>",
		INTLIT:     "<integer>",
		DECIMALLIT: "<decimal>",
		CHARLIT:    "<char>",
		STRINGLIT:  "<string>",
		BOOLEANLIT: "<boolean>",

		IF:       "if",
		ELSE:     "else",
		THEN:     "then",
		FUNCTION: "func",
		WHILE:    "while",
		VAR:      "var",
		PRINT:    "print",
		RETURN:   "return",
		INTEGER:  "Integer",
		DECIMAL:  "Decimal",
		STRING:   "String",
		CHAR:     "Char",
		BOOLEAN:  "Boolean",
		TRUE:     "true",
		FALSE:    "false",

		OP_PLUS:   "+",
		OP_MINUS:  "-",
		OP_MULTI:  "*",
		OP_DIVIDE: "/",
		OP_EQ:     "==",
		OP_NOTEQ:  "!=",
		OP_LT:     "<",
		OP_LTE:    "<=",
		OP_GT:     ">",
		OP_GTE:    ">=",
		NOT:       "!",

		ASSIGN:    "=",
		SEMICOLON: ";",
		COLON:     ":",
		L_BRACKET: "(",
		R_BRACKET: ")",
		L_BRACE:   "{",
		R_BRACE:   "}",
		ERROR:     "<error>",
		EOF:       "<eof>",
	}
)
