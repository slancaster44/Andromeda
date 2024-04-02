package assembler

import "andromeda/toolchain/assembler/tokenizer"

func (a *Assembler) passTwo() {
	a.curLoc = 0
	a.pc = 0
	a.outerDef = ""

	tok, err := a.CurTok()
	for err == nil {
		switch tok.ID {
		case tokenizer.TOK_DIR:
			a.handleDirectivePassTwo(tok.Contents)
		case tokenizer.TOK_NEWLINE:
			a.ConsumeTok()
		case tokenizer.TOK_DEF:
			a.handleDef()
		case tokenizer.TOK_SUBDEF:
			a.handleSubDef()
		case tokenizer.TOK_INS:
			a.handleInstructionPassTwo()
		default:
			a.AddErrorf("unable to handle token '%v' in pass two\n", tok)
			a.ConsumeTok()
		}

		tok, err = a.CurTok()
	}
}
