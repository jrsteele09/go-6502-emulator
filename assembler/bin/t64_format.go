package bin

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"
	"strings"

	"github.com/jrsteele09/go-6502-emulator/assembler"
	"github.com/jrsteele09/go-6502-emulator/utils"
)

var _ BinaryFormat = (*T64Format)(nil)

// T64Format implements the BinaryFormat interface for Commodore 64 tape archives
// T64 format: Tape archive format that can contain multiple PRG files
type T64Format struct {
	// TapeName is the name of the tape (up to 24 characters)
	TapeName string
	// FileName is the name of the program file (up to 16 characters)
	FileName string
	// MaxEntries is the maximum number of entries in the tape
	MaxEntries uint16
	// UsedEntries is the number of used entries
	UsedEntries uint16
}

// T64 format constants
const (
	T64HeaderSize  = 64 // Size of T64 header
	T64EntrySize   = 32 // Size of each directory entry
	T64MaxTapeName = 24 // Maximum tape name length
	T64MaxFilename = 16 // Maximum filename length
	T64Signature   = "C64 tape image file"
	T64Version     = 0x0100 // Version 1.0
)

// NewT64Format creates a new T64 format generator
func NewT64Format(tapeName string, maxEntries uint16) *T64Format {
	if len(tapeName) > T64MaxTapeName {
		tapeName = tapeName[:T64MaxTapeName]
	}
	if maxEntries == 0 {
		maxEntries = 30 // Default max entries
	}
	return &T64Format{
		TapeName:    tapeName,
		FileName:    "PROGRAM", // Default filename
		MaxEntries:  maxEntries,
		UsedEntries: 0,
	}
}

// NewT64FormatWithFilename creates a new T64 format generator with custom filename
func NewT64FormatWithFilename(tapeName, fileName string, maxEntries uint16) *T64Format {
	if len(tapeName) > T64MaxTapeName {
		tapeName = tapeName[:T64MaxTapeName]
	}
	if len(fileName) > T64MaxFilename {
		fileName = fileName[:T64MaxFilename]
	}
	if fileName == "" {
		fileName = "PROGRAM" // Default filename
	}
	if maxEntries == 0 {
		maxEntries = 30 // Default max entries
	}
	return &T64Format{
		TapeName:    tapeName,
		FileName:    fileName,
		MaxEntries:  maxEntries,
		UsedEntries: 0,
	}
}

// CreateFile creates a T64 tape archive file from assembled segments
func (t *T64Format) CreateFile(filename string, segments []assembler.AssembledData, verbose bool) error {
	data, err := t.CreateData(segments)
	if err != nil {
		return err
	}

	if verbose {
		fmt.Printf("T64 tape name: %s\n", t.TapeName)
		fmt.Printf("T64 max entries: %d\n", t.MaxEntries)
		fmt.Printf("T64 used entries: %d\n", t.UsedEntries)
		fmt.Printf("T64 archive size: %d bytes\n", len(data))
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create T64 file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write T64 data: %w", err)
	}

	return nil
}

// CreateData creates T64 format data in memory from assembled segments
func (t *T64Format) CreateData(segments []assembler.AssembledData) ([]byte, error) {
	if len(segments) == 0 {
		return nil, fmt.Errorf("no segments to convert to T64")
	}

	// Create PRG data from segments
	prgData, loadAddr, err := t.createPRGFromSegments(segments)
	if err != nil {
		return nil, err
	}

	// Calculate total size needed
	headerSize := T64HeaderSize
	directorySize := int(t.MaxEntries) * T64EntrySize
	dataSize := len(prgData)
	totalSize := headerSize + directorySize + dataSize

	// Create T64 archive
	archive := make([]byte, totalSize)

	// Write header
	t.writeHeader(archive)

	// Write directory entry
	t.writeDirectoryEntry(archive, t.FileName, loadAddr, uint16(len(prgData)), headerSize+directorySize)

	// Write program data
	copy(archive[headerSize+directorySize:], prgData)

	// Update used entries count
	t.UsedEntries = 1

	// Update the used entries in the header
	archive[36] = byte(t.UsedEntries & 0xFF)
	archive[37] = byte(t.UsedEntries >> 8)

	return archive, nil
}

