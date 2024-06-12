package cpu

import (
	"testing"

	"github.com/jrsteele09/go-6502-emulator/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type addressSize = uint16

type setupFunc = func(p *Cpu) int
type assertionFunc = func(t *testing.T, p *Cpu, name string)

type InstructionTest struct {
	testName string
	setup    setupFunc
	assert   assertionFunc
}

const startAddress = uint16(0xD000)
const stackAddress = uint16(0x01FF)

func executeTests(t *testing.T, tests []InstructionTest) {
	for _, test := range tests {
		m := memory.NewMemory[uint16](64 * 1024)
		cpu := NewCpu(m)
		noOfOps := test.setup(cpu)
		cpu.Reg.PC = startAddress

		for i := 0; i < noOfOps; i++ {
			var complete Completed
			for !complete {
				c, err := cpu.Execute()
				require.NoError(t, err)
				complete = c
			}
		}
		test.assert(t, cpu, test.testName)
	}
}

func TestForDebugging(t *testing.T) {
	m := memory.NewMemory[uint16](64 * 1024)
	cpu := NewCpu(m)
	cpu.Reg.PC = startAddress

	cpu.mem.Write(startAddress, 0x91, 0x05)
	cpu.mem.Write(0x0005, 0x01, 0x50)
	cpu.Reg.A = 0xff
	cpu.Reg.Y = 0x05

	var complete Completed
	for !complete {
		c, err := cpu.Execute()
		require.NoError(t, err)
		complete = c
	}

	assert.Equal(t, byte(0xff), cpu.mem.Read(0x5006))
	assert.Equal(t, uint64(6), cpu.cycles)
}

func TestNmiInterruptHandling(t *testing.T) {
	m := memory.NewMemory[uint16](64 * 1024)
	cpu := NewCpu(m)
	cpu.Reg.PC = startAddress

	cpu.mem.Write(startAddress, 0xA9, 0x05)
	cpu.mem.Write(nmiVector, 0xAD, 0xDE)

	var complete Completed
	setInterrupt := true
	// Execute Instruction and then set interrupt
	for !complete {
		c, err := cpu.Execute()
		if setInterrupt {
			cpu.Nmi()
			setInterrupt = false
		}
		require.NoError(t, err)
		complete = c
	}

	// Execute the interrupt handler
	complete = false
	for !complete {
		c, err := cpu.Execute()
		require.NoError(t, err)
		complete = c
	}

	assert.Equal(t, uint16(0xDEAD), cpu.Reg.PC)
	assert.Equal(t, uint8(0xD0), cpu.mem.Read(stackAddress))
	assert.Equal(t, uint8(0x02), cpu.mem.Read(stackAddress-1))
	assert.Equal(t, true, cpu.Reg.IsSet(InterruptDisableFlag))
	assert.Equal(t, false, cpu.nmi)
	assert.Equal(t, false, cpu.irq)
}

func TestIrqInterruptHandling(t *testing.T) {
	m := memory.NewMemory[addressSize](64 * 1024)
	cpu := NewCpu(m)
	cpu.Reg.PC = startAddress

	cpu.mem.Write(startAddress, 0xA9, 0x05)
	cpu.mem.Write(irqVector, 0xAD, 0xDE)

	var complete Completed
	setInterrupt := true
	cpu.Reg.SetStatus(InterruptDisableFlag, false)
	// Execute Instruction and then set interrupt
	for !complete {
		c, err := cpu.Execute()
		if setInterrupt {
			cpu.Irq()
			setInterrupt = false
		}
		require.NoError(t, err)
		complete = c
	}

	// Execute the interrupt handler
	complete = false
	for !complete {
		c, err := cpu.Execute()
		require.NoError(t, err)
		complete = c
	}

	assert.Equal(t, uint16(0xDEAD), cpu.Reg.PC)
	assert.Equal(t, uint8(0xD0), cpu.mem.Read(stackAddress))
	assert.Equal(t, uint8(0x02), cpu.mem.Read(stackAddress-1))
	assert.Equal(t, true, cpu.Reg.IsSet(InterruptDisableFlag))
	assert.Equal(t, false, cpu.nmi)
	assert.Equal(t, false, cpu.irq)
}

func TestADC(t *testing.T) {
	var tests = []InstructionTest{
		{"TestADCImmediate", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x69, 0x01)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, false, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(OverflowFlag), name)
		}},
		{"TestADCImmediateWithCarry", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x69, 0x03, 0x69, 0xFF)
			return 2
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x02), p.Reg.A, name)
			assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, uint64(4), p.cycles, name)
		}},

		{"TestADCImmediateNegativeFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x69, 0x80)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x80), p.Reg.A)
			assert.Equal(t, uint64(2), p.cycles)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag))
		}},
		{"TestADCImmediateZeroFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x69, 0x00)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x00), p.Reg.A, name)
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(ZeroFlag), name)
		}},
		{"TestADCImmediateCarryFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x69, 0x01)
			p.Reg.A = 0xFF
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x00), p.Reg.A, name)
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
		}},
		{"TestADCImmediateOverFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x69, 0x02)
			p.Reg.A = 127
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x81), p.Reg.A, name)
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(OverflowFlag), name)
		}},
		{"TestADCImmediateOverFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x69, 0x02)
			p.Reg.A = 127
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x81), p.Reg.A, name)
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(OverflowFlag), name)
		}},

		{"TestADCZeropage", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x65, 0x80)
			p.mem.Write(0x0080, 01)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(3), p.cycles, name)
		}},
		{"TestADCZeropageX", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x75, 0x80)
			p.mem.Write(0x0081, 01)
			p.Reg.X = 0x01
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(4), p.cycles, name)
		}},
		{"TestADCAbsolute", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x6D, 0x80, 0x00)
			p.mem.Write(0x0080, 01)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(4), p.cycles, name)
		}},
		{"TestADCAbsoluteX", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x7D, 0x80, 0x00)
			p.mem.Write(0x0081, 01)
			p.Reg.X = 0x01
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(4), p.cycles, name)
		}},
		{"TestADCAbsoluteXPageOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x7D, 0x01, 0x00)
			p.mem.Write(0x0100, 01)
			p.Reg.X = 0xFF
			p.Reg.A = 0x01
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x02), p.Reg.A, name)
			assert.Equal(t, uint64(5), p.cycles, name)
		}},
		{"TestADCAbsoluteY", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x79, 0x80, 0x00)
			p.mem.Write(0x0081, 01)
			p.Reg.Y = 0x01
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(4), p.cycles, name)

		}},
		{"TestADCAbsoluteYPageOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x79, 0x01, 0x00)
			p.mem.Write(0x0100, 01)
			p.Reg.Y = 0xFF
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A)
			assert.Equal(t, uint64(5), p.cycles)
		}},
		{"TestADCIndexedIndirect", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x61, 0x05)
			p.mem.Write(0x000A, 0x10, 0x50)
			p.mem.Write(0x5010, 0x01)
			p.Reg.X = 0x05
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(6), p.cycles, name)
		}},
		{"TestADCIndexedIndirectPageOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x61, 0x01)
			p.mem.Write(0x0000, 0x10, 0x50)
			p.mem.Write(0x5010, 0x01)
			p.Reg.X = 0xFF
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint8(1), p.Reg.A, name)
			assert.Equal(t, uint64(6), p.cycles, name)
		}},
		{"TestADCIndirectIndexed", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x71, 0x05) // ADC ($0x05),Y
			p.mem.Write(0x0005, 0x10, 0x50)
			p.mem.Write(0x5015, 0x01)
			p.Reg.Y = 0x05
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(5), p.cycles, name)
		}},
		{"TestADCIndirectIndexedPageOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x71, 0x05)
			p.mem.Write(0x0005, 0x01, 0x50)
			p.mem.Write(0x5100, 0x01)
			p.Reg.Y = 0xFF
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(6), p.cycles, name)
		}},
		{"TestADCDecimalMode", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x69, 5)
			p.Reg.SetStatus(DecimalFlag, true)
			p.Reg.A = 9
			p.Reg.SetStatus(CarryFlag, false)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint8(0x14), p.Reg.A, name)
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, false, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(OverflowFlag), name)
		}},
		{"TestADCDecimalMode", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x69, 0x39)
			p.Reg.SetStatus(DecimalFlag, true)
			p.Reg.A = 0x49
			p.Reg.SetStatus(CarryFlag, false)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x88), p.Reg.A, name)
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, true, p.Reg.IsSet(OverflowFlag), name)
		}},
		{"TestADCDecimalMode", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x69, 0x49)
			p.Reg.SetStatus(DecimalFlag, true)
			p.Reg.A = 0x69
			p.Reg.SetStatus(CarryFlag, false)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x18), p.Reg.A, name)
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, false, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(OverflowFlag), name)
		}},
		{"TestADCDecimalMode", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x69, 0x08)
			p.Reg.SetStatus(DecimalFlag, true)
			p.Reg.A = 0x06
			p.Reg.SetStatus(CarryFlag, false)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x14), p.Reg.A, name)
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, false, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(OverflowFlag), name)
		}},
		{"TestADCDecimalMode", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x69, 0x01)
			p.Reg.SetStatus(DecimalFlag, true)
			p.Reg.A = 0x99
			p.Reg.SetStatus(CarryFlag, false)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x00), p.Reg.A, name)
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, false, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, true, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(OverflowFlag), name)
		}},
	}
	executeTests(t, tests)
}

