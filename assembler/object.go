package assembler

import (
	"fmt"
	"toolchain/instruction"
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
				out += fmt.Sprintf("%s:%-14s; Defined as 0x%04X\n", k, "", v)
			}
		}

		out += fmt.Sprintf("\t%-16v; 0x%04X\n", v, curLocation)
		curLocation += 1
	}

	return out
}
