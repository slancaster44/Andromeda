package assemble

import (
	"andromeda/toolchain/assembler/tokenizer"
	"andromeda/toolchain/instruction"
	"fmt"
)

func (a *AssemblyContext) handleInstruction() {
	mnemonicTok, _ := a.curToken()
	mnemonic := mnemonicTok.Contents
	opcode, ok := instruction.StringOpcodeMap[mnemonic]

	if !ok {
		a.insertError(fmt.Errorf("Unexpected opcode '%s'\n", mnemonic))
	}

	var ins instruction.Instruction
	if opcode == instruction.HALT || opcode == instruction.NOP {
		ins = instruction.NewInstruction(opcode, 0, 0)
	} else {
		a.checkAndConsume(tokenizer.TOK_DOT, ".")
		a.checkAndConsumeByID(tokenizer.TOK_ADDR_MODE)

		amTok, _ := a.curToken()
		addrMode, ok := instruction.StringAddressingMap[amTok.Contents]
		if !ok {
			a.insertError(fmt.Errorf("Unexpected addressing mode '%s'\n", amTok.Contents))
		}

		number, err := a.getNumber(true)
		a.insertError(err)

		ins = instruction.NewInstruction(opcode, addrMode, int(number))

		if ins.IsJmp() && addrMode == instruction.AM_IMM {
			number = number - uint64(a.curAddress())
			ins = instruction.NewInstruction(opcode, addrMode, int(number))
		}
	}

	obj := a.getCurrentObject()
	obj.Code = append(obj.Code, ins)
}
