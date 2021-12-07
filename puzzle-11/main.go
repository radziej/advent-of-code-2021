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

	var fishes []int
	lines := readLines(workingDirectory + "/puzzle-11/input.txt")
	for _, token := range strings.Split(lines[0], ",") {
		number, err := strconv.Atoi(token)
		if err != nil {
			log.Fatal(err)
		}
		fishes = append(fishes, number)
	}
	fmt.Printf("day 0, fishes: %v\n", len(fishes))

	for day := 0; day < 80; day++ {
		newFishes := 0
		for i := 0; i < len(fishes); i++ {
			if fishes[i] > 0 {
				fishes[i]--
			} else {
				fishes[i] = 6
				newFishes++
			}
		}
		moreFishes := make([]int, newFishes)
		for i := range moreFishes {
			moreFishes[i] = 8
		}
		fishes = append(fishes, moreFishes...)
	}
	fmt.Printf("final number of fishes: %v\n", len(fishes))

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

func advanceTimer(timer int) int {
	timer--
	if timer < 0 {
		return 6
	}
	return timer
}
