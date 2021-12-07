package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type BingoBoard struct {
	Fields  [5][5]int
	Matches [5][5]bool
}

func (b *BingoBoard) Feed(line string) error {
	rowIndex := -1
	for i, row := range b.Fields {
		if row[0] == 0 && row[1] == 0 {
			rowIndex = i
		}
	}
	if rowIndex == -1 {
		return errors.New("board is full, cannot feed line")
	}
	for i, element := range strings.Fields(line) {
		number, err := strconv.Atoi(element)
		if err != nil {
			log.Fatal(err)
		}
		b.Fields[rowIndex][i] = number
	}
	return nil
}

func (b BingoBoard) String() string {
	representation := ""
	for _, row := range b.Fields {
		representation += fmt.Sprintf("%v\n", row)
	}
	return representation
}

func (b *BingoBoard) Mark(draw int) {
	for i := 0; i < len(b.Fields); i++ {
		for j := 0; j < len(b.Fields[i]); j++ {
			if b.Fields[i][j] == draw {
				b.Matches[i][j] = true
			}
		}
	}
}

func (b *BingoBoard) IsVictorious() bool {
	// Rows
	for i := 0; i < len(b.Matches); i++ {
		matches := 0
		for j := 0; j < len(b.Matches[0]); j++ {
			if b.Matches[i][j] {
				matches++
			}
		}
		if matches == 5 {
			return true
		}
	}

	// Columns
	for j := 0; j < len(b.Matches[0]); j++ {
		matches := 0
		for i := 0; i < len(b.Matches); i++ {
			if b.Matches[i][j] {
				matches++
			}
		}
		if matches == 5 {
			return true
		}
	}
	return false
}

func (b *BingoBoard) Score() int {
	score := 0
	for i := 0; i < len(b.Fields); i++ {
		for j := 0; j < len(b.Fields[i]); j++ {
			if !b.Matches[i][j] {
				score += b.Fields[i][j]
			}
		}
	}
	return score
}

func main() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var draws []int
	var boards []BingoBoard

	numLines := -1
	currentBoard := BingoBoard{}
	for line := range readLines(workingDirectory + "/puzzle-07/input.txt") {
		numLines++
		if line == "" {
			continue
		}

		// First line is drawn numbers
		if numLines == 0 {
			for _, element := range strings.Split(line, ",") {
				number, err := strconv.Atoi(element)
				if err != nil {
					log.Fatal(err)
				}
				draws = append(draws, number)
			}
			continue
		}

		// Create a new board and fill up with 5 lines
		if currentBoard.Feed(line) != nil {
			boards = append(boards, currentBoard)
			currentBoard = BingoBoard{}
			err := currentBoard.Feed(line)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	//for _, draw := range []int{26, 85, 63, 25, 86} {
	for _, draw := range draws {
		for i := range boards {
			boards[i].Mark(draw)
			if boards[i].IsVictorious() {
				fmt.Printf("score: %v\n", boards[i].Score())
				fmt.Printf("draw: %v\n", draw)
				fmt.Printf("final score: %v\n", draw*boards[i].Score())
				os.Exit(0)
			}
		}
	}
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
