; ModuleID = 'gc.c'
source_filename = "gc.c"
target datalayout = "e-m:o-i64:64-i128:128-n32:64-S128"
target triple = "arm64-apple-macosx14.0.0"

%struct.GarbageCollector = type { ptr, i8, ptr, i64 }
%struct.Allocation = type { ptr, i64, i8, ptr, ptr }
%struct.AllocationMap = type { i64, i64, double, double, double, i64, i64, ptr }

@__stderrp = external global ptr, align 8
@.str = private unnamed_addr constant [60 x i8] c"[%s] %s:%s:%d: Ignoring request to free unknown pointer %p\0A\00", align 1
@log_level_strings = external global [0 x ptr], align 8
@__func__.gc_free = private unnamed_addr constant [8 x i8] c"gc_free\00", align 1
@.str.1 = private unnamed_addr constant [5 x i8] c"gc.c\00", align 1
@gc = common global %struct.GarbageCollector zeroinitializer, align 8

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define ptr @gc_malloc(ptr noundef %0, i64 noundef %1) #0 {
  %3 = alloca ptr, align 8
  %4 = alloca i64, align 8
  store ptr %0, ptr %3, align 8
  store i64 %1, ptr %4, align 8
  %5 = load ptr, ptr %3, align 8
  %6 = load i64, ptr %4, align 8
  %7 = call ptr @gc_malloc_ext(ptr noundef %5, i64 noundef %6, ptr noundef null)
  ret ptr %7
}

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define ptr @gc_malloc_ext(ptr noundef %0, i64 noundef %1, ptr noundef %2) #0 {
  %4 = alloca ptr, align 8
  %5 = alloca i64, align 8
  %6 = alloca ptr, align 8
  store ptr %0, ptr %4, align 8
  store i64 %1, ptr %5, align 8
  store ptr %2, ptr %6, align 8
  %7 = load ptr, ptr %4, align 8
  %8 = load i64, ptr %5, align 8
  %9 = load ptr, ptr %6, align 8
  %10 = call ptr @gc_allocate(ptr noundef %7, i64 noundef 0, i64 noundef %8, ptr noundef %9)
  ret ptr %10
}

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define ptr @gc_malloc_static(ptr noundef %0, i64 noundef %1, ptr noundef %2) #0 {
  %4 = alloca ptr, align 8
  %5 = alloca i64, align 8
  %6 = alloca ptr, align 8
  %7 = alloca ptr, align 8
  store ptr %0, ptr %4, align 8
  store i64 %1, ptr %5, align 8
  store ptr %2, ptr %6, align 8
  %8 = load ptr, ptr %4, align 8
  %9 = load i64, ptr %5, align 8
  %10 = load ptr, ptr %6, align 8
  %11 = call ptr @gc_malloc_ext(ptr noundef %8, i64 noundef %9, ptr noundef %10)
  store ptr %11, ptr %7, align 8
  %12 = load ptr, ptr %4, align 8
  %13 = load ptr, ptr %7, align 8
  call void @gc_make_root(ptr noundef %12, ptr noundef %13)
  %14 = load ptr, ptr %7, align 8
  ret ptr %14
}

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define internal void @gc_make_root(ptr noundef %0, ptr noundef %1) #0 {
  %3 = alloca ptr, align 8
  %4 = alloca ptr, align 8
  %5 = alloca ptr, align 8
  store ptr %0, ptr %3, align 8
  store ptr %1, ptr %4, align 8
  %6 = load ptr, ptr %3, align 8
  %7 = getelementptr inbounds %struct.GarbageCollector, ptr %6, i32 0, i32 0
  %8 = load ptr, ptr %7, align 8
  %9 = load ptr, ptr %4, align 8
  %10 = call ptr @gc_allocation_map_get(ptr noundef %8, ptr noundef %9)
  store ptr %10, ptr %5, align 8
  %11 = load ptr, ptr %5, align 8
  %12 = icmp ne ptr %11, null
  br i1 %12, label %13, label %20

13:                                               ; preds = %2
  %14 = load ptr, ptr %5, align 8
  %15 = getelementptr inbounds %struct.Allocation, ptr %14, i32 0, i32 2
  %16 = load i8, ptr %15, align 8
  %17 = sext i8 %16 to i32
  %18 = or i32 %17, 1
  %19 = trunc i32 %18 to i8
  store i8 %19, ptr %15, align 8
  br label %20

20:                                               ; preds = %13, %2
  ret void
}

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define ptr @gc_make_static(ptr noundef %0, ptr noundef %1) #0 {
  %3 = alloca ptr, align 8
  %4 = alloca ptr, align 8
  store ptr %0, ptr %3, align 8
  store ptr %1, ptr %4, align 8
  %5 = load ptr, ptr %3, align 8
  %6 = load ptr, ptr %4, align 8
  call void @gc_make_root(ptr noundef %5, ptr noundef %6)
  %7 = load ptr, ptr %4, align 8
  ret ptr %7
}

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define internal ptr @gc_allocate(ptr noundef %0, i64 noundef %1, i64 noundef %2, ptr noundef %3) #0 {
  %5 = alloca ptr, align 8
  %6 = alloca i64, align 8
  %7 = alloca i64, align 8
  %8 = alloca ptr, align 8
  %9 = alloca i64, align 8
  %10 = alloca ptr, align 8
  %11 = alloca i64, align 8
  %12 = alloca ptr, align 8
  store ptr %0, ptr %5, align 8
  store i64 %1, ptr %6, align 8
  store i64 %2, ptr %7, align 8
  store ptr %3, ptr %8, align 8
  %13 = load ptr, ptr %5, align 8
  %14 = call zeroext i1 @gc_needs_sweep(ptr noundef %13)
  br i1 %14, label %15, label %25

15:                                               ; preds = %4
  %16 = load ptr, ptr %5, align 8
  %17 = getelementptr inbounds %struct.GarbageCollector, ptr %16, i32 0, i32 1
  %18 = load i8, ptr %17, align 8
  %19 = trunc i8 %18 to i1
  br i1 %19, label %25, label %20

20:                                               ; preds = %15
  %21 = load ptr, ptr %5, align 8
  %22 = call i64 @gc_run(ptr noundef %21)
  store i64 %22, ptr %9, align 8
  br label %23

23:                                               ; preds = %20
  br label %24

24:                                               ; preds = %23
  br label %25

25:                                               ; preds = %24, %15, %4
  %26 = load i64, ptr %6, align 8
  %27 = load i64, ptr %7, align 8
  %28 = call ptr @gc_mcalloc(i64 noundef %26, i64 noundef %27)
  store ptr %28, ptr %10, align 8
  %29 = load i64, ptr %6, align 8
  %30 = icmp ne i64 %29, 0
  br i1 %30, label %31, label %35

31:                                               ; preds = %25
  %32 = load i64, ptr %6, align 8
  %33 = load i64, ptr %7, align 8
  %34 = mul i64 %32, %33
  br label %37

35:                                               ; preds = %25
  %36 = load i64, ptr %7, align 8
  br label %37

37:                                               ; preds = %35, %31
  %38 = phi i64 [ %34, %31 ], [ %36, %35 ]
  store i64 %38, ptr %11, align 8
  %39 = load ptr, ptr %10, align 8
  %40 = icmp ne ptr %39, null
  br i1 %40, label %60, label %41

41:                                               ; preds = %37
  %42 = load ptr, ptr %5, align 8
  %43 = getelementptr inbounds %struct.GarbageCollector, ptr %42, i32 0, i32 1
  %44 = load i8, ptr %43, align 8
  %45 = trunc i8 %44 to i1
  br i1 %45, label %60, label %46

46:                                               ; preds = %41
  %47 = call ptr @__error()
  %48 = load i32, ptr %47, align 4
  %49 = icmp eq i32 %48, 35
  br i1 %49, label %54, label %50

50:                                               ; preds = %46
  %51 = call ptr @__error()
  %52 = load i32, ptr %51, align 4
  %53 = icmp eq i32 %52, 12
  br i1 %53, label %54, label %60

54:                                               ; preds = %50, %46
  %55 = load ptr, ptr %5, align 8
  %56 = call i64 @gc_run(ptr noundef %55)
  %57 = load i64, ptr %6, align 8
  %58 = load i64, ptr %7, align 8
  %59 = call ptr @gc_mcalloc(i64 noundef %57, i64 noundef %58)
  store ptr %59, ptr %10, align 8
  br label %60

60:                                               ; preds = %54, %50, %41, %37
  %61 = load ptr, ptr %10, align 8
  %62 = icmp ne ptr %61, null
  br i1 %62, label %63, label %84

63:                                               ; preds = %60
  br label %64

64:                                               ; preds = %63
  br label %65

65:                                               ; preds = %64
  %66 = load ptr, ptr %5, align 8
  %67 = getelementptr inbounds %struct.GarbageCollector, ptr %66, i32 0, i32 0
  %68 = load ptr, ptr %67, align 8
  %69 = load ptr, ptr %10, align 8
  %70 = load i64, ptr %11, align 8
  %71 = load ptr, ptr %8, align 8
  %72 = call ptr @gc_allocation_map_put(ptr noundef %68, ptr noundef %69, i64 noundef %70, ptr noundef %71)
  store ptr %72, ptr %12, align 8
  %73 = load ptr, ptr %12, align 8
  %74 = icmp ne ptr %73, null
  br i1 %74, label %75, label %81

75:                                               ; preds = %65
  br label %76

76:                                               ; preds = %75
  br label %77

77:                                               ; preds = %76
  %78 = load ptr, ptr %12, align 8
  %79 = getelementptr inbounds %struct.Allocation, ptr %78, i32 0, i32 0
  %80 = load ptr, ptr %79, align 8
  store ptr %80, ptr %10, align 8
  br label %83

81:                                               ; preds = %65
  %82 = load ptr, ptr %10, align 8
  call void @free(ptr noundef %82)
  store ptr null, ptr %10, align 8
  br label %83

83:                                               ; preds = %81, %77
  br label %84

84:                                               ; preds = %83, %60
  %85 = load ptr, ptr %10, align 8
  ret ptr %85
}

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define ptr @gc_calloc(ptr noundef %0, i64 noundef %1, i64 noundef %2) #0 {
  %4 = alloca ptr, align 8
  %5 = alloca i64, align 8
  %6 = alloca i64, align 8
  store ptr %0, ptr %4, align 8
  store i64 %1, ptr %5, align 8
  store i64 %2, ptr %6, align 8
  %7 = load ptr, ptr %4, align 8
  %8 = load i64, ptr %5, align 8
  %9 = load i64, ptr %6, align 8
  %10 = call ptr @gc_calloc_ext(ptr noundef %7, i64 noundef %8, i64 noundef %9, ptr noundef null)
  ret ptr %10
}

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define ptr @gc_calloc_ext(ptr noundef %0, i64 noundef %1, i64 noundef %2, ptr noundef %3) #0 {
  %5 = alloca ptr, align 8
  %6 = alloca i64, align 8
  %7 = alloca i64, align 8
  %8 = alloca ptr, align 8
  store ptr %0, ptr %5, align 8
  store i64 %1, ptr %6, align 8
  store i64 %2, ptr %7, align 8
  store ptr %3, ptr %8, align 8
  %9 = load ptr, ptr %5, align 8
  %10 = load i64, ptr %6, align 8
  %11 = load i64, ptr %7, align 8
  %12 = load ptr, ptr %8, align 8
  %13 = call ptr @gc_allocate(ptr noundef %9, i64 noundef %10, i64 noundef %11, ptr noundef %12)
  ret ptr %13
}

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define ptr @gc_realloc(ptr noundef %0, ptr noundef %1, i64 noundef %2) #0 {
  %4 = alloca ptr, align 8
  %5 = alloca ptr, align 8
  %6 = alloca ptr, align 8
  %7 = alloca i64, align 8
  %8 = alloca ptr, align 8
  %9 = alloca ptr, align 8
  %10 = alloca ptr, align 8
  %11 = alloca ptr, align 8
  store ptr %0, ptr %5, align 8
  store ptr %1, ptr %6, align 8
  store i64 %2, ptr %7, align 8
  %12 = load ptr, ptr %5, align 8
  %13 = getelementptr inbounds %struct.GarbageCollector, ptr %12, i32 0, i32 0
  %14 = load ptr, ptr %13, align 8
  %15 = load ptr, ptr %6, align 8
  %16 = call ptr @gc_allocation_map_get(ptr noundef %14, ptr noundef %15)
  store ptr %16, ptr %8, align 8
  %17 = load ptr, ptr %6, align 8
  %18 = icmp ne ptr %17, null
  br i1 %18, label %19, label %24

