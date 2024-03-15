package vm_test

import (
	"math"
	"testing"
	"toolchain/instruction"
	"toolchain/vm"
)

func TestLD(t *testing.T) {
	mem := []int16{
		int16(instruction.LDI(2)), //0
		int16(instruction.LDR(3)), //1
		int16(instruction.LDM(4)), //2
		32,                        //3
		5,                         //4
		-11,                       //5
	}

	v := vm.NewVM(mem)

	v.SingleStep()
	if v.Accumulator != 2 {
		t.Fatalf("Expected '2' in accumulator got '%d'", v.Accumulator)
	}

	v.SingleStep()
	if v.Accumulator != 32 {
		t.Fatalf("Expected 32 in accumulator got '%d'", v.Accumulator)
	}

	v.SingleStep()
	if v.Accumulator != -11 {
		t.Fatalf("Expected -11 in accumulator got '%d'", v.Accumulator)
	}
}

func TestStore(t *testing.T) {
	location := uint16(0xFFEE)
	constant := int16(-2)

	mem := []int16{
		instruction.LDI(int8(constant)).ToInt16(),
		instruction.STR(2).ToInt16(),
		int16(location),
	}

	v := vm.NewVM(mem)
	v.SingleStep()
	v.SingleStep()

	result := int16(v.Memory[location])
	if result != constant {
		t.Fatalf("Got wrong constant from memory store '0x%X'", result)
	}
}

func TestJSR(t *testing.T) {
	mem := []int16{
		instruction.NewNop().ToInt16(),
		instruction.JSRI(2).ToInt16(),
		instruction.NewHalt().ToInt16(),
		instruction.ADDI(1).ToInt16(),
		instruction.NewHalt().ToInt16(),
	}

	v := vm.NewVM(mem)
	v.Run()

	if v.PC != uint16(len(mem)) {
		t.Fatalf("JSR: expected 0x%X in PC, got 0x%X", uint16(len(mem)-1), v.PC)
	}

	if v.Accumulator != 2 {
		t.Fatalf("Expected '2' in accumulator, got '%d'", v.Accumulator)
	}
}

func TestJMP(t *testing.T) {
	mem := make([]int16, 1024*64)
	mem[0x0000] = instruction.JMPR(1).ToInt16()
	mem[0x0001] = 0x0FEC
	mem[0x0FEC] = instruction.LDI(-3).ToInt16()
	mem[0x0FED] = instruction.NewHalt().ToInt16()

	v := vm.NewVM(mem)
	v.Run()

	if v.PC != 0x0FEE {
		t.Fatalf("JMP: Expected '0x0FED' in PC, got 0x%X", v.PC)
	}

	if v.Accumulator != -3 {
		t.Fatalf("JMP: Expected -3 in accumulator, got '%d'", v.Accumulator)
	}

}

func TestArithmetic(t *testing.T) {
	//Test Series: 4 + -2 NAND 110 XOR 0xFFEE - 11
	//Result should be 8

	xor_const := uint16(0xFFEE)

	mem := []int16{
		instruction.LDI(4).ToInt16(),
		instruction.ADDI(-2).ToInt16(),
		instruction.NANDI(110).ToInt16(),
		instruction.XORR(6).ToInt16(),
		instruction.SUBI(11).ToInt16(),
		instruction.NewHalt().ToInt16(),
		int16(xor_const),
	}

	v := vm.NewVM(mem)
	v.Run()

	if v.PC != 6 {
		t.Fatalf("Failed to halt")
	}

	if v.Accumulator != 8 {
		t.Fatalf("Expected 32 in accumulator got '%d'", v.Accumulator)
	}
}

func TestInvalidInstructionTrap(t *testing.T) {
	mem := make([]int16, 1024*64)
	for i := 0; i < 1024*64; i++ {
		mem[i] = int16(instruction.NewNop())
	}

	mem[0xFFEE] = math.MaxInt16 //Insert invalid instructon

	v := vm.NewVM(mem)
	for i := 0; i <= 0xFFEE; i++ {
		v.SingleStep()
	}

	if v.PC != 2 {
		t.Fatalf("Failed to trigger invalid instruction trap, pc at address '0x%X'", v.PC)
	}

	if uint16(v.Accumulator) != 0xFFEF {
		t.Fatalf("Failed to load accumulator on invalid instruction trap, acc='0x%X'", v.Accumulator)
	}
}
