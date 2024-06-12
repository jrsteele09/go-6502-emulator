package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCarryFlag(t *testing.T) {
	regs := NewRegisters()
	regs.SetCarryFlag(0x80, 0x1)
	assert.Equal(t, true, regs.IsSet(CarryFlag))
	regs.SetCarryFlag(0x7F, 0x1)
	assert.Equal(t, true, regs.IsSet(CarryFlag))
	regs.SetCarryFlag(0x7E, 0x80)
	assert.Equal(t, false, regs.IsSet(CarryFlag))
}

func TestZeroFlag(t *testing.T) {
	regs := NewRegisters()
	regs.SetZeroFlag(0x00)
	assert.Equal(t, true, regs.IsSet(ZeroFlag))
	regs.SetZeroFlag(0x01)
	assert.Equal(t, false, regs.IsSet(ZeroFlag))
}

func TestNegativeFlag(t *testing.T) {
	regs := NewRegisters()
	regs.SetNegativeFlag(0xFF)
	assert.Equal(t, true, regs.IsSet(NegativeFlag))
	regs.SetNegativeFlag(0x01)
	assert.Equal(t, false, regs.IsSet(NegativeFlag))
}

func TestOverFlowFlag(t *testing.T) {
	regs := NewRegisters()
	var m byte = 127
	var n byte = 1
	result := m + n
	regs.SetOverflowFlag(m, n, result, true)
	assert.Equal(t, true, regs.IsSet(OverflowFlag))

	m = 128
	n = 1
	result = m + n
	regs.SetOverflowFlag(m, n, result, true)
	assert.Equal(t, false, regs.IsSet(OverflowFlag))

	m = 255
	n = 1
	result = m + n
	regs.SetOverflowFlag(m, n, result, true)
	assert.Equal(t, false, regs.IsSet(OverflowFlag))

	m = 80
	n = 80
	result = m + n
	regs.SetOverflowFlag(m, n, result, true)
	assert.Equal(t, true, regs.IsSet(OverflowFlag))

	m = 0
	n = 80
	result = m - n
	regs.SetOverflowFlag(m, n, result, true)
	assert.Equal(t, true, regs.IsSet(OverflowFlag))

	m = 10
	n = 10
	result = m - n
	regs.SetOverflowFlag(m, n, result, true)
	assert.Equal(t, false, regs.IsSet(OverflowFlag))

}
