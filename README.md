# 6502 Assembler and Debugger

[![Go Report Card](https://goreportcard.com/badge/github.com/jrsteele09/go-6502-emulator)](https://goreportcard.com/report/github.com/jrsteele09/go-6502-emulator)
[![GoDoc](https://pkg.go.dev/badge/github.com/jrsteele09/go-6502-emulator)](https://pkg.go.dev/github.com/jrsteele09/go-6502-emulator)

A simple 6502 assembler and debugger written in Go. This project aims to provide everything you need to write, assemble, and debug 6502 assembly programs.

The assembler will continue to evolve over time, but the vision is to be compatible with as many 6502 assemblers as possible.

## Features

**Assembler:**
- Full 6502 instruction set support
- PRG file output format
- Comprehensive error reporting with line numbers
- Support for labels, constants, and expressions

**Debugger:**
- Interactive REPL-style debugger
- Memory examination and modification
- Disassembly with PC and breakpoint markers
- Single-step execution with register display
- Breakpoint management
- Program loading and execution control

## Installation

### Prerequisites
- Go 1.24 or later

### Building from Source

```bash
git clone https://github.com/jrsteele09/go-6502-emulator.git
cd go-6502-emulator
go build -o asm6502 ./cmd/assembler
go build -o debug6502 ./cmd/debugger
```

Place them in an appropriate location for your system and configure a path to that location.

## Assembler Usage

### Basic Usage

```bash
# Simple example
asm6502 -i program.s

# Assemble with custom output filename
asm6502 -i program.s -o myprogram.prg

# Verbose output showing assembly progress
asm6502 -i program.s -v

# Show help
asm6502 -h
```

## Debugger Usage

### Starting the Debugger

```bash
debug6502
```

This opens an interactive debugger session with a helpful prompt showing the current program counter.

### Basic Debugger Commands

```
Available Commands:
  H, HELP, ?        - Show help
  Q, QUIT, EXIT     - Exit debugger
  R, REGISTERS      - Show CPU registers
  D [addr] [count]  - Disassemble memory (default: PC, 10 instructions)
  M [addr] [count]  - Memory hex dump (default: $0000, 16 bytes)
  L <filename>      - Load PRG file into memory
  G [addr]          - Go/Run from address (default: current PC)
  S [count]         - Step instruction(s) (default: 1)
  B <addr>          - Set breakpoint at address
  BR, BREAKPOINTS   - List all breakpoints
  C <addr>          - Clear breakpoint at address
  Z <start> <end>   - Zero memory range
  F <start> <end> <val> - Fill memory range with value
  T <src> <dest> <len> - Transfer memory block
```

### Example Debugging Session

```bash
$ debug6502
╔══════════════════════════════════════════════════════════════╗
║                    6502 Debugger v1.0.0                      ║
╚══════════════════════════════════════════════════════════════╝

. $0000> L hello.prg        # Load your assembled program
Loaded PRG file: hello.prg
  Segment 1: $1000 to $11FF (512 bytes)
Total: 512 bytes loaded
PC set to $1000

. $1000> R                  # Show registers
Registers:
  A: $00  X: $00  Y: $00  PC: $1000  S: $FF
  Flags: $24 (%00100100) (..1..I..)  NV1BDIZC

. $1000> D                  # Disassemble from current PC
Disassembly from $1000:

> $1000: A2 00      LDX #$00
  $1002: A9 20      LDA #$20
  $1004: 9D 00 04   STA $0400,X
  $1007: 9D 00 05   STA $0500,X
  $100A: 9D 00 06   STA $0600,X
  $100D: 9D E8 07   STA $07E8,X
  $1010: E8         INX
  $1011: D0 F1      BNE $1004

. $1000> B $1010            # Set breakpoint
Breakpoint set at $1010

. $1000> G                  # Run program
Running from $1000...
Breakpoint hit at $1010
Next: $1010: E8         INX

. $1010> S                  # Step one instruction
Executing: $1010: E8         INX
Registers:
  A: $20  X: $01  Y: $00  PC: $1011  S: $FF
  Flags: $24 (%00100100) (..1..I..)  NV1BDIZC

. $1011> Q                  # Quit debugger
Quitting...
```

## Assembly Language Features

The assembler supports:

### Basic Assembly
```assembly
.ORG $1000             ; Set program origin

START:                 ; Labels for code organization
    LDA #$42          ; Load accumulator with hex value
    STA $D020         ; Store to memory location
    LDX #$10          ; Load X register
    
LOOP:
    DEX               ; Decrement X register
    BNE LOOP          ; Branch if not equal to zero
    RTS               ; Return from subroutine
```

### Data Definition
```assembly
.ORG $2000
    .BYTE $01, $02, $03    ; Define byte data
    .WORD $1234            ; Define 16-bit word
```

### Supported Features
- All standard 6502 instructions and addressing modes
- Labels and local symbols
- Hex literals ($42) and decimal numbers (66)
- Comments using semicolons (;)
- `.ORG` directive for setting assembly origin
- `.BYTE` and `.WORD` directives for data definition

## Examples

Check out the `examples/` directory for sample programs demonstrating:
- Basic assembly programming
- Memory operations and addressing modes
- Nested loops and branching

To build and run the examples:

```bash
# Assemble an example
asm6502 -i examples/hello.asm -o hello.prg

# Load it in the debugger
debug6502
> L hello.prg
> D
> S    # Step through the program
```

## License

MIT License. See [LICENSE](LICENSE) for more information.
