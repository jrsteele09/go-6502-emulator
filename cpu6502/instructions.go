package cpu6502

import "github.com/jrsteele09/go-65xx-emulator/cpu65xxx"

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

func OpCodes(p *Cpu) []*cpu65xxx.OpCodeDef {
	opCodes := make([]*cpu65xxx.OpCodeDef, 0xFF)
	id := cpu65xxx.NewInstructionDefinition(cpu65xxx.GetAddressingMode)
	opCodes[0x69] = id.Instruction(cpu65xxx.Mnemonic(ADC, cpu65xxx.ImmediateModeStr), 2, p.GetAdc)
	opCodes[0x65] = id.Instruction(cpu65xxx.Mnemonic(ADC, cpu65xxx.ZeropageModeStr), 3, p.GetAdc)
	opCodes[0x75] = id.Instruction(cpu65xxx.Mnemonic(ADC, cpu65xxx.ZeropageXModeStr), 4, p.GetAdc)
	opCodes[0x6D] = id.Instruction(cpu65xxx.Mnemonic(ADC, cpu65xxx.AbsoluteModeStr), 4, p.GetAdc)
	opCodes[0x7D] = id.Instruction(cpu65xxx.Mnemonic(ADC, cpu65xxx.AbsoluteIndexedXModeStr), 4, p.GetAdc)
	opCodes[0x79] = id.Instruction(cpu65xxx.Mnemonic(ADC, cpu65xxx.AbsoluteIndexedYModeStr), 4, p.GetAdc)
	opCodes[0x61] = id.Instruction(cpu65xxx.Mnemonic(ADC, cpu65xxx.IndexedIndirectModeStr), 6, p.GetAdc)
	opCodes[0x71] = id.Instruction(cpu65xxx.Mnemonic(ADC, cpu65xxx.IndirectIndexedModeStr), 5, p.GetAdc)
	opCodes[0x29] = id.Instruction(cpu65xxx.Mnemonic(AND, cpu65xxx.ImmediateModeStr), 2, p.GetAnd)
	opCodes[0x25] = id.Instruction(cpu65xxx.Mnemonic(AND, cpu65xxx.ZeropageModeStr), 3, p.GetAnd)
	opCodes[0x35] = id.Instruction(cpu65xxx.Mnemonic(AND, cpu65xxx.ZeropageXModeStr), 4, p.GetAnd)
	opCodes[0x2D] = id.Instruction(cpu65xxx.Mnemonic(AND, cpu65xxx.AbsoluteModeStr), 4, p.GetAnd)
	opCodes[0x3D] = id.Instruction(cpu65xxx.Mnemonic(AND, cpu65xxx.AbsoluteIndexedXModeStr), 4, p.GetAnd)
	opCodes[0x39] = id.Instruction(cpu65xxx.Mnemonic(AND, cpu65xxx.AbsoluteIndexedYModeStr), 4, p.GetAnd)
	opCodes[0x21] = id.Instruction(cpu65xxx.Mnemonic(AND, cpu65xxx.IndexedIndirectModeStr), 6, p.GetAnd)
	opCodes[0x31] = id.Instruction(cpu65xxx.Mnemonic(AND, cpu65xxx.IndirectIndexedModeStr), 5, p.GetAnd)
	opCodes[0x0A] = id.Instruction(cpu65xxx.Mnemonic(ASL, cpu65xxx.AccumulatorModeStr), 2, p.GetAsl)
	opCodes[0x06] = id.Instruction(cpu65xxx.Mnemonic(ASL, cpu65xxx.ZeropageModeStr), 5, p.GetAsl)
	opCodes[0x16] = id.Instruction(cpu65xxx.Mnemonic(ASL, cpu65xxx.ZeropageXModeStr), 6, p.GetAsl)
	opCodes[0x0E] = id.Instruction(cpu65xxx.Mnemonic(ASL, cpu65xxx.AbsoluteModeStr), 6, p.GetAsl)
	opCodes[0x1E] = id.Instruction(cpu65xxx.Mnemonic(ASL, cpu65xxx.AbsoluteIndexedXModeStr), 7, p.GetAsl)
	opCodes[0x90] = id.Instruction(cpu65xxx.Mnemonic(BCC, cpu65xxx.RelativeModeStr), 2, p.GetBcc)
	opCodes[0xB0] = id.Instruction(cpu65xxx.Mnemonic(BCS, cpu65xxx.RelativeModeStr), 2, p.GetBcs)
	opCodes[0xF0] = id.Instruction(cpu65xxx.Mnemonic(BEQ, cpu65xxx.RelativeModeStr), 2, p.GetBeq)
	opCodes[0x24] = id.Instruction(cpu65xxx.Mnemonic(BIT, cpu65xxx.ZeropageModeStr), 3, p.GetBit)
	opCodes[0x2C] = id.Instruction(cpu65xxx.Mnemonic(BIT, cpu65xxx.AbsoluteModeStr), 4, p.GetBit)
	opCodes[0x30] = id.Instruction(cpu65xxx.Mnemonic(BMI, cpu65xxx.RelativeModeStr), 2, p.GetBmi)
	opCodes[0xD0] = id.Instruction(cpu65xxx.Mnemonic(BNE, cpu65xxx.RelativeModeStr), 2, p.GetBne)
	opCodes[0x10] = id.Instruction(cpu65xxx.Mnemonic(BPL, cpu65xxx.RelativeModeStr), 2, p.GetBpl)
	opCodes[0x00] = id.Instruction(cpu65xxx.Mnemonic(BRK, cpu65xxx.ImpliedModeStr), 7, p.GetBrk)
	opCodes[0x50] = id.Instruction(cpu65xxx.Mnemonic(BVC, cpu65xxx.RelativeModeStr), 2, p.GetBvc)
	opCodes[0x70] = id.Instruction(cpu65xxx.Mnemonic(BVS, cpu65xxx.RelativeModeStr), 2, p.GetBvs)
	opCodes[0x18] = id.Instruction(cpu65xxx.Mnemonic(CLC, cpu65xxx.ImpliedModeStr), 2, p.GetClc)
	opCodes[0xD8] = id.Instruction(cpu65xxx.Mnemonic(CLD, cpu65xxx.ImpliedModeStr), 2, p.GetCld)
	opCodes[0x58] = id.Instruction(cpu65xxx.Mnemonic(CLI, cpu65xxx.ImpliedModeStr), 2, p.GetCli)
	opCodes[0xB8] = id.Instruction(cpu65xxx.Mnemonic(CLV, cpu65xxx.ImpliedModeStr), 2, p.GetClv)
	opCodes[0xC9] = id.Instruction(cpu65xxx.Mnemonic(CMP, cpu65xxx.ImmediateModeStr), 2, p.GetCmp)
	opCodes[0xC5] = id.Instruction(cpu65xxx.Mnemonic(CMP, cpu65xxx.ZeropageModeStr), 3, p.GetCmp)
	opCodes[0xD5] = id.Instruction(cpu65xxx.Mnemonic(CMP, cpu65xxx.ZeropageXModeStr), 4, p.GetCmp)
	opCodes[0xCD] = id.Instruction(cpu65xxx.Mnemonic(CMP, cpu65xxx.AbsoluteModeStr), 4, p.GetCmp)
	opCodes[0xDD] = id.Instruction(cpu65xxx.Mnemonic(CMP, cpu65xxx.AbsoluteIndexedXModeStr), 4, p.GetCmp)
	opCodes[0xD9] = id.Instruction(cpu65xxx.Mnemonic(CMP, cpu65xxx.AbsoluteIndexedYModeStr), 4, p.GetCmp)
	opCodes[0xC1] = id.Instruction(cpu65xxx.Mnemonic(CMP, cpu65xxx.IndexedIndirectModeStr), 6, p.GetCmp)
	opCodes[0xD1] = id.Instruction(cpu65xxx.Mnemonic(CMP, cpu65xxx.IndirectIndexedModeStr), 5, p.GetCmp)
	opCodes[0xE0] = id.Instruction(cpu65xxx.Mnemonic(CPX, cpu65xxx.ImmediateModeStr), 2, p.GetCpx)
	opCodes[0xE4] = id.Instruction(cpu65xxx.Mnemonic(CPX, cpu65xxx.ZeropageModeStr), 3, p.GetCpx)
	opCodes[0xEC] = id.Instruction(cpu65xxx.Mnemonic(CPX, cpu65xxx.AbsoluteModeStr), 4, p.GetCpx)
	opCodes[0xC0] = id.Instruction(cpu65xxx.Mnemonic(CPY, cpu65xxx.ImmediateModeStr), 2, p.GetCpy)
	opCodes[0xC4] = id.Instruction(cpu65xxx.Mnemonic(CPY, cpu65xxx.ZeropageModeStr), 3, p.GetCpy)
	opCodes[0xCC] = id.Instruction(cpu65xxx.Mnemonic(CPY, cpu65xxx.AbsoluteModeStr), 4, p.GetCpy)
	opCodes[0xC6] = id.Instruction(cpu65xxx.Mnemonic(DEC, cpu65xxx.ZeropageModeStr), 5, p.GetDec)
	opCodes[0xD6] = id.Instruction(cpu65xxx.Mnemonic(DEC, cpu65xxx.ZeropageXModeStr), 6, p.GetDec)
	opCodes[0xCE] = id.Instruction(cpu65xxx.Mnemonic(DEC, cpu65xxx.AbsoluteModeStr), 6, p.GetDec)
	opCodes[0xDE] = id.Instruction(cpu65xxx.Mnemonic(DEC, cpu65xxx.AbsoluteIndexedXModeStr), 7, p.GetDec)
	opCodes[0xCA] = id.Instruction(cpu65xxx.Mnemonic(DEX, cpu65xxx.ImpliedModeStr), 2, p.GetDex)
	opCodes[0x88] = id.Instruction(cpu65xxx.Mnemonic(DEY, cpu65xxx.ImpliedModeStr), 2, p.GetDey)
	opCodes[0x49] = id.Instruction(cpu65xxx.Mnemonic(EOR, cpu65xxx.ImmediateModeStr), 2, p.GetEor)
	opCodes[0x45] = id.Instruction(cpu65xxx.Mnemonic(EOR, cpu65xxx.ZeropageModeStr), 3, p.GetEor)
	opCodes[0x55] = id.Instruction(cpu65xxx.Mnemonic(EOR, cpu65xxx.ZeropageXModeStr), 4, p.GetEor)
	opCodes[0x4D] = id.Instruction(cpu65xxx.Mnemonic(EOR, cpu65xxx.AbsoluteModeStr), 4, p.GetEor)
	opCodes[0x5D] = id.Instruction(cpu65xxx.Mnemonic(EOR, cpu65xxx.AbsoluteIndexedXModeStr), 4, p.GetEor)
	opCodes[0x59] = id.Instruction(cpu65xxx.Mnemonic(EOR, cpu65xxx.AbsoluteIndexedYModeStr), 4, p.GetEor)
	opCodes[0x41] = id.Instruction(cpu65xxx.Mnemonic(EOR, cpu65xxx.IndexedIndirectModeStr), 6, p.GetEor)
	opCodes[0x51] = id.Instruction(cpu65xxx.Mnemonic(EOR, cpu65xxx.IndirectIndexedModeStr), 5, p.GetEor)
	opCodes[0xE6] = id.Instruction(cpu65xxx.Mnemonic(INC, cpu65xxx.ZeropageModeStr), 5, p.GetInc)
	opCodes[0xF6] = id.Instruction(cpu65xxx.Mnemonic(INC, cpu65xxx.ZeropageXModeStr), 6, p.GetInc)
	opCodes[0xEE] = id.Instruction(cpu65xxx.Mnemonic(INC, cpu65xxx.AbsoluteModeStr), 6, p.GetInc)
	opCodes[0xFE] = id.Instruction(cpu65xxx.Mnemonic(INC, cpu65xxx.AbsoluteIndexedXModeStr), 7, p.GetInc)
	opCodes[0xE8] = id.Instruction(cpu65xxx.Mnemonic(INX, cpu65xxx.ImpliedModeStr), 2, p.GetInx)
	opCodes[0xC8] = id.Instruction(cpu65xxx.Mnemonic(INY, cpu65xxx.ImpliedModeStr), 2, p.GetIny)
	opCodes[0x4C] = id.Instruction(cpu65xxx.Mnemonic(JMP, cpu65xxx.AbsoluteModeStr), 3, p.GetJmp)
	opCodes[0x6C] = id.Instruction(cpu65xxx.Mnemonic(JMP, cpu65xxx.AbsoluteIndirectModeStr), 5, p.GetJmp)
	opCodes[0x20] = id.Instruction(cpu65xxx.Mnemonic(JSR, cpu65xxx.AbsoluteModeStr), 6, p.GetJsr)
	opCodes[0xA9] = id.Instruction(cpu65xxx.Mnemonic(LDA, cpu65xxx.ImmediateModeStr), 2, p.GetLda)
	opCodes[0xA5] = id.Instruction(cpu65xxx.Mnemonic(LDA, cpu65xxx.ZeropageModeStr), 3, p.GetLda)
	opCodes[0xB5] = id.Instruction(cpu65xxx.Mnemonic(LDA, cpu65xxx.ZeropageXModeStr), 4, p.GetLda)
	opCodes[0xAD] = id.Instruction(cpu65xxx.Mnemonic(LDA, cpu65xxx.AbsoluteModeStr), 4, p.GetLda)
	opCodes[0xBD] = id.Instruction(cpu65xxx.Mnemonic(LDA, cpu65xxx.AbsoluteIndexedXModeStr), 4, p.GetLda)
	opCodes[0xB9] = id.Instruction(cpu65xxx.Mnemonic(LDA, cpu65xxx.AbsoluteIndexedYModeStr), 4, p.GetLda)
	opCodes[0xA1] = id.Instruction(cpu65xxx.Mnemonic(LDA, cpu65xxx.IndexedIndirectModeStr), 6, p.GetLda)
	opCodes[0xB1] = id.Instruction(cpu65xxx.Mnemonic(LDA, cpu65xxx.IndirectIndexedModeStr), 5, p.GetLda)
	opCodes[0xA2] = id.Instruction(cpu65xxx.Mnemonic(LDX, cpu65xxx.ImmediateModeStr), 2, p.GetLdx)
	opCodes[0xA6] = id.Instruction(cpu65xxx.Mnemonic(LDX, cpu65xxx.ZeropageModeStr), 3, p.GetLdx)
	opCodes[0xB6] = id.Instruction(cpu65xxx.Mnemonic(LDX, cpu65xxx.ZeropageYModeStr), 4, p.GetLdx)
	opCodes[0xAE] = id.Instruction(cpu65xxx.Mnemonic(LDX, cpu65xxx.AbsoluteModeStr), 4, p.GetLdx)
	opCodes[0xBE] = id.Instruction(cpu65xxx.Mnemonic(LDX, cpu65xxx.AbsoluteIndexedXModeStr), 4, p.GetLdx)
	opCodes[0xA0] = id.Instruction(cpu65xxx.Mnemonic(LDY, cpu65xxx.ImmediateModeStr), 2, p.GetLdy)
	opCodes[0xA4] = id.Instruction(cpu65xxx.Mnemonic(LDY, cpu65xxx.ZeropageModeStr), 3, p.GetLdy)
	opCodes[0xB4] = id.Instruction(cpu65xxx.Mnemonic(LDY, cpu65xxx.ZeropageXModeStr), 4, p.GetLdy)
	opCodes[0xAC] = id.Instruction(cpu65xxx.Mnemonic(LDY, cpu65xxx.AbsoluteModeStr), 4, p.GetLdy)
	opCodes[0xBC] = id.Instruction(cpu65xxx.Mnemonic(LDY, cpu65xxx.AbsoluteIndexedXModeStr), 4, p.GetLdy)
	opCodes[0x4A] = id.Instruction(cpu65xxx.Mnemonic(LSR, cpu65xxx.AccumulatorModeStr), 2, p.GetLsr)
	opCodes[0x46] = id.Instruction(cpu65xxx.Mnemonic(LSR, cpu65xxx.ZeropageModeStr), 5, p.GetLsr)
	opCodes[0x56] = id.Instruction(cpu65xxx.Mnemonic(LSR, cpu65xxx.ZeropageXModeStr), 6, p.GetLsr)
	opCodes[0x4E] = id.Instruction(cpu65xxx.Mnemonic(LSR, cpu65xxx.AbsoluteModeStr), 6, p.GetLsr)
	opCodes[0x5E] = id.Instruction(cpu65xxx.Mnemonic(LSR, cpu65xxx.AbsoluteIndexedXModeStr), 7, p.GetLsr)
	opCodes[0xEA] = id.Instruction(cpu65xxx.Mnemonic(NOP, cpu65xxx.ImpliedModeStr), 2, p.GetNop)
	opCodes[0x09] = id.Instruction(cpu65xxx.Mnemonic(ORA, cpu65xxx.ImmediateModeStr), 2, p.GetOra)
	opCodes[0x05] = id.Instruction(cpu65xxx.Mnemonic(ORA, cpu65xxx.ZeropageModeStr), 3, p.GetOra)
	opCodes[0x15] = id.Instruction(cpu65xxx.Mnemonic(ORA, cpu65xxx.ZeropageXModeStr), 4, p.GetOra)
	opCodes[0x0D] = id.Instruction(cpu65xxx.Mnemonic(ORA, cpu65xxx.AbsoluteModeStr), 4, p.GetOra)
	opCodes[0x1D] = id.Instruction(cpu65xxx.Mnemonic(ORA, cpu65xxx.AbsoluteIndexedXModeStr), 4, p.GetOra)
	opCodes[0x19] = id.Instruction(cpu65xxx.Mnemonic(ORA, cpu65xxx.AbsoluteIndexedYModeStr), 4, p.GetOra)
	opCodes[0x01] = id.Instruction(cpu65xxx.Mnemonic(ORA, cpu65xxx.IndexedIndirectModeStr), 6, p.GetOra)
	opCodes[0x11] = id.Instruction(cpu65xxx.Mnemonic(ORA, cpu65xxx.IndirectIndexedModeStr), 5, p.GetOra)
	opCodes[0x48] = id.Instruction(cpu65xxx.Mnemonic(PHA, cpu65xxx.ImpliedModeStr), 3, p.GetPha)
	opCodes[0x68] = id.Instruction(cpu65xxx.Mnemonic(PLA, cpu65xxx.ImpliedModeStr), 4, p.GetPla)
	opCodes[0x08] = id.Instruction(cpu65xxx.Mnemonic(PHP, cpu65xxx.ImpliedModeStr), 3, p.GetPhp)
	opCodes[0x28] = id.Instruction(cpu65xxx.Mnemonic(PLP, cpu65xxx.ImpliedModeStr), 4, p.GetPlp)
	opCodes[0x2A] = id.Instruction(cpu65xxx.Mnemonic(ROL, cpu65xxx.AccumulatorModeStr), 2, p.GetRol)
	opCodes[0x26] = id.Instruction(cpu65xxx.Mnemonic(ROL, cpu65xxx.ZeropageModeStr), 5, p.GetRol)
	opCodes[0x36] = id.Instruction(cpu65xxx.Mnemonic(ROL, cpu65xxx.ZeropageXModeStr), 6, p.GetRol)
	opCodes[0x2E] = id.Instruction(cpu65xxx.Mnemonic(ROL, cpu65xxx.AbsoluteModeStr), 6, p.GetRol)
	opCodes[0x3E] = id.Instruction(cpu65xxx.Mnemonic(ROL, cpu65xxx.AbsoluteIndexedXModeStr), 7, p.GetRol)
	opCodes[0x6A] = id.Instruction(cpu65xxx.Mnemonic(ROR, cpu65xxx.AccumulatorModeStr), 2, p.GetRor)
	opCodes[0x66] = id.Instruction(cpu65xxx.Mnemonic(ROR, cpu65xxx.ZeropageModeStr), 5, p.GetRor)
	opCodes[0x76] = id.Instruction(cpu65xxx.Mnemonic(ROR, cpu65xxx.ZeropageXModeStr), 6, p.GetRor)
	opCodes[0x6E] = id.Instruction(cpu65xxx.Mnemonic(ROR, cpu65xxx.AbsoluteModeStr), 6, p.GetRor)
	opCodes[0x7E] = id.Instruction(cpu65xxx.Mnemonic(ROR, cpu65xxx.AbsoluteIndexedXModeStr), 7, p.GetRor)
	opCodes[0x40] = id.Instruction(cpu65xxx.Mnemonic(RTI, cpu65xxx.ImpliedModeStr), 6, p.GetRti)
	opCodes[0x60] = id.Instruction(cpu65xxx.Mnemonic(RTS, cpu65xxx.ImpliedModeStr), 6, p.GetRts)
	opCodes[0xE9] = id.Instruction(cpu65xxx.Mnemonic(SBC, cpu65xxx.ImmediateModeStr), 2, p.GetSbc)
	opCodes[0xE5] = id.Instruction(cpu65xxx.Mnemonic(SBC, cpu65xxx.ZeropageModeStr), 3, p.GetSbc)
	opCodes[0xF5] = id.Instruction(cpu65xxx.Mnemonic(SBC, cpu65xxx.ZeropageXModeStr), 4, p.GetSbc)
	opCodes[0xED] = id.Instruction(cpu65xxx.Mnemonic(SBC, cpu65xxx.AbsoluteModeStr), 4, p.GetSbc)
	opCodes[0xFD] = id.Instruction(cpu65xxx.Mnemonic(SBC, cpu65xxx.AbsoluteIndexedXModeStr), 4, p.GetSbc)
	opCodes[0xF9] = id.Instruction(cpu65xxx.Mnemonic(SBC, cpu65xxx.AbsoluteIndexedYModeStr), 4, p.GetSbc)
	opCodes[0xE1] = id.Instruction(cpu65xxx.Mnemonic(SBC, cpu65xxx.IndexedIndirectModeStr), 6, p.GetSbc)
	opCodes[0xF1] = id.Instruction(cpu65xxx.Mnemonic(SBC, cpu65xxx.IndirectIndexedModeStr), 5, p.GetSbc)
	opCodes[0x38] = id.Instruction(cpu65xxx.Mnemonic(SEC, cpu65xxx.IndirectIndexedModeStr), 2, p.GetSec)
	opCodes[0xF8] = id.Instruction(cpu65xxx.Mnemonic(SED, cpu65xxx.IndirectIndexedModeStr), 2, p.GetSed)
	opCodes[0x78] = id.Instruction(cpu65xxx.Mnemonic(SEI, cpu65xxx.IndirectIndexedModeStr), 2, p.GetSei)
	opCodes[0x85] = id.Instruction(cpu65xxx.Mnemonic(STA, cpu65xxx.ZeropageModeStr), 3, p.GetSta)
	opCodes[0x95] = id.Instruction(cpu65xxx.Mnemonic(STA, cpu65xxx.ZeropageXModeStr), 4, p.GetSta)
	opCodes[0x8D] = id.Instruction(cpu65xxx.Mnemonic(STA, cpu65xxx.AbsoluteModeStr), 4, p.GetSta)
	opCodes[0x9D] = id.Instruction(cpu65xxx.Mnemonic(STA, cpu65xxx.AbsoluteIndexedXModeStr), 5, p.GetSta)
	opCodes[0x99] = id.Instruction(cpu65xxx.Mnemonic(STA, cpu65xxx.AbsoluteIndexedYModeStr), 5, p.GetSta)
	opCodes[0x81] = id.Instruction(cpu65xxx.Mnemonic(STA, cpu65xxx.IndexedIndirectModeStr), 6, p.GetSta)
	opCodes[0x91] = id.Instruction(cpu65xxx.Mnemonic(STA, cpu65xxx.IndirectIndexedModeStr), 6, p.GetSta)
	opCodes[0x86] = id.Instruction(cpu65xxx.Mnemonic(STX, cpu65xxx.ZeropageModeStr), 3, p.GetStx)
	opCodes[0x96] = id.Instruction(cpu65xxx.Mnemonic(STX, cpu65xxx.ZeropageYModeStr), 4, p.GetStx)
	opCodes[0x8E] = id.Instruction(cpu65xxx.Mnemonic(STX, cpu65xxx.AbsoluteModeStr), 4, p.GetStx)
	opCodes[0x84] = id.Instruction(cpu65xxx.Mnemonic(STY, cpu65xxx.ZeropageModeStr), 3, p.GetSty)
	opCodes[0x94] = id.Instruction(cpu65xxx.Mnemonic(STY, cpu65xxx.ZeropageXModeStr), 4, p.GetSty)
	opCodes[0x8C] = id.Instruction(cpu65xxx.Mnemonic(STY, cpu65xxx.AbsoluteModeStr), 4, p.GetSty)
	opCodes[0xAA] = id.Instruction(cpu65xxx.Mnemonic(TAX, cpu65xxx.ImpliedModeStr), 2, p.GetTax)
	opCodes[0xA8] = id.Instruction(cpu65xxx.Mnemonic(TAY, cpu65xxx.ImpliedModeStr), 2, p.GetTay)
	opCodes[0xBA] = id.Instruction(cpu65xxx.Mnemonic(TSX, cpu65xxx.ImpliedModeStr), 2, p.GetTsx)
	opCodes[0x8A] = id.Instruction(cpu65xxx.Mnemonic(TXA, cpu65xxx.ImpliedModeStr), 2, p.GetTxa)
	opCodes[0x9A] = id.Instruction(cpu65xxx.Mnemonic(TXS, cpu65xxx.ImpliedModeStr), 2, p.GetTxs)
	opCodes[0x98] = id.Instruction(cpu65xxx.Mnemonic(TYA, cpu65xxx.ImpliedModeStr), 2, p.GetTya)
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

func (p *Cpu) GetAdc(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	load := opcode.AddressingMode.Load(p, false)

	var adcMode = map[bool]func(b byte){
		true: func(b byte) {
			// BCD Mode
			carryFlag := false
			lowNibble := (p.Reg.A & 0x0F) + (b & 0x0F)
			if p.Reg.IsSet(cpu65xxx.CarryFlag) {
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
			p.Reg.SetStatus(cpu65xxx.CarryFlag, carryFlag)
		},
		false: func(b byte) {
			// Binary Mode
			m := p.Reg.A
			r := m + b
			if p.Reg.IsSet(cpu65xxx.CarryFlag) {
				r++
			}
			p.Reg.A = r

			p.Reg.SetCarryFlag(m, p.Reg.A)
		},
	}

	return func() (cpu65xxx.Completed, error) {
		b, completed := load()
		if !completed {
			return false, nil
		}
		m := p.Reg.A
		adcMode[p.Reg.IsSet(cpu65xxx.DecimalFlag)](b)
		p.Reg.SetZeroFlag(p.Reg.A)
		p.Reg.SetOverflowFlag(m, b, p.Reg.A, true)
		p.Reg.SetNegativeFlag(p.Reg.A)
		return true, nil
	}
}

func (p *Cpu) GetAnd(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	load := opcode.AddressingMode.Load(p, false)

	return func() (cpu65xxx.Completed, error) {
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

func (p *Cpu) GetAsl(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	load := opcode.AddressingMode.Load(p, true)
	store := opcode.AddressingMode.Store(p, true)
	result := byte(0x00)

	return func() (cpu65xxx.Completed, error) {
		readByte, _ := load()
		result = readByte << 1
		p.Reg.SetNegativeFlag(result)
		p.Reg.SetZeroFlag(result)
		p.Reg.SetCarryFlag(readByte, result)
		store(result)
		return true, nil
	}
}

func (p *Cpu) GetBranch(opcode cpu65xxx.OpCodeDef, flag cpu65xxx.StatusFlag, state bool) cpu65xxx.InstructionFunc {
	load := opcode.AddressingMode.Load(p, true)
	var complete cpu65xxx.Completed = false
	readByte := byte(0x00)
	var newPC uint16
	var overflow bool

	return func() (cpu65xxx.Completed, error) {
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
func (p *Cpu) GetBcc(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	return p.GetBranch(opcode, cpu65xxx.CarryFlag, false)
}

func (p *Cpu) GetBcs(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	return p.GetBranch(opcode, cpu65xxx.CarryFlag, true)
}

func (p *Cpu) GetBeq(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	return p.GetBranch(opcode, cpu65xxx.ZeroFlag, true)
}

func (p *Cpu) GetBit(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	load := opcode.AddressingMode.Load(p, false)

	return func() (cpu65xxx.Completed, error) {
		b, completed := load()
		if !completed {
			return false, nil
		}
		result := b & p.Reg.A
		bit7 := (b >> 7) & 0x01
		bit6 := (b >> 6) & 0x01
		p.Reg.SetStatus(cpu65xxx.NegativeFlag, (bit7 == 1))
		p.Reg.SetStatus(cpu65xxx.OverflowFlag, (bit6 == 1))
		p.Reg.SetZeroFlag(result)
		return true, nil
	}
}

func (p *Cpu) GetBmi(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	return p.GetBranch(opcode, cpu65xxx.NegativeFlag, true)
}

func (p *Cpu) GetBne(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	return p.GetBranch(opcode, cpu65xxx.ZeroFlag, false)
}

func (p *Cpu) GetBpl(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	return p.GetBranch(opcode, cpu65xxx.NegativeFlag, false)
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

func (p *Cpu) GetBrk(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	//	load, _ := p.immediateMode(false) // This is a bit odd, but BRK actually reads the next byte
	//	cycle := 1
	var highPC byte
	var lowPC byte
	return func() (cpu65xxx.Completed, error) {
		p.Reg.PC++
		p.Push(byte(p.Reg.PC >> 8))
		p.Push(byte(p.Reg.PC & 0xff))
		p.Reg.SetStatus(cpu65xxx.BreakFlag, true)
		p.Push(byte(p.Reg.Status))
		lowPC = p.mem.Read(irqVector)
		highPC = p.mem.Read(irqVector + 1)
		p.Reg.PC = (uint16(highPC) << 8) | uint16(lowPC)
		return true, nil
	}
}

func (p *Cpu) GetBvc(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	return p.GetBranch(opcode, cpu65xxx.OverflowFlag, false)
}

func (p *Cpu) GetBvs(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	return p.GetBranch(opcode, cpu65xxx.OverflowFlag, true)
}

func (p *Cpu) GetClc(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	return func() (cpu65xxx.Completed, error) {
		p.Reg.SetStatus(cpu65xxx.CarryFlag, false)
		return true, nil
	}
}

func (p *Cpu) GetCld(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	return func() (cpu65xxx.Completed, error) {
		p.Reg.SetStatus(cpu65xxx.DecimalFlag, false)
		return true, nil
	}
}

func (p *Cpu) GetCli(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	return func() (cpu65xxx.Completed, error) {
		p.Reg.SetStatus(cpu65xxx.InterruptDisableFlag, false)
		return true, nil
	}
}

func (p *Cpu) GetClv(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	return func() (cpu65xxx.Completed, error) {
		p.Reg.SetStatus(cpu65xxx.OverflowFlag, false)
		return true, nil
	}
}

func (p *Cpu) GetCompare(opcode cpu65xxx.OpCodeDef, regValue uint8) cpu65xxx.InstructionFunc {
	load := opcode.AddressingMode.Load(p, false)

	return func() (cpu65xxx.Completed, error) {
		b, completed := load()
		if !completed {
			return false, nil
		}

		result := (regValue - b)
		bit7 := (result >> 7) & 0x01

		p.Reg.SetStatus(cpu65xxx.ZeroFlag, (result == 0))
		p.Reg.SetStatus(cpu65xxx.NegativeFlag, (bit7 == 1))
		p.Reg.SetStatus(cpu65xxx.CarryFlag, (regValue > b))
		return true, nil
	}
}

func (p *Cpu) GetCmp(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	return p.GetCompare(opcode, p.Reg.A)
}

func (p *Cpu) GetCpx(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	return p.GetCompare(opcode, p.Reg.X)
}

func (p *Cpu) GetCpy(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	return p.GetCompare(opcode, p.Reg.Y)
}

func (p *Cpu) GetDec(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	load := opcode.AddressingMode.Load(p, true)
	store := opcode.AddressingMode.Store(p, true)

	return func() (cpu65xxx.Completed, error) {
		readByte, _ := load()
		readByte--

		p.Reg.SetNegativeFlag(readByte)
		p.Reg.SetZeroFlag(readByte)

		store(readByte)
		return true, nil
	}
}

func (p *Cpu) GetDex(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	return func() (cpu65xxx.Completed, error) {
		p.Reg.X--
		p.Reg.SetNegativeFlag(p.Reg.X)
		p.Reg.SetZeroFlag(p.Reg.X)
		return true, nil
	}
}

func (p *Cpu) GetDey(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	// LoadAm, _ := am()
	return func() (cpu65xxx.Completed, error) {
		p.Reg.Y--
		p.Reg.SetNegativeFlag(p.Reg.Y)
		p.Reg.SetZeroFlag(p.Reg.Y)
		return true, nil
	}
}

func (p *Cpu) GetEor(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	load := opcode.AddressingMode.Load(p, false)

	return func() (cpu65xxx.Completed, error) {
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

func (p *Cpu) GetInc(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	load := opcode.AddressingMode.Load(p, true)
	store := opcode.AddressingMode.Store(p, true)

	return func() (cpu65xxx.Completed, error) {
		readByte, _ := load()
		readByte++

		p.Reg.SetNegativeFlag(readByte)
		p.Reg.SetZeroFlag(readByte)

		store(readByte)
		return true, nil
	}
}

func (p *Cpu) GetInx(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	return func() (cpu65xxx.Completed, error) {
		p.Reg.X++
		p.Reg.SetNegativeFlag(p.Reg.X)
		p.Reg.SetZeroFlag(p.Reg.X)
		return true, nil
	}
}

func (p *Cpu) GetIny(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	return func() (cpu65xxx.Completed, error) {
		p.Reg.Y++
		p.Reg.SetNegativeFlag(p.Reg.Y)
		p.Reg.SetZeroFlag(p.Reg.Y)
		return true, nil
	}
}

func (p *Cpu) GetJmp(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	return func() (cpu65xxx.Completed, error) {
		address := opcode.AddressingMode.Address(p)
		p.Reg.PC = address
		return true, nil
	}
}

func (p *Cpu) GetJsr(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	return func() (cpu65xxx.Completed, error) {
		p.Push(byte((p.Reg.PC & 0xFF00) >> 8))
		p.Push(byte(p.Reg.PC & 0x00FF))
		address := opcode.AddressingMode.Address(p)
		p.Reg.PC = address
		return true, nil
	}
}

func (p *Cpu) GetLda(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	load := opcode.AddressingMode.Load(p, false)

	return func() (cpu65xxx.Completed, error) {
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

func (p *Cpu) GetLdx(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	load := opcode.AddressingMode.Load(p, false)

	return func() (cpu65xxx.Completed, error) {
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

func (p *Cpu) GetLdy(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	load := opcode.AddressingMode.Load(p, false)

	return func() (cpu65xxx.Completed, error) {
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

func (p *Cpu) GetLsr(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	load := opcode.AddressingMode.Load(p, true)
	store := opcode.AddressingMode.Store(p, true)

	// LoadAm, _ := am()
	return func() (cpu65xxx.Completed, error) {
		b, completed := load()
		if !completed {
			return false, nil
		}

		bit0 := b & 0x1

		b = b >> 1
		p.Reg.SetZeroFlag(b)
		p.Reg.SetStatus(cpu65xxx.CarryFlag, bit0 == 1)

		store(b)
		return true, nil
	}
}

func (p *Cpu) GetNop(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	return func() (cpu65xxx.Completed, error) {
		return true, nil
	}
}

func (p *Cpu) GetOra(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	load := opcode.AddressingMode.Load(p, false)

	return func() (cpu65xxx.Completed, error) {
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

func (p *Cpu) GetPha(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	return func() (cpu65xxx.Completed, error) {
		p.Push(p.Reg.A)
		return true, nil
	}
}

func (p *Cpu) GetPhp(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	return func() (cpu65xxx.Completed, error) {
		p.Push(p.Reg.S)
		return true, nil
	}
}

func (p *Cpu) GetPla(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	return func() (cpu65xxx.Completed, error) {
		p.Reg.A = p.Pop()
		return true, nil
	}
}

func (p *Cpu) GetPlp(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	// LoadAm, _ := am()
	return func() (cpu65xxx.Completed, error) {
		p.Reg.Status = p.Pop()
		return true, nil
	}
}

func (p *Cpu) GetRol(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	load := opcode.AddressingMode.Load(p, true)
	store := opcode.AddressingMode.Store(p, true)

	return func() (cpu65xxx.Completed, error) {
		b, complete := load()
		if !complete {
			return false, nil
		}

		carry := (b & 0x80) != 0
		b = (b << 1)
		if p.Reg.IsSet(cpu65xxx.CarryFlag) {
			b = b | 0x01
		}

		p.Reg.SetStatus(cpu65xxx.CarryFlag, carry)
		p.Reg.SetNegativeFlag(b)
		p.Reg.SetZeroFlag(b)

		store(b)
		return true, nil
	}
}

func (p *Cpu) GetRor(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	load := opcode.AddressingMode.Load(p, true)
	store := opcode.AddressingMode.Store(p, true)

	return func() (cpu65xxx.Completed, error) {
		b, complete := load()
		if !complete {
			return false, nil
		}

		carry := (b & 0x01) != 0
		b = (b >> 1)
		if p.Reg.IsSet(cpu65xxx.CarryFlag) {
			b = b | 0x80
		}

		p.Reg.SetStatus(cpu65xxx.CarryFlag, carry)
		p.Reg.SetNegativeFlag(b)
		p.Reg.SetZeroFlag(b)

		store(b)
		return true, nil
	}
}

func (p *Cpu) GetRti(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	return func() (cpu65xxx.Completed, error) {
		p.Reg.Status = p.Pop()
		lowBytePC := p.Pop()
		hiBytePC := p.Pop()
		p.Reg.PC = (uint16(lowBytePC) | (uint16(hiBytePC) << 8))
		p.Reg.SetStatus(cpu65xxx.BreakFlag, false)
		return true, nil
	}
}

func (p *Cpu) GetRts(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	return func() (cpu65xxx.Completed, error) {
		lowBytePC := p.Pop()
		hiBytePC := p.Pop()
		p.Reg.PC = (uint16(lowBytePC) | (uint16(hiBytePC) << 8))
		return true, nil
	}
}

func (p *Cpu) GetSbc(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	load := opcode.AddressingMode.Load(p, false)

	var sbcMode = map[bool]func(b byte){
		true: func(b byte) {
			// BCD Mode
			carry := byte(1)
			if !p.Reg.IsSet(cpu65xxx.CarryFlag) {
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
			p.Reg.SetStatus(cpu65xxx.CarryFlag, int8(highNibble) >= 0)
		},

		false: func(b byte) {
			// Binary Mode
			m := p.Reg.A
			c := byte(0)
			if p.Reg.IsSet(cpu65xxx.CarryFlag) {
				c = 1
			}
			r := m - b - (1 - c)
			p.Reg.A = r

			// Update the Carry flag
			p.Reg.SetStatus(cpu65xxx.CarryFlag, m >= (b+(1-c)))
		},
	}

	return func() (cpu65xxx.Completed, error) {
		b, completed := load()
		if !completed {
			return false, nil
		}
		m := p.Reg.A
		sbcMode[p.Reg.IsSet(cpu65xxx.DecimalFlag)](b)
		p.Reg.SetZeroFlag(p.Reg.A)
		p.Reg.SetNegativeFlag(p.Reg.A)
		p.Reg.SetOverflowFlag(m, b, p.Reg.A, false)
		return true, nil
	}
}

func (p *Cpu) GetSec(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	return func() (cpu65xxx.Completed, error) {
		p.Reg.SetStatus(cpu65xxx.CarryFlag, true)
		return true, nil
	}
}

func (p *Cpu) GetSed(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	return func() (cpu65xxx.Completed, error) {
		p.Reg.SetStatus(cpu65xxx.DecimalFlag, true)
		return true, nil
	}
}

func (p *Cpu) GetSei(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	return func() (cpu65xxx.Completed, error) {
		p.Reg.SetStatus(cpu65xxx.InterruptDisableFlag, true)
		return true, nil
	}
}

func (p *Cpu) GetSta(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	store := opcode.AddressingMode.Store(p, true)
	return func() (cpu65xxx.Completed, error) {
		store(p.Reg.A)
		return true, nil
	}
}

func (p *Cpu) GetStx(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	store := opcode.AddressingMode.Store(p, true)
	return func() (cpu65xxx.Completed, error) {
		store(p.Reg.X)
		return true, nil
	}
}

func (p *Cpu) GetSty(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	store := opcode.AddressingMode.Store(p, true)
	return func() (cpu65xxx.Completed, error) {
		store(p.Reg.Y)
		return true, nil
	}
}

func (p *Cpu) GetTax(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	return func() (cpu65xxx.Completed, error) {
		p.Reg.X = p.Reg.A
		p.Reg.SetZeroFlag(p.Reg.X)
		p.Reg.SetNegativeFlag(p.Reg.X)
		return true, nil
	}
}

func (p *Cpu) GetTay(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	return func() (cpu65xxx.Completed, error) {
		p.Reg.Y = p.Reg.A
		p.Reg.SetZeroFlag(p.Reg.Y)
		p.Reg.SetNegativeFlag(p.Reg.Y)
		return true, nil
	}
}

func (p *Cpu) GetTsx(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	return func() (cpu65xxx.Completed, error) {
		p.Reg.X = p.Reg.S
		p.Reg.SetZeroFlag(p.Reg.X)
		p.Reg.SetNegativeFlag(p.Reg.X)
		return true, nil
	}
}

func (p *Cpu) GetTxa(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	// LoadAm, _ := am()
	return func() (cpu65xxx.Completed, error) {
		p.Reg.A = p.Reg.X
		p.Reg.SetZeroFlag(p.Reg.A)
		p.Reg.SetNegativeFlag(p.Reg.A)
		return true, nil
	}
}

func (p *Cpu) GetTxs(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	return func() (cpu65xxx.Completed, error) {
		p.Reg.S = p.Reg.X
		p.Reg.SetZeroFlag(p.Reg.S)
		p.Reg.SetNegativeFlag(p.Reg.S)
		return true, nil
	}
}

func (p *Cpu) GetTya(opcode cpu65xxx.OpCodeDef) cpu65xxx.InstructionFunc {
	// LoadAm, _ := am()
	return func() (cpu65xxx.Completed, error) {
		p.Reg.A = p.Reg.Y
		p.Reg.SetZeroFlag(p.Reg.A)
		p.Reg.SetNegativeFlag(p.Reg.A)
		return true, nil
	}
}
