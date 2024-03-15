dim l as long, d as Double, b as Boolean, da as Date, s as String, c as Currency    
let l = 10: let d = 3.141592: let da = #2010/12/31#
let s="Allo le monde": 
let c = 100.0$
let b = true    
print "variables:"
print l, ", ";
print d, ", ";
print b
print da, ", ";
print s, ", ";
print c

const c1 as Long = 10, c2 as double = 11.5 , c3 as string = "Hello", c4 as boolean = true, c5 as date = #0001/01/01 00:00:05#
print "constants:"

print "c1:", c1 , ", c2:", c2 , ", c3:" , c3 , ", c4:",c4, ", c5:",c5

print "literals:"
print "world", true, 1.23456, 10, #2010/12/31#

Const c0 As Long = 10



// type conversion
dim i as integer, s as single
dim l as long, d as double

'  convert from larger to smaller
let i = 10: let s = 100.01
'  convert from smaller to larger
let l = i: let d = s

print i
print s
print l
print d

let i = 1234 let l = 1234567890
'  convert from int to float
let s = i + 1: let d = l + 1

print i
print s
print l
print d


let s = 12.34 let d = 12345678.90
'  convert from float to int
let i = s : let l = d

print i
print s
print l
print d

' invalid implicit conversion to string
' dim si as string, ss as string, sl as string, sd as string

' let s = 12.34 let d = 12345678.90
' let l = 1234567890 let i = 1234

' let si = i: let ss = s: let sl = l: let sd = d
' print si
' print ss
' print sl
' print sd

let d = 12345678.90
'  convert from float to int
let s = d
let l = s
let i = l

print i
print s
print l
print d

sub ifBranch()
    dim d As Long
    let d = 10
    If c0 <> d Then
        print "c0 is Not d"
    Else
        print "c0 is d"
    End If
end Sub

call ifBranch()

function times2(n as long)
   times2 = n * 2
end function

print time2(10)