package instruction

const (
	NOP uint8 = iota
	LD
	STORE
	ADD
	NAND
	XOR
	SUB
	JSR
	JMP
	JNZ
	JNS
	HALT
)

var OpcodeStringMap = map[uint8]string{
	HALT:  "halt",
	NOP:   "nop",
	LD:    "lda",
	STORE: "sta",
	ADD:   "add",
	NAND:  "nand",
	XOR:   "xor",
	SUB:   "sub",
	JSR:   "jsr",
	JMP:   "jmp",
	JNZ:   "jnz",
	JNS:   "jns",
}

var StringOpcodeMap = map[string]uint8{
	"halt": HALT,
	"nop":  NOP,
	"lda":  LD,
	"sta":  STORE,
	"add":  ADD,
	"nand": NAND,
	"xor":  XOR,
	"sub":  SUB,
	"jsr":  JSR,
	"jmp":  JMP,
	"jnz":  JNZ,
	"jns":  JNS,
}
