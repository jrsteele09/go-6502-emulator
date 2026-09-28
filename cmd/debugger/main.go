package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"

	"github.com/ergochat/readline"
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
	lineEditor  *readline.Instance
	lastCommand string
}

func NewDebuggerRepl() (*DebuggerRepl, error) {
	lineEditor, err := readline.NewFromConfig(&readline.Config{
		HistoryLimit:      500,
		HistorySearchFold: true,
		InterruptPrompt:   "^C",
		EOFPrompt:         "^D",
	})
	if err != nil {
		return nil, err
	}

	return &DebuggerRepl{
		debugger:   debugger.NewDebugger(),
		lineEditor: lineEditor,
	}, nil
}

func main() {
	repl, err := NewDebuggerRepl()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to initialise debugger input: %v\n", err)
		os.Exit(1)
	}
	defer repl.lineEditor.Close()

	interrupts := make(chan os.Signal, 1)
	signal.Notify(interrupts, os.Interrupt)
	defer signal.Stop(interrupts)

	go func() {
		for range interrupts {
			repl.debugger.Stop()
		}
	}()

	// Auto-load any files passed on the command line
	if len(os.Args) > 1 {
		repl.AutoLoad(os.Args[1:])
	}
	repl.Run()
}

// Run starts the debugger REPL
func (r *DebuggerRepl) Run() {
	r.printBanner()
	r.showHelp()

	for {
		r.lineEditor.SetPrompt(r.prompt())
		line, err := r.lineEditor.ReadLine()
		if errors.Is(err, readline.ErrInterrupt) {
			continue
		}
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			fmt.Printf("%sInput error: %v%s\n", Red, err, Reset)
			break
		}

		input := strings.TrimSpace(line)
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

// AutoLoad loads a list of PRG files on startup before entering the REPL
func (r *DebuggerRepl) AutoLoad(files []string) {
	for _, f := range files {
		out := r.debugger.LoadPRG(f)
		fmt.Print(colorizeOutput(out))
	}
}

// printBanner displays the debugger banner
func (r *DebuggerRepl) printBanner() {
	fmt.Printf("%s%s", Bold, Cyan)
	fmt.Println("╔══════════════════════════════════════════════════════════════╗")
	fmt.Println("║                     6502 Debugger v0.6                       ║")
	fmt.Println("╚══════════════════════════════════════════════════════════════╝")
	fmt.Printf("%s", Reset)
	fmt.Println()
}

// prompt returns the debugger prompt with the current program counter.
func (r *DebuggerRepl) prompt() string {
	pc := r.debugger.GetCPU().Registers().PC
	return fmt.Sprintf("%s.%s $%04X%s> ", Bold, Green, pc, Reset)
}

// processCommand processes a user command and returns true if exit is requested
func (r *DebuggerRepl) processCommand(input string) bool {
	if value, matched := parseProgramCounterCommand(input); matched {
		if value == "" {
			fmt.Printf("%sUsage: PC=$C000%s\n", Red, Reset)
		} else {
			fmt.Print(colorizeOutput(r.debugger.SetProgramCounter(value)))
		}
		return false
	}

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
		if len(args) > 0 && isVerboseRegisterArgument(args[0]) {
			output = r.debugger.ShowRegistersVerbose()
		}
		fmt.Print(colorizeOutput(output))
	case "RV", "RD", "REGISTERS-VERBOSE", "REGISTERS-DESCRIPTIVE":
		output := r.debugger.ShowRegistersVerbose()
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

func isVerboseRegisterArgument(argument string) bool {
	switch strings.ToUpper(argument) {
	case "V", "VERBOSE", "D", "DESCRIPTIVE", "DETAILED":
		return true
	default:
		return false
	}
}

// parseProgramCounterCommand accepts monitor-style PC=$C000 assignments as
// well as the whitespace variants PC = $C000 and PC $C000.
func parseProgramCounterCommand(input string) (string, bool) {
	input = strings.TrimSpace(input)
	if len(input) < 2 || !strings.EqualFold(input[:2], "PC") {
		return "", false
	}
	if len(input) > 2 && input[2] != '=' && input[2] != ' ' && input[2] != '\t' {
		return "", false
	}

	value := strings.TrimSpace(input[2:])
	if strings.HasPrefix(value, "=") {
		value = strings.TrimSpace(value[1:])
	}
	return value, true
}

// showHelp displays available commands
func (r *DebuggerRepl) showHelp() {
	fmt.Printf("%s%sAvailable Commands:%s\n", Bold, Yellow, Reset)
	fmt.Printf("%s  H, HELP, ?%s        - Show this help\n", Cyan, Reset)
	fmt.Printf("%s  Q, QUIT, EXIT%s     - Exit debugger\n", Cyan, Reset)
	fmt.Printf("%s  R, REGISTERS%s      - Show CPU registers\n", Cyan, Reset)
	fmt.Printf("%s  RV, RD%s            - Show descriptive registers and named status flags\n", Cyan, Reset)
	fmt.Printf("%s  PC=<address>%s      - Set the program counter (for example PC=$C000)\n", Cyan, Reset)
	fmt.Printf("%s  D [addr] [count]%s  - Disassemble memory (default: PC, 10 instructions)\n", Cyan, Reset)
	fmt.Printf("%s  M [addr] [count]%s  - Memory hex dump (default: $0000, 16 bytes)\n", Cyan, Reset)
	fmt.Printf("%s  L <filename>%s      - Load PRG file into memory\n", Cyan, Reset)
	fmt.Printf("%s  G [addr]%s          - Go/Run from address (default: current PC)\n", Cyan, Reset)
	fmt.Printf("%s  S [count]%s         - Step instruction(s); stop before BRK/top-level RTS\n", Cyan, Reset)
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
	fmt.Printf("  - Up/Down browse command history; Left/Right edit the current line\n")
	fmt.Printf("  - Home/End or Ctrl+A/Ctrl+E jump to the start/end of the line\n")
	fmt.Printf("  - Option+Left/Right or Alt+B/Alt+F move by one word\n")
	fmt.Printf("  - Use Ctrl+C to stop a running program and return to the prompt\n")
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
		} else if strings.HasPrefix(line, "Registers") || strings.Contains(line, "Disassembly") || strings.Contains(line, "Memory dump") {
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
