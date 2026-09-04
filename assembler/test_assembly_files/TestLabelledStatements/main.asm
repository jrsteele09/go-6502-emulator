org $1400
start lda #$01
next: nop
value db $02,$03
ptr dw start
load lda value
space ds 2
colon_label: lda #$ff
end jsr start
end start
