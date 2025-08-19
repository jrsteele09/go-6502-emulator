package cpu

import (
    "testing"

    "github.com/stretchr/testify/assert"
)

// Tests for undocumented/illegal opcodes implemented in illegal_instructions.go
func TestIllegalOpcodes(t *testing.T) {
    var tests = []InstructionTest{
        {"LAX zeropage loads A and X, sets Z/N", func(p *CPU) int {
            // LAX $80 (0xA7)
            p.mem.Write(startAddress, 0xA7, 0x80)
            p.mem.Write(0x0080, 0x42)
            return 1
        }, func(t *testing.T, p *CPU, name string) {
            assert.Equal(t, byte(0x42), p.Reg.A, name)
            assert.Equal(t, byte(0x42), p.Reg.X, name)
            assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
            assert.Equal(t, false, p.Reg.IsSet(NegativeFlag), name)
            assert.Equal(t, uint64(3), p.cycles, name)
        }},

        {"SAX zeropage stores A&X without changing flags", func(p *CPU) int {
            // SAX $10 (0x87)
            p.mem.Write(startAddress, 0x87, 0x10)
            p.Reg.A = 0xF0
            p.Reg.X = 0x0F
            return 1
        }, func(t *testing.T, p *CPU, name string) {
            assert.Equal(t, byte(0x00), p.mem.Read(0x0010), name)
            assert.Equal(t, uint64(3), p.cycles, name)
        }},

        {"SLO absolute shifts mem left and ORA into A", func(p *CPU) int {
            // SLO $D020 (0x0F)
            p.mem.Write(startAddress, 0x0F, 0x20, 0xD0)
            p.mem.Write(0xD020, 0x40)
            p.Reg.A = 0x01
            return 1
        }, func(t *testing.T, p *CPU, name string) {
            assert.Equal(t, byte(0x80), p.mem.Read(0xD020), name+" mem")
            assert.Equal(t, byte(0x81), p.Reg.A, name+" A")
            assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name+" C")
            assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name+" N")
            assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name+" Z")
            assert.Equal(t, uint64(6), p.cycles, name)
        }},

        {"RLA zeropage rotates left and ANDs into A", func(p *CPU) int {
            // RLA $80 (0x27)
            p.mem.Write(startAddress, 0x27, 0x80)
            p.mem.Write(0x0080, 0x80)
            p.Reg.A = 0xFF
            p.Reg.SetStatus(CarryFlag, false)
            return 1
        }, func(t *testing.T, p *CPU, name string) {
            assert.Equal(t, byte(0x00), p.mem.Read(0x0080), name+" mem")
            assert.Equal(t, byte(0x00), p.Reg.A, name+" A")
            assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name+" C from old bit7")
            assert.Equal(t, true, p.Reg.IsSet(ZeroFlag), name+" Z")
            assert.Equal(t, false, p.Reg.IsSet(NegativeFlag), name+" N")
            assert.Equal(t, uint64(5), p.cycles, name)
        }},

        {"SRE zeropage shifts right and EORs into A", func(p *CPU) int {
            // SRE $40 (0x47)
            p.mem.Write(startAddress, 0x47, 0x40)
            p.mem.Write(0x0040, 0x03)
            p.Reg.A = 0x01
            return 1
        }, func(t *testing.T, p *CPU, name string) {
            assert.Equal(t, byte(0x01), p.mem.Read(0x0040), name+" mem")
            assert.Equal(t, byte(0x00), p.Reg.A, name+" A")
            assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name+" C from bit0")
            assert.Equal(t, true, p.Reg.IsSet(ZeroFlag), name+" Z")
            assert.Equal(t, uint64(5), p.cycles, name)
        }},

    {"RRA zeropage rotates right and ADCs into A (carry becomes ADC carry-in)", func(p *CPU) int {
            // RRA $20 (0x67)
            p.mem.Write(startAddress, 0x67, 0x20)
            p.mem.Write(0x0020, 0x01)
            p.Reg.A = 0x00
            p.Reg.SetStatus(CarryFlag, false) // ensures ROR produces res=0x00, carryOut=1 which becomes ADC carry-in
            return 1
        }, func(t *testing.T, p *CPU, name string) {
            assert.Equal(t, byte(0x00), p.mem.Read(0x0020), name+" mem after ROR")
            assert.Equal(t, byte(0x01), p.Reg.A, name+" A after ADC")
            // After ADC, carry reflects the addition result, not ROR's carry-out
            assert.Equal(t, false, p.Reg.IsSet(CarryFlag), name+" C after ADC")
            assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name+" Z")
            assert.Equal(t, false, p.Reg.IsSet(NegativeFlag), name+" N")
            assert.Equal(t, uint64(5), p.cycles, name)
        }},

        {"DCP zeropage decrements and compares with A", func(p *CPU) int {
            // DCP $11 (0xC7)
            p.mem.Write(startAddress, 0xC7, 0x11)
            p.mem.Write(0x0011, 0x00)
            p.Reg.A = 0x80
            return 1
        }, func(t *testing.T, p *CPU, name string) {
            assert.Equal(t, byte(0xFF), p.mem.Read(0x0011), name+" mem")
            assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
            assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
            assert.Equal(t, false, p.Reg.IsSet(CarryFlag), name)
            assert.Equal(t, uint64(5), p.cycles, name)
        }},

        {"ISC zeropage increments and SBCs from A", func(p *CPU) int {
            // ISC $12 (0xE7)
            p.mem.Write(startAddress, 0xE7, 0x12)
            p.mem.Write(0x0012, 0xFF) // becomes 0x00
            p.Reg.A = 0x02
            p.Reg.SetStatus(CarryFlag, true)
            return 1
        }, func(t *testing.T, p *CPU, name string) {
            assert.Equal(t, byte(0x00), p.mem.Read(0x0012), name+" mem")
            assert.Equal(t, byte(0x02), p.Reg.A, name+" A")
            assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
            assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
            assert.Equal(t, false, p.Reg.IsSet(NegativeFlag), name)
            assert.Equal(t, uint64(5), p.cycles, name)
        }},

        {"ANC immediate ANDs and moves bit7 to C", func(p *CPU) int {
            // ANC #$FF (0x0B)
            p.mem.Write(startAddress, 0x0B, 0xFF)
            p.Reg.A = 0x81
            return 1
        }, func(t *testing.T, p *CPU, name string) {
            assert.Equal(t, byte(0x81), p.Reg.A, name)
            assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
            assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
            assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
            assert.Equal(t, uint64(2), p.cycles, name)
        }},

        {"ALR immediate AND then LSR A", func(p *CPU) int {
            // ALR #$FF (0x4B)
            p.mem.Write(startAddress, 0x4B, 0xFF)
            p.Reg.A = 0x03
            return 1
        }, func(t *testing.T, p *CPU, name string) {
            assert.Equal(t, byte(0x01), p.Reg.A, name)
            assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
            assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
            assert.Equal(t, false, p.Reg.IsSet(NegativeFlag), name)
            assert.Equal(t, uint64(2), p.cycles, name)
        }},

        {"ARR immediate AND then ROR A (approx flags)", func(p *CPU) int {
            // ARR #$FF (0x6B)
            p.mem.Write(startAddress, 0x6B, 0xFF)
            p.Reg.A = 0xFF
            p.Reg.SetStatus(CarryFlag, true) // carry-in to ROR A
            return 1
        }, func(t *testing.T, p *CPU, name string) {
            assert.Equal(t, byte(0xFF), p.Reg.A, name)
            assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
            assert.Equal(t, false, p.Reg.IsSet(OverflowFlag), name)
            assert.Equal(t, true, p.Reg.IsSet(NegativeFlag), name)
            assert.Equal(t, false, p.Reg.IsSet(ZeroFlag), name)
            assert.Equal(t, uint64(2), p.cycles, name)
        }},

        {"XAA immediate sets A = X & imm", func(p *CPU) int {
            // XAA #$F0 (0x8B)
            p.mem.Write(startAddress, 0x8B, 0xF0)
            p.Reg.X = 0x0F
            return 1
        }, func(t *testing.T, p *CPU, name string) {
            assert.Equal(t, byte(0x00), p.Reg.A, name)
            assert.Equal(t, true, p.Reg.IsSet(ZeroFlag), name)
            assert.Equal(t, false, p.Reg.IsSet(NegativeFlag), name)
            assert.Equal(t, uint64(2), p.cycles, name)
        }},

        {"SBC* immediate alias behaves like SBC #imm", func(p *CPU) int {
            // SBC* #$01 (0xEB)
            p.mem.Write(startAddress, 0xEB, 0x01)
            p.Reg.A = 0x10
            p.Reg.SetStatus(CarryFlag, true)
            return 1
        }, func(t *testing.T, p *CPU, name string) {
            assert.Equal(t, byte(0x0F), p.Reg.A, name)
            assert.Equal(t, true, p.Reg.IsSet(CarryFlag), name)
            assert.Equal(t, uint64(2), p.cycles, name)
        }},

        {"NOP* implied consumes 2 cycles and does nothing", func(p *CPU) int {
            // NOP* implied (0x1A)
            p.mem.Write(startAddress, 0x1A)
            return 1
        }, func(t *testing.T, p *CPU, name string) {
            assert.Equal(t, uint64(2), p.cycles, name)
        }},
    }
    executeTests(t, tests)
}
