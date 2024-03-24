package vm_test

import (
	"andromeda/toolchain/instruction"
	"andromeda/toolchain/vm"
	"math"
	"testing"
)

func TestLD(t *testing.T) {
	mem := make([]int16, 64*1024) /* TODO: Auto increment & decrement */
	mem[0] = instruction.NewInstruction(instruction.LD, instruction.AM_IMM, 2).ToInt16()
	mem[1] = instruction.NewInstruction(instruction.LD, instruction.AM_DIR, 3).ToInt16()
	mem[2] = instruction.NewInstruction(instruction.LD, instruction.AM_IND, 4).ToInt16()
	mem[5] = -11
	mem[0xFF03] = 32
	mem[0xFF04] = 5

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
	constant := -2

	mem := make([]int16, 64*1024)
	mem[0] = instruction.NewInstruction(instruction.LD, instruction.AM_IMM, constant).ToInt16()
	mem[1] = instruction.NewInstruction(instruction.STORE, instruction.AM_IND, 2).ToInt16()
	mem[0xFF02] = int16(location)

	v := vm.NewVM(mem)
	v.SingleStep()
	v.SingleStep()

	result := int(v.Memory[location])
	if result != constant {
		t.Fatalf("Got wrong constant from memory store, '0x%X' found at '0x%X'", result, location)
	}
}

func TestJSR(t *testing.T) {
	mem := []int16{
		instruction.NewInstruction(instruction.NOP, 0, 0).ToInt16(),
		instruction.NewInstruction(instruction.JSR, instruction.AM_IMM, 3).ToInt16(),
		instruction.NewInstruction(instruction.HALT, 0, 0).ToInt16(),
		instruction.NewInstruction(instruction.LD, instruction.AM_IMM, 1).ToInt16(),
		instruction.NewInstruction(instruction.HALT, 0, 0).ToInt16(),
	}

	v := vm.NewVM(mem)
	v.Run()

	if v.PC != uint16(len(mem)) {
		t.Fatalf("JSR: expected 0x%X in PC, got 0x%X", uint16(len(mem)), v.PC)
	}

	if v.Accumulator != 1 {
		t.Fatalf("Expected '1' in accumulator, got '%d'", v.Accumulator)
	}
}

func TestJMP(t *testing.T) {
	mem := make([]int16, 1024*64)
	mem[0x0000] = instruction.NewInstruction(instruction.JMP, instruction.AM_DIR, 1).ToInt16()
	mem[0xFF01] = 0x0FEC
	mem[0x0FEC] = instruction.NewInstruction(instruction.LD, instruction.AM_IMM, -3).ToInt16()
	mem[0x0FED] = instruction.NewInstruction(instruction.HALT, 0, 0).ToInt16()

	v := vm.NewVM(mem)
	v.Run()

	if v.PC != 0x0FEE {
		t.Fatalf("JMP: Expected '0x0FED' in PC, got 0x%X", v.PC)
	}

	if v.Accumulator != -3 {
		t.Fatalf("JMP: Expected -3 in accumulator, got '%d'", v.Accumulator)
	}
}

func TestJNZ(t *testing.T) {
	mem := make([]int16, 1024*64)
	mem[0] = instruction.NewInstruction(instruction.LD, instruction.AM_IMM, 1).ToInt16()
	mem[1] = instruction.NewInstruction(instruction.JNZ, instruction.AM_OFF, 2).ToInt16()
	mem[2] = instruction.NewInstruction(instruction.LD, instruction.AM_IMM, 300).ToInt16()
	mem[3] = instruction.NewHalt().ToInt16()

	v := vm.NewVM(mem)
	v.Run()

	if v.Accumulator != 1 {
		t.Fatalf("JNZ: Expected 1 in accumualtor got '%d'\n", v.Accumulator)
	}

}

func TestArithmetic(t *testing.T) {
	//Test Series: 4 + -2 NAND 110 XOR 0xFFEE - 11
	//Result should be 8

	mem := []int16{
		instruction.NewInstruction(instruction.LD, instruction.AM_IMM, 4).ToInt16(),
		instruction.NewInstruction(instruction.ADD, instruction.AM_IMM, -2).ToInt16(),
		instruction.NewInstruction(instruction.NAND, instruction.AM_IMM, 110).ToInt16(),
		instruction.NewInstruction(instruction.XOR, instruction.AM_IMM, 0xFFEE).ToInt16(),
		instruction.NewInstruction(instruction.SUB, instruction.AM_IMM, 11).ToInt16(),
		instruction.NewInstruction(instruction.HALT, 0, 0).ToInt16(),
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

func TestRelativeAddressing(t *testing.T) {
	mem := []int16{
		instruction.NewNop().ToInt16(),
		instruction.NewNop().ToInt16(),
		instruction.NewInstruction(instruction.LD, instruction.AM_REL, 2).ToInt16(),
		instruction.NewHalt().ToInt16(),
		-11,
	}

	v := vm.NewVM(mem)
	v.Run()

	if v.Accumulator != -11 {
		t.Fatalf("Unexpected value in accumulator got '%d' expected '%d'\n", v.Accumulator, -11)
	}
}

func TestInvalidInstructionTrap(t *testing.T) {
	mem := make([]int16, 1024*64)
	for i := 0; i < 1024*64; i++ {
		mem[i] = int16(instruction.NewInstruction(instruction.NOP, 0, 0))
	}

	mem[0xFFEE] = math.MaxInt16 //Insert invalid instruction

	v := vm.NewVM(mem)
	for i := 0; i <= 0xFFEE; i++ {
		v.SingleStep()
	}

	if v.PC != 3 {
		t.Fatalf("Failed to trigger invalid instruction trap, pc at address '0x%X'", v.PC)
	}

	if uint16(v.Accumulator) != 0xFFEF {
		t.Fatalf("Failed to load accumulator on invalid instruction trap, acc='0x%X'", v.Accumulator)
	}
}
