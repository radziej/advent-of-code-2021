package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Current working directory: %v\n", workingDirectory)

	lines := readAllLines(workingDirectory + "/puzzle-39/input.txt")
	//lines := readAllLines(workingDirectory + "/puzzle-39/test.txt")

	enhancePattern := make(map[int]bool, 512)
	for i, r := range lines[0] {
		if string(r) == "#" {
			enhancePattern[i] = true
		} else {
			enhancePattern[i] = false
		}
	}

	grid := make([][]bool, len(lines)-2)
	for y, line := range lines[2:] {
		grid[y] = make([]bool, len(line))
		for x, r := range line {
			if string(r) == "#" {
				grid[y][x] = true
			} else {
				grid[y][x] = false
			}
		}
	}
	fmt.Println(SPrintGrid(grid))

	for step := 0; step < 2; step++ {
		// Set up next grid by expanding by one row/column on each side
		nextGrid := make([][]bool, len(grid)+2)
		for i := range nextGrid {
			nextGrid[i] = make([]bool, len(grid[0])+2)
		}

		for i := range nextGrid {
			for j := range nextGrid[i] {
				// Shift iterators back by the expanded row/column
				nextGrid[i][j] = enhanceAt(&grid, &enhancePattern, j-1, i-1, enhancePattern[0] && step%2 != 0)
			}
		}
		grid = nextGrid
		fmt.Println(SPrintGrid(grid))
	}

	litPixels := 0
	for _, row := range grid {
		for _, value := range row {
			if value {
				litPixels++
			}
		}
	}
	fmt.Println("Total amount of lit pixels in fully enhanced image:", litPixels)
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

func SPrintGrid(grid [][]bool) string {
	representation := ""
	for _, row := range grid {
		line := make([]string, len(row))
		index := 0
		for i, field := range row {
			if field {
				line[i] = "#"
			} else {
				line[i] = "."
			}
			index++
		}
		representation += strings.Join(line, "") + "\n"
	}
	return representation
}

func enhanceAt(grid *[][]bool, pattern *map[int]bool, x int, y int, void bool) bool {
	binaryRepresentation := make([]string, 9)
	for i, yOffset := range [3]int{y - 1, y, y + 1} {
		for j, xOffset := range [3]int{x - 1, x, x + 1} {
			value := false
			if yOffset < 0 || yOffset >= len(*grid) || xOffset < 0 || xOffset >= len((*grid)[0]) {
				value = void
			} else {
				value = (*grid)[yOffset][xOffset]
			}

			if value {
				binaryRepresentation[3*i+j] = "1"
			} else {
				binaryRepresentation[3*i+j] = "0"
			}
		}
	}
	number, err := strconv.ParseInt(strings.Join(binaryRepresentation, ""), 2, 64)
	if x == 2 && y == 2 {
		fmt.Println(strings.Join(binaryRepresentation, ""), number)
	}
	if err != nil {
		log.Fatal(err)
	}
	return (*pattern)[int(number)]
}
