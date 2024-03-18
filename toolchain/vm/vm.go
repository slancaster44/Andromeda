package vm

import (
	"andromeda/toolchain/instruction"
	"fmt"
)

type VM struct {
	Memory      [64 * 1024]int16
	Accumulator int16
	PC          uint16
	Operations  map[uint8]func(int16)
	HFF         bool
}

func NewVM(memory_image []int16) *VM {
	vm := &VM{}

	copy(vm.Memory[:], memory_image)
	vm.Accumulator = 0
	vm.PC = 0
	vm.HFF = false

	vm.Operations = map[uint8]func(int16){
		instruction.NOP:  vm.NOP,
		instruction.LD:   vm.LD,
		instruction.ADD:  vm.ADD,
		instruction.NAND: vm.NAND,
		instruction.XOR:  vm.XOR,
		instruction.SUB:  vm.SUB,
		instruction.JSR:  vm.JSR,
		instruction.JMP:  vm.JMP,
		instruction.JNS:  vm.JNS,
		instruction.JNZ:  vm.JNZ,
		instruction.HALT: vm.HALT,
	}

	return vm
}

func (v *VM) NOP(i int16) {}

func (v *VM) LD(i int16) {
	v.Accumulator = i
}

func (v *VM) ADD(i int16) {
	v.Accumulator += i
}

func (v *VM) NAND(i int16) {
	v.Accumulator &= i
	v.Accumulator = ^v.Accumulator
}

func (v *VM) XOR(i int16) {
	v.Accumulator ^= i
}

func (v *VM) SUB(i int16) {
	v.Accumulator -= i
}

func (v *VM) JSR(i int16) {
	v.Accumulator = int16(v.PC)
	v.JMP(i)
}

func (v *VM) JMP(i int16) {
	ins := instruction.Instruction(v.Memory[v.PC-1])
	if ins.AddressingMode() == instruction.AM_IMM {
		v.PC += uint16(i - 1)
	} else {
		v.PC = uint16(i - 1)
	}
}

func (v *VM) JNZ(i int16) {
	if v.Accumulator != 0 {
		v.JMP(i)
	}
}

func (v *VM) JNS(i int16) {
	if v.Accumulator >= 0 {
		v.JMP(i)
	}
}

func (v *VM) HALT(i int16) {
	v.HFF = true
}

func (v *VM) Store() {
	i := instruction.Instruction(v.Memory[v.PC])

	if i.AddressingMode() == instruction.AM_IMM {
		v.Memory[v.Accumulator] = i.Immediate()
	} else if i.AddressingMode() == instruction.AM_DIR {
		v.Memory[i.Address()] = v.Accumulator
	} else if i.AddressingMode() == instruction.AM_IND {
		v.Memory[uint16(v.Memory[i.Address()])] = v.Accumulator
	} else {
		v.Memory[uint16(v.Memory[i.Address()])] = v.Accumulator
		v.Memory[i.Address()]++
	}
}

func (v *VM) InvalidInstructionTrap() {
	fmt.Printf("Invalid Instruction at '0x%X'\n", v.PC-1)
	v.Accumulator = int16(v.PC + 1)
	v.PC = 2
}

func (v *VM) SingleStep() {
	i := instruction.Instruction(v.Memory[v.PC])

	if !i.IsValid() {
		v.InvalidInstructionTrap()
	} else if i.Opcode() == instruction.STORE {
		v.Store()
	} else {
		op := v.Operations[i.Opcode()]
		var val int16

		if i.AddressingMode() == instruction.AM_IMM {
			val = i.Immediate()
		} else if i.AddressingMode() == instruction.AM_DIR {
			val = v.Memory[i.Address()]
		} else if i.AddressingMode() == instruction.AM_IND {
			val = v.Memory[uint16(v.Memory[i.Address()])]
		} else if i.AddressingMode() == instruction.AM_INC {
			val = v.Memory[uint16(v.Memory[i.Address()])]
			v.Memory[i.Address()]++
		} else {
			v.InvalidInstructionTrap()
		}

		op(val)
	}
	v.PC++

}

func (v *VM) Run() {
	for !v.HFF {
		v.SingleStep()
	}
}
