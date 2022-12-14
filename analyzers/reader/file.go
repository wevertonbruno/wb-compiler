package reader

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type File struct {
	src             string
	reader          *bufio.Reader
	currentPosition Position
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
		currentPosition: Position{
			Line:   1,
			Column: 0,
		}}
}

func (f *File) Read() byte {
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

func (f *File) calculatePosition(b byte) {
	if b == EOL {
		f.currentPosition.Line++
		f.currentPosition.Column = 0
	} else {
		f.currentPosition.Column++
	}
}

func (f *File) Peek() byte {
	b, err := f.reader.Peek(1)
	if err != nil {
		if err == io.EOF {
			return EOF
		}
		panic(fmt.Sprintf(cannotReadError, err))
	}
	return b[0]
}

func (f *File) CurrentPosition() Position {
	return f.currentPosition
}
