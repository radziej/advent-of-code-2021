package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

var adjacency = [...]int{-1, 0, 1}

func main() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var energies [10][10]int
	row := 0
	for line := range readLines(workingDirectory + "/puzzle-21/input.txt") {
		for column, character := range line {
			number, err := strconv.Atoi(string(character))
			if err != nil {
				log.Fatal(err)
			}
			energies[row][column] = number
		}
		row++
	}

	totalFlashes := 0
	for step := 1; step <= 100; step++ {
		var flashed [10][10]bool
		for row := range energies {
			for column := range energies[row] {
				raiseEnergy(&energies, &flashed, row, column)
			}
		}
		for row := range flashed {
			for column := range flashed[row] {
				if flashed[row][column] {
					totalFlashes++
				}
			}
		}

		if step%10 == 0 {
			fmt.Printf("Total flashes after %v steps: %v\n", step, totalFlashes)
		}
	}

	fmt.Printf("Total flashes: %v\n", totalFlashes)
	//printState(energies)
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

func printState(energies [10][10]int) {
	for _, row := range energies {
		for _, field := range row {
			fmt.Printf("%v", field)
		}
		fmt.Println("")
	}
}

func raiseEnergy(energies *[10][10]int, flashed *[10][10]bool, row int, column int) {
	// Octopi only flash once
	if flashed[row][column] {
		return
	}
	// Raise energy level
	if energies[row][column] < 9 {
		(*energies)[row][column]++
	} else {
		(*energies)[row][column] = 0
		(*flashed)[row][column] = true
		// Raise adjacent energy levels
		for _, rowOffset := range adjacency {
			for _, columnOffset := range adjacency {
				if row+rowOffset >= 0 && row+rowOffset < 10 && column+columnOffset >= 0 && column+columnOffset < 10 {
					raiseEnergy(energies, flashed, row+rowOffset, column+columnOffset)
				}
			}
		}
	}
}