19:                                               ; preds = %3
  %20 = load ptr, ptr %8, align 8
  %21 = icmp ne ptr %20, null
  br i1 %21, label %24, label %22

22:                                               ; preds = %19
  %23 = call ptr @__error()
  store i32 22, ptr %23, align 4
  store ptr null, ptr %4, align 8
  br label %69

24:                                               ; preds = %19, %3
  %25 = load ptr, ptr %6, align 8
  %26 = load i64, ptr %7, align 8
  %27 = call ptr @realloc(ptr noundef %25, i64 noundef %26) #10
  store ptr %27, ptr %9, align 8
  %28 = load ptr, ptr %9, align 8
  %29 = icmp ne ptr %28, null
  br i1 %29, label %31, label %30

30:                                               ; preds = %24
  store ptr null, ptr %4, align 8
  br label %69

31:                                               ; preds = %24
  %32 = load ptr, ptr %6, align 8
  %33 = icmp ne ptr %32, null
  br i1 %33, label %44, label %34

34:                                               ; preds = %31
  %35 = load ptr, ptr %5, align 8
  %36 = getelementptr inbounds %struct.GarbageCollector, ptr %35, i32 0, i32 0
  %37 = load ptr, ptr %36, align 8
  %38 = load ptr, ptr %9, align 8
  %39 = load i64, ptr %7, align 8
  %40 = call ptr @gc_allocation_map_put(ptr noundef %37, ptr noundef %38, i64 noundef %39, ptr noundef null)
  store ptr %40, ptr %10, align 8
  %41 = load ptr, ptr %10, align 8
  %42 = getelementptr inbounds %struct.Allocation, ptr %41, i32 0, i32 0
  %43 = load ptr, ptr %42, align 8
  store ptr %43, ptr %4, align 8
  br label %69

44:                                               ; preds = %31
  %45 = load ptr, ptr %6, align 8
  %46 = load ptr, ptr %9, align 8
  %47 = icmp eq ptr %45, %46
  br i1 %47, label %48, label %52

48:                                               ; preds = %44
  %49 = load i64, ptr %7, align 8
  %50 = load ptr, ptr %8, align 8
  %51 = getelementptr inbounds %struct.Allocation, ptr %50, i32 0, i32 1
  store i64 %49, ptr %51, align 8
  br label %67

52:                                               ; preds = %44
  %53 = load ptr, ptr %8, align 8
  %54 = getelementptr inbounds %struct.Allocation, ptr %53, i32 0, i32 3
  %55 = load ptr, ptr %54, align 8
  store ptr %55, ptr %11, align 8
  %56 = load ptr, ptr %5, align 8
  %57 = getelementptr inbounds %struct.GarbageCollector, ptr %56, i32 0, i32 0
  %58 = load ptr, ptr %57, align 8
  %59 = load ptr, ptr %6, align 8
  call void @gc_allocation_map_remove(ptr noundef %58, ptr noundef %59, i1 noundef zeroext true)
  %60 = load ptr, ptr %5, align 8
  %61 = getelementptr inbounds %struct.GarbageCollector, ptr %60, i32 0, i32 0
  %62 = load ptr, ptr %61, align 8
  %63 = load ptr, ptr %9, align 8
  %64 = load i64, ptr %7, align 8
  %65 = load ptr, ptr %11, align 8
  %66 = call ptr @gc_allocation_map_put(ptr noundef %62, ptr noundef %63, i64 noundef %64, ptr noundef %65)
  br label %67

67:                                               ; preds = %52, %48
  %68 = load ptr, ptr %9, align 8
  store ptr %68, ptr %4, align 8
  br label %69

69:                                               ; preds = %67, %34, %30, %22
  %70 = load ptr, ptr %4, align 8
  ret ptr %70
}

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define internal ptr @gc_allocation_map_get(ptr noundef %0, ptr noundef %1) #0 {
  %3 = alloca ptr, align 8
  %4 = alloca ptr, align 8
  %5 = alloca ptr, align 8
  %6 = alloca i64, align 8
  %7 = alloca ptr, align 8
  store ptr %0, ptr %4, align 8
  store ptr %1, ptr %5, align 8
  %8 = load ptr, ptr %5, align 8
  %9 = call i64 @gc_hash(ptr noundef %8)
  %10 = load ptr, ptr %4, align 8
  %11 = getelementptr inbounds %struct.AllocationMap, ptr %10, i32 0, i32 0
  %12 = load i64, ptr %11, align 8
  %13 = urem i64 %9, %12
  store i64 %13, ptr %6, align 8
  %14 = load ptr, ptr %4, align 8
  %15 = getelementptr inbounds %struct.AllocationMap, ptr %14, i32 0, i32 7
  %16 = load ptr, ptr %15, align 8
  %17 = load i64, ptr %6, align 8
  %18 = getelementptr inbounds ptr, ptr %16, i64 %17
  %19 = load ptr, ptr %18, align 8
  store ptr %19, ptr %7, align 8
  br label %20

20:                                               ; preds = %31, %2
  %21 = load ptr, ptr %7, align 8
  %22 = icmp ne ptr %21, null
  br i1 %22, label %23, label %35

23:                                               ; preds = %20
  %24 = load ptr, ptr %7, align 8
  %25 = getelementptr inbounds %struct.Allocation, ptr %24, i32 0, i32 0
  %26 = load ptr, ptr %25, align 8
  %27 = load ptr, ptr %5, align 8
  %28 = icmp eq ptr %26, %27
  br i1 %28, label %29, label %31

29:                                               ; preds = %23
  %30 = load ptr, ptr %7, align 8
  store ptr %30, ptr %3, align 8
  br label %36

31:                                               ; preds = %23
  %32 = load ptr, ptr %7, align 8
  %33 = getelementptr inbounds %struct.Allocation, ptr %32, i32 0, i32 4
  %34 = load ptr, ptr %33, align 8
  store ptr %34, ptr %7, align 8
  br label %20, !llvm.loop !6

35:                                               ; preds = %20
  store ptr null, ptr %3, align 8
  br label %36

36:                                               ; preds = %35, %29
  %37 = load ptr, ptr %3, align 8
  ret ptr %37
}

declare ptr @__error() #1

; Function Attrs: allocsize(1)
declare ptr @realloc(ptr noundef, i64 noundef) #2

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define internal ptr @gc_allocation_map_put(ptr noundef %0, ptr noundef %1, i64 noundef %2, ptr noundef %3) #0 {
  %5 = alloca ptr, align 8
  %6 = alloca ptr, align 8
  %7 = alloca ptr, align 8
  %8 = alloca i64, align 8
  %9 = alloca ptr, align 8
  %10 = alloca i64, align 8
  %11 = alloca ptr, align 8
  %12 = alloca ptr, align 8
  %13 = alloca ptr, align 8
  %14 = alloca ptr, align 8
  store ptr %0, ptr %6, align 8
  store ptr %1, ptr %7, align 8
  store i64 %2, ptr %8, align 8
  store ptr %3, ptr %9, align 8
  %15 = load ptr, ptr %7, align 8
  %16 = call i64 @gc_hash(ptr noundef %15)
  %17 = load ptr, ptr %6, align 8
  %18 = getelementptr inbounds %struct.AllocationMap, ptr %17, i32 0, i32 0
  %19 = load i64, ptr %18, align 8
  %20 = urem i64 %16, %19
  store i64 %20, ptr %10, align 8
  br label %21

21:                                               ; preds = %4
  br label %22

22:                                               ; preds = %21
  %23 = load ptr, ptr %7, align 8
  %24 = load i64, ptr %8, align 8
  %25 = load ptr, ptr %9, align 8
  %26 = call ptr @gc_allocation_new(ptr noundef %23, i64 noundef %24, ptr noundef %25)
  store ptr %26, ptr %11, align 8
  %27 = load ptr, ptr %6, align 8
  %28 = getelementptr inbounds %struct.AllocationMap, ptr %27, i32 0, i32 7
  %29 = load ptr, ptr %28, align 8
  %30 = load i64, ptr %10, align 8
  %31 = getelementptr inbounds ptr, ptr %29, i64 %30
  %32 = load ptr, ptr %31, align 8
  store ptr %32, ptr %12, align 8
  store ptr null, ptr %13, align 8
  br label %33

33:                                               ; preds = %66, %22
  %34 = load ptr, ptr %12, align 8
  %35 = icmp ne ptr %34, null
  br i1 %35, label %36, label %71

36:                                               ; preds = %33
  %37 = load ptr, ptr %12, align 8
  %38 = getelementptr inbounds %struct.Allocation, ptr %37, i32 0, i32 0
  %39 = load ptr, ptr %38, align 8
  %40 = load ptr, ptr %7, align 8
  %41 = icmp eq ptr %39, %40
  br i1 %41, label %42, label %66

42:                                               ; preds = %36
  %43 = load ptr, ptr %12, align 8
  %44 = getelementptr inbounds %struct.Allocation, ptr %43, i32 0, i32 4
  %45 = load ptr, ptr %44, align 8
  %46 = load ptr, ptr %11, align 8
  %47 = getelementptr inbounds %struct.Allocation, ptr %46, i32 0, i32 4
  store ptr %45, ptr %47, align 8
  %48 = load ptr, ptr %13, align 8
  %49 = icmp ne ptr %48, null
  br i1 %49, label %57, label %50

50:                                               ; preds = %42
  %51 = load ptr, ptr %11, align 8
  %52 = load ptr, ptr %6, align 8
  %53 = getelementptr inbounds %struct.AllocationMap, ptr %52, i32 0, i32 7
  %54 = load ptr, ptr %53, align 8
  %55 = load i64, ptr %10, align 8
  %56 = getelementptr inbounds ptr, ptr %54, i64 %55
  store ptr %51, ptr %56, align 8
  br label %61

57:                                               ; preds = %42
  %58 = load ptr, ptr %11, align 8
  %59 = load ptr, ptr %13, align 8
  %60 = getelementptr inbounds %struct.Allocation, ptr %59, i32 0, i32 4
  store ptr %58, ptr %60, align 8
  br label %61

