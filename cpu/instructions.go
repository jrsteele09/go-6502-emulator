package cpu

const (
	adcStr = "ADC"
	andStr = "AND"
	aslStr = "ASL"
	bccStr = "BCC"
	bcsStr = "BCS"
	beqStr = "BEQ"
	bitStr = "BIT"
	bmiStr = "BMI"
	bneStr = "BNE"
	bplStr = "BPL"
	brkStr = "BRK"
	bvcStr = "BVC"
	bvsStr = "BVS"
	clcStr = "CLC"
	cldStr = "CLD"
	cliStr = "CLI"
	clvStr = "CLV"
	cmpStr = "CMP"
	cpxStr = "CPX"
	cpyStr = "CPY"
	decStr = "DEC"
	dexStr = "DEX"
	deyStr = "DEY"
	eorStr = "EOR"
	incStr = "INC"
	inxStr = "INX"
	inyStr = "INY"
	jmpStr = "JMP"
	jsrStr = "JSR"
	ldaStr = "LDA"
	ldxStr = "LDX"
	ldyStr = "LDY"
	lsrStr = "LSR"
	nopStr = "NOP"
	oraStr = "ORA"
	phaStr = "PHA"
	phpStr = "PHP"
	plaStr = "PLA"
	plpStr = "PLP"
	rolStr = "ROL"
	rorStr = "ROR"
	rtiStr = "RTI"
	rtsStr = "RTS"
	sbcStr = "SBC"
	secStr = "SEC"
	sedStr = "SED"
	seiStr = "SEI"
	staStr = "STA"
	stxStr = "STX"
	styStr = "STY"
	taxStr = "TAX"
	tayStr = "TAY"
	tsxStr = "TSX"
	txaStr = "TXA"
	txsStr = "TXS"
	tyaStr = "TYA"
)

