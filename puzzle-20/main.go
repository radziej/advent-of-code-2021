package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type Chunk struct {
	Opener   string
	Closer   string
	Parent   *Chunk
	Children []*Chunk
}

var openers = map[string]string{"(": ")", "[": "]", "{": "}", "<": ">"}
var closers = map[string]string{")": "(", "]": "[", "}": "{", ">": "<"}
var scoring = map[string]int{")": 1, "]": 2, "}": 3, ">": 4}

func main() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var rows [][]*Chunk
	//fmt.Println(workingDirectory)
	//for _, line := range []string{"[({(<(())[]>[[{[]{<()<>>", "[(()[<>])]({[<{<<[]>>(", "(((({<>}<{<{<>}{[]{[]{}", "{<[[]]>}<{[{[{[]{()[[[]", "<{([{{}}[<[[[<>{}]]]>[]]"} {
	for line := range readLines(workingDirectory + "/puzzle-19/input.txt") {
		rows = append(rows, []*Chunk{})
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
					rows[len(rows)-1] = append(rows[len(rows)-1], currentChunk)
				}
			} else if opener, ok := closers[character]; ok {
				// Attempt to close chunk for closing bracket
				if currentChunk.Opener == opener {
					currentChunk.Closer = character
					currentChunk = currentChunk.Parent
				} else {
					// Remove corrupted row
					rows = rows[:len(rows)-1]
					break
				}
			}
		}
	}

	var scores []int
	for _, row := range rows {
		currentChunk := row[len(row)-1]
		var missingClosers []string
		for {
			if currentChunk.Closer == "" {
				missingClosers = append(missingClosers, openers[currentChunk.Opener])
			}
			if currentChunk.Closer != "" || len(currentChunk.Children) == 0 {
				break
			} else {
				currentChunk = currentChunk.Children[len(currentChunk.Children)-1]
			}
		}

		score := 0
		for i := len(missingClosers) - 1; i >= 0; i-- {
			//fmt.Print(missingClosers[i])
			score = score*5 + scoring[missingClosers[i]]
		}
		//fmt.Println("")
		scores = append(scores, score)
	}

	fmt.Printf("Number of scores: %v\n", len(scores))
	sort.Ints(scores)
	//fmt.Printf("Scores: %v\n", scores)
	fmt.Printf("Median score: %v\n", scores[len(scores)/2])
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
