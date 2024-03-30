source_filename = "irgen/llvm/compile"

%struct.GarbageCollector = type { i8*, i8, i8*, i64 }

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
@l = global i64 0
@d = global double 0.0
@b = global i1 false
@da = global float 0.0
@s = global float 0.0
@s0 = global i8* null
@c = global float 0.0
@i = global i32 0
@c1 = constant i64 10
@c2 = constant double 11.5
@c3 = constant [6 x i8] c"Hello\00"
@c4 = constant i1 true
@c5 = constant float 0x417E285040000000
@c0 = constant i64 10
@s1 = global i8* null
@s2 = global i8* null
@s3 = global i8* null
@k0 = constant i64 10
@k1 = constant i64 20
@i11 = global i32 0
@addBy = global i32 0
@FloatAddBy = global float 0.0
@.global_b8_1 = constant i64 3
@.global_b8_0 = constant i64 3
@b8 = global [9 x i32] zeroinitializer
@.global_a_0 = constant i64 2
@a = global [2 x i64] zeroinitializer
@curr = global float 0.0
@lon = global i64 0
@sng = global float 0.0
@dbl = global double 0.0
@str = global i8* null
@int = global i32 0
@bool = global i1 false
@.global_curra_0 = constant i64 2
@curra = global [2 x float] zeroinitializer
@.global_lona_0 = constant i64 2
@lona = global [2 x i64] zeroinitializer
@.global_snga_0 = constant i64 2
@snga = global [2 x float] zeroinitializer
@.global_dbla_0 = constant i64 2
@dbla = global [2 x double] zeroinitializer
@.global_stra_0 = constant i64 2
@stra = global [2 x i8*] zeroinitializer
@.global_inta_0 = constant i64 2
@inta = global [2 x i32] zeroinitializer
@.global_boola_0 = constant i64 2
@boola = global [2 x i1] zeroinitializer
@currcg = constant float 18.0
@loncg = constant i64 28
@sngcg = constant float 0x40656570A0000000
@dblcg = constant double 69.0
@strcg = constant [12 x i8] c"hello world\00"
@intcg = constant i32 72
@boolcg = constant i1 false
@datecg = constant float 0x42012E0BE0000000
@.global_arrayf_1 = constant i64 2
@.global_arrayf_0 = constant i64 3
@arrayf = global [6 x i32] zeroinitializer
@.global_arrayd_0 = global i64 0
@.declareLocalArrays_arrayf_1 = constant i64 3
@.declareLocalArrays_arrayf_0 = constant i64 2
@.declareLocalArrays_arrayd_0 = global i64 0
@_1 = constant [14 x i8] c"Allo le monde\00"
@_2 = constant [6 x i8] c"hello\00"
@_3 = constant [6 x i8] c"world\00"
@_4 = constant [2 x i8] c" \00"
@_5 = constant [3 x i8] c" !\00"
@_6 = constant [7 x i8] c"Hello \00"
@_7 = constant [8 x i8] c"world !\00"
@_8 = constant [6 x i8] c"hello\00"
@_9 = constant [6 x i8] c"hello\00"

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

define void @ifBranch() {
0:
	%d = alloca i64
	store i64 10, i64* %d
	%1 = load i64, i64* @c0
	%2 = load i64, i64* %d
	%3 = icmp ne i64 %1, %2
	br i1 %3, label %4, label %8

4:
	%5 = alloca [12 x i8]
	store [12 x i8] c"c0 is Not d\00", [12 x i8]* %5
	%6 = call i32 (i8*, ...) @printf([12 x i8]* %5)
	%7 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	br label %12

8:
	%9 = alloca [8 x i8]
	store [8 x i8] c"c0 is d\00", [8 x i8]* %9
	%10 = call i32 (i8*, ...) @printf([8 x i8]* %9)
	%11 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	br label %12

12:
	ret void
}

define void @while1() {
0:
	%i = alloca i32
	%1 = trunc i64 20 to i32
	store i32 %1, i32* %i
	br label %2

2:
	%3 = load i32, i32* %i
	%4 = sext i32 %3 to i64
	%5 = icmp sgt i64 %4, 10
	br i1 %5, label %6, label %15

6:
	%7 = load i32, i32* %i
	%8 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %8
	%9 = call i32 (i8*, ...) @printf([3 x i8]* %8, i32 %7)
	%10 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%11 = load i32, i32* %i
	%12 = sext i32 %11 to i64
	%13 = sub i64 %12, 1
	%14 = trunc i64 %13 to i32
	store i32 %14, i32* %i
	br label %2

15:
	ret void
}

define i64 @times(i32 %m, i64 %n) {
0:
	%1 = alloca i32
	store i32 %m, i32* %1
	%2 = alloca i64
	store i64 %n, i64* %2
	%3 = alloca i64
	store i64 0, i64* %3
	%4 = load i64, i64* %2
	%5 = load i32, i32* %1
	%6 = sext i32 %5 to i64
	%7 = mul i64 %4, %6
	store i64 %7, i64* %3
	%8 = load i64, i64* %3
	ret i64 %8
}

define i64 @times2(i64 %n) {
0:
	%1 = alloca i64
	store i64 %n, i64* %1
	%2 = alloca i64
	store i64 0, i64* %2
	%3 = load i64, i64* %1
	%4 = mul i64 %3, 2
	store i64 %4, i64* %2
	%5 = load i64, i64* %2
	ret i64 %5
}

define void @listParmArray(i32* %a, i32 %a_size) {
0:
	%1 = alloca i32*
	store i32* %a, i32** %1
	%2 = alloca i32
	store i32 %a_size, i32* %2
	%3 = sext i32 %a_size to i64
	%4 = icmp ult i64 0, 0
	br i1 %4, label %5, label %7

5:
	store i32 2, i32* @.ErrorNumber
	%6 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

7:
	%8 = icmp uge i64 0, %3
	br i1 %8, label %5, label %9

9:
	%10 = getelementptr i32, i32* %a, i64 0
	%11 = load i32, i32* %10
	%12 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %12
	%13 = call i32 (i8*, ...) @printf([3 x i8]* %12, i32 %11)
	%14 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%15 = sext i32 %a_size to i64
	%16 = icmp ult i64 1, 0
	br i1 %16, label %17, label %19

17:
	store i32 2, i32* @.ErrorNumber
	%18 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

19:
	%20 = icmp uge i64 1, %15
	br i1 %20, label %17, label %21

21:
	%22 = getelementptr i32, i32* %a, i64 1
	%23 = load i32, i32* %22
	%24 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %24
	%25 = call i32 (i8*, ...) @printf([3 x i8]* %24, i32 %23)
	%26 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%27 = sext i32 %a_size to i64
	%28 = icmp ult i64 2, 0
	br i1 %28, label %29, label %31

29:
	store i32 2, i32* @.ErrorNumber
	%30 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

31:
	%32 = icmp uge i64 2, %27
	br i1 %32, label %29, label %33

33:
	%34 = getelementptr i32, i32* %a, i64 2
	%35 = load i32, i32* %34
	%36 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %36
	%37 = call i32 (i8*, ...) @printf([3 x i8]* %36, i32 %35)
	%38 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	ret void
}

define void @addParmArray(i32* %a, i32 %a_size) {
0:
	%1 = alloca i32*
	store i32* %a, i32** %1
	%2 = alloca i32
	store i32 %a_size, i32* %2
	%3 = sext i32 %a_size to i64
	%4 = icmp ult i64 0, 0
	br i1 %4, label %5, label %7

5:
	store i32 2, i32* @.ErrorNumber
	%6 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

7:
	%8 = icmp uge i64 0, %3
	br i1 %8, label %5, label %9

9:
	%10 = getelementptr i32, i32* %a, i64 0
	%11 = load i32, i32* %10
	%12 = sext i32 %a_size to i64
	%13 = icmp ult i64 1, 0
	br i1 %13, label %14, label %16

14:
	store i32 2, i32* @.ErrorNumber
	%15 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

16:
	%17 = icmp uge i64 1, %12
	br i1 %17, label %14, label %18

18:
	%19 = getelementptr i32, i32* %a, i64 1
	%20 = load i32, i32* %19
	%21 = add i32 %11, %20
	%22 = sext i32 %a_size to i64
	%23 = icmp ult i64 2, 0
	br i1 %23, label %24, label %26

24:
	store i32 2, i32* @.ErrorNumber
	%25 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

26:
	%27 = icmp uge i64 2, %22
	br i1 %27, label %24, label %28

28:
	%29 = getelementptr i32, i32* %a, i64 2
	%30 = load i32, i32* %29
	%31 = add i32 %21, %30
	%32 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %32
	%33 = call i32 (i8*, ...) @printf([3 x i8]* %32, i32 %31)
	%34 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	ret void
}

define i32 @addDefaultValue(i32 %a, i32 %b) {
0:
	%1 = alloca i32
	store i32 %a, i32* %1
	%2 = alloca i32
	store i32 %b, i32* %2
	%3 = alloca i32
	store i32 0, i32* %3
	%4 = load i32, i32* %1
	%5 = load i32, i32* %2
	%6 = add i32 %4, %5
	store i32 %6, i32* %3
	%7 = load i32, i32* %3
	ret i32 %7
}

define void @addByRef(i32* %a) {
0:
	%1 = alloca i32*
	store i32* %a, i32** %1
	%2 = load i32*, i32** %1
	%3 = load i32, i32* %2
	%4 = sext i32 %3 to i64
	%5 = add i64 %4, 1
	%6 = trunc i64 %5 to i32
	%7 = load i32*, i32** %1
	store i32 %6, i32* %7
	%8 = alloca [16 x i8]
	store [16 x i8] c"in addByRef A: \00", [16 x i8]* %8
	%9 = call i32 (i8*, ...) @printf([16 x i8]* %8)
	%10 = load i32*, i32** %1
	%11 = load i32, i32* %10
	%12 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %12
	%13 = call i32 (i8*, ...) @printf([3 x i8]* %12, i32 %11)
	%14 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	ret void
}

define void @addByVal(i32 %a) {
0:
	%1 = alloca i32
	store i32 %a, i32* %1
	%2 = load i32, i32* %1
	%3 = sext i32 %2 to i64
	%4 = add i64 %3, 1
	%5 = trunc i64 %4 to i32
	store i32 %5, i32* %1
	%6 = alloca [16 x i8]
	store [16 x i8] c"in addByVal A: \00", [16 x i8]* %6
	%7 = call i32 (i8*, ...) @printf([16 x i8]* %6)
	%8 = load i32, i32* %1
	%9 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %9
	%10 = call i32 (i8*, ...) @printf([3 x i8]* %9, i32 %8)
	%11 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	ret void
}

