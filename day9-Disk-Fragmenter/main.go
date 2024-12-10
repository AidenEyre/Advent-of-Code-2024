package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

type filesystem struct {
	fs        []int
	freeSpace [][2]int
	usedSpace [][2]int
}

func main() {
	start := time.Now() // get time

	filesystem := filesystem{}
	var err error

	filesystem.fs, err = loadDataIntoSlice("puzzle-data.txt")
	if err != nil {
		panic(err)
	}
	filesystem.expandDiskMap()
	filesystem.compact()
	checksum := calcChecksum(filesystem.fs)

	fmt.Printf("Filesystem checksum is: '%d'\n", checksum)

	elapsed := time.Since(start) // Calculate elapsed time
	fmt.Printf("Execution time: %d ms\n", elapsed.Milliseconds())
}

func loadDataIntoSlice(filename string) ([]int, error) {
	var dataSlice = []int{}

	file, err := os.Open(filename)
	if err != nil {
		return []int{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		dataSlice = make([]int, len(line))
		for i, number := range line {
			dataSlice[i] = int(number - '0')
		}
	}

	return dataSlice, scanner.Err()
}

func (f *filesystem) expandDiskMap() {
	fsExpanded := []int{}
	fileId := 0

	for fsIndex := range f.fs {
		current := f.fs[fsIndex]
		if isEven(fsIndex) {
			used := [2]int{
				len(fsExpanded),
				current,
			}
			f.usedSpace = append(f.usedSpace, used)
			for block := 0; block < current; block++ {
				fsExpanded = append(fsExpanded, fileId)
			}
			fileId++
		} else {
			free := [2]int{
				len(fsExpanded),
				current,
			}
			f.freeSpace = append(f.freeSpace, free)
			for block := 0; block < current; block++ {
				fsExpanded = append(fsExpanded, 0)
			}
		}
	}

	f.fs = fsExpanded
}

func isEven(number int) bool {
	return number&1 == 0
}

func (f *filesystem) compact() {
	fs := make([]int, len(f.fs))
	copy(fs, f.fs)

	for j := len(f.usedSpace) - 1; j >= 0; j-- {
		for i := 0; i < len(f.freeSpace)-1; i++ {
			usedStart := f.usedSpace[j][0]
			usedCount := f.usedSpace[j][1]
			freeStart := f.freeSpace[i][0]
			freeCount := f.freeSpace[i][1]

			if freeStart > usedStart {
				break
			}
			if freeCount < usedCount {
				continue
			}

			for k := 0; k < usedCount; k++ {
				fs[freeStart+k] = fs[usedStart+k]
				fs[usedStart+k] = 0
			}

			f.freeSpace[i][1] -= usedCount
			f.freeSpace[i][0] += usedCount
			f.usedSpace = append(f.usedSpace[:j], f.usedSpace[j+1:]...)
			break
		}
	}

	f.fs = fs
}

func calcChecksum(input []int) int {
	checksum := 0

	for position := range input {
		checksum += input[position] * position
	}

	return checksum
}
