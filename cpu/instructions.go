package cpu

const (
	ADC = "ADC"
	AND = "AND"
	ASL = "ASL"
	BCC = "BCC"
	BCS = "BCS"
	BEQ = "BEQ"
	BIT = "BIT"
	BMI = "BMI"
	BNE = "BNE"
	BPL = "BPL"
	BRK = "BRK"
	BVC = "BVC"
	BVS = "BVS"
	CLC = "CLC"
	CLD = "CLD"
	CLI = "CLI"
	CLV = "CLV"
	CMP = "CMP"
	CPX = "CPX"
	CPY = "CPY"
	DEC = "DEC"
	DEX = "DEX"
	DEY = "DEY"
	EOR = "EOR"
	INC = "INC"
	INX = "INX"
	INY = "INY"
	JMP = "JMP"
	JSR = "JSR"
	LDA = "LDA"
	LDX = "LDX"
	LDY = "LDY"
	LSR = "LSR"
	NOP = "NOP"
	ORA = "ORA"
	PHA = "PHA"
	PHP = "PHP"
	PLA = "PLA"
	PLP = "PLP"
	ROL = "ROL"
	ROR = "ROR"
	RTI = "RTI"
	RTS = "RTS"
	SBC = "SBC"
	SEC = "SEC"
	SED = "SED"
	SEI = "SEI"
	STA = "STA"
	STX = "STX"
	STY = "STY"
	TAX = "TAX"
	TAY = "TAY"
	TSX = "TSX"
	TXA = "TXA"
	TXS = "TXS"
	TYA = "TYA"
)

