package filer

import (
	"bufio"
	"os"
	"sync"
)

type Reader struct {
	file      *os.File
	inputChan chan string
	wg        *sync.WaitGroup
}

func NewReader(inputFilePath string, inputChan chan string, wg *sync.WaitGroup) (*Reader, error) {
	inputFile, err := os.Open(inputFilePath)
	if err != nil {
		return nil, err
	}

	return &Reader{
		file:      inputFile,
		inputChan: inputChan,
		wg:        wg,
	}, nil
}

func (r *Reader) Close() {
	r.file.Close()
}

func (r *Reader) Read() ([]string, error) {
	lines := make([]string, 0, 10)
	scanner := bufio.NewScanner(r.file)

	for scanner.Scan() {
		if scanner.Text() == "" {
			continue
		}

		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
