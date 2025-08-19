package debugger

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jrsteele09/go-6502-emulator/assembler/bin"
	"github.com/jrsteele09/go-6502-emulator/cpu"
	"github.com/jrsteele09/go-6502-emulator/memory"
)

// Debugger represents the 6502 debugger core functionality
type Debugger struct {
	cpu            cpu.CPU6502
	memory         *memory.Memory[uint16]
	disassembler   *Disassembler
	breakpoints    map[uint16]bool
	running        bool
	lastDisasmAddr uint16
}

// NewDebugger creates a new 6502 debugger instance
func NewDebugger() *Debugger {
	mem := memory.NewMemory[uint16](64 * 1024) // 64KB memory
	cpuInstance := cpu.NewCPU(mem)
	opcodes := cpuInstance.OpCodes()
	disasm := NewDisassembler(mem, opcodes)

	return &Debugger{
		cpu:            cpuInstance,
		memory:         mem,
		disassembler:   disasm,
		breakpoints:    make(map[uint16]bool),
		running:        false,
		lastDisasmAddr: 0,
	}
}

// GetCPU returns the CPU instance for external access
func (d *Debugger) GetCPU() cpu.CPU6502 {
	return d.cpu
}

// GetMemory returns the memory instance for external access
func (d *Debugger) GetMemory() *memory.Memory[uint16] {
	return d.memory
}

// GetBreakpoints returns the breakpoints map for external access
func (d *Debugger) GetBreakpoints() map[uint16]bool {
	return d.breakpoints
}

// IsRunning returns whether the debugger is currently running
func (d *Debugger) IsRunning() bool {
	return d.running
}

// SetRunning sets the running state
func (d *Debugger) SetRunning(running bool) {
	d.running = running
}

// GetLastDisasmAddr returns the last disassembly address
func (d *Debugger) GetLastDisasmAddr() uint16 {
	return d.lastDisasmAddr
}

// SetLastDisasmAddr sets the last disassembly address
func (d *Debugger) SetLastDisasmAddr(addr uint16) {
	d.lastDisasmAddr = addr
}

// parseAddress parses an address string (hex or decimal)
func (d *Debugger) ParseAddress(addr string) (uint16, error) {
	if addr == "" {
		return 0, fmt.Errorf("address required")
	}

	if strings.HasPrefix(addr, "$") {
		// Hex address
		val, err := strconv.ParseUint(addr[1:], 16, 16)
		return uint16(val), err
	} else {
		// Decimal address
		val, err := strconv.ParseUint(addr, 10, 16)
		return uint16(val), err
	}
}

// parseValue parses a value string (hex or decimal)
func (d *Debugger) ParseValue(val string) (uint8, error) {
	if val == "" {
		return 0, fmt.Errorf("value required")
	}

	if strings.HasPrefix(val, "$") {
		// Hex value
		v, err := strconv.ParseUint(val[1:], 16, 8)
		return uint8(v), err
	} else {
		// Decimal value
		v, err := strconv.ParseUint(val, 10, 8)
		return uint8(v), err
	}
}

// formatAddress formats an address for display
func (d *Debugger) FormatAddress(addr uint16) string {
	return fmt.Sprintf("$%04X", addr)
}

// formatByte formats a byte for display
func (d *Debugger) FormatByte(b uint8) string {
	return fmt.Sprintf("$%02X", b)
}

// ShowRegisters displays the CPU registers
func (d *Debugger) ShowRegisters() string {
	regs := d.cpu.Registers()

	result := "Registers:\n"
	result += fmt.Sprintf("  A: %s  X: %s  Y: %s  PC: %s  S: %s\n",
		d.FormatByte(regs.A),
		d.FormatByte(regs.X),
		d.FormatByte(regs.Y),
		d.FormatAddress(regs.PC),
		d.FormatByte(regs.S))

	// Status flags
	status := regs.Status

	// Convert to binary string
	binary := fmt.Sprintf("%08b", status)

	// Create flag representation
	flags := ""
	if status&byte(cpu.NegativeFlag) != 0 {
		flags += "N"
	} else {
		flags += "."
	}
	if status&byte(cpu.OverflowFlag) != 0 {
		flags += "V"
	} else {
		flags += "."
	}
	// Unused flag (bit 5) - always set on 6502
	flags += "1"
	if status&byte(cpu.BreakFlag) != 0 {
		flags += "B"
	} else {
		flags += "."
	}
	if status&byte(cpu.DecimalFlag) != 0 {
		flags += "D"
	} else {
		flags += "."
	}
	if status&byte(cpu.InterruptDisableFlag) != 0 {
		flags += "I"
	} else {
		flags += "."
	}
	if status&byte(cpu.ZeroFlag) != 0 {
		flags += "Z"
	} else {
		flags += "."
	}
	if status&byte(cpu.CarryFlag) != 0 {
		flags += "C"
	} else {
		flags += "."
	}

	result += fmt.Sprintf("  Flags: %s (%%%s) (%s)  NV1BDIZC\n",
		d.FormatByte(status), binary, flags)

	return result
}

