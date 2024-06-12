package cpu

import (
	"fmt"

	"github.com/jrsteele09/go-6502-emulator/memory"
)

type Completed bool

type Cpu6502 interface {
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

const (
	resetVectorAddr  uint16 = 0xFFFC
	irqVector        uint16 = 0xFFFE
	nmiVector        uint16 = 0xFFFA
	stackPageAddress uint16 = 0x0100
)

type HaltExecution interface {
	Stop()
	Resume()
}

type Cpu struct {
	Reg               *Registers
	mem               memory.MemoryFunctions[uint16]
	opCodes           []*OpCodeDef
	cycles            uint64
	instructionCycles int
	instruction       InstructionFunc
	operands          []byte
	irq               bool
	nmi               bool
	halted            bool
}

var _ Cpu6502 = &Cpu{}

func NewCpu(m memory.MemoryFunctions[uint16]) *Cpu {
	cpu := &Cpu{mem: m, Reg: NewRegisters()}
	cpu.opCodes = OpCodes(cpu)
	cpu.Reg.SetStatus(UnusedFlag, true)
	cpu.Reg.S = 0xff
	cpu.irq = false
	cpu.nmi = false
	cpu.Reset()
	return cpu
}

func (p *Cpu) Registers() *Registers {
	return p.Reg
}

func (p *Cpu) Memory() memory.MemoryFunctions[uint16] {
	return p.mem
}

func (p *Cpu) Operands() []byte {
	return p.operands
}

func (p *Cpu) Stop() {
	p.halted = true
}

func (p *Cpu) Resume() {
	p.halted = false
}

func (p *Cpu) Execute() (Completed, error) {
	if p.halted {
		return false, nil
	}
	p.cycles++
	if p.instructionCycles > 0 {
		p.instructionCycles--
		return false, nil
	}
	completed, err := p.instruction()
	if completed {
		if p.checkInterrupts() {
			p.instruction = p.interruptInstruction
			p.instructionCycles = 7
		} else {
			p.instruction = p.readOpCode
			p.instructionCycles = 0
		}
	}
	return completed, err
}

func (p *Cpu) checkInterrupts() bool {
	if p.nmi {
		return true
	} else if p.irq && !p.Reg.IsSet(InterruptDisableFlag) {
		return true
	}
	return false
}

func (p *Cpu) readOpCode() (Completed, error) {
	opCode := p.NextByte()
	opCodeDef := p.opCodes[opCode]
	if opCodeDef == nil {
		return true, fmt.Errorf("unknown opCode: %x", opCode)
	}
	p.operands = make([]byte, opCodeDef.Bytes)
	p.instructionCycles = (opCodeDef.Cycles - 2) // Take two off for reading op code + next cycle

	for i := 0; i < opCodeDef.Bytes-1; i++ {
		p.operands[i] = p.NextByte()
	}
	p.instruction = opCodeDef.ExecGetter(*opCodeDef)
	return false, nil
}

func (p *Cpu) interruptInstruction() (Completed, error) {
	p.interruptStackPush()
	p.Reg.SetStatus(InterruptDisableFlag, true)
	var PCH byte
	var PCL byte
	if p.nmi {
		p.nmi = false
		p.irq = false
		PCL = p.mem.Read(uint16(nmiVector))
		PCH = p.mem.Read(uint16(nmiVector + 1))
	} else if p.irq {
		p.irq = false
		PCL = p.mem.Read(uint16(irqVector))
		PCH = p.mem.Read(uint16(irqVector + 1))
	}
	p.Reg.PC = (uint16(PCH) << 8) + uint16(PCL)
	return true, nil
}

func (p *Cpu) interruptStackPush() {
	p.Push(byte(p.Reg.PC >> 8))
	p.Push(byte(p.Reg.PC & 0xff))
	p.Reg.SetStatus(BreakFlag, false)
	p.Push(byte(p.Reg.Status))
	p.Reg.SetStatus(InterruptDisableFlag, true)
}

func (p *Cpu) NextByte() byte {
	b := p.mem.Read(uint16(p.Reg.PC))
	p.Reg.PC++
	return b
}

func (p *Cpu) Push(b byte) {
	a := stackPageAddress + uint16(p.Reg.S)
	p.mem.Write(uint16(a), b)
	p.Reg.S--
}

func (p *Cpu) Pop() byte {
	a := stackPageAddress + uint16(p.Reg.S)
	b := p.mem.Read(uint16(a))
	p.Reg.S++
	return b
}

func (p *Cpu) Nmi() {
	p.nmi = true
}

func (p *Cpu) Irq() {
	p.irq = !p.Reg.IsSet(InterruptDisableFlag)
}

func (p *Cpu) Reset() {
	p.Reg.SetStatus(InterruptDisableFlag, true)
	resetVecLow := p.mem.Read(uint16(resetVectorAddr))
	resetVecHigh := p.mem.Read(uint16(resetVectorAddr + 1))
	p.Reg.PC = (uint16(resetVecHigh) << 8) | uint16(resetVecLow)
	p.instruction = p.readOpCode
	p.instructionCycles = 0
}
