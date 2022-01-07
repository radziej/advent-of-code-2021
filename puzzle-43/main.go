package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type Vector3D [3]int

func main() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Current working directory: %v\n", workingDirectory)
	lines := readAllLines(workingDirectory + "/puzzle-43/input.txt")

	signedInts := regexp.MustCompile(`[-]?\d+`)
	cubes := make(map[Vector3D]bool)
	for _, line := range lines {
		// Read and convert ranges
		coordinateStrings := signedInts.FindAllString(line, 6)
		xRange := [2]int{parseOrDie(coordinateStrings[0]), parseOrDie(coordinateStrings[1])}
		yRange := [2]int{parseOrDie(coordinateStrings[2]), parseOrDie(coordinateStrings[3])}
		zRange := [2]int{parseOrDie(coordinateStrings[4]), parseOrDie(coordinateStrings[5])}
		// Ignoring ranges outside defined limits
		if xRange[0] < -50 || xRange[1] > 50 || yRange[0] < -50 || yRange[1] > 50 || zRange[0] < -50 || zRange[1] > 50 {
			fmt.Println("Skipping", line)
			continue
		}
		// Turn on/off cubes
		if line[0:2] == "on" {
			for x := xRange[0]; x <= xRange[1]; x++ {
				for y := yRange[0]; y <= yRange[1]; y++ {
					for z := zRange[0]; z <= zRange[1]; z++ {
						cubes[Vector3D{x, y, z}] = true
					}
				}
			}
		} else if line[0:3] == "off" {
			for x := xRange[0]; x <= xRange[1]; x++ {
				for y := yRange[0]; y <= yRange[1]; y++ {
					for z := zRange[0]; z <= zRange[1]; z++ {
						if _, ok := cubes[Vector3D{x, y, z}]; ok {
							delete(cubes, Vector3D{x, y, z})
						}
					}
				}
			}
		} else {
			panic(fmt.Sprintf("unrecognized line format: %v", line))
		}
	}

	fmt.Println("Total of", len(cubes), "cubes are turned on.")
}

func readAllLines(p string) []string {
	file, err := os.Open(p)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func parseOrDie(s string) int {
	number, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	return number
}
