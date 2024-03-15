package assembler

import (
	"testing"
	"toolchain/instruction"
)

func TestOrigin(t *testing.T) {
	text := "org(0x0012)"
	context := NewAssemblyContext()
	objects := context.Assemble(Tokenize(text))
	_, eofErr := context.curToken()

	if len(objects) != 1 {
		t.Errorf("Expected only one object")
	} else if objects[0].Origin != 0x0012 {
		t.Errorf("Expected origin of '0x0012' got '0x%X'\n", objects[0].Origin)
	} else if eofErr == nil {
		t.Errorf("Expected all tokens to be consumed\n")
	}

	for _, v := range objects[0].Errors {
		t.Error(v)
	}
}

func TestInstruction(t *testing.T) {
	text :=
		`
		org (0x1000)
		add 1
		add %2
		add (%3)
		`

	expectedInstruction := []instruction.Instruction{
		instruction.ADDI(1),
		instruction.ADDR(2),
		instruction.ADDM(3),
	}

	context := NewAssemblyContext()
	objects := context.Assemble(Tokenize(text))
	_, eofErr := context.curToken()

	if len(objects) != 1 {
		t.Errorf("Expected only one object")
	} else if objects[0].Origin != 0x1000 {
		t.Errorf("Expected origin of '0x1000' got '0x%X'\n", objects[0].Origin)
	} else if eofErr == nil {
		t.Errorf("Expected all tokens to be consumed\n")
	}

	for _, v := range objects[0].Errors {
		t.Error(v)
	}

	for i, v := range expectedInstruction {
		if objects[0].Code[i] != v {
			t.Fatalf("Did not generate expected instructions\n")
		}
	}

	t.Log(objects[0])
}

func TestDefines(t *testing.T) {
	text :=
		`
		org(0x0000)
	def entry:
		lda	1
		add 2
	def exit:
		sta %0x70
		`
	expectedInstruction := []instruction.Instruction{
		instruction.LDI(1),
		instruction.ADDI(2),
		instruction.STR(0x70),
	}

	expectedLabel := map[string]uint16{
		"entry": 0x0000,
		"exit":  0x0002,
	}

	context := NewAssemblyContext()
	objects := context.Assemble(Tokenize(text))
	_, eofErr := context.curToken()

	if len(objects) != 1 {
		t.Errorf("Expected only one object")
	} else if objects[0].Origin != 0x0000 {
		t.Errorf("Expected origin of '0x1000' got '0x%X'\n", objects[0].Origin)
	} else if eofErr == nil {
		t.Errorf("Expected all tokens to be consumed\n")
	}

	t.Logf("Number of errors %d\n", len(objects[0].Errors))
	for _, v := range objects[0].Errors {
		t.Error(v)
	}

	for i, v := range expectedInstruction {
		if objects[0].Code[i] != v {
			t.Fatalf("Did not generate expected instructions\n")
		}
	}

	for i, v := range expectedLabel {
		t.Logf("%s - 0x%X", i, v)
		if objects[0].Labels[i] != v {
			t.Fatalf("Did not generate expected labels got '0x%X'\n", objects[0].Labels[i])
		}
	}

	t.Log(objects[0])
}
