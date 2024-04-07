@str = global [9 x i8] c"==> %d \0A\00"

declare i32 @printf(i8* %0, ...)

define i32 @varia(i32 %".values_0)", i32* %values) {
entry:
	%total = alloca i32
	%count = alloca i32
	%i = alloca i32
	store i32 %".values_0)", i32* %count
	store i32 0, i32* %i
	store i32 0, i32* %total
	br label %loop.cond

loop.cond:
	%0 = load i32, i32* %i
	%1 = load i32, i32* %count
	%2 = icmp ult i32 %0, %1
	br i1 %2, label %loop.body, label %loop.end

loop.body:
	%3 = getelementptr i32, i32* %values, i32 %0
	%4 = load i32, i32* %3
	%5 = load i32, i32* %total
	%6 = add i32 %4, %5
	store i32 %6, i32* %total
	%7 = load i32, i32* %i
	%8 = add i32 %7, 1
	store i32 %8, i32* %i
	br label %loop.cond

loop.end:
	%9 = load i32, i32* %total
	ret i32 %9
}

define i32 @main() {
0:
	%1 = alloca [2 x i32]
	%2 = getelementptr i32, [2 x i32]* %1, i64 0
	store i32 10, i32* %2
	%3 = getelementptr i32, [2 x i32]* %1, i64 1
	store i32 20, i32* %3
	%4 = load i32, i32* %2
	%5 = call i32 (i8*, ...) @printf([9 x i8]* @str, i32 %4)
	%6 = load i32, i32* %3
	%7 = call i32 (i8*, ...) @printf([9 x i8]* @str, i32 %6)
	%8 = call i32 @varia(i32 2, [2 x i32]* %1)
	%9 = call i32 (i8*, ...) @printf([9 x i8]* @str, i32 %8)
	ret i32 0
}