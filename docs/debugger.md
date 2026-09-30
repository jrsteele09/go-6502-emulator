# 6502 Debugger User Manual

The debugger provides an interactive command prompt for loading PRG and BIN
files, inspecting 6502 state, disassembling memory, stepping instructions,
running to breakpoints, and editing memory.

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

A BIN file does not contain a load address, so supply one with `-address`:

```bash
debug6502 -address 4096 program.bin
debug6502 -address C000 program.bin
```

A bare address containing `A` through `F` is hexadecimal. If a `$` prefix is
used on the command line, quote the address so that the shell does not treat it
as a shell variable.

Display the version injected when the command was built:

```bash
debug6502 -version
```

Files are loaded from left to right. A PRG file's two-byte load address or the
address supplied for a BIN file determines where its bytes are placed. After
loading, the program counter is set to that address, so with multiple files it
finishes at the load address of the last file. Later files replace earlier
bytes if their address ranges overlap. The command-line address applies to
every BIN file in the list.

PRG and BIN files can be loaded directly. D64 and T64 containers are not
accepted by the debugger.

## The prompt

The prompt disassembles the instruction at the current program counter. It
shows the address, raw instruction bytes, and decoded instruction:

```text
$1000: A9 42      LDA #$42 >
```

The prompt refreshes after every command, so changes made by stepping, running,
loading a program, or assigning PC are immediately visible. An unrecognised
opcode is shown as `???`.

Commands and their long-form names are case-insensitive. The prompt keeps up
to 500 commands in memory for the current session and provides these editing
keys:

| Key | Action |
| --- | --- |
| Up or Ctrl+P | Show the previous command in history. |
| Down or Ctrl+N | Show the next command in history. |
| Left or Ctrl+B | Move left one character. |
| Right or Ctrl+F | Move right one character. |
| Home or Ctrl+A | Move to the beginning of the line. |
| End or Ctrl+E | Move to the end of the line. |
| Option/Alt+Left or Alt+B | Move left one word. |
| Option/Alt+Right or Alt+F | Move right one word. |
| Ctrl+R | Search backwards through command history. |
| Ctrl+W | Delete the preceding word. |
| Ctrl+U | Delete to the beginning of the line. |
| Ctrl+K | Delete to the end of the line. |

macOS Command+Arrow behavior is controlled by the terminal emulator. It works
when the terminal maps those shortcuts to Home/End or Ctrl+A/Ctrl+E; otherwise
use those portable bindings directly. Pressing Enter on an empty line repeats
the last non-empty command.

Debugger output is separated from the command and the following prompt by a
blank line. A `PC` assignment is the exception: it produces no status message
and the next prompt immediately shows the instruction at the new address.

Addresses accept `$`-prefixed hexadecimal, bare hexadecimal containing `A`
through `F`, or unsigned decimal:

```text
$1000
C000
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
| `RV`, `RD` | `REGISTERS-VERBOSE`, `REGISTERS-DESCRIPTIVE` | Display registers in a descriptive column and expand every status flag. |
| `PC=address` | — | Set the program counter. `PC address` and `PC = address` are also accepted. |
| `D [address] [count]` | `DISASSEMBLE` | Disassemble instructions. The default count is 10. |
| `M [address] [count]` | `MEMORY` | Display a hexadecimal and ASCII memory dump. Defaults to `$0000` and 16 bytes. |
| `L filename [address]` | `LOAD` | Load a PRG, or load a BIN at the required address, and set PC to its load address. |
| `G [address]` | `GO` | Run from an address, or from the current PC when omitted. |
| `S [count]` | `STEP` | Execute one or more instructions. Stops before `BRK` or a top-level `RTS`; the default count is 1. |
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
$0000: 00         BRK > L program.prg

Loaded PRG file: program.prg
  Segment 1: $1000 to $100A (11 bytes)
Total: 11 bytes loaded
PC set to $1000
```

Loading changes PC but does not reset A, X, Y, the stack pointer, status flags,
memory outside the loaded range, or existing breakpoints. Start a new debugger
session when a completely clean machine state is required.

