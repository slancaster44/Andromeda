#!/usr/bin/env bash
go run ../toolchain/assembler/main/main.go ./src/main.aas 0 512
mkdir -p bin
mv *.bin ./bin/forth.bin