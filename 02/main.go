package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const input = `1,0,0,3,1,1,2,3,1,3,4,3,1,5,0,3,2,9,1,19,1,19,5,23,1,23,5,27,2,27,10,31,1,31,9,35,1,35,5,39,1,6,39,43,2,9,43,47,1,5,47,51,2,6,51,55,1,5,55,59,2,10,59,63,1,63,6,67,2,67,6,71,2,10,71,75,1,6,75,79,2,79,9,83,1,83,5,87,1,87,9,91,1,91,9,95,1,10,95,99,1,99,13,103,2,6,103,107,1,107,5,111,1,6,111,115,1,9,115,119,1,119,9,123,2,123,10,127,1,6,127,131,2,131,13,135,1,13,135,139,1,9,139,143,1,9,143,147,1,147,13,151,1,151,9,155,1,155,13,159,1,6,159,163,1,13,163,167,1,2,167,171,1,171,13,0,99,2,0,14,0`

func main() {

	for noun := 0; noun < 100; noun++ {
		for verb := 0; verb < 100; verb++ {
			icc := NewIntCodeComputer(input, noun, verb)
			res := icc.Run()
			if res == 19690720 {
				log.Printf("noun = '%d', verb = '%d' => '%d'", noun, verb, res)
			}
		}
	}
}

type IntCodeComputer struct {
	Memory []int
}

func (icc IntCodeComputer) Run() int {
	var step int

	for instructionPtr := 0; instructionPtr < len(icc.Memory); instructionPtr += step {
		opCode := icc.Memory[instructionPtr]

		if opCode == 99 {
			return icc.Memory[0]
		}

		instruction, ok := Instructions[opCode]
		if !ok {
			fmt.Printf("encountered unknown OpCode '%d', existing", opCode)
		}

		step = instruction.NumValues + 1

		args := icc.Memory[instructionPtr+1 : instructionPtr+step]

		instruction.Fn(icc.Memory, args)
	}

	log.Fatal("didn't find opcode 99")
	return 0
}

func NewIntCodeComputer(input string, noun int, verb int) *IntCodeComputer {
	ss := strings.Split(input, ",")
	var ii []int

	for _, s := range ss {
		i, err := strconv.Atoi(s)
		if err != nil {
			log.Fatal(err)
		}
		ii = append(ii, i)
	}

	// Replace
	ii[1] = noun
	ii[2] = verb

	return &IntCodeComputer{
		Memory: ii,
	}
}

type InstructionFunc func(memory []int, params []int)

type Instruction struct {
	Fn        InstructionFunc
	NumValues int
}

// Instructions is lookup opcode -> Instruction
var Instructions = map[int]Instruction{
	1: Instruction{
		Fn: func(memory []int, params []int) {
			t1 := memory[params[0]]
			t2 := memory[params[1]]

			memory[params[2]] = t1 + t2
		},
		NumValues: 3,
	},
	2: Instruction{
		Fn: func(memory []int, params []int) {
			f1 := memory[params[0]]
			f2 := memory[params[1]]

			memory[params[2]] = f1 * f2
		},
		NumValues: 3,
	},
	99: Instruction{
		Fn: func(memory []int, params []int) {
			log.Println(memory)
			os.Exit(0)
		},
		NumValues: 0,
	},
}
