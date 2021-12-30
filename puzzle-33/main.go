package main

import (
	"fmt"
)

type Vector [2]int

func main() {
	// This solution assumes that the target area does not overlap with the origin of the x axis and is to be found in
	// the positive x direction and negative y direction.

	// Skipping trivial input retrieval from text file
	targetX := [2]int{70, 96}
	targetY := [2]int{-179, -124}
	// Test ranges
	//targetX := [2]int{20, 30}
	//targetY := [2]int{-10, -5}

	origin := Vector{0, 0}
	maxY := origin[1]
	for vx := 1; vx <= targetX[1]; vx++ {
		// Descent is symmetric; positive v_y always returns to 0 with the next step being -v_y - 1
		for vy := targetY[0]; vy <= abs(targetY[0]); vy++ {
			// Setting initial conditions of simulation
			position := origin
			velocity := Vector{vx, vy}
			tentativeY := position[1]

			for step := 0; ; step++ {
				if tentativeY < position[1] {
					tentativeY = position[1]
				}

				if position[0] >= targetX[0] && position[0] <= targetX[1] && position[1] >= targetY[0] && position[1] <= targetY[1] {
					if maxY < tentativeY {
						maxY = tentativeY
						fmt.Printf("Found better initial velocity v=(%v, %v) with y_max=%v\n", vx, vy, maxY)
					}
				}

				if velocity[0] == 0 && position[0] < targetX[0] {
					//fmt.Println("Undershot in x direction")
					break
				}
				if position[0] > targetX[1] {
					//fmt.Println("Overshot in x direction")
					break
				} else if velocity[1] < 0 && position[1] < targetY[0] {
					//fmt.Println("Overshot in y direction")
					break
				}
				// Advance position by velocity
				position = Vector{position[0] + velocity[0], position[1] + velocity[1]}
				velocity = Vector{stepX(velocity[0]), velocity[1] - 1}
				//time.Sleep(5 * time.Millisecond)
				//fmt.Println(position, velocity)
			}
			if position[1] > targetY[1] {
				//fmt.Printf("Descent too slow for vx=%v at vy=%v with final position %v\n", vx, vy, position)
				break
			}
		}
	}
}

func contains(slice []int, value int) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func stepX(x int) int {
	if x > 0 {
		return x - 1
	} else if x < 0 {
		return x + 1
	}
	return x
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
