package assembler_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/jrsteele09/go-6502-emulator/assembler"
	"github.com/jrsteele09/go-6502-emulator/cpu"
	"github.com/jrsteele09/go-6502-emulator/debugger"
	"github.com/jrsteele09/go-6502-emulator/memory"
	"github.com/jrsteele09/go-6502-emulator/utils"
	"github.com/stretchr/testify/require"
)

const (
	createExpectedResults = false
	expectedResultsFolder = "./testasm/expected_results"
)

func createHardware() (*memory.Memory[uint16], *cpu.CPU) {
	mem := memory.NewMemory[uint16](64 * 1024) // 64KB memory
	cpu := cpu.NewCPU(mem, false)
	return mem, cpu
}

func TestAssemble_TestIncludeFile(t *testing.T) {
	// SETUP
	_, cpu := createHardware()
	asm := assembler.New(cpu.OpCodes())
	resolver := utils.NewOSFileResolver("./testasm/TestIncludeFile")

	// ASSEMBLE
	segments, err := asm.AssembleFile("testinclude.asm", resolver)

	// ASSERT ASSEMBLED RESULTS
	require.NoError(t, err, "AssembleFile failed")
	require.Len(t, segments, 1, "Expected exactly one segment")
	require.Equal(t, uint16(0xC000), segments[0].StartAddress, "Expected start address $C000")

	// ASSERT DISSASSEMBLY
	disassembleAndCompare(t, segments)
}

func TestAssemble_Snake(t *testing.T) {
	// SETUP
	_, cpu := createHardware()
	asm := assembler.New(cpu.OpCodes())
	resolver := utils.NewOSFileResolver("./testasm/TestSnakeAssembly")

	// ASSEMBLE
	segments, err := asm.AssembleFile("snake.asm", resolver)

	// ASSERT ASSEMBLED RESULTS
	require.NoError(t, err, "AssembleFile failed")
	require.Len(t, segments, 1, "Expected exactly one segment")
	require.Equal(t, uint16(0x4000), segments[0].StartAddress, "Expected start address $C000")

	// ASSERT DISSASSEMBLY
	disassembleAndCompare(t, segments)
}

func TestAssemble_TestBorderFlashLoop(t *testing.T) {
	// SETUP
	_, cpu := createHardware()
	asm := assembler.New(cpu.OpCodes())
	resolver := utils.NewOSFileResolver("./testasm/TestBorderFlashLoop")

	// ASSEMBLE
	segments, err := asm.AssembleFile("screen.asm", resolver)

	// ASSERT ASSEMBLED RESULTS
	require.NoError(t, err, "AssembleFile failed")
	require.Len(t, segments, 1, "Expected exactly one segment")
	require.Equal(t, uint16(0xC000), segments[0].StartAddress, "Expected start address $C000")

	// ASSERT DISASSEMBLY
	disassembleAndCompare(t, segments)
}

func TestAssemble_MultiSegmentAssembly(t *testing.T) {
	// SETUP
	_, cpu := createHardware()
	asm := assembler.New(cpu.OpCodes())
	resolver := utils.NewOSFileResolver("./testasm/TestMultiSegmentAssembly")

	// ASSEMBLE
	segments, err := asm.AssembleFile("main.asm", resolver)

	// ASSERT ASSEMBLED RESULTS
	require.NoError(t, err, "AssembleFile failed")
	require.Len(t, segments, 2, "Expected exactly one segment")

	// ASSERT DISASSEMBLY
	disassembleAndCompare(t, segments)
}

func TestAssemble_OpenTheBorder(t *testing.T) {
	// SETUP
	_, cpu := createHardware()
	asm := assembler.New(cpu.OpCodes())
	resolver := utils.NewOSFileResolver("./testasm/TestOpenTheBorderAssembly")

	// ASSEMBLE
	segments, err := asm.AssembleFile("temp.asm", resolver)

	// ASSERT ASSEMBLED RESULTS
	require.NoError(t, err, "AssembleFile failed")
	require.Len(t, segments, 1, "Expected exactly one segment")

	// ASSERT DISASSEMBLY
	disassembleAndCompare(t, segments)
}

func disassembleAndCompare(t *testing.T, segments []assembler.AssembledData) {
	mem, cpu := createHardware()
	writeSegmentsToMemory(mem, segments)
	dbg := debugger.NewDisassembler(mem, cpu.OpCodes())

	for i, segment := range segments {
		testFile := fmt.Sprintf("%s/%s.%d.txt", expectedResultsFolder, t.Name(), i)
		start := segment.StartAddress
		end := segment.StartAddress + uint16(len(segment.Data.Bytes())) - 1
		output := disassembleFromAddress(dbg, start, end)
		if createExpectedResults {
			os.WriteFile(testFile, []byte(output), 0644)
		}
		expectedResults, err := os.ReadFile(testFile)
		require.NoError(t, err, "Failed to read expected results file")
		compareDisassembly(t, string(expectedResults), output)
	}
}

func writeSegmentsToMemory(mem *memory.Memory[uint16], segments []assembler.AssembledData) {
	for _, segment := range segments {
		for i, b := range segment.Data.Bytes() {
			mem.Write(segment.StartAddress+uint16(i), b)
		}
	}
}

func disassembleFromAddress(dbg *debugger.Disassembler, startAddr uint16, endAddr uint16) string {
	var strBuffer strings.Builder
	address := startAddr
	for {
		output, b := dbg.Disassemble(address)
		strBuffer.WriteString(fmt.Sprintf("%s\n", output))
		address += uint16(b)
		if address > endAddr {
			break
		}
	}
	return strBuffer.String()
}

func compareDisassembly(t *testing.T, expected, actual string) {
	expectedArray := strings.Split(strings.TrimSpace(expected), "\n")
	actualArray := strings.Split(strings.TrimSpace(actual), "\n")
	for i := range expectedArray {
		require.Equal(t, expectedArray[i], actualArray[i], "Mismatch at line %d: expected %q, got %q", i+1, expectedArray[i], actualArray[i])
	}
}
