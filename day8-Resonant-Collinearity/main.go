package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"slices"
	"strings"
	"time"
)

// This is an ugly solution... I spent too much time stuck on something silly
// though, calling it good.

type cityScanner struct {
	cityMap        [][]string
	antennaCoords  map[string][][]int
	antinodeCoords [][]int
}

func main() {
	start := time.Now() // get time

	scanner := cityScanner{
		antennaCoords: make(map[string][][]int),
	}

	cityMap, err := loadDataIntoSlice("puzzle-data.txt")
	if err != nil {
		panic(err)
	}
	scanner.cityMap = cityMap
	scanner.logAntennaLocations()
	scanner.findUniqueAntinodes()
	fmt.Printf("Unique antinode locations: '%d'\n", len(scanner.antinodeCoords))

	elapsed := time.Since(start) // Calculate elapsed time
	fmt.Printf("Execution time: %d ms\n", elapsed.Milliseconds())
}

func loadDataIntoSlice(filename string) ([][]string, error) {
	var dataSlice = make([][]string, 0, 50)
	for i := range dataSlice {
		dataSlice[i] = make([]string, 0, 50)
	}

	file, err := os.Open(filename)
	if err != nil {
		return [][]string{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		dataSlice = append(dataSlice, strings.Split(line, ""))
	}

	return dataSlice, scanner.Err()
}

func (cs *cityScanner) logAntennaLocations() {
	for y := range cs.cityMap {
		for x := range cs.cityMap[y] {
			key := cs.cityMap[y][x]
			if key == "." {
				continue
			}
			cs.antennaCoords[key] = append(cs.antennaCoords[key], []int{x, y})
		}
	}
}

func (cs *cityScanner) findUniqueAntinodes() {
	for _, antennas := range cs.antennaCoords {
		for i := range antennas {
			for j := range antennas {
				a1, a2 := antennas[i], antennas[j]
				if slices.Equal(a1, a2) {
					continue
				}
				cs.calculateAntinodeLocations(a1, a2)
			}
		}
	}
}

func (cs *cityScanner) calculateAntinodeLocations(a1, a2 []int) {
	movementX := math.Abs(float64(a1[0] - a2[0]))
	movementY := math.Abs(float64(a1[1] - a2[1]))
	invalidCount := 0
	multiplier := 1

	for {
		a1New := make([]int, 2)
		a2New := make([]int, 2)
		if a2[0] >= a1[0] && a2[1] <= a1[1] {
			a1New[0] = a1[0] + (int(movementX) * multiplier)
			a1New[1] = a1[1] - (int(movementY) * multiplier)
			a2New[0] = a2[0] - (int(movementX) * multiplier)
			a2New[1] = a2[1] + (int(movementY) * multiplier)
		}
		if a2[0] >= a1[0] && a2[1] >= a1[1] {
			a1New[0] = a1[0] + (int(movementX) * multiplier)
			a1New[1] = a1[1] + (int(movementY) * multiplier)
			a2New[0] = a2[0] - (int(movementX) * multiplier)
			a2New[1] = a2[1] - (int(movementY) * multiplier)
		}
		if a2[0] <= a1[0] && a2[1] >= a1[1] {
			a1New[0] = a1[0] - (int(movementX) * multiplier)
			a1New[1] = a1[1] + (int(movementY) * multiplier)
			a2New[0] = a2[0] + (int(movementX) * multiplier)
			a2New[1] = a2[1] - (int(movementY) * multiplier)
		}
		if a2[0] <= a1[0] && a2[1] <= a1[1] {
			a1New[0] = a1[0] - (int(movementX) * multiplier)
			a1New[1] = a1[1] - (int(movementY) * multiplier)
			a2New[0] = a2[0] + (int(movementX) * multiplier)
			a2New[1] = a2[1] + (int(movementY) * multiplier)
		}
		a1Valid := cs.isAntinodeValid(a1New)
		a2Valid := cs.isAntinodeValid(a2New)
		if a1Valid {
			invalidCount = 0
			cs.antinodeCoords = append(cs.antinodeCoords, a1New)
		}
		if a2Valid {
			invalidCount = 0
			cs.antinodeCoords = append(cs.antinodeCoords, a2New)
		}

		if !a1Valid && !a2Valid {
			invalidCount++
		}
		if invalidCount > 10 {
			return
		}

		multiplier++
	}
}

func (cs *cityScanner) isAntinodeValid(node []int) bool {
	if !cs.inBounds(node) {
		return false
	}
	if containsCoordsSlice(cs.antinodeCoords, node) {
		return false
	}

	return true
}

func (cs *cityScanner) inBounds(location []int) bool {
	min := 0
	max := len(cs.cityMap) - 1
	x := location[0]
	y := location[1]

	if x < min || y < min || x > max || y > max {
		return false
	}
	return true
}

func containsCoordsSlice(slice [][]int, target []int) bool {
	for _, item := range slice {
		if item[0] == target[0] && item[1] == target[1] {
			return true
		}
	}
	return false
}
