@jump_buffer = global [48 x i32] zeroinitializer
@errorNumber = global i32 0
@errorMessage = global [256 x i8] zeroinitializer
@.str0 = global [22 x i8] c"Division by zero: %d\0A\00"

declare i32 @setjmp(i32* %0)

declare void @longjmp(i32* %0, i32 %1)

declare i32 @printf(i8* %0, ...)

declare i8* @strcpy(i8* %dst, i8* %src)

define void @throwException() {
0:
	call void @longjmp([48 x i32]* @jump_buffer, i32 1)
	unreachable
}

define void @function_that_might_throw_exception() {
0:
	%1 = call i8* @strcpy([256 x i8]* @errorMessage, [22 x i8]* @.str0)
	store i32 17, i32* @errorNumber
	call void @throwException()
	ret void
}

define i32 @main() {
0:
	%1 = call i32 @setjmp([48 x i32]* @jump_buffer)
	%2 = icmp eq i32 %1, 0
	br i1 %2, label %normalCall, label %exception

exception:
	%3 = load i32, i32* @errorNumber
	%4 = call i32 (i8*, ...) @printf([256 x i8]* @errorMessage, i32 %3)
	ret i32 1

normalCall:
	call void @function_that_might_throw_exception()
	br label %end

end:
	ret i32 0
}