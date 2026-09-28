# 6502 Instruction Set Manual

This manual documents the 6502 instruction set. It covers every
documented instruction and, in a separate appendix, commonly encountered
undocumented opcodes.

Cycle counts are the minimum hardware cycle counts. A value in **Extra cycles**
describes the condition that adds cycles:

- **Page crossed: +1** — indexed address calculation moved into another 256-byte page.
- **Branch taken: +1; page crossed: +2 total** — a taken branch costs one extra cycle, or two when its destination is on another page.
- **None** — the instruction has a fixed cycle count.

## Processor Status (P) register

The Processor Status register, abbreviated **P**, contains six persistent
flags. Status bytes are conventionally displayed as `N V 1 B D I Z C`, from
bit 7 down to bit 0:

| Bit | Symbol | Name | Set (`1`) when |
|---:|:---:|---|---|
| 7 | `N` | Negative | The instruction-defined result has bit 7 set. |
| 6 | `V` | Overflow | A signed addition or subtraction cannot be represented in the range -128 to 127, or `BIT` copied a set operand bit 6. |
| 5 | `1` | Reserved | A status byte written to the stack or read on the data bus normally has this bit set. It is not a writable flag. |
| 4 | `B` | Break marker | The stacked status came from `BRK` or `PHP`. It is not a flag stored in P. |
| 3 | `D` | Decimal mode | `ADC` and `SBC` perform packed-BCD arithmetic. |
| 2 | `I` | Interrupt disable | Maskable IRQ recognition is disabled. NMI is unaffected. |
| 1 | `Z` | Zero | The instruction-defined result is zero. |
| 0 | `C` | Carry | Addition produced a carry, subtraction required no borrow, or a shift/rotate moved out a one bit. |

Only `N`, `V`, `D`, `I`, `Z`, and `C` are persistent flags. The `B` symbol is
useful when examining a status byte on the stack: `BRK` and `PHP` push it as
one, while a hardware IRQ or NMI pushes it as zero. Bit 5 is pushed as one.
`PLP` and `RTI` restore the six real flags; the pulled values in bits 5 and 4
do not create reserved or break latches.

Every instruction below has a **Processor Status (P)** note split into plain-language actions:

- **Checks** means the instruction reads the flag's current value to decide what to do, but checking alone does not change it.
- **Changes** means the instruction recalculates the flag; it may end up set or clear.
- **Sets** or **clears** means the instruction forces that flag to one or zero.
- **Restores** means a persistent flag is loaded from the hardware stack.
- **No flags checked or changed** means the entire P register is left alone. Any flag not named in a note is also unchanged.

For binary arithmetic, `ADC` sets `V` when its two inputs have the same sign
and the result has the opposite sign. `SBC` sets `V` when A and the operand
have different signs and the result's sign differs from A. In decimal mode,
the 6502's `N`, `V`, and `Z` values come from intermediate binary results,
not simply from the final BCD-adjusted accumulator; use `C` as the valid
decimal carry/no-borrow indication.

### Status-flag quick reference

This table is an index; each instruction entry below gives the exact condition
for setting or clearing every affected flag. A dash means no persistent flag
changes.

| Instructions | Flags read | Flags changed |
|---|---|---|
| `ADC`, `SBC` | `C`, `D` | `N`, `V`, `Z`, `C` |
| `AND`, `EOR`, `ORA`, `LDA`, `LDX`, `LDY`, `TAX`, `TAY`, `TSX`, `TXA`, `TYA`, `DEC`, `DEX`, `DEY`, `INC`, `INX`, `INY`, `PLA` | — | `N`, `Z` |
| `ASL`, `LSR` | — | `N`, `Z`, `C` |
| `ROL`, `ROR` | `C` | `N`, `Z`, `C` |
| `BIT` | — | `N`, `V`, `Z` |
| `CMP`, `CPX`, `CPY` | — | `N`, `Z`, `C` |
| `BCC`, `BCS` | `C` | — |
| `BEQ`, `BNE` | `Z` | — |
| `BMI`, `BPL` | `N` | — |
| `BVC`, `BVS` | `V` | — |
| `CLC`, `CLD`, `CLI`, `CLV` | — | clears `C`, `D`, `I`, or `V` respectively |
| `SEC`, `SED`, `SEI` | — | sets `C`, `D`, or `I` respectively |
| `BRK` | — | sets `I`; pushes status with `B=1` and bit 5 set |
| `PHP` | — | pushes status with `B=1` and bit 5 set; no live flag changes |
| `PLP`, `RTI` | — | restores `N`, `V`, `D`, `I`, `Z`, `C` |
| `JMP`, `JSR`, `NOP`, `PHA`, `RTS`, `STA`, `STX`, `STY`, `TXS` | — | — |

## Addressing notation

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

## Instruction Set

### ADC — Add with carry

Adds the operand and the carry flag to the accumulator. It updates carry, zero, negative, and overflow; decimal mode uses BCD arithmetic.

