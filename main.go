package main

import (
	"github.com/popsul/gones/bus"
	"github.com/popsul/gones/cpu"
	"github.com/popsul/gones/ppu"
	"github.com/popsul/gones/reader"
	"os"
)

type Nes struct {
	ppu *ppu.Ppu
	ram *bus.Ram
	characterMem *bus.Ram
	programPom *bus.Rom
	dma *cpu.Dma
	interrupts *cpu.Interrupts
}

func NewNes(rom reader.NesRom) *Nes {
	nes := new(Nes)
	// todo: keyboard
	nes.ram = bus.NewRam(2048)

	nes.characterMem = bus.NewRam(0x4000)
	nes.characterMem.Fill(rom.Character)

	nes.programPom = bus.NewRom(rom.Program)

	// todo: ppuBus
	nes.interrupts = cpu.NewInterrupts()

	// todo: Допилить
	nes.ppu = ppu.NewPpu()
	nes.dma = cpu.NewDma(*nes.ram, *nes.ppu)

	// todo: cpuBus
	// todo: CPU

	return nes
}

func main() {
	var nesFile = os.Args[1]
	println("input file: ", nesFile)

	_ = reader.ReadRom(nesFile)
}
