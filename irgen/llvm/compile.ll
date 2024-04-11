source_filename = "irgen/llvm/compile"

%struct.GarbageCollector = type { i8*, i8, i8*, i64 }
%LongArrayType = type { i32, i64* }
%IntegerArrayType = type { i32, i32* }
%SingleArrayType = type { i32, float* }
%DoubleArrayType = type { i32, double* }
%StringArrayType = type { i32, i8* }
%BooleanArrayType = type { i32, i1* }

@gc = external global %struct.GarbageCollector, align 8
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
@.divisionByZero = private global [18 x i8] c"Division by zero\0A\00"
@.arrayIndexOutOfBounds = private global [27 x i8] c"Array index out of bounds\0A\00"

declare void @gc_start(%struct.GarbageCollector* noundef %gc, i8* noundef %base_stack)

declare i64 @gc_stop(%struct.GarbageCollector* noundef %gc)

declare i8* @gc_malloc(%struct.GarbageCollector* noundef %gc, i8* noundef %size)

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

define void @AddArrayNumbers() {
0:
	%.a_1 = alloca i32
	store i32 2, i32* %.a_1
	%.a_0 = alloca i32
	store i32 3, i32* %.a_0
	%1 = alloca [6 x i32]
	%2 = alloca i32
	store i32 0, i32* %2
	br label %3

3:
	%4 = load i32, i32* %2
	%5 = icmp ult i32 %4, 6
	br i1 %5, label %6, label %10

6:
	%7 = getelementptr i32, [6 x i32]* %1, i32 %4
	store i32 0, i32* %7
	%8 = load i32, i32* %2
	%9 = add i32 %8, 1
	store i32 %9, i32* %2
	br label %3

10:
	ret void
}

define i32 @main(i32 %argc, i8** %argv) {
0:
	%1 = alloca i32
	store i32 %argc, i32* %1
	call void @gc_start(%struct.GarbageCollector* @gc, i32* %1)
	%2 = call i32 @setjmp([48 x i32]* @.JumpBuffer)
	%3 = icmp eq i32 %2, 0
	br i1 %3, label %normalCall, label %exception

normalCall:
	call void @.main()
	br label %end

exception:
	%4 = load i32, i32* @.ErrorNumber
	%5 = call i32 (i8*, ...) @printf([256 x i8]* @.ErrorMessage)
	ret i32 %4

end:
	%6 = call i64 @gc_stop(%struct.GarbageCollector* @gc)
	ret i32 0
}

define void @.main() {
0:
	ret void
}
