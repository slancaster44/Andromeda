package instruction_test

import (
	"andromeda/toolchain/instruction"
	"testing"
)

func TestNewInstruction(t *testing.T) {
	ins := instruction.NewInstruction(instruction.LD, instruction.AM_IND, -23)
	if ins.Opcode() != instruction.LD {
		t.Errorf("Opcode Write. Expected LDA (%016b) got %016b for %016b", instruction.LD, ins.Opcode(), ins)
	}

	if ins.Immediate() != int16(-23) {
		t.Errorf("Immediate Write. Expected -23 (%d) got %d", int16(-23), ins.Immediate())
	}

	if ins.AddressingMode() != instruction.AM_IND {
		t.Errorf("Expected (%b) for addressing mode got %b\n", instruction.AM_IND, ins.AddressingMode())
	}

}
