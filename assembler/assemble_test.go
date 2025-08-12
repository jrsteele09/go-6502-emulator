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

	asmCode := `
	.ORG $1000
	
	.EQU SCREEN_BASE = $0400
	
	start:
		LDA #$01
		STA SCREEN_BASE
		JMP *
		
	data:
		.BYTE $48, $45, $4C
		.WORD start
	`

	dataSegments, err := assembler.Assemble(strings.NewReader(asmCode))
	require.NoError(t, err)
	require.NotNil(t, dataSegments)
	require.GreaterOrEqual(t, len(dataSegments), 1)

	// Debug: print all segments
	for i, seg := range dataSegments {
		t.Logf("Segment %d: StartAddress=0x%04x, DataLen=%d", i, seg.StartAddress, len(seg.Data))
	}

	// Find the segment that starts at $1000
	var targetSegmentIndex = -1
	for i, seg := range dataSegments {
		if seg.StartAddress == 0x1000 {
			targetSegmentIndex = i
			break
		}
	}
	require.NotEqual(t, -1, targetSegmentIndex, "Should have a segment starting at $1000")
	require.Greater(t, len(dataSegments[targetSegmentIndex].Data), 0)
}

func TestAssemblerBasic(t *testing.T) {

	m := memory.NewMemory[uint16](64 * 1024)
	c := cpu.NewCPU(m)

	assembler := assembler.New(cpu.OpCodes(c))

	asmCode := `BNE *-125`

	dataSegments, err := assembler.Assemble(strings.NewReader(asmCode))
	require.NoError(t, err)
	require.NotNil(t, dataSegments)
}

func TestAssemblerBneWithLabel(t *testing.T) {

	m := memory.NewMemory[uint16](64 * 1024)
	c := cpu.NewCPU(m)

	assembler := assembler.New(cpu.OpCodes(c))

	asmCode := `
	LOOP:
    DEX         ; Decrement X register
    BNE LOOP    ; Branch if not equal to zero
	`

	dataSegments, err := assembler.Assemble(strings.NewReader(asmCode))
	require.NoError(t, err)
	require.NotNil(t, dataSegments)
}

func TestAssemblerBneWithLabelOutOfRange(t *testing.T) {

	m := memory.NewMemory[uint16](64 * 1024)
	c := cpu.NewCPU(m)

	assembler := assembler.New(cpu.OpCodes(c))

	asmCode := `
	LOOP:
    DEX   
	*=$c000
    BNE LOOP    ; Branch if not equal to zero
	`

	dataSegments, err := assembler.Assemble(strings.NewReader(asmCode))
	require.Error(t, err)
	require.Nil(t, dataSegments)
}

