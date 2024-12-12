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
	stones = blink(stones, 25)

	fmt.Printf("Stonecount: '%d'\n", len(stones))

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

func blink(stones []int, count int) []int {
	newStones := make([]int, len(stones))
	copy(newStones, stones)

	fmt.Println(newStones)
	for i := 1; i <= count; i++ {
		newStones = changeStones(newStones)
	}

	return newStones
}

func changeStones(stones []int) []int {
	changedStones := []int{}
	for _, stone := range stones {
		if stone == 0 {
			changedStones = append(changedStones, handleZero())
			continue
		}
		if !isEven(stone) {
			changedStones = append(changedStones, handleDefault(stone))
			continue
		}
		splitStones := handleEven(stone)
		for _, splitStone := range splitStones {
			changedStones = append(changedStones, splitStone)
		}
	}

	return changedStones
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