**Processor Status (P):** Checks Carry (C), adding an extra 1 when it is set, and checks Decimal mode (D) to select binary or BCD arithmetic. Changes Carry (C), Zero (Z), Overflow (V), and Negative (N). In binary mode, Carry is set when the unsigned sum exceeds 255; Zero is set when the 8-bit result is zero; Negative copies result bit 7; Overflow is set when two operands with the same sign produce a result with the opposite sign. In decimal mode, see the decimal-flag warning above.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `69` | Immediate | `ADC #$nn` | 2 | 2 | None |
| `65` | Zero page | `ADC $nn` | 2 | 3 | None |
| `75` | Zero page,X | `ADC $nn,X` | 2 | 4 | None |
| `6D` | Absolute | `ADC $nnnn` | 3 | 4 | None |
| `7D` | Absolute,X | `ADC $nnnn,X` | 3 | 4 | Page crossed: +1 |
| `79` | Absolute,Y | `ADC $nnnn,Y` | 3 | 4 | Page crossed: +1 |
| `61` | Indexed indirect | `ADC ($nn,X)` | 2 | 6 | None |
| `71` | Indirect indexed | `ADC ($nn),Y` | 2 | 5 | Page crossed: +1 |

### AND — Logical AND

ANDs the operand with the accumulator and updates zero and negative.

**Processor Status (P):** Changes Zero (Z) and Negative (N) from the new A value. Zero is set when A is zero; Negative copies A bit 7.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `29` | Immediate | `AND #$nn` | 2 | 2 | None |
| `25` | Zero page | `AND $nn` | 2 | 3 | None |
| `35` | Zero page,X | `AND $nn,X` | 2 | 4 | None |
| `2D` | Absolute | `AND $nnnn` | 3 | 4 | None |
| `3D` | Absolute,X | `AND $nnnn,X` | 3 | 4 | Page crossed: +1 |
| `39` | Absolute,Y | `AND $nnnn,Y` | 3 | 4 | Page crossed: +1 |
| `21` | Indexed indirect | `AND ($nn,X)` | 2 | 6 | None |
| `31` | Indirect indexed | `AND ($nn),Y` | 2 | 5 | Page crossed: +1 |

### ASL — Arithmetic shift left

Shifts the accumulator or memory left by one bit. Bit 7 moves into carry and zero/negative reflect the result.

**Processor Status (P):** Changes Carry (C), Zero (Z), and Negative (N). Carry receives the old bit 7; Zero is set when the shifted result is zero; Negative copies result bit 7.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `0A` | Accumulator | `ASL A` | 1 | 2 | None |
| `06` | Zero page | `ASL $nn` | 2 | 5 | None |
| `16` | Zero page,X | `ASL $nn,X` | 2 | 6 | None |
| `0E` | Absolute | `ASL $nnnn` | 3 | 6 | None |
| `1E` | Absolute,X | `ASL $nnnn,X` | 3 | 7 | None |

### BCC — Branch if carry clear

Branches when the carry flag is clear.

**Processor Status (P):** Checks Carry (C). Branches when Carry is clear; changes no flags.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `90` | Relative | `BCC $relative` | 2 | 2 | Branch taken: +1; page crossed: +2 total |

### BCS — Branch if carry set

Branches when the carry flag is set.

**Processor Status (P):** Checks Carry (C). Branches when Carry is set; changes no flags.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `B0` | Relative | `BCS $relative` | 2 | 2 | Branch taken: +1; page crossed: +2 total |

### BEQ — Branch if equal

Branches when the zero flag is set.

**Processor Status (P):** Checks Zero (Z). Branches when Zero is set; changes no flags.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `F0` | Relative | `BEQ $relative` | 2 | 2 | Branch taken: +1; page crossed: +2 total |

### BIT — Bit test

Tests the accumulator against memory without changing either value. Zero reflects `A AND operand`; bits 7 and 6 of memory become negative and overflow.

**Processor Status (P):** Changes Zero (Z), Overflow (V), and Negative (N). Zero is set when `A AND operand` is zero; Overflow copies operand bit 6; Negative copies operand bit 7.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `24` | Zero page | `BIT $nn` | 2 | 3 | None |
| `2C` | Absolute | `BIT $nnnn` | 3 | 4 | None |

### BMI — Branch if minus

Branches when the negative flag is set.

**Processor Status (P):** Checks Negative (N). Branches when Negative is set; changes no flags.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `30` | Relative | `BMI $relative` | 2 | 2 | Branch taken: +1; page crossed: +2 total |

### BNE — Branch if not equal

Branches when the zero flag is clear.

**Processor Status (P):** Checks Zero (Z). Branches when Zero is clear; changes no flags.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `D0` | Relative | `BNE $relative` | 2 | 2 | Branch taken: +1; page crossed: +2 total |

### BPL — Branch if plus

Branches when the negative flag is clear.

**Processor Status (P):** Checks Negative (N). Branches when Negative is clear; changes no flags.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `10` | Relative | `BPL $relative` | 2 | 2 | Branch taken: +1; page crossed: +2 total |

### BRK — Software interrupt

Pushes the address following `BRK`'s padding byte and a status byte, sets interrupt disable, and loads the IRQ/BRK vector from `$FFFE-$FFFF`. The pushed status has its break marker set.

**Processor Status (P):** Sets Interrupt disable (I). The status byte pushed to the stack has `B=1` and bit 5 set; its six persistent flags have their pre-`BRK` values because Interrupt disable is set after the push. There is no live Break flag.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `00` | Implied | `BRK` | 1 | 7 | None |

### BVC — Branch if overflow clear

Branches when the overflow flag is clear.

