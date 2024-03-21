source_filename = "irgen/llvm/compile.ll"

@vbEmpty = constant [1 x i8] c"\00"
@vbCR = constant [2 x i8] c"\0D\00"
@vbLF = constant [2 x i8] c"\0A\00"
@vbCrLf = constant [3 x i8] c"\0D\0A\00"
@vbTab = constant [2 x i8] c"\09\00"
@true = constant [5 x i8] c"True\00"
@false = constant [6 x i8] c"False\00"
@.JumpBuffer = global [48 x i32] zeroinitializer
@.ErrorNumber = global i32 0
@.ErrorMessage = global [256 x i8] zeroinitializer
@.divisionByZero = global [18 x i8] c"Division by zero\0A\00"
@.arrayIndexOutOfBounds = global [27 x i8] c"Array index out of bounds\0A\00"

declare i8* @malloc(i64 %size)

declare i8* @calloc(i64 %n, i64 %size)

declare void @free(i8* %ptr)

declare i32 @printf(i8* %format, ...)

declare i32 @puts(i8* %s)

declare i32 @scanf(i8* %format, ...)

declare i8* @strcpy(i8* %dst, i8* %src)

declare i8* @strcat(i8* %dst, i8* %src)

declare i32 @sscanf(i8* %str, i8* %format, i8* %dst)

declare i32 @strlen(i8* %str)

declare i32 @sprintf(i8* %str, i8* %format, ...)

declare void @exit(i32 %status)

declare i32 @setjmp(i32* %0)

declare void @longjmp(i32* %0, i32 %1)

define void @.throwException() {
0:
	call void @longjmp([48 x i32]* @.JumpBuffer, i32 1)
	unreachable
}

define i32 @main() {
0:
	%1 = call i32 @setjmp([48 x i32]* @.JumpBuffer)
	%2 = icmp eq i32 %1, 0
	br i1 %2, label %normalCall, label %exception

normalCall:
	call void @.main()
	br label %end

exception:
	%3 = load i32, i32* @.ErrorNumber
	%4 = call i32 (i8*, ...) @printf([256 x i8]* @.ErrorMessage)
	ret i32 %3

end:
	ret i32 0
}

define void @.main() {
0:
	%1 = sitofp i64 3 to double
	%2 = fcmp oeq double %1, 0.0
	br i1 %2, label %3, label %5

3:
	store i32 1, i32* @.ErrorNumber
	%4 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [18 x i8]* @.divisionByZero)
	call void @.throwException()
	unreachable

5:
	%6 = fdiv double 1.0, %1
	%7 = alloca [4 x i8]
	store [4 x i8] c"%lf\00", [4 x i8]* %7
	%8 = call i32 (i8*, ...) @printf([4 x i8]* %7, double %6)
	%9 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%10 = sitofp i64 3 to double
	%11 = fcmp oeq double 2.0, 0.0
	br i1 %11, label %12, label %14

12:
	store i32 1, i32* @.ErrorNumber
	%13 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [18 x i8]* @.divisionByZero)
	call void @.throwException()
	unreachable

14:
	%15 = fdiv double %10, 2.0
	%16 = alloca [4 x i8]
	store [4 x i8] c"%lf\00", [4 x i8]* %16
	%17 = call i32 (i8*, ...) @printf([4 x i8]* %16, double %15)
	%18 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%19 = icmp eq i64 2, 0
	br i1 %19, label %20, label %22

20:
	store i32 1, i32* @.ErrorNumber
	%21 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [18 x i8]* @.divisionByZero)
	call void @.throwException()
	unreachable

22:
	%23 = sdiv i64 3, 2
	%24 = alloca [4 x i8]
	store [4 x i8] c"%ld\00", [4 x i8]* %24
	%25 = call i32 (i8*, ...) @printf([4 x i8]* %24, i64 %23)
	%26 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%27 = icmp eq i64 0, 0
	br i1 %27, label %28, label %30

28:
	store i32 1, i32* @.ErrorNumber
	%29 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [18 x i8]* @.divisionByZero)
	call void @.throwException()
	unreachable

30:
	%31 = sdiv i64 1, 0
	%32 = alloca [4 x i8]
	store [4 x i8] c"%ld\00", [4 x i8]* %32
	%33 = call i32 (i8*, ...) @printf([4 x i8]* %32, i64 %31)
	%34 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	ret void
}
