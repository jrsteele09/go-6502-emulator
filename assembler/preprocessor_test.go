package assembler

import (
	"strings"
	"testing"
)

func TestPreprocessor_BasicInclude(t *testing.T) {
	files := map[string]string{
		"main.asm": `LDA #$10
#include "sub.asm"
STA $1000`,
		"sub.asm": `CMP #$20
BEQ end`,
	}

	resolver := NewMemoryFileResolver(files)
	preprocessor := NewPreprocessor(resolver)

	result, err := preprocessor.Process(strings.NewReader(files["main.asm"]))
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	output, err := readAll(result)
	if err != nil {
		t.Fatalf("Failed to read result: %v", err)
	}

	expected := `LDA #$10
CMP #$20
BEQ end
STA $1000
`

	if output != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output)
	}
}

func TestPreprocessor_DotInclude(t *testing.T) {
	files := map[string]string{
		"main.asm": `LDA #$10
.include "sub.asm"
STA $1000`,
		"sub.asm": `CMP #$20
BEQ end`,
	}

	resolver := NewMemoryFileResolver(files)
	preprocessor := NewPreprocessor(resolver)

	result, err := preprocessor.Process(strings.NewReader(files["main.asm"]))
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	output, err := readAll(result)
	if err != nil {
		t.Fatalf("Failed to read result: %v", err)
	}

	expected := `LDA #$10
CMP #$20
BEQ end
STA $1000
`

	if output != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output)
	}
}

func TestPreprocessor_UppercaseInclude(t *testing.T) {
	files := map[string]string{
		"main.asm": `LDA #$10
.INCLUDE "sub.asm"
STA $1000`,
		"sub.asm": `CMP #$20
BEQ end`,
	}

	resolver := NewMemoryFileResolver(files)
	preprocessor := NewPreprocessor(resolver)

	result, err := preprocessor.Process(strings.NewReader(files["main.asm"]))
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	output, err := readAll(result)
	if err != nil {
		t.Fatalf("Failed to read result: %v", err)
	}

	expected := `LDA #$10
CMP #$20
BEQ end
STA $1000
`

	if output != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output)
	}
}

func TestPreprocessor_NestedIncludes(t *testing.T) {
	files := map[string]string{
		"main.asm": `LDA #$10
#include "level1.asm"
STA $1000`,
		"level1.asm": `CMP #$20
#include "level2.asm"
BEQ end`,
		"level2.asm": `INX
INY`,
	}

	resolver := NewMemoryFileResolver(files)
	preprocessor := NewPreprocessor(resolver)

	result, err := preprocessor.Process(strings.NewReader(files["main.asm"]))
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	output, err := readAll(result)
	if err != nil {
		t.Fatalf("Failed to read result: %v", err)
	}

	expected := `LDA #$10
CMP #$20
INX
INY
BEQ end
STA $1000
`

	if output != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output)
	}
}

func TestPreprocessor_CircularInclude(t *testing.T) {
	files := map[string]string{
		"main.asm": `LDA #$10
#include "circular.asm"
STA $1000`,
		"circular.asm": `CMP #$20
#include "circular.asm"
BEQ end`,
	}

	resolver := NewMemoryFileResolver(files)
	preprocessor := NewPreprocessor(resolver)

	_, err := preprocessor.Process(strings.NewReader(files["main.asm"]))
	if err == nil {
		t.Fatal("Expected circular include error, got none")
	}

	if !strings.Contains(err.Error(), "circular include") {
		t.Errorf("Expected circular include error, got: %v", err)
	}
}

func TestPreprocessor_MaxDepthExceeded(t *testing.T) {
	files := map[string]string{
		"main.asm":   `#include "level1.asm"`,
		"level1.asm": `#include "level2.asm"`,
		"level2.asm": `#include "level3.asm"`,
		"level3.asm": `#include "level4.asm"`,
		"level4.asm": `#include "level5.asm"`,
		"level5.asm": `LDA #$10`,
	}

	resolver := NewMemoryFileResolver(files)
	preprocessor := NewPreprocessor(resolver)
	preprocessor.SetMaxDepth(3) // Set low depth limit

	_, err := preprocessor.Process(strings.NewReader(files["main.asm"]))
	if err == nil {
		t.Fatal("Expected max depth error, got none")
	}

	if !strings.Contains(err.Error(), "maximum include depth") {
		t.Errorf("Expected max depth error, got: %v", err)
	}
}

