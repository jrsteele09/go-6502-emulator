; This is garbage code, but it tests simple "+" "-" labels
            *=$c000
-           jmp +++
            bne -
+++         lda #1
-           ldx #2
            ldx #3
            jmp +
            bne -
+           rts
