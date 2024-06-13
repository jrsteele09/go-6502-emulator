// Package memory provides types and functions for managing memory in the 6502 emulator.
package memory

// AddressSize defines an interface for address sizes that can be used in memory operations.
type AddressSize interface {
	~uint16 | ~uint32 | ~uint64
}

// Operations defines the interface for memory operations such as read and write.
type Operations[AZ AddressSize] interface {
	Write(address AZ, data ...byte)
	Read(address AZ) byte
}

// Memory represents a block of memory with a specific size.
type Memory[AZ AddressSize] struct {
	size  uint64
	bytes []byte
}

// NewMemory creates a new Memory instance with the specified size.
func NewMemory[AZ AddressSize](n uint64) *Memory[AZ] {
	return &Memory[AZ]{bytes: make([]byte, n), size: n}
}

// Write writes data to the specified address in memory.
func (m *Memory[AZ]) Write(address AZ, data ...byte) {
	for i := 0; i < len(data); i++ {
		m.bytes[(uint64(address)+uint64(i))%m.size] = data[i]
	}
}

// Read reads a byte from the specified address in memory.
func (m *Memory[AZ]) Read(address AZ) byte {
	return m.bytes[uint64(address)%m.size]
}
