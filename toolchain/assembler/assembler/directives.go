package assembler

import (
	"andromeda/toolchain/assembler/tokenizer"
	"errors"
	"os"
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
	case "equ":
		a.handleEqu()
	case "include":
		a.handleIncludePassOne()
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
	case "str":
		a.handleStrPassTwo()
	case "equ":
		a.handleEqu()
	case "include":
		a.handleIncludePassTwo()
	default:
		a.AddErrorf("unable to handle directive '%s'\n", directive)
	}

	a.CheckAndConsumeToken(errors.New("expected ')'"), tokenizer.TOK_RPAREN)
}

func (a *Assembler) handleOrg() {
	orgLoc := a.getNumber()
	a.pc = orgLoc
}

func (a *Assembler) handleEqu() {
	a.CheckCurToken(errors.New("expected identifier"), tokenizer.TOK_IDENT)
	tok, _ := a.CurTok()
	ident := tok.Contents
	a.ConsumeTok()

	a.CheckAndConsumeToken(errors.New("expected ','"), tokenizer.TOK_COMMA)
	number := a.getNumber()

	a.AddLabel(ident, number)
}

func (a *Assembler) handleDwPassOne() {
	a.getNumberPassOne() //ignore return, we'll do this again in the second pass
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

func (a *Assembler) handleIncludePassOne() {
	a.CheckCurToken(errors.New("expected string literal"), tokenizer.TOK_STR)
	tok, _ := a.CurTok()
	a.ConsumeTok()

	bytes, err := os.ReadFile(tok.Contents)
	a.AddError(err)
	fileContents := string(bytes)

	newTokens := tokenizer.Tokenize(fileContents, tok.Contents)
	leftTokens := a.tokens[:a.curLoc+1]
	rightTokens := a.tokens[a.curLoc+1:]

	tokens := append(newTokens, rightTokens...)
	a.tokens = append(leftTokens, tokens...)
}

func (a *Assembler) handleIncludePassTwo() {
	a.CheckAndConsumeToken(errors.New("expected string literal"), tokenizer.TOK_STR)
}

func (a *Assembler) handleStrPassTwo() {
	a.CheckCurToken(errors.New("expected string literal"), tokenizer.TOK_STR)
	tok, _ := a.CurTok()
	a.ConsumeTok()

	for i := a.pc; i < a.pc+uint16(len(tok.Contents)); i++ {
		a.Code[i] = uint16(tok.Contents[i-a.pc])
	}
	a.pc += uint16(len(tok.Contents))
}