func TestAND(t *testing.T) {
	var tests = []InstructionTest{
		{"TestANDImmediate", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x29, 0x01)
			p.Reg.A = 0xFF
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, false, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(OverflowFlag), name)
		}},
		{"TestANDImmediateNegativeFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x29, 0x80)
			p.Reg.A = 0x80
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x80), p.Reg.A)
			assert.Equal(t, uint64(2), p.cycles)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag))
		}},
		{"TestANDImmediateZeroFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x29, 0x00)
			p.Reg.A = 0xFF
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x00), p.Reg.A, name)
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(ZeroFlag), name)
		}},
		{"TestANDZeropage", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x25, 0x80)
			p.mem.Write(0x0080, 01)
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(3), p.cycles, name)
		}},
		{"TestANDZeropageX", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x35, 0x80)
			p.mem.Write(0x0081, 01)
			p.Reg.X = 0x01
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(4), p.cycles, name)
		}},
		{"TestANDAbsolute", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x2d, 0x80, 0x00)
			p.mem.Write(0x0080, 01)
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(4), p.cycles, name)
		}},
		{"TestANDAbsoluteX", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x3d, 0x80, 0x00)
			p.mem.Write(0x0081, 01)
			p.Reg.X = 0x01
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(4), p.cycles, name)
		}},
		{"TestANDAbsoluteXPageOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x3d, 0x01, 0x00)
			p.mem.Write(0x0100, 01)
			p.Reg.X = 0xFF
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(5), p.cycles, name)
		}},
		{"TestANDAbsoluteY", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x39, 0x80, 0x00)
			p.mem.Write(0x0081, 01)
			p.Reg.Y = 0x01
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(4), p.cycles, name)

		}},
		{"TestANDAbsoluteYPageOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x39, 0x01, 0x00)
			p.mem.Write(0x0100, 01)
			p.Reg.A = 0xff
			p.Reg.Y = 0xFF
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(5), p.cycles, name)
		}},
		{"TestANDIndexedIndirect", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x21, 0x05)
			p.mem.Write(0x000A, 0x10, 0x50)
			p.mem.Write(0x5010, 0x01)
			p.Reg.X = 0x05
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(6), p.cycles, name)
		}},
		{"TestANDIndexedIndirectPageOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x21, 0x01)
			p.mem.Write(0x0000, 0x10, 0x50)
			p.mem.Write(0x5010, 0x01)
			p.Reg.A = 0xff
			p.Reg.X = 0xFF
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(6), p.cycles, name)
		}},
		{"TestANDIndirectIndexed", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x31, 0x05)
			p.mem.Write(0x0005, 0x10, 0x50)
			p.mem.Write(0x5015, 0x01)
			p.Reg.A = 0xff
			p.Reg.Y = 0x05
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(5), p.cycles, name)
		}},
		{"TestANDIndirectIndexedPageOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x31, 0x05)
			p.mem.Write(0x0005, 0x01, 0x50)
			p.mem.Write(0x5100, 0x01)
			p.Reg.A = 0xff
			p.Reg.Y = 0xFF
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(6), p.cycles, name)
		}},
	}
	executeTests(t, tests)
}

func TestASL(t *testing.T) {
	var tests = []InstructionTest{
		{"TestASLAccumulator", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x0A)
			p.Reg.A = 0x7F
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0xFE), p.Reg.A, name)
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
		}},
		{"TestASLZeropage", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x06, 0x01)
			p.mem.Write(0x0000, 0x00, 0x7F)
			return 1

		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0xFE), p.mem.Read(0x0001), name)
			assert.Equal(t, uint64(5), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
		}},
		{"TestASLZeropageX", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x16, 0x00)
			p.mem.Write(0x0000, 0x00, 0x7F)
			p.Reg.X = 0x01
			return 1

		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0xFE), p.mem.Read(0x0001), name)
			assert.Equal(t, uint64(6), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
		}},
		{"TestASLAbsolute", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x0E, 0x21, 0xD0)
			p.mem.Write(0xD020, 0x00, 0x7F)
			return 1

		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0xFE), p.mem.Read(0xD021), name+" result")
			assert.Equal(t, uint64(6), p.cycles, name+" clock cycles")
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name+" NegativeFlag")
		}},
		{"TestASLAbsoluteX", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x1E, 0x20, 0xD0)
			p.mem.Write(0xD020, 0x00, 0x00, 0x00, 0x7F)
			p.Reg.X = 3
			return 1

		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0xFE), p.mem.Read(0xD023), name+" result")
			assert.Equal(t, uint64(7), p.cycles, name+" clock cycles")
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name+" NegativeFlag")
		}},
	}
	executeTests(t, tests)
}

func TestBCC(t *testing.T) {
	var tests = []InstructionTest{
		{"TestBCC -1", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x90, 0xFE)
			p.Reg.SetStatus(CarryFlag, false)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint16(0xD000), p.Reg.PC, name)
			assert.Equal(t, uint64(2), p.cycles, name)
		}},
		{"TestBCC +5", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x90, 0x05)
			p.Reg.SetStatus(CarryFlag, false)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint16(0xD007), p.Reg.PC, name)
			assert.Equal(t, uint64(2), p.cycles, name)
		}},
		{"TestBCC -5", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x90, 0xf9)
			p.Reg.SetStatus(CarryFlag, false)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint16(0xcffb), p.Reg.PC, name)
			assert.Equal(t, uint64(3), p.cycles, name)
		}},
	}
	executeTests(t, tests)
}

func TestBCS(t *testing.T) {
	var tests = []InstructionTest{
		{"TestBCS -1", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xB0, 0xFE)
			p.Reg.SetStatus(CarryFlag, true)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint16(0xD000), p.Reg.PC, name)
			assert.Equal(t, uint64(2), p.cycles, name)
		}},
		{"TestBCS +5", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xB0, 0x05)
			p.Reg.SetStatus(CarryFlag, true)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint16(0xD007), p.Reg.PC, name)
			assert.Equal(t, uint64(2), p.cycles, name)
		}},
		{"TestBCS -5", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xB0, 0xf9)
			p.Reg.SetStatus(CarryFlag, true)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint16(0xcffb), p.Reg.PC, name)
			assert.Equal(t, uint64(3), p.cycles, name)
		}},
	}
	executeTests(t, tests)
}

func TestBEQ(t *testing.T) {
	var tests = []InstructionTest{
		{"TestBCC -1", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xF0, 0xFE)
			p.Reg.SetStatus(ZeroFlag, true)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint16(0xD000), p.Reg.PC, name)
			assert.Equal(t, uint64(2), p.cycles, name)
		}},
		{"TestBCC -1 Zero Flag false", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xF0, 0xFE)
			p.Reg.SetStatus(ZeroFlag, false)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint16(0xD002), p.Reg.PC, name)
			assert.Equal(t, uint64(2), p.cycles, name)
		}},

		{"TestBCC +5", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xF0, 0x05)
			p.Reg.SetStatus(ZeroFlag, true)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint16(0xD007), p.Reg.PC, name)
			assert.Equal(t, uint64(2), p.cycles, name)
		}},
		{"TestBCC -5", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xF0, 0xf9)
			p.Reg.SetStatus(ZeroFlag, true)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint16(0xcffb), p.Reg.PC, name)
			assert.Equal(t, uint64(3), p.cycles, name)
		}},
	}
	executeTests(t, tests)
}

func TestBMI(t *testing.T) {
	var tests = []InstructionTest{
		{"TestBMI -1", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x30, 0xFE)
			p.Reg.SetStatus(NegativeFlag, true)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint16(0xD000), p.Reg.PC, name)
			assert.Equal(t, uint64(2), p.cycles, name)
		}},
		{"TestBMI -1 Negative Flag false", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x30, 0xFE)
			p.Reg.SetStatus(NegativeFlag, false)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint16(0xD002), p.Reg.PC, name)
			assert.Equal(t, uint64(2), p.cycles, name)
		}},

		{"TestBMI +5", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x30, 0x05)
			p.Reg.SetStatus(NegativeFlag, true)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint16(0xD007), p.Reg.PC, name)
			assert.Equal(t, uint64(2), p.cycles, name)
		}},
		{"TestBMI -5", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x30, 0xf9)
			p.Reg.SetStatus(NegativeFlag, true)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint16(0xcffb), p.Reg.PC, name)
			assert.Equal(t, uint64(3), p.cycles, name)
		}},
	}
	executeTests(t, tests)
}

func TestBNE(t *testing.T) {
	var tests = []InstructionTest{
		{"TestBNE -1", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xD0, 0xFE)
			p.Reg.SetStatus(ZeroFlag, false)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint16(0xD000), p.Reg.PC, name)
			assert.Equal(t, uint64(2), p.cycles, name)
		}},
		{"TestBNE -1 Zero Flag false", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xD0, 0xFE)
			p.Reg.SetStatus(ZeroFlag, true)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint16(0xD002), p.Reg.PC, name)
			assert.Equal(t, uint64(2), p.cycles, name)
		}},

		{"TestBNE +5", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xD0, 0x05)
			p.Reg.SetStatus(ZeroFlag, false)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint16(0xD007), p.Reg.PC, name)
			assert.Equal(t, uint64(2), p.cycles, name)
		}},
		{"TestBNE -5", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xD0, 0xf9)
			p.Reg.SetStatus(ZeroFlag, false)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint16(0xcffb), p.Reg.PC, name)
			assert.Equal(t, uint64(3), p.cycles, name)
		}},
	}
	executeTests(t, tests)
}

