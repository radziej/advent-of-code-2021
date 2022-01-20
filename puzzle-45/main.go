package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strings"
)

const NumAmphipods int = 4 * 2

var Homes = map[string]int{
	"A": 3,
	"B": 5,
	"C": 7,
	"D": 9,
}

type Vector2D [2]int

type Amphipod struct {
	Position Vector2D
	Value    string
}

type State [NumAmphipods]Amphipod

func (s State) String() string {
	grid := [][]string{
		{"#", "#", "#", "#", "#", "#", "#", "#", "#", "#", "#", "#", "#"},
		{"#", ".", ".", ".", ".", ".", ".", ".", ".", ".", ".", ".", "#"},
		{"#", "#", "#", ".", "#", ".", "#", ".", "#", ".", "#", "#", "#"},
		{" ", " ", "#", ".", "#", ".", "#", ".", "#", ".", "#", " ", " "},
		{" ", " ", "#", "#", "#", "#", "#", "#", "#", "#", "#", " ", " "},
	}
	for _, pod := range s {
		grid[pod.Position[1]][pod.Position[0]] = pod.Value
	}

	representation := strings.Join(grid[0], "")
	for _, row := range grid[1:] {
		representation += "\n" + strings.Join(row, "")
	}
	return representation
}

func (s *State) sort() {
	sort.Slice((*s)[:], func(i, j int) bool {
		if (*s)[i].Position[0] == (*s)[j].Position[0] {
			return (*s)[i].Position[1] < (*s)[j].Position[1]
		}
		return (*s)[i].Position[0] < (*s)[j].Position[0]
	})
}

type Graph struct {
	Nodes map[State]int
	Edges map[State]map[State]int
}

func (g *Graph) SetNode(state State, cost int) {
	g.Nodes[state] = cost
	if _, ok := g.Edges[state]; !ok {
		g.Edges[state] = make(map[State]int)
	}
}

func main() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Current working directory: %v\n", workingDirectory)
	//lines := readAllLines(workingDirectory + "/puzzle-45/test.txt")
	lines := readAllLines(workingDirectory + "/puzzle-45/input.txt")

	// (x, y) with valid positions for amphipods
	validPositions := []Vector2D{
		{1, 1},
		{2, 1},
		{3, 2},
		{3, 3},
		{4, 1},
		{5, 2},
		{5, 3},
		{6, 1},
		{7, 2},
		{7, 3},
		{8, 1},
		{9, 2},
		{9, 3},
		{10, 1},
		{11, 1},
	}

	var currentState State
	index := 0
	for _, vec := range validPositions {
		if occupant := lines[vec[1]][vec[0]]; occupant >= 'A' && occupant <= 'D' {
			currentState[index] = Amphipod{vec, string(occupant)}
			index++
		}
	}
	fmt.Println(currentState)

	// :)
	graph := Graph{make(map[State]int), make(map[State]map[State]int)}
	graph.SetNode(currentState, 0)
	exploreStateGraph(&graph, currentState)
	fmt.Println("Total of", len(graph.Nodes), "states to consider.")

	// Determining shortest path
	target := State{
		{Vector2D{3, 2}, "A"},
		{Vector2D{3, 3}, "A"},
		{Vector2D{5, 2}, "B"},
		{Vector2D{5, 3}, "B"},
		{Vector2D{7, 2}, "C"},
		{Vector2D{7, 3}, "C"},
		{Vector2D{9, 2}, "D"},
		{Vector2D{9, 3}, "D"},
	}
	if value, ok := graph.Nodes[target]; ok {
		fmt.Println(value)
	} else {
		fmt.Println("Missing target")
	}
	predecessors := shortestPath(&graph, map[State]int{currentState: 0})

	fmt.Println("-- A path that requires minimum energy --")
	path := []State{target}
	for predecessors[path[0]] != (State{}) {
		path = append([]State{predecessors[path[0]]}, path...)
	}
	for _, s := range path {
		fmt.Println(s)
		fmt.Println("")
	}
	fmt.Println("Least energy required for organizing Amphipods:", graph.Nodes[target])
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