**Processor Status (P):** Checks Overflow (V). Branches when Overflow is clear; changes no flags.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `50` | Relative | `BVC $relative` | 2 | 2 | Branch taken: +1; page crossed: +2 total |

### BVS — Branch if overflow set

Branches when the overflow flag is set.

**Processor Status (P):** Checks Overflow (V). Branches when Overflow is set; changes no flags.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `70` | Relative | `BVS $relative` | 2 | 2 | Branch taken: +1; page crossed: +2 total |

### CLC — Clear carry

Clears the carry flag.

**Processor Status (P):** Clears Carry (C). No other flags change.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `18` | Implied | `CLC` | 1 | 2 | None |

### CLD — Clear decimal mode

Clears the decimal flag, selecting binary arithmetic for ADC and SBC.

**Processor Status (P):** Clears Decimal mode (D). No other flags change.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `D8` | Implied | `CLD` | 1 | 2 | None |

### CLI — Clear interrupt disable

Clears the interrupt-disable flag, allowing an asserted IRQ line to be serviced at an instruction boundary.

**Processor Status (P):** Clears Interrupt disable (I). No other flags change.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `58` | Implied | `CLI` | 1 | 2 | None |

### CLV — Clear overflow

Clears the overflow flag.

**Processor Status (P):** Clears Overflow (V). No other flags change.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `B8` | Implied | `CLV` | 1 | 2 | None |

### CMP — Compare accumulator

Subtracts the operand from the accumulator for flag purposes without storing the result. Carry means `A >= operand`; zero means equality.

**Processor Status (P):** Changes Carry (C), Zero (Z), and Negative (N) as though `A - operand` were calculated. Carry is set when `A >= operand`; Zero is set when the 8-bit difference is zero; Negative copies difference bit 7. A itself is not changed.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `C9` | Immediate | `CMP #$nn` | 2 | 2 | None |
| `C5` | Zero page | `CMP $nn` | 2 | 3 | None |
| `D5` | Zero page,X | `CMP $nn,X` | 2 | 4 | None |
| `CD` | Absolute | `CMP $nnnn` | 3 | 4 | None |
| `DD` | Absolute,X | `CMP $nnnn,X` | 3 | 4 | Page crossed: +1 |
| `D9` | Absolute,Y | `CMP $nnnn,Y` | 3 | 4 | Page crossed: +1 |
| `C1` | Indexed indirect | `CMP ($nn,X)` | 2 | 6 | None |
| `D1` | Indirect indexed | `CMP ($nn),Y` | 2 | 5 | Page crossed: +1 |

### CPX — Compare X register

Compares X with the operand without changing X. Carry means `X >= operand`; zero means equality.

**Processor Status (P):** Changes Carry (C), Zero (Z), and Negative (N) as though `X - operand` were calculated, using the same rules as `CMP`. X itself is not changed.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `E0` | Immediate | `CPX #$nn` | 2 | 2 | None |
| `E4` | Zero page | `CPX $nn` | 2 | 3 | None |
| `EC` | Absolute | `CPX $nnnn` | 3 | 4 | None |

### CPY — Compare Y register

Compares Y with the operand without changing Y. Carry means `Y >= operand`; zero means equality.

**Processor Status (P):** Changes Carry (C), Zero (Z), and Negative (N) as though `Y - operand` were calculated, using the same rules as `CMP`. Y itself is not changed.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `C0` | Immediate | `CPY #$nn` | 2 | 2 | None |
| `C4` | Zero page | `CPY $nn` | 2 | 3 | None |
| `CC` | Absolute | `CPY $nnnn` | 3 | 4 | None |

### DEC — Decrement memory

Subtracts one from a memory byte and updates zero and negative.

**Processor Status (P):** Changes Zero (Z) and Negative (N) from the decremented byte. Zero is set when the result is zero; Negative copies result bit 7.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `C6` | Zero page | `DEC $nn` | 2 | 5 | None |
| `D6` | Zero page,X | `DEC $nn,X` | 2 | 6 | None |
| `CE` | Absolute | `DEC $nnnn` | 3 | 6 | None |
| `DE` | Absolute,X | `DEC $nnnn,X` | 3 | 7 | None |

### DEX — Decrement X

Subtracts one from X and updates zero and negative.

**Processor Status (P):** Changes Zero (Z) and Negative (N) from the new X value. Zero is set when X is zero; Negative copies X bit 7.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `CA` | Implied | `DEX` | 1 | 2 | None |

### DEY — Decrement Y

Subtracts one from Y and updates zero and negative.

**Processor Status (P):** Changes Zero (Z) and Negative (N) from the new Y value. Zero is set when Y is zero; Negative copies Y bit 7.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `88` | Implied | `DEY` | 1 | 2 | None |

### EOR — Exclusive OR

Exclusive-ORs the operand with the accumulator and updates zero and negative.

**Processor Status (P):** Changes Zero (Z) and Negative (N) from the new A value. Zero is set when A is zero; Negative copies A bit 7.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `49` | Immediate | `EOR #$nn` | 2 | 2 | None |
| `45` | Zero page | `EOR $nn` | 2 | 3 | None |
| `55` | Zero page,X | `EOR $nn,X` | 2 | 4 | None |
| `4D` | Absolute | `EOR $nnnn` | 3 | 4 | None |
| `5D` | Absolute,X | `EOR $nnnn,X` | 3 | 4 | Page crossed: +1 |
| `59` | Absolute,Y | `EOR $nnnn,Y` | 3 | 4 | Page crossed: +1 |
| `41` | Indexed indirect | `EOR ($nn,X)` | 2 | 6 | None |
| `51` | Indirect indexed | `EOR ($nn),Y` | 2 | 5 | Page crossed: +1 |

