package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"
)

func main() {
	start := time.Now() // get time

	result := 0
	validUpdates := [][]int{}
	rules, err := loadDataIntoSlice("puzzle-data-rules.txt", "|")
	if err != nil {
		panic(err)
	}
	updates, err := loadDataIntoSlice("puzzle-data-updates.txt", ",")
	if err != nil {
		panic(err)
	}

	for i := range updates {
		rulesAdjacencyList := getAdjacencyListFromSlice(updates[i], rules)
		topologicalOrder := topologicalDFSSort(rulesAdjacencyList)

		isValid := validateUpdateOrder(topologicalOrder, updates[i])
		if isValid {
			validUpdates = append(validUpdates, updates[i])
		}
	}

	for i := range validUpdates {
		result += getMiddlePageNumber(validUpdates[i])
	}

	fmt.Printf("Sum of middle numbers in valid updates is: '%d'\n", result)

	elapsed := time.Since(start) // Calculate elapsed time
	fmt.Printf("Execution time: %d ms\n", elapsed.Milliseconds())
}

func loadDataIntoSlice(filename string, delimiter string) ([][]int, error) {
	var dataSlice = [][]int{}

	file, err := os.Open(filename)
	if err != nil {
		return [][]int{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		stringsSlice := strings.Split(line, delimiter)
		intSlice := make([]int, len(stringsSlice))

		intSlice, err = convertSliceStringToInt(stringsSlice)
		dataSlice = append(dataSlice, intSlice)
		if err != nil {
			return [][]int{}, fmt.Errorf("failed to parse line, '%s': %w", line, err)
		}
	}
	if err != nil {
		return [][]int{}, fmt.Errorf("failed to parse file: %s", err)
	}

	return dataSlice, nil
}

func convertSliceStringToInt(input []string) ([]int, error) {
	result := make([]int, len(input))
	var err error
	for i := range input {
		result[i], err = strconv.Atoi(input[i])
		if err != nil {
			return []int{}, fmt.Errorf("failed to convert string to int: %w", err)
		}
	}
	return result, nil
}

func getAdjacencyListFromSlice(input []int, rules [][]int) map[int][]int {
	scrapedList := [][]int{}

	for i := range input {
		for j := range rules {
			if input[i] == rules[j][0] {
				scrapedList = append(scrapedList, rules[j])
			}
		}
	}

	return createAdjacencyList(scrapedList)
}

func createAdjacencyList(input [][]int) map[int][]int {
	result := make(map[int][]int)

	for i := range input {
		result = addToMap(input[i][0], input[i][1], result)
	}

	return result
}

func addToMap(key, value int, inputMap map[int][]int) map[int][]int {
	_, exists := inputMap[key]
	if !exists {
		inputMap[key] = []int{}
	}

	inputMap[key] = append(inputMap[key], value)
	return inputMap
}

func topologicalDFSSort(graph map[int][]int) []int {
	visited := make(map[int]bool)
	order := []int{}

	for node := range graph {
		if !visited[node] {
			handleDFS(node, visited, &order, graph)
		}
	}

	slices.Reverse(order)

	return order
}

func handleDFS(node int, visited map[int]bool, order *[]int, graph map[int][]int) {
	stack := []int{node}

	for len(stack) > 0 {
		currentNode := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if visited[currentNode] {
			continue
		}

		visited[currentNode] = true

		for _, neighbor := range graph[currentNode] {
			if !visited[neighbor] {
				handleDFS(neighbor, visited, order, graph)
			}
		}

		*order = append(*order, currentNode)
	}
}

func validateUpdateOrder(topologicalOrder, update []int) bool {
	previous := 0
	orderIndex := make(map[int]int, len(topologicalOrder))
	for i, node := range topologicalOrder {
		orderIndex[node] = i
	}

	for i := range update {
		current := orderIndex[update[i]]
		if current < previous {
			return false
		}
		previous = current
	}

	return true
}

func getMiddlePageNumber(pages []int) int {
	numPages := len(pages)
	middleIndex := numPages / 2
	return pages[middleIndex]
}
