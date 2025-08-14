.ORG $1000

LDA #$42        ; Load accumulator with hex 42
STA $D020       ; Store in border color register (C64)

LDX #$10        ; Load X register with 16
LOOP:
    DEX         ; Decrement X register
    BNE LOOP    ; Branch if not equal to zero
    BEQ END
RTS             ; Return from subroutine
END:
    RTS
