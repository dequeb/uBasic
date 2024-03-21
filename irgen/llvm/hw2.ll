@s = global [7 x i8] c"Hello \00"
@_1 = global [6 x i8] c"world\00"

declare i32 @printf(i8* %format, ...)

define i32 @main() {
0:
	%1 = call i32 (i8*, ...) @printf(i8* getelementptr ([7 x i8], [7 x i8]* @s, i64 0, i64 0))
	%2 = call i32 (i8*, ...) @printf(i8* getelementptr ([6 x i8], [6 x i8]* @_1, i64 0, i64 0))
	%3 = call i32 (i8*, ...) @printf([2 x i8] c"\0A\00")
	ret i32 0
}
