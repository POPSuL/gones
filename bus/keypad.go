package bus

import "fmt"

type Keypad struct {
	isSet        bool
	index        uint
	keyRegisters [8]bool
	keyBuffer    [8]bool
}

func NewKeypad() *Keypad {
	keypad := new(Keypad)
	return keypad
}

func (K *Keypad) DumpBuffer() {
	for _, x := range K.keyBuffer {
		fmt.Printf("%t ", x)
	}
	fmt.Printf("\n")
}

func (K *Keypad) DumpRegisters() {
	for _, x := range K.keyRegisters {
		fmt.Printf("%t ", x)
	}
	fmt.Printf("\n")
}

func (K *Keypad) Read() bool {
	//K.DumpRegisters()
	//K.DumpBuffer()
	k := K.keyRegisters[K.index]
	K.index++
	if K.index >= uint(len(K.keyRegisters)) {
		K.index = 0
	}
	return k
}

func (K *Keypad) Write(data byte) {
	if data&0x01 > 0 {
		K.isSet = true
	} else if K.isSet && (data&0x01) < 1 {
		K.isSet = false
		K.index = 0
		//K.keyRegisters = K.keyBuffer
	}
}

func (K *Keypad) KeyDown(key uint) {
	if key < 0x08 {
		K.keyRegisters[key] = true
	}
}

func (K *Keypad) KeyUp(key uint) {
	if key < 0x08 {
		K.keyRegisters[key] = false
	}
}
