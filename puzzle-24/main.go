package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"
)

type Graph struct {
	Nodes map[string]Node
	Edges map[string]map[string]Edge
}

func (g *Graph) Node(id string) Node {
	return g.Nodes[id]
}

func (g *Graph) SetEdge(uid, vid string) {
	if _, ok := g.Nodes[uid]; !ok {
		if isUpper(uid) {
			g.Nodes[uid] = Node{uid, "large"}
		} else {
			g.Nodes[uid] = Node{uid, "small"}
		}
		g.Edges[uid] = make(map[string]Edge)
	}
	if _, ok := g.Nodes[vid]; !ok {
		if isUpper(vid) {
			g.Nodes[vid] = Node{vid, "large"}
		} else {
			g.Nodes[vid] = Node{vid, "small"}
		}
		g.Edges[vid] = make(map[string]Edge)
	}

	g.Edges[uid][vid] = Edge{uid, vid}
	g.Edges[vid][uid] = Edge{vid, uid}
}

type Node struct {
	Id   string
	Size string
}

type Edge struct {
	SourceId string
	TargetId string
}

func main() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Current working directory: %v\n", workingDirectory)

	graph := Graph{make(map[string]Node), make(map[string]map[string]Edge)}
	//example := []string{"start-A", "start-b", "A-c", "A-b", "b-d", "A-end", "b-end"}
	//for _, line := range example {
	for line := range readLines(workingDirectory + "/puzzle-24/input.txt") {
		ids := strings.SplitN(line, "-", 2)
		graph.SetEdge(ids[0], ids[1])
	}
	for source, targets := range graph.Edges {
		fmt.Printf("%v ->", source)
		for target, _ := range targets {
			fmt.Printf(" %v", target)
		}
		fmt.Println("")
	}

	paths := distinctPaths(&graph, []Node{Node{"start", "small"}}, [][]Node{})
	//for _, p := range paths {
	//	fmt.Println(p)
	//}
	fmt.Printf("Number of distinct paths: %v\n", len(paths))
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

func isUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func distinctPaths(g *Graph, path []Node, paths [][]Node) [][]Node {
	// Accept path if traversal is complete with "end" node
	if path[len(path)-1] == (Node{"end", "small"}) {
		paths = append(paths, path)
		return paths
	}

	for vid, _ := range (*g).Edges[path[len(path)-1].Id] {
		vnode := (*g).Nodes[vid]
		// Special condition: "start" may only be visited once
		if vnode.Id == "start" {
			continue
		}
		// Special conditions: large caves can be visited any time and "start" may only be visited once
		if vnode.Size == "large" || !contains(path, vnode) || (contains(path, vnode) && max(countSmalls(path)) < 2) {
			// Each branch (potentially) results in a distinct path
			branch := make([]Node, len(path))
			copy(branch, path)
			branch = append(branch, vnode)
			paths = distinctPaths(g, branch, paths)
		}
	}
	return paths
}

func contains(nodes []Node, node Node) bool {
	for _, n := range nodes {
		if n == node {
			return true
		}
	}
	return false
}

func max(counts map[string]int) int {
	maximum := 0
	for _, v := range counts {
		if v > maximum {
			maximum = v
		}
	}
	return maximum
}

func countSmalls(nodes []Node) map[string]int {
	counts := make(map[string]int)
	for _, n := range nodes {
		if n.Size == "large" {
			continue
		}
		if _, ok := counts[n.Id]; ok {
			counts[n.Id]++
		} else {
			counts[n.Id] = 1
		}
	}
	return counts
}
