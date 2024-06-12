package cpu

const (
	ByteAddressing = "nn"
	WordAddressing = "nnnn"
)

type LoadAddress func() (byte, Completed)
type StoreAddress func(b byte) Completed

type AddressingModeType string

type AddressingMode interface {
	Store(cpu Cpu6502, ignoreExtraCycle bool) StoreAddress
	Load(cpu Cpu6502, ignoreExtraCycle bool) LoadAddress
	Address(cpu Cpu6502) uint16
}

const (
	ImpliedModeStr          AddressingModeType = ""
	AbsoluteIndirectModeStr AddressingModeType = "(nnnn)"
	AbsoluteModeStr         AddressingModeType = "nnnn"
	AbsoluteIndexedXModeStr AddressingModeType = "nnnn,X"
	AbsoluteIndexedYModeStr AddressingModeType = "nnnn,Y"
	AccumulatorModeStr      AddressingModeType = "A"
	ZeropageModeStr         AddressingModeType = "nn"
	ZeropageXModeStr        AddressingModeType = "nn,X"
	ZeropageYModeStr        AddressingModeType = "nn,Y"
	IndexedIndirectModeStr  AddressingModeType = "(nn,X)"
	IndirectIndexedModeStr  AddressingModeType = "(nn),Y"
	ImmediateModeStr        AddressingModeType = "#nn"
	RelativeModeStr         AddressingModeType = "+nn"
)

type AbsoluteIndirectMode struct{}
type AbsoluteMode struct{}
type AbsoluteXMode struct{}
type AbsoluteYMode struct{}
type AccumulatorMode struct{}
type ZeropageMode struct{}
type ZeropageXMode struct{}
type ZeropageYMode struct{}
type IndexedIndirectMode struct{}
type IndirectIndexedMode struct{}
type ImmediateMode struct{}
type RelativeMode struct{}

func absoluteAddress(cpu Cpu6502) uint16 {
	operands := cpu.Operands()
	memAddress := uint16(operands[0])
	memAddress |= uint16(operands[1]) << 8
	return memAddress
}

func absoluteXAddress(cpu Cpu6502, ignoreExtraCycle bool) (uint16, bool) {
	extraCycle := false
	operands := cpu.Operands()
	lsb := uint16(operands[0]) + uint16(cpu.Registers().X)
	address := (uint16(operands[1]) << 8) + lsb
	if !ignoreExtraCycle && lsb > 0xFF {
		extraCycle = true
	}
	return uint16(address), extraCycle
}

func absoluteYAddress(cpu Cpu6502, ignoreExtraCycle bool) (uint16, bool) {
	extraCycle := false
	operands := cpu.Operands()
	lsb := uint16(operands[0]) + uint16(cpu.Registers().Y)
	address := (uint16(operands[1]) << 8) + lsb
	if !ignoreExtraCycle && lsb > 0xFF {
		extraCycle = true
	}
	return uint16(address), extraCycle
}

func zeropageXAddress(cpu Cpu6502) uint16 {
	return (uint16(cpu.Operands()[0]) + uint16(cpu.Registers().X)) & 0xFF
}

func zeropageYAddress(cpu Cpu6502) uint16 {
	return (uint16(cpu.Operands()[0]) + uint16(cpu.Registers().Y)) & 0xFF
}

func indexedIndirectAddress(cpu Cpu6502) uint16 {
	mem := cpu.Memory()
	operands := cpu.Operands()

	zeropageAddress := uint16(operands[0] + cpu.Registers().X)
	lsb := (mem.Read(zeropageAddress))
	msb := (mem.Read(zeropageAddress + 1))
	return ((uint16(msb) << 8) | uint16(lsb))
}

func indirectIndexedAddress(cpu Cpu6502, ignoreExtraCycle bool) (uint16, bool) {
	mem := cpu.Memory()
	zeropageAddress := uint16(cpu.Operands()[0])
	lsb := mem.Read(zeropageAddress)
	msb := mem.Read(zeropageAddress + 1)
	newLsb := lsb + cpu.Registers().Y
	extraCycle := false
	if !ignoreExtraCycle && newLsb < lsb {
		msb += 1
		extraCycle = true
	}

	address := (uint16(msb) << 8) + uint16(newLsb)
	return address, extraCycle
}

