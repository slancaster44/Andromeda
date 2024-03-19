package instruction

import "fmt"

/* Instruction Format
 * 16 bits wide, numbered [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15]
 * Imm -> Bits 8-15 (inclusive), sign extended to 16 bits
 * Addressing Mode [6, 7], Addressing mode of the instruction
 * Opcode -> [0-5] (inclusive); Indicates operation
 */

type Instruction uint16

func NewInstruction(opcode uint8, addressing_mode byte, immediate int) Instruction {
	return (Instruction((opcode<<3)|addressing_mode) << 8) | Instruction(uint8(immediate))
}

func NewHalt() Instruction {
	return NewInstruction(HALT, 0, 0)
}

func NewNop() Instruction {
	return NewInstruction(NOP, 0, 0)
}

func (i Instruction) Opcode() uint8 {
	return uint8(i >> 11)
}

func (i Instruction) Immediate() int16 {
	value := int16(int8(0x00FF & i))
	return value
}

func (i Instruction) Address() uint16 {
	return 0xFF00 + uint16(i.RightByte())
}

func (i Instruction) IsJmp() bool {
	return i.Opcode() == JNS || i.Opcode() == JNZ || i.Opcode() == JMP || i.Opcode() == JSR
}

func (i Instruction) AddressingMode() byte {
	return i.LeftByte() & 0b111
}

func (i Instruction) LeftByte() uint8 {
	return uint8((i >> 8) & 0x00FF)
}

func (i Instruction) RightByte() uint8 {
	return uint8(i & 0x00FF)
}

func (i Instruction) ToInt16() int16 {
	return int16(i)
}

func (i Instruction) IsValid() bool {
	if i.Opcode() == HALT {
		return i == NewInstruction(HALT, 0, 0)
	} else if i.Opcode() == NOP {
		return i == NewInstruction(NOP, 0, 0)
	}

	maxOp := uint8(0)
	for _, op := range StringOpcodeMap {
		if maxOp < op {
			maxOp = op
		}
	}

	maxAm := uint8(0)
	for _, am := range StringAddressingMap {
		if maxAm < am {
			maxAm = am
		}
	}

	maxIns := NewInstruction(maxOp, maxAm, 0)

	return i.LeftByte() <= maxIns.LeftByte()
}

func (i Instruction) String() string {
	if i.Opcode() == HALT {
		return "hlt"
	} else if i.Opcode() == NOP {
		return "nop"
	} else if i.AddressingMode() == AM_IMM {
		return fmt.Sprintf("%s.%s %d", OpcodeStringMap[i.Opcode()], AddressingStringMap[i.AddressingMode()], i.Immediate())
	}

	return fmt.Sprintf("%s.%s 0x%02X", OpcodeStringMap[i.Opcode()], AddressingStringMap[i.AddressingMode()], i.Immediate())
}
