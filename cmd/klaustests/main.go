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

	"github.com/jrsteele09/go-6502-emulator/cpu"
	"github.com/jrsteele09/go-6502-emulator/debugger"
	"github.com/jrsteele09/go-6502-emulator/memory"
)

const (
	functionalTestBinURL = "https://raw.githubusercontent.com/Klaus2m5/6502_65C02_functional_tests/master/bin_files/6502_functional_test.bin"
	functionalTestLstURL = "https://raw.githubusercontent.com/Klaus2m5/6502_65C02_functional_tests/master/bin_files/6502_functional_test.lst"

	defaultCacheDir        = "testdata/klaus"
	functionalTestBinName  = "6502_functional_test.bin"
	functionalTestLstName  = "6502_functional_test.lst"
	defaultStartAddress    = uint16(0x0400)
	defaultSuccessAddress  = uint16(0x3469)
	defaultMaxInstructions = uint64(900_000_000)
)

func main() {
	cacheDir := flag.String("cache-dir", defaultCacheDir, "directory used to cache downloaded Klaus test artifacts")
	forceDownload := flag.Bool("force", false, "download test artifacts even when cached files exist")
	startAddress := flag.String("start", formatAddress(defaultStartAddress), "program start address")
	successAddress := flag.String("pass", formatAddress(defaultSuccessAddress), "success trap address")
	maxInstructions := flag.Uint64("max-instructions", defaultMaxInstructions, "maximum completed instructions before failing")
	flag.Parse()

	start, err := parseAddress(*startAddress)
	if err != nil {
		fatalf("invalid -start address: %v", err)
	}
	pass, err := parseAddress(*successAddress)
	if err != nil {
		fatalf("invalid -pass address: %v", err)
	}

	binPath, lstPath, err := ensureFunctionalTestArtifacts(*cacheDir, *forceDownload)
	if err != nil {
		fatalf("failed to prepare Klaus functional test artifacts: %v", err)
	}

	fmt.Printf("Using Klaus functional test binary: %s\n", binPath)
	fmt.Printf("Using Klaus functional test listing: %s\n", lstPath)
	fmt.Printf("Start: %s  Pass trap: %s  Max instructions: %d\n", formatAddress(start), formatAddress(pass), *maxInstructions)

	if err := runFunctionalTest(binPath, start, pass, *maxInstructions); err != nil {
		fatalf("%v", err)
	}

	fmt.Println("PASS: Klaus 6502 functional test reached the success trap")
}

func ensureFunctionalTestArtifacts(cacheDir string, force bool) (string, string, error) {
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		return "", "", err
	}

	binPath := filepath.Join(cacheDir, functionalTestBinName)
	lstPath := filepath.Join(cacheDir, functionalTestLstName)

	if err := downloadIfNeeded(functionalTestBinURL, binPath, force); err != nil {
		return "", "", err
	}
	if err := downloadIfNeeded(functionalTestLstURL, lstPath, force); err != nil {
		return "", "", err
	}

	return binPath, lstPath, nil
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

func runFunctionalTest(binPath string, startAddress, successAddress uint16, maxInstructions uint64) error {
	image, err := os.ReadFile(binPath)
	if err != nil {
		return err
	}
	if len(image) > 64*1024 {
		return fmt.Errorf("test binary is %d bytes; expected at most 65536", len(image))
	}

	mem := memory.NewMemory[uint16](64 * 1024)
	mem.Write(0, image...)

	testCPU := cpu.NewCPU(mem, false)
	testCPU.Registers().PC = startAddress

	disasm := debugger.NewDisassembler(mem, testCPU.OpCodes())
	lastPC := testCPU.Registers().PC

	for instructions := uint64(0); instructions < maxInstructions; instructions++ {
		currentPC := testCPU.Registers().PC
		completed := cpu.Completed(false)
		for !completed {
			done, err := testCPU.Execute()
			if err != nil {
				return fmt.Errorf("execution error at %s after %d instructions: %w\n%s", formatAddress(currentPC), instructions, err, formatCPUState(testCPU, disasm))
			}
			completed = done
		}

		pc := testCPU.Registers().PC
		if pc == lastPC {
			if pc == successAddress {
				return nil
			}
			return fmt.Errorf("test trapped at %s after %d instructions\n%s", formatAddress(pc), instructions+1, formatCPUState(testCPU, disasm))
		}
		lastPC = pc
	}

	return fmt.Errorf("timed out after %d completed instructions\n%s", maxInstructions, formatCPUState(testCPU, disasm))
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