func TestBPL(t *testing.T) {
	var tests = []InstructionTest{
		{"TestBPL -1", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x10, 0xFE)
			p.Reg.SetStatus(NegativeFlag, false)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint16(0xD000), p.Reg.PC, name)
			assert.Equal(t, uint64(2), p.cycles, name)
		}},
		{"TestBPL -1 Negative Flag true", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x10, 0xFE)
			p.Reg.SetStatus(NegativeFlag, true)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint16(0xD002), p.Reg.PC, name)
			assert.Equal(t, uint64(2), p.cycles, name)
		}},

		{"TestBPL +5", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x10, 0x05)
			p.Reg.SetStatus(NegativeFlag, false)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint16(0xD007), p.Reg.PC, name)
			assert.Equal(t, uint64(2), p.cycles, name)
		}},
		{"TestBPL -5", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x10, 0xf9)
			p.Reg.SetStatus(NegativeFlag, false)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint16(0xcffb), p.Reg.PC, name)
			assert.Equal(t, uint64(3), p.cycles, name)
		}},
	}
	executeTests(t, tests)
}

func TestBVC(t *testing.T) {
	var tests = []InstructionTest{
		{"TestBVC -1", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x50, 0xFE)
			p.Reg.SetStatus(OverflowFlag, false)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint16(0xD000), p.Reg.PC, name)
			assert.Equal(t, uint64(2), p.cycles, name)
		}},
		{"TestBVC -1 Overflow Flag true", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x50, 0xFE)
			p.Reg.SetStatus(OverflowFlag, true)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint16(0xD002), p.Reg.PC, name)
			assert.Equal(t, uint64(2), p.cycles, name)
		}},

		{"TestBVC +5", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x50, 0x05)
			p.Reg.SetStatus(OverflowFlag, false)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint16(0xD007), p.Reg.PC, name)
			assert.Equal(t, uint64(2), p.cycles, name)
		}},
		{"TestBVC -5", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x50, 0xf9)
			p.Reg.SetStatus(OverflowFlag, false)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint16(0xcffb), p.Reg.PC, name)
			assert.Equal(t, uint64(3), p.cycles, name)
		}},
	}
	executeTests(t, tests)
}

func TestBVS(t *testing.T) {
	var tests = []InstructionTest{
		{"TestBVS -1", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x70, 0xFE)
			p.Reg.SetStatus(OverflowFlag, false)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint16(0xD002), p.Reg.PC, name)
			assert.Equal(t, uint64(2), p.cycles, name)
		}},
		{"TestBVS -1 Overflow Flag true", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x70, 0xFE)
			p.Reg.SetStatus(OverflowFlag, true)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint16(0xD000), p.Reg.PC, name)
			assert.Equal(t, uint64(2), p.cycles, name)
		}},

		{"TestBVS +5", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x70, 0x05)
			p.Reg.SetStatus(OverflowFlag, true)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint16(0xD007), p.Reg.PC, name)
			assert.Equal(t, uint64(2), p.cycles, name)
		}},
		{"TestBVS -5", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x70, 0xf9)
			p.Reg.SetStatus(OverflowFlag, true)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint16(0xcffb), p.Reg.PC, name)
			assert.Equal(t, uint64(3), p.cycles, name)
		}},
	}
	executeTests(t, tests)
}

func TestBIT(t *testing.T) {
	var tests = []InstructionTest{
		{"TestBIT ZeroPage", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x24, 0xFF)
			p.mem.Write(0x00ff, 0xc0)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag))
			assert.Equal(t, true, p.Reg.IsSet(OverflowFlag))
			assert.Equal(t, uint64(3), p.cycles, name)
		}},
		{"TestBIT absolute", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x2C, 0x00FF)
			p.mem.Write(0x00ff, 0xc0)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag))
			assert.Equal(t, true, p.Reg.IsSet(OverflowFlag))
			assert.Equal(t, uint64(4), p.cycles, name)
		}},
	}
	executeTests(t, tests)
}

func TestBRK(t *testing.T) {
	var tests = []InstructionTest{
		{"TestBRK", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x00)
			p.mem.Write(irqVector, 0x12)
			p.mem.Write(irqVector+1, 0xF0)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint16(0xF012), p.Reg.PC)
			assert.Equal(t, true, p.Reg.IsSet(BreakFlag))
			assert.Equal(t, uint8(0xD0), p.mem.Read(stackAddress)) // Check values on stack
			assert.Equal(t, uint8(0x02), p.mem.Read(stackAddress-1))
			assert.Equal(t, uint8(BreakFlag), p.mem.Read(stackAddress-2)&BreakFlag)
			assert.Equal(t, true, p.Reg.IsSet(BreakFlag))
			assert.Equal(t, uint64(7), p.cycles, name)
		}},
	}
	executeTests(t, tests)

}

func TestCLC(t *testing.T) {
	var tests = []InstructionTest{
		{"TestCLC", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x18, 0x00)
			p.Reg.SetStatus(CarryFlag, true)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, false, p.Reg.IsSet(CarryFlag))
			assert.Equal(t, uint64(2), p.cycles)
		}},
	}
	executeTests(t, tests)
}

func TestCLD(t *testing.T) {
	var tests = []InstructionTest{
		{"TestCLD", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xD8, 0x00)
			p.Reg.SetStatus(DecimalFlag, true)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, false, p.Reg.IsSet(DecimalFlag))
			assert.Equal(t, uint64(2), p.cycles)
		}},
	}
	executeTests(t, tests)
}

func TestCLI(t *testing.T) {
	var tests = []InstructionTest{
		{"TestCLI", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x58, 0x00)
			p.Reg.SetStatus(InterruptDisableFlag, true)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, false, p.Reg.IsSet(InterruptDisableFlag))
			assert.Equal(t, uint64(2), p.cycles)
		}},
	}
	executeTests(t, tests)
}

func TestCLV(t *testing.T) {
	var tests = []InstructionTest{
		{"TestCLV", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xB8, 0x00)
			p.Reg.SetStatus(OverflowFlag, true)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, false, p.Reg.IsSet(OverflowFlag))
			assert.Equal(t, uint64(2), p.cycles)
		}},
	}
	executeTests(t, tests)
}

/*
   Compare Result	N	Z	C
A, X, or Y < Memory	*	0	0
A, X, or Y = Memory	0	1	1
A, X, or Y > Memory	*	0	1

* The N flag will be bit 7 of A, X, or Y - Memory

N = ((A - Memory) & 0x80) == 0x80
Z = (A - Memory) = 0x00
C = A >= Memory
*/

func TestCMP(t *testing.T) {
	var tests = []InstructionTest{
		{"TestCMPImmediate A > Immediate", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xC9, 0x01)
			p.Reg.A = 0xFF
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
		}},
		{"TestCMPImmediateNegativeFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xC9, 0x81)
			p.Reg.A = 0x80
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(2), p.cycles)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(CarryFlag), name)

		}},
		{"TestCMPImmediateZeroFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xC9, 0xFE)
			p.Reg.A = 0xFE
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, false, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, true, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(CarryFlag), name)
		}},
		{"TestCMPZeropage", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xC5, 0x80)
			p.mem.Write(0x0080, 01)
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(3), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)

		}},
		{"TestCMPZeropageX", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xD5, 0x80)
			p.mem.Write(0x0081, 01)
			p.Reg.X = 0x01
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(4), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
		}},
		{"TestCMPAbsolute", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xCD, 0x80, 0x00)
			p.mem.Write(0x0080, 01)
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(4), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
		}},
		{"TestCMPAbsoluteX", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xDD, 0x80, 0x00)
			p.mem.Write(0x0081, 01)
			p.Reg.X = 0x01
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(4), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
		}},
		{"TestCMPAbsoluteXPageOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xDD, 0x01, 0x00)
			p.mem.Write(0x0100, 01)
			p.Reg.X = 0xFF
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(5), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
		}},
		{"TestCMPAbsoluteY", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xD9, 0x80, 0x00)
			p.mem.Write(0x0081, 01)
			p.Reg.Y = 0x01
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(4), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)

		}},
		{"TestCMPAbsoluteYPageOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xD9, 0x01, 0x00)
			p.mem.Write(0x0100, 01)
			p.Reg.A = 0xff
			p.Reg.Y = 0xFF
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(5), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
		}},
		{"TestCMPIndexedIndirect", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xC1, 0x05)
			p.mem.Write(0x000A, 0x10, 0x50)
			p.mem.Write(0x5010, 0x01)
			p.Reg.X = 0x05
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(6), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
		}},
		{"TestCMPIndexedIndirectPageOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xC1, 0x01)
			p.mem.Write(0x0000, 0x10, 0x50)
			p.mem.Write(0x5010, 0x01)
			p.Reg.A = 0xff
			p.Reg.X = 0xFF
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(6), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
		}},
		{"TestCMPIndirectIndexed", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xD1, 0x05)
			p.mem.Write(0x0005, 0x10, 0x50)
			p.mem.Write(0x5015, 0x01)
			p.Reg.A = 0xff
			p.Reg.Y = 0x05
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(5), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
		}},
		{"TestCMPIndirectIndexedPageOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xD1, 0x05)
			p.mem.Write(0x0005, 0x01, 0x50)
			p.mem.Write(0x5100, 0x01)
			p.Reg.A = 0xff
			p.Reg.Y = 0xFF
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(6), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
		}},
	}
	executeTests(t, tests)
}

