package analyzers

import (
	"compiler/analyzers/reader"
	"compiler/analyzers/token"
	"fmt"
)

const (
	lexicalError      = "lexical error. %v"
	commentLineSymbol = '#'
	newLineSymbol     = '\n'
)

type Lexer struct {
	file      *reader.File
	debugMode bool

	currentChar     byte
	currentSpelling string
}

func NewLexer(src string, debugMode bool) *Lexer {
	file := reader.NewFile(src)
	lexer := &Lexer{
		file:            file,
		debugMode:       debugMode,
		currentChar:     byte(0),
		currentSpelling: "",
	}
	lexer.next()
	return lexer
}

func (l *Lexer) GetToken() (token.Token, error) {
	l.skipWhiteSpace()
	l.skipComment()
	defer l.next()

	switch l.currentChar {
	case '+':
		return token.NewToken(token.OP_PLUS, l.file.CurrentPosition), nil
	case '-':
		return token.NewToken(token.OP_MINUS, l.file.CurrentPosition), nil
	case '*':
		return token.NewToken(token.OP_MULTI, l.file.CurrentPosition), nil
	case '/':
		return token.NewToken(token.OP_DIVIDE, l.file.CurrentPosition), nil
	case '(':
		return token.NewToken(token.L_BRACKET, l.file.CurrentPosition), nil
	case ')':
		return token.NewToken(token.R_BRACKET, l.file.CurrentPosition), nil
	case '{':
		return token.NewToken(token.L_BRACE, l.file.CurrentPosition), nil
	case '}':
		return token.NewToken(token.R_BRACE, l.file.CurrentPosition), nil
	case reader.EOL:
		return token.NewToken(token.NEWLINE, l.file.CurrentPosition), nil
	case reader.EOF:
		return token.NewToken(token.EOF, l.file.CurrentPosition), nil
	case '=':
		//Check if next character is an eq symbol
		if l.peek() == '=' {
			l.next()
			return token.NewToken(token.OP_EQ, l.file.CurrentPosition), nil
		} else {
			return token.NewToken(token.ASSIGN, l.file.CurrentPosition), nil
		}
	case '>':
		if l.peek() == '=' {
			l.next()
			return token.NewToken(token.OP_GTE, l.file.CurrentPosition), nil
		} else {
			return token.NewToken(token.OP_GT, l.file.CurrentPosition), nil
		}
	case '<':
		if l.peek() == '=' {
			l.next()
			return token.NewToken(token.OP_LTE, l.file.CurrentPosition), nil
		} else {
			return token.NewToken(token.OP_LT, l.file.CurrentPosition), nil
		}
	case '!':
		if l.peek() == '=' {
			l.next()
			return token.NewToken(token.OP_NOTEQ, l.file.CurrentPosition), nil
		} else {
			return token.Token{}, l.abort("Expected !=, got !" + string(l.peek()))
		}
	case '"':
		var str []byte
		pos := l.file.CurrentPosition
		if l.peek() != '"' {
			l.next()
			for l.currentChar != '"' {
				str = append(str, l.currentChar)
				l.next()
			}
		}
		return token.NewTokenString(token.STRING, string(str), pos), nil
	default:
		if isDigit(l.currentChar) { // Check for numbers
			var number []byte
			pos := l.file.CurrentPosition
			number = append(number, l.currentChar)
			for isDigit(l.peek()) {
				l.next()
				number = append(number, l.currentChar)
			}
			if l.peek() == '.' {
				l.next()
				number = append(number, l.currentChar)
				if !isDigit(l.peek()) {
					return token.Token{}, l.abort("Illegal character in number: " + string(number) + string(l.peek()))
				}
			}
			for isDigit(l.peek()) {
				l.next()
				number = append(number, l.currentChar)
			}
			return token.NewTokenString(token.NUMBER, string(number), pos), nil
		} else if isAlpha(l.currentChar) { // Check for Identifiers
			var id []byte
			pos := l.file.CurrentPosition
			id = append(id, l.currentChar)
			for isAlphaNumeric(l.peek()) {
				l.next()
				id = append(id, l.currentChar)
			}
			return token.NewTokenString(token.IDENTIFIER, string(id), pos), nil
		} else {
			return token.Token{}, l.abort("Unknown token: " + string(l.currentChar))
		}
	}
}

func (l *Lexer) next() {
	l.currentChar = l.file.Read()
}

func (l *Lexer) peek() byte {
	return l.file.Peek()
}

func (l *Lexer) abort(message string) error {
	return fmt.Errorf(lexicalError, message)
}

func (l *Lexer) skipWhiteSpace() {
	for isWhiteSpace(l.currentChar) {
		l.next()
	}
}

func isWhiteSpace(b byte) bool {
	for _, v := range []byte{' ', '\t', '\r'} {
		if v == b {
			return true
		}
	}
	return false
}

func (l *Lexer) skipComment() {
	if l.currentChar == commentLineSymbol {
		for l.currentChar != newLineSymbol {
			l.next()
		}
	}
}

func isDigit(b byte) bool {
	zero := byte('0')
	nine := byte('9')

	return b >= zero && b <= nine
}

func isAlpha(b byte) bool {
	a := byte('a')
	z := byte('z')
	A := byte('A')
	Z := byte('Z')
	return b >= a && b <= z || b >= A && b <= Z
}

func isAlphaNumeric(b byte) bool {
	return isDigit(b) || isAlpha(b)
}
