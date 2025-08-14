package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jrsteele09/go-6502-emulator/debugger"
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

type DebuggerRepl struct {
	debugger    *debugger.Debugger
	scanner     *bufio.Scanner
	lastCommand string
}

func NewDebuggerRepl() *DebuggerRepl {
	return &DebuggerRepl{
		debugger: debugger.NewDebugger(),
		scanner:  bufio.NewScanner(os.Stdin),
	}
}

func main() {
	repl := NewDebuggerRepl()
	repl.Run()
}

// Run starts the debugger REPL
func (r *DebuggerRepl) Run() {
	r.printBanner()
	r.showHelp()

	for {
		r.printPrompt()

		if !r.scanner.Scan() {
			break // EOF or error
		}

		input := strings.TrimSpace(r.scanner.Text())
		if input == "" {
			input = r.lastCommand // Repeat last command
		} else {
			r.lastCommand = input
		}

		if input == "" {
			continue
		}

		if r.processCommand(input) {
			break // Exit requested
		}
	}

	fmt.Printf("%sQuitting...%s\n", Yellow, Reset)
}

// printBanner displays the debugger banner
func (r *DebuggerRepl) printBanner() {
	fmt.Printf("%s%s", Bold, Cyan)
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                    6502 Debugger v1.0.0                      ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Printf("%s", Reset)
	fmt.Println()
}

// printPrompt displays the debugger prompt
func (r *DebuggerRepl) printPrompt() {
	pc := r.debugger.GetCPU().Registers().PC
	fmt.Printf("%s.%s $%04X%s> ", Bold, Green, pc, Reset)
}

// processCommand processes a user command and returns true if exit is requested
func (r *DebuggerRepl) processCommand(input string) bool {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return false
	}

	command := strings.ToUpper(parts[0])
	args := parts[1:]

	switch command {
	case "H", "HELP", "?":
		r.showHelp()
	case "Q", "QUIT", "EXIT":
		return true
	case "D", "DISASSEMBLE":
		output := r.debugger.Disassemble(args)
		fmt.Print(colorizeOutput(output))
	case "M", "MEMORY":
		output := r.debugger.HexDump(args)
		fmt.Print(colorizeOutput(output))
	case "R", "REGISTERS":
		output := r.debugger.ShowRegisters()
		fmt.Print(colorizeOutput(output))
	case "L", "LOAD":
		if len(args) == 0 {
			fmt.Printf("%sUsage: L <filename>%s\n", Red, Reset)
		} else {
			output := r.debugger.LoadPRG(args[0])
			fmt.Print(colorizeOutput(output))
		}
	case "G", "GO":
		output := r.debugger.Go(args)
		fmt.Print(colorizeOutput(output))
	case "S", "STEP":
		output := r.debugger.Step(args)
		fmt.Print(colorizeOutput(output))
	case "B", "BREAK":
		output := r.debugger.SetBreakpoint(args)
		fmt.Print(colorizeOutput(output))
	case "BR", "BREAKPOINTS":
		output := r.debugger.ListBreakpoints()
		fmt.Print(colorizeOutput(output))
	case "C", "CLEAR":
		output := r.debugger.ClearBreakpoint(args)
		fmt.Print(colorizeOutput(output))
	case "Z", "ZERO":
		output := r.debugger.ZeroMemory(args)
		fmt.Print(colorizeOutput(output))
	case "F", "FILL":
		output := r.debugger.FillMemory(args)
		fmt.Print(colorizeOutput(output))
	case "T", "TRANSFER":
		output := r.debugger.TransferMemory(args)
		fmt.Print(colorizeOutput(output))
	default:
		fmt.Printf("%sUnknown command: %s. Type 'H' for help.%s\n", Red, command, Reset)
	}

	return false
}

// showHelp displays available commands
func (r *DebuggerRepl) showHelp() {
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

// colorizeOutput adds color formatting to debugger output
func colorizeOutput(output string) string {
	// Add basic colorization for common patterns
	lines := strings.Split(output, "\n")
	for i, line := range lines {
		if strings.Contains(line, "Error") {
			lines[i] = Red + line + Reset
		} else if strings.Contains(line, "Loaded") || strings.Contains(line, "set") || strings.Contains(line, "cleared") {
			lines[i] = Green + line + Reset
		} else if strings.Contains(line, "Registers:") || strings.Contains(line, "Disassembly") || strings.Contains(line, "Memory dump") {
			lines[i] = Bold + Yellow + line + Reset
		} else if strings.Contains(line, ">") {
			// Current PC marker
			lines[i] = strings.Replace(line, ">", Green+">"+Reset, 1)
		} else if strings.Contains(line, "*") {
			// Breakpoint marker
			lines[i] = strings.Replace(line, "*", Red+"*"+Reset, 1)
		}
	}
	return strings.Join(lines, "\n")
}