func OpCodes(p *Cpu) []*OpCodeDef {
	opCodes := make([]*OpCodeDef, 0xFF)
	id := NewInstructionDefinition(GetAddressingMode)
	opCodes[0x69] = id.Instruction(Mnemonic(ADC, ImmediateModeStr), 2, p.GetAdc)
	opCodes[0x65] = id.Instruction(Mnemonic(ADC, ZeropageModeStr), 3, p.GetAdc)
	opCodes[0x75] = id.Instruction(Mnemonic(ADC, ZeropageXModeStr), 4, p.GetAdc)
	opCodes[0x6D] = id.Instruction(Mnemonic(ADC, AbsoluteModeStr), 4, p.GetAdc)
	opCodes[0x7D] = id.Instruction(Mnemonic(ADC, AbsoluteIndexedXModeStr), 4, p.GetAdc)
	opCodes[0x79] = id.Instruction(Mnemonic(ADC, AbsoluteIndexedYModeStr), 4, p.GetAdc)
	opCodes[0x61] = id.Instruction(Mnemonic(ADC, IndexedIndirectModeStr), 6, p.GetAdc)
	opCodes[0x71] = id.Instruction(Mnemonic(ADC, IndirectIndexedModeStr), 5, p.GetAdc)
	opCodes[0x29] = id.Instruction(Mnemonic(AND, ImmediateModeStr), 2, p.GetAnd)
	opCodes[0x25] = id.Instruction(Mnemonic(AND, ZeropageModeStr), 3, p.GetAnd)
	opCodes[0x35] = id.Instruction(Mnemonic(AND, ZeropageXModeStr), 4, p.GetAnd)
	opCodes[0x2D] = id.Instruction(Mnemonic(AND, AbsoluteModeStr), 4, p.GetAnd)
	opCodes[0x3D] = id.Instruction(Mnemonic(AND, AbsoluteIndexedXModeStr), 4, p.GetAnd)
	opCodes[0x39] = id.Instruction(Mnemonic(AND, AbsoluteIndexedYModeStr), 4, p.GetAnd)
	opCodes[0x21] = id.Instruction(Mnemonic(AND, IndexedIndirectModeStr), 6, p.GetAnd)
	opCodes[0x31] = id.Instruction(Mnemonic(AND, IndirectIndexedModeStr), 5, p.GetAnd)
	opCodes[0x0A] = id.Instruction(Mnemonic(ASL, AccumulatorModeStr), 2, p.GetAsl)
	opCodes[0x06] = id.Instruction(Mnemonic(ASL, ZeropageModeStr), 5, p.GetAsl)
	opCodes[0x16] = id.Instruction(Mnemonic(ASL, ZeropageXModeStr), 6, p.GetAsl)
	opCodes[0x0E] = id.Instruction(Mnemonic(ASL, AbsoluteModeStr), 6, p.GetAsl)
	opCodes[0x1E] = id.Instruction(Mnemonic(ASL, AbsoluteIndexedXModeStr), 7, p.GetAsl)
	opCodes[0x90] = id.Instruction(Mnemonic(BCC, RelativeModeStr), 2, p.GetBcc)
	opCodes[0xB0] = id.Instruction(Mnemonic(BCS, RelativeModeStr), 2, p.GetBcs)
	opCodes[0xF0] = id.Instruction(Mnemonic(BEQ, RelativeModeStr), 2, p.GetBeq)
	opCodes[0x24] = id.Instruction(Mnemonic(BIT, ZeropageModeStr), 3, p.GetBit)
	opCodes[0x2C] = id.Instruction(Mnemonic(BIT, AbsoluteModeStr), 4, p.GetBit)
	opCodes[0x30] = id.Instruction(Mnemonic(BMI, RelativeModeStr), 2, p.GetBmi)
	opCodes[0xD0] = id.Instruction(Mnemonic(BNE, RelativeModeStr), 2, p.GetBne)
	opCodes[0x10] = id.Instruction(Mnemonic(BPL, RelativeModeStr), 2, p.GetBpl)
	opCodes[0x00] = id.Instruction(Mnemonic(BRK, ImpliedModeStr), 7, p.GetBrk)
	opCodes[0x50] = id.Instruction(Mnemonic(BVC, RelativeModeStr), 2, p.GetBvc)
	opCodes[0x70] = id.Instruction(Mnemonic(BVS, RelativeModeStr), 2, p.GetBvs)
	opCodes[0x18] = id.Instruction(Mnemonic(CLC, ImpliedModeStr), 2, p.GetClc)
	opCodes[0xD8] = id.Instruction(Mnemonic(CLD, ImpliedModeStr), 2, p.GetCld)
	opCodes[0x58] = id.Instruction(Mnemonic(CLI, ImpliedModeStr), 2, p.GetCli)
	opCodes[0xB8] = id.Instruction(Mnemonic(CLV, ImpliedModeStr), 2, p.GetClv)
	opCodes[0xC9] = id.Instruction(Mnemonic(CMP, ImmediateModeStr), 2, p.GetCmp)
	opCodes[0xC5] = id.Instruction(Mnemonic(CMP, ZeropageModeStr), 3, p.GetCmp)
	opCodes[0xD5] = id.Instruction(Mnemonic(CMP, ZeropageXModeStr), 4, p.GetCmp)
	opCodes[0xCD] = id.Instruction(Mnemonic(CMP, AbsoluteModeStr), 4, p.GetCmp)
	opCodes[0xDD] = id.Instruction(Mnemonic(CMP, AbsoluteIndexedXModeStr), 4, p.GetCmp)
	opCodes[0xD9] = id.Instruction(Mnemonic(CMP, AbsoluteIndexedYModeStr), 4, p.GetCmp)
	opCodes[0xC1] = id.Instruction(Mnemonic(CMP, IndexedIndirectModeStr), 6, p.GetCmp)
	opCodes[0xD1] = id.Instruction(Mnemonic(CMP, IndirectIndexedModeStr), 5, p.GetCmp)
	opCodes[0xE0] = id.Instruction(Mnemonic(CPX, ImmediateModeStr), 2, p.GetCpx)
	opCodes[0xE4] = id.Instruction(Mnemonic(CPX, ZeropageModeStr), 3, p.GetCpx)
	opCodes[0xEC] = id.Instruction(Mnemonic(CPX, AbsoluteModeStr), 4, p.GetCpx)
	opCodes[0xC0] = id.Instruction(Mnemonic(CPY, ImmediateModeStr), 2, p.GetCpy)
	opCodes[0xC4] = id.Instruction(Mnemonic(CPY, ZeropageModeStr), 3, p.GetCpy)
	opCodes[0xCC] = id.Instruction(Mnemonic(CPY, AbsoluteModeStr), 4, p.GetCpy)
	opCodes[0xC6] = id.Instruction(Mnemonic(DEC, ZeropageModeStr), 5, p.GetDec)
	opCodes[0xD6] = id.Instruction(Mnemonic(DEC, ZeropageXModeStr), 6, p.GetDec)
	opCodes[0xCE] = id.Instruction(Mnemonic(DEC, AbsoluteModeStr), 6, p.GetDec)
	opCodes[0xDE] = id.Instruction(Mnemonic(DEC, AbsoluteIndexedXModeStr), 7, p.GetDec)
	opCodes[0xCA] = id.Instruction(Mnemonic(DEX, ImpliedModeStr), 2, p.GetDex)
	opCodes[0x88] = id.Instruction(Mnemonic(DEY, ImpliedModeStr), 2, p.GetDey)
	opCodes[0x49] = id.Instruction(Mnemonic(EOR, ImmediateModeStr), 2, p.GetEor)
	opCodes[0x45] = id.Instruction(Mnemonic(EOR, ZeropageModeStr), 3, p.GetEor)
	opCodes[0x55] = id.Instruction(Mnemonic(EOR, ZeropageXModeStr), 4, p.GetEor)
	opCodes[0x4D] = id.Instruction(Mnemonic(EOR, AbsoluteModeStr), 4, p.GetEor)
	opCodes[0x5D] = id.Instruction(Mnemonic(EOR, AbsoluteIndexedXModeStr), 4, p.GetEor)
	opCodes[0x59] = id.Instruction(Mnemonic(EOR, AbsoluteIndexedYModeStr), 4, p.GetEor)
	opCodes[0x41] = id.Instruction(Mnemonic(EOR, IndexedIndirectModeStr), 6, p.GetEor)
	opCodes[0x51] = id.Instruction(Mnemonic(EOR, IndirectIndexedModeStr), 5, p.GetEor)
	opCodes[0xE6] = id.Instruction(Mnemonic(INC, ZeropageModeStr), 5, p.GetInc)
	opCodes[0xF6] = id.Instruction(Mnemonic(INC, ZeropageXModeStr), 6, p.GetInc)
	opCodes[0xEE] = id.Instruction(Mnemonic(INC, AbsoluteModeStr), 6, p.GetInc)
	opCodes[0xFE] = id.Instruction(Mnemonic(INC, AbsoluteIndexedXModeStr), 7, p.GetInc)
	opCodes[0xE8] = id.Instruction(Mnemonic(INX, ImpliedModeStr), 2, p.GetInx)
	opCodes[0xC8] = id.Instruction(Mnemonic(INY, ImpliedModeStr), 2, p.GetIny)
	opCodes[0x4C] = id.Instruction(Mnemonic(JMP, AbsoluteModeStr), 3, p.GetJmp)
	opCodes[0x6C] = id.Instruction(Mnemonic(JMP, AbsoluteIndirectModeStr), 5, p.GetJmp)
	opCodes[0x20] = id.Instruction(Mnemonic(JSR, AbsoluteModeStr), 6, p.GetJsr)
	opCodes[0xA9] = id.Instruction(Mnemonic(LDA, ImmediateModeStr), 2, p.GetLda)
	opCodes[0xA5] = id.Instruction(Mnemonic(LDA, ZeropageModeStr), 3, p.GetLda)
	opCodes[0xB5] = id.Instruction(Mnemonic(LDA, ZeropageXModeStr), 4, p.GetLda)
	opCodes[0xAD] = id.Instruction(Mnemonic(LDA, AbsoluteModeStr), 4, p.GetLda)
	opCodes[0xBD] = id.Instruction(Mnemonic(LDA, AbsoluteIndexedXModeStr), 4, p.GetLda)
	opCodes[0xB9] = id.Instruction(Mnemonic(LDA, AbsoluteIndexedYModeStr), 4, p.GetLda)
	opCodes[0xA1] = id.Instruction(Mnemonic(LDA, IndexedIndirectModeStr), 6, p.GetLda)
	opCodes[0xB1] = id.Instruction(Mnemonic(LDA, IndirectIndexedModeStr), 5, p.GetLda)
	opCodes[0xA2] = id.Instruction(Mnemonic(LDX, ImmediateModeStr), 2, p.GetLdx)
	opCodes[0xA6] = id.Instruction(Mnemonic(LDX, ZeropageModeStr), 3, p.GetLdx)
	opCodes[0xB6] = id.Instruction(Mnemonic(LDX, ZeropageYModeStr), 4, p.GetLdx)
	opCodes[0xAE] = id.Instruction(Mnemonic(LDX, AbsoluteModeStr), 4, p.GetLdx)
	opCodes[0xBE] = id.Instruction(Mnemonic(LDX, AbsoluteIndexedXModeStr), 4, p.GetLdx)
	opCodes[0xA0] = id.Instruction(Mnemonic(LDY, ImmediateModeStr), 2, p.GetLdy)
	opCodes[0xA4] = id.Instruction(Mnemonic(LDY, ZeropageModeStr), 3, p.GetLdy)
	opCodes[0xB4] = id.Instruction(Mnemonic(LDY, ZeropageXModeStr), 4, p.GetLdy)
	opCodes[0xAC] = id.Instruction(Mnemonic(LDY, AbsoluteModeStr), 4, p.GetLdy)
	opCodes[0xBC] = id.Instruction(Mnemonic(LDY, AbsoluteIndexedXModeStr), 4, p.GetLdy)
	opCodes[0x4A] = id.Instruction(Mnemonic(LSR, AccumulatorModeStr), 2, p.GetLsr)
	opCodes[0x46] = id.Instruction(Mnemonic(LSR, ZeropageModeStr), 5, p.GetLsr)
	opCodes[0x56] = id.Instruction(Mnemonic(LSR, ZeropageXModeStr), 6, p.GetLsr)
	opCodes[0x4E] = id.Instruction(Mnemonic(LSR, AbsoluteModeStr), 6, p.GetLsr)
	opCodes[0x5E] = id.Instruction(Mnemonic(LSR, AbsoluteIndexedXModeStr), 7, p.GetLsr)
	opCodes[0xEA] = id.Instruction(Mnemonic(NOP, ImpliedModeStr), 2, p.GetNop)
	opCodes[0x09] = id.Instruction(Mnemonic(ORA, ImmediateModeStr), 2, p.GetOra)
	opCodes[0x05] = id.Instruction(Mnemonic(ORA, ZeropageModeStr), 3, p.GetOra)
	opCodes[0x15] = id.Instruction(Mnemonic(ORA, ZeropageXModeStr), 4, p.GetOra)
	opCodes[0x0D] = id.Instruction(Mnemonic(ORA, AbsoluteModeStr), 4, p.GetOra)
	opCodes[0x1D] = id.Instruction(Mnemonic(ORA, AbsoluteIndexedXModeStr), 4, p.GetOra)
	opCodes[0x19] = id.Instruction(Mnemonic(ORA, AbsoluteIndexedYModeStr), 4, p.GetOra)
	opCodes[0x01] = id.Instruction(Mnemonic(ORA, IndexedIndirectModeStr), 6, p.GetOra)
	opCodes[0x11] = id.Instruction(Mnemonic(ORA, IndirectIndexedModeStr), 5, p.GetOra)
	opCodes[0x48] = id.Instruction(Mnemonic(PHA, ImpliedModeStr), 3, p.GetPha)
	opCodes[0x68] = id.Instruction(Mnemonic(PLA, ImpliedModeStr), 4, p.GetPla)
	opCodes[0x08] = id.Instruction(Mnemonic(PHP, ImpliedModeStr), 3, p.GetPhp)
	opCodes[0x28] = id.Instruction(Mnemonic(PLP, ImpliedModeStr), 4, p.GetPlp)
	opCodes[0x2A] = id.Instruction(Mnemonic(ROL, AccumulatorModeStr), 2, p.GetRol)
	opCodes[0x26] = id.Instruction(Mnemonic(ROL, ZeropageModeStr), 5, p.GetRol)
	opCodes[0x36] = id.Instruction(Mnemonic(ROL, ZeropageXModeStr), 6, p.GetRol)
	opCodes[0x2E] = id.Instruction(Mnemonic(ROL, AbsoluteModeStr), 6, p.GetRol)
	opCodes[0x3E] = id.Instruction(Mnemonic(ROL, AbsoluteIndexedXModeStr), 7, p.GetRol)
	opCodes[0x6A] = id.Instruction(Mnemonic(ROR, AccumulatorModeStr), 2, p.GetRor)
	opCodes[0x66] = id.Instruction(Mnemonic(ROR, ZeropageModeStr), 5, p.GetRor)
	opCodes[0x76] = id.Instruction(Mnemonic(ROR, ZeropageXModeStr), 6, p.GetRor)
	opCodes[0x6E] = id.Instruction(Mnemonic(ROR, AbsoluteModeStr), 6, p.GetRor)
	opCodes[0x7E] = id.Instruction(Mnemonic(ROR, AbsoluteIndexedXModeStr), 7, p.GetRor)
	opCodes[0x40] = id.Instruction(Mnemonic(RTI, ImpliedModeStr), 6, p.GetRti)
	opCodes[0x60] = id.Instruction(Mnemonic(RTS, ImpliedModeStr), 6, p.GetRts)
	opCodes[0xE9] = id.Instruction(Mnemonic(SBC, ImmediateModeStr), 2, p.GetSbc)
	opCodes[0xE5] = id.Instruction(Mnemonic(SBC, ZeropageModeStr), 3, p.GetSbc)
	opCodes[0xF5] = id.Instruction(Mnemonic(SBC, ZeropageXModeStr), 4, p.GetSbc)
	opCodes[0xED] = id.Instruction(Mnemonic(SBC, AbsoluteModeStr), 4, p.GetSbc)
	opCodes[0xFD] = id.Instruction(Mnemonic(SBC, AbsoluteIndexedXModeStr), 4, p.GetSbc)
	opCodes[0xF9] = id.Instruction(Mnemonic(SBC, AbsoluteIndexedYModeStr), 4, p.GetSbc)
	opCodes[0xE1] = id.Instruction(Mnemonic(SBC, IndexedIndirectModeStr), 6, p.GetSbc)
	opCodes[0xF1] = id.Instruction(Mnemonic(SBC, IndirectIndexedModeStr), 5, p.GetSbc)
	opCodes[0x38] = id.Instruction(Mnemonic(SEC, IndirectIndexedModeStr), 2, p.GetSec)
	opCodes[0xF8] = id.Instruction(Mnemonic(SED, IndirectIndexedModeStr), 2, p.GetSed)
	opCodes[0x78] = id.Instruction(Mnemonic(SEI, IndirectIndexedModeStr), 2, p.GetSei)
	opCodes[0x85] = id.Instruction(Mnemonic(STA, ZeropageModeStr), 3, p.GetSta)
	opCodes[0x95] = id.Instruction(Mnemonic(STA, ZeropageXModeStr), 4, p.GetSta)
	opCodes[0x8D] = id.Instruction(Mnemonic(STA, AbsoluteModeStr), 4, p.GetSta)
	opCodes[0x9D] = id.Instruction(Mnemonic(STA, AbsoluteIndexedXModeStr), 5, p.GetSta)
	opCodes[0x99] = id.Instruction(Mnemonic(STA, AbsoluteIndexedYModeStr), 5, p.GetSta)
	opCodes[0x81] = id.Instruction(Mnemonic(STA, IndexedIndirectModeStr), 6, p.GetSta)
	opCodes[0x91] = id.Instruction(Mnemonic(STA, IndirectIndexedModeStr), 6, p.GetSta)
	opCodes[0x86] = id.Instruction(Mnemonic(STX, ZeropageModeStr), 3, p.GetStx)
	opCodes[0x96] = id.Instruction(Mnemonic(STX, ZeropageYModeStr), 4, p.GetStx)
	opCodes[0x8E] = id.Instruction(Mnemonic(STX, AbsoluteModeStr), 4, p.GetStx)
	opCodes[0x84] = id.Instruction(Mnemonic(STY, ZeropageModeStr), 3, p.GetSty)
	opCodes[0x94] = id.Instruction(Mnemonic(STY, ZeropageXModeStr), 4, p.GetSty)
	opCodes[0x8C] = id.Instruction(Mnemonic(STY, AbsoluteModeStr), 4, p.GetSty)
	opCodes[0xAA] = id.Instruction(Mnemonic(TAX, ImpliedModeStr), 2, p.GetTax)
	opCodes[0xA8] = id.Instruction(Mnemonic(TAY, ImpliedModeStr), 2, p.GetTay)
	opCodes[0xBA] = id.Instruction(Mnemonic(TSX, ImpliedModeStr), 2, p.GetTsx)
	opCodes[0x8A] = id.Instruction(Mnemonic(TXA, ImpliedModeStr), 2, p.GetTxa)
	opCodes[0x9A] = id.Instruction(Mnemonic(TXS, ImpliedModeStr), 2, p.GetTxs)
	opCodes[0x98] = id.Instruction(Mnemonic(TYA, ImpliedModeStr), 2, p.GetTya)
	return opCodes
}

