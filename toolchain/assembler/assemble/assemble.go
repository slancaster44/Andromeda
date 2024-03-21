package assemble

import (
	"andromeda/toolchain/assembler/object"
	"andromeda/toolchain/assembler/tokenizer"
	"errors"
	"fmt"
)

type AssemblyContext struct {
	objects []*object.CodeObject

	tokens      []tokenizer.Token
	curLocation int

	curOuterLabel string
	backpatches   map[string][]uint16
}

func NewAssemblyContext() *AssemblyContext {
	return &AssemblyContext{
		objects:       []*object.CodeObject{},
		backpatches:   map[string][]uint16{},
		curLocation:   0,
		curOuterLabel: "",
	}
}

func (a *AssemblyContext) curToken() (tokenizer.Token, error) {
	if a.curLocation >= len(a.tokens) {
		return tokenizer.Token{}, errors.New("Out of tokens\n")
	}

	return a.tokens[a.curLocation], nil
}

func (a *AssemblyContext) nextToken() (tokenizer.Token, error) {
	a.curLocation++
	return a.curToken()
}

func (a *AssemblyContext) lastToken() (tokenizer.Token, error) {
	a.curLocation--
	return a.curToken()
}

func (a *AssemblyContext) getCurrentObject() *object.CodeObject {
	if len(a.objects) == 0 {
		panic(fmt.Errorf("No origin previously defined\n"))
	}

	return a.objects[len(a.objects)-1]
}

func (a *AssemblyContext) curAddress() uint16 {
	obj := a.getCurrentObject()
	return obj.Origin + uint16(len(obj.Code))
}

func (a *AssemblyContext) insertError(e error) {
	if e != nil {
		obj := a.getCurrentObject()
		obj.Errors[a.curAddress()] = e
	}
}

func (a *AssemblyContext) Assemble(toks []tokenizer.Token) []*object.CodeObject {
	if len(a.tokens) == 0 {
		a.tokens = toks
	} else {
		a.curLocation++
		front := a.tokens[:a.curLocation]
		back := a.tokens[a.curLocation:]
		a.tokens = append(append(front, toks...), back...)
	}

	tok, err := a.curToken()
	for err == nil {
		switch tok.ID {
		case tokenizer.TOK_DIR:
			a.handleDirective()
		case tokenizer.TOK_INS:
			a.handleInstruction()
		case tokenizer.TOK_DEF:
			a.handleDefines()
		case tokenizer.TOK_SUBDEF:
			a.handleSubdefines()
		case tokenizer.TOK_NEWLINE:
		default:
			a.insertError(fmt.Errorf("Unable to assemble token '%s'\n", tok.Contents))
		}

		tok, err = a.nextToken()
	}

	return a.objects
}

func (a *AssemblyContext) checkAndConsumeByID(id tokenizer.TokenID) {
	tok, err := a.nextToken()
	if err != nil || tok.ID != id {
		a.insertError(fmt.Errorf("Unexpected token '%s'\n", tok.Contents))
	}
}

func (a *AssemblyContext) checkAndConsume(id tokenizer.TokenID, contents string) {
	tok, err := a.nextToken()
	if err != nil {
		a.insertError(err)
	} else if tok.ID != id || tok.Contents != contents {
		a.insertError(fmt.Errorf("Expected '%s' got '%s'\n", contents, tok.Contents))
	}
}
