# 6502 Debugger User Manual

The debugger provides an interactive command prompt for loading Commodore PRG
files, inspecting 6502 state, disassembling memory, stepping instructions,
running to breakpoints, and editing memory.

It uses a flat 64 KiB address space. It does not provide C64 ROM routines or
hardware peripherals, so programs that depend on KERNAL calls, the VIC-II,
CIAs, or other memory-mapped devices cannot run as they would on a complete
C64.

## Starting the debugger

Start with empty memory:

```bash
debug6502
```

Pass one or more PRG files to load them before the prompt opens:

```bash
debug6502 program.prg
debug6502 code.prg data.prg
```

Files are loaded from left to right. Each file's two-byte PRG load address
determines where its payload is placed. After loading, the program counter is
set to that file's load address, so with multiple files it finishes at the
load address of the last file. Later files replace earlier bytes if their
address ranges overlap.

Only PRG files can be loaded directly. D64 and T64 containers are not accepted
by the debugger.

## The prompt

The prompt displays the current program counter:

```text
. $1000>
```

Commands and their long-form names are case-insensitive. Pressing Enter on an
empty line repeats the last non-empty command.

Addresses accept either `$`-prefixed hexadecimal or unsigned decimal:

```text
$1000
4096
```

The `0x1000` form is not accepted. Command counts and transfer lengths are
decimal. Fill values may be `$`-prefixed hexadecimal or decimal and must fit
in one byte.

## Command reference

| Command | Long form | Description |
| --- | --- | --- |
| `H` | `HELP`, `?` | Show command help. |
| `Q` | `QUIT`, `EXIT` | Exit the debugger. |
| `R` | `REGISTERS` | Display CPU registers and status flags. |
| `D [address] [count]` | `DISASSEMBLE` | Disassemble instructions. The default count is 10. |
| `M [address] [count]` | `MEMORY` | Display a hexadecimal and ASCII memory dump. Defaults to `$0000` and 16 bytes. |
| `L filename` | `LOAD` | Load a PRG file and set PC to its load address. |
| `G [address]` | `GO` | Run from an address, or from the current PC when omitted. |
| `S [count]` | `STEP` | Execute one or more instructions. The default count is 1. |
| `B address` | `BREAK` | Set a breakpoint. |
| `BR` | `BREAKPOINTS` | List active breakpoints. |
| `C address` | `CLEAR` | Clear a breakpoint. |
| `Z start end` | `ZERO` | Fill an inclusive address range with zero. |
| `F start end value` | `FILL` | Fill an inclusive address range with a byte value. |
| `T source destination length` | `TRANSFER` | Copy a decimal number of bytes. |

Arguments are separated by whitespace. A filename containing spaces cannot be
passed to `L`.

## Loading programs

Load a PRG after entering the debugger:

```text
. $0000> L program.prg
Loaded PRG file: program.prg
  Segment 1: $1000 to $100A (11 bytes)
Total: 11 bytes loaded
PC set to $1000
```

Loading changes PC but does not reset A, X, Y, the stack pointer, status flags,
memory outside the loaded range, or existing breakpoints. Start a new debugger
session when a completely clean machine state is required.

## Inspecting registers

Use `R` to display all registers:

```text
. $1000> R
Registers:
  A: $00  X: $00  Y: $00  PC: $1000  S: $FF
  Flags: $24 (%00100100) (..1..I..)  NV1BDIZC
```

The flag characters are shown in bit order:

| Flag | Meaning |
| --- | --- |
| `N` | Negative |
| `V` | Overflow |
| `1` | Unused status bit |
| `B` | Break |
| `D` | Decimal mode |
| `I` | Interrupt disable |
| `Z` | Zero |
| `C` | Carry |

A letter means the flag is set; a period means it is clear. There is currently
no command for directly changing registers or flags.

## Disassembling memory

Disassemble ten instructions from the current position:

```text
D
```

Specify a start address and decimal instruction count:

```text
D $1000 20
D 4096 20
```

The first `D` after startup or loading begins at PC. Subsequent `D` commands
without an address continue after the previous listing rather than returning
to PC. Supply the address explicitly to restart elsewhere.

The listing uses `>` for the current PC and `*` for a breakpoint:

```text
> $1000: A9 42      LDA #$42
  $1002: 8D 00 02   STA $0200
* $1005: A2 03      LDX #$03
```

If an address is both PC and a breakpoint, `*` is shown. An unrecognised opcode
is displayed as `???` and advances the listing by one byte. The debugger uses
the documented 6502 instruction set; undocumented opcodes are not executable
in the debugger.

## Examining memory

Display memory from an address for a decimal byte count:

```text
M $0200 32
M 512 32
```

Output is aligned to 16-byte rows and includes an ASCII representation.
Non-printable bytes appear as periods. A row is displayed in full even when
the requested address or count covers only part of it.

## Stepping instructions

Execute the instruction at PC:

```text
S
```

Execute several instructions:

```text
S 5
```

Each step shows the instruction and the registers after it completes. During a
multi-step command, execution stops early if the resulting PC has a
breakpoint. A single `S` still executes when PC is currently on a breakpoint,
which is useful for moving past one before continuing with `G`.

## Running and breakpoints

Set a breakpoint, run, and inspect the stopped instruction:

```text
B $1007
G $1000
D $1007 5
R
```

`G` checks for a breakpoint before executing the instruction at that address.
When it stops, PC still points at the breakpoint instruction. To continue past
it, either step once or clear it:

```text
S
G
```

or:

```text
C $1007
G
```

List and clear breakpoints with:

```text
BR
C $1007
```

Running stops at a breakpoint or an execution error. Press Ctrl+C while `G` is
running to request a stop. The current instruction completes, then control
returns to the debugger prompt with CPU and memory state preserved. `BRK`,
`RTS`, and an infinite loop do not automatically return to the prompt.

## Editing memory

Zero an inclusive range:

```text
Z $2000 $20FF
```

Fill an inclusive range with a byte:

```text
F $0400 $07E7 $20
F 1024 2023 32
```

Set a single byte by using the same start and end address:

```text
F $D020 $D020 $06
```

Copy a block using a decimal length:

```text
T $2000 $3000 256
```

Transfer copies forward one byte at a time. Avoid overlapping source and
destination ranges, particularly when the destination begins inside and above
the source range. Memory operations do not guard against wrapping through
`$FFFF`; keep every requested range within the 64 KiB address space. In
particular, do not use `$FFFF` as the inclusive end address of `Z` or `F`.

## Example debugging session

Suppose `loop.prg` loads at `$1000` and contains:

```asm
org $1000
    lda #$42
    sta $0200
    ldx #3
loop:
    dex
    bne loop
    brk
```

Start the debugger with the program already loaded:

```bash
debug6502 loop.prg
```

Inspect the program, stop at the loop, and step through an iteration:

```text
. $1000> D $1000 6
. $1000> B $1007
. $1000> G
Breakpoint hit at $1007
Next: $1007: CA         DEX

. $1007> R
. $1007> S
. $1008> R
. $1008> G
Breakpoint hit at $1007
```

Examine the value written by `STA $0200`:

```text
. $1007> M $0200 1
```

To run to the final `BRK`, replace the loop breakpoint with one at `$100A`:

```text
. $1007> C $1007
. $1007> B $100A
. $1007> G
```
