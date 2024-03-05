@str = global [15 x i8] c"Hello, world!\0A\00"

declare i32 @puts(i8* %0)

define i32 @main() {
0:
        %1 = call i32 @puts(i8* getelementptr ([15 x i8], [15 x i8]* @str, i64 0, i64 0))
        ret i32 0
}
