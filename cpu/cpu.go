// Package cpu provides types and functions for emulating the 6502 CPU, including CPU operations and instruction execution.
package cpu

import (
	"fmt"

	"github.com/jrsteele09/go-6502-emulator/memory"
)

// Completed represents whether an instruction has completed execution.
type Completed bool
type BinaryOrDecimalMode bool

const (
	BCDMode    BinaryOrDecimalMode = true
	BinaryMode                     = false
)

// CPU6502 interface defines the methods required to emulate the 6502 CPU.
type CPU6502 interface {
	Stop()
	Resume()
	Execute() (Completed, error)
	Nmi()
	Irq()
	SetIRQ(asserted bool)
	Reset()
	Push(b byte)
	Pop() byte
	Registers() *Registers
	Memory() memory.Operations[uint16]
	Operands() []byte
	OpCodes() []OpCodeDef
}

const (
	resetVectorAddr     uint16 = 0xFFFC
	irqVector           uint16 = 0xFFFE
	nmiVector           uint16 = 0xFFFA
	stackPageAddress    uint16 = 0x0100
	maxInstructionBytes        = 3
)

// HaltExecution interface defines methods to stop and resume CPU execution.
type HaltExecution interface {
	Stop()
	Resume()
}

// CPU represents the 6502 CPU with registers, memory, and opcode definitions.
type CPU struct {
	Reg    *Registers
	mem    memory.Operations[uint16]
	cycles uint64

	// Keep the mutable execution state together on the hot side of CPU. The
	// much larger, read-mostly opcode metadata table is deliberately last.
	execute executionState
	irq     bool
	nmi     bool
	halted  bool
	opCodes [256]OpCodeDef
}

type executionStage byte

const (
	executionFetch executionStage = iota
	executionInstruction
	executionInterrupt
)

// executionState contains all transient state for the instruction currently
// being executed. Selecting a new instruction replaces this value, ensuring
// scratch state cannot leak between instructions or survive Reset.
type executionState struct {
	opcode            *OpCodeDef
	instructionCycles int
	address           uint16

	// Everything below is instruction-local scratch space. selectInstruction,
	// interrupt selection, and Reset replace the whole structure.
	operands      [maxInstructionBytes - 1]byte
	stage         executionStage
	operandLength uint8
	phase         uint8
	pageCrossed   bool
}

// Ensure Cpu implements the Cpu6502 interface.
var _ CPU6502 = &CPU{}

