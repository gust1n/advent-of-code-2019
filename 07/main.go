package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
)

const intcode = `3,8,1001,8,10,8,105,1,0,0,21,34,47,72,81,94,175,256,337,418,99999,3,9,102,3,9,9,1001,9,3,9,4,9,99,3,9,101,4,9,9,1002,9,5,9,4,9,99,3,9,1001,9,5,9,1002,9,5,9,1001,9,2,9,1002,9,5,9,101,5,9,9,4,9,99,3,9,102,2,9,9,4,9,99,3,9,1001,9,4,9,102,4,9,9,4,9,99,3,9,102,2,9,9,4,9,3,9,101,2,9,9,4,9,3,9,1002,9,2,9,4,9,3,9,102,2,9,9,4,9,3,9,1001,9,2,9,4,9,3,9,101,2,9,9,4,9,3,9,101,2,9,9,4,9,3,9,102,2,9,9,4,9,3,9,101,2,9,9,4,9,3,9,1001,9,1,9,4,9,99,3,9,102,2,9,9,4,9,3,9,1001,9,2,9,4,9,3,9,101,1,9,9,4,9,3,9,1001,9,2,9,4,9,3,9,102,2,9,9,4,9,3,9,101,2,9,9,4,9,3,9,1001,9,1,9,4,9,3,9,1001,9,2,9,4,9,3,9,1001,9,2,9,4,9,3,9,1002,9,2,9,4,9,99,3,9,101,2,9,9,4,9,3,9,101,2,9,9,4,9,3,9,101,2,9,9,4,9,3,9,1001,9,1,9,4,9,3,9,1002,9,2,9,4,9,3,9,101,1,9,9,4,9,3,9,102,2,9,9,4,9,3,9,101,1,9,9,4,9,3,9,1001,9,2,9,4,9,3,9,102,2,9,9,4,9,99,3,9,1001,9,1,9,4,9,3,9,102,2,9,9,4,9,3,9,101,2,9,9,4,9,3,9,102,2,9,9,4,9,3,9,101,1,9,9,4,9,3,9,102,2,9,9,4,9,3,9,1001,9,2,9,4,9,3,9,101,2,9,9,4,9,3,9,1002,9,2,9,4,9,3,9,101,2,9,9,4,9,99,3,9,102,2,9,9,4,9,3,9,1001,9,2,9,4,9,3,9,1002,9,2,9,4,9,3,9,101,2,9,9,4,9,3,9,101,2,9,9,4,9,3,9,1002,9,2,9,4,9,3,9,1001,9,1,9,4,9,3,9,101,2,9,9,4,9,3,9,1001,9,1,9,4,9,3,9,1002,9,2,9,4,9,99`

func main() {
	var highestOut int

	for _, perm := range makePermutations([]int{5, 6, 7, 8, 9}) {
		// Create five amplifiers
		amps := make([]*IntCodeComputer, 5)
		for i := 0; i < 5; i++ {
			icc := NewIntCodeComputer(intcode)
			amps[i] = icc
		}

		var wg sync.WaitGroup // keep track of running amps

		// run each amp and connect it's output to the next input
		for i := range amps {
			amp := amps[i]
			nextIdx := (i + 1) % len(amps)

			// run the amp
			wg.Add(1)
			go func() {
				amp.Run(amps[nextIdx].input) // write to the next (wrapping) channels input
				wg.Done()
			}()

			// send phase as first input
			phase := perm[i]
			amp.input <- phase
		}

		// send initial input to the first amp
		amps[0].input <- 0

		wg.Wait() // wait for all amps to finish

		lastOut := amps[4].highestOut
		if lastOut > highestOut {
			highestOut = lastOut
		}
	}

	log.Println("highest", highestOut)
}

// IntCodeComputer describes a computer running intcode programs
type IntCodeComputer struct {
	memory         []int
	input          chan int
	output         chan int
	instructionPtr int
	running        bool
	highestOut     int
}

// Read returns current memory position, based on parameter mode
func (icc *IntCodeComputer) Read(mode int) int {
	val := icc.memory[icc.instructionPtr]

	if mode == 0 {
		val = icc.memory[val]
	}

	icc.instructionPtr++
	return val
}

// Write writes passed value to position in memory
func (icc *IntCodeComputer) Write(pos int, val int) {
	icc.memory[pos] = val
}

