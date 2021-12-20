package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
)

type Point struct {
	Row    int
	Column int
}

func main() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Current working directory: %v\n", workingDirectory)

	var grid [][]int
	for line := range readLines(workingDirectory + "/puzzle-29/input.txt") {
		row := make([]int, len(line))
		for i, character := range line {
			number, err := strconv.Atoi(string(character))
			if err != nil {
				log.Fatal(err)
			}
			row[i] = number
		}
		grid = append(grid, row)
	}

	// Example data
	//grid = [][]int{
	//	{1, 1, 6, 3, 7, 5, 1, 7, 4, 2},
	//	{1, 3, 8, 1, 3, 7, 3, 6, 7, 2},
	//	{2, 1, 3, 6, 5, 1, 1, 3, 2, 8},
	//	{3, 6, 9, 4, 9, 3, 1, 5, 6, 9},
	//	{7, 4, 6, 3, 4, 1, 7, 1, 1, 1},
	//	{1, 3, 1, 9, 1, 2, 8, 1, 3, 7},
	//	{1, 3, 5, 9, 9, 1, 2, 4, 2, 1},
	//	{3, 1, 2, 5, 4, 2, 1, 6, 3, 9},
	//	{1, 2, 9, 3, 1, 3, 8, 5, 2, 1},
	//	{2, 3, 1, 1, 9, 4, 4, 5, 8, 1},
	//}

	//for _, row := range grid {
	//	fmt.Println(row)
	//}

	//path := recursiveSearch(grid, Point{len(grid) - 1, len(grid[0]) - 1}, []Point{{0, 0}}, []Point{})
	//fmt.Println(path)
	//for _, step := range path {
	//	fmt.Printf(" %v", grid[step.Row][step.Column])
	//}
	//fmt.Println("")
	//fmt.Println("Total risk of path", pathRisk(grid, path))
	lowestCost := dynamicSearch(grid, Point{len(grid) - 1, len(grid[0]) - 1})
	fmt.Println("Optimal cost:", lowestCost)
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

func dynamicSearch(grid [][]int, target Point) int {
	height := len(grid)
	width := len(grid[0])
	costs := make([][]int, height)
	unvisited := make(map[Point]bool)
	for i := range grid {
		costs[i] = make([]int, width)
		for j := range costs[i] {
			costs[i][j] = math.MaxInt
			unvisited[Point{i, j}] = true
		}
	}
	costs[0][0] = 0

	for len(unvisited) > 0 {
		// Find node with the lowest tentative cost
		point := target
		pointCost := math.MaxInt64
		for p, _ := range unvisited {
			if costs[p.Row][p.Column] < pointCost {
				pointCost = costs[p.Row][p.Column]
				point = p
			}
		}

		// Left
		row := point.Row
		column := point.Column
		cost := costs[row][column]
		if row-1 >= 0 && costs[row-1][column] > cost+grid[row-1][column] {
			costs[row-1][column] = cost + grid[row-1][column]
		}
		// Right
		if row+1 < height && costs[row+1][column] > cost+grid[row+1][column] {
			costs[row+1][column] = cost + grid[row+1][column]
		}
		// Up
		if column-1 >= 0 && costs[row][column-1] > cost+grid[row][column-1] {
			costs[row][column-1] = cost + grid[row][column-1]
		}
		// Down
		if column+1 < width && costs[row][column+1] > cost+grid[row][column+1] {
			costs[row][column+1] = cost + grid[row][column+1]
		}
		delete(unvisited, point)
	}

	return costs[target.Row][target.Column]
}

func recursiveSearch(grid [][]int, target Point, path []Point, shortestPath []Point) []Point {
	position := path[len(path)-1]
	risk := pathRisk(grid, path)
	shortestRisk := pathRisk(grid, shortestPath)

	// Success
	if position == target {
		if risk >= shortestRisk {
			return shortestPath
		} else {
			fmt.Printf("Found shorter path with risk of: %v\n", risk)
			return path
		}
	}

	// Early return with too high risk?
	if risk >= shortestRisk {
		return shortestPath
	}

	// Figure out available pathing options
	var options []Point
	for _, option := range []Point{{position.Row - 1, position.Column}, {position.Row + 1, position.Column}, {position.Row, position.Column - 1}, {position.Row, position.Column + 1}} {
		valid := true
		// Must remain within grid bounds
		if option.Row < 0 || option.Row >= len(grid) || option.Column < 0 || option.Column >= len(grid[0]) {
			valid = false
			continue
		}
		for _, step := range path[:len(path)-1] {
			// Loops only increase risk
			if step == option {
				valid = false
				break
			}
			// Adjacent paths could have been reached earlier, thus only increase risk
			if step.Row-1 == option.Row && step.Column == option.Column ||
				step.Row+1 == option.Row && step.Column == option.Column ||
				step.Row == option.Row && step.Column-1 == option.Column ||
				step.Row == option.Row && step.Column+1 == option.Column {
				valid = false
				break
			}
		}
		if valid {
			options = append(options, option)
		}
	}

	// Prioritize options by distance
	sort.Slice(options, func(i, j int) bool {
		iDistance := math.Sqrt(math.Pow(float64(options[i].Row-target.Row), 2) + math.Pow(float64(options[i].Column-target.Column), 2))
		jDistance := math.Sqrt(math.Pow(float64(options[j].Row-target.Row), 2) + math.Pow(float64(options[j].Column-target.Column), 2))
		return iDistance < jDistance
	})

	for _, option := range options {
		branch := make([]Point, len(path)+1)
		copy(branch, path)
		branch[len(branch)-1] = option
		shortestPath = recursiveSearch(grid, target, branch, shortestPath)
	}
	return shortestPath
}

func pathRisk(grid [][]int, path []Point) int {
	if len(path) == 0 {
		return math.MaxInt
	}
	totalRisk := 0
	for _, step := range path[1:] {
		totalRisk += grid[step.Row][step.Column]
	}
	return totalRisk
}