61:                                               ; preds = %57, %50
  %62 = load ptr, ptr %12, align 8
  call void @gc_allocation_delete(ptr noundef %62)
  br label %63

63:                                               ; preds = %61
  br label %64

64:                                               ; preds = %63
  %65 = load ptr, ptr %11, align 8
  store ptr %65, ptr %5, align 8
  br label %104

66:                                               ; preds = %36
  %67 = load ptr, ptr %12, align 8
  store ptr %67, ptr %13, align 8
  %68 = load ptr, ptr %12, align 8
  %69 = getelementptr inbounds %struct.Allocation, ptr %68, i32 0, i32 4
  %70 = load ptr, ptr %69, align 8
  store ptr %70, ptr %12, align 8
  br label %33, !llvm.loop !8

71:                                               ; preds = %33
  %72 = load ptr, ptr %6, align 8
  %73 = getelementptr inbounds %struct.AllocationMap, ptr %72, i32 0, i32 7
  %74 = load ptr, ptr %73, align 8
  %75 = load i64, ptr %10, align 8
  %76 = getelementptr inbounds ptr, ptr %74, i64 %75
  %77 = load ptr, ptr %76, align 8
  store ptr %77, ptr %12, align 8
  %78 = load ptr, ptr %12, align 8
  %79 = load ptr, ptr %11, align 8
  %80 = getelementptr inbounds %struct.Allocation, ptr %79, i32 0, i32 4
  store ptr %78, ptr %80, align 8
  %81 = load ptr, ptr %11, align 8
  %82 = load ptr, ptr %6, align 8
  %83 = getelementptr inbounds %struct.AllocationMap, ptr %82, i32 0, i32 7
  %84 = load ptr, ptr %83, align 8
  %85 = load i64, ptr %10, align 8
  %86 = getelementptr inbounds ptr, ptr %84, i64 %85
  store ptr %81, ptr %86, align 8
  %87 = load ptr, ptr %6, align 8
  %88 = getelementptr inbounds %struct.AllocationMap, ptr %87, i32 0, i32 6
  %89 = load i64, ptr %88, align 8
  %90 = add i64 %89, 1
  store i64 %90, ptr %88, align 8
  br label %91

91:                                               ; preds = %71
  br label %92

92:                                               ; preds = %91
  %93 = load ptr, ptr %11, align 8
  %94 = getelementptr inbounds %struct.Allocation, ptr %93, i32 0, i32 0
  %95 = load ptr, ptr %94, align 8
  store ptr %95, ptr %14, align 8
  %96 = load ptr, ptr %6, align 8
  %97 = call zeroext i1 @gc_allocation_map_resize_to_fit(ptr noundef %96)
  br i1 %97, label %98, label %102

98:                                               ; preds = %92
  %99 = load ptr, ptr %6, align 8
  %100 = load ptr, ptr %14, align 8
  %101 = call ptr @gc_allocation_map_get(ptr noundef %99, ptr noundef %100)
  store ptr %101, ptr %11, align 8
  br label %102

102:                                              ; preds = %98, %92
  %103 = load ptr, ptr %11, align 8
  store ptr %103, ptr %5, align 8
  br label %104

104:                                              ; preds = %102, %64
  %105 = load ptr, ptr %5, align 8
  ret ptr %105
}

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define internal void @gc_allocation_map_remove(ptr noundef %0, ptr noundef %1, i1 noundef zeroext %2) #0 {
  %4 = alloca ptr, align 8
  %5 = alloca ptr, align 8
  %6 = alloca i8, align 1
  %7 = alloca i64, align 8
  %8 = alloca ptr, align 8
  %9 = alloca ptr, align 8
  %10 = alloca ptr, align 8
  store ptr %0, ptr %4, align 8
  store ptr %1, ptr %5, align 8
  %11 = zext i1 %2 to i8
  store i8 %11, ptr %6, align 1
  %12 = load ptr, ptr %5, align 8
  %13 = call i64 @gc_hash(ptr noundef %12)
  %14 = load ptr, ptr %4, align 8
  %15 = getelementptr inbounds %struct.AllocationMap, ptr %14, i32 0, i32 0
  %16 = load i64, ptr %15, align 8
  %17 = urem i64 %13, %16
  store i64 %17, ptr %7, align 8
  %18 = load ptr, ptr %4, align 8
  %19 = getelementptr inbounds %struct.AllocationMap, ptr %18, i32 0, i32 7
  %20 = load ptr, ptr %19, align 8
  %21 = load i64, ptr %7, align 8
  %22 = getelementptr inbounds ptr, ptr %20, i64 %21
  %23 = load ptr, ptr %22, align 8
  store ptr %23, ptr %8, align 8
  store ptr null, ptr %9, align 8
  br label %24

24:                                               ; preds = %62, %3
  %25 = load ptr, ptr %8, align 8
  %26 = icmp ne ptr %25, null
  br i1 %26, label %27, label %64

27:                                               ; preds = %24
  %28 = load ptr, ptr %8, align 8
  %29 = getelementptr inbounds %struct.Allocation, ptr %28, i32 0, i32 4
  %30 = load ptr, ptr %29, align 8
  store ptr %30, ptr %10, align 8
  %31 = load ptr, ptr %8, align 8
  %32 = getelementptr inbounds %struct.Allocation, ptr %31, i32 0, i32 0
  %33 = load ptr, ptr %32, align 8
  %34 = load ptr, ptr %5, align 8
  %35 = icmp eq ptr %33, %34
  br i1 %35, label %36, label %60

36:                                               ; preds = %27
  %37 = load ptr, ptr %9, align 8
  %38 = icmp ne ptr %37, null
  br i1 %38, label %48, label %39

39:                                               ; preds = %36
  %40 = load ptr, ptr %8, align 8
  %41 = getelementptr inbounds %struct.Allocation, ptr %40, i32 0, i32 4
  %42 = load ptr, ptr %41, align 8
  %43 = load ptr, ptr %4, align 8
  %44 = getelementptr inbounds %struct.AllocationMap, ptr %43, i32 0, i32 7
  %45 = load ptr, ptr %44, align 8
  %46 = load i64, ptr %7, align 8
  %47 = getelementptr inbounds ptr, ptr %45, i64 %46
  store ptr %42, ptr %47, align 8
  br label %54

48:                                               ; preds = %36
  %49 = load ptr, ptr %8, align 8
  %50 = getelementptr inbounds %struct.Allocation, ptr %49, i32 0, i32 4
  %51 = load ptr, ptr %50, align 8
  %52 = load ptr, ptr %9, align 8
  %53 = getelementptr inbounds %struct.Allocation, ptr %52, i32 0, i32 4
  store ptr %51, ptr %53, align 8
  br label %54

54:                                               ; preds = %48, %39
  %55 = load ptr, ptr %8, align 8
  call void @gc_allocation_delete(ptr noundef %55)
  %56 = load ptr, ptr %4, align 8
  %57 = getelementptr inbounds %struct.AllocationMap, ptr %56, i32 0, i32 6
  %58 = load i64, ptr %57, align 8
  %59 = add i64 %58, -1
  store i64 %59, ptr %57, align 8
  br label %62

60:                                               ; preds = %27
  %61 = load ptr, ptr %8, align 8
  store ptr %61, ptr %9, align 8
  br label %62

62:                                               ; preds = %60, %54
  %63 = load ptr, ptr %10, align 8
  store ptr %63, ptr %8, align 8
  br label %24, !llvm.loop !9

64:                                               ; preds = %24
  %65 = load i8, ptr %6, align 1
  %66 = trunc i8 %65 to i1
  br i1 %66, label %67, label %70

67:                                               ; preds = %64
  %68 = load ptr, ptr %4, align 8
  %69 = call zeroext i1 @gc_allocation_map_resize_to_fit(ptr noundef %68)
  br label %70

70:                                               ; preds = %67, %64
  ret void
}

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define void @gc_free(ptr noundef %0, ptr noundef %1) #0 {
  %3 = alloca ptr, align 8
  %4 = alloca ptr, align 8
  %5 = alloca ptr, align 8
  store ptr %0, ptr %3, align 8
  store ptr %1, ptr %4, align 8
  %6 = load ptr, ptr %3, align 8
  %7 = getelementptr inbounds %struct.GarbageCollector, ptr %6, i32 0, i32 0
  %8 = load ptr, ptr %7, align 8
  %9 = load ptr, ptr %4, align 8
  %10 = call ptr @gc_allocation_map_get(ptr noundef %8, ptr noundef %9)
  store ptr %10, ptr %5, align 8
  %11 = load ptr, ptr %5, align 8
  %12 = icmp ne ptr %11, null
  br i1 %12, label %13, label %29

13:                                               ; preds = %2
  %14 = load ptr, ptr %5, align 8
  %15 = getelementptr inbounds %struct.Allocation, ptr %14, i32 0, i32 3
  %16 = load ptr, ptr %15, align 8
  %17 = icmp ne ptr %16, null
  br i1 %17, label %18, label %23

18:                                               ; preds = %13
  %19 = load ptr, ptr %5, align 8
  %20 = getelementptr inbounds %struct.Allocation, ptr %19, i32 0, i32 3
  %21 = load ptr, ptr %20, align 8
  %22 = load ptr, ptr %4, align 8
  call void %21(ptr noundef %22)
  br label %23

23:                                               ; preds = %18, %13
  %24 = load ptr, ptr %4, align 8
  call void @free(ptr noundef %24)
  %25 = load ptr, ptr %3, align 8
  %26 = getelementptr inbounds %struct.GarbageCollector, ptr %25, i32 0, i32 0
  %27 = load ptr, ptr %26, align 8
  %28 = load ptr, ptr %4, align 8
  call void @gc_allocation_map_remove(ptr noundef %27, ptr noundef %28, i1 noundef zeroext true)
  br label %36

29:                                               ; preds = %2
  br label %30

30:                                               ; preds = %29
  %31 = load ptr, ptr @__stderrp, align 8
  %32 = load ptr, ptr getelementptr inbounds ([0 x ptr], ptr @log_level_strings, i64 0, i64 1), align 8
  %33 = load ptr, ptr %4, align 8
  %34 = call i32 (ptr, ptr, ...) @fprintf(ptr noundef %31, ptr noundef @.str, ptr noundef %32, ptr noundef @__func__.gc_free, ptr noundef @.str.1, i32 noundef 466, ptr noundef %33)
  br label %35

35:                                               ; preds = %30
  br label %36

36:                                               ; preds = %35, %23
  ret void
}

declare void @free(ptr noundef) #1

