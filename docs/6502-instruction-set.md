# 6502 Instruction Set

A complete reference to the 6502 instructions, status flags, opcodes, addressing modes, sizes, and cycle counts, including common undocumented opcodes.

## Contents

| Section | Description |
|---|---|
| [CPU registers](#cpu-registers) | The accumulator, index registers, stack pointer, program counter, and processor status register |
| [Addressing notation](#addressing-notation) | Operand syntax and addressing modes |
| [Instructions](#instructions) | Documented instructions, status effects, opcodes, sizes, and cycle counts |
| [Undocumented instructions](#undocumented-instructions) | Common undocumented opcodes |

## CPU registers

| Register | Name | Size | Purpose |
|:---:|---|---:|---|
| `A` | Accumulator | 8 bits | Holds operands and results for arithmetic, logic, load, store, and shift operations. |
| `X` | X index register | 8 bits | Supplies an index for indexed addressing and can be used as a counter or temporary value. |
| `Y` | Y index register | 8 bits | Supplies an index for indexed addressing and can be used as a counter or temporary value. |
| `S` | Stack pointer | 8 bits | Selects the next location in the hardware stack at `$0100`–`$01FF`. |
| `PC` | Program counter | 16 bits | Holds the address of the next instruction byte to fetch. |
| `P` | Processor status | 8 bits | Holds the condition and control flags. |

### Processor status register

The processor status register, **P**, stores six flags. Its 8-bit status-byte
layout is:

| Bit | Symbol | Name | Notes |
|---:|:---:|---|---|
| 7 | `N` | Negative | Copies bit 7 of the result. |
| 6 | `V` | Overflow | `ADC` and `SBC`: signed overflow. `BIT`: copies bit 6 of the tested memory value. |
| 5 | — | Unused |  |
| 4 | `B` | Break | Not stored in P. `PHP` and `BRK` push B=1; IRQ and NMI push B=0. |
| 3 | `D` | Decimal mode | `SED` sets D=1; `CLD` sets D=0. `ADC` and `SBC` use BCD when D=1 and binary when D=0. |
| 2 | `I` | Interrupt disable | `SEI` sets I=1; `CLI` sets I=0. IRQ is disabled when I=1 and enabled when I=0; NMI is unaffected. |
| 1 | `Z` | Zero | Set when the result is zero. |
| 0 | `C` | Carry | Addition: carry out. Subtraction: no borrow. Shifts and rotates: bit shifted out. |

## Addressing notation

The instruction tables use the following notation for operands and addressing modes:

| Notation | Addressing mode | Meaning |
|---|---|---|
| — | Implied | The operand is implied by the instruction. |
| `A` | Accumulator | The accumulator is the operand. |
| `#$nn` | Immediate | The next byte is the value. |
| `$nn` | Zero page | An 8-bit address in page zero. |
| `$nn,X` | Zero page,X | Add X to the zero-page address; the result wraps within page zero. |
| `$nn,Y` | Zero page,Y | Add Y to the zero-page address; the result wraps within page zero. |
| `$nnnn` | Absolute | A complete 16-bit address. |
| `$nnnn,X` | Absolute,X | Add X to a 16-bit base address. |
| `$nnnn,Y` | Absolute,Y | Add Y to a 16-bit base address. |
| `($nn,X)` | Indexed indirect | Add X in page zero, then read the 16-bit target address from page zero. |
| `($nn),Y` | Indirect indexed | Read a 16-bit base address from page zero, then add Y. |
| `($nnnn)` | Indirect | Read the jump target through a 16-bit pointer. |
| `$relative` | Relative | A signed 8-bit displacement from the address following the instruction. |

## Instructions

### Instruction index

| Range | Instructions |
|---|---|
| A–C | [ADC](#adc--add-with-carry) · [AND](#and--logical-and) · [ASL](#asl--arithmetic-shift-left) · [BCC](#bcc--branch-if-carry-clear) · [BCS](#bcs--branch-if-carry-set) · [BEQ](#beq--branch-if-equal) · [BIT](#bit--bit-test) · [BMI](#bmi--branch-if-minus) · [BNE](#bne--branch-if-not-equal) · [BPL](#bpl--branch-if-plus) · [BRK](#brk--software-interrupt) · [BVC](#bvc--branch-if-overflow-clear) · [BVS](#bvs--branch-if-overflow-set) · [CLC](#clc--clear-carry) · [CLD](#cld--clear-decimal-mode) · [CLI](#cli--clear-interrupt-disable) · [CLV](#clv--clear-overflow) · [CMP](#cmp--compare-accumulator) · [CPX](#cpx--compare-x-register) · [CPY](#cpy--compare-y-register) |
| D–L | [DEC](#dec--decrement-memory) · [DEX](#dex--decrement-x) · [DEY](#dey--decrement-y) · [EOR](#eor--exclusive-or) · [INC](#inc--increment-memory) · [INX](#inx--increment-x) · [INY](#iny--increment-y) · [JMP](#jmp--jump) · [JSR](#jsr--jump-to-subroutine) · [LDA](#lda--load-accumulator) · [LDX](#ldx--load-x) · [LDY](#ldy--load-y) · [LSR](#lsr--logical-shift-right) |
| N–S | [NOP](#nop--no-operation) · [ORA](#ora--logical-inclusive-or) · [PHA](#pha--push-accumulator) · [PHP](#php--push-processor-status) · [PLA](#pla--pull-accumulator) · [PLP](#plp--pull-processor-status) · [ROL](#rol--rotate-left) · [ROR](#ror--rotate-right) · [RTI](#rti--return-from-interrupt) · [RTS](#rts--return-from-subroutine) · [SBC](#sbc--subtract-with-carry) · [SEC](#sec--set-carry) · [SED](#sed--set-decimal-mode) · [SEI](#sei--set-interrupt-disable) · [STA](#sta--store-accumulator) · [STX](#stx--store-x) · [STY](#sty--store-y) |
| T–Z | [TAX](#tax--transfer-accumulator-to-x) · [TAY](#tay--transfer-accumulator-to-y) · [TSX](#tsx--transfer-stack-pointer-to-x) · [TXA](#txa--transfer-x-to-accumulator) · [TXS](#txs--transfer-x-to-stack-pointer) · [TYA](#tya--transfer-y-to-accumulator) |

### ADC — Add with carry

Adds the operand and the carry flag to the accumulator. It updates carry, zero, negative, and overflow; decimal mode uses BCD arithmetic.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | In binary mode, copies bit 7 of the 8-bit result. In decimal mode it comes from an intermediate binary result. |
| `V` | Overflow | Set or cleared | In binary mode, set when two inputs with the same sign produce a result with the opposite sign. Decimal-mode results come from an intermediate binary result. |
| `D` | Decimal mode | Checked | Selects packed-BCD arithmetic when set and binary arithmetic when clear. |
| `Z` | Zero | Set or cleared | In binary mode, set when the 8-bit result is zero. In decimal mode it comes from an intermediate binary result. |
| `C` | Carry | Checked, then set or cleared | Its old value is added to A and the operand. Its new value is set when the unsigned sum produces a carry; in decimal mode it is the decimal carry. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `69` | Immediate | `ADC #$nn` | 2 | 2 |  |
| `65` | Zero page | `ADC $nn` | 2 | 3 |  |
| `75` | Zero page,X | `ADC $nn,X` | 2 | 4 |  |
| `6D` | Absolute | `ADC $nnnn` | 3 | 4 |  |
| `7D` | Absolute,X | `ADC $nnnn,X` | 3 | 4 | Page boundary crossed: +1 cycle |
| `79` | Absolute,Y | `ADC $nnnn,Y` | 3 | 4 | Page boundary crossed: +1 cycle |
| `61` | Indexed indirect | `ADC ($nn,X)` | 2 | 6 |  |
| `71` | Indirect indexed | `ADC ($nn),Y` | 2 | 5 | Page boundary crossed: +1 cycle |

### AND — Logical AND

ANDs the operand with the accumulator and updates zero and negative.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | Set when the new A value has bit 7 set; otherwise cleared. |
| `Z` | Zero | Set or cleared | Set when the new A value is zero; otherwise cleared. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `29` | Immediate | `AND #$nn` | 2 | 2 |  |
| `25` | Zero page | `AND $nn` | 2 | 3 |  |
| `35` | Zero page,X | `AND $nn,X` | 2 | 4 |  |
| `2D` | Absolute | `AND $nnnn` | 3 | 4 |  |
| `3D` | Absolute,X | `AND $nnnn,X` | 3 | 4 | Page boundary crossed: +1 cycle |
| `39` | Absolute,Y | `AND $nnnn,Y` | 3 | 4 | Page boundary crossed: +1 cycle |
| `21` | Indexed indirect | `AND ($nn,X)` | 2 | 6 |  |
| `31` | Indirect indexed | `AND ($nn),Y` | 2 | 5 | Page boundary crossed: +1 cycle |

### ASL — Arithmetic shift left

Shifts the accumulator or memory left by one bit. Bit 7 moves into carry and zero/negative reflect the result.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | Copies bit 7 of the shifted result. |
| `Z` | Zero | Set or cleared | Set when the shifted result is zero; otherwise cleared. |
| `C` | Carry | Set or cleared | Receives bit 7 of the original value. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `0A` | Accumulator | `ASL A` | 1 | 2 |  |
| `06` | Zero page | `ASL $nn` | 2 | 5 |  |
| `16` | Zero page,X | `ASL $nn,X` | 2 | 6 |  |
| `0E` | Absolute | `ASL $nnnn` | 3 | 6 |  |
| `1E` | Absolute,X | `ASL $nnnn,X` | 3 | 7 |  |

### BCC — Branch if carry clear

Branches when the carry flag is clear.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `C` | Carry | Checked | The branch is taken when Carry is clear. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `90` | Relative | `BCC $relative` | 2 | 2 | Branch taken: +1 cycle; page boundary crossed: +2 cycles total |

### BCS — Branch if carry set

Branches when the carry flag is set.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `C` | Carry | Checked | The branch is taken when Carry is set. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `B0` | Relative | `BCS $relative` | 2 | 2 | Branch taken: +1 cycle; page boundary crossed: +2 cycles total |

### BEQ — Branch if equal

Branches when the zero flag is set.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `Z` | Zero | Checked | The branch is taken when Zero is set. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `F0` | Relative | `BEQ $relative` | 2 | 2 | Branch taken: +1 cycle; page boundary crossed: +2 cycles total |

### BIT — Bit test

Tests the accumulator against memory without changing either value. Zero reflects `A AND operand`; bits 7 and 6 of memory become negative and overflow.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | Copies bit 7 of the memory operand. |
| `V` | Overflow | Set or cleared | Copies bit 6 of the memory operand. |
| `Z` | Zero | Set or cleared | Set when A AND the memory operand is zero; otherwise cleared. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `24` | Zero page | `BIT $nn` | 2 | 3 |  |
| `2C` | Absolute | `BIT $nnnn` | 3 | 4 |  |

### BMI — Branch if minus

Branches when the negative flag is set.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Checked | The branch is taken when Negative is set. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `30` | Relative | `BMI $relative` | 2 | 2 | Branch taken: +1 cycle; page boundary crossed: +2 cycles total |

### BNE — Branch if not equal

Branches when the zero flag is clear.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `Z` | Zero | Checked | The branch is taken when Zero is clear. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `D0` | Relative | `BNE $relative` | 2 | 2 | Branch taken: +1 cycle; page boundary crossed: +2 cycles total |

### BPL — Branch if plus

Branches when the negative flag is clear.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Checked | The branch is taken when Negative is clear. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `10` | Relative | `BPL $relative` | 2 | 2 | Branch taken: +1 cycle; page boundary crossed: +2 cycles total |

### BRK — Software interrupt

Pushes the address following `BRK`'s padding byte and a status byte, sets interrupt disable, and loads the IRQ/BRK vector from `$FFFE-$FFFF`. The pushed status byte has `B` set to 1.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `P` | Processor status | Pushed | Writes a status byte to the stack with `B` set to 1. |
| `I` | Interrupt disable | Set | Set after the status byte is pushed, preventing maskable IRQ recognition. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `00` | Implied | `BRK` | 1 | 7 |  |

### BVC — Branch if overflow clear

Branches when the overflow flag is clear.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `V` | Overflow | Checked | The branch is taken when Overflow is clear. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `50` | Relative | `BVC $relative` | 2 | 2 | Branch taken: +1 cycle; page boundary crossed: +2 cycles total |

### BVS — Branch if overflow set

Branches when the overflow flag is set.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `V` | Overflow | Checked | The branch is taken when Overflow is set. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `70` | Relative | `BVS $relative` | 2 | 2 | Branch taken: +1 cycle; page boundary crossed: +2 cycles total |

### CLC — Clear carry

Clears the carry flag.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `C` | Carry | Cleared | Forced to zero. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `18` | Implied | `CLC` | 1 | 2 |  |

### CLD — Clear decimal mode

Clears the decimal flag, selecting binary arithmetic for ADC and SBC.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `D` | Decimal mode | Cleared | Forced to zero, selecting binary arithmetic for ADC and SBC. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `D8` | Implied | `CLD` | 1 | 2 |  |

### CLI — Clear interrupt disable

Clears the interrupt-disable flag, allowing an asserted IRQ line to be serviced at an instruction boundary.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `I` | Interrupt disable | Cleared | Forced to zero, allowing maskable IRQ recognition. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `58` | Implied | `CLI` | 1 | 2 |  |

### CLV — Clear overflow

Clears the overflow flag.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `V` | Overflow | Cleared | Forced to zero. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `B8` | Implied | `CLV` | 1 | 2 |  |

### CMP — Compare accumulator

Subtracts the operand from the accumulator for flag purposes without storing the result. Carry means `A >= operand`; zero means equality.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | Set when bit 7 of the 8-bit result of A minus the operand is set; otherwise cleared. |
| `Z` | Zero | Set or cleared | Set when A equals the operand; otherwise cleared. |
| `C` | Carry | Set or cleared | Set when A is greater than or equal to the operand; cleared when a borrow is required. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `C9` | Immediate | `CMP #$nn` | 2 | 2 |  |
| `C5` | Zero page | `CMP $nn` | 2 | 3 |  |
| `D5` | Zero page,X | `CMP $nn,X` | 2 | 4 |  |
| `CD` | Absolute | `CMP $nnnn` | 3 | 4 |  |
| `DD` | Absolute,X | `CMP $nnnn,X` | 3 | 4 | Page boundary crossed: +1 cycle |
| `D9` | Absolute,Y | `CMP $nnnn,Y` | 3 | 4 | Page boundary crossed: +1 cycle |
| `C1` | Indexed indirect | `CMP ($nn,X)` | 2 | 6 |  |
| `D1` | Indirect indexed | `CMP ($nn),Y` | 2 | 5 | Page boundary crossed: +1 cycle |

### CPX — Compare X register

Compares X with the operand without changing X. Carry means `X >= operand`; zero means equality.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | Set when bit 7 of the 8-bit result of X minus the operand is set; otherwise cleared. |
| `Z` | Zero | Set or cleared | Set when X equals the operand; otherwise cleared. |
| `C` | Carry | Set or cleared | Set when X is greater than or equal to the operand; cleared when a borrow is required. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `E0` | Immediate | `CPX #$nn` | 2 | 2 |  |
| `E4` | Zero page | `CPX $nn` | 2 | 3 |  |
| `EC` | Absolute | `CPX $nnnn` | 3 | 4 |  |

### CPY — Compare Y register

Compares Y with the operand without changing Y. Carry means `Y >= operand`; zero means equality.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | Set when bit 7 of the 8-bit result of Y minus the operand is set; otherwise cleared. |
| `Z` | Zero | Set or cleared | Set when Y equals the operand; otherwise cleared. |
| `C` | Carry | Set or cleared | Set when Y is greater than or equal to the operand; cleared when a borrow is required. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `C0` | Immediate | `CPY #$nn` | 2 | 2 |  |
| `C4` | Zero page | `CPY $nn` | 2 | 3 |  |
| `CC` | Absolute | `CPY $nnnn` | 3 | 4 |  |

### DEC — Decrement memory

Subtracts one from a memory byte and updates zero and negative.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | Set when the decremented memory value has bit 7 set; otherwise cleared. |
| `Z` | Zero | Set or cleared | Set when the decremented memory value is zero; otherwise cleared. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `C6` | Zero page | `DEC $nn` | 2 | 5 |  |
| `D6` | Zero page,X | `DEC $nn,X` | 2 | 6 |  |
| `CE` | Absolute | `DEC $nnnn` | 3 | 6 |  |
| `DE` | Absolute,X | `DEC $nnnn,X` | 3 | 7 |  |

### DEX — Decrement X

Subtracts one from X and updates zero and negative.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | Set when the new X value has bit 7 set; otherwise cleared. |
| `Z` | Zero | Set or cleared | Set when the new X value is zero; otherwise cleared. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `CA` | Implied | `DEX` | 1 | 2 |  |

### DEY — Decrement Y

Subtracts one from Y and updates zero and negative.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | Set when the new Y value has bit 7 set; otherwise cleared. |
| `Z` | Zero | Set or cleared | Set when the new Y value is zero; otherwise cleared. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `88` | Implied | `DEY` | 1 | 2 |  |

### EOR — Exclusive OR

Exclusive-ORs the operand with the accumulator and updates zero and negative.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | Set when the new A value has bit 7 set; otherwise cleared. |
| `Z` | Zero | Set or cleared | Set when the new A value is zero; otherwise cleared. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `49` | Immediate | `EOR #$nn` | 2 | 2 |  |
| `45` | Zero page | `EOR $nn` | 2 | 3 |  |
| `55` | Zero page,X | `EOR $nn,X` | 2 | 4 |  |
| `4D` | Absolute | `EOR $nnnn` | 3 | 4 |  |
| `5D` | Absolute,X | `EOR $nnnn,X` | 3 | 4 | Page boundary crossed: +1 cycle |
| `59` | Absolute,Y | `EOR $nnnn,Y` | 3 | 4 | Page boundary crossed: +1 cycle |
| `41` | Indexed indirect | `EOR ($nn,X)` | 2 | 6 |  |
| `51` | Indirect indexed | `EOR ($nn),Y` | 2 | 5 | Page boundary crossed: +1 cycle |

### INC — Increment memory

Adds one to a memory byte and updates zero and negative.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | Set when the incremented memory value has bit 7 set; otherwise cleared. |
| `Z` | Zero | Set or cleared | Set when the incremented memory value is zero; otherwise cleared. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `E6` | Zero page | `INC $nn` | 2 | 5 |  |
| `F6` | Zero page,X | `INC $nn,X` | 2 | 6 |  |
| `EE` | Absolute | `INC $nnnn` | 3 | 6 |  |
| `FE` | Absolute,X | `INC $nnnn,X` | 3 | 7 |  |

### INX — Increment X

Adds one to X and updates zero and negative.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | Set when the new X value has bit 7 set; otherwise cleared. |
| `Z` | Zero | Set or cleared | Set when the new X value is zero; otherwise cleared. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `E8` | Implied | `INX` | 1 | 2 |  |

### INY — Increment Y

Adds one to Y and updates zero and negative.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | Set when the new Y value has bit 7 set; otherwise cleared. |
| `Z` | Zero | Set or cleared | Set when the new Y value is zero; otherwise cleared. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `C8` | Implied | `INY` | 1 | 2 |  |

### JMP — Jump

Loads the program counter with the target address. Indirect JMP uses the 6502 page-wrap behavior when the pointer ends in `$FF`.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| — | — | — | No status flags are affected. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `4C` | Absolute | `JMP $nnnn` | 3 | 3 |  |
| `6C` | Indirect | `JMP ($nnnn)` | 3 | 5 |  |

### JSR — Jump to subroutine

Pushes the address immediately before the next instruction, then jumps to the absolute target. RTS returns to the following instruction.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| — | — | — | No status flags are affected. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `20` | Absolute | `JSR $nnnn` | 3 | 6 |  |

### LDA — Load accumulator

Loads the operand into A and updates zero and negative.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | Set when the value loaded into A has bit 7 set; otherwise cleared. |
| `Z` | Zero | Set or cleared | Set when the value loaded into A is zero; otherwise cleared. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `A9` | Immediate | `LDA #$nn` | 2 | 2 |  |
| `A5` | Zero page | `LDA $nn` | 2 | 3 |  |
| `B5` | Zero page,X | `LDA $nn,X` | 2 | 4 |  |
| `AD` | Absolute | `LDA $nnnn` | 3 | 4 |  |
| `BD` | Absolute,X | `LDA $nnnn,X` | 3 | 4 | Page boundary crossed: +1 cycle |
| `B9` | Absolute,Y | `LDA $nnnn,Y` | 3 | 4 | Page boundary crossed: +1 cycle |
| `A1` | Indexed indirect | `LDA ($nn,X)` | 2 | 6 |  |
| `B1` | Indirect indexed | `LDA ($nn),Y` | 2 | 5 | Page boundary crossed: +1 cycle |

### LDX — Load X

Loads the operand into X and updates zero and negative.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | Set when the value loaded into X has bit 7 set; otherwise cleared. |
| `Z` | Zero | Set or cleared | Set when the value loaded into X is zero; otherwise cleared. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `A2` | Immediate | `LDX #$nn` | 2 | 2 |  |
| `A6` | Zero page | `LDX $nn` | 2 | 3 |  |
| `B6` | Zero page,Y | `LDX $nn,Y` | 2 | 4 |  |
| `AE` | Absolute | `LDX $nnnn` | 3 | 4 |  |
| `BE` | Absolute,Y | `LDX $nnnn,Y` | 3 | 4 | Page boundary crossed: +1 cycle |

### LDY — Load Y

Loads the operand into Y and updates zero and negative.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | Set when the value loaded into Y has bit 7 set; otherwise cleared. |
| `Z` | Zero | Set or cleared | Set when the value loaded into Y is zero; otherwise cleared. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `A0` | Immediate | `LDY #$nn` | 2 | 2 |  |
| `A4` | Zero page | `LDY $nn` | 2 | 3 |  |
| `B4` | Zero page,X | `LDY $nn,X` | 2 | 4 |  |
| `AC` | Absolute | `LDY $nnnn` | 3 | 4 |  |
| `BC` | Absolute,X | `LDY $nnnn,X` | 3 | 4 | Page boundary crossed: +1 cycle |

### LSR — Logical shift right

Shifts the accumulator or memory right by one bit. Bit 0 moves into carry, bit 7 becomes zero, and zero/negative are updated.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Cleared | Always cleared because the shift inserts zero into result bit 7. |
| `Z` | Zero | Set or cleared | Set when the shifted result is zero; otherwise cleared. |
| `C` | Carry | Set or cleared | Receives bit 0 of the original value. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `4A` | Accumulator | `LSR A` | 1 | 2 |  |
| `46` | Zero page | `LSR $nn` | 2 | 5 |  |
| `56` | Zero page,X | `LSR $nn,X` | 2 | 6 |  |
| `4E` | Absolute | `LSR $nnnn` | 3 | 6 |  |
| `5E` | Absolute,X | `LSR $nnnn,X` | 3 | 7 |  |

### NOP — No operation

Performs no state-changing operation.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| — | — | — | No status flags are affected. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `EA` | Implied | `NOP` | 1 | 2 |  |

### ORA — Logical inclusive OR

ORs the operand with the accumulator and updates zero and negative.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | Set when the new A value has bit 7 set; otherwise cleared. |
| `Z` | Zero | Set or cleared | Set when the new A value is zero; otherwise cleared. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `09` | Immediate | `ORA #$nn` | 2 | 2 |  |
| `05` | Zero page | `ORA $nn` | 2 | 3 |  |
| `15` | Zero page,X | `ORA $nn,X` | 2 | 4 |  |
| `0D` | Absolute | `ORA $nnnn` | 3 | 4 |  |
| `1D` | Absolute,X | `ORA $nnnn,X` | 3 | 4 | Page boundary crossed: +1 cycle |
| `19` | Absolute,Y | `ORA $nnnn,Y` | 3 | 4 | Page boundary crossed: +1 cycle |
| `01` | Indexed indirect | `ORA ($nn,X)` | 2 | 6 |  |
| `11` | Indirect indexed | `ORA ($nn),Y` | 2 | 5 | Page boundary crossed: +1 cycle |

### PHA — Push accumulator

Pushes A onto the hardware stack and decrements the stack pointer.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| — | — | — | No status flags are affected. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `48` | Implied | `PHA` | 1 | 3 |  |

### PHP — Push processor status

Pushes a copy of the processor status to the stack with `B` set to 1 in the pushed byte.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `P` | Processor status | Pushed | Writes a status byte to the stack with `B` set to 1. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `08` | Implied | `PHP` | 1 | 3 |  |

### PLA — Pull accumulator

Pulls A from the hardware stack and updates zero and negative.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | Set when the value pulled into A has bit 7 set; otherwise cleared. |
| `Z` | Zero | Set or cleared | Set when the value pulled into A is zero; otherwise cleared. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `68` | Implied | `PLA` | 1 | 4 |  |

### PLP — Pull processor status

Pulls the status register from the hardware stack.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Restored | Loaded from bit 7 of the status byte pulled from the stack. |
| `V` | Overflow | Restored | Loaded from bit 6 of the status byte pulled from the stack. |
| `D` | Decimal mode | Restored | Loaded from bit 3 of the status byte pulled from the stack. |
| `I` | Interrupt disable | Restored | Loaded from bit 2 of the status byte pulled from the stack. |
| `Z` | Zero | Restored | Loaded from bit 1 of the status byte pulled from the stack. |
| `C` | Carry | Restored | Loaded from bit 0 of the status byte pulled from the stack. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `28` | Implied | `PLP` | 1 | 4 |  |

### ROL — Rotate left

Rotates the accumulator or memory left through carry. Old bit 7 enters carry and the old carry enters bit 0.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | Copies bit 7 of the rotated result. |
| `Z` | Zero | Set or cleared | Set when the rotated result is zero; otherwise cleared. |
| `C` | Carry | Checked, then set or cleared | The old Carry enters result bit 0; the original bit 7 becomes the new Carry. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `2A` | Accumulator | `ROL A` | 1 | 2 |  |
| `26` | Zero page | `ROL $nn` | 2 | 5 |  |
| `36` | Zero page,X | `ROL $nn,X` | 2 | 6 |  |
| `2E` | Absolute | `ROL $nnnn` | 3 | 6 |  |
| `3E` | Absolute,X | `ROL $nnnn,X` | 3 | 7 |  |

### ROR — Rotate right

Rotates the accumulator or memory right through carry. Old bit 0 enters carry and the old carry enters bit 7.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | Copies bit 7 of the rotated result. |
| `Z` | Zero | Set or cleared | Set when the rotated result is zero; otherwise cleared. |
| `C` | Carry | Checked, then set or cleared | The old Carry enters result bit 7; the original bit 0 becomes the new Carry. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `6A` | Accumulator | `ROR A` | 1 | 2 |  |
| `66` | Zero page | `ROR $nn` | 2 | 5 |  |
| `76` | Zero page,X | `ROR $nn,X` | 2 | 6 |  |
| `6E` | Absolute | `ROR $nnnn` | 3 | 6 |  |
| `7E` | Absolute,X | `ROR $nnnn,X` | 3 | 7 |  |

### RTI — Return from interrupt

Pulls status and the program counter from the stack, resuming the interrupted program.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Restored | Loaded from bit 7 of the status byte pulled from the stack. |
| `V` | Overflow | Restored | Loaded from bit 6 of the status byte pulled from the stack. |
| `D` | Decimal mode | Restored | Loaded from bit 3 of the status byte pulled from the stack. |
| `I` | Interrupt disable | Restored | Loaded from bit 2 of the status byte pulled from the stack. |
| `Z` | Zero | Restored | Loaded from bit 1 of the status byte pulled from the stack. |
| `C` | Carry | Restored | Loaded from bit 0 of the status byte pulled from the stack. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `40` | Implied | `RTI` | 1 | 6 |  |

### RTS — Return from subroutine

Pulls the saved address from the stack, adds one, and resumes after the corresponding JSR.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| — | — | — | No status flags are affected. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `60` | Implied | `RTS` | 1 | 6 |  |

### SBC — Subtract with carry

Subtracts the operand and inverse carry from the accumulator. It updates carry, zero, negative, and overflow; decimal mode uses BCD arithmetic.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | In binary mode, copies bit 7 of the 8-bit result. In decimal mode it comes from an intermediate binary result. |
| `V` | Overflow | Set or cleared | In binary mode, set when A and the operand have different signs and the result's sign differs from A. Decimal-mode results come from an intermediate binary result. |
| `D` | Decimal mode | Checked | Selects packed-BCD arithmetic when set and binary arithmetic when clear. |
| `Z` | Zero | Set or cleared | In binary mode, set when the 8-bit result is zero. In decimal mode it comes from an intermediate binary result. |
| `C` | Carry | Checked, then set or cleared | A set old Carry means no incoming borrow; a clear old Carry subtracts one extra. The new Carry is set when no borrow is required. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `E9` | Immediate | `SBC #$nn` | 2 | 2 |  |
| `E5` | Zero page | `SBC $nn` | 2 | 3 |  |
| `F5` | Zero page,X | `SBC $nn,X` | 2 | 4 |  |
| `ED` | Absolute | `SBC $nnnn` | 3 | 4 |  |
| `FD` | Absolute,X | `SBC $nnnn,X` | 3 | 4 | Page boundary crossed: +1 cycle |
| `F9` | Absolute,Y | `SBC $nnnn,Y` | 3 | 4 | Page boundary crossed: +1 cycle |
| `E1` | Indexed indirect | `SBC ($nn,X)` | 2 | 6 |  |
| `F1` | Indirect indexed | `SBC ($nn),Y` | 2 | 5 | Page boundary crossed: +1 cycle |

### SEC — Set carry

Sets the carry flag.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `C` | Carry | Set | Forced to one. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `38` | Implied | `SEC` | 1 | 2 |  |

### SED — Set decimal mode

Sets the decimal flag, selecting BCD arithmetic for ADC and SBC.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `D` | Decimal mode | Set | Forced to one, selecting packed-BCD arithmetic for ADC and SBC. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `F8` | Implied | `SED` | 1 | 2 |  |

### SEI — Set interrupt disable

Sets the interrupt-disable flag so an asserted IRQ line is not serviced. NMI is unaffected.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `I` | Interrupt disable | Set | Forced to one, preventing maskable IRQ recognition. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `78` | Implied | `SEI` | 1 | 2 |  |

### STA — Store accumulator

Stores A in memory. Store instructions have fixed timing even when indexed addressing crosses a page.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| — | — | — | No status flags are affected. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `85` | Zero page | `STA $nn` | 2 | 3 |  |
| `95` | Zero page,X | `STA $nn,X` | 2 | 4 |  |
| `8D` | Absolute | `STA $nnnn` | 3 | 4 |  |
| `9D` | Absolute,X | `STA $nnnn,X` | 3 | 5 |  |
| `99` | Absolute,Y | `STA $nnnn,Y` | 3 | 5 |  |
| `81` | Indexed indirect | `STA ($nn,X)` | 2 | 6 |  |
| `91` | Indirect indexed | `STA ($nn),Y` | 2 | 6 |  |

### STX — Store X

Stores X in memory.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| — | — | — | No status flags are affected. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `86` | Zero page | `STX $nn` | 2 | 3 |  |
| `96` | Zero page,Y | `STX $nn,Y` | 2 | 4 |  |
| `8E` | Absolute | `STX $nnnn` | 3 | 4 |  |

### STY — Store Y

Stores Y in memory.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| — | — | — | No status flags are affected. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `84` | Zero page | `STY $nn` | 2 | 3 |  |
| `94` | Zero page,X | `STY $nn,X` | 2 | 4 |  |
| `8C` | Absolute | `STY $nnnn` | 3 | 4 |  |

### TAX — Transfer accumulator to X

Copies A into X and updates zero and negative.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | Set when the new X value has bit 7 set; otherwise cleared. |
| `Z` | Zero | Set or cleared | Set when the new X value is zero; otherwise cleared. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `AA` | Implied | `TAX` | 1 | 2 |  |

### TAY — Transfer accumulator to Y

Copies A into Y and updates zero and negative.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | Set when the new Y value has bit 7 set; otherwise cleared. |
| `Z` | Zero | Set or cleared | Set when the new Y value is zero; otherwise cleared. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `A8` | Implied | `TAY` | 1 | 2 |  |

### TSX — Transfer stack pointer to X

Copies the stack pointer into X and updates zero and negative.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | Set when the new X value has bit 7 set; otherwise cleared. |
| `Z` | Zero | Set or cleared | Set when the new X value is zero; otherwise cleared. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `BA` | Implied | `TSX` | 1 | 2 |  |

### TXA — Transfer X to accumulator

Copies X into A and updates zero and negative.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | Set when the new A value has bit 7 set; otherwise cleared. |
| `Z` | Zero | Set or cleared | Set when the new A value is zero; otherwise cleared. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `8A` | Implied | `TXA` | 1 | 2 |  |

### TXS — Transfer X to stack pointer

Copies X into the stack pointer.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| — | — | — | No status flags are affected. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `9A` | Implied | `TXS` | 1 | 2 |  |

### TYA — Transfer Y to accumulator

Copies Y into A and updates zero and negative.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | Set when the new A value has bit 7 set; otherwise cleared. |
| `Z` | Zero | Set or cleared | Set when the new A value is zero; otherwise cleared. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `98` | Implied | `TYA` | 1 | 2 |  |

## Undocumented Instructions

### Undocumented instruction index

| Range | Instructions |
|---|---|
| A–L | [ALR](#alr--and-then-logical-shift-right-undocumented) · [ANC](#anc--and-and-copy-negative-to-carry-undocumented) · [ARR](#arr--and-then-rotate-right-undocumented) · [DCP](#dcp--decrement-then-compare-undocumented) · [DOP](#dop--two-byte-no-operation-undocumented) · [ISC](#isc--increment-then-subtract-with-carry-undocumented) · [LAX](#lax--load-accumulator-and-x-undocumented) |
| N–X | [NOP*](#nop--single-byte-no-operation-opcodes-undocumented) · [RLA](#rla--rotate-left-then-and-undocumented) · [RRA](#rra--rotate-right-then-add-with-carry-undocumented) · [SAX](#sax--store-a-and-x-undocumented) · [SBC*](#sbc--immediate-sbc-alias-undocumented) · [SKB](#skb--skip-byte-through-zero-page-undocumented) · [SKW](#skw--skip-byte-through-zero-pagex-undocumented) · [SLO](#slo--shift-left-then-or-undocumented) · [SRE](#sre--shift-right-then-exclusive-or-undocumented) · [TOP](#top--three-byte-no-operation-undocumented) · [XAA](#xaa--transfer-x-and-immediate-to-accumulator-undocumented) |

These opcodes were not specified by MOS Technology. Names are community
conventions, some behaviours depend on the chip revision and electrical
conditions, and they are not guaranteed to behave identically on every 6502.

### ALR — AND then logical shift right (Undocumented)

ANDs an immediate value with the accumulator, shifts the result right, and stores it in the accumulator. Bit 0 moves into carry.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Cleared | Always cleared because the final right shift inserts zero into A bit 7. |
| `Z` | Zero | Set or cleared | Set when the final A value is zero; otherwise cleared. |
| `C` | Carry | Set or cleared | Receives bit 0 of the intermediate A AND operand value. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `4B` | Immediate | `ALR #$nn` | 2 | 2 |  |

### ANC — AND and copy negative to carry (Undocumented)

ANDs an immediate value with the accumulator, then copies result bit 7 into both the negative and carry flags.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | Copies bit 7 of the final A value. |
| `Z` | Zero | Set or cleared | Set when the final A value is zero; otherwise cleared. |
| `C` | Carry | Set or cleared | Copies bit 7 of the final A value, so it matches Negative. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `0B`, `2B` | Immediate | `ANC #$nn` | 2 | 2 |  |

### ARR — AND then rotate right (Undocumented)

ANDs an immediate value with the accumulator, then rotates right through carry. In binary mode, carry and overflow are derived from the rotated value as described below. Decimal-mode behaviour is unusual and should not be treated as portable.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | Copies bit 7 of the final A value. |
| `V` | Overflow | Set or cleared | In binary mode, set when result bits 6 and 5 differ; decimal-mode behaviour can be unpredictable. |
| `D` | Decimal mode | Checked | Affects the instruction's arithmetic when set. |
| `Z` | Zero | Set or cleared | Set when the final A value is zero; otherwise cleared. |
| `C` | Carry | Checked, then set or cleared | The old Carry enters result bit 7. In binary mode, result bit 6 becomes the new Carry. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `6B` | Immediate | `ARR #$nn` | 2 | 2 |  |

### DCP — Decrement then compare (Undocumented)

Decrements memory, then compares the new value with the accumulator as CMP would.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | Set when bit 7 of A minus the decremented memory value is set; otherwise cleared. |
| `Z` | Zero | Set or cleared | Set when A equals the decremented memory value; otherwise cleared. |
| `C` | Carry | Set or cleared | Set when A is greater than or equal to the decremented memory value; otherwise cleared. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `C7` | Zero page | `DCP $nn` | 2 | 5 |  |
| `D7` | Zero page,X | `DCP $nn,X` | 2 | 6 |  |
| `CF` | Absolute | `DCP $nnnn` | 3 | 6 |  |
| `DF` | Absolute,X | `DCP $nnnn,X` | 3 | 7 |  |
| `DB` | Absolute,Y | `DCP $nnnn,Y` | 3 | 7 |  |
| `C3` | Indexed indirect | `DCP ($nn,X)` | 2 | 8 |  |
| `D3` | Indirect indexed | `DCP ($nn),Y` | 2 | 8 |  |

### DOP — Two-byte no operation (Undocumented)

Consumes an immediate operand without otherwise changing CPU state.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| — | — | — | No status flags are affected. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `80`, `82`, `C2`, `E2` | Immediate | `DOP #$nn` | 2 | 2 |  |

### ISC — Increment then subtract with carry (Undocumented)

Increments memory, then subtracts the new value from the accumulator as `SBC` would.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | In binary mode, copies bit 7 of A after subtracting the incremented memory value. In decimal mode it comes from an intermediate binary result. |
| `V` | Overflow | Set or cleared | In binary mode, set when A and the incremented memory value have different signs and the result's sign differs from A. |
| `D` | Decimal mode | Checked | Selects decimal or binary subtraction. |
| `Z` | Zero | Set or cleared | In binary mode, set when A after the subtraction is zero. In decimal mode it comes from an intermediate binary result. |
| `C` | Carry | Checked, then set or cleared | Used and updated as by SBC; set after the subtraction when no borrow is required. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `E7` | Zero page | `ISC $nn` | 2 | 5 |  |
| `F7` | Zero page,X | `ISC $nn,X` | 2 | 6 |  |
| `EF` | Absolute | `ISC $nnnn` | 3 | 6 |  |
| `FF` | Absolute,X | `ISC $nnnn,X` | 3 | 7 |  |
| `FB` | Absolute,Y | `ISC $nnnn,Y` | 3 | 7 |  |
| `E3` | Indexed indirect | `ISC ($nn,X)` | 2 | 8 |  |
| `F3` | Indirect indexed | `ISC ($nn),Y` | 2 | 8 |  |

### LAX — Load accumulator and X (Undocumented)

Loads the same operand into both A and X, then updates zero and negative.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | Set when the value loaded into A and X has bit 7 set; otherwise cleared. |
| `Z` | Zero | Set or cleared | Set when the value loaded into A and X is zero; otherwise cleared. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `A7` | Zero page | `LAX $nn` | 2 | 3 |  |
| `B7` | Zero page,Y | `LAX $nn,Y` | 2 | 4 |  |
| `AF` | Absolute | `LAX $nnnn` | 3 | 4 |  |
| `BF` | Absolute,Y | `LAX $nnnn,Y` | 3 | 4 | Page boundary crossed: +1 cycle |
| `A3` | Indexed indirect | `LAX ($nn,X)` | 2 | 6 |  |
| `B3` | Indirect indexed | `LAX ($nn),Y` | 2 | 5 | Page boundary crossed: +1 cycle |

### NOP* — Single-byte no operation opcodes (Undocumented)

These opcodes behave like the documented `NOP` on the 6502.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| — | — | — | No status flags are affected. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `1A`, `3A`, `5A`, `7A`, `DA`, `FA` | Implied | `NOP*` | 1 | 2 |  |

### RLA — Rotate left then AND (Undocumented)

Rotates memory left through carry, then ANDs the new memory value into the accumulator.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | Copies bit 7 of A after it is ANDed with the rotated memory value. |
| `Z` | Zero | Set or cleared | Set when the final A value is zero; otherwise cleared. |
| `C` | Carry | Checked, then set or cleared | The old Carry enters memory bit 0; the original memory bit 7 becomes the new Carry. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `27` | Zero page | `RLA $nn` | 2 | 5 |  |
| `37` | Zero page,X | `RLA $nn,X` | 2 | 6 |  |
| `2F` | Absolute | `RLA $nnnn` | 3 | 6 |  |
| `3F` | Absolute,X | `RLA $nnnn,X` | 3 | 7 |  |
| `3B` | Absolute,Y | `RLA $nnnn,Y` | 3 | 7 |  |
| `23` | Indexed indirect | `RLA ($nn,X)` | 2 | 8 |  |
| `33` | Indirect indexed | `RLA ($nn),Y` | 2 | 8 |  |

### RRA — Rotate right then add with carry (Undocumented)

Rotates memory right through carry, then adds the new memory value to the accumulator.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | In binary mode, copies bit 7 of A after adding the rotated memory value. In decimal mode it comes from an intermediate binary result. |
| `V` | Overflow | Set or cleared | In binary mode, set when A and the rotated memory value have the same sign and the result has the opposite sign. |
| `D` | Decimal mode | Checked | Selects decimal or binary addition. |
| `Z` | Zero | Set or cleared | In binary mode, set when A after the addition is zero. In decimal mode it comes from an intermediate binary result. |
| `C` | Carry | Checked, then set or cleared | The old Carry enters memory bit 7; the original memory bit 0 becomes ADC's carry-in, and ADC produces the final Carry. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `67` | Zero page | `RRA $nn` | 2 | 5 |  |
| `77` | Zero page,X | `RRA $nn,X` | 2 | 6 |  |
| `6F` | Absolute | `RRA $nnnn` | 3 | 6 |  |
| `7F` | Absolute,X | `RRA $nnnn,X` | 3 | 7 |  |
| `7B` | Absolute,Y | `RRA $nnnn,Y` | 3 | 7 |  |
| `63` | Indexed indirect | `RRA ($nn,X)` | 2 | 8 |  |
| `73` | Indirect indexed | `RRA ($nn),Y` | 2 | 8 |  |

### SAX — Store A AND X (Undocumented)

Stores `A AND X` in memory without changing either register.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| — | — | — | No status flags are affected. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `87` | Zero page | `SAX $nn` | 2 | 3 |  |
| `97` | Zero page,Y | `SAX $nn,Y` | 2 | 4 |  |
| `8F` | Absolute | `SAX $nnnn` | 3 | 4 |  |
| `83` | Indexed indirect | `SAX ($nn,X)` | 2 | 6 |  |

### SBC* — Immediate SBC alias (Undocumented)

Performs the same immediate subtract-with-carry operation as opcode `E9`.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | In binary mode, copies bit 7 of the 8-bit result. In decimal mode it comes from an intermediate binary result. |
| `V` | Overflow | Set or cleared | In binary mode, set when A and the operand have different signs and the result's sign differs from A. Decimal-mode results come from an intermediate binary result. |
| `D` | Decimal mode | Checked | Selects packed-BCD arithmetic when set and binary arithmetic when clear. |
| `Z` | Zero | Set or cleared | In binary mode, set when the 8-bit result is zero. In decimal mode it comes from an intermediate binary result. |
| `C` | Carry | Checked, then set or cleared | A set old Carry means no incoming borrow; a clear old Carry subtracts one extra. The new Carry is set when no borrow is required. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `EB` | Immediate | `SBC* #$nn` | 2 | 2 |  |

### SKB — Skip byte through zero page (Undocumented)

Reads a zero-page operand and otherwise behaves as a no operation.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| — | — | — | No status flags are affected. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `04`, `44`, `64` | Zero page | `SKB $nn` | 2 | 3 |  |

### SKW — Skip byte through zero page,X (Undocumented)

Reads a zero-page,X operand and otherwise behaves as a no operation.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| — | — | — | No status flags are affected. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `14`, `34`, `54`, `74`, `D4`, `F4` | Zero page,X | `SKW $nn,X` | 2 | 4 |  |

### SLO — Shift left then OR (Undocumented)

Shifts memory left, then ORs the new memory value into the accumulator.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | Copies bit 7 of A after it is ORed with the shifted memory value. |
| `Z` | Zero | Set or cleared | Set when the final A value is zero; otherwise cleared. |
| `C` | Carry | Set or cleared | Receives bit 7 of the original memory value. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `07` | Zero page | `SLO $nn` | 2 | 5 |  |
| `17` | Zero page,X | `SLO $nn,X` | 2 | 6 |  |
| `0F` | Absolute | `SLO $nnnn` | 3 | 6 |  |
| `1F` | Absolute,X | `SLO $nnnn,X` | 3 | 7 |  |
| `1B` | Absolute,Y | `SLO $nnnn,Y` | 3 | 7 |  |
| `03` | Indexed indirect | `SLO ($nn,X)` | 2 | 8 |  |
| `13` | Indirect indexed | `SLO ($nn),Y` | 2 | 8 |  |

### SRE — Shift right then exclusive OR (Undocumented)

Shifts memory right, then exclusive-ORs the new memory value into the accumulator.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | Copies bit 7 of A after it is XORed with the shifted memory value. |
| `Z` | Zero | Set or cleared | Set when the final A value is zero; otherwise cleared. |
| `C` | Carry | Set or cleared | Receives bit 0 of the original memory value. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `47` | Zero page | `SRE $nn` | 2 | 5 |  |
| `57` | Zero page,X | `SRE $nn,X` | 2 | 6 |  |
| `4F` | Absolute | `SRE $nnnn` | 3 | 6 |  |
| `5F` | Absolute,X | `SRE $nnnn,X` | 3 | 7 |  |
| `5B` | Absolute,Y | `SRE $nnnn,Y` | 3 | 7 |  |
| `43` | Indexed indirect | `SRE ($nn,X)` | 2 | 8 |  |
| `53` | Indirect indexed | `SRE ($nn),Y` | 2 | 8 |  |

### TOP — Three-byte no operation (Undocumented)

Reads an absolute operand and otherwise behaves as a no operation. The absolute,X forms take one extra cycle when indexing crosses a page.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| — | — | — | No status flags are affected. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `0C` | Absolute | `TOP $nnnn` | 3 | 4 |  |
| `1C`, `3C`, `5C`, `7C`, `DC`, `FC` | Absolute,X | `TOP $nnnn,X` | 3 | 4 | Page boundary crossed: +1 cycle |

### XAA — Transfer X AND immediate to accumulator (Undocumented)

On the 6502 this instruction combines X, an immediate operand, and an internal bus value, then stores the result in A. Its result is electrically unstable and can vary with chip revision, temperature, and supply voltage; no single deterministic formula is reliable.

#### Status register

| Symbol | Name | Action | Description |
|:---:|---|---|---|
| `N` | Negative | Set or cleared | Set when the value placed in A has bit 7 set; otherwise cleared. |
| `Z` | Zero | Set or cleared | Set when the value placed in A is zero; otherwise cleared. |

#### Opcodes

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Cycle penalty |
|:---:|---|---|---:|---:|---|
| `8B` | Immediate | `XAA #$nn` | 2 | 2 |  |
