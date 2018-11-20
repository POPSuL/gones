package bus

type Rom struct {
	data []byte
}

func NewRom(r []byte) *Rom {
	rom := new(Rom)
	rom.data = r
	return rom
}

func (R *Rom) Size() uint {
	return uint(len(R.data))
}

func (R *Rom) Read(addr uint) byte {
	return R.data[addr]
}