### INC — Increment memory

Adds one to a memory byte and updates zero and negative.

**Processor Status (P):** Changes Zero (Z) and Negative (N) from the incremented byte. Zero is set when the result is zero; Negative copies result bit 7.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `E6` | Zero page | `INC $nn` | 2 | 5 | None |
| `F6` | Zero page,X | `INC $nn,X` | 2 | 6 | None |
| `EE` | Absolute | `INC $nnnn` | 3 | 6 | None |
| `FE` | Absolute,X | `INC $nnnn,X` | 3 | 7 | None |

### INX — Increment X

Adds one to X and updates zero and negative.

**Processor Status (P):** Changes Zero (Z) and Negative (N) from the new X value. Zero is set when X is zero; Negative copies X bit 7.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `E8` | Implied | `INX` | 1 | 2 | None |

### INY — Increment Y

Adds one to Y and updates zero and negative.

**Processor Status (P):** Changes Zero (Z) and Negative (N) from the new Y value. Zero is set when Y is zero; Negative copies Y bit 7.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `C8` | Implied | `INY` | 1 | 2 | None |

### JMP — Jump

Loads the program counter with the target address. Indirect JMP uses the 6502 page-wrap behavior when the pointer ends in `$FF`.

**Processor Status (P):** No flags checked or changed.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `4C` | Absolute | `JMP $nnnn` | 3 | 3 | None |
| `6C` | Indirect | `JMP ($nnnn)` | 3 | 5 | None |

### JSR — Jump to subroutine

Pushes the address immediately before the next instruction, then jumps to the absolute target. RTS returns to the following instruction.

**Processor Status (P):** No flags checked or changed.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `20` | Absolute | `JSR $nnnn` | 3 | 6 | None |

### LDA — Load accumulator

Loads the operand into A and updates zero and negative.

**Processor Status (P):** Changes Zero (Z) and Negative (N) from the loaded A value. Zero is set when A is zero; Negative copies A bit 7.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `A9` | Immediate | `LDA #$nn` | 2 | 2 | None |
| `A5` | Zero page | `LDA $nn` | 2 | 3 | None |
| `B5` | Zero page,X | `LDA $nn,X` | 2 | 4 | None |
| `AD` | Absolute | `LDA $nnnn` | 3 | 4 | None |
| `BD` | Absolute,X | `LDA $nnnn,X` | 3 | 4 | Page crossed: +1 |
| `B9` | Absolute,Y | `LDA $nnnn,Y` | 3 | 4 | Page crossed: +1 |
| `A1` | Indexed indirect | `LDA ($nn,X)` | 2 | 6 | None |
| `B1` | Indirect indexed | `LDA ($nn),Y` | 2 | 5 | Page crossed: +1 |

### LDX — Load X

Loads the operand into X and updates zero and negative.

**Processor Status (P):** Changes Zero (Z) and Negative (N) from the loaded X value. Zero is set when X is zero; Negative copies X bit 7.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `A2` | Immediate | `LDX #$nn` | 2 | 2 | None |
| `A6` | Zero page | `LDX $nn` | 2 | 3 | None |
| `B6` | Zero page,Y | `LDX $nn,Y` | 2 | 4 | None |
| `AE` | Absolute | `LDX $nnnn` | 3 | 4 | None |
| `BE` | Absolute,Y | `LDX $nnnn,Y` | 3 | 4 | Page crossed: +1 |

### LDY — Load Y

Loads the operand into Y and updates zero and negative.

**Processor Status (P):** Changes Zero (Z) and Negative (N) from the loaded Y value. Zero is set when Y is zero; Negative copies Y bit 7.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `A0` | Immediate | `LDY #$nn` | 2 | 2 | None |
| `A4` | Zero page | `LDY $nn` | 2 | 3 | None |
| `B4` | Zero page,X | `LDY $nn,X` | 2 | 4 | None |
| `AC` | Absolute | `LDY $nnnn` | 3 | 4 | None |
| `BC` | Absolute,X | `LDY $nnnn,X` | 3 | 4 | Page crossed: +1 |

### LSR — Logical shift right

Shifts the accumulator or memory right by one bit. Bit 0 moves into carry, bit 7 becomes zero, and zero/negative are updated.

**Processor Status (P):** Changes Carry (C) and Zero (Z), and clears Negative (N). Carry receives the old bit 0; Zero is set when the shifted result is zero. Negative is always clear because the shift inserts zero into result bit 7.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `4A` | Accumulator | `LSR A` | 1 | 2 | None |
| `46` | Zero page | `LSR $nn` | 2 | 5 | None |
| `56` | Zero page,X | `LSR $nn,X` | 2 | 6 | None |
| `4E` | Absolute | `LSR $nnnn` | 3 | 6 | None |
| `5E` | Absolute,X | `LSR $nnnn,X` | 3 | 7 | None |

### NOP — No operation

Performs no state-changing operation.

