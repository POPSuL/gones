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
	keypad     *bus.Keypad
}

func NewCpuBus(ram *bus.Ram, prgRom *bus.Rom, ppu *ppu.Ppu, alu *apu.Apu, keypad *bus.Keypad, dma *Dma) *CpuBus {
	cb := new(CpuBus)
	cb.ram = ram
	cb.programRom = prgRom
	cb.ppu = ppu
	cb.dma = dma
	cb.keypad = keypad
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
		// TODO Add 2P
		if CB.keypad.Read() {
			data = 1
		}
		data = 0
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
		} else {
			//CB.apu.Write(addr-0x4000, data)
		}
	}
}
