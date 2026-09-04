ROM_vectors = 1
load_data_direct = 1
data_segment = $200
IRQ_bit = 3
PASS_VALUE = $42
FAIL_VALUE = $99
intdis equ %00000100
MASK equ $ff&~intdis

EMIT_STATUS macro VALUE, TARGET
    lda #VALUE
    sta TARGET
endm

LOAD_IMMEDIATE macro
    lda #\1
endm

org $1200
if (data_segment & $ff) = 0
    EMIT_STATUS PASS_VALUE, $C010
else
    EMIT_STATUS FAIL_VALUE, $C010
endif

if (load_data_direct = 1) & (ROM_vectors = 1)
    db $A5
else
    db $00
endif

if (1 << IRQ_bit) = 8
    db $08
endif

if 'R'-3 = $4f
    db $4f
endif

if MASK = $fb
    LOAD_IMMEDIATE MASK
endif
