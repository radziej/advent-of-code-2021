package main

import "fmt"

// map[sumOfRolls]frequency
var rolls = map[int]int{
	// 1, 1, 1 = 3
	3: 1,
	// 2, 1, 1 = 4
	// 1, 2, 1 = 4
	// 1, 1, 2 = 4
	4: 3,
	// 3, 1, 1 = 5
	// 1, 3, 1 = 5
	// 1, 1, 3 = 5
	// 2, 2, 1 = 5
	// 2, 1, 2 = 5
	// 1, 2, 2 = 5
	5: 6,
	// 2, 2, 2 = 6
	// 3, 2, 1 = 6
	// 1, 3, 2 = 6
	// 2, 1, 3 = 6
	// 2, 3, 1 = 6
	// 1, 2, 3 = 6
	// 3, 1, 2 = 6
	6: 7,
	// 3, 3, 1 = 7
	// 3, 1, 3 = 7
	// 1, 3, 3 = 7
	// 3, 2, 2 = 7
	// 2, 3, 2 = 7
	// 2, 2, 3 = 7
	7: 6,
	// 3, 3, 2 = 8
	// 3, 2, 3 = 8
	// 2, 3, 3 = 8
	8: 3,
	// 3, 3, 3 = 9
	9: 1,
}

func main() {
	// Initial scores and positions
	score1 := 0
	space1 := 8 // Input
	//space1 := 4 // Test
	score2 := 0
	space2 := 4 // Input
	//space2 := 8 // Test

	victories := map[int]int{
		1: 0,
		2: 0,
	}
	play(0, 1, score1, space1, score2, space2, &victories)

	fmt.Println("Victories player 1:", victories[1])
	fmt.Println("Victories player 2:", victories[2])
}

//14598271321389
//444356092776315
//9223372036854775807

func play(turn int, branches int, score1 int, space1 int, score2 int, space2 int, victories *map[int]int) *map[int]int {
	// Win/return condition
	if score1 >= 21 || score2 >= 21 {
		if score1 > score2 {
			(*victories)[1] += branches
		} else {
			(*victories)[2] += branches
		}
		return victories
	}

	// Play another (set of) turn(s)
	if turn%2 == 0 {
		for roll, frequency := range rolls {
			tmpSpace1 := advanceSteps(space1, roll)
			tmpScore1 := score1 + tmpSpace1
			victories = play(turn+1, branches*frequency, tmpScore1, tmpSpace1, score2, space2, victories)
		}
	} else {
		for roll, frequency := range rolls {
			tmpSpace2 := advanceSteps(space2, roll)
			tmpScore2 := score2 + tmpSpace2
			victories = play(turn+1, branches*frequency, score1, space1, tmpScore2, tmpSpace2, victories)
		}
	}
	return victories
}

func advanceSteps(space, steps int) int {
	result := space + steps
	for result > 10 {
		result -= 10
	}
	return result
}
