; little demo to open up the border
; for win2c64 by Aart Bik
; https://www.aartbik.com/

scroly   =   $d011
raster   =   $d012
vicirq   =   $d019
irqmsk   =   $d01a
ciaicr   =   $dc0d
ci2icr   =   $dd0d
garbage  =   $3fff

;
; encode SYS 2064 line
; in BASIC program space
;
        .org  $0801
        .byte     $0c $08 $0a $00 $9e $20 $32
        .byte $30 $36 $34 $00 $00 $00 $00 $00

lab2064 sei            ; disable irq
        ldx #$7f       ;
        stx ciaicr     ; disable timer irq CIA 1
        stx ci2icr     ; disable timer irq CIA 2
        ldx #$01       ;
        stx irqmsk     ; enable raster irq
        ldx #<nearend  ;
        stx $0314      ;
        ldx #>nearend  ;
        stx $0315      ; set handler
        ldx #$1b       ;
        stx scroly     ; 25 rows
        ldx #$f9       ;
        stx raster     ; irq at raster $f9
        ldx #$00       ;
        stx garbage    ; clear "garbage"
        cli            ; enable irq
        rts            ; done
;
; NEAR END HANDLER
;
nearend ldx #$13       ;
        stx scroly     ; 24 rows
delay   inx            ;
        bne delay      ;
        ldx #$1b       ;
        stx scroly     ; 25 rows
        ldx #$01       ;
        stx vicirq     ; ack irq
        jmp $ea31      ; continue