// Disassemble shows disassembled instructions
func (d *Debugger) Disassemble(args []string) string {
	var startAddr uint16
	var count int = 10

	if len(args) > 0 {
		if addr, err := d.ParseAddress(args[0]); err == nil {
			startAddr = addr
		} else {
			return fmt.Sprintf("Error parsing address: %v\n", err)
		}
	} else {
		// No address specified, continue from where we left off
		if d.lastDisasmAddr == 0 {
			startAddr = d.cpu.Registers().PC
		} else {
			startAddr = d.lastDisasmAddr
		}
	}

	if len(args) > 1 {
		if c, err := strconv.Atoi(args[1]); err == nil && c > 0 {
			count = c
		}
	}

	result := fmt.Sprintf("Disassembly from %s:\n\n", d.FormatAddress(startAddr))

	addr := startAddr
	for i := 0; i < count; i++ {
		// Check if this is the current PC
		marker := " "
		if addr == d.cpu.Registers().PC {
			marker = ">"
		}

		// Check if there's a breakpoint here
		if d.breakpoints[addr] {
			marker = "*"
		}

		instruction, length := d.disassembler.Disassemble(addr)
		result += fmt.Sprintf("%-1s %s\n", marker, instruction)
		addr += uint16(length)
	}

	// Remember where we left off for next time
	d.lastDisasmAddr = addr

	return result
}

// HexDump shows a hex dump of memory
func (d *Debugger) HexDump(args []string) string {
	var startAddr uint16 = 0x0000
	var count int = 16

	if len(args) > 0 {
		if addr, err := d.ParseAddress(args[0]); err == nil {
			startAddr = addr
		} else {
			return fmt.Sprintf("Error parsing address: %v\n", err)
		}
	}

	if len(args) > 1 {
		if c, err := strconv.Atoi(args[1]); err == nil && c > 0 {
			count = c
		}
	}

	result := fmt.Sprintf("Memory dump from %s:\n\n", d.FormatAddress(startAddr))

	// Round down to 16-byte boundary for nice display
	displayStart := startAddr & 0xFFF0

	for i := 0; i < ((count+15)/16)*16; i += 16 {
		addr := displayStart + uint16(i)

		// Address column
		result += fmt.Sprintf("%04X: ", addr)

		// Hex bytes
		hex := ""
		ascii := ""
		for j := 0; j < 16; j++ {
			byteAddr := addr + uint16(j)
			b := d.memory.Read(byteAddr)

			if byteAddr >= startAddr && byteAddr < startAddr+uint16(count) {
				hex += fmt.Sprintf("%02X ", b)
			} else {
				hex += fmt.Sprintf("%02X ", b)
			}

			// ASCII representation
			if b >= 32 && b <= 126 {
				ascii += string(b)
			} else {
				ascii += "."
			}
		}

		result += fmt.Sprintf("%-48s %s\n", hex, ascii)

		if addr >= startAddr+uint16(count) {
			break
		}
	}

	return result
}

