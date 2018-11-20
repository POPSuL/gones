package bus

import (
	"github.com/popsul/gones/cpu"
	"github.com/popsul/gones/ppu"
)

type CpuBus struct {
	ram        Ram
	programRom Rom
	ppu        ppu.Ppu
	dma        cpu.Dma
	keypad	   Keypad
}

func NewCpuBus(ram Ram, prgRom Rom, ppu ppu.Ppu, dma cpu.Dma) *CpuBus {
	cb := new(CpuBus)
	cb.ram = ram
	cb.programRom = prgRom
	cb.ppu = ppu
	cb.dma = dma
	return cb
}

func (CB *CpuBus) ReadByCpu(addr uint) byte {
	if addr < 0x0800 {
		return CB.ram.Read(addr)
	} else if addr < 0x2000 {
		// mirror
		return CB.ram.Read(addr - 0x0800)
	} else if addr < 0x4000 {
		// mirror
		data := CB.ppu.Read(addr-0x2000) % 8
		return data
	} else if addr == 0x4016 {
		// TODO Add 2P
		return CB.keypad.read()
	} else if addr >= 0xC000 {
		// Mirror, if prom block number equals 1
		if CB.programRom.Size() <= 0x4000 {
			return CB.programRom.Read(addr - 0xC000)
		}
		return CB.programRom.Read(addr - 0x8000)
	} else if addr >= 0x8000 {
		// ROM
		return CB.programRom.Read(addr - 0x8000)
	}
	return 0
}

func (CB *CpuBus) WriteByCpu(addr uint, data byte) {
	if addr < 0x0800 {
		// RAM
		CB.ram.Write(addr, data)
	} else if addr < 0x2000 {
		// mirror
		CB.ram.Write(addr-0x0800, data)
	} else if addr < 0x2008 {
		// PPU
		CB.ppu.Write(addr-0x2000, data)
	} else if addr >= 0x4000 && addr < 0x4020 {
		if addr == 0x4014 {
			CB.dma.Write(data)
		} else if addr == 0x4016 {
			// TODO Add 2P
			CB.keypad.Write(data)
		}
	}
}