func TestAssemblerComprehensive(t *testing.T) {
	m := memory.NewMemory[uint16](64 * 1024)
	c := cpu.NewCPU(m)

	assembler := assembler.New(cpu.OpCodes(c))

	// Test program that exercises all enhanced assembler features
	asmCode := `; Test program for enhanced 6502 assembler
.ORG $1000

; Define some constants
.EQU SCREEN_BASE = $0400
.EQU COLOR_WHITE = $01

; Define some data
start:
    .BYTE $01, $02, $03
    .TEXT "HELLO WORLD"

; Some code
main:
    LDA #COLOR_WHITE
    STA SCREEN_BASE
    LDX #$00
    LDA #$48
    STA SCREEN_BASE
    JMP *       ; Infinite loop (jump to current address)

data:
    .BYTE $48, $45, $4C  ; "HEL" in ASCII
    
end:
    ; Reserve 16 bytes of space
    .DS 16
    .WORD start, end     ; Put WORD directive after labels are defined
`

	dataSegments, err := assembler.Assemble(strings.NewReader(asmCode))
	require.NoError(t, err)
	require.NotNil(t, dataSegments)
	require.GreaterOrEqual(t, len(dataSegments), 1)

	// Debug: print all segments
	for i, seg := range dataSegments {
		t.Logf("Segment %d: StartAddress=0x%04x, DataLen=%d", i, seg.StartAddress, len(seg.Data))
	}

	// Find the segment that starts at $1000
	var targetSegmentIndex = -1
	for i, seg := range dataSegments {
		if seg.StartAddress == 0x1000 {
			targetSegmentIndex = i
			break
		}
	}
	require.NotEqual(t, -1, targetSegmentIndex, "Should have a segment starting at $1000")

	segment := dataSegments[targetSegmentIndex]
	require.Greater(t, len(segment.Data), 0)

	// Test that the assembled code contains expected data
	data := segment.Data

	// Verify .BYTE $01, $02, $03 at the start
	require.GreaterOrEqual(t, len(data), 3, "Should have at least 3 bytes for initial .BYTE directive")
	require.Equal(t, uint8(0x01), data[0], "First byte should be 0x01")
	require.Equal(t, uint8(0x02), data[1], "Second byte should be 0x02")
	require.Equal(t, uint8(0x03), data[2], "Third byte should be 0x03")

	// Verify .TEXT "HELLO WORLD" is present
	// The text should start after the initial bytes
	expectedText := "HELLO WORLD"
	textStartIndex := 3 // After 3 bytes
	require.GreaterOrEqual(t, len(data), textStartIndex+len(expectedText), "Should have enough data for text")

	actualText := string(data[textStartIndex : textStartIndex+len(expectedText)])
	require.Equal(t, expectedText, actualText, "Text should match 'HELLO WORLD'")

	t.Logf("Successfully assembled program with %d bytes of data", len(data))
	t.Logf("Program starts at address 0x%04x", segment.StartAddress)

	// Verify that instructions and additional data are present
	// The program should have: 3 bytes + 11 chars + instructions + 3 data bytes + 16 reserved bytes + 4 word bytes
	expectedMinSize := 3 + len(expectedText) + 10 // Minimum expected with some instructions
	require.Greater(t, len(data), expectedMinSize, "Should have instructions and additional data")
}

func TestAssemblerProgramCounterSet(t *testing.T) {
	m := memory.NewMemory[uint16](64 * 1024)
	c := cpu.NewCPU(m)

	assembler := assembler.New(cpu.OpCodes(c))

	// Test that *= works the same as .ORG
	asmCode := `*= $2000

start:
    LDA #$42
    STA $0400
    JMP *
    
data:
    .BYTE $01, $02, $03
`

	dataSegments, err := assembler.Assemble(strings.NewReader(asmCode))
	require.NoError(t, err)
	require.NotNil(t, dataSegments)
	require.GreaterOrEqual(t, len(dataSegments), 1)

	// Debug: print all segments
	for i, seg := range dataSegments {
		t.Logf("Segment %d: StartAddress=0x%04x, DataLen=%d", i, seg.StartAddress, len(seg.Data))
	}

	// Find the segment that starts at $2000
	var targetSegmentIndex = -1
	for i, seg := range dataSegments {
		if seg.StartAddress == 0x2000 {
			targetSegmentIndex = i
			break
		}
	}
	require.NotEqual(t, -1, targetSegmentIndex, "Should have a segment starting at $2000")

	segment := dataSegments[targetSegmentIndex]
	require.Greater(t, len(segment.Data), 0)

	// Test that the assembled code contains expected data
	data := segment.Data

	// Should have instructions and data
	require.Greater(t, len(data), 5, "Should have instructions and data")

	t.Logf("Successfully assembled program with *= directive: %d bytes at address 0x%04x", len(data), segment.StartAddress)
}

