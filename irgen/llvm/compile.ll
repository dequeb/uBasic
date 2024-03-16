@vbEmpty = constant [1 x i8] c"\00"
@vbCR = constant [2 x i8] c"\0D\00"
@vbLF = constant [2 x i8] c"\0A\00"
@vbCrLf = constant [3 x i8] c"\0D\0A\00"
@vbTab = constant [2 x i8] c"\09\00"
@true = constant [5 x i8] c"True\00"
@false = constant [6 x i8] c"False\00"
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
@c5 = constant double 5.0
@c0 = constant i64 10
@s1 = global i8* null
@s2 = global i8* null
@s3 = global i8* null
@k0 = constant i64 10
@k1 = constant i64 20
@i11 = global i32 0
@_1 = constant [14 x i8] c"Allo le monde\00"
@_2 = constant [6 x i8] c"hello\00"
@_3 = constant [6 x i8] c"world\00"
@_4 = constant [2 x i8] c" \00"
@_5 = constant [3 x i8] c" !\00"
@_6 = constant [7 x i8] c"Hello \00"
@_7 = constant [8 x i8] c"world !\00"

declare i8* @malloc(i64 %size)

declare void @free(i8* %ptr)

declare i32 @printf(i8* %format, ...)

declare i32 @puts(i8* %s)

declare i32 @scanf(i8* %format, ...)

declare i8* @strcpy(i8* %dst, i8* %src)

declare i8* @strcat(i8* %dst, i8* %src)

declare i32 @sscanf(i8* %str, i8* %format, i8* %dst)

declare i32 @strlen(i8* %str)

declare i32 @sprintf(i8* %str, i8* %format, ...)

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
	%4 = icmp sgt i32 %3, 10
	br i1 %4, label %5, label %12

5:
	%6 = load i32*, i32* %i
	%7 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %7
	%8 = call i32 (i8*, ...) @printf([3 x i8]* %7, i32* %6)
	%9 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%10 = load i32, i32* %i
	%11 = sub i32 %10, 1
	store i32 %11, i32* %i
	br label %2

12:
	ret void
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

