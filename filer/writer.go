package filer

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Writer struct {
	outputFile  *os.File
	outputChan  chan int
	bufioWriter *bufio.Writer
}

func NewWriter(outputFilePath string, outputChan chan int) (*Writer, error) {
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Println(err)
	}

	writer := bufio.NewWriter(outputFile)

	return &Writer{
		outputFile:  outputFile,
		outputChan:  outputChan,
		bufioWriter: writer,
	}, nil
}

func (w *Writer) Close() {
	w.outputFile.Close()
}

func (w *Writer) Write(output int) error {
	item := []byte(fmt.Sprintf("%s\n", strconv.Itoa(output)))

	if _, err := w.bufioWriter.Write(item); err != nil {
		return err
	}

	return nil
}

func (w *Writer) Flush() error {
	if err := w.bufioWriter.Flush(); err != nil {
		return err
	}

	return nil
}
