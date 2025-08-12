# 6502 Assembler Command Line Tool

A command-line assembler for the 6502 processor that outputs Commodore 64 compatible PRG files.

## Features

- ✅ Full 6502 instruction set support
- ✅ Include directive support (`#include`, `.include`, `.INCLUDE`)
- ✅ Nested includes with circular detection
- ✅ PRG file output format
- ✅ Comprehensive error handling
- ✅ Verbose output mode

## Installation

```bash
go build -o asm6502
```

## Usage

```bash
./asm6502 [options] -i <input.asm>
```

### Options

- `-i string` - Input assembly file (required)
- `-o string` - Output PRG file (default: input filename with .prg extension)
- `-v` - Verbose output
- `-h` - Show help
- `-version` - Show version

### Examples

```bash
# Basic assembly
./asm6502 -i game.asm

# Custom output file
./asm6502 -i game.asm -o output.prg

# Verbose assembly
./asm6502 -i game.asm -v
```

## Assembly Syntax

### Basic Instructions
```assembly
    LDA #$42        ; Load accumulator with immediate value
    STA $2000       ; Store accumulator at memory location
    LDX #$10        ; Load X register
    LDY #$20        ; Load Y register
    INX             ; Increment X
    INY             ; Increment Y
    JMP start       ; Jump to label
```

### Directives
```assembly
    .ORG $1000      ; Set origin address
```

### Include Files
```assembly
    #include "file.asm"     ; C-style include
    .include "file.asm"     ; Traditional assembler include
    .INCLUDE "file.asm"     ; Uppercase variant
```

### Labels
```assembly
start:
    LDA #$00
    JMP start       ; Jump back to start
```

## Output Format

The assembler generates Commodore 64 compatible PRG files:
- First 2 bytes: Load address (little-endian)
- Remaining bytes: Program data

## Error Handling

The assembler provides detailed error messages for:
- Missing or invalid files
- Syntax errors
- Circular includes
- Invalid instructions or addressing modes

## Example Program

**main.asm:**
```assembly
; Simple 6502 program
    .ORG $1000

start:
    LDA #$42
    STA $2000
    #include "subroutines.asm"
    JMP start
```

**subroutines.asm:**
```assembly
; Utility subroutines
clear_screen:
    LDA #$20
    STA $0400
    RTS
```

**Assembly:**
```bash
./asm6502 -i main.asm -v
```

**Output:**
```
6502 Assembler v1.0.0
Input file:  main.asm
Output file: main.prg

Assembling...
PRG load address: $1000
Program size: 15 bytes
  Writing segment at $1000: 15 bytes
Assembly completed successfully!
Generated 1 segment(s)
  Segment 1: 15 bytes at $1000
```

## Supported Instructions

Full 6502 instruction set including:
- Load/Store: LDA, LDX, LDY, STA, STX, STY
- Arithmetic: ADC, SBC, INC, DEC, INX, INY, DEX, DEY
- Logic: AND, ORA, EOR
- Shifts: ASL, LSR, ROL, ROR
- Branches: BEQ, BNE, BCS, BCC, BMI, BPL, BVS, BVC
- Jumps: JMP, JSR, RTS
- Stack: PHA, PLA, PHP, PLP
- Status: CLC, SEC, CLD, SED, CLI, SEI, CLV
- Compare: CMP, CPX, CPY
- Transfer: TAX, TAY, TXA, TYA, TSX, TXS
- System: BRK, RTI, NOP

## Addressing Modes

- Immediate: `#$42`
- Absolute: `$1000`
- Zero Page: `$42`
- Indexed: `$1000,X` / `$1000,Y`
- Indirect: `($1000)` / `($42,X)` / `($42),Y`
- Relative: Used automatically for branches
- Implied: No operand needed
