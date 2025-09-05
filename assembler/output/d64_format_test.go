package output

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/jrsteele09/go-6502-emulator/assembler"
	"github.com/jrsteele09/go-6502-emulator/utils"
)

func TestD64Format_CreateData(t *testing.T) {
	d64 := NewD64Format("TEST DISK", "01")

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
			result, err := d64.CreateData(tt.segments)
			if (err != nil) != tt.wantErr {
				t.Errorf("D64Format.CreateData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Verify D64 format
				if len(result) != D64TotalSize {
					t.Errorf("D64Format.CreateData() result length = %d, expected %d", len(result), D64TotalSize)
				}

				// Verify BAM header
				bamOffset := d64.trackSectorToOffset(18, 0)
				if result[bamOffset] != 18 || result[bamOffset+1] != 1 {
					t.Errorf("D64Format.CreateData() BAM header incorrect: track=%d, sector=%d", result[bamOffset], result[bamOffset+1])
				}

				// Verify disk name in BAM
				diskNameOffset := bamOffset + 0x90
				diskNameInBAM := string(result[diskNameOffset : diskNameOffset+len(d64.DiskName)])
				if diskNameInBAM != d64.DiskName {
					t.Errorf("D64Format.CreateData() disk name in BAM = %s, expected %s", diskNameInBAM, d64.DiskName)
				}

				// Verify disk ID in BAM
				diskIDOffset := bamOffset + 0xA2
				diskIDInBAM := string(result[diskIDOffset : diskIDOffset+D64DiskIDLength])
				if diskIDInBAM != d64.DiskID {
					t.Errorf("D64Format.CreateData() disk ID in BAM = %s, expected %s", diskIDInBAM, d64.DiskID)
				}

				// Verify directory structure
				dirOffset := d64.trackSectorToOffset(18, 1)
				if result[dirOffset] != 0 || result[dirOffset+1] != 255 {
					t.Errorf("D64Format.CreateData() directory header incorrect: track=%d, sector=%d", result[dirOffset], result[dirOffset+1])
				}

				// Verify directory entry exists
				entryOffset := dirOffset + 2
				if result[entryOffset] != 0x82 { // PRG file type
					t.Errorf("D64Format.CreateData() directory entry file type = %02X, expected 82", result[entryOffset])
				}

				// Verify program file exists at track 1, sector 0
				fileOffset := d64.trackSectorToOffset(1, 0)
				if result[fileOffset+2] == 0 && result[fileOffset+3] == 0 {
					t.Error("D64Format.CreateData() program file data appears to be empty")
				}
			}
		})
	}
}

func TestD64Format_CreateFile(t *testing.T) {
	d64 := NewD64Format("TEST DISK", "01")
	segments := []assembler.AssembledData{
		{StartAddress: 0x1000, Data: utils.Value(bytes.NewBuffer([]byte{0xA9, 0x10, 0x60}))}, // LDA #$10, RTS
	}

	// Create temporary file
	tmpFile := "test_output.d64"
	defer os.Remove(tmpFile)

	err := d64.CreateFile(tmpFile, segments, false)
	if err != nil {
		t.Fatalf("D64Format.CreateFile() error = %v", err)
	}

	// Read the file back
	data, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	// Verify file size
	if len(data) != D64TotalSize {
		t.Errorf("D64Format.CreateFile() file size = %d, expected %d", len(data), D64TotalSize)
	}

	// Verify it's a valid D64 (check BAM)
	bamOffset := d64.trackSectorToOffset(18, 0)
	if data[bamOffset] != 18 || data[bamOffset+1] != 1 {
		t.Errorf("D64Format.CreateFile() invalid BAM header")
	}
}

func TestD64Format_CreateFileNoSegments(t *testing.T) {
	d64 := NewD64Format("TEST DISK", "01")
	err := d64.CreateFile("dummy.d64", []assembler.AssembledData{}, false)
	if err == nil {
		t.Error("D64Format.CreateFile() with no segments should return error")
	}
}

func TestD64Format_ImplementsBinaryFormat(t *testing.T) {
	var _ BinaryFormat = (*D64Format)(nil)
	// If this compiles, the interface is implemented correctly
}

