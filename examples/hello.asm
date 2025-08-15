.include "constants.inc"

        .ORG $1000
        LDA #TEST_VALUE     ; Load accumulator with test value
        STA BORDER_COLOR    ; Store in border color register

        LDX #LOOP_COUNT     ; Load X register with loop count
LOOP    DEX         ; Decrement X register
        BNE LOOP    ; Branch if not equal to zero
        BEQ END
        RTS         ; Return from subroutine
END:    RTS
