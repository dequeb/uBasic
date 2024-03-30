; @str = global [13 x i8] c"Hello, %lf!\0A\00"

; declare i32 @printf(i8* %0, ...)

; define i32 @main() {
; 0:
; 	%1 = fadd double 0x3FF199999999999A, 2.0
; 	%2 = load i8*, [13 x i8]* @str
; 	%3 = call i32 (i8*, ...) @printf(i8* %2, double %1)
; 	ret i32 0
; }

; @str = global [13 x i8] c"Hello, %lf!\0A\00"

; declare i32 @printf(i8* %0, ...)

; define i32 @main() {
; 0:
; 	%1 = fadd double 0x3FEFAE147AE147AE, 2.0
; 	%2 = call i32 (i8*, ...) @printf([13 x i8]* @str, double %1)
; 	ret i32 0
; }

@str = global [13 x i8] c"Hello, %lf!\0A\00"

declare i32 @printf(i8* %0, ...)

define i32 @main() {
0:
	%1 = fadd float 0x3FEFAE1480000000, 0x40100A3D80000000
	%2 = call i32 (i8*, ...) @printf([13 x i8]* @str, float %1)
	ret i32 0
}