define void @FloatAddByRef(float* %a) {
0:
	%1 = alloca float*
	store float* %a, float** %1
	%2 = load float*, float** %1
	%3 = load float, float* %2
	%4 = fpext float %3 to double
	%5 = fadd double %4, 0x3FEFAE147AE147AE
	%6 = fptrunc double %5 to float
	%7 = load float*, float** %1
	store float %6, float* %7
	%8 = alloca [21 x i8]
	store [21 x i8] c"in FloatAddByRef A: \00", [21 x i8]* %8
	%9 = call i32 (i8*, ...) @printf([21 x i8]* %8)
	%10 = load float*, float** %1
	%11 = load float, float* %10
	%12 = alloca [3 x i8]
	store [3 x i8] c"%f\00", [3 x i8]* %12
	%13 = call i32 (i8*, ...) @printf([3 x i8]* %12, float %11)
	%14 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	ret void
}

define void @FloatAddByVal(float %a) {
0:
	%1 = alloca float
	store float %a, float* %1
	%2 = load float, float* %1
	%3 = fpext float %2 to double
	%4 = fadd double %3, 4.0
	%5 = fptrunc double %4 to float
	store float %5, float* %1
	%6 = alloca [21 x i8]
	store [21 x i8] c"in FloatAddByVal A: \00", [21 x i8]* %6
	%7 = call i32 (i8*, ...) @printf([21 x i8]* %6)
	%8 = load float, float* %1
	%9 = alloca [3 x i8]
	store [3 x i8] c"%f\00", [3 x i8]* %9
	%10 = call i32 (i8*, ...) @printf([3 x i8]* %9, float %8)
	%11 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	ret void
}

define void @setLocalConstants() {
0:
	%currcl = alloca float
	store float 36.0, float* %currcl
	%loncl = alloca i64
	store i64 46, i64* %loncl
	%sngcl = alloca float
	store float 0x40632570A0000000, float* %sngcl
	%dblcl = alloca double
	store double 87.0, double* %dblcl
	%strcl = alloca [25 x i8]
	store [25 x i8] c"hello world! hello world\00", [25 x i8]* %strcl
	%intcl = alloca i32
	store i32 90, i32* %intcl
	%boolcl = alloca i1
	store i1 false, i1* %boolcl
	%datecl = alloca float
	store float 0x42012E0BE0000000, float* %datecl
	%1 = load float, float* %currcl
	%2 = alloca [3 x i8]
	store [3 x i8] c"%f\00", [3 x i8]* %2
	%3 = call i32 (i8*, ...) @printf([3 x i8]* %2, float %1)
	%4 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%5 = load i64, i64* %loncl
	%6 = alloca [4 x i8]
	store [4 x i8] c"%ld\00", [4 x i8]* %6
	%7 = call i32 (i8*, ...) @printf([4 x i8]* %6, i64 %5)
	%8 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%9 = load float, float* %sngcl
	%10 = alloca [3 x i8]
	store [3 x i8] c"%f\00", [3 x i8]* %10
	%11 = call i32 (i8*, ...) @printf([3 x i8]* %10, float %9)
	%12 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%13 = load double, double* %dblcl
	%14 = alloca [4 x i8]
	store [4 x i8] c"%lf\00", [4 x i8]* %14
	%15 = call i32 (i8*, ...) @printf([4 x i8]* %14, double %13)
	%16 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%17 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %17
	%18 = call i32 (i8*, ...) @printf([3 x i8]* %17, [25 x i8]* %strcl)
	%19 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%20 = load i32, i32* %intcl
	%21 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %21
	%22 = call i32 (i8*, ...) @printf([3 x i8]* %21, i32 %20)
	%23 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%24 = load i1, i1* %boolcl
	%25 = icmp eq i1 %24, true
	%26 = select i1 %25, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%27 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %27
	%28 = call i32 (i8*, ...) @printf([3 x i8]* %27, [5 x i8]* %26)
	%29 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%30 = load float, float* %datecl
	%31 = alloca [3 x i8]
	store [3 x i8] c"%f\00", [3 x i8]* %31
	%32 = call i32 (i8*, ...) @printf([3 x i8]* %31, float %30)
	%33 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	ret void
}

