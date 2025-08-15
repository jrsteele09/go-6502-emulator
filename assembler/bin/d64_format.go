package bin

import (
	"bytes"
	"fmt"
	"os"

	"github.com/jrsteele09/go-6502-emulator/assembler"
	"github.com/jrsteele09/go-6502-emulator/utils"
)

var _ BinaryFormat = (*D64Format)(nil)

// D64Format implements the BinaryFormat interface for Commodore 64 disk images
// D64 format: Standard 1541 disk image format with 683 blocks (174,848 bytes)
type D64Format struct {
	// DiskName is the name of the disk (up to 16 characters)
	DiskName string
	// DiskID is a 2-character disk identifier
	DiskID string
	// StartTrack is the track where files start (default: 1)
	StartTrack uint8
	// StartSector is the sector where files start (default: 0)
	StartSector uint8
}

// D64 disk constants
const (
	D64TotalSize     = 174848 // Total size of D64 disk image
	D64BlockSize     = 256    // Size of each block
	D64TotalBlocks   = 683    // Total number of blocks
	D64DirectorySize = 8      // Number of directory blocks
	D64BAMSize       = 1      // Number of BAM (Block Allocation Map) blocks
	D64MaxFilename   = 16     // Maximum filename length
	D64MaxDiskName   = 16     // Maximum disk name length
	D64DiskIDLength  = 2      // Disk ID length
)

// NewD64Format creates a new D64 format generator
func NewD64Format(diskName, diskID string) *D64Format {
	if len(diskName) > D64MaxDiskName {
		diskName = diskName[:D64MaxDiskName]
	}
	if len(diskID) != D64DiskIDLength {
		diskID = "01" // Default disk ID
	}
	return &D64Format{
		DiskName:    diskName,
		DiskID:      diskID,
		StartTrack:  1,
		StartSector: 0,
	}
}

// CreateFile creates a D64 disk image file from assembled segments
func (d *D64Format) CreateFile(filename string, segments []assembler.AssembledData, verbose bool) error {
	data, err := d.CreateData(segments)
	if err != nil {
		return err
	}

	if verbose {
		fmt.Printf("D64 disk name: %s\n", d.DiskName)
		fmt.Printf("D64 disk ID: %s\n", d.DiskID)
		fmt.Printf("D64 image size: %d bytes\n", len(data))
		fmt.Printf("Program segments: %d\n", len(segments))
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create D64 file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write D64 data: %w", err)
	}

	return nil
}

// CreateData creates D64 format data in memory from assembled segments
func (d *D64Format) CreateData(segments []assembler.AssembledData) ([]byte, error) {
	if len(segments) == 0 {
		return nil, fmt.Errorf("no segments to convert to D64")
	}

	// Create empty D64 disk image
	diskImage := make([]byte, D64TotalSize)

	// Initialize BAM (Block Allocation Map) at track 18, sector 0
	d.initializeBAM(diskImage)

	// Initialize directory at track 18, sector 1
	d.initializeDirectory(diskImage)

	// Write program segments as files on the disk
	err := d.writeSegmentsToD64(diskImage, segments)
	if err != nil {
		return nil, err
	}

	return diskImage, nil
}

// initializeBAM initializes the Block Allocation Map
func (d *D64Format) initializeBAM(diskImage []byte) {
	// BAM is located at track 18, sector 0
	bamOffset := d.trackSectorToOffset(18, 0)

	// BAM header
	diskImage[bamOffset] = 18     // Track of first directory sector
	diskImage[bamOffset+1] = 1    // Sector of first directory sector
	diskImage[bamOffset+2] = 0x41 // DOS version ('A')

	// Initialize all tracks as free (simplified)
	for i := 4; i < 0x90; i++ {
		diskImage[bamOffset+i] = 0xFF // Mark all sectors as free
	}

	// Disk name (padded with spaces)
	copy(diskImage[bamOffset+0x90:bamOffset+0x90+D64MaxDiskName], d.padString(d.DiskName, D64MaxDiskName))

	// Disk ID
	copy(diskImage[bamOffset+0xA2:bamOffset+0xA2+D64DiskIDLength], d.DiskID)

	// DOS type
	diskImage[bamOffset+0xA5] = '2'
	diskImage[bamOffset+0xA6] = 'A'
}

// initializeDirectory initializes the directory
func (d *D64Format) initializeDirectory(diskImage []byte) {
	// Directory starts at track 18, sector 1
	dirOffset := d.trackSectorToOffset(18, 1)

	// Directory header (points to next directory sector, or 0,255 if last)
	diskImage[dirOffset] = 0     // No next directory track
	diskImage[dirOffset+1] = 255 // No next directory sector (end marker)
}

