package cpu_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/jrsteele09/go-6502-emulator/assembler"
	"github.com/jrsteele09/go-6502-emulator/cpu"
	"github.com/jrsteele09/go-6502-emulator/debugger"
	"github.com/jrsteele09/go-6502-emulator/memory"
	"github.com/stretchr/testify/require"
)

const (
	functionalTestBinURL = "https://raw.githubusercontent.com/Klaus2m5/6502_65C02_functional_tests/master/bin_files/6502_functional_test.bin"
	functionalTestLstURL = "https://raw.githubusercontent.com/Klaus2m5/6502_65C02_functional_tests/master/bin_files/6502_functional_test.lst"
	decimalTestSourceURL = "https://raw.githubusercontent.com/Klaus2m5/6502_65C02_functional_tests/master/6502_decimal_test.a65"

	klausCacheDir          = "testdata/klaus"
	functionalTestBinName  = "6502_functional_test.bin"
	functionalTestLstName  = "6502_functional_test.lst"
	decimalTestSourceName  = "6502_decimal_test.a65"
	functionalStart        = uint16(0x0400)
	functionalStop         = uint16(0x3469)
	functionalInstructions = uint64(30_646_177)
	functionalCycles       = uint64(96_241_367)
	decimalStart           = uint16(0x0200)
	decimalStop            = uint16(0x024B)
	decimalErrorAddress    = uint16(0x000C)
	decimalStopOpcode      = byte(0xDB)
	decimalInstructions    = uint64(307_166)
	decimalCycles          = uint64(984_995)
	maxKlausInstructions   = uint64(900_000_000)
)

type klausResult struct {
	instructions uint64
	cycles       uint64
	stopAddress  uint16
}

func TestKlausFunctional(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Klaus integration test in short mode")
	}
	binPath := klausArtifact(t, functionalTestBinName, functionalTestBinURL)
	lstPath := klausArtifact(t, functionalTestLstName, functionalTestLstURL)
	require.NoError(t, verifyKlausPassTrap(lstPath, functionalStop))

	result, err := runKlausFunctional(binPath, functionalStart, functionalStop, maxKlausInstructions)
	require.NoError(t, err)
	require.Equal(t, functionalStop, result.stopAddress)
	require.Equal(t, functionalInstructions, result.instructions)
	require.Equal(t, functionalCycles, result.cycles)
}

func TestKlausDecimal(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping Klaus integration test in short mode")
	}
	sourcePath := klausArtifact(t, decimalTestSourceName, decimalTestSourceURL)

	result, err := runKlausDecimal(sourcePath, maxKlausInstructions)
	require.NoError(t, err)
	require.Equal(t, decimalStop, result.stopAddress)
	require.Equal(t, decimalInstructions, result.instructions)
	require.Equal(t, decimalCycles, result.cycles)
}

func klausArtifact(t *testing.T, name, url string) string {
	t.Helper()
	require.NoError(t, os.MkdirAll(klausCacheDir, 0o755))
	path := filepath.Join(klausCacheDir, name)
	if info, err := os.Stat(path); err == nil && info.Size() > 0 {
		return path
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	t.Cleanup(cancel)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoErrorf(t, err, "download Klaus artifact %s", url)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode, "download Klaus artifact %s", url)

	tmpPath := path + ".tmp"
	out, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	require.NoError(t, err)
	_, copyErr := io.Copy(out, resp.Body)
	closeErr := out.Close()
	require.NoError(t, copyErr)
	require.NoError(t, closeErr)
	require.NoError(t, os.Rename(tmpPath, path))
	return path
}

