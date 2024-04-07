@ten = global i64 1
@str = global [13 x i8] c"Value: %ld!\0A\00"

declare i32 @printf(i8* %0, ...)

define i32 @main() {
0:
	%1 = load i64, i64* @ten
	%2 = icmp eq i64 %1, 1
	br i1 %2, label %if.true243, label %if.false243

if.true243:
	%3 = sub i64 %1, 1
	br label %if.end243

if.false243:
	%4 = add i64 %1, 1
	br label %if.end243

if.end243:
	%5 = phi i64 [ %3, %if.true243 ], [ %4, %if.false243 ]
	store i64 %5, i64* @ten
	%6 = load i64, i64* getelementptr (i64, i64* @ten, i64 0)
	%7 = call i32 (i8*, ...) @printf(i8* getelementptr ([13 x i8], [13 x i8]* @str, i64 0, i64 0), i64 %6)
	ret i32 0
}



; @true = global [5 x i8] c"True\00"
; @false = global [6 x i8] c"False\00"

; declare i32 @puts(i8* %str)

; define i8* @_fromCharToStringBoolean_(i1 %value) {
; 0:
; 	%1 = icmp eq i1 %value, true
; 	br i1 %1, label %if.true, label %if.false

; if.true:
; 	br label %if.end

; if.false:
; 	br label %if.end

; if.end:
; 	%2 = phi [5 x i8]* [ getelementptr ([5 x i8], [5 x i8]* @true, i64 0), %if.true ], [ getelementptr ([6 x i8], [6 x i8]* @false, i64 0), %if.false ]
; 	ret [5 x i8]* %2
; }

; define i32 @main() {
; 0:
; 	%1 = call i8* @_fromCharToStringBoolean_(i1 true)
; 	%2 = call i32 @puts(i8* %1)
; 	%3 = call i8* @_fromCharToStringBoolean_(i1 false)
; 	%4 = call i32 @puts(i8* %3)
; 	ret i32 0
; }
