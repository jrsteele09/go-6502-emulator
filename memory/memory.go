package memory

type AddressSize interface {
	~uint16 | ~uint32 | ~uint64
}

type MemoryFunctions[AZ AddressSize] interface {
	Write(address AZ, data ...byte)
	Read(address AZ) byte
}

type Memory[AZ AddressSize] struct {
	size  uint64
	bytes []byte
}

func NewMemory[AZ AddressSize](n uint64) *Memory[AZ] {
	return &Memory[AZ]{bytes: make([]byte, n), size: n}
}

func (m *Memory[AZ]) Write(address AZ, data ...byte) {
	for i := 0; i < len(data); i++ {
		m.bytes[(uint64(address)+uint64(i))%m.size] = data[i]
	}
}

func (m *Memory[AZ]) Read(address AZ) byte {
	return m.bytes[uint64(address)%m.size]
}
