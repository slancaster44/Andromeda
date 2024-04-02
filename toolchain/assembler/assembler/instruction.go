package assembler

import (
	"andromeda/toolchain/assembler/tokenizer"
	"andromeda/toolchain/instruction"
	"errors"
)

func (a *Assembler) handleInstructionPassOne() {
	a.pc++

	tok, _ := a.CurTok() //Err handled before we get here
	if tok.Contents == "hlt" || tok.Contents == "nop" {
		a.ConsumeTok()
		return
	}

	a.ConsumeTok()
	a.CheckAndConsumeToken(errors.New("expected '.'"), tokenizer.TOK_DOT)
	a.CheckAndConsumeToken(errors.New("expected addressing mode identifier"), tokenizer.TOK_ADDR_MODE)

	a.ConsumeTok() //Consume whatever the argument is

	tok, err := a.CurTok() //Deal with the cases where there might be more
	if err == nil {
		if tok.ID == tokenizer.TOK_DOT {
			a.ConsumeTok()
			a.CheckAndConsumeToken(errors.New("expected identifier after '.'"), tokenizer.TOK_IDENT)
		} else if tok.ID == tokenizer.TOK_MINUS {
			a.getNumber() //Let the get number function handle checking
		}
	}
}

func (a *Assembler) handleInstructionPassTwo() {
	tok, _ := a.CurTok()
	opcode, ok := instruction.StringOpcodeMap[tok.Contents]
	if !ok {
		a.AddErrorf("unknown operand '%s'", tok.Contents)
		return
	}
	a.ConsumeTok()

	var mode byte
	var number int
	if opcode != instruction.HLT && opcode != instruction.NOP {
		a.CheckAndConsumeToken(errors.New("expected '.'"), tokenizer.TOK_DOT)
		a.CheckCurToken(errors.New("expected addressing mode"), tokenizer.TOK_ADDR_MODE)
		tok, _ = a.CurTok()
		mode, ok = instruction.StringAddressingMap[tok.Contents]
		if !ok {
			a.AddErrorf("unkown addressing mode '%s'", tok.Contents)
			return
		}

		a.ConsumeTok()
		number = int(a.getNumber())

		if mode == instruction.AM_OFF {
			number -= int(a.pc)
		}
	}

	a.Code[a.pc] = uint16(instruction.NewInstruction(opcode, mode, number))
	a.pc++
}
