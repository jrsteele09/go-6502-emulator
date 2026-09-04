FEATURE = 1

if FEATURE = 1
if 1
VALUE = $42
else
VALUE = $98
endif
else
VALUE = $99
endif

LOAD_STORE macro VALUE, TARGET
    lda #VALUE
    sta TARGET
endm

MARK macro
    nop
endm

org $1000
LOAD_STORE $42, $C000
MARK

if 0
    db $FF
else
    db $7E
endif
ds 1
db $01
