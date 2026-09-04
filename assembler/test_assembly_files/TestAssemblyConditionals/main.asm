org $1600

data_start
        db $01, $02, $03
data_end

bss_start
        ds 3
bss_end

        if (data_end - data_start) != (bss_end - bss_start)
            ERROR ERROR ERROR
        else
            lda #lo(data_end - data_start)
        endif

        if (bss_end - bss_start) = 3
            ldx #hi(data_start)
        endif

        end data_start