func TestAssemblerProgramCounterSetVsOrg(t *testing.T) {
	m := memory.NewMemory[uint16](64 * 1024)
	c := cpu.NewCPU(m)

	assembler := assembler.New(cpu.OpCodes(c))

	// Test program using .ORG
	asmCodeOrg := `.ORG $3000

start:
    LDA #$42
    STA $0400
    JMP start
    
    .BYTE $01, $02, $03
`

	// Test program using *=
	asmCodeStar := `*= $3000

start:
    LDA #$42
    STA $0400
    JMP start
    
    .BYTE $01, $02, $03
`

	// Assemble both versions
	dataSegmentsOrg, err := assembler.Assemble(strings.NewReader(asmCodeOrg))
	require.NoError(t, err)
	require.NotNil(t, dataSegmentsOrg)
	require.GreaterOrEqual(t, len(dataSegmentsOrg), 1)

	dataSegmentsStar, err := assembler.Assemble(strings.NewReader(asmCodeStar))
	require.NoError(t, err)
	require.NotNil(t, dataSegmentsStar)
	require.GreaterOrEqual(t, len(dataSegmentsStar), 1)

	// Both should produce identical results
	require.Equal(t, len(dataSegmentsOrg), len(dataSegmentsStar), "Should have same number of segments")

	// Find the segments that start at $3000
	var orgSegmentIndex, starSegmentIndex = -1, -1
	for i, seg := range dataSegmentsOrg {
		if seg.StartAddress == 0x3000 {
			orgSegmentIndex = i
			break
		}
	}
	for i, seg := range dataSegmentsStar {
		if seg.StartAddress == 0x3000 {
			starSegmentIndex = i
			break
		}
	}

	require.NotEqual(t, -1, orgSegmentIndex, "Should have .ORG segment at $3000")
	require.NotEqual(t, -1, starSegmentIndex, "Should have *= segment at $3000")

	orgSegment := dataSegmentsOrg[orgSegmentIndex]
	starSegment := dataSegmentsStar[starSegmentIndex]

	// Compare the assembled data
	require.Equal(t, orgSegment.StartAddress, starSegment.StartAddress, "Start addresses should be identical")
	require.Equal(t, len(orgSegment.Data), len(starSegment.Data), "Data lengths should be identical")
	require.Equal(t, orgSegment.Data, starSegment.Data, "Assembled data should be identical")

	t.Logf(".ORG assembled %d bytes at 0x%04x", len(orgSegment.Data), orgSegment.StartAddress)
	t.Logf("*= assembled %d bytes at 0x%04x", len(starSegment.Data), starSegment.StartAddress)
	t.Log("Both .ORG and *= produce identical results")
}

func TestAssemblerStringDirectives(t *testing.T) {
	m := memory.NewMemory[uint16](64 * 1024)
	c := cpu.NewCPU(m)

	assembler := assembler.New(cpu.OpCodes(c))

	// Test all string directive variants
	asmCode := `*= $4000

text_data:
    .TEXT "Hello"
    
string_data:
    .STRING "World"
    
str_data:
    .STR "Test"
    
asc_data:
    .ASC "ASCII"
    
asciiz_data:
    .ASCIIZ "NullTerm"
`

	dataSegments, err := assembler.Assemble(strings.NewReader(asmCode))
	require.NoError(t, err)
	require.NotNil(t, dataSegments)
	require.GreaterOrEqual(t, len(dataSegments), 1)

	// Find the segment that starts at $4000
	var targetSegmentIndex = -1
	for i, seg := range dataSegments {
		if seg.StartAddress == 0x4000 {
			targetSegmentIndex = i
			break
		}
	}
	require.NotEqual(t, -1, targetSegmentIndex, "Should have a segment starting at $4000")

	segment := dataSegments[targetSegmentIndex]
	data := segment.Data

	// Expected data: "Hello" + "World" + "Test" + "ASCII" + "NullTerm" + null terminator
	expectedStrings := []string{"Hello", "World", "Test", "ASCII", "NullTerm"}
	expectedSize := len("Hello") + len("World") + len("Test") + len("ASCII") + len("NullTerm") + 1 // +1 for null terminator

	require.Equal(t, expectedSize, len(data), "Should contain all string data")

	// Verify the strings are concatenated correctly
	offset := 0
	for i, expectedStr := range expectedStrings {
		actualStr := string(data[offset : offset+len(expectedStr)])
		require.Equal(t, expectedStr, actualStr, "String %d should match", i)
		offset += len(expectedStr)
	}

	// Verify null terminator for .ASCIIZ
	require.Equal(t, byte(0), data[len(data)-1], "Should end with null terminator from .ASCIIZ")

	t.Logf("Successfully assembled string directives: %d bytes at address 0x%04x", len(data), segment.StartAddress)
	t.Logf("String data: %v", data)
}

func createTestAssembler() *assembler.Assembler {
	m := memory.NewMemory[uint16](64 * 1024)
	c := cpu.NewCPU(m)
	return assembler.New(cpu.OpCodes(c))
}