define void @declareLocalArrays() {
0:
	%arrayf = alloca [6 x i32]
	%1 = alloca i32
	store i32 0, i32* %1
	br label %2

2:
	%3 = load i32, i32* %1
	%4 = icmp ult i32 %3, 6
	br i1 %4, label %5, label %9

5:
	%6 = getelementptr i32, [6 x i32]* %arrayf, i32 %3
	store i32 0, i32* %6
	%7 = load i32, i32* %1
	%8 = add i32 %7, 1
	store i32 %8, i32* %1
	br label %2

9:
	%arrayd = alloca float*
	%10 = load i64, i64* @.declareLocalArrays_arrayf_0
	%11 = icmp ult i64 0, 0
	br i1 %11, label %12, label %14

12:
	store i32 2, i32* @.ErrorNumber
	%13 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

14:
	%15 = icmp uge i64 0, %10
	br i1 %15, label %12, label %16

16:
	%17 = load i64, i64* @.declareLocalArrays_arrayf_1
	%18 = icmp ult i64 1, 0
	br i1 %18, label %19, label %21

19:
	store i32 2, i32* @.ErrorNumber
	%20 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

21:
	%22 = icmp uge i64 1, %17
	br i1 %22, label %19, label %23

23:
	%24 = mul i64 0, %17
	%25 = add i64 %24, 1
	%26 = getelementptr [6 x i32], [6 x i32]* %arrayf, i64 0, i64 %25
	%27 = trunc i64 4 to i32
	store i32 %27, i32* %26
	%28 = load i64, i64* @.declareLocalArrays_arrayf_0
	%29 = icmp ult i64 0, 0
	br i1 %29, label %30, label %32

30:
	store i32 2, i32* @.ErrorNumber
	%31 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

32:
	%33 = icmp uge i64 0, %28
	br i1 %33, label %30, label %34

34:
	%35 = load i64, i64* @.declareLocalArrays_arrayf_1
	%36 = icmp ult i64 1, 0
	br i1 %36, label %37, label %39

37:
	store i32 2, i32* @.ErrorNumber
	%38 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

39:
	%40 = icmp uge i64 1, %35
	br i1 %40, label %37, label %41

41:
	%42 = mul i64 0, %35
	%43 = add i64 %42, 1
	%44 = getelementptr [6 x i32], [6 x i32]* %arrayf, i64 0, i64 %43
	%45 = load i32, i32* %44
	%46 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %46
	%47 = call i32 (i8*, ...) @printf([3 x i8]* %46, i32 %45)
	%48 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%49 = load i64, i64* @.declareLocalArrays_arrayf_0
	%50 = icmp ult i64 1, 0
	br i1 %50, label %51, label %53

51:
	store i32 2, i32* @.ErrorNumber
	%52 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

53:
	%54 = icmp uge i64 1, %49
	br i1 %54, label %51, label %55

55:
	%56 = load i64, i64* @.declareLocalArrays_arrayf_1
	%57 = icmp ult i64 1, 0
	br i1 %57, label %58, label %60

58:
	store i32 2, i32* @.ErrorNumber
	%59 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

60:
	%61 = icmp uge i64 1, %56
	br i1 %61, label %58, label %62

62:
	%63 = mul i64 1, %56
	%64 = add i64 %63, 1
	%65 = getelementptr [6 x i32], [6 x i32]* %arrayf, i64 0, i64 %64
	%66 = load i32, i32* %65
	%67 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %67
	%68 = call i32 (i8*, ...) @printf([3 x i8]* %67, i32 %66)
	%69 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
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
	store i64 10, i64* @l
	store double 0x400921FAFC8B007A, double* @d
	store float 0x42012E0BE0000000, float* @da
	%1 = call i32 @strlen([14 x i8]* @_1)
	%2 = call i8* @gc_malloc(%struct.GarbageCollector* @gc, i32 %1)
	store i8* %2, i8** @s0
	%3 = call i8* @strcpy(i8* %2, [14 x i8]* @_1)
	store float 100.0, float* @c
	store i1 true, i1* @b
	%4 = alloca [11 x i8]
	store [11 x i8] c"variables:\00", [11 x i8]* %4
	%5 = call i32 (i8*, ...) @printf([11 x i8]* %4)
	%6 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%7 = load i64, i64* @l
	%8 = alloca [4 x i8]
	store [4 x i8] c"%ld\00", [4 x i8]* %8
	%9 = call i32 (i8*, ...) @printf([4 x i8]* %8, i64 %7)
	%10 = alloca [3 x i8]
	store [3 x i8] c", \00", [3 x i8]* %10
	%11 = call i32 (i8*, ...) @printf([3 x i8]* %10)
	%12 = load double, double* @d
	%13 = alloca [4 x i8]
	store [4 x i8] c"%lf\00", [4 x i8]* %13
	%14 = call i32 (i8*, ...) @printf([4 x i8]* %13, double %12)
	%15 = alloca [3 x i8]
	store [3 x i8] c", \00", [3 x i8]* %15
	%16 = call i32 (i8*, ...) @printf([3 x i8]* %15)
	%17 = load i1, i1* @b
	%18 = icmp eq i1 %17, true
	%19 = select i1 %18, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%20 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %20
	%21 = call i32 (i8*, ...) @printf([3 x i8]* %20, [5 x i8]* %19)
	%22 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%23 = load float, float* @da
	%24 = alloca [3 x i8]
	store [3 x i8] c"%f\00", [3 x i8]* %24
	%25 = call i32 (i8*, ...) @printf([3 x i8]* %24, float %23)
	%26 = alloca [3 x i8]
	store [3 x i8] c", \00", [3 x i8]* %26
	%27 = call i32 (i8*, ...) @printf([3 x i8]* %26)
	%28 = load i8*, i8** getelementptr (i8*, i8** @s0, i32 0)
	%29 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %29
	%30 = call i32 (i8*, ...) @printf([3 x i8]* %29, i8* %28)
	%31 = alloca [3 x i8]
	store [3 x i8] c", \00", [3 x i8]* %31
	%32 = call i32 (i8*, ...) @printf([3 x i8]* %31)
	%33 = load float, float* @c
	%34 = alloca [3 x i8]
	store [3 x i8] c"%f\00", [3 x i8]* %34
	%35 = call i32 (i8*, ...) @printf([3 x i8]* %34, float %33)
	%36 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%37 = alloca [11 x i8]
	store [11 x i8] c"constants:\00", [11 x i8]* %37
	%38 = call i32 (i8*, ...) @printf([11 x i8]* %37)
	%39 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%40 = alloca [4 x i8]
	store [4 x i8] c"c1:\00", [4 x i8]* %40
	%41 = call i32 (i8*, ...) @printf([4 x i8]* %40)
	%42 = load i64, i64* @c1
	%43 = alloca [4 x i8]
	store [4 x i8] c"%ld\00", [4 x i8]* %43
	%44 = call i32 (i8*, ...) @printf([4 x i8]* %43, i64 %42)
	%45 = alloca [6 x i8]
	store [6 x i8] c", c2:\00", [6 x i8]* %45
	%46 = call i32 (i8*, ...) @printf([6 x i8]* %45)
	%47 = load double, double* @c2
	%48 = alloca [4 x i8]
	store [4 x i8] c"%lf\00", [4 x i8]* %48
	%49 = call i32 (i8*, ...) @printf([4 x i8]* %48, double %47)
	%50 = alloca [6 x i8]
	store [6 x i8] c", c3:\00", [6 x i8]* %50
	%51 = call i32 (i8*, ...) @printf([6 x i8]* %50)
	%52 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %52
	%53 = call i32 (i8*, ...) @printf([3 x i8]* %52, [6 x i8]* @c3)
	%54 = alloca [6 x i8]
	store [6 x i8] c", c4:\00", [6 x i8]* %54
	%55 = call i32 (i8*, ...) @printf([6 x i8]* %54)
	%56 = load i1, i1* @c4
	%57 = icmp eq i1 %56, true
	%58 = select i1 %57, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%59 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %59
	%60 = call i32 (i8*, ...) @printf([3 x i8]* %59, [5 x i8]* %58)
	%61 = alloca [6 x i8]
	store [6 x i8] c", c5:\00", [6 x i8]* %61
	%62 = call i32 (i8*, ...) @printf([6 x i8]* %61)
	%63 = load float, float* @c5
	%64 = alloca [3 x i8]
	store [3 x i8] c"%f\00", [3 x i8]* %64
	%65 = call i32 (i8*, ...) @printf([3 x i8]* %64, float %63)
	%66 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%67 = alloca [10 x i8]
	store [10 x i8] c"literals:\00", [10 x i8]* %67
	%68 = call i32 (i8*, ...) @printf([10 x i8]* %67)
	%69 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%70 = alloca [6 x i8]
	store [6 x i8] c"world\00", [6 x i8]* %70
	%71 = call i32 (i8*, ...) @printf([6 x i8]* %70)
	%72 = alloca [5 x i8]
	store [5 x i8] c"True\00", [5 x i8]* %72
	%73 = call i32 (i8*, ...) @printf([5 x i8]* %72)
	%74 = alloca [8 x i8]
	store [8 x i8] c"1.23456\00", [8 x i8]* %74
	%75 = call i32 (i8*, ...) @printf([8 x i8]* %74)
	%76 = alloca [3 x i8]
	store [3 x i8] c"10\00", [3 x i8]* %76
	%77 = call i32 (i8*, ...) @printf([3 x i8]* %76)
	%78 = alloca [11 x i8]
	store [11 x i8] c"2010/12/31\00", [11 x i8]* %78
	%79 = call i32 (i8*, ...) @printf([11 x i8]* %78)
	%80 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%81 = trunc i64 10 to i32
	store i32 %81, i32* @i
	%82 = fptrunc double 0x405900A3D70A3D71 to float
	store float %82, float* @s
	%83 = load i64, i32* @i
	store i64 %83, i64* @l
	%84 = load double, float* @s
	store double %84, double* @d
	%85 = load i32, i32* @i
	%86 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %86
	%87 = call i32 (i8*, ...) @printf([3 x i8]* %86, i32 %85)
	%88 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%89 = load float, float* @s
	%90 = alloca [3 x i8]
	store [3 x i8] c"%f\00", [3 x i8]* %90
	%91 = call i32 (i8*, ...) @printf([3 x i8]* %90, float %89)
	%92 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%93 = load i64, i64* @l
	%94 = alloca [4 x i8]
	store [4 x i8] c"%ld\00", [4 x i8]* %94
	%95 = call i32 (i8*, ...) @printf([4 x i8]* %94, i64 %93)
	%96 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%97 = load double, double* @d
	%98 = alloca [4 x i8]
	store [4 x i8] c"%lf\00", [4 x i8]* %98
	%99 = call i32 (i8*, ...) @printf([4 x i8]* %98, double %97)
	%100 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%101 = trunc i64 1234 to i32
	store i32 %101, i32* @i
	store i64 1234567890, i64* @l
	%102 = load i32, i32* @i
	%103 = sext i32 %102 to i64
	%104 = add i64 %103, 1
	%105 = sitofp i64 %104 to float
	store float %105, float* @s
	%106 = load i64, i64* @l
	%107 = add i64 %106, 1
	%108 = sitofp i64 %107 to double
	store double %108, double* @d
	%109 = load i32, i32* @i
	%110 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %110
	%111 = call i32 (i8*, ...) @printf([3 x i8]* %110, i32 %109)
	%112 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%113 = load float, float* @s
	%114 = alloca [3 x i8]
	store [3 x i8] c"%f\00", [3 x i8]* %114
	%115 = call i32 (i8*, ...) @printf([3 x i8]* %114, float %113)
	%116 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%117 = load i64, i64* @l
	%118 = alloca [4 x i8]
	store [4 x i8] c"%ld\00", [4 x i8]* %118
	%119 = call i32 (i8*, ...) @printf([4 x i8]* %118, i64 %117)
	%120 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%121 = load double, double* @d
	%122 = alloca [4 x i8]
	store [4 x i8] c"%lf\00", [4 x i8]* %122
	%123 = call i32 (i8*, ...) @printf([4 x i8]* %122, double %121)
	%124 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%125 = fptrunc double 0x4028AE147AE147AE to float
	store float %125, float* @s
	store double 0x41678C29DCCCCCCD, double* @d
	%126 = load i32, float* @s
	store i32 %126, i32* @i
	%127 = load i64, double* @d
	store i64 %127, i64* @l
	%128 = load i32, i32* @i
	%129 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %129
	%130 = call i32 (i8*, ...) @printf([3 x i8]* %129, i32 %128)
	%131 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%132 = load float, float* @s
	%133 = alloca [3 x i8]
	store [3 x i8] c"%f\00", [3 x i8]* %133
	%134 = call i32 (i8*, ...) @printf([3 x i8]* %133, float %132)
	%135 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%136 = load i64, i64* @l
	%137 = alloca [4 x i8]
	store [4 x i8] c"%ld\00", [4 x i8]* %137
	%138 = call i32 (i8*, ...) @printf([4 x i8]* %137, i64 %136)
	%139 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%140 = load double, double* @d
	%141 = alloca [4 x i8]
	store [4 x i8] c"%lf\00", [4 x i8]* %141
	%142 = call i32 (i8*, ...) @printf([4 x i8]* %141, double %140)
	%143 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	store double 0x41678C29DCCCCCCD, double* @d
	%144 = load float, double* @d
	store float %144, float* @s
	%145 = load i64, float* @s
	store i64 %145, i64* @l
	%146 = load i32, i64* @l
	store i32 %146, i32* @i
	%147 = load i32, i32* @i
	%148 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %148
	%149 = call i32 (i8*, ...) @printf([3 x i8]* %148, i32 %147)
	%150 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%151 = load float, float* @s
	%152 = alloca [3 x i8]
	store [3 x i8] c"%f\00", [3 x i8]* %152
	%153 = call i32 (i8*, ...) @printf([3 x i8]* %152, float %151)
	%154 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%155 = load i64, i64* @l
	%156 = alloca [4 x i8]
	store [4 x i8] c"%ld\00", [4 x i8]* %156
	%157 = call i32 (i8*, ...) @printf([4 x i8]* %156, i64 %155)
	%158 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%159 = load double, double* @d
	%160 = alloca [4 x i8]
	store [4 x i8] c"%lf\00", [4 x i8]* %160
	%161 = call i32 (i8*, ...) @printf([4 x i8]* %160, double %159)
	%162 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	call void @ifBranch()
	%163 = call i32 @strlen([6 x i8]* @_2)
	%164 = call i8* @gc_malloc(%struct.GarbageCollector* @gc, i32 %163)
	store i8* %164, i8** @s1
	%165 = call i8* @strcpy(i8* %164, [6 x i8]* @_2)
	%166 = call i32 @strlen([6 x i8]* @_3)
	%167 = call i8* @gc_malloc(%struct.GarbageCollector* @gc, i32 %166)
	store i8* %167, i8** @s2
	%168 = call i8* @strcpy(i8* %167, [6 x i8]* @_3)
	%169 = call i32 @strlen(i8** @s1)
	%170 = call i32 @strlen([2 x i8]* @_4)
	%171 = add i32 %169, %170
	%172 = call i8* @gc_malloc(%struct.GarbageCollector* @gc, i32 %171)
	%173 = call i8* @strcpy(i8* %172, i8** @s1)
	%174 = call i8* @strcat(i8* %172, [2 x i8]* @_4)
	%175 = call i32 @strlen(i8* %172)
	%176 = call i32 @strlen(i8** @s2)
	%177 = add i32 %175, %176
	%178 = call i8* @gc_malloc(%struct.GarbageCollector* @gc, i32 %177)
	%179 = call i8* @strcpy(i8* %178, i8* %172)
	%180 = call i8* @strcat(i8* %178, i8** @s2)
	%181 = call i32 @strlen(i8* %178)
	%182 = call i32 @strlen([3 x i8]* @_5)
	%183 = add i32 %181, %182
	%184 = call i8* @gc_malloc(%struct.GarbageCollector* @gc, i32 %183)
	%185 = call i8* @strcpy(i8* %184, i8* %178)
	%186 = call i8* @strcat(i8* %184, [3 x i8]* @_5)
	%187 = call i32 @strlen(i8* %184)
	%188 = call i8* @gc_malloc(%struct.GarbageCollector* @gc, i32 %187)
	store i8* %188, i8** @s3
	%189 = call i8* @strcpy(i8* %188, i8* %184)
	%190 = load i8*, i8** getelementptr (i8*, i8** @s3, i32 0)
	%191 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %191
	%192 = call i32 (i8*, ...) @printf([3 x i8]* %191, i8* %190)
	%193 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%194 = call i32 @strlen([7 x i8]* @_6)
	%195 = call i32 @strlen([8 x i8]* @_7)
	%196 = add i32 %194, %195
	%197 = call i8* @gc_malloc(%struct.GarbageCollector* @gc, i32 %196)
	%198 = call i8* @strcpy(i8* %197, [7 x i8]* @_6)
	%199 = call i8* @strcat(i8* %197, [8 x i8]* @_7)
	%200 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %200
	%201 = call i32 (i8*, ...) @printf([3 x i8]* %200, i8* %197)
	%202 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%203 = add i64 1, 2
	%204 = alloca [4 x i8]
	store [4 x i8] c"%ld\00", [4 x i8]* %204
	%205 = call i32 (i8*, ...) @printf([4 x i8]* %204, i64 %203)
	%206 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%207 = fmul double 0x4025CCCCCCCCCCCD, 0x3FEF5C28F5C28F5C
	%208 = alloca [4 x i8]
	store [4 x i8] c"%lf\00", [4 x i8]* %208
	%209 = call i32 (i8*, ...) @printf([4 x i8]* %208, double %207)
	%210 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%211 = sitofp i64 1 to double
	%212 = fadd double %211, 0x4007D70A3D70A3D7
	%213 = alloca [4 x i8]
	store [4 x i8] c"%lf\00", [4 x i8]* %213
	%214 = call i32 (i8*, ...) @printf([4 x i8]* %213, double %212)
	%215 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	br i1 true, label %216, label %217

216:
	br label %217

217:
	%218 = phi i1 [ false, %0 ], [ false, %216 ]
	%219 = icmp eq i1 %218, true
	%220 = select i1 %219, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%221 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %221
	%222 = call i32 (i8*, ...) @printf([3 x i8]* %221, [5 x i8]* %220)
	%223 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	br i1 true, label %224, label %225

224:
	br label %225

225:
	%226 = phi i1 [ false, %217 ], [ true, %224 ]
	%227 = icmp eq i1 %226, true
	%228 = select i1 %227, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%229 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %229
	%230 = call i32 (i8*, ...) @printf([3 x i8]* %229, [5 x i8]* %228)
	%231 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	br i1 false, label %232, label %233

232:
	br label %233

233:
	%234 = phi i1 [ false, %225 ], [ false, %232 ]
	%235 = icmp eq i1 %234, true
	%236 = select i1 %235, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%237 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %237
	%238 = call i32 (i8*, ...) @printf([3 x i8]* %237, [5 x i8]* %236)
	%239 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	br i1 false, label %240, label %241

240:
	br label %241

241:
	%242 = phi i1 [ false, %233 ], [ true, %240 ]
	%243 = icmp eq i1 %242, true
	%244 = select i1 %243, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%245 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %245
	%246 = call i32 (i8*, ...) @printf([3 x i8]* %245, [5 x i8]* %244)
	%247 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	br i1 true, label %249, label %248

248:
	br label %249

249:
	%250 = phi i1 [ true, %241 ], [ false, %248 ]
	%251 = icmp eq i1 %250, true
	%252 = select i1 %251, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%253 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %253
	%254 = call i32 (i8*, ...) @printf([3 x i8]* %253, [5 x i8]* %252)
	%255 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	br i1 true, label %257, label %256

256:
	br label %257

257:
	%258 = phi i1 [ true, %249 ], [ true, %256 ]
	%259 = icmp eq i1 %258, true
	%260 = select i1 %259, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%261 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %261
	%262 = call i32 (i8*, ...) @printf([3 x i8]* %261, [5 x i8]* %260)
	%263 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	br i1 false, label %265, label %264

264:
	br label %265

265:
	%266 = phi i1 [ true, %257 ], [ false, %264 ]
	%267 = icmp eq i1 %266, true
	%268 = select i1 %267, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%269 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %269
	%270 = call i32 (i8*, ...) @printf([3 x i8]* %269, [5 x i8]* %268)
	%271 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	br i1 false, label %273, label %272

272:
	br label %273

273:
	%274 = phi i1 [ true, %265 ], [ true, %272 ]
	%275 = icmp eq i1 %274, true
	%276 = select i1 %275, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%277 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %277
	%278 = call i32 (i8*, ...) @printf([3 x i8]* %277, [5 x i8]* %276)
	%279 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%280 = xor i1 true, true
	%281 = icmp eq i1 %280, true
	%282 = select i1 %281, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%283 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %283
	%284 = call i32 (i8*, ...) @printf([3 x i8]* %283, [5 x i8]* %282)
	%285 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%286 = fptosi float 2.25 to i64
	%287 = mul i64 2, %286
	%288 = alloca [3 x i8]
	store [3 x i8] c"%f\00", [3 x i8]* %288
	%289 = call i32 (i8*, ...) @printf([3 x i8]* %288, i64 %287)
	%290 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%291 = icmp eq i64 2, 0
	br i1 %291, label %292, label %294

292:
	store i32 1, i32* @.ErrorNumber
	%293 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [18 x i8]* @.divisionByZero)
	call void @.throwException()
	unreachable