declare i32 @fprintf(ptr noundef, ptr noundef, ...) #1

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define void @gc_start(ptr noundef %0, ptr noundef %1) #0 {
  %3 = alloca ptr, align 8
  %4 = alloca ptr, align 8
  store ptr %0, ptr %3, align 8
  store ptr %1, ptr %4, align 8
  %5 = load ptr, ptr %3, align 8
  %6 = load ptr, ptr %4, align 8
  call void @gc_start_ext(ptr noundef %5, ptr noundef %6, i64 noundef 1024, i64 noundef 1024, double noundef 2.000000e-01, double noundef 8.000000e-01, double noundef 5.000000e-01)
  ret void
}

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define void @gc_start_ext(ptr noundef %0, ptr noundef %1, i64 noundef %2, i64 noundef %3, double noundef %4, double noundef %5, double noundef %6) #0 {
  %8 = alloca ptr, align 8
  %9 = alloca ptr, align 8
  %10 = alloca i64, align 8
  %11 = alloca i64, align 8
  %12 = alloca double, align 8
  %13 = alloca double, align 8
  %14 = alloca double, align 8
  %15 = alloca double, align 8
  %16 = alloca double, align 8
  store ptr %0, ptr %8, align 8
  store ptr %1, ptr %9, align 8
  store i64 %2, ptr %10, align 8
  store i64 %3, ptr %11, align 8
  store double %4, ptr %12, align 8
  store double %5, ptr %13, align 8
  store double %6, ptr %14, align 8
  %17 = load double, ptr %12, align 8
  %18 = fcmp ogt double %17, 0.000000e+00
  br i1 %18, label %19, label %21

19:                                               ; preds = %7
  %20 = load double, ptr %12, align 8
  br label %22

21:                                               ; preds = %7
  br label %22

22:                                               ; preds = %21, %19
  %23 = phi double [ %20, %19 ], [ 2.000000e-01, %21 ]
  store double %23, ptr %15, align 8
  %24 = load double, ptr %13, align 8
  %25 = fcmp ogt double %24, 0.000000e+00
  br i1 %25, label %26, label %28

26:                                               ; preds = %22
  %27 = load double, ptr %13, align 8
  br label %29

28:                                               ; preds = %22
  br label %29

29:                                               ; preds = %28, %26
  %30 = phi double [ %27, %26 ], [ 8.000000e-01, %28 ]
  store double %30, ptr %16, align 8
  %31 = load double, ptr %14, align 8
  %32 = fcmp ogt double %31, 0.000000e+00
  br i1 %32, label %33, label %35

33:                                               ; preds = %29
  %34 = load double, ptr %14, align 8
  br label %36

35:                                               ; preds = %29
  br label %36

36:                                               ; preds = %35, %33
  %37 = phi double [ %34, %33 ], [ 5.000000e-01, %35 ]
  store double %37, ptr %14, align 8
  %38 = load ptr, ptr %8, align 8
  %39 = getelementptr inbounds %struct.GarbageCollector, ptr %38, i32 0, i32 1
  store i8 0, ptr %39, align 8
  %40 = load ptr, ptr %9, align 8
  %41 = load ptr, ptr %8, align 8
  %42 = getelementptr inbounds %struct.GarbageCollector, ptr %41, i32 0, i32 2
  store ptr %40, ptr %42, align 8
  %43 = load i64, ptr %10, align 8
  %44 = load i64, ptr %11, align 8
  %45 = icmp ult i64 %43, %44
  br i1 %45, label %46, label %48

46:                                               ; preds = %36
  %47 = load i64, ptr %11, align 8
  br label %50

48:                                               ; preds = %36
  %49 = load i64, ptr %10, align 8
  br label %50

50:                                               ; preds = %48, %46
  %51 = phi i64 [ %47, %46 ], [ %49, %48 ]
  store i64 %51, ptr %10, align 8
  %52 = load i64, ptr %11, align 8
  %53 = load i64, ptr %10, align 8
  %54 = load double, ptr %14, align 8
  %55 = load double, ptr %15, align 8
  %56 = load double, ptr %16, align 8
  %57 = call ptr @gc_allocation_map_new(i64 noundef %52, i64 noundef %53, double noundef %54, double noundef %55, double noundef %56)
  %58 = load ptr, ptr %8, align 8
  %59 = getelementptr inbounds %struct.GarbageCollector, ptr %58, i32 0, i32 0
  store ptr %57, ptr %59, align 8
  br label %60

60:                                               ; preds = %50
  br label %61

61:                                               ; preds = %60
  ret void
}

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define internal ptr @gc_allocation_map_new(i64 noundef %0, i64 noundef %1, double noundef %2, double noundef %3, double noundef %4) #0 {
  %6 = alloca i64, align 8
  %7 = alloca i64, align 8
  %8 = alloca double, align 8
  %9 = alloca double, align 8
  %10 = alloca double, align 8
  %11 = alloca ptr, align 8
  store i64 %0, ptr %6, align 8
  store i64 %1, ptr %7, align 8
  store double %2, ptr %8, align 8
  store double %3, ptr %9, align 8
  store double %4, ptr %10, align 8
  %12 = call ptr @malloc(i64 noundef 64) #11
  store ptr %12, ptr %11, align 8
  %13 = load i64, ptr %6, align 8
  %14 = call i64 @next_prime(i64 noundef %13)
  %15 = load ptr, ptr %11, align 8
  %16 = getelementptr inbounds %struct.AllocationMap, ptr %15, i32 0, i32 1
  store i64 %14, ptr %16, align 8
  %17 = load i64, ptr %7, align 8
  %18 = call i64 @next_prime(i64 noundef %17)
  %19 = load ptr, ptr %11, align 8
  %20 = getelementptr inbounds %struct.AllocationMap, ptr %19, i32 0, i32 0
  store i64 %18, ptr %20, align 8
  %21 = load ptr, ptr %11, align 8
  %22 = getelementptr inbounds %struct.AllocationMap, ptr %21, i32 0, i32 0
  %23 = load i64, ptr %22, align 8
  %24 = load ptr, ptr %11, align 8
  %25 = getelementptr inbounds %struct.AllocationMap, ptr %24, i32 0, i32 1
  %26 = load i64, ptr %25, align 8
  %27 = icmp ult i64 %23, %26
  br i1 %27, label %28, label %34

28:                                               ; preds = %5
  %29 = load ptr, ptr %11, align 8
  %30 = getelementptr inbounds %struct.AllocationMap, ptr %29, i32 0, i32 1
  %31 = load i64, ptr %30, align 8
  %32 = load ptr, ptr %11, align 8
  %33 = getelementptr inbounds %struct.AllocationMap, ptr %32, i32 0, i32 0
  store i64 %31, ptr %33, align 8
  br label %34

34:                                               ; preds = %28, %5
  %35 = load double, ptr %8, align 8
  %36 = load ptr, ptr %11, align 8
  %37 = getelementptr inbounds %struct.AllocationMap, ptr %36, i32 0, i32 4
  store double %35, ptr %37, align 8
  %38 = load double, ptr %8, align 8
  %39 = load ptr, ptr %11, align 8
  %40 = getelementptr inbounds %struct.AllocationMap, ptr %39, i32 0, i32 0
  %41 = load i64, ptr %40, align 8
  %42 = uitofp i64 %41 to double
  %43 = fmul double %38, %42
  %44 = fptosi double %43 to i32
  %45 = sext i32 %44 to i64
  %46 = load ptr, ptr %11, align 8
  %47 = getelementptr inbounds %struct.AllocationMap, ptr %46, i32 0, i32 5
  store i64 %45, ptr %47, align 8
  %48 = load double, ptr %9, align 8
  %49 = load ptr, ptr %11, align 8
  %50 = getelementptr inbounds %struct.AllocationMap, ptr %49, i32 0, i32 2
  store double %48, ptr %50, align 8
  %51 = load double, ptr %10, align 8
  %52 = load ptr, ptr %11, align 8
  %53 = getelementptr inbounds %struct.AllocationMap, ptr %52, i32 0, i32 3
  store double %51, ptr %53, align 8
  %54 = load ptr, ptr %11, align 8
  %55 = getelementptr inbounds %struct.AllocationMap, ptr %54, i32 0, i32 0
  %56 = load i64, ptr %55, align 8
  %57 = call ptr @calloc(i64 noundef %56, i64 noundef 8) #12
  %58 = load ptr, ptr %11, align 8
  %59 = getelementptr inbounds %struct.AllocationMap, ptr %58, i32 0, i32 7
  store ptr %57, ptr %59, align 8
  %60 = load ptr, ptr %11, align 8
  %61 = getelementptr inbounds %struct.AllocationMap, ptr %60, i32 0, i32 6
  store i64 0, ptr %61, align 8
  br label %62

62:                                               ; preds = %34
  br label %63

63:                                               ; preds = %62
  %64 = load ptr, ptr %11, align 8
  ret ptr %64
}

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define void @gc_pause(ptr noundef %0) #0 {
  %2 = alloca ptr, align 8
  store ptr %0, ptr %2, align 8
  %3 = load ptr, ptr %2, align 8
  %4 = getelementptr inbounds %struct.GarbageCollector, ptr %3, i32 0, i32 1
  store i8 1, ptr %4, align 8
  ret void
}

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define void @gc_resume(ptr noundef %0) #0 {
  %2 = alloca ptr, align 8
  store ptr %0, ptr %2, align 8
  %3 = load ptr, ptr %2, align 8
  %4 = getelementptr inbounds %struct.GarbageCollector, ptr %3, i32 0, i32 1
  store i8 0, ptr %4, align 8
  ret void
}

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define void @gc_mark_alloc(ptr noundef %0, ptr noundef %1) #0 {
  %3 = alloca ptr, align 8
  %4 = alloca ptr, align 8
  %5 = alloca ptr, align 8
  %6 = alloca ptr, align 8
  store ptr %0, ptr %3, align 8
  store ptr %1, ptr %4, align 8
  %7 = load ptr, ptr %3, align 8
  %8 = getelementptr inbounds %struct.GarbageCollector, ptr %7, i32 0, i32 0
  %9 = load ptr, ptr %8, align 8
  %10 = load ptr, ptr %4, align 8
  %11 = call ptr @gc_allocation_map_get(ptr noundef %9, ptr noundef %10)
  store ptr %11, ptr %5, align 8
  %12 = load ptr, ptr %5, align 8
  %13 = icmp ne ptr %12, null
  br i1 %13, label %14, label %56

14:                                               ; preds = %2
  %15 = load ptr, ptr %5, align 8
  %16 = getelementptr inbounds %struct.Allocation, ptr %15, i32 0, i32 2
  %17 = load i8, ptr %16, align 8
  %18 = sext i8 %17 to i32
  %19 = and i32 %18, 2
  %20 = icmp ne i32 %19, 0
  br i1 %20, label %56, label %21

21:                                               ; preds = %14
  br label %22

22:                                               ; preds = %21
  br label %23

23:                                               ; preds = %22
  %24 = load ptr, ptr %5, align 8
  %25 = getelementptr inbounds %struct.Allocation, ptr %24, i32 0, i32 2
  %26 = load i8, ptr %25, align 8
  %27 = sext i8 %26 to i32
  %28 = or i32 %27, 2
  %29 = trunc i32 %28 to i8
  store i8 %29, ptr %25, align 8
  br label %30

30:                                               ; preds = %23
  br label %31

31:                                               ; preds = %30
  %32 = load ptr, ptr %5, align 8
  %33 = getelementptr inbounds %struct.Allocation, ptr %32, i32 0, i32 0
  %34 = load ptr, ptr %33, align 8
  store ptr %34, ptr %6, align 8
  br label %35

35:                                               ; preds = %52, %31
  %36 = load ptr, ptr %6, align 8
  %37 = load ptr, ptr %5, align 8
  %38 = getelementptr inbounds %struct.Allocation, ptr %37, i32 0, i32 0
  %39 = load ptr, ptr %38, align 8
  %40 = load ptr, ptr %5, align 8
  %41 = getelementptr inbounds %struct.Allocation, ptr %40, i32 0, i32 1
  %42 = load i64, ptr %41, align 8
  %43 = getelementptr inbounds i8, ptr %39, i64 %42
  %44 = getelementptr inbounds i8, ptr %43, i64 -8
  %45 = icmp ule ptr %36, %44
  br i1 %45, label %46, label %55

