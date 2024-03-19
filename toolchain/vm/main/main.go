package main

import (
	vm2 "andromeda/toolchain/vm"
	"encoding/binary"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./vm [rom_file]")
	}

	bytes, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	var bin []int16
	for i := 0; i < len(bytes); i += 2 {
		bin = append(bin, int16(binary.BigEndian.Uint16(bytes[i:])))
	}

	fmt.Println("Loading Binary at 0x0000...")
	for i, v := range bin {
		if i%16 == 0 {
			fmt.Printf("\n0x%04X:\t", i)
		}

		fmt.Printf("0x%04X ", v)
	}
	fmt.Print("\n\n")

	vm := vm2.NewVM(bin)
	vm.Debug()
}
