package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

var directions = [4][2]int{
	{-1, 0}, // Up
	{0, 1},  // Right
	{1, 0},  // Down
	{0, -1}, // Left
}

type Lab struct {
	labMap                 [][]string
	guardCoords            []int
	distinctGuardLocations int
}

func main() {
	start := time.Now() // get time

	lab := Lab{
		distinctGuardLocations: 1,
	}
	var err error
	lab.labMap, err = loadDataIntoSlice("puzzle-data.txt", "|")
	if err != nil {
		panic(err)
	}
	lab.findGuard()
	lab.walkGuard()

	fmt.Printf("Number of distinct guard locations: %d\n", lab.distinctGuardLocations)

	elapsed := time.Since(start) // Calculate elapsed time
	fmt.Printf("Execution time: %d ms\n", elapsed.Milliseconds())
}

func loadDataIntoSlice(filename string, delimiter string) ([][]string, error) {
	var dataSlice = [][]string{}

	file, err := os.Open(filename)
	if err != nil {
		return [][]string{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		dataSlice = append(dataSlice, strings.Split(line, ""))
		if err != nil {
			return [][]string{}, fmt.Errorf("failed to parse line, '%s': %w", line, err)
		}
	}
	if err != nil {
		return [][]string{}, fmt.Errorf("failed to parse file: %s", err)
	}

	return dataSlice, nil
}

func (l *Lab) findGuard() {
	for row := range l.labMap {
		for column := range l.labMap[row] {
			if l.labMap[row][column] == "^" {
				l.guardCoords = []int{row, column}
				return
			}
		}
	}
}

func (l *Lab) walkGuard() {
	i := 0

	for {
		dir := directions[i%len(directions)]
		row := l.guardCoords[0] + dir[0]
		column := l.guardCoords[1] + dir[1]

		if !l.inBoundsFromPoint(row, column) {
			return
		}

		nextLocation := l.labMap[row][column]
		if nextLocation == "#" {
			i++
			continue
		}
		if nextLocation == "." {
			l.distinctGuardLocations++
			l.labMap[row][column] = "X"
		}

		l.guardCoords[0] = row
		l.guardCoords[1] = column
	}
}

func (l *Lab) inBoundsFromPoint(newX, newY int) bool {
	min := 0
	max := len(l.labMap) - 1

	if newX < min || newY < min || newX > max || newY > max {
		return false
	}
	return true
}
