package bin

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/jrsteele09/go-6502-emulator/assembler"
	"github.com/jrsteele09/go-6502-emulator/utils"
)

func TestT64Format_CreateData(t *testing.T) {
	t64 := NewT64Format("TEST TAPE", 30)

	tests := []struct {
		name     string
		segments []assembler.AssembledData
		wantErr  bool
	}{
		{
			name: "single segment",
			segments: []assembler.AssembledData{
				{StartAddress: 0x1000, Data: utils.Value(bytes.NewBuffer([]byte{0xA9, 0x10, 0x60}))}, // LDA #$10, RTS
			},
			wantErr: false,
		},
		{
			name: "multiple segments",
			segments: []assembler.AssembledData{
				{StartAddress: 0x1000, Data: utils.Value(bytes.NewBuffer([]byte{0xA9, 0x10}))}, // LDA #$10
				{StartAddress: 0x1005, Data: utils.Value(bytes.NewBuffer([]byte{0x60}))},       // RTS
			},
			wantErr: false,
		},
		{
			name:     "no segments",
			segments: []assembler.AssembledData{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := t64.CreateData(tt.segments)
			if (err != nil) != tt.wantErr {
				t.Errorf("T64Format.CreateData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Verify T64 format structure
				expectedMinSize := T64HeaderSize + int(t64.MaxEntries)*T64EntrySize + 5 // Minimum size with small program
				if len(result) < expectedMinSize {
					t.Errorf("T64Format.CreateData() result length = %d, expected at least %d", len(result), expectedMinSize)
				}

				// Verify T64 signature
				signature := string(result[0:len(T64Signature)])
				if signature != T64Signature {
					t.Errorf("T64Format.CreateData() signature = %s, expected %s", signature, T64Signature)
				}

				// Verify version
				version := uint16(result[32]) | (uint16(result[33]) << 8)
				if version != T64Version {
					t.Errorf("T64Format.CreateData() version = 0x%04X, expected 0x%04X", version, T64Version)
				}

				// Verify max entries
				maxEntries := uint16(result[34]) | (uint16(result[35]) << 8)
				if maxEntries != t64.MaxEntries {
					t.Errorf("T64Format.CreateData() max entries = %d, expected %d", maxEntries, t64.MaxEntries)
				}

				// Verify used entries (should be 1 after creation)
				usedEntries := uint16(result[36]) | (uint16(result[37]) << 8)
				expectedUsedEntries := uint16(1) // We create one file entry
				if usedEntries != expectedUsedEntries {
					t.Errorf("T64Format.CreateData() used entries = %d, expected %d", usedEntries, expectedUsedEntries)
				}

				// Verify tape name
				tapeName := string(result[40 : 40+len(t64.TapeName)])
				if tapeName != t64.TapeName {
					t.Errorf("T64Format.CreateData() tape name = %s, expected %s", tapeName, t64.TapeName)
				}

				// Verify directory entry
				entryOffset := T64HeaderSize
				if result[entryOffset] != 1 { // Entry type
					t.Errorf("T64Format.CreateData() entry type = %d, expected 1", result[entryOffset])
				}
				if result[entryOffset+1] != 0x82 { // File type (PRG)
					t.Errorf("T64Format.CreateData() file type = 0x%02X, expected 0x82", result[entryOffset+1])
				}

				// Verify start address matches first segment
				startAddr := uint16(result[entryOffset+2]) | (uint16(result[entryOffset+3]) << 8)
				expectedStartAddr := tt.segments[0].StartAddress
				if len(tt.segments) > 1 {
					// Find the lowest start address
					for _, segment := range tt.segments {
						if segment.StartAddress < expectedStartAddr {
							expectedStartAddr = segment.StartAddress
						}
					}
				}
				if startAddr != expectedStartAddr {
					t.Errorf("T64Format.CreateData() start address = 0x%04X, expected 0x%04X", startAddr, expectedStartAddr)
				}

				// Verify program data exists
				dataOffset := int(result[entryOffset+8]) | (int(result[entryOffset+9]) << 8) |
					(int(result[entryOffset+10]) << 16) | (int(result[entryOffset+11]) << 24)
				if dataOffset >= len(result) {
					t.Errorf("T64Format.CreateData() data offset %d exceeds archive size %d", dataOffset, len(result))
				}

				// Verify PRG format in data (load address)
				if len(result) > dataOffset+1 {
					loadAddr := uint16(result[dataOffset]) | (uint16(result[dataOffset+1]) << 8)
					if loadAddr != expectedStartAddr {
						t.Errorf("T64Format.CreateData() PRG load address = 0x%04X, expected 0x%04X", loadAddr, expectedStartAddr)
					}
				}
			}
		})
	}
}

func TestT64Format_CreateFile(t *testing.T) {
	t64 := NewT64Format("TEST TAPE", 30)
	segments := []assembler.AssembledData{
		{StartAddress: 0x1000, Data: utils.Value(bytes.NewBuffer([]byte{0xA9, 0x10, 0x60}))}, // LDA #$10, RTS
	}

	// Create temporary file
	tmpFile := "test_output.t64"
	defer os.Remove(tmpFile)

	err := t64.CreateFile(tmpFile, segments, false)
	if err != nil {
		t.Fatalf("T64Format.CreateFile() error = %v", err)
	}

	// Read the file back
	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	// Verify T64 signature
	if len(data) < len(T64Signature) {
		t.Fatalf("T64Format.CreateFile() file too small")
	}

	signature := string(data[0:len(T64Signature)])
	if signature != T64Signature {
		t.Errorf("T64Format.CreateFile() signature = %s, expected %s", signature, T64Signature)
	}

	// Verify it contains our program data
	// The load address should be present in the PRG data section
	found := false
	for i := 0; i < len(data)-1; i++ {
		if data[i] == 0x00 && data[i+1] == 0x10 { // Load address 0x1000 in little-endian
			found = true
			break
		}
	}
	if !found {
		t.Error("T64Format.CreateFile() load address not found in archive")
	}
}

func TestT64Format_CreateFileNoSegments(t *testing.T) {
	t64 := NewT64Format("TEST TAPE", 30)
	err := t64.CreateFile("dummy.t64", []assembler.AssembledData{}, false)
	if err == nil {
		t.Error("T64Format.CreateFile() with no segments should return error")
	}
}

func TestT64Format_ImplementsBinaryFormat(t *testing.T) {
	var _ BinaryFormat = (*T64Format)(nil)
	// If this compiles, the interface is implemented correctly
}

func TestNewT64Format(t *testing.T) {
	tests := []struct {
		name            string
		tapeName        string
		maxEntries      uint16
		expectedName    string
		expectedEntries uint16
	}{
		{
			name:            "normal parameters",
			tapeName:        "TEST TAPE",
			maxEntries:      30,
			expectedName:    "TEST TAPE",
			expectedEntries: 30,
		},
		{
			name:            "long tape name",
			tapeName:        "THIS IS A VERY LONG TAPE NAME THAT EXCEEDS LIMITS",
			maxEntries:      50,
			expectedName:    "THIS IS A VERY LONG TAPE", // Truncated to 24 chars
			expectedEntries: 50,
		},
		{
			name:            "zero max entries",
			tapeName:        "TEST",
			maxEntries:      0,
			expectedName:    "TEST",
			expectedEntries: 30, // Default
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t64 := NewT64Format(tt.tapeName, tt.maxEntries)
			if t64.TapeName != tt.expectedName {
				t.Errorf("NewT64Format() tape name = %s, expected %s", t64.TapeName, tt.expectedName)
			}
			if t64.MaxEntries != tt.expectedEntries {
				t.Errorf("NewT64Format() max entries = %d, expected %d", t64.MaxEntries, tt.expectedEntries)
			}
		})
	}
}

