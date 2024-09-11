package assembler_test

import (
	"strings"
	"testing"

	"github.com/jrsteele09/go-6502-emulator/assembler"
	"github.com/jrsteele09/go-6502-emulator/cpu"
	"github.com/jrsteele09/go-6502-emulator/memory"
	"github.com/stretchr/testify/require"
)

func TestAssembler(t *testing.T) {

	m := memory.NewMemory[uint16](64 * 1024)
	c := cpu.NewCPU(m)

	assembler := assembler.New(cpu.OpCodes(c))
	// asmCode := `
	// // LDA #-5
	// BNE *-127

	// LDA #255
	// STA $0200
	// LDA #$05
	// STA $0201
	// BRK
	// `

	asmCode := `
	;BNE *-127`

	err := assembler.Assemble(strings.NewReader(asmCode))
	require.NoError(t, err)

}
