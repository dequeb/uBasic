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
@.bzzz_1 = constant i64 2
@.bzzz_0 = constant i64 2
@bzzz = global [4 x i32] zeroinitializer

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
	%.a_1 = alloca i64
	store i64 2, i64* %.a_1
	%.a_0 = alloca i64
	store i64 3, i64* %.a_0
	%1 = alloca [6 x i32]
	%2 = alloca i64
	store i64 0, i64* %2
	br label %3

3:
	%4 = load i64, i64* %2
	%5 = icmp ult i64 %4, 6
	br i1 %5, label %6, label %10

6:
	%7 = getelementptr i32, [6 x i32]* %1, i64 %4
	store i32 0, i32* %7
	%8 = load i64, i64* %2
	%9 = add i64 %8, 1
	store i64 %9, i64* %2
	br label %3

10:
	%11 = load i64, i64* %.a_0
	%12 = icmp uge i64 0, %11
	br i1 %12, label %13, label %15

13:
	store i32 2, i32* @.ErrorNumber
	%14 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

15:
	%16 = icmp ult i64 0, 0
	br i1 %16, label %13, label %17

17:
	%18 = load i64, i64* %.a_1
	%19 = icmp uge i64 1, %18
	br i1 %19, label %20, label %22

20:
	store i32 2, i32* @.ErrorNumber
	%21 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

22:
	%23 = icmp ult i64 1, 0
	br i1 %23, label %20, label %24

24:
	%25 = mul i64 0, %18
	%26 = add i64 %25, 1
	%27 = getelementptr [6 x i32], [6 x i32]* %1, i64 0, i64 %26
	%28 = trunc i64 1 to i32
	store i32 %28, i32* %27
	%29 = load i64, i64* %.a_0
	%30 = icmp uge i64 1, %29
	br i1 %30, label %31, label %33

31:
	store i32 2, i32* @.ErrorNumber
	%32 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

33:
	%34 = icmp ult i64 1, 0
	br i1 %34, label %31, label %35

35:
	%36 = load i64, i64* %.a_1
	%37 = icmp uge i64 0, %36
	br i1 %37, label %38, label %40

38:
	store i32 2, i32* @.ErrorNumber
	%39 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

40:
	%41 = icmp ult i64 0, 0
	br i1 %41, label %38, label %42

42:
	%43 = mul i64 1, %36
	%44 = add i64 %43, 0
	%45 = getelementptr [6 x i32], [6 x i32]* %1, i64 0, i64 %44
	%46 = trunc i64 2 to i32
	store i32 %46, i32* %45
	%47 = load i64, i64* %.a_0
	%48 = icmp uge i64 1, %47
	br i1 %48, label %49, label %51

49:
	store i32 2, i32* @.ErrorNumber
	%50 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

51:
	%52 = icmp ult i64 1, 0
	br i1 %52, label %49, label %53

53:
	%54 = load i64, i64* %.a_1
	%55 = icmp uge i64 1, %54
	br i1 %55, label %56, label %58

56:
	store i32 2, i32* @.ErrorNumber
	%57 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

58:
	%59 = icmp ult i64 1, 0
	br i1 %59, label %56, label %60

60:
	%61 = mul i64 1, %54
	%62 = add i64 %61, 1
	%63 = getelementptr [6 x i32], [6 x i32]* %1, i64 0, i64 %62
	%64 = trunc i64 3 to i32
	store i32 %64, i32* %63
	%65 = load i64, i64* @.bzzz_0
	%66 = icmp uge i64 0, %65
	br i1 %66, label %67, label %69

67:
	store i32 2, i32* @.ErrorNumber
	%68 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

69:
	%70 = icmp ult i64 0, 0
	br i1 %70, label %67, label %71

71:
	%72 = load i64, i64* @.bzzz_1
	%73 = icmp uge i64 0, %72
	br i1 %73, label %74, label %76

74:
	store i32 2, i32* @.ErrorNumber
	%75 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

76:
	%77 = icmp ult i64 0, 0
	br i1 %77, label %74, label %78

78:
	%79 = mul i64 0, %72
	%80 = add i64 %79, 0
	%81 = getelementptr [4 x i32], [4 x i32]* @bzzz, i64 0, i64 %80
	%82 = trunc i64 4 to i32
	store i32 %82, i32* %81
	%83 = load i64, i64* %.a_0
	%84 = icmp uge i64 0, %83
	br i1 %84, label %85, label %87

85:
	store i32 2, i32* @.ErrorNumber
	%86 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

87:
	%88 = icmp ult i64 0, 0
	br i1 %88, label %85, label %89

89:
	%90 = load i64, i64* %.a_1
	%91 = icmp uge i64 0, %90
	br i1 %91, label %92, label %94

92:
	store i32 2, i32* @.ErrorNumber
	%93 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

94:
	%95 = icmp ult i64 0, 0
	br i1 %95, label %92, label %96

