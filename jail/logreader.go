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

func (l *logReader) readLine() error {
	line, err := l.reader.ReadString('\n')
	if err != nil {
		return err
	}

	if line != "" {
		l.lines <- line
	}
	return nil
}

func (l *logReader) run() {
	for {
		err := l.readLine()
		if err != nil {
			l.errors <- err
			break
		}
	}
}
