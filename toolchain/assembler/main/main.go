package main

import (
	"andromeda/toolchain/assembler/assembler"
	"andromeda/toolchain/assembler/tokenizer"
	"encoding/binary"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func readFile(filename string) string {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return string(bytes)
}

func main() {

	if len(os.Args) != 4 {
		fmt.Println("Usage: [filename] [origin] [rom_size_in_bytes]")
		os.Exit(1)
	}

	filename := os.Args[1]
	text := strings.ToLower(readFile(filename))

	tokens := tokenizer.Tokenize(text)
	output := assembler.Assemble(tokens)
	code := output.Code

	origin, err := strconv.ParseUint(os.Args[2], 10, 16)
	if err != nil {
		panic(err)
	}

	top, err := strconv.ParseUint(os.Args[3], 10, 16)
	if err != nil {
		panic(err)
	}

	var fileBuffer []byte
	for i := origin; i <= top; i++ {
		buffer := make([]byte, 2)
		binary.BigEndian.PutUint16(buffer, code[i])
		fileBuffer = append(fileBuffer, buffer...)
	}

	err = os.WriteFile("a.bin", fileBuffer, 0664)
	if err != nil {
		panic(err)
	}
}