294:
	%295 = sdiv i64 2, 2
	%296 = alloca [4 x i8]
	store [4 x i8] c"%ld\00", [4 x i8]* %296
	%297 = call i32 (i8*, ...) @printf([4 x i8]* %296, i64 %295)
	%298 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%299 = icmp slt i64 2, 3
	%300 = icmp eq i1 %299, true
	%301 = select i1 %300, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%302 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %302
	%303 = call i32 (i8*, ...) @printf([3 x i8]* %302, [5 x i8]* %301)
	%304 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%305 = icmp sgt i64 2, 3
	%306 = icmp eq i1 %305, true
	%307 = select i1 %306, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%308 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %308
	%309 = call i32 (i8*, ...) @printf([3 x i8]* %308, [5 x i8]* %307)
	%310 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%311 = icmp sle i64 2, 3
	%312 = icmp eq i1 %311, true
	%313 = select i1 %312, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%314 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %314
	%315 = call i32 (i8*, ...) @printf([3 x i8]* %314, [5 x i8]* %313)
	%316 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%317 = icmp sge i64 2, 3
	%318 = icmp eq i1 %317, true
	%319 = select i1 %318, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%320 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %320
	%321 = call i32 (i8*, ...) @printf([3 x i8]* %320, [5 x i8]* %319)
	%322 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%323 = icmp eq i64 2, 3
	%324 = icmp eq i1 %323, true
	%325 = select i1 %324, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%326 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %326
	%327 = call i32 (i8*, ...) @printf([3 x i8]* %326, [5 x i8]* %325)
	%328 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%329 = icmp ne i64 2, 3
	%330 = icmp eq i1 %329, true
	%331 = select i1 %330, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%332 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %332
	%333 = call i32 (i8*, ...) @printf([3 x i8]* %332, [5 x i8]* %331)
	%334 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%335 = fcmp ogt double 2.5, 3.0
	%336 = icmp eq i1 %335, true
	%337 = select i1 %336, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%338 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %338
	%339 = call i32 (i8*, ...) @printf([3 x i8]* %338, [5 x i8]* %337)
	%340 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%341 = fcmp olt double 2.5, 3.0
	%342 = icmp eq i1 %341, true
	%343 = select i1 %342, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%344 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %344
	%345 = call i32 (i8*, ...) @printf([3 x i8]* %344, [5 x i8]* %343)
	%346 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%347 = fcmp ole double 2.5, 3.0
	%348 = icmp eq i1 %347, true
	%349 = select i1 %348, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%350 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %350
	%351 = call i32 (i8*, ...) @printf([3 x i8]* %350, [5 x i8]* %349)
	%352 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%353 = fcmp oge double 2.5, 3.0
	%354 = icmp eq i1 %353, true
	%355 = select i1 %354, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%356 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %356
	%357 = call i32 (i8*, ...) @printf([3 x i8]* %356, [5 x i8]* %355)
	%358 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%359 = fcmp oeq double 2.5, 3.0
	%360 = icmp eq i1 %359, true
	%361 = select i1 %360, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%362 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %362
	%363 = call i32 (i8*, ...) @printf([3 x i8]* %362, [5 x i8]* %361)
	%364 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%365 = fcmp one double 2.5, 3.0
	%366 = icmp eq i1 %365, true
	%367 = select i1 %366, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%368 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %368
	%369 = call i32 (i8*, ...) @printf([3 x i8]* %368, [5 x i8]* %367)
	%370 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%371 = trunc i64 10 to i32
	store i32 %371, i32* @i11
	br label %372

372:
	%373 = load i32, i32* @i11
	%374 = sext i32 %373 to i64
	%375 = icmp sgt i64 %374, 0
	br i1 %375, label %376, label %385

376:
	%377 = load i32, i32* @i11
	%378 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %378
	%379 = call i32 (i8*, ...) @printf([3 x i8]* %378, i32 %377)
	%380 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%381 = load i32, i32* @i11
	%382 = sext i32 %381 to i64
	%383 = sub i64 %382, 1
	%384 = trunc i64 %383 to i32
	store i32 %384, i32* @i11
	br label %372