func TestPreprocessor_FileNotFound(t *testing.T) {
	files := map[string]string{
		"main.asm": `LDA #$10
#include "missing.asm"
STA $1000`,
	}

	resolver := NewMemoryFileResolver(files)
	preprocessor := NewPreprocessor(resolver)

	_, err := preprocessor.Process(strings.NewReader(files["main.asm"]))
	if err == nil {
		t.Fatal("Expected file not found error, got none")
	}

	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("Expected file not found error, got: %v", err)
	}
}

func TestPreprocessor_SingleQuotes(t *testing.T) {
	files := map[string]string{
		"main.asm": `LDA #$10
#include 'sub.asm'
STA $1000`,
		"sub.asm": `CMP #$20`,
	}

	resolver := NewMemoryFileResolver(files)
	preprocessor := NewPreprocessor(resolver)

	result, err := preprocessor.Process(strings.NewReader(files["main.asm"]))
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	output, err := readAll(result)
	if err != nil {
		t.Fatalf("Failed to read result: %v", err)
	}

	expected := `LDA #$10
CMP #$20
STA $1000
`

	if output != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output)
	}
}

func TestPreprocessor_UnquotedPath(t *testing.T) {
	files := map[string]string{
		"main.asm": `LDA #$10
#include sub.asm
STA $1000`,
		"sub.asm": `CMP #$20`,
	}

	resolver := NewMemoryFileResolver(files)
	preprocessor := NewPreprocessor(resolver)

	result, err := preprocessor.Process(strings.NewReader(files["main.asm"]))
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	output, err := readAll(result)
	if err != nil {
		t.Fatalf("Failed to read result: %v", err)
	}

	expected := `LDA #$10
CMP #$20
STA $1000
`

	if output != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output)
	}
}

func TestPreprocessor_MixedIncludeStyles(t *testing.T) {
	files := map[string]string{
		"main.asm": `LDA #$10
#include "hash.asm"
.include "dot.asm"
.INCLUDE "upper.asm"
STA $1000`,
		"hash.asm": `; Hash include
CMP #$20`,
		"dot.asm": `; Dot include
INX`,
		"upper.asm": `; Upper include
INY`,
	}

	resolver := NewMemoryFileResolver(files)
	preprocessor := NewPreprocessor(resolver)

	result, err := preprocessor.Process(strings.NewReader(files["main.asm"]))
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	output, err := readAll(result)
	if err != nil {
		t.Fatalf("Failed to read result: %v", err)
	}

	expected := `LDA #$10
; Hash include
CMP #$20
; Dot include
INX
; Upper include
INY
STA $1000
`

	if output != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output)
	}
}

func TestPreprocessor_EmptyFile(t *testing.T) {
	files := map[string]string{
		"main.asm": `LDA #$10
#include "empty.asm"
STA $1000`,
		"empty.asm": ``,
	}

	resolver := NewMemoryFileResolver(files)
	preprocessor := NewPreprocessor(resolver)

	result, err := preprocessor.Process(strings.NewReader(files["main.asm"]))
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	output, err := readAll(result)
	if err != nil {
		t.Fatalf("Failed to read result: %v", err)
	}

	expected := `LDA #$10
STA $1000
`

	if output != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output)
	}
}

func TestPreprocessor_WhitespaceHandling(t *testing.T) {
	files := map[string]string{
		"main.asm": `LDA #$10
    #include    "sub.asm"    
STA $1000`,
		"sub.asm": `CMP #$20`,
	}

	resolver := NewMemoryFileResolver(files)
	preprocessor := NewPreprocessor(resolver)

	result, err := preprocessor.Process(strings.NewReader(files["main.asm"]))
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	output, err := readAll(result)
	if err != nil {
		t.Fatalf("Failed to read result: %v", err)
	}

	expected := `LDA #$10
CMP #$20
STA $1000
`

	if output != expected {
		t.Errorf("Expected:\n%s\nGot:\n%s", expected, output)
	}
}

// Helper function to read all content from a reader
func readAll(r interface{}) (string, error) {
	if sr, ok := r.(*strings.Reader); ok {
		content := make([]byte, sr.Len())
		_, err := sr.Read(content)
		return string(content), err
	}
	return "", nil
}
