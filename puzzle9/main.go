package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type Segment struct {
	x1, y1, x2, y2 int
}

func isHorizontal(s Segment) bool {
	if s.y1 == s.y2 {
		return true
	} else {
		return false
	}
}

func isVertical(s Segment) bool {
	if s.x1 == s.x2 {
		return true
	} else {
		return false
	}
}

func main() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var segments []Segment
	pattern := regexp.MustCompile(`\d+`)
	for line := range readLines(workingDirectory + "/puzzle10/input.txt") {
		matches := pattern.FindAllString(line, -1)
		if len(matches) != 4 {
			for i, m := range matches {
				fmt.Printf("%v: %v\n", i, m)
			}
			log.Fatalf("invalid line format: %v\n", line)
		}

		numbers := [4]int{}
		for i, match := range matches {
			number, err := strconv.Atoi(match)
			if err != nil {
				log.Fatal(err)
			}
			numbers[i] = number
		}
		segments = append(segments, Segment{numbers[0], numbers[1], numbers[2], numbers[3]})
	}

	// Dynamic determination of grid size? Overkill.
	grid := [1000][1000]int{}
	for _, segment := range segments {
		if isHorizontal(segment) {
			for i := min(segment.x1, segment.x2); i <= max(segment.x1, segment.x2); i++ {
				grid[i][segment.y1]++
			}
		} else if isVertical(segment) {
			for i := min(segment.y1, segment.y2); i <= max(segment.y1, segment.y2); i++ {
				grid[segment.x1][i]++
			}
		} else {
			// By definition, remaining lines must be diagonal
			width := segment.x2 - segment.x1
			widthSign := sign(width)
			height := segment.y2 - segment.y1
			heightSign := sign(height)
			for i := 0; i <= max(widthSign*width, heightSign*height); i++ {
				grid[segment.x1+widthSign*i][segment.y1+heightSign*i]++
			}
		}
	}

	fmt.Printf("Grid points with 2 or more overlaps: %v\n", aboveThreshold(grid, 2))
}

func sign(a int) int {
	if a > 0 {
		return 1
	} else if a == 0 {
		return 0
	} else {
		return -1
	}
}

func min(a, b int) int {
	if a <= b {
		return a
	} else {
		return b
	}
}

func max(a, b int) int {
	if a >= b {
		return a
	} else {
		return b
	}
}

func aboveThreshold(grid [1000][1000]int, threshold int) int {
	count := 0
	for _, row := range grid {
		for _, field := range row {
			if field >= threshold {
				count++
			}
		}
	}
	return count
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