// Instructions describes the IntCodeComputers set of instructions
func (icc *IntCodeComputer) Instructions() map[int]InstructionFunc {
	return map[int]InstructionFunc{
		// Addition
		1: func(paramModes []int) {
			t1 := icc.Read(paramModes[0])
			t2 := icc.Read(paramModes[1])
			pos := icc.Read(1)

			icc.Write(pos, t1+t2)
		},
		// Multiplication
		2: func(paramModes []int) {
			t1 := icc.Read(paramModes[0])
			t2 := icc.Read(paramModes[1])
			pos := icc.Read(1)

			icc.Write(pos, t1*t2)
		},
		// Read input and write to memory
		3: func(paramModes []int) {
			pos := icc.Read(1)
			icc.Write(pos, <-icc.input)
		},
		// Write output
		4: func(paramModes []int) {
			out := icc.Read(paramModes[0])
			if out > icc.highestOut {
				icc.highestOut = out
			}
			icc.output <- out
		},
		// Jump if true
		5: func(paramModes []int) {
			jump := icc.Read(paramModes[0])
			loc := icc.Read(paramModes[1])
			if jump != 0 {
				icc.instructionPtr = loc
			}
		},
		// Jump if false
		6: func(paramModes []int) {
			jump := icc.Read(paramModes[0])
			loc := icc.Read(paramModes[1])
			if jump == 0 {
				icc.instructionPtr = loc
			}
		},
		// Less than
		7: func(paramModes []int) {
			var (
				p1  = icc.Read(paramModes[0])
				p2  = icc.Read(paramModes[1])
				pos = icc.Read(1)
				res int
			)
			if p1 < p2 {
				res = 1
			}
			icc.Write(pos, res)
		},
		// Equals
		8: func(paramModes []int) {
			var (
				p1  = icc.Read(paramModes[0])
				p2  = icc.Read(paramModes[1])
				pos = icc.Read(1)
				res int
			)
			if p1 == p2 {
				res = 1
			}
			icc.Write(pos, res)
		},
		// Halt
		99: func(paramModes []int) {
			icc.running = false
		},
	}
}

// Run starts the IntCodeComputer, returns a channel that closes when program halts
func (icc *IntCodeComputer) Run(outputCh chan int) {
	icc.running = true
	icc.output = outputCh

	for icc.running {
		opCode, paramModes := parseInstruction(icc.memory[icc.instructionPtr])

		instruction, ok := icc.Instructions()[opCode]
		if !ok {
			log.Fatalf("encountered unknown OpCode '%d', existing\n\n", opCode)
		}

		// found the opcode, move the instruction pointer
		icc.instructionPtr++

		// run the instruction
		instruction(paramModes)
	}
}

// NewIntCodeComputer initializes a new IntCodeComputer
func NewIntCodeComputer(intcode string) *IntCodeComputer {
	ss := strings.Split(intcode, ",")
	var ii []int

	for _, s := range ss {
		i, err := strconv.Atoi(s)
		if err != nil {
			log.Fatal(err)
		}
		ii = append(ii, i)
	}

	return &IntCodeComputer{
		memory: ii,
		input:  make(chan int, 1),
	}
}

// InstructionFunc describes a function for executing an instruction
type InstructionFunc func(paramModes []int)

// parseInstruction returns opcode and each arguments mode (ltr)
func parseInstruction(input int) (int, []int) {
	// Zero-pad to make sure we've got all five theoretical positions
	str := fmt.Sprintf("%05d", input)

	// parse opcode (last two digits)
	opcode, err := strconv.Atoi(fmt.Sprintf("%s%s", string(str[3]), string(str[4])))
	if err != nil {
		log.Fatal(err)
	}

	// determine arg modes
	paramModes := make([]int, 3)

	// reverse the slice as it's read RTL
	paramModes[0] = parseParamMode(str[2])
	paramModes[1] = parseParamMode(str[1])
	paramModes[2] = parseParamMode(str[0])

	return opcode, paramModes
}

func parseParamMode(b byte) int {
	i, err := strconv.Atoi(string(b))
	if err != nil {
		log.Fatal(err)
	}

	return i
}

func makePermutations(arr []int) [][]int {
	var helper func([]int, int)
	res := [][]int{}

	helper = func(arr []int, n int) {
		if n == 1 {
			tmp := make([]int, len(arr))
			copy(tmp, arr)
			res = append(res, tmp)
		} else {
			for i := 0; i < n; i++ {
				helper(arr, n-1)
				if n%2 == 1 {
					tmp := arr[i]
					arr[i] = arr[n-1]
					arr[n-1] = tmp
				} else {
					tmp := arr[0]
					arr[0] = arr[n-1]
					arr[n-1] = tmp
				}
			}
		}
	}
	helper(arr, len(arr))
	return res
}
