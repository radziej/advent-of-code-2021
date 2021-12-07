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

	fishes := map[int]uint{
		0: 0,
		1: 0,
		2: 0,
		3: 0,
		4: 0,
		5: 0,
		6: 0,
		7: 0,
		8: 0,
	}
	lines := readLines(workingDirectory + "/puzzle12/input.txt")
	for _, token := range strings.Split(lines[0], ",") {
		number, err := strconv.Atoi(token)
		if err != nil {
			log.Fatal(err)
		}
		fishes[number]++
	}
	fmt.Printf("day 0, fishes %v\n", sum(fishes))

	for day := 1; day <= 256; day++ {
		resettingFishes := fishes[0]
		for time := 0; time < len(fishes)-1; time++ {
			fishes[time] = fishes[time+1]
		}
		// Resetting fishes with 0 timer and adding new spawns
		fishes[6] += resettingFishes
		fishes[8] = resettingFishes
		if day%10 == 0 {
			fmt.Printf("day %v, fishes %v\n", day, sum(fishes))
		}
	}
	fmt.Printf("final number of fishes: %v\n", sum(fishes))
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

func sum(mapping map[int]uint) uint {
	var total uint = 0
	for _, value := range mapping {
		total += value
	}
	return total
}
