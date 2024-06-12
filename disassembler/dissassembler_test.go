package disassembler

import (
	"testing"

	cpu "github.com/jrsteele09/go-6502-emulator/cpu"
	"github.com/jrsteele09/go-6502-emulator/memory"
	"github.com/stretchr/testify/assert"
)

func TestDisassembler(t *testing.T) {
	m := memory.NewMemory[uint16](64 * 1024)
	p := cpu.NewCpu(m)
	m.Write(0xC000,
		0x02,
		0xA9, 0x01,
		0xA9, 0x80,
		0xA9, 0x00,
		0xA5, 0x80,
		0xB5, 0x80,
		0xAD, 0x80, 0x00,
		0xBD, 0x80, 0x00,
		0xBD, 0x01, 0x00,
		0xB9, 0x80, 0x00,
		0xB9, 0x01, 0x00,
		0xA1, 0x05,
		0xB1, 0x05,
	) //0xc01d

	expectedDissassmbledCode := "C000:  02         ???\nC001:  A9 01      LDA #$01\nC003:  A9 80      LDA #$80\nC005:  A9 00      LDA #$00\nC007:  A5 80      LDA $80\nC009:  B5 80      LDA $80,X\nC00B:  AD 80 00   LDA $0080\nC00E:  BD 80 00   LDA $0080,X\nC011:  BD 01 00   LDA $0001,X\nC014:  B9 80 00   LDA $0080,Y\nC017:  B9 01 00   LDA $0001,Y\nC01A:  A1 05      LDA ($05,X)\nC01C:  B1 05      LDA ($05),Y\n"

	dissassembler := NewDissassembler(m, cpu.OpCodes(p))

	dissassembledCode := ""

	address := uint16(0xC000)
	for address < uint16(0xC01D) {
		line, bytes := dissassembler.Disassemble(address)
		dissassembledCode += line + "\n"
		address += uint16(bytes)
	}

	assert.Equal(t, expectedDissassmbledCode, dissassembledCode)
}

func TestDisassemblerRelativeAddressingMode(t *testing.T) {
	m := memory.NewMemory[uint16](64 * 1024)
	p := cpu.NewCpu(m)
	m.Write(0xC000, 0xF0, byte(0xAF))
	dissassembler := NewDissassembler(m, cpu.OpCodes(p))

	address := uint16(0xC000)
	dissassembledCode, _ := dissassembler.Disassemble(address)

	assert.Equal(t, "C000:  F0 AF      BEQ BFAF", dissassembledCode)

}
