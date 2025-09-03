// Package cpu provides types and functions for emulating the 6502 CPU, including CPU operations and instruction execution.
package cpu

import (
	"fmt"

	"github.com/jrsteele09/go-6502-emulator/memory"
)

// Completed represents whether an instruction has completed execution.
type Completed bool

// CPU6502 interface defines the methods required to emulate the 6502 CPU.
type CPU6502 interface {
	Stop()
	Resume()
	Execute() (Completed, error)
	Nmi()
	Irq()
	Reset()
	Push(b byte)
	Pop() byte
	Registers() *Registers
	Memory() memory.Operations[uint16]
	Operands() []byte
	OpCodes() []*OpCodeDef
}

const (
	resetVectorAddr  uint16 = 0xFFFC
	irqVector        uint16 = 0xFFFE
	nmiVector        uint16 = 0xFFFA
	stackPageAddress uint16 = 0x0100
)

// HaltExecution interface defines methods to stop and resume CPU execution.
type HaltExecution interface {
	Stop()
	Resume()
}

// CPU represents the 6502 CPU with registers, memory, and opcode definitions.
type CPU struct {
	Reg               *Registers
	mem               memory.Operations[uint16]
	opCodes           []*OpCodeDef
	cycles            uint64
	instructionCycles int
	instruction       InstructionFunc
	operands          []byte
	irq               bool
	nmi               bool
	halted            bool
}

// Ensure Cpu implements the Cpu6502 interface.
var _ CPU6502 = &CPU{}

// NewCPU creates a new Cpu instance with the provided memory functions.
func NewCPU(m memory.Operations[uint16], useIllegalOpCodes bool) *CPU {
	cpu := &CPU{mem: m, Reg: NewRegisters()}
	cpu.opCodes = createOpCodes(cpu)

	if useIllegalOpCodes {
		// Attach undocumented/illegal opcodes
		addIllegalOpCodes(cpu)
	}

	cpu.Reg.SetStatus(UnusedFlag, true)
	cpu.Reg.S = 0xff
	cpu.irq = false
	cpu.nmi = false
	cpu.Reset()
	return cpu
}

func (p *CPU) OpCodes() []*OpCodeDef {
	return p.opCodes
}

// Registers returns the CPU's registers.
func (p *CPU) Registers() *Registers {
	return p.Reg
}

// Memory returns the memory functions used by the CPU.
func (p *CPU) Memory() memory.Operations[uint16] {
	return p.mem
}

// Operands returns the operands for the current instruction.
func (p *CPU) Operands() []byte {
	return p.operands
}

// Stop halts the CPU's execution.
func (p *CPU) Stop() {
	p.halted = true
}

// Resume resumes the CPU's execution.
func (p *CPU) Resume() {
	p.halted = false
}

// Execute executes the current instruction and returns whether it is completed and any error encountered.
func (p *CPU) Execute() (Completed, error) {
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

func (p *CPU) checkInterrupts() bool {
	if p.nmi {
		return true
	} else if p.irq && !p.Reg.IsSet(InterruptDisableFlag) {
		return true
	}
	return false
}

func (p *CPU) readOpCode() (Completed, error) {
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

func (p *CPU) interruptInstruction() (Completed, error) {
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

func (p *CPU) interruptStackPush() {
	p.Push(byte(p.Reg.PC >> 8))
	p.Push(byte(p.Reg.PC & 0xff))
	p.Reg.SetStatus(BreakFlag, false)
	p.Push(byte(p.Reg.Status))
	p.Reg.SetStatus(InterruptDisableFlag, true)
}

// NextByte reads the next byte from memory and increments the program counter.
func (p *CPU) NextByte() byte {
	b := p.mem.Read(uint16(p.Reg.PC))
	p.Reg.PC++
	return b
}

// Push pushes a byte onto the stack.
func (p *CPU) Push(b byte) {
	a := stackPageAddress + uint16(p.Reg.S)
	p.mem.Write(uint16(a), b)
	p.Reg.S--
}

// Pop pops a byte from the stack.
func (p *CPU) Pop() byte {
	a := stackPageAddress + uint16(p.Reg.S)
	b := p.mem.Read(uint16(a))
	p.Reg.S++
	return b
}

// Nmi triggers a non-maskable interrupt.
func (p *CPU) Nmi() {
	p.nmi = true
}

// Irq triggers an interrupt request.
func (p *CPU) Irq() {
	p.irq = !p.Reg.IsSet(InterruptDisableFlag)
}

// Reset resets the CPU to its initial state.
func (p *CPU) Reset() {
	p.Reg.SetStatus(InterruptDisableFlag, true)
	resetVecLow := p.mem.Read(uint16(resetVectorAddr))
	resetVecHigh := p.mem.Read(uint16(resetVectorAddr + 1))
	p.Reg.PC = (uint16(resetVecHigh) << 8) | uint16(resetVecLow)
	p.instruction = p.readOpCode
	p.instructionCycles = 0
}
