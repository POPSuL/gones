package apu

import "C"
import (
	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/popsul/gones/bus"
	"github.com/popsul/gones/common"
	"github.com/popsul/gones/interrupts"
	"time"
)

var CounterTable = [...]uint{
	0x0A, 0xFE, 0x14, 0x02, 0x28, 0x04, 0x50, 0x06,
	0xA0, 0x08, 0x3C, 0x0A, 0x0E, 0x0C, 0x1A, 0x0E,
	0x0C, 0x10, 0x18, 0x12, 0x30, 0x14, 0x60, 0x16,
	0xC0, 0x18, 0x48, 0x1A, 0x10, 0x1C, 0x20, 0x1E,
}

type Generator interface {
	GetStreamer() beep.Streamer
	Write(byte, byte)
}

type Apu struct {
	ram           bus.Ram
	interrupts    *interrupts.Interrupts
	cycle         uint
	step          uint
	registers     [0x18]byte
	noise         Generator
	triangle      Generator
	square0       Generator
	square1       Generator
	enableIrq     bool
	sequencerMode bool
}

func NewApu(interrupts *interrupts.Interrupts) *Apu {
	a := &Apu{
		ram:        *bus.NewRam(0x1f),
		interrupts: interrupts,
		noise:      NewNoise(),
		triangle:   NewTriangle(),
		square0:    NewSquare(),
		square1:    NewSquare(),
	}
	sr := beep.SampleRate(22050)
	speaker.Init(sr, sr.N(time.Second/10))
	speaker.Play(a.noise.GetStreamer(), a.triangle.GetStreamer())
	return a
}

func (A *Apu) Write(addr uint, data byte) {
	//fmt.Printf("AP: 0x%02x - 0x%02x\n", addr, data)
	if addr <= 0x03 {
		A.square0.Write(byte(addr&0xff), data)
	} else if addr <= 0x07 {
		A.square1.Write(byte((addr-0x04)&0xff), data)
	} else if addr <= 0x0b {
		A.triangle.Write(byte((addr-0x08)&0xff), data)
	} else if addr <= 0x0f {
		A.noise.Write(byte((addr-0x0c)&0xff), data)
	} else if addr == 0x17 {
		A.sequencerMode = common.I2b(uint(data & 0x80))
		A.registers[addr] = data
		A.enableIrq = common.I2b(uint(data & 0x40))
	}
	//A.ram.Write(addr, data)
}

func (A *Apu) updateEnvelope() {

}

func (A *Apu) updateSweepAndLengthCounter() {

}

func (A *Apu) updateBySequenceMode0() {
	A.updateEnvelope()
	if A.step%2 == 1 {
		A.updateSweepAndLengthCounter()
	}
	A.step++
	if A.step == 4 {
		if A.enableIrq {
			A.interrupts.AssertIrq()
		}
		A.step = 0
	}
}

func (A *Apu) updateBySequenceMode1() {

}

func (A *Apu) Run(cycle uint) {
	A.cycle += cycle
	if A.cycle < common.FrameCounterRate {
		return
	}
	A.cycle -= common.FrameCounterRate
	//fmt.Printf("AP: SQ %t\n", A.sequencerMode)
	if A.sequencerMode {
		A.updateBySequenceMode1()
	} else {
		A.updateBySequenceMode0()
	}
}