define i32 @main(i32 %argc, i8** %argv) {
0:
	store i64 10, i64* @l
	store double 0x400921FAFC8B007A, double* @d
	store float 0x42012E0BE0000000, float* @da
	%1 = load i8*, i8** @s0
	%2 = call i32 @strlen([14 x i8]* @_1)
	%3 = call i8* @malloc(i32 %2)
	store i8* %3, i8** @s0
	%4 = call i8* @strcpy(i8* %3, [14 x i8]* @_1)
	store float 100.0, float* @c
	store i1 true, i1* @b
	%5 = alloca [11 x i8]
	store [11 x i8] c"variables:\00", [11 x i8]* %5
	%6 = call i32 (i8*, ...) @printf([11 x i8]* %5)
	%7 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%8 = load i64, i64* @l
	%9 = alloca [4 x i8]
	store [4 x i8] c"%ld\00", [4 x i8]* %9
	%10 = call i32 (i8*, ...) @printf([4 x i8]* %9, i64 %8)
	%11 = alloca [3 x i8]
	store [3 x i8] c", \00", [3 x i8]* %11
	%12 = call i32 (i8*, ...) @printf([3 x i8]* %11)
	%13 = load double, double* @d
	%14 = alloca [4 x i8]
	store [4 x i8] c"%lf\00", [4 x i8]* %14
	%15 = call i32 (i8*, ...) @printf([4 x i8]* %14, double %13)
	%16 = alloca [3 x i8]
	store [3 x i8] c", \00", [3 x i8]* %16
	%17 = call i32 (i8*, ...) @printf([3 x i8]* %16)
	%18 = load i1, i1* @b
	%19 = icmp eq i1 %18, true
	%20 = select i1 %19, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%21 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %21
	%22 = call i32 (i8*, ...) @printf([3 x i8]* %21, [5 x i8]* %20)
	%23 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%24 = load float, float* @da
	%25 = alloca [3 x i8]
	store [3 x i8] c"%f\00", [3 x i8]* %25
	%26 = call i32 (i8*, ...) @printf([3 x i8]* %25, float %24)
	%27 = alloca [3 x i8]
	store [3 x i8] c", \00", [3 x i8]* %27
	%28 = call i32 (i8*, ...) @printf([3 x i8]* %27)
	%29 = load i8*, i8** getelementptr (i8*, i8** @s0, i32 0)
	%30 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %30
	%31 = call i32 (i8*, ...) @printf([3 x i8]* %30, i8* %29)
	%32 = alloca [3 x i8]
	store [3 x i8] c", \00", [3 x i8]* %32
	%33 = call i32 (i8*, ...) @printf([3 x i8]* %32)
	%34 = load float, float* @c
	%35 = alloca [3 x i8]
	store [3 x i8] c"%f\00", [3 x i8]* %35
	%36 = call i32 (i8*, ...) @printf([3 x i8]* %35, float %34)
	%37 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%38 = alloca [11 x i8]
	store [11 x i8] c"constants:\00", [11 x i8]* %38
	%39 = call i32 (i8*, ...) @printf([11 x i8]* %38)
	%40 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%41 = alloca [4 x i8]
	store [4 x i8] c"c1:\00", [4 x i8]* %41
	%42 = call i32 (i8*, ...) @printf([4 x i8]* %41)
	%43 = load i64, i64* @c1
	%44 = alloca [4 x i8]
	store [4 x i8] c"%ld\00", [4 x i8]* %44
	%45 = call i32 (i8*, ...) @printf([4 x i8]* %44, i64 %43)
	%46 = alloca [6 x i8]
	store [6 x i8] c", c2:\00", [6 x i8]* %46
	%47 = call i32 (i8*, ...) @printf([6 x i8]* %46)
	%48 = load double, double* @c2
	%49 = alloca [4 x i8]
	store [4 x i8] c"%lf\00", [4 x i8]* %49
	%50 = call i32 (i8*, ...) @printf([4 x i8]* %49, double %48)
	%51 = alloca [6 x i8]
	store [6 x i8] c", c3:\00", [6 x i8]* %51
	%52 = call i32 (i8*, ...) @printf([6 x i8]* %51)
	%53 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %53
	%54 = call i32 (i8*, ...) @printf([3 x i8]* %53, [6 x i8]* @c3)
	%55 = alloca [6 x i8]
	store [6 x i8] c", c4:\00", [6 x i8]* %55
	%56 = call i32 (i8*, ...) @printf([6 x i8]* %55)
	%57 = load i1, i1* @c4
	%58 = icmp eq i1 %57, true
	%59 = select i1 %58, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%60 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %60
	%61 = call i32 (i8*, ...) @printf([3 x i8]* %60, [5 x i8]* %59)
	%62 = alloca [6 x i8]
	store [6 x i8] c", c5:\00", [6 x i8]* %62
	%63 = call i32 (i8*, ...) @printf([6 x i8]* %62)
	%64 = load double, double* @c5
	%65 = alloca [3 x i8]
	store [3 x i8] c"%f\00", [3 x i8]* %65
	%66 = call i32 (i8*, ...) @printf([3 x i8]* %65, double %64)
	%67 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%68 = alloca [10 x i8]
	store [10 x i8] c"literals:\00", [10 x i8]* %68
	%69 = call i32 (i8*, ...) @printf([10 x i8]* %68)
	%70 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%71 = alloca [6 x i8]
	store [6 x i8] c"world\00", [6 x i8]* %71
	%72 = call i32 (i8*, ...) @printf([6 x i8]* %71)
	%73 = alloca [5 x i8]
	store [5 x i8] c"True\00", [5 x i8]* %73
	%74 = call i32 (i8*, ...) @printf([5 x i8]* %73)
	%75 = alloca [8 x i8]
	store [8 x i8] c"1.23456\00", [8 x i8]* %75
	%76 = call i32 (i8*, ...) @printf([8 x i8]* %75)
	%77 = alloca [3 x i8]
	store [3 x i8] c"10\00", [3 x i8]* %77
	%78 = call i32 (i8*, ...) @printf([3 x i8]* %77)
	%79 = alloca [11 x i8]
	store [11 x i8] c"2010/12/31\00", [11 x i8]* %79
	%80 = call i32 (i8*, ...) @printf([11 x i8]* %79)
	%81 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%82 = trunc i64 10 to i32
	store i32 %82, i32* @i
	%83 = fptrunc double 0x405900A3D70A3D71 to float
	store float %83, float* @s
	%84 = load i32, i32* @i
	%85 = sext i32 %84 to i64
	store i64 %85, i64* @l
	%86 = load float, float* @s
	%87 = fpext float %86 to double
	store double %87, double* @d
	%88 = load i32, i32* @i
	%89 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %89
	%90 = call i32 (i8*, ...) @printf([3 x i8]* %89, i32 %88)
	%91 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%92 = load float, float* @s
	%93 = alloca [3 x i8]
	store [3 x i8] c"%f\00", [3 x i8]* %93
	%94 = call i32 (i8*, ...) @printf([3 x i8]* %93, float %92)
	%95 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%96 = load i64, i64* @l
	%97 = alloca [4 x i8]
	store [4 x i8] c"%ld\00", [4 x i8]* %97
	%98 = call i32 (i8*, ...) @printf([4 x i8]* %97, i64 %96)
	%99 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%100 = load double, double* @d
	%101 = alloca [4 x i8]
	store [4 x i8] c"%lf\00", [4 x i8]* %101
	%102 = call i32 (i8*, ...) @printf([4 x i8]* %101, double %100)
	%103 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%104 = trunc i64 1234 to i32
	store i32 %104, i32* @i
	store i64 1234567890, i64* @l
	%105 = load i32, i32* @i
	%106 = add i32 %105, 1
	%107 = sitofp i32 %106 to float
	store float %107, float* @s
	%108 = load i64, i64* @l
	%109 = add i64 %108, 1
	%110 = sitofp i64 %109 to double
	store double %110, double* @d
	%111 = load i32, i32* @i
	%112 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %112
	%113 = call i32 (i8*, ...) @printf([3 x i8]* %112, i32 %111)
	%114 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%115 = load float, float* @s
	%116 = alloca [3 x i8]
	store [3 x i8] c"%f\00", [3 x i8]* %116
	%117 = call i32 (i8*, ...) @printf([3 x i8]* %116, float %115)
	%118 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%119 = load i64, i64* @l
	%120 = alloca [4 x i8]
	store [4 x i8] c"%ld\00", [4 x i8]* %120
	%121 = call i32 (i8*, ...) @printf([4 x i8]* %120, i64 %119)
	%122 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%123 = load double, double* @d
	%124 = alloca [4 x i8]
	store [4 x i8] c"%lf\00", [4 x i8]* %124
	%125 = call i32 (i8*, ...) @printf([4 x i8]* %124, double %123)
	%126 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%127 = fptrunc double 0x4028AE147AE147AE to float
	store float %127, float* @s
	store double 0x41678C29DCCCCCCD, double* @d
	%128 = load float, float* @s
	%129 = fptosi float %128 to i32
	store i32 %129, i32* @i
	%130 = load double, double* @d
	%131 = fptosi double %130 to i64
	store i64 %131, i64* @l
	%132 = load i32, i32* @i
	%133 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %133
	%134 = call i32 (i8*, ...) @printf([3 x i8]* %133, i32 %132)
	%135 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%136 = load float, float* @s
	%137 = alloca [3 x i8]
	store [3 x i8] c"%f\00", [3 x i8]* %137
	%138 = call i32 (i8*, ...) @printf([3 x i8]* %137, float %136)
	%139 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%140 = load i64, i64* @l
	%141 = alloca [4 x i8]
	store [4 x i8] c"%ld\00", [4 x i8]* %141
	%142 = call i32 (i8*, ...) @printf([4 x i8]* %141, i64 %140)
	%143 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%144 = load double, double* @d
	%145 = alloca [4 x i8]
	store [4 x i8] c"%lf\00", [4 x i8]* %145
	%146 = call i32 (i8*, ...) @printf([4 x i8]* %145, double %144)
	%147 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	store double 0x41678C29DCCCCCCD, double* @d
	%148 = load double, double* @d
	%149 = fptrunc double %148 to float
	store float %149, float* @s
	%150 = load float, float* @s
	%151 = fptosi float %150 to i64
	store i64 %151, i64* @l
	%152 = load i64, i64* @l
	%153 = trunc i64 %152 to i32
	store i32 %153, i32* @i
	%154 = load i32, i32* @i
	%155 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %155
	%156 = call i32 (i8*, ...) @printf([3 x i8]* %155, i32 %154)
	%157 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%158 = load float, float* @s
	%159 = alloca [3 x i8]
	store [3 x i8] c"%f\00", [3 x i8]* %159
	%160 = call i32 (i8*, ...) @printf([3 x i8]* %159, float %158)
	%161 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%162 = load i64, i64* @l
	%163 = alloca [4 x i8]
	store [4 x i8] c"%ld\00", [4 x i8]* %163
	%164 = call i32 (i8*, ...) @printf([4 x i8]* %163, i64 %162)
	%165 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%166 = load double, double* @d
	%167 = alloca [4 x i8]
	store [4 x i8] c"%lf\00", [4 x i8]* %167
	%168 = call i32 (i8*, ...) @printf([4 x i8]* %167, double %166)
	%169 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	call void @ifBranch()
	%170 = load i8*, i8** @s1
	%171 = call i32 @strlen([6 x i8]* @_2)
	%172 = call i8* @malloc(i32 %171)
	store i8* %172, i8** @s1
	%173 = call i8* @strcpy(i8* %172, [6 x i8]* @_2)
	%174 = load i8*, i8** @s2
	%175 = call i32 @strlen([6 x i8]* @_3)
	%176 = call i8* @malloc(i32 %175)
	store i8* %176, i8** @s2
	%177 = call i8* @strcpy(i8* %176, [6 x i8]* @_3)
	%178 = load i8*, i8** @s1
	%179 = call i32 @strlen(i8* %178)
	%180 = call i32 @strlen([2 x i8]* @_4)
	%181 = add i32 %179, %180
	%182 = call i8* @malloc(i32 %181)
	%183 = call i8* @strcpy(i8* %182, i8* %178)
	%184 = call i8* @strcat(i8* %182, [2 x i8]* @_4)
	%185 = load i8*, i8** @s2
	%186 = call i32 @strlen(i8* %182)
	%187 = call i32 @strlen(i8* %185)
	%188 = add i32 %186, %187
	%189 = call i8* @malloc(i32 %188)
	%190 = call i8* @strcpy(i8* %189, i8* %182)
	%191 = call i8* @strcat(i8* %189, i8* %185)
	%192 = call i32 @strlen(i8* %189)
	%193 = call i32 @strlen([3 x i8]* @_5)
	%194 = add i32 %192, %193
	%195 = call i8* @malloc(i32 %194)
	%196 = call i8* @strcpy(i8* %195, i8* %189)
	%197 = call i8* @strcat(i8* %195, [3 x i8]* @_5)
	%198 = load i8*, i8** @s3
	%199 = call i32 @strlen(i8* %195)
	%200 = call i8* @malloc(i32 %199)
	store i8* %200, i8** @s3
	%201 = call i8* @strcpy(i8* %200, i8* %195)
	%202 = load i8*, i8** getelementptr (i8*, i8** @s3, i32 0)
	%203 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %203
	%204 = call i32 (i8*, ...) @printf([3 x i8]* %203, i8* %202)
	%205 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%206 = call i32 @strlen([7 x i8]* @_6)
	%207 = call i32 @strlen([8 x i8]* @_7)
	%208 = add i32 %206, %207
	%209 = call i8* @malloc(i32 %208)
	%210 = call i8* @strcpy(i8* %209, [7 x i8]* @_6)
	%211 = call i8* @strcat(i8* %209, [8 x i8]* @_7)
	%212 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %212
	%213 = call i32 (i8*, ...) @printf([3 x i8]* %212, i8* %209)
	call void @free(i8* %209)
	%214 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%215 = add i64 1, 2
	%216 = alloca [4 x i8]
	store [4 x i8] c"%ld\00", [4 x i8]* %216
	%217 = call i32 (i8*, ...) @printf([4 x i8]* %216, i64 %215)
	%218 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%219 = fmul double 0x4025CCCCCCCCCCCD, 0x3FEF5C28F5C28F5C
	%220 = alloca [4 x i8]
	store [4 x i8] c"%lf\00", [4 x i8]* %220
	%221 = call i32 (i8*, ...) @printf([4 x i8]* %220, double %219)
	%222 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%223 = sitofp i64 1 to double
	%224 = fadd double %223, 0x4007D70A3D70A3D7
	%225 = alloca [4 x i8]
	store [4 x i8] c"%lf\00", [4 x i8]* %225
	%226 = call i32 (i8*, ...) @printf([4 x i8]* %225, double %224)
	%227 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	br i1 true, label %228, label %229

228:
	br label %229

229:
	%230 = phi i1 [ false, %0 ], [ false, %228 ]
	%231 = icmp eq i1 %230, true
	%232 = select i1 %231, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%233 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %233
	%234 = call i32 (i8*, ...) @printf([3 x i8]* %233, [5 x i8]* %232)
	%235 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	br i1 true, label %236, label %237

236:
	br label %237

237:
	%238 = phi i1 [ false, %229 ], [ true, %236 ]
	%239 = icmp eq i1 %238, true
	%240 = select i1 %239, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%241 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %241
	%242 = call i32 (i8*, ...) @printf([3 x i8]* %241, [5 x i8]* %240)
	%243 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	br i1 false, label %244, label %245

244:
	br label %245

245:
	%246 = phi i1 [ false, %237 ], [ false, %244 ]
	%247 = icmp eq i1 %246, true
	%248 = select i1 %247, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%249 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %249
	%250 = call i32 (i8*, ...) @printf([3 x i8]* %249, [5 x i8]* %248)
	%251 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	br i1 false, label %252, label %253

252:
	br label %253

253:
	%254 = phi i1 [ false, %245 ], [ true, %252 ]
	%255 = icmp eq i1 %254, true
	%256 = select i1 %255, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%257 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %257
	%258 = call i32 (i8*, ...) @printf([3 x i8]* %257, [5 x i8]* %256)
	%259 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	br i1 true, label %261, label %260

260:
	br label %261

261:
	%262 = phi i1 [ true, %253 ], [ false, %260 ]
	%263 = icmp eq i1 %262, true
	%264 = select i1 %263, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%265 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %265
	%266 = call i32 (i8*, ...) @printf([3 x i8]* %265, [5 x i8]* %264)
	%267 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	br i1 true, label %269, label %268

268:
	br label %269

269:
	%270 = phi i1 [ true, %261 ], [ true, %268 ]
	%271 = icmp eq i1 %270, true
	%272 = select i1 %271, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%273 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %273
	%274 = call i32 (i8*, ...) @printf([3 x i8]* %273, [5 x i8]* %272)
	%275 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	br i1 false, label %277, label %276

276:
	br label %277

277:
	%278 = phi i1 [ true, %269 ], [ false, %276 ]
	%279 = icmp eq i1 %278, true
	%280 = select i1 %279, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%281 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %281
	%282 = call i32 (i8*, ...) @printf([3 x i8]* %281, [5 x i8]* %280)
	%283 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	br i1 false, label %285, label %284

284:
	br label %285

285:
	%286 = phi i1 [ true, %277 ], [ true, %284 ]
	%287 = icmp eq i1 %286, true
	%288 = select i1 %287, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%289 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %289
	%290 = call i32 (i8*, ...) @printf([3 x i8]* %289, [5 x i8]* %288)
	%291 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%292 = xor i1 true, true
	%293 = icmp eq i1 %292, true
	%294 = select i1 %293, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%295 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %295
	%296 = call i32 (i8*, ...) @printf([3 x i8]* %295, [5 x i8]* %294)
	%297 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%298 = sitofp i64 2 to float
	%299 = fmul float %298, 2.25
	%300 = alloca [3 x i8]
	store [3 x i8] c"%f\00", [3 x i8]* %300
	%301 = call i32 (i8*, ...) @printf([3 x i8]* %300, float %299)
	%302 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%303 = sdiv i64 2, 2
	%304 = alloca [4 x i8]
	store [4 x i8] c"%ld\00", [4 x i8]* %304
	%305 = call i32 (i8*, ...) @printf([4 x i8]* %304, i64 %303)
	%306 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%307 = icmp slt i64 2, 3
	%308 = icmp eq i1 %307, true
	%309 = select i1 %308, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%310 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %310
	%311 = call i32 (i8*, ...) @printf([3 x i8]* %310, [5 x i8]* %309)
	%312 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%313 = icmp sgt i64 2, 3
	%314 = icmp eq i1 %313, true
	%315 = select i1 %314, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%316 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %316
	%317 = call i32 (i8*, ...) @printf([3 x i8]* %316, [5 x i8]* %315)
	%318 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%319 = icmp sle i64 2, 3
	%320 = icmp eq i1 %319, true
	%321 = select i1 %320, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%322 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %322
	%323 = call i32 (i8*, ...) @printf([3 x i8]* %322, [5 x i8]* %321)
	%324 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%325 = icmp sge i64 2, 3
	%326 = icmp eq i1 %325, true
	%327 = select i1 %326, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%328 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %328
	%329 = call i32 (i8*, ...) @printf([3 x i8]* %328, [5 x i8]* %327)
	%330 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%331 = icmp eq i64 2, 3
	%332 = icmp eq i1 %331, true
	%333 = select i1 %332, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%334 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %334
	%335 = call i32 (i8*, ...) @printf([3 x i8]* %334, [5 x i8]* %333)
	%336 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%337 = icmp ne i64 2, 3
	%338 = icmp eq i1 %337, true
	%339 = select i1 %338, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%340 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %340
	%341 = call i32 (i8*, ...) @printf([3 x i8]* %340, [5 x i8]* %339)
	%342 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%343 = fcmp ogt double 2.5, 3.0
	%344 = icmp eq i1 %343, true
	%345 = select i1 %344, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%346 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %346
	%347 = call i32 (i8*, ...) @printf([3 x i8]* %346, [5 x i8]* %345)
	%348 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%349 = fcmp olt double 2.5, 3.0
	%350 = icmp eq i1 %349, true
	%351 = select i1 %350, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%352 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %352
	%353 = call i32 (i8*, ...) @printf([3 x i8]* %352, [5 x i8]* %351)
	%354 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%355 = fcmp ole double 2.5, 3.0
	%356 = icmp eq i1 %355, true
	%357 = select i1 %356, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%358 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %358
	%359 = call i32 (i8*, ...) @printf([3 x i8]* %358, [5 x i8]* %357)
	%360 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%361 = fcmp oge double 2.5, 3.0
	%362 = icmp eq i1 %361, true
	%363 = select i1 %362, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%364 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %364
	%365 = call i32 (i8*, ...) @printf([3 x i8]* %364, [5 x i8]* %363)
	%366 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%367 = fcmp oeq double 2.5, 3.0
	%368 = icmp eq i1 %367, true
	%369 = select i1 %368, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%370 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %370
	%371 = call i32 (i8*, ...) @printf([3 x i8]* %370, [5 x i8]* %369)
	%372 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%373 = fcmp one double 2.5, 3.0
	%374 = icmp eq i1 %373, true
	%375 = select i1 %374, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
	%376 = alloca [3 x i8]
	store [3 x i8] c"%s\00", [3 x i8]* %376
	%377 = call i32 (i8*, ...) @printf([3 x i8]* %376, [5 x i8]* %375)
	%378 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%379 = trunc i64 10 to i32
	store i32 %379, i32* @i11
	br label %380

380:
	%381 = load i32, i32* @i11
	%382 = icmp sgt i32 %381, 0
	br i1 %382, label %383, label %390

383:
	%384 = load i32, i32* @i11
	%385 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %385
	%386 = call i32 (i8*, ...) @printf([3 x i8]* %385, i32 %384)
	%387 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%388 = load i32, i32* @i11
	%389 = sub i32 %388, 1
	store i32 %389, i32* @i11
	br label %380

390:
	call void @while1()
	%391 = call i64 @times2(i64 10)
	%392 = alloca [4 x i8]
	store [4 x i8] c"%ld\00", [4 x i8]* %392
	%393 = call i32 (i8*, ...) @printf([4 x i8]* %392, i64 %391)
	%394 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	ret i32 0
}
