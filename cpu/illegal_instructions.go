package cpu

// addIllegalOpCodes populates undocumented/illegal 6502 opcodes commonly used on the C64 (6510).
// It mutates the provided opcode table in place.
func addIllegalOpCodes(p *CPU) {
	id := NewInstruction(getAddressingMode)

	// LAX
	p.opCodes[0xA7] = id.Instruction(Mnemonic("LAX", ZeropageModeStr), 3, lax)
	p.opCodes[0xB7] = id.Instruction(Mnemonic("LAX", ZeropageYModeStr), 4, lax)
	p.opCodes[0xAF] = id.Instruction(Mnemonic("LAX", AbsoluteModeStr), 4, lax)
	p.opCodes[0xBF] = id.Instruction(Mnemonic("LAX", AbsoluteIndexedYModeStr), 4, lax)
	p.opCodes[0xA3] = id.Instruction(Mnemonic("LAX", IndexedIndirectModeStr), 6, lax)
	p.opCodes[0xB3] = id.Instruction(Mnemonic("LAX", IndirectIndexedModeStr), 5, lax)

	// SAX (AXS/AAX)
	p.opCodes[0x87] = id.Instruction(Mnemonic("SAX", ZeropageModeStr), 3, sax)
	p.opCodes[0x97] = id.Instruction(Mnemonic("SAX", ZeropageYModeStr), 4, sax)
	p.opCodes[0x8F] = id.Instruction(Mnemonic("SAX", AbsoluteModeStr), 4, sax)
	p.opCodes[0x83] = id.Instruction(Mnemonic("SAX", IndexedIndirectModeStr), 6, sax)

	// SLO, RLA, SRE, RRA
	p.opCodes[0x07] = id.Instruction(Mnemonic("SLO", ZeropageModeStr), 5, slo)
	p.opCodes[0x17] = id.Instruction(Mnemonic("SLO", ZeropageXModeStr), 6, slo)
	p.opCodes[0x0F] = id.Instruction(Mnemonic("SLO", AbsoluteModeStr), 6, slo)
	p.opCodes[0x1F] = id.Instruction(Mnemonic("SLO", AbsoluteIndexedXModeStr), 7, slo)
	p.opCodes[0x1B] = id.Instruction(Mnemonic("SLO", AbsoluteIndexedYModeStr), 7, slo)
	p.opCodes[0x03] = id.Instruction(Mnemonic("SLO", IndexedIndirectModeStr), 8, slo)
	p.opCodes[0x13] = id.Instruction(Mnemonic("SLO", IndirectIndexedModeStr), 8, slo)

	p.opCodes[0x27] = id.Instruction(Mnemonic("RLA", ZeropageModeStr), 5, rla)
	p.opCodes[0x37] = id.Instruction(Mnemonic("RLA", ZeropageXModeStr), 6, rla)
	p.opCodes[0x2F] = id.Instruction(Mnemonic("RLA", AbsoluteModeStr), 6, rla)
	p.opCodes[0x3F] = id.Instruction(Mnemonic("RLA", AbsoluteIndexedXModeStr), 7, rla)
	p.opCodes[0x3B] = id.Instruction(Mnemonic("RLA", AbsoluteIndexedYModeStr), 7, rla)
	p.opCodes[0x23] = id.Instruction(Mnemonic("RLA", IndexedIndirectModeStr), 8, rla)
	p.opCodes[0x33] = id.Instruction(Mnemonic("RLA", IndirectIndexedModeStr), 8, rla)

	p.opCodes[0x47] = id.Instruction(Mnemonic("SRE", ZeropageModeStr), 5, sre)
	p.opCodes[0x57] = id.Instruction(Mnemonic("SRE", ZeropageXModeStr), 6, sre)
	p.opCodes[0x4F] = id.Instruction(Mnemonic("SRE", AbsoluteModeStr), 6, sre)
	p.opCodes[0x5F] = id.Instruction(Mnemonic("SRE", AbsoluteIndexedXModeStr), 7, sre)
	p.opCodes[0x5B] = id.Instruction(Mnemonic("SRE", AbsoluteIndexedYModeStr), 7, sre)
	p.opCodes[0x43] = id.Instruction(Mnemonic("SRE", IndexedIndirectModeStr), 8, sre)
	p.opCodes[0x53] = id.Instruction(Mnemonic("SRE", IndirectIndexedModeStr), 8, sre)

	p.opCodes[0x67] = id.Instruction(Mnemonic("RRA", ZeropageModeStr), 5, rra)
	p.opCodes[0x77] = id.Instruction(Mnemonic("RRA", ZeropageXModeStr), 6, rra)
	p.opCodes[0x6F] = id.Instruction(Mnemonic("RRA", AbsoluteModeStr), 6, rra)
	p.opCodes[0x7F] = id.Instruction(Mnemonic("RRA", AbsoluteIndexedXModeStr), 7, rra)
	p.opCodes[0x7B] = id.Instruction(Mnemonic("RRA", AbsoluteIndexedYModeStr), 7, rra)
	p.opCodes[0x63] = id.Instruction(Mnemonic("RRA", IndexedIndirectModeStr), 8, rra)
	p.opCodes[0x73] = id.Instruction(Mnemonic("RRA", IndirectIndexedModeStr), 8, rra)

	// DCP, ISC
	p.opCodes[0xC7] = id.Instruction(Mnemonic("DCP", ZeropageModeStr), 5, dcp)
	p.opCodes[0xD7] = id.Instruction(Mnemonic("DCP", ZeropageXModeStr), 6, dcp)
	p.opCodes[0xCF] = id.Instruction(Mnemonic("DCP", AbsoluteModeStr), 6, dcp)
	p.opCodes[0xDF] = id.Instruction(Mnemonic("DCP", AbsoluteIndexedXModeStr), 7, dcp)
	p.opCodes[0xDB] = id.Instruction(Mnemonic("DCP", AbsoluteIndexedYModeStr), 7, dcp)
	p.opCodes[0xC3] = id.Instruction(Mnemonic("DCP", IndexedIndirectModeStr), 8, dcp)
	p.opCodes[0xD3] = id.Instruction(Mnemonic("DCP", IndirectIndexedModeStr), 8, dcp)

	p.opCodes[0xE7] = id.Instruction(Mnemonic("ISC", ZeropageModeStr), 5, isc)
	p.opCodes[0xF7] = id.Instruction(Mnemonic("ISC", ZeropageXModeStr), 6, isc)
	p.opCodes[0xEF] = id.Instruction(Mnemonic("ISC", AbsoluteModeStr), 6, isc)
	p.opCodes[0xFF] = id.Instruction(Mnemonic("ISC", AbsoluteIndexedXModeStr), 7, isc)
	p.opCodes[0xFB] = id.Instruction(Mnemonic("ISC", AbsoluteIndexedYModeStr), 7, isc)
	p.opCodes[0xE3] = id.Instruction(Mnemonic("ISC", IndexedIndirectModeStr), 8, isc)
	p.opCodes[0xF3] = id.Instruction(Mnemonic("ISC", IndirectIndexedModeStr), 8, isc)

	// ANC, ALR, ARR, XAA, and SBC duplicate
	p.opCodes[0x0B] = id.Instruction(Mnemonic("ANC", ImmediateModeStr), 2, anc)
	p.opCodes[0x2B] = id.Instruction(Mnemonic("ANC", ImmediateModeStr), 2, anc)
	p.opCodes[0x4B] = id.Instruction(Mnemonic("ALR", ImmediateModeStr), 2, alr)
	p.opCodes[0x6B] = id.Instruction(Mnemonic("ARR", ImmediateModeStr), 2, arr)
	p.opCodes[0x8B] = id.Instruction(Mnemonic("XAA", ImmediateModeStr), 2, xaa)

	// SBC immediate alias (illegal) - keep distinct mnemonic to avoid overriding standard SBC # at 0xE9
	p.opCodes[0xEB] = id.Instruction(Mnemonic("SBC*", ImmediateModeStr), 2, sbc)

	// Multi-byte NOPs and variants (use distinct mnemonics to avoid overriding canonical NOP 0xEA)
	// Single-byte NOP variants (implied)
	for _, opc := range []byte{0x1A, 0x3A, 0x5A, 0x7A, 0xDA, 0xFA} {
		p.opCodes[opc] = id.Instruction(Mnemonic("NOP*", ImpliedModeStr), 2, nop)
	}
	// Two-byte NOPs (aka DOP) with immediate operand
	for _, opc := range []byte{0x80, 0x82, 0xC2, 0xE2} {
		p.opCodes[opc] = id.Instruction(Mnemonic("DOP", ImmediateModeStr), 2, nop)
	}
	// Three-byte NOPs on ABS
	p.opCodes[0x0C] = id.Instruction(Mnemonic("TOP", AbsoluteModeStr), 4, nop)
	// Three-byte NOPs on ABS,X (aka TOP)
	for _, opc := range []byte{0x1C, 0x3C, 0x5C, 0x7C, 0xDC, 0xFC} {
		p.opCodes[opc] = id.Instruction(Mnemonic("TOP", AbsoluteIndexedXModeStr), 4, nop)
	}
	// Zero-page and Zero-page,X NOP-like
	for _, opc := range []byte{0x04, 0x44, 0x64} {
		p.opCodes[opc] = id.Instruction(Mnemonic("SKB", ZeropageModeStr), 3, nop)
	}
	for _, opc := range []byte{0x14, 0x34, 0x54, 0x74, 0xD4, 0xF4} {
		p.opCodes[opc] = id.Instruction(Mnemonic("SKW", ZeropageXModeStr), 4, nop)
	}
}
