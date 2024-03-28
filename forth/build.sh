go run ../toolchain/assembler/main/main.go main.aas 0 512
mkdir -p bin
mv *.bin ./bin/forth.bin
mv *.lst ./bin/forth.lst
go run ../toolchain/vm/main/main.go ./bin/forth.bin