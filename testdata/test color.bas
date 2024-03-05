dim a() as long ' this is a comment
redim preserve a(1) : print "Hello, World! ";
print #2021-10-30#
print 2.2, ", ", true, ", ", false, ", ", format(3.4$, "%.2f$"), ", ", 2, ", ", format(4$, "%.2f$")
' this is a comment
dim e as Variant, f as variant
Enum fa 
    aa
    ab
    ac
End Enum ' this is a comment
enum fb 
    ba  
    bb
    bc
End Enum
print fa.aa
if false then 
    stop 
elseif true then 
    ' this is a comment
    print "true"
else 
    print "false"
end if 
Let f = fb.bb
Let e = fa.ab
let f = e
print e, f
const g1 as variant = 3
print g1
function g() as long
    let g = 5
end function
function h() as long
    let h = 3
    if 3 == h OR h == 4 or g() == 5 then 
        let h = 4
        exit function
    else 
        let h = 5
    end if
end function
sub i()
    do while true
        print "in loop"
        exit do
    loop
end sub
print h()
call i()
Let e = 2 *4 + 5 / 9 -(45 div 5) mod 3
let f = "hello" & "world"
print e, f

