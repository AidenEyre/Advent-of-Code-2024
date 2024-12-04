package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func main() {
	results := 0
	instructions, err := loadDataIntoSlice("puzzle-data.txt")
	if err != nil {
		panic(err)
	}

	for i := range instructions {
		scrapedInstructions, err := scrapeInstructionPairs(instructions[i])
		if err != nil {
			panic(err)
		}
		results += runInstructions(scrapedInstructions)
	}
	fmt.Printf("Instruction results are: '%d'\n", results)
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

func scrapeInstructionPairs(corruptedCode string) ([][2]int, error) {
	r := regexp.MustCompile(`mul\((\d+),(\d+)\)`)
	matches := r.FindAllStringSubmatch(corruptedCode, -1)

	var instructionPairs [][2]int
	var matchedIntPair [2]int
	var err error
	for _, match := range matches {
		if len(match) >= 3 {
			matchedIntPair, err = convertMatchToInt(match)
			instructionPairs = append(instructionPairs, matchedIntPair)
		}
		if err != nil {
			return [][2]int{}, fmt.Errorf("failed to scrape instruction pairs: %w", err)
		}
	}

	return instructionPairs, nil
}

func convertMatchToInt(match []string) ([2]int, error) {
	firstNum, err := strconv.Atoi(match[1])
	if err != nil {
		return [2]int{}, fmt.Errorf("failed to convert string to int: %w", err)
	}
	secondNum, err := strconv.Atoi(match[2])
	if err != nil {
		return [2]int{}, fmt.Errorf("failed to convert string to int: %w", err)
	}
	return [2]int{firstNum, secondNum}, nil
}

func runInstructions(instructions [][2]int) int {
	result := 0
	for i := range instructions {
		result += (instructions[i][0] * instructions[i][1])
	}
	return result
}
