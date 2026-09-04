org $1500
MASK = 251
target lda #($ff-'B'+1)
       and #MASK&$0f
       ora #(1<<3|1)
       cmp #'K'
       cmp #('F'^$aa)
       lda #hi(target)
       ldx #lo(target)
bytes  db MASK&$ff, 'A'+1, ~$fc&$ff
words  dw target, target+1
       ds 1<<0
       db $ee
endpad db $ff
       end target
