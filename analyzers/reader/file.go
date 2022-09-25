package reader

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

const (
	cannotOpenFileError = "Can not open the file: %v"
	cannotReadError     = "Can not read the file: %v"
	EOL                 = '\n'
	EOF                 = '\u0000'
)

type Position struct {
	Line   int
	Column int
}

type File struct {
	src             string
	reader          *bufio.Reader
	CurrentPosition Position
	lastRead        byte
}

func NewFile(src string) *File {
	file, err := os.Open(src)
	if err != nil {
		panic(fmt.Sprintf(cannotOpenFileError, err))
	}

	return &File{
		src:    src,
		reader: bufio.NewReader(file),
		CurrentPosition: Position{
			Line:   1,
			Column: 0,
		}}
}

func (s *File) Read() byte {
	b, err := s.reader.ReadByte()
	if err != nil {
		if err == io.EOF {
			return EOF
		}
		panic(fmt.Sprintf(cannotReadError, err))
	}
	s.calculatePosition(b)
	return b
}

func (s *File) calculatePosition(b byte) {
	if b == EOL {
		s.CurrentPosition.Line++
		s.CurrentPosition.Column = 0
	} else {
		s.CurrentPosition.Column++
	}
}

func (s *File) Peek() byte {
	b, err := s.reader.Peek(1)
	if err != nil {
		if err == io.EOF {
			return EOF
		}
		panic(fmt.Sprintf(cannotReadError, err))
	}
	return b[0]
}
