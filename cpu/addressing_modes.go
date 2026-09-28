package cpu

const (
	// ByteAddressing defines a two-character byte addressing mode representation.
	ByteAddressing = "nn"

	// WordAddressing defines a four-character word addressing mode representation.
	WordAddressing = "nnnn"
)

// AddressingModeType represents a string type for addressing mode representation.
type AddressingModeType string

// AddressingMode defines an interface for various addressing modes in the 6502 CPU.
// Each addressing mode can store and load values and calculate the effective address.
type AddressingMode interface {
	// Store writes a byte to the calculated address.
	// The ignoreExtraCycle parameter determines whether to ignore the extra cycle typically required for certain addressing modes.
	Store(cpu *CPU, b byte, ignoreExtraCycle bool) Completed

	// Load reads a byte from the calculated address.
	// The ignoreExtraCycle parameter determines whether to ignore the extra cycle typically required for certain addressing modes.
	Load(cpu *CPU, ignoreExtraCycle bool) (byte, Completed)

	// Address calculates and returns the effective address for the addressing mode.
	Address(cpu *CPU) uint16
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

func absoluteAddress(cpu *CPU) uint16 {
	operands := cpu.execute.operands
	memAddress := uint16(operands[0])
	memAddress |= uint16(operands[1]) << 8
	return memAddress
}

func absoluteXAddress(cpu *CPU, ignoreExtraCycle bool) (uint16, bool) {
	extraCycle := false
	operands := cpu.execute.operands
	lsb := uint16(operands[0]) + uint16(cpu.Reg.X)
	address := (uint16(operands[1]) << 8) + lsb
	if !ignoreExtraCycle && lsb > 0xFF {
		extraCycle = true
	}
	return uint16(address), extraCycle
}

func absoluteYAddress(cpu *CPU, ignoreExtraCycle bool) (uint16, bool) {
	extraCycle := false
	operands := cpu.execute.operands
	lsb := uint16(operands[0]) + uint16(cpu.Reg.Y)
	address := (uint16(operands[1]) << 8) + lsb
	if !ignoreExtraCycle && lsb > 0xFF {
		extraCycle = true
	}
	return uint16(address), extraCycle
}

func zeropageXAddress(cpu *CPU) uint16 {
	return (uint16(cpu.execute.operands[0]) + uint16(cpu.Reg.X)) & 0xFF
}

func zeropageYAddress(cpu *CPU) uint16 {
	return (uint16(cpu.execute.operands[0]) + uint16(cpu.Reg.Y)) & 0xFF
}

func indexedIndirectAddress(cpu *CPU) uint16 {
	mem := cpu.mem
	operands := cpu.execute.operands

	zeropageAddress := uint16(operands[0] + cpu.Reg.X)
	lsb := (mem.Read(zeropageAddress))
	msb := (mem.Read((zeropageAddress + 1) & 0x00FF))
	return ((uint16(msb) << 8) | uint16(lsb))
}

func indirectIndexedAddress(cpu *CPU, ignoreExtraCycle bool) (uint16, bool) {
	mem := cpu.mem
	zeropageAddress := uint16(cpu.execute.operands[0])
	lsb := mem.Read(zeropageAddress)
	msb := mem.Read((zeropageAddress + 1) & 0x00FF)
	baseAddress := (uint16(msb) << 8) | uint16(lsb)
	address := baseAddress + uint16(cpu.Reg.Y)
	extraCycle := !ignoreExtraCycle && (baseAddress&0xFF00) != (address&0xFF00)
	return address, extraCycle
}

func (m absoluteIndirectMode) Store(_ *CPU, _ byte, _ bool) Completed {
	return true
}

func (m absoluteIndirectMode) Load(_ *CPU, _ bool) (byte, Completed) {
	return 0x00, true
}

func (m absoluteIndirectMode) Address(cpu *CPU) uint16 {
	absoluteAddress := absoluteAddress(cpu)
	mem := cpu.mem
	lsb := mem.Read(absoluteAddress)
	msbAddress := (absoluteAddress & 0xFF00) | ((absoluteAddress + 1) & 0x00FF)
	msb := mem.Read(msbAddress)
	return (uint16(msb) << 8) + uint16(lsb)
}

// AbsoluteMode
func (m absoluteMode) Store(cpu *CPU, b byte, _ bool) Completed {
	cpu.mem.Write(absoluteAddress(cpu), b)
	return true
}

func (m absoluteMode) Load(cpu *CPU, _ bool) (byte, Completed) {
	return cpu.mem.Read(absoluteAddress(cpu)), true
}

func (m absoluteMode) Address(cpu *CPU) uint16 {
	return absoluteAddress(cpu)
}

// AbsoluteXMode
func (m absoluteXMode) Store(cpu *CPU, b byte, ignoreExtraCycle bool) Completed {
	if cpu.execute.phase != 0 {
		cpu.execute.phase = 0
		cpu.mem.Write(cpu.execute.address, b)
		return true
	}
	address, extraCycle := absoluteXAddress(cpu, ignoreExtraCycle)
	if extraCycle {
		cpu.execute.address = address
		cpu.execute.phase = 1
		return false
	}
	cpu.mem.Write(address, b)
	return true
}

func (m absoluteXMode) Load(cpu *CPU, ignoreExtraCycle bool) (byte, Completed) {
	if cpu.execute.phase != 0 {
		cpu.execute.phase = 0
		return cpu.mem.Read(cpu.execute.address), true
	}
	address, extraCycle := absoluteXAddress(cpu, ignoreExtraCycle)
	if extraCycle {
		cpu.execute.address = address
		cpu.execute.phase = 1
		return 0x00, false
	}
	return cpu.mem.Read(address), true
}

func (m absoluteXMode) Address(cpu *CPU) uint16 {
	address, _ := absoluteXAddress(cpu, true)
	return address
}

// AbsoluteYMode
func (m absoluteYMode) Store(cpu *CPU, b byte, ignoreExtraCycle bool) Completed {
	if cpu.execute.phase != 0 {
		cpu.execute.phase = 0
		cpu.mem.Write(cpu.execute.address, b)
		return true
	}
	address, extraCycle := absoluteYAddress(cpu, ignoreExtraCycle)
	if extraCycle {
		cpu.execute.address = address
		cpu.execute.phase = 1
		return false
	}
	cpu.mem.Write(address, b)
	return true
}

func (m absoluteYMode) Load(cpu *CPU, ignoreExtraCycle bool) (byte, Completed) {
	if cpu.execute.phase != 0 {
		cpu.execute.phase = 0
		return cpu.mem.Read(cpu.execute.address), true
	}
	address, extraCycle := absoluteYAddress(cpu, ignoreExtraCycle)
	if extraCycle {
		cpu.execute.address = address
		cpu.execute.phase = 1
		return 0x00, false
	}
	return cpu.mem.Read(address), true
}

func (m absoluteYMode) Address(cpu *CPU) uint16 {
	address, _ := absoluteYAddress(cpu, true)
	return address
}

// AccumulatorMode
func (m accumulatorMode) Store(cpu *CPU, b byte, _ bool) Completed {
	cpu.Reg.A = b
	return true
}

func (m accumulatorMode) Load(cpu *CPU, _ bool) (byte, Completed) {
	return cpu.Reg.A, true
}

func (m accumulatorMode) Address(_ *CPU) uint16 {
	return 0x0000
}

// Zeropage
func (m zeropageMode) Store(cpu *CPU, b byte, _ bool) Completed {
	cpu.mem.Write(uint16(cpu.execute.operands[0]), b)
	return true
}

func (m zeropageMode) Load(cpu *CPU, _ bool) (byte, Completed) {
	return cpu.mem.Read(uint16(cpu.execute.operands[0])), true
}

func (m zeropageMode) Address(cpu *CPU) uint16 {
	return uint16(cpu.execute.operands[0])
}

// ZeropageXMode
func (m zeropageXMode) Store(cpu *CPU, b byte, _ bool) Completed {
	cpu.mem.Write(zeropageXAddress(cpu), b)
	return true
}

func (m zeropageXMode) Load(cpu *CPU, _ bool) (byte, Completed) {
	return cpu.mem.Read(zeropageXAddress(cpu)), true
}

func (m zeropageXMode) Address(cpu *CPU) uint16 {
	return zeropageXAddress(cpu)
}

// ZeropageYMode
func (m zeropageYMode) Store(cpu *CPU, b byte, _ bool) Completed {
	cpu.mem.Write(zeropageYAddress(cpu), b)
	return true
}

func (m zeropageYMode) Load(cpu *CPU, _ bool) (byte, Completed) {
	return cpu.mem.Read(zeropageYAddress(cpu)), true
}

func (m zeropageYMode) Address(cpu *CPU) uint16 {
	return zeropageYAddress(cpu)
}

// IndexedIndirectMode
func (m indexedIndirectMode) Load(cpu *CPU, _ bool) (byte, Completed) {
	return cpu.mem.Read(indexedIndirectAddress(cpu)), true
}

func (m indexedIndirectMode) Store(cpu *CPU, b byte, _ bool) Completed {
	cpu.mem.Write(indexedIndirectAddress(cpu), b)
	return true
}

func (m indexedIndirectMode) Address(cpu *CPU) uint16 {
	return indexedIndirectAddress(cpu)
}

// IndirectIndexedMode
func (m indirectIndexedMode) Load(cpu *CPU, ignoreExtraCycle bool) (byte, Completed) {
	if cpu.execute.phase != 0 {
		cpu.execute.phase = 0
		return cpu.mem.Read(cpu.execute.address), true
	}
	address, extraCycle := indirectIndexedAddress(cpu, ignoreExtraCycle)
	if extraCycle {
		cpu.execute.address = address
		cpu.execute.phase = 1
		return 0x00, false
	}
	return cpu.mem.Read(address), true
}

func (m indirectIndexedMode) Store(cpu *CPU, b byte, ignoreExtraCycle bool) Completed {
	if cpu.execute.phase != 0 {
		cpu.execute.phase = 0
		cpu.mem.Write(cpu.execute.address, b)
		return true
	}
	address, extraCycle := indirectIndexedAddress(cpu, ignoreExtraCycle)
	if extraCycle {
		cpu.execute.address = address
		cpu.execute.phase = 1
		return false
	}
	cpu.mem.Write(address, b)
	return true
}

func (m indirectIndexedMode) Address(cpu *CPU) uint16 {
	address, _ := indirectIndexedAddress(cpu, true)
	return address
}

// ImmediateMode
func (m immediateMode) Load(cpu *CPU, _ bool) (byte, Completed) {
	return cpu.execute.operands[0], true
}

func (m immediateMode) Store(_ *CPU, _ byte, _ bool) Completed {
	return true
}

func (m immediateMode) Address(_ *CPU) uint16 {
	return 0x000
}

// RelativeMode
func (m relativeMode) Store(_ *CPU, _ byte, _ bool) Completed {
	return true
}

func (m relativeMode) Load(cpu *CPU, _ bool) (byte, Completed) {
	return cpu.execute.operands[0], true
}

func (m relativeMode) Address(_ *CPU) uint16 {
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
