package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jrsteele09/go-6502-emulator/assembler"
	"github.com/jrsteele09/go-6502-emulator/assembler/bin"
	"github.com/jrsteele09/go-6502-emulator/cpu"
	"github.com/jrsteele09/go-6502-emulator/memory"
)

const (
	version = "1.0.0"
)

func main() {
	var (
		inputFile    = flag.String("i", "", "Input assembly file (required)")
		outputFile   = flag.String("o", "", "Output file (default: input filename with appropriate extension)")
		outputFormat = flag.String("f", "prg", "Output format: prg, d64, or t64 (default: prg)")
		showHelp     = flag.Bool("h", false, "Show help")
		showVer      = flag.Bool("version", false, "Show version")
		verbose      = flag.Bool("v", false, "Verbose output")
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "6502 Assembler v%s\n\n", version)
		fmt.Fprintf(os.Stderr, "Usage: %s [options] -i <input.asm>\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nOutput Formats:\n")
		fmt.Fprintf(os.Stderr, "  prg  - Commodore 64 PRG file (default)\n")
		fmt.Fprintf(os.Stderr, "  d64  - Commodore 64 disk image (1541 format)\n")
		fmt.Fprintf(os.Stderr, "  t64  - Commodore 64 tape archive\n")
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -i game.asm                    # Output to game.prg\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -i game.asm -f d64             # Output to game.d64\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -i game.asm -o disk.d64 -f d64 # Output to disk.d64\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -i game.asm -f t64 -v          # Output to game.t64 with verbose\n", os.Args[0])
	}

	flag.Parse()

	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}

	if *showVer {
		fmt.Printf("6502 Assembler v%s\n", version)
		os.Exit(0)
	}

	if *inputFile == "" {
		fmt.Fprintf(os.Stderr, "Error: Input file is required\n\n")
		flag.Usage()
		os.Exit(1)
	}

	// Validate output format
	*outputFormat = strings.ToLower(*outputFormat)
	validFormats := map[string]bool{"prg": true, "d64": true, "t64": true}
	if !validFormats[*outputFormat] {
		fmt.Fprintf(os.Stderr, "Error: Invalid output format '%s'. Valid formats: prg, d64, t64\n", *outputFormat)
		os.Exit(1)
	}

	// Check if input file exists
	if _, err := os.Stat(*inputFile); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Input file '%s' does not exist\n", *inputFile)
		os.Exit(1)
	}

	// Determine output file name
	if *outputFile == "" {
		ext := filepath.Ext(*inputFile)
		baseName := strings.TrimSuffix(*inputFile, ext)
		*outputFile = baseName + "." + *outputFormat
	}

	if *verbose {
		fmt.Printf("6502 Assembler v%s\n", version)
		fmt.Printf("Input file:    %s\n", *inputFile)
		fmt.Printf("Output file:   %s\n", *outputFile)
		fmt.Printf("Output format: %s\n", strings.ToUpper(*outputFormat))
		fmt.Println()
	}

	// Create assembler with full 6502 instruction set
	opcodes := createOpcodes()
	asm := assembler.New(opcodes)

	// Open input file
	file, err := os.Open(*inputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to open input file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Set up file resolver for includes
	baseDir := filepath.Dir(*inputFile)
	resolver := assembler.NewOSFileResolver(baseDir)

	if *verbose {
		fmt.Println("Assembling...")
	}

	// Assemble with preprocessor support for includes
	segments, err := asm.AssembleWithPreprocessor(file, resolver)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Assembly failed: %v\n", err)
		os.Exit(1)
	}

	if len(segments) == 0 {
		fmt.Fprintf(os.Stderr, "Error: No code generated\n")
		os.Exit(1)
	}

	// Create binary format based on output format flag
	var format bin.BinaryFormat
	switch *outputFormat {
	case "prg":
		format = bin.NewPRGFormat()
	case "d64":
		// Extract base name for disk name
		diskName := strings.ToUpper(filepath.Base(strings.TrimSuffix(*outputFile, filepath.Ext(*outputFile))))
		if len(diskName) > 16 {
			diskName = diskName[:16]
		}
		format = bin.NewD64Format(diskName, "01")
	case "t64":
		// Extract base name for tape name
		tapeName := strings.ToUpper(filepath.Base(strings.TrimSuffix(*outputFile, filepath.Ext(*outputFile))))
		if len(tapeName) > 24 {
			tapeName = tapeName[:24]
		}
		format = bin.NewT64Format(tapeName, 30)
	}

	err = format.CreateFile(*outputFile, segments, *verbose)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating %s file: %v\n", strings.ToUpper(*outputFormat), err)
		os.Exit(1)
	}

	if *verbose {
		fmt.Printf("Assembly completed successfully!\n")
		fmt.Printf("Generated %d segment(s)\n", len(segments))
		for i, segment := range segments {
			fmt.Printf("  Segment %d: %d bytes at $%04X\n", i+1, len(segment.Data.Bytes()), segment.StartAddress)
		}
	} else {
		fmt.Printf("Assembly successful: %s -> %s\n", *inputFile, *outputFile)
	}
}

// createOpcodes creates the full 6502 instruction set
func createOpcodes() []*cpu.OpCodeDef {
	mem := memory.NewMemory[uint16](64 * 1024)
	testCPU := cpu.NewCPU(mem)
	return cpu.OpCodes(testCPU)
}
