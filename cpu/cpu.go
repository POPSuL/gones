package cpu

import (
	"errors"
	"fmt"
	. "github.com/popsul/gones/common"
	"github.com/popsul/gones/interrupts"
)

const cpuClock = 1789772.5

var instrNumber = 0

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
		return uint(C.bus.ReadByCpu(addr)) | (uint(C.bus.ReadByCpu(addr+1)) << 8)
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
	status := byte(B2i(C.registers.P.Negative)<<7 |
		B2i(C.registers.P.Overflow)<<6 |
		B2i(C.registers.P.Reserved)<<5 |
		B2i(C.registers.P.BreakMode)<<4 |
		B2i(C.registers.P.DecimalMode)<<3 |
		B2i(C.registers.P.Interrupt)<<2 |
		B2i(C.registers.P.Zero)<<1 |
		B2i(C.registers.P.Carry))
	C.Push(status)
}

func (C *Cpu) popStatus() {
	status := C.Pop()
	C.registers.P.Negative = I2b(status & 0x80)
	C.registers.P.Overflow = I2b(status & 0x40)
	C.registers.P.Reserved = I2b(status & 0x20)
	C.registers.P.BreakMode = I2b(status & 0x10)
	C.registers.P.DecimalMode = I2b(status & 0x08)
	C.registers.P.Interrupt = I2b(status & 0x04)
	C.registers.P.Zero = I2b(status & 0x02)
	C.registers.P.Carry = I2b(status & 0x01)
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
			B2i((addr&0xff00) != (C.registers.PC&0xFF00)),
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
		additionalCycle := B2i((addr & 0xFF00) != ((addr + C.registers.X) & 0xFF00))
		return newAddrOrDataAndAdditionalCycle((addr+C.registers.X)&0xFFFF, additionalCycle)
	case AbsoluteY:
		addr := C.Fetch(C.registers.PC, true)
		additionalCycle := B2i((addr & 0xFF00) != ((addr + C.registers.Y) & 0xFF00))
		return newAddrOrDataAndAdditionalCycle((addr+C.registers.Y)&0xFFFF, additionalCycle)
	case PreIndexedIndirect:
		baseAddr := (C.Fetch(C.registers.PC, false) + C.registers.X) & 0xFF
		addr := C.Read(baseAddr, false) + (C.Read((baseAddr+1)&0xFF, false) << 8)
		return newAddrOrDataAndAdditionalCycle(
			addr&0xFFFF,
			B2i((addr&0xFF00) != (baseAddr&0xFF00)),
		)
	case PostIndexedIndirect:
		addrOrData := C.Fetch(C.registers.PC, false)
		baseAddr := C.Read(addrOrData, false) + (C.Read((addrOrData+1)&0xFF, false) << 8)
		addr := baseAddr + C.registers.Y
		return newAddrOrDataAndAdditionalCycle(
			addr&0xFFFF,
			B2i((addr&0xFF00) != (baseAddr&0xFF00)),
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

func (C *Cpu) dumpRegisters() {
	fmt.Printf(
		"0b%d%d%d%d%d%d%d%d SP: 0x%04x PC: 0x%04x A: 0x%02x X: 0x%02x Y: 0x%02x\n",
		B2i(C.registers.P.DecimalMode),
		B2i(C.registers.P.Zero),
		B2i(C.registers.P.Negative),
		B2i(C.registers.P.Interrupt),
		B2i(C.registers.P.BreakMode),
		B2i(C.registers.P.Carry),
		B2i(C.registers.P.Overflow),
		B2i(C.registers.P.Reserved),
		C.registers.SP,
		C.registers.PC,
		C.registers.A,
		C.registers.X,
		C.registers.Y,
	)
}

func (C *Cpu) execOpCode(op uint, dataInfo AddrOrDataAndAdditionalCycle) {
	opInfo := OpCodes[op]
	addrOrData := dataInfo.addrOrData
	mode := opInfo.Addressing
	instrNumber++
	//if instrNumber > 200000 {
	fmt.Printf(
		"OP %d (%s) ADDR: 0x%04x (%s)\n",
		instrNumber,
		opInfo.BaseName,
		dataInfo.addrOrData,
		AddressingName[opInfo.Addressing],
	)

	C.dumpRegisters()
	//}

	C.hasBranched = false
	switch opInfo.BaseName {
	case "LDA":
		C.registers.A = B2ix(mode == Immediate, addrOrData, C.Read(addrOrData, false))
		fmt.Printf("0x%04x 0x%04x\n", addrOrData, C.registers.A)
		C.registers.P.Negative = I2b(C.registers.A & 0x80)
		C.registers.P.Zero = !I2b(C.registers.A)
		break
	case "LDX":
		C.registers.X = B2ix(mode == Immediate, addrOrData, C.Read(addrOrData, false))
		C.registers.P.Negative = I2b(C.registers.X & 0x80)
		C.registers.P.Zero = !I2b(C.registers.X)
		break
	case "LDY":
		C.registers.Y = B2ix(mode == Immediate, addrOrData, C.Read(addrOrData, false))
		C.registers.P.Negative = I2b(C.registers.Y & 0x80)
		C.registers.P.Zero = !I2b(C.registers.Y)
		break
	case "STA":
		//fmt.Printf("STA %d %d\n", addrOrData, C.registers.A)
		C.Write(addrOrData, byte(C.registers.A))
		break
	case "STX":
		C.Write(addrOrData, byte(C.registers.X))
		break
	case "STY":
		C.Write(addrOrData, byte(C.registers.Y))
		break
	case "TAX":
		C.registers.X = C.registers.A
		C.registers.P.Negative = I2b(C.registers.X & 0x80)
		C.registers.P.Zero = !I2b(C.registers.X)
		break
	case "TAY":
		C.registers.Y = C.registers.A
		C.registers.P.Negative = I2b(C.registers.Y & 0x80)
		C.registers.P.Zero = !I2b(C.registers.Y)
		break
	case "TSX":
		C.registers.X = C.registers.SP & 0xFF
		C.registers.P.Negative = I2b(C.registers.X & 0x80)
		C.registers.P.Zero = !I2b(C.registers.X)
		break
	case "TXA":
		C.registers.A = C.registers.X
		C.registers.P.Negative = I2b(C.registers.A & 0x80)
		C.registers.P.Zero = !I2b(C.registers.A)
		break
	case "TXS":
		C.registers.SP = C.registers.X + 0x0100
		break
	case "TYA":
		C.registers.A = C.registers.Y
		C.registers.P.Negative = I2b(C.registers.A & 0x80)
		C.registers.P.Zero = !I2b(C.registers.A)
		break
	case "ADC":
		data := B2ix(mode == Immediate, addrOrData, C.Read(addrOrData, false))
		operated := data + C.registers.A + B2i(C.registers.P.Carry)
		overflow := !(((C.registers.A ^ data) & 0x80) != 0) && ((C.registers.A^operated)&0x80) != 0
		C.registers.P.Overflow = overflow
		C.registers.P.Carry = operated > 0xFF
		C.registers.P.Negative = I2b(operated & 0x80)
		C.registers.P.Zero = !I2b(operated & 0xFF)
		C.registers.A = operated & 0xFF
		break
	case "AND":
		data := B2ix(mode == Immediate, addrOrData, C.Read(addrOrData, false))
		operated := data & C.registers.A
		C.registers.P.Negative = I2b(operated & 0x80)
		C.registers.P.Zero = !I2b(operated)
		C.registers.A = operated & 0xFF
		break
	case "ASL":
		if mode == Accumulator {
			acc := C.registers.A
			C.registers.P.Carry = !!I2b(acc & 0x80)
			C.registers.A = (acc << 1) & 0xFF
			C.registers.P.Zero = !I2b(C.registers.A)
			C.registers.P.Negative = I2b(C.registers.A & 0x80)
		} else {
			data := C.Read(addrOrData, false)
			C.registers.P.Carry = I2b(data & 0x80)
			shifted := (data << 1) & 0xFF
			C.Write(addrOrData, byte(shifted))
			C.registers.P.Zero = !I2b(shifted)
			C.registers.P.Negative = I2b(shifted & 0x80)
		}
		break
	case "BIT":
		data := C.Read(addrOrData, false)
		C.registers.P.Negative = I2b(data & 0x80)
		C.registers.P.Overflow = I2b(data & 0x40)
		C.registers.P.Zero = !I2b(C.registers.A & data)
		break
	case "CMP":
		data := B2ix(mode == Immediate, addrOrData, C.Read(addrOrData, false))
		//fmt.Printf("CMP: %d\n", data)
		compared := int(C.registers.A) - int(data)
		C.registers.P.Carry = compared >= 0
		C.registers.P.Negative = I2b(uint(compared & 0x80))
		C.registers.P.Zero = !I2b(uint(compared & 0xff))
		break
	case "CPX":
		data := B2ix(mode == Immediate, addrOrData, C.Read(addrOrData, false))
		compared := int(C.registers.X) - int(data)
		//fmt.Printf("CPX: %d %d\n", compared, uint(compared)&0x80)
		C.registers.P.Carry = compared >= 0
		C.registers.P.Negative = I2b(uint(compared) & 0x80)
		//fmt.Printf("NEG: %d\n", uint(compared)&0xff)
		C.registers.P.Zero = !I2b(uint(compared) & 0xff)
		break
	case "CPY":
		data := B2ix(mode == Immediate, addrOrData, C.Read(addrOrData, false))
		compared := int(C.registers.Y) - int(data)
		C.registers.P.Carry = compared >= 0
		C.registers.P.Negative = I2b(uint(compared & 0x80))
		C.registers.P.Zero = !I2b(uint(compared & 0xff))
		break
	case "DEC":
		data := (C.Read(addrOrData, false) - 1) & 0xFF
		C.registers.P.Negative = I2b(data & 0x80)
		C.registers.P.Zero = !I2b(data)
		C.Write(addrOrData, byte(data))
		break
	case "DEX":
		C.registers.X = (C.registers.X - 1) & 0xFF
		C.registers.P.Negative = I2b(C.registers.X & 0x80)
		C.registers.P.Zero = !I2b(C.registers.X)
		break
	case "DEY":
		C.registers.Y = (C.registers.Y - 1) & 0xFF
		C.registers.P.Negative = I2b(C.registers.Y & 0x80)
		C.registers.P.Zero = !I2b(C.registers.Y)
		break
	case "EOR":
		data := B2ix(mode == Immediate, addrOrData, C.Read(addrOrData, false))
		operated := data ^ C.registers.A
		C.registers.P.Negative = I2b(operated & 0x80)
		C.registers.P.Zero = !I2b(operated)
		C.registers.A = operated & 0xFF
		break
	case "INC":
		data := (C.Read(addrOrData, false) + 1) & 0xFF
		C.registers.P.Negative = I2b(data & 0x80)
		C.registers.P.Zero = !I2b(data)
		C.Write(addrOrData, byte(data))
		break
	case "INX":
		C.registers.X = (C.registers.X + 1) & 0xFF
		C.registers.P.Negative = I2b(C.registers.X & 0x80)
		C.registers.P.Zero = !I2b(C.registers.X)
		break
	case "INY":
		C.registers.Y = (C.registers.Y + 1) & 0xFF
		C.registers.P.Negative = I2b(C.registers.Y & 0x80)
		C.registers.P.Zero = !I2b(C.registers.Y)
		break
	case "LSR":
		if mode == Accumulator {
			acc := C.registers.A & 0xFF
			C.registers.P.Carry = I2b(acc & 0x01)
			C.registers.A = acc >> 1
			C.registers.P.Zero = !I2b(C.registers.A)
		} else {
			data := C.Read(addrOrData, false)
			C.registers.P.Carry = I2b(data & 0x01)
			C.registers.P.Zero = !I2b(data >> 1)
			C.Write(addrOrData, byte(data>>1))
		}
		C.registers.P.Negative = false
		break
	case "ORA":
		data := B2ix(mode == Immediate, addrOrData, C.Read(addrOrData, false))
		operated := data | C.registers.A
		C.registers.P.Negative = I2b(operated & 0x80)
		C.registers.P.Zero = !I2b(operated)
		C.registers.A = operated & 0xFF
		break
	case "ROL":
		if mode == Accumulator {
			acc := C.registers.A
			C.registers.A = (acc<<1)&0xFF | B2ix(C.registers.P.Carry, 0x01, 0x00)
			C.registers.P.Carry = I2b(acc & 0x80)
			C.registers.P.Zero = !I2b(C.registers.A)
			C.registers.P.Negative = I2b(C.registers.A & 0x80)
		} else {
			data := C.Read(addrOrData, false)
			writeData := (data<<1 | B2i(C.registers.P.Carry)) & 0xFF
			C.Write(addrOrData, byte(writeData))
			C.registers.P.Carry = !!I2b(data & 0x80)
			C.registers.P.Zero = !I2b(writeData)
			C.registers.P.Negative = I2b(writeData & 0x80)
		}
		break
	case "ROR":
		if mode == Accumulator {
			acc := C.registers.A
			C.registers.A = acc>>1 | B2ix(C.registers.P.Carry, 0x80, 0x00)
			C.registers.P.Carry = I2b(acc & 0x01)
			C.registers.P.Zero = !I2b(C.registers.A)
			C.registers.P.Negative = I2b(C.registers.A & 0x80)
		} else {
			data := C.Read(addrOrData, false)
			writeData := data>>1 | B2ix(C.registers.P.Carry, 0x80, 0x00)
			C.Write(addrOrData, byte(writeData))
			C.registers.P.Carry = I2b(data & 0x01)
			C.registers.P.Zero = !I2b(writeData)
			C.registers.P.Negative = I2b(writeData & 0x80)
		}
		break
	case "SBC":
		data := B2ix(mode == Immediate, addrOrData, C.Read(addrOrData, false))
		operated := C.registers.A - data - B2ix(C.registers.P.Carry, 0, 1)
		overflow := ((C.registers.A^operated)&0x80) != 0 && ((C.registers.A^data)&0x80) != 0
		C.registers.P.Overflow = overflow
		C.registers.P.Carry = operated >= 0
		C.registers.P.Negative = I2b(operated & 0x80)
		C.registers.P.Zero = !I2b(operated & 0xFF)
		C.registers.A = operated & 0xFF
		break
	case "PHA":
		C.Push(byte(C.registers.A))
		break
	case "PHP":
		C.registers.P.BreakMode = true
		C.pushStatus()
		break
	case "PLA":
		C.registers.A = C.Pop()
		C.registers.P.Negative = I2b(C.registers.A & 0x80)
		C.registers.P.Zero = !I2b(C.registers.A)
		break
	case "PLP":
		C.popStatus()
		C.registers.P.Reserved = true
		break
	case "JMP":
		C.registers.PC = addrOrData
		break
	case "JSR":
		pc := C.registers.PC - 1
		C.Push(byte((pc >> 8) & 0xFF))
		C.Push(byte(pc & 0xFF))
		C.registers.PC = addrOrData
		break
	case "RTS":
		C.PopPC()
		C.registers.PC++
		break
	case "RTI":
		C.popStatus()
		C.PopPC()
		C.registers.P.Reserved = true
		break
	case "BCC":
		if !C.registers.P.Carry {
			C.Branch(addrOrData)
		}
		break
	case "BCS":
		if C.registers.P.Carry {
			C.Branch(addrOrData)
		}
		break
	case "BEQ":
		if C.registers.P.Zero {
			C.Branch(addrOrData)
		}
		break
	case "BMI":
		if C.registers.P.Negative {
			C.Branch(addrOrData)
		}
		break
	case "BNE":
		if !C.registers.P.Zero {
			C.Branch(addrOrData)
		}
		break
	case "BPL":
		if !C.registers.P.Negative {
			C.Branch(addrOrData)
		}
		break
	case "BVS":
		if C.registers.P.Overflow {
			C.Branch(addrOrData)
		}
		break
	case "BVC":
		if !C.registers.P.Overflow {
			C.Branch(addrOrData)
		}
		break
	case "CLD":
		C.registers.P.DecimalMode = false
		break
	case "CLC":
		C.registers.P.Carry = false
		break
	case "CLI":
		C.registers.P.Interrupt = false
		break
	case "CLV":
		C.registers.P.Overflow = false
		break
	case "SEC":
		C.registers.P.Carry = true
		break
	case "SEI":
		C.registers.P.Interrupt = true
		break
	case "SED":
		C.registers.P.DecimalMode = true
		break
	case "BRK":
		interrupt := C.registers.P.Interrupt
		C.registers.PC++
		C.Push(byte((C.registers.PC >> 8) & 0xFF))
		C.Push(byte(C.registers.PC & 0xFF))
		C.registers.P.BreakMode = true
		C.pushStatus()
		C.registers.P.Interrupt = true
		// Ignore interrupt when already set.
		if !interrupt {
			C.registers.PC = C.Read(0xFFFE, true)
		}
		C.registers.PC--
		break
	case "NOP":
		break
	// Unofficial OpCode
	case "NOPD":
		C.registers.PC++
		break
	case "NOPI":
		C.registers.PC += 2
		break
	case "LAX":
		data := C.Read(addrOrData, false)
		C.registers.A, C.registers.X = data, data
		C.registers.P.Negative = I2b(C.registers.A & 0x80)
		C.registers.P.Zero = !I2b(C.registers.A)
		break
	case "SAX":
		operated := C.registers.A & C.registers.X
		C.Write(addrOrData, byte(operated))
		break
	case "DCP":
		operated := (C.Read(addrOrData, false) - 1) & 0xFF
		C.registers.P.Negative = I2b(((C.registers.A - operated) & 0x1FF) & 0x80)
		C.registers.P.Zero = !I2b((C.registers.A - operated) & 0x1FF)
		C.Write(addrOrData, byte(operated))
		break
	case "ISB":
		data := (C.Read(addrOrData, false) + 1) & 0xFF
		operated := (^data & 0xFF) + C.registers.A + B2i(C.registers.P.Carry)
		overflow := !(((C.registers.A ^ data) & 0x80) != 0) && ((C.registers.A^operated)&0x80) != 0
		C.registers.P.Overflow = overflow
		C.registers.P.Carry = operated > 0xFF
		C.registers.P.Negative = I2b(operated & 0x80)
		C.registers.P.Zero = !I2b(operated & 0xFF)
		C.registers.A = operated & 0xFF
		C.Write(addrOrData, byte(data))
		break
	case "SLO":
		data := C.Read(addrOrData, false)
		C.registers.P.Carry = I2b(data & 0x80)
		data = (data << 1) & 0xFF
		C.registers.A |= data
		C.registers.P.Negative = I2b(C.registers.A & 0x80)
		C.registers.P.Zero = !I2b(C.registers.A & 0xFF)
		C.Write(addrOrData, byte(data))
		break
	case "RLA":
		data := (C.Read(addrOrData, false) << 1) + B2i(C.registers.P.Carry)
		C.registers.P.Carry = I2b(data & 0x100)
		C.registers.A = (data & C.registers.A) & 0xFF
		C.registers.P.Negative = I2b(C.registers.A & 0x80)
		C.registers.P.Zero = !I2b(C.registers.A & 0xFF)
		C.Write(addrOrData, byte(data))
		break
	case "SRE":
		data := C.Read(addrOrData, false)
		C.registers.P.Carry = I2b(data & 0x01)
		data >>= 1
		C.registers.A ^= data
		C.registers.P.Negative = I2b(C.registers.A & 0x80)
		C.registers.P.Zero = !I2b(C.registers.A & 0xFF)
		C.Write(addrOrData, byte(data))
		break
	case "RRA":
		data := C.Read(addrOrData, false)
		carry := data & 0x01
		data = (data >> 1) | B2ix(C.registers.P.Carry, 0x80, 0x00)
		operated := data + C.registers.A + carry
		overflow := !(((C.registers.A ^ data) & 0x80) != 0) && ((C.registers.A^operated)&0x80) != 0
		C.registers.P.Overflow = overflow
		C.registers.P.Negative = I2b(operated & 0x80)
		C.registers.P.Zero = !I2b(operated & 0xFF)
		C.registers.A = operated & 0xFF
		C.registers.P.Carry = operated > 0xFF
		C.Write(addrOrData, byte(data))
		break
	default:
		panic(errors.New(fmt.Sprintf("Unknown  opcode %d (%s detected)\n", op, opInfo.BaseName)))
	}
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
	C.execOpCode(opcode, data)
	//fmt.Printf(
	//	"%d+%d+%d=%d\n",
	//	ocp.Cycle,
	//	data.additionalCycle,
	//	B2i(C.hasBranched),
	//	ocp.Cycle + data.additionalCycle + B2i(C.hasBranched))
	return ocp.Cycle + data.additionalCycle + B2i(C.hasBranched)
}

func (C *Cpu) Dump() {
	for {
		pc := C.registers.PC
		opcode := C.Fetch(C.registers.PC, false)
		ocp := OpCodes[opcode]
		data := C.getAddrOrDataWithAdditionalCycle(ocp.Addressing)
		fmt.Printf("0x%04x\t%s\t%04x\n", pc, ocp.FullName, data.addrOrData)
	}
}
