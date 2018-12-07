package interrupts

type Interrupts struct {
	nmi bool
	irq bool
}

func NewInterrupts() *Interrupts {
	return new(Interrupts).init()
}

func (I *Interrupts) init() *Interrupts {
	I.nmi = false
	I.irq = false
	return I
}

func (I *Interrupts) IsNmiAssert() bool {
	return I.nmi
}

func (I *Interrupts) IsIrqAssert() bool {
	return I.irq
}

func (I *Interrupts) AssertNmi() {
	I.nmi = true
}

func (I *Interrupts) ReleaseNmi() {
	I.nmi = false
}

func (I *Interrupts) AssertIrq() {
	I.irq = true
}

func (I *Interrupts) ReleaseIrq() {
	I.irq = false
}
