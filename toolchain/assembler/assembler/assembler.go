package assembler

import (
	"andromeda/toolchain/assembler/tokenizer"
	"errors"
	"fmt"
)

type Assembler struct {
	Labels map[string]uint16
	Errors []error

	outerDef string

	pc         uint16
	tokens     []tokenizer.Token
	curLoc     int
	passNumber int

	Code [1024 * 64]uint16
}

func Assemble(tokens []tokenizer.Token) *Assembler {
	p := &Assembler{
		Labels: map[string]uint16{},
		Errors: []error{},

		pc:       0,
		tokens:   tokens,
		curLoc:   0,
		outerDef: "",
	}

	p.passOne()
	p.passTwo()

	return p
}

//////* Some Basic Utility Functions *//////

func (a *Assembler) PrintErrors() {
	for _, err := range a.Errors {
		fmt.Printf("%s", err.Error())
	}
}

func (a *Assembler) AddError(err error) {
	tok, tokErr := a.CurTok()
	if tokErr == nil && err != nil {
		err = fmt.Errorf("(%s:%d) -- %v", tok.Filename, tok.LineNumber, err)
	}

	if err != nil {
		a.Errors = append(a.Errors, err)
	}
}

func (a *Assembler) AddErrorf(format string, args ...any) {
	a.AddError(fmt.Errorf(format, args...))
}

func (a *Assembler) CurTok() (tokenizer.Token, error) {
	if a.curLoc >= len(a.tokens) {
		return tokenizer.Token{}, errors.New("out of tokens")
	}

	return a.tokens[a.curLoc], nil
}

func (a *Assembler) ConsumeTok() {
	a.curLoc++
}

func (a *Assembler) CheckCurToken(err error, ids ...tokenizer.TokenID) {
	isGood := false
	tok, tokErr := a.CurTok()
	if tokErr != nil {
		a.AddError(tokErr)
	}

	for _, id := range ids {
		if id == tok.ID {
			isGood = true
		}
	}

	if !isGood {
		a.AddError(fmt.Errorf("error: %v, at token '%s'\n", err, tok.Contents))
	}
}

func (a *Assembler) CheckAndConsumeToken(err error, ids ...tokenizer.TokenID) {
	a.CheckCurToken(err, ids...)
	a.ConsumeTok()
}

func (a *Assembler) AddLabel(identifier string, value uint16) {
	_, ok := a.Labels[identifier]
	if ok && a.passNumber == 1 {
		a.AddErrorf("redefined label '%s'", identifier)
	}

	a.Labels[identifier] = value
}
