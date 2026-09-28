# 6502 assembler language reference

This document describes the syntax the assembler currently accepts. It is a
reference to implemented behaviour, rather than a list of syntax that might be
added later. For the available instructions, opcodes, byte counts, and cycle
counts, see the [6502 instruction set](instruction-set.md).

## Quick example

```asm
SCREEN = $0400
COLOUR EQU $d800

STORE macro value, address
    lda #value
    sta address
endm

org $1000

start:
    STORE $41, SCREEN
    STORE 1, COLOUR

-   dex
    bne -

    rts
```

Mnemonics, directives, macro names, and preprocessor keywords are
case-insensitive. Named labels and constants are case-sensitive, so `Loop` and
`loop` are different symbols.

## Source lines and comments

The assembler accepts one statement per line. A label may appear by itself or
before an instruction or assembler directive on the same line.

```asm
loop:       dex             ; label with a colon
again       iny             ; colonless label
value       db $01, $02     ; label before a directive
done                        ; label-only line
```

Three comment forms are supported:

```asm
lda #1          ; comment to end of line
lda #2          // comment to end of line
/* a block comment
   may span lines */
```

Comment markers inside quoted text are not treated as comments.

## Instructions and addressing

The assembler's instruction set is constructed from the CPU opcode table.
The command-line assembler enables the full table, including the undocumented
opcodes implemented by the emulator. Instruction mnemonics are
case-insensitive.

| Addressing mode | Example |
| --- | --- |
| Implied | `clc` |
| Accumulator | `asl a` |
| Immediate | `lda #$40` |
| Zero page | `lda $40` |
| Zero page,X | `lda $40,x` |
| Zero page,Y | `ldx $40,y` |
| Absolute | `lda $c000` |
| Absolute,X | `lda $c000,x` |
| Absolute,Y | `lda $c000,y` |
| Indexed indirect | `lda ($40,x)` |
| Indirect indexed | `lda ($40),y` |
| Indirect | `jmp ($c000)` |
| Relative | `bne loop` |

For a known numeric address at or below `$ff`, the assembler chooses a
zero-page form when that mnemonic provides one. Otherwise it chooses the
absolute form. A forward-referenced named label is conservatively sized as a
16-bit address during layout, except for branch instructions, which are always
relative.

Branch displacements must fit in the 6502's signed 8-bit range, from -128 to
+127 bytes relative to the instruction following the branch.

## Labels

### Named labels

A named label starts with a letter or underscore and may then contain letters,
digits, or underscores. Both colon and colonless forms are accepted:

```asm
start:
    lda #0

loop    inx
        bne loop
```

Named labels may be referenced before or after their definition:

```asm
    jmp finished
    nop
finished:
    rts
```

Labels may be used in expressions, data directives, and instruction operands:

```asm
table:
    db 1, 2, 3

    lda table+1
    dw table
    lda #<table
    ldx #>table
```

A named label may only be defined once and may not use the same name as a
constant.

### Anonymous `+` and `-` labels

`+` defines a forward anonymous label and `-` defines a backward anonymous
label. A reference resolves to the nearest matching definition in the required
direction:

```asm
-   dex
    bne -          ; nearest preceding "-"

    jmp +          ; nearest following "+"
    nop
+   rts
```

Repeated signs form separate label names. Thus `+`, `++`, and `+++` are three
independent namespaces, as are `-`, `--`, and `---`:

```asm
-       jmp +++
        bne -
+++     lda #1
```

Anonymous labels can be defined more than once; each reference selects the
closest definition of exactly the same number of signs.

### Current address

`*` in an expression is the program counter at that source position:

```asm
here = *
    dw *
    jmp *
```

Branches also accept explicit relative displacements such as `bne *+4` and
`bne *-6`.

### `!1`-style local labels

`!1`, `!1+`, and `!1-` local-label syntax is **not currently supported**.
Although `!` is recognised as a lexer token, the assembler does not define or
resolve bang-number labels. Use a named label or the supported `+`/`-`
anonymous labels instead.

Do not confuse `!1` with `\1`: the latter is a positional macro argument and
is supported inside a macro body.

## Constants and `EQU`

All of the following constant forms are supported:

```asm
BASE = $40
FIRST EQU BASE + 1
SECOND .EQU FIRST + 1
.EQU THIRD, SECOND + 1
.EQU FOURTH = THIRD + 1
```

`EQU` and `.EQU` are case-insensitive. Constant names follow the same
identifier spelling rules as named labels. Constants are processed in source
order, so a constant expression should only refer to constants already
defined.

Constants are reassignable. This is useful for source-time counters:

```asm
test_number = 0
db test_number
test_number = test_number + 1
db test_number
```

The following equivalent variable directive is also accepted:

```asm
var VALUE = $42
.var OTHER = VALUE + 1
```

Constants and labels share a namespace during assembly and cannot be defined
with the same spelling.

## Expressions

