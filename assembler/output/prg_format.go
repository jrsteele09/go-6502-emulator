package output

import (
	"bytes"
	"fmt"
	"os"

	"github.com/jrsteele09/go-6502-emulator/assembler"
	"github.com/jrsteele09/go-6502-emulator/utils"
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
// Note: PRG format requires contiguous memory, so all segments must be adjacent
func (p *PRGFormat) CreateFile(filename string, segments []assembler.AssembledData, verbose bool) error {
	// Build PRG bytes using CreateData, then write them to disk
	prgData, err := p.CreateData(segments)
	if err != nil {
		return err
	}

	if verbose && len(prgData) >= 2 {
		loadAddr := uint16(prgData[0]) | (uint16(prgData[1]) << 8)
		fmt.Printf("PRG load address: $%04X\n", loadAddr)
		fmt.Printf("Program size: %d bytes\n", len(prgData)-2)
	}

	// Write to file
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	if _, err = file.Write(prgData); err != nil {
		return fmt.Errorf("failed to write PRG data: %w", err)
	}

	return nil
}

// validateAndMergeContiguousSegments checks if segments are contiguous and merges them
// Returns the merged data, load address, and error if segments have gaps
func (p *PRGFormat) validateAndMergeContiguousSegments(segments []assembler.AssembledData) ([]byte, uint16, error) {
	if len(segments) == 0 {
		return nil, 0, fmt.Errorf("no segments provided")
	}

	// Sort segments by start address to check for contiguity
	sortedSegments := make([]assembler.AssembledData, len(segments))
	copy(sortedSegments, segments)

	// Simple bubble sort by start address
	for i := 0; i < len(sortedSegments); i++ {
		for j := i + 1; j < len(sortedSegments); j++ {
			if sortedSegments[i].StartAddress > sortedSegments[j].StartAddress {
				sortedSegments[i], sortedSegments[j] = sortedSegments[j], sortedSegments[i]
			}
		}
	}

	loadAddr := sortedSegments[0].StartAddress
	var mergedData []byte

	// Merge all segments, filling any gaps with zeros. Overlaps are errors.
	expectedNextAddr := loadAddr
	for i, segment := range sortedSegments {
		// Overlap check (shouldn't happen due to sort and logic)
		if i > 0 && segment.StartAddress < expectedNextAddr {
			return nil, 0, fmt.Errorf("segments overlap: segment at $%04X overlaps with previous segment", segment.StartAddress)
		}

		if segment.StartAddress > expectedNextAddr {
			gap := int(segment.StartAddress - expectedNextAddr)
			// Fill gap with zeros
			mergedData = append(mergedData, make([]byte, gap)...)
		}

		// Append segment data
		segBytes := segment.Data.Bytes()
		mergedData = append(mergedData, segBytes...)
		expectedNextAddr = segment.StartAddress + uint16(len(segBytes))
	}

	return mergedData, loadAddr, nil
}

// CreateData creates PRG format data in memory from assembled segments
// Returns the PRG data as a byte slice (2-byte load address + program data)
// Note: PRG format requires contiguous memory, so all segments must be adjacent
func (p *PRGFormat) CreateData(segments []assembler.AssembledData) ([]byte, error) {
	if len(segments) == 0 {
		return nil, fmt.Errorf("no segments to convert")
	}

	if len(segments) > 1 {
		return nil, fmt.Errorf("multiple segments not supported in PRG format")
	}

	// Validate segments are contiguous
	contiguousData, loadAddr, err := p.validateAndMergeContiguousSegments(segments)
	if err != nil {
		return nil, fmt.Errorf("PRG format error: %w", err)
	}

	// Create output buffer: 2 bytes for load address + program data
	output := make([]byte, 2+len(contiguousData))

	// Write load address (little-endian)
	output[0] = byte(loadAddr & 0xFF)
	output[1] = byte(loadAddr >> 8)

	// Copy program data
	copy(output[2:], contiguousData)

	return output, nil
}

// LoadFile loads a PRG file and returns assembled segments
func (p *PRGFormat) LoadFile(filename string, verbose bool) ([]assembler.AssembledData, error) {
	// Read the file
	if verbose {
		fmt.Printf("Loading PRG file: %s\n", filename)
	}
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
			Data:         utils.Value(bytes.NewBuffer(programData)),
		},
	}

	return segments, nil
}
