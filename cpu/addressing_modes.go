package cpu

const (
	// ByteAddressing defines a two-character byte addressing mode representation.
	ByteAddressing = "nn"

	// WordAddressing defines a four-character word addressing mode representation.
	WordAddressing = "nnnn"
)

// LoadAddress defines a function type that loads a byte and indicates if the operation is completed.
type LoadAddress func() (byte, Completed)

// StoreAddress defines a function type that stores a byte and indicates if the operation is completed.
type StoreAddress func(b byte) Completed

// AddressingModeType represents a string type for addressing mode representation.
type AddressingModeType string

// AddressingMode defines an interface for various addressing modes in the 6502 CPU.
// Each addressing mode can store and load values and calculate the effective address.
type AddressingMode interface {
	// Store returns a function that stores a byte at the calculated address.
	// The ignoreExtraCycle parameter determines whether to ignore the extra cycle typically required for certain addressing modes.
	Store(cpu CPU6502, ignoreExtraCycle bool) StoreAddress

	// Load returns a function that loads a byte from the calculated address.
	// The ignoreExtraCycle parameter determines whether to ignore the extra cycle typically required for certain addressing modes.
	Load(cpu CPU6502, ignoreExtraCycle bool) LoadAddress

	// Address calculates and returns the effective address for the addressing mode.
	Address(cpu CPU6502) uint16
}

const (
	// ImpliedModeStr is the implied addressing mode.
	ImpliedModeStr AddressingModeType = ""
	// AbsoluteIndirectModeStr is the absolute indirect addressing mode.
	AbsoluteIndirectModeStr AddressingModeType = "(nnnn)"
	// AbsoluteModeStr is the absolute addressing mode.
	AbsoluteModeStr AddressingModeType = "nnnn"
	// AbsoluteIndexedXModeStr is the absolute indexed X addressing mode.
	AbsoluteIndexedXModeStr AddressingModeType = "nnnn,X"
	// AbsoluteIndexedYModeStr is the absolute indexed Y addressing mode.
	AbsoluteIndexedYModeStr AddressingModeType = "nnnn,Y"
	// AccumulatorModeStr is the accumulator addressing mode.
	AccumulatorModeStr AddressingModeType = "A"
	// ZeropageModeStr is the zeropage addressing mode.
	ZeropageModeStr AddressingModeType = "nn"
	// ZeropageXModeStr is the zeropage indexed X addressing mode.
	ZeropageXModeStr AddressingModeType = "nn,X"
	// ZeropageYModeStr is the zeropage indexed Y addressing mode.
	ZeropageYModeStr AddressingModeType = "nn,Y"
	// IndexedIndirectModeStr is the indexed indirect addressing mode.
	IndexedIndirectModeStr AddressingModeType = "(nn,X)"
	// IndirectIndexedModeStr is the indirect indexed addressing mode.
	IndirectIndexedModeStr AddressingModeType = "(nn),Y"
	// ImmediateModeStr is the immediate addressing mode.
	ImmediateModeStr AddressingModeType = "#nn"
	// RelativeModeStr is the relative addressing mode.
	RelativeModeStr AddressingModeType = "*nn"
)

type absoluteIndirectMode struct{}
type absoluteMode struct{}
type absoluteXMode struct{}
type absoluteYMode struct{}
type accumulatorMode struct{}
type zeropageMode struct{}
type zeropageXMode struct{}
type zeropageYMode struct{}
type indexedIndirectMode struct{}
type indirectIndexedMode struct{}
type immediateMode struct{}
type relativeMode struct{}

func absoluteAddress(cpu CPU6502) uint16 {
	operands := cpu.Operands()
	memAddress := uint16(operands[0])
	memAddress |= uint16(operands[1]) << 8
	return memAddress
}

func absoluteXAddress(cpu CPU6502, ignoreExtraCycle bool) (uint16, bool) {
	extraCycle := false
	operands := cpu.Operands()
	lsb := uint16(operands[0]) + uint16(cpu.Registers().X)
	address := (uint16(operands[1]) << 8) + lsb
	if !ignoreExtraCycle && lsb > 0xFF {
		extraCycle = true
	}
	return uint16(address), extraCycle
}

func absoluteYAddress(cpu CPU6502, ignoreExtraCycle bool) (uint16, bool) {
	extraCycle := false
	operands := cpu.Operands()
	lsb := uint16(operands[0]) + uint16(cpu.Registers().Y)
	address := (uint16(operands[1]) << 8) + lsb
	if !ignoreExtraCycle && lsb > 0xFF {
		extraCycle = true
	}
	return uint16(address), extraCycle
}

func zeropageXAddress(cpu CPU6502) uint16 {
	return (uint16(cpu.Operands()[0]) + uint16(cpu.Registers().X)) & 0xFF
}

func zeropageYAddress(cpu CPU6502) uint16 {
	return (uint16(cpu.Operands()[0]) + uint16(cpu.Registers().Y)) & 0xFF
}

