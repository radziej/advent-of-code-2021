package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Vector [3]int

func (v Vector) Add(addend Vector) Vector {
	return Vector{v[0] + addend[0], v[1] + addend[1], v[2] + addend[2]}
}

func (v Vector) Sub(subtrahend Vector) Vector {
	return Vector{v[0] - subtrahend[0], v[1] - subtrahend[1], v[2] - subtrahend[2]}
}

type Matrix [3][3]int

func (m Matrix) Dot(v Vector) Vector {
	return Vector{
		m[0][0]*v[0] + m[0][1]*v[1] + m[0][2]*v[2],
		m[1][0]*v[0] + m[1][1]*v[1] + m[1][2]*v[2],
		m[2][0]*v[0] + m[2][1]*v[1] + m[2][2]*v[2],
	}
}

type Graph struct {
	Name  string
	Nodes map[Vector]bool
	Edges map[Vector]map[Vector]int
}

func (g Graph) AllNodes() []Vector {
	keys := make([]Vector, len(g.Nodes))
	i := 0
	for k, _ := range g.Nodes {
		keys[i] = k
		i++
	}
	return keys
}

//func (g Graph) String() string {
//	s := fmt.Sprintf()
//	for _, e
//}

func NewGraph() Graph {
	return Graph{"", make(map[Vector]bool), make(map[Vector]map[Vector]int)}
}

func main() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Current working directory: %v\n", workingDirectory)

	var scannersRaw []Graph
	var currentScanner Graph
	//for line := range readLines(workingDirectory + "/puzzle-37/test.txt") {
	for line := range readLines(workingDirectory + "/puzzle-37/input.txt") {
		if line == "" {
			continue
		}
		if line[:3] == "---" {
			if currentScanner.Nodes != nil {
				scannersRaw = append(scannersRaw, currentScanner)
			}
			fmt.Println("Parsing", line)
			currentScanner = NewGraph()
			currentScanner.Name = line
			continue
		}

		vec := Vector{0, 0, 0}
		for i, field := range strings.Split(line, ",") {
			number, err := strconv.Atoi(field)
			if err != nil {
				log.Fatal(err)
			}
			vec[i] = number
		}
		currentScanner.Nodes[vec] = true
	}
	scannersRaw = append(scannersRaw, currentScanner)
	//fmt.Println(scannersRaw[0])

	// Calculate edges of fully connected graph
	for i := range scannersRaw {
		edges := fullyConnectedEdges(scannersRaw[i].AllNodes())
		scannersRaw[i].Edges = edges
	}
	//fmt.Println(scannersRaw[0].Nodes)
	//n1 := scannersRaw[0].AllNodes()[0]
	//n2 := scannersRaw[0].AllNodes()[1]
	//fmt.Println(scannersRaw[0].Edges[n1][n2], "vs", scannersRaw[0].Edges[n2][n1])

	// Defining first beacon as origin of coordinate system
	var scanners []Graph
	scanners = append(scanners, scannersRaw[0])

	// Attempt transformation of scanner to unified coordinate system until all are transformed
	remainingScanners := scannersRaw[1:]
	rawIndex := 0
	// Also collect scanner origins
	var scannerPositions []Vector
	for len(remainingScanners) > 0 {
		rawScanner := remainingScanners[rawIndex]
		for _, scanner := range scanners {
			// Find overlapping scanner (in unified coordinate system) by matching node distances
			overlap := make(map[Vector]Vector)
			for node2, links2 := range scanner.Edges {
				for node1, links1 := range rawScanner.Edges {
					matches := findMatches(links1, links2)
					//fmt.Println(len(matches))
					// Assuming a matching pair of nodes, there must be at least 11 additional overlapping nodes
					if len(matches) >= 11 {
						overlap[node1] = node2
					}
				}
			}

			// Transform scanner to unified coordinate system
			if len(overlap) >= 12 {
				fmt.Println("Found", len(overlap), "overlapping beacons between", rawScanner.Name, "and", scanner.Name)
				//fmt.Println(overlap)
				// Need vector pair with different component magnitudes to determine rotation
				rawVectors := make([]Vector, len(overlap))
				i := 0
				for vec, _ := range overlap {
					rawVectors[i] = vec
					i++
				}
				rawVector1, rawVector2 := differenceWithUniqueComponents(rawVectors)
				transformation := determineTransformation(rawVector1.Sub(rawVector2), overlap[rawVector1].Sub(overlap[rawVector2]))
				//fmt.Println("Transformation", transformation)
				origin := overlap[rawVector1].Sub(transformation.Dot(rawVector1))
				scannerPositions = append(scannerPositions, origin)
				//fmt.Println("New origin at", origin)

				// Construct scanner in unified coordinate system
				unifiedScanner := NewGraph()
				unifiedScanner.Name = rawScanner.Name
				for _, node := range rawScanner.AllNodes() {
					unifiedScanner.Nodes[origin.Add(transformation.Dot(node))] = true
				}
				unifiedScanner.Edges = fullyConnectedEdges(unifiedScanner.AllNodes())
				scanners = append(scanners, unifiedScanner)

				// Reset for next raw scanner
				//fmt.Println(len(remainingScanners), "remaining scanners")
				remainingScanners = RemoveAtIndex(remainingScanners, rawIndex)
				fmt.Println(len(remainingScanners), "remaining scanners")
				rawIndex = -1
				break
			}
		}
		rawIndex++
	}

	max := 0
	for i := 0; i < len(scannerPositions)-1; i++ {
		for j := i + 1; j < len(scannerPositions); j++ {
			manhattenDistance := abs(scannerPositions[i][0]-scannerPositions[j][0]) +
				abs(scannerPositions[i][1]-scannerPositions[j][1]) +
				abs(scannerPositions[i][2]-scannerPositions[j][2])
			if max < manhattenDistance {
				max = manhattenDistance
			}
		}
	}
	fmt.Println("Maximum distance between scanners of", max)
	//var scanners []Vector
	//var beacons []Vector
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

