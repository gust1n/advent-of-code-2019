package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const intcode = `3,225,1,225,6,6,1100,1,238,225,104,0,2,136,183,224,101,-5304,224,224,4,224,1002,223,8,223,1001,224,6,224,1,224,223,223,1101,72,47,225,1101,59,55,225,1101,46,75,225,1101,49,15,224,101,-64,224,224,4,224,1002,223,8,223,1001,224,5,224,1,224,223,223,102,9,210,224,1001,224,-270,224,4,224,1002,223,8,223,1001,224,2,224,1,223,224,223,101,14,35,224,101,-86,224,224,4,224,1002,223,8,223,101,4,224,224,1,224,223,223,1102,40,74,224,1001,224,-2960,224,4,224,1002,223,8,223,101,5,224,224,1,224,223,223,1101,10,78,225,1001,39,90,224,1001,224,-149,224,4,224,102,8,223,223,1001,224,4,224,1,223,224,223,1002,217,50,224,1001,224,-1650,224,4,224,1002,223,8,223,1001,224,7,224,1,224,223,223,1102,68,8,225,1,43,214,224,1001,224,-126,224,4,224,102,8,223,223,101,3,224,224,1,224,223,223,1102,88,30,225,1102,18,80,225,1102,33,28,225,4,223,99,0,0,0,677,0,0,0,0,0,0,0,0,0,0,0,1105,0,99999,1105,227,247,1105,1,99999,1005,227,99999,1005,0,256,1105,1,99999,1106,227,99999,1106,0,265,1105,1,99999,1006,0,99999,1006,227,274,1105,1,99999,1105,1,280,1105,1,99999,1,225,225,225,1101,294,0,0,105,1,0,1105,1,99999,1106,0,300,1105,1,99999,1,225,225,225,1101,314,0,0,106,0,0,1105,1,99999,108,677,677,224,102,2,223,223,1005,224,329,1001,223,1,223,1107,677,226,224,102,2,223,223,1006,224,344,1001,223,1,223,108,226,226,224,102,2,223,223,1005,224,359,1001,223,1,223,1108,677,226,224,102,2,223,223,1006,224,374,101,1,223,223,108,677,226,224,102,2,223,223,1006,224,389,1001,223,1,223,107,226,226,224,102,2,223,223,1005,224,404,1001,223,1,223,8,226,226,224,102,2,223,223,1006,224,419,101,1,223,223,1107,677,677,224,102,2,223,223,1006,224,434,1001,223,1,223,1107,226,677,224,1002,223,2,223,1006,224,449,101,1,223,223,7,677,677,224,1002,223,2,223,1006,224,464,1001,223,1,223,1108,226,677,224,1002,223,2,223,1005,224,479,1001,223,1,223,8,677,226,224,1002,223,2,223,1005,224,494,101,1,223,223,7,226,677,224,102,2,223,223,1005,224,509,101,1,223,223,1008,677,226,224,102,2,223,223,1006,224,524,101,1,223,223,8,226,677,224,1002,223,2,223,1006,224,539,1001,223,1,223,1007,677,677,224,102,2,223,223,1005,224,554,101,1,223,223,107,226,677,224,1002,223,2,223,1005,224,569,1001,223,1,223,1108,677,677,224,1002,223,2,223,1006,224,584,1001,223,1,223,1008,226,226,224,1002,223,2,223,1005,224,599,101,1,223,223,1008,677,677,224,102,2,223,223,1005,224,614,101,1,223,223,7,677,226,224,1002,223,2,223,1005,224,629,1001,223,1,223,107,677,677,224,1002,223,2,223,1006,224,644,101,1,223,223,1007,226,677,224,1002,223,2,223,1005,224,659,1001,223,1,223,1007,226,226,224,102,2,223,223,1005,224,674,101,1,223,223,4,223,99,226`

var debug bool

func main() {
	// debug = true
	icc := NewIntCodeComputer(intcode, 5)
	icc.Run()
}

// IntCodeComputer describes a computer running intcode programs
type IntCodeComputer struct {
	memory         []int
	Input          int
	instructionPtr int
}

// Read returns current memory position, based on parameter mode
func (icc *IntCodeComputer) Read(mode int) int {
	val := icc.memory[icc.instructionPtr]

	if mode == 0 {
		val = icc.memory[val]
	}
	if debug {
		log.Printf("val at '%d': '%d'\n", icc.instructionPtr, val)
	}

	icc.instructionPtr++
	return val
}

// Write writes passed value to position in memory
func (icc *IntCodeComputer) Write(pos int, val int) {
	if debug {
		log.Printf("writing to memory pos '%d': '%d'\n", pos, val)
	}
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
			icc.Write(pos, icc.Input)
		},
		// Write output
		4: func(paramModes []int) {
			val := icc.Read(paramModes[0])
			fmt.Println(val)
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
			os.Exit(0)
		},
	}
}

// Run starts the IntCodeComputer
func (icc IntCodeComputer) Run() {
	for {
		opCode, paramModes := parseInstruction(icc.memory[icc.instructionPtr])

		if debug {
			log.Printf("-> pnt at '%d', running opcode '%d' with parameter modes '%v'\n", icc.instructionPtr, opCode, paramModes)
		}

		instruction, ok := icc.Instructions()[opCode]
		if !ok {
			fmt.Printf("encountered unknown OpCode '%d', existing\n\n", opCode)
			os.Exit(1)
		}

		// found the opcode, move the instruction pointer
		icc.instructionPtr++

		instruction(paramModes)
	}
}

// NewIntCodeComputer initializes a new IntCodeComputer
func NewIntCodeComputer(intcode string, input int) *IntCodeComputer {
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
		Input:  input,
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