func TestDuplicateLabels(t *testing.T) {
	assembler := createTestAssembler()

	program := `
		.ORG $1000
	start:
		NOP
	start:  ; Duplicate label - should error
		LDA #$FF
	`

	reader := strings.NewReader(program)
	_, err := assembler.Assemble(reader)

	if err == nil {
		t.Fatal("Expected error for duplicate label, but got none")
	}

	if !strings.Contains(err.Error(), "duplicate label 'start'") {
		t.Errorf("Expected duplicate label error, got: %v", err)
	}
}

func TestDuplicateVariables(t *testing.T) {
	assembler := createTestAssembler()

	program := `
		.ORG $1000
		.EQU myVar = $FF
		.EQU myVar = $AA  ; Duplicate variable - should error
		NOP
	`

	reader := strings.NewReader(program)
	_, err := assembler.Assemble(reader)

	if err == nil {
		t.Fatal("Expected error for duplicate variable, but got none")
	}

	if !strings.Contains(err.Error(), "duplicate variable 'myVar'") {
		t.Errorf("Expected duplicate variable error, got: %v", err)
	}
}

func TestLabelVariableConflict(t *testing.T) {
	assembler := createTestAssembler()

	program := `
		.ORG $1000
		.EQU mySymbol = $FF
	mySymbol:  ; Conflicts with variable - should error
		NOP
	`

	reader := strings.NewReader(program)
	_, err := assembler.Assemble(reader)

	if err == nil {
		t.Fatal("Expected error for label/variable conflict, but got none")
	}

	if !strings.Contains(err.Error(), "conflicts with existing variable") {
		t.Errorf("Expected label/variable conflict error, got: %v", err)
	}
}

func TestVariableLabelConflict(t *testing.T) {
	assembler := createTestAssembler()

	program := `
		.ORG $1000
	mySymbol:
		NOP
		.EQU mySymbol = $FF  ; Conflicts with label - should error
	`

	reader := strings.NewReader(program)
	_, err := assembler.Assemble(reader)

	if err == nil {
		t.Fatal("Expected error for variable/label conflict, but got none")
	}

	if !strings.Contains(err.Error(), "conflicts with existing label") {
		t.Errorf("Expected variable/label conflict error, got: %v", err)
	}
}

func TestValidUniqueSymbols(t *testing.T) {
	assembler := createTestAssembler()

	program := `
		.ORG $1000
		.EQU var1 = $FF
		.EQU var2 = $AA
	label1:
		NOP
	label2:
		LDA var1
		LDX var2
		JMP label1
	`

	reader := strings.NewReader(program)
	segments, err := assembler.Assemble(reader)

	if err != nil {
		t.Fatalf("Expected successful assembly with unique symbols, got error: %v", err)
	}

	if len(segments) != 1 {
		t.Errorf("Expected 1 segment, got %d", len(segments))
	}

	if segments[0].StartAddress != 0x1000 {
		t.Errorf("Expected start address 0x1000, got 0x%X", segments[0].StartAddress)
	}

	// Verify that we assembled some bytes
	if len(segments[0].Data) == 0 {
		t.Error("Expected assembled data, got empty segment")
	}
}

func TestStateResetBetweenAssemblies(t *testing.T) {
	assembler := createTestAssembler()

	// First assembly with a label
	program1 := `
		.ORG $1000
	testLabel:
		NOP
	`

	reader1 := strings.NewReader(program1)
	_, err := assembler.Assemble(reader1)
	if err != nil {
		t.Fatalf("First assembly failed: %v", err)
	}

	// Second assembly should not conflict with first (state should be reset)
	program2 := `
		.ORG $2000
	testLabel:  ; Same label name - should be OK due to state reset
		LDA #$FF
	`

	reader2 := strings.NewReader(program2)
	segments, err := assembler.Assemble(reader2)
	if err != nil {
		t.Fatalf("Second assembly failed - state not properly reset: %v", err)
	}

	if len(segments) != 1 {
		t.Errorf("Expected 1 segment, got %d", len(segments))
	}

	if segments[0].StartAddress != 0x2000 {
		t.Errorf("Expected start address 0x2000, got 0x%X", segments[0].StartAddress)
	}
}
