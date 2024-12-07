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

var ops = map[string]func(int, int) int{
	"+": func(a, b int) int { return a + b },
	"*": func(a, b int) int { return a * b },
	"||": func(a, b int) int {
		digits := int(math.Log10(float64(b)) + 1)
		result := a*int(math.Pow(10, float64(digits))) + b
		return result
	},
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
		availableOperators:     []string{"+", "*", "||"},
	}

	targets, numbers, err := loadDataIntoSlice("puzzle-data.txt")
	if err != nil {
		panic(err)
	}
	calibrator.targets = targets
	calibrator.numbers = numbers
	calibrator.calibrate()

	fmt.Printf("Total calibration result: '%d'\n", calibrator.totalCalibrationResult)

	elapsed := time.Since(start) // Calculate elapsed time
	fmt.Printf("Execution time: %d ms\n", elapsed.Milliseconds())
}

func loadDataIntoSlice(filename string) ([]int, [][]int, error) {
	var numbers = make([][]int, 0, 850)
	for i := range numbers {
		numbers[i] = make([]int, 0, 30)
	}
	var targets = make([]int, 0, 850)

	file, err := os.Open(filename)
	if err != nil {
		return []int{}, [][]int{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		splitTargetFromNumbers := strings.SplitN(line, ": ", 2)
		target, _ := strconv.Atoi(splitTargetFromNumbers[0])
		numberStrs := strings.Fields(splitTargetFromNumbers[1])
		numberSet := make([]int, len(numberStrs))
		for i, str := range numberStrs {
			numberSet[i], _ = strconv.Atoi(str)
		}

		targets = append(targets, target)
		numbers = append(numbers, numberSet)
	}

	return targets, numbers, scanner.Err()
}

func (c *calibrator3000) calibrate() {
	results := make(chan int, len(c.targets))
	work := make(chan int, len(c.targets))
	workerCount := 8 // my # of perf cores

	for w := 0; w < workerCount; w++ {
		go func() {
			for i := range work {
				go c.concurrentCalibrate(c.targets[i], c.numbers[i], results)
			}
		}()
	}

	for i := range c.targets {
		work <- i
	}
	close(work)

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

	for _, operator := range c.availableOperators {
		op := ops[operator]
		nextVal := op(value, numbers[index])

		if nextVal > target {
			continue
		}

		if c.recEvaluate(numbers, target, index+1, nextVal) {
			return true
		}
	}

	return false
}
