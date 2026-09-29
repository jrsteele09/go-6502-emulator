package cpu

func loadOperand(p *CPU, opcode *OpCodeDef, ignoreExtraCycle bool) (byte, Completed) {
	return opcode.AddressingMode.Load(p, ignoreExtraCycle)
}

func storeOperand(p *CPU, opcode *OpCodeDef, value byte) Completed {
	return opcode.AddressingMode.Store(p, value, true)
}

func adc(p *CPU, opcode *OpCodeDef) (Completed, error) {
	b, completed := loadOperand(p, opcode, false)
	if !completed {
		return false, nil
	}
	m := p.Reg.A
	carry := uint16(0)
	if p.Reg.IsSet(CarryFlag) {
		carry = 1
	}
	binarySum := uint16(m) + uint16(b) + carry
	flagResult := byte(binarySum)
	if p.Reg.IsSet(DecimalFlag) {
		lowSum := uint16(p.Reg.A&0x0F) + uint16(b&0x0F) + carry
		highSum := uint16(p.Reg.A&0xF0) + uint16(b&0xF0)
		if lowSum > 0x09 {
			lowSum += 0x06
			highSum += 0x10
		}
		flagResult = byte(highSum | (lowSum & 0x0F))
		sum := highSum | (lowSum & 0x0F)
		if highSum >= 0xA0 {
			sum += 0x60
		}
		p.Reg.A = byte(sum)
		p.Reg.SetStatus(CarryFlag, sum > 0xFF)
	} else {
		p.Reg.A = byte(binarySum)
		p.Reg.SetStatus(CarryFlag, binarySum > 0xFF)
	}
	p.Reg.SetZeroFlag(byte(binarySum))
	p.Reg.SetOverflowFlag(m, b, flagResult, true)
	p.Reg.SetNegativeFlag(flagResult)
	return true, nil
}

func and(p *CPU, opcode *OpCodeDef) (Completed, error) {
	b, completed := loadOperand(p, opcode, false)
	if !completed {
		return false, nil
	}
	p.Reg.A &= b
	p.Reg.SetNegativeFlag(p.Reg.A)
	p.Reg.SetZeroFlag(p.Reg.A)
	return true, nil
}

func asl(p *CPU, opcode *OpCodeDef) (Completed, error) {
	b, _ := loadOperand(p, opcode, true)
	result := b << 1
	p.Reg.SetNegativeFlag(result)
	p.Reg.SetZeroFlag(result)
	p.Reg.SetCarryFlag(b, result)
	storeOperand(p, opcode, result)
	return true, nil
}

func branch(p *CPU, opcode *OpCodeDef, flag StatusFlag, state bool) (Completed, error) {
	switch p.execute.phase {
	case 1:
		p.execute.phase = 2
		if p.execute.pageCrossed {
			return false, nil
		}
		p.Reg.PC = p.execute.address
		return true, nil
	case 2:
		p.Reg.PC = p.execute.address
		return true, nil
	}
	if p.Reg.IsSet(flag) != state {
		return true, nil
	}
	offset, complete := loadOperand(p, opcode, true)
	if !complete {
		return false, nil
	}
	p.execute.address, p.execute.pageCrossed = p.addPCOffset(offset)
	p.execute.phase = 1
	return false, nil
}

func bcc(p *CPU, op *OpCodeDef) (Completed, error) { return branch(p, op, CarryFlag, false) }
func bcs(p *CPU, op *OpCodeDef) (Completed, error) { return branch(p, op, CarryFlag, true) }
func beq(p *CPU, op *OpCodeDef) (Completed, error) { return branch(p, op, ZeroFlag, true) }
func bmi(p *CPU, op *OpCodeDef) (Completed, error) { return branch(p, op, NegativeFlag, true) }
func bne(p *CPU, op *OpCodeDef) (Completed, error) { return branch(p, op, ZeroFlag, false) }
func bpl(p *CPU, op *OpCodeDef) (Completed, error) { return branch(p, op, NegativeFlag, false) }
func bvc(p *CPU, op *OpCodeDef) (Completed, error) { return branch(p, op, OverflowFlag, false) }
func bvs(p *CPU, op *OpCodeDef) (Completed, error) { return branch(p, op, OverflowFlag, true) }

