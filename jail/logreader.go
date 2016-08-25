package jail

import (
	"bufio"
	"io"
	"os"
)

type LogReader struct {
	filename string
	file     *os.File
	offset   int64
	reader   *bufio.Reader
	Lines    chan string
	Errors   chan error
}

func NewLogReader(filename string, offset int64) *LogReader {
	return &LogReader{
		filename: filename,
		offset:   offset,
		Lines:    make(chan string),
		Errors:   make(chan error),
	}
}

func (l *LogReader) readLine() error {
	line, err := l.reader.ReadString('\n')
	if err != nil {
		return err
	}

	if line != "" {
		l.Lines <- line
	}
	return nil
}

func (l *LogReader) Run() {
	file, err := os.Open(l.filename)
	defer file.Close()

	if err != nil && !os.IsNotExist(err) {
		l.Errors <- err
	}

	if err == nil {
		l.file = file
		l.reader = bufio.NewReader(file)
	}

	for {
		er = l.readLine()
		if er != nil {
			l.Errors <- er
			break
		}
	}
}

func (l *LogReadeer) Reset() {
	l.offset = 0
}
