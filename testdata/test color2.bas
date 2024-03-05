resume next ' should be ignored
on error goto error1
print 1/0
error1:
print "skipped error"
resume next
on error resume next
print 1/0
on error goto 0 


do Until True
    print "do until"
    exit do
loop
do while false
    print "do while"
    exit do
loop
do 
    print "do"
    exit do
loop until true
do 
    print "do"
    exit do
loop while false
dim a(2) as integer
let a(0) = 1
let a(1) = 2
Dim x as integer
for each x in a
    print x
    exit for
next
dim b() as currency
erase b

dim y as single
select case y
    case 3
        print "3"
    case 4
        debug.print "4"
    case else
        MsgBox "else"
end select
