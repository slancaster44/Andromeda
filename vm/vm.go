package vm

import (
	"fmt"
	"toolchain/instruction"
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
		instruction.NOP:   vm.NOP,
		instruction.LD:    vm.LD,
		instruction.STORE: vm.STORE,
		instruction.ADD:   vm.ADD,
		instruction.NAND:  vm.NAND,
		instruction.XOR:   vm.XOR,
		instruction.SUB:   vm.SUB,
		instruction.JSR:   vm.JSR,
		instruction.JMP:   vm.JMP,
		instruction.JNS:   vm.JNS,
		instruction.JNZ:   vm.JNZ,
		instruction.HALT:  vm.HALT,
	}

	return vm
}

func (v *VM) NOP(i int16) {}

func (v *VM) LD(i int16) {
	v.Accumulator = i
}

func (v *VM) STORE(i int16) {
	v.Memory[uint16(i)] = v.Accumulator
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
	if ins.IsImmediate() {
		v.PC += uint16(i)
	} else {
		v.PC = uint16(i)
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

func (v *VM) InvalidInstructionTrap() {
	fmt.Printf("Invalid Instruction at '0x%X'\n", v.PC-1)
	v.Accumulator = int16(v.PC)
	v.PC = 2
}

func (v *VM) SingleStep() {
	i := instruction.Instruction(v.Memory[v.PC])
	v.PC++

	if !i.IsValid() {
		v.InvalidInstructionTrap()
	} else {
		op := v.Operations[i.Opcode()]
		var val int16

		if i.IsImmediate() {
			val = i.Immediate()
		} else if i.IsDirect() {
			val = v.Memory[i.Immediate()]
		} else if i.IsIndirect() {
			val = v.Memory[v.Memory[i.Immediate()]]
		} else {
			v.InvalidInstructionTrap()
		}

		op(val)
	}

}

func (v *VM) Run() {
	for !v.HFF {
		v.SingleStep()
	}
}
