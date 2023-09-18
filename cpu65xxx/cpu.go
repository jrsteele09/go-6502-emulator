package cpu65xxx

import "github.com/jrsteele09/go-65xx-emulator/memory"

type Completed bool

type Cpu interface {
	Stop()
	Resume()
	Execute() (Completed, error)
	Nmi()
	Irq()
	Reset()
	Push(b byte)
	Pop() byte
	Registers() *Registers
	Memory() memory.MemoryFunctions[uint16]
	Operands() []byte
}
