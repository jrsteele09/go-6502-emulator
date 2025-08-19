package bin

import (
	"bytes"
	"os"
	"testing"

	"github.com/jrsteele09/go-6502-emulator/assembler"
	"github.com/jrsteele09/go-6502-emulator/utils"
)

func TestPRGFormat_CreateData(t *testing.T) {
	prg := NewPRGFormat()

	tests := []struct {
		name     string
		segments []assembler.AssembledData
		expected []byte
		wantErr  bool
	}{
		{
			name: "single segment",
			segments: []assembler.AssembledData{
				{StartAddress: 0x1000, Data: utils.Value(bytes.NewBuffer([]byte{0xA9, 0x10}))}, // LDA #$10
			},
			expected: []byte{0x00, 0x10, 0xA9, 0x10}, // Load addr $1000 + data
			wantErr:  false,
		},
		{
			name: "multiple segments",
			segments: []assembler.AssembledData{
				{StartAddress: 0x1000, Data: utils.Value(bytes.NewBuffer([]byte{0xA9, 0x10}))}, // LDA #$10
				{StartAddress: 0x1005, Data: utils.Value(bytes.NewBuffer([]byte{0x60}))},       // RTS
			},
			expected: nil,
			wantErr:  true, // PRG format does not support multiple segments
		},
		{
			name: "lowest address first",
			segments: []assembler.AssembledData{
				{StartAddress: 0x2000, Data: utils.Value(bytes.NewBuffer([]byte{0x60}))},       // RTS
				{StartAddress: 0x1000, Data: utils.Value(bytes.NewBuffer([]byte{0xA9, 0x10}))}, // LDA #$10
			},
			expected: nil,
			wantErr:  true, // Multiple segments now error
		},
		{
			name:     "no segments",
			segments: []assembler.AssembledData{},
			expected: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := prg.CreateData(tt.segments)
			if (err != nil) != tt.wantErr {
				t.Errorf("PRGFormat.CreateData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// For simple tests, compare directly if expected is provided
				if tt.expected != nil {
					if len(result) != len(tt.expected) {
						t.Errorf("PRGFormat.CreateData() result length = %d, expected %d", len(result), len(tt.expected))
						return
					}
					for i, b := range tt.expected {
						if result[i] != b {
							t.Errorf("PRGFormat.CreateData() result[%d] = %02X, expected %02X", i, result[i], b)
						}
					}
				}
			}
		})
	}
}

func TestPRGFormat_CreateFile(t *testing.T) {
	prg := NewPRGFormat()
	segments := []assembler.AssembledData{
		{StartAddress: 0x1000, Data: utils.Value(bytes.NewBuffer([]byte{0xA9, 0x10}))}, // LDA #$10
	}

	// Create temporary file
	tmpFile := "test_output.prg"
	defer os.Remove(tmpFile)

	err := prg.CreateFile(tmpFile, segments, false)
	if err != nil {
		t.Fatalf("PRGFormat.CreateFile() error = %v", err)
	}

	// Read the file back
	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	expected := []byte{0x00, 0x10, 0xA9, 0x10} // Load addr + data
	if len(data) != len(expected) {
		t.Errorf("File length = %d, expected %d", len(data), len(expected))
		return
	}

	for i, b := range expected {
		if data[i] != b {
			t.Errorf("File data[%d] = %02X, expected %02X", i, data[i], b)
		}
	}
}

func TestPRGFormat_CreateFileNoSegments(t *testing.T) {
	prg := NewPRGFormat()
	err := prg.CreateFile("dummy.prg", []assembler.AssembledData{}, false)
	if err == nil {
		t.Error("PRGFormat.CreateFile() with no segments should return error")
	}
}
