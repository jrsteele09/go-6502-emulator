# 6502 Assembler and Debugger

[![Go Reference](https://pkg.go.dev/badge/github.com/jrsteele09/go-6502-emulator)](https://pkg.go.dev/github.com/jrsteele09/go-6502-emulator)

A simple 6502 assembler and debugger written in Go. This project aims to provide everything you need to write, assemble, and debug 6502 assembly programs.

The assembler will continue to evolve over time, but the vision is to be compatible with as many 6502 assemblers as possible.

## Features

**Assembler:**
- Full 6502 instruction set support
- PRG, D64, and T64 file output formats
- Comprehensive error reporting with line numbers
- Support for labels, constants, expressions, and include files

**Debugger:**
- Interactive REPL-style debugger
- Memory examination and modification
- Disassembly with PC and breakpoint markers
- Single-step execution with register display
- Breakpoint management
- PRG loading and execution control

**Testing:**
- Functional testing with the [Klaus 6502 functional tests and decimal mode test](https://github.com/Klaus2m5/6502_65C02_functional_tests)

## Installation

### Prerequisites
- Go 1.26 or later

### Building from Source

```bash
git clone https://github.com/jrsteele09/go-6502-emulator.git
cd go-6502-emulator
go build -o asm6502 ./cmd/assembler
go build -o debug6502 ./cmd/debugger
```

Place them in an appropriate location for your system and configure a path to that location.

## Assembler Usage

See the [assembler guide](assembler/README.md) for the Go API, supported source
syntax, and an internal pipeline diagram.

### Basic Usage

```bash
# Simple example (outputs program.prg)
asm6502 -i program.s

# Assemble with custom output filename
asm6502 -i program.s -o myprogram.prg

# Assemble to D64 disk image
asm6502 -i program.s -f d64

# Assemble to T64 tape archive with custom program name
asm6502 -i program.s -f t64 -n MYPROG

# Verbose output showing assembly progress
asm6502 -i program.s -v

# Show version
asm6502 -version

# Show help
asm6502 -h
```

## Debugger Usage

### Starting the Debugger

```bash
debug6502
```

This opens an interactive debugger session, shows command help, and displays a prompt with the current program counter.

You can also auto-load one or more PRG files on startup:

```bash
debug6502 program.prg
```

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
║                     6502 Debugger v0.6                       ║
╚══════════════════════════════════════════════════════════════╝

. $0000> L basic.prg        # Load your assembled program
Loaded PRG file: basic.prg
  Segment 1: $1000 to $100A (11 bytes)
Total: 11 bytes loaded
PC set to $1000

. $1000> R                  # Show registers
Registers:
  A: $00  X: $00  Y: $00  PC: $1000  S: $FF
  Flags: $24 (%00100100) (..1..I..)  NV1BDIZC

. $1000> D                  # Disassemble from current PC
Disassembly from $1000:

> $1000: A9 42      LDA #$42
  $1002: 8D 20 D0   STA $D020
  $1005: A2 10      LDX #$10
  $1007: CA         DEX
  $1008: D0 FD      BNE $1007
  $100A: 60         RTS

. $1000> S                  # Step one instruction
Executing: $1000: A9 42      LDA #$42
Registers:
  A: $42  X: $00  Y: $00  PC: $1002  S: $FF
  Flags: $24 (%00100100) (..1..I..)  NV1BDIZC

. $1002> B $1007            # Set breakpoint
Breakpoint set at $1007

. $1002> G                  # Run program
Running from $1002... (Ctrl+C to break)

Breakpoint hit at $1007
Next: $1007: CA         DEX

. $1007> Q                  # Quit debugger
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
    .TEXT "HELLO"          ; Define string data
    .ASCIIZ "READY"        ; Define null-terminated string data
```

### Supported Features
- All standard 6502 instructions and addressing modes
- Labels and `+`/`-` shorthand labels
- Constants with `NAME = value` or `.VAR NAME = value`
- Expressions using arithmetic, label references, and constants
- Hex (`$42`), binary (`%01000010`), and decimal (`66`) literals
- Low-byte (`<LABEL`) and high-byte (`>LABEL`) address extraction
- Comments using semicolons (`;`), `//`, or `/* ... */`
- Include files with `#include "file.asm"` or `.include "file.asm"`
- Origin directives with `.ORG $1000` or `* = $1000`
- Data directives: `.BYTE`, `.DB`, `.WORD`, `.DW`, `.TEXT`, `.STRING`, `.STR`, `.ASC`, `.ASCIIZ`, and `.DS`

## License

MIT License. See [LICENSE](LICENSE) for more information.
