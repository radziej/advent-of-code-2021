package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var grid [][]int
	for line := range readLines(workingDirectory + "/puzzle-17/input.txt") {
		var digits []int
		for _, character := range line {
			digit, err := strconv.Atoi(string(character))
			if err != nil {
				log.Fatal(err)
			}
			digits = append(digits, digit)
		}
		grid = append(grid, digits)
	}

	riskLevels := make([][]int, len(grid))
	for i := range grid {
		riskLevels[i] = make([]int, len(grid[0]))
	}
	for row := range grid {
		for column := range grid[row] {
			if isLowPoint(&grid, row, column) {
				riskLevels[row][column] = grid[row][column] + 1
			} else {
				riskLevels[row][column] = 0
			}
		}
	}

	count := 0
	risk := 0
	for row := range riskLevels {
		for column := range riskLevels[row] {
			if riskLevels[row][column] > 0 {
				count++
				risk += riskLevels[row][column]
			}
		}
	}
	fmt.Printf("Number of low points %v, with total risk of %v\n", count, risk)
}

func readLines(p string) chan string {
	channel := make(chan string, 1)

	go func() {
		defer close(channel)

		file, err := os.Open(p)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			channel <- scanner.Text()
		}
	}()
	return channel
}

func isLowPoint(grid *[][]int, row int, column int) bool {
	value := (*grid)[row][column]
	if row-1 >= 0 && value >= (*grid)[row-1][column] {
		return false
	} else if row+1 < len(*grid) && value >= (*grid)[row+1][column] {
		return false
	} else if column-1 >= 0 && value >= (*grid)[row][column-1] {
		return false
	} else if column+1 < len((*grid)[row]) && value >= (*grid)[row][column+1] {
		return false
	}
	return true
}
