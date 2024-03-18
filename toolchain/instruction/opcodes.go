package instruction

const (
	JMP uint8 = iota
	NOP
	LD
	STORE
	ADD
	NAND
	XOR
	SUB
	JSR
	JNZ
	JNS
	HALT
)

var OpcodeStringMap = map[uint8]string{
	HALT:  "hlt",
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
	"hlt":  HALT,
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
