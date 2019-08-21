package common

func B2i(b bool) uint16 {
	if b {
		return 1
	}
	return 0
}

func B2ix(b bool, true uint16, false uint16) uint16 {
	if b {
		return true
	}
	return false
}

func I2ix(i uint, true uint, false uint) uint {
	if i > 0 {
		return true
	}
	return false
}

func I2b(i uint) bool {
	return i > 0
}

func Int2b(i int) bool {
	return i > 0
}
