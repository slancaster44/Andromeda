package assembler

import (
	"andromeda/toolchain/assembler/tokenizer"
	"errors"
	"strconv"
)

func (a *Assembler) getNumber() uint16 {
	tok, err := a.CurTok()
	a.AddError(err)

	var out uint64
	switch tok.ID {
	case tokenizer.TOK_HEX_INT:
		out, err = strconv.ParseUint(tok.Contents, 16, 16)
	case tokenizer.TOK_DEC_INT:
		out, err = strconv.ParseUint(tok.Contents, 10, 16)
	case tokenizer.TOK_BIN_INT:
		out, err = strconv.ParseUint(tok.Contents, 2, 16)
	case tokenizer.TOK_MINUS:
		a.ConsumeTok()
		return uint16(-int16(a.getNumber()))
	case tokenizer.TOK_DOLLAR:
		out = uint64(a.pc)
	case tokenizer.TOK_CARROT:
		val, _ := a.Labels[a.lastOuterDef]
		out = uint64(val)
	case tokenizer.TOK_IDENT:
		return a.handleLabel()
	case tokenizer.TOK_CHAR:
		out = uint64(tok.Contents[0])
	default:
		a.AddErrorf("expected constant integer value,  got '%s'\n", tok.Contents)
	}
	a.AddError(err)

	a.ConsumeTok()
	return uint16(out)
}

func (a *Assembler) handleLabel() uint16 {
	tok, err := a.CurTok()
	a.AddError(err)

	ident := tok.Contents

	a.ConsumeTok()

	tok2, err := a.CurTok()
	if err == nil && tok2.ID == tokenizer.TOK_DOT {
		a.ConsumeTok()
		a.CheckCurToken(errors.New("expected secondary identifier after '.'"), tokenizer.TOK_IDENT)
		tok3, _ := a.CurTok() //Err checked by CheckCurToken
		ident += "." + tok3.Contents
		a.ConsumeTok()
	}

	value, ok := a.Labels[ident]
	if !ok {
		a.AddErrorf("undefined label '%s'", ident)
	}

	return value
}

func (a *Assembler) getNumberPassOne() uint16 {
	tok, err := a.CurTok()
	a.AddError(err)

	var out uint64
	switch tok.ID {
	case tokenizer.TOK_HEX_INT:
		out, err = strconv.ParseUint(tok.Contents, 16, 16)
	case tokenizer.TOK_DEC_INT:
		out, err = strconv.ParseUint(tok.Contents, 10, 16)
	case tokenizer.TOK_BIN_INT:
		out, err = strconv.ParseUint(tok.Contents, 2, 16)
	case tokenizer.TOK_MINUS:
		a.ConsumeTok()
		return uint16(-int16(a.getNumber()))
	case tokenizer.TOK_DOLLAR:
		out = uint64(a.pc)
	case tokenizer.TOK_IDENT:
		return a.handleLabelPassOne()
	case tokenizer.TOK_CHAR:
		out = uint64(tok.Contents[0])
	case tokenizer.TOK_CARROT:
		val, _ := a.Labels[a.lastOuterDef]
		out = uint64(val)
	default:
		a.AddErrorf("Expected constant integer value,  got '%s'\n", tok.Contents)
	}
	a.AddError(err)

	a.ConsumeTok()
	return uint16(out)
}

func (a *Assembler) handleLabelPassOne() uint16 {
	tok, err := a.CurTok()
	a.AddError(err)

	ident := tok.Contents

	a.ConsumeTok()

	tok2, err := a.CurTok()
	if err == nil && tok2.ID == tokenizer.TOK_DOT {
		a.ConsumeTok()
		a.CheckCurToken(errors.New("expected secondary identifier after '.'"), tokenizer.TOK_IDENT)
		tok3, _ := a.CurTok() //Err checked by CheckCurToken
		ident += "." + tok3.Contents
		a.ConsumeTok()
	}

	return 0
}
