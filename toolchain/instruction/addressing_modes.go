package instruction

const (
	AM_IMM byte = 0b000
	AM_DIR      = 0b001
	AM_REL      = 0b010
	AM_OFF      = 0b011
	AM_IND      = 0b100
	AM_INC      = 0b101
	AM_DEC      = 0b110
)

var AddressingStringMap = map[byte]string{
	AM_IMM: "imm",
	AM_DIR: "dir",
	AM_IND: "ind",
	AM_INC: "inc",
	AM_DEC: "dec",
	AM_REL: "rel",
	AM_OFF: "off",
}

var StringAddressingMap = map[string]byte{
	"imm": AM_IMM,
	"dir": AM_DIR,
	"rel": AM_REL,
	"ind": AM_IND,
	"inc": AM_INC,
	"dec": AM_DEC,
	"off": AM_OFF,
}