// createOpCodes returns a map of the actual 6502 opcode numbers to OpCodeDef instances.
func createOpCodes(p *CPU) []*OpCodeDef {
	opCodes := make([]*OpCodeDef, 256)
	id := NewInstruction(getAddressingMode)
	opCodes[0x69] = id.Instruction(Mnemonic(adcStr, ImmediateModeStr), 2, p.adc)
	opCodes[0x65] = id.Instruction(Mnemonic(adcStr, ZeropageModeStr), 3, p.adc)
	opCodes[0x75] = id.Instruction(Mnemonic(adcStr, ZeropageXModeStr), 4, p.adc)
	opCodes[0x6D] = id.Instruction(Mnemonic(adcStr, AbsoluteModeStr), 4, p.adc)
	opCodes[0x7D] = id.Instruction(Mnemonic(adcStr, AbsoluteIndexedXModeStr), 4, p.adc)
	opCodes[0x79] = id.Instruction(Mnemonic(adcStr, AbsoluteIndexedYModeStr), 4, p.adc)
	opCodes[0x61] = id.Instruction(Mnemonic(adcStr, IndexedIndirectModeStr), 6, p.adc)
	opCodes[0x71] = id.Instruction(Mnemonic(adcStr, IndirectIndexedModeStr), 5, p.adc)
	opCodes[0x29] = id.Instruction(Mnemonic(andStr, ImmediateModeStr), 2, p.and)
	opCodes[0x25] = id.Instruction(Mnemonic(andStr, ZeropageModeStr), 3, p.and)
	opCodes[0x35] = id.Instruction(Mnemonic(andStr, ZeropageXModeStr), 4, p.and)
	opCodes[0x2D] = id.Instruction(Mnemonic(andStr, AbsoluteModeStr), 4, p.and)
	opCodes[0x3D] = id.Instruction(Mnemonic(andStr, AbsoluteIndexedXModeStr), 4, p.and)
	opCodes[0x39] = id.Instruction(Mnemonic(andStr, AbsoluteIndexedYModeStr), 4, p.and)
	opCodes[0x21] = id.Instruction(Mnemonic(andStr, IndexedIndirectModeStr), 6, p.and)
	opCodes[0x31] = id.Instruction(Mnemonic(andStr, IndirectIndexedModeStr), 5, p.and)
	opCodes[0x0A] = id.Instruction(Mnemonic(aslStr, AccumulatorModeStr), 2, p.asl)
	opCodes[0x06] = id.Instruction(Mnemonic(aslStr, ZeropageModeStr), 5, p.asl)
	opCodes[0x16] = id.Instruction(Mnemonic(aslStr, ZeropageXModeStr), 6, p.asl)
	opCodes[0x0E] = id.Instruction(Mnemonic(aslStr, AbsoluteModeStr), 6, p.asl)
	opCodes[0x1E] = id.Instruction(Mnemonic(aslStr, AbsoluteIndexedXModeStr), 7, p.asl)
	opCodes[0x90] = id.Instruction(Mnemonic(bccStr, RelativeModeStr), 2, p.bcc)
	opCodes[0xB0] = id.Instruction(Mnemonic(bcsStr, RelativeModeStr), 2, p.bcs)
	opCodes[0xF0] = id.Instruction(Mnemonic(beqStr, RelativeModeStr), 2, p.beq)
	opCodes[0x24] = id.Instruction(Mnemonic(bitStr, ZeropageModeStr), 3, p.bit)
	opCodes[0x2C] = id.Instruction(Mnemonic(bitStr, AbsoluteModeStr), 4, p.bit)
	opCodes[0x30] = id.Instruction(Mnemonic(bmiStr, RelativeModeStr), 2, p.bmi)
	opCodes[0xD0] = id.Instruction(Mnemonic(bneStr, RelativeModeStr), 2, p.bne)
	opCodes[0x10] = id.Instruction(Mnemonic(bplStr, RelativeModeStr), 2, p.bpl)
	opCodes[0x00] = id.Instruction(Mnemonic(brkStr, ImpliedModeStr), 7, p.brk)
	opCodes[0x50] = id.Instruction(Mnemonic(bvcStr, RelativeModeStr), 2, p.bvc)
	opCodes[0x70] = id.Instruction(Mnemonic(bvsStr, RelativeModeStr), 2, p.bvs)
	opCodes[0x18] = id.Instruction(Mnemonic(clcStr, ImpliedModeStr), 2, p.clc)
	opCodes[0xD8] = id.Instruction(Mnemonic(cldStr, ImpliedModeStr), 2, p.cld)
	opCodes[0x58] = id.Instruction(Mnemonic(cliStr, ImpliedModeStr), 2, p.cli)
	opCodes[0xB8] = id.Instruction(Mnemonic(clvStr, ImpliedModeStr), 2, p.clv)
	opCodes[0xC9] = id.Instruction(Mnemonic(cmpStr, ImmediateModeStr), 2, p.cmp)
	opCodes[0xC5] = id.Instruction(Mnemonic(cmpStr, ZeropageModeStr), 3, p.cmp)
	opCodes[0xD5] = id.Instruction(Mnemonic(cmpStr, ZeropageXModeStr), 4, p.cmp)
	opCodes[0xCD] = id.Instruction(Mnemonic(cmpStr, AbsoluteModeStr), 4, p.cmp)
	opCodes[0xDD] = id.Instruction(Mnemonic(cmpStr, AbsoluteIndexedXModeStr), 4, p.cmp)
	opCodes[0xD9] = id.Instruction(Mnemonic(cmpStr, AbsoluteIndexedYModeStr), 4, p.cmp)
	opCodes[0xC1] = id.Instruction(Mnemonic(cmpStr, IndexedIndirectModeStr), 6, p.cmp)
	opCodes[0xD1] = id.Instruction(Mnemonic(cmpStr, IndirectIndexedModeStr), 5, p.cmp)
	opCodes[0xE0] = id.Instruction(Mnemonic(cpxStr, ImmediateModeStr), 2, p.cpx)
	opCodes[0xE4] = id.Instruction(Mnemonic(cpxStr, ZeropageModeStr), 3, p.cpx)
	opCodes[0xEC] = id.Instruction(Mnemonic(cpxStr, AbsoluteModeStr), 4, p.cpx)
	opCodes[0xC0] = id.Instruction(Mnemonic(cpyStr, ImmediateModeStr), 2, p.cpy)
	opCodes[0xC4] = id.Instruction(Mnemonic(cpyStr, ZeropageModeStr), 3, p.cpy)
	opCodes[0xCC] = id.Instruction(Mnemonic(cpyStr, AbsoluteModeStr), 4, p.cpy)
	opCodes[0xC6] = id.Instruction(Mnemonic(decStr, ZeropageModeStr), 5, p.dec)
	opCodes[0xD6] = id.Instruction(Mnemonic(decStr, ZeropageXModeStr), 6, p.dec)
	opCodes[0xCE] = id.Instruction(Mnemonic(decStr, AbsoluteModeStr), 6, p.dec)
	opCodes[0xDE] = id.Instruction(Mnemonic(decStr, AbsoluteIndexedXModeStr), 7, p.dec)
	opCodes[0xCA] = id.Instruction(Mnemonic(dexStr, ImpliedModeStr), 2, p.dex)
	opCodes[0x88] = id.Instruction(Mnemonic(deyStr, ImpliedModeStr), 2, p.dey)
	opCodes[0x49] = id.Instruction(Mnemonic(eorStr, ImmediateModeStr), 2, p.eor)
	opCodes[0x45] = id.Instruction(Mnemonic(eorStr, ZeropageModeStr), 3, p.eor)
	opCodes[0x55] = id.Instruction(Mnemonic(eorStr, ZeropageXModeStr), 4, p.eor)
	opCodes[0x4D] = id.Instruction(Mnemonic(eorStr, AbsoluteModeStr), 4, p.eor)
	opCodes[0x5D] = id.Instruction(Mnemonic(eorStr, AbsoluteIndexedXModeStr), 4, p.eor)
	opCodes[0x59] = id.Instruction(Mnemonic(eorStr, AbsoluteIndexedYModeStr), 4, p.eor)
	opCodes[0x41] = id.Instruction(Mnemonic(eorStr, IndexedIndirectModeStr), 6, p.eor)
	opCodes[0x51] = id.Instruction(Mnemonic(eorStr, IndirectIndexedModeStr), 5, p.eor)
	opCodes[0xE6] = id.Instruction(Mnemonic(incStr, ZeropageModeStr), 5, p.inc)
	opCodes[0xF6] = id.Instruction(Mnemonic(incStr, ZeropageXModeStr), 6, p.inc)
	opCodes[0xEE] = id.Instruction(Mnemonic(incStr, AbsoluteModeStr), 6, p.inc)
	opCodes[0xFE] = id.Instruction(Mnemonic(incStr, AbsoluteIndexedXModeStr), 7, p.inc)
	opCodes[0xE8] = id.Instruction(Mnemonic(inxStr, ImpliedModeStr), 2, p.inx)
	opCodes[0xC8] = id.Instruction(Mnemonic(inyStr, ImpliedModeStr), 2, p.iny)
	opCodes[0x4C] = id.Instruction(Mnemonic(jmpStr, AbsoluteModeStr), 3, p.jmp)
	opCodes[0x6C] = id.Instruction(Mnemonic(jmpStr, AbsoluteIndirectModeStr), 5, p.jmp)
	opCodes[0x20] = id.Instruction(Mnemonic(jsrStr, AbsoluteModeStr), 6, p.jsr)
	opCodes[0xA9] = id.Instruction(Mnemonic(ldaStr, ImmediateModeStr), 2, p.lda)
	opCodes[0xA5] = id.Instruction(Mnemonic(ldaStr, ZeropageModeStr), 3, p.lda)
	opCodes[0xB5] = id.Instruction(Mnemonic(ldaStr, ZeropageXModeStr), 4, p.lda)
	opCodes[0xAD] = id.Instruction(Mnemonic(ldaStr, AbsoluteModeStr), 4, p.lda)
	opCodes[0xBD] = id.Instruction(Mnemonic(ldaStr, AbsoluteIndexedXModeStr), 4, p.lda)
	opCodes[0xB9] = id.Instruction(Mnemonic(ldaStr, AbsoluteIndexedYModeStr), 4, p.lda)
	opCodes[0xA1] = id.Instruction(Mnemonic(ldaStr, IndexedIndirectModeStr), 6, p.lda)
	opCodes[0xB1] = id.Instruction(Mnemonic(ldaStr, IndirectIndexedModeStr), 5, p.lda)
	opCodes[0xA2] = id.Instruction(Mnemonic(ldxStr, ImmediateModeStr), 2, p.ldx)
	opCodes[0xA6] = id.Instruction(Mnemonic(ldxStr, ZeropageModeStr), 3, p.ldx)
	opCodes[0xB6] = id.Instruction(Mnemonic(ldxStr, ZeropageYModeStr), 4, p.ldx)
	opCodes[0xAE] = id.Instruction(Mnemonic(ldxStr, AbsoluteModeStr), 4, p.ldx)
	opCodes[0xBE] = id.Instruction(Mnemonic(ldxStr, AbsoluteIndexedXModeStr), 4, p.ldx)
	opCodes[0xA0] = id.Instruction(Mnemonic(ldyStr, ImmediateModeStr), 2, p.ldy)
	opCodes[0xA4] = id.Instruction(Mnemonic(ldyStr, ZeropageModeStr), 3, p.ldy)
	opCodes[0xB4] = id.Instruction(Mnemonic(ldyStr, ZeropageXModeStr), 4, p.ldy)
	opCodes[0xAC] = id.Instruction(Mnemonic(ldyStr, AbsoluteModeStr), 4, p.ldy)
	opCodes[0xBC] = id.Instruction(Mnemonic(ldyStr, AbsoluteIndexedXModeStr), 4, p.ldy)
	opCodes[0x4A] = id.Instruction(Mnemonic(lsrStr, AccumulatorModeStr), 2, p.lsr)
	opCodes[0x46] = id.Instruction(Mnemonic(lsrStr, ZeropageModeStr), 5, p.lsr)
	opCodes[0x56] = id.Instruction(Mnemonic(lsrStr, ZeropageXModeStr), 6, p.lsr)
	opCodes[0x4E] = id.Instruction(Mnemonic(lsrStr, AbsoluteModeStr), 6, p.lsr)
	opCodes[0x5E] = id.Instruction(Mnemonic(lsrStr, AbsoluteIndexedXModeStr), 7, p.lsr)
	opCodes[0xEA] = id.Instruction(Mnemonic(nopStr, ImpliedModeStr), 2, p.nop)
	opCodes[0x09] = id.Instruction(Mnemonic(oraStr, ImmediateModeStr), 2, p.ora)
	opCodes[0x05] = id.Instruction(Mnemonic(oraStr, ZeropageModeStr), 3, p.ora)
	opCodes[0x15] = id.Instruction(Mnemonic(oraStr, ZeropageXModeStr), 4, p.ora)
	opCodes[0x0D] = id.Instruction(Mnemonic(oraStr, AbsoluteModeStr), 4, p.ora)
	opCodes[0x1D] = id.Instruction(Mnemonic(oraStr, AbsoluteIndexedXModeStr), 4, p.ora)
	opCodes[0x19] = id.Instruction(Mnemonic(oraStr, AbsoluteIndexedYModeStr), 4, p.ora)
	opCodes[0x01] = id.Instruction(Mnemonic(oraStr, IndexedIndirectModeStr), 6, p.ora)
	opCodes[0x11] = id.Instruction(Mnemonic(oraStr, IndirectIndexedModeStr), 5, p.ora)
	opCodes[0x48] = id.Instruction(Mnemonic(phaStr, ImpliedModeStr), 3, p.pha)
	opCodes[0x68] = id.Instruction(Mnemonic(plaStr, ImpliedModeStr), 4, p.pla)
	opCodes[0x08] = id.Instruction(Mnemonic(phpStr, ImpliedModeStr), 3, p.php)
	opCodes[0x28] = id.Instruction(Mnemonic(plpStr, ImpliedModeStr), 4, p.plp)
	opCodes[0x2A] = id.Instruction(Mnemonic(rolStr, AccumulatorModeStr), 2, p.rol)
	opCodes[0x26] = id.Instruction(Mnemonic(rolStr, ZeropageModeStr), 5, p.rol)
	opCodes[0x36] = id.Instruction(Mnemonic(rolStr, ZeropageXModeStr), 6, p.rol)
	opCodes[0x2E] = id.Instruction(Mnemonic(rolStr, AbsoluteModeStr), 6, p.rol)
	opCodes[0x3E] = id.Instruction(Mnemonic(rolStr, AbsoluteIndexedXModeStr), 7, p.rol)
	opCodes[0x6A] = id.Instruction(Mnemonic(rorStr, AccumulatorModeStr), 2, p.ror)
	opCodes[0x66] = id.Instruction(Mnemonic(rorStr, ZeropageModeStr), 5, p.ror)
	opCodes[0x76] = id.Instruction(Mnemonic(rorStr, ZeropageXModeStr), 6, p.ror)
	opCodes[0x6E] = id.Instruction(Mnemonic(rorStr, AbsoluteModeStr), 6, p.ror)
	opCodes[0x7E] = id.Instruction(Mnemonic(rorStr, AbsoluteIndexedXModeStr), 7, p.ror)
	opCodes[0x40] = id.Instruction(Mnemonic(rtiStr, ImpliedModeStr), 6, p.rti)
	opCodes[0x60] = id.Instruction(Mnemonic(rtsStr, ImpliedModeStr), 6, p.rts)
	opCodes[0xE9] = id.Instruction(Mnemonic(sbcStr, ImmediateModeStr), 2, p.sbc)
	opCodes[0xE5] = id.Instruction(Mnemonic(sbcStr, ZeropageModeStr), 3, p.sbc)
	opCodes[0xF5] = id.Instruction(Mnemonic(sbcStr, ZeropageXModeStr), 4, p.sbc)
	opCodes[0xED] = id.Instruction(Mnemonic(sbcStr, AbsoluteModeStr), 4, p.sbc)
	opCodes[0xFD] = id.Instruction(Mnemonic(sbcStr, AbsoluteIndexedXModeStr), 4, p.sbc)
	opCodes[0xF9] = id.Instruction(Mnemonic(sbcStr, AbsoluteIndexedYModeStr), 4, p.sbc)
	opCodes[0xE1] = id.Instruction(Mnemonic(sbcStr, IndexedIndirectModeStr), 6, p.sbc)
	opCodes[0xF1] = id.Instruction(Mnemonic(sbcStr, IndirectIndexedModeStr), 5, p.sbc)
	opCodes[0x38] = id.Instruction(Mnemonic(secStr, ImpliedModeStr), 2, p.sec)
	opCodes[0xF8] = id.Instruction(Mnemonic(sedStr, ImpliedModeStr), 2, p.sed)
	opCodes[0x78] = id.Instruction(Mnemonic(seiStr, ImpliedModeStr), 2, p.sei)
	opCodes[0x85] = id.Instruction(Mnemonic(staStr, ZeropageModeStr), 3, p.sta)
	opCodes[0x95] = id.Instruction(Mnemonic(staStr, ZeropageXModeStr), 4, p.sta)
	opCodes[0x8D] = id.Instruction(Mnemonic(staStr, AbsoluteModeStr), 4, p.sta)
	opCodes[0x9D] = id.Instruction(Mnemonic(staStr, AbsoluteIndexedXModeStr), 5, p.sta)
	opCodes[0x99] = id.Instruction(Mnemonic(staStr, AbsoluteIndexedYModeStr), 5, p.sta)
	opCodes[0x81] = id.Instruction(Mnemonic(staStr, IndexedIndirectModeStr), 6, p.sta)
	opCodes[0x91] = id.Instruction(Mnemonic(staStr, IndirectIndexedModeStr), 6, p.sta)
	opCodes[0x86] = id.Instruction(Mnemonic(stxStr, ZeropageModeStr), 3, p.stx)
	opCodes[0x96] = id.Instruction(Mnemonic(stxStr, ZeropageYModeStr), 4, p.stx)
	opCodes[0x8E] = id.Instruction(Mnemonic(stxStr, AbsoluteModeStr), 4, p.stx)
	opCodes[0x84] = id.Instruction(Mnemonic(styStr, ZeropageModeStr), 3, p.sty)
	opCodes[0x94] = id.Instruction(Mnemonic(styStr, ZeropageXModeStr), 4, p.sty)
	opCodes[0x8C] = id.Instruction(Mnemonic(styStr, AbsoluteModeStr), 4, p.sty)
	opCodes[0xAA] = id.Instruction(Mnemonic(taxStr, ImpliedModeStr), 2, p.tax)
	opCodes[0xA8] = id.Instruction(Mnemonic(tayStr, ImpliedModeStr), 2, p.tay)
	opCodes[0xBA] = id.Instruction(Mnemonic(tsxStr, ImpliedModeStr), 2, p.tsx)
	opCodes[0x8A] = id.Instruction(Mnemonic(txaStr, ImpliedModeStr), 2, p.txa)
	opCodes[0x9A] = id.Instruction(Mnemonic(txsStr, ImpliedModeStr), 2, p.txs)
	opCodes[0x98] = id.Instruction(Mnemonic(tyaStr, ImpliedModeStr), 2, p.tya)
	return opCodes
}

