package main

import (
	"fmt"
)

const program = `3,8,1005,8,336,1106,0,11,0,0,0,104,1,104,0,3,8,102,-1,8,10,1001,10,1,10,4,10,108,1,8,10,4,10,101,0,8,28,1006,0,36,1,2,5,10,1006,0,57,1006,0,68,3,8,102,-1,8,10,1001,10,1,10,4,10,108,0,8,10,4,10,1002,8,1,63,2,6,20,10,1,106,7,10,2,9,0,10,3,8,102,-1,8,10,101,1,10,10,4,10,108,1,8,10,4,10,102,1,8,97,1006,0,71,3,8,1002,8,-1,10,101,1,10,10,4,10,108,1,8,10,4,10,1002,8,1,122,2,105,20,10,3,8,1002,8,-1,10,1001,10,1,10,4,10,108,0,8,10,4,10,101,0,8,148,2,1101,12,10,1006,0,65,2,1001,19,10,3,8,102,-1,8,10,1001,10,1,10,4,10,108,0,8,10,4,10,101,0,8,181,3,8,1002,8,-1,10,1001,10,1,10,4,10,1008,8,0,10,4,10,1002,8,1,204,2,7,14,10,2,1005,20,10,1006,0,19,3,8,102,-1,8,10,101,1,10,10,4,10,108,1,8,10,4,10,102,1,8,236,1006,0,76,1006,0,28,1,1003,10,10,1006,0,72,3,8,1002,8,-1,10,101,1,10,10,4,10,108,0,8,10,4,10,102,1,8,271,1006,0,70,2,107,20,10,1006,0,81,3,8,1002,8,-1,10,1001,10,1,10,4,10,108,1,8,10,4,10,1002,8,1,303,2,3,11,10,2,9,1,10,2,1107,1,10,101,1,9,9,1007,9,913,10,1005,10,15,99,109,658,104,0,104,1,21101,0,387508441896,1,21102,1,353,0,1106,0,457,21101,0,937151013780,1,21101,0,364,0,1105,1,457,3,10,104,0,104,1,3,10,104,0,104,0,3,10,104,0,104,1,3,10,104,0,104,1,3,10,104,0,104,0,3,10,104,0,104,1,21102,179490040923,1,1,21102,411,1,0,1105,1,457,21101,46211964123,0,1,21102,422,1,0,1106,0,457,3,10,104,0,104,0,3,10,104,0,104,0,21101,838324716308,0,1,21101,0,445,0,1106,0,457,21102,1,868410610452,1,21102,1,456,0,1106,0,457,99,109,2,22101,0,-1,1,21101,40,0,2,21101,0,488,3,21101,478,0,0,1106,0,521,109,-2,2105,1,0,0,1,0,0,1,109,2,3,10,204,-1,1001,483,484,499,4,0,1001,483,1,483,108,4,483,10,1006,10,515,1101,0,0,483,109,-2,2105,1,0,0,109,4,2101,0,-1,520,1207,-3,0,10,1006,10,538,21101,0,0,-3,22102,1,-3,1,21202,-2,1,2,21101,0,1,3,21101,557,0,0,1105,1,562,109,-4,2105,1,0,109,5,1207,-3,1,10,1006,10,585,2207,-4,-2,10,1006,10,585,22101,0,-4,-4,1106,0,653,21201,-4,0,1,21201,-3,-1,2,21202,-2,2,3,21102,604,1,0,1106,0,562,21202,1,1,-4,21101,0,1,-1,2207,-4,-2,10,1006,10,623,21102,0,1,-1,22202,-2,-1,-2,2107,0,-3,10,1006,10,645,21202,-1,1,1,21101,0,645,0,106,0,520,21202,-2,-1,-2,22201,-4,-2,-4,109,-5,2105,1,0`

type direction int

const (
	up direction = iota
	right
	down
	left
)

func (d direction) String() string {
	switch d {
	case up:
		return "up"
	case right:
		return "right"
	case down:
		return "down"
	case left:
		return "left"
	}
	panic("unknown direction")
}

func main() {
	grid := [50][100]int{}
	painted := [1000][1000]int{}
	x := 0
	y := 0
	var dir direction

	// first panel should be white (part 2)
	grid[x][y] = 1

	getCurrentColor := func() int {
		return grid[x][y]
	}

	icc := NewComputer(program)
	output := make(chan int)
	go icc.Run(output)

	for {
		// pass current color as input
		c := getCurrentColor()
		fmt.Println("sending current color", c)
		icc.SendInput(c)

		// first output is new paint
		newPaint, ok := <-output
		if !ok { // program done
			break
		}
		fmt.Println("new paint", newPaint)
		grid[y][x] = newPaint
		painted[y][x] = 1

		// second output is turn direction
		turn, ok := <-output
		if !ok { // program done
			break
		}
		fmt.Printf("making turn (%d) from %s", turn, dir)
		dir = makeTurn(dir, turn)
		fmt.Printf(", got %s\n", dir)

		// step once
		x, y = step(grid, x, y, dir)
		fmt.Printf("new pos [%d, %d]\n", x, y)
	}

	// count how many were painted
	// var totalPainted int
	// for _, line := range painted {
	// 	for _, col := range line {
	// 		totalPainted = totalPainted + col
	// 	}
	// }

	// fmt.Println("total painted", totalPainted)

	for _, line := range grid {
		for _, col := range line {
			if col == 0 {
				fmt.Printf(" ")
			} else {
				fmt.Printf("#")
			}
		}
		fmt.Printf("\n")
	}
}

// returns new x, y
func step(grid [50][100]int, x int, y int, dir direction) (int, int) {

	switch dir {
	case up:
		return x, y - 1
	case right:
		return x + 1, y
	case down:
		return x, y + 1
	case left:
		return x - 1, y
	}

	return x, y
}

func makeTurn(currentDir direction, turnInstruction int) direction {
	switch currentDir {
	case up:
		if turnInstruction == 0 {
			return left
		}
		return right
	case right:
		if turnInstruction == 0 {
			return up
		}
		return down
	case down:
		if turnInstruction == 0 {
			return right
		}
		return left
	case left:
		if turnInstruction == 0 {
			return down
		}
		return up
	}

	return currentDir
}