385:
	call void @while1()
	%386 = trunc i64 2 to i32
	%387 = add i64 10, 2
	%388 = call i64 @times(i32 %386, i64 %387)
	%389 = alloca [4 x i8]
	store [4 x i8] c"%ld\00", [4 x i8]* %389
	%390 = call i32 (i8*, ...) @printf([4 x i8]* %389, i64 %388)
	%391 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%392 = trunc i64 3 to i32
	%393 = call i64 @times(i32 %392, i64 5)
	%394 = alloca [4 x i8]
	store [4 x i8] c"%ld\00", [4 x i8]* %394
	%395 = call i32 (i8*, ...) @printf([4 x i8]* %394, i64 %393)
	%396 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%397 = alloca [3 x i32]
	%398 = trunc i64 2 to i32
	%399 = getelementptr i32, [3 x i32]* %397, i64 0
	store i32 %398, i32* %399
	%400 = trunc i64 5 to i32
	%401 = getelementptr i32, [3 x i32]* %397, i64 1
	store i32 %400, i32* %401
	%402 = trunc i64 9 to i32
	%403 = getelementptr i32, [3 x i32]* %397, i64 2
	store i32 %402, i32* %403
	call void @listParmArray([3 x i32]* %397, i64 3)
	%404 = alloca [3 x i32]
	%405 = trunc i64 17 to i32
	%406 = getelementptr i32, [3 x i32]* %404, i64 0
	store i32 %405, i32* %406
	%407 = sub i64 0, 19
	%408 = trunc i64 %407 to i32
	%409 = getelementptr i32, [3 x i32]* %404, i64 1
	store i32 %408, i32* %409
	%410 = trunc i64 25 to i32
	%411 = getelementptr i32, [3 x i32]* %404, i64 2
	store i32 %410, i32* %411
	call void @addParmArray([3 x i32]* %404, i64 3)
	%412 = trunc i64 10 to i32
	%413 = trunc i64 4 to i32
	%414 = call i32 @addDefaultValue(i32 %412, i32 %413)
	%415 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %415
	%416 = call i32 (i8*, ...) @printf([3 x i8]* %415, i32 %414)
	%417 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%418 = trunc i64 10 to i32
	%419 = trunc i64 2 to i32
	%420 = call i32 @addDefaultValue(i32 %418, i32 %419)
	%421 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %421
	%422 = call i32 (i8*, ...) @printf([3 x i8]* %421, i32 %420)
	%423 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%424 = trunc i64 1 to i32
	%425 = trunc i64 2 to i32
	%426 = call i32 @addDefaultValue(i32 %424, i32 %425)
	%427 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %427
	%428 = call i32 (i8*, ...) @printf([3 x i8]* %427, i32 %426)
	%429 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%430 = trunc i64 1 to i32
	store i32 %430, i32* @addBy
	%431 = load i32, i32* @addBy
	%432 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %432
	%433 = call i32 (i8*, ...) @printf([3 x i8]* %432, i32 %431)
	%434 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	call void @addByRef(i32* @addBy)
	%435 = load i32, i32* @addBy
	%436 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %436
	%437 = call i32 (i8*, ...) @printf([3 x i8]* %436, i32 %435)
	%438 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%439 = load i32, i32* @addBy
	%440 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %440
	%441 = call i32 (i8*, ...) @printf([3 x i8]* %440, i32 %439)
	%442 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%443 = load i32, i32* @addBy
	call void @addByVal(i32 %443)
	%444 = load i32, i32* @addBy
	%445 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %445
	%446 = call i32 (i8*, ...) @printf([3 x i8]* %445, i32 %444)
	%447 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%448 = fptrunc double 3.0 to float
	store float %448, float* @FloatAddBy
	%449 = load float, float* @FloatAddBy
	%450 = alloca [3 x i8]
	store [3 x i8] c"%f\00", [3 x i8]* %450
	%451 = call i32 (i8*, ...) @printf([3 x i8]* %450, float %449)
	%452 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	call void @FloatAddByRef(float* @FloatAddBy)
	%453 = load float, float* @FloatAddBy
	%454 = alloca [3 x i8]
	store [3 x i8] c"%f\00", [3 x i8]* %454
	%455 = call i32 (i8*, ...) @printf([3 x i8]* %454, float %453)
	%456 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%457 = load float, float* @FloatAddBy
	%458 = alloca [3 x i8]
	store [3 x i8] c"%f\00", [3 x i8]* %458
	%459 = call i32 (i8*, ...) @printf([3 x i8]* %458, float %457)
	%460 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%461 = load float, float* @FloatAddBy
	call void @FloatAddByVal(float %461)
	%462 = load float, float* @FloatAddBy
	%463 = alloca [3 x i8]
	store [3 x i8] c"%f\00", [3 x i8]* %463
	%464 = call i32 (i8*, ...) @printf([3 x i8]* %463, float %462)
	%465 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%466 = sitofp i64 3 to double
	%467 = fcmp oeq double %466, 0.0
	br i1 %467, label %468, label %470

468:
	store i32 1, i32* @.ErrorNumber
	%469 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [18 x i8]* @.divisionByZero)
	call void @.throwException()
	unreachable

470:
	%471 = fdiv double 1.0, %466
	%472 = alloca [4 x i8]
	store [4 x i8] c"%lf\00", [4 x i8]* %472
	%473 = call i32 (i8*, ...) @printf([4 x i8]* %472, double %471)
	%474 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%475 = sitofp i64 3 to double
	%476 = fcmp oeq double 2.0, 0.0
	br i1 %476, label %477, label %479

477:
	store i32 1, i32* @.ErrorNumber
	%478 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [18 x i8]* @.divisionByZero)
	call void @.throwException()
	unreachable

479:
	%480 = fdiv double %475, 2.0
	%481 = alloca [4 x i8]
	store [4 x i8] c"%lf\00", [4 x i8]* %481
	%482 = call i32 (i8*, ...) @printf([4 x i8]* %481, double %480)
	%483 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%484 = icmp eq i64 2, 0
	br i1 %484, label %485, label %487

485:
	store i32 1, i32* @.ErrorNumber
	%486 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [18 x i8]* @.divisionByZero)
	call void @.throwException()
	unreachable

487:
	%488 = sdiv i64 3, 2
	%489 = alloca [4 x i8]
	store [4 x i8] c"%ld\00", [4 x i8]* %489
	%490 = call i32 (i8*, ...) @printf([4 x i8]* %489, i64 %488)
	%491 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%492 = load i64, i64* @.global_b8_0
	%493 = icmp ult i64 0, 0
	br i1 %493, label %494, label %496

494:
	store i32 2, i32* @.ErrorNumber
	%495 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

496:
	%497 = icmp uge i64 0, %492
	br i1 %497, label %494, label %498

498:
	%499 = load i64, i64* @.global_b8_1
	%500 = icmp ult i64 0, 0
	br i1 %500, label %501, label %503

501:
	store i32 2, i32* @.ErrorNumber
	%502 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

503:
	%504 = icmp uge i64 0, %499
	br i1 %504, label %501, label %505

505:
	%506 = mul i64 0, %499
	%507 = add i64 %506, 0
	%508 = getelementptr [9 x i32], [9 x i32]* @b8, i64 0, i64 %507
	%509 = trunc i64 2 to i32
	store i32 %509, i32* %508
	%510 = load i64, i64* @.global_b8_0
	%511 = icmp ult i64 0, 0
	br i1 %511, label %512, label %514

512:
	store i32 2, i32* @.ErrorNumber
	%513 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

514:
	%515 = icmp uge i64 0, %510
	br i1 %515, label %512, label %516

516:
	%517 = load i64, i64* @.global_b8_1
	%518 = icmp ult i64 1, 0
	br i1 %518, label %519, label %521

519:
	store i32 2, i32* @.ErrorNumber
	%520 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

521:
	%522 = icmp uge i64 1, %517
	br i1 %522, label %519, label %523

523:
	%524 = mul i64 0, %517
	%525 = add i64 %524, 1
	%526 = getelementptr [9 x i32], [9 x i32]* @b8, i64 0, i64 %525
	%527 = trunc i64 3 to i32
	store i32 %527, i32* %526
	%528 = load i64, i64* @.global_b8_0
	%529 = icmp ult i64 0, 0
	br i1 %529, label %530, label %532

530:
	store i32 2, i32* @.ErrorNumber
	%531 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

532:
	%533 = icmp uge i64 0, %528
	br i1 %533, label %530, label %534

534:
	%535 = load i64, i64* @.global_b8_1
	%536 = icmp ult i64 2, 0
	br i1 %536, label %537, label %539

537:
	store i32 2, i32* @.ErrorNumber
	%538 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

539:
	%540 = icmp uge i64 2, %535
	br i1 %540, label %537, label %541

541:
	%542 = mul i64 0, %535
	%543 = add i64 %542, 2
	%544 = getelementptr [9 x i32], [9 x i32]* @b8, i64 0, i64 %543
	%545 = trunc i64 4 to i32
	store i32 %545, i32* %544
	%546 = load i64, i64* @.global_b8_0
	%547 = icmp ult i64 1, 0
	br i1 %547, label %548, label %550

548:
	store i32 2, i32* @.ErrorNumber
	%549 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

550:
	%551 = icmp uge i64 1, %546
	br i1 %551, label %548, label %552

552:
	%553 = load i64, i64* @.global_b8_1
	%554 = icmp ult i64 0, 0
	br i1 %554, label %555, label %557

555:
	store i32 2, i32* @.ErrorNumber
	%556 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

557:
	%558 = icmp uge i64 0, %553
	br i1 %558, label %555, label %559

559:
	%560 = mul i64 1, %553
	%561 = add i64 %560, 0
	%562 = getelementptr [9 x i32], [9 x i32]* @b8, i64 0, i64 %561
	%563 = trunc i64 5 to i32
	store i32 %563, i32* %562
	%564 = load i64, i64* @.global_b8_0
	%565 = icmp ult i64 1, 0
	br i1 %565, label %566, label %568

566:
	store i32 2, i32* @.ErrorNumber
	%567 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

568:
	%569 = icmp uge i64 1, %564
	br i1 %569, label %566, label %570

570:
	%571 = load i64, i64* @.global_b8_1
	%572 = icmp ult i64 1, 0
	br i1 %572, label %573, label %575

573:
	store i32 2, i32* @.ErrorNumber
	%574 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

575:
	%576 = icmp uge i64 1, %571
	br i1 %576, label %573, label %577

577:
	%578 = mul i64 1, %571
	%579 = add i64 %578, 1
	%580 = getelementptr [9 x i32], [9 x i32]* @b8, i64 0, i64 %579
	%581 = trunc i64 6 to i32
	store i32 %581, i32* %580
	%582 = load i64, i64* @.global_b8_0
	%583 = icmp ult i64 1, 0
	br i1 %583, label %584, label %586

584:
	store i32 2, i32* @.ErrorNumber
	%585 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

586:
	%587 = icmp uge i64 1, %582
	br i1 %587, label %584, label %588

588:
	%589 = load i64, i64* @.global_b8_1
	%590 = icmp ult i64 1, 0
	br i1 %590, label %591, label %593

591:
	store i32 2, i32* @.ErrorNumber
	%592 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

593:
	%594 = icmp uge i64 1, %589
	br i1 %594, label %591, label %595

595:
	%596 = mul i64 1, %589
	%597 = add i64 %596, 1
	%598 = getelementptr [9 x i32], [9 x i32]* @b8, i64 0, i64 %597
	%599 = trunc i64 7 to i32
	store i32 %599, i32* %598
	%600 = load i64, i64* @.global_b8_0
	%601 = icmp ult i64 1, 0
	br i1 %601, label %602, label %604

