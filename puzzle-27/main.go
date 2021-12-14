package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

func main() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Current working directory: %v\n", workingDirectory)

	var template []string
	rules := make(map[[2]string]string)
	toggle := true
	for line := range readLines(workingDirectory + "/puzzle-27/input.txt") {
		if line == "" {
			toggle = false
			continue
		}
		if toggle {
			// Read in template
			for _, character := range line {
				template = append(template, string(character))
			}
		} else {
			// Read in rules
			rule := strings.SplitN(line, " -> ", 2)
			rules[[2]string{string(rule[0][0]), string(rule[0][1])}] = rule[1]
		}
	}
	fmt.Println(template)
	fmt.Println(rules)

	for step := 1; step <= 10; step++ {
		splice := template
		for index := len(template) - 2; index >= 0; index-- {
			pair := [2]string{template[index], template[index+1]}
			if value, ok := rules[pair]; ok {
				//fmt.Println(splice[:index+2])
				//fmt.Println(splice[index+1:])
				splice = append(splice[:index+2], splice[index+1:]...)
				splice[index+1] = value
				//fmt.Println(splice)
			}
		}
		template = splice
	}

	count := make(map[string]int)
	for _, value := range template {
		if _, ok := count[value]; ok {
			count[value]++
		} else {
			count[value] = 1
		}
	}

	min := math.MaxInt
	max := 0
	for _, value := range count {
		if min > value {
			min = value
		}
		if max < value {
			max = value
		}
	}
	fmt.Printf("Max - min: %v - %v = %v", max, min, max-min)
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
