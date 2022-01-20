package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Current working directory: %v\n", workingDirectory)
	lines := readAllLines(workingDirectory + "/puzzle-49/input.txt")
	//lines := readAllLines(workingDirectory + "/puzzle-49/test.txt")

	numRows := len(lines)
	numColumns := len(lines[0])
	grid := NewGrid(numRows, numColumns)
	for line := range lines {
		for position, character := range lines[line] {
			grid[line][position] = character
		}
	}

	for step := 1; ; step++ {
		//fmt.Println("")
		//fmt.Println(grid)
		halfMoved, firstMoves := moveHorizontally(grid)
		fullyMoved, secondMoves := moveVertically(halfMoved)
		grid = fullyMoved
		if firstMoves+secondMoves == 0 {
			fmt.Println("No more movement after", step, "steps.")
			break
		}
	}
}

func readAllLines(p string) []string {
	file, err := os.Open(p)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

type Grid [][]rune

func (g Grid) String() string {
	lines := make([]string, len(g))
	for row := range g {
		lines[row] = string(g[row])
	}
	return strings.Join(lines, "\n")
}

func NewGrid(rows, columns int) Grid {
	grid := make(Grid, rows)
	for i := 0; i < rows; i++ {
		grid[i] = make([]rune, columns)
	}
	return grid
}

func moveHorizontally(grid Grid) (Grid, int) {
	next := NewGrid(len(grid), len(grid[0]))
	moves := 0
	for row := 0; row < len(grid); row++ {
		for column := 0; column < len(grid[row]); column++ {
			if next[row][column] != 0 { // Skip already visited fields
				continue
			} else if grid[row][column] != '>' { // Transfer irrelevant fields
				next[row][column] = grid[row][column]
				continue
			}

			nextColumn := column + 1
			if nextColumn >= len(grid[row]) {
				nextColumn = 0
			}
			if grid[row][nextColumn] == '.' {
				next[row][column] = '.'
				next[row][nextColumn] = '>'
				moves++
			} else {
				next[row][column] = '>'
			}
		}
	}
	return next, moves
}

func moveVertically(grid Grid) (Grid, int) {
	next := NewGrid(len(grid), len(grid[0]))
	moves := 0
	for row := 0; row < len(grid); row++ {
		for column := 0; column < len(grid[row]); column++ {
			if next[row][column] != 0 { // Skip already visited fields
				continue
			} else if grid[row][column] != 'v' { // Transfer irrelevant fields
				next[row][column] = grid[row][column]
				continue
			}

			nextRow := row + 1
			if nextRow >= len(grid) {
				nextRow = 0
			}
			if grid[nextRow][column] == '.' {
				next[row][column] = '.'
				next[nextRow][column] = 'v'
				moves++
			} else {
				next[row][column] = 'v'
			}
		}
	}
	return next, moves
}
