package output

import "github.com/jrsteele09/go-6502-emulator/assembler"

// BinaryFormat defines the interface for binary format generators
type BinaryFormat interface {
	// CreateFile creates a binary file from assembled segments
	CreateFile(filename string, segments []assembler.AssembledData, verbose bool) error
	// CreateData creates binary format data in memory from assembled segments
	CreateData(segments []assembler.AssembledData) ([]byte, error)
	// LoadFile loads a binary file and returns assembled segments
	LoadFile(filename string, verbose bool) ([]assembler.AssembledData, error)
}