602:
	store i32 2, i32* @.ErrorNumber
	%603 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

604:
	%605 = icmp uge i64 1, %600
	br i1 %605, label %602, label %606

606:
	%607 = load i64, i64* @.global_b8_1
	%608 = icmp ult i64 2, 0
	br i1 %608, label %609, label %611

609:
	store i32 2, i32* @.ErrorNumber
	%610 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

611:
	%612 = icmp uge i64 2, %607
	br i1 %612, label %609, label %613

613:
	%614 = mul i64 1, %607
	%615 = add i64 %614, 2
	%616 = getelementptr [9 x i32], [9 x i32]* @b8, i64 0, i64 %615
	%617 = trunc i64 8 to i32
	store i32 %617, i32* %616
	%618 = load i64, i64* @.global_b8_0
	%619 = icmp ult i64 2, 0
	br i1 %619, label %620, label %622

620:
	store i32 2, i32* @.ErrorNumber
	%621 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

622:
	%623 = icmp uge i64 2, %618
	br i1 %623, label %620, label %624

624:
	%625 = load i64, i64* @.global_b8_1
	%626 = icmp ult i64 1, 0
	br i1 %626, label %627, label %629

627:
	store i32 2, i32* @.ErrorNumber
	%628 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

629:
	%630 = icmp uge i64 1, %625
	br i1 %630, label %627, label %631

631:
	%632 = mul i64 2, %625
	%633 = add i64 %632, 1
	%634 = getelementptr [9 x i32], [9 x i32]* @b8, i64 0, i64 %633
	%635 = trunc i64 9 to i32
	store i32 %635, i32* %634
	%636 = load i64, i64* @.global_b8_0
	%637 = icmp ult i64 1, 0
	br i1 %637, label %638, label %640

638:
	store i32 2, i32* @.ErrorNumber
	%639 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

640:
	%641 = icmp uge i64 1, %636
	br i1 %641, label %638, label %642

642:
	%643 = load i64, i64* @.global_b8_1
	%644 = icmp ult i64 2, 0
	br i1 %644, label %645, label %647

645:
	store i32 2, i32* @.ErrorNumber
	%646 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

647:
	%648 = icmp uge i64 2, %643
	br i1 %648, label %645, label %649

649:
	%650 = mul i64 1, %643
	%651 = add i64 %650, 2
	%652 = getelementptr [9 x i32], [9 x i32]* @b8, i64 0, i64 %651
	%653 = load i32, i32* %652
	%654 = load i64, i64* @.global_b8_0
	%655 = icmp ult i64 2, 0
	br i1 %655, label %656, label %658

656:
	store i32 2, i32* @.ErrorNumber
	%657 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

658:
	%659 = icmp uge i64 2, %654
	br i1 %659, label %656, label %660

660:
	%661 = load i64, i64* @.global_b8_1
	%662 = icmp ult i64 1, 0
	br i1 %662, label %663, label %665

663:
	store i32 2, i32* @.ErrorNumber
	%664 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

665:
	%666 = icmp uge i64 1, %661
	br i1 %666, label %663, label %667

667:
	%668 = mul i64 2, %661
	%669 = add i64 %668, 1
	%670 = getelementptr [9 x i32], [9 x i32]* @b8, i64 0, i64 %669
	%671 = load i32, i32* %670
	%672 = add i32 %653, %671
	%673 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %673
	%674 = call i32 (i8*, ...) @printf([3 x i8]* %673, i32 %672)
	%675 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%676 = load i64, i64* @.global_b8_0
	%677 = icmp ult i64 0, 0
	br i1 %677, label %678, label %680

678:
	store i32 2, i32* @.ErrorNumber
	%679 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

680:
	%681 = icmp uge i64 0, %676
	br i1 %681, label %678, label %682

682:
	%683 = load i64, i64* @.global_b8_1
	%684 = icmp ult i64 0, 0
	br i1 %684, label %685, label %687

685:
	store i32 2, i32* @.ErrorNumber
	%686 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

687:
	%688 = icmp uge i64 0, %683
	br i1 %688, label %685, label %689

689:
	%690 = mul i64 0, %683
	%691 = add i64 %690, 0
	%692 = getelementptr [9 x i32], [9 x i32]* @b8, i64 0, i64 %691
	%693 = load i32, i32* %692
	%694 = load i64, i64* @.global_b8_0
	%695 = icmp ult i64 1, 0
	br i1 %695, label %696, label %698

696:
	store i32 2, i32* @.ErrorNumber
	%697 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

698:
	%699 = icmp uge i64 1, %694
	br i1 %699, label %696, label %700

700:
	%701 = load i64, i64* @.global_b8_1
	%702 = icmp ult i64 1, 0
	br i1 %702, label %703, label %705

703:
	store i32 2, i32* @.ErrorNumber
	%704 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

705:
	%706 = icmp uge i64 1, %701
	br i1 %706, label %703, label %707

707:
	%708 = mul i64 1, %701
	%709 = add i64 %708, 1
	%710 = getelementptr [9 x i32], [9 x i32]* @b8, i64 0, i64 %709
	%711 = load i32, i32* %710
	%712 = sub i32 %693, %711
	%713 = load i64, i64* @.global_b8_0
	%714 = icmp ult i64 2, 0
	br i1 %714, label %715, label %717

715:
	store i32 2, i32* @.ErrorNumber
	%716 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

717:
	%718 = icmp uge i64 2, %713
	br i1 %718, label %715, label %719

719:
	%720 = load i64, i64* @.global_b8_1
	%721 = icmp ult i64 1, 0
	br i1 %721, label %722, label %724

722:
	store i32 2, i32* @.ErrorNumber
	%723 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

724:
	%725 = icmp uge i64 1, %720
	br i1 %725, label %722, label %726

726:
	%727 = mul i64 2, %720
	%728 = add i64 %727, 1
	%729 = getelementptr [9 x i32], [9 x i32]* @b8, i64 0, i64 %728
	%730 = load i32, i32* %729
	%731 = load i64, i64* @.global_b8_0
	%732 = icmp ult i64 1, 0
	br i1 %732, label %733, label %735

733:
	store i32 2, i32* @.ErrorNumber
	%734 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

735:
	%736 = icmp uge i64 1, %731
	br i1 %736, label %733, label %737

737:
	%738 = load i64, i64* @.global_b8_1
	%739 = icmp ult i64 2, 0
	br i1 %739, label %740, label %742

740:
	store i32 2, i32* @.ErrorNumber
	%741 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

742:
	%743 = icmp uge i64 2, %738
	br i1 %743, label %740, label %744

744:
	%745 = mul i64 1, %738
	%746 = add i64 %745, 2
	%747 = getelementptr [9 x i32], [9 x i32]* @b8, i64 0, i64 %746
	%748 = load i32, i32* %747
	%749 = sub i32 0, %748
	%750 = mul i32 %730, %749
	%751 = add i32 %712, %750
	%752 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %752
	%753 = call i32 (i8*, ...) @printf([3 x i8]* %752, i32 %751)
	%754 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%755 = load i64, i64* @.global_b8_0
	%756 = icmp ult i64 1, 0
	br i1 %756, label %757, label %759

757:
	store i32 2, i32* @.ErrorNumber
	%758 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

759:
	%760 = icmp uge i64 1, %755
	br i1 %760, label %757, label %761

761:
	%762 = load i64, i64* @.global_b8_1
	%763 = icmp ult i64 2, 0
	br i1 %763, label %764, label %766

764:
	store i32 2, i32* @.ErrorNumber
	%765 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

766:
	%767 = icmp uge i64 2, %762
	br i1 %767, label %764, label %768

768:
	%769 = mul i64 1, %762
	%770 = add i64 %769, 2
	%771 = getelementptr [9 x i32], [9 x i32]* @b8, i64 0, i64 %770
	%772 = load i32, i32* %771
	%773 = load i64, i64* @.global_b8_0
	%774 = icmp ult i64 1, 0
	br i1 %774, label %775, label %777

775:
	store i32 2, i32* @.ErrorNumber
	%776 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

777:
	%778 = icmp uge i64 1, %773
	br i1 %778, label %775, label %779

779:
	%780 = load i64, i64* @.global_b8_1
	%781 = icmp ult i64 1, 0
	br i1 %781, label %782, label %784

782:
	store i32 2, i32* @.ErrorNumber
	%783 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

784:
	%785 = icmp uge i64 1, %780
	br i1 %785, label %782, label %786

786:
	%787 = mul i64 1, %780
	%788 = add i64 %787, 1
	%789 = getelementptr [9 x i32], [9 x i32]* @b8, i64 0, i64 %788
	%790 = load i32, i32* %789
	%791 = load i64, i64* @.global_b8_0
	%792 = icmp ult i64 2, 0
	br i1 %792, label %793, label %795

793:
	store i32 2, i32* @.ErrorNumber
	%794 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

795:
	%796 = icmp uge i64 2, %791
	br i1 %796, label %793, label %797

797:
	%798 = load i64, i64* @.global_b8_1
	%799 = icmp ult i64 1, 0
	br i1 %799, label %800, label %802

800:
	store i32 2, i32* @.ErrorNumber
	%801 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

802:
	%803 = icmp uge i64 1, %798
	br i1 %803, label %800, label %804

804:
	%805 = mul i64 2, %798
	%806 = add i64 %805, 1
	%807 = getelementptr [9 x i32], [9 x i32]* @b8, i64 0, i64 %806
	%808 = load i32, i32* %807
	%809 = load i64, i64* @.global_b8_0
	%810 = icmp ult i64 1, 0
	br i1 %810, label %811, label %813

811:
	store i32 2, i32* @.ErrorNumber
	%812 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

813:
	%814 = icmp uge i64 1, %809
	br i1 %814, label %811, label %815

815:
	%816 = load i64, i64* @.global_b8_1
	%817 = icmp ult i64 2, 0
	br i1 %817, label %818, label %820

818:
	store i32 2, i32* @.ErrorNumber
	%819 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

820:
	%821 = icmp uge i64 2, %816
	br i1 %821, label %818, label %822

822:
	%823 = mul i64 1, %816
	%824 = add i64 %823, 2
	%825 = getelementptr [9 x i32], [9 x i32]* @b8, i64 0, i64 %824
	%826 = load i32, i32* %825
	%827 = mul i32 %808, %826
	%828 = add i32 %790, %827
	%829 = sub i32 %772, %828
	%830 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %830
	%831 = call i32 (i8*, ...) @printf([3 x i8]* %830, i32 %829)
	%832 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%833 = load i64, i64* @.global_b8_0
	%834 = icmp ult i64 1, 0
	br i1 %834, label %835, label %837

