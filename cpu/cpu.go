package cpu

import (
	"errors"
	"fmt"
	"github.com/popsul/gones/interrupts"
)

const cpuClock = 1789772.5

func b2i(b bool) uint {
	if b {
		return 1
	}
	return 0
}

func i2b(i uint) bool {
	return i > 0
}

type Cpu struct {
	bus         *CpuBus
	interrupts  *interrupts.Interrupts
	registers   *Registers
	hasBranched bool
}

type AddrOrDataAndAdditionalCycle struct {
	addrOrData, additionalCycle uint
}

func newAddrOrDataAndAdditionalCycle(addrOrData uint, additionalCycle uint) AddrOrDataAndAdditionalCycle {
	return AddrOrDataAndAdditionalCycle{
		addrOrData,
		additionalCycle,
	}
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

func (C *Cpu) Fetch(addr uint, asWord bool) uint {
	if asWord {
		C.registers.PC += 2
	} else {
		C.registers.PC += 1
	}
	return C.Read(addr, asWord)
}

func (C *Cpu) Read(addr uint, asWord bool) uint {
	addr &= 0xFFFF
	if asWord {
		return uint(C.bus.ReadByCpu(addr)) | uint(C.bus.ReadByCpu(addr+1))<<8
	}
	return uint(C.bus.ReadByCpu(addr))
}

func (C *Cpu) Write(addr uint, data byte) {
	C.bus.WriteByCpu(addr, data)
}

func (C *Cpu) Push(data byte) {
	C.Write(0x100|(C.registers.SP&0xFF), data)
	C.registers.SP--
}

func (C *Cpu) Pop() uint {
	C.registers.SP++
	return C.Read(0x100|(C.registers.SP&0xFF), false)
}

func (C *Cpu) Branch(addr uint) {
	C.registers.PC = addr
	C.hasBranched = true
}

func (C *Cpu) pushStatus() {
	status := byte(b2i(C.registers.P.Negative)<<7 |
		b2i(C.registers.P.Overflow)<<6 |
		b2i(C.registers.P.Reserved)<<5 |
		b2i(C.registers.P.BreakMode)<<4 |
		b2i(C.registers.P.DecimalMode)<<3 |
		b2i(C.registers.P.Interrupt)<<2 |
		b2i(C.registers.P.Zero)<<1 |
		b2i(C.registers.P.Carry))
	C.Push(status)
}

func (C *Cpu) popStatus() {
	status := C.Pop()
	C.registers.P.Negative = i2b(status & 0x80)
	C.registers.P.Overflow = i2b(status & 0x40)
	C.registers.P.Reserved = i2b(status & 0x20)
	C.registers.P.BreakMode = i2b(status & 0x10)
	C.registers.P.DecimalMode = i2b(status & 0x08)
	C.registers.P.Interrupt = i2b(status & 0x04)
	C.registers.P.Zero = i2b(status & 0x02)
	C.registers.P.Carry = i2b(status & 0x01)
}

func (C *Cpu) PopPC() {
	C.registers.PC = C.Pop()
	C.registers.PC += C.Pop() << 8
}

func (C *Cpu) ProcessNmi() {
	C.interrupts.ReleaseNmi()
	C.registers.P.BreakMode = false
	C.Push(byte((C.registers.PC >> 8) & 0xFF))
	C.Push(byte(C.registers.PC & 0xFF))
	C.pushStatus()
	C.registers.P.Interrupt = true
	C.registers.PC = C.Read(0xFFFA, true)
}

func (C *Cpu) processIrq() {
	if C.registers.P.Interrupt {
		return
	}
	C.interrupts.ReleaseIrq()
	C.registers.P.BreakMode = false
	C.Push(byte((C.registers.PC >> 8) & 0xFF))
	C.Push(byte(C.registers.PC & 0xFF))
	C.pushStatus()
	C.registers.P.Interrupt = true
	C.registers.PC = C.Read(0xFFFE, false)
}

func (C *Cpu) getAddrOrDataWithAdditionalCycle(mode Addressing) AddrOrDataAndAdditionalCycle {
	switch mode {
	case Accumulator:
		return newAddrOrDataAndAdditionalCycle(0x00, 0)
	case Implied:
		return newAddrOrDataAndAdditionalCycle(0x00, 0)
	case Immediate:
		return newAddrOrDataAndAdditionalCycle(C.Fetch(C.registers.PC, false), 0)
	case Relative:
		baseAddr := C.Fetch(C.registers.PC, false)
		var addr uint
		if baseAddr < 0x80 {
			addr = baseAddr + C.registers.PC
		} else {
			addr = baseAddr + C.registers.PC - 256
		}
		return newAddrOrDataAndAdditionalCycle(
			addr,
			b2i((addr&0xff00) != (C.registers.PC&0xFF00)),
		)
	case ZeroPage:
		return newAddrOrDataAndAdditionalCycle(C.Fetch(C.registers.PC, false), 0)
	case ZeroPageX:
		addr := C.Fetch(C.registers.PC, false)
		return newAddrOrDataAndAdditionalCycle(
			(addr+C.registers.X)&0xff,
			0,
		)
	case ZeroPageY:
		addr := C.Fetch(C.registers.PC, false)
		return newAddrOrDataAndAdditionalCycle(addr+C.registers.Y&0xff, 0)
	case Absolute:
		return newAddrOrDataAndAdditionalCycle(C.Fetch(C.registers.PC, true), 0)
	case AbsoluteX:
		addr := C.Fetch(C.registers.PC, true)
		additionalCycle := b2i((addr & 0xFF00) != ((addr + C.registers.X) & 0xFF00))
		return newAddrOrDataAndAdditionalCycle((addr+C.registers.X)&0xFFFF, additionalCycle)
	case AbsoluteY:
		addr := C.Fetch(C.registers.PC, true)
		additionalCycle := b2i((addr & 0xFF00) != ((addr + C.registers.Y) & 0xFF00))
		return newAddrOrDataAndAdditionalCycle((addr+C.registers.Y)&0xFFFF, additionalCycle)
	case PreIndexedIndirect:
		baseAddr := (C.Fetch(C.registers.PC, false) + C.registers.X) & 0xFF
		addr := C.Read(baseAddr, false) + (C.Read((baseAddr+1)&0xFF, false) << 8)
		return newAddrOrDataAndAdditionalCycle(
			addr&0xFFFF,
			b2i((addr&0xFF00) != (baseAddr&0xFF00)),
		)
	case PostIndexedIndirect:
		addrOrData := C.Fetch(C.registers.PC, false)
		baseAddr := C.Read(addrOrData, false) + (C.Read((addrOrData+1)&0xFF, false) << 8)
		addr := baseAddr + C.registers.Y
		return newAddrOrDataAndAdditionalCycle(
			addr&0xFFFF,
			b2i((addr&0xFF00) != (baseAddr&0xFF00)),
		)
	case IndirectAbsolute:
		addrOrData := C.Fetch(C.registers.PC, true)
		addr := C.Read(addrOrData, false) +
			(C.Read((addrOrData&0xFF00)|(((addrOrData&0xFF)+1)&0xFF), false) << 8)
		return newAddrOrDataAndAdditionalCycle(addr&0xFFFF, 0)
	default:
		fmt.Printf("Mode: %d\n", mode)
		panic(errors.New("Unknown addressing mode detected."))
	}
}

func (C *Cpu) execOpCode(op uint) {
	opInfo := OpCodes[op]
	data := C.getAddrOrDataWithAdditionalCycle(opInfo.Addressing)
	fmt.Printf("OP %02x (%s) ADDR: %04x (%02x)\n", op, opInfo.BaseName, data.addrOrData, opInfo.Addressing)
}

func (C *Cpu) Run() uint {
	if C.interrupts.IsNmiAssert() {
		C.ProcessNmi()
	}
	if C.interrupts.IsIrqAssert() {
		C.processIrq()
	}

	opcode := C.Fetch(C.registers.PC, false)
	ocp := OpCodes[opcode]
	data := C.getAddrOrDataWithAdditionalCycle(ocp.Addressing)
	C.execOpCode(opcode)
	return ocp.Cycle + data.additionalCycle + b2i(C.hasBranched)
}
