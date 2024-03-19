package assemble

import (
	"andromeda/toolchain/assembler/tokenizer"
	"errors"
	"fmt"
)

func (a *AssemblyContext) getLabel(shouldBackpatch bool) uint64 {
	a.checkAndConsumeByID(tokenizer.TOK_IDENT)
	identTok, err := a.curToken()
	a.insertError(err)

	ident := identTok.Contents

	maybeDot, err := a.nextToken()
	if maybeDot.ID == tokenizer.TOK_DOT {
		a.checkAndConsumeByID(tokenizer.TOK_IDENT)
		second, err := a.curToken()
		a.insertError(err)
		ident += "." + second.Contents
	} else {
		_, err := a.lastToken()
		if err != nil {
			a.insertError(err)
		}
	}

	addr, ok := a.getCurrentObject().Labels[ident]
	if !ok && shouldBackpatch {
		lst, lstOk := a.backpatches[ident]
		if lstOk {
			a.backpatches[ident] = append(lst, a.curAddress())
			return 0
		} else {
			a.backpatches[ident] = []uint16{a.curAddress()}
			return uint64(addr)
		}
	} else if !shouldBackpatch && !ok {
		a.insertError(errors.New("label not yet defined"))
		return 0
	}

	return uint64(addr)
}

func (a *AssemblyContext) handleSubdefines() {
	a.checkAndConsumeByID(tokenizer.TOK_IDENT)
	tokIdent, _ := a.curToken()
	obj := a.getCurrentObject()

	ident := a.curOuterLabel + "." + tokIdent.Contents

	_, ok := obj.Labels[ident]
	if ok {
		a.insertError(fmt.Errorf("Redfined label '%s'\n", ident))
	}

	obj.Labels[ident] = a.curAddress()
	a.updatePatches(ident)

	a.checkAndConsume(tokenizer.TOK_COLON, ":")
}

func (a *AssemblyContext) handleDefines() {
	a.checkAndConsumeByID(tokenizer.TOK_IDENT)
	tokIdent, _ := a.curToken()
	obj := a.getCurrentObject()

	_, ok := obj.Labels[tokIdent.Contents]
	if ok {
		a.insertError(fmt.Errorf("Redfined label '%s'\n", tokIdent.Contents))
	}

	obj.Labels[tokIdent.Contents] = a.curAddress()
	a.curOuterLabel = tokIdent.Contents
	a.updatePatches(tokIdent.Contents)

	a.checkAndConsume(tokenizer.TOK_COLON, ":")
}

func (a *AssemblyContext) updatePatches(label string) {
	patches, ok := a.backpatches[label]
	if !ok {
		return
	}

	for _, patch := range patches {
		for _, object := range a.objects {
			object.InsertPatch(patch, a.curAddress())
		}
	}
}