Load a raw BIN file by supplying its destination address:

```text
$0000: 00         BRK > L program.bin $1000

Loaded BIN file: program.bin
  Segment 1: $1000 to $100A (11 bytes)
Total: 11 bytes loaded
PC set to $1000
```

## Inspecting registers

Use `R` to display all registers:

```text
$1000: A9 42      LDA #$42 > R

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

A letter means the flag is set; a period means it is clear. Set the program
counter with monitor-style assignment syntax:

```text
PC=$C000
```

Spaces around the equals sign are optional. `PC $C000` and decimal addresses
such as `PC=49152` are also accepted. The command prints no confirmation;
the updated address and instruction appear in the next prompt. Setting PC also
resets the implicit starting position used by the next `D` command. There is
currently no command for directly changing the other registers or flags.

For a column-oriented view with full names, use `RV` or `RD`. The forms
`R V`, `R D`, `R VERBOSE`, `R DESCRIPTIVE`, and `R DETAILED` are also
accepted:

```text
Registers:
  Accumulator          A    $80  (128)
  X index register     X    $10  (16)
  Y index register     Y    $20  (32)
  Stack pointer        S    $FD  (253)
  Program counter      PC   $C000  (49152)
  Processor status     P    $AA  (%10101010)

Processor status flags:
  Negative             N    set
  Overflow             V    clear
  Unused               U    set
  Break                B    clear
  Decimal mode         D    set
  Interrupt disable    I    clear
  Zero                 Z    set
  Carry                C    clear
```

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

A single step shows the registers after the instruction executes. A multi-step
command additionally prefixes each instruction with `Step n:` before showing
its resulting registers. Execution stops early if the resulting PC has a
breakpoint.

Like `G`, `S` stops before executing `BRK` or a top-level `RTS` reached with an
empty hardware stack (`S` is `$FF`). PC remains on the terminating instruction
and the stack is unchanged. An `RTS` with a JSR return address on the stack is
executed normally. A single `S` still executes an ordinary instruction when PC
is currently on a breakpoint, which is useful for moving past one before
continuing with `G`.

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

Running and stepping stop at a breakpoint or execution error as appropriate,
and both commands stop before a `BRK` or a top-level `RTS` reached while the
hardware stack is empty (`S` is `$FF`). PC remains on the terminating
instruction and the stack is unchanged. An `RTS` with a return address on the
stack executes normally.

Press Ctrl+C while `G` is running to request a stop. The current instruction
completes, then control returns to the debugger prompt with CPU and memory
state preserved. An infinite loop does not automatically return to the prompt.

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
$1000: A9 42      LDA #$42 > D $1000 6

> $1000: A9 42      LDA #$42
  $1002: 8D 00 02   STA $0200
  $1005: A2 03      LDX #$03
  $1007: CA         DEX
  $1008: D0 FD      BNE $1007
  $100A: 00         BRK

$1000: A9 42      LDA #$42 > B $1007

Breakpoint set at $1007

$1000: A9 42      LDA #$42 > G

Running from $1000... (Ctrl+C to stop)

Breakpoint hit at $1007
Next: $1007: CA         DEX

$1007: CA         DEX > R

  A: $42  X: $03  Y: $00  PC: $1007  S: $FF
  Flags: $24 (%00100100) (..1..I..)  NV1BDIZC

$1007: CA         DEX > S

  A: $42  X: $02  Y: $00  PC: $1008  S: $FF
  Flags: $24 (%00100100) (..1..I..)  NV1BDIZC

$1008: D0 FD      BNE $1007 > G

Running from $1008... (Ctrl+C to stop)

Breakpoint hit at $1007
Next: $1007: CA         DEX
```

Examine the value written by `STA $0200`:

```text
$1007: CA         DEX > M $0200 1
```

To run to the final `BRK`, replace the loop breakpoint with one at `$100A`:

```text
$1007: CA         DEX > C $1007
$1007: CA         DEX > B $100A
$1007: CA         DEX > G
```