func bit(p *CPU, opcode *OpCodeDef) (Completed, error) {
	b, completed := loadOperand(p, opcode, false)
	if !completed {
		return false, nil
	}
	p.Reg.SetStatus(NegativeFlag, b&0x80 != 0)
	p.Reg.SetStatus(OverflowFlag, b&0x40 != 0)
	p.Reg.SetZeroFlag(b & p.Reg.A)
	return true, nil
}

func brk(p *CPU, _ *OpCodeDef) (Completed, error) {
	p.Reg.PC++
	p.Push(byte(p.Reg.PC >> 8))
	p.Push(byte(p.Reg.PC))
	p.Push(p.Reg.StatusForStack(true))
	p.Reg.SetStatus(InterruptDisableFlag, true)
	low := p.mem.Read(irqVector)
	high := p.mem.Read(irqVector + 1)
	p.Reg.PC = uint16(high)<<8 | uint16(low)
	return true, nil
}

func clc(p *CPU, _ *OpCodeDef) (Completed, error) {
	p.Reg.SetStatus(CarryFlag, false)
	return true, nil
}
func cld(p *CPU, _ *OpCodeDef) (Completed, error) {
	p.Reg.SetStatus(DecimalFlag, false)
	return true, nil
}
func cli(p *CPU, _ *OpCodeDef) (Completed, error) {
	p.Reg.SetStatus(InterruptDisableFlag, false)
	return true, nil
}
func clv(p *CPU, _ *OpCodeDef) (Completed, error) {
	p.Reg.SetStatus(OverflowFlag, false)
	return true, nil
}

func compare(p *CPU, opcode *OpCodeDef, register RegisterEnum) (Completed, error) {
	b, completed := loadOperand(p, opcode, false)
	if !completed {
		return false, nil
	}
	var value byte
	switch register {
	case AReg:
		value = p.Reg.A
	case XReg:
		value = p.Reg.X
	case YReg:
		value = p.Reg.Y
	case SReg:
		value = p.Reg.S
	case PCReg:
		value = byte(p.Reg.PC)
	case StatusReg:
		value = p.Reg.Status
	}
	result := value - b
	p.Reg.SetStatus(ZeroFlag, result == 0)
	p.Reg.SetStatus(NegativeFlag, result&0x80 != 0)
	p.Reg.SetStatus(CarryFlag, value >= b)
	return true, nil
}

func cmp(p *CPU, op *OpCodeDef) (Completed, error) { return compare(p, op, AReg) }
func cpx(p *CPU, op *OpCodeDef) (Completed, error) { return compare(p, op, XReg) }
func cpy(p *CPU, op *OpCodeDef) (Completed, error) { return compare(p, op, YReg) }

func dec(p *CPU, opcode *OpCodeDef) (Completed, error) {
	b, _ := loadOperand(p, opcode, true)
	b--
	p.Reg.SetNegativeFlag(b)
	p.Reg.SetZeroFlag(b)
	storeOperand(p, opcode, b)
	return true, nil
}
func dex(p *CPU, _ *OpCodeDef) (Completed, error) {
	p.Reg.X--
	p.Reg.SetNegativeFlag(p.Reg.X)
	p.Reg.SetZeroFlag(p.Reg.X)
	return true, nil
}
func dey(p *CPU, _ *OpCodeDef) (Completed, error) {
	p.Reg.Y--
	p.Reg.SetNegativeFlag(p.Reg.Y)
	p.Reg.SetZeroFlag(p.Reg.Y)
	return true, nil
}

func eor(p *CPU, opcode *OpCodeDef) (Completed, error) {
	b, completed := loadOperand(p, opcode, false)
	if !completed {
		return false, nil
	}
	p.Reg.A ^= b
	p.Reg.SetNegativeFlag(p.Reg.A)
	p.Reg.SetZeroFlag(p.Reg.A)
	return true, nil
}
func inc(p *CPU, opcode *OpCodeDef) (Completed, error) {
	b, _ := loadOperand(p, opcode, true)
	b++
	p.Reg.SetNegativeFlag(b)
	p.Reg.SetZeroFlag(b)
	storeOperand(p, opcode, b)
	return true, nil
}
func inx(p *CPU, _ *OpCodeDef) (Completed, error) {
	p.Reg.X++
	p.Reg.SetNegativeFlag(p.Reg.X)
	p.Reg.SetZeroFlag(p.Reg.X)
	return true, nil
}
func iny(p *CPU, _ *OpCodeDef) (Completed, error) {
	p.Reg.Y++
	p.Reg.SetNegativeFlag(p.Reg.Y)
	p.Reg.SetZeroFlag(p.Reg.Y)
	return true, nil
}

