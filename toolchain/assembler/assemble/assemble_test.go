package assemble

import (
	"andromeda/toolchain/assembler/object"
	"andromeda/toolchain/assembler/tokenizer"
	"andromeda/toolchain/instruction"
	"reflect"
	"testing"
)

func TestOrigin(t *testing.T) {
	text := "org(0x0012)"
	context := NewAssemblyContext()
	objects := context.Assemble(tokenizer.Tokenize(text))
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
		add.imm 1
		add.dir 2
		add.ind 3
		add.inc 4
		`

	expectedInstruction := []instruction.Instruction{
		instruction.NewInstruction(instruction.ADD, instruction.AM_IMM, 1),
		instruction.NewInstruction(instruction.ADD, instruction.AM_DIR, 2),
		instruction.NewInstruction(instruction.ADD, instruction.AM_IND, 3),
		instruction.NewInstruction(instruction.ADD, instruction.AM_INC, 4),
	}

	context := NewAssemblyContext()
	objects := context.Assemble(tokenizer.Tokenize(text))
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
		lda.imm	1
		add.imm 2
	subdef exit:
		sta.dir 0x70
		`
	expectedInstruction := []instruction.Instruction{
		instruction.NewInstruction(instruction.LD, instruction.AM_IMM, 1),
		instruction.NewInstruction(instruction.ADD, instruction.AM_IMM, 2),
		instruction.NewInstruction(instruction.STORE, instruction.AM_DIR, 0x70),
	}

	expectedLabel := map[string]uint16{
		"entry":      0x0000,
		"entry.exit": 0x0002,
	}

	context := NewAssemblyContext()
	objects := context.Assemble(tokenizer.Tokenize(text))
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

func TestLabels(t *testing.T) {
	text :=
		`org (0x0000)
	def loop_back:
		lda.dir constant
		jmp.rel loop_back
		jmp.rel forward
	def constant: dw(11)
	def forward: hlt`

	expectedInstruction := []instruction.Instruction{
		instruction.NewInstruction(instruction.LD, instruction.AM_DIR, 3),
		instruction.NewInstruction(instruction.JMP, instruction.AM_REL, -1),
		instruction.NewInstruction(instruction.JMP, instruction.AM_REL, 2),
		instruction.Instruction(11),
		instruction.NewHalt(),
	}

	expectedLabel := map[string]uint16{
		"loop_back": 0x0000,
		"constant":  0x0003,
		"forward":   0x0004,
	}

	context := NewAssemblyContext()
	objects := context.Assemble(tokenizer.Tokenize(text))
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
			t.Errorf("Did not generate expected instructions\n")
		}
	}

	for i, v := range expectedLabel {
		t.Logf("%s - 0x%X", i, v)
		if objects[0].Labels[i] != v {
			t.Errorf("Did not generate expected labels got '0x%X'\n", objects[0].Labels[i])
		}
	}

	t.Log(objects[0])
	t.Log(objects[0])
}

func TestMultipleObjects(t *testing.T) {
	text :=
		`
		org(0x0000)
		lda.imm 1
		org(0x0003)
		lda.imm 2
		`

	expected1 := &object.CodeObject{
		Code: []instruction.Instruction{
			instruction.NewInstruction(instruction.LD, instruction.AM_IMM, 1),
		},
		Origin: 0,
		Labels: map[string]uint16{},
		Errors: map[uint16]error{},
	}

	expected2 := &object.CodeObject{
		Code: []instruction.Instruction{
			instruction.NewInstruction(instruction.LD, instruction.AM_IMM, 2),
		},
		Origin: 3,
		Labels: map[string]uint16{},
		Errors: map[uint16]error{},
	}

	tokens := tokenizer.Tokenize(text)
	objects := NewAssemblyContext().Assemble(tokens)

	if len(objects) != 2 {
		t.Fatalf("Expected 2 objects, got '%d'\n", len(objects))
	}

	if !reflect.DeepEqual(expected1, objects[0]) {
		t.Log("Expected does not match first object")
		t.Fatalf("%v\n%v", expected1, objects[0])
	}

	if !reflect.DeepEqual(expected2, objects[1]) {
		t.Log("Expected does not match first object")
		t.Fatalf("%v\n%v", expected2, objects[1])
	}
}

func TestDollar(t *testing.T) {
	text :=
		`
		org(0x0012)
		lda.imm $+
		`

	expected1 := &object.CodeObject{
		Code: []instruction.Instruction{
			instruction.NewInstruction(instruction.LD, instruction.AM_IMM, 0x13),
		},
		Origin: 0x12,
		Labels: map[string]uint16{},
		Errors: map[uint16]error{},
	}

	tokens := tokenizer.Tokenize(text)
	objects := NewAssemblyContext().Assemble(tokens)

	if !reflect.DeepEqual(expected1, objects[0]) {
		t.Log("Expected does not match first object")
		t.Fatalf("%v\n%v", expected1, objects[0])
	}
}
