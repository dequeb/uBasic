@array = global [2 x i32] zeroinitializer
@str = global [8 x i8] c"==> %d\0A\00"

declare i32 @printf(i8* %0, ...)

define i32 @main() {
0:
	store i32 10, i32* getelementptr ([2 x i32], [2 x i32]* @array, i64 0, i64 0)
	store i32 20, i32* getelementptr ([2 x i32], [2 x i32]* @array, i64 0, i64 1)
	%1 = load i32, i32* getelementptr ([2 x i32], [2 x i32]* @array, i64 0, i64 0)
	%2 = call i32 (i8*, ...) @printf(i8* getelementptr ([8 x i8], [8 x i8]* @str, i64 0, i64 0), i32 %1)
	%3 = load i32, i32* getelementptr ([2 x i32], [2 x i32]* @array, i64 0, i64 1)
	%4 = call i32 (i8*, ...) @printf(i8* getelementptr ([8 x i8], [8 x i8]* @str, i64 0, i64 0), i32 %3)
	%5 = load i32, i32* getelementptr ([2 x i32], [2 x i32]* @array, i64 0, i64 0)
	%6 = load i32, i32* getelementptr ([2 x i32], [2 x i32]* @array, i64 0, i64 1)
	%7 = add i32 %5, %6
	%8 = call i32 (i8*, ...) @printf(i8* getelementptr ([8 x i8], [8 x i8]* @str, i64 0, i64 0), i32 %7)
	ret i32 0
}
