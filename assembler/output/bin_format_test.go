package output

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/jrsteele09/go-6502-emulator/assembler"
	"github.com/jrsteele09/go-6502-emulator/utils"
	"github.com/stretchr/testify/require"
)

func TestBINFormatCreateData(t *testing.T) {
	format := NewBINFormat(0)
	segments := []assembler.AssembledData{
		{StartAddress: 0x1003, Data: utils.Value(bytes.NewBuffer([]byte{0x60}))},
		{StartAddress: 0x1000, Data: utils.Value(bytes.NewBuffer([]byte{0xA9, 0x42}))},
	}

	data, err := format.CreateData(segments)

	require.NoError(t, err)
	require.Equal(t, []byte{0xA9, 0x42, 0x00, 0x60}, data)
}

func TestBINFormatCreateDataRejectsOverlappingSegments(t *testing.T) {
	format := NewBINFormat(0)
	segments := []assembler.AssembledData{
		{StartAddress: 0x1000, Data: utils.Value(bytes.NewBuffer([]byte{0xA9, 0x42}))},
		{StartAddress: 0x1001, Data: utils.Value(bytes.NewBuffer([]byte{0x60}))},
	}

	_, err := format.CreateData(segments)

	require.ErrorContains(t, err, "segments overlap")
}

func TestBINFormatCreateAndLoadFile(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "program.bin")
	writer := NewBINFormat(0)
	segments := []assembler.AssembledData{
		{StartAddress: 0x2000, Data: utils.Value(bytes.NewBuffer([]byte{0xA9, 0x42, 0x60}))},
	}
	require.NoError(t, writer.CreateFile(filename, segments, false))

	data, err := os.ReadFile(filename)
	require.NoError(t, err)
	require.Equal(t, []byte{0xA9, 0x42, 0x60}, data)

	loaded, err := NewBINFormat(0x2000).LoadFile(filename, false)
	require.NoError(t, err)
	require.Len(t, loaded, 1)
	require.Equal(t, uint16(0x2000), loaded[0].StartAddress)
	require.Equal(t, data, loaded[0].Data.Bytes())
}

func TestBINFormatLoadFileRejectsDataBeyondMemory(t *testing.T) {
	filename := filepath.Join(t.TempDir(), "program.bin")
	require.NoError(t, os.WriteFile(filename, []byte{0xA9, 0x42}, 0o644))

	_, err := NewBINFormat(0xFFFF).LoadFile(filename, false)

	require.ErrorContains(t, err, "does not fit")
}
