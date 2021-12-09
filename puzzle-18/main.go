package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

func main() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var grid [][]int
	for line := range readLines(workingDirectory + "/puzzle-18/input.txt") {
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

	var basinSizes []int
	for row := range riskLevels {
		for column := range riskLevels[row] {
			if riskLevels[row][column] > 0 {
				// Low points are defined to be part of unique basins -> no filtering needed
				basinSizes = append(basinSizes, len(cartographBasin(&grid, Point{row, column}, []Point{})))
			}
		}
	}

	fmt.Printf("Number of basins %v\n", len(basinSizes))
	sort.Ints(basinSizes)
	product := 1
	for _, size := range basinSizes[len(basinSizes)-3:] {
		product *= size
	}
	fmt.Printf("Top 3 largest basins %v, product %v\n", basinSizes[len(basinSizes)-3:], product)
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

type Point struct {
	row, column int
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

func cartographBasin(grid *[][]int, point Point, result []Point) []Point {
	if contains(result, point) {
		// Not revisiting old points
		return result
	} else if (*grid)[point.row][point.column] >= 9 {
		// Points of height 9 are by definition not part of basins
		return result
	} else {
		result = append(result, point)
	}

	if point.row-1 >= 0 {
		result = cartographBasin(grid, Point{point.row - 1, point.column}, result)
	}
	if point.row+1 < len(*grid) {
		result = cartographBasin(grid, Point{point.row + 1, point.column}, result)
	}
	if point.column-1 >= 0 {
		result = cartographBasin(grid, Point{point.row, point.column - 1}, result)
	}
	if point.column+1 < len((*grid)[point.row]) {
		result = cartographBasin(grid, Point{point.row, point.column + 1}, result)
	}
	return result
}

func contains(s []Point, p Point) bool {
	for _, e := range s {
		if e == p {
			return true
		}
	}
	return false
}