46:                                               ; preds = %35
  br label %47

47:                                               ; preds = %46
  br label %48

48:                                               ; preds = %47
  %49 = load ptr, ptr %3, align 8
  %50 = load ptr, ptr %6, align 8
  %51 = load ptr, ptr %50, align 8
  call void @gc_mark_alloc(ptr noundef %49, ptr noundef %51)
  br label %52

52:                                               ; preds = %48
  %53 = load ptr, ptr %6, align 8
  %54 = getelementptr inbounds i8, ptr %53, i32 1
  store ptr %54, ptr %6, align 8
  br label %35, !llvm.loop !10

55:                                               ; preds = %35
  br label %56

56:                                               ; preds = %55, %14, %2
  ret void
}

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define void @gc_mark_stack(ptr noundef %0) #0 {
  %2 = alloca ptr, align 8
  %3 = alloca ptr, align 8
  %4 = alloca ptr, align 8
  %5 = alloca ptr, align 8
  store ptr %0, ptr %2, align 8
  br label %6

6:                                                ; preds = %1
  br label %7

7:                                                ; preds = %6
  %8 = call ptr @llvm.frameaddress.p0(i32 0)
  store ptr %8, ptr %3, align 8
  %9 = load ptr, ptr %2, align 8
  %10 = getelementptr inbounds %struct.GarbageCollector, ptr %9, i32 0, i32 2
  %11 = load ptr, ptr %10, align 8
  store ptr %11, ptr %4, align 8
  %12 = load ptr, ptr %3, align 8
  store ptr %12, ptr %5, align 8
  br label %13

13:                                               ; preds = %22, %7
  %14 = load ptr, ptr %5, align 8
  %15 = load ptr, ptr %4, align 8
  %16 = getelementptr inbounds i8, ptr %15, i64 -8
  %17 = icmp ule ptr %14, %16
  br i1 %17, label %18, label %25

18:                                               ; preds = %13
  %19 = load ptr, ptr %2, align 8
  %20 = load ptr, ptr %5, align 8
  %21 = load ptr, ptr %20, align 8
  call void @gc_mark_alloc(ptr noundef %19, ptr noundef %21)
  br label %22

22:                                               ; preds = %18
  %23 = load ptr, ptr %5, align 8
  %24 = getelementptr inbounds i8, ptr %23, i32 1
  store ptr %24, ptr %5, align 8
  br label %13, !llvm.loop !11

25:                                               ; preds = %13
  ret void
}

; Function Attrs: nocallback nofree nosync nounwind readnone willreturn
declare ptr @llvm.frameaddress.p0(i32 immarg) #3

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define void @gc_mark_roots(ptr noundef %0) #0 {
  %2 = alloca ptr, align 8
  %3 = alloca i64, align 8
  %4 = alloca ptr, align 8
  store ptr %0, ptr %2, align 8
  br label %5

5:                                                ; preds = %1
  br label %6

6:                                                ; preds = %5
  store i64 0, ptr %3, align 8
  br label %7

7:                                                ; preds = %46, %6
  %8 = load i64, ptr %3, align 8
  %9 = load ptr, ptr %2, align 8
  %10 = getelementptr inbounds %struct.GarbageCollector, ptr %9, i32 0, i32 0
  %11 = load ptr, ptr %10, align 8
  %12 = getelementptr inbounds %struct.AllocationMap, ptr %11, i32 0, i32 0
  %13 = load i64, ptr %12, align 8
  %14 = icmp ult i64 %8, %13
  br i1 %14, label %15, label %49

15:                                               ; preds = %7
  %16 = load ptr, ptr %2, align 8
  %17 = getelementptr inbounds %struct.GarbageCollector, ptr %16, i32 0, i32 0
  %18 = load ptr, ptr %17, align 8
  %19 = getelementptr inbounds %struct.AllocationMap, ptr %18, i32 0, i32 7
  %20 = load ptr, ptr %19, align 8
  %21 = load i64, ptr %3, align 8
  %22 = getelementptr inbounds ptr, ptr %20, i64 %21
  %23 = load ptr, ptr %22, align 8
  store ptr %23, ptr %4, align 8
  br label %24

24:                                               ; preds = %41, %15
  %25 = load ptr, ptr %4, align 8
  %26 = icmp ne ptr %25, null
  br i1 %26, label %27, label %45

27:                                               ; preds = %24
  %28 = load ptr, ptr %4, align 8
  %29 = getelementptr inbounds %struct.Allocation, ptr %28, i32 0, i32 2
  %30 = load i8, ptr %29, align 8
  %31 = sext i8 %30 to i32
  %32 = and i32 %31, 1
  %33 = icmp ne i32 %32, 0
  br i1 %33, label %34, label %41

34:                                               ; preds = %27
  br label %35

35:                                               ; preds = %34
  br label %36

36:                                               ; preds = %35
  %37 = load ptr, ptr %2, align 8
  %38 = load ptr, ptr %4, align 8
  %39 = getelementptr inbounds %struct.Allocation, ptr %38, i32 0, i32 0
  %40 = load ptr, ptr %39, align 8
  call void @gc_mark_alloc(ptr noundef %37, ptr noundef %40)
  br label %41

41:                                               ; preds = %36, %27
  %42 = load ptr, ptr %4, align 8
  %43 = getelementptr inbounds %struct.Allocation, ptr %42, i32 0, i32 4
  %44 = load ptr, ptr %43, align 8
  store ptr %44, ptr %4, align 8
  br label %24, !llvm.loop !12

45:                                               ; preds = %24
  br label %46

46:                                               ; preds = %45
  %47 = load i64, ptr %3, align 8
  %48 = add i64 %47, 1
  store i64 %48, ptr %3, align 8
  br label %7, !llvm.loop !13

49:                                               ; preds = %7
  ret void
}

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define void @gc_mark(ptr noundef %0) #0 {
  %2 = alloca ptr, align 8
  %3 = alloca ptr, align 8
  %4 = alloca [48 x i32], align 4
  store ptr %0, ptr %2, align 8
  br label %5

5:                                                ; preds = %1
  br label %6

6:                                                ; preds = %5
  %7 = load ptr, ptr %2, align 8
  call void @gc_mark_roots(ptr noundef %7)
  store volatile ptr @gc_mark_stack, ptr %3, align 8
  call void @llvm.memset.p0.i64(ptr align 4 %4, i8 0, i64 192, i1 false)
  %8 = getelementptr inbounds [48 x i32], ptr %4, i64 0, i64 0
  %9 = call i32 @setjmp(ptr noundef %8) #13
  %10 = load volatile ptr, ptr %3, align 8
  %11 = load ptr, ptr %2, align 8
  call void %10(ptr noundef %11)
  ret void
}

; Function Attrs: argmemonly nocallback nofree nounwind willreturn writeonly
declare void @llvm.memset.p0.i64(ptr nocapture writeonly, i8, i64, i1 immarg) #4

; Function Attrs: returns_twice
declare i32 @setjmp(ptr noundef) #5

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define i64 @gc_sweep(ptr noundef %0) #0 {
  %2 = alloca ptr, align 8
  %3 = alloca i64, align 8
  %4 = alloca i64, align 8
  %5 = alloca ptr, align 8
  %6 = alloca ptr, align 8
  store ptr %0, ptr %2, align 8
  br label %7

7:                                                ; preds = %1
  br label %8

8:                                                ; preds = %7
  store i64 0, ptr %3, align 8
  store i64 0, ptr %4, align 8
  br label %9

9:                                                ; preds = %83, %8
  %10 = load i64, ptr %4, align 8
  %11 = load ptr, ptr %2, align 8
  %12 = getelementptr inbounds %struct.GarbageCollector, ptr %11, i32 0, i32 0
  %13 = load ptr, ptr %12, align 8
  %14 = getelementptr inbounds %struct.AllocationMap, ptr %13, i32 0, i32 0
  %15 = load i64, ptr %14, align 8
  %16 = icmp ult i64 %10, %15
  br i1 %16, label %17, label %86

17:                                               ; preds = %9
  %18 = load ptr, ptr %2, align 8
  %19 = getelementptr inbounds %struct.GarbageCollector, ptr %18, i32 0, i32 0
  %20 = load ptr, ptr %19, align 8
  %21 = getelementptr inbounds %struct.AllocationMap, ptr %20, i32 0, i32 7
  %22 = load ptr, ptr %21, align 8
  %23 = load i64, ptr %4, align 8
  %24 = getelementptr inbounds ptr, ptr %22, i64 %23
  %25 = load ptr, ptr %24, align 8
  store ptr %25, ptr %5, align 8
  store ptr null, ptr %6, align 8
  br label %26

26:                                               ; preds = %81, %17
  %27 = load ptr, ptr %5, align 8
  %28 = icmp ne ptr %27, null
  br i1 %28, label %29, label %82

29:                                               ; preds = %26
  %30 = load ptr, ptr %5, align 8
  %31 = getelementptr inbounds %struct.Allocation, ptr %30, i32 0, i32 2
  %32 = load i8, ptr %31, align 8
  %33 = sext i8 %32 to i32
  %34 = and i32 %33, 2
  %35 = icmp ne i32 %34, 0
  br i1 %35, label %36, label %48

36:                                               ; preds = %29
  br label %37

37:                                               ; preds = %36
  br label %38

38:                                               ; preds = %37
  %39 = load ptr, ptr %5, align 8
  %40 = getelementptr inbounds %struct.Allocation, ptr %39, i32 0, i32 2
  %41 = load i8, ptr %40, align 8
  %42 = sext i8 %41 to i32
  %43 = and i32 %42, -3
  %44 = trunc i32 %43 to i8
  store i8 %44, ptr %40, align 8
  %45 = load ptr, ptr %5, align 8
  %46 = getelementptr inbounds %struct.Allocation, ptr %45, i32 0, i32 4
  %47 = load ptr, ptr %46, align 8
  store ptr %47, ptr %5, align 8
  br label %81

48:                                               ; preds = %29
  br label %49

49:                                               ; preds = %48
  br label %50

50:                                               ; preds = %49
  %51 = load ptr, ptr %5, align 8
  %52 = getelementptr inbounds %struct.Allocation, ptr %51, i32 0, i32 1
  %53 = load i64, ptr %52, align 8
  %54 = load i64, ptr %3, align 8
  %55 = add i64 %54, %53
  store i64 %55, ptr %3, align 8
  %56 = load ptr, ptr %5, align 8
  %57 = getelementptr inbounds %struct.Allocation, ptr %56, i32 0, i32 3
  %58 = load ptr, ptr %57, align 8
  %59 = icmp ne ptr %58, null
  br i1 %59, label %60, label %67

60:                                               ; preds = %50
  %61 = load ptr, ptr %5, align 8
  %62 = getelementptr inbounds %struct.Allocation, ptr %61, i32 0, i32 3
  %63 = load ptr, ptr %62, align 8
  %64 = load ptr, ptr %5, align 8
  %65 = getelementptr inbounds %struct.Allocation, ptr %64, i32 0, i32 0
  %66 = load ptr, ptr %65, align 8
  call void %63(ptr noundef %66)
  br label %67

