package ppu

import "github.com/popsul/gones/bus"

type Palette struct {
	ram bus.Ram
}

func NewPalette() *Palette {
	palette := new(Palette)
	palette.ram = *bus.NewRam(0x20)
	return palette
}

func (P *Palette) IsSpriteMirror(addr uint16) bool {
	return (addr == 0x10) || (addr == 0x14) || (addr == 0x18) || (addr == 0x1c)
}
func (P *Palette) IsBackgroundMirror(addr uint16) bool {
	return (addr == 0x04) || (addr == 0x08) || (addr == 0x0c)
}
func (P *Palette) Read() []byte {
	ret := make([]byte, P.ram.Size())
	var i uint16
	for i = 0; i < P.ram.Size(); i++ {
		if P.IsSpriteMirror(i) {
			ret[i] = P.ram.Read(i - 0x10)
		} else if P.IsBackgroundMirror(i) {
			ret[i] = P.ram.Read(0x00)
		} else {
			ret[i] = P.ram.Read(i)
		}
	}
	return ret
}
func (P *Palette) GetPaletteAddr(addr uint16) uint16 {
	var mirrorDowned = (addr & 0xFF) % 0x20
	//NOTE: 0x3f10, 0x3f14, 0x3f18, 0x3f1c is mirror of 0x3f00, 0x3f04, 0x3f08, 0x3f0c
	if P.IsSpriteMirror(mirrorDowned) {
		return mirrorDowned - 0x10
	} else {
		return mirrorDowned
	}
}

func (P *Palette) Write(addr uint16, data byte) {
	//$this->ram[$this->getPaletteAddr($addr)] = $data;
	P.ram.Write(P.GetPaletteAddr(addr), data)
}
