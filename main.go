package main

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"tbc/filer"
)

func main() {
	wg := &sync.WaitGroup{}

	// канал для записи строк из файла
	inputChan := make(chan string)

	// канал для записи суммы строк из файла
	outputChan := make(chan int)

	// создаём читателя
	fileReader, err := filer.NewReader("input.txt", inputChan, wg)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer fileReader.Close()

	// создаём писателя
	outputFile, err := filer.NewWriter("output.txt", outputChan)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer outputFile.Close()

	wg.Add(1)

	// горутина для прослушки канала с суммами строк
	go func() {
		defer wg.Done()

		// читаем из канала, пока он открыт
		for output := range outputChan {
			fmt.Println(fmt.Sprintf("Got output: %d", output))

			// запись суммы строк в файл
			if err := outputFile.Write(output); err != nil {
				return
			}
		}
	}()

	wg.Add(1)
	timeToSleep := 1

	// горутина для прослушки канала со строками из файла
	go func() {
		defer wg.Done()

		// читаем из канала, пока он открыт
		for input := range inputChan {
			wg.Add(1)

			// горутина для суммирования строк из файла
			go sumLine(input, outputChan, wg, timeToSleep)

			// снотворное для сохранения последовательности строк
			timeToSleep++
		}
	}()

	// вычитка строк из файла и запись из в слайс
	lines, err := fileReader.Read()
	if err != nil {
		fmt.Println(err)
		return
	}

	// запись трок из в слайса в канал
	for _, line := range lines {
		inputChan <- line
		fmt.Println(fmt.Sprintf("Wrote line: %s", line))
	}

	close(inputChan)

	go func() {
		wg.Wait()
		close(outputChan)
	}()

	time.Sleep(time.Duration(timeToSleep) * time.Second)

	// пуляем буферные данные файла в настоящий файл
	if err := outputFile.Flush(); err != nil {
		return
	}

	fmt.Println("done")

	return
}

func sumLine(lineNums string, outputChan chan int, wg *sync.WaitGroup, timeToSleep int) {
	defer wg.Done()

	fmt.Println(fmt.Sprintf("Received input line: %s", lineNums))

	time.Sleep(time.Duration(timeToSleep) * time.Second)

	lineNumsSlice := strings.Split(lineNums, " ")
	total := 0

	for _, lineNum := range lineNumsSlice {
		lineNumInt, _ := strconv.Atoi(lineNum)

		total += lineNumInt
	}

	outputChan <- total
}
