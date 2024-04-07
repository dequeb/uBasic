@true = global [5 x i8] c"True\00"
@false = global [6 x i8] c"False\00"

declare i32 @puts(i8* %str)

define i8* @_fromCharToStringBoolean_(i1 %value) {
0:
        %1 = icmp eq i1 %value, true
        %2 = select i1 %1, [5 x i8]* getelementptr ([5 x i8], [5 x i8]* @true, i64 0), [6 x i8]* getelementptr ([6 x i8], [6 x i8]* @false, i64 0)
        ret [5 x i8]* %2
}

define i32 @main() {
0:
        %1 = call i8* @_fromCharToStringBoolean_(i1 true)
        %2 = call i32 @puts(i8* %1)
        %3 = call i8* @_fromCharToStringBoolean_(i1 false)
        %4 = call i32 @puts(i8* %3)
        ret i32 0
}