func indexedIndirectAddress(cpu CPU6502) uint16 {
	mem := cpu.Memory()
	operands := cpu.Operands()

	zeropageAddress := uint16(operands[0] + cpu.Registers().X)
	lsb := (mem.Read(zeropageAddress))
	msb := (mem.Read(zeropageAddress + 1))
	return ((uint16(msb) << 8) | uint16(lsb))
}

func indirectIndexedAddress(cpu CPU6502, ignoreExtraCycle bool) (uint16, bool) {
	mem := cpu.Memory()
	zeropageAddress := uint16(cpu.Operands()[0])
	lsb := mem.Read(zeropageAddress)
	msb := mem.Read(zeropageAddress + 1)
	newLsb := lsb + cpu.Registers().Y
	extraCycle := false
	if !ignoreExtraCycle && newLsb < lsb {
		msb++
		extraCycle = true
	}

	address := (uint16(msb) << 8) + uint16(newLsb)
	return address, extraCycle
}

func (m absoluteIndirectMode) Store(_ CPU6502, _ bool) StoreAddress {
	return func(_ byte) Completed { return true }
}

func (m absoluteIndirectMode) Load(_ CPU6502, _ bool) LoadAddress {
	return func() (byte, Completed) { return 0x00, true }
}

func (m absoluteIndirectMode) Address(cpu CPU6502) uint16 {
	absoluteAddress := absoluteAddress(cpu)
	mem := cpu.Memory()
	lsb := mem.Read(absoluteAddress)
	msb := mem.Read(absoluteAddress + 1)
	return (uint16(msb) << 8) + uint16(lsb)
}

// AbsoluteMode
func (m absoluteMode) Store(cpu CPU6502, _ bool) StoreAddress {
	mem := cpu.Memory()
	return func(b byte) Completed {
		mem.Write(absoluteAddress(cpu), b)
		return true
	}
}

func (m absoluteMode) Load(cpu CPU6502, _ bool) LoadAddress {
	mem := cpu.Memory()
	return func() (byte, Completed) { // Load
		return mem.Read(absoluteAddress(cpu)), true
	}
}

func (m absoluteMode) Address(cpu CPU6502) uint16 {
	return absoluteAddress(cpu)
}

// AbsoluteXMode
func (m absoluteXMode) Store(cpu CPU6502, ignoreExtraCycle bool) StoreAddress {
	address := uint16(0x0000)
	extraCycle := false
	mem := cpu.Memory()

	return func(b byte) Completed {
		if extraCycle {
			mem.Write(address, b)
			return true
		}
		if address, extraCycle = absoluteXAddress(cpu, ignoreExtraCycle); extraCycle {
			return false
		}
		mem.Write(address, b)
		return true
	}
}

func (m absoluteXMode) Load(cpu CPU6502, ignoreExtraCycle bool) LoadAddress {
	address := uint16(0x0000)
	extraCycle := false
	result := byte(0x00)
	mem := cpu.Memory()

	return func() (byte, Completed) {
		if extraCycle {
			return mem.Read(address), true
		}

		if address, extraCycle = absoluteXAddress(cpu, ignoreExtraCycle); extraCycle {
			return 0x00, false
		}
		result = mem.Read(address)
		return result, true
	}
}

func (m absoluteXMode) Address(cpu CPU6502) uint16 {
	address, _ := absoluteXAddress(cpu, true)
	return address
}

// AbsoluteYMode
func (m absoluteYMode) Store(cpu CPU6502, ignoreExtraCycle bool) StoreAddress {
	address := uint16(0x0000)
	extraCycle := false
	mem := cpu.Memory()

	return func(b byte) Completed {
		if extraCycle {
			mem.Write(address, b)
			return true
		}

		if address, extraCycle = absoluteYAddress(cpu, ignoreExtraCycle); extraCycle {
			return false
		}

		mem.Write(address, b)
		return true
	}
}

func (m absoluteYMode) Load(cpu CPU6502, ignoreExtraCycle bool) LoadAddress {
	address := uint16(0x0000)
	extraCycle := false
	mem := cpu.Memory()

	return func() (byte, Completed) {
		if extraCycle {
			return mem.Read(address), true
		}
		if address, extraCycle = absoluteYAddress(cpu, ignoreExtraCycle); extraCycle {
			return 0x00, false
		}
		return mem.Read(address), true
	}
}

func (m absoluteYMode) Address(cpu CPU6502) uint16 {
	address, _ := absoluteYAddress(cpu, true)
	return address
}

// AccumulatorMode
func (m accumulatorMode) Store(cpu CPU6502, _ bool) StoreAddress {
	registers := cpu.Registers()

	return func(b byte) Completed {
		registers.A = b
		return true
	}
}

func (m accumulatorMode) Load(cpu CPU6502, _ bool) LoadAddress {
	return func() (byte, Completed) { return cpu.Registers().A, true }
}

func (m accumulatorMode) Address(_ CPU6502) uint16 {
	return 0x0000
}

// Zeropage
func (m zeropageMode) Store(cpu CPU6502, _ bool) StoreAddress {
	mem := cpu.Memory()

	return func(b byte) Completed {
		mem.Write(uint16(cpu.Operands()[0]), b)
		return true
	}
}