func jmp(p *CPU, opcode *OpCodeDef) (Completed, error) {
	p.Reg.PC = opcode.AddressingMode.Address(p)
	return true, nil
}
func jsr(p *CPU, opcode *OpCodeDef) (Completed, error) {
	returnAddress := p.Reg.PC - 1
	p.Push(byte(returnAddress >> 8))
	p.Push(byte(returnAddress))
	p.Reg.PC = opcode.AddressingMode.Address(p)
	return true, nil
}

func lda(p *CPU, opcode *OpCodeDef) (Completed, error) {
	b, completed := loadOperand(p, opcode, false)
	if !completed {
		return false, nil
	}
	p.Reg.A = b
	p.Reg.SetNegativeFlag(b)
	p.Reg.SetZeroFlag(b)
	return true, nil
}
func ldx(p *CPU, opcode *OpCodeDef) (Completed, error) {
	b, completed := loadOperand(p, opcode, false)
	if !completed {
		return false, nil
	}
	p.Reg.X = b
	p.Reg.SetNegativeFlag(b)
	p.Reg.SetZeroFlag(b)
	return true, nil
}
func ldy(p *CPU, opcode *OpCodeDef) (Completed, error) {
	b, completed := loadOperand(p, opcode, false)
	if !completed {
		return false, nil
	}
	p.Reg.Y = b
	p.Reg.SetNegativeFlag(b)
	p.Reg.SetZeroFlag(b)
	return true, nil
}

func lsr(p *CPU, opcode *OpCodeDef) (Completed, error) {
	b, completed := loadOperand(p, opcode, true)
	if !completed {
		return false, nil
	}
	carry := b&1 != 0
	b >>= 1
	p.Reg.SetZeroFlag(b)
	p.Reg.SetStatus(NegativeFlag, false)
	p.Reg.SetStatus(CarryFlag, carry)
	storeOperand(p, opcode, b)
	return true, nil
}

func nop(p *CPU, _ *OpCodeDef) (Completed, error) { return true, nil }

func ora(p *CPU, opcode *OpCodeDef) (Completed, error) {
	b, completed := loadOperand(p, opcode, false)
	if !completed {
		return false, nil
	}
	p.Reg.A |= b
	p.Reg.SetNegativeFlag(p.Reg.A)
	p.Reg.SetZeroFlag(p.Reg.A)
	return true, nil
}

func pha(p *CPU, _ *OpCodeDef) (Completed, error) { p.Push(p.Reg.A); return true, nil }
func php(p *CPU, _ *OpCodeDef) (Completed, error) {
	p.Push(p.Reg.StatusForStack(true))
	return true, nil
}
func pla(p *CPU, _ *OpCodeDef) (Completed, error) {
	p.Reg.A = p.Pop()
	p.Reg.SetZeroFlag(p.Reg.A)
	p.Reg.SetNegativeFlag(p.Reg.A)
	return true, nil
}
func plp(p *CPU, _ *OpCodeDef) (Completed, error) {
	p.Reg.RestoreStatus(p.Pop())
	return true, nil
}

func rol(p *CPU, opcode *OpCodeDef) (Completed, error) {
	b, completed := loadOperand(p, opcode, true)
	if !completed {
		return false, nil
	}
	carry := b&0x80 != 0
	b <<= 1
	if p.Reg.IsSet(CarryFlag) {
		b |= 1
	}
	p.Reg.SetStatus(CarryFlag, carry)
	p.Reg.SetNegativeFlag(b)
	p.Reg.SetZeroFlag(b)
	storeOperand(p, opcode, b)
	return true, nil
}
func ror(p *CPU, opcode *OpCodeDef) (Completed, error) {
	b, completed := loadOperand(p, opcode, true)
	if !completed {
		return false, nil
	}
	carry := b&1 != 0
	b >>= 1
	if p.Reg.IsSet(CarryFlag) {
		b |= 0x80
	}
	p.Reg.SetStatus(CarryFlag, carry)
	p.Reg.SetNegativeFlag(b)
	p.Reg.SetZeroFlag(b)
	storeOperand(p, opcode, b)
	return true, nil
}

