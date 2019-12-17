package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

const intcode = `1102,34463338,34463338,63,1007,63,34463338,63,1005,63,53,1101,3,0,1000,109,988,209,12,9,1000,209,6,209,3,203,0,1008,1000,1,63,1005,63,65,1008,1000,2,63,1005,63,904,1008,1000,0,63,1005,63,58,4,25,104,0,99,4,0,104,0,99,4,17,104,0,99,0,0,1102,0,1,1020,1102,29,1,1001,1101,0,28,1016,1102,1,31,1011,1102,1,396,1029,1101,26,0,1007,1101,0,641,1026,1101,466,0,1023,1101,30,0,1008,1102,1,22,1003,1101,0,35,1019,1101,0,36,1018,1102,1,37,1012,1102,1,405,1028,1102,638,1,1027,1102,33,1,1000,1102,1,27,1002,1101,21,0,1017,1101,0,20,1015,1101,0,34,1005,1101,0,23,1010,1102,25,1,1013,1101,39,0,1004,1101,32,0,1009,1101,0,38,1006,1101,0,473,1022,1102,1,1,1021,1101,0,607,1024,1102,1,602,1025,1101,24,0,1014,109,22,21108,40,40,-9,1005,1013,199,4,187,1105,1,203,1001,64,1,64,1002,64,2,64,109,-17,2102,1,4,63,1008,63,32,63,1005,63,229,4,209,1001,64,1,64,1105,1,229,1002,64,2,64,109,9,21108,41,44,1,1005,1015,245,1105,1,251,4,235,1001,64,1,64,1002,64,2,64,109,4,1206,3,263,1105,1,269,4,257,1001,64,1,64,1002,64,2,64,109,-8,21102,42,1,5,1008,1015,42,63,1005,63,291,4,275,1105,1,295,1001,64,1,64,1002,64,2,64,109,-13,1208,6,22,63,1005,63,313,4,301,1105,1,317,1001,64,1,64,1002,64,2,64,109,24,21107,43,44,-4,1005,1017,339,4,323,1001,64,1,64,1105,1,339,1002,64,2,64,109,-5,2107,29,-8,63,1005,63,361,4,345,1001,64,1,64,1105,1,361,1002,64,2,64,109,-4,2101,0,-3,63,1008,63,32,63,1005,63,387,4,367,1001,64,1,64,1106,0,387,1002,64,2,64,109,13,2106,0,3,4,393,1001,64,1,64,1105,1,405,1002,64,2,64,109,-27,2102,1,2,63,1008,63,35,63,1005,63,425,1105,1,431,4,411,1001,64,1,64,1002,64,2,64,109,5,1202,2,1,63,1008,63,31,63,1005,63,455,1001,64,1,64,1106,0,457,4,437,1002,64,2,64,109,19,2105,1,1,1001,64,1,64,1105,1,475,4,463,1002,64,2,64,109,-6,21102,44,1,1,1008,1017,45,63,1005,63,499,1001,64,1,64,1105,1,501,4,481,1002,64,2,64,109,6,1205,-2,513,1106,0,519,4,507,1001,64,1,64,1002,64,2,64,109,-17,1207,-1,40,63,1005,63,537,4,525,1106,0,541,1001,64,1,64,1002,64,2,64,109,-8,1201,9,0,63,1008,63,38,63,1005,63,567,4,547,1001,64,1,64,1106,0,567,1002,64,2,64,109,-3,2101,0,6,63,1008,63,32,63,1005,63,591,1001,64,1,64,1105,1,593,4,573,1002,64,2,64,109,22,2105,1,8,4,599,1106,0,611,1001,64,1,64,1002,64,2,64,109,8,1206,-4,625,4,617,1105,1,629,1001,64,1,64,1002,64,2,64,109,3,2106,0,0,1106,0,647,4,635,1001,64,1,64,1002,64,2,64,109,-29,2107,27,9,63,1005,63,667,1001,64,1,64,1106,0,669,4,653,1002,64,2,64,109,7,1207,-4,28,63,1005,63,689,1001,64,1,64,1105,1,691,4,675,1002,64,2,64,109,-7,2108,30,3,63,1005,63,711,1001,64,1,64,1105,1,713,4,697,1002,64,2,64,109,17,21101,45,0,-5,1008,1010,45,63,1005,63,735,4,719,1106,0,739,1001,64,1,64,1002,64,2,64,109,-9,1202,-2,1,63,1008,63,39,63,1005,63,765,4,745,1001,64,1,64,1106,0,765,1002,64,2,64,109,10,21101,46,0,-5,1008,1011,48,63,1005,63,785,1106,0,791,4,771,1001,64,1,64,1002,64,2,64,109,-10,1208,0,36,63,1005,63,811,1001,64,1,64,1105,1,813,4,797,1002,64,2,64,109,7,1205,8,827,4,819,1105,1,831,1001,64,1,64,1002,64,2,64,109,-15,2108,27,4,63,1005,63,853,4,837,1001,64,1,64,1106,0,853,1002,64,2,64,109,14,1201,-3,0,63,1008,63,30,63,1005,63,877,1001,64,1,64,1106,0,879,4,859,1002,64,2,64,109,11,21107,47,46,-5,1005,1018,899,1001,64,1,64,1105,1,901,4,885,4,64,99,21101,0,27,1,21101,0,915,0,1105,1,922,21201,1,31783,1,204,1,99,109,3,1207,-2,3,63,1005,63,964,21201,-2,-1,1,21101,0,942,0,1106,0,922,21201,1,0,-1,21201,-2,-3,1,21101,0,957,0,1105,1,922,22201,1,-1,-2,1106,0,968,22102,1,-2,-2,109,-3,2105,1,0`

func main() {
	icc := NewIntCodeComputer(intcode)
	icc.input <- 2
	icc.Run()
}

// IntCodeComputer describes a computer running intcode programs
type IntCodeComputer struct {
	memory         []int
	relativeBase   int
	input          chan int
	instructionPtr int
	running        bool
}

// Read returns current memory position, based on parameter mode
func (icc *IntCodeComputer) Read(mode int) int {
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
func (icc *IntCodeComputer) Write(mode int, val int) {
	// read the value at memory position
	pos := icc.Read(1)

	if mode == 2 {
		pos = pos + icc.relativeBase
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
			icc.Write(paramModes[0], <-icc.input)
		},
		// Write output
		4: func(paramModes []int) {
			out := icc.Read(paramModes[0])
			fmt.Println(out)
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

// Run starts the IntCodeComputer, returns a channel that closes when program halts
func (icc *IntCodeComputer) Run() {
	icc.running = true

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
}

// NewIntCodeComputer initializes a new IntCodeComputer
func NewIntCodeComputer(intcode string) *IntCodeComputer {
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