// NewCPU creates a new Cpu instance with the provided memory functions.
func NewCPU(m memory.Operations[uint16], useIllegalOpCodes bool) *CPU {
	cpu := &CPU{mem: m, Reg: NewRegisters()}
	cpu.opCodes = createOpCodes()

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

func (p *CPU) OpCodes() []OpCodeDef {
	return p.opCodes[:]
}

// Registers returns the CPU's registers.
func (p *CPU) Registers() *Registers {
	return p.Reg
}

// Memory returns the memory functions used by the CPU.
func (p *CPU) Memory() memory.Operations[uint16] {
	return p.mem
}

// Operands returns a view of the operands for the current instruction.
// The returned slice is reused when the next opcode is decoded and must not be retained.
func (p *CPU) Operands() []byte {
	return p.execute.operands[:p.execute.operandLength]
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
	if p.execute.instructionCycles > 0 {
		p.execute.instructionCycles--
		return false, nil
	}
	completed, err := p.executionStage()
	if err != nil {
		return completed, err
	}
	if !completed {
		return false, nil
	}
	if p.checkInterrupts() {
		p.execute = executionState{stage: executionInterrupt, instructionCycles: 7}
	} else {
		p.execute = executionState{stage: executionFetch}
	}
	return true, nil
}

//go:inline
func (p *CPU) executionStage() (Completed, error) {
	var completed Completed
	var err error
	switch p.execute.stage {
	case executionFetch:
		completed, err = p.readOpCode()
	case executionInstruction:
		completed, err = p.execute.opcode.Execute(p, p.execute.opcode)
	case executionInterrupt:
		completed, err = p.interruptInstruction()
	}
	if err != nil {
		return completed, err
	}
	return completed, nil
}

//go:inline
func (p *CPU) checkInterrupts() bool {
	if p.nmi {
		return true
	} else if p.irq && !p.Reg.IsSet(InterruptDisableFlag) {
		return true
	}
	return false
}

//go:inline
func (p *CPU) readOpCode() (Completed, error) {
	opCode := p.NextByte()
	opCodeDef := &p.opCodes[opCode]
	if opCodeDef.Execute == nil {
		return true, fmt.Errorf("unknown opCode: %x", opCode)
	}
	p.setExecutionState(opCodeDef)

	for i := 0; i < opCodeDef.Bytes-1; i++ {
		p.execute.operands[i] = p.NextByte()
	}
	return false, nil
}

//go:inline
func (p *CPU) setExecutionState(opCodeDef *OpCodeDef) {
	p.execute = executionState{
		opcode:            opCodeDef,
		stage:             executionInstruction,
		instructionCycles: opCodeDef.Cycles - 2, // Subtract 2 cycles for the fetch and decode stages
		operandLength:     uint8(opCodeDef.Bytes - 1),
	}
}

//go:inline
func (p *CPU) interruptInstruction() (Completed, error) {
	p.interruptStackPush()
	p.Reg.SetStatus(InterruptDisableFlag, true)
	var PCH byte
	var PCL byte
	if p.nmi {
		p.nmi = false
		PCL = p.mem.Read(uint16(nmiVector))
		PCH = p.mem.Read(uint16(nmiVector + 1))
	} else if p.irq {
		PCL = p.mem.Read(uint16(irqVector))
		PCH = p.mem.Read(uint16(irqVector + 1))
	}
	p.Reg.PC = (uint16(PCH) << 8) + uint16(PCL)
	return true, nil
}

//go:inline
func (p *CPU) interruptStackPush() {
	p.Push(byte(p.Reg.PC >> 8))
	p.Push(byte(p.Reg.PC & 0xff))
	p.Reg.SetStatus(BreakFlag, false)
	p.Push(p.Reg.Status | byte(UnusedFlag))
}

// NextByte reads the next byte from memory and increments the program counter.
//
//go:inline
func (p *CPU) NextByte() byte {
	b := p.mem.Read(uint16(p.Reg.PC))
	p.Reg.PC++
	return b
}

// Push pushes a byte onto the stack.
//
//go:inline
func (p *CPU) Push(b byte) {
	a := stackPageAddress + uint16(p.Reg.S)
	p.mem.Write(uint16(a), b)
	p.Reg.S--
}

// Pop pops a byte from the stack.
//
//go:inline
func (p *CPU) Pop() byte {
	p.Reg.S++
	a := stackPageAddress + uint16(p.Reg.S)
	b := p.mem.Read(uint16(a))
	return b
}

// Nmi triggers a non-maskable interrupt.
//
//go:inline
func (p *CPU) Nmi() {
	p.nmi = true
}

// Irq triggers an interrupt request.
//
//go:inline
func (p *CPU) Irq() {
	p.SetIRQ(true)
}

// SetIRQ sets the level of the maskable interrupt input. The line remains at
// this level until the caller changes it; the interrupt-disable flag only
// controls whether an asserted line can be serviced.
func (p *CPU) SetIRQ(asserted bool) {
	p.irq = asserted
}

// Reset resets the CPU to its initial state.
//
//go:inline
func (p *CPU) Reset() {
	p.Reg.SetStatus(InterruptDisableFlag, true)
	resetVecLow := p.mem.Read(uint16(resetVectorAddr))
	resetVecHigh := p.mem.Read(uint16(resetVectorAddr + 1))
	p.Reg.PC = (uint16(resetVecHigh) << 8) | uint16(resetVecLow)
	p.execute = executionState{stage: executionFetch}
}
