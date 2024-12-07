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
	labMap                       [][]string
	guardCoords                  []int
	distinctGuardLocations       int
	distinctObstructionLocations int
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
	lab.WalkGuard(0, false)

	fmt.Printf("Number of distinct guard locations: %d\n", lab.distinctGuardLocations)
	fmt.Printf("Number of distinct obstruction locations: %d\n", lab.distinctObstructionLocations)

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

func (l *Lab) inBoundsFromPoint(newX, newY int) bool {
	min := 0
	max := len(l.labMap) - 1

	if newX < min || newY < min || newX > max || newY > max {
		return false
	}
	return true
}

///// Part 2
//
// Not sure how I feel about this, but it runs under 50ms on my machine. Couldn't
// get it to track loop on only the xCount so I left both.
//
// It's messy today. It's also Friday, I'm calling it.

func (l *Lab) WalkGuard(dirCount int, simulated bool) string {
	steps := 0
	xCount := 0

	for {
		dir := directions[dirCount%len(directions)]
		row := l.guardCoords[0] + dir[0]
		column := l.guardCoords[1] + dir[1]

		if !l.inBoundsFromPoint(row, column) {
			return ""
		}

		if xCount > 1000 {
			return "loop"
		}
		if steps > 7000 {
			return "loop"
		}

		nextLocation := l.labMap[row][column]
		if nextLocation == "#" {
			dirCount++
			continue
		}
		if !simulated && nextLocation == "." {
			copyLab := l.copyLab()
			copyLab.labMap[row][column] = "#"
			loop := copyLab.WalkGuard(dirCount, true)
			if loop == "loop" {
				l.distinctObstructionLocations++
			}
		}
		if nextLocation == "." && !simulated {
			xCount = 0
			l.distinctGuardLocations++
			l.labMap[row][column] = "X"
		}
		if nextLocation == "X" && simulated {
			xCount++
		}

		l.guardCoords[0] = row
		l.guardCoords[1] = column
		steps++
	}
}

func (l *Lab) copyLab() Lab {
	mapCopy := make([][]string, len(l.labMap))
	copy(mapCopy, l.labMap)
	guardCoordsCopy := make([]int, 2)
	copy(guardCoordsCopy, l.guardCoords)

	labCopy := Lab{
		labMap:      mapCopy,
		guardCoords: guardCoordsCopy,
	}

	return labCopy
}