func runKlausFunctional(binPath string, startAddress, successAddress uint16, maxInstructions uint64) (klausResult, error) {
	image, err := os.ReadFile(binPath)
	if err != nil {
		return klausResult{}, err
	}
	if len(image) > 64*1024 {
		return klausResult{}, fmt.Errorf("test binary is %d bytes; expected at most 65536", len(image))
	}

	mem := memory.NewMemory[uint16](64 * 1024)
	mem.Write(0, image...)
	processor := cpu.NewCPU(mem, false)
	processor.Registers().PC = startAddress
	disasm := debugger.NewDisassembler(mem, processor.OpCodes())
	lastPC := processor.Registers().PC
	result := klausResult{}

	for result.instructions < maxInstructions {
		currentPC := processor.Registers().PC
		cycles, err := executeKlausInstruction(processor)
		if err != nil {
			return klausResult{}, fmt.Errorf("execution error at %s after %d instructions / %d cycles: %w\n%s", formatKlausAddress(currentPC), result.instructions, result.cycles, err, formatKlausCPUState(processor, disasm))
		}
		result.instructions++
		result.cycles += cycles

		pc := processor.Registers().PC
		if pc == lastPC {
			if pc == successAddress {
				result.stopAddress = pc
				return result, nil
			}
			return klausResult{}, fmt.Errorf("test trapped at %s after %d instructions / %d cycles\n%s", formatKlausAddress(pc), result.instructions, result.cycles, formatKlausCPUState(processor, disasm))
		}
		lastPC = pc
	}

	return klausResult{}, fmt.Errorf("timed out after %d completed instructions / %d cycles\n%s", maxInstructions, result.cycles, formatKlausCPUState(processor, disasm))
}

func runKlausDecimal(sourcePath string, maxInstructions uint64) (klausResult, error) {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return klausResult{}, fmt.Errorf("open source test: %w", err)
	}
	defer sourceFile.Close()

	mem := memory.NewMemory[uint16](64 * 1024)
	processor := cpu.NewCPU(mem, false)
	segments, err := assembler.New(processor.OpCodes()).Assemble(sourceFile, sourcePath)
	if err != nil {
		return klausResult{}, fmt.Errorf("assemble source test: %w", err)
	}
	for _, segment := range segments {
		mem.Write(segment.StartAddress, segment.Data.Bytes()...)
	}

	processor.Registers().PC = decimalStart
	disasm := debugger.NewDisassembler(mem, processor.OpCodes())
	result := klausResult{}
	for result.instructions < maxInstructions {
		currentPC := processor.Registers().PC
		if mem.Read(currentPC) == decimalStopOpcode {
			if mem.Read(decimalErrorAddress) != 0 {
				return klausResult{}, fmt.Errorf("decimal test reported ERROR=%d at %s after %d instructions / %d cycles\n%s", mem.Read(decimalErrorAddress), formatKlausAddress(currentPC), result.instructions, result.cycles, formatKlausCPUState(processor, disasm))
			}
			result.stopAddress = currentPC
			return result, nil
		}

		cycles, err := executeKlausInstruction(processor)
		if err != nil {
			return klausResult{}, fmt.Errorf("execution error at %s after %d instructions / %d cycles: %w\n%s", formatKlausAddress(currentPC), result.instructions, result.cycles, err, formatKlausCPUState(processor, disasm))
		}
		result.instructions++
		result.cycles += cycles
	}

	return klausResult{}, fmt.Errorf("decimal test timed out after %d completed instructions / %d cycles\n%s", maxInstructions, result.cycles, formatKlausCPUState(processor, disasm))
}

func executeKlausInstruction(processor *cpu.CPU) (uint64, error) {
	var cycles uint64
	completed := cpu.Completed(false)
	for !completed {
		done, err := processor.Execute()
		cycles++
		if err != nil {
			return cycles, err
		}
		completed = done
	}
	return cycles, nil
}

func verifyKlausPassTrap(lstPath string, passAddress uint16) error {
	listing, err := os.ReadFile(lstPath)
	if err != nil {
		return err
	}
	passTrap := fmt.Sprintf("%04x : 4c%02x%02x", passAddress, byte(passAddress), byte(passAddress>>8))
	if !strings.Contains(strings.ToLower(string(listing)), passTrap) {
		return fmt.Errorf("expected %s to contain success trap %q", lstPath, passTrap)
	}
	return nil
}

func formatKlausAddress(address uint16) string {
	return fmt.Sprintf("$%04X", address)
}

func formatKlausCPUState(processor *cpu.CPU, disasm *debugger.Disassembler) string {
	registers := processor.Registers()
	instruction, _ := disasm.Disassemble(registers.PC)
	return fmt.Sprintf(
		"PC=%s A=$%02X X=$%02X Y=$%02X S=$%02X P=$%02X\nNext: %s",
		formatKlausAddress(registers.PC),
		registers.A,
		registers.X,
		registers.Y,
		registers.S,
		registers.Status,
		instruction,
	)
}