func TestCPX(t *testing.T) {
	var tests = []InstructionTest{
		{"TestCPXImmediate A > Immediate", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xE0, 0x01)
			p.Reg.X = 0xFF
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
		}},
		{"TestCPXImmediateNegativeFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xE0, 0x81)
			p.Reg.X = 0x80
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(2), p.cycles)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(CarryFlag), name)

		}},
		{"TestCPXImmediateZeroFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xE0, 0xFE)
			p.Reg.X = 0xFE
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, false, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, true, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(CarryFlag), name)
		}},
		{"TestCPXZeropage", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xE4, 0x80)
			p.mem.Write(0x0080, 01)
			p.Reg.X = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(3), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
		}},
		{"TestCPXAbsolute", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xEC, 0x80, 0x00)
			p.mem.Write(0x0080, 01)
			p.Reg.X = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(4), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
		}},
	}
	executeTests(t, tests)
}

func TestCPY(t *testing.T) {
	var tests = []InstructionTest{
		{"TestCPYImmediate A > Immediate", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xC0, 0x01)
			p.Reg.Y = 0xFF
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
		}},
		{"TestCPYImmediateNegativeFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xC0, 0x81)
			p.Reg.Y = 0x80
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(2), p.cycles)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(CarryFlag), name)

		}},
		{"TestCPYImmediateZeroFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xC0, 0xFE)
			p.Reg.Y = 0xFE
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, false, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, true, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(CarryFlag), name)
		}},
		{"TestCPYZeropage", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xC4, 0x80)
			p.mem.Write(0x0080, 01)
			p.Reg.Y = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(3), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
		}},
		{"TestCPYAbsolute", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xCC, 0x80, 0x00)
			p.mem.Write(0x0080, 01)
			p.Reg.Y = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(4), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
		}},
	}
	executeTests(t, tests)
}

func TestDEC(t *testing.T) {
	var tests = []InstructionTest{
		{"TestDecZeropage", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xC6, 0xFF)
			p.mem.Write(0x00FF, 1)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(5), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, uint8(0x00), p.mem.Read(0x00FF))
		}},
		{"TestDecZeropageX", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xD6, 0xFF)
			p.mem.Write(0x0000, 1)
			p.Reg.X = 0x01
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(6), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, uint8(0x00), p.mem.Read(0x0100))
		}},
		{"TestDecAbsolute", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xCE, 0xFF, 0xFF)
			p.mem.Write(0xFFFF, 0xFF)
			p.Reg.X = 0x01
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(6), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, uint8(0xFE), p.mem.Read(0xFFFF))
		}},
		{"TestDecAbsoluteX", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xDE, 0xFF, 0x10)
			p.mem.Write(0x1100, 0xFF)
			p.Reg.X = 0x01
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(7), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, uint8(0xFE), p.mem.Read(0x1100))
		}},
	}
	executeTests(t, tests)
}

func TestRegDecrement(t *testing.T) {
	var tests = []InstructionTest{
		{"TestDex", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xCA)
			p.Reg.X = 3
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, uint8(2), p.Reg.X, name)
		}},
		{"TestDexOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xCA)
			p.Reg.X = 0x00
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, uint8(0xFF), p.Reg.X, name)
		}},
		{"TestDey", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x88)
			p.Reg.Y = 3
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, uint8(2), p.Reg.Y, name)
		}},
		{"TestDeyOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x88)
			p.Reg.Y = 0x00
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, uint8(0xFF), p.Reg.Y, name)
		}},
	}
	executeTests(t, tests)
}

func TestEOR(t *testing.T) {
	var tests = []InstructionTest{
		{"TestEORImmediate", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x49, 0x01)
			p.Reg.A = 0xFF
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0xfe), p.Reg.A, name)
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(OverflowFlag), name)
		}},
		{"TestEORImmediateNegativeFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x49, 0x80)
			p.Reg.A = 0x0
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x80), p.Reg.A)
			assert.Equal(t, uint64(2), p.cycles)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag))
		}},
		{"TestEORImmediateZeroFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x49, 0xFF)
			p.Reg.A = 0xFF
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x00), p.Reg.A, name)
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(ZeroFlag), name)
		}},
		{"TestEORZeropage", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x45, 0x80)
			p.mem.Write(0x0080, 01)
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0xFE), p.Reg.A, name)
			assert.Equal(t, uint64(3), p.cycles, name)
		}},
		{"TestEORZeropageX", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x55, 0x80)
			p.mem.Write(0x0081, 01)
			p.Reg.X = 0x01
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0xfe), p.Reg.A, name)
			assert.Equal(t, uint64(4), p.cycles, name)
		}},
		{"TestEORAbsolute", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x4D, 0x80, 0x00)
			p.mem.Write(0x0080, 01)
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0xfe), p.Reg.A, name)
			assert.Equal(t, uint64(4), p.cycles, name)
		}},
		{"TestEORAbsoluteX", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x5D, 0x80, 0x00)
			p.mem.Write(0x0081, 01)
			p.Reg.X = 0x01
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0xfe), p.Reg.A, name)
			assert.Equal(t, uint64(4), p.cycles, name)
		}},
		{"TestEORAbsoluteXPageOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x5D, 0x01, 0x00)
			p.mem.Write(0x0100, 01)
			p.Reg.X = 0xFF
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0xfe), p.Reg.A, name)
			assert.Equal(t, uint64(5), p.cycles, name)
		}},
		{"TestEORAbsoluteY", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x59, 0x80, 0x00)
			p.mem.Write(0x0081, 01)
			p.Reg.Y = 0x01
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0xfe), p.Reg.A, name)
			assert.Equal(t, uint64(4), p.cycles, name)

		}},
		{"TestEORAbsoluteYPageOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x59, 0x01, 0x00)
			p.mem.Write(0x0100, 01)
			p.Reg.A = 0xff
			p.Reg.Y = 0xFF
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0xfe), p.Reg.A, name)
			assert.Equal(t, uint64(5), p.cycles, name)
		}},
		{"TestEORIndexedIndirect", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x41, 0x05)
			p.mem.Write(0x000A, 0x10, 0x50)
			p.mem.Write(0x5010, 0xfe)
			p.Reg.X = 0x05
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(6), p.cycles, name)
		}},
		{"TestEORIndexedIndirectPageOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x41, 0x01)
			p.mem.Write(0x0000, 0x10, 0x50)
			p.mem.Write(0x5010, 0x01)
			p.Reg.A = 0xff
			p.Reg.X = 0xFF
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0xfe), p.Reg.A, name)
			assert.Equal(t, uint64(6), p.cycles, name)
		}},
		{"TestEORIndirectIndexed", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x51, 0x05)
			p.mem.Write(0x0005, 0x10, 0x50)
			p.mem.Write(0x5015, 0x01)
			p.Reg.A = 0xff
			p.Reg.Y = 0x05
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0xfe), p.Reg.A, name)
			assert.Equal(t, uint64(5), p.cycles, name)
		}},
		{"TestEORIndirectIndexedPageOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x51, 0x05)
			p.mem.Write(0x0005, 0x01, 0x50)
			p.mem.Write(0x5100, 0x01)
			p.Reg.A = 0xff
			p.Reg.Y = 0xFF
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0xfe), p.Reg.A, name)
			assert.Equal(t, uint64(6), p.cycles, name)
		}},
	}
	executeTests(t, tests)
}

func TestINC(t *testing.T) {
	var tests = []InstructionTest{
		{"TestIncZeropage", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xE6, 0xFF)
			p.mem.Write(0x00FF, 0xFF)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(5), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, uint8(0x00), p.mem.Read(0x00FF))
		}},
		{"TestIncZeropageX", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xf6, 0xFF)
			p.mem.Write(0x0000, 0xFF)
			p.Reg.X = 0x01
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(6), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, uint8(0x00), p.mem.Read(0x0100))
		}},
		{"TestIncAbsolute", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xEE, 0xFF, 0xFF)
			p.mem.Write(0xFFFF, 0xFD)
			p.Reg.X = 0x01
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(6), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, uint8(0xFE), p.mem.Read(0xFFFF))
		}},
		{"TestIncAbsoluteX", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xFE, 0xFF, 0x10)
			p.mem.Write(0x1100, 0xFD)
			p.Reg.X = 0x01
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(7), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, uint8(0xFE), p.mem.Read(0x1100))
		}},
	}
	executeTests(t, tests)
}

func TestRegIncrement(t *testing.T) {
	var tests = []InstructionTest{
		{"TestInx", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xE8)
			p.Reg.X = 3
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, uint8(4), p.Reg.X, name)
		}},
		{"TestInxOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xE8)
			p.Reg.X = 0xFF
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, uint8(0x00), p.Reg.X, name)
		}},
		{"TestIny", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xC8)
			p.Reg.Y = 3
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, uint8(4), p.Reg.Y, name)
		}},
		{"TestInyOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xC8)
			p.Reg.Y = 0xFF
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, uint8(0x00), p.Reg.Y, name)
		}},
	}
	executeTests(t, tests)
}
func TestJump(t *testing.T) {
	var tests = []InstructionTest{
		{"TestJMP", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x4C, 0x10, 0x50)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(3), p.cycles, name)
			assert.Equal(t, uint16(0x5010), p.Reg.PC)
		}},
		{"TestJMP", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x6C, 0x00, 0x51)
			p.mem.Write(0x5100, 0x010, 0x50)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(5), p.cycles, name)
			assert.Equal(t, uint16(0x5010), p.Reg.PC)
		}},
		{"TestJSR", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x20, 0x00, 0x51)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(6), p.cycles, name)
			assert.Equal(t, uint16(0x5100), p.Reg.PC)
			assert.Equal(t, byte(0xD0), p.mem.Read(stackAddress))
			assert.Equal(t, byte(0x03), p.mem.Read(stackAddress-1))
		}},
	}
	executeTests(t, tests)
}

