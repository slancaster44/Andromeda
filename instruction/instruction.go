package instruction

import "fmt"

/* Instruction Format
 * 16 bits wide, numbered [0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15]
 * Imm -> Bits 8-15 (inclusive), sign extended to 16 bits
 * Immediate Bit -> 7; Set indicates immedate is used as such, not as a register address
 * Indirect Bit -> 6; Set indicates indirect register access
 * Opcode -> [0-5] (inclusive); Indicates operation
 */

type Instruction uint16

func NewInstruction(opcode uint8, immediate int8, is_imm bool, is_ind bool) Instruction {
	out := Instruction(opcode) << 10
	out |= 0x00FF & Instruction(immediate)

	if is_imm {
		out |= 0x0100
	} else if is_ind {
		out |= 0x0200
	}

	return out
}

func (i Instruction) Opcode() uint8 {
	return uint8(i >> 10)
}

func (i Instruction) Immediate() int16 {
	value := int16(int8(0x00FF & i))
	return value
}

func (i Instruction) Address() uint16 {
	return uint16(i&0x00FF) + 0xFF00
}

func (i Instruction) IsImmediate() bool {
	return (0x0100 & i) != 0
}

func (i Instruction) IsDirect() bool {
	return !i.IsIndirect()
}

func (i Instruction) IsIndirect() bool {
	return (0x0200 & i) != 0
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

var validInstructions []Instruction = []Instruction{
	LDI(0), LDR(0), LDM(0),
	STR(0), STM(0),
	ADDI(0), ADDR(0), ADDM(0),
	NANDI(0), NANDR(0), NANDM(0),
	XORI(0), XORR(0), XORM(0),
	SUBI(0), SUBR(0), SUBM(0),
	JSRI(0), JSRR(0), JSRM(0),
	JMPI(0), JMPR(0), JMPM(0),
	JNZI(0), JNZR(0), JNZM(0),
	JNSI(0), JNSR(0), JNSM(0),
	NewHalt(), NewNop(),
}

//Note: The left byte of an instruction contains
// all the opcode and addressing information, but
// not the immediate value. This way, we can test
// that both the opcode, and the addressing mode
// is valid without considering the immediate
func (i Instruction) IsValid() bool {
	for _, x := range validInstructions {
		if x.LeftByte() == i.LeftByte() {
			return true
		}
	}

	return false
}

func (i Instruction) String() string {
	op_str := OpcodeStringMap[i.Opcode()]
	if i.Opcode() == HALT || i.Opcode() == NOP {
		return op_str
	} else if i.IsImmediate() { //eg: jsr	2
		return fmt.Sprintf("%s\t%d", op_str, i.Immediate())
	} else if i.IsDirect() { //eg: jsr %0x11
		return fmt.Sprintf("%s\t%%0x%X", op_str, i.Immediate())
	} else if i.IsIndirect() { //eg: jsr (%0x11)
		return fmt.Sprintf("%s\t(%%0x%X)", op_str, i.Immediate())
	}

	panic("Cannot format invalid instruction")
}