// LoadPRG loads a PRG file into memory
func (d *Debugger) LoadPRG(filename string) string {
	// Use the PRG format loader
	prgFormat := bin.NewPRGFormat()
	segments, err := prgFormat.LoadFile(filename, false)
	if err != nil {
		return fmt.Sprintf("Error loading PRG file: %v\n", err)
	}

	if len(segments) == 0 {
		return "No segments found in PRG file\n"
	}

	result := fmt.Sprintf("Loaded PRG file: %s\n", filename)

	totalBytes := 0
	for i, segment := range segments {
		// Load segment into memory
		for j, b := range segment.Data.Bytes() {
			d.memory.Write(segment.StartAddress+uint16(j), b)
		}

		totalBytes += len(segment.Data.Bytes())
		result += fmt.Sprintf("  Segment %d: %s to %s (%d bytes)\n",
			i+1,
			d.FormatAddress(segment.StartAddress),
			d.FormatAddress(segment.StartAddress+uint16(len(segment.Data.Bytes()))-1),
			len(segment.Data.Bytes()))
	}

	result += fmt.Sprintf("Total: %d bytes loaded\n", totalBytes)

	// Set PC to first segment's start address
	if len(segments) > 0 {
		d.cpu.Registers().PC = segments[0].StartAddress
		d.lastDisasmAddr = segments[0].StartAddress // Reset disassembly position to PC
		result += fmt.Sprintf("PC set to %s\n", d.FormatAddress(segments[0].StartAddress))
	}

	return result
}

// Step executes one or more instructions
func (d *Debugger) Step(args []string) string {
	count := 1
	if len(args) > 0 {
		if c, err := strconv.Atoi(args[0]); err == nil && c > 0 {
			count = c
		}
	}

	result := ""
	for i := 0; i < count; i++ {
		pc := d.cpu.Registers().PC
		instruction, _ := d.disassembler.Disassemble(pc)

		// Only show step number if stepping multiple instructions
		if count > 1 {
			result += fmt.Sprintf("Step %d: %s\n", i+1, instruction)
		} else {
			result += fmt.Sprintf("Executing: %s\n", instruction)
		}

		completed := cpu.Completed(false)
		var err error
		for !completed {
			completed, err = d.cpu.Execute()
			if err != nil {
				result += fmt.Sprintf("Execution error: %v\n", err)
				break
			}
		}

		if !bool(completed) {
			result += "Instruction not completed\n"
		}

		// Show registers after step
		result += d.ShowRegisters()

		// Check for breakpoints
		newPC := d.cpu.Registers().PC
		if d.breakpoints[newPC] && i < count-1 {
			result += fmt.Sprintf("Breakpoint hit at %s\n", d.FormatAddress(newPC))
			break
		}
	}

	return result
}

// Go runs the program from the specified address
func (d *Debugger) Go(args []string) string {
	if len(args) > 0 {
		if addr, err := d.ParseAddress(args[0]); err == nil {
			d.cpu.Registers().PC = addr
		} else {
			return fmt.Sprintf("Error parsing address: %v\n", err)
		}
	}

	startPC := d.cpu.Registers().PC
	result := fmt.Sprintf("Running from %s... (Ctrl+C to break)\n", d.FormatAddress(startPC))

	d.running = true
	instructionCount := 0

	for d.running {
		pc := d.cpu.Registers().PC

		// Check for breakpoint
		if d.breakpoints[pc] {
			result += fmt.Sprintf("\nBreakpoint hit at %s\n", d.FormatAddress(pc))
			instruction, _ := d.disassembler.Disassemble(pc)
			result += fmt.Sprintf("Next: %s\n", instruction)
			d.running = false
			break
		}

		completed, err := d.cpu.Execute()
		if err != nil {
			result += fmt.Sprintf("\nExecution error at %s: %v\n", d.FormatAddress(pc), err)
			d.running = false
			break
		}

		if !bool(completed) {
			// Instruction needs more cycles, continue
			continue
		}

		instructionCount++

		// Check for infinite loops or runaway execution
		if instructionCount%100000 == 0 {
			result += fmt.Sprintf("\nExecuted %d instructions. Press Ctrl+C to break.\n", instructionCount)
		}
	}

	if d.running {
		d.running = false
		result += "\nExecution stopped\n"
	}

	return result
}

// SetBreakpoint sets a breakpoint at the specified address
func (d *Debugger) SetBreakpoint(args []string) string {
	if len(args) == 0 {
		return "Usage: B <address>\n"
	}

	addr, err := d.ParseAddress(args[0])
	if err != nil {
		return fmt.Sprintf("Error parsing address: %v\n", err)
	}

	d.breakpoints[addr] = true
	return fmt.Sprintf("Breakpoint set at %s\n", d.FormatAddress(addr))
}