func TestLDA(t *testing.T) {
	var tests = []InstructionTest{
		{"TestLDAImmediate", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xA9, 0x01)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, false, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
		}},
		{"TestLDAImmediateNegativeFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xA9, 0x80)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x80), p.Reg.A)
			assert.Equal(t, uint64(2), p.cycles)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag))
		}},
		{"TestLDAImmediateZeroFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xA9, 0x00)
			p.Reg.A = 0xFF
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x00), p.Reg.A, name)
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(ZeroFlag), name)
		}},
		{"TestLDAZeropage", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xA5, 0x80)
			p.mem.Write(0x0080, 01)
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(3), p.cycles, name)
		}},
		{"TestLDAZeropageX", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xB5, 0x80)
			p.mem.Write(0x0081, 01)
			p.Reg.X = 0x01
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(4), p.cycles, name)
		}},
		{"TestLDAAbsolute", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xAD, 0x80, 0x00)
			p.mem.Write(0x0080, 01)
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(4), p.cycles, name)
		}},
		{"TestLDAAbsoluteX", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xBD, 0x80, 0x00)
			p.mem.Write(0x0081, 01)
			p.Reg.X = 0x01
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(4), p.cycles, name)
		}},
		{"TestLDAAbsoluteXPageOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xBD, 0x01, 0x00)
			p.mem.Write(0x0100, 01)
			p.Reg.X = 0xFF
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(5), p.cycles, name)
		}},
		{"TestLDAAbsoluteY", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xB9, 0x80, 0x00)
			p.mem.Write(0x0081, 01)
			p.Reg.Y = 0x01
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(4), p.cycles, name)

		}},
		{"TestLDAAbsoluteYPageOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xB9, 0x01, 0x00)
			p.mem.Write(0x0100, 01)
			p.Reg.A = 0xff
			p.Reg.Y = 0xFF
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(5), p.cycles, name)
		}},
		{"TestLDAIndexedIndirect", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xA1, 0x05)
			p.mem.Write(0x000A, 0x10, 0x50)
			p.mem.Write(0x5010, 0x01)
			p.Reg.X = 0x05
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(6), p.cycles, name)
		}},
		{"TestLDAIndexedIndirectPageOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xA1, 0x01)
			p.mem.Write(0x0000, 0x10, 0x50)
			p.mem.Write(0x5010, 0x01)
			p.Reg.A = 0xff
			p.Reg.X = 0xFF
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(6), p.cycles, name)
		}},
		{"TestLDAIndirectIndexed", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xB1, 0x05)
			p.mem.Write(0x0005, 0x10, 0x50)
			p.mem.Write(0x5015, 0x01)
			p.Reg.A = 0xff
			p.Reg.Y = 0x05
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(5), p.cycles, name)
		}},
		{"TestLDAIndirectIndexedPageOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xB1, 0x05)
			p.mem.Write(0x0005, 0x01, 0x50)
			p.mem.Write(0x5100, 0x01)
			p.Reg.A = 0xff
			p.Reg.Y = 0xFF
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(6), p.cycles, name)
		}},
	}
	executeTests(t, tests)
}

func TestLDX(t *testing.T) {
	var tests = []InstructionTest{
		{"TestLDXImmediate", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xA2, 0x01)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.X, name)
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, false, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
		}},
		{"TestLDXImmediateNegativeFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xA2, 0x80)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x80), p.Reg.X)
			assert.Equal(t, uint64(2), p.cycles)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag))
		}},
		{"TestLDXImmediateZeroFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xA2, 0x00)
			p.Reg.X = 0xFF
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x00), p.Reg.X, name)
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(ZeroFlag), name)
		}},
		{"TestLDXZeropage", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xA6, 0x80)
			p.mem.Write(0x0080, 01)
			p.Reg.X = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.X, name)
			assert.Equal(t, uint64(3), p.cycles, name)
		}},
		{"TestLDXZeropageY", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xB6, 0x80)
			p.mem.Write(0x0081, 01)
			p.Reg.Y = 0x01
			p.Reg.X = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.X, name)
			assert.Equal(t, uint64(4), p.cycles, name)
		}},
		{"TestLDXAbsolute", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xAE, 0x80, 0x00)
			p.mem.Write(0x0080, 01)
			p.Reg.X = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.X, name)
			assert.Equal(t, uint64(4), p.cycles, name)
		}},
		{"TestLDXAbsoluteY", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xBE, 0x80, 0x00)
			p.mem.Write(0x0081, 01)
			p.Reg.X = 0x01
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.X, name)
			assert.Equal(t, uint64(4), p.cycles, name)
		}},
		{"TestLDXAbsoluteYPageOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xBE, 0x01, 0x00)
			p.mem.Write(0x0100, 01)
			p.Reg.X = 0xFF
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.X, name)
			assert.Equal(t, uint64(5), p.cycles, name)
		}},
	}
	executeTests(t, tests)
}

func TestLDY(t *testing.T) {
	var tests = []InstructionTest{
		{"TestLDYImmediate", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xA0, 0x01)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.Y, name)
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, false, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
		}},
		{"TestLDYImmediateNegativeFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xA0, 0x80)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x80), p.Reg.Y)
			assert.Equal(t, uint64(2), p.cycles)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag))
		}},
		{"TestLDYImmediateZeroFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xA0, 0x00)
			p.Reg.A = 0xFF
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x00), p.Reg.Y, name)
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(ZeroFlag), name)
		}},
		{"TestLDYZeropage", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xA4, 0x80)
			p.mem.Write(0x0080, 01)
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.Y, name)
			assert.Equal(t, uint64(3), p.cycles, name)
		}},
		{"TestLDYZeropageX", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xB4, 0x80)
			p.mem.Write(0x0081, 01)
			p.Reg.X = 0x01
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.Y, name)
			assert.Equal(t, uint64(4), p.cycles, name)
		}},
		{"TestLDYAbsolute", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xAC, 0x80, 0x00)
			p.mem.Write(0x0080, 01)
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.Y, name)
			assert.Equal(t, uint64(4), p.cycles, name)
		}},
		{"TestLDYAbsoluteX", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xBC, 0x80, 0x00)
			p.mem.Write(0x0081, 01)
			p.Reg.X = 0x01
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.Y, name)
			assert.Equal(t, uint64(4), p.cycles, name)
		}},
		{"TestLDYAbsoluteXPageOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xBC, 0x01, 0x00)
			p.mem.Write(0x0100, 01)
			p.Reg.X = 0xFF
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.Y, name)
			assert.Equal(t, uint64(5), p.cycles, name)
		}},
	}
	executeTests(t, tests)
}

func TestLSR(t *testing.T) {
	var tests = []InstructionTest{
		{"TestLSRImmediate", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x4A)
			p.Reg.A = 0x2
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, false, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
		}},
		{"TestLSRImmediateZeroAndCarryFlagsSet", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x4A)
			p.Reg.A = 0x1
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, byte(0x00), p.Reg.A, name)
			assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, true, p.Reg.IsSet(ZeroFlag), name)

		}},
		{"TestLSRZeropage", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x46, 0x80)
			p.mem.Write(0x0080, 02)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(5), p.cycles, name)
			assert.Equal(t, uint8(0x01), p.mem.Read(0x0080), name)
			assert.Equal(t, false, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
		}},
		{"TestLSRZeropageZeroAndCarryFlagsSet", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x46, 0x80)
			p.mem.Write(0x0080, 01)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(5), p.cycles, name)
			assert.Equal(t, uint8(0x00), p.mem.Read(0x0080), name)
			assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, true, p.Reg.IsSet(ZeroFlag), name)
		}},

		{"TestLSRZeropageX", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x56, 0x80)
			p.mem.Write(0x0084, 02)
			p.Reg.X = 0x04
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(6), p.cycles, name)
			assert.Equal(t, uint8(0x01), p.mem.Read(0x0084), name)
			assert.Equal(t, false, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
		}},
		{"TestLSRAbsolute", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x4E, 0x00, 0xC0)
			p.mem.Write(0xC000, 02)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(6), p.cycles, name)
			assert.Equal(t, uint8(0x01), p.mem.Read(0xC000), name)
			assert.Equal(t, false, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
		}},
		{"TestLSRAbsoluteX", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x5E, 0x00, 0xC0)
			p.mem.Write(0xC004, 02)
			p.Reg.X = 0x04
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(7), p.cycles, name)
			assert.Equal(t, uint8(0x01), p.mem.Read(0xC004), name)
			assert.Equal(t, false, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
		}},
	}
	executeTests(t, tests)
}

func TestNOP(t *testing.T) {
	var tests = []InstructionTest{
		{"TestLSRImmediate", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xEA)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(2), p.cycles, name)
		}},
	}
	executeTests(t, tests)

}