67:                                               ; preds = %60, %50
  %68 = load ptr, ptr %5, align 8
  %69 = getelementptr inbounds %struct.Allocation, ptr %68, i32 0, i32 0
  %70 = load ptr, ptr %69, align 8
  call void @free(ptr noundef %70)
  %71 = load ptr, ptr %5, align 8
  %72 = getelementptr inbounds %struct.Allocation, ptr %71, i32 0, i32 4
  %73 = load ptr, ptr %72, align 8
  store ptr %73, ptr %6, align 8
  %74 = load ptr, ptr %2, align 8
  %75 = getelementptr inbounds %struct.GarbageCollector, ptr %74, i32 0, i32 0
  %76 = load ptr, ptr %75, align 8
  %77 = load ptr, ptr %5, align 8
  %78 = getelementptr inbounds %struct.Allocation, ptr %77, i32 0, i32 0
  %79 = load ptr, ptr %78, align 8
  call void @gc_allocation_map_remove(ptr noundef %76, ptr noundef %79, i1 noundef zeroext false)
  %80 = load ptr, ptr %6, align 8
  store ptr %80, ptr %5, align 8
  br label %81

81:                                               ; preds = %67, %38
  br label %26, !llvm.loop !14

82:                                               ; preds = %26
  br label %83

83:                                               ; preds = %82
  %84 = load i64, ptr %4, align 8
  %85 = add i64 %84, 1
  store i64 %85, ptr %4, align 8
  br label %9, !llvm.loop !15

86:                                               ; preds = %9
  %87 = load ptr, ptr %2, align 8
  %88 = getelementptr inbounds %struct.GarbageCollector, ptr %87, i32 0, i32 0
  %89 = load ptr, ptr %88, align 8
  %90 = call zeroext i1 @gc_allocation_map_resize_to_fit(ptr noundef %89)
  %91 = load i64, ptr %3, align 8
  ret i64 %91
}

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define internal zeroext i1 @gc_allocation_map_resize_to_fit(ptr noundef %0) #0 {
  %2 = alloca i1, align 1
  %3 = alloca ptr, align 8
  %4 = alloca double, align 8
  store ptr %0, ptr %3, align 8
  %5 = load ptr, ptr %3, align 8
  %6 = call double @gc_allocation_map_load_factor(ptr noundef %5)
  store double %6, ptr %4, align 8
  %7 = load double, ptr %4, align 8
  %8 = load ptr, ptr %3, align 8
  %9 = getelementptr inbounds %struct.AllocationMap, ptr %8, i32 0, i32 3
  %10 = load double, ptr %9, align 8
  %11 = fcmp ogt double %7, %10
  br i1 %11, label %12, label %21

12:                                               ; preds = %1
  br label %13

13:                                               ; preds = %12
  br label %14

14:                                               ; preds = %13
  %15 = load ptr, ptr %3, align 8
  %16 = load ptr, ptr %3, align 8
  %17 = getelementptr inbounds %struct.AllocationMap, ptr %16, i32 0, i32 0
  %18 = load i64, ptr %17, align 8
  %19 = mul i64 %18, 2
  %20 = call i64 @next_prime(i64 noundef %19)
  call void @gc_allocation_map_resize(ptr noundef %15, i64 noundef %20)
  store i1 true, ptr %2, align 1
  br label %37

21:                                               ; preds = %1
  %22 = load double, ptr %4, align 8
  %23 = load ptr, ptr %3, align 8
  %24 = getelementptr inbounds %struct.AllocationMap, ptr %23, i32 0, i32 2
  %25 = load double, ptr %24, align 8
  %26 = fcmp olt double %22, %25
  br i1 %26, label %27, label %36

27:                                               ; preds = %21
  br label %28

28:                                               ; preds = %27
  br label %29

29:                                               ; preds = %28
  %30 = load ptr, ptr %3, align 8
  %31 = load ptr, ptr %3, align 8
  %32 = getelementptr inbounds %struct.AllocationMap, ptr %31, i32 0, i32 0
  %33 = load i64, ptr %32, align 8
  %34 = udiv i64 %33, 2
  %35 = call i64 @next_prime(i64 noundef %34)
  call void @gc_allocation_map_resize(ptr noundef %30, i64 noundef %35)
  store i1 true, ptr %2, align 1
  br label %37

36:                                               ; preds = %21
  store i1 false, ptr %2, align 1
  br label %37

37:                                               ; preds = %36, %29, %14
  %38 = load i1, ptr %2, align 1
  ret i1 %38
}

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define void @gc_unroot_roots(ptr noundef %0) #0 {
  %2 = alloca ptr, align 8
  %3 = alloca i64, align 8
  %4 = alloca ptr, align 8
  store ptr %0, ptr %2, align 8
  br label %5

5:                                                ; preds = %1
  br label %6

6:                                                ; preds = %5
  store i64 0, ptr %3, align 8
  br label %7

7:                                                ; preds = %46, %6
  %8 = load i64, ptr %3, align 8
  %9 = load ptr, ptr %2, align 8
  %10 = getelementptr inbounds %struct.GarbageCollector, ptr %9, i32 0, i32 0
  %11 = load ptr, ptr %10, align 8
  %12 = getelementptr inbounds %struct.AllocationMap, ptr %11, i32 0, i32 0
  %13 = load i64, ptr %12, align 8
  %14 = icmp ult i64 %8, %13
  br i1 %14, label %15, label %49

15:                                               ; preds = %7
  %16 = load ptr, ptr %2, align 8
  %17 = getelementptr inbounds %struct.GarbageCollector, ptr %16, i32 0, i32 0
  %18 = load ptr, ptr %17, align 8
  %19 = getelementptr inbounds %struct.AllocationMap, ptr %18, i32 0, i32 7
  %20 = load ptr, ptr %19, align 8
  %21 = load i64, ptr %3, align 8
  %22 = getelementptr inbounds ptr, ptr %20, i64 %21
  %23 = load ptr, ptr %22, align 8
  store ptr %23, ptr %4, align 8
  br label %24

24:                                               ; preds = %41, %15
  %25 = load ptr, ptr %4, align 8
  %26 = icmp ne ptr %25, null
  br i1 %26, label %27, label %45

27:                                               ; preds = %24
  %28 = load ptr, ptr %4, align 8
  %29 = getelementptr inbounds %struct.Allocation, ptr %28, i32 0, i32 2
  %30 = load i8, ptr %29, align 8
  %31 = sext i8 %30 to i32
  %32 = and i32 %31, 1
  %33 = icmp ne i32 %32, 0
  br i1 %33, label %34, label %41

34:                                               ; preds = %27
  %35 = load ptr, ptr %4, align 8
  %36 = getelementptr inbounds %struct.Allocation, ptr %35, i32 0, i32 2
  %37 = load i8, ptr %36, align 8
  %38 = sext i8 %37 to i32
  %39 = and i32 %38, -2
  %40 = trunc i32 %39 to i8
  store i8 %40, ptr %36, align 8
  br label %41

41:                                               ; preds = %34, %27
  %42 = load ptr, ptr %4, align 8
  %43 = getelementptr inbounds %struct.Allocation, ptr %42, i32 0, i32 4
  %44 = load ptr, ptr %43, align 8
  store ptr %44, ptr %4, align 8
  br label %24, !llvm.loop !16

45:                                               ; preds = %24
  br label %46

46:                                               ; preds = %45
  %47 = load i64, ptr %3, align 8
  %48 = add i64 %47, 1
  store i64 %48, ptr %3, align 8
  br label %7, !llvm.loop !17

49:                                               ; preds = %7
  ret void
}

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define i64 @gc_stop(ptr noundef %0) #0 {
  %2 = alloca ptr, align 8
  %3 = alloca i64, align 8
  store ptr %0, ptr %2, align 8
  %4 = load ptr, ptr %2, align 8
  call void @gc_unroot_roots(ptr noundef %4)
  %5 = load ptr, ptr %2, align 8
  %6 = call i64 @gc_sweep(ptr noundef %5)
  store i64 %6, ptr %3, align 8
  %7 = load ptr, ptr %2, align 8
  %8 = getelementptr inbounds %struct.GarbageCollector, ptr %7, i32 0, i32 0
  %9 = load ptr, ptr %8, align 8
  call void @gc_allocation_map_delete(ptr noundef %9)
  %10 = load i64, ptr %3, align 8
  ret i64 %10
}

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define internal void @gc_allocation_map_delete(ptr noundef %0) #0 {
  %2 = alloca ptr, align 8
  %3 = alloca ptr, align 8
  %4 = alloca ptr, align 8
  %5 = alloca i64, align 8
  store ptr %0, ptr %2, align 8
  br label %6

6:                                                ; preds = %1
  br label %7

7:                                                ; preds = %6
  store i64 0, ptr %5, align 8
  br label %8

8:                                                ; preds = %34, %7
  %9 = load i64, ptr %5, align 8
  %10 = load ptr, ptr %2, align 8
  %11 = getelementptr inbounds %struct.AllocationMap, ptr %10, i32 0, i32 0
  %12 = load i64, ptr %11, align 8
  %13 = icmp ult i64 %9, %12
  br i1 %13, label %14, label %37

14:                                               ; preds = %8
  %15 = load ptr, ptr %2, align 8
  %16 = getelementptr inbounds %struct.AllocationMap, ptr %15, i32 0, i32 7
  %17 = load ptr, ptr %16, align 8
  %18 = load i64, ptr %5, align 8
  %19 = getelementptr inbounds ptr, ptr %17, i64 %18
  %20 = load ptr, ptr %19, align 8
  store ptr %20, ptr %3, align 8
  %21 = icmp ne ptr %20, null
  br i1 %21, label %22, label %33

22:                                               ; preds = %14
  br label %23

23:                                               ; preds = %26, %22
  %24 = load ptr, ptr %3, align 8
  %25 = icmp ne ptr %24, null
  br i1 %25, label %26, label %32

26:                                               ; preds = %23
  %27 = load ptr, ptr %3, align 8
  store ptr %27, ptr %4, align 8
  %28 = load ptr, ptr %3, align 8
  %29 = getelementptr inbounds %struct.Allocation, ptr %28, i32 0, i32 4
  %30 = load ptr, ptr %29, align 8
  store ptr %30, ptr %3, align 8
  %31 = load ptr, ptr %4, align 8
  call void @gc_allocation_delete(ptr noundef %31)
  br label %23, !llvm.loop !18

32:                                               ; preds = %23
  br label %33

33:                                               ; preds = %32, %14
  br label %34

34:                                               ; preds = %33
  %35 = load i64, ptr %5, align 8
  %36 = add i64 %35, 1
  store i64 %36, ptr %5, align 8
  br label %8, !llvm.loop !19

37:                                               ; preds = %8
  %38 = load ptr, ptr %2, align 8
  %39 = getelementptr inbounds %struct.AllocationMap, ptr %38, i32 0, i32 7
  %40 = load ptr, ptr %39, align 8
  call void @free(ptr noundef %40)
  %41 = load ptr, ptr %2, align 8
  call void @free(ptr noundef %41)
  ret void
}

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define i64 @gc_run(ptr noundef %0) #0 {
  %2 = alloca ptr, align 8
  store ptr %0, ptr %2, align 8
  br label %3

3:                                                ; preds = %1
  br label %4

