package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	start := time.Now() // get time

	var err error
	stones, err := loadDataIntoSlice("puzzle-data.txt")
	if err != nil {
		panic(err)
	}

	stonesMap := map[int]int{}
	for _, stone := range stones {
		stonesMap[stone] = 1
	}

	stoneCount := simulateBlinks(stonesMap, 75)

	fmt.Printf("Stonecount: '%d'\n", stoneCount)

	elapsed := time.Since(start) // Calculate elapsed time
	fmt.Printf("Execution time: %d ms\n", elapsed.Milliseconds())
}

func loadDataIntoSlice(filename string) ([]int, error) {
	var stones = []int{}

	file, err := os.Open(filename)
	if err != nil {
		return []int{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		numbers := strings.Split(line, " ")
		for _, number := range numbers {
			stone, _ := strconv.Atoi(number)
			stones = append(stones, stone)
		}
	}

	return stones, scanner.Err()
}

func simulateBlinks(stones map[int]int, blinks int) int {
	for i := 1; i <= blinks; i++ {
		stones = processStones(stones)
	}

	// Sum all stone counts
	totalStones := 0
	for _, count := range stones {
		totalStones += count
	}
	return totalStones
}

func processStones(stones map[int]int) map[int]int {
	newStones := make(map[int]int)

	for stone, count := range stones {
		if stone == 0 {
			newStones[1] += count
		} else if isEven(stone) {
			split := handleEven(stone)
			newStones[split[0]] += count
			newStones[split[1]] += count
		} else {
			newStones[stone*2024] += count
		}
	}

	return newStones
}

// RULE FUNCTIONS
/////////////////

func handleZero() int {
	return 1
}

func isEven(stone int) bool {
	numDigits := int(math.Log10(float64(stone)) + 1)
	if numDigits%2 == 0 {
		return true
	}
	return false
}

func handleEven(stone int) []int {
	base, divisor := 10, 10

	for stone/divisor > divisor {
		divisor *= base
	}

	firstHalf := stone / divisor
	secondHalf := stone % divisor

	return []int{firstHalf, secondHalf}
}

func handleDefault(stone int) int {
	return stone * 2024
}