// AbsoluteIndirectMode
func (m AbsoluteIndirectMode) Store(cpu Cpu6502, ignoreExtraCycle bool) StoreAddress {
	return func(b byte) Completed { return true }
}

func (m AbsoluteIndirectMode) Load(cpu Cpu6502, ignoreExtraCycle bool) LoadAddress {
	return func() (byte, Completed) { return 0x00, true }
}

func (m AbsoluteIndirectMode) Address(cpu Cpu6502) uint16 {
	absoluteAddress := absoluteAddress(cpu)
	mem := cpu.Memory()
	lsb := mem.Read(absoluteAddress)
	msb := mem.Read(absoluteAddress + 1)
	return (uint16(msb) << 8) + uint16(lsb)
}

// AbsoluteMode
func (m AbsoluteMode) Store(cpu Cpu6502, ignoreExtraCycle bool) StoreAddress {
	mem := cpu.Memory()
	return func(b byte) Completed {
		mem.Write(absoluteAddress(cpu), b)
		return true
	}
}

func (m AbsoluteMode) Load(cpu Cpu6502, ignoreExtraCycle bool) LoadAddress {
	mem := cpu.Memory()
	return func() (byte, Completed) { // Load
		return mem.Read(absoluteAddress(cpu)), true
	}
}

func (m AbsoluteMode) Address(cpu Cpu6502) uint16 {
	return absoluteAddress(cpu)
}

