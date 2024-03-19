package object

import (
	"andromeda/toolchain/instruction"
	"fmt"
)

type CodeObject struct {
	Code   []instruction.Instruction
	Labels map[string]uint16
	Errors map[uint16]error
	Origin uint16
}

func NewCodeObject(origin uint16) *CodeObject {
	return &CodeObject{
		Code:   []instruction.Instruction{},
		Labels: map[string]uint16{},
		Errors: map[uint16]error{},
		Origin: origin,
	}
}

func (c *CodeObject) String() string {
	curLocation := c.Origin
	out := fmt.Sprintf("\n\torg\t0x%X\n", c.Origin)
	for _, v := range c.Code {
		err, ok := c.Errors[curLocation]
		if ok {
			fmt.Printf("\n\n\t%v\n\n", err)
		}

		for k, v := range c.Labels {
			if v == curLocation {
				out += fmt.Sprintf("\n%s: ; Defined as 0x%04X\n", k, v)
			}
		}

		out += fmt.Sprintf("\t%-24v; 0x%04X (0x%04X)\n", v, curLocation, v.ToInt16())
		curLocation += 1
	}

	return out
}

func (c *CodeObject) InsertPatch(location uint16, program_counter uint16) {
	offset := int(location - c.Origin)
	if offset < len(c.Code)-1 && offset >= 0 {
		i := c.Code[offset]
		if i.AddressingMode() == instruction.AM_IMM && i.IsJmp() {
			v := uint16(int(program_counter) - int(location))
			c.Code[offset] = instruction.Instruction((uint16(i) & 0xFF00) + v)
		} else {
			c.Code[offset] = instruction.Instruction(i + instruction.Instruction(0x00FF&program_counter))
		}
	}
}

func (c *CodeObject) DoesOverlap(co *CodeObject) bool {
	cEnd := int(c.Origin) + len(c.Code) - 1
	coEnd := int(co.Origin) + len(co.Code) - 1
	return (cEnd < int(co.Origin)) && (coEnd < int(c.Origin))
}
