package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

var directions = [4][2]int{
	{1, 0},
	{-1, 0},
	{0, -1},
	{0, 1},
}

type hikingGuide struct {
	topoMap    [][]int // y, x grid
	trailheads [][]int // x, y coords
	peakCoords [][]int
	totalScore int
}

func main() {
	start := time.Now() // get time

	guide := hikingGuide{}
	var err error

	guide.topoMap, err = loadDataIntoSlice("puzzle-data.txt")
	if err != nil {
		panic(err)
	}
	guide.findTrailheads()
	guide.getAllTrailScores()

	fmt.Printf("total trailhead scores: '%d'\n", guide.totalScore)

	elapsed := time.Since(start) // Calculate elapsed time
	fmt.Printf("Execution time: %d ms\n", elapsed.Milliseconds())
}

func loadDataIntoSlice(filename string) ([][]int, error) {
	var dataGrid = [][]int{}

	file, err := os.Open(filename)
	if err != nil {
		return [][]int{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		row := make([]int, len(line))
		for i, char := range line {
			row[i] = int(char - '0')
		}
		dataGrid = append(dataGrid, row)
	}

	return dataGrid, scanner.Err()
}

func (hg *hikingGuide) findTrailheads() {
	for y := range hg.topoMap {
		for x := range hg.topoMap[y] {
			if hg.topoMap[y][x] == 0 {
				trailhead := []int{x, y, 0}
				hg.trailheads = append(hg.trailheads, trailhead)
			}
		}
	}
}

func (hg *hikingGuide) getAllTrailScores() {
	for _, th := range hg.trailheads {
		hg.peakCoords = [][]int{}
		hg.totalScore += hg.getTrailheadScore(th)
	}
}

func (hg *hikingGuide) getTrailheadScore(coords []int) int {
	target := 9
	score := 0

	if !hg.inBounds(coords) {
		return 0
	}
	if containsCoordsSlice(hg.peakCoords, coords) {
		return 0
	}
	if hg.topoMap[coords[1]][coords[0]] == target {
		// Uncomment this for part 1 answer.
		// hg.peakCoords = append(hg.peakCoords, coords)

		return 1
	}

	for _, dir := range directions {
		newX := coords[0] + dir[0]
		newY := coords[1] + dir[1]
		newCoords := []int{newX, newY}
		if !hg.isGradual(coords, newCoords) {
			continue
		}
		score += hg.getTrailheadScore(newCoords)
	}

	return score
}

func (hg *hikingGuide) inBounds(coords []int) bool {
	min := 0
	max := len(hg.topoMap) - 1

	if coords[0] < min || coords[1] < min || coords[0] > max || coords[1] > max {
		return false
	}
	return true
}

func (hg *hikingGuide) isGradual(coords, newCoords []int) bool {
	if !hg.inBounds(newCoords) {
		return false
	}

	height := hg.topoMap[coords[1]][coords[0]]
	newHeight := hg.topoMap[newCoords[1]][newCoords[0]]

	if newHeight-height == 1 {
		return true
	}

	return false
}

func containsCoordsSlice(slice [][]int, target []int) bool {
	for _, item := range slice {
		if item[0] == target[0] && item[1] == target[1] {
			return true
		}
	}
	return false
}
