package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type Number struct {
	Value    int
	Parent   *Number
	Children []*Number
}

func (n Number) String() string {
	if n.Value != -1 {
		return fmt.Sprint(n.Value)
	}
	if n.Children != nil {
		var strs []string
		for _, c := range n.Children {
			strs = append(strs, (*c).String())
		}
		return "[" + strings.Join(strs, ", ") + "]"
	}
	return fmt.Sprintf("Invalid number: %v, %v", n.Value, n.Children)
}

func main() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Current working directory: %v\n", workingDirectory)

	var numbers []*Number
	parsedNumbers := 0
	for line := range readLines(workingDirectory + "/puzzle-36/input.txt") {
		//for _, line := range []string{
		//	// Larger example
		//	//"[[[0,[4,5]],[0,0]],[[[4,5],[2,6]],[9,5]]]",
		//	//"[7,[[[3,7],[4,3]],[[6,3],[8,8]]]]",
		//	//"[[2,[[0,8],[3,4]]],[[[6,7],1],[7,[1,6]]]]",
		//	//"[[[[2,4],7],[6,[0,5]]],[[[6,8],[2,8]],[[2,1],[4,5]]]]",
		//	//"[7,[5,[[3,8],[1,4]]]]",
		//	//"[[2,[2,2]],[8,[8,1]]]",
		//	//"[2,9]",
		//	//"[1,[[[9,3],9],[[9,0],[0,7]]]]",
		//	//"[[[5,[7,4]],7],1]",
		//	//"[[[[4,2],2],6],[8,7]]",
		//	// Full example
		//	"[[[0,[5,8]],[[1,7],[9,6]]],[[4,[1,2]],[[1,4],2]]]",
		//	"[[[5,[2,8]],4],[5,[[9,9],0]]]",
		//	"[6,[[[6,2],[5,6]],[[7,6],[4,7]]]]",
		//	"[[[6,[0,7]],[0,9]],[4,[9,[9,0]]]]",
		//	"[[[7,[6,4]],[3,[1,3]]],[[[5,5],1],9]]",
		//	"[[6,[[7,3],[3,2]]],[[[3,8],[5,7]],4]]",
		//	"[[[[5,4],[7,7]],8],[[8,3],8]]",
		//	"[[9,3],[[9,9],[6,[4,9]]]]",
		//	"[[2,[[7,7],7]],[[5,8],[[9,3],[0,2]]]]",
		//	"[[[[5,2],5],[8,[3,7]]],[[5,[7,5]],[4,4]]]",
		//} {
		buffer := ""
		root := &Number{-1, nil, nil}
		number := root
		for _, r := range line[1:] {
			character := string(r)
			switch character {
			case "[":
				newNumber := &Number{-1, number, nil}
				number.Children = append(number.Children, newNumber)
				number = newNumber
			case "]":
				if buffer != "" {
					number.Children = append(number.Children, &Number{TokenToInt(buffer), number, nil})
				}
				number = number.Parent
				buffer = ""
			case ",":
				if buffer != "" {
					number.Children = append(number.Children, &Number{TokenToInt(buffer), number, nil})
				}
				buffer = ""
			default:
				buffer += character
			}
		}
		fmt.Printf("Parsed %v number %v\n", parsedNumbers, root)
		parsedNumbers++
		numbers = append(numbers, root)
	}

	number := copyNumber(numbers[0], nil)
	for _, n := range numbers[1:] {
		number = add(number, copyNumber(n, nil))
		//fmt.Println(number)
		reduced := false
		for !reduced {
			reduced = reduce(number)
			//fmt.Println("Reduced to", number)
		}
	}
	fmt.Println("Number:", number)
	fmt.Println("Magnitude:", magnitude(number))

	maxMagnitude := 0
	for i := 0; i < len(numbers); i++ {
		for j := 0; j < len(numbers); j++ {
			if i == j {
				continue
			}
			// Create deep copies to avoid manipulating addends
			a := copyNumber(numbers[i], nil)
			b := copyNumber(numbers[j], nil)
			// Add second addend and reduce
			number := add(a, b)
			reduced := false
			for !reduced {
				reduced = reduce(number)
			}

			if mag := magnitude(number); mag > maxMagnitude {
				fmt.Println(number)
				fmt.Printf("Found larger magnitude of %v using %v: %v and %v: %v\n", mag, i, numbers[i], j, numbers[j])
				maxMagnitude = mag
			}
		}
	}
	fmt.Println("Max magnitude:", maxMagnitude)
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

