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

	var coordinates [][2]int
	var instructions [][2]int
	//example := []string{"start-A", "start-b", "A-c", "A-b", "b-d", "A-end", "b-end"}
	//for _, line := range example {
	pointToggle := true
	for line := range readLines(workingDirectory + "/puzzle-26/input.txt") {
		if line == "" {
			pointToggle = false
			continue
		}
		if pointToggle {
			// Read in coordinates
			numbers := strings.SplitN(line, ",", 2)
			var coordinate [2]int
			for i, n := range numbers {
				number, err := strconv.Atoi(n)
				if err != nil {
					log.Fatal(err)
				}
				coordinate[i] = number
			}
			coordinates = append(coordinates, coordinate)
		} else {
			// Read in instructions
			var instruction [2]int
			if string(line[11]) == "x" {
				instruction[0] = 0
			} else if string(line[11]) == "y" {
				instruction[0] = 1
			} else {
				log.Fatalf("unknown folding direction: %v", line)
			}
			number, err := strconv.Atoi(line[13:])
			if err != nil {
				log.Fatal(err)
			}
			instruction[1] = number
			instructions = append(instructions, instruction)
		}
	}
	//fmt.Println(coordinates)
	fmt.Println(instructions)

	xMax, yMax := 0, 0
	for _, coordinate := range coordinates {
		if coordinate[0] >= xMax {
			xMax = coordinate[0] + 1
		}
		if coordinate[1] >= yMax {
			yMax = coordinate[1] + 1
		}
	}

	grid := make([][]int, yMax)
	for i := range grid {
		grid[i] = make([]int, xMax)
	}
	for _, coordinates := range coordinates {
		// Note that y -> row, x -> column
		grid[coordinates[1]][coordinates[0]] = 1
	}
	fmt.Printf("Initial dots: %v\n", len(coordinates))
	fmt.Printf("Width: %v, Height: %v\n", len(grid[0]), len(grid))

	folded := &grid
	for _, instruction := range instructions {
		if instruction[0] == 0 {
			folded = foldHorizontally(folded)
		} else if instruction[0] == 1 {
			folded = foldVertically(folded)
		}
	}
	fmt.Printf("Width: %v, height: %v\n", len((*folded)[0]), len(*folded))

	//totalDots := 0
	for _, row := range *folded {
		for _, field := range row {
			if field > 0 {
				fmt.Print("#")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Println("")
	}
	//fmt.Printf("Total dots after first fold: %v\n", totalDots)
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

func foldHorizontally(grid *[][]int) *[][]int {
	folded := make([][]int, len(*grid))
	halfWidth := len((*grid)[0]) / 2
	for i := range folded {
		folded[i] = make([]int, halfWidth)
	}

	for row := range folded {
		for i := 0; i < len(folded[row]); i++ {
			j := len((*grid)[row]) - 1 - i
			folded[row][i] = (*grid)[row][i] + (*grid)[row][j]
		}
	}

	return &folded
}

func foldVertically(grid *[][]int) *[][]int {
	folded := make([][]int, len(*grid)/2)
	for i := range folded {
		folded[i] = make([]int, len((*grid)[i]))
	}

	for column := range folded[0] {
		for i := 0; i < len(folded); i++ {
			j := len((*grid)) - 1 - i
			folded[i][column] = (*grid)[i][column] + (*grid)[j][column]
		}
	}

	return &folded
}
