package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"strings"
	"sync"
	"time"
)

const iso8601 = "2006-01-02T15:04:05.00Z"

func main() {

	s := time.Now()
	args := os.Args[1:]
	if len(args) != 6 { // for format  LogExtractor.exe -f "From Time" -t "To Time" -i "Log file directory location"
		fmt.Println("Please give proper command line arguments")
		return
	}
	startTimeArg := args[1]
	finishTimeArg := args[3]
	fileDirectory := args[5]
	rangeStart, err := time.Parse(iso8601, startTimeArg)
	if err != nil {
		fmt.Println("Could not able to parse the start time", startTimeArg)
		return
	}
	rangeFinish, err := time.Parse(iso8601, finishTimeArg)
	if err != nil {
		fmt.Println("Could not able to parse the finish time", finishTimeArg)
		return
	}
	items, _ := ioutil.ReadDir(fileDirectory)
	var wg sync.WaitGroup
	wg.Add(len(items))
	for _, item := range items {
	//	fmt.Println(item.Name())
		go func(filename string, rangeStart time.Time, rangeFinish time.Time) {
			ProcessFile(filename, rangeStart, rangeFinish)
			wg.Done()
		}(fileDirectory+"/"+item.Name(), rangeStart, rangeFinish)
	}
	wg.Wait()
	fmt.Println("\nTime taken - ", time.Since(s))
}

func ProcessFile(fileName string, rangeStart time.Time, rangeFinish time.Time) {
	file, err := os.Open(fileName)
	defer file.Close()
	if err != nil {
		fmt.Println("couldn't read the file", err)
		return
	}

	firstLineSize := getFirstLineSize(file)
	firstLine := make([]byte, firstLineSize)
	_, err = file.ReadAt(firstLine, 0)
	if err != nil {
		fmt.Println("Could not able to read first line with offset", 0, "and firstLine size", firstLineSize)
		return
	}

	firstLogSlice := strings.SplitN(string(firstLine), ",", 2)
	firstLogCreationTimeString := firstLogSlice[0]
	firstLogCreationTime, err := time.Parse(iso8601, firstLogCreationTimeString)
	if err != nil {
		fmt.Println("can not able to parse time : ", err)
	}
	if firstLogCreationTime.After(rangeFinish) {
		return
	}

	Process(file, rangeStart, rangeFinish)
}

func getFirstLineSize(file *os.File) int {
	offset := int64(0)
	firstLineSize := 0
	for {
		b := make([]byte, 1)
		n, err := file.ReadAt(b, offset)
		if err != nil {
			fmt.Println("Error reading file ", err)
			break
		}
		char := string(b[0])
		if char == "\n" {
			break
		}
		offset++
		firstLineSize += n
	}
	return firstLineSize
}

func Process(f *os.File, start time.Time, end time.Time) error {

	linesPool := sync.Pool{New: func() interface{} {
		lines := make([]byte, 250*1024)
		return lines
	}}

	stringPool := sync.Pool{New: func() interface{} {
		lines := ""
		return lines
	}}

	r := bufio.NewReader(f)

	var wg sync.WaitGroup

	for {
		buf := linesPool.Get().([]byte)

		n, err := r.Read(buf)
		buf = buf[:n]

		if n == 0 {
			if err != nil {
				fmt.Println(err)
				break
			}
			if err == io.EOF {
				break
			}
			return err
		}

		nextUntillNewline, err := r.ReadBytes('\n')

		if err != io.EOF {
			buf = append(buf, nextUntillNewline...)
		}

		wg.Add(1)
		go func() {
			ProcessChunk(buf, &linesPool, &stringPool, start, end)
			wg.Done()
		}()

	}

	wg.Wait()
	return nil
}

func ProcessChunk(chunk []byte, linesPool *sync.Pool, stringPool *sync.Pool, start time.Time, end time.Time) {

	var wg sync.WaitGroup

	logs := stringPool.Get().(string)
	logs = string(chunk)

	linesPool.Put(chunk)

	logsSlice := strings.Split(logs, "\n")

	stringPool.Put(logs)

	chunkSize := 300
	n := len(logsSlice)
	noOfThread := n / chunkSize

	if n%chunkSize != 0 {
		noOfThread++
	}

	for i := 0; i < (noOfThread); i++ {

		wg.Add(1)
		go func(s int, e int) {
			defer wg.Done()
			for i := s; i < e; i++ {
				text := logsSlice[i]
				if len(text) == 0 {
					continue
				}
				logSlice := strings.SplitN(text, ",", 2)
				logCreationTimeString := logSlice[0]

				logCreationTime, err := time.Parse(iso8601, logCreationTimeString)
				if err != nil {
					fmt.Printf("\n Could not able to parse the time :%s for log : %v", logCreationTimeString, text)
					return
				}

				if logCreationTime.After(start) && logCreationTime.Before(end) {
					fmt.Println(text)
				}
			}

		}(i*chunkSize, int(math.Min(float64((i+1)*chunkSize), float64(len(logsSlice)))))
	}

	wg.Wait()
	logsSlice = nil
}