func exploreStateGraph(graph *Graph, state State) {
	// Using a map instead of an array to avoid unnecessary iterations
	positions := make(map[Vector2D]string, len(state))
	for _, pod := range state {
		positions[pod.Position] = pod.Value
	}

	var states []State
	for index, pod := range state {
		if pod.Position[1] > 1 { // Apply rules for inside any home
			// No additional moves when at final destination
			if pod.Position[0] == Homes[pod.Value] {
				if pod.Position[1] == 3 {
					// In first position
					continue
				} else if pod.Position[1] == 2 && positions[Vector2D{pod.Position[0], 3}] == pod.Value {
					// In second position
					continue
				}
			}

			// Moving up/out of home
			exitBlocked := false
			for y := pod.Position[1] - 1; y >= 1; y-- {
				if _, ok := positions[Vector2D{pod.Position[0], y}]; ok {
					exitBlocked = true
					break
				}
			}
			if exitBlocked {
				continue
			}

			// Moving left
			for x := pod.Position[0] - 1; x >= 1; x-- {
				if _, ok := positions[Vector2D{x, 1}]; ok {
					break
				} else if !ok && x != 3 && x != 5 && x != 7 && x != 9 {
					newPod := Amphipod{Vector2D{x, 1}, pod.Value}
					newState := State{}
					copy(newState[:], state[:])
					newState[index] = newPod
					newState.sort()
					if _, ok := (*graph).Nodes[newState]; !ok {
						(*graph).SetNode(newState, math.MaxInt64)
						states = append(states, newState)
					}
					(*graph).Edges[state][newState] = CalculateTransitionCost(pod, newPod)
				}
			}

			// Moving right
			for x := pod.Position[0] + 1; x <= 11; x++ {
				if _, ok := positions[Vector2D{x, 1}]; ok {
					break
				} else if !ok && x != 3 && x != 5 && x != 7 && x != 9 {
					newPod := Amphipod{Vector2D{x, 1}, pod.Value}
					newState := State{}
					copy(newState[:], state[:])
					newState[index] = newPod
					newState.sort()
					if _, ok := (*graph).Nodes[newState]; !ok {
						(*graph).SetNode(newState, math.MaxInt64)
						states = append(states, newState)
					}
					(*graph).Edges[state][newState] = CalculateTransitionCost(pod, newPod)
				}
			}
		} else { // Apply rules for hallway
			// Find position in target home to move into if not blocked
			xHome := Homes[pod.Value]
			var target Vector2D
			//homeBlocked := false
			if _, ok := positions[Vector2D{xHome, 2}]; ok {
				continue
				//homeBlocked = true
			} else {
				target = Vector2D{xHome, 2}
			}
			if value, ok := positions[Vector2D{xHome, 3}]; ok {
				if value != pod.Value {
					continue
					//homeBlocked = true
				}
			} else {
				target = Vector2D{xHome, 3}
			}
			//for y := 2; y <= 3; y++ {
			//	if value, ok := positions[Vector2D{xHome, y}]; ok {
			//		if value != pod.Value {
			//			homeBlocked = true
			//			break
			//		}
			//	} else {
			//		target = Vector2D{xHome, y}
			//	}
			//}
			//if homeBlocked {
			//	continue
			//}

			// Check if path to target is clear
			direction := 0
			if xHome > pod.Position[0] {
				direction = 1
			} else if xHome == pod.Position[0] {
				panic(fmt.Sprintf("state is not possible: %v", state))
			} else {
				direction = -1
			}
			blocked := false
			for x := pod.Position[0] + direction; x != xHome; {
				if _, ok := positions[Vector2D{x, 1}]; ok {
					blocked = true
					break
				}
				x += direction
			}
			if blocked {
				continue
			}

			newPod := Amphipod{target, pod.Value}
			newState := State{}
			copy(newState[:], state[:])
			newState[index] = newPod
			newState.sort()
			if _, ok := (*graph).Nodes[newState]; !ok {
				(*graph).SetNode(newState, math.MaxInt64)
				states = append(states, newState)
			}
			(*graph).Edges[state][newState] = CalculateTransitionCost(pod, newPod)
		}
	}
	for _, nextState := range states {
		exploreStateGraph(graph, nextState)
	}
}

func CalculateTransitionCost(from, to Amphipod) int {
	distance := absInt(to.Position[0]-from.Position[0]) + absInt(to.Position[1]-from.Position[1])
	switch from.Value {
	case "A":
		return distance * 1
	case "B":
		return distance * 10
	case "C":
		return distance * 100
	case "D":
		return distance * 1000
	default:
		panic(fmt.Sprintf("unknown value: %v", from.Value))
	}
}

func absInt(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func shortestPath(graph *Graph, tentatives map[State]int) map[State]State {
	unvisited := make(map[State]struct{}, len((*graph).Nodes))
	predecessor := make(map[State]State, len((*graph).Nodes))
	for node, _ := range (*graph).Nodes {
		unvisited[node] = struct{}{}
		predecessor[node] = State{}
	}

	for len(tentatives) > 0 {
		if len(unvisited)%1000 == 0 {
			fmt.Println(len(unvisited), "remaining unvisited nodes")
		}
		// Find item with lowest tentative cost
		state := State{}
		currentCost := math.MaxInt64
		for s, c := range tentatives {
			if currentCost > c {
				currentCost = c
				state = s
			}
		}
		delete(tentatives, state)
		delete(unvisited, state)

		for neighbor, transitionCost := range (*graph).Edges[state] {
			// Do not visit same node twice
			if _, ok := unvisited[neighbor]; !ok {
				continue
			}

			totalCost := currentCost + transitionCost
			if neighborCost := (*graph).Nodes[neighbor]; neighborCost > totalCost {
				(*graph).Nodes[neighbor] = totalCost
				predecessor[neighbor] = state
			}
			tentatives[neighbor] = (*graph).Nodes[neighbor]
		}
	}

	return predecessor
}