96:
	%97 = mul i64 0, %90
	%98 = add i64 %97, 0
	%99 = getelementptr [6 x i32], [6 x i32]* %1, i64 0, i64 %98
	%100 = load i32, i32* %99
	%101 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %101
	%102 = call i32 (i8*, ...) @printf([3 x i8]* %101, i32 %100)
	%103 = load i64, i64* %.a_0
	%104 = icmp uge i64 0, %103
	br i1 %104, label %105, label %107

105:
	store i32 2, i32* @.ErrorNumber
	%106 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

107:
	%108 = icmp ult i64 0, 0
	br i1 %108, label %105, label %109

109:
	%110 = load i64, i64* %.a_1
	%111 = icmp uge i64 1, %110
	br i1 %111, label %112, label %114

112:
	store i32 2, i32* @.ErrorNumber
	%113 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

114:
	%115 = icmp ult i64 1, 0
	br i1 %115, label %112, label %116

116:
	%117 = mul i64 0, %110
	%118 = add i64 %117, 1
	%119 = getelementptr [6 x i32], [6 x i32]* %1, i64 0, i64 %118
	%120 = load i32, i32* %119
	%121 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %121
	%122 = call i32 (i8*, ...) @printf([3 x i8]* %121, i32 %120)
	%123 = load i64, i64* %.a_0
	%124 = icmp uge i64 1, %123
	br i1 %124, label %125, label %127

125:
	store i32 2, i32* @.ErrorNumber
	%126 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

127:
	%128 = icmp ult i64 1, 0
	br i1 %128, label %125, label %129

129:
	%130 = load i64, i64* %.a_1
	%131 = icmp uge i64 0, %130
	br i1 %131, label %132, label %134

132:
	store i32 2, i32* @.ErrorNumber
	%133 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

134:
	%135 = icmp ult i64 0, 0
	br i1 %135, label %132, label %136

136:
	%137 = mul i64 1, %130
	%138 = add i64 %137, 0
	%139 = getelementptr [6 x i32], [6 x i32]* %1, i64 0, i64 %138
	%140 = load i32, i32* %139
	%141 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %141
	%142 = call i32 (i8*, ...) @printf([3 x i8]* %141, i32 %140)
	%143 = load i64, i64* %.a_0
	%144 = icmp uge i64 1, %143
	br i1 %144, label %145, label %147

145:
	store i32 2, i32* @.ErrorNumber
	%146 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

147:
	%148 = icmp ult i64 1, 0
	br i1 %148, label %145, label %149

149:
	%150 = load i64, i64* %.a_1
	%151 = icmp uge i64 1, %150
	br i1 %151, label %152, label %154

152:
	store i32 2, i32* @.ErrorNumber
	%153 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

154:
	%155 = icmp ult i64 1, 0
	br i1 %155, label %152, label %156

156:
	%157 = mul i64 1, %150
	%158 = add i64 %157, 1
	%159 = getelementptr [6 x i32], [6 x i32]* %1, i64 0, i64 %158
	%160 = load i32, i32* %159
	%161 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %161
	%162 = call i32 (i8*, ...) @printf([3 x i8]* %161, i32 %160)
	%163 = load i64, i64* @.bzzz_0
	%164 = icmp uge i64 0, %163
	br i1 %164, label %165, label %167

165:
	store i32 2, i32* @.ErrorNumber
	%166 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

167:
	%168 = icmp ult i64 0, 0
	br i1 %168, label %165, label %169

169:
	%170 = load i64, i64* @.bzzz_1
	%171 = icmp uge i64 0, %170
	br i1 %171, label %172, label %174

172:
	store i32 2, i32* @.ErrorNumber
	%173 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

174:
	%175 = icmp ult i64 0, 0
	br i1 %175, label %172, label %176

176:
	%177 = mul i64 0, %170
	%178 = add i64 %177, 0
	%179 = getelementptr [4 x i32], [4 x i32]* @bzzz, i64 0, i64 %178
	%180 = load i32, i32* %179
	%181 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %181
	%182 = call i32 (i8*, ...) @printf([3 x i8]* %181, i32 %180)
	%183 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
	%184 = load i64, i64* %.a_0
	%185 = icmp uge i64 0, %184
	br i1 %185, label %186, label %188

186:
	store i32 2, i32* @.ErrorNumber
	%187 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

188:
	%189 = icmp ult i64 0, 0
	br i1 %189, label %186, label %190

190:
	%191 = load i64, i64* %.a_1
	%192 = icmp uge i64 0, %191
	br i1 %192, label %193, label %195

193:
	store i32 2, i32* @.ErrorNumber
	%194 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

195:
	%196 = icmp ult i64 0, 0
	br i1 %196, label %193, label %197

197:
	%198 = mul i64 0, %191
	%199 = add i64 %198, 0
	%200 = getelementptr [6 x i32], [6 x i32]* %1, i64 0, i64 %199
	%201 = load i32, i32* %200
	%202 = load i64, i64* %.a_0
	%203 = icmp uge i64 0, %202
	br i1 %203, label %204, label %206

