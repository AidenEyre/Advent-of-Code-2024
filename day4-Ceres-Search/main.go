package main

import (
	"bufio"
	"fmt"
	"os"
)

var directions = map[string][2]int{
	"right":     {1, 0},
	"left":      {-1, 0},
	"down":      {0, -1},
	"up":        {0, 1},
	"downRight": {1, -1},
	"upLeft":    {-1, 1},
	"downLeft":  {-1, -1},
	"upRight":   {1, 1},
}

type WordsearchSolver5000 struct {
	wordMatrix     [][]byte
	widthAndLength int
	WordLength     int
	Word           string
}

func main() {
	wordsearch, err := loadDataIntoSlice("puzzle-data.txt")
	if err != nil {
		panic(err)
	}
	solver := &WordsearchSolver5000{
		wordMatrix:     wordsearch,
		widthAndLength: len(wordsearch),
		Word:           "XMAS",
	}
	wordsFound := solver.searchForWords(solver.Word[0], "regular")

	fmt.Printf("WordsearchSolver5000 found '%d' instances of '%s'", wordsFound, solver.Word)
}

func loadDataIntoSlice(filename string) ([][]byte, error) {
	var dataSlice [][]byte

	file, err := os.Open(filename)
	if err != nil {
		return [][]byte{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		dataSlice = append(dataSlice, []byte(line))
	}
	if err != nil {
		return [][]byte{}, fmt.Errorf("failed to parse file: %s", err)
	}

	return dataSlice, nil
}

func (ws *WordsearchSolver5000) inBoundsFromPoint(x, y int, direction string) bool {
	max := ws.widthAndLength - 1
	min := 0

	dir, exists := directions[direction]
	if !exists {
		fmt.Printf("Invalid direction: %s\n", direction)
		return false
	}

	steps := len(ws.Word) - 1
	endX := x + dir[0]*steps
	endY := y + dir[1]*steps

	if endX < min || endY < min || endX > max || endY > max {
		return false
	}
	return true
}

func (ws *WordsearchSolver5000) searchForWords(startChar byte, searchType string) int {
	totalWords := 0

	for x := range ws.wordMatrix {
		for y := range ws.wordMatrix[x] {
			if ws.wordMatrix[x][y] != startChar {
				continue
			}

			if searchType == "regular" {
				totalWords += ws.searchStraightFromPoint(x, y)
			}
			if searchType == "cross" {
				totalWords += ws.searchCrossFromPoint(x, y)
			}
		}
	}

	return totalWords
}

func (ws *WordsearchSolver5000) searchStraightFromPoint(x, y int) int {
	count := 0
	for direction := range directions {
		if ws.inBoundsFromPoint(x, y, direction) && ws.checkForMatchStraight(x, y, direction) {
			count++
		}
	}
	return count
}

func (ws *WordsearchSolver5000) checkForMatchStraight(x, y int, direction string) bool {
	steps := len(ws.Word)
	dir, exists := directions[direction]
	if !exists {
		fmt.Printf("Invalid direction: %s\n", direction)
		return false
	}

	// fmt.Printf("\nDIRECTION: %s\n(%d,%d) - %s\n", direction, x, y, string(ws.wordMatrix[x][y]))
	for i := 1; i < steps; i++ {
		x = x + dir[0]
		y = y + dir[1]
		// fmt.Printf("(%d,%d) - %s\n", x, y, string(ws.wordMatrix[x][y]))
		if ws.wordMatrix[x][y] != ws.Word[i] {
			return false
		}
	}
	// fmt.Println("HERE")
	return true
}

func (ws *WordsearchSolver5000) searchCrossFromPoint(x, y int) int {
	return 0
}
