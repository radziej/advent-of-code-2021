package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Chunk struct {
	Opener   string
	Closer   string
	Parent   *Chunk
	Children []*Chunk
}

var openers = map[string]bool{"(": true, "[": true, "{": true, "<": true}
var closers = map[string]string{")": "(", "]": "[", "}": "{", ">": "<"}
var scoring = map[string]int{")": 3, "]": 57, "}": 1197, ">": 25137}

func main() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var corruptions []int
	for line := range readLines(workingDirectory + "/puzzle-19/input.txt") {
		var currentChunk *Chunk
		for _, character := range strings.Split(line, "") {
			if _, ok := openers[character]; ok {
				// Create new chunk for opening bracket
				if currentChunk != nil {
					newChunk := &Chunk{character, "", currentChunk, []*Chunk{}}
					currentChunk.Children = append(currentChunk.Children, newChunk)
					currentChunk = newChunk
				} else {
					currentChunk = &Chunk{character, "", nil, []*Chunk{}}
				}
			} else if opener, ok := closers[character]; ok {
				// Attempt to close chunk for closing bracket
				if currentChunk.Opener == opener {
					currentChunk.Closer = character
					currentChunk = currentChunk.Parent
				} else {
					corruptions = append(corruptions, scoring[character])
					break
				}
			}
		}
	}

	totalCorruption := 0
	for _, corruption := range corruptions {
		totalCorruption += corruption
	}
	fmt.Printf("Number of corrupted lines: %v\n", len(corruptions))
	fmt.Printf("Total corruption score: %v\n", totalCorruption)
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
