package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jrsteele09/go-6502-emulator/assembler/bin"
	"github.com/jrsteele09/go-6502-emulator/cpu"
	"github.com/jrsteele09/go-6502-emulator/disassembler"
	"github.com/jrsteele09/go-6502-emulator/memory"
)

// ANSI color codes
const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	Gray    = "\033[37m"
	Bold    = "\033[1m"
	Dim     = "\033[2m"
)

// Debugger represents the 6502 debugger REPL
type Debugger struct {
	cpu          cpu.CPU6502
	memory       *memory.Memory[uint16]
	disassembler *disassembler.Disassembler
	breakpoints  map[uint16]bool
	running      bool
	lastCommand  string
	scanner      *bufio.Scanner
}

// NewDebugger creates a new 6502 debugger instance
func NewDebugger() *Debugger {
	mem := memory.NewMemory[uint16](64 * 1024) // 64KB memory
	cpuInstance := cpu.NewCPU(mem)
	opcodes := cpu.OpCodes(cpuInstance)
	disasm := disassembler.NewDisassembler(mem, opcodes)

	return &Debugger{
		cpu:          cpuInstance,
		memory:       mem,
		disassembler: disasm,
		breakpoints:  make(map[uint16]bool),
		running:      false,
		scanner:      bufio.NewScanner(os.Stdin),
	}
}

func main() {
	debugger := NewDebugger()
	debugger.Run()
}

// Run starts the debugger REPL
func (d *Debugger) Run() {
	d.printBanner()
	d.showHelp()

	for {
		d.printPrompt()

		if !d.scanner.Scan() {
			break // EOF or error
		}

		input := strings.TrimSpace(d.scanner.Text())
		if input == "" {
			input = d.lastCommand // Repeat last command
		} else {
			d.lastCommand = input
		}

		if input == "" {
			continue
		}

		if d.processCommand(input) {
			break // Exit requested
		}
	}

	fmt.Printf("%sGoodbye!%s\n", Yellow, Reset)
}

// printBanner displays the debugger banner
func (d *Debugger) printBanner() {
	fmt.Printf("%s%s", Bold, Cyan)
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    6502 Debugger v1.0.0                     ║")
	fmt.Println("║              Commodore 128 Monitor Style REPL               ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Printf("%s", Reset)
	fmt.Println()
}

// printPrompt displays the debugger prompt
func (d *Debugger) printPrompt() {
	pc := d.cpu.Registers().PC
	fmt.Printf("%s.%s %04X%s> ", Bold, Green, pc, Reset)
}

// processCommand processes a user command and returns true if exit is requested
func (d *Debugger) processCommand(input string) bool {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return false
	}

	command := strings.ToUpper(parts[0])
	args := parts[1:]

	switch command {
	case "H", "HELP", "?":
		d.showHelp()
	case "Q", "QUIT", "EXIT":
		return true
	case "D", "DISASSEMBLE":
		d.disassemble(args)
	case "M", "MEMORY":
		d.hexDump(args)
	case "R", "REGISTERS":
		d.showRegisters()
	case "L", "LOAD":
		d.loadPRG(args)
	case "G", "GO":
		d.go_(args)
	case "S", "STEP":
		d.step(args)
	case "B", "BREAK":
		d.setBreakpoint(args)
	case "BR", "BREAKPOINTS":
		d.listBreakpoints()
	case "C", "CLEAR":
		d.clearBreakpoint(args)
	case "Z", "ZERO":
		d.zeroMemory(args)
	case "F", "FILL":
		d.fillMemory(args)
	case "T", "TRANSFER":
		d.transferMemory(args)
	default:
		fmt.Printf("%sUnknown command: %s. Type 'H' for help.%s\n", Red, command, Reset)
	}

	return false
}

