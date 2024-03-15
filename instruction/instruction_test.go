package instruction_test

import (
	"testing"
	"toolchain/instruction"
)

func TestNewInstruction(t *testing.T) {
	ins := instruction.NewInstruction(instruction.LD, -23, true, false)
	if ins.Opcode() != instruction.LD {
		t.Fatalf("Opcode Write. Expected LDA (%016b) got %016b for %016b", instruction.LD, ins.Opcode(), ins)
	}

	if ins.Immediate() != int16(-23) {
		t.Fatalf("Immediate Write. Expected -23 (%d) got %d", int16(-23), ins.Immediate())
	}

	if !ins.IsImmediate() {
		t.Fatalf("Addressing Bit Write. Expected 'true' for IsImmediate")
	}

	if ins.IsIndirect() {
		t.Fatalf("Addressing Bit Write. Expected 'false' for IsIndirect")
	}

	ins = instruction.NewInstruction(instruction.LD, -11, false, true)
	if ins.IsImmediate() {
		t.Fatalf("Addressing Bit Write. Expected 'true' for IsImmediate")
	}

	if !ins.IsIndirect() {
		t.Fatalf("Addressing Bit Write. Expected 'false' for IsIndirect")
	}
}

func TestLDM(t *testing.T) {
	i := instruction.LDM(3)
	if i.IsImmediate() {
		t.Fatalf("LDM tested as immediate")
	}

	if i.IsDirect() {
		t.Fatalf("LDM tested as direct")
	}

	if !i.IsIndirect() {
		t.Fatalf("LDM did not test as indirect")
	}
}
