package output

import (
	"bytes"
	"fmt"
	"os"
	"sort"

	"github.com/jrsteele09/go-6502-emulator/assembler"
	"github.com/jrsteele09/go-6502-emulator/utils"
)

var _ BinaryFormat = (*BINFormat)(nil)

// BINFormat reads and writes headerless binary files. LoadAddress is used only
// when reading because a BIN file does not contain an address of its own.
type BINFormat struct {
	LoadAddress uint16
}

// NewBINFormat creates a raw binary formatter. loadAddress is used by LoadFile.
func NewBINFormat(loadAddress uint16) *BINFormat {
	return &BINFormat{LoadAddress: loadAddress}
}

// CreateFile writes assembled bytes without a header or load address.
func (b *BINFormat) CreateFile(filename string, segments []assembler.AssembledData, verbose bool) error {
	data, startAddress, err := b.createData(segments)
	if err != nil {
		return err
	}

	if err := os.WriteFile(filename, data, 0o644); err != nil {
		return fmt.Errorf("failed to write BIN data: %w", err)
	}

	if verbose {
		fmt.Printf("BIN start address: $%04X\n", startAddress)
		fmt.Printf("BIN size: %d bytes\n", len(data))
	}
	return nil
}

// CreateData returns assembled bytes without a header or load address.
func (b *BINFormat) CreateData(segments []assembler.AssembledData) ([]byte, error) {
	data, _, err := b.createData(segments)
	return data, err
}

func (b *BINFormat) createData(segments []assembler.AssembledData) ([]byte, uint16, error) {
	if len(segments) == 0 {
		return nil, 0, fmt.Errorf("no segments to convert to BIN")
	}

	sortedSegments := append([]assembler.AssembledData(nil), segments...)
	sort.SliceStable(sortedSegments, func(i, j int) bool {
		return sortedSegments[i].StartAddress < sortedSegments[j].StartAddress
	})

	startAddress := sortedSegments[0].StartAddress
	nextAddress := uint32(startAddress)
	var data []byte
	for _, segment := range sortedSegments {
		segmentStart := uint32(segment.StartAddress)
		segmentData := segment.Data.Bytes()
		segmentEnd := segmentStart + uint32(len(segmentData))
		if segmentEnd > 1<<16 {
			return nil, 0, fmt.Errorf("segment at $%04X extends beyond $FFFF", segment.StartAddress)
		}
		if segmentStart < nextAddress {
			return nil, 0, fmt.Errorf("segments overlap at $%04X", segment.StartAddress)
		}

		data = append(data, make([]byte, int(segmentStart-nextAddress))...)
		data = append(data, segmentData...)
		nextAddress = segmentEnd
	}

	return data, startAddress, nil
}

// LoadFile reads raw bytes and places them in a segment at LoadAddress.
func (b *BINFormat) LoadFile(filename string, verbose bool) ([]assembler.AssembledData, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read BIN file: %w", err)
	}
	if len(data) == 0 {
		return nil, fmt.Errorf("BIN file is empty")
	}
	if int(b.LoadAddress)+len(data) > 1<<16 {
		return nil, fmt.Errorf("BIN file does not fit at $%04X", b.LoadAddress)
	}

	if verbose {
		fmt.Printf("BIN load address: $%04X\n", b.LoadAddress)
		fmt.Printf("BIN size: %d bytes\n", len(data))
	}

	return []assembler.AssembledData{{
		StartAddress: b.LoadAddress,
		Data:         utils.Value(bytes.NewBuffer(data)),
	}}, nil
}
