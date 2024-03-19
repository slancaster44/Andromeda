package main

import (
	"andromeda/toolchain/assembler/assemble"
	"andromeda/toolchain/assembler/link"
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
		fmt.Println("Usage: [filename] [rom_origin_in_hex] [rom_size_in_bytes]")
		os.Exit(1)
	}

	filename := os.Args[1]
	text := strings.ToLower(readFile(filename))

	tokens := tokenizer.Tokenize(text)
	objects := assemble.NewAssemblyContext().Assemble(tokens)

	lst, err := os.Create("a.lst")
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(lst)

	for _, o := range objects {
		_, err := lst.WriteString(o.String())
		if err != nil {
			panic(err)
		}
	}

	origin, err := strconv.ParseUint(os.Args[2], 16, 16)
	if err != nil {
		panic(err)
	}

	romSize, err := strconv.ParseUint(os.Args[3], 10, 16)
	if err != nil {
		panic(err)
	}

	romImage, err := link.NewLinkerContext(uint16(origin), uint16(romSize)).Link(objects)
	if err != nil {
		panic(err)
	}

	f, err := os.Create("a.bin")
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	err = binary.Write(f, binary.BigEndian, romImage)
	if err != nil {
		panic(err)
	}

}
