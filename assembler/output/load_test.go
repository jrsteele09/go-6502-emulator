package output

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/jrsteele09/go-6502-emulator/assembler"
	"github.com/jrsteele09/go-6502-emulator/utils"
)

func TestPRGLoadFile(t *testing.T) {
	// Create test data
	testData := []assembler.AssembledData{
		{
			StartAddress: 0x1000,
			Data:         utils.Value(bytes.NewBuffer([]byte{0xA9, 0x42, 0x8D, 0x20, 0xD0, 0x60})), // LDA #$42, STA $D020, RTS
		},
	}

	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "prg_load_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create PRG file
	prgFile := filepath.Join(tempDir, "test.prg")
	prg := NewPRGFormat()
	err = prg.CreateFile(prgFile, testData, false)
	if err != nil {
		t.Fatalf("Failed to create PRG file: %v", err)
	}

	// Load the file back
	loadedSegments, err := prg.LoadFile(prgFile, false)
	if err != nil {
		t.Fatalf("Failed to load PRG file: %v", err)
	}

	// Verify loaded data
	if len(loadedSegments) != 1 {
		t.Fatalf("Expected 1 segment, got %d", len(loadedSegments))
	}

	segment := loadedSegments[0]
	if segment.StartAddress != testData[0].StartAddress {
		t.Errorf("Expected start address $%04X, got $%04X", testData[0].StartAddress, segment.StartAddress)
	}

	if len(segment.Data.Bytes()) != len(testData[0].Data.Bytes()) {
		t.Errorf("Expected %d bytes, got %d", len(testData[0].Data.Bytes()), len(segment.Data.Bytes()))
	}

	for i, b := range testData[0].Data.Bytes() {
		if segment.Data.Bytes()[i] != b {
			t.Errorf("Data mismatch at byte %d: expected $%02X, got $%02X", i, b, segment.Data.Bytes()[i])
		}
	}
}

func TestPRGLoadFileInvalid(t *testing.T) {
	prg := NewPRGFormat()

	// Test non-existent file
	_, err := prg.LoadFile("nonexistent.prg", false)
	if err == nil {
		t.Error("Expected error for non-existent file")
	}

	// Create empty file
	tempDir, err := os.MkdirTemp("", "prg_load_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	emptyFile := filepath.Join(tempDir, "empty.prg")
	err = os.WriteFile(emptyFile, []byte{}, 0644)
	if err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	}

	// Test empty file
	_, err = prg.LoadFile(emptyFile, false)
	if err == nil {
		t.Error("Expected error for empty file")
	}

	// Test file with only 1 byte (need at least 2 for load address)
	oneByteFile := filepath.Join(tempDir, "onebyte.prg")
	err = os.WriteFile(oneByteFile, []byte{0x00}, 0644)
	if err != nil {
		t.Fatalf("Failed to create one-byte file: %v", err)
	}

	_, err = prg.LoadFile(oneByteFile, false)
	if err == nil {
		t.Error("Expected error for file with only 1 byte")
	}
}

func TestD64LoadFile(t *testing.T) {
	// Create test data
	testData := []assembler.AssembledData{
		{
			StartAddress: 0x1000,
			Data:         utils.Value(bytes.NewBuffer([]byte{0xA9, 0x42, 0x8D, 0x20, 0xD0, 0x60})), // LDA #$42, STA $D020, RTS
		},
	}

	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "d64_load_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create D64 file
	d64File := filepath.Join(tempDir, "test.d64")
	d64 := NewD64Format("TEST DISK", "TD")
	err = d64.CreateFile(d64File, testData, false)
	if err != nil {
		t.Fatalf("Failed to create D64 file: %v", err)
	}

	// Load the file back
	loadedSegments, err := d64.LoadFile(d64File, false)
	if err != nil {
		t.Fatalf("Failed to load D64 file: %v", err)
	}

	// Verify loaded data
	if len(loadedSegments) != 1 {
		t.Fatalf("Expected 1 segment, got %d", len(loadedSegments))
	}

	segment := loadedSegments[0]
	if segment.StartAddress != testData[0].StartAddress {
		t.Errorf("Expected start address $%04X, got $%04X", testData[0].StartAddress, segment.StartAddress)
	}

	if len(segment.Data.Bytes()) != len(testData[0].Data.Bytes()) {
		t.Errorf("Expected %d bytes, got %d", len(testData[0].Data.Bytes()), len(segment.Data.Bytes()))
	}

	for i, b := range testData[0].Data.Bytes() {
		if segment.Data.Bytes()[i] != b {
			t.Errorf("Data mismatch at byte %d: expected $%02X, got $%02X", i, b, segment.Data.Bytes()[i])
		}
	}
}

