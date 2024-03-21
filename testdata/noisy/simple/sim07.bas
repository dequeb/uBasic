sub sub1(byref a as long) 
    let a = 1
    exit sub
    let a = 2
end sub

function func1() as long
    let func1 = 1
    exit function
    let func1 = 2
end function

dim a as long
call sub1(a)
print a
print func1()
print "this is the end."
