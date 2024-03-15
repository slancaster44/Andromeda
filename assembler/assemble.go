package assembler

import (
	"errors"
	"fmt"
	"strconv"
	"toolchain/instruction"
)

type AssemblyContext struct {
	objects []*CodeObject

	tokens      []Token
	curLocation int

	curOuterLabel string
	backpatches   map[string][]uint16
}

func NewAssemblyContext() *AssemblyContext {
	return &AssemblyContext{
		objects:       []*CodeObject{},
		backpatches:   map[string][]uint16{},
		curLocation:   0,
		curOuterLabel: "",
	}
}

func (a *AssemblyContext) curToken() (Token, error) {
	if a.curLocation >= len(a.tokens) {
		return Token{}, errors.New("Out of tokens\n")
	}

	return a.tokens[a.curLocation], nil
}

func (a *AssemblyContext) nextToken() (Token, error) {
	a.curLocation++
	return a.curToken()
}

func (a *AssemblyContext) lastToken() (Token, error) {
	a.curLocation--
	return a.curToken()
}

func (a *AssemblyContext) getCurrentObject() *CodeObject {
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

func (a *AssemblyContext) Assemble(toks []Token) []*CodeObject {
	a.tokens = append(a.tokens, toks...)

	tok, err := a.curToken()
	for err == nil {
		switch tok.ID {
		case TOK_DIR:
			a.handleDirective()
		case TOK_INS:
			a.handleInstruction()
		case TOK_DEF:
			a.handleDefines()
		case TOK_NEWLINE:
		default:
			a.insertError(fmt.Errorf("Unable to assemble token '%s'\n", tok.Contents))
		}

		tok, err = a.nextToken()
	}

	return a.objects
}

func (a *AssemblyContext) handleDirective() {
	tok, _ := a.curToken()
	switch tok.Contents {
	case "org":
		a.handleOrigin()
	default:
		a.insertError(fmt.Errorf("Unkown directive '%s'\n", tok.Contents))
	}
}

func (a *AssemblyContext) handleOrigin() {
	a.checkAndConsume(TOK_LPAREN, "(")

	number, err := a.getNumber()
	a.insertError(err)

	a.checkAndConsume(TOK_RPAREN, ")")

	a.objects = append(a.objects, NewCodeObject(uint16(number)))
}

func (a *AssemblyContext) handleInstruction() {
	tok, _ := a.curToken()
	opcode := instruction.StringOpcodeMap[tok.Contents]

	tok, err := a.nextToken()
	a.insertError(err)

	is_ind := tok.ID == TOK_LPAREN
	if is_ind {
		tok, err = a.nextToken()
		a.insertError(err)
	}

	is_imm := tok.ID != TOK_PERCENT && !is_ind
	if !is_imm && tok.ID != TOK_PERCENT {
		a.insertError(fmt.Errorf("Expected '%%' got '%s'\n", tok.Contents))
	} else if is_imm {
		_, err = a.lastToken()
		a.insertError(err)
	}

	number, err := a.getNumber()
	a.insertError(err)

	if is_ind {
		a.checkAndConsume(TOK_RPAREN, ")")
	}

	ins := instruction.NewInstruction(opcode, int8(number), is_imm, is_ind)
	a.getCurrentObject().Code = append(a.getCurrentObject().Code, ins)
}

func (a *AssemblyContext) checkAndConsumeByID(id TokenID) {
	tok, err := a.nextToken()
	if err != nil || tok.ID != id {
		a.insertError(fmt.Errorf("Unexpected token '%s'\n", tok.Contents))
	}
}

func (a *AssemblyContext) handleDefines() {
	a.checkAndConsumeByID(TOK_IDENT)
	tok_ident, _ := a.curToken()
	obj := a.getCurrentObject()

	obj.Labels[tok_ident.Contents] = a.curAddress()
	a.curOuterLabel = tok_ident.Contents

	a.checkAndConsume(TOK_COLON, ":")
}

func (a *AssemblyContext) checkAndConsume(id TokenID, contents string) {
	tok, err := a.nextToken()
	if err != nil {
		a.insertError(err)
	} else if tok.ID != id || tok.Contents != contents {
		a.insertError(fmt.Errorf("Expected '%s' got '%s'\n", contents, tok.Contents))
	}
}

func (a *AssemblyContext) getNumber() (uint64, error) {

	tok, err := a.nextToken()
	var number uint64
	if err != nil {
		return 0, fmt.Errorf("Expected number, got eof\n")
	} else if tok.ID == TOK_HEX_INT {
		number, err = strconv.ParseUint(tok.Contents, 16, 16)
	} else if tok.ID == TOK_DEC_INT {
		number, err = strconv.ParseUint(tok.Contents, 10, 16)
	} else if tok.ID == TOK_BIN_INT {
		number, err = strconv.ParseUint(tok.Contents, 2, 16)
	} else {
		return 1, fmt.Errorf("Expected number got '%s'\n", tok.Contents)
	}

	return number, nil
}