func TestD64LoadFileMultipleFiles(t *testing.T) {
	// Create test data with multiple segments
	testData := []assembler.AssembledData{
		{
			StartAddress: 0x1000,
			Data:         utils.Value(bytes.NewBuffer([]byte{0xA9, 0x42, 0x8D, 0x20, 0xD0, 0x60})), // LDA #$42, STA $D020, RTS
		},
		{
			StartAddress: 0x2000,
			Data:         utils.Value(bytes.NewBuffer([]byte{0xA9, 0x00, 0x8D, 0x21, 0xD0, 0x60})), // LDA #$00, STA $D021, RTS
		},
	}

	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "d64_load_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create D64 file with multiple files
	d64File := filepath.Join(tempDir, "multi.d64")
	d64 := NewD64Format("MULTI DISK", "MD")
	err = d64.CreateFile(d64File, testData, false)
	if err != nil {
		t.Fatalf("Failed to create D64 file: %v", err)
	}

	// Load the file back
	loadedSegments, err := d64.LoadFile(d64File, false)
	if err != nil {
		t.Fatalf("Failed to load D64 file: %v", err)
	}

	// D64 combines multiple segments into one PRG file, so expect 1 segment
	if len(loadedSegments) != 1 {
		t.Fatalf("Expected 1 combined segment, got %d", len(loadedSegments))
	}

	// Check the combined segment
	segment := loadedSegments[0]
	expectedStartAddr := uint16(0x1000) // Lowest start address
	if segment.StartAddress != expectedStartAddr {
		t.Errorf("Expected start address $%04X, got $%04X", expectedStartAddr, segment.StartAddress)
	}

	// The combined data should have the first segment at offset 0
	// and the second segment at offset 0x1000 (0x2000 - 0x1000)
	expectedSize := 0x1000 + len(testData[1].Data.Bytes()) // gap + second segment
	if len(segment.Data.Bytes()) != expectedSize {
		t.Errorf("Expected %d bytes, got %d", expectedSize, len(segment.Data.Bytes()))
	}

	// Check first segment data (at start of combined data)
	for i, b := range testData[0].Data.Bytes() {
		if segment.Data.Bytes()[i] != b {
			t.Errorf("First segment: Data mismatch at byte %d: expected $%02X, got $%02X", i, b, segment.Data.Bytes()[i])
		}
	}

	// Check second segment data (at offset 0x1000)
	secondSegmentOffset := 0x1000
	for i, b := range testData[1].Data.Bytes() {
		if segment.Data.Bytes()[secondSegmentOffset+i] != b {
			t.Errorf("Second segment: Data mismatch at byte %d: expected $%02X, got $%02X", i, b, segment.Data.Bytes()[secondSegmentOffset+i])
		}
	}
}

func TestT64LoadFile(t *testing.T) {
	// Create test data
	testData := []assembler.AssembledData{
		{
			StartAddress: 0x1000,
			Data:         utils.Value(bytes.NewBuffer([]byte{0xA9, 0x42, 0x8D, 0x20, 0xD0, 0x60})), // LDA #$42, STA $D020, RTS
		},
	}

	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "t64_load_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create T64 file
	t64File := filepath.Join(tempDir, "test.t64")
	t64Format := NewT64Format("TEST TAPE", 30)
	err = t64Format.CreateFile(t64File, testData, false)
	if err != nil {
		t.Fatalf("Failed to create T64 file: %v", err)
	}

	// Load the file back
	loadedSegments, err := t64Format.LoadFile(t64File, false)
	if err != nil {
		t.Fatalf("Failed to load T64 file: %v", err)
	}

	// Verify loaded data
	if len(loadedSegments) != 1 {
		t.Fatalf("Expected 1 segment, got %d", len(loadedSegments))
	}

	segment := loadedSegments[0]
	if segment.StartAddress != testData[0].StartAddress {
		t.Errorf("Expected start address $%04X, got $%04X", testData[0].StartAddress, segment.StartAddress)
	}

	if len(segment.Data.Bytes()) != len(testData[0].Data.Bytes()) {
		t.Errorf("Expected %d bytes, got %d", len(testData[0].Data.Bytes()), len(segment.Data.Bytes()))
	}

	for i, b := range testData[0].Data.Bytes() {
		if segment.Data.Bytes()[i] != b {
			t.Errorf("Data mismatch at byte %d: expected $%02X, got $%02X", i, b, segment.Data.Bytes()[i])
		}
	}
}