func TestT64Format_CreatePRGFromSegments(t *testing.T) {
	t64 := NewT64Format("TEST", 30)

	tests := []struct {
		name             string
		segments         []assembler.AssembledData
		expectedLoadAddr uint16
		expectedSize     int
	}{
		{
			name: "single segment",
			segments: []assembler.AssembledData{
				{StartAddress: 0x1000, Data: utils.Value(bytes.NewBuffer([]byte{0xA9, 0x10}))}}, // LDA #$10
			expectedLoadAddr: 0x1000,
			expectedSize:     4, // 2 bytes load address + 2 bytes data
		},
		{
			name: "multiple segments",
			segments: []assembler.AssembledData{
				{StartAddress: 0x1000, Data: utils.Value(bytes.NewBuffer([]byte{0xA9, 0x10}))}, // LDA #$10
				{StartAddress: 0x1005, Data: utils.Value(bytes.NewBuffer([]byte{0x60}))},       // RTS
			},
			expectedLoadAddr: 0x1000,
			expectedSize:     8, // 2 bytes load address + 6 bytes (including gap)
		},
		{
			name: "segments out of order",
			segments: []assembler.AssembledData{
				{StartAddress: 0x2000, Data: utils.Value(bytes.NewBuffer([]byte{0x60}))},       // RTS
				{StartAddress: 0x1000, Data: utils.Value(bytes.NewBuffer([]byte{0xA9, 0x10}))}, // LDA #$10
			},
			expectedLoadAddr: 0x1000, // Should use lowest address
			expectedSize:     4099,   // 2 bytes load address + large gap
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prgData, loadAddr, err := t64.createPRGFromSegments(tt.segments)
			if err != nil {
				t.Errorf("createPRGFromSegments() error = %v", err)
				return
			}

			if loadAddr != tt.expectedLoadAddr {
				t.Errorf("createPRGFromSegments() load address = 0x%04X, expected 0x%04X", loadAddr, tt.expectedLoadAddr)
			}

			if len(prgData) != tt.expectedSize {
				t.Errorf("createPRGFromSegments() PRG size = %d, expected %d", len(prgData), tt.expectedSize)
			}

			// Verify load address in PRG data
			if len(prgData) >= 2 {
				prgLoadAddr := uint16(prgData[0]) | (uint16(prgData[1]) << 8)
				if prgLoadAddr != tt.expectedLoadAddr {
					t.Errorf("createPRGFromSegments() PRG load address = 0x%04X, expected 0x%04X", prgLoadAddr, tt.expectedLoadAddr)
				}
			}
		})
	}
}

