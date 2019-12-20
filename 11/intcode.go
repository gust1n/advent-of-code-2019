package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

// Computer describes a computer running intcode programs
type Computer struct {
	memory         []int
	relativeBase   int
	input          chan int
	staticInput    int
	output         chan int
	instructionPtr int
	running        bool
}

// InstructionFunc describes a function for executing an instruction
type InstructionFunc func(paramModes []int)

// SendInput sends input to the computer
func (icc *Computer) SendInput(i int) {
	// icc.input <- i
	icc.staticInput = i
}

// Read returns current memory position, based on parameter mode
func (icc *Computer) Read(mode int) int {
	val := icc.memory[icc.instructionPtr]

	if mode == 0 {
		val = icc.memory[val]
	} else if mode == 2 {
		val = icc.memory[val+icc.relativeBase]
	}

	icc.instructionPtr++
	return val
}

// Write writes passed value to position in memory
func (icc *Computer) Write(mode int, val int) {
	// read the value at memory position
	pos := icc.Read(1)

	if mode == 2 {
		pos = pos + icc.relativeBase
	}

	icc.memory[pos] = val
}

// Instructions describes the Computers set of instructions
func (icc *Computer) Instructions() map[int]InstructionFunc {
	return map[int]InstructionFunc{
		// Addition
		1: func(paramModes []int) {
			t1 := icc.Read(paramModes[0])
			t2 := icc.Read(paramModes[1])

			icc.Write(paramModes[2], t1+t2)
		},
		// Multiplication
		2: func(paramModes []int) {
			t1 := icc.Read(paramModes[0])
			t2 := icc.Read(paramModes[1])

			icc.Write(paramModes[2], t1*t2)
		},
		// Read input and write to memory
		3: func(paramModes []int) {
			icc.Write(paramModes[0], icc.staticInput)
			// icc.Write(paramModes[0], <-icc.input)
		},
		// Write output
		4: func(paramModes []int) {
			out := icc.Read(paramModes[0])
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
				res int
			)
			if p1 < p2 {
				res = 1
			}
			icc.Write(paramModes[2], res)
		},
		// Equals
		8: func(paramModes []int) {
			var (
				p1  = icc.Read(paramModes[0])
				p2  = icc.Read(paramModes[1])
				res int
			)
			if p1 == p2 {
				res = 1
			}
			icc.Write(paramModes[2], res)
		},
		// Adjust relative base
		9: func(paramModes []int) {
			diff := icc.Read(paramModes[0])
			icc.relativeBase = icc.relativeBase + diff
		},
		// Halt
		99: func(paramModes []int) {
			icc.running = false
		},
	}
}

// Run starts the Computer, returns a channel that closes when program halts
func (icc *Computer) Run(outputCh chan int) {
	icc.running = true
	icc.output = outputCh

	for icc.running {
		opCode, paramModes := parseInstruction(icc.memory[icc.instructionPtr])

		instruction, ok := icc.Instructions()[opCode]
		if !ok {
			log.Printf("encountered unknown OpCode '%d', existing\n\n", opCode)
			break
		}

		// found the opcode, move the instruction pointer
		icc.instructionPtr++

		// run the instruction
		instruction(paramModes)
	}
	close(outputCh)
}

// NewComputer initializes a new Computer
func NewComputer(intcode string) *Computer {
	ss := strings.Split(intcode, ",")

	// Initialize with preallocated memory
	ii := make([]int, 10000)

	// Read program into memory
	for idx, s := range ss {
		i, err := strconv.Atoi(s)
		if err != nil {
			log.Fatal(err)
		}
		ii[idx] = i
	}

	return &Computer{
		memory: ii,
		input:  make(chan int),
	}
}

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
