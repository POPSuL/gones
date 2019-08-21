package main

import (
	"github.com/popsul/gones/apu"
	"github.com/popsul/gones/bus"
	"github.com/popsul/gones/cpu"
	"github.com/popsul/gones/interrupts"
	"github.com/popsul/gones/ppu"
	"github.com/popsul/gones/reader"
	"os"
	"runtime"
)

type Nes struct {
	ppu *ppu.Ppu

	cpu        *cpu.Cpu
	dma        *cpu.Dma
	interrupts *interrupts.Interrupts

	cpuBus       *cpu.CpuBus
	ram          *bus.Ram
	ppuBus       *bus.PpuBus
	characterMem *bus.Ram
	programPom   *bus.Rom
	keypad       *bus.Keypad

	renderer *ppu.Renderer
}

func NewNes(rom *reader.NesRom) *Nes {
	nes := new(Nes)

	nes.keypad = bus.NewKeypad()
	nes.ram = bus.NewRam(2048)

	nes.characterMem = bus.NewRam(0x4000)
	nes.characterMem.Fill(rom.Character)

	nes.programPom = bus.NewRom(rom.Program)

	nes.ppuBus = bus.NewPpuBus(nes.characterMem)
	nes.interrupts = interrupts.NewInterrupts()

	nes.ppu = ppu.NewPpu(nes.ppuBus, nes.interrupts, rom.HorizontalMirror)
	nes.dma = cpu.NewDma(nes.ram, nes.ppu)

	a := apu.NewApu()

	nes.cpuBus = cpu.NewCpuBus(nes.ram, nes.programPom, nes.ppu, a, nes.keypad, nes.dma)
	nes.cpu = cpu.NewCpu(nes.cpuBus, nes.interrupts)
	nes.cpu.Reset()

	nes.renderer = ppu.NewRenderer(nes.keypad)

	return nes
}

func (N *Nes) Frame() {
	for true {
		var cycle uint = 0
		if N.dma.IsDmaProcessing() {
			N.dma.Run()
			cycle = 514
		}
		cycle += N.cpu.Run()
		renderingData := N.ppu.Run(cycle * 3)
		if renderingData != nil {
			//fmt.Printf("RenderingData is not nil!\n")
			//	N.cpu.bus->keypad->fetch();
			N.renderer.Render(renderingData)
			break
		}
	}
}

func (N *Nes) Dump() {
	N.cpu.Dump()
}

func main() {
	var nesFile = os.Args[1]
	println("input file: ", nesFile)

	rom := reader.ReadRom(nesFile)
	nes := NewNes(rom)
	//nes.Dump()
	//return
	for true {
		nes.Frame()
		runtime.Gosched()
	}
}
