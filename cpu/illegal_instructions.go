package cpu

// addIllegalOpCodes populates undocumented/illegal 6502 opcodes commonly used on the C64 (6510).
// It mutates the provided opcode table in place.
func addIllegalOpCodes(p *CPU, opCodes []*OpCodeDef) {
	id := NewInstruction(getAddressingMode)

	// Helpers implementing illegal ops inline using existing addressing mode loaders/stores
	lax := func(op OpCodeDef) InstructionFunc {
		load := op.AddressingMode.Load(p, false)
		return func() (Completed, error) {
			b, completed := load()
			if !completed {
				return false, nil
			}
			p.Reg.A = b
			p.Reg.X = b
			p.Reg.SetZeroFlag(b)
			p.Reg.SetNegativeFlag(b)
			return true, nil
		}
	}

	sax := func(op OpCodeDef) InstructionFunc {
		store := op.AddressingMode.Store(p, true)
		return func() (Completed, error) {
			store(p.Reg.A & p.Reg.X)
			return true, nil
		}
	}

	// Read-Modify-Write combos
	slo := func(op OpCodeDef) InstructionFunc { // ASL mem, then ORA
		load := op.AddressingMode.Load(p, true)
		store := op.AddressingMode.Store(p, true)
		return func() (Completed, error) {
			b, _ := load()
			res := b << 1
			p.Reg.SetCarryFlag(b, res)
			store(res)
			p.Reg.A = p.Reg.A | res
			p.Reg.SetZeroFlag(p.Reg.A)
			p.Reg.SetNegativeFlag(p.Reg.A)
			return true, nil
		}
	}

	rla := func(op OpCodeDef) InstructionFunc { // ROL mem, then AND
		load := op.AddressingMode.Load(p, true)
		store := op.AddressingMode.Store(p, true)
		return func() (Completed, error) {
			b, _ := load()
			carry := (b & 0x80) != 0
			res := (b << 1)
			if p.Reg.IsSet(CarryFlag) {
				res |= 0x01
			}
			p.Reg.SetStatus(CarryFlag, carry)
			store(res)
			p.Reg.A = p.Reg.A & res
			p.Reg.SetZeroFlag(p.Reg.A)
			p.Reg.SetNegativeFlag(p.Reg.A)
			return true, nil
		}
	}

	sre := func(op OpCodeDef) InstructionFunc { // LSR mem, then EOR
		load := op.AddressingMode.Load(p, true)
		store := op.AddressingMode.Store(p, true)
		return func() (Completed, error) {
			b, _ := load()
			bit0 := b & 0x01
			res := b >> 1
			p.Reg.SetStatus(CarryFlag, bit0 == 1)
			store(res)
			p.Reg.A = p.Reg.A ^ res
			p.Reg.SetZeroFlag(p.Reg.A)
			p.Reg.SetNegativeFlag(p.Reg.A)
			return true, nil
		}
	}

	rra := func(op OpCodeDef) InstructionFunc { // ROR mem, then ADC
		load := op.AddressingMode.Load(p, true)
		store := op.AddressingMode.Store(p, true)
		return func() (Completed, error) {
			b, _ := load()
			carryOut := (b & 0x01) != 0
			res := b >> 1
			if p.Reg.IsSet(CarryFlag) {
				res |= 0x80
			}
			p.Reg.SetStatus(CarryFlag, carryOut)
			store(res)

			// ADC in binary (most emulators ignore BCD here)
			m := p.Reg.A
			r := m + res
			if p.Reg.IsSet(CarryFlag) {
				r++
			}
			p.Reg.A = r
			p.Reg.SetCarryFlag(m, p.Reg.A)
			p.Reg.SetZeroFlag(p.Reg.A)
			p.Reg.SetOverflowFlag(m, res, p.Reg.A, true)
			p.Reg.SetNegativeFlag(p.Reg.A)
			return true, nil
		}
	}

	dcp := func(op OpCodeDef) InstructionFunc { // DEC mem, then CMP
		load := op.AddressingMode.Load(p, true)
		store := op.AddressingMode.Store(p, true)
		return func() (Completed, error) {
			b, _ := load()
			b--
			store(b)
			// CMP A, b
			result := p.Reg.A - b
			p.Reg.SetStatus(ZeroFlag, result == 0)
			p.Reg.SetStatus(NegativeFlag, (result&0x80) != 0)
			p.Reg.SetStatus(CarryFlag, p.Reg.A >= b)
			return true, nil
		}
	}

	isc := func(op OpCodeDef) InstructionFunc { // INC mem, then SBC
		load := op.AddressingMode.Load(p, true)
		store := op.AddressingMode.Store(p, true)
		return func() (Completed, error) {
			b, _ := load()
			b++
			store(b)
			// SBC in binary
			m := p.Reg.A
			c := byte(0)
			if p.Reg.IsSet(CarryFlag) {
				c = 1
			}
			r := m - b - (1 - c)
			p.Reg.A = r
			p.Reg.SetStatus(CarryFlag, m >= (b+(1-c)))
			p.Reg.SetZeroFlag(p.Reg.A)
			p.Reg.SetNegativeFlag(p.Reg.A)
			p.Reg.SetOverflowFlag(m, b, p.Reg.A, false)
			return true, nil
		}
	}

	anc := func(op OpCodeDef) InstructionFunc { // AND imm, C = bit7
		load := op.AddressingMode.Load(p, false)
		return func() (Completed, error) {
			b, completed := load()
			if !completed {
				return false, nil
			}
			p.Reg.A = p.Reg.A & b
			p.Reg.SetZeroFlag(p.Reg.A)
			p.Reg.SetNegativeFlag(p.Reg.A)
			p.Reg.SetStatus(CarryFlag, (p.Reg.A&0x80) != 0)
			return true, nil
		}
	}

	alr := func(op OpCodeDef) InstructionFunc { // AND imm, then LSR A
		load := op.AddressingMode.Load(p, false)
		return func() (Completed, error) {
			b, completed := load()
			if !completed {
				return false, nil
			}
			m := p.Reg.A & b
			carry := (m & 0x01) != 0
			m = m >> 1
			p.Reg.A = m
			p.Reg.SetStatus(CarryFlag, carry)
			p.Reg.SetZeroFlag(p.Reg.A)
			p.Reg.SetNegativeFlag(p.Reg.A)
			return true, nil
		}
	}

	arr := func(op OpCodeDef) InstructionFunc { // AND imm, then ROR A (approx flags)
		load := op.AddressingMode.Load(p, false)
		return func() (Completed, error) {
			b, completed := load()
			if !completed {
				return false, nil
			}
			m := p.Reg.A & b
			carryIn := byte(0)
			if p.Reg.IsSet(CarryFlag) {
				carryIn = 0x80
			}
			res := (m >> 1) | carryIn
			p.Reg.A = res
			// C approximated as bit 6 of m
			p.Reg.SetStatus(CarryFlag, (m&0x40) != 0)
			// V approximated from bits 5^6 of result
			v := ((res >> 5) & 1) ^ ((res >> 6) & 1)
			p.Reg.SetStatus(OverflowFlag, v == 1)
			p.Reg.SetZeroFlag(p.Reg.A)
			p.Reg.SetNegativeFlag(p.Reg.A)
			return true, nil
		}
	}

	xaa := func(op OpCodeDef) InstructionFunc { // A = X & imm (very unstable on real HW)
		load := op.AddressingMode.Load(p, false)
		return func() (Completed, error) {
			b, completed := load()
			if !completed {
				return false, nil
			}
			p.Reg.A = p.Reg.X & b
			p.Reg.SetZeroFlag(p.Reg.A)
			p.Reg.SetNegativeFlag(p.Reg.A)
			return true, nil
		}
	}

	// LAX
	opCodes[0xA7] = id.Instruction(Mnemonic("LAX", ZeropageModeStr), 3, lax)
	opCodes[0xB7] = id.Instruction(Mnemonic("LAX", ZeropageYModeStr), 4, lax)
	opCodes[0xAF] = id.Instruction(Mnemonic("LAX", AbsoluteModeStr), 4, lax)
	opCodes[0xBF] = id.Instruction(Mnemonic("LAX", AbsoluteIndexedYModeStr), 4, lax)
	opCodes[0xA3] = id.Instruction(Mnemonic("LAX", IndexedIndirectModeStr), 6, lax)
	opCodes[0xB3] = id.Instruction(Mnemonic("LAX", IndirectIndexedModeStr), 5, lax)

	// SAX (AXS/AAX)
	opCodes[0x87] = id.Instruction(Mnemonic("SAX", ZeropageModeStr), 3, sax)
	opCodes[0x97] = id.Instruction(Mnemonic("SAX", ZeropageYModeStr), 4, sax)
	opCodes[0x8F] = id.Instruction(Mnemonic("SAX", AbsoluteModeStr), 4, sax)
	opCodes[0x83] = id.Instruction(Mnemonic("SAX", IndexedIndirectModeStr), 6, sax)

	// SLO, RLA, SRE, RRA
	opCodes[0x07] = id.Instruction(Mnemonic("SLO", ZeropageModeStr), 5, slo)
	opCodes[0x17] = id.Instruction(Mnemonic("SLO", ZeropageXModeStr), 6, slo)
	opCodes[0x0F] = id.Instruction(Mnemonic("SLO", AbsoluteModeStr), 6, slo)
	opCodes[0x1F] = id.Instruction(Mnemonic("SLO", AbsoluteIndexedXModeStr), 7, slo)
	opCodes[0x1B] = id.Instruction(Mnemonic("SLO", AbsoluteIndexedYModeStr), 7, slo)
	opCodes[0x03] = id.Instruction(Mnemonic("SLO", IndexedIndirectModeStr), 8, slo)
	opCodes[0x13] = id.Instruction(Mnemonic("SLO", IndirectIndexedModeStr), 8, slo)

	opCodes[0x27] = id.Instruction(Mnemonic("RLA", ZeropageModeStr), 5, rla)
	opCodes[0x37] = id.Instruction(Mnemonic("RLA", ZeropageXModeStr), 6, rla)
	opCodes[0x2F] = id.Instruction(Mnemonic("RLA", AbsoluteModeStr), 6, rla)
	opCodes[0x3F] = id.Instruction(Mnemonic("RLA", AbsoluteIndexedXModeStr), 7, rla)
	opCodes[0x3B] = id.Instruction(Mnemonic("RLA", AbsoluteIndexedYModeStr), 7, rla)
	opCodes[0x23] = id.Instruction(Mnemonic("RLA", IndexedIndirectModeStr), 8, rla)
	opCodes[0x33] = id.Instruction(Mnemonic("RLA", IndirectIndexedModeStr), 8, rla)

	opCodes[0x47] = id.Instruction(Mnemonic("SRE", ZeropageModeStr), 5, sre)
	opCodes[0x57] = id.Instruction(Mnemonic("SRE", ZeropageXModeStr), 6, sre)
	opCodes[0x4F] = id.Instruction(Mnemonic("SRE", AbsoluteModeStr), 6, sre)
	opCodes[0x5F] = id.Instruction(Mnemonic("SRE", AbsoluteIndexedXModeStr), 7, sre)
	opCodes[0x5B] = id.Instruction(Mnemonic("SRE", AbsoluteIndexedYModeStr), 7, sre)
	opCodes[0x43] = id.Instruction(Mnemonic("SRE", IndexedIndirectModeStr), 8, sre)
	opCodes[0x53] = id.Instruction(Mnemonic("SRE", IndirectIndexedModeStr), 8, sre)

	opCodes[0x67] = id.Instruction(Mnemonic("RRA", ZeropageModeStr), 5, rra)
	opCodes[0x77] = id.Instruction(Mnemonic("RRA", ZeropageXModeStr), 6, rra)
	opCodes[0x6F] = id.Instruction(Mnemonic("RRA", AbsoluteModeStr), 6, rra)
	opCodes[0x7F] = id.Instruction(Mnemonic("RRA", AbsoluteIndexedXModeStr), 7, rra)
	opCodes[0x7B] = id.Instruction(Mnemonic("RRA", AbsoluteIndexedYModeStr), 7, rra)
	opCodes[0x63] = id.Instruction(Mnemonic("RRA", IndexedIndirectModeStr), 8, rra)
	opCodes[0x73] = id.Instruction(Mnemonic("RRA", IndirectIndexedModeStr), 8, rra)

	// DCP, ISC
	opCodes[0xC7] = id.Instruction(Mnemonic("DCP", ZeropageModeStr), 5, dcp)
	opCodes[0xD7] = id.Instruction(Mnemonic("DCP", ZeropageXModeStr), 6, dcp)
	opCodes[0xCF] = id.Instruction(Mnemonic("DCP", AbsoluteModeStr), 6, dcp)
	opCodes[0xDF] = id.Instruction(Mnemonic("DCP", AbsoluteIndexedXModeStr), 7, dcp)
	opCodes[0xDB] = id.Instruction(Mnemonic("DCP", AbsoluteIndexedYModeStr), 7, dcp)
	opCodes[0xC3] = id.Instruction(Mnemonic("DCP", IndexedIndirectModeStr), 8, dcp)
	opCodes[0xD3] = id.Instruction(Mnemonic("DCP", IndirectIndexedModeStr), 8, dcp)

	opCodes[0xE7] = id.Instruction(Mnemonic("ISC", ZeropageModeStr), 5, isc)
	opCodes[0xF7] = id.Instruction(Mnemonic("ISC", ZeropageXModeStr), 6, isc)
	opCodes[0xEF] = id.Instruction(Mnemonic("ISC", AbsoluteModeStr), 6, isc)
	opCodes[0xFF] = id.Instruction(Mnemonic("ISC", AbsoluteIndexedXModeStr), 7, isc)
	opCodes[0xFB] = id.Instruction(Mnemonic("ISC", AbsoluteIndexedYModeStr), 7, isc)
	opCodes[0xE3] = id.Instruction(Mnemonic("ISC", IndexedIndirectModeStr), 8, isc)
	opCodes[0xF3] = id.Instruction(Mnemonic("ISC", IndirectIndexedModeStr), 8, isc)

	// ANC, ALR, ARR, XAA, and SBC duplicate
	opCodes[0x0B] = id.Instruction(Mnemonic("ANC", ImmediateModeStr), 2, anc)
	opCodes[0x2B] = id.Instruction(Mnemonic("ANC", ImmediateModeStr), 2, anc)
	opCodes[0x4B] = id.Instruction(Mnemonic("ALR", ImmediateModeStr), 2, alr)
	opCodes[0x6B] = id.Instruction(Mnemonic("ARR", ImmediateModeStr), 2, arr)
	opCodes[0x8B] = id.Instruction(Mnemonic("XAA", ImmediateModeStr), 2, xaa)

	// SBC immediate alias (illegal) - keep distinct mnemonic to avoid overriding standard SBC # at 0xE9
	opCodes[0xEB] = id.Instruction(Mnemonic("SBC*", ImmediateModeStr), 2, p.sbc)

	// Multi-byte NOPs and variants (use distinct mnemonics to avoid overriding canonical NOP 0xEA)
	// Single-byte NOP variants (implied)
	for _, opc := range []byte{0x1A, 0x3A, 0x5A, 0x7A, 0xDA, 0xFA} {
		opCodes[opc] = id.Instruction(Mnemonic("NOP*", ImpliedModeStr), 2, p.nop)
	}
	// Two-byte NOPs (aka DOP) with immediate operand
	for _, opc := range []byte{0x80, 0x82, 0xC2, 0xE2} {
		opCodes[opc] = id.Instruction(Mnemonic("DOP", ImmediateModeStr), 2, p.nop)
	}
	// Three-byte NOPs on ABS
	opCodes[0x0C] = id.Instruction(Mnemonic("TOP", AbsoluteModeStr), 4, p.nop)
	// Three-byte NOPs on ABS,X (aka TOP)
	for _, opc := range []byte{0x1C, 0x3C, 0x5C, 0x7C, 0xDC, 0xFC} {
		opCodes[opc] = id.Instruction(Mnemonic("TOP", AbsoluteIndexedXModeStr), 4, p.nop)
	}
	// Zero-page and Zero-page,X NOP-like
	for _, opc := range []byte{0x04, 0x44, 0x64} {
		opCodes[opc] = id.Instruction(Mnemonic("SKB", ZeropageModeStr), 3, p.nop)
	}
	for _, opc := range []byte{0x14, 0x34, 0x54, 0x74, 0xD4, 0xF4} {
		opCodes[opc] = id.Instruction(Mnemonic("SKW", ZeropageXModeStr), 4, p.nop)
	}
}
