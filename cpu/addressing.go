package cpu

// Addressing

const (
	Immediate           = iota
	ZeroPage            = iota
	Relative            = iota
	Implied             = iota
	Absolute            = iota
	Accumulator         = iota
	ZeroPageX           = iota
	ZeroPageY           = iota
	AbsoluteX           = iota
	AbsoluteY           = iota
	PreIndexedIndirect  = iota
	PostIndexedIndirect = iota
	IndirectAbsolute    = iota
)