func (p *Cpu) addPCOffset(b byte) (uint16, bool) {
	pcLsb := uint16(p.Reg.PC & 0x00FF)
	offset := int8(b)
	pcLsb += uint16(offset)
	carryFlag := (pcLsb > uint16(0xFF))
	newPC := p.Reg.PC + uint16(offset)
	return newPC, carryFlag
}

func (p *Cpu) GetAdc(opcode OpCodeDef) InstructionFunc {
	load := opcode.AddressingMode.Load(p, false)

	var adcMode = map[bool]func(b byte){
		true: func(b byte) {
			// BCD Mode
			carryFlag := false
			lowNibble := (p.Reg.A & 0x0F) + (b & 0x0F)
			if p.Reg.IsSet(CarryFlag) {
				lowNibble++
			}
			highNibble := (p.Reg.A & 0xF0) + (b & 0xF0)

			if lowNibble > 0x09 {
				lowNibble += 0x06
			}

			highNibble += (lowNibble & 0xF0) // Add carry from low nibble, if any

			if highNibble > 0x90 {
				highNibble += 0x60 // Decimal adjustment for the high nibble
				carryFlag = true
			}

			result := uint16(lowNibble&0x0F) + uint16(highNibble)
			p.Reg.A = byte(result & 0xFF)
			p.Reg.SetStatus(CarryFlag, carryFlag)
		},
		false: func(b byte) {
			// Binary Mode
			m := p.Reg.A
			r := m + b
			if p.Reg.IsSet(CarryFlag) {
				r++
			}
			p.Reg.A = r

			p.Reg.SetCarryFlag(m, p.Reg.A)
		},
	}

	return func() (Completed, error) {
		b, completed := load()
		if !completed {
			return false, nil
		}
		m := p.Reg.A
		adcMode[p.Reg.IsSet(DecimalFlag)](b)
		p.Reg.SetZeroFlag(p.Reg.A)
		p.Reg.SetOverflowFlag(m, b, p.Reg.A, true)
		p.Reg.SetNegativeFlag(p.Reg.A)
		return true, nil
	}
}