func (p *CPU) addPCOffset(b byte) (uint16, bool) {
	pcLsb := uint16(p.Reg.PC & 0x00FF)
	offset := int8(b)
	pcLsb += uint16(offset)
	carryFlag := (pcLsb > uint16(0xFF))
	newPC := p.Reg.PC + uint16(offset)
	return newPC, carryFlag
}

func (p *CPU) adc(opcode OpCodeDef) InstructionFunc {
	load := opcode.AddressingMode.Load(p, false)

	var adcMode = map[BinaryOrDecimalMode]func(b byte){
		BCDMode: func(b byte) {
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
		BinaryMode: func(b byte) {
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
		adcMode[BinaryOrDecimalMode(p.Reg.IsSet(DecimalFlag))](b)
		p.Reg.SetZeroFlag(p.Reg.A)
		p.Reg.SetOverflowFlag(m, b, p.Reg.A, true)
		p.Reg.SetNegativeFlag(p.Reg.A)
		return true, nil
	}
}

func (p *CPU) and(opcode OpCodeDef) InstructionFunc {
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

func (p *CPU) asl(opcode OpCodeDef) InstructionFunc {
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

func (p *CPU) branch(opcode OpCodeDef, flag StatusFlag, state bool) InstructionFunc {
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
func (p *CPU) bcc(opcode OpCodeDef) InstructionFunc {
	return p.branch(opcode, CarryFlag, false)
}

func (p *CPU) bcs(opcode OpCodeDef) InstructionFunc {
	return p.branch(opcode, CarryFlag, true)
}

func (p *CPU) beq(opcode OpCodeDef) InstructionFunc {
	return p.branch(opcode, ZeroFlag, true)
}

func (p *CPU) bit(opcode OpCodeDef) InstructionFunc {
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

func (p *CPU) bmi(opcode OpCodeDef) InstructionFunc {
	return p.branch(opcode, NegativeFlag, true)
}

func (p *CPU) bne(opcode OpCodeDef) InstructionFunc {
	return p.branch(opcode, ZeroFlag, false)
}

func (p *CPU) bpl(opcode OpCodeDef) InstructionFunc {
	return p.branch(opcode, NegativeFlag, false)
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

func (p *CPU) brk(_ OpCodeDef) InstructionFunc {
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

func (p *CPU) bvc(opcode OpCodeDef) InstructionFunc {
	return p.branch(opcode, OverflowFlag, false)
}

func (p *CPU) bvs(opcode OpCodeDef) InstructionFunc {
	return p.branch(opcode, OverflowFlag, true)
}

func (p *CPU) clc(_ OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Reg.SetStatus(CarryFlag, false)
		return true, nil
	}
}

func (p *CPU) cld(_ OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Reg.SetStatus(DecimalFlag, false)
		return true, nil
	}
}

func (p *CPU) cli(_ OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Reg.SetStatus(InterruptDisableFlag, false)
		return true, nil
	}
}

func (p *CPU) clv(_ OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Reg.SetStatus(OverflowFlag, false)
		return true, nil
	}
}

func (p *CPU) compare(opcode OpCodeDef, regValue uint8) InstructionFunc {
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

func (p *CPU) cmp(opcode OpCodeDef) InstructionFunc {
	return p.compare(opcode, p.Reg.A)
}

func (p *CPU) cpx(opcode OpCodeDef) InstructionFunc {
	return p.compare(opcode, p.Reg.X)
}

func (p *CPU) cpy(opcode OpCodeDef) InstructionFunc {
	return p.compare(opcode, p.Reg.Y)
}

func (p *CPU) dec(opcode OpCodeDef) InstructionFunc {
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

func (p *CPU) dex(_ OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Reg.X--
		p.Reg.SetNegativeFlag(p.Reg.X)
		p.Reg.SetZeroFlag(p.Reg.X)
		return true, nil
	}
}

func (p *CPU) dey(_ OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Reg.Y--
		p.Reg.SetNegativeFlag(p.Reg.Y)
		p.Reg.SetZeroFlag(p.Reg.Y)
		return true, nil
	}
}

func (p *CPU) eor(opcode OpCodeDef) InstructionFunc {
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

func (p *CPU) inc(opcode OpCodeDef) InstructionFunc {
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

func (p *CPU) inx(_ OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Reg.X++
		p.Reg.SetNegativeFlag(p.Reg.X)
		p.Reg.SetZeroFlag(p.Reg.X)
		return true, nil
	}
}

func (p *CPU) iny(_ OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Reg.Y++
		p.Reg.SetNegativeFlag(p.Reg.Y)
		p.Reg.SetZeroFlag(p.Reg.Y)
		return true, nil
	}
}

func (p *CPU) jmp(opcode OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		address := opcode.AddressingMode.Address(p)
		p.Reg.PC = address
		return true, nil
	}
}

func (p *CPU) jsr(opcode OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Push(byte((p.Reg.PC & 0xFF00) >> 8))
		p.Push(byte(p.Reg.PC & 0x00FF))
		address := opcode.AddressingMode.Address(p)
		p.Reg.PC = address
		return true, nil
	}
}

func (p *CPU) lda(opcode OpCodeDef) InstructionFunc {
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

func (p *CPU) ldx(opcode OpCodeDef) InstructionFunc {
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

func (p *CPU) ldy(opcode OpCodeDef) InstructionFunc {
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

func (p *CPU) lsr(opcode OpCodeDef) InstructionFunc {
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

func (p *CPU) nop(_ OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		return true, nil
	}
}

func (p *CPU) ora(opcode OpCodeDef) InstructionFunc {
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

func (p *CPU) pha(_ OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Push(p.Reg.A)
		return true, nil
	}
}

func (p *CPU) php(_ OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Push(p.Reg.S)
		return true, nil
	}
}

func (p *CPU) pla(_ OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Reg.A = p.Pop()
		return true, nil
	}
}

func (p *CPU) plp(_ OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Reg.Status = p.Pop()
		return true, nil
	}
}

func (p *CPU) rol(opcode OpCodeDef) InstructionFunc {
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

func (p *CPU) ror(opcode OpCodeDef) InstructionFunc {
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

func (p *CPU) rti(_ OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Reg.Status = p.Pop()
		lowBytePC := p.Pop()
		hiBytePC := p.Pop()
		p.Reg.PC = (uint16(lowBytePC) | (uint16(hiBytePC) << 8))
		p.Reg.SetStatus(BreakFlag, false)
		return true, nil
	}
}

func (p *CPU) rts(_ OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		lowBytePC := p.Pop()
		hiBytePC := p.Pop()
		p.Reg.PC = (uint16(lowBytePC) | (uint16(hiBytePC) << 8))
		return true, nil
	}
}

func (p *CPU) sbc(opcode OpCodeDef) InstructionFunc {
	load := opcode.AddressingMode.Load(p, false)

	var sbcMode = map[BinaryOrDecimalMode]func(b byte){
		BCDMode: func(b byte) {
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

		BinaryMode: func(b byte) {
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
		sbcMode[BinaryOrDecimalMode(p.Reg.IsSet(DecimalFlag))](b)
		p.Reg.SetZeroFlag(p.Reg.A)
		p.Reg.SetNegativeFlag(p.Reg.A)
		p.Reg.SetOverflowFlag(m, b, p.Reg.A, false)
		return true, nil
	}
}

func (p *CPU) sec(_ OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Reg.SetStatus(CarryFlag, true)
		return true, nil
	}
}

func (p *CPU) sed(_ OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Reg.SetStatus(DecimalFlag, true)
		return true, nil
	}
}

func (p *CPU) sei(_ OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Reg.SetStatus(InterruptDisableFlag, true)
		return true, nil
	}
}

func (p *CPU) sta(opcode OpCodeDef) InstructionFunc {
	store := opcode.AddressingMode.Store(p, true)
	return func() (Completed, error) {
		store(p.Reg.A)
		return true, nil
	}
}

func (p *CPU) stx(opcode OpCodeDef) InstructionFunc {
	store := opcode.AddressingMode.Store(p, true)
	return func() (Completed, error) {
		store(p.Reg.X)
		return true, nil
	}
}

func (p *CPU) sty(opcode OpCodeDef) InstructionFunc {
	store := opcode.AddressingMode.Store(p, true)
	return func() (Completed, error) {
		store(p.Reg.Y)
		return true, nil
	}
}

func (p *CPU) tax(_ OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Reg.X = p.Reg.A
		p.Reg.SetZeroFlag(p.Reg.X)
		p.Reg.SetNegativeFlag(p.Reg.X)
		return true, nil
	}
}

