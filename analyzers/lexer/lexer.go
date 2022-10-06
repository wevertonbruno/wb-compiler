package lexer

import (
	"fmt"
	"github.com/wevertonbruno/wb-compiler/analyzers/reader"
	"github.com/wevertonbruno/wb-compiler/analyzers/token"
)

const (
	lexicalError      = "lexical error. %v"
	commentLineSymbol = '#'
	newLineSymbol     = '\n'
)

type Lexer struct {
	reader reader.Reader

	currentChar     byte
	currentSpelling string
}

func NewLexer(reader reader.Reader) *Lexer {
	lexer := &Lexer{
		reader:          reader,
		currentChar:     byte(0),
		currentSpelling: "",
	}
	lexer.next()
	return lexer
}

func (l *Lexer) GetToken() token.Token {
	l.skipWhiteSpace()
	l.skipComment()
	defer l.next()

	switch l.currentChar {
	case '+':
		return token.NewToken(token.OP_PLUS, l.reader.CurrentPosition())
	case '-':
		return token.NewToken(token.OP_MINUS, l.reader.CurrentPosition())
	case '*':
		return token.NewToken(token.OP_MULTI, l.reader.CurrentPosition())
	case '/':
		return token.NewToken(token.OP_DIVIDE, l.reader.CurrentPosition())
	case '(':
		return token.NewToken(token.L_BRACKET, l.reader.CurrentPosition())
	case ')':
		return token.NewToken(token.R_BRACKET, l.reader.CurrentPosition())
	case '{':
		return token.NewToken(token.L_BRACE, l.reader.CurrentPosition())
	case '}':
		return token.NewToken(token.R_BRACE, l.reader.CurrentPosition())
	case ':':
		return token.NewToken(token.COLON, l.reader.CurrentPosition())
	case ';':
		return token.NewToken(token.SEMICOLON, l.reader.CurrentPosition())
	case ',':
		return token.NewToken(token.COMMA, l.reader.CurrentPosition())
	case reader.EOL:
		return token.NewToken(token.NEWLINE, l.reader.CurrentPosition())
	case reader.EOF:
		return token.NewToken(token.EOF, l.reader.CurrentPosition())
	case '=':
		//Check if next character is an eq symbol
		if l.peek() == '=' {
			l.next()
			return token.NewToken(token.OP_EQ, l.reader.CurrentPosition())
		} else {
			return token.NewToken(token.ASSIGN, l.reader.CurrentPosition())
		}
	case '>':
		if l.peek() == '=' {
			l.next()
			return token.NewToken(token.OP_GTE, l.reader.CurrentPosition())
		} else {
			return token.NewToken(token.OP_GT, l.reader.CurrentPosition())
		}
	case '<':
		if l.peek() == '=' {
			l.next()
			return token.NewToken(token.OP_LTE, l.reader.CurrentPosition())
		} else {
			return token.NewToken(token.OP_LT, l.reader.CurrentPosition())
		}
	case '!':
		if l.peek() == '=' {
			l.next()
			return token.NewToken(token.OP_NOTEQ, l.reader.CurrentPosition())
		} else {
			return token.NewToken(token.NOT, l.reader.CurrentPosition())
		}
	case '"':
		var str []byte
		pos := l.reader.CurrentPosition()
		if l.peek() != '"' {
			l.next()
			for l.currentChar != '"' {
				str = append(str, l.currentChar)
				l.next()
			}
		}
		return token.NewTokenString(token.STRINGLIT, string(str), pos)
	default:
		if isDigit(l.currentChar) { // Check for numbers
			var number []byte
			pos := l.reader.CurrentPosition()
			number = append(number, l.currentChar)
			for isDigit(l.peek()) {
				l.next()
				number = append(number, l.currentChar)
			}
			if l.peek() == '.' {
				l.next()
				number = append(number, l.currentChar)
				if !isDigit(l.peek()) {
					l.abort("Illegal character in number: " + string(number) + string(l.peek()))
				}
				for isDigit(l.peek()) {
					l.next()
					number = append(number, l.currentChar)
				}
				return token.NewTokenString(token.DECIMALLIT, string(number), pos)
			} else {
				return token.NewTokenString(token.INTLIT, string(number), pos)
			}
		} else if isAlpha(l.currentChar) { // Check for Identifiers
			var id []byte
			pos := l.reader.CurrentPosition()
			id = append(id, l.currentChar)
			for isAlphaNumeric(l.peek()) {
				l.next()
				id = append(id, l.currentChar)
			}
			return token.NewTokenString(token.IDENTIFIER, string(id), pos)
		} else {
			l.abort("Unknown token: " + string(l.currentChar))
		}
	}
	return token.Token{}
}

func (l *Lexer) next() {
	l.currentChar = l.reader.Read()
}

func (l *Lexer) peek() byte {
	return l.reader.Peek()
}

func (l *Lexer) abort(message string) error {
	panic(fmt.Errorf(lexicalError, message))
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