func (m zeropageMode) Load(cpu CPU6502, _ bool) LoadAddress {
	mem := cpu.Memory()

	return func() (byte, Completed) {
		return mem.Read(uint16(cpu.Operands()[0])), true
	}
}

func (m zeropageMode) Address(cpu CPU6502) uint16 {
	return uint16(cpu.Operands()[0])
}

// ZeropageXMode
func (m zeropageXMode) Store(cpu CPU6502, _ bool) StoreAddress {
	mem := cpu.Memory()
	return func(b byte) Completed {
		mem.Write(zeropageXAddress(cpu), b)
		return true
	}
}

func (m zeropageXMode) Load(cpu CPU6502, _ bool) LoadAddress {
	mem := cpu.Memory()
	return func() (byte, Completed) {
		return mem.Read(zeropageXAddress(cpu)), true
	}
}

func (m zeropageXMode) Address(cpu CPU6502) uint16 {
	return zeropageXAddress(cpu)
}

// ZeropageYMode
func (m zeropageYMode) Store(cpu CPU6502, _ bool) StoreAddress {
	mem := cpu.Memory()
	return func(b byte) Completed {
		mem.Write(zeropageYAddress(cpu), b)
		return true
	}
}

func (m zeropageYMode) Load(cpu CPU6502, _ bool) LoadAddress {
	mem := cpu.Memory()
	return func() (byte, Completed) {
		return mem.Read(zeropageYAddress(cpu)), true
	}
}

func (m zeropageYMode) Address(cpu CPU6502) uint16 {
	return zeropageYAddress(cpu)
}

// IndexedIndirectMode
func (m indexedIndirectMode) Load(cpu CPU6502, _ bool) LoadAddress {
	mem := cpu.Memory()
	return func() (byte, Completed) {
		return mem.Read(indexedIndirectAddress(cpu)), true
	}
}

func (m indexedIndirectMode) Store(cpu CPU6502, _ bool) StoreAddress {
	mem := cpu.Memory()
	return func(b byte) Completed {
		mem.Write(indexedIndirectAddress(cpu), b)
		return true
	}
}

func (m indexedIndirectMode) Address(cpu CPU6502) uint16 {
	return indexedIndirectAddress(cpu)
}

// IndirectIndexedMode
func (m indirectIndexedMode) Load(cpu CPU6502, ignoreExtraCycle bool) LoadAddress {
	extraCycle := false
	address := uint16(0x0000)
	mem := cpu.Memory()

	return func() (byte, Completed) {
		if extraCycle {
			return mem.Read(address), true
		}
		if address, extraCycle = indirectIndexedAddress(cpu, ignoreExtraCycle); extraCycle {
			return 0x00, false
		}
		return mem.Read(address), true
	}
}

func (m indirectIndexedMode) Store(cpu CPU6502, ignoreExtraCycle bool) StoreAddress {
	extraCycle := false
	address := uint16(0x0000)
	mem := cpu.Memory()

	return func(b byte) Completed {
		if extraCycle {
			mem.Write(address, b)
			return true
		}
		if address, extraCycle = indirectIndexedAddress(cpu, ignoreExtraCycle); extraCycle {
			return false
		}
		mem.Write(address, b)
		return true
	}
}

func (m indirectIndexedMode) Address(cpu CPU6502) uint16 {
	address, _ := indirectIndexedAddress(cpu, true)
	return address
}

// ImmediateMode
func (m immediateMode) Load(cpu CPU6502, _ bool) LoadAddress {
	return func() (byte, Completed) {
		operands := cpu.Operands()
		return operands[0], true
	}
}

func (m immediateMode) Store(_ CPU6502, _ bool) StoreAddress {
	return func(_ byte) Completed { return true }
}

func (m immediateMode) Address(_ CPU6502) uint16 {
	return 0x000
}

// RelativeMode
func (m relativeMode) Store(_ CPU6502, _ bool) StoreAddress {
	return func(_ byte) Completed { return true }
}

func (m relativeMode) Load(cpu CPU6502, _ bool) LoadAddress {
	return func() (byte, Completed) {
		operands := cpu.Operands()
		return operands[0], true
	}
}

func (m relativeMode) Address(_ CPU6502) uint16 {
	return 0x000
}

func getAddressingMode(am AddressingModeType) AddressingMode {
	return map[AddressingModeType]AddressingMode{
		AbsoluteIndirectModeStr: absoluteIndirectMode{},
		AbsoluteModeStr:         absoluteMode{},
		AbsoluteIndexedXModeStr: absoluteXMode{},
		AbsoluteIndexedYModeStr: absoluteYMode{},
		AccumulatorModeStr:      accumulatorMode{},
		ZeropageModeStr:         zeropageMode{},
		ZeropageXModeStr:        zeropageXMode{},
		ZeropageYModeStr:        zeropageYMode{},
		IndexedIndirectModeStr:  indexedIndirectMode{},
		IndirectIndexedModeStr:  indirectIndexedMode{},
		ImmediateModeStr:        immediateMode{},
		RelativeModeStr:         relativeMode{},
	}[am]
}
