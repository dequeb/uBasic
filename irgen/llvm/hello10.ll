@ten = global i64 10
@str = global [13 x i8] c"Hello, %ld!\0A\00"

declare i32 @printf(i8* %0, ...)

define i32 @main() {
0:
        %1 = alloca i64
        store i64 10, i64* %1
        %2 = load i64, i64* getelementptr (i64, i64* @ten, i64 0)
        %3 = call i32 (i8*, ...) @printf(i8* getelementptr ([13 x i8], [13 x i8]* @str, i64 0, i64 0), i64 10)
        %4 = call i32 (i8*, ...) @printf(i8* getelementptr ([13 x i8], [13 x i8]* @str, i64 0, i64 0), i64* getelementptr (i64, i64* @ten, i64 0))
        %5 = call i32 (i8*, ...) @printf(i8* getelementptr ([13 x i8], [13 x i8]* @str, i64 0, i64 0), i64 %2)
        ret i32 0
}