func (p *Cpu) GetAnd(opcode OpCodeDef) InstructionFunc {
	load := opcode.AddressingMode.Load(p, false)

	return func() (Completed, error) {
		b, completed := load()
		if !completed {
			return false, nil
		}
		p.Reg.A = b & p.Reg.A
		p.Reg.SetNegativeFlag(p.Reg.A)
		p.Reg.SetZeroFlag(p.Reg.A)
		return true, nil
	}
}

func (p *Cpu) GetAsl(opcode OpCodeDef) InstructionFunc {
	load := opcode.AddressingMode.Load(p, true)
	store := opcode.AddressingMode.Store(p, true)
	result := byte(0x00)

	return func() (Completed, error) {
		readByte, _ := load()
		result = readByte << 1
		p.Reg.SetNegativeFlag(result)
		p.Reg.SetZeroFlag(result)
		p.Reg.SetCarryFlag(readByte, result)
		store(result)
		return true, nil
	}
}

func (p *Cpu) GetBranch(opcode OpCodeDef, flag StatusFlag, state bool) InstructionFunc {
	load := opcode.AddressingMode.Load(p, true)
	var complete Completed = false
	readByte := byte(0x00)
	var newPC uint16
	var overflow bool

	return func() (Completed, error) {
		if overflow {
			p.Reg.PC = newPC
			return true, nil
		}
		if p.Reg.IsSet(flag) != state {
			return true, nil
		}
		readByte, complete = load()
		if !complete {
			return false, nil
		}
		newPC, overflow = p.addPCOffset(readByte)
		if overflow {
			return false, nil
		}
		p.Reg.PC = newPC
		return true, nil
	}
}

