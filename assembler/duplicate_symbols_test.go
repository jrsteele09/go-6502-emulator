package assembler

import (
	"strings"
	"testing"

	"github.com/jrsteele09/go-6502-emulator/cpu"
	"github.com/jrsteele09/go-6502-emulator/memory"
)

func createTestAssembler() *Assembler {
	m := memory.NewMemory[uint16](64 * 1024)
	c := cpu.NewCPU(m)
	return New(cpu.OpCodes(c))
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