// writeSegmentsToD64 writes program segments as files to the D64 disk
func (d *D64Format) writeSegmentsToD64(diskImage []byte, segments []assembler.AssembledData) error {
	// For simplicity, we'll write all segments as one PRG file
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

	// Write file to disk starting at track 1, sector 0
	err := d.writeFileToD64(diskImage, prgData, "PROGRAM", 1, 0)
	if err != nil {
		return err
	}

	// Add directory entry
	d.addDirectoryEntry(diskImage, "PROGRAM", 1, 0, len(prgData))

	return nil
}

// writeFileToD64 writes file data to the D64 disk at specified track/sector
func (d *D64Format) writeFileToD64(diskImage []byte, data []byte, filename string, startTrack, startSector uint8) error {
	currentTrack := startTrack
	currentSector := startSector
	dataOffset := 0

	for dataOffset < len(data) {
		blockOffset := d.trackSectorToOffset(currentTrack, currentSector)

		// How much data can we write in this block?
		remainingData := len(data) - dataOffset
		blockDataSize := remainingData
		if blockDataSize > 254 { // Reserve 2 bytes for track/sector link
			blockDataSize = 254
		}

		// Write track/sector link (next block)
		if remainingData > 254 {
			// More data follows, link to next sector
			nextTrack, nextSector := d.getNextSector(currentTrack, currentSector)
			diskImage[blockOffset] = nextTrack
			diskImage[blockOffset+1] = nextSector
		} else {
			// Last block
			diskImage[blockOffset] = 0
			diskImage[blockOffset+1] = byte(blockDataSize + 1) // +1 because 0 means 256 bytes
		}

		// Write data
		copy(diskImage[blockOffset+2:blockOffset+2+blockDataSize], data[dataOffset:dataOffset+blockDataSize])

		dataOffset += blockDataSize
		currentTrack, currentSector = d.getNextSector(currentTrack, currentSector)
	}

	return nil
}

// addDirectoryEntry adds a file entry to the directory
func (d *D64Format) addDirectoryEntry(diskImage []byte, filename string, startTrack, startSector uint8, fileSize int) {
	// Find first available directory entry (starts at track 18, sector 1, offset 2)
	dirOffset := d.trackSectorToOffset(18, 1) + 2

	// Directory entry format (32 bytes per entry, 8 entries per sector)
	entryOffset := dirOffset // Use first entry for simplicity

	// File type (PRG = $82)
	diskImage[entryOffset] = 0x82

	// Track/sector of first file block
	diskImage[entryOffset+1] = startTrack
	diskImage[entryOffset+2] = startSector

	// Filename (16 bytes, padded with shifted spaces $A0)
	paddedFilename := d.padStringShifted(filename, D64MaxFilename)
	copy(diskImage[entryOffset+3:entryOffset+3+D64MaxFilename], paddedFilename)

	// File size in blocks (low byte, high byte)
	blocks := (fileSize + 253) / 254 // Round up to nearest block
	diskImage[entryOffset+28] = byte(blocks & 0xFF)
	diskImage[entryOffset+29] = byte(blocks >> 8)
}

// Helper functions

// trackSectorToOffset converts track/sector to byte offset in D64 image
func (d *D64Format) trackSectorToOffset(track, sector uint8) int {
	// Track 1-17: 21 sectors each
	// Track 18-24: 19 sectors each
	// Track 25-30: 18 sectors each
	// Track 31-35: 17 sectors each

	offset := 0
	for t := uint8(1); t < track; t++ {
		if t <= 17 {
			offset += 21 * D64BlockSize
		} else if t <= 24 {
			offset += 19 * D64BlockSize
		} else if t <= 30 {
			offset += 18 * D64BlockSize
		} else {
			offset += 17 * D64BlockSize
		}
	}

	offset += int(sector) * D64BlockSize
	return offset
}

// getNextSector returns the next available sector
func (d *D64Format) getNextSector(track, sector uint8) (uint8, uint8) {
	// Simple linear allocation (skip track 18 which is reserved for directory/BAM)
	nextSector := sector + 1
	nextTrack := track

	sectorsPerTrack := d.getSectorsPerTrack(track)
	if nextSector >= sectorsPerTrack {
		nextSector = 0
		nextTrack++
		if nextTrack == 18 {
			nextTrack = 19 // Skip directory track
		}
	}

	return nextTrack, nextSector
}

// getSectorsPerTrack returns the number of sectors for a given track
func (d *D64Format) getSectorsPerTrack(track uint8) uint8 {
	if track <= 17 {
		return 21
	} else if track <= 24 {
		return 19
	} else if track <= 30 {
		return 18
	} else {
		return 17
	}
}

// padString pads a string with spaces to specified length
func (d *D64Format) padString(s string, length int) []byte {
	result := make([]byte, length)
	copy(result, []byte(s))
	for i := len(s); i < length; i++ {
		result[i] = ' '
	}
	return result
}

