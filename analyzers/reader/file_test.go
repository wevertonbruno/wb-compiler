package reader

import (
	"os"
	"testing"
)

const (
	testFile = "reader_test.txt"
)

func TestSourceFile_Read(t *testing.T) {
	file := createTestFile("Hello", t)
	defer file.Close()
	source := NewFile(testFile)
	text := []byte("Hello")
	i := 0
	for b := source.Read(); b != EOF; b = source.Read() {
		if text[i] != b {
			t.Errorf("%v is not equals to %v", text[i], b)
		}
		i++
	}
	deleteTestFile()
}

func TestSourceFile_Peek(t *testing.T) {
	file := createTestFile("Hello", t)
	defer file.Close()
	source := NewFile(testFile)
	text := []byte("Hello")
	i := 0
	for b := source.Peek(); b != EOF; b = source.Peek() {
		source.Read()
		if text[i] != b {
			t.Errorf("%v is not equals to %v", text[i], b)
		}
		i++
	}
	deleteTestFile()
}

func createTestFile(text string, t *testing.T) *os.File {
	file, err := os.Create(testFile)
	if err != nil {
		t.Errorf("Cannot create test file: %v", err)
	}

	_, err = file.WriteString(text)
	if err != nil {
		t.Errorf("Cannot write in test file: %v", err)
	}
	return file
}

func deleteTestFile() {
	err := os.Remove(testFile)
	if err != nil {
		println("Cannot delete file: ", err)
	}
}