// A Branch not taken requires 2 cycles
func (p *Cpu) GetBcc(opcode OpCodeDef) InstructionFunc {
	return p.GetBranch(opcode, CarryFlag, false)
}

func (p *Cpu) GetBcs(opcode OpCodeDef) InstructionFunc {
	return p.GetBranch(opcode, CarryFlag, true)
}

func (p *Cpu) GetBeq(opcode OpCodeDef) InstructionFunc {
	return p.GetBranch(opcode, ZeroFlag, true)
}

func (p *Cpu) GetBit(opcode OpCodeDef) InstructionFunc {
	load := opcode.AddressingMode.Load(p, false)

	return func() (Completed, error) {
		b, completed := load()
		if !completed {
			return false, nil
		}
		result := b & p.Reg.A
		bit7 := (b >> 7) & 0x01
		bit6 := (b >> 6) & 0x01
		p.Reg.SetStatus(NegativeFlag, (bit7 == 1))
		p.Reg.SetStatus(OverflowFlag, (bit6 == 1))
		p.Reg.SetZeroFlag(result)
		return true, nil
	}
}

func (p *Cpu) GetBmi(opcode OpCodeDef) InstructionFunc {
	return p.GetBranch(opcode, NegativeFlag, true)
}

func (p *Cpu) GetBne(opcode OpCodeDef) InstructionFunc {
	return p.GetBranch(opcode, ZeroFlag, false)
}