// padStringShifted pads a string with shifted spaces ($A0) for PETSCII
func (d *D64Format) padStringShifted(s string, length int) []byte {
	result := make([]byte, length)
	copy(result, []byte(s))
	for i := len(s); i < length; i++ {
		result[i] = 0xA0 // Shifted space in PETSCII
	}
	return result
}

// LoadFile loads a D64 disk image and returns assembled segments from PRG files
func (d *D64Format) LoadFile(filename string, verbose bool) ([]assembler.AssembledData, error) {
	// Read the file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read D64 file: %w", err)
	}

	if len(data) != D64TotalSize {
		return nil, fmt.Errorf("invalid D64 file size: %d bytes, expected %d", len(data), D64TotalSize)
	}

	// Read disk information from BAM
	bamOffset := d.trackSectorToOffset(18, 0)
	diskName := string(data[bamOffset+0x90 : bamOffset+0x90+D64MaxDiskName])
	diskID := string(data[bamOffset+0xA2 : bamOffset+0xA2+D64DiskIDLength])

	// Update our format info
	d.DiskName = diskName
	d.DiskID = diskID

	if verbose {
		fmt.Printf("D64 disk name: %s\n", diskName)
		fmt.Printf("D64 disk ID: %s\n", diskID)
	}

	// Read directory to find PRG files
	segments := make([]assembler.AssembledData, 0)
	dirOffset := d.trackSectorToOffset(18, 1)

	// Process directory entries (8 entries per sector, starting at offset 2)
	for i := 0; i < 8; i++ {
		entryOffset := dirOffset + 2 + (i * 32)

		// Check if entry is used (file type != 0)
		fileType := data[entryOffset]
		if fileType == 0 {
			continue // Empty entry
		}

		// Check if it's a PRG file (type & 0x07 == 2, and not deleted)
		if (fileType&0x07) != 2 || (fileType&0x80) == 0 {
			continue // Not a PRG file or deleted
		}

		// Get file track/sector
		fileTrack := data[entryOffset+1]
		fileSector := data[entryOffset+2]

		// Get filename
		filename := string(data[entryOffset+3 : entryOffset+3+D64MaxFilename])
		// Remove padding (shifted spaces 0xA0)
		for j := len(filename) - 1; j >= 0 && (filename[j] == ' ' || byte(filename[j]) == 0xA0); j-- {
			filename = filename[:j]
		}

		if verbose {
			fmt.Printf("Found PRG file: %s at track %d, sector %d\n", filename, fileTrack, fileSector)
		}

		// Read the file data
		fileData, err := d.readFileFromD64(data, fileTrack, fileSector)
		if err != nil {
			if verbose {
				fmt.Printf("Error reading file %s: %v\n", filename, err)
			}
			continue
		}

		// Parse PRG data (first 2 bytes are load address)
		if len(fileData) >= 2 {
			loadAddr := uint16(fileData[0]) | (uint16(fileData[1]) << 8)
			programData := fileData[2:]

			if verbose {
				fmt.Printf("  Load address: $%04X, size: %d bytes\n", loadAddr, len(programData))
			}

			segment := assembler.AssembledData{
				StartAddress: loadAddr,
				Data:         utils.Value(bytes.NewBuffer(programData)),
			}
			segments = append(segments, segment)
		}
	}

	if len(segments) == 0 {
		return nil, fmt.Errorf("no PRG files found in D64 image")
	}

	return segments, nil
}

// readFileFromD64 reads a file from the D64 disk image following the block chain
func (d *D64Format) readFileFromD64(diskData []byte, startTrack, startSector uint8) ([]byte, error) {
	var fileData []byte
	currentTrack := startTrack
	currentSector := startSector

	// Follow the file block chain
	for {
		if currentTrack == 0 {
			break // End of file
		}

		blockOffset := d.trackSectorToOffset(currentTrack, currentSector)
		if blockOffset+256 > len(diskData) {
			return nil, fmt.Errorf("invalid track/sector: %d/%d", currentTrack, currentSector)
		}

		// Read next track/sector link
		nextTrack := diskData[blockOffset]
		nextSector := diskData[blockOffset+1]

		var dataSize int
		if nextTrack == 0 {
			// Last block, next sector indicates data size
			dataSize = int(nextSector) - 1 // -1 because 0 means 256 bytes
			if dataSize < 0 {
				dataSize = 255
			}
		} else {
			// Not last block, use full 254 bytes
			dataSize = 254
		}

		// Copy data from this block (excluding track/sector link)
		blockData := diskData[blockOffset+2 : blockOffset+2+dataSize]
		fileData = append(fileData, blockData...)

		// Move to next block
		currentTrack = nextTrack
		currentSector = nextSector
	}

	return fileData, nil
}
