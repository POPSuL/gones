package main

import (
	"github.com/popsul/gones/apu"
	"github.com/popsul/gones/bus"
	"github.com/popsul/gones/common"
	"github.com/popsul/gones/cpu"
	"github.com/popsul/gones/interrupts"
	"github.com/popsul/gones/ppu"
	"github.com/popsul/gones/reader"
	"os"
	"runtime"
	"time"
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
	keypad1      *bus.Keypad
	keypad2      *bus.Keypad
	apu          *apu.Apu

	renderer *ppu.Renderer
}

func NewNes(rom *reader.NesRom) *Nes {
	nes := new(Nes)

	nes.keypad1 = bus.NewKeypad()
	nes.keypad2 = bus.NewKeypad()
	nes.ram = bus.NewRam(2048)

	nes.characterMem = bus.NewRam(0x4000)
	nes.characterMem.Fill(rom.Character)

	nes.programPom = bus.NewRom(rom.Program)

	nes.ppuBus = bus.NewPpuBus(nes.characterMem)
	nes.interrupts = interrupts.NewInterrupts()

	nes.ppu = ppu.NewPpu(nes.ppuBus, nes.interrupts, rom.HorizontalMirror)
	nes.dma = cpu.NewDma(nes.ram, nes.ppu)

	nes.apu = apu.NewApu(nes.interrupts)

	nes.cpuBus = cpu.NewCpuBus(nes.ram, nes.programPom, nes.ppu, nes.apu, nes.keypad1, nes.keypad2, nes.dma)
	nes.cpu = cpu.NewCpu(nes.cpuBus, nes.interrupts)
	nes.cpu.Reset()

	nes.renderer = ppu.NewRenderer(nes.keypad1)

	return nes
}

func (N *Nes) Frame(deadline float64) {
	allowedCycles := deadline / 1000 / 1000 / 1000 * float64(common.CpuClock)
	for allowedCycles > 0 {
		var cycle uint = 0
		if N.dma.IsDmaProcessing() {
			N.dma.Run()
			cycle = 514
		}
		cpuCycles := N.cpu.Run()
		allowedCycles -= float64(cpuCycles)
		cycle += cpuCycles
		renderingData := N.ppu.Run(cycle * 3)
		N.apu.Run(cycle)
		if renderingData != nil {
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
	timestamp := time.Now().UnixNano()
	for true {
		now := time.Now().UnixNano()
		nes.Frame(float64(now - timestamp))
		timestamp = now
		runtime.Gosched()
	}
}