func (p *Cpu) GetBpl(opcode OpCodeDef) InstructionFunc {
	return p.GetBranch(opcode, NegativeFlag, false)
}

// cc	addr	data
// --	----	----
// 1	PC	00	;BRK opcode
// 2	PC+1	??	;the padding byte, ignored by the CPU
// 3	S	PCH	;high byte of PC
// 4	S-1	PCL	;low byte of PC
// 5	S-2	P	;status flags with B flag set
// 6	FFFE	??	;low byte of target address
// 7	FFFF	??	;high byte of target address

func (p *Cpu) GetBrk(opcode OpCodeDef) InstructionFunc {
	//	load, _ := p.immediateMode(false) // This is a bit odd, but BRK actually reads the next byte
	//	cycle := 1
	var highPC byte
	var lowPC byte
	return func() (Completed, error) {
		p.Reg.PC++
		p.Push(byte(p.Reg.PC >> 8))
		p.Push(byte(p.Reg.PC & 0xff))
		p.Reg.SetStatus(BreakFlag, true)
		p.Push(byte(p.Reg.Status))
		lowPC = p.mem.Read(irqVector)
		highPC = p.mem.Read(irqVector + 1)
		p.Reg.PC = (uint16(highPC) << 8) | uint16(lowPC)
		return true, nil
	}
}

func (p *Cpu) GetBvc(opcode OpCodeDef) InstructionFunc {
	return p.GetBranch(opcode, OverflowFlag, false)
}

func (p *Cpu) GetBvs(opcode OpCodeDef) InstructionFunc {
	return p.GetBranch(opcode, OverflowFlag, true)
}

func (p *Cpu) GetClc(opcode OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Reg.SetStatus(CarryFlag, false)
		return true, nil
	}
}

func (p *Cpu) GetCld(opcode OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Reg.SetStatus(DecimalFlag, false)
		return true, nil
	}
}

func (p *Cpu) GetCli(opcode OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Reg.SetStatus(InterruptDisableFlag, false)
		return true, nil
	}
}

func (p *Cpu) GetClv(opcode OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Reg.SetStatus(OverflowFlag, false)
		return true, nil
	}
}

func (p *Cpu) GetCompare(opcode OpCodeDef, regValue uint8) InstructionFunc {
	load := opcode.AddressingMode.Load(p, false)

	return func() (Completed, error) {
		b, completed := load()
		if !completed {
			return false, nil
		}

		result := (regValue - b)
		bit7 := (result >> 7) & 0x01

		p.Reg.SetStatus(ZeroFlag, (result == 0))
		p.Reg.SetStatus(NegativeFlag, (bit7 == 1))
		p.Reg.SetStatus(CarryFlag, (regValue > b))
		return true, nil
	}
}

func (p *Cpu) GetCmp(opcode OpCodeDef) InstructionFunc {
	return p.GetCompare(opcode, p.Reg.A)
}

func (p *Cpu) GetCpx(opcode OpCodeDef) InstructionFunc {
	return p.GetCompare(opcode, p.Reg.X)
}

func (p *Cpu) GetCpy(opcode OpCodeDef) InstructionFunc {
	return p.GetCompare(opcode, p.Reg.Y)
}

func (p *Cpu) GetDec(opcode OpCodeDef) InstructionFunc {
	load := opcode.AddressingMode.Load(p, true)
	store := opcode.AddressingMode.Store(p, true)

	return func() (Completed, error) {
		readByte, _ := load()
		readByte--

		p.Reg.SetNegativeFlag(readByte)
		p.Reg.SetZeroFlag(readByte)

		store(readByte)
		return true, nil
	}
}

func (p *Cpu) GetDex(opcode OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Reg.X--
		p.Reg.SetNegativeFlag(p.Reg.X)
		p.Reg.SetZeroFlag(p.Reg.X)
		return true, nil
	}
}

func (p *Cpu) GetDey(opcode OpCodeDef) InstructionFunc {
	// LoadAm, _ := am()
	return func() (Completed, error) {
		p.Reg.Y--
		p.Reg.SetNegativeFlag(p.Reg.Y)
		p.Reg.SetZeroFlag(p.Reg.Y)
		return true, nil
	}
}

func (p *Cpu) GetEor(opcode OpCodeDef) InstructionFunc {
	load := opcode.AddressingMode.Load(p, false)

	return func() (Completed, error) {
		b, completed := load()
		if !completed {
			return false, nil
		}
		p.Reg.A = b ^ p.Reg.A
		p.Reg.SetNegativeFlag(p.Reg.A)
		p.Reg.SetZeroFlag(p.Reg.A)
		return true, nil
	}
}

func (p *Cpu) GetInc(opcode OpCodeDef) InstructionFunc {
	load := opcode.AddressingMode.Load(p, true)
	store := opcode.AddressingMode.Store(p, true)

	return func() (Completed, error) {
		readByte, _ := load()
		readByte++

		p.Reg.SetNegativeFlag(readByte)
		p.Reg.SetZeroFlag(readByte)

		store(readByte)
		return true, nil
	}
}

