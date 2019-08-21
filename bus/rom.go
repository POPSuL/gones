package bus

type Rom struct {
	data []byte
}

func NewRom(r []byte) *Rom {
	rom := new(Rom)
	rom.data = r
	return rom
}

func (R *Rom) Size() uint16 {
	return uint16(len(R.data))
}

func (R *Rom) Read(addr uint16) byte {
	return R.data[addr]
}