// showHelp displays available commands
func (d *Debugger) showHelp() {
	fmt.Printf("%s%sAvailable Commands:%s\n", Bold, Yellow, Reset)
	fmt.Printf("%s  H, HELP, ?%s        - Show this help\n", Cyan, Reset)
	fmt.Printf("%s  Q, QUIT, EXIT%s     - Exit debugger\n", Cyan, Reset)
	fmt.Printf("%s  R, REGISTERS%s      - Show CPU registers\n", Cyan, Reset)
	fmt.Printf("%s  D [addr] [count]%s  - Disassemble memory (default: PC, 10 instructions)\n", Cyan, Reset)
	fmt.Printf("%s  M [addr] [count]%s  - Memory hex dump (default: $0000, 16 bytes)\n", Cyan, Reset)
	fmt.Printf("%s  L <filename>%s      - Load PRG file into memory\n", Cyan, Reset)
	fmt.Printf("%s  G [addr]%s          - Go/Run from address (default: current PC)\n", Cyan, Reset)
	fmt.Printf("%s  S [count]%s         - Step instruction(s) (default: 1)\n", Cyan, Reset)
	fmt.Printf("%s  B <addr>%s          - Set breakpoint at address\n", Cyan, Reset)
	fmt.Printf("%s  BR, BREAKPOINTS%s   - List all breakpoints\n", Cyan, Reset)
	fmt.Printf("%s  C <addr>%s          - Clear breakpoint at address\n", Cyan, Reset)
	fmt.Printf("%s  Z <start> <end>%s   - Zero memory range\n", Cyan, Reset)
	fmt.Printf("%s  F <start> <end> <val>%s - Fill memory range with value\n", Cyan, Reset)
	fmt.Printf("%s  T <src> <dest> <len>%s - Transfer memory block\n", Cyan, Reset)
	fmt.Println()
	fmt.Printf("%s%sNotes:%s\n", Bold, Yellow, Reset)
	fmt.Printf("  - Addresses can be in hex ($1000) or decimal (4096)\n")
	fmt.Printf("  - Press Enter to repeat last command\n")
	fmt.Printf("  - Use Ctrl+C to break running program\n")
	fmt.Println()
}