4:                                                ; preds = %3
  %5 = load ptr, ptr %2, align 8
  call void @gc_mark(ptr noundef %5)
  %6 = load ptr, ptr %2, align 8
  %7 = call i64 @gc_sweep(ptr noundef %6)
  ret i64 %7
}

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define ptr @gc_strdup(ptr noundef %0, ptr noundef %1) #0 {
  %3 = alloca ptr, align 8
  %4 = alloca ptr, align 8
  %5 = alloca ptr, align 8
  %6 = alloca i64, align 8
  %7 = alloca ptr, align 8
  store ptr %0, ptr %4, align 8
  store ptr %1, ptr %5, align 8
  %8 = load ptr, ptr %5, align 8
  %9 = call i64 @strlen(ptr noundef %8)
  %10 = add i64 %9, 1
  store i64 %10, ptr %6, align 8
  %11 = load ptr, ptr %4, align 8
  %12 = load i64, ptr %6, align 8
  %13 = call ptr @gc_malloc(ptr noundef %11, i64 noundef %12)
  store ptr %13, ptr %7, align 8
  %14 = load ptr, ptr %7, align 8
  %15 = icmp eq ptr %14, null
  br i1 %15, label %16, label %17

16:                                               ; preds = %2
  store ptr null, ptr %3, align 8
  br label %24

17:                                               ; preds = %2
  %18 = load ptr, ptr %7, align 8
  %19 = load ptr, ptr %5, align 8
  %20 = load i64, ptr %6, align 8
  %21 = load ptr, ptr %7, align 8
  %22 = call i64 @llvm.objectsize.i64.p0(ptr %21, i1 false, i1 true, i1 false)
  %23 = call ptr @__memcpy_chk(ptr noundef %18, ptr noundef %19, i64 noundef %20, i64 noundef %22) #14
  store ptr %23, ptr %3, align 8
  br label %24

24:                                               ; preds = %17, %16
  %25 = load ptr, ptr %3, align 8
  ret ptr %25
}

declare i64 @strlen(ptr noundef) #1

; Function Attrs: nounwind
declare ptr @__memcpy_chk(ptr noundef, ptr noundef, i64 noundef, i64 noundef) #6

; Function Attrs: nocallback nofree nosync nounwind readnone speculatable willreturn
declare i64 @llvm.objectsize.i64.p0(ptr, i1 immarg, i1 immarg, i1 immarg) #7

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define internal zeroext i1 @gc_needs_sweep(ptr noundef %0) #0 {
  %2 = alloca ptr, align 8
  store ptr %0, ptr %2, align 8
  %3 = load ptr, ptr %2, align 8
  %4 = getelementptr inbounds %struct.GarbageCollector, ptr %3, i32 0, i32 0
  %5 = load ptr, ptr %4, align 8
  %6 = getelementptr inbounds %struct.AllocationMap, ptr %5, i32 0, i32 6
  %7 = load i64, ptr %6, align 8
  %8 = load ptr, ptr %2, align 8
  %9 = getelementptr inbounds %struct.GarbageCollector, ptr %8, i32 0, i32 0
  %10 = load ptr, ptr %9, align 8
  %11 = getelementptr inbounds %struct.AllocationMap, ptr %10, i32 0, i32 5
  %12 = load i64, ptr %11, align 8
  %13 = icmp ugt i64 %7, %12
  ret i1 %13
}

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define internal ptr @gc_mcalloc(i64 noundef %0, i64 noundef %1) #0 {
  %3 = alloca ptr, align 8
  %4 = alloca i64, align 8
  %5 = alloca i64, align 8
  store i64 %0, ptr %4, align 8
  store i64 %1, ptr %5, align 8
  %6 = load i64, ptr %4, align 8
  %7 = icmp ne i64 %6, 0
  br i1 %7, label %11, label %8

8:                                                ; preds = %2
  %9 = load i64, ptr %5, align 8
  %10 = call ptr @malloc(i64 noundef %9) #11
  store ptr %10, ptr %3, align 8
  br label %15

11:                                               ; preds = %2
  %12 = load i64, ptr %4, align 8
  %13 = load i64, ptr %5, align 8
  %14 = call ptr @calloc(i64 noundef %12, i64 noundef %13) #12
  store ptr %14, ptr %3, align 8
  br label %15

15:                                               ; preds = %11, %8
  %16 = load ptr, ptr %3, align 8
  ret ptr %16
}

; Function Attrs: allocsize(0)
declare ptr @malloc(i64 noundef) #8

; Function Attrs: allocsize(0,1)
declare ptr @calloc(i64 noundef, i64 noundef) #9

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define internal i64 @gc_hash(ptr noundef %0) #0 {
  %2 = alloca ptr, align 8
  store ptr %0, ptr %2, align 8
  %3 = load ptr, ptr %2, align 8
  %4 = ptrtoint ptr %3 to i64
  %5 = lshr i64 %4, 3
  ret i64 %5
}

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define internal ptr @gc_allocation_new(ptr noundef %0, i64 noundef %1, ptr noundef %2) #0 {
  %4 = alloca ptr, align 8
  %5 = alloca i64, align 8
  %6 = alloca ptr, align 8
  %7 = alloca ptr, align 8
  store ptr %0, ptr %4, align 8
  store i64 %1, ptr %5, align 8
  store ptr %2, ptr %6, align 8
  %8 = call ptr @malloc(i64 noundef 40) #11
  store ptr %8, ptr %7, align 8
  %9 = load ptr, ptr %4, align 8
  %10 = load ptr, ptr %7, align 8
  %11 = getelementptr inbounds %struct.Allocation, ptr %10, i32 0, i32 0
  store ptr %9, ptr %11, align 8
  %12 = load i64, ptr %5, align 8
  %13 = load ptr, ptr %7, align 8
  %14 = getelementptr inbounds %struct.Allocation, ptr %13, i32 0, i32 1
  store i64 %12, ptr %14, align 8
  %15 = load ptr, ptr %7, align 8
  %16 = getelementptr inbounds %struct.Allocation, ptr %15, i32 0, i32 2
  store i8 0, ptr %16, align 8
  %17 = load ptr, ptr %6, align 8
  %18 = load ptr, ptr %7, align 8
  %19 = getelementptr inbounds %struct.Allocation, ptr %18, i32 0, i32 3
  store ptr %17, ptr %19, align 8
  %20 = load ptr, ptr %7, align 8
  %21 = getelementptr inbounds %struct.Allocation, ptr %20, i32 0, i32 4
  store ptr null, ptr %21, align 8
  %22 = load ptr, ptr %7, align 8
  ret ptr %22
}

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define internal void @gc_allocation_delete(ptr noundef %0) #0 {
  %2 = alloca ptr, align 8
  store ptr %0, ptr %2, align 8
  %3 = load ptr, ptr %2, align 8
  call void @free(ptr noundef %3)
  ret void
}

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define internal i64 @next_prime(i64 noundef %0) #0 {
  %2 = alloca i64, align 8
  store i64 %0, ptr %2, align 8
  br label %3

3:                                                ; preds = %7, %1
  %4 = load i64, ptr %2, align 8
  %5 = call zeroext i1 @is_prime(i64 noundef %4)
  %6 = xor i1 %5, true
  br i1 %6, label %7, label %10

7:                                                ; preds = %3
  %8 = load i64, ptr %2, align 8
  %9 = add i64 %8, 1
  store i64 %9, ptr %2, align 8
  br label %3, !llvm.loop !20

10:                                               ; preds = %3
  %11 = load i64, ptr %2, align 8
  ret i64 %11
}

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define internal zeroext i1 @is_prime(i64 noundef %0) #0 {
  %2 = alloca i1, align 1
  %3 = alloca i64, align 8
  %4 = alloca i64, align 8
  store i64 %0, ptr %3, align 8
  %5 = load i64, ptr %3, align 8
  %6 = icmp ule i64 %5, 3
  br i1 %6, label %7, label %10

7:                                                ; preds = %1
  %8 = load i64, ptr %3, align 8
  %9 = icmp ugt i64 %8, 1
  store i1 %9, ptr %2, align 1
  br label %43

10:                                               ; preds = %1
  %11 = load i64, ptr %3, align 8
  %12 = urem i64 %11, 2
  %13 = icmp eq i64 %12, 0
  br i1 %13, label %18, label %14

14:                                               ; preds = %10
  %15 = load i64, ptr %3, align 8
  %16 = urem i64 %15, 3
  %17 = icmp eq i64 %16, 0
  br i1 %17, label %18, label %19

18:                                               ; preds = %14, %10
  store i1 false, ptr %2, align 1
  br label %43

19:                                               ; preds = %14
  store i64 5, ptr %4, align 8
  br label %20

20:                                               ; preds = %39, %19
  %21 = load i64, ptr %4, align 8
  %22 = load i64, ptr %4, align 8
  %23 = mul i64 %21, %22
  %24 = load i64, ptr %3, align 8
  %25 = icmp ule i64 %23, %24
  br i1 %25, label %26, label %42

26:                                               ; preds = %20
  %27 = load i64, ptr %3, align 8
  %28 = load i64, ptr %4, align 8
  %29 = urem i64 %27, %28
  %30 = icmp eq i64 %29, 0
  br i1 %30, label %37, label %31

31:                                               ; preds = %26
  %32 = load i64, ptr %3, align 8
  %33 = load i64, ptr %4, align 8
  %34 = add i64 %33, 2
  %35 = urem i64 %32, %34
  %36 = icmp eq i64 %35, 0
  br i1 %36, label %37, label %38

37:                                               ; preds = %31, %26
  store i1 false, ptr %2, align 1
  br label %43

38:                                               ; preds = %31
  br label %39

39:                                               ; preds = %38
  %40 = load i64, ptr %4, align 8
  %41 = add i64 %40, 6
  store i64 %41, ptr %4, align 8
  br label %20, !llvm.loop !21

42:                                               ; preds = %20
  store i1 true, ptr %2, align 1
  br label %43

43:                                               ; preds = %42, %37, %18, %7
  %44 = load i1, ptr %2, align 1
  ret i1 %44
}

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define internal double @gc_allocation_map_load_factor(ptr noundef %0) #0 {
  %2 = alloca ptr, align 8
  store ptr %0, ptr %2, align 8
  %3 = load ptr, ptr %2, align 8
  %4 = getelementptr inbounds %struct.AllocationMap, ptr %3, i32 0, i32 6
  %5 = load i64, ptr %4, align 8
  %6 = uitofp i64 %5 to double
  %7 = load ptr, ptr %2, align 8
  %8 = getelementptr inbounds %struct.AllocationMap, ptr %7, i32 0, i32 0
  %9 = load i64, ptr %8, align 8
  %10 = uitofp i64 %9 to double
  %11 = fdiv double %6, %10
  ret double %11
}

; Function Attrs: noinline nounwind optnone ssp uwtable(sync)
define internal void @gc_allocation_map_resize(ptr noundef %0, i64 noundef %1) #0 {
  %3 = alloca ptr, align 8
  %4 = alloca i64, align 8
  %5 = alloca ptr, align 8
  %6 = alloca i64, align 8
  %7 = alloca ptr, align 8
  %8 = alloca ptr, align 8
  %9 = alloca i64, align 8
  store ptr %0, ptr %3, align 8
  store i64 %1, ptr %4, align 8
  %10 = load i64, ptr %4, align 8
  %11 = load ptr, ptr %3, align 8
  %12 = getelementptr inbounds %struct.AllocationMap, ptr %11, i32 0, i32 1
  %13 = load i64, ptr %12, align 8
  %14 = icmp ule i64 %10, %13
  br i1 %14, label %15, label %16

