VALUE = 1 // inline slash comment after a constant

;COMMENTED macro
;    lda #$ff
;    endm

/*
BLOCKED macro
    lda #$ee
    endm

if 1
    db $ff
endif
*/

if VALUE = 1 ; inline semicolon comment after a conditional
REAL macro ARG /* block comment after macro declaration */
    lda #ARG ; preserve ordinary source comments
endm // slash comment after macro terminator
else
UNUSED macro
    lda #$00
endm
endif

org $1300
REAL $7f
db $2f // inline slash comment after ordinary source
db $3b ; inline semicolon comment after ordinary source
.text ";///*"
