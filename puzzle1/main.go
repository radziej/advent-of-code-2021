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

	data := parseFile(wd + "/puzzle1/input.txt")
	//fmt.Println(data)
	//fmt.Println(countIncreases([]int{0, 1, 2, 3}))
	//fmt.Println(countIncreases([]int{0, 0, 0, 1, 0, 1}))
	fmt.Println(countIncreases(data))
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