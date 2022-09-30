package reader

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

type Reader interface {
	Read() byte
	Peek() byte
	CurrentPosition() Position
}
