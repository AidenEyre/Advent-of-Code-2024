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
		safe, err := isSafe(reports[i])
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

func isSafe(report string) (bool, error) {
	openedReport, err := stringToIntSlice(report)
	if err != nil {
		return false, fmt.Errorf("failed to analyze report: %w", err)
	}

	shouldIncrease := false
	if (openedReport[0] - openedReport[1]) > 0 {
		shouldIncrease = true
	}

	for i := 0; i < len(openedReport)-1; i++ {
		itemOne := openedReport[i]
		itemTwo := openedReport[i+1]
		isSafe := analyzePair(itemOne, itemTwo, shouldIncrease)

		if !isSafe {
			return false, nil
		}
	}

	return true, nil
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
