package assembler

import (
	"andromeda/toolchain/assembler/tokenizer"
)

func (a *Assembler) passOne() {
	tok, err := a.CurTok()
	for err == nil {
		switch tok.ID {
		case tokenizer.TOK_DIR:
			a.handleDirectivePassOne(tok.Contents)
		case tokenizer.TOK_NEWLINE:
			a.ConsumeTok()
		case tokenizer.TOK_DEF:
			a.handleDef()
		case tokenizer.TOK_SUBDEF:
			a.handleSubDef()
		case tokenizer.TOK_INS:
			a.handleInstructionPassOne()
		default:
			a.AddErrorf("unable to handle token '%v' in pass one\n", tok)
			a.ConsumeTok()
		}

		tok, err = a.CurTok()
	}
}
