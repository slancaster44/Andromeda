package link

import (
	"andromeda/toolchain/assembler/assemble"
	"andromeda/toolchain/assembler/tokenizer"
	"andromeda/toolchain/instruction"
	"reflect"
	"testing"
)

func TestLink(t *testing.T) {
	text :=
		`
		org(0x0000)
		lda.imm 1
		org(0x0003)
		lda.imm 2
		`

	expected := make([]instruction.Instruction, 20)
	expected[0] = instruction.NewInstruction(instruction.LD, instruction.AM_IMM, 1)
	expected[3] = instruction.NewInstruction(instruction.LD, instruction.AM_IMM, 2)

	tokens := tokenizer.Tokenize(text)
	objects := assemble.NewAssemblyContext().Assemble(tokens)
	output, err := NewLinkerContext(0x0000, 20).Link(objects)

	if err != nil {
		t.Fatalf("Got assembly error '%v'\n", err)
	}

	if !reflect.DeepEqual(output, expected) {
		t.Fatalf("Expected does not match output\n\n%v\n\n%v\n", expected, output)
	}
}
