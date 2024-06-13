package cpu

// StatusFlag represents the various status flags of the 6502 CPU.
type StatusFlag byte

const (
	// CarryFlag indicates whether an arithmetic operation has resulted in a carry out of the most significant bit.
	CarryFlag StatusFlag = 1 << iota // Bit 1

	// ZeroFlag indicates whether the result of an operation is zero.
	ZeroFlag = 1 << iota // Bit 2

	// InterruptDisableFlag disables interrupts when set.
	InterruptDisableFlag = 1 << iota // Bit 3

	// DecimalFlag enables Binary-Coded Decimal (BCD) mode when set.
	DecimalFlag = 1 << iota // Bit 4

	// BreakFlag indicates a BRK instruction has been executed.
	BreakFlag = 1 << iota // Bit 5

	// UnusedFlag is an unused flag in the 6502 status register.
	UnusedFlag = 1 << iota // Bit 6

	// OverflowFlag indicates whether an arithmetic operation has resulted in an overflow.
	OverflowFlag = 1 << iota // Bit 7

	// NegativeFlag indicates whether the result of an operation is negative.
	NegativeFlag = 1 << iota // Bit 8
)

// Registers holds the CPU registers including the accumulator (A), index registers (X and Y), stack pointer (S), program counter (PC), and status flags.
type Registers struct {
	A, X, Y, S byte
	PC         uint16
	Status     byte
}

// NewRegisters creates a new Registers instance.
func NewRegisters() *Registers {
	return &Registers{}
}

// SetStatus sets or clears the specified status flag.
func (r *Registers) SetStatus(f StatusFlag, on bool) {
	if on {
		r.Status = r.Status | byte(f)
	} else {
		r.Status = r.Status &^ byte(f)
	}
}

// IsSet checks if the specified status flag is set.
func (r *Registers) IsSet(f StatusFlag) bool {
	return (r.Status & byte(f)) != 0
}

// SetCarryFlag sets the carry flag based on the comparison of oldValue and result.
func (r *Registers) SetCarryFlag(oldValue, result byte) {
	r.SetStatus(CarryFlag, oldValue > result)
}

// SetZeroFlag sets the zero flag if the result is zero.
func (r *Registers) SetZeroFlag(result byte) {
	r.SetStatus(ZeroFlag, result == 0)
}

// SetNegativeFlag sets the negative flag if the result's most significant bit is set.
func (r *Registers) SetNegativeFlag(result byte) {
	r.SetStatus(NegativeFlag, (result&0x80) != 0)
}

// SetOverflowFlag sets the overflow flag based on the result of an addition or subtraction operation.
func (r *Registers) SetOverflowFlag(m, n, result byte, isAddition bool) {
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

	if isAddition {
		r.SetStatus(OverflowFlag, ((m^result)&(n^result))&0x80 != 0)
	} else {
		r.SetStatus(OverflowFlag, ((m^result)&(m^n))&0x80 != 0)
	}
}
