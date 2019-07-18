package cpu

import (
	"github.com/popsul/gones/bus"
	"github.com/popsul/gones/ppu"
)

type Dma struct {
	isProcessing bool
	ramAddr      uint
	ram          *bus.Ram
	ppu          *ppu.Ppu
	addr         uint
	cycle        uint
}

func NewDma(ram *bus.Ram, ppu *ppu.Ppu) *Dma {
	return new(Dma).init(ram, ppu)
}

func (D *Dma) init(ram *bus.Ram, ppu *ppu.Ppu) *Dma {
	dma := new(Dma)
	dma.isProcessing = false
	dma.ramAddr = 0x0000
	dma.ram = ram
	dma.ppu = ppu
	dma.addr = 0x0000
	dma.cycle = 0x0000
	return dma
}

func (D *Dma) IsDmaProcessing() bool {
	return D.isProcessing
}

func (D *Dma) Run() {
	if !D.isProcessing {
		return
	}

	for i := uint(0); i < 0x100; i++ {
		D.ppu.TransferSprite(i, D.ram.Read(D.ramAddr+i))
	}
	D.isProcessing = false
}

func (D *Dma) Write(data byte) {
	D.ramAddr = uint(data << 8)
	D.isProcessing = true
}
