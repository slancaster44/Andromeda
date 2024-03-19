package assemble

import (
	"andromeda/toolchain/assembler/object"
	"andromeda/toolchain/assembler/tokenizer"
	"andromeda/toolchain/instruction"
	"fmt"
)

func (a *AssemblyContext) handleDirective() {
	tok, _ := a.curToken()
	switch tok.Contents {
	case "org":
		a.handleOrigin()
	case "dw":
		a.handleDataword()
	case "equ":
		a.handleEqu()
	default:
		a.insertError(fmt.Errorf("Unkown directive '%s'\n", tok.Contents))
	}
}

func (a *AssemblyContext) handleOrigin() {
	a.checkAndConsume(tokenizer.TOK_LPAREN, "(")

	number := a.getNumber(false)

	a.checkAndConsume(tokenizer.TOK_RPAREN, ")")

	a.objects = append(a.objects, object.NewCodeObject(uint16(number)))
}

func (a *AssemblyContext) handleDataword() {
	a.checkAndConsume(tokenizer.TOK_LPAREN, "(")

	number := a.getNumber(false)

	a.checkAndConsume(tokenizer.TOK_RPAREN, ")")

	obj := a.getCurrentObject()
	obj.Code = append(obj.Code, instruction.Instruction(number))
}

func (a *AssemblyContext) handleEqu() {
	a.checkAndConsume(tokenizer.TOK_LPAREN, "(")
	a.checkAndConsumeByID(tokenizer.TOK_IDENT)

	identTok, _ := a.curToken()
	ident := identTok.Contents

	a.checkAndConsume(tokenizer.TOK_COMMA, ",")

	number := a.getNumber(false)

	a.checkAndConsume(tokenizer.TOK_RPAREN, ")")

	obj := a.getCurrentObject()
	obj.Labels[ident] = uint16(number)
	a.updatePatches(ident)
}
