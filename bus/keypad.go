package bus

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

func (K *Keypad) Read() bool {
	return false
}

func (K *Keypad) Write(data byte) {
	if data&0x01 > 0 {
		K.isSet = true
	} else {
		K.isSet = false
		K.index = 0
		K.keyRegisters = K.keyBuffer
	}
}
