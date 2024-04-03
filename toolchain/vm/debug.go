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
	fmt.Printf("Accumulator: %d (0x%04X)\n", v.Accumulator, uint16(v.Accumulator))
	fmt.Printf("Program Counter: 0x%04X (%s)\n", v.PC, instruction.Instruction(v.Memory[v.PC]))
	fmt.Printf("Halt Flip Flop: %t\n", v.HFF)
}

func getStartStopArgs(input string) (int, int) {
	start := getNumberAtIndex(input, 1)
	stop := getNumberAtIndex(input, 2)
	return start, stop
}

func getNumberAtIndex(input string, index int) int {
	splitInput := strings.Split(input, " ")
	if len(splitInput) <= index {
		fmt.Println("Not enough arguments")
		return -1
	}

	value := strings.TrimSpace(splitInput[index])
	integer, err := strconv.ParseUint(value, 16, 16)
	if err != nil {
		fmt.Print(err)
		return -1
	}

	return int(integer)
}

func (v *VM) PrintMemory(start, stop uint64) {
	fmt.Print("       |")
	for i := 0; i < 16; i++ {
		fmt.Printf("  0%X  |", i)
	}

	for i := start; i <= stop; i++ {
		if i%16 == 0 {
			fmt.Printf("\n0x%04X:\t", uint16(i))
		}

		fmt.Printf("0x%04X ", uint16(v.Memory[i]))
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
			v.PrintMemory(uint64(start), uint64(stop))
		} else if input[0] == 'j' {
			val := getNumberAtIndex(input, 1)
			if val != -1 {
				v.PC = uint16(val)
			}
		} else if input[0] == 'f' {
			val := getNumberAtIndex(input, 1)
			if val != -1 {
				for i := 0; i < val; i++ {
					v.SingleStep()
				}
			}
		} else if input[0] == 'q' {
			break
		} else {
			fmt.Println("Unknown command")
		}
	}
}
