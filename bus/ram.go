package bus

type Ram struct {
	size uint
	data []byte
}

func NewRam(size uint) *Ram {
	ram := new(Ram)
	ram.size = size
	ram.Reset()
	return ram
}

func (R *Ram) Reset()  {
	R.data = make([]byte, R.size)
}

func (R *Ram) Write(addr uint, val byte) {
	R.data[addr] = byte(val)
}

func (R *Ram) Read(addr uint) byte {
	return R.data[addr]
}
func (R *Ram) Fill(bytes []byte) {
	R.data = bytes
}

func (R *Ram) Size() uint {
	return R.size
}