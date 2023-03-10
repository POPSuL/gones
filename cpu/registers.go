package cpu

// Register's names
const (
	RA  = 0x00
	RX  = 0x01
	RY  = 0x02
	RP  = 0x03
	RSP = 0x04
	RPC = 0x05
)

type Registers struct {
	A, X, Y, SP, PC uint
	P               *Status
}

type Status struct {
	Negative    bool
	Overflow    bool
	Reserved    bool
	BreakMode   bool
	DecimalMode bool
	Interrupt   bool
	Zero        bool
	Carry       bool
}

func NewStatus(
	negative bool,
	overflow bool,
	reserved bool,
	breakMode bool,
	decimalMode bool,
	interrupt bool,
	zero bool,
	carry bool,
) *Status {
	return &Status{
		Negative:    negative,
		Overflow:    overflow,
		Reserved:    reserved,
		BreakMode:   breakMode,
		DecimalMode: decimalMode,
		Interrupt:   interrupt,
		Zero:        zero,
		Carry:       carry,
	}
}

func NewRegisters() *Registers {
	return &Registers{
		0x00,
		0x00,
		0x00,
		0x01fd,
		0x0000,
		NewStatus(
			false,
			false,
			true,
			true,
			false,
			true,
			false,
			false,
		),
	}
}
