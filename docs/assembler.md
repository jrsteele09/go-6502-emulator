# 6502 Assembler User Manual

This manual describes the accepted source syntax and command-line behaviour.
For instruction meanings, opcodes, byte counts, and cycle counts, see the
[6502 instruction set](6502-instruction-set.md).

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
case-insensitive. Named labels, constants, and variables are case-sensitive,
so `Loop` and `loop` are different symbols.

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

The assembler accepts the documented 6502 instructions and the undocumented
instructions listed in the [instruction-set reference](6502-instruction-set.md).
Instruction mnemonics are case-insensitive.

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

An immediate operand always emits one byte.
For a non-immediate numeric expression, values from -128 through 255 select an
8-bit operand and values from -32768 through 65535 select a 16-bit operand. An
8-bit address operand is only valid when the instruction has a matching
zero-page addressing mode. There is no syntax for forcing an absolute encoding
of a small numeric value.

For a named label that is already known, an address at or below `$ff` selects
zero-page addressing when the instruction provides it; otherwise it selects
absolute addressing. Branch labels always select relative addressing.

When an address expression begins with a named symbol, operand width is chosen
from that symbol's value before any following arithmetic. Keep such arithmetic
within the same 8-bit or 16-bit range. If it crosses `$ff`, beginning the
expression with a literal, such as `0 + BASE + 1`, makes the final expression
value determine the width.

Avoid forward references to labels that will reside in zero page. Their
address is initially calculated using an absolute-size instruction, but the
final instruction can select the shorter zero-page form. Define zero-page
addresses before use, normally with `=` or `EQU`.

Branch displacements must fit in the 6502's signed 8-bit range, from -128 to
+127 bytes relative to the instruction following the branch.

## Labels

### Named labels

A named label starts with an ASCII letter or underscore and may then contain
letters, digits, or underscores. Both colon and colonless forms are accepted:

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
constant or variable.

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

## Constants and Variables

Constants can be defined using any of the supported `EQU` forms:

```asm
FIRST EQU $41
SECOND .EQU FIRST + 1
.EQU THIRD, SECOND + 1
.EQU FOURTH = THIRD + 1
.EQU FIFTH FOURTH + 1
```

`EQU` and `.EQU` are case-insensitive. Variables use assignment syntax, either
directly or through the `var` directive:

```asm
BASE = $40
var VALUE = $42
.var OTHER = VALUE + 1
```

Variables are reassignable, which is useful for source-time counters:

```asm
test_number = 0
db test_number
test_number = test_number + 1
db test_number
```

Constant and variable names follow the same identifier spelling rules as named
labels. Definitions are processed in source order, so an expression should
only refer to symbols already defined. Constants, variables, and labels share
a namespace: a constant or variable name cannot also be used as a label.

Despite the conventional distinction, all constant and variable forms can be
redefined. Use `EQU` to communicate that a value is intended to remain fixed,
and use `=` or `var` when reassignment is intended.

An `EQU` expression may use literals and earlier symbols declared with `EQU`
or bare `=`. It cannot refer to a label or to a symbol declared only with
`var`/`.var`. The `=` and `var` forms may use a previously defined label.

## Expressions

Expressions are accepted in instruction operands, constants, variables,
conditionals, origins, and numeric data directives.

### Literals and symbols

| Form | Meaning | Example |
| --- | --- | --- |
| Decimal | Base-10 integer | `123` |
| Hexadecimal | `$` prefix | `$c000` |
| Binary | `%` prefix | `%10101010` |
| Character | Single quoted character | `'A'` |
| Symbol | Constant, variable, or label | `table` |
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

Macro names are case-insensitive; parameter names are case-sensitive.
Expansions may invoke other macros and may contain constants, variables,
includes, and conditionals. Recursive/nested expansion is limited to 32
levels. Variadic parameters and macro-local label syntax are not supported.

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

Conditions may use constants, variables, and labels defined earlier in the
source:

```asm
org $1000
start:
if start = $1000
    db $42
endif
```

A forward label cannot be used to decide a conditional and is reported as
undefined. Content inside an inactive branch is ignored: macros and symbols in
it are not defined, and included files in it are not opened.

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
| Variable | `var`, `.var` | Define or update a variable using `NAME = expression`. |
| End marker | `end`, `.end` | Accept and ignore the rest of this source line. |
| Constant | `EQU` variants | Define a source constant; see [Constants and Variables](#constants-and-variables). |
| Include | `#include`, `.include`, `#import` | Include another source file; only available through file assembly. |

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
appends one zero byte. Text is emitted as its raw UTF-8 bytes; no ASCII-to-
PETSCII conversion is performed.

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
Words such as `BSS` and `CODE` are **not directives**. When used alone, they
are ordinary colonless labels. Use `org` to start data at a different address.

### `END`

`end`/`.end` is accepted for source compatibility. It consumes the remainder
of its line but does not stop reading the file and does not record an entry
point. Thus the operand in `end start` is currently ignored.

Because colonless labels are supported, bare `end` is treated as a label when
it is followed by another statement on the same line; otherwise it is the
directive. `.end` is always the directive.

## Includes

Includes support three directive spellings and quoted or unquoted paths:

```asm
#include "constants.asm"
.include 'macros.asm'
#import data/tables.asm
```

Includes are expanded in source order. Included files share macros, constants,
and variables with the including file, so definitions from an included file
are available after the include.

```asm
// macros.asm
EMIT macro value
    db value
endm

// main.asm
.include "macros.asm"
EMIT $42
```

Relative include paths are resolved from the directory containing the main
input file, including paths written inside nested includes. Circular includes
are rejected, and include nesting is limited to 10 levels.

`#importonce` may be placed inside an included file to ensure that file is
processed only once. After its first inclusion, later `#include`, `.include`,
or `#import` directives using the same path spelling are skipped. Equivalent
spellings such as `library.asm` and `./library.asm` are treated as different
paths.

## Command-line use

Assemble a source file with:

```bash
./asm6502 -i program.asm
```

Without `-o`, this produces `program.prg`. Use `-f` to select another output
format:

```bash
./asm6502 -i program.asm -o game.prg
./asm6502 -i program.asm -f bin
./asm6502 -i program.asm -f d64 -n GAME
./asm6502 -i program.asm -f t64 -n GAME
```

| Option | Purpose |
| --- | --- |
| `-i` | Input assembly file; required |
| `-o` | Output filename; defaults from the input name |
| `-f` | Output format: `prg`, `bin`, `d64`, or `t64` |
| `-n` | Program name stored in D64 or T64 output |
| `-verbose` | Print segment and output details |
| `-h` | Show help |
| `-v`, `-version` | Show the assembler version |

BIN output contains raw bytes without a header or load address. Multiple
segments are ordered by address and gaps are filled with zero bytes; overlapping
segments are rejected. PRG output currently accepts one segment. D64 and T64
output combine multiple segments into one loadable image and zero-fill gaps.
D64 and T64 do not reject overlapping segments; data from a later segment
replaces earlier data in the overlapping range.