func TestORA(t *testing.T) {
	var tests = []InstructionTest{
		{"TestORAImmediate", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x09, 0x01)
			p.Reg.A = 0x00
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, false, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(OverflowFlag), name)
		}},
		{"TestORAImmediateNegativeFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x09, 0x80)
			p.Reg.A = 0x00
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x80), p.Reg.A)
			assert.Equal(t, uint64(2), p.cycles)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag))
		}},
		{"TestORAImmediateZeroFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x09, 0x00)
			p.Reg.A = 0x00
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x00), p.Reg.A, name)
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(ZeroFlag), name)
		}},
		{"TestORAZeropage", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x05, 0x80)
			p.mem.Write(0x0080, 01)
			p.Reg.A = 0x00
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(3), p.cycles, name)
		}},
		{"TestORAZeropageX", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x15, 0x80)
			p.mem.Write(0x0081, 01)
			p.Reg.X = 0x01
			p.Reg.A = 0x00
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(4), p.cycles, name)
		}},
		{"TestORAAbsolute", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x0d, 0x80, 0x00)
			p.mem.Write(0x0080, 01)
			p.Reg.A = 0x00
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(4), p.cycles, name)
		}},
		{"TestORAAbsoluteX", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x1d, 0x80, 0x00)
			p.mem.Write(0x0081, 01)
			p.Reg.X = 0x01
			p.Reg.A = 0x00
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(4), p.cycles, name)
		}},
		{"TestORAAbsoluteXPageOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x1d, 0x01, 0x00)
			p.mem.Write(0x0100, 01)
			p.Reg.X = 0xFF
			p.Reg.A = 0x00
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(5), p.cycles, name)
		}},
		{"TestORAAbsoluteY", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x19, 0x80, 0x00)
			p.mem.Write(0x0081, 01)
			p.Reg.Y = 0x01
			p.Reg.A = 0x00
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(4), p.cycles, name)

		}},
		{"TestORAAbsoluteYPageOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x19, 0x01, 0x00)
			p.mem.Write(0x0100, 01)
			p.Reg.A = 0x00
			p.Reg.Y = 0xFF
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(5), p.cycles, name)
		}},
		{"TestORAIndexedIndirect", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x01, 0x05)
			p.mem.Write(0x000A, 0x10, 0x50)
			p.mem.Write(0x5010, 0x01)
			p.Reg.X = 0x05
			p.Reg.A = 0x00
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(6), p.cycles, name)
		}},
		{"TestORAIndexedIndirectPageOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x01, 0x01)
			p.mem.Write(0x0000, 0x10, 0x50)
			p.mem.Write(0x5010, 0x01)
			p.Reg.A = 0x00
			p.Reg.X = 0xFF
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(6), p.cycles, name)
		}},
		{"TestORAIndirectIndexed", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x11, 0x05)
			p.mem.Write(0x0005, 0x10, 0x50)
			p.mem.Write(0x5015, 0x01)
			p.Reg.A = 0x00
			p.Reg.Y = 0x05
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(5), p.cycles, name)
		}},
		{"TestORAIndirectIndexedPageOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x11, 0x05)
			p.mem.Write(0x0005, 0x01, 0x50)
			p.mem.Write(0x5100, 0x01)
			p.Reg.A = 0x00
			p.Reg.Y = 0xFF
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, uint64(6), p.cycles, name)
		}},
	}
	executeTests(t, tests)
}

func TestPHA(t *testing.T) {
	var tests = []InstructionTest{
		{"TestPHA", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x48)
			p.Reg.A = 0xFF
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint8(0xFF), p.mem.Read(stackAddress)) // Check values on stack
			assert.Equal(t, uint64(3), p.cycles, name)
		}},
	}
	executeTests(t, tests)
}

func TestPHP(t *testing.T) {
	var tests = []InstructionTest{
		{"TestPHP", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x08)
			p.Reg.Status = 0xFF
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint8(0xFF), p.mem.Read(stackAddress)) // Check values on stack
			assert.Equal(t, uint64(3), p.cycles, name)
		}},
	}
	executeTests(t, tests)
}

func TestPLA(t *testing.T) {
	var tests = []InstructionTest{
		{"TestPLA", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x68)
			p.mem.Write(stackAddress, 0xFF)
			p.Reg.A = 0x00
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint8(0xFF), p.Reg.A) // Check values on stack
			assert.Equal(t, uint64(4), p.cycles, name)
		}},
	}
	executeTests(t, tests)
}

func TestPLP(t *testing.T) {
	var tests = []InstructionTest{
		{"TestPLP", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x28)
			p.mem.Write(stackAddress, 0xFF)
			p.Reg.Status = 0x00
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint8(0xFF), p.Reg.Status) // Check values on stack
			assert.Equal(t, uint64(4), p.cycles, name)
		}},
	}
	executeTests(t, tests)
}

func TestROL(t *testing.T) {
	var tests = []InstructionTest{
		{"TestROLAccumulator", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x2A)
			p.Reg.A = 0x01
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, byte(0x02), p.Reg.A, name)
			assert.Equal(t, false, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
		}},
		{"TestROLAccumulatorCarryFlags", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x2A)
			p.Reg.A = 0x80
			p.Reg.SetStatus(CarryFlag, true)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
		}},
		{"TestROLAccumulatorCarryFlagAndZeroFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x2A)
			p.Reg.A = 0x80
			p.Reg.SetStatus(CarryFlag, false)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, byte(0x00), p.Reg.A, name)
			assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, true, p.Reg.IsSet(ZeroFlag), name)
		}},
		{"TestROLZeropage", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x26, 0x80)
			p.mem.Write(0x0080, 0x01)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(5), p.cycles, name)
			assert.Equal(t, uint8(0x02), p.mem.Read(0x0080), name)
			assert.Equal(t, false, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
		}},
		{"TestROLZeropageX", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x36, 0x80)
			p.mem.Write(0x0084, 0x01)
			p.Reg.X = 0x04
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(6), p.cycles, name)
			assert.Equal(t, uint8(0x02), p.mem.Read(0x0084), name)
			assert.Equal(t, false, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
		}},
		{"TestROLAbsolute", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x2E, 0x00, 0xC0)
			p.mem.Write(0xC000, 0x01)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(6), p.cycles, name)
			assert.Equal(t, uint8(0x02), p.mem.Read(0xC000), name)
			assert.Equal(t, false, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
		}},
		{"TestROLAbsoluteX", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x3E, 0x00, 0xC0)
			p.mem.Write(0xC004, 0x01)
			p.Reg.X = 0x04
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(7), p.cycles, name)
			assert.Equal(t, uint8(0x02), p.mem.Read(0xC004), name)
			assert.Equal(t, false, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
		}},
	}
	executeTests(t, tests)
}

func TestROR(t *testing.T) {
	var tests = []InstructionTest{
		{"TestRORAccumulator", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x6A)
			p.Reg.A = 0x80
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, byte(0x40), p.Reg.A, name)
			assert.Equal(t, false, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
		}},
		{"TestRORAccumulatorCarryFlags", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x6A)
			p.Reg.A = 0x01
			p.Reg.SetStatus(CarryFlag, true)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, byte(0x80), p.Reg.A, name)
			assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
		}},
		{"TestRORAccumulatorCarryFlagAndZeroFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x6A)
			p.Reg.A = 0x01
			p.Reg.SetStatus(CarryFlag, false)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, byte(0x00), p.Reg.A, name)
			assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, true, p.Reg.IsSet(ZeroFlag), name)
		}},
		{"TestRORZeropage", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x66, 0x80)
			p.mem.Write(0x0080, 0x80)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(5), p.cycles, name)
			assert.Equal(t, uint8(0x40), p.mem.Read(0x0080), name)
			assert.Equal(t, false, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
		}},
		{"TestRORZeropageX", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x76, 0x80)
			p.mem.Write(0x0084, 0x80)
			p.Reg.X = 0x04
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(6), p.cycles, name)
			assert.Equal(t, uint8(0x40), p.mem.Read(0x0084), name)
			assert.Equal(t, false, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
		}},
		{"TestRORAbsolute", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x6E, 0x00, 0xC0)
			p.mem.Write(0xC000, 0x80)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(6), p.cycles, name)
			assert.Equal(t, uint8(0x40), p.mem.Read(0xC000), name)
			assert.Equal(t, false, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
		}},
		{"TestRORAbsoluteX", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x7E, 0x00, 0xC0)
			p.mem.Write(0xC004, 0x80)
			p.Reg.X = 0x04
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(7), p.cycles, name)
			assert.Equal(t, uint8(0x40), p.mem.Read(0xC004), name)
			assert.Equal(t, false, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
		}},
	}
	executeTests(t, tests)
}

func TestRTI(t *testing.T) {
	var tests = []InstructionTest{
		{"TestRORAccumulator", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x40)
			p.mem.Write(stackAddress, 0xC0)
			p.mem.Write(stackAddress-1, 0x00)
			p.mem.Write(stackAddress-2, 0xFF)
			p.Reg.S = uint8(0xFF - 2)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(6), p.cycles, name)
			assert.Equal(t, uint16(0xC000), p.Reg.PC, name)
			assert.Equal(t, false, p.Reg.IsSet(BreakFlag), name)
		}},
	}
	executeTests(t, tests)
}

func TestRTS(t *testing.T) {
	var tests = []InstructionTest{
		{"TestRORAccumulator", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x60)
			p.mem.Write(stackAddress, 0xC0)
			p.mem.Write(stackAddress-1, 0x00)
			p.Reg.S = uint8(0xFF - 1)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(6), p.cycles, name)
			assert.Equal(t, uint16(0xC000), p.Reg.PC, name)
		}},
	}
	executeTests(t, tests)
}