func TestNewD64Format(t *testing.T) {
	tests := []struct {
		name         string
		diskName     string
		diskID       string
		expectedName string
		expectedID   string
	}{
		{
			name:         "normal parameters",
			diskName:     "TEST DISK",
			diskID:       "01",
			expectedName: "TEST DISK",
			expectedID:   "01",
		},
		{
			name:         "long disk name",
			diskName:     "THIS IS A VERY LONG DISK NAME",
			diskID:       "01",
			expectedName: "THIS IS A VERY L", // Truncated to 16 chars
			expectedID:   "01",
		},
		{
			name:         "invalid disk ID",
			diskName:     "TEST",
			diskID:       "ABC", // Too long
			expectedName: "TEST",
			expectedID:   "01", // Default
		},
		{
			name:         "empty disk ID",
			diskName:     "TEST",
			diskID:       "",
			expectedName: "TEST",
			expectedID:   "01", // Default
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d64 := NewD64Format(tt.diskName, tt.diskID)
			if d64.DiskName != tt.expectedName {
				t.Errorf("NewD64Format() disk name = %s, expected %s", d64.DiskName, tt.expectedName)
			}
			if d64.DiskID != tt.expectedID {
				t.Errorf("NewD64Format() disk ID = %s, expected %s", d64.DiskID, tt.expectedID)
			}
		})
	}
}

func TestD64Format_TrackSectorToOffset(t *testing.T) {
	d64 := NewD64Format("TEST", "01")

	tests := []struct {
		track          uint8
		sector         uint8
		expectedOffset int
	}{
		{track: 1, sector: 0, expectedOffset: 0},                  // First block
		{track: 1, sector: 1, expectedOffset: 256},                // Second block
		{track: 2, sector: 0, expectedOffset: 21 * 256},           // First block of track 2
		{track: 18, sector: 0, expectedOffset: 17*21*256 + 0*256}, // BAM location
		{track: 18, sector: 1, expectedOffset: 17*21*256 + 1*256}, // Directory location
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("track_%d_sector_%d", tt.track, tt.sector), func(t *testing.T) {
			offset := d64.trackSectorToOffset(tt.track, tt.sector)
			if offset != tt.expectedOffset {
				t.Errorf("trackSectorToOffset(%d, %d) = %d, expected %d", tt.track, tt.sector, offset, tt.expectedOffset)
			}
		})
	}
}

func TestD64Format_GetSectorsPerTrack(t *testing.T) {
	d64 := NewD64Format("TEST", "01")

	tests := []struct {
		track    uint8
		expected uint8
	}{
		{track: 1, expected: 21},
		{track: 17, expected: 21},
		{track: 18, expected: 19},
		{track: 24, expected: 19},
		{track: 25, expected: 18},
		{track: 30, expected: 18},
		{track: 31, expected: 17},
		{track: 35, expected: 17},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("track_%d", tt.track), func(t *testing.T) {
			sectors := d64.getSectorsPerTrack(tt.track)
			if sectors != tt.expected {
				t.Errorf("getSectorsPerTrack(%d) = %d, expected %d", tt.track, sectors, tt.expected)
			}
		})
	}
}

func TestD64Format_PadString(t *testing.T) {
	d64 := NewD64Format("TEST", "01")

	tests := []struct {
		input    string
		length   int
		expected string
	}{
		{input: "TEST", length: 8, expected: "TEST    "},
		{input: "HELLO", length: 5, expected: "HELLO"},
		{input: "A", length: 3, expected: "A  "},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := d64.padString(tt.input, tt.length)
			resultStr := string(result)
			if resultStr != tt.expected {
				t.Errorf("padString(%s, %d) = %s, expected %s", tt.input, tt.length, resultStr, tt.expected)
			}
		})
	}
}

func TestD64Format_PadStringShifted(t *testing.T) {
	d64 := NewD64Format("TEST", "01")

	result := d64.padStringShifted("TEST", 8)

	// Check that the string part is correct
	if string(result[:4]) != "TEST" {
		t.Errorf("padStringShifted() string part incorrect: got %s", string(result[:4]))
	}

	// Check that padding is shifted spaces (0xA0)
	for i := 4; i < 8; i++ {
		if result[i] != 0xA0 {
			t.Errorf("padStringShifted() padding byte %d = 0x%02X, expected 0xA0", i, result[i])
		}
	}
}
