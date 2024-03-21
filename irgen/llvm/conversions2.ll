; ModuleID = 'conversions2.c'
source_filename = "conversions2.c"
target datalayout = "e-m:o-i64:64-i128:128-n32:64-S128"
target triple = "arm64-apple-macosx14.0.0"

@i = common global i32 0, align 4
@f = common global float 0.000000e+00, align 4
@l = common global i64 0, align 8
@d = common global double 0.000000e+00, align 8
@.str = private unnamed_addr constant [30 x i8] c"i: %i, l: %li, d: %lf, f: %f\0A\00", align 1

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define i32 @main() #0 {
  %1 = alloca i32, align 4
  store i32 0, ptr %1, align 4
  store i32 10, ptr @i, align 4
  %2 = load i32, ptr @i, align 4
  %3 = sitofp i32 %2 to float
  store float %3, ptr @f, align 4
  store i64 1234567890, ptr @l, align 8
  %4 = load i64, ptr @l, align 8
  %5 = sitofp i64 %4 to double
  store double %5, ptr @d, align 8
  %6 = load i32, ptr @i, align 4
  %7 = load i64, ptr @l, align 8
  %8 = load double, ptr @d, align 8
  %9 = load float, ptr @f, align 4
  %10 = fpext float %9 to double
  %11 = call i32 (ptr, ...) @printf(ptr noundef @.str, i32 noundef %6, i64 noundef %7, double noundef %8, double noundef %10)
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
