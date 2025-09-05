        *=$c000
loop    lda data,X
        bne loop
data    .byte 1,2