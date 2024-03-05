
on error goto step1
dim a as integer
let a = 1/0
step1:
print a
on error resume next
print 1/0
print "error 2"
on error resume step2
print 1/0
stop
step2:
on error goto 0
if true then
    on error goto step3:
    print 1/0
    stop
end if
step3:
print "successfull"
stop