func rti(p *CPU, _ *OpCodeDef) (Completed, error) {
	p.Reg.RestoreStatus(p.Pop())
	low, high := p.Pop(), p.Pop()
	p.Reg.PC = uint16(low) | uint16(high)<<8
	return true, nil
}
func rts(p *CPU, _ *OpCodeDef) (Completed, error) {
	low, high := p.Pop(), p.Pop()
	p.Reg.PC = (uint16(low) | uint16(high)<<8) + 1
	return true, nil
}

func sbc(p *CPU, opcode *OpCodeDef) (Completed, error) {
	b, completed := loadOperand(p, opcode, false)
	if !completed {
		return false, nil
	}
	m := p.Reg.A
	carry := 0
	if p.Reg.IsSet(CarryFlag) {
		carry = 1
	}
	diff := int(m) - int(b) - (1 - carry)
	flagResult := byte(diff)
	if p.Reg.IsSet(DecimalFlag) {
		result := diff
		if int(p.Reg.A&0x0F)-int(b&0x0F)-(1-carry) < 0 {
			result -= 0x06
		}
		if diff < 0 {
			result -= 0x60
		}
		p.Reg.A = byte(result)
		p.Reg.SetStatus(CarryFlag, diff >= 0)
	} else {
		p.Reg.A = flagResult
		p.Reg.SetStatus(CarryFlag, diff >= 0)
	}
	p.Reg.SetZeroFlag(flagResult)
	p.Reg.SetNegativeFlag(flagResult)
	p.Reg.SetOverflowFlag(m, b, flagResult, false)
	return true, nil
}

func sec(p *CPU, _ *OpCodeDef) (Completed, error) {
	p.Reg.SetStatus(CarryFlag, true)
	return true, nil
}
func sed(p *CPU, _ *OpCodeDef) (Completed, error) {
	p.Reg.SetStatus(DecimalFlag, true)
	return true, nil
}
func sei(p *CPU, _ *OpCodeDef) (Completed, error) {
	p.Reg.SetStatus(InterruptDisableFlag, true)
	return true, nil
}

func sta(p *CPU, op *OpCodeDef) (Completed, error) { storeOperand(p, op, p.Reg.A); return true, nil }
func stx(p *CPU, op *OpCodeDef) (Completed, error) { storeOperand(p, op, p.Reg.X); return true, nil }
func sty(p *CPU, op *OpCodeDef) (Completed, error) { storeOperand(p, op, p.Reg.Y); return true, nil }

func tax(p *CPU, _ *OpCodeDef) (Completed, error) {
	p.Reg.X = p.Reg.A
	p.Reg.SetZeroFlag(p.Reg.X)
	p.Reg.SetNegativeFlag(p.Reg.X)
	return true, nil
}
func tay(p *CPU, _ *OpCodeDef) (Completed, error) {
	p.Reg.Y = p.Reg.A
	p.Reg.SetZeroFlag(p.Reg.Y)
	p.Reg.SetNegativeFlag(p.Reg.Y)
	return true, nil
}
func tsx(p *CPU, _ *OpCodeDef) (Completed, error) {
	p.Reg.X = p.Reg.S
	p.Reg.SetZeroFlag(p.Reg.X)
	p.Reg.SetNegativeFlag(p.Reg.X)
	return true, nil
}
func txa(p *CPU, _ *OpCodeDef) (Completed, error) {
	p.Reg.A = p.Reg.X
	p.Reg.SetZeroFlag(p.Reg.A)
	p.Reg.SetNegativeFlag(p.Reg.A)
	return true, nil
}
func txs(p *CPU, _ *OpCodeDef) (Completed, error) { p.Reg.S = p.Reg.X; return true, nil }
func tya(p *CPU, _ *OpCodeDef) (Completed, error) {
	p.Reg.A = p.Reg.Y
	p.Reg.SetZeroFlag(p.Reg.A)
	p.Reg.SetNegativeFlag(p.Reg.A)
	return true, nil
}