func TokenToInt(t string) int {
	integer, err := strconv.Atoi(t)
	if err != nil {
		log.Fatal(err)
	}
	return integer
}

func add(a *Number, b *Number) *Number {
	number := &Number{-1, nil, []*Number{a, b}}
	number.Children[0].Parent = number
	number.Children[1].Parent = number
	return number
}

func reduce(n *Number) (reduced bool) {
	// Explode leftmost pair that is nested inside 4 other pairs
	if path := findNested(n, []int{}); len(path) != 0 {
		//fmt.Println("Exploding", getChild(n, path), "at", path)
		explode(n, path)
		return false
	}

	// Split leftmost individual value that is equal or greater than 10
	if hasSplit := split(n); hasSplit {
		return false
	}
	return true
}

func findNested(n *Number, path []int) []int {
	if len(path) == 4 || n.Children == nil {
		return path
	}

	for i := range n.Children {
		if n.Children[i].Children != nil {
			deeper := make([]int, len(path)+1)
			copy(deeper, path)
			deeper[len(path)] = i
			//fmt.Println(deeper)
			if deeperPath := findNested(n.Children[i], deeper); len(deeperPath) == 4 {
				return deeperPath
			}
		}
	}
	return []int{}
}

// The branching numbers are assumed to follow a *binary* tree hierarchy. If this is not the case, the function can be
// repurposed by avoiding hard-coded indices.
func explode(n *Number, path []int) {
	// Retrieve and replace number to explode
	parent := getChild(n, path[:len(path)-1])
	explodedNumber := parent.Children[path[len(path)-1]]
	parent.Children[path[len(path)-1]] = &Number{0, parent, nil}

	// Left
	var leftPath []int
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == 1 {
			leftPath = make([]int, len(path[:i]))
			copy(leftPath, path[:i])
			leftPath = append(leftPath, 0)
			break
		}
	}
	if leftPath != nil {
		leftNumber := getChild(n, leftPath)
		for leftNumber.Children != nil {
			leftNumber = leftNumber.Children[1]
		}
		leftNumber.Value += explodedNumber.Children[0].Value
	}

	// Right
	var rightPath []int
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == 0 {
			rightPath = make([]int, len(path[:i]))
			copy(rightPath, path[:i])
			rightPath = append(rightPath, 1)
			break
		}
	}
	if rightPath != nil {
		rightNumber := getChild(n, rightPath)
		for rightNumber.Children != nil {
			rightNumber = rightNumber.Children[0]
		}
		rightNumber.Value += explodedNumber.Children[1].Value
	}
}

func getChild(n *Number, path []int) *Number {
	for _, step := range path {
		n = n.Children[step]
	}
	return n
}

func split(n *Number) bool {
	if n.Children == nil {
		return false
	}

	for i := range n.Children {
		if n.Children[i].Value >= 10 {
			//fmt.Println("Splitting", n.Children[i])
			number := &Number{-1, n, nil}
			number.Children = append(number.Children, &Number{int(math.Floor(float64(n.Children[i].Value) / 2.0)), number, nil})
			number.Children = append(number.Children, &Number{int(math.Ceil(float64(n.Children[i].Value) / 2.0)), number, nil})
			n.Children[i] = number
			return true
		} else if n.Children[i].Children != nil {
			if result := split(n.Children[i]); result {
				return result
			}
		}
	}
	return false
}

func magnitude(n *Number) int {
	if n.Children == nil {
		return n.Value
	}
	return 3*magnitude(n.Children[0]) + 2*magnitude(n.Children[1])
}

func copyNumber(number *Number, parent *Number) *Number {
	if number.Children == nil {
		return &Number{number.Value, parent, nil}
	}
	numberCopy := &Number{-1, parent, nil}
	for _, c := range number.Children {
		numberCopy.Children = append(numberCopy.Children, copyNumber(c, numberCopy))
	}
	return numberCopy
}
