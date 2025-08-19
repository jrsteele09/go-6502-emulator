package assembler

import (
	"strings"
	"testing"

	"github.com/jrsteele09/go-6502-emulator/cpu"
	"github.com/jrsteele09/go-6502-emulator/memory"
)

// createTestOpcodes creates a real 6502 opcode set for testing
func createTestOpcodes() []*cpu.OpCodeDef {
	mem := memory.NewMemory[uint16](64 * 1024)
	testCPU := cpu.NewCPU(mem)
	return testCPU.OpCodes()
}

func TestAssembler_WithPreprocessorIncludes(t *testing.T) {
	// Use real CPU opcodes instead of manually defining them
	opcodes := createTestOpcodes()
	assembler := New(opcodes)

	// Test files with include directives
	files := map[string]string{
		"main.asm": `; Main program
LDA #$10
#include "subroutines.asm"
.include "data.asm"
STA $1000`,

		"subroutines.asm": `; Subroutines
compare:
    CMP #$20
    #include "utils.asm"`,

		"utils.asm": `; Utility functions
increment:
    INX
    INY`,

		"data.asm": `; Data section
.INCLUDE "constants.asm"`,

		"constants.asm": `; Constants
; value: .EQU $FF`,
	}

	resolver := NewMemoryFileResolver(files)

	// Test the assembly with preprocessor
	segments, err := assembler.AssembleWithPreprocessor(strings.NewReader(files["main.asm"]), resolver)
	if err != nil {
		t.Fatalf("Assembly failed: %v", err)
	}
	if len(segments) == 0 {
		t.Fatal("Expected at least one segment, got none")
	}
	firstSegment := segments[0]
	if len(firstSegment.Data.Bytes()) < 7 {
		t.Errorf("Expected at least 7 bytes of assembled code, got %d", len(firstSegment.Data.Bytes()))
	}
	if firstSegment.Data.Bytes()[0] != 0xA9 {
		t.Errorf("Expected first byte to be 0xA9 (LDA immediate), got 0x%02X", firstSegment.Data.Bytes()[0])
	}
	if firstSegment.Data.Bytes()[1] != 0x10 {
		t.Errorf("Expected second byte to be 0x10, got 0x%02X", firstSegment.Data.Bytes()[1])
	}
	t.Logf("Successfully assembled %d segments with preprocessing", len(segments))
	t.Logf("First segment: %d bytes at address 0x%04X", len(firstSegment.Data.Bytes()), firstSegment.StartAddress)
}

func TestAssembler_PreprocessorCircularIncludeError(t *testing.T) {
	opcodes := createTestOpcodes()
	assembler := New(opcodes)

	// Test files with circular include
	files := map[string]string{
		"main.asm": `LDA #$10
#include "circular.asm"`,
		"circular.asm": `CMP #$20
#include "circular.asm"`,
	}

	resolver := NewMemoryFileResolver(files)

	_, err := assembler.AssembleWithPreprocessor(strings.NewReader(files["main.asm"]), resolver)
	if err == nil {
		t.Fatal("Expected circular include error, got none")
	}

	if !strings.Contains(err.Error(), "circular include") {
		t.Errorf("Expected circular include error, got: %v", err)
	}
}

func TestAssembler_PreprocessorFileNotFoundError(t *testing.T) {
	opcodes := createTestOpcodes()
	assembler := New(opcodes)

	// Test files with missing include
	files := map[string]string{
		"main.asm": `LDA #$10
#include "missing.asm"`,
	}

	resolver := NewMemoryFileResolver(files)

	_, err := assembler.AssembleWithPreprocessor(strings.NewReader(files["main.asm"]), resolver)
	if err == nil {
		t.Fatal("Expected file not found error, got none")
	}

	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("Expected file not found error, got: %v", err)
	}
}

func TestAssembler_PreprocessorNestedIncludes(t *testing.T) {
	// Use real CPU opcodes
	opcodes := createTestOpcodes()
	assembler := New(opcodes)

	// Test nested includes
	files := map[string]string{
		"main.asm": `LDA #$10
.include "level1.asm"`,
		"level1.asm": `CMP #$20
.INCLUDE "level2.asm"`,
		"level2.asm": `INX
INY`,
	}

	resolver := NewMemoryFileResolver(files)

	segments, err := assembler.AssembleWithPreprocessor(strings.NewReader(files["main.asm"]), resolver)
	if err != nil {
		t.Fatalf("Assembly failed: %v", err)
	}
	if len(segments) == 0 {
		t.Fatal("Expected at least one segment, got none")
	}
	firstSegment := segments[0]
	if len(firstSegment.Data.Bytes()) != 6 {
		t.Errorf("Expected 6 bytes of assembled code, got %d", len(firstSegment.Data.Bytes()))
	}
	expected := []byte{0xA9, 0x10, 0xC9, 0x20, 0xE8, 0xC8}
	actual := firstSegment.Data.Bytes()
	for i, expectedByte := range expected {
		if i >= len(actual) {
			t.Errorf("Missing byte at position %d", i)
			continue
		}
		if actual[i] != expectedByte {
			t.Errorf("At position %d: expected 0x%02X, got 0x%02X", i, expectedByte, actual[i])
		}
	}
	t.Logf("Successfully assembled nested includes with %d bytes", len(firstSegment.Data.Bytes()))
}

func TestAssembler_BackwardCompatibility(t *testing.T) {
	// Test that the original Assemble method still works without preprocessing
	opcodes := createTestOpcodes()
	assembler := New(opcodes)

	sourceCode := `LDA #$10`

	segments, err := assembler.Assemble(strings.NewReader(sourceCode))
	if err != nil {
		t.Fatalf("Assembly failed: %v", err)
	}
	if len(segments) == 0 {
		t.Fatal("Expected at least one segment, got none")
	}
	firstSegment := segments[0]
	if len(firstSegment.Data.Bytes()) != 2 {
		t.Errorf("Expected 2 bytes, got %d", len(firstSegment.Data.Bytes()))
	}
	actual := firstSegment.Data.Bytes()
	if actual[0] != 0xA9 || actual[1] != 0x10 {
		t.Errorf("Expected [0xA9, 0x10], got [0x%02X, 0x%02X]", actual[0], actual[1])
	}
	t.Log("Backward compatibility confirmed - original Assemble method works")
}
