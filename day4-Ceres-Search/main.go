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
	xy             [2]int
}

func main() {
	wordsearch, err := loadDataIntoSlice("puzzle-data.txt")
	if err != nil {
		panic(err)
	}
	solver := &WordsearchSolver5000{
		wordMatrix:     wordsearch,
		widthAndLength: len(wordsearch),
		Word:           "MAS",
		// Word:           "XMAS",
	}
	// wordsFoundStraight := solver.searchForWords(solver.Word[0], "regular")
	wordsFoundCross := solver.searchForWords(solver.Word[1], "cross")

	// fmt.Printf("WordsearchSolver5000 found '%d' instances of '%s' STRAIGHT", wordsFoundStraight, solver.Word)
	fmt.Printf("WordsearchSolver5000 found '%d' instances of '%s' CROSS", wordsFoundCross, solver.Word)
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

func (ws *WordsearchSolver5000) inBoundsFromPoint(direction string) bool {
	max := ws.widthAndLength - 1
	min := 0
	steps := len(ws.Word) - 1

	if direction == "all" {
		return ws.inBoundsAll(min, max)
	}

	dir := directions[direction]
	endX := ws.xy[0] + dir[0]*steps
	endY := ws.xy[1] + dir[1]*steps
	if endX < min || endY < min || endX > max || endY > max {
		return false
	}
	return true
}

func (ws *WordsearchSolver5000) inBoundsAll(min, max int) bool {
	for direction := range directions {
		dir := directions[direction]
		tmpX := ws.xy[0] + dir[0]
		tmpY := ws.xy[1] + dir[1]
		if tmpX < min || tmpY < min || tmpX > max || tmpY > max {
			return false
		}
	}
	return true
}

// searchForWords ... Solves crossword puzzle
//
//   - searchType=regular - 8 directions: horizontal, vertical, and diagonal.
//   - searchType=cross - requires the word to appear twice in a cross pattern with
//     startChar being the middle.
func (ws *WordsearchSolver5000) searchForWords(startChar byte, searchType string) int {
	totalWords := 0

	for x := range ws.wordMatrix {
		for y := range ws.wordMatrix[x] {
			if ws.wordMatrix[x][y] != startChar {
				continue
			}
			ws.xy = [2]int{x, y}

			if searchType == "regular" {
				totalWords += ws.searchStraightFromPoint()
			}
			if searchType == "cross" {
				totalWords += ws.searchCrossFromPoint()
			}
		}
	}

	return totalWords
}

func (ws *WordsearchSolver5000) searchStraightFromPoint() int {
	count := 0
	for direction := range directions {
		if ws.inBoundsFromPoint(direction) && ws.checkForMatchStraight(direction) {
			count++
		}
	}
	return count
}

func (ws *WordsearchSolver5000) checkForMatchStraight(direction string) bool {
	steps := len(ws.Word)
	dir := directions[direction]

	tmpX := ws.xy[0]
	tmpY := ws.xy[1]
	for i := 1; i < steps; i++ {
		tmpX = tmpX + dir[0]
		tmpY = tmpY + dir[1]
		if ws.wordMatrix[tmpX][tmpY] != ws.Word[i] {
			return false
		}
	}
	return true
}

func (ws *WordsearchSolver5000) searchCrossFromPoint() int {
	if !ws.inBoundsFromPoint("all") {
		return 0
	}

	firstChar := ws.Word[0]
	lastChar := ws.Word[len(ws.Word)-1]

	topLeft := ws.getCorner("upLeft")
	if topLeft != firstChar && topLeft != lastChar {
		return 0
	}

	bottomRight := ws.getCorner("downRight")
	if bottomRight != ws.getOppositeWordChar(topLeft) {
		return 0
	}

	topRight := ws.getCorner("upRight")
	if topRight != firstChar && topRight != lastChar {
		return 0
	}

	bottomLeft := ws.getCorner("downLeft")
	if bottomLeft != ws.getOppositeWordChar(topRight) {
		return 0
	}

	return 1
}

func (ws *WordsearchSolver5000) getCorner(direction string) byte {
	dir := directions[direction]
	tmpX := ws.xy[0] + dir[0]
	tmpY := ws.xy[1] + dir[1]
	return ws.wordMatrix[tmpX][tmpY]
}

func (ws *WordsearchSolver5000) getOppositeWordChar(input byte) byte {
	firstChar := ws.Word[0]
	lastChar := ws.Word[len(ws.Word)-1]

	if input == firstChar {
		return lastChar
	}
	return firstChar
}
