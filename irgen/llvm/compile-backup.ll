%struct.GarbageCollector = type { ptr, i8, ptr, i64 }
@gc = external global %struct.GarbageCollector, align 8

@vbEmpty = constant [1 x i8] c"\00"
declare void @gc_start(ptr noundef, ptr noundef) #1
declare i64 @gc_stop(ptr noundef) #1
declare ptr @gc_malloc(ptr noundef, i64 noundef) #1

define i32 @main(i32 %argc, i8** %argv) {
0:
	%1 = alloca i32
	store i32 %argc, i32* %1
  	call void @gc_start(ptr noundef @gc, ptr noundef %1)
  
  	%2 = alloca [14 x i8]
 	%3 = call i64 @gc_stop(ptr noundef @gc)
	ret i32 0
}
