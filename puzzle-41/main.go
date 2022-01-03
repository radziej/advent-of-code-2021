package main

import "fmt"

type DeterministicDie struct {
	state int
}

func (dd *DeterministicDie) incrementState() {
	dd.state++
	if dd.state > 100 {
		dd.state = 1
	}
}

func (dd *DeterministicDie) Generate(size int) []int {
	result := make([]int, size)
	for i := range result {
		result[i] = dd.state
		dd.incrementState()
	}
	return result
}

func main() {
	// Initial scores and positions
	score1 := 0
	space1 := 8
	//space1 := 4  // Testing space
	score2 := 0
	space2 := 4
	//space2 := 8 // Testing space

	rolls := 0
	die := DeterministicDie{1}
	for turn := 0; ; turn++ {
		// Simulate rolls of turn
		stepsForward := 0
		for _, value := range die.Generate(3) {
			stepsForward += value
		}
		rolls += 3
		//fmt.Println(stepsForward, "steps forward")

		// Assign rolls to single player per turn
		if turn%2 == 0 {
			space1 = advanceSteps(space1, stepsForward)
			score1 += space1
		} else {
			space2 = advanceSteps(space2, stepsForward)
			score2 += space2
		}

		// Check win condition
		if score1 >= 1000 || score2 >= 1000 {
			if score1 < score2 {
				fmt.Println("Losing score", score1, "*", rolls, "=", score1*rolls)
			} else {
				fmt.Println("Losing score", score2, "*", rolls, "=", score2*rolls)
			}
			break
		}

		fmt.Println(score1, space1, "--", score2, space2)
	}
}

func advanceSteps(space, steps int) int {
	result := space + steps
	for result > 10 {
		result -= 10
	}
	return result
}