// ClearBreakpoint clears a breakpoint at the specified address
func (d *Debugger) ClearBreakpoint(args []string) string {
	if len(args) == 0 {
		return "Usage: C <address>\n"
	}

	addr, err := d.ParseAddress(args[0])
	if err != nil {
		return fmt.Sprintf("Error parsing address: %v\n", err)
	}

	if d.breakpoints[addr] {
		delete(d.breakpoints, addr)
		return fmt.Sprintf("Breakpoint cleared at %s\n", d.FormatAddress(addr))
	} else {
		return fmt.Sprintf("No breakpoint at %s\n", d.FormatAddress(addr))
	}
}

// ListBreakpoints lists all active breakpoints
func (d *Debugger) ListBreakpoints() string {
	if len(d.breakpoints) == 0 {
		return "No breakpoints set\n"
	}

	result := "Active Breakpoints:\n"
	for addr := range d.breakpoints {
		result += fmt.Sprintf("  %s\n", d.FormatAddress(addr))
	}

	return result
}

// ZeroMemory zeros a range of memory
func (d *Debugger) ZeroMemory(args []string) string {
	if len(args) < 2 {
		return "Usage: Z <start_address> <end_address>\n"
	}

	startAddr, err := d.ParseAddress(args[0])
	if err != nil {
		return fmt.Sprintf("Error parsing start address: %v\n", err)
	}

	endAddr, err := d.ParseAddress(args[1])
	if err != nil {
		return fmt.Sprintf("Error parsing end address: %v\n", err)
	}

	if endAddr < startAddr {
		return "End address must be >= start address\n"
	}

	count := int(endAddr - startAddr + 1)
	for addr := startAddr; addr <= endAddr; addr++ {
		d.memory.Write(addr, 0x00)
	}

	return fmt.Sprintf("Zeroed %d bytes from %s to %s\n",
		count, d.FormatAddress(startAddr), d.FormatAddress(endAddr))
}

// FillMemory fills a range of memory with a value
func (d *Debugger) FillMemory(args []string) string {
	if len(args) < 3 {
		return "Usage: F <start_address> <end_address> <value>\n"
	}

	startAddr, err := d.ParseAddress(args[0])
	if err != nil {
		return fmt.Sprintf("Error parsing start address: %v\n", err)
	}

	endAddr, err := d.ParseAddress(args[1])
	if err != nil {
		return fmt.Sprintf("Error parsing end address: %v\n", err)
	}

	value, err := d.ParseValue(args[2])
	if err != nil {
		return fmt.Sprintf("Error parsing value: %v\n", err)
	}

	if endAddr < startAddr {
		return "End address must be >= start address\n"
	}

	count := int(endAddr - startAddr + 1)
	for addr := startAddr; addr <= endAddr; addr++ {
		d.memory.Write(addr, value)
	}

	return fmt.Sprintf("Filled %d bytes from %s to %s with %s\n",
		count, d.FormatAddress(startAddr), d.FormatAddress(endAddr), d.FormatByte(value))
}

// TransferMemory transfers a block of memory
func (d *Debugger) TransferMemory(args []string) string {
	if len(args) < 3 {
		return "Usage: T <source_address> <dest_address> <length>\n"
	}

	srcAddr, err := d.ParseAddress(args[0])
	if err != nil {
		return fmt.Sprintf("Error parsing source address: %v\n", err)
	}

	destAddr, err := d.ParseAddress(args[1])
	if err != nil {
		return fmt.Sprintf("Error parsing destination address: %v\n", err)
	}

	length, err := strconv.ParseUint(args[2], 10, 16)
	if err != nil {
		return fmt.Sprintf("Error parsing length: %v\n", err)
	}

	// Copy the memory block
	for i := uint16(0); i < uint16(length); i++ {
		value := d.memory.Read(srcAddr + i)
		d.memory.Write(destAddr+i, value)
	}

	return fmt.Sprintf("Transferred %d bytes from %s to %s\n",
		length, d.FormatAddress(srcAddr), d.FormatAddress(destAddr))
}
