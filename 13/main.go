package main

import (
	"fmt"
)

const (
	empty = iota
	wall
	block
	paddle
	ball
)

func compare(a, b int) int {
	if a == b {
		return 0
	} else if a > b {
		return -1
	} else {
		return 1
	}
}

func main() {
	grid := [23][45]int{}

	paintGrid := func(score int) {
		for _, line := range grid {
			for _, col := range line {
				switch col {
				case empty:
					fmt.Printf(" ")
				case wall:
					fmt.Printf("|")
				case block:
					fmt.Printf("X")
				case paddle:
					fmt.Printf("_")
				case ball:
					fmt.Printf("o")
				}
			}
			fmt.Printf("\n")
		}
		fmt.Printf("score: %d \n\n", score)
	}

	var (
		score   int
		ballX   int
		paddleX int
	)

	// run intcode async
	icc := NewComputer(program)
	output := make(chan int, 10)
	inputReq := make(chan struct{})
	done := make(chan struct{})
	go func() {
		icc.Run(inputReq, output)
		close(done)
	}()

	for {
		select {
		case <-done: // program done
			paintGrid(score)
			return
		case <-inputReq: // program wants input
			diff := compare(paddleX, ballX)
			icc.SendInput(diff)

		case x := <-output: // progam sending output
			y := <-output
			tileType := <-output
			if (x == -1) && (y == 0) { // getting score command
				score = tileType
				continue
			}
			grid[y][x] = tileType
			if tileType == ball {
				ballX = x
			} else if tileType == paddle {
				paddleX = x
			}
		}
	}
}
