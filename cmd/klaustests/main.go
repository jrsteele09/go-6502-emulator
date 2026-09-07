package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jrsteele09/go-6502-emulator/assembler"
	"github.com/jrsteele09/go-6502-emulator/cpu"
	"github.com/jrsteele09/go-6502-emulator/debugger"
	"github.com/jrsteele09/go-6502-emulator/memory"
)

const (
	functionalTestBinURL = "https://raw.githubusercontent.com/Klaus2m5/6502_65C02_functional_tests/master/bin_files/6502_functional_test.bin"
	functionalTestLstURL = "https://raw.githubusercontent.com/Klaus2m5/6502_65C02_functional_tests/master/bin_files/6502_functional_test.lst"
	decimalTestSourceURL = "https://raw.githubusercontent.com/Klaus2m5/6502_65C02_functional_tests/master/6502_decimal_test.a65"

	defaultCacheDir        = "testdata/klaus"
	functionalTestBinName  = "6502_functional_test.bin"
	functionalTestLstName  = "6502_functional_test.lst"
	decimalTestSourceName  = "6502_decimal_test.a65"
	defaultStartAddress    = uint16(0x0400)
	defaultSuccessAddress  = uint16(0x3469)
	defaultMaxInstructions = uint64(900_000_000)
	decimalTestStart       = uint16(0x0200)
	decimalTestError       = uint16(0x000C)
	decimalTestStopOpcode  = byte(0xDB)
	functionalInstructions = uint64(30_646_177)
	functionalCycles       = uint64(96_241_367)
)

type runResult struct {
	Instructions uint64
	Cycles       uint64
	StopAddress  uint16
}

func main() {
	cacheDir := flag.String("cache-dir", defaultCacheDir, "directory used to cache downloaded Klaus test artifacts")
	source := flag.String("source", "", "assemble and run a source test instead of the downloaded functional test")
	forceDownload := flag.Bool("force", false, "download test artifacts even when cached files exist")
	startAddress := flag.String("start", formatAddress(defaultStartAddress), "program start address")
	successAddress := flag.String("pass", formatAddress(defaultSuccessAddress), "success trap address")
	maxInstructions := flag.Uint64("max-instructions", defaultMaxInstructions, "maximum completed instructions before failing")
	flag.Parse()

	if *source != "" {
		result, err := runDecimalSourceTest(*source, *maxInstructions)
		if err != nil {
			fatalf("%v", err)
		}
		fmt.Printf("PASS: assembled decimal source test completed at %s after %d instructions / %d cycles\n", formatAddress(result.StopAddress), result.Instructions, result.Cycles)
		return
	}

	start, err := parseAddress(*startAddress)
	if err != nil {
		fatalf("invalid -start address: %v", err)
	}
	pass, err := parseAddress(*successAddress)
	if err != nil {
		fatalf("invalid -pass address: %v", err)
	}

	binPath, lstPath, decimalSourcePath, err := ensureKlausTestArtifacts(*cacheDir, *forceDownload)
	if err != nil {
		fatalf("failed to prepare Klaus functional test artifacts: %v", err)
	}
	if err := verifyFunctionalPassTrap(lstPath, pass); err != nil {
		fatalf("functional test listing does not match configured pass trap: %v", err)
	}

	fmt.Printf("Using Klaus functional test binary: %s\n", binPath)
	fmt.Printf("Using Klaus functional test listing: %s\n", lstPath)
	fmt.Printf("Start: %s  Pass trap: %s  Max instructions: %d\n", formatAddress(start), formatAddress(pass), *maxInstructions)

	functionalResult, err := runFunctionalTest(binPath, start, pass, *maxInstructions)
	if err != nil {
		fatalf("%v", err)
	}
	if start == defaultStartAddress && pass == defaultSuccessAddress {
		if functionalResult.Instructions != functionalInstructions || functionalResult.Cycles != functionalCycles {
			fatalf("functional test timing mismatch: got %d instructions / %d cycles; expected %d instructions / %d cycles",
				functionalResult.Instructions, functionalResult.Cycles, functionalInstructions, functionalCycles)
		}
	}

	fmt.Printf("PASS: Klaus 6502 functional test reached %s after %d instructions / %d cycles\n", formatAddress(functionalResult.StopAddress), functionalResult.Instructions, functionalResult.Cycles)
	fmt.Printf("Using Klaus decimal test source: %s\n", decimalSourcePath)
	decimalResult, err := runDecimalSourceTest(decimalSourcePath, *maxInstructions)
	if err != nil {
		fatalf("decimal source test failed: %v", err)
	}

	fmt.Printf("PASS: Klaus decimal mode test completed at %s after %d instructions / %d cycles\n", formatAddress(decimalResult.StopAddress), decimalResult.Instructions, decimalResult.Cycles)
}