func TestT64LoadFileMultipleFiles(t *testing.T) {
	// Create test data with multiple segments
	testData := []assembler.AssembledData{
		{
			StartAddress: 0x1000,
			Data:         utils.Value(bytes.NewBuffer([]byte{0xA9, 0x42, 0x8D, 0x20, 0xD0, 0x60})), // LDA #$42, STA $D020, RTS
		},
		{
			StartAddress: 0x2000,
			Data:         utils.Value(bytes.NewBuffer([]byte{0xA9, 0x00, 0x8D, 0x21, 0xD0, 0x60})), // LDA #$00, STA $D021, RTS
		},
	}

	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "t64_load_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create T64 file with multiple files
	t64File := filepath.Join(tempDir, "multi.t64")
	t64Format := NewT64Format("MULTI TAPE", 30)
	err = t64Format.CreateFile(t64File, testData, false)
	if err != nil {
		t.Fatalf("Failed to create T64 file: %v", err)
	}

	// Load the file back
	loadedSegments, err := t64Format.LoadFile(t64File, false)
	if err != nil {
		t.Fatalf("Failed to load T64 file: %v", err)
	}

	// T64 combines multiple segments into one file, so expect 1 segment
	if len(loadedSegments) != 1 {
		t.Fatalf("Expected 1 combined segment, got %d", len(loadedSegments))
	}

	// Check the combined segment
	segment := loadedSegments[0]
	expectedStartAddr := uint16(0x1000) // Lowest start address
	if segment.StartAddress != expectedStartAddr {
		t.Errorf("Expected start address $%04X, got $%04X", expectedStartAddr, segment.StartAddress)
	}

	// The combined data should have the first segment at offset 0
	// and the second segment at offset 0x1000 (0x2000 - 0x1000)
	expectedSize := 0x1000 + len(testData[1].Data.Bytes()) // gap + second segment
	if len(segment.Data.Bytes()) != expectedSize {
		t.Errorf("Expected %d bytes, got %d", expectedSize, len(segment.Data.Bytes()))
	}

	// Check first segment data (at start of combined data)
	for i, b := range testData[0].Data.Bytes() {
		if segment.Data.Bytes()[i] != b {
			t.Errorf("First segment: Data mismatch at byte %d: expected $%02X, got $%02X", i, b, segment.Data.Bytes()[i])
		}
	}

	// Check second segment data (at offset 0x1000)
	secondSegmentOffset := 0x1000
	for i, b := range testData[1].Data.Bytes() {
		if segment.Data.Bytes()[secondSegmentOffset+i] != b {
			t.Errorf("Second segment: Data mismatch at byte %d: expected $%02X, got $%02X", i, b, segment.Data.Bytes()[secondSegmentOffset+i])
		}
	}
}

func TestLoadFileInvalidFormats(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "load_invalid_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test D64 with invalid size
	d64 := NewD64Format("TEST", "TD")
	invalidD64 := filepath.Join(tempDir, "invalid.d64")
	err = os.WriteFile(invalidD64, []byte{0x00, 0x01, 0x02}, 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid D64: %v", err)
	}

	_, err = d64.LoadFile(invalidD64, false)
	if err == nil {
		t.Error("Expected error for invalid D64 size")
	}

	// Test T64 with invalid signature
	t64Format := NewT64Format("TEST", 30)
	invalidT64 := filepath.Join(tempDir, "invalid.t64")
	err = os.WriteFile(invalidT64, make([]byte, 100), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid T64: %v", err)
	}

	_, err = t64Format.LoadFile(invalidT64, false)
	if err == nil {
		t.Error("Expected error for invalid T64 signature")
	}
}
