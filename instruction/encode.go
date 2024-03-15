package instruction

//Load an immediate value into the accumulator
// lda 	<signed 8-bit value>
func LDI(value int8) Instruction {
	return NewInstruction(LD, value, true, false)
}

//Load the contents of a register into the accumulator
//lda %<signed 8-bit value>
func LDR(value int8) Instruction {
	return NewInstruction(LD, value, false, false)
}

//Load the contents of memory (pointed to by a register)
//into the accumulator
//lda (%<signed 8-bit value>)
func LDM(value int8) Instruction {
	return NewInstruction(LD, value, false, true)
}

//Store the accumulator to a register
//sta %<signed 8-bit value>
func STR(value int8) Instruction {
	return NewInstruction(STORE, value, false, false)
}

//Store accumulator into memory (pointed to by a register)
//sta (%<signed 8-bit value>)
func STM(value int8) Instruction {
	return NewInstruction(STORE, value, false, true)
}

//Add an immedate value to the accumulator
//add <signed 8-bit value>
func ADDI(value int8) Instruction {
	return NewInstruction(ADD, value, true, false)
}

//Add the contents of a register to the accumulator
//add %<signed 8-bit value>
func ADDR(value int8) Instruction {
	return NewInstruction(ADD, value, false, false)
}

//Add the contens of memory (pointed to by a register)
//add (%<signed 8-bit value>)
func ADDM(value int8) Instruction {
	return NewInstruction(ADD, value, false, true)
}

//NAND an immedate value to the accumulator
//nand <signed 8-bit value>
func NANDI(value int8) Instruction {
	return NewInstruction(NAND, value, true, false)
}

//Nand the contents of a register to the accumulator
//nand %<signed 8-bit value>
func NANDR(value int8) Instruction {
	return NewInstruction(NAND, value, false, false)
}

//nand the contens of memory (pointed to by a register)
//nand (%<signed 8-bit value>)
func NANDM(value int8) Instruction {
	return NewInstruction(NAND, value, false, true)
}

//xor an immedate value to the accumulator
//xor <signed 8-bit value>
func XORI(value int8) Instruction {
	return NewInstruction(XOR, value, true, false)
}

//xor the contents of a register to the accumulator
//xor %<signed 8-bit value>
func XORR(value int8) Instruction {
	return NewInstruction(XOR, value, false, false)
}

//xor the contens of memory (pointed to by a register)
//xor (%<signed 8-bit value>)
func XORM(value int8) Instruction {
	return NewInstruction(XOR, value, false, true)
}

//sub an immedate value to the accumulator
//sub <signed 8-bit value>
func SUBI(value int8) Instruction {
	return NewInstruction(SUB, value, true, false)
}

//sub the contents of a register to the accumulator
//sub %<signed 8-bit value>
func SUBR(value int8) Instruction {
	return NewInstruction(SUB, value, false, false)
}

//sub the contens of memory (pointed to by a register)
//sub (%<signed 8-bit value>)
func SUBM(value int8) Instruction {
	return NewInstruction(SUB, value, false, true)
}

//Jump subroutine to PC+offset
//jsr <signed 8-bit value>
func JSRI(value int8) Instruction {
	return NewInstruction(JSR, value, true, false)
}

//Jump subroutine to location in register
//jsr %<signed 8-bit value>
func JSRR(value int8) Instruction {
	return NewInstruction(JSR, value, false, false)
}

//Jump subroutine to a location in memory
//jsr (%<signed 8-bit value>)
func JSRM(value int8) Instruction {
	return NewInstruction(JSR, value, false, true)
}

//Jump to PC+offset
//jmp <signed 8-bit value>
func JMPI(value int8) Instruction {
	return NewInstruction(JMP, value, true, false)
}

//Jump to the address in the given register
//jmp %<signed 8-bit value>
func JMPR(value int8) Instruction {
	return NewInstruction(JMP, value, false, false)
}

//Jmp to the address in the contens of memory (pointed to by a register)
//jmp (%<signed 8-bit value>)
func JMPM(value int8) Instruction {
	return NewInstruction(JMP, value, false, true)
}

//Jump to PC+offset if acc!=0
//jnz <signed 8-bit value>
func JNZI(value int8) Instruction {
	return NewInstruction(JNZ, value, true, false)
}

//Jump to the address in the given register if acc!=0
//jnz %<signed 8-bit value>
func JNZR(value int8) Instruction {
	return NewInstruction(JNZ, value, false, false)
}

//Jmp to the address in the contents of memory
//(pointed to by a register) if acc != 0
//jnz (%<signed 8-bit value>)
func JNZM(value int8) Instruction {
	return NewInstruction(JNZ, value, false, true)
}

//Jump to PC+offset if acc > 0
//jns <signed 8-bit value>
func JNSI(value int8) Instruction {
	return NewInstruction(JNS, value, true, false)
}

//Jump to the address in the given register if acc > 0
//jns %<signed 8-bit value>
func JNSR(value int8) Instruction {
	return NewInstruction(JNS, value, false, false)
}

//Jmp to the address in the contents of memory
//(pointed to by a register) if acc > 0
//jns (%<signed 8-bit value>)
func JNSM(value int8) Instruction {
	return NewInstruction(JNS, value, false, true)
}

func NewHalt() Instruction {
	return NewInstruction(HALT, 0, false, false)
}

func NewNop() Instruction {
	return NewInstruction(NOP, 0, false, false)
}
