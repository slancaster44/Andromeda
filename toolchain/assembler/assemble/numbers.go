package assemble

import (
	"andromeda/toolchain/assembler/tokenizer"
	"fmt"
	"strconv"
)

func (a *AssemblyContext) getNumber(shouldBackpatch bool) uint64 {

	tok, err := a.nextToken()
	var number uint64
	if err != nil {
		a.insertError(fmt.Errorf("Expected number, got eof\n"))
		return 0
	} else if tok.ID == tokenizer.TOK_HEX_INT {
		number, err = strconv.ParseUint(tok.Contents, 16, 16)
	} else if tok.ID == tokenizer.TOK_DEC_INT {
		number, err = strconv.ParseUint(tok.Contents, 10, 16)
	} else if tok.ID == tokenizer.TOK_BIN_INT {
		number, err = strconv.ParseUint(tok.Contents, 2, 16)
	} else if tok.ID == tokenizer.TOK_IDENT {
		_, err := a.lastToken()
		if err != nil {
			a.insertError(err)
			return 2
		}

		return a.getLabel(shouldBackpatch)
	} else {
		a.insertError(fmt.Errorf("Expected number got '%s'\n", tok.Contents))
		return 1
	}

	return number
}
