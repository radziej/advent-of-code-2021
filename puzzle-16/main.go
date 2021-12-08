package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	var signals [][]string
	var outputs [][]string
	for line := range readLines(workingDirectory + "/puzzle-16/input.txt") {
		parts := strings.Split(line, "|")
		signals = append(signals, strings.Fields(parts[0]))
		outputs = append(outputs, strings.Fields(parts[1]))
	}

	//outputs = [][]string{{"cdfeb", "fcadb", "cdfeb", "cdbaf"}}
	//for i, signal := range [][]string{{"acedgfb", "cdfbe", "gcdfa", "fbcad", "dab", "cefabd", "cdfgeb", "eafb", "cagedb", "ab"}} {
	total := 0
	for i, signal := range signals {
		numberSets := translateSegments(signal)
		//for i := 0; i < 10; i++ {
		//	fmt.Printf("%v:", i)
		//	for key2, _ := range numberSets[i] {
		//		fmt.Printf(" %v", key2)
		//	}
		//	fmt.Println("")
		//}

		var digits []int
		for _, output := range outputs[i] {
			set := toSet(output)
			for number, referenceSet := range numberSets {
				if equal(set, referenceSet) {
					digits = append(digits, number)
					break
				}
			}
		}
		total += joinInts(digits)
	}
	fmt.Println(total)
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

func translateSegments(patterns []string) map[int]map[string]bool {
	translation := make(map[int]map[string]bool)

	// Phase 1 - identify unique patterns (1, 4, 7, 8)
	var remainder []string
	for _, pattern := range patterns {
		if len(pattern) == 2 {
			translation[1] = toSet(pattern)
		} else if len(pattern) == 3 {
			translation[7] = toSet(pattern)
		} else if len(pattern) == 4 {
			translation[4] = toSet(pattern)
		} else if len(pattern) == 7 {
			translation[8] = toSet(pattern)
		} else {
			remainder = append(remainder, pattern)
		}
	}

	// Phase 2 - identify unique overlaps (3, 6)
	patterns = remainder
	remainder = nil
	for _, pattern := range patterns {
		set := toSet(pattern)
		if len(pattern) == 5 && len(intersection(set, translation[1])) == 2 {
			translation[3] = set
		} else if len(pattern) == 6 && len(intersection(set, translation[1])) == 1 {
			translation[6] = set
		} else {
			remainder = append(remainder, pattern)
		}
	}

	// Phase 3 - identify remainder (0, 2, 5, 9)
	patterns = remainder
	remainder = nil
	for _, pattern := range patterns {
		set := toSet(pattern)
		if len(pattern) == 5 && len(intersection(set, translation[4])) == 2 {
			translation[2] = set
		} else if len(pattern) == 5 && len(intersection(set, translation[4])) == 3 {
			translation[5] = set
		} else if len(pattern) == 6 && len(intersection(set, translation[3])) == 4 {
			translation[0] = set
		} else if len(pattern) == 6 && len(intersection(set, translation[3])) == 5 {
			translation[9] = set
		} else {
			remainder = append(remainder, pattern)
		}
	}

	return translation
}

func toSet(s string) map[string]bool {
	result := make(map[string]bool)
	for _, character := range strings.Split(s, "") {
		result[character] = true
	}
	return result
}

func intersection(a map[string]bool, b map[string]bool) map[string]bool {
	result := make(map[string]bool)
	for key, _ := range a {
		if _, ok := b[key]; ok {
			result[key] = true
		}
	}
	return result
}

func equal(a map[string]bool, b map[string]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for key, _ := range a {
		if _, ok := b[key]; !ok {
			return false
		}
	}
	return true
}

func joinInts(s []int) int {
	result := 0
	power := 1
	for i := len(s) - 1; i >= 0; i-- {
		result += s[i] * power
		power *= 10
	}
	return result
}
