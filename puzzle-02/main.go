package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	data := parseFile(wd + "/puzzle-02/input.txt")
	//fmt.Println(data)
	//fmt.Println(slidingWindowSum([]int{0, 1, 2, 3, 4}, 3))
	fmt.Println(countIncreases(slidingWindowSum(data, 3)))
}

func parseFile(p string) []int {
	file, err := os.Open(p)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var data []int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		number, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, number)
	}

	return data
}

func sum(data []int) int {
	sum := 0
	for _, value := range data {
		sum += value
	}
	return sum
}

func slidingWindowSum(data []int, size int) []int {
	var windows []int
	for i := size; i <= len(data); i++ {
		windows = append(windows, sum(data[i-size:i]))
	}
	return windows
}

func countIncreases(data []int) int {
	increases := 0
	previousValue := data[0]
	for _, value := range data {
		if value > previousValue {
			increases++
		}
		previousValue = value
	}
	return increases
}