Expressions are accepted in instruction operands, constants, conditionals,
origins, and numeric data directives.

### Literals and symbols

| Form | Meaning | Example |
| --- | --- | --- |
| Decimal | Base-10 integer | `123` |
| Hexadecimal | `$` prefix | `$c000` |
| Binary | `%` prefix | `%10101010` |
| Character | Single quoted character | `'A'` |
| Symbol | Constant or label | `table` |
| Current address | Current program counter | `*` |

Character literals support `\n`, `\r`, `\t`, `\\`, and `\'` escapes.

### Operators

The operators below are listed from lowest to highest precedence. Parentheses
may be used to make evaluation explicit.

| Precedence | Operators | Meaning |
| --- | --- | --- |
| 1 | `=`, `==`, `!=`, `<`, `<=`, `>`, `>=` | Comparison; returns 0 or 1 |
| 2 | `\|` | Bitwise OR |
| 3 | `^` | Bitwise XOR |
| 4 | `&` | Bitwise AND |
| 5 | `<<`, `>>` | Left and right shift |
| 6 | `+`, `-` | Addition and subtraction |
| 7 | `*`, `/` | Multiplication and integer division |
| 8 | unary `-`, `~`, `<`, `>` | Negate, complement, low byte, high byte |

`<expression` selects the low byte and `>expression` selects the high byte.
The equivalent `lo` and `hi` functions accept either form shown here:

```asm
lda #<table
ldx #>table
db lo(table), hi(table)
db lo $1234, hi $1234
```

Immediate expressions may be parenthesised:

```asm
lda #(BASE + 3)
```

There is no modulo operator. Division by zero and a negative shift count are
reported as errors.

## Macros

Macros are source-level text expansions. Both declaration styles are
supported:

```asm
LOAD_STORE macro value, address
    lda #value
    sta address
endm

macro CLEAR address
    lda #0
    sta address
endm
```

Parameters in a declaration may be separated by commas or whitespace. Macro
arguments at an invocation are comma-separated:

```asm
LOAD_STORE $42, $d020
CLEAR $d021
```

Named parameters may be referenced as a bare word or with a leading
backslash:

```asm
PAIR macro first, second
    db first, \second
endm
```

Positional arguments use `\1`, `\2`, and so on:

```asm
LOAD_IMMEDIATE macro
    lda #\1
endm

LOAD_IMMEDIATE $42
```

When named parameters are declared, the invocation must provide exactly that
many arguments. When no parameters are declared, the highest positional
reference determines the minimum required argument count.

Macro names are case-insensitive. Expansions may invoke other macros and may
contain constants, includes, and conditionals. Recursive/nested expansion is
limited to 32 levels. Variadic parameters and macro-local label syntax are not
implemented.

## Conditional assembly

Conditional assembly uses `if`, optional `else`, and `endif`:

```asm
CPU_TYPE = 0

if CPU_TYPE != 1
    db $01
else
    db $02
endif
```

Conditions use the normal expression syntax. Zero is false and any non-zero
value is true. Conditionals may be nested.

Constants are evaluated during source preprocessing. A condition may also use
a label that has already been assigned when the layout pass reaches the
condition:

```asm
org $1000
start:
if start = $1000
    db $42
endif
```

A forward label cannot be used to decide a conditional because the selected
branch could itself change that label's address; it is reported as undefined.
Content inside an inactive branch is not expanded, and an include in an
inactive branch is not opened.

Only the bare forms `if`, `else`, and `endif` are supported. `.if`, `ifdef`,
`ifndef`, and `elseif` are not implemented.

## Assembler directives

Except where noted, assembler directives may be written with or without a
leading period and are case-insensitive.

