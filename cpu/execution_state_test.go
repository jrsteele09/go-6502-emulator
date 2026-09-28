package cpu

import (
	"testing"

	"github.com/jrsteele09/go-6502-emulator/memory"
	"github.com/stretchr/testify/require"
)

func TestResetClearsInstructionScratchpad(t *testing.T) {
	mem := memory.NewMemory[uint16](64 * 1024)
	mem.Write(resetVectorAddr, 0x00, 0x02)
	mem.Write(0x0200, 0xBD, 0xFF, 0x20) // LDA $20FF,X; crosses into $2100.
	mem.Write(0x2100, 0x11)

	processor := NewCPU(mem, false)
	processor.Reg.X = 1

	for range 8 {
		completed, err := processor.Execute()
		require.NoError(t, err)
		require.False(t, bool(completed))
		if processor.execute.phase != 0 {
			break
		}
	}
	require.NotZero(t, processor.execute.phase)
	require.Equal(t, uint16(0x2100), processor.execute.address)

	mem.Write(resetVectorAddr, 0x00, 0x03)
	mem.Write(0x0300, 0xBD, 0x00, 0x20) // Same opcode, no page crossing.
	mem.Write(0x2001, 0x42)
	processor.Reset()

	require.Equal(t, executionFetch, processor.execute.stage)
	require.Zero(t, processor.execute.phase)
	require.Zero(t, processor.execute.address)
	require.Nil(t, processor.execute.opcode)

	for range 8 {
		completed, err := processor.Execute()
		require.NoError(t, err)
		if completed {
			break
		}
	}
	require.Equal(t, byte(0x42), processor.Reg.A)
}

func TestOpcodeSelectionInitializesInstructionScratchpad(t *testing.T) {
	mem := memory.NewMemory[uint16](64 * 1024)
	mem.Write(resetVectorAddr, 0x00, 0x02)
	mem.Write(0x0200, 0xEA)
	processor := NewCPU(mem, false)

	processor.execute.phase = 2
	processor.execute.address = 0xCAFE
	processor.execute.pageCrossed = true

	completed, err := processor.Execute()
	require.NoError(t, err)
	require.False(t, bool(completed))
	require.Equal(t, executionInstruction, processor.execute.stage)
	require.Zero(t, processor.execute.phase)
	require.Zero(t, processor.execute.address)
	require.False(t, processor.execute.pageCrossed)
}
