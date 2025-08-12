# 6502 Debugger

A comprehensive REPL debugger for the 6502 processor, inspired by the Commodore 128 monitor functionality. Features colorized terminal output and complete debugging capabilities.

## Features

- ✅ **Full 6502 CPU Emulation** - Complete register set and memory operations
- ✅ **PRG File Loading** - Load Commodore 64 PRG files into memory
- ✅ **Interactive Disassembly** - Real-time disassembly with syntax highlighting
- ✅ **Memory Hex Dump** - Colorized memory viewing with ASCII representation
- ✅ **Breakpoint Support** - Set, list, and clear breakpoints
- ✅ **Single Stepping** - Step through instructions one by one
- ✅ **Program Execution** - Run programs with breakpoint support
- ✅ **Memory Operations** - Fill, zero, and transfer memory blocks
- ✅ **Register Display** - View CPU registers and status flags
- ✅ **Colorized Output** - Terminal colors for better readability
- ✅ **Command History** - Repeat last command by pressing Enter

## Installation

```bash
cd cmd/debugger
go build -o debug6502
```

## Usage

```bash
./debug6502
```

## Commands

### Basic Commands
- `H`, `HELP`, `?` - Show help information
- `Q`, `QUIT`, `EXIT` - Exit debugger
- `R`, `REGISTERS` - Display CPU registers and status flags

### Memory Operations
- `M [addr] [count]` - Memory hex dump (default: $0000, 16 bytes)
- `Z <start> <end>` - Zero memory range
- `F <start> <end> <val>` - Fill memory range with value
- `T <src> <dest> <len>` - Transfer memory block

### Disassembly and Execution
- `D [addr] [count]` - Disassemble memory (default: PC, 10 instructions)
- `L <filename>` - Load PRG file into memory
- `G [addr]` - Go/Run from address (default: current PC)
- `S [count]` - Step instruction(s) (default: 1)

### Breakpoints
- `B <addr>` - Set breakpoint at address
- `BR`, `BREAKPOINTS` - List all breakpoints
- `C <addr>` - Clear breakpoint at address

## Address Formats

Addresses can be specified in two formats:
- **Hexadecimal**: `$1000`, `$FF00`
- **Decimal**: `4096`, `65280`

## Example Session

```
╔══════════════════════════════════════════════════════════════╗
║                    6502 Debugger v1.0.0                     ║
║              Commodore 128 Monitor Style REPL               ║
╚══════════════════════════════════════════════════════════════╝

. 0000> L test.prg
Loaded PRG file: test.prg
  Segment 1: $1000 to $100A (11 bytes)
Total: 11 bytes loaded
PC set to $1000

. 1000> R
Registers:
  A: $00  X: $00  Y: $00  PC: $1000  S: $FF
  Flags: $24 (.....I..)  NV.BDIZC

. 1000> D
Disassembly from $1000:

> 1000:  A9 42      LDA #$42
  1002:  8D 20 D0   STA $D020
  1005:  A2 10      LDX #$10
  1007:  CA         DEX
  1008:  D0 FD      BNE $1007
  100A:  60         RTS

. 1000> B $1008
Breakpoint set at $1008

. 1000> G
Running from $1000... (Ctrl+C to break)

Breakpoint hit at $1008
Next: 1008:  D0 FD      BNE $1007

. 1008> S
Step 1: 1008:  D0 FD      BNE $1007
Registers:
  A: $42  X: $0F  Y: $00  PC: $1007  S: $FF
  Flags: $24 (.....I..)  NV.BDIZC
```

## Color Scheme

The debugger uses ANSI colors for enhanced readability:

- **Cyan** - Commands, labels, and addresses
- **Green** - Success messages and current PC indicator
- **Red** - Errors and breakpoint indicators
- **Yellow** - Section headers and warnings
- **Gray** - ASCII representation in hex dumps
- **Dim** - Secondary information

## Register Display

The register display shows:
- **A, X, Y** - Accumulator and index registers
- **PC** - Program Counter (16-bit)
- **S** - Stack Pointer
- **Flags** - Status register with individual flag breakdown

Status flags are displayed as:
- **N** - Negative
- **V** - Overflow
- **B** - Break
- **D** - Decimal mode
- **I** - Interrupt disable
- **Z** - Zero
- **C** - Carry

## Memory Operations

### Hex Dump Format
```
1000: A9 42 8D 20 D0 A2 10 CA D0 FD 60 00 00 00 00 00  .B. ......`.....
```
- Address column (cyan)
- Hex bytes
- ASCII representation (gray)

### Disassembly Format
```
> 1000:  A9 42      LDA #$42
  1002:  8D 20 D0   STA $D020
* 1008:  D0 FD      BNE $1007
```
- `>` indicates current PC
- `*` indicates breakpoint
- Address, hex bytes, and assembly instruction

## Breakpoint System

- Set breakpoints at any address using `B <addr>`
- Program execution stops when PC reaches a breakpoint
- List all breakpoints with `BR`
- Clear individual breakpoints with `C <addr>`
- Breakpoints are preserved across program runs

## File Loading

The debugger supports loading PRG files:
- Automatically sets PC to the first segment's load address
- Displays memory layout after loading
- Shows total bytes loaded
- Supports multiple segments (though rare in PRG files)

## Keyboard Controls

- **Enter** - Repeat last command
- **Ctrl+C** - Break running program (returns to prompt)
- Standard line editing with backspace

## Error Handling

The debugger provides clear error messages for:
- Invalid addresses
- File loading errors
- Invalid command syntax
- Memory access violations

## Technical Details

- **Memory**: 64KB address space ($0000-$FFFF)
- **CPU**: Full 6502 instruction set support
- **Stack**: Starts at $FF (page $01)
- **Reset Vector**: $FFFC (standard 6502)
- **File Formats**: PRG (Commodore program files)

## Building from Source

Requires Go 1.19 or later:

```bash
cd /path/to/go-6502-emulator
go build -o debug6502 ./cmd/debugger
```

## Integration

The debugger uses the same CPU and memory components as the main emulator, ensuring 100% compatibility with assembled programs from the included assembler.
