package jail

import (
	"bufio"
	"fmt"
	"os"
)

type logReader struct {
	filename string
	file     *os.File
	reader   *bufio.Reader
	lines    chan string
	errors   chan error
}

func newLogReader(filename string) *logReader {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	r := bufio.NewReader(f)

	return &logReader{
		filename: filename,
		file:     f,
		reader:   r,
		lines:    make(chan string),
		errors:   make(chan error),
	}
}

func (l *logReader) readLine() {
	line, err := l.reader.ReadString('\n')
	if err != nil {
		go func() {
			l.errors <- err
		}()
	}

	if line != "" {
		go func() {
			l.lines <- line
		}()
	}
}
func (l *logReader) reset() {
	f, err := os.Open(l.filename)
	if err != nil {
		fmt.Println(err)
	}
	r := bufio.NewReader(f)
	l.reader.Reset(r)
}
