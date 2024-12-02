package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

type SignificantLocations struct {
	listOne   []int
	listTwo   []int
	distApart []int
}

func main() {
	sl := &SignificantLocations{}

	err := sl.loadData("puzzle-data.txt")
	if err != nil {
		panic(err)
	}

	// Part 1
	sl.sortLists()
	sl.calcDist()
	totalDist := sl.totalDist()
	fmt.Printf("Total Distance is '%d'\n", totalDist)

	// Part 2
	similarityScore := sl.calcSimilarityScore()
	fmt.Printf("Similarity Score is '%d'\n", similarityScore)
}

func (sl *SignificantLocations) loadData(filename string) error {
	// Read data and sort into two slices.
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var itemOne int
	var itemTwo int
	for scanner.Scan() {
		line := scanner.Text()
		itemOne, itemTwo, err = sl.parseLine(line)

		sl.listOne = append(sl.listOne, itemOne)
		sl.listTwo = append(sl.listTwo, itemTwo)
	}

	if err != nil {
		return fmt.Errorf("failed to parse file: %s", err)
	}
	return nil
}

func (sl *SignificantLocations) parseLine(line string) (int, int, error) {
	splitLine := strings.Fields(line)

	itemOne, err := strconv.Atoi(splitLine[0])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to convert string to int: %w", err)
	}
	itemTwo, err := strconv.Atoi(splitLine[1])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to convert string to int: %w", err)
	}

	return itemOne, itemTwo, nil
}

func (sl *SignificantLocations) sortLists() {
	sort.Ints(sl.listOne)
	sort.Ints(sl.listTwo)
}

func (sl *SignificantLocations) calcDist() {
	sl.distApart = make([]int, len(sl.listOne))
	for i := range sl.listOne {
		absDiff := math.Abs(float64(sl.listOne[i] - sl.listTwo[i]))
		sl.distApart[i] = int(absDiff)
	}
}

func (sl *SignificantLocations) totalDist() int {
	totalDist := 0
	for i := range sl.distApart {
		totalDist += sl.distApart[i]
	}
	return totalDist
}

func (sl *SignificantLocations) calcSimilarityScore() int {
	similarityScore := 0
	for i := range sl.listOne {
		occurrences := getOccurrences(sl.listTwo, sl.listOne[i])
		similarityScore += sl.listOne[i] * occurrences
	}
	return similarityScore
}

func getOccurrences(list []int, item int) int {
	occurrences := 0
	for i := range list {
		if list[i] == item {
			occurrences++
		}
	}
	return occurrences
}
