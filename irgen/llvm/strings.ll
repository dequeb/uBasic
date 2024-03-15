@str = global i8* null
@.str0 = global [15 x i8] c"Hello, World!\0A\00"

declare i8* @strcpy(i8* %dst, i8* %src)

declare i32 @puts(i8* %s)

declare i8* @malloc(i32 %size)

declare void @free(i8* %ptr)

define i32 @main() {
0:
        %1 = call i8* @malloc(i32 15)
        store i8* %1, i8** @str
        %2 = load i8*, i8** @str
        %3 = call i8* @strcpy(i8* %2, [15 x i8]* @.str0)
        %4 = load i8*, i8** @str
        %5 = call i32 @puts(i8* %4)
        call void @free(i8* %4)
        ret i32 0
}