**Processor Status (P):** No flags checked or changed.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `EA` | Implied | `NOP` | 1 | 2 | None |

### ORA — Logical inclusive OR

ORs the operand with the accumulator and updates zero and negative.

**Processor Status (P):** Changes Zero (Z) and Negative (N) from the new A value. Zero is set when A is zero; Negative copies A bit 7.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `09` | Immediate | `ORA #$nn` | 2 | 2 | None |
| `05` | Zero page | `ORA $nn` | 2 | 3 | None |
| `15` | Zero page,X | `ORA $nn,X` | 2 | 4 | None |
| `0D` | Absolute | `ORA $nnnn` | 3 | 4 | None |
| `1D` | Absolute,X | `ORA $nnnn,X` | 3 | 4 | Page crossed: +1 |
| `19` | Absolute,Y | `ORA $nnnn,Y` | 3 | 4 | Page crossed: +1 |
| `01` | Indexed indirect | `ORA ($nn,X)` | 2 | 6 | None |
| `11` | Indirect indexed | `ORA ($nn),Y` | 2 | 5 | Page crossed: +1 |

### PHA — Push accumulator

Pushes A onto the hardware stack and decrements the stack pointer.

**Processor Status (P):** No flags checked or changed.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `48` | Implied | `PHA` | 1 | 3 | None |

### PHP — Push processor status

Pushes a status byte with the break marker and reserved bit set.

**Processor Status (P):** Changes no live flags. In the copy pushed to the stack, `B=1` and bit 5 is set; every other bit copies its corresponding persistent flag.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `08` | Implied | `PHP` | 1 | 3 | None |

### PLA — Pull accumulator

Pulls A from the hardware stack and updates zero and negative.

**Processor Status (P):** Changes Zero (Z) and Negative (N) from the value pulled into A. Zero is set when A is zero; Negative copies A bit 7.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `68` | Implied | `PLA` | 1 | 4 | None |

### PLP — Pull processor status

Pulls the status register from the hardware stack.

**Processor Status (P):** Restores Negative (N), Overflow (V), Decimal mode (D), Interrupt disable (I), Zero (Z), and Carry (C) from the pulled byte. Pulled bits 5 and 4 do not become persistent flags.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `28` | Implied | `PLP` | 1 | 4 | None |

### ROL — Rotate left

Rotates the accumulator or memory left through carry. Old bit 7 enters carry and the old carry enters bit 0.

**Processor Status (P):** Checks the old Carry (C), then changes Carry (C), Zero (Z), and Negative (N). The old Carry enters result bit 0; the old operand bit 7 becomes the new Carry; Zero is set when the result is zero; Negative copies result bit 7.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `2A` | Accumulator | `ROL A` | 1 | 2 | None |
| `26` | Zero page | `ROL $nn` | 2 | 5 | None |
| `36` | Zero page,X | `ROL $nn,X` | 2 | 6 | None |
| `2E` | Absolute | `ROL $nnnn` | 3 | 6 | None |
| `3E` | Absolute,X | `ROL $nnnn,X` | 3 | 7 | None |

### ROR — Rotate right

Rotates the accumulator or memory right through carry. Old bit 0 enters carry and the old carry enters bit 7.

**Processor Status (P):** Checks the old Carry (C), then changes Carry (C), Zero (Z), and Negative (N). The old Carry enters result bit 7; the old operand bit 0 becomes the new Carry; Zero is set when the result is zero; Negative copies result bit 7.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `6A` | Accumulator | `ROR A` | 1 | 2 | None |
| `66` | Zero page | `ROR $nn` | 2 | 5 | None |
| `76` | Zero page,X | `ROR $nn,X` | 2 | 6 | None |
| `6E` | Absolute | `ROR $nnnn` | 3 | 6 | None |
| `7E` | Absolute,X | `ROR $nnnn,X` | 3 | 7 | None |

### RTI — Return from interrupt

Pulls status and the program counter from the stack, resuming the interrupted program.

**Processor Status (P):** Restores Negative (N), Overflow (V), Decimal mode (D), Interrupt disable (I), Zero (Z), and Carry (C) from the pulled status byte. Pulled bits 5 and 4 do not become persistent flags.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `40` | Implied | `RTI` | 1 | 6 | None |

### RTS — Return from subroutine

Pulls the saved address from the stack, adds one, and resumes after the corresponding JSR.

**Processor Status (P):** No flags checked or changed.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `60` | Implied | `RTS` | 1 | 6 | None |

### SBC — Subtract with carry

Subtracts the operand and inverse carry from the accumulator. It updates carry, zero, negative, and overflow; decimal mode uses BCD arithmetic.

**Processor Status (P):** Checks Carry (C): a clear Carry subtracts one extra, while a set Carry does not. It checks Decimal mode (D) to select binary or BCD arithmetic. Changes Carry (C), Zero (Z), Overflow (V), and Negative (N). In binary mode, Carry is set when no borrow is required; Zero is set when the 8-bit result is zero; Negative copies result bit 7; Overflow is set when A and the operand have different signs and the result's sign differs from A. In decimal mode, see the decimal-flag warning above.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `E9` | Immediate | `SBC #$nn` | 2 | 2 | None |
| `E5` | Zero page | `SBC $nn` | 2 | 3 | None |
| `F5` | Zero page,X | `SBC $nn,X` | 2 | 4 | None |
| `ED` | Absolute | `SBC $nnnn` | 3 | 4 | None |
| `FD` | Absolute,X | `SBC $nnnn,X` | 3 | 4 | Page crossed: +1 |
| `F9` | Absolute,Y | `SBC $nnnn,Y` | 3 | 4 | Page crossed: +1 |
| `E1` | Indexed indirect | `SBC ($nn,X)` | 2 | 6 | None |
| `F1` | Indirect indexed | `SBC ($nn),Y` | 2 | 5 | Page crossed: +1 |

