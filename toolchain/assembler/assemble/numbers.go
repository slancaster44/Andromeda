package assemble

import (
	"andromeda/toolchain/assembler/tokenizer"
	"fmt"
	"strconv"
)

func (a *AssemblyContext) getNumber(shouldBackpatch bool) (uint64, error) {

	tok, err := a.nextToken()
	var number uint64
	if err != nil {
		return 0, fmt.Errorf("Expected number, got eof\n")
	} else if tok.ID == tokenizer.TOK_HEX_INT {
		number, err = strconv.ParseUint(tok.Contents, 16, 16)
	} else if tok.ID == tokenizer.TOK_DEC_INT {
		number, err = strconv.ParseUint(tok.Contents, 10, 16)
	} else if tok.ID == tokenizer.TOK_BIN_INT {
		number, err = strconv.ParseUint(tok.Contents, 2, 16)
	} else if tok.ID == tokenizer.TOK_IDENT {
		_, err := a.lastToken()
		if err != nil {
			return 2, err
		}

		return a.getLabel(shouldBackpatch)
	} else {
		return 1, fmt.Errorf("Expected number got '%s'\n", tok.Contents)
	}

	return number, nil
}
