package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func main() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var positions []int
	lines := readLines(workingDirectory + "/puzzle13/input.txt")
	for _, token := range strings.Split(lines[0], ",") {
		number, err := strconv.Atoi(token)
		if err != nil {
			log.Fatal(err)
		}
		positions = append(positions, number)
	}

	costs := make(map[int]int)
	for p := min(positions); p <= max(positions); p++ {
		cost := 0
		for _, position := range positions {
			cost += abs(p - position)
		}
		costs[p] = cost
	}

	optimalPosition, optimalCost := -1, math.MaxInt64
	for position, cost := range costs {
		if cost < optimalCost {
			optimalPosition = position
			optimalCost = cost
		}
	}
	fmt.Printf("Optimal position %v at cost of %v fuel", optimalPosition, optimalCost)
}

func readLines(p string) []string {
	file, err := os.Open(p)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func max(values []int) int {
	maximum := values[0]
	for _, value := range values[1:] {
		if maximum < value {
			maximum = value
		}
	}
	return maximum
}

func min(values []int) int {
	minimum := values[0]
	for _, value := range values[1:] {
		if minimum > value {
			minimum = value
		}
	}
	return minimum
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}