func fullyConnectedEdges(nodes []Vector) map[Vector]map[Vector]int {
	edges := make(map[Vector]map[Vector]int)
	for _, n1 := range nodes {
		for _, n2 := range nodes {
			if n1 == n2 {
				continue
			}
			if _, ok := edges[n1]; !ok {
				edges[n1] = make(map[Vector]int)
			}
			edges[n1][n2] = (n1[0]-n2[0])*(n1[0]-n2[0]) + (n1[1]-n2[1])*(n1[1]-n2[1]) + (n1[2]-n2[2])*(n1[2]-n2[2])
		}
	}
	return edges
}

func findMatches(a map[Vector]int, b map[Vector]int) [][2]Vector {
	var matches [][2]Vector
	for va, da := range a {
		for vb, db := range b {
			if da == db {
				//fmt.Println("Found tentative match", da, db)
				matches = append(matches, [2]Vector{va, vb})
			}
		}
	}
	return matches
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func differenceWithUniqueComponents(vectors []Vector) (Vector, Vector) {
	for _, v1 := range vectors {
		for _, v2 := range vectors {
			if v1 == v2 {
				continue
			}

			diff := v1.Sub(v2)
			// Zero cannot be used to determine the sign
			for _, c := range diff {
				if c == 0 {
					continue
				}
			}
			fmt.Println(v1, v2)
			if abs(diff[0]) != abs(diff[1]) && abs(diff[0]) != abs(diff[2]) && abs(diff[1]) != abs(diff[2]) {
				return v1, v2
			}
		}
	}
	return Vector{0, 0, 0}, Vector{0, 0, 0}
}

func determineTransformation(from, to Vector) Matrix {
	//fmt.Println(from, "->", to)
	var mat Matrix
	for i, target := range to {
		for j, source := range from {
			if abs(source) == abs(target) {
				mat[i][j] = target / source
			}
		}
	}
	return mat
}

func RemoveAtIndex(graph []Graph, index int) []Graph {
	//fmt.Println(len(graph), index)
	if index == 0 {
		return graph[1:]
	} else if index == len(graph)-1 {
		return graph[:len(graph)-1]
	} else {
		return append(graph[:index], graph[index+1:]...)
	}
}
