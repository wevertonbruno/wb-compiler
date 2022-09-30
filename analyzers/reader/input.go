package reader

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Input struct {
	reader          *bufio.Reader
	currentPosition Position
	lastRead        byte
}

func NewInput(input string) *Input {
	stringReader := strings.NewReader(input)
	return &Input{
		reader: bufio.NewReader(stringReader),
		currentPosition: Position{
			Line:   1,
			Column: 0,
		}}
}

func (f *Input) Read() byte {
	b, err := f.reader.ReadByte()
	if err != nil {
		if err == io.EOF {
			return EOF
		}
		panic(fmt.Sprintf(cannotReadError, err))
	}
	f.calculatePosition(b)
	return b
}

func (f *Input) calculatePosition(b byte) {
	if b == EOL {
		f.currentPosition.Line++
		f.currentPosition.Column = 0
	} else {
		f.currentPosition.Column++
	}
}

func (f *Input) Peek() byte {
	b, err := f.reader.Peek(1)
	if err != nil {
		if err == io.EOF {
			return EOF
		}
		panic(fmt.Sprintf(cannotReadError, err))
	}
	return b[0]
}

func (f *Input) CurrentPosition() Position {
	return f.currentPosition
}
