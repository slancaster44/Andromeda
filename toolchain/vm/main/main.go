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

	vm := vm2.NewVM(bin)
	vm.Debug()
}
