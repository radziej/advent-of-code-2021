package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var signals [][]string
	var outputs [][]string
	for line := range readLines(workingDirectory + "/puzzle-15/input.txt") {
		parts := strings.Split(line, "|")
		signals = append(signals, strings.Fields(parts[0]))
		outputs = append(outputs, strings.Fields(parts[1]))
	}

	segmentFrequency := map[int]int{
		1: 0,
		2: 0,
		3: 0,
		4: 0,
		5: 0,
		6: 0,
		7: 0,
	}
	for _, output := range outputs {
		for _, pattern := range output {
			segmentFrequency[len(pattern)]++
		}
	}
	fmt.Printf("Ones   -> 2 Segments: %v\n", segmentFrequency[2])
	fmt.Printf("Fours  -> 4 Segments: %v\n", segmentFrequency[4])
	fmt.Printf("Sevens -> 3 Segments: %v\n", segmentFrequency[3])
	fmt.Printf("Eights -> 7 Segments: %v\n", segmentFrequency[7])
	fmt.Printf("Sum of frequencies: %v\n", segmentFrequency[2]+segmentFrequency[4]+segmentFrequency[3]+segmentFrequency[7])
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
