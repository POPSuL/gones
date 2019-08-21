package bus

import "fmt"

type Ram struct {
	size uint16
	data []byte
}

func NewRam(size uint16) *Ram {
	return &Ram{
		size: size,
		data: make([]byte, size),
	}
}

func (R *Ram) Write(addr uint16, val byte) {
	R.data[addr] = val
}

func (R *Ram) Read(addr uint16) byte {
	return R.data[addr]
}

func (R *Ram) Fill(bytes []byte) {
	println("Fill", len(bytes))
	R.data = bytes
}

func (R *Ram) Size() uint16 {
	return R.size
}

func (R *Ram) Dump() {
	for _, x := range R.data {
		fmt.Printf("0x%02x ", x)
	}
	fmt.Printf("\n")
}
