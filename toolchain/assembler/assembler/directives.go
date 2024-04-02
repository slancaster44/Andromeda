package assembler

import (
	"andromeda/toolchain/assembler/tokenizer"
	"errors"
)

func (a *Assembler) handleDirectivePassOne(directive string) {
	a.ConsumeTok() //Skip over directive name token
	a.CheckAndConsumeToken(errors.New("expected '('"), tokenizer.TOK_LPAREN)

	switch directive {
	case "org":
		a.handleOrg()
	case "dw":
		a.handleDwPassOne()
	case "str":
		a.handleStrPassOne()
	default:
		a.AddErrorf("unable to handle directive '%s'\n", directive)
	}

	a.CheckAndConsumeToken(errors.New("expected ')'"), tokenizer.TOK_RPAREN)
}

func (a *Assembler) handleDirectivePassTwo(directive string) {
	a.ConsumeTok() //Skip over directive name token
	a.CheckAndConsumeToken(errors.New("expected '('"), tokenizer.TOK_LPAREN)

	switch directive {
	case "org":
		a.handleOrg()
	case "dw":
		a.handleDwPassTwo()
	default:
		a.AddErrorf("unable to handle directive '%s'\n", directive)
	}

	a.CheckAndConsumeToken(errors.New("expected ')'"), tokenizer.TOK_RPAREN)
}

func (a *Assembler) handleOrg() {
	orgLoc := a.getNumber()
	a.pc = orgLoc
}

func (a *Assembler) handleDwPassOne() {
	a.getNumber() //ignore return, we'll do this again in the second pass
	a.pc++
}

func (a *Assembler) handleDwPassTwo() {
	word := a.getNumber()
	a.Code[a.pc] = uint16(word)
	a.pc++
}

func (a *Assembler) handleStrPassOne() {
	a.CheckCurToken(errors.New("expected string literal"), tokenizer.TOK_STR)
	tok, _ := a.CurTok()
	a.pc += uint16(len(tok.Contents))
	a.ConsumeTok()
}
