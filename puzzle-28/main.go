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

	template := make(map[[2]string]int)
	counts := make(map[string]int)
	rules := make(map[[2]string]string)
	toggle := true
	for line := range readLines(workingDirectory + "/puzzle-28/input.txt") {
		if line == "" {
			toggle = false
			continue
		}
		if toggle {
			// Read in template
			for i := 0; i < len(line)-1; i++ {
				pair := [2]string{string(line[i]), string(line[i+1])}
				upsertTemplate(&template, pair, 1)
				upsertCount(&counts, string(line[i]), 1)
			}
			// Special case because pair loop skips last element
			upsertCount(&counts, string(line[len(line)-1]), 1)
		} else {
			// Read in rules
			rule := strings.SplitN(line, " -> ", 2)
			rules[[2]string{string(rule[0][0]), string(rule[0][1])}] = rule[1]
		}
	}
	fmt.Println(template)
	fmt.Println(counts)
	fmt.Println(rules)

	for step := 1; step <= 40; step++ {
		splice := make(map[[2]string]int)
		for pair, pairCount := range template {
			if insert, ok := rules[pair]; ok {
				upsertTemplate(&splice, [2]string{pair[0], insert}, pairCount)
				upsertTemplate(&splice, [2]string{insert, pair[1]}, pairCount)
				upsertCount(&counts, insert, pairCount)
			} else {
				upsertTemplate(&splice, pair, pairCount)
			}
		}
		template = splice
	}
	//fmt.Println(template)

	//
	min := math.MaxInt
	max := 0
	for _, count := range counts {
		if min > count {
			min = count
		}
		if max < count {
			max = count
		}
	}
	//fmt.Println(math.MaxInt)
	//fmt.Println(2188189693529)
	//fmt.Println(max)

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

func upsertTemplate(m *map[[2]string]int, key [2]string, value int) {
	if _, ok := (*m)[key]; ok {
		(*m)[key] += value
	} else {
		(*m)[key] = value
	}
}

func upsertCount(m *map[string]int, key string, value int) {
	if _, ok := (*m)[key]; ok {
		(*m)[key] += value
	} else {
		(*m)[key] = value
	}
}