// createPRGFromSegments creates PRG format data from segments
func (t *T64Format) createPRGFromSegments(segments []assembler.AssembledData) ([]byte, uint16, error) {
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
		endAddr := segment.StartAddress + uint16(len(segment.Data.Bytes()))
		if endAddr > maxAddr {
			maxAddr = endAddr
		}
	}

	// Create PRG data (load address + program data)
	totalSize := int(maxAddr - loadAddr)
	prgData := make([]byte, 2+totalSize)

	// Write load address (little-endian)
	prgData[0] = byte(loadAddr & 0xFF)
	prgData[1] = byte(loadAddr >> 8)

	// Copy segment data to appropriate positions
	for _, segment := range segments {
		offset := int(segment.StartAddress - loadAddr)
		copy(prgData[2+offset:], segment.Data.Bytes())
	}

	return prgData, loadAddr, nil
}

// writeHeader writes the T64 header
func (t *T64Format) writeHeader(archive []byte) {
	// T64 signature (32 bytes)
	copy(archive[0:], t.padString(T64Signature, 32))

	// Version (2 bytes, little-endian)
	archive[32] = byte(T64Version & 0xFF)
	archive[33] = byte(T64Version >> 8)

	// Maximum number of entries (2 bytes, little-endian)
	archive[34] = byte(t.MaxEntries & 0xFF)
	archive[35] = byte(t.MaxEntries >> 8)

	// Used entries (2 bytes, little-endian) - will be updated later
	archive[36] = byte(t.UsedEntries & 0xFF)
	archive[37] = byte(t.UsedEntries >> 8)

	// Reserved (2 bytes)
	archive[38] = 0
	archive[39] = 0

	// Tape name (24 bytes, padded with spaces)
	copy(archive[40:], t.padString(t.TapeName, T64MaxTapeName))
}

// writeDirectoryEntry writes a directory entry for the program
func (t *T64Format) writeDirectoryEntry(archive []byte, filename string, loadAddr, fileSize uint16, dataOffset int) {
	entryOffset := T64HeaderSize // First directory entry

	// Entry type (1 = Normal tape file)
	archive[entryOffset] = 1

	// File type (PRG = $82)
	archive[entryOffset+1] = 0x82

	// Start address (2 bytes, little-endian)
	archive[entryOffset+2] = byte(loadAddr & 0xFF)
	archive[entryOffset+3] = byte(loadAddr >> 8)

	// End address (2 bytes, little-endian)
	// In T64 format, end address is the last byte address (inclusive)
	// File size includes the 2-byte load address, so actual program size is fileSize - 2
	// Some emulators expect end address to be exclusive (last address + 1)
	// We use the standard inclusive format: endAddr = startAddr + programSize - 1
	endAddr := loadAddr + (fileSize - 2) - 1 // -1 because end address is inclusive
	archive[entryOffset+4] = byte(endAddr & 0xFF)
	archive[entryOffset+5] = byte(endAddr >> 8)

	// Reserved (2 bytes)
	archive[entryOffset+6] = 0
	archive[entryOffset+7] = 0

	// Offset to file data (4 bytes, little-endian)
	archive[entryOffset+8] = byte(dataOffset & 0xFF)
	archive[entryOffset+9] = byte((dataOffset >> 8) & 0xFF)
	archive[entryOffset+10] = byte((dataOffset >> 16) & 0xFF)
	archive[entryOffset+11] = byte((dataOffset >> 24) & 0xFF)

	// Reserved (4 bytes)
	for i := 12; i < 16; i++ {
		archive[entryOffset+i] = 0
	}

	// Filename (16 bytes, padded with spaces)
	copy(archive[entryOffset+16:], t.padString(filename, T64MaxFilename))
}

// padString pads a string with spaces to specified length
func (t *T64Format) padString(s string, length int) []byte {
	result := make([]byte, length)
	copy(result, []byte(s))
	for i := len(s); i < length; i++ {
		result[i] = ' '
	}
	return result
}

