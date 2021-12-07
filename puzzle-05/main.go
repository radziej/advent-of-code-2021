package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Each line, i.e., diagnostic report consists of 12 binary digits
	zeroes := [12]int{}
	ones := [12]int{}

	//fmt.Println(wd)
	//movements := []Movement{{"down", 10}, {"down", 10}, {"up", 5}, {"forward", 10}}
	for line := range readLines(workingDirectory + "/puzzle-05/input.txt") {
		countColumns(&zeroes, &ones, line)
	}
	fmt.Printf("Zeroes: %v\n", zeroes)
	fmt.Printf("Ones: %v\n", ones)
	//fmt.Println(bitsToInteger([]int{0, 0, 0, 1}))
	//fmt.Println(bitsToInteger([]int{1, 0, 0, 1}))
	gamma, epsilon := determineGammaEpsilon(&zeroes, &ones)
	fmt.Printf("Gamma, Epsilon: %v, %v\n", gamma, epsilon)
	fmt.Printf("Gamma * Epsilon: %v\n", gamma*epsilon)
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

func countColumns(zeroes *[12]int, ones *[12]int, digits string) {
	for position, character := range digits {
		if string(character) == "0" {
			(*zeroes)[position]++
		} else if string(character) == "1" {
			(*ones)[position]++
		} else {
			log.Fatalf("unrecognized bit: %v\n", string(character))
		}
	}
}

func determineGammaEpsilon(zeroes *[12]int, ones *[12]int) (int, int) {
	var gammaBinary [12]int
	for i := 0; i < len(zeroes); i++ {
		if zeroes[i] <= ones[i] {
			gammaBinary[i] = 1
		} else {
			gammaBinary[i] = 0
		}
	}

	var epsilonBinary [12]int
	for i, value := range gammaBinary {
		if value == 1 {
			epsilonBinary[i] = 0
		} else {
			epsilonBinary[i] = 1
		}
	}
	return bitsToInteger(gammaBinary[:]), bitsToInteger(epsilonBinary[:])
}

func bitsToInteger(bits []int) int {
	result := 0
	for i := 0; i < len(bits); i++ {
		result += bits[len(bits)-1-i] << i
	}
	return result
}
