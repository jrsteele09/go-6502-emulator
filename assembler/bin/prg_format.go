package bin

import (
	"fmt"
	"os"

	"github.com/jrsteele09/go-6502-emulator/assembler"
)

var _ BinaryFormat = (*PRGFormat)(nil)

// PRGFormat implements the BinaryFormat interface for Commodore 64 PRG files
type PRGFormat struct{}

// NewPRGFormat creates a new PRG format generator
func NewPRGFormat() *PRGFormat {
	return &PRGFormat{}
}

// CreateFile creates a Commodore 64 PRG file from assembled segments
// PRG format: 2-byte load address (little-endian) followed by program data
func (p *PRGFormat) CreateFile(filename string, segments []assembler.AssembledData, verbose bool) error {
	if len(segments) == 0 {
		return fmt.Errorf("no segments to write")
	}

	// Find the lowest start address for the load address
	loadAddr := segments[0].StartAddress
	for _, segment := range segments {
		if segment.StartAddress < loadAddr {
			loadAddr = segment.StartAddress
		}
	}

	// Calculate total size needed
	maxAddr := uint16(0)
	for _, segment := range segments {
		endAddr := segment.StartAddress + uint16(len(segment.Data))
		if endAddr > maxAddr {
			maxAddr = endAddr
		}
	}

	// Create output buffer: 2 bytes for load address + program data
	totalSize := int(maxAddr - loadAddr)
	output := make([]byte, 2+totalSize)

	// Write load address (little-endian)
	output[0] = byte(loadAddr & 0xFF)
	output[1] = byte(loadAddr >> 8)

	if verbose {
		fmt.Printf("PRG load address: $%04X\n", loadAddr)
		fmt.Printf("Program size: %d bytes\n", totalSize)
	}

	// Copy segment data to appropriate positions
	for _, segment := range segments {
		offset := int(segment.StartAddress - loadAddr)
		copy(output[2+offset:], segment.Data)

		if verbose {
			fmt.Printf("  Writing segment at $%04X: %d bytes\n", segment.StartAddress, len(segment.Data))
		}
	}

	// Write to file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(output)
	if err != nil {
		return fmt.Errorf("failed to write PRG data: %w", err)
	}

	return nil
}

// CreateData creates PRG format data in memory from assembled segments
// Returns the PRG data as a byte slice (2-byte load address + program data)
func (p *PRGFormat) CreateData(segments []assembler.AssembledData) ([]byte, error) {
	if len(segments) == 0 {
		return nil, fmt.Errorf("no segments to convert")
	}

	// Find the lowest start address for the load address
	loadAddr := segments[0].StartAddress
	for _, segment := range segments {
		if segment.StartAddress < loadAddr {
			loadAddr = segment.StartAddress
		}
	}

	// Calculate total size needed
	maxAddr := uint16(0)
	for _, segment := range segments {
		endAddr := segment.StartAddress + uint16(len(segment.Data))
		if endAddr > maxAddr {
			maxAddr = endAddr
		}
	}

	// Create output buffer: 2 bytes for load address + program data
	totalSize := int(maxAddr - loadAddr)
	output := make([]byte, 2+totalSize)

	// Write load address (little-endian)
	output[0] = byte(loadAddr & 0xFF)
	output[1] = byte(loadAddr >> 8)

	// Copy segment data to appropriate positions
	for _, segment := range segments {
		offset := int(segment.StartAddress - loadAddr)
		copy(output[2+offset:], segment.Data)
	}

	return output, nil
}

// LoadFile loads a PRG file and returns assembled segments
func (p *PRGFormat) LoadFile(filename string, verbose bool) ([]assembler.AssembledData, error) {
	// Read the file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read PRG file: %w", err)
	}

	if len(data) < 2 {
		return nil, fmt.Errorf("PRG file too small (minimum 2 bytes for load address)")
	}

	// Extract load address (little-endian)
	loadAddr := uint16(data[0]) | (uint16(data[1]) << 8)

	// Extract program data
	programData := data[2:]

	if verbose {
		fmt.Printf("PRG load address: $%04X\n", loadAddr)
		fmt.Printf("Program size: %d bytes\n", len(programData))
	}

	// Create a single segment with all the data
	segments := []assembler.AssembledData{
		{
			StartAddress: loadAddr,
			Data:         programData,
		},
	}

	return segments, nil
}
