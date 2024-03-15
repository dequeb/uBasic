	.section	__TEXT,__text,regular,pure_instructions
	.build_version macos, 14, 0
	.globl	_main                           ; -- Begin function main
	.p2align	2
_main:                                  ; @main
	.cfi_startproc
; %bb.0:
	sub	sp, sp, #80
	stp	x22, x21, [sp, #32]             ; 16-byte Folded Spill
	stp	x20, x19, [sp, #48]             ; 16-byte Folded Spill
	stp	x29, x30, [sp, #64]             ; 16-byte Folded Spill
	.cfi_def_cfa_offset 80
	.cfi_offset w30, -8
	.cfi_offset w29, -16
	.cfi_offset w19, -24
	.cfi_offset w20, -32
	.cfi_offset w21, -40
	.cfi_offset w22, -48
	adrp	x8, _i@PAGE
	mov	w9, #10                         ; =0xa
	mov	x11, #3758096384                ; =0xe0000000
	mov	w10, #1311                      ; =0x51f
	movk	x11, #163, lsl #32
	adrp	x20, _s@PAGE
	movk	w10, #17096, lsl #16
	adrp	x21, _l@PAGE
	adrp	x22, _d@PAGE
	movk	x11, #16473, lsl #48
	str	w9, [x8, _i@PAGEOFF]
	mov	w8, #25637                      ; =0x6425
	add	x0, sp, #29
	str	w10, [x20, _s@PAGEOFF]
	str	x9, [x21, _l@PAGEOFF]
	str	x11, [x22, _d@PAGEOFF]
	sturh	w8, [sp, #29]
	strb	wzr, [sp, #31]
	str	x9, [sp]
	bl	_printf
Lloh0:
	adrp	x19, _vbEmpty@PAGE
Lloh1:
	add	x19, x19, _vbEmpty@PAGEOFF
	mov	x0, x19
	bl	_puts
	ldr	s0, [x20, _s@PAGEOFF]
	mov	w8, #26149                      ; =0x6625
	add	x0, sp, #26
	strb	wzr, [sp, #28]
	fcvt	d0, s0
	strh	w8, [sp, #26]
	str	d0, [sp]
	bl	_printf
	mov	x0, x19
	bl	_puts
	mov	w8, #27685                      ; =0x6c25
	ldr	x9, [x21, _l@PAGEOFF]
	movk	w8, #100, lsl #16
	add	x0, sp, #22
	str	x9, [sp]
	stur	w8, [sp, #22]
	bl	_printf
	mov	x0, x19
	bl	_puts
	mov	w8, #27685                      ; =0x6c25
	ldr	d0, [x22, _d@PAGEOFF]
	movk	w8, #102, lsl #16
	add	x0, sp, #18
	str	d0, [sp]
	stur	w8, [sp, #18]
	bl	_printf
	mov	x0, x19
	bl	_puts
	ldp	x29, x30, [sp, #64]             ; 16-byte Folded Reload
	mov	w0, wzr
	ldp	x20, x19, [sp, #48]             ; 16-byte Folded Reload
	ldp	x22, x21, [sp, #32]             ; 16-byte Folded Reload
	add	sp, sp, #80
	ret
	.loh AdrpAdd	Lloh0, Lloh1
	.cfi_endproc
                                        ; -- End function
	.section	__TEXT,__const
	.globl	_vbEmpty                        ; @vbEmpty
_vbEmpty:
	.space	1

	.globl	_vbCR                           ; @vbCR
_vbCR:
	.asciz	"\r"

	.globl	_vbLF                           ; @vbLF
_vbLF:
	.asciz	"\n"

	.globl	_vbCrLf                         ; @vbCrLf
_vbCrLf:
	.asciz	"\r\n"

	.globl	_vbTab                          ; @vbTab
_vbTab:
	.asciz	"\t"

	.globl	_true                           ; @true
_true:
	.asciz	"True"

	.globl	_false                          ; @false
_false:
	.asciz	"False"

	.globl	_i                              ; @i
.zerofill __DATA,__common,_i,4,2
	.globl	_s                              ; @s
.zerofill __DATA,__common,_s,4,2
	.globl	_l                              ; @l
.zerofill __DATA,__common,_l,8,3
	.globl	_d                              ; @d
.zerofill __DATA,__common,_d,8,3
.subsections_via_symbols
