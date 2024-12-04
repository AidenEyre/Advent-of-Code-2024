package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
)

func main() {
	instructions, err := loadDataIntoSlice("puzzle-data.txt")
	if err != nil {
		panic(err)
	}

	var combinedInstructions [][2]int
	for i := range instructions {
		scrapedInstructions, err := scrapeInstructionPairs(instructions[i])
		if err != nil {
			panic(err)
		}
		combinedInstructions = append(combinedInstructions, scrapedInstructions...)
	}
	results := runInstructions(combinedInstructions)

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
	r := regexp.MustCompile(`mul\((\d+),(\d+)\)|do\(\)|don't\(\)`)
	matches := r.FindAllStringSubmatch(corruptedCode, -1)

	var instructionPairs [][2]int
	var err error
	var pair [2]int
	for _, match := range matches {
		pair, err = handleMatch(match)
		if err != nil {
			return [][2]int{}, fmt.Errorf("failed to scrape instruction pairs: %w", err)
		}
		instructionPairs = append(instructionPairs, pair)
	}

	return instructionPairs, nil
}

func handleMatch(match []string) ([2]int, error) {
	var err error
	var matchedIntPair [2]int
	if match[0] == "do()" {
		return [2]int{0, 0}, nil
	} else if match[0] == "don't()" {
		return [2]int{-1, -1}, nil
	} else if len(match) >= 3 && match[1] != "" && match[2] != "" {
		matchedIntPair, err = convertMatchToInt(match)
		return matchedIntPair, nil
	}
	if err != nil {
		return [2]int{}, fmt.Errorf("failed to handle match: %w", err)
	}
	return [2]int{}, fmt.Errorf("failed to handle match")
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
	enabled := true
	result := 0
	for i := range instructions {
		if instructions[i][0] == 0 && instructions[i][1] == 0 {
			enabled = true
		}
		if instructions[i][0] == -1 && instructions[i][1] == -1 {
			enabled = false
		}
		if !enabled {
			continue
		}
		result += (instructions[i][0] * instructions[i][1])
	}
	return result
}
