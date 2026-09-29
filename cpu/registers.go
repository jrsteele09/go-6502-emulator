package cpu

// StatusFlag represents the various status flags of the 6502 CPU.
type StatusFlag byte

const (
	// CarryFlag indicates whether an arithmetic operation has resulted in a carry out of the most significant bit.
	CarryFlag StatusFlag = 1 << iota // Bit 0

	// ZeroFlag indicates whether the result of an operation is zero.
	ZeroFlag = 1 << iota // Bit 1

	// InterruptDisableFlag disables interrupts when set.
	InterruptDisableFlag = 1 << iota // Bit 2

	// DecimalFlag enables Binary-Coded Decimal (BCD) mode when set.
	DecimalFlag = 1 << iota // Bit 3

	// BreakFlag is the B marker synthesized in status bytes pushed by BRK and PHP.
	// It is not stored as a live processor flag.
	BreakFlag = 1 << iota // Bit 4

	// UnusedFlag is the unused bit synthesized in pushed status bytes. The
	// register representation keeps it set for conventional status displays.
	UnusedFlag = 1 << iota // Bit 5

	// OverflowFlag indicates whether an arithmetic operation has resulted in an overflow.
	OverflowFlag = 1 << iota // Bit 6

	// NegativeFlag indicates whether the result of an operation is negative.
	NegativeFlag = 1 << iota // Bit 7
)

// Registers holds the CPU registers including the accumulator (A), index registers (X and Y), stack pointer (S), program counter (PC), and status flags.
type Registers struct {
	A, X, Y, S byte
	PC         uint16
	Status     byte
}

type RegisterEnum byte

const (
	AReg RegisterEnum = iota
	XReg
	YReg
	SReg
	PCReg
	StatusReg
)

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

// StatusForStack returns the status byte written by PHP, BRK, IRQ, or NMI.
// B is synthesized for software pushes and clear for hardware interrupts.
func (r *Registers) StatusForStack(breakMarker bool) byte {
	status := (r.Status &^ byte(BreakFlag)) | byte(UnusedFlag)
	if breakMarker {
		status |= byte(BreakFlag)
	}
	return status
}

// RestoreStatus restores the six stored flags from a stacked status byte.
// B is ignored because it is not a live flag; the conventional unused bit is
// kept set in the register representation.
func (r *Registers) RestoreStatus(status byte) {
	r.Status = (status &^ byte(BreakFlag)) | byte(UnusedFlag)
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
