package cpu

// Addressing

type Addressing int

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

var AddressingName = map[Addressing]string{
	Immediate:           "immediate",
	ZeroPage:            "zeroPage",
	Relative:            "relative",
	Implied:             "implied",
	Absolute:            "absolute",
	Accumulator:         "accumulator",
	ZeroPageX:           "zeropagex",
	ZeroPageY:           "zeropagey",
	AbsoluteX:           "absoluteX",
	AbsoluteY:           "absoluteY",
	PreIndexedIndirect:  "preIndexedIndirect",
	PostIndexedIndirect: "postIndexedIndirect",
	IndirectAbsolute:    "indirectAbsolute",
}

var AddressingNameShort = map[Addressing]string{
	Immediate:           "imm",
	ZeroPage:            "zer",
	Relative:            "rel",
	Implied:             "imp",
	Absolute:            "abs",
	Accumulator:         "acc",
	ZeroPageX:           "zerX",
	ZeroPageY:           "zerY",
	AbsoluteX:           "absX",
	AbsoluteY:           "absY",
	PreIndexedIndirect:  "preII",
	PostIndexedIndirect: "posII",
	IndirectAbsolute:    "indAbs",
}