15:                                               ; preds = %2
  br label %91

16:                                               ; preds = %2
  br label %17

17:                                               ; preds = %16
  br label %18

18:                                               ; preds = %17
  %19 = load i64, ptr %4, align 8
  %20 = call ptr @calloc(i64 noundef %19, i64 noundef 8) #12
  store ptr %20, ptr %5, align 8
  store i64 0, ptr %6, align 8
  br label %21

21:                                               ; preds = %59, %18
  %22 = load i64, ptr %6, align 8
  %23 = load ptr, ptr %3, align 8
  %24 = getelementptr inbounds %struct.AllocationMap, ptr %23, i32 0, i32 0
  %25 = load i64, ptr %24, align 8
  %26 = icmp ult i64 %22, %25
  br i1 %26, label %27, label %62

27:                                               ; preds = %21
  %28 = load ptr, ptr %3, align 8
  %29 = getelementptr inbounds %struct.AllocationMap, ptr %28, i32 0, i32 7
  %30 = load ptr, ptr %29, align 8
  %31 = load i64, ptr %6, align 8
  %32 = getelementptr inbounds ptr, ptr %30, i64 %31
  %33 = load ptr, ptr %32, align 8
  store ptr %33, ptr %7, align 8
  br label %34

34:                                               ; preds = %37, %27
  %35 = load ptr, ptr %7, align 8
  %36 = icmp ne ptr %35, null
  br i1 %36, label %37, label %58

37:                                               ; preds = %34
  %38 = load ptr, ptr %7, align 8
  %39 = getelementptr inbounds %struct.Allocation, ptr %38, i32 0, i32 4
  %40 = load ptr, ptr %39, align 8
  store ptr %40, ptr %8, align 8
  %41 = load ptr, ptr %7, align 8
  %42 = getelementptr inbounds %struct.Allocation, ptr %41, i32 0, i32 0
  %43 = load ptr, ptr %42, align 8
  %44 = call i64 @gc_hash(ptr noundef %43)
  %45 = load i64, ptr %4, align 8
  %46 = urem i64 %44, %45
  store i64 %46, ptr %9, align 8
  %47 = load ptr, ptr %5, align 8
  %48 = load i64, ptr %9, align 8
  %49 = getelementptr inbounds ptr, ptr %47, i64 %48
  %50 = load ptr, ptr %49, align 8
  %51 = load ptr, ptr %7, align 8
  %52 = getelementptr inbounds %struct.Allocation, ptr %51, i32 0, i32 4
  store ptr %50, ptr %52, align 8
  %53 = load ptr, ptr %7, align 8
  %54 = load ptr, ptr %5, align 8
  %55 = load i64, ptr %9, align 8
  %56 = getelementptr inbounds ptr, ptr %54, i64 %55
  store ptr %53, ptr %56, align 8
  %57 = load ptr, ptr %8, align 8
  store ptr %57, ptr %7, align 8
  br label %34, !llvm.loop !22

58:                                               ; preds = %34
  br label %59

59:                                               ; preds = %58
  %60 = load i64, ptr %6, align 8
  %61 = add i64 %60, 1
  store i64 %61, ptr %6, align 8
  br label %21, !llvm.loop !23

62:                                               ; preds = %21
  %63 = load ptr, ptr %3, align 8
  %64 = getelementptr inbounds %struct.AllocationMap, ptr %63, i32 0, i32 7
  %65 = load ptr, ptr %64, align 8
  call void @free(ptr noundef %65)
  %66 = load i64, ptr %4, align 8
  %67 = load ptr, ptr %3, align 8
  %68 = getelementptr inbounds %struct.AllocationMap, ptr %67, i32 0, i32 0
  store i64 %66, ptr %68, align 8
  %69 = load ptr, ptr %5, align 8
  %70 = load ptr, ptr %3, align 8
  %71 = getelementptr inbounds %struct.AllocationMap, ptr %70, i32 0, i32 7
  store ptr %69, ptr %71, align 8
  %72 = load ptr, ptr %3, align 8
  %73 = getelementptr inbounds %struct.AllocationMap, ptr %72, i32 0, i32 6
  %74 = load i64, ptr %73, align 8
  %75 = uitofp i64 %74 to double
  %76 = load ptr, ptr %3, align 8
  %77 = getelementptr inbounds %struct.AllocationMap, ptr %76, i32 0, i32 4
  %78 = load double, ptr %77, align 8
  %79 = load ptr, ptr %3, align 8
  %80 = getelementptr inbounds %struct.AllocationMap, ptr %79, i32 0, i32 0
  %81 = load i64, ptr %80, align 8
  %82 = load ptr, ptr %3, align 8
  %83 = getelementptr inbounds %struct.AllocationMap, ptr %82, i32 0, i32 6
  %84 = load i64, ptr %83, align 8
  %85 = sub i64 %81, %84
  %86 = uitofp i64 %85 to double
  %87 = call double @llvm.fmuladd.f64(double %78, double %86, double %75)
  %88 = fptoui double %87 to i64
  %89 = load ptr, ptr %3, align 8
  %90 = getelementptr inbounds %struct.AllocationMap, ptr %89, i32 0, i32 5
  store i64 %88, ptr %90, align 8
  br label %91

91:                                               ; preds = %62, %15
  ret void
}

; Function Attrs: nocallback nofree nosync nounwind readnone speculatable willreturn
declare double @llvm.fmuladd.f64(double, double, double) #7

attributes #0 = { noinline nounwind optnone ssp uwtable(sync) "frame-pointer"="non-leaf" "min-legal-vector-width"="0" "no-trapping-math"="true" "probe-stack"="__chkstk_darwin" "stack-protector-buffer-size"="8" "target-cpu"="apple-m1" "target-features"="+aes,+crc,+crypto,+dotprod,+fp-armv8,+fp16fml,+fullfp16,+lse,+neon,+ras,+rcpc,+rdm,+sha2,+sha3,+sm4,+v8.1a,+v8.2a,+v8.3a,+v8.4a,+v8.5a,+v8a,+zcm,+zcz" }
attributes #1 = { "frame-pointer"="non-leaf" "no-trapping-math"="true" "probe-stack"="__chkstk_darwin" "stack-protector-buffer-size"="8" "target-cpu"="apple-m1" "target-features"="+aes,+crc,+crypto,+dotprod,+fp-armv8,+fp16fml,+fullfp16,+lse,+neon,+ras,+rcpc,+rdm,+sha2,+sha3,+sm4,+v8.1a,+v8.2a,+v8.3a,+v8.4a,+v8.5a,+v8a,+zcm,+zcz" }
attributes #2 = { allocsize(1) "frame-pointer"="non-leaf" "no-trapping-math"="true" "probe-stack"="__chkstk_darwin" "stack-protector-buffer-size"="8" "target-cpu"="apple-m1" "target-features"="+aes,+crc,+crypto,+dotprod,+fp-armv8,+fp16fml,+fullfp16,+lse,+neon,+ras,+rcpc,+rdm,+sha2,+sha3,+sm4,+v8.1a,+v8.2a,+v8.3a,+v8.4a,+v8.5a,+v8a,+zcm,+zcz" }
attributes #3 = { nocallback nofree nosync nounwind readnone willreturn }
attributes #4 = { argmemonly nocallback nofree nounwind willreturn writeonly }
attributes #5 = { returns_twice "frame-pointer"="non-leaf" "no-trapping-math"="true" "probe-stack"="__chkstk_darwin" "stack-protector-buffer-size"="8" "target-cpu"="apple-m1" "target-features"="+aes,+crc,+crypto,+dotprod,+fp-armv8,+fp16fml,+fullfp16,+lse,+neon,+ras,+rcpc,+rdm,+sha2,+sha3,+sm4,+v8.1a,+v8.2a,+v8.3a,+v8.4a,+v8.5a,+v8a,+zcm,+zcz" }
attributes #6 = { nounwind "frame-pointer"="non-leaf" "no-trapping-math"="true" "probe-stack"="__chkstk_darwin" "stack-protector-buffer-size"="8" "target-cpu"="apple-m1" "target-features"="+aes,+crc,+crypto,+dotprod,+fp-armv8,+fp16fml,+fullfp16,+lse,+neon,+ras,+rcpc,+rdm,+sha2,+sha3,+sm4,+v8.1a,+v8.2a,+v8.3a,+v8.4a,+v8.5a,+v8a,+zcm,+zcz" }
attributes #7 = { nocallback nofree nosync nounwind readnone speculatable willreturn }
attributes #8 = { allocsize(0) "frame-pointer"="non-leaf" "no-trapping-math"="true" "probe-stack"="__chkstk_darwin" "stack-protector-buffer-size"="8" "target-cpu"="apple-m1" "target-features"="+aes,+crc,+crypto,+dotprod,+fp-armv8,+fp16fml,+fullfp16,+lse,+neon,+ras,+rcpc,+rdm,+sha2,+sha3,+sm4,+v8.1a,+v8.2a,+v8.3a,+v8.4a,+v8.5a,+v8a,+zcm,+zcz" }
attributes #9 = { allocsize(0,1) "frame-pointer"="non-leaf" "no-trapping-math"="true" "probe-stack"="__chkstk_darwin" "stack-protector-buffer-size"="8" "target-cpu"="apple-m1" "target-features"="+aes,+crc,+crypto,+dotprod,+fp-armv8,+fp16fml,+fullfp16,+lse,+neon,+ras,+rcpc,+rdm,+sha2,+sha3,+sm4,+v8.1a,+v8.2a,+v8.3a,+v8.4a,+v8.5a,+v8a,+zcm,+zcz" }
attributes #10 = { allocsize(1) }
attributes #11 = { allocsize(0) }
attributes #12 = { allocsize(0,1) }
attributes #13 = { returns_twice }
attributes #14 = { nounwind }

!llvm.module.flags = !{!0, !1, !2, !3, !4}
!llvm.ident = !{!5}

!0 = !{i32 2, !"SDK Version", [2 x i32] [i32 14, i32 4]}
!1 = !{i32 1, !"wchar_size", i32 4}
!2 = !{i32 8, !"PIC Level", i32 2}
!3 = !{i32 7, !"uwtable", i32 1}
!4 = !{i32 7, !"frame-pointer", i32 1}
!5 = !{!"Apple clang version 15.0.0 (clang-1500.3.9.4)"}
!6 = distinct !{!6, !7}
!7 = !{!"llvm.loop.mustprogress"}
!8 = distinct !{!8, !7}
!9 = distinct !{!9, !7}
!10 = distinct !{!10, !7}
!11 = distinct !{!11, !7}
!12 = distinct !{!12, !7}
!13 = distinct !{!13, !7}
!14 = distinct !{!14, !7}
!15 = distinct !{!15, !7}
!16 = distinct !{!16, !7}
!17 = distinct !{!17, !7}
!18 = distinct !{!18, !7}
!19 = distinct !{!19, !7}
!20 = distinct !{!20, !7}
!21 = distinct !{!21, !7}
!22 = distinct !{!22, !7}
!23 = distinct !{!23, !7}
