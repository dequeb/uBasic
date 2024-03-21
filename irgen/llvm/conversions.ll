; ModuleID = 'conversions.c'
source_filename = "conversions.c"
target datalayout = "e-m:o-i64:64-i128:128-n32:64-S128"
target triple = "arm64-apple-macosx14.0.0"

@i = common global i32 0, align 4
@l = common global i64 0, align 8
@f = common global float 0.000000e+00, align 4
@d = common global double 0.000000e+00, align 8
@.str = private unnamed_addr constant [30 x i8] c"i: %i, l: %li, d: %lf, f: %f\0A\00", align 1

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define i32 @main() #0 {
  %1 = alloca i32, align 4
  store i32 0, ptr %1, align 4
  store i32 10, ptr @i, align 4
  %2 = load i32, ptr @i, align 4
  %3 = sext i32 %2 to i64
  store i64 %3, ptr @l, align 8
  %4 = load i64, ptr @l, align 8
  %5 = trunc i64 %4 to i32
  store i32 %5, ptr @i, align 4
  store float 0x4041801CC0000000, ptr @f, align 4
  %6 = load float, ptr @f, align 4
  %7 = fpext float %6 to double
  store double %7, ptr @d, align 8
  %8 = load double, ptr @d, align 8
  %9 = fptrunc double %8 to float
  store float %9, ptr @f, align 4
  %10 = load i32, ptr @i, align 4
  %11 = load i64, ptr @l, align 8
  %12 = load double, ptr @d, align 8
  %13 = load float, ptr @f, align 4
  %14 = fpext float %13 to double
  %15 = call i32 (ptr, ...) @printf(ptr noundef @.str, i32 noundef %10, i64 noundef %11, double noundef %12, double noundef %14)
  ret i32 0
}

declare i32 @printf(ptr noundef, ...) #1

attributes #0 = { noinline nounwind optnone ssp uwtable(sync) "frame-pointer"="non-leaf" "min-legal-vector-width"="0" "no-trapping-math"="true" "probe-stack"="__chkstk_darwin" "stack-protector-buffer-size"="8" "target-cpu"="apple-m1" "target-features"="+aes,+crc,+crypto,+dotprod,+fp-armv8,+fp16fml,+fullfp16,+lse,+neon,+ras,+rcpc,+rdm,+sha2,+sha3,+sm4,+v8.1a,+v8.2a,+v8.3a,+v8.4a,+v8.5a,+v8a,+zcm,+zcz" }
attributes #1 = { "frame-pointer"="non-leaf" "no-trapping-math"="true" "probe-stack"="__chkstk_darwin" "stack-protector-buffer-size"="8" "target-cpu"="apple-m1" "target-features"="+aes,+crc,+crypto,+dotprod,+fp-armv8,+fp16fml,+fullfp16,+lse,+neon,+ras,+rcpc,+rdm,+sha2,+sha3,+sm4,+v8.1a,+v8.2a,+v8.3a,+v8.4a,+v8.5a,+v8a,+zcm,+zcz" }

!llvm.module.flags = !{!0, !1, !2, !3, !4}
!llvm.ident = !{!5}

!0 = !{i32 2, !"SDK Version", [2 x i32] [i32 14, i32 0]}
!1 = !{i32 1, !"wchar_size", i32 4}
!2 = !{i32 8, !"PIC Level", i32 2}
!3 = !{i32 7, !"uwtable", i32 1}
!4 = !{i32 7, !"frame-pointer", i32 1}
!5 = !{!"Apple clang version 15.0.0 (clang-1500.0.40.1)"}