// See https://www.righto.com/2012/12/the-6502-overflow-flag-explained.html#:~:text=The%206502%20has%20a%20SBC,the%20carry%20flag%20is%20used.
func TestSBC(t *testing.T) {
	var tests = []InstructionTest{
		{"TestSBCImmediate", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xE9, 0x01)
			p.Reg.A = 0xFF
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0xFD), p.Reg.A, name) // Carry flag not set
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(OverflowFlag), name)
		}},
		{"TestSBCImmediateWithCarryFlagSet", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xE9, 0x01)
			p.Reg.SetStatus(CarryFlag, true)
			p.Reg.A = 0x00
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0xFF), p.Reg.A, name)
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag))
		}},
		{"TestSBCImmediateZeroFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xE9, 0x01)
			p.Reg.SetStatus(CarryFlag, true)
			p.Reg.A = 0x01
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x00), p.Reg.A, name)
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(ZeroFlag), name)
		}},
		{"TestSBCImmediateCarryFlagNotSet", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xE9, 176)
			p.Reg.SetStatus(CarryFlag, true)
			p.Reg.A = 80
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(160), p.Reg.A, name)
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, false, p.Reg.IsSet(CarryFlag), name)
		}},
		{"TestSBCImmediateCarryFlagSet", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xE9, 176)
			p.Reg.SetStatus(CarryFlag, true)
			p.Reg.A = 208
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(32), p.Reg.A, name)
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
		}},

		{"TestSBCImmediateOverFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xE9, 5)
			p.Reg.A = 130
			p.Reg.SetStatus(CarryFlag, true)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(125), p.Reg.A, name)
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(OverflowFlag), name)
			assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
		}},
		{"TestSBCZeropage", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xE5, 0x80)
			p.mem.Write(0x0080, 01)
			p.Reg.SetStatus(CarryFlag, true)
			p.Reg.A = 10
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x09), p.Reg.A, name)
			assert.Equal(t, uint64(3), p.cycles, name)
		}},
		{"TestSBCZeropageX", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xF5, 0x80)
			p.mem.Write(0x0081, 01)
			p.Reg.SetStatus(CarryFlag, true)
			p.Reg.X = 0x01
			p.Reg.A = 10
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x09), p.Reg.A, name)
			assert.Equal(t, uint64(4), p.cycles, name)
		}},
		{"TestSBCAbsolute", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xED, 0x80, 0x00)
			p.mem.Write(0x0080, 01)
			p.Reg.SetStatus(CarryFlag, true)
			p.Reg.A = 10
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(9), p.Reg.A, name)
			assert.Equal(t, uint64(4), p.cycles, name)
		}},
		{"TestSBCAbsoluteX", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xFD, 0x80, 0x00)
			p.mem.Write(0x0081, 01)
			p.Reg.SetStatus(CarryFlag, true)
			p.Reg.X = 0x01
			p.Reg.A = 10
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x09), p.Reg.A, name)
			assert.Equal(t, uint64(4), p.cycles, name)
		}},
		{"TestSBCAbsoluteXPageOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xFD, 0x01, 0x00)
			p.mem.Write(0x0100, 01)
			p.Reg.SetStatus(CarryFlag, true)
			p.Reg.X = 0xFF
			p.Reg.A = 10
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x09), p.Reg.A, name)
			assert.Equal(t, uint64(5), p.cycles, name)
		}},
		{"TestSBCAbsoluteY", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xF9, 0x80, 0x00)
			p.mem.Write(0x0081, 01)
			p.Reg.SetStatus(CarryFlag, true)
			p.Reg.Y = 0x01
			p.Reg.A = 10
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x09), p.Reg.A, name)
			assert.Equal(t, uint64(4), p.cycles, name)

		}},
		{"TestSBCAbsoluteYPageOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xF9, 0x01, 0x00)
			p.mem.Write(0x0100, 01)
			p.Reg.SetStatus(CarryFlag, true)
			p.Reg.Y = 0xFF
			p.Reg.A = 10
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x09), p.Reg.A)
			assert.Equal(t, uint64(5), p.cycles)
		}},
		{"TestSDCIndexedIndirect", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xE1, 0x05)
			p.mem.Write(0x000A, 0x10, 0x50)
			p.mem.Write(0x5010, 0x01)
			p.Reg.SetStatus(CarryFlag, true)
			p.Reg.X = 0x05
			p.Reg.A = 10
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x09), p.Reg.A, name)
			assert.Equal(t, uint64(6), p.cycles, name)
		}},
		{"TestSDCIndexedIndirectPageOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xE1, 0x01)
			p.mem.Write(0x0000, 0x10, 0x50)
			p.mem.Write(0x5010, 0x01)
			p.Reg.SetStatus(CarryFlag, true)
			p.Reg.X = 0xFF
			p.Reg.A = 10
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint8(9), p.Reg.A, name)
			assert.Equal(t, uint64(6), p.cycles, name)
		}},
		{"TestSBCIndirectIndexed", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xF1, 0x05) // ADC ($0x05),Y
			p.mem.Write(0x0005, 0x10, 0x50)
			p.mem.Write(0x5015, 0x01)
			p.Reg.SetStatus(CarryFlag, true)
			p.Reg.Y = 0x05
			p.Reg.A = 10
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x09), p.Reg.A, name)
			assert.Equal(t, uint64(5), p.cycles, name)
		}},
		{"TestSBCIndirectIndexedPageOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xF1, 0x05)
			p.mem.Write(0x0005, 0x01, 0x50)
			p.mem.Write(0x5100, 0x01)
			p.Reg.SetStatus(CarryFlag, true)
			p.Reg.Y = 0xFF
			p.Reg.A = 10
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x09), p.Reg.A, name)
			assert.Equal(t, uint64(6), p.cycles, name)
		}},
		{"TestSBCDecimalMode", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xE9, 5)
			p.Reg.SetStatus(DecimalFlag, true)
			p.Reg.A = 9
			p.Reg.SetStatus(CarryFlag, true)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint8(4), p.Reg.A, name)
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, false, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(OverflowFlag), name)
		}},
		{"TestSBCDecimalMode", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xE9, 0x21)
			p.Reg.SetStatus(DecimalFlag, true)
			p.Reg.SetStatus(CarryFlag, true)
			p.Reg.A = 0x19
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x98), p.Reg.A, name)
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(OverflowFlag), name)
		}},
		{"TestSBCDecimalMode", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xE9, 0x21)
			p.Reg.SetStatus(DecimalFlag, true)
			p.Reg.A = 0x19
			p.Reg.SetStatus(CarryFlag, false)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x97), p.Reg.A, name)
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(OverflowFlag), name)
		}},
		{"TestADCDecimalMode", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xE9, 0x08)
			p.Reg.SetStatus(DecimalFlag, true)
			p.Reg.SetStatus(CarryFlag, true)
			p.Reg.A = 0x06
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x98), p.Reg.A, name)
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(OverflowFlag), name)
		}},
		{"TestADCDecimalMode", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xE9, 0x11)
			p.Reg.SetStatus(DecimalFlag, true)
			p.Reg.SetStatus(CarryFlag, true)
			p.Reg.A = 0x90
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x79), p.Reg.A, name)
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, false, p.Reg.IsSet(NegativeFlag), name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
			assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
			assert.Equal(t, true, p.Reg.IsSet(OverflowFlag), name)
		}},
	}
	executeTests(t, tests)
}

func TestSEC(t *testing.T) {
	var tests = []InstructionTest{
		{"TestSEC", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x38)
			p.Reg.SetStatus(CarryFlag, false)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
		}},
	}
	executeTests(t, tests)
}

func TestSED(t *testing.T) {
	var tests = []InstructionTest{
		{"TestSED", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xF8)
			p.Reg.SetStatus(DecimalFlag, false)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(DecimalFlag), name)
		}},
	}
	executeTests(t, tests)
}

func TestSEI(t *testing.T) {
	var tests = []InstructionTest{
		{"TestSEI", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x78)
			p.Reg.SetStatus(InterruptDisableFlag, false)
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, uint64(2), p.cycles, name)
			assert.Equal(t, true, p.Reg.IsSet(InterruptDisableFlag), name)
		}},
	}
	executeTests(t, tests)
}