204:
	store i32 2, i32* @.ErrorNumber
	%205 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

206:
	%207 = icmp ult i64 0, 0
	br i1 %207, label %204, label %208

208:
	%209 = load i64, i64* %.a_1
	%210 = icmp uge i64 1, %209
	br i1 %210, label %211, label %213

211:
	store i32 2, i32* @.ErrorNumber
	%212 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

213:
	%214 = icmp ult i64 1, 0
	br i1 %214, label %211, label %215

215:
	%216 = mul i64 0, %209
	%217 = add i64 %216, 1
	%218 = getelementptr [6 x i32], [6 x i32]* %1, i64 0, i64 %217
	%219 = load i32, i32* %218
	%220 = add i32 %201, %219
	%221 = load i64, i64* %.a_0
	%222 = icmp uge i64 1, %221
	br i1 %222, label %223, label %225

223:
	store i32 2, i32* @.ErrorNumber
	%224 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

225:
	%226 = icmp ult i64 1, 0
	br i1 %226, label %223, label %227

227:
	%228 = load i64, i64* %.a_1
	%229 = icmp uge i64 0, %228
	br i1 %229, label %230, label %232

230:
	store i32 2, i32* @.ErrorNumber
	%231 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

232:
	%233 = icmp ult i64 0, 0
	br i1 %233, label %230, label %234

234:
	%235 = mul i64 1, %228
	%236 = add i64 %235, 0
	%237 = getelementptr [6 x i32], [6 x i32]* %1, i64 0, i64 %236
	%238 = load i32, i32* %237
	%239 = add i32 %220, %238
	%240 = load i64, i64* %.a_0
	%241 = icmp uge i64 1, %240
	br i1 %241, label %242, label %244

242:
	store i32 2, i32* @.ErrorNumber
	%243 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

244:
	%245 = icmp ult i64 1, 0
	br i1 %245, label %242, label %246

246:
	%247 = load i64, i64* %.a_1
	%248 = icmp uge i64 1, %247
	br i1 %248, label %249, label %251

249:
	store i32 2, i32* @.ErrorNumber
	%250 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

251:
	%252 = icmp ult i64 1, 0
	br i1 %252, label %249, label %253

253:
	%254 = mul i64 1, %247
	%255 = add i64 %254, 1
	%256 = getelementptr [6 x i32], [6 x i32]* %1, i64 0, i64 %255
	%257 = load i32, i32* %256
	%258 = add i32 %239, %257
	%259 = load i64, i64* @.bzzz_0
	%260 = icmp uge i64 0, %259
	br i1 %260, label %261, label %263

261:
	store i32 2, i32* @.ErrorNumber
	%262 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

263:
	%264 = icmp ult i64 0, 0
	br i1 %264, label %261, label %265

265:
	%266 = load i64, i64* @.bzzz_1
	%267 = icmp uge i64 0, %266
	br i1 %267, label %268, label %270

268:
	store i32 2, i32* @.ErrorNumber
	%269 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

270:
	%271 = icmp ult i64 0, 0
	br i1 %271, label %268, label %272

272:
	%273 = mul i64 0, %266
	%274 = add i64 %273, 0
	%275 = getelementptr [4 x i32], [4 x i32]* @bzzz, i64 0, i64 %274
	%276 = load i32, i32* %275
	%277 = add i32 %258, %276
	%278 = alloca [3 x i8]
	store [3 x i8] c"%d\00", [3 x i8]* %278
	%279 = call i32 (i8*, ...) @printf([3 x i8]* %278, i32 %277)
	%280 = call i32 @puts(i8* getelementptr ([1 x i8], [1 x i8]* @vbEmpty, i32 0, i32 0))
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
	%1 = load i64, i64* @.bzzz_0
	%2 = icmp uge i64 0, %1
	br i1 %2, label %3, label %5

3:
	store i32 2, i32* @.ErrorNumber
	%4 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

5:
	%6 = icmp ult i64 0, 0
	br i1 %6, label %3, label %7

7:
	%8 = load i64, i64* @.bzzz_1
	%9 = icmp uge i64 0, %8
	br i1 %9, label %10, label %12

10:
	store i32 2, i32* @.ErrorNumber
	%11 = call i8* @strcpy([256 x i8]* @.ErrorMessage, [27 x i8]* @.arrayIndexOutOfBounds)
	call void @.throwException()
	unreachable

12:
	%13 = icmp ult i64 0, 0
	br i1 %13, label %10, label %14

14:
	%15 = mul i64 0, %8
	%16 = add i64 %15, 0
	%17 = getelementptr [4 x i32], [4 x i32]* @bzzz, i64 0, i64 %16
	%18 = trunc i64 4 to i32
	store i32 %18, i32* %17
	call void @AddArrayNumbers()
	ret void
}