| Directive | Accepted names | Description |
| --- | --- | --- |
| Origin | `org`, `.org` | Set the program counter and begin a segment at that address. |
| Origin assignment | `* = expression` | Alias for setting the origin. |
| Byte data | `byte`, `.byte`, `db`, `.db` | Emit one byte for each expression. |
| Word data | `word`, `.word`, `dw`, `.dw` | Emit each 16-bit expression low byte first. |
| Text | `text`, `string`, `str`, `asc` and dotted forms | Emit the bytes of one quoted string without a terminator. |
| Zero-terminated text | `asciiz`, `.asciiz` | Emit one quoted string followed by `$00`. |
| Data space | `ds`, `.ds` | Reserve the requested number of bytes, emitted as zeroes. |
| Variable | `var`, `.var` | Define or update a constant using `NAME = expression`. |
| End marker | `end`, `.end` | Accept and ignore the rest of this source line. |
| Constant | `EQU` variants | Define a source constant; see [Constants and `EQU`](#constants-and-equ). |
| Include | `#include`, `.include` | Include another source file; only available through file assembly. |
| Import once | `#importonce` | Prevent subsequent processing of the current included file. |

### Numeric data

Byte and word directives accept expression lists. Commas are recommended;
whitespace-separated values are also accepted:

```asm
db $01, $02, <table, >table
byte $03 $04
dw $1234, table
```

Byte values must fit in 8 bits and word values in 16 bits. Words are emitted in
6502 little-endian order.

### Text data

```asm
text "READY"
string "PRESS FIRE"
str "OK"
asc "ABC"
asciiz "DONE"
```

The first four forms are aliases and do not append a terminator. `asciiz`
appends one zero byte.

### Origins, segments, and storage

```asm
org 0
workspace ds 17

org $0200
start:
    rts
```

Each change of origin finishes the current output segment and starts another
when bytes are emitted. `ds` contributes real zero-filled bytes to that
segment; it is not an uninitialised linker reservation. Origins must fit in the
16-bit address space, and `ds` rejects negative sizes.

There is no separate linker, relocation system, or named section system.
Words such as `BSS` and `CODE` are **not directives**: when used alone as in
the Klaus decimal test, they are accepted as ordinary colonless labels. The
following `org` directives are what actually create its `$0000` and `$0200`
segments.

### `END`

`end`/`.end` is accepted for source compatibility. It consumes the remainder
of its line but does not stop reading the file and does not record an entry
point. Thus the operand in `end start` is currently ignored.

Because colonless labels are supported, `end` is treated as a label when it is
followed by another statement on the same line; otherwise it is the directive.

## Includes and `#importonce`

Includes support either spelling and quoted or unquoted paths:

```asm
#include "constants.asm"
.include 'macros.asm'
.include data/tables.asm
```

Includes are expanded in source order. Included files share the caller's macro
and constant context, so a macro defined in an included file is available
after the include in the parent file.

```asm
// macros.asm
EMIT macro value
    db value
endm

// main.asm
.include "macros.asm"
EMIT $42
```

Place `#importonce` in an included file to skip later includes of the same
resolved path. Circular includes are rejected, and include nesting is limited
to 10 levels.

Includes require `Assembler.AssembleFile` and a `utils.FileResolver`. The
reader-based `Assembler.Assemble` method has no resolver and therefore does not
expand includes.

## Klaus decimal-mode source compatibility

The repository's Klaus/Bruce Clark decimal-mode source is assembled as an
integration test. Its relevant syntax is supported as follows:

| Source construct | Behaviour |
| --- | --- |
| `name = value` | Reassignable source constant |
| `end_of_test macro ... endm` | Parameterless macro |
| `if cputype != 1 ... endif` | Conditional expression |
| `BSS` and `CODE` | Ordinary labels, not section directives |
| `org 0` and `org $200` | Create the data and code segments |
| `N1 ds 1` | Colonless label plus zero-filled storage |
| `TEST ldy #1` | Colonless label plus instruction |
| `N2H+1`, `N2H,x` | Label arithmetic and indexed addressing |
| `end TEST` | Compatibility end marker; `TEST` is ignored |

The test expects a 17-byte zero-filled segment at `$0000` and a code segment at
`$0200`, then executes the assembled test against the emulator.

## Explicitly unsupported syntax

The following commonly seen assembler features are not currently implemented:

- `!1`, `!1+`, or `!1-` local labels
- `.if`, `ifdef`, `ifndef`, and `elseif`
- variadic macros or macro-local label declarations
- named sections such as true `BSS`, `CODE`, or `.segment` directives
- `.align`, `.fill`, `.res`, and `.incbin`
- relocation, linking, exports, and imports of symbols
- a modulo operator

## Command-line use

Build and invoke the assembler from the repository root:

```bash
go build -o asm6502 ./cmd/assembler
./asm6502 -i program.asm
```

| Option | Purpose |
| --- | --- |
| `-i` | Input assembly file; required |
| `-o` | Output filename; defaults from the input name |
| `-f` | Output format: `prg`, `d64`, or `t64` |
| `-n` | Program name stored in D64 or T64 output |
| `-v` | Print segment and output details |
| `-h` | Show help |
| `-version` | Show the assembler version |

PRG output currently accepts one segment. D64 and T64 output combine multiple
segments into one loadable image and zero-fill gaps. Those two container
writers do not currently diagnose overlapping segments; later segment data is
copied over earlier data in the overlap.

## Go API

Create an assembler from the CPU opcode table:

```go
mem := memory.NewMemory[uint16](64 * 1024)
processor := cpu.NewCPU(mem, true)
asm := assembler.New(processor.OpCodes())
```

Use `Assemble` for an `io.Reader` that has no includes:

```go
segments, err := asm.Assemble(strings.NewReader(source), "program.asm")
```

Use `AssembleFile` with a resolver for include support:

```go
resolver := utils.NewOSFileResolver("./asm")
segments, err := asm.AssembleFile("main.asm", resolver)
```

Both methods reset the label, constant, anonymous-label, and program-counter
state before each assembly. The result is a slice of `AssembledData`, with a
start address and byte buffer for each origin-created segment.