func TestSTA(t *testing.T) {
	var tests = []InstructionTest{
		{"TestSTAZeropage", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x85, 0x80)
			p.mem.Write(0x0080, 00)
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0xff), p.mem.Read(0x0080), name)
			assert.Equal(t, uint64(3), p.cycles, name)
		}},
		{"TestSTAZeropageX", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x95, 0x80)
			p.Reg.X = 0x01
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0xff), p.mem.Read(0x0081), name)
			assert.Equal(t, uint64(4), p.cycles, name)
		}},
		{"TestSTAAbsolute", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x8D, 0x80, 0x00)
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0xff), p.mem.Read(0x0080), name)
			assert.Equal(t, uint64(4), p.cycles, name)
		}},
		{"TestSTAAbsoluteX", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x9D, 0x80, 0x00)
			p.Reg.X = 0x01
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0xff), p.mem.Read(0x0081), name)
			assert.Equal(t, uint64(5), p.cycles, name)
		}},
		{"TestSTAAbsoluteXPageOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x9D, 0x01, 0x00)
			p.Reg.X = 0xFF
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0xff), p.mem.Read(0x0100), name)
			assert.Equal(t, uint64(5), p.cycles, name)
		}},
		{"TestSTAAbsoluteY", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x99, 0x80, 0x00)
			p.Reg.Y = 0x01
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0xff), p.mem.Read(0x0081), name)
			assert.Equal(t, uint64(5), p.cycles, name)

		}},
		{"TestSTAAbsoluteYPageOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x99, 0x01, 0x00)
			p.Reg.A = 0xff
			p.Reg.Y = 0xFF
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0xff), p.mem.Read(0x0100), name)
			assert.Equal(t, uint64(5), p.cycles, name)
		}},
		{"TestSTAIndexedIndirect", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x81, 0x05)
			p.mem.Write(0x000A, 0x10, 0x50)
			p.Reg.X = 0x05
			p.Reg.A = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0xff), p.mem.Read(0x5010), name)
			assert.Equal(t, uint64(6), p.cycles, name)
		}},
		{"TestSTAIndexedIndirectPageOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x81, 0x01)
			p.mem.Write(0x0000, 0x10, 0x50)
			p.Reg.A = 0xff
			p.Reg.X = 0xFF
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0xff), p.mem.Read(0x5010), name)
			assert.Equal(t, uint64(6), p.cycles, name)
		}},
		{"TestSTAIndirectIndexed", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x91, 0x05)
			p.mem.Write(0x0005, 0x10, 0x50)
			p.Reg.A = 0xff
			p.Reg.Y = 0x05
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0xff), p.mem.Read(0x5015), name)
			assert.Equal(t, uint64(6), p.cycles, name)
		}},
		{"TestSTAIndirectIndexedPageOverflow", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x91, 0x05)
			p.mem.Write(0x0005, 0x01, 0x50)
			p.Reg.A = 0xff
			p.Reg.Y = 0x05
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0xff), p.mem.Read(0x5006), name)
			assert.Equal(t, uint64(6), p.cycles, name)
		}},
	}
	executeTests(t, tests)
}

func TestSTX(t *testing.T) {
	var tests = []InstructionTest{
		{"TestSTXZeropage", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x86, 0x80)
			p.mem.Write(0x0080, 00)
			p.Reg.X = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0xff), p.mem.Read(0x0080), name)
			assert.Equal(t, uint64(3), p.cycles, name)
		}},
		{"TestSTXZeropageY", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x96, 0x80)
			p.Reg.Y = 0x01
			p.Reg.X = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0xff), p.mem.Read(0x0081), name)
			assert.Equal(t, uint64(4), p.cycles, name)
		}},
		{"TestSTXAbsolute", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x8E, 0x80, 0x00)
			p.Reg.X = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0xff), p.mem.Read(0x0080), name)
			assert.Equal(t, uint64(4), p.cycles, name)
		}},
	}
	executeTests(t, tests)
}

func TestSTY(t *testing.T) {
	var tests = []InstructionTest{
		{"TestSTYZeropage", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x84, 0x80)
			p.mem.Write(0x0080, 00)
			p.Reg.Y = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0xff), p.mem.Read(0x0080), name)
			assert.Equal(t, uint64(3), p.cycles, name)
		}},
		{"TestSTYZeropageX", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x94, 0x80)
			p.Reg.X = 0x01
			p.Reg.Y = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0xff), p.mem.Read(0x0081), name)
			assert.Equal(t, uint64(4), p.cycles, name)
		}},
		{"TestSTYAbsolute", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x8C, 0x80, 0x00)
			p.Reg.Y = 0xff
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0xff), p.mem.Read(0x0080), name)
			assert.Equal(t, uint64(4), p.cycles, name)
		}},
	}
	executeTests(t, tests)
}

func TestTAX(t *testing.T) {
	var tests = []InstructionTest{
		{"TestTAX", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xAA)
			p.Reg.A = 0x01
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.X, name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag))
			assert.Equal(t, false, p.Reg.IsSet(NegativeFlag))
			assert.Equal(t, uint64(2), p.cycles, name)
		}},
		{"TestTAXNegativeFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xAA)
			p.Reg.A = 0x80
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x80), p.Reg.X, name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag))
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag))
			assert.Equal(t, uint64(2), p.cycles, name)
		}},
		{"TestTAXZeroFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xAA)
			p.Reg.A = 0x00
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x00), p.Reg.X, name)
			assert.Equal(t, true, p.Reg.IsSet(ZeroFlag))
			assert.Equal(t, false, p.Reg.IsSet(NegativeFlag))
			assert.Equal(t, uint64(2), p.cycles, name)
		}},
	}
	executeTests(t, tests)
}

func TestTAY(t *testing.T) {
	var tests = []InstructionTest{
		{"TestTAY", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xA8)
			p.Reg.A = 0x01
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.Y, name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag))
			assert.Equal(t, false, p.Reg.IsSet(NegativeFlag))
			assert.Equal(t, uint64(2), p.cycles, name)
		}},
		{"TestTAYNegativeFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xA8)
			p.Reg.A = 0x80
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x80), p.Reg.Y, name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag))
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag))
			assert.Equal(t, uint64(2), p.cycles, name)
		}},
		{"TestTAYZeroFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xA8)
			p.Reg.A = 0x00
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x00), p.Reg.Y, name)
			assert.Equal(t, true, p.Reg.IsSet(ZeroFlag))
			assert.Equal(t, false, p.Reg.IsSet(NegativeFlag))
			assert.Equal(t, uint64(2), p.cycles, name)
		}},
	}
	executeTests(t, tests)
}

func TestTSX(t *testing.T) {
	var tests = []InstructionTest{
		{"TestTSX", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xBA)
			p.Reg.S = 0x01
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.X, name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag))
			assert.Equal(t, false, p.Reg.IsSet(NegativeFlag))
			assert.Equal(t, uint64(2), p.cycles, name)
		}},
		{"TestTSXNegativeFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xBA)
			p.Reg.S = 0x80
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x80), p.Reg.X, name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag))
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag))
			assert.Equal(t, uint64(2), p.cycles, name)
		}},
		{"TestTSXZeroFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0xBA)
			p.Reg.S = 0x00
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x00), p.Reg.X, name)
			assert.Equal(t, true, p.Reg.IsSet(ZeroFlag))
			assert.Equal(t, false, p.Reg.IsSet(NegativeFlag))
			assert.Equal(t, uint64(2), p.cycles, name)
		}},
	}
	executeTests(t, tests)
}

func TestTXA(t *testing.T) {
	var tests = []InstructionTest{
		{"TestTXA", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x8A)
			p.Reg.X = 0x01
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag))
			assert.Equal(t, false, p.Reg.IsSet(NegativeFlag))
			assert.Equal(t, uint64(2), p.cycles, name)
		}},
		{"TestTXANegativeFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x8A)
			p.Reg.X = 0x80
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x80), p.Reg.A, name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag))
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag))
			assert.Equal(t, uint64(2), p.cycles, name)
		}},
		{"TestTXAZeroFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x8A)
			p.Reg.X = 0x00
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x00), p.Reg.A, name)
			assert.Equal(t, true, p.Reg.IsSet(ZeroFlag))
			assert.Equal(t, false, p.Reg.IsSet(NegativeFlag))
			assert.Equal(t, uint64(2), p.cycles, name)
		}},
	}
	executeTests(t, tests)
}

func TestTXS(t *testing.T) {
	var tests = []InstructionTest{
		{"TestTXS", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x9A)
			p.Reg.X = 0x01
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.S, name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag))
			assert.Equal(t, false, p.Reg.IsSet(NegativeFlag))
			assert.Equal(t, uint64(2), p.cycles, name)
		}},
		{"TestTXSNegativeFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x9A)
			p.Reg.X = 0x80
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x80), p.Reg.S, name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag))
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag))
			assert.Equal(t, uint64(2), p.cycles, name)
		}},
		{"TestTXSZeroFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x9A)
			p.Reg.X = 0x00
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x00), p.Reg.S, name)
			assert.Equal(t, true, p.Reg.IsSet(ZeroFlag))
			assert.Equal(t, false, p.Reg.IsSet(NegativeFlag))
			assert.Equal(t, uint64(2), p.cycles, name)
		}},
	}
	executeTests(t, tests)
}

func TestTYA(t *testing.T) {
	var tests = []InstructionTest{
		{"TestTYA", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x98)
			p.Reg.Y = 0x01
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x01), p.Reg.A, name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag))
			assert.Equal(t, false, p.Reg.IsSet(NegativeFlag))
			assert.Equal(t, uint64(2), p.cycles, name)
		}},
		{"TestTYANegativeFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x98)
			p.Reg.Y = 0x80
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x80), p.Reg.A, name)
			assert.Equal(t, false, p.Reg.IsSet(ZeroFlag))
			assert.Equal(t, true, p.Reg.IsSet(NegativeFlag))
			assert.Equal(t, uint64(2), p.cycles, name)
		}},
		{"TestTYAZeroFlag", func(p *Cpu) int {
			p.mem.Write(startAddress, 0x98)
			p.Reg.Y = 0x00
			return 1
		}, func(t *testing.T, p *Cpu, name string) {
			assert.Equal(t, byte(0x00), p.Reg.A, name)
			assert.Equal(t, true, p.Reg.IsSet(ZeroFlag))
			assert.Equal(t, false, p.Reg.IsSet(NegativeFlag))
			assert.Equal(t, uint64(2), p.cycles, name)
		}},
	}
	executeTests(t, tests)
}

// p.opCodes[0x4C] = p.NewInstructionDefinition(InstructionStr(JMP, absoluteModeStr), 3)
// p.opCodes[0x6C] = p.NewInstructionDefinition(InstructionStr(JMP, absoluteIndirectModeStr), 5)
// p.opCodes[0x20] = p.NewInstructionDefinition(InstructionStr(JSR, absoluteModeStr), 6)
