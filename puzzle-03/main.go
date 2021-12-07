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
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(wd)
	//movements := []Movement{{"down", 10}, {"down", 10}, {"forward", 10}, {"up", 5}}
	movements := parseFile(wd + "/puzzle-03/input.txt")
	position := Position{0, 0}
	for _, movement := range movements {
		move(&position, &movement)
	}
	fmt.Printf("Final position: %v\n", position)
	fmt.Printf("Multiplied position: %v\n", position.x*position.y)
}

type Movement struct {
	direction string
	amount    int
}

type Position struct {
	x int
	y int
}

func parseFile(p string) []Movement {
	file, err := os.Open(p)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var data []Movement
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fields := strings.SplitN(scanner.Text(), " ", 2)

		number, err := strconv.Atoi(fields[1])
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, Movement{fields[0], number})
	}

	return data
}

func move(position *Position, movement *Movement) {
	if movement.direction == "forward" {
		position.x += movement.amount
	} else if movement.direction == "up" {
		position.y -= movement.amount
	} else if movement.direction == "down" {
		position.y += movement.amount
	} else {
		log.Fatalf("unrecognized direction: %v", movement.direction)
	}
}
