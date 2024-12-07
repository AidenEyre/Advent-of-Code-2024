package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

var ops = map[string]func(int, int) int{
	"+": func(a, b int) int { return a + b },
	"-": func(a, b int) int { return a - b },
	"*": func(a, b int) int { return a * b },
	"/": func(a, b int) int { return a / b },
	"%": func(a, b int) int { return a % b },
}

type calibrator3000 struct {
	targets                []int
	numbers                [][]int
	availableOperators     []string
	totalCalibrationResult int
}

func main() {
	start := time.Now() // get time

	calibrator := calibrator3000{
		totalCalibrationResult: 0,
		availableOperators:     []string{"+", "*"},
	}

	targets, numbers, err := loadDataIntoSlice("puzzle-data.txt")
	if err != nil {
		panic(err)
	}
	calibrator.targets = targets
	calibrator.numbers = numbers
	calibrator.calibrate()

	for i := range targets {
		fmt.Printf("%d: %v\n", calibrator.targets[i], calibrator.numbers[i])
	}

	fmt.Printf("Total calibration result: '%d'\n", calibrator.totalCalibrationResult)

	elapsed := time.Since(start) // Calculate elapsed time
	fmt.Printf("Execution time: %d ms\n", elapsed.Milliseconds())
}

func loadDataIntoSlice(filename string) ([]int, [][]int, error) {
	var numbers = [][]int{}
	var targets = []int{}

	file, err := os.Open(filename)
	if err != nil {
		return []int{}, [][]int{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		splitTargetFromNumbers := strings.Split(line, ": ")
		target, _ := strconv.Atoi(splitTargetFromNumbers[0])
		numberSet := stringSliceToInt(strings.Split(splitTargetFromNumbers[1], " "))

		targets = append(targets, target)
		numbers = append(numbers, numberSet)
	}

	return targets, numbers, nil
}

func stringSliceToInt(input []string) []int {
	intSlice := make([]int, len(input))
	for i := range input {
		converted, _ := strconv.Atoi(input[i])
		intSlice[i] = converted
	}
	return intSlice
}

func (c *calibrator3000) calibrate() {
	results := make(chan int, len(c.targets))

	for i := range c.targets {
		go c.concurrentCalibrate(c.targets[i], c.numbers[i], results)
	}
	for range c.targets {
		c.totalCalibrationResult += <-results
	}

	close(results)
}

func (c *calibrator3000) concurrentCalibrate(target int, numbers []int, results chan int) {
	if c.recEvaluate(numbers, target, 1, numbers[0]) {
		results <- target
	} else {
		results <- 0
	}
}

func (c *calibrator3000) recEvaluate(numbers []int, target, index, value int) bool {
	if index == len(numbers) {
		return value == target
	}
	if value > target {
		return false
	}

	for i := range c.availableOperators {
		nextVal := ops[c.availableOperators[i]](value, numbers[index])
		if c.recEvaluate(numbers, target, index+1, nextVal) {
			return true
		}
	}

	return false
}