### SEC — Set carry

Sets the carry flag.

**Processor Status (P):** Sets Carry (C). No other flags change.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `38` | Implied | `SEC` | 1 | 2 | None |

### SED — Set decimal mode

Sets the decimal flag, selecting BCD arithmetic for ADC and SBC.

**Processor Status (P):** Sets Decimal mode (D). No other flags change.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `F8` | Implied | `SED` | 1 | 2 | None |

### SEI — Set interrupt disable

Sets the interrupt-disable flag so an asserted IRQ line is not serviced. NMI is unaffected.

**Processor Status (P):** Sets Interrupt disable (I). No other flags change.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `78` | Implied | `SEI` | 1 | 2 | None |

### STA — Store accumulator

Stores A in memory. Store instructions have fixed timing even when indexed addressing crosses a page.

**Processor Status (P):** No flags checked or changed.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `85` | Zero page | `STA $nn` | 2 | 3 | None |
| `95` | Zero page,X | `STA $nn,X` | 2 | 4 | None |
| `8D` | Absolute | `STA $nnnn` | 3 | 4 | None |
| `9D` | Absolute,X | `STA $nnnn,X` | 3 | 5 | None |
| `99` | Absolute,Y | `STA $nnnn,Y` | 3 | 5 | None |
| `81` | Indexed indirect | `STA ($nn,X)` | 2 | 6 | None |
| `91` | Indirect indexed | `STA ($nn),Y` | 2 | 6 | None |

### STX — Store X

Stores X in memory.

**Processor Status (P):** No flags checked or changed.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `86` | Zero page | `STX $nn` | 2 | 3 | None |
| `96` | Zero page,Y | `STX $nn,Y` | 2 | 4 | None |
| `8E` | Absolute | `STX $nnnn` | 3 | 4 | None |

### STY — Store Y

Stores Y in memory.

**Processor Status (P):** No flags checked or changed.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `84` | Zero page | `STY $nn` | 2 | 3 | None |
| `94` | Zero page,X | `STY $nn,X` | 2 | 4 | None |
| `8C` | Absolute | `STY $nnnn` | 3 | 4 | None |

### TAX — Transfer accumulator to X

Copies A into X and updates zero and negative.

**Processor Status (P):** Changes Zero (Z) and Negative (N) from the new X value. Zero is set when X is zero; Negative copies X bit 7.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `AA` | Implied | `TAX` | 1 | 2 | None |

### TAY — Transfer accumulator to Y

Copies A into Y and updates zero and negative.

**Processor Status (P):** Changes Zero (Z) and Negative (N) from the new Y value. Zero is set when Y is zero; Negative copies Y bit 7.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `A8` | Implied | `TAY` | 1 | 2 | None |

### TSX — Transfer stack pointer to X

Copies the stack pointer into X and updates zero and negative.

**Processor Status (P):** Changes Zero (Z) and Negative (N) from the new X value. Zero is set when X is zero; Negative copies X bit 7.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `BA` | Implied | `TSX` | 1 | 2 | None |

### TXA — Transfer X to accumulator

Copies X into A and updates zero and negative.

**Processor Status (P):** Changes Zero (Z) and Negative (N) from the new A value. Zero is set when A is zero; Negative copies A bit 7.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `8A` | Implied | `TXA` | 1 | 2 | None |

### TXS — Transfer X to stack pointer

Copies X into the stack pointer. It does not update flags.

**Processor Status (P):** No flags checked or changed.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `9A` | Implied | `TXS` | 1 | 2 | None |

### TYA — Transfer Y to accumulator

Copies Y into A and updates zero and negative.

**Processor Status (P):** Changes Zero (Z) and Negative (N) from the new A value. Zero is set when A is zero; Negative copies A bit 7.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `98` | Implied | `TYA` | 1 | 2 | None |

## Undocumented Instructions

These opcodes were not specified by MOS Technology. Names are community
conventions, some behaviours depend on the chip revision and electrical
conditions, and they are not guaranteed to behave identically on every 6502.

### ALR — AND then logical shift right (Undocumented)

ANDs an immediate value with the accumulator, shifts the result right, and stores it in the accumulator. Bit 0 moves into carry.

**Processor Status (P):** Changes Carry (C) and Zero (Z), and clears Negative (N). Carry receives bit 0 of the intermediate `A AND operand` value; Zero is set when final A is zero. Negative is always clear after the right shift.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `4B` | Immediate | `ALR #$nn` | 2 | 2 | None |

### ANC — AND and copy negative to carry (Undocumented)

ANDs an immediate value with the accumulator, then copies result bit 7 into both the negative and carry flags.

**Processor Status (P):** Changes Carry (C), Zero (Z), and Negative (N). Carry and Negative both copy final A bit 7; Zero is set when final A is zero.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `0B`, `2B` | Immediate | `ANC #$nn` | 2 | 2 | None |

