package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const intcode = `3,225,1,225,6,6,1100,1,238,225,104,0,2,136,183,224,101,-5304,224,224,4,224,1002,223,8,223,1001,224,6,224,1,224,223,223,1101,72,47,225,1101,59,55,225,1101,46,75,225,1101,49,15,224,101,-64,224,224,4,224,1002,223,8,223,1001,224,5,224,1,224,223,223,102,9,210,224,1001,224,-270,224,4,224,1002,223,8,223,1001,224,2,224,1,223,224,223,101,14,35,224,101,-86,224,224,4,224,1002,223,8,223,101,4,224,224,1,224,223,223,1102,40,74,224,1001,224,-2960,224,4,224,1002,223,8,223,101,5,224,224,1,224,223,223,1101,10,78,225,1001,39,90,224,1001,224,-149,224,4,224,102,8,223,223,1001,224,4,224,1,223,224,223,1002,217,50,224,1001,224,-1650,224,4,224,1002,223,8,223,1001,224,7,224,1,224,223,223,1102,68,8,225,1,43,214,224,1001,224,-126,224,4,224,102,8,223,223,101,3,224,224,1,224,223,223,1102,88,30,225,1102,18,80,225,1102,33,28,225,4,223,99,0,0,0,677,0,0,0,0,0,0,0,0,0,0,0,1105,0,99999,1105,227,247,1105,1,99999,1005,227,99999,1005,0,256,1105,1,99999,1106,227,99999,1106,0,265,1105,1,99999,1006,0,99999,1006,227,274,1105,1,99999,1105,1,280,1105,1,99999,1,225,225,225,1101,294,0,0,105,1,0,1105,1,99999,1106,0,300,1105,1,99999,1,225,225,225,1101,314,0,0,106,0,0,1105,1,99999,108,677,677,224,102,2,223,223,1005,224,329,1001,223,1,223,1107,677,226,224,102,2,223,223,1006,224,344,1001,223,1,223,108,226,226,224,102,2,223,223,1005,224,359,1001,223,1,223,1108,677,226,224,102,2,223,223,1006,224,374,101,1,223,223,108,677,226,224,102,2,223,223,1006,224,389,1001,223,1,223,107,226,226,224,102,2,223,223,1005,224,404,1001,223,1,223,8,226,226,224,102,2,223,223,1006,224,419,101,1,223,223,1107,677,677,224,102,2,223,223,1006,224,434,1001,223,1,223,1107,226,677,224,1002,223,2,223,1006,224,449,101,1,223,223,7,677,677,224,1002,223,2,223,1006,224,464,1001,223,1,223,1108,226,677,224,1002,223,2,223,1005,224,479,1001,223,1,223,8,677,226,224,1002,223,2,223,1005,224,494,101,1,223,223,7,226,677,224,102,2,223,223,1005,224,509,101,1,223,223,1008,677,226,224,102,2,223,223,1006,224,524,101,1,223,223,8,226,677,224,1002,223,2,223,1006,224,539,1001,223,1,223,1007,677,677,224,102,2,223,223,1005,224,554,101,1,223,223,107,226,677,224,1002,223,2,223,1005,224,569,1001,223,1,223,1108,677,677,224,1002,223,2,223,1006,224,584,1001,223,1,223,1008,226,226,224,1002,223,2,223,1005,224,599,101,1,223,223,1008,677,677,224,102,2,223,223,1005,224,614,101,1,223,223,7,677,226,224,1002,223,2,223,1005,224,629,1001,223,1,223,107,677,677,224,1002,223,2,223,1006,224,644,101,1,223,223,1007,226,677,224,1002,223,2,223,1005,224,659,1001,223,1,223,1007,226,226,224,102,2,223,223,1005,224,674,101,1,223,223,4,223,99,226`

func main() {
	icc := NewIntCodeComputer(intcode, 1)
	icc.Run()
}

// IntCodeComputer describes a computer running intcode programs
type IntCodeComputer struct {
	memory         []int
	Input          int
	instructionPtr int
}

