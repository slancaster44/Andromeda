package vm

import (
	"andromeda/toolchain/instruction"
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

func getInput() string {
	fmt.Print(">> ")
	input, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	return input
}

func (v *VM) PrintStatus() {
	fmt.Printf("Accumulator: 0x%04X\n", v.Accumulator)
	fmt.Printf("Program Counter: 0x%04X (%s)\n", v.PC, instruction.Instruction(v.Memory[v.PC]))
	fmt.Printf("Halt Flip Flop: %t\n", v.HFF)
}

func getStartStopArgs(input string) (uint64, uint64) {
	splitInput := strings.Split(input, " ")
	if len(splitInput) != 3 {
		fmt.Println("Expected two arguments, a start and stop address")
		return 0, 0
	}

	startStr := strings.TrimSpace(splitInput[1])
	start, err := strconv.ParseUint(startStr, 16, 16)
	if err != nil {
		fmt.Print(err)
		return 0, 0
	}

	stopStr := strings.TrimSpace(splitInput[2])
	stop, err := strconv.ParseUint(stopStr, 16, 16)
	if err != nil {
		fmt.Print(err)
		return 0, 0
	}

	return start, stop
}

func (v *VM) PrintMemory(start, stop uint64) {
	fmt.Print("       |")
	for i := 0; i < 16; i++ {
		fmt.Printf("  0%X  |", i)
	}

	for i := start; i <= stop; i++ {
		if i%16 == 0 {
			fmt.Printf("\n0x%04X:\t", i)
		}

		fmt.Printf("0x%04X ", v.Memory[i])
	}
	fmt.Printf("\n\n")
}

func (v *VM) Debug() {
	for {
		v.PrintStatus()

		input := getInput()
		if input[0] == 's' {
			v.SingleStep()
		} else if input[0] == 'm' {
			start, stop := getStartStopArgs(input)
			v.PrintMemory(start, stop)
		} else if input[0] == 'q' {
			break
		} else {
			fmt.Println("Unknown command")
		}
	}
}
