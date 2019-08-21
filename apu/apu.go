package apu

import "github.com/popsul/gones/bus"

type Apu struct {
	ram bus.Ram
}

func NewApu() *Apu {
	return &Apu{
		ram: *bus.NewRam(0x1f),
	}
}

func (A *Apu) Write(addr uint16, data byte) {
	A.ram.Write(addr, data)
}