// Instructions describes the IntCodeComputers set of instructions
func (icc *IntCodeComputer) Instructions() map[int]Instruction {
	return map[int]Instruction{
		// Addition
		1: Instruction{
			Fn: func(params []InstructionArg) {
				t1 := params[0].MemoryValue(icc.memory)
				t2 := params[1].MemoryValue(icc.memory)

				// log.Printf("writing '%d' to position '%d'", t1+t2, params[2].Value)
				icc.memory[params[2].Value] = t1 + t2
			},
			NumValues: 3,
		},
		// Multiplication
		2: Instruction{
			Fn: func(params []InstructionArg) {
				f1 := params[0].MemoryValue(icc.memory)
				f2 := params[1].MemoryValue(icc.memory)

				icc.memory[params[2].Value] = f1 * f2
			},
			NumValues: 3,
		},
		// Read input
		3: Instruction{
			Fn: func(params []InstructionArg) {
				icc.memory[params[0].Value] = icc.Input
			},
			NumValues: 1,
		},
		// Write output
		4: Instruction{
			Fn: func(params []InstructionArg) {
				fmt.Println(params[0].MemoryValue(icc.memory))
			},
			NumValues: 1,
		},
		// Jump if true
		5: Instruction{
			Fn: func(params []InstructionArg) {
				if params[0].MemoryValue(icc.memory) != 0 {
					icc.instructionPtr = params[1].MemoryValue(icc.memory)
				}
			},
			NumValues: 2,
		},
		// Jump if false
		6: Instruction{
			Fn: func(params []InstructionArg) {
				if params[0].MemoryValue(icc.memory) == 0 {
					icc.instructionPtr = params[1].MemoryValue(icc.memory)
				}
			},
			NumValues: 2,
		},
		// Less than
		7: Instruction{
			Fn: func(params []InstructionArg) {
				var res int
				if params[0].MemoryValue(icc.memory) < params[1].MemoryValue(icc.memory) {
					res = 1
				}
				icc.memory[params[2].Value] = res
			},
			NumValues: 3,
		},
		// Equals
		8: Instruction{
			Fn: func(params []InstructionArg) {
				var res int
				if params[0].MemoryValue(icc.memory) == params[1].MemoryValue(icc.memory) {
					res = 1
				}
				icc.memory[params[2].Value] = res
			},
			NumValues: 3,
		},
		// Halt
		99: Instruction{
			Fn: func(params []InstructionArg) {
				log.Println(icc.memory)
				os.Exit(0)
			},
			NumValues: 0,
		},
	}
}

// Run starts the IntCodeComputer
func (icc IntCodeComputer) Run() {
	var step int

	for icc.instructionPtr < len(icc.memory) {
		opCode, paramModes := parseInstruction(icc.memory[icc.instructionPtr])

		if opCode == 99 {
			return
		}

		instruction, ok := icc.Instructions()[opCode]
		if !ok {
			fmt.Printf("encountered unknown OpCode '%d', existing", opCode)
		}

		step = instruction.NumValues + 1

		argValues := icc.memory[icc.instructionPtr+1 : icc.instructionPtr+step]

		// pair each arg with parameter mode
		var args []InstructionArg

		for i, argValue := range argValues {
			args = append(args, InstructionArg{
				Mode:  paramModes[i],
				Value: argValue,
			})
		}

		// log.Printf("running opcode '%d', with args '%v'\n", opCode, args)

		instruction.Fn(args)
		icc.instructionPtr += step
	}

	log.Fatal("didn't find opcode 99")
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
type InstructionFunc func(args []InstructionArg)

// Instruction describes an intcode instruction
type Instruction struct {
	Fn        InstructionFunc
	NumValues int
}

// InstructionArg describes an argument passed to an Instruction
type InstructionArg struct {
	Mode  int
	Value int
}

// MemoryValue returns the value from memory for the argument, based on parameter mode
func (a InstructionArg) MemoryValue(memory []int) int {
	if a.Mode == 1 {
		return a.Value
	}

	return memory[a.Value]
}

// parseInstruction returns opcode and each arguments mode (ltr)
func parseInstruction(input int) (int, [3]int) {
	// Zero-pad to make sure we've got all five theoretical positions
	str := fmt.Sprintf("%05d", input)

	// parse opcode (last two digits)
	opcode, err := strconv.Atoi(fmt.Sprintf("%s%s", string(str[3]), string(str[4])))
	if err != nil {
		log.Fatal(err)
	}

	// determine arg modes
	var paramModes [3]int

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
