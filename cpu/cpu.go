package cpu

import (
	"fmt"
	"github.com/popsul/gones/interrupts"
)

type Cpu struct {
	bus         *CpuBus
	interrupts  *interrupts.Interrupts
	registers   *Registers
	hasBranched bool
}

func NewCpu(bus *CpuBus, interrupts *interrupts.Interrupts) *Cpu {
	cpu := new(Cpu)
	cpu.bus = bus
	cpu.interrupts = interrupts
	cpu.registers = NewRegisters()
	cpu.hasBranched = false
	return cpu
}

func (C *Cpu) Reset() {
	C.registers = NewRegisters()
	// TODO: flownes set 0x8000 to PC when read(0xfffc) fails.
	C.registers.PC = uint(C.Read(0xfffc, true))
	fmt.Printf("Initial pc: %04x\n", C.registers.PC)
}

func (C *Cpu) Read(addr uint, asWord bool) uint {
	addr &= 0xFFFF
	if asWord {
		return uint(C.bus.ReadByCpu(addr)) | uint(C.bus.ReadByCpu(addr+1))<<8
	}
	return uint(C.bus.ReadByCpu(addr))
}
