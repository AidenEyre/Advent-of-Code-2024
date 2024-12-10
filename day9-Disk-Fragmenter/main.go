package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	start := time.Now() // get time

	filesystem, err := loadDataIntoSlice("puzzle-data.txt")
	if err != nil {
		panic(err)
	}
	fsCompacted := compatFilesystem(filesystem)
	checksum := calcChecksum(fsCompacted)

	fmt.Printf("Filesystem checksum is: '%d'\n", checksum)

	elapsed := time.Since(start) // Calculate elapsed time
	fmt.Printf("Execution time: %d ms\n", elapsed.Milliseconds())
}

func loadDataIntoSlice(filename string) ([]string, error) {
	var dataSlice = []string{}

	file, err := os.Open(filename)
	if err != nil {
		return []string{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		dataSlice = strings.Split(line, "")
	}

	return dataSlice, scanner.Err()
}

func compatFilesystem(fs []string) []int {
	fsCompacted := []int{}
	fileId := 0
	freeSpace := []int{}

	for fsIndex := range fs {
		current, _ := strconv.Atoi(fs[fsIndex])
		if isEven(fsIndex) {
			for block := 0; block < current; block++ {
				fsCompacted = append(fsCompacted, fileId)
			}
			fileId++
		} else {
			for block := 0; block < current; block++ {
				freeSpace = append(freeSpace, len(fsCompacted))
				fsCompacted = append(fsCompacted, -1)
			}
		}
	}

	freeSpaceIndex := len(freeSpace)
	for fs := 0; fs < freeSpaceIndex; fs++ {
		lastIndex := len(fsCompacted) - 1

		for fsCompacted[lastIndex] == -1 {
			lastIndex--
			freeSpaceIndex--
		}

		fsCompacted[freeSpace[fs]] = fsCompacted[lastIndex]
		fsCompacted = fsCompacted[:lastIndex]
	}

	return fsCompacted
}

func isEven(number int) bool {
	return number&1 == 0
}

func calcChecksum(input []int) int {
	checksum := 0

	for position := range input {
		checksum += input[position] * position
	}

	return checksum
}