835:
	store i32 2, i32* @.ErrorNumber
	%836 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

837:
	%838 = icmp uge i64 1, %833
	br i1 %838, label %835, label %839

839:
	%840 = load i64, i64* @.global_b8_1
	%841 = icmp ult i64 2, 0
	br i1 %841, label %842, label %844

842:
	store i32 2, i32* @.ErrorNumber
	%843 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

844:
	%845 = icmp uge i64 2, %840
	br i1 %845, label %842, label %846

846:
	%847 = mul i64 1, %840
	%848 = add i64 %847, 2
	%849 = getelementptr [9 x i32], [9 x i32]* @b8, i64 0, i64 %848
	%850 = load i32, i32* %849
	%851 = sext i32 %850 to i64
	%852 = mul i64 8, %851
	%853 = load i64, i64* @.global_b8_0
	%854 = icmp ult i64 1, 0
	br i1 %854, label %855, label %857

855:
	store i32 2, i32* @.ErrorNumber
	%856 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

857:
	%858 = icmp uge i64 1, %853
	br i1 %858, label %855, label %859

859:
	%860 = load i64, i64* @.global_b8_1
	%861 = icmp ult i64 1, 0
	br i1 %861, label %862, label %864

862:
	store i32 2, i32* @.ErrorNumber
	%863 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

864:
	%865 = icmp uge i64 1, %860
	br i1 %865, label %862, label %866

866:
	%867 = mul i64 1, %860
	%868 = add i64 %867, 1
	%869 = getelementptr [9 x i32], [9 x i32]* @b8, i64 0, i64 %868
	%870 = load i32, i32* %869
	%871 = load i64, i64* @.global_b8_0
	%872 = icmp ult i64 2, 0
	br i1 %872, label %873, label %875

873:
	store i32 2, i32* @.ErrorNumber
	%874 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

875:
	%876 = icmp uge i64 2, %871
	br i1 %876, label %873, label %877

877:
	%878 = load i64, i64* @.global_b8_1
	%879 = icmp ult i64 1, 0
	br i1 %879, label %880, label %882

880:
	store i32 2, i32* @.ErrorNumber
	%881 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

882:
	%883 = icmp uge i64 1, %878
	br i1 %883, label %880, label %884

884:
	%885 = mul i64 2, %878
	%886 = add i64 %885, 1
	%887 = getelementptr [9 x i32], [9 x i32]* @b8, i64 0, i64 %886
	%888 = load i32, i32* %887
	%889 = load i64, i64* @.global_b8_0
	%890 = icmp ult i64 1, 0
	br i1 %890, label %891, label %893

891:
	store i32 2, i32* @.ErrorNumber
	%892 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

893:
	%894 = icmp uge i64 1, %889
	br i1 %894, label %891, label %895

895:
	%896 = load i64, i64* @.global_b8_1
	%897 = icmp ult i64 2, 0
	br i1 %897, label %898, label %900

898:
	store i32 2, i32* @.ErrorNumber
	%899 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

900:
	%901 = icmp uge i64 2, %896
	br i1 %901, label %898, label %902

902:
	%903 = mul i64 1, %896
	%904 = add i64 %903, 2
	%905 = getelementptr [9 x i32], [9 x i32]* @b8, i64 0, i64 %904
	%906 = load i32, i32* %905
	%907 = mul i32 %888, %906
	%908 = add i32 %870, %907
	%909 = sext i32 %908 to i64
	%910 = sub i64 %852, %909
	%911 = alloca [4 x i8]
	store [4 x i8] c"%ld\00", [4 x i8]* %911
	%912 = call i32 (i8*, ...) @printf([4 x i8]* %911, i64 %910)
	%913 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%914 = load i64, i64* @.global_a_0
	%915 = icmp ult i64 1, 0
	br i1 %915, label %916, label %918

916:
	store i32 2, i32* @.ErrorNumber
	%917 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

918:
	%919 = icmp uge i64 1, %914
	br i1 %919, label %916, label %920

920:
	%921 = getelementptr [2 x i64], [2 x i64]* @a, i64 0, i64 1
	store i64 2, i64* %921
	%922 = load i64, i64* @.global_a_0
	%923 = icmp ult i64 0, 0
	br i1 %923, label %924, label %926

924:
	store i32 2, i32* @.ErrorNumber
	%925 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

926:
	%927 = icmp uge i64 0, %922
	br i1 %927, label %924, label %928

928:
	%929 = getelementptr [2 x i64], [2 x i64]* @a, i64 0, i64 0
	store i64 23, i64* %929
	%930 = load i64, i64* @.global_a_0
	%931 = icmp ult i64 1, 0
	br i1 %931, label %932, label %934

932:
	store i32 2, i32* @.ErrorNumber
	%933 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

934:
	%935 = icmp uge i64 1, %930
	br i1 %935, label %932, label %936

936:
	%937 = getelementptr [2 x i64], [2 x i64]* @a, i64 0, i64 1
	%938 = load i64, i64* %937
	%939 = alloca [4 x i8]
	store [4 x i8] c"%ld\00", [4 x i8]* %939
	%940 = call i32 (i8*, ...) @printf([4 x i8]* %939, i64 %938)
	%941 = alloca [3 x i8]
	store [3 x i8] c", \00", [3 x i8]* %941
	%942 = call i32 (i8*, ...) @printf([3 x i8]* %941)
	%943 = load i64, i64* @.global_a_0
	%944 = icmp ult i64 0, 0
	br i1 %944, label %945, label %947

945:
	store i32 2, i32* @.ErrorNumber
	%946 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

947:
	%948 = icmp uge i64 0, %943
	br i1 %948, label %945, label %949

949:
	%950 = getelementptr [2 x i64], [2 x i64]* @a, i64 0, i64 0
	%951 = load i64, i64* %950
	%952 = alloca [4 x i8]
	store [4 x i8] c"%ld\00", [4 x i8]* %952
	%953 = call i32 (i8*, ...) @printf([4 x i8]* %952, i64 %951)
	%954 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	store float 1.0, float* @curr
	%955 = fptosi double 0x4025FAE147AE147B to i64
	store i64 %955, i64* @lon
	%956 = fptrunc double 0x3FF199999999999A to float
	store float %956, float* @sng
	store double 0x3FF4CCCCCCCCCCCD, double* @dbl
	%957 = call i32 @strlen([6 x i8]* @_8)
	%958 = call i8* @gc_malloc(%struct.GarbageCollector* @gc, i32 %957)
	store i8* %958, i8** @str
	%959 = call i8* @strcpy(i8* %958, [6 x i8]* @_8)
	%960 = trunc i64 1 to i32
	store i32 %960, i32* @int
	store i1 false, i1* @bool
	%961 = load float, float* @curr
	%962 = alloca [3 x i8]
	store [3 x i8] c"%f\00", [3 x i8]* %962
	%963 = call i32 (i8*, ...) @printf([3 x i8]* %962, float %961)
	%964 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%965 = load i64, i64* @lon
	%966 = alloca [4 x i8]
	store [4 x i8] c"%ld\00", [4 x i8]* %966
	%967 = call i32 (i8*, ...) @printf([4 x i8]* %966, i64 %965)
	%968 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%969 = load float, float* @sng
	%970 = alloca [3 x i8]
	store [3 x i8] c"%f\00", [3 x i8]* %970
	%971 = call i32 (i8*, ...) @printf([3 x i8]* %970, float %969)
	%972 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%973 = load double, double* @dbl
	%974 = alloca [4 x i8]
	store [4 x i8] c"%lf\00", [4 x i8]* %974
	%975 = call i32 (i8*, ...) @printf([4 x i8]* %974, double %973)
	%976 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%977 = load i8*, i8** getelementptr (i8*, i8** @str, i32 0)
	%978 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %978
	%979 = call i32 (i8*, ...) @printf([3 x i8]* %978, i8* %977)
	%980 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%981 = load i32, i32* @int
	%982 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %982
	%983 = call i32 (i8*, ...) @printf([3 x i8]* %982, i32 %981)
	%984 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%985 = load i1, i1* @bool
	%986 = icmp eq i1 %985, true
	%987 = select i1 %986, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%988 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %988
	%989 = call i32 (i8*, ...) @printf([3 x i8]* %988, [5 x i8]* %987)
	%990 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%991 = load i64, i64* @.global_curra_0
	%992 = icmp ult i64 1, 0
	br i1 %992, label %993, label %995

993:
	store i32 2, i32* @.ErrorNumber
	%994 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

995:
	%996 = icmp uge i64 1, %991
	br i1 %996, label %993, label %997

997:
	%998 = getelementptr [2 x float], [2 x float]* @curra, i64 0, i64 1
	store float 0x3FF0CCCCC0000000, float* %998
	%999 = load i64, i64* @.global_lona_0
	%1000 = icmp ult i64 1, 0
	br i1 %1000, label %1001, label %1003

1001:
	store i32 2, i32* @.ErrorNumber
	%1002 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

1003:
	%1004 = icmp uge i64 1, %999
	br i1 %1004, label %1001, label %1005

1005:
	%1006 = getelementptr [2 x i64], [2 x i64]* @lona, i64 0, i64 1
	store i64 123456789, i64* %1006
	%1007 = load i64, i64* @.global_snga_0
	%1008 = icmp ult i64 1, 0
	br i1 %1008, label %1009, label %1011

1009:
	store i32 2, i32* @.ErrorNumber
	%1010 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

1011:
	%1012 = icmp uge i64 1, %1007
	br i1 %1012, label %1009, label %1013

1013:
	%1014 = getelementptr [2 x float], [2 x float]* @snga, i64 0, i64 1
	%1015 = fptrunc double 0x3FF199999999999A to float
	store float %1015, float* %1014
	%1016 = load i64, i64* @.global_dbla_0
	%1017 = icmp ult i64 1, 0
	br i1 %1017, label %1018, label %1020

1018:
	store i32 2, i32* @.ErrorNumber
	%1019 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

1020:
	%1021 = icmp uge i64 1, %1016
	br i1 %1021, label %1018, label %1022

1022:
	%1023 = getelementptr [2 x double], [2 x double]* @dbla, i64 0, i64 1
	store double 0x3FF4CCCCCCCCCCCD, double* %1023
	%1024 = load i64, i64* @.global_stra_0
	%1025 = icmp ult i64 1, 0
	br i1 %1025, label %1026, label %1028