### ARR — AND then rotate right (Undocumented)

ANDs an immediate value with the accumulator, then rotates right through carry. In binary mode, carry and overflow are derived from the rotated value as described below. Decimal-mode behaviour is unusual and should not be treated as portable.

**Processor Status (P):** Checks the old Carry (C) and Decimal mode (D), then changes Carry (C), Zero (Z), Overflow (V), and Negative (N). In binary mode, the old Carry enters result bit 7; the new Carry copies result bit 6; Overflow is the XOR of result bits 5 and 6; Zero and Negative reflect final A. Decimal-mode behaviour can be unpredictable.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `6B` | Immediate | `ARR #$nn` | 2 | 2 | None |

### DCP — Decrement then compare (Undocumented)

Decrements memory, then compares the new value with the accumulator as CMP would.

**Processor Status (P):** Changes Carry (C), Zero (Z), and Negative (N) from `A - decremented operand`, using the same rules as `CMP`. Any flag results from the intermediate decrement are discarded.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `C7` | Zero page | `DCP $nn` | 2 | 5 | None |
| `D7` | Zero page,X | `DCP $nn,X` | 2 | 6 | None |
| `CF` | Absolute | `DCP $nnnn` | 3 | 6 | None |
| `DF` | Absolute,X | `DCP $nnnn,X` | 3 | 7 | None |
| `DB` | Absolute,Y | `DCP $nnnn,Y` | 3 | 7 | None |
| `C3` | Indexed indirect | `DCP ($nn,X)` | 2 | 8 | None |
| `D3` | Indirect indexed | `DCP ($nn),Y` | 2 | 8 | None |

### DOP — Two-byte no operation (Undocumented)

Consumes an immediate operand without otherwise changing CPU state.

**Processor Status (P):** No flags checked or changed.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `80`, `82`, `C2`, `E2` | Immediate | `DOP #$nn` | 2 | 2 | None |

### ISC — Increment then subtract with carry (Undocumented)

Increments memory, then subtracts the new value from the accumulator as `SBC` would.

**Processor Status (P):** Checks Carry (C) and Decimal mode (D), then changes Carry (C), Zero (Z), Overflow (V), and Negative (N) as `SBC` does. A clear Carry subtracts one extra; Carry is set when no borrow is required. In decimal mode, the caveat for `SBC` also applies.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `E7` | Zero page | `ISC $nn` | 2 | 5 | None |
| `F7` | Zero page,X | `ISC $nn,X` | 2 | 6 | None |
| `EF` | Absolute | `ISC $nnnn` | 3 | 6 | None |
| `FF` | Absolute,X | `ISC $nnnn,X` | 3 | 7 | None |
| `FB` | Absolute,Y | `ISC $nnnn,Y` | 3 | 7 | None |
| `E3` | Indexed indirect | `ISC ($nn,X)` | 2 | 8 | None |
| `F3` | Indirect indexed | `ISC ($nn),Y` | 2 | 8 | None |

### LAX — Load accumulator and X (Undocumented)

Loads the same operand into both A and X, then updates zero and negative.

**Processor Status (P):** Changes Zero (Z) and Negative (N) from the value loaded into A and X. Zero is set when the value is zero; Negative copies its bit 7.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `A7` | Zero page | `LAX $nn` | 2 | 3 | None |
| `B7` | Zero page,Y | `LAX $nn,Y` | 2 | 4 | None |
| `AF` | Absolute | `LAX $nnnn` | 3 | 4 | None |
| `BF` | Absolute,Y | `LAX $nnnn,Y` | 3 | 4 | Page crossed: +1 |
| `A3` | Indexed indirect | `LAX ($nn,X)` | 2 | 6 | None |
| `B3` | Indirect indexed | `LAX ($nn),Y` | 2 | 5 | Page crossed: +1 |

### NOP* — Single-byte no operation opcodes (Undocumented)

These opcodes behave like the documented `NOP` on the 6502.

**Processor Status (P):** No flags checked or changed.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `1A`, `3A`, `5A`, `7A`, `DA`, `FA` | Implied | `NOP*` | 1 | 2 | None |

### RLA — Rotate left then AND (Undocumented)

Rotates memory left through carry, then ANDs the new memory value into the accumulator.

**Processor Status (P):** Checks the old Carry (C), then changes Carry (C), Zero (Z), and Negative (N). The old Carry enters memory bit 0; the old memory bit 7 becomes the new Carry; Zero and Negative reflect final A.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `27` | Zero page | `RLA $nn` | 2 | 5 | None |
| `37` | Zero page,X | `RLA $nn,X` | 2 | 6 | None |
| `2F` | Absolute | `RLA $nnnn` | 3 | 6 | None |
| `3F` | Absolute,X | `RLA $nnnn,X` | 3 | 7 | None |
| `3B` | Absolute,Y | `RLA $nnnn,Y` | 3 | 7 | None |
| `23` | Indexed indirect | `RLA ($nn,X)` | 2 | 8 | None |
| `33` | Indirect indexed | `RLA ($nn),Y` | 2 | 8 | None |

### RRA — Rotate right then add with carry (Undocumented)

Rotates memory right through carry, then adds the new memory value to the accumulator.

