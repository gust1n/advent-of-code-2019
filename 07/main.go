package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
)

const intcode = `3,8,1001,8,10,8,105,1,0,0,21,34,47,72,81,94,175,256,337,418,99999,3,9,102,3,9,9,1001,9,3,9,4,9,99,3,9,101,4,9,9,1002,9,5,9,4,9,99,3,9,1001,9,5,9,1002,9,5,9,1001,9,2,9,1002,9,5,9,101,5,9,9,4,9,99,3,9,102,2,9,9,4,9,99,3,9,1001,9,4,9,102,4,9,9,4,9,99,3,9,102,2,9,9,4,9,3,9,101,2,9,9,4,9,3,9,1002,9,2,9,4,9,3,9,102,2,9,9,4,9,3,9,1001,9,2,9,4,9,3,9,101,2,9,9,4,9,3,9,101,2,9,9,4,9,3,9,102,2,9,9,4,9,3,9,101,2,9,9,4,9,3,9,1001,9,1,9,4,9,99,3,9,102,2,9,9,4,9,3,9,1001,9,2,9,4,9,3,9,101,1,9,9,4,9,3,9,1001,9,2,9,4,9,3,9,102,2,9,9,4,9,3,9,101,2,9,9,4,9,3,9,1001,9,1,9,4,9,3,9,1001,9,2,9,4,9,3,9,1001,9,2,9,4,9,3,9,1002,9,2,9,4,9,99,3,9,101,2,9,9,4,9,3,9,101,2,9,9,4,9,3,9,101,2,9,9,4,9,3,9,1001,9,1,9,4,9,3,9,1002,9,2,9,4,9,3,9,101,1,9,9,4,9,3,9,102,2,9,9,4,9,3,9,101,1,9,9,4,9,3,9,1001,9,2,9,4,9,3,9,102,2,9,9,4,9,99,3,9,1001,9,1,9,4,9,3,9,102,2,9,9,4,9,3,9,101,2,9,9,4,9,3,9,102,2,9,9,4,9,3,9,101,1,9,9,4,9,3,9,102,2,9,9,4,9,3,9,1001,9,2,9,4,9,3,9,101,2,9,9,4,9,3,9,1002,9,2,9,4,9,3,9,101,2,9,9,4,9,99,3,9,102,2,9,9,4,9,3,9,1001,9,2,9,4,9,3,9,1002,9,2,9,4,9,3,9,101,2,9,9,4,9,3,9,101,2,9,9,4,9,3,9,1002,9,2,9,4,9,3,9,1001,9,1,9,4,9,3,9,101,2,9,9,4,9,3,9,1001,9,1,9,4,9,3,9,1002,9,2,9,4,9,99`

var debug bool

func main() {
	debug = true

	var highest int
	perms := permutations([]int{0, 1, 2, 3, 4})

	for _, perm := range perms { // try all permutations
		amps := make([]*IntCodeComputer, 5)

		// Create five amplifiers
		for i := 0; i < 5; i++ {
			icc := NewIntCodeComputer(intcode)
			amps[i] = icc
		}

		var wg sync.WaitGroup
		lastOutput := make(chan int)

		// run each amp and connect it's output to the next input
		for i := range amps {
			amp := amps[i]

			var outputCh chan int
			nextIdx := (i + 1) % len(amps)
			if nextIdx == 0 {
				log.Printf("connecting amp '%d's output to last output", i)
				outputCh = lastOutput
			} else {
				log.Printf("connecting amp '%d's output to amp '%d's input", i, nextIdx)
				outputCh = amps[nextIdx].input // write to the next channels input
			}

			// run the amp
			go func() {
				wg.Add(1)
				amp.Run(outputCh)
				wg.Done()
			}()

			// send phase as first input
			phase := perm[i]
			amp.input <- phase
		}

		// send initial input to the first amp
		amps[0].input <- 0

		// wait for last output
		out := <-lastOutput
		if out > highest {
			highest = out
		}

		// if debug {
		// 	log.Printf("-> in to '%d', phase '%d', input '%d'", i, phase, input)
		// }

		// Wait for all to finish
		// wg.Wait()
	}

	log.Println("highest", highest)
}

// IntCodeComputer describes a computer running intcode programs
type IntCodeComputer struct {
	memory         []int
	input          chan int
	output         chan int
	instructionPtr int
	running        bool
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
			icc.output <- icc.Read(paramModes[0])
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
func (icc IntCodeComputer) Run(outputCh chan int) {
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
		input:  make(chan int, 15),
		// output: make(chan int, 5),
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

func permutations(arr []int) [][]int {
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
