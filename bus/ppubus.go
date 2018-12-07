package bus

type PpuBus struct {
	ram *Ram
}

func NewPpuBus(ram *Ram) *PpuBus {
	ppuBus := new(PpuBus)
	ppuBus.ram = ram
	return ppuBus
}

func (p *PpuBus) ReadByPpu(addr uint) byte {
	return p.ram.Read(addr)
}

func (p *PpuBus) WriteByPpu(addr uint, value byte) {
	p.ram.Write(addr, value)
}