func (p *Cpu) GetInx(opcode OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Reg.X++
		p.Reg.SetNegativeFlag(p.Reg.X)
		p.Reg.SetZeroFlag(p.Reg.X)
		return true, nil
	}
}

func (p *Cpu) GetIny(opcode OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Reg.Y++
		p.Reg.SetNegativeFlag(p.Reg.Y)
		p.Reg.SetZeroFlag(p.Reg.Y)
		return true, nil
	}
}

func (p *Cpu) GetJmp(opcode OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		address := opcode.AddressingMode.Address(p)
		p.Reg.PC = address
		return true, nil
	}
}

func (p *Cpu) GetJsr(opcode OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Push(byte((p.Reg.PC & 0xFF00) >> 8))
		p.Push(byte(p.Reg.PC & 0x00FF))
		address := opcode.AddressingMode.Address(p)
		p.Reg.PC = address
		return true, nil
	}
}

func (p *Cpu) GetLda(opcode OpCodeDef) InstructionFunc {
	load := opcode.AddressingMode.Load(p, false)

	return func() (Completed, error) {
		b, completed := load()
		if !completed {
			return false, nil
		}
		p.Reg.A = b
		p.Reg.SetNegativeFlag(p.Reg.A)
		p.Reg.SetZeroFlag(p.Reg.A)
		return true, nil
	}
}

func (p *Cpu) GetLdx(opcode OpCodeDef) InstructionFunc {
	load := opcode.AddressingMode.Load(p, false)

	return func() (Completed, error) {
		b, completed := load()
		if !completed {
			return false, nil
		}
		p.Reg.X = b
		p.Reg.SetNegativeFlag(p.Reg.X)
		p.Reg.SetZeroFlag(p.Reg.X)
		return true, nil
	}
}

func (p *Cpu) GetLdy(opcode OpCodeDef) InstructionFunc {
	load := opcode.AddressingMode.Load(p, false)

	return func() (Completed, error) {
		b, completed := load()
		if !completed {
			return false, nil
		}
		p.Reg.Y = b
		p.Reg.SetNegativeFlag(p.Reg.Y)
		p.Reg.SetZeroFlag(p.Reg.Y)
		return true, nil
	}
}

func (p *Cpu) GetLsr(opcode OpCodeDef) InstructionFunc {
	load := opcode.AddressingMode.Load(p, true)
	store := opcode.AddressingMode.Store(p, true)

	// LoadAm, _ := am()
	return func() (Completed, error) {
		b, completed := load()
		if !completed {
			return false, nil
		}

		bit0 := b & 0x1

		b = b >> 1
		p.Reg.SetZeroFlag(b)
		p.Reg.SetStatus(CarryFlag, bit0 == 1)

		store(b)
		return true, nil
	}
}

func (p *Cpu) GetNop(opcode OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		return true, nil
	}
}

func (p *Cpu) GetOra(opcode OpCodeDef) InstructionFunc {
	load := opcode.AddressingMode.Load(p, false)

	return func() (Completed, error) {
		b, completed := load()
		if !completed {
			return false, nil
		}
		p.Reg.A = b | p.Reg.A
		p.Reg.SetNegativeFlag(p.Reg.A)
		p.Reg.SetZeroFlag(p.Reg.A)
		return true, nil
	}
}

func (p *Cpu) GetPha(opcode OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Push(p.Reg.A)
		return true, nil
	}
}

func (p *Cpu) GetPhp(opcode OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Push(p.Reg.S)
		return true, nil
	}
}

func (p *Cpu) GetPla(opcode OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Reg.A = p.Pop()
		return true, nil
	}
}

func (p *Cpu) GetPlp(opcode OpCodeDef) InstructionFunc {
	// LoadAm, _ := am()
	return func() (Completed, error) {
		p.Reg.Status = p.Pop()
		return true, nil
	}
}

func (p *Cpu) GetRol(opcode OpCodeDef) InstructionFunc {
	load := opcode.AddressingMode.Load(p, true)
	store := opcode.AddressingMode.Store(p, true)

	return func() (Completed, error) {
		b, complete := load()
		if !complete {
			return false, nil
		}

		carry := (b & 0x80) != 0
		b = (b << 1)
		if p.Reg.IsSet(CarryFlag) {
			b = b | 0x01
		}

		p.Reg.SetStatus(CarryFlag, carry)
		p.Reg.SetNegativeFlag(b)
		p.Reg.SetZeroFlag(b)

		store(b)
		return true, nil
	}
}

func (p *Cpu) GetRor(opcode OpCodeDef) InstructionFunc {
	load := opcode.AddressingMode.Load(p, true)
	store := opcode.AddressingMode.Store(p, true)

	return func() (Completed, error) {
		b, complete := load()
		if !complete {
			return false, nil
		}

		carry := (b & 0x01) != 0
		b = (b >> 1)
		if p.Reg.IsSet(CarryFlag) {
			b = b | 0x80
		}

		p.Reg.SetStatus(CarryFlag, carry)
		p.Reg.SetNegativeFlag(b)
		p.Reg.SetZeroFlag(b)

		store(b)
		return true, nil
	}
}

func (p *Cpu) GetRti(opcode OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Reg.Status = p.Pop()
		lowBytePC := p.Pop()
		hiBytePC := p.Pop()
		p.Reg.PC = (uint16(lowBytePC) | (uint16(hiBytePC) << 8))
		p.Reg.SetStatus(BreakFlag, false)
		return true, nil
	}
}

