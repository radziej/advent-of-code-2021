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

func main() {
	workingDirectory, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Current working directory: %v\n", workingDirectory)
	lines := readAllLines(workingDirectory + "/puzzle-48/input.txt")
	alu := NewArithmeticLogicUnit()

	// Tests
	//lines = []string{
	//	"inp x",
	//	"mul x -1",
	//}
	//alu.LoadProgram(lines)
	//fmt.Println("Expecting x = -10:", alu.ExecuteProgram([]int{10}))
	//alu.ResetRegisters()
	//
	//lines = []string{
	//	"inp z",
	//	"inp x",
	//	"mul z 3",
	//	"eql z x",
	//}
	//alu.LoadProgram(lines)
	//fmt.Println("Expecting z = 0:", alu.ExecuteProgram([]int{10, 10}))
	//alu.ResetRegisters()
	//fmt.Println("Expecting z = 1:", alu.ExecuteProgram([]int{10, 30}))
	//alu.ResetRegisters()
	//
	//lines = []string{
	//	"inp w",
	//	"add z w",
	//	"mod z 2",
	//	"div w 2",
	//	"add y w",
	//	"mod y 2",
	//	"div w 2",
	//	"add x w",
	//	"mod x 2",
	//	"div w 2",
	//	"mod w 2",
	//}
	//alu.LoadProgram(lines)
	//fmt.Println("Expecting w, x, y, z = 1:", alu.ExecuteProgram([]int{15}))
	//alu.ResetRegisters()
	//fmt.Println("Expecting w, x, y, z = 0:", alu.ExecuteProgram([]int{0}))
	//alu.ResetRegisters()

	// Brute force :)
	//alu.LoadProgram(lines)
	//for modelNumber := 99999999999999; modelNumber >= 11111111111111; modelNumber-- {
	//	if modelNumber%1000000 == 0 {
	//		fmt.Println("Testing model number:", modelNumber)
	//	}
	//
	//	skip := false
	//	input := make([]int, 14)
	//	for bound, v := range strings.Split(strconv.Itoa(modelNumber), "") {
	//		input[bound], _ = strconv.Atoi(v)
	//		if input[bound] == 0 {
	//			skip = true
	//			break
	//		}
	//	}
	//	if skip {
	//		continue
	//	}
	//	if result := alu.ExecuteProgram(input); result["z"] == 0 {
	//		fmt.Println("Highest accepted model number:", modelNumber)
	//		break
	//	}
	//}
	//alu.ResetRegisters()

	// Find boundaries for each part of the program
	var bounds [14]int
	index := 0
	for i, line := range lines {
		if line[0:3] == "inp" {
			bounds[index] = i
			index++
		}
	}

	var stack [][2]int
	pairs := make(map[[2]int][][2]int, 14/2)
	powerOf10 := 13
	for _, bound := range bounds {
		aFields := strings.Fields(lines[bound+5])
		a, err := strconv.Atoi(aFields[2])
		if err != nil {
			log.Fatal(err)
		}

		bFields := strings.Fields(lines[bound+15])
		b, err := strconv.Atoi(bFields[2])
		if err != nil {
			log.Fatal(err)
		}

		if a > 0 { // Push to stack
			stack = append(stack, [2]int{powerOf10, b})
		} else { // Pop from stack
			element := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			// Determine valid digit pairs
			var validPairs [][2]int
			for x := 1; x <= 9; x++ {
				for y := 1; y <= 9; y++ {
					if x+element[1]+a == y {
						validPairs = append(validPairs, [2]int{x, y})
					}
				}
			}
			pairs[[2]int{element[0], powerOf10}] = validPairs
		}
		powerOf10--
	}
	fmt.Println(pairs)

	// Chose highest (consider first and second digit in this order) pairs for largest model number
	modelNumber := 0
	for indices, values := range pairs {
		modelNumber += int(math.Pow10(indices[0])) * values[0][0]
		modelNumber += int(math.Pow10(indices[1])) * values[0][1]
	}
	// Verify result
	alu.LoadProgram(lines)
	splitModelNumber := make([]int, len(bounds))
	for i, s := range strings.Split(strconv.Itoa(modelNumber), "") {
		number, _ := strconv.Atoi(s)
		splitModelNumber[i] = number
	}
	if result := alu.ExecuteProgram(splitModelNumber); result["z"] == 0 {
		fmt.Println("Selected valid model number:", modelNumber)
	} else {
		fmt.Println("Selected model number is invalid:", modelNumber)
	}
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

type Instruction struct {
	Function func(string, int)
	Register string
	Value    string
}

type ArithmeticLogicUnit struct {
	registers map[string]int
	program   []Instruction
	input     []int
}

func (alu *ArithmeticLogicUnit) LoadProgram(instructions []string) {
	alu.program = make([]Instruction, len(instructions))
	for i, instruction := range instructions {
		fields := strings.Fields(instruction)
		switch fields[0] {
		case "inp":
			alu.program[i] = Instruction{alu.FeedInput, fields[1], ""}
		case "add":
			alu.program[i] = Instruction{alu.Add, fields[1], fields[2]}
		case "mul":
			alu.program[i] = Instruction{alu.Mul, fields[1], fields[2]}
		case "div":
			alu.program[i] = Instruction{alu.Div, fields[1], fields[2]}
		case "mod":
			alu.program[i] = Instruction{alu.Mod, fields[1], fields[2]}
		case "eql":
			alu.program[i] = Instruction{alu.Equal, fields[1], fields[2]}
		}
	}
}

func (alu *ArithmeticLogicUnit) ExecuteProgram(input []int) map[string]int {
	alu.input = input
	for _, instruction := range alu.program {
		if instruction.Value == "" ||
			instruction.Value == "w" ||
			instruction.Value == "x" ||
			instruction.Value == "y" ||
			instruction.Value == "z" {
			instruction.Function(instruction.Register, alu.registers[instruction.Value])
		} else {
			number, err := strconv.Atoi(instruction.Value)
			if err != nil {
				panic(err)
			}
			instruction.Function(instruction.Register, number)
		}
	}

	// Reset registers and return copy of result
	result := make(map[string]int, len(alu.registers))
	for k, v := range alu.registers {
		result[k] = v
	}
	return result
}

func (alu *ArithmeticLogicUnit) ResetRegisters() {
	for key := range alu.registers {
		alu.registers[key] = 0
	}
}

func (alu *ArithmeticLogicUnit) FeedInput(register string, _ int) {
	alu.registers[register] = alu.input[0]
	alu.input = alu.input[1:]
}

func (alu *ArithmeticLogicUnit) Add(register string, value int) {
	alu.registers[register] += value
}

func (alu *ArithmeticLogicUnit) Mul(register string, value int) {
	alu.registers[register] *= value
}

func (alu *ArithmeticLogicUnit) Div(register string, value int) {
	alu.registers[register] /= value
}

func (alu *ArithmeticLogicUnit) Mod(register string, value int) {
	alu.registers[register] %= value
}

func (alu *ArithmeticLogicUnit) Equal(register string, value int) {
	if alu.registers[register] == value {
		alu.registers[register] = 1
	} else {
		alu.registers[register] = 0
	}
}

func NewArithmeticLogicUnit() ArithmeticLogicUnit {
	alu := ArithmeticLogicUnit{
		map[string]int{"w": 0, "x": 0, "y": 0, "z": 0},
		nil,
		nil,
	}
	return alu
}