// AbsoluteXMode
func (m AbsoluteXMode) Store(cpu Cpu6502, ignoreExtraCycle bool) StoreAddress {
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

func (m AbsoluteXMode) Load(cpu Cpu6502, ignoreExtraCycle bool) LoadAddress {
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

func (m AbsoluteXMode) Address(cpu Cpu6502) uint16 {
	address, _ := absoluteXAddress(cpu, true)
	return address
}

// AbsoluteYMode
func (m AbsoluteYMode) Store(cpu Cpu6502, ignoreExtraCycle bool) StoreAddress {
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

func (m AbsoluteYMode) Load(cpu Cpu6502, ignoreExtraCycle bool) LoadAddress {
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

func (m AbsoluteYMode) Address(cpu Cpu6502) uint16 {
	address, _ := absoluteYAddress(cpu, true)
	return address
}

// AccumulatorMode
func (m AccumulatorMode) Store(cpu Cpu6502, ignoreExtraCycle bool) StoreAddress {
	registers := cpu.Registers()

	return func(b byte) Completed {
		registers.A = b
		return true
	}
}

func (m AccumulatorMode) Load(cpu Cpu6502, ignoreExtraCycle bool) LoadAddress {
	return func() (byte, Completed) { return cpu.Registers().A, true }
}

func (m AccumulatorMode) Address(cpu Cpu6502) uint16 {
	return 0x0000
}

// Zeropage
func (m ZeropageMode) Store(cpu Cpu6502, ignoreExtraCycle bool) StoreAddress {
	mem := cpu.Memory()

	return func(b byte) Completed {
		mem.Write(uint16(cpu.Operands()[0]), b)
		return true
	}
}

func (m ZeropageMode) Load(cpu Cpu6502, ignoreExtraCycle bool) LoadAddress {
	mem := cpu.Memory()

	return func() (byte, Completed) {
		return mem.Read(uint16(cpu.Operands()[0])), true
	}
}

func (m ZeropageMode) Address(cpu Cpu6502) uint16 {
	return uint16(cpu.Operands()[0])
}

// ZeropageXMode
func (m ZeropageXMode) Store(cpu Cpu6502, ignoreExtraCycle bool) StoreAddress {
	mem := cpu.Memory()
	return func(b byte) Completed {
		mem.Write(zeropageXAddress(cpu), b)
		return true
	}
}

func (m ZeropageXMode) Load(cpu Cpu6502, ignoreExtraCycle bool) LoadAddress {
	mem := cpu.Memory()
	return func() (byte, Completed) {
		return mem.Read(zeropageXAddress(cpu)), true
	}
}

func (m ZeropageXMode) Address(cpu Cpu6502) uint16 {
	return zeropageXAddress(cpu)
}

// ZeropageYMode
func (m ZeropageYMode) Store(cpu Cpu6502, ignoreExtraCycle bool) StoreAddress {
	mem := cpu.Memory()
	return func(b byte) Completed {
		mem.Write(zeropageYAddress(cpu), b)
		return true
	}
}

func (m ZeropageYMode) Load(cpu Cpu6502, ignoreExtraCycle bool) LoadAddress {
	mem := cpu.Memory()
	return func() (byte, Completed) {
		return mem.Read(zeropageYAddress(cpu)), true
	}
}

func (m ZeropageYMode) Address(cpu Cpu6502) uint16 {
	return zeropageYAddress(cpu)
}

// IndexedIndirectMode
func (m IndexedIndirectMode) Load(cpu Cpu6502, ignoreExtraCycle bool) LoadAddress {
	mem := cpu.Memory()
	return func() (byte, Completed) {
		return mem.Read(indexedIndirectAddress(cpu)), true
	}
}

func (m IndexedIndirectMode) Store(cpu Cpu6502, ignoreExtraCycle bool) StoreAddress {
	mem := cpu.Memory()
	return func(b byte) Completed {
		mem.Write(indexedIndirectAddress(cpu), b)
		return true
	}
}

func (m IndexedIndirectMode) Address(cpu Cpu6502) uint16 {
	return indexedIndirectAddress(cpu)
}

// IndirectIndexedMode
func (m IndirectIndexedMode) Load(cpu Cpu6502, ignoreExtraCycle bool) LoadAddress {
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

func (m IndirectIndexedMode) Store(cpu Cpu6502, ignoreExtraCycle bool) StoreAddress {
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

func (m IndirectIndexedMode) Address(cpu Cpu6502) uint16 {
	address, _ := indirectIndexedAddress(cpu, true)
	return address
}

// ImmediateMode
func (m ImmediateMode) Load(cpu Cpu6502, ignoreExtraCycle bool) LoadAddress {
	return func() (byte, Completed) {
		operands := cpu.Operands()
		return operands[0], true
	}
}

func (m ImmediateMode) Store(cpu Cpu6502, ignoreExtraCycle bool) StoreAddress {
	return func(b byte) Completed { return true }
}

func (m ImmediateMode) Address(cpu Cpu6502) uint16 {
	return 0x000
}

// RelativeMode
func (m RelativeMode) Store(cpu Cpu6502, ignoreExtraCycle bool) StoreAddress {
	return func(b byte) Completed { return true }
}

func (m RelativeMode) Load(cpu Cpu6502, ignoreExtraCycle bool) LoadAddress {
	return func() (byte, Completed) {
		operands := cpu.Operands()
		return operands[0], true
	}
}

func (m RelativeMode) Address(cpu Cpu6502) uint16 {
	return 0x000
}

func GetAddressingMode(am AddressingModeType) AddressingMode {
	return map[AddressingModeType]AddressingMode{
		AbsoluteIndirectModeStr: AbsoluteIndirectMode{},
		AbsoluteModeStr:         AbsoluteMode{},
		AbsoluteIndexedXModeStr: AbsoluteXMode{},
		AbsoluteIndexedYModeStr: AbsoluteYMode{},
		AccumulatorModeStr:      AccumulatorMode{},
		ZeropageModeStr:         ZeropageMode{},
		ZeropageXModeStr:        ZeropageXMode{},
		ZeropageYModeStr:        ZeropageYMode{},
		IndexedIndirectModeStr:  IndexedIndirectMode{},
		IndirectIndexedModeStr:  IndirectIndexedMode{},
		ImmediateModeStr:        ImmediateMode{},
		RelativeModeStr:         RelativeMode{},
	}[am]
}
