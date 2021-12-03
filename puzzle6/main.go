package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Each line, i.e., diagnostic report consists of 12 binary digits
	var readings [][12]int
	for line := range readLines(workingDirectory + "/puzzle6/input.txt") {
		var reading [12]int
		for pos, char := range line {
			reading[pos], err = strconv.Atoi(string(char))
			if err != nil {
				log.Fatal(err)
			}
		}
		readings = append(readings, reading)
	}

	oxygenReadings := readings
	for position := 0; position < 12; position++ {
		if len(oxygenReadings) < 2 {
			break
		}

		counts := count(oxygenReadings, position)
		if counts[1] >= counts[0] {
			oxygenReadings = filterForDigit(oxygenReadings, 1, position)
		} else {
			oxygenReadings = filterForDigit(oxygenReadings, 0, position)
		}
	}
	fmt.Println(oxygenReadings)

	carbondioxideReadings := readings
	for position := 0; position < 12; position++ {
		if len(carbondioxideReadings) < 2 {
			break
		}

		counts := count(carbondioxideReadings, position)
		if counts[1] < counts[0] {
			carbondioxideReadings = filterForDigit(carbondioxideReadings, 1, position)
		} else {
			carbondioxideReadings = filterForDigit(carbondioxideReadings, 0, position)
		}
	}
	fmt.Println(carbondioxideReadings)

	oxygen := bitsToInteger(oxygenReadings[0][:])
	carbondioxide := bitsToInteger(carbondioxideReadings[0][:])
	fmt.Printf("Oxygen: %v\n", oxygen)
	fmt.Printf("Carbon dioxide: %v\n", carbondioxide)
	fmt.Printf("Life support: %v\n", oxygen * carbondioxide)
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

func count(readings [][12]int, position int) map[int]int {
	counts := make(map[int]int)
	for _, reading := range readings {
		counts[reading[position]]++
	}
	return counts
}

func filterForDigit(readings [][12]int, digit int, position int) [][12]int {
	var result [][12]int
	for _, reading := range readings {
		if reading[position] == digit {
			result = append(result, reading)
		}
	}
	return result
}

func bitsToInteger(bits []int) int {
	result := 0
	for i := 0; i < len(bits); i++ {
		result += bits[len(bits) - 1 - i] << i
	}
	return result
}