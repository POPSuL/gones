package cpu

import (
	"github.com/popsul/gones/apu"
	"github.com/popsul/gones/bus"
	"github.com/popsul/gones/ppu"
)

type CpuBus struct {
	ram        *bus.Ram
	programRom *bus.Rom
	ppu        *ppu.Ppu
	dma        *Dma
	keypad1    *bus.Keypad
	keypad2    *bus.Keypad
	apu        *apu.Apu
}

func NewCpuBus(ram *bus.Ram, prgRom *bus.Rom, ppu *ppu.Ppu, apu *apu.Apu, keypad1 *bus.Keypad, keypad2 *bus.Keypad, dma *Dma) *CpuBus {
	cb := new(CpuBus)
	cb.ram = ram
	cb.programRom = prgRom
	cb.ppu = ppu
	cb.dma = dma
	cb.keypad1 = keypad1
	cb.keypad2 = keypad2
	cb.apu = apu
	return cb
}

func (CB *CpuBus) ReadByCpu(addr uint) byte {
	var data byte = 0
	if addr < 0x0800 {
		data = CB.ram.Read(addr)
	} else if addr < 0x2000 {
		// mirror
		data = CB.ram.Read(addr - 0x0800)
	} else if addr < 0x4000 {
		// mirror
		data = CB.ppu.Read((addr - 0x2000) % 8)
	} else if addr == 0x4016 {
		if CB.keypad1.Read() {
			data = 1
		}
	} else if addr == 0x4017 {
		if CB.keypad2.Read() {
			data = 1
		}
	} else if addr >= 0xC000 {
		// Mirror, if prom block number equals 1
		if CB.programRom.Size() <= 0x4000 {
			data = CB.programRom.Read(addr - 0xC000)
		} else {
			data = CB.programRom.Read(addr - 0x8000)
		}
	} else if addr >= 0x8000 {
		// ROM
		data = CB.programRom.Read(addr - 0x8000)
	}

	return data
}

func (CB *CpuBus) WriteByCpu(addr uint, data byte) {
	if addr < 0x2000 {
		CB.ram.Write(addr%0x0800, data)
	} else if addr < 0x2008 {
		// PPU
		CB.ppu.Write(addr-0x2000, data)
	} else if addr >= 0x4000 && addr < 0x4020 {
		if addr == 0x4014 {
			CB.dma.Write(data)
		} else if addr == 0x4016 {
			CB.keypad1.Write(data)
		} else if addr == 0x4017 {
			CB.keypad2.Write(data)
		} else {
			CB.apu.Write(addr-0x4000, data)
		}
	}
}