1026:
	store i32 2, i32* @.ErrorNumber
	%1027 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

1028:
	%1029 = icmp uge i64 1, %1024
	br i1 %1029, label %1026, label %1030

1030:
	%1031 = getelementptr [2 x i8*], [2 x i8*]* @stra, i64 0, i64 1
	%1032 = call i32 @strlen([6 x i8]* @_9)
	%1033 = call i8* @gc_malloc(%struct.GarbageCollector* @gc, i32 %1032)
	store i8* %1033, i8** %1031
	%1034 = call i8* @strcpy(i8* %1033, [6 x i8]* @_9)
	%1035 = load i64, i64* @.global_inta_0
	%1036 = icmp ult i64 1, 0
	br i1 %1036, label %1037, label %1039

1037:
	store i32 2, i32* @.ErrorNumber
	%1038 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

1039:
	%1040 = icmp uge i64 1, %1035
	br i1 %1040, label %1037, label %1041

1041:
	%1042 = getelementptr [2 x i32], [2 x i32]* @inta, i64 0, i64 1
	%1043 = trunc i64 1 to i32
	store i32 %1043, i32* %1042
	%1044 = load i64, i64* @.global_boola_0
	%1045 = icmp ult i64 1, 0
	br i1 %1045, label %1046, label %1048

1046:
	store i32 2, i32* @.ErrorNumber
	%1047 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

1048:
	%1049 = icmp uge i64 1, %1044
	br i1 %1049, label %1046, label %1050

1050:
	%1051 = getelementptr [2 x i1], [2 x i1]* @boola, i64 0, i64 1
	store i1 false, i1* %1051
	%1052 = load i64, i64* @.global_curra_0
	%1053 = icmp ult i64 1, 0
	br i1 %1053, label %1054, label %1056

1054:
	store i32 2, i32* @.ErrorNumber
	%1055 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

1056:
	%1057 = icmp uge i64 1, %1052
	br i1 %1057, label %1054, label %1058

1058:
	%1059 = getelementptr [2 x float], [2 x float]* @curra, i64 0, i64 1
	%1060 = load float, float* %1059
	%1061 = alloca [3 x i8]
	store [3 x i8] c"%f\00", [3 x i8]* %1061
	%1062 = call i32 (i8*, ...) @printf([3 x i8]* %1061, float %1060)
	%1063 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%1064 = load i64, i64* @.global_lona_0
	%1065 = icmp ult i64 1, 0
	br i1 %1065, label %1066, label %1068

1066:
	store i32 2, i32* @.ErrorNumber
	%1067 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

1068:
	%1069 = icmp uge i64 1, %1064
	br i1 %1069, label %1066, label %1070

1070:
	%1071 = getelementptr [2 x i64], [2 x i64]* @lona, i64 0, i64 1
	%1072 = load i64, i64* %1071
	%1073 = alloca [4 x i8]
	store [4 x i8] c"%ld\00", [4 x i8]* %1073
	%1074 = call i32 (i8*, ...) @printf([4 x i8]* %1073, i64 %1072)
	%1075 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%1076 = load i64, i64* @.global_snga_0
	%1077 = icmp ult i64 1, 0
	br i1 %1077, label %1078, label %1080

1078:
	store i32 2, i32* @.ErrorNumber
	%1079 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

1080:
	%1081 = icmp uge i64 1, %1076
	br i1 %1081, label %1078, label %1082

1082:
	%1083 = getelementptr [2 x float], [2 x float]* @snga, i64 0, i64 1
	%1084 = load float, float* %1083
	%1085 = alloca [3 x i8]
	store [3 x i8] c"%f\00", [3 x i8]* %1085
	%1086 = call i32 (i8*, ...) @printf([3 x i8]* %1085, float %1084)
	%1087 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%1088 = load i64, i64* @.global_dbla_0
	%1089 = icmp ult i64 1, 0
	br i1 %1089, label %1090, label %1092

1090:
	store i32 2, i32* @.ErrorNumber
	%1091 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

1092:
	%1093 = icmp uge i64 1, %1088
	br i1 %1093, label %1090, label %1094

1094:
	%1095 = getelementptr [2 x double], [2 x double]* @dbla, i64 0, i64 1
	%1096 = load double, double* %1095
	%1097 = alloca [4 x i8]
	store [4 x i8] c"%lf\00", [4 x i8]* %1097
	%1098 = call i32 (i8*, ...) @printf([4 x i8]* %1097, double %1096)
	%1099 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%1100 = load i64, i64* @.global_stra_0
	%1101 = icmp ult i64 1, 0
	br i1 %1101, label %1102, label %1104

1102:
	store i32 2, i32* @.ErrorNumber
	%1103 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

1104:
	%1105 = icmp uge i64 1, %1100
	br i1 %1105, label %1102, label %1106

1106:
	%1107 = getelementptr [2 x i8*], [2 x i8*]* @stra, i64 0, i64 1
	%1108 = load i8*, i8** %1107
	%1109 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %1109
	%1110 = call i32 (i8*, ...) @printf([3 x i8]* %1109, i8* %1108)
	%1111 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%1112 = load i64, i64* @.global_inta_0
	%1113 = icmp ult i64 1, 0
	br i1 %1113, label %1114, label %1116

1114:
	store i32 2, i32* @.ErrorNumber
	%1115 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

1116:
	%1117 = icmp uge i64 1, %1112
	br i1 %1117, label %1114, label %1118

1118:
	%1119 = getelementptr [2 x i32], [2 x i32]* @inta, i64 0, i64 1
	%1120 = load i32, i32* %1119
	%1121 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %1121
	%1122 = call i32 (i8*, ...) @printf([3 x i8]* %1121, i32 %1120)
	%1123 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%1124 = load i64, i64* @.global_boola_0
	%1125 = icmp ult i64 1, 0
	br i1 %1125, label %1126, label %1128

1126:
	store i32 2, i32* @.ErrorNumber
	%1127 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

1128:
	%1129 = icmp uge i64 1, %1124
	br i1 %1129, label %1126, label %1130

1130:
	%1131 = getelementptr [2 x i1], [2 x i1]* @boola, i64 0, i64 1
	%1132 = load i1, i1* %1131
	%1133 = icmp eq i1 %1132, true
	%1134 = select i1 %1133, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%1135 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %1135
	%1136 = call i32 (i8*, ...) @printf([3 x i8]* %1135, [5 x i8]* %1134)
	%1137 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%1138 = load float, float* @currcg
	%1139 = alloca [3 x i8]
	store [3 x i8] c"%f\00", [3 x i8]* %1139
	%1140 = call i32 (i8*, ...) @printf([3 x i8]* %1139, float %1138)
	%1141 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%1142 = load i64, i64* @loncg
	%1143 = alloca [4 x i8]
	store [4 x i8] c"%ld\00", [4 x i8]* %1143
	%1144 = call i32 (i8*, ...) @printf([4 x i8]* %1143, i64 %1142)
	%1145 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%1146 = load float, float* @sngcg
	%1147 = alloca [3 x i8]
	store [3 x i8] c"%f\00", [3 x i8]* %1147
	%1148 = call i32 (i8*, ...) @printf([3 x i8]* %1147, float %1146)
	%1149 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%1150 = load double, double* @dblcg
	%1151 = alloca [4 x i8]
	store [4 x i8] c"%lf\00", [4 x i8]* %1151
	%1152 = call i32 (i8*, ...) @printf([4 x i8]* %1151, double %1150)
	%1153 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%1154 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %1154
	%1155 = call i32 (i8*, ...) @printf([3 x i8]* %1154, [12 x i8]* @strcg)
	%1156 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%1157 = load i32, i32* @intcg
	%1158 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %1158
	%1159 = call i32 (i8*, ...) @printf([3 x i8]* %1158, i32 %1157)
	%1160 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%1161 = load i1, i1* @boolcg
	%1162 = icmp eq i1 %1161, true
	%1163 = select i1 %1162, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%1164 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %1164
	%1165 = call i32 (i8*, ...) @printf([3 x i8]* %1164, [5 x i8]* %1163)
	%1166 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%1167 = load float, float* @datecg
	%1168 = alloca [3 x i8]
	store [3 x i8] c"%f\00", [3 x i8]* %1168
	%1169 = call i32 (i8*, ...) @printf([3 x i8]* %1168, float %1167)
	%1170 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	call void @setLocalConstants()
	%1171 = load i64, i64* @.global_arrayf_0
	%1172 = icmp ult i64 2, 0
	br i1 %1172, label %1173, label %1175

1173:
	store i32 2, i32* @.ErrorNumber
	%1174 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

1175:
	%1176 = icmp uge i64 2, %1171
	br i1 %1176, label %1173, label %1177

1177:
	%1178 = load i64, i64* @.global_arrayf_1
	%1179 = icmp ult i64 0, 0
	br i1 %1179, label %1180, label %1182

1180:
	store i32 2, i32* @.ErrorNumber
	%1181 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

1182:
	%1183 = icmp uge i64 0, %1178
	br i1 %1183, label %1180, label %1184

1184:
	%1185 = mul i64 2, %1178
	%1186 = add i64 %1185, 0
	%1187 = getelementptr [6 x i32], [6 x i32]* @arrayf, i64 0, i64 %1186
	%1188 = trunc i64 3 to i32
	store i32 %1188, i32* %1187
	call void @declareLocalArrays()
	%1189 = load i64, i64* @.global_arrayf_0
	%1190 = icmp ult i64 2, 0
	br i1 %1190, label %1191, label %1193

1191:
	store i32 2, i32* @.ErrorNumber
	%1192 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

1193:
	%1194 = icmp uge i64 2, %1189
	br i1 %1194, label %1191, label %1195

1195:
	%1196 = load i64, i64* @.global_arrayf_1
	%1197 = icmp ult i64 0, 0
	br i1 %1197, label %1198, label %1200

1198:
	store i32 2, i32* @.ErrorNumber
	%1199 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

1200:
	%1201 = icmp uge i64 0, %1196
	br i1 %1201, label %1198, label %1202

1202:
	%1203 = mul i64 2, %1196
	%1204 = add i64 %1203, 0
	%1205 = getelementptr [6 x i32], [6 x i32]* @arrayf, i64 0, i64 %1204
	%1206 = load i32, i32* %1205
	%1207 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %1207
	%1208 = call i32 (i8*, ...) @printf([3 x i8]* %1207, i32 %1206)
	%1209 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	ret void
}
