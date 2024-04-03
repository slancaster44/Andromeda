package assembler

import (
	"andromeda/toolchain/assembler/tokenizer"
	"errors"
)

func (a *Assembler) handleDef() {
	a.ConsumeTok() //Move over 'def' token
	a.CheckCurToken(errors.New("expected identifier"), tokenizer.TOK_IDENT)

	tok, err := a.CurTok()
	a.AddError(err)

	ident := tok.Contents
	a.AddLabel(ident, a.pc)
	a.outerDef = ident

	a.ConsumeTok()
	a.CheckAndConsumeToken(errors.New("expected ':'"), tokenizer.TOK_COLON)
}

func (a *Assembler) handleSubDef() {
	a.ConsumeTok() //Move over 'def' token
	a.CheckCurToken(errors.New("expected identifier"), tokenizer.TOK_IDENT)

	tok, err := a.CurTok()
	a.AddError(err)

	ident := a.outerDef + "." + tok.Contents
	a.AddLabel(ident, a.pc)

	a.ConsumeTok()
	a.CheckAndConsumeToken(errors.New("expected ':'"), tokenizer.TOK_COLON)
}
