package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

// Cuboid Three pairs of coordinates [x1, x2], [y1, y2], [z1, z2] defining the ranges of a cube
type Cuboid [3][2]int

func (c Cuboid) Cubes() int {
	return (c[0][1] - c[0][0]) * (c[1][1] - c[1][0]) * (c[2][1] - c[2][0])
}

func main() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Current working directory: %v\n", workingDirectory)
	//lines := readAllLines(workingDirectory + "/puzzle-44/input.txt")
	lines := readAllLines(workingDirectory + "/puzzle-44/input.txt")

	//fmt.Println(Cuboid{{0, 10}, {0, 10}, {0, 10}}.Cubes())
	//fmt.Println(Cuboid{{-20, -10}, {-20, -10}, {-20, -10}}.Cubes())

	var cuboids []Cuboid
	signedInts := regexp.MustCompile(`[-]?\d+`)
	for _, line := range lines {
		// Parse cuboid
		coordinateStrings := signedInts.FindAllString(line, 6)
		newCuboid := Cuboid{
			// Careful here; intervals include both bounds, so we must extend one of them
			{parseOrDie(coordinateStrings[0]), parseOrDie(coordinateStrings[1]) + 1},
			{parseOrDie(coordinateStrings[2]), parseOrDie(coordinateStrings[3]) + 1},
			{parseOrDie(coordinateStrings[4]), parseOrDie(coordinateStrings[5]) + 1},
		}
		//// Ignoring ranges outside defined limits
		//if newCuboid[0][0] < -50 || newCuboid[0][1] > 50 || newCuboid[1][0] < -50 || newCuboid[1][1] > 50 || newCuboid[2][0] < -50 || newCuboid[2][1] > 50 {
		//	//fmt.Println("Skipping", line)
		//	continue
		//}

		// Ensure new cuboid is covering unique volume by removing overlaps
		var nextCuboids []Cuboid
		for _, cuboid := range cuboids {
			if isOverlapping(cuboid, newCuboid) {
				overlap := findOverlap(cuboid, newCuboid)
				//fmt.Println(cuboid, "overlaps with", newCuboid, "at", overlap)
				cuboidSplits := splitAndRemove(cuboid, overlap)
				//fmt.Println(len(cuboidSplits), "splits:", cuboidSplits)
				if sumCubes(cuboidSplits)+overlap.Cubes() != cuboid.Cubes() {
					fmt.Println(sumCubes(cuboidSplits), "+", overlap.Cubes(), "=", sumCubes(cuboidSplits)+overlap.Cubes(), "vs", cuboid.Cubes())
				}
				for _, cuboidSplit := range cuboidSplits {
					nextCuboids = append(nextCuboids, cuboidSplit)
				}
			} else {
				nextCuboids = append(nextCuboids, cuboid)
			}
		}

		// Turn on/off cuboids
		if line[0:2] == "on" {
			nextCuboids = append(nextCuboids, newCuboid)
		} else if line[0:3] == "off" {
			// Nothing to do here
		} else {
			panic(fmt.Sprintf("unrecognized line format: %v", line))
		}
		cuboids = nextCuboids

		//fmt.Println("Total of", sumCubes(cuboids), "cubes in", len(cuboids), "cuboids are turned on.")
	}

	fmt.Println("Total amount of", sumCubes(cuboids), "cubes in", len(cuboids), "cuboids.")
}

func sumCubes(cuboids []Cuboid) int {
	cubes := 0
	for _, c := range cuboids {
		cubes += c.Cubes()
	}
	return cubes
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

func isOverlapping(a, b Cuboid) bool {
	for axis := range a {
		if a[axis][0] > b[axis][1] || a[axis][1] < b[axis][0] {
			return false
		}
	}
	return true
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Returns the overlap between two cuboids, assuming that they do have an overlap
func findOverlap(a, b Cuboid) Cuboid {
	return Cuboid{
		{max(a[0][0], b[0][0]), min(a[0][1], b[0][1])},
		{max(a[1][0], b[1][0]), min(a[1][1], b[1][1])},
		{max(a[2][0], b[2][0]), min(a[2][1], b[2][1])},
	}
}

func splitAndRemove(from Cuboid, remove Cuboid) []Cuboid {
	// Expecting 1 to 3 bounds for each of the 3 axis
	var bounds [3][][2]int
	for axis := range from {
		// Bound to the left of removal cuboid
		if from[axis][0] != remove[axis][0] {
			bounds[axis] = append(bounds[axis], [2]int{from[axis][0], remove[axis][0]})
		}
		// Bounds of the removal cuboid
		bounds[axis] = append(bounds[axis], remove[axis])
		// Bound to the right of removal cuboid
		if from[axis][1] != remove[axis][1] {
			bounds[axis] = append(bounds[axis], [2]int{remove[axis][1], from[axis][1]})
		}
	}

	// Construct cuboids from sets of bounds
	var cuboids []Cuboid
	for _, xBound := range bounds[0] {
		for _, yBound := range bounds[1] {
			for _, zBound := range bounds[2] {
				if cuboid := (Cuboid{xBound, yBound, zBound}); cuboid != remove {
					cuboids = append(cuboids, cuboid)
				}
			}
		}
	}
	return cuboids
}