// parseAddress parses an address string (hex or decimal)
func (d *Debugger) parseAddress(addr string) (uint16, error) {
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
func (d *Debugger) parseValue(val string) (uint8, error) {
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
func (d *Debugger) formatAddress(addr uint16) string {
	return fmt.Sprintf("$%04X", addr)
}

// formatByte formats a byte for display
func (d *Debugger) formatByte(b uint8) string {
	return fmt.Sprintf("$%02X", b)
}

// showRegisters displays the CPU registers
func (d *Debugger) showRegisters() {
	regs := d.cpu.Registers()

	fmt.Printf("%s%sRegisters:%s\n", Bold, Yellow, Reset)
	fmt.Printf("  %sA:%s %s  %sX:%s %s  %sY:%s %s  %sPC:%s %s  %sS:%s %s\n",
		Cyan, Reset, d.formatByte(regs.A),
		Cyan, Reset, d.formatByte(regs.X),
		Cyan, Reset, d.formatByte(regs.Y),
		Cyan, Reset, d.formatAddress(regs.PC),
		Cyan, Reset, d.formatByte(regs.S))

	// Status flags
	status := regs.Status
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
	flags += "." // Unused flag
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

	fmt.Printf("  %sFlags:%s %s (%s)  NV.BDIZC\n",
		Cyan, Reset, d.formatByte(status), flags)
	fmt.Println()
}

// disassemble shows disassembled instructions
func (d *Debugger) disassemble(args []string) {
	var startAddr uint16 = d.cpu.Registers().PC
	var count int = 10

	if len(args) > 0 {
		if addr, err := d.parseAddress(args[0]); err == nil {
			startAddr = addr
		} else {
			fmt.Printf("%sError parsing address: %v%s\n", Red, err, Reset)
			return
		}
	}

	if len(args) > 1 {
		if c, err := strconv.Atoi(args[1]); err == nil && c > 0 {
			count = c
		}
	}

	fmt.Printf("%s%sDisassembly from %s:%s\n", Bold, Yellow, d.formatAddress(startAddr), Reset)
	fmt.Println()

	addr := startAddr
	for i := 0; i < count; i++ {
		// Check if this is the current PC
		prefix := "  "
		if addr == d.cpu.Registers().PC {
			prefix = fmt.Sprintf("%s>%s", Green, Reset)
		}

		// Check if there's a breakpoint here
		if d.breakpoints[addr] {
			prefix = fmt.Sprintf("%s*%s", Red, Reset)
		}

		instruction, length := d.disassembler.Disassemble(addr)
		fmt.Printf("%s %s\n", prefix, instruction)
		addr += uint16(length)
	}
	fmt.Println()
}

// hexDump shows a hex dump of memory
func (d *Debugger) hexDump(args []string) {
	var startAddr uint16 = 0x0000
	var count int = 16

	if len(args) > 0 {
		if addr, err := d.parseAddress(args[0]); err == nil {
			startAddr = addr
		} else {
			fmt.Printf("%sError parsing address: %v%s\n", Red, err, Reset)
			return
		}
	}

	if len(args) > 1 {
		if c, err := strconv.Atoi(args[1]); err == nil && c > 0 {
			count = c
		}
	}

	fmt.Printf("%s%sMemory dump from %s:%s\n", Bold, Yellow, d.formatAddress(startAddr), Reset)
	fmt.Println()

	// Round down to 16-byte boundary for nice display
	displayStart := startAddr & 0xFFF0

	for i := 0; i < ((count+15)/16)*16; i += 16 {
		addr := displayStart + uint16(i)

		// Address column
		fmt.Printf("%s%04X:%s ", Cyan, addr, Reset)

		// Hex bytes
		hex := ""
		ascii := ""
		for j := 0; j < 16; j++ {
			byteAddr := addr + uint16(j)
			b := d.memory.Read(byteAddr)

			if byteAddr >= startAddr && byteAddr < startAddr+uint16(count) {
				hex += fmt.Sprintf("%02X ", b)
			} else {
				hex += fmt.Sprintf("%s%02X%s ", Dim, b, Reset)
			}

			// ASCII representation
			if b >= 32 && b <= 126 {
				ascii += string(b)
			} else {
				ascii += "."
			}
		}

		fmt.Printf("%-48s %s%s%s\n", hex, Gray, ascii, Reset)

		if addr >= startAddr+uint16(count) {
			break
		}
	}
	fmt.Println()
}

// loadPRG loads a PRG file into memory
func (d *Debugger) loadPRG(args []string) {
	if len(args) == 0 {
		fmt.Printf("%sUsage: L <filename>%s\n", Red, Reset)
		return
	}

	filename := args[0]

	// Use the PRG format loader
	prgFormat := bin.NewPRGFormat()
	segments, err := prgFormat.LoadFile(filename, false)
	if err != nil {
		fmt.Printf("%sError loading PRG file: %v%s\n", Red, err, Reset)
		return
	}

	if len(segments) == 0 {
		fmt.Printf("%sNo segments found in PRG file%s\n", Red, Reset)
		return
	}

	fmt.Printf("%s%sLoaded PRG file: %s%s\n", Bold, Green, filename, Reset)

	totalBytes := 0
	for i, segment := range segments {
		// Load segment into memory
		for j, b := range segment.Data {
			d.memory.Write(segment.StartAddress+uint16(j), b)
		}

		totalBytes += len(segment.Data)
		fmt.Printf("  %sSegment %d:%s %s to %s (%d bytes)\n",
			Cyan, i+1, Reset,
			d.formatAddress(segment.StartAddress),
			d.formatAddress(segment.StartAddress+uint16(len(segment.Data))-1),
			len(segment.Data))
	}

	fmt.Printf("%sTotal: %d bytes loaded%s\n", Green, totalBytes, Reset)

	// Set PC to first segment's start address
	if len(segments) > 0 {
		d.cpu.Registers().PC = segments[0].StartAddress
		fmt.Printf("%sPC set to %s%s\n", Green, d.formatAddress(segments[0].StartAddress), Reset)
	}
	fmt.Println()
}

// step executes one or more instructions
func (d *Debugger) step(args []string) {
	count := 1
	if len(args) > 0 {
		if c, err := strconv.Atoi(args[0]); err == nil && c > 0 {
			count = c
		}
	}

	for i := 0; i < count; i++ {
		pc := d.cpu.Registers().PC
		instruction, _ := d.disassembler.Disassemble(pc)

		fmt.Printf("%sStep %d:%s %s\n", Yellow, i+1, Reset, instruction)

		completed, err := d.cpu.Execute()
		if err != nil {
			fmt.Printf("%sExecution error: %v%s\n", Red, err, Reset)
			break
		}

		if !bool(completed) {
			fmt.Printf("%sInstruction not completed%s\n", Yellow, Reset)
		}

		// Show registers after step
		d.showRegisters()

		// Check for breakpoints
		newPC := d.cpu.Registers().PC
		if d.breakpoints[newPC] && i < count-1 {
			fmt.Printf("%sBreakpoint hit at %s%s\n", Red, d.formatAddress(newPC), Reset)
			break
		}
	}
}

// go_ runs the program from the specified address
func (d *Debugger) go_(args []string) {
	if len(args) > 0 {
		if addr, err := d.parseAddress(args[0]); err == nil {
			d.cpu.Registers().PC = addr
		} else {
			fmt.Printf("%sError parsing address: %v%s\n", Red, err, Reset)
			return
		}
	}

	startPC := d.cpu.Registers().PC
	fmt.Printf("%sRunning from %s... (Ctrl+C to break)%s\n", Green, d.formatAddress(startPC), Reset)

	d.running = true
	instructionCount := 0

	for d.running {
		pc := d.cpu.Registers().PC

		// Check for breakpoint
		if d.breakpoints[pc] {
			fmt.Printf("\n%sBreakpoint hit at %s%s\n", Red, d.formatAddress(pc), Reset)
			instruction, _ := d.disassembler.Disassemble(pc)
			fmt.Printf("Next: %s\n", instruction)
			d.running = false
			break
		}

		completed, err := d.cpu.Execute()
		if err != nil {
			fmt.Printf("\n%sExecution error at %s: %v%s\n", Red, d.formatAddress(pc), err, Reset)
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
			fmt.Printf("\n%sExecuted %d instructions. Press Ctrl+C to break.%s\n", Yellow, instructionCount, Reset)
		}
	}

	if d.running {
		d.running = false
		fmt.Printf("\n%sExecution stopped%s\n", Yellow, Reset)
	}
}

// setBreakpoint sets a breakpoint at the specified address
func (d *Debugger) setBreakpoint(args []string) {
	if len(args) == 0 {
		fmt.Printf("%sUsage: B <address>%s\n", Red, Reset)
		return
	}

	addr, err := d.parseAddress(args[0])
	if err != nil {
		fmt.Printf("%sError parsing address: %v%s\n", Red, err, Reset)
		return
	}

	d.breakpoints[addr] = true
	fmt.Printf("%sBreakpoint set at %s%s\n", Green, d.formatAddress(addr), Reset)
}

// clearBreakpoint clears a breakpoint at the specified address
func (d *Debugger) clearBreakpoint(args []string) {
	if len(args) == 0 {
		fmt.Printf("%sUsage: C <address>%s\n", Red, Reset)
		return
	}

	addr, err := d.parseAddress(args[0])
	if err != nil {
		fmt.Printf("%sError parsing address: %v%s\n", Red, err, Reset)
		return
	}

	if d.breakpoints[addr] {
		delete(d.breakpoints, addr)
		fmt.Printf("%sBreakpoint cleared at %s%s\n", Green, d.formatAddress(addr), Reset)
	} else {
		fmt.Printf("%sNo breakpoint at %s%s\n", Yellow, d.formatAddress(addr), Reset)
	}
}

// listBreakpoints lists all active breakpoints
func (d *Debugger) listBreakpoints() {
	if len(d.breakpoints) == 0 {
		fmt.Printf("%sNo breakpoints set%s\n", Yellow, Reset)
		return
	}

	fmt.Printf("%s%sActive Breakpoints:%s\n", Bold, Yellow, Reset)
	for addr := range d.breakpoints {
		fmt.Printf("  %s%s%s\n", Red, d.formatAddress(addr), Reset)
	}
	fmt.Println()
}

// zeroMemory zeros a range of memory
func (d *Debugger) zeroMemory(args []string) {
	if len(args) < 2 {
		fmt.Printf("%sUsage: Z <start_address> <end_address>%s\n", Red, Reset)
		return
	}

	startAddr, err := d.parseAddress(args[0])
	if err != nil {
		fmt.Printf("%sError parsing start address: %v%s\n", Red, err, Reset)
		return
	}

	endAddr, err := d.parseAddress(args[1])
	if err != nil {
		fmt.Printf("%sError parsing end address: %v%s\n", Red, err, Reset)
		return
	}

	if endAddr < startAddr {
		fmt.Printf("%sEnd address must be >= start address%s\n", Red, Reset)
		return
	}

	count := int(endAddr - startAddr + 1)
	for addr := startAddr; addr <= endAddr; addr++ {
		d.memory.Write(addr, 0x00)
	}

	fmt.Printf("%sZeroed %d bytes from %s to %s%s\n",
		Green, count, d.formatAddress(startAddr), d.formatAddress(endAddr), Reset)
}

// fillMemory fills a range of memory with a value
func (d *Debugger) fillMemory(args []string) {
	if len(args) < 3 {
		fmt.Printf("%sUsage: F <start_address> <end_address> <value>%s\n", Red, Reset)
		return
	}

	startAddr, err := d.parseAddress(args[0])
	if err != nil {
		fmt.Printf("%sError parsing start address: %v%s\n", Red, err, Reset)
		return
	}

	endAddr, err := d.parseAddress(args[1])
	if err != nil {
		fmt.Printf("%sError parsing end address: %v%s\n", Red, err, Reset)
		return
	}

	value, err := d.parseValue(args[2])
	if err != nil {
		fmt.Printf("%sError parsing value: %v%s\n", Red, err, Reset)
		return
	}

	if endAddr < startAddr {
		fmt.Printf("%sEnd address must be >= start address%s\n", Red, Reset)
		return
	}

	count := int(endAddr - startAddr + 1)
	for addr := startAddr; addr <= endAddr; addr++ {
		d.memory.Write(addr, value)
	}

	fmt.Printf("%sFilled %d bytes from %s to %s with %s%s\n",
		Green, count, d.formatAddress(startAddr), d.formatAddress(endAddr), d.formatByte(value), Reset)
}

// transferMemory transfers a block of memory
func (d *Debugger) transferMemory(args []string) {
	if len(args) < 3 {
		fmt.Printf("%sUsage: T <source_address> <dest_address> <length>%s\n", Red, Reset)
		return
	}

	srcAddr, err := d.parseAddress(args[0])
	if err != nil {
		fmt.Printf("%sError parsing source address: %v%s\n", Red, err, Reset)
		return
	}

	destAddr, err := d.parseAddress(args[1])
	if err != nil {
		fmt.Printf("%sError parsing destination address: %v%s\n", Red, err, Reset)
		return
	}

	length, err := strconv.ParseUint(args[2], 10, 16)
	if err != nil {
		fmt.Printf("%sError parsing length: %v%s\n", Red, err, Reset)
		return
	}

	// Copy the memory block
	for i := uint16(0); i < uint16(length); i++ {
		value := d.memory.Read(srcAddr + i)
		d.memory.Write(destAddr+i, value)
	}

	fmt.Printf("%sTransferred %d bytes from %s to %s%s\n",
		Green, length, d.formatAddress(srcAddr), d.formatAddress(destAddr), Reset)
}
