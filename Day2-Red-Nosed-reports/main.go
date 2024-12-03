package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	var reports = []string{}
	var safeCount = 0

	reports, err := loadDataIntoSlice("puzzle-data.txt")
	if err != nil {
		panic(err)
	}

	for i := range reports {
		report, err := stringToIntSlice(reports[i])
		if err != nil {
			panic(err)
		}
		safe, err := isSafe(report, false)
		if err != nil {
			panic(err)
		}
		if safe {
			safeCount++
		}
	}

	fmt.Printf("Number of safe reports: '%d'", safeCount)
}

func loadDataIntoSlice(filename string) ([]string, error) {
	var dataSlice = []string{}

	file, err := os.Open(filename)
	if err != nil {
		return []string{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		dataSlice = append(dataSlice, line)
	}
	if err != nil {
		return []string{}, fmt.Errorf("failed to parse file: %s", err)
	}

	return dataSlice, nil
}

func isSafe(report []int, isRecursivePass bool) (bool, error) {
	var err error
	shouldIncrease := false
	if (report[0] - report[1]) > 0 {
		shouldIncrease = true
	}

	for i := 0; i < len(report)-1; i++ {
		itemOne := report[i]
		itemTwo := report[i+1]
		isSafe := analyzePair(itemOne, itemTwo, shouldIncrease)

		if isSafe {
			continue
		}

		if !isRecursivePass {
			isSafe, err = testToleration(report, i, i+1)
			if err != nil {
				return false, fmt.Errorf("failed to analyze report: %w", err)
			}
		}
		if isSafe {
			return true, nil
		}
		return false, nil
	}
	return true, nil
}

func testToleration(report []int, firstIndex, secondIndex int) (bool, error) {
	safe := false
	var err error

	// I gave up. After hours of working at this, I realized, I was missing the
	// case where removing the first item would solve the error in item 2/3.
	// For example, `57 56 57 59 60 63 64 65`. This is because it was assuming
	// it was counting down... It's late, so I'm brute forcing!
	for i := range report {
		safe, err = isSafe(removeAtIndex(report, i), true)
		if err != nil {
			return false, err
		}
		if safe {
			return true, nil
		}
	}

	return false, nil
}

func analyzePair(a, b int, shouldIncrease bool) bool {
	absDiff, isIncreasing := absAndSign(a, b)

	if absDiff > 3 || absDiff == 0 {
		return false
	}

	if isIncreasing != shouldIncrease {
		return false
	}

	return true
}

func stringToIntSlice(input string) ([]int, error) {
	var err error
	inputStrings := strings.Fields(input)
	result := make([]int, len(inputStrings))
	for i := range inputStrings {
		result[i], err = strconv.Atoi(inputStrings[i])
		if err != nil {
			return []int{}, fmt.Errorf("failed to convert []string to []int: %w", err)
		}
	}
	return result, nil
}

func absAndSign(a, b int) (int, bool) {
	diff := a - b
	isIncreasing := diff > 0
	if diff < 0 {
		diff = -diff
	}
	return diff, isIncreasing
}

func removeAtIndex(slice []int, s int) []int {
	newSlice := make([]int, len(slice))
	copy(newSlice, slice)
	return append(newSlice[:s], newSlice[s+1:]...)
}