// LoadFile loads a T64 tape archive and returns assembled segments from contained files
func (t *T64Format) LoadFile(filename string, verbose bool) ([]assembler.AssembledData, error) {
	// Read the file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read T64 file: %w", err)
	}

	if len(data) < T64HeaderSize {
		return nil, fmt.Errorf("invalid T64 file: too small (minimum %d bytes)", T64HeaderSize)
	}

	// Verify T64 signature
	signature := string(data[0:32])
	if !strings.HasPrefix(signature, "C64") {
		return nil, fmt.Errorf("invalid T64 file: bad signature")
	}

	// Read header
	maxEntries := binary.LittleEndian.Uint16(data[34:36])
	usedEntries := binary.LittleEndian.Uint16(data[36:38])
	tapeName := string(data[40:64])

	// Update our format info
	t.TapeName = tapeName

	if verbose {
		fmt.Printf("T64 tape name: %s\n", tapeName)
		fmt.Printf("Used entries: %d, Max entries: %d\n", usedEntries, maxEntries)
	}

	// Process directory entries
	segments := make([]assembler.AssembledData, 0)

	for i := uint16(0); i < usedEntries && i < maxEntries; i++ {
		entryOffset := T64HeaderSize + (i * T64EntrySize)

		if int(entryOffset)+T64EntrySize > len(data) {
			break
		}

		// Read directory entry
		entryType := data[entryOffset]
		fileType := data[entryOffset+1]
		startAddr := binary.LittleEndian.Uint16(data[entryOffset+2 : entryOffset+4])
		endAddr := binary.LittleEndian.Uint16(data[entryOffset+4 : entryOffset+6])
		fileOffset := binary.LittleEndian.Uint32(data[entryOffset+8 : entryOffset+12])
		fileName := string(data[entryOffset+16 : entryOffset+32])

		// Remove padding
		fileName = strings.TrimRight(fileName, "\x00 ")

		// Skip empty or non-PRG entries
		if entryType == 0 || fileType != 0x82 { // 0x82 = PRG file
			continue
		}

		if verbose {
			fmt.Printf("Found PRG file: %s at offset $%08X\n", fileName, fileOffset)
			fmt.Printf("  Start: $%04X, End: $%04X\n", startAddr, endAddr)
		}

		// Read file data
		// In T64 format:
		// - startAddr is the load address where program will be placed
		// - endAddr is the last address of the loaded program (inclusive)
		// - File data = load address (2 bytes) + program data
		// - Program data size = endAddr - startAddr + 1 (because endAddr is inclusive)
		programDataSize := endAddr - startAddr + 1

		// Total file size is load address (2 bytes) + program data
		totalFileSize := programDataSize + 2

		if fileOffset+uint32(totalFileSize) > uint32(len(data)) {
			if verbose {
				fmt.Printf("  Error: file data extends beyond archive (need %d bytes)\n", totalFileSize)
			}
			continue
		}

		fullFileData := data[fileOffset : fileOffset+uint32(totalFileSize)]

		// Extract the actual program data (skip the first 2 bytes which are the load address)
		if len(fullFileData) < 2 {
			if verbose {
				fmt.Printf("  Error: file too small to contain load address\n")
			}
			continue
		}

		// Verify that the load address in the file matches the directory entry
		fileLoadAddr := uint16(fullFileData[0]) | (uint16(fullFileData[1]) << 8)
		if fileLoadAddr != startAddr {
			if verbose {
				fmt.Printf("  Warning: file load address $%04X doesn't match directory start address $%04X\n", fileLoadAddr, startAddr)
			}
		}

		fileData := fullFileData[2:] // Skip load address bytes

		if verbose {
			fmt.Printf("  Size: %d bytes\n", len(fileData))
		}

		segment := assembler.AssembledData{
			StartAddress: startAddr,
			Data:         utils.Value(bytes.NewBuffer(fileData)),
		}
		segments = append(segments, segment)
	}

	if len(segments) == 0 {
		return nil, fmt.Errorf("no PRG files found in T64 archive")
	}

	return segments, nil
}
