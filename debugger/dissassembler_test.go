package debugger

import (
	"strings"
	"testing"

	cpu "github.com/jrsteele09/go-6502-emulator/cpu"
	"github.com/jrsteele09/go-6502-emulator/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDisassembler(t *testing.T) {
	m := memory.NewMemory[uint16](64 * 1024)
	p := cpu.NewCPU(m, false)
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

	expectedDissassmbledCode :=
		`$C000: 02         ???
$C001: A9 01      LDA #$01
$C003: A9 80      LDA #$80
$C005: A9 00      LDA #$00
$C007: A5 80      LDA $80
$C009: B5 80      LDA $80,X
$C00B: AD 80 00   LDA $0080
$C00E: BD 80 00   LDA $0080,X
$C011: BD 01 00   LDA $0001,X
$C014: B9 80 00   LDA $0080,Y
$C017: B9 01 00   LDA $0001,Y
$C01A: A1 05      LDA ($05,X)
$C01C: B1 05      LDA ($05),Y`

	dissassembler := NewDisassembler(m, p.OpCodes())
	dissassembledCode := ""

	address := uint16(0xC000)
	for address < uint16(0xC01D) {
		line, bytes := dissassembler.Disassemble(address)
		dissassembledCode += line + "\n"
		address += uint16(bytes)
	}
	compareDisassembly(t, expectedDissassmbledCode, dissassembledCode)
}

func TestDisassemblerRelativeAddressingMode(t *testing.T) {
	m := memory.NewMemory[uint16](64 * 1024)
	p := cpu.NewCPU(m, false)
	m.Write(0xC000, 0xF0, byte(0xAF))
	dissassembler := NewDisassembler(m, p.OpCodes())

	address := uint16(0xC000)
	dissassembledCode, _ := dissassembler.Disassemble(address)

	assert.Equal(t, "$C000: F0 AF      BEQ $BFB1", dissassembledCode)

}

func compareDisassembly(t *testing.T, expected, actual string) {
	expectedArray := strings.Split(strings.TrimSpace(expected), "\n")
	actualArray := strings.Split(strings.TrimSpace(actual), "\n")
	for i := range expectedArray {
		require.Equal(t, expectedArray[i], actualArray[i], "Mismatch at line %d: expected %q, got %q", i+1, expectedArray[i], actualArray[i])
	}
}