func (p *CPU) tay(_ OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Reg.Y = p.Reg.A
		p.Reg.SetZeroFlag(p.Reg.Y)
		p.Reg.SetNegativeFlag(p.Reg.Y)
		return true, nil
	}
}

func (p *CPU) tsx(_ OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Reg.X = p.Reg.S
		p.Reg.SetZeroFlag(p.Reg.X)
		p.Reg.SetNegativeFlag(p.Reg.X)
		return true, nil
	}
}

func (p *CPU) txa(_ OpCodeDef) InstructionFunc {
	// LoadAm, _ := am()
	return func() (Completed, error) {
		p.Reg.A = p.Reg.X
		p.Reg.SetZeroFlag(p.Reg.A)
		p.Reg.SetNegativeFlag(p.Reg.A)
		return true, nil
	}
}

func (p *CPU) txs(_ OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Reg.S = p.Reg.X
		p.Reg.SetZeroFlag(p.Reg.S)
		p.Reg.SetNegativeFlag(p.Reg.S)
		return true, nil
	}
}

func (p *CPU) tya(_ OpCodeDef) InstructionFunc {
	return func() (Completed, error) {
		p.Reg.A = p.Reg.Y
		p.Reg.SetZeroFlag(p.Reg.A)
		p.Reg.SetNegativeFlag(p.Reg.A)
		return true, nil
	}
}
