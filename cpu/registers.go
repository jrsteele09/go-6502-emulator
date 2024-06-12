package cpu

type StatusFlag byte

const (
	CarryFlag            StatusFlag = 1 << iota // Bit 1
	ZeroFlag                        = 1 << iota // Bit 2
	InterruptDisableFlag            = 1 << iota // Bit 3
	DecimalFlag                     = 1 << iota // Bit 4
	BreakFlag                       = 1 << iota // Bit 5
	UnusedFlag                      = 1 << iota // Bit 6
	OverflowFlag                    = 1 << iota // Bit 7
	NegativeFlag                    = 1 << iota // Bit 8
)

type Registers struct {
	A, X, Y, S byte
	PC         uint16
	Status     byte
}

func NewRegisters() *Registers {
	return &Registers{}
}

func (r *Registers) SetStatus(f StatusFlag, on bool) {
	if on {
		r.Status = r.Status | byte(f)
	} else {
		r.Status = r.Status &^ byte(f)
	}
}

func (r *Registers) IsSet(f StatusFlag) bool {
	return (r.Status & byte(f)) != 0
}

func (r *Registers) SetCarryFlag(oldValue, result byte) {
	r.SetStatus(CarryFlag, oldValue > result)
}

func (r *Registers) SetZeroFlag(result byte) {
	r.SetStatus(ZeroFlag, result == 0)
}

func (r *Registers) SetNegativeFlag(result byte) {
	r.SetStatus(NegativeFlag, (result&0x80) != 0)
}

/*
In the 6502 processor, the Overflow flag is commonly used with the ADC (Add with Carry) and SBC (Subtract with Carry) instructions. Here's what it's useful for in these contexts:

Addition (ADC): In signed arithmetic, adding two positive numbers should result in another
positive number and adding two negative numbers should result in another negative number.
If this is not the case (i.e., adding two positive numbers yields a negative result,
or adding two negative numbers yields a positive result), then an overflow has occurred.

Subtraction (SBC): Similar to addition, in signed arithmetic, subtracting a negative number from a
positive number should yield a positive number, and subtracting a positive number from a negative
number should yield a negative result. If this is not the case, an overflow has occurred.
*/
func (r *Registers) SetOverflowFlag(m, n, result byte, isAddition bool) {
	if isAddition {
		r.SetStatus(OverflowFlag, ((m^result)&(n^result))&0x80 != 0)
	} else {
		r.SetStatus(OverflowFlag, ((m^result)&(m^n))&0x80 != 0)
	}
}