**Processor Status (P):** Checks the old Carry (C) and Decimal mode (D), then changes Carry (C), Zero (Z), Overflow (V), and Negative (N). The old Carry enters memory bit 7, and old memory bit 0 becomes the subsequent addition's carry-in. The final flags follow `ADC`; the decimal-mode caveat therefore applies when D is set.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `67` | Zero page | `RRA $nn` | 2 | 5 | None |
| `77` | Zero page,X | `RRA $nn,X` | 2 | 6 | None |
| `6F` | Absolute | `RRA $nnnn` | 3 | 6 | None |
| `7F` | Absolute,X | `RRA $nnnn,X` | 3 | 7 | None |
| `7B` | Absolute,Y | `RRA $nnnn,Y` | 3 | 7 | None |
| `63` | Indexed indirect | `RRA ($nn,X)` | 2 | 8 | None |
| `73` | Indirect indexed | `RRA ($nn),Y` | 2 | 8 | None |

### SAX — Store A AND X (Undocumented)

Stores `A AND X` in memory without changing either register.

**Processor Status (P):** No flags checked or changed.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `87` | Zero page | `SAX $nn` | 2 | 3 | None |
| `97` | Zero page,Y | `SAX $nn,Y` | 2 | 4 | None |
| `8F` | Absolute | `SAX $nnnn` | 3 | 4 | None |
| `83` | Indexed indirect | `SAX ($nn,X)` | 2 | 6 | None |

### SBC* — Immediate SBC alias (Undocumented)

Performs the same immediate subtract-with-carry operation as opcode `E9`.

**Processor Status (P):** Checks Carry (C) and Decimal mode (D), then changes Carry (C), Zero (Z), Overflow (V), and Negative (N), exactly as documented for `SBC`.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `EB` | Immediate | `SBC* #$nn` | 2 | 2 | None |

### SKB — Skip byte through zero page (Undocumented)

Reads a zero-page operand and otherwise behaves as a no operation.

**Processor Status (P):** No flags checked or changed.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `04`, `44`, `64` | Zero page | `SKB $nn` | 2 | 3 | None |

### SKW — Skip byte through zero page,X (Undocumented)

Reads a zero-page,X operand and otherwise behaves as a no operation.

**Processor Status (P):** No flags checked or changed.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `14`, `34`, `54`, `74`, `D4`, `F4` | Zero page,X | `SKW $nn,X` | 2 | 4 | None |

### SLO — Shift left then OR (Undocumented)

Shifts memory left, then ORs the new memory value into the accumulator.

**Processor Status (P):** Changes Carry (C), Zero (Z), and Negative (N). Carry receives old memory bit 7; Zero is set when final A is zero; Negative copies final A bit 7.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `07` | Zero page | `SLO $nn` | 2 | 5 | None |
| `17` | Zero page,X | `SLO $nn,X` | 2 | 6 | None |
| `0F` | Absolute | `SLO $nnnn` | 3 | 6 | None |
| `1F` | Absolute,X | `SLO $nnnn,X` | 3 | 7 | None |
| `1B` | Absolute,Y | `SLO $nnnn,Y` | 3 | 7 | None |
| `03` | Indexed indirect | `SLO ($nn,X)` | 2 | 8 | None |
| `13` | Indirect indexed | `SLO ($nn),Y` | 2 | 8 | None |

### SRE — Shift right then exclusive OR (Undocumented)

Shifts memory right, then exclusive-ORs the new memory value into the accumulator.

**Processor Status (P):** Changes Carry (C), Zero (Z), and Negative (N). Carry receives old memory bit 0; Zero is set when final A is zero; Negative copies final A bit 7.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `47` | Zero page | `SRE $nn` | 2 | 5 | None |
| `57` | Zero page,X | `SRE $nn,X` | 2 | 6 | None |
| `4F` | Absolute | `SRE $nnnn` | 3 | 6 | None |
| `5F` | Absolute,X | `SRE $nnnn,X` | 3 | 7 | None |
| `5B` | Absolute,Y | `SRE $nnnn,Y` | 3 | 7 | None |
| `43` | Indexed indirect | `SRE ($nn,X)` | 2 | 8 | None |
| `53` | Indirect indexed | `SRE ($nn),Y` | 2 | 8 | None |

### TOP — Three-byte no operation (Undocumented)

Reads an absolute operand and otherwise behaves as a no operation. The absolute,X forms take one extra cycle when indexing crosses a page.

**Processor Status (P):** No flags checked or changed.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `0C` | Absolute | `TOP $nnnn` | 3 | 4 | None |
| `1C`, `3C`, `5C`, `7C`, `DC`, `FC` | Absolute,X | `TOP $nnnn,X` | 3 | 4 | Page crossed: +1 |

### XAA — Transfer X AND immediate to accumulator (Undocumented)

On the 6502 this instruction combines X, an immediate operand, and an internal bus value, then stores the result in A. Its result is electrically unstable and can vary with chip revision, temperature, and supply voltage; no single deterministic formula is reliable.

**Processor Status (P):** Changes Zero (Z) and Negative (N) from the value placed in A. Zero is set when A is zero; Negative copies A bit 7.

| Opcode | Addressing mode | Syntax | Bytes | Cycles | Extra cycles |
|---:|---|---|---:|---:|---|
| `8B` | Immediate | `XAA #$nn` | 2 | 2 | None |
