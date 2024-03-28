package vm

import (
	"andromeda/toolchain/instruction"
	"fmt"
	"os"
	"os/exec"
)

type VM struct {
	Memory        [64 * 1024]int16
	Accumulator   int16
	PC            uint16
	Operations    map[uint8]func(int16)
	HFF           bool
	OutputDevices map[uint16]func(int16)
	InputDevices  map[uint16]func() int16
}

func NewVM(memory_image []int16) *VM {
	vm := &VM{}

	copy(vm.Memory[:], memory_image)
	vm.Accumulator = 0
	vm.PC = 0
	vm.HFF = false

	vm.Operations = map[uint8]func(int16){
		instruction.NOP: vm.NOP,
		instruction.LDA: vm.LD,
		instruction.ADD: vm.ADD,
		instruction.NND: vm.NAND,
		instruction.XOR: vm.XOR,
		instruction.SUB: vm.SUB,
		instruction.JSR: vm.JSR,
		instruction.JMP: vm.JMP,
		instruction.JNS: vm.JNS,
		instruction.JNZ: vm.JNZ,
		instruction.HLT: vm.HALT,
		instruction.OUT: vm.OUT,
		instruction.INP: vm.IN,
	}

	vm.InputDevices = map[uint16]func() int16{
		0x0000: vm.SerialIn,
	}

	vm.OutputDevices = map[uint16]func(int16){
		0x0000: vm.SerialOut,
	}

	return vm
}

func (v *VM) SerialIn() int16 {
	err := exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	if err != nil {
		panic(err)
	}

	err = exec.Command("stty", "-F", "/dev/tty", "echo").Run()
	if err != nil {
		panic(err)
	}

	outBuf := make([]byte, 1)
	_, err = os.Stdin.Read(outBuf)
	if err != nil {
		panic(err)
	}

	return int16(outBuf[0])
}

func (v *VM) SerialOut(i int16) {
	fmt.Printf("%c", byte(i))
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
	v.PC = uint16(i) - 1 //PC increment happens after the execution of the instruction
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

func (v *VM) OUT(i int16) {
	fn := v.OutputDevices[uint16(i)]
	fn(v.Accumulator)
}

func (v *VM) IN(i int16) {
	fn := v.InputDevices[uint16(i)]
	v.Accumulator = fn()
}

func (v *VM) Store() {
	i := instruction.Instruction(v.Memory[v.PC])

	if i.AddressingMode() == instruction.AM_IMM {
		v.Memory[i.Immediate()] = v.Accumulator
	} else if i.AddressingMode() == instruction.AM_REL {
		v.Memory[int(v.PC)+int(i.Immediate())] = v.Accumulator
	} else if i.AddressingMode() == instruction.AM_DIR {
		v.Memory[i.Address()] = v.Accumulator
	} else if i.AddressingMode() == instruction.AM_IND {
		v.Memory[uint16(v.Memory[i.Address()])] = v.Accumulator
	} else if i.AddressingMode() == instruction.AM_INC {
		v.Memory[uint16(v.Memory[i.Address()])] = v.Accumulator
		v.Memory[i.Address()]++
	} else {
		v.Memory[i.Address()]--
		v.Memory[uint16(v.Memory[i.Address()])] = v.Accumulator
	}
}

func (v *VM) InvalidInstructionTrap() {
	fmt.Printf("Invalid Instruction at '0x%X'\n", v.PC)
	v.Accumulator = int16(v.PC + 1)
	v.PC = 1
}

func (v *VM) SingleStep() {
	i := instruction.Instruction(v.Memory[v.PC])

	if v.HFF {
		return
	} else if !i.IsValid() {
		v.InvalidInstructionTrap()
	} else if i.Opcode() == instruction.STA {
		v.Store()
	} else {
		op := v.Operations[i.Opcode()]
		var val int16

		if i.AddressingMode() == instruction.AM_IMM {
			val = i.Immediate()
		} else if i.AddressingMode() == instruction.AM_DIR {
			val = v.Memory[i.Address()]
		} else if i.AddressingMode() == instruction.AM_OFF {
			val = int16(int(v.PC) + int(i.Immediate()))
		} else if i.AddressingMode() == instruction.AM_REL {
			val = v.Memory[int(v.PC)+int(i.Immediate())]
		} else if i.AddressingMode() == instruction.AM_IND {
			val = v.Memory[uint16(v.Memory[i.Address()])]
		} else if i.AddressingMode() == instruction.AM_INC {
			val = v.Memory[uint16(v.Memory[i.Address()])]
			v.Memory[i.Address()]++
		} else if i.AddressingMode() == instruction.AM_DEC {
			v.Memory[i.Address()]--
			val = v.Memory[uint16(v.Memory[i.Address()])]
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