func TestT64Format_PadString(t *testing.T) {
	t64 := NewT64Format("TEST", 30)

	tests := []struct {
		input    string
		length   int
		expected string
	}{
		{input: "TEST", length: 8, expected: "TEST    "},
		{input: "HELLO", length: 5, expected: "HELLO"},
		{input: "A", length: 3, expected: "A  "},
		{input: "", length: 4, expected: "    "},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s_len_%d", tt.input, tt.length), func(t *testing.T) {
			result := t64.padString(tt.input, tt.length)
			resultStr := string(result)
			if resultStr != tt.expected {
				t.Errorf("padString(%s, %d) = %s, expected %s", tt.input, tt.length, resultStr, tt.expected)
			}
			if len(result) != tt.length {
				t.Errorf("padString(%s, %d) length = %d, expected %d", tt.input, tt.length, len(result), tt.length)
			}
		})
	}
}

func TestT64Format_WriteHeader(t *testing.T) {
	t64 := NewT64Format("MY TAPE", 50)
	archive := make([]byte, T64HeaderSize)

	t64.writeHeader(archive)

	// Verify signature
	signature := string(archive[0:len(T64Signature)])
	if signature != T64Signature {
		t.Errorf("writeHeader() signature = %s, expected %s", signature, T64Signature)
	}

	// Verify version
	version := uint16(archive[32]) | (uint16(archive[33]) << 8)
	if version != T64Version {
		t.Errorf("writeHeader() version = 0x%04X, expected 0x%04X", version, T64Version)
	}

	// Verify max entries
	maxEntries := uint16(archive[34]) | (uint16(archive[35]) << 8)
	if maxEntries != t64.MaxEntries {
		t.Errorf("writeHeader() max entries = %d, expected %d", maxEntries, t64.MaxEntries)
	}

	// Verify tape name
	tapeName := string(archive[40 : 40+len(t64.TapeName)])
	if tapeName != t64.TapeName {
		t.Errorf("writeHeader() tape name = %s, expected %s", tapeName, t64.TapeName)
	}
}

func TestT64Format_WriteDirectoryEntry(t *testing.T) {
	t64 := NewT64Format("TEST", 30)
	totalSize := T64HeaderSize + int(t64.MaxEntries)*T64EntrySize + 100
	archive := make([]byte, totalSize)

	filename := "TESTPROG"
	loadAddr := uint16(0x1000)
	fileSize := uint16(50)
	dataOffset := T64HeaderSize + int(t64.MaxEntries)*T64EntrySize

	t64.writeDirectoryEntry(archive, filename, loadAddr, fileSize, dataOffset)

	entryOffset := T64HeaderSize

	// Verify entry type
	if archive[entryOffset] != 1 {
		t.Errorf("writeDirectoryEntry() entry type = %d, expected 1", archive[entryOffset])
	}

	// Verify file type
	if archive[entryOffset+1] != 0x82 {
		t.Errorf("writeDirectoryEntry() file type = 0x%02X, expected 0x82", archive[entryOffset+1])
	}

	// Verify start address
	startAddr := uint16(archive[entryOffset+2]) | (uint16(archive[entryOffset+3]) << 8)
	if startAddr != loadAddr {
		t.Errorf("writeDirectoryEntry() start address = 0x%04X, expected 0x%04X", startAddr, loadAddr)
	}

	// Verify end address
	// End address should be: loadAddr + (fileSize - 2) - 1 = 0x1000 + 48 - 1 = 0x102F
	expectedEndAddr := loadAddr + (fileSize - 2) - 1
	endAddr := uint16(archive[entryOffset+4]) | (uint16(archive[entryOffset+5]) << 8)
	if endAddr != expectedEndAddr {
		t.Errorf("writeDirectoryEntry() end address = 0x%04X, expected 0x%04X", endAddr, expectedEndAddr)
	}

	// Verify data offset
	actualDataOffset := int(archive[entryOffset+8]) | (int(archive[entryOffset+9]) << 8) |
		(int(archive[entryOffset+10]) << 16) | (int(archive[entryOffset+11]) << 24)
	if actualDataOffset != dataOffset {
		t.Errorf("writeDirectoryEntry() data offset = %d, expected %d", actualDataOffset, dataOffset)
	}

	// Verify filename
	storedFilename := string(archive[entryOffset+16 : entryOffset+16+len(filename)])
	if storedFilename != filename {
		t.Errorf("writeDirectoryEntry() filename = %s, expected %s", storedFilename, filename)
	}
}