func runDecimalSourceTest(sourcePath string, maxInstructions uint64) (runResult, error) {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return runResult{}, fmt.Errorf("open source test: %w", err)
	}
	defer sourceFile.Close()

	mem := memory.NewMemory[uint16](64 * 1024)
	testCPU := cpu.NewCPU(mem, false)
	segments, err := assembler.New(testCPU.OpCodes()).Assemble(sourceFile, sourcePath)
	if err != nil {
		return runResult{}, fmt.Errorf("assemble source test: %w", err)
	}
	for _, segment := range segments {
		mem.Write(segment.StartAddress, segment.Data.Bytes()...)
	}

	testCPU.Registers().PC = decimalTestStart
	disasm := debugger.NewDisassembler(mem, testCPU.OpCodes())
	result := runResult{}
	for result.Instructions < maxInstructions {
		currentPC := testCPU.Registers().PC
		if mem.Read(currentPC) == decimalTestStopOpcode {
			if mem.Read(decimalTestError) != 0 {
				return runResult{}, fmt.Errorf("source test reported ERROR=%d at %s after %d instructions / %d cycles (N1=$%02X N2=$%02X DA=$%02X AR=$%02X CF=$%02X)\n%s", mem.Read(decimalTestError), formatAddress(currentPC), result.Instructions, result.Cycles, mem.Read(0x0000), mem.Read(0x0001), mem.Read(0x0004), mem.Read(0x0006), mem.Read(0x000A), formatCPUState(testCPU, disasm))
			}
			result.StopAddress = currentPC
			return result, nil
		}

		cycles, err := executeInstruction(testCPU)
		if err != nil {
			return runResult{}, fmt.Errorf("execution error at %s after %d instructions / %d cycles: %w\n%s", formatAddress(currentPC), result.Instructions, result.Cycles, err, formatCPUState(testCPU, disasm))
		}
		result.Instructions++
		result.Cycles += cycles
	}

	return runResult{}, fmt.Errorf("source test timed out after %d completed instructions / %d cycles\n%s", maxInstructions, result.Cycles, formatCPUState(testCPU, disasm))
}

func ensureKlausTestArtifacts(cacheDir string, force bool) (string, string, string, error) {
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return "", "", "", err
	}

	binPath := filepath.Join(cacheDir, functionalTestBinName)
	lstPath := filepath.Join(cacheDir, functionalTestLstName)
	decimalSourcePath := filepath.Join(cacheDir, decimalTestSourceName)

	if err := downloadIfNeeded(functionalTestBinURL, binPath, force); err != nil {
		return "", "", "", err
	}
	if err := downloadIfNeeded(functionalTestLstURL, lstPath, force); err != nil {
		return "", "", "", err
	}
	if err := downloadIfNeeded(decimalTestSourceURL, decimalSourcePath, force); err != nil {
		return "", "", "", err
	}

	return binPath, lstPath, decimalSourcePath, nil
}

func downloadIfNeeded(url, path string, force bool) error {
	if !force {
		if info, err := os.Stat(path); err == nil && info.Size() > 0 {
			return nil
		} else if err != nil && !errors.Is(err, os.ErrNotExist) {
			return err
		}
	}

	fmt.Printf("Downloading %s\n", url)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GET %s: %s", url, resp.Status)
	}

	tmpPath := path + ".tmp"
	out, err := os.OpenFile(tmpPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, resp.Body); err != nil {
		out.Close()
		return err
	}
	if err := out.Close(); err != nil {
		return err
	}

	return os.Rename(tmpPath, path)
}

func runFunctionalTest(binPath string, startAddress, successAddress uint16, maxInstructions uint64) (runResult, error) {
	image, err := os.ReadFile(binPath)
	if err != nil {
		return runResult{}, err
	}
	if len(image) > 64*1024 {
		return runResult{}, fmt.Errorf("test binary is %d bytes; expected at most 65536", len(image))
	}

	mem := memory.NewMemory[uint16](64 * 1024)
	mem.Write(0, image...)

	testCPU := cpu.NewCPU(mem, false)
	testCPU.Registers().PC = startAddress

	disasm := debugger.NewDisassembler(mem, testCPU.OpCodes())
	lastPC := testCPU.Registers().PC
	result := runResult{}

	for result.Instructions < maxInstructions {
		currentPC := testCPU.Registers().PC
		cycles, err := executeInstruction(testCPU)
		if err != nil {
			return runResult{}, fmt.Errorf("execution error at %s after %d instructions / %d cycles: %w\n%s", formatAddress(currentPC), result.Instructions, result.Cycles, err, formatCPUState(testCPU, disasm))
		}
		result.Instructions++
		result.Cycles += cycles

		pc := testCPU.Registers().PC
		if pc == lastPC {
			if pc == successAddress {
				result.StopAddress = pc
				return result, nil
			}
			return runResult{}, fmt.Errorf("test trapped at %s after %d instructions / %d cycles\n%s", formatAddress(pc), result.Instructions, result.Cycles, formatCPUState(testCPU, disasm))
		}
		lastPC = pc
	}

	return runResult{}, fmt.Errorf("timed out after %d completed instructions / %d cycles\n%s", maxInstructions, result.Cycles, formatCPUState(testCPU, disasm))
}

func executeInstruction(testCPU *cpu.CPU) (uint64, error) {
	var cycles uint64
	completed := cpu.Completed(false)
	for !completed {
		done, err := testCPU.Execute()
		cycles++
		if err != nil {
			return cycles, err
		}
		completed = done
	}
	return cycles, nil
}

func verifyFunctionalPassTrap(lstPath string, passAddress uint16) error {
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

func parseAddress(s string) (uint16, error) {
	value := strings.TrimSpace(s)
	value = strings.TrimPrefix(value, "$")
	if strings.HasPrefix(value, "0x") || strings.HasPrefix(value, "0X") {
		value = value[2:]
	}
	parsed, err := strconv.ParseUint(value, 16, 16)
	if err != nil {
		return 0, err
	}
	return uint16(parsed), nil
}

func formatAddress(address uint16) string {
	return fmt.Sprintf("$%04X", address)
}

func formatCPUState(testCPU *cpu.CPU, disasm *debugger.Disassembler) string {
	reg := testCPU.Registers()
	instruction, _ := disasm.Disassemble(reg.PC)
	return fmt.Sprintf(
		"PC=%s A=$%02X X=$%02X Y=$%02X S=$%02X P=$%02X\nNext: %s",
		formatAddress(reg.PC),
		reg.A,
		reg.X,
		reg.Y,
		reg.S,
		reg.Status,
		instruction,
	)
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "klaustests: "+format+"\n", args...)
	os.Exit(1)
}
