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

// This feels rough. I spent too much time stuck on something silly though,
// calling it.

type cityScanner struct {
	cityMap                 [][]string
	antennaCoords           map[string][][]int
	uniqueAntinodeLocations int
}

func main() {
	start := time.Now() // get time

	scanner := cityScanner{
		antennaCoords:           make(map[string][][]int),
		uniqueAntinodeLocations: 0,
	}

	cityMap, err := loadDataIntoSlice("puzzle-data.txt")
	if err != nil {
		panic(err)
	}
	scanner.cityMap = cityMap
	scanner.logAntennaLocations()
	scanner.findUniqueAntinodes()

	// scanner.printAntennaMap()
	// scanner.printCityMap()
	fmt.Printf("Unique antinode locations: '%d'\n", scanner.uniqueAntinodeLocations)

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
	antinodeCoords := [][]int{}

	for _, antennas := range cs.antennaCoords {
		for i := range antennas {
			for j := range antennas {
				a1, a2 := antennas[i], antennas[j]
				if slices.Equal(a1, a2) {
					continue
				}
				antinodes := findPotentialAntinodeLocations(a1, a2)
				for i := range antinodes {
					if cs.inBounds(antinodes[i]) && !containsCoordsSlice(antinodeCoords, antinodes[i]) {
						antinodeCoords = append(antinodeCoords, antinodes[i])
						// fmt.Printf("Antenna: %s potential: (%d,%d)\n", key, antinodes[i][0], antinodes[i][1])
						cs.uniqueAntinodeLocations++
					}
				}
			}
		}
	}
}

func findPotentialAntinodeLocations(a1, a2 []int) [][]int {
	movementX := math.Abs(float64(a1[0] - a2[0]))
	movementY := math.Abs(float64(a1[1] - a2[1]))

	// I don't like how I did this, there's got to be a better way. I spent too
	// much time on this today tough
	if a2[0] >= a1[0] && a2[1] <= a1[1] {
		return [][]int{
			{
				a2[0] + int(movementX),
				a2[1] - int(movementY),
			},
			{
				a1[0] - int(movementX),
				a1[1] + int(movementY),
			},
		}
	}
	if a2[0] >= a1[0] && a2[1] >= a1[1] {
		return [][]int{
			{
				a2[0] + int(movementX),
				a2[1] + int(movementY),
			},
			{
				a1[0] - int(movementX),
				a1[1] - int(movementY),
			},
		}
	}
	if a2[0] <= a1[0] && a2[1] >= a1[1] {
		return [][]int{
			{
				a2[0] - int(movementX),
				a2[1] + int(movementY),
			},
			{
				a1[0] + int(movementX),
				a1[1] - int(movementY),
			},
		}
	}
	if a2[0] <= a1[0] && a2[1] <= a1[1] {
		return [][]int{
			{
				a2[0] - int(movementX),
				a2[1] - int(movementY),
			},
			{
				a1[0] + int(movementX),
				a1[1] + int(movementY),
			},
		}
	}

	return [][]int{}
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

// For testing
func (cs *cityScanner) printAntennaMap() {
	for key, locations := range cs.antennaCoords {
		fmt.Printf("Antenna Type '%s':\n", key)
		for _, coord := range locations {
			fmt.Printf("  - Location: (%d, %d)\n", coord[0], coord[1])
		}
	}
}

func (cs *cityScanner) printCityMap() {
	for i := range cs.cityMap {
		fmt.Println(cs.cityMap[i])
	}
}