func (p *Cpu) GetRts(opcode OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		lowBytePC := p.Pop()
		hiBytePC := p.Pop()
		p.Reg.PC = (uint16(lowBytePC) | (uint16(hiBytePC) << 8))
		return true, nil
	}
}

func (p *Cpu) GetSbc(opcode OpCodeDef) InstructionFunc {
	load := opcode.AddressingMode.Load(p, false)

	var sbcMode = map[bool]func(b byte){
		true: func(b byte) {
			// BCD Mode
			carry := byte(1)
			if !p.Reg.IsSet(CarryFlag) {
				carry = 0
			}

			// Extract high and low nibbles
			accLow := p.Reg.A & 0x0F
			accHigh := p.Reg.A & 0xF0
			bLow := b & 0x0F
			bHigh := b & 0xF0

			// Perform BCD subtraction on low nibble
			lowNibble := accLow - bLow - (1 - carry)
			borrow := byte(0)
			if int8(lowNibble) < 0 {
				lowNibble += 10 // Adjust for BCD
				borrow = 0x10   // Generate a borrow for the high nibble
			}

			// Perform BCD subtraction on high nibble
			highNibble := accHigh - bHigh - borrow
			if int8(highNibble) < 0 {
				highNibble += 0xA0 // Adjust for BCD
			}

			// Combine high and low nibbles
			result := (highNibble & 0xF0) + (lowNibble & 0x0F)

			p.Reg.A = byte(result)

			// Set or clear the Carry flag
			p.Reg.SetStatus(CarryFlag, int8(highNibble) >= 0)
		},

		false: func(b byte) {
			// Binary Mode
			m := p.Reg.A
			c := byte(0)
			if p.Reg.IsSet(CarryFlag) {
				c = 1
			}
			r := m - b - (1 - c)
			p.Reg.A = r

			// Update the Carry flag
			p.Reg.SetStatus(CarryFlag, m >= (b+(1-c)))
		},
	}

	return func() (Completed, error) {
		b, completed := load()
		if !completed {
			return false, nil
		}
		m := p.Reg.A
		sbcMode[p.Reg.IsSet(DecimalFlag)](b)
		p.Reg.SetZeroFlag(p.Reg.A)
		p.Reg.SetNegativeFlag(p.Reg.A)
		p.Reg.SetOverflowFlag(m, b, p.Reg.A, false)
		return true, nil
	}
}

func (p *Cpu) GetSec(opcode OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Reg.SetStatus(CarryFlag, true)
		return true, nil
	}
}

func (p *Cpu) GetSed(opcode OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Reg.SetStatus(DecimalFlag, true)
		return true, nil
	}
}

func (p *Cpu) GetSei(opcode OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Reg.SetStatus(InterruptDisableFlag, true)
		return true, nil
	}
}

func (p *Cpu) GetSta(opcode OpCodeDef) InstructionFunc {
	store := opcode.AddressingMode.Store(p, true)
	return func() (Completed, error) {
		store(p.Reg.A)
		return true, nil
	}
}

func (p *Cpu) GetStx(opcode OpCodeDef) InstructionFunc {
	store := opcode.AddressingMode.Store(p, true)
	return func() (Completed, error) {
		store(p.Reg.X)
		return true, nil
	}
}

func (p *Cpu) GetSty(opcode OpCodeDef) InstructionFunc {
	store := opcode.AddressingMode.Store(p, true)
	return func() (Completed, error) {
		store(p.Reg.Y)
		return true, nil
	}
}

func (p *Cpu) GetTax(opcode OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Reg.X = p.Reg.A
		p.Reg.SetZeroFlag(p.Reg.X)
		p.Reg.SetNegativeFlag(p.Reg.X)
		return true, nil
	}
}

func (p *Cpu) GetTay(opcode OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Reg.Y = p.Reg.A
		p.Reg.SetZeroFlag(p.Reg.Y)
		p.Reg.SetNegativeFlag(p.Reg.Y)
		return true, nil
	}
}

func (p *Cpu) GetTsx(opcode OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Reg.X = p.Reg.S
		p.Reg.SetZeroFlag(p.Reg.X)
		p.Reg.SetNegativeFlag(p.Reg.X)
		return true, nil
	}
}

func (p *Cpu) GetTxa(opcode OpCodeDef) InstructionFunc {
	// LoadAm, _ := am()
	return func() (Completed, error) {
		p.Reg.A = p.Reg.X
		p.Reg.SetZeroFlag(p.Reg.A)
		p.Reg.SetNegativeFlag(p.Reg.A)
		return true, nil
	}
}

func (p *Cpu) GetTxs(opcode OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Reg.S = p.Reg.X
		p.Reg.SetZeroFlag(p.Reg.S)
		p.Reg.SetNegativeFlag(p.Reg.S)
		return true, nil
	}
}

func (p *Cpu) GetTya(opcode OpCodeDef) InstructionFunc {
	// LoadAm, _ := am()
	return func() (Completed, error) {
		p.Reg.A = p.Reg.Y
		p.Reg.SetZeroFlag(p.Reg.A)
		p.Reg.SetNegativeFlag(p.Reg.A)
		return true, nil
	}
}
