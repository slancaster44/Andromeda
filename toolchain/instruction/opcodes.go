package instruction

const (
	NOP uint8 = 0b00000
	LDA       = 0b00001
	STA       = 0b00010
	ADD       = 0b00011
	NND       = 0b00100
	XOR       = 0b00101
	SUB       = 0b00110
	JMP       = 0b00111
	JSR       = 0b01000
	JNZ       = 0b01001
	JNS       = 0b01010
	HLT       = 0b01011
	INP       = 0b01100
	OUT       = 0b01101
)

var OpcodeStringMap = map[uint8]string{
	HLT: "hlt",
	NOP: "nop",
	LDA: "lda",
	STA: "sta",
	ADD: "add",
	NND: "nnd",
	XOR: "xor",
	SUB: "sub",
	JSR: "jsr",
	JMP: "jmp",
	JNZ: "jnz",
	JNS: "jns",
	INP: "inp",
	OUT: "out",
}

var StringOpcodeMap = map[string]uint8{
	"hlt": HLT,
	"nop": NOP,
	"lda": LDA,
	"sta": STA,
	"add": ADD,
	"nnd": NND,
	"xor": XOR,
	"sub": SUB,
	"jsr": JSR,
	"jmp": JMP,
	"jnz": JNZ,
	"jns": JNS,
	"out": OUT,
	"inp": INP,
}
