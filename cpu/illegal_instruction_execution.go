package cpu

func lax(p *CPU, op *OpCodeDef) (Completed, error) {
	b, completed := loadOperand(p, op, false)
	if !completed {
		return false, nil
	}
	p.Reg.A, p.Reg.X = b, b
	p.Reg.SetZeroFlag(b)
	p.Reg.SetNegativeFlag(b)
	return true, nil
}
func sax(p *CPU, op *OpCodeDef) (Completed, error) {
	storeOperand(p, op, p.Reg.A&p.Reg.X)
	return true, nil
}

func slo(p *CPU, op *OpCodeDef) (Completed, error) {
	b, _ := loadOperand(p, op, true)
	result := b << 1
	p.Reg.SetCarryFlag(b, result)
	storeOperand(p, op, result)
	p.Reg.A |= result
	p.Reg.SetZeroFlag(p.Reg.A)
	p.Reg.SetNegativeFlag(p.Reg.A)
	return true, nil
}
func rla(p *CPU, op *OpCodeDef) (Completed, error) {
	b, _ := loadOperand(p, op, true)
	carry := b&0x80 != 0
	result := b << 1
	if p.Reg.IsSet(CarryFlag) {
		result |= 1
	}
	p.Reg.SetStatus(CarryFlag, carry)
	storeOperand(p, op, result)
	p.Reg.A &= result
	p.Reg.SetZeroFlag(p.Reg.A)
	p.Reg.SetNegativeFlag(p.Reg.A)
	return true, nil
}
func sre(p *CPU, op *OpCodeDef) (Completed, error) {
	b, _ := loadOperand(p, op, true)
	carry := b&1 != 0
	result := b >> 1
	p.Reg.SetStatus(CarryFlag, carry)
	storeOperand(p, op, result)
	p.Reg.A ^= result
	p.Reg.SetZeroFlag(p.Reg.A)
	p.Reg.SetNegativeFlag(p.Reg.A)
	return true, nil
}
func rra(p *CPU, op *OpCodeDef) (Completed, error) {
	b, _ := loadOperand(p, op, true)
	carryOut := b&1 != 0
	result := b >> 1
	if p.Reg.IsSet(CarryFlag) {
		result |= 0x80
	}
	p.Reg.SetStatus(CarryFlag, carryOut)
	storeOperand(p, op, result)
	m := p.Reg.A
	r := m + result
	if p.Reg.IsSet(CarryFlag) {
		r++
	}
	p.Reg.A = r
	p.Reg.SetCarryFlag(m, r)
	p.Reg.SetZeroFlag(r)
	p.Reg.SetOverflowFlag(m, result, r, true)
	p.Reg.SetNegativeFlag(r)
	return true, nil
}
func dcp(p *CPU, op *OpCodeDef) (Completed, error) {
	b, _ := loadOperand(p, op, true)
	b--
	storeOperand(p, op, b)
	result := p.Reg.A - b
	p.Reg.SetStatus(ZeroFlag, result == 0)
	p.Reg.SetStatus(NegativeFlag, result&0x80 != 0)
	p.Reg.SetStatus(CarryFlag, p.Reg.A >= b)
	return true, nil
}
func isc(p *CPU, op *OpCodeDef) (Completed, error) {
	b, _ := loadOperand(p, op, true)
	b++
	storeOperand(p, op, b)
	m, carry := p.Reg.A, byte(0)
	if p.Reg.IsSet(CarryFlag) {
		carry = 1
	}
	r := m - b - (1 - carry)
	p.Reg.A = r
	p.Reg.SetStatus(CarryFlag, m >= b+(1-carry))
	p.Reg.SetZeroFlag(r)
	p.Reg.SetNegativeFlag(r)
	p.Reg.SetOverflowFlag(m, b, r, false)
	return true, nil
}
func anc(p *CPU, op *OpCodeDef) (Completed, error) {
	b, completed := loadOperand(p, op, false)
	if !completed {
		return false, nil
	}
	p.Reg.A &= b
	p.Reg.SetZeroFlag(p.Reg.A)
	p.Reg.SetNegativeFlag(p.Reg.A)
	p.Reg.SetStatus(CarryFlag, p.Reg.A&0x80 != 0)
	return true, nil
}
func alr(p *CPU, op *OpCodeDef) (Completed, error) {
	b, completed := loadOperand(p, op, false)
	if !completed {
		return false, nil
	}
	result := p.Reg.A & b
	carry := result&1 != 0
	result >>= 1
	p.Reg.A = result
	p.Reg.SetStatus(CarryFlag, carry)
	p.Reg.SetZeroFlag(result)
	p.Reg.SetNegativeFlag(result)
	return true, nil
}
func arr(p *CPU, op *OpCodeDef) (Completed, error) {
	b, completed := loadOperand(p, op, false)
	if !completed {
		return false, nil
	}
	masked := p.Reg.A & b
	carryIn := byte(0)
	if p.Reg.IsSet(CarryFlag) {
		carryIn = 0x80
	}
	result := masked>>1 | carryIn
	p.Reg.A = result
	p.Reg.SetStatus(CarryFlag, masked&0x40 != 0)
	v := result>>5&1 ^ result>>6&1
	p.Reg.SetStatus(OverflowFlag, v == 1)
	p.Reg.SetZeroFlag(result)
	p.Reg.SetNegativeFlag(result)
	return true, nil
}
func xaa(p *CPU, op *OpCodeDef) (Completed, error) {
	b, completed := loadOperand(p, op, false)
	if !completed {
		return false, nil
	}
	p.Reg.A = p.Reg.X & b
	p.Reg.SetZeroFlag(p.Reg.A)
	p.Reg.SetNegativeFlag(p.Reg.A)
	return true, nil
}
