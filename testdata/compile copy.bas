Dim l As long, d As Double, b As Boolean, da As Date, _
s As single, s0 As String, c As Currency, _
i As integer
Let l = 10: Let d = 3.141592: Let da = #2010/12/31#
Let s0="Allo le monde": 
Let c = 100.0$
Let b = True    
print "variables:"
print l, ", ";
print d, ", ";
print b
print da, ", ";
print s0, ", ";
print c

Const c1 As Long = 10, c2 As double = 11.5 , c3 As string = "Hello", c4 As boolean = True, c5 As date = #0001/01/01 00:00:05#
print "constants:"

print "c1:", c1 , ", c2:", c2 , ", c3:" , c3 , ", c4:",c4, ", c5:",c5

print "literals:"
print "world", True, 1.23456, 10, #2010/12/31#

Const c0 As Long = 10



' type conversion
'  convert from larger To smaller
Let i = 10: Let s = 100.01
'  convert from smaller To larger
Let l = i: Let d = s

print i
print s
print l
print d

Let i = 1234 Let l = 1234567890
'  convert from int To float
Let s = i + 1: Let d = l + 1

print i
print s
print l
print d


Let s = 12.34 Let d = 12345678.90
'  convert from float To int
Let i = s : Let l = d

print i
print s
print l
print d

' invalid implicit conversion To string
' Dim si As string, ss As string, sl As string, sd As string

' Let s = 12.34 Let d = 12345678.90
' Let l = 1234567890 Let i = 1234

' Let si = i: Let ss = s: Let sl = l: Let sd = d
' print si
' print ss
' print sl
' print sd

Let d = 12345678.90
'  convert from float To int
Let s = d
Let l = s
Let i = l

print i
print s
print l
print d

Sub ifBranch()
    Dim d As Long
    Let d = 10
    If c0 <> d Then
        print "c0 is Not d"
    Else
        print "c0 is d"
    End If
end Sub

Call ifBranch()

Dim s1 As string, s2 As string
Let s1 = "hello"
Let s2 = "world"
Dim s3 As string

Let s3 = s1 & " " & s2 & " !"
print s3
print "Hello " & "world !"
print 1 + 2
print 10.9 *0.98
print 1 + 2.98
print True And False
print True And True
print False And False
print False And True
print True Or False
print True Or True
print False Or False
print False Or True
print Not True

print 2 * 2.25$
print 2 / 2
print 2 < 3
print 2 > 3
print 2 <= 3
print 2 >= 3
print 2 == 3
print 2 <> 3
print 2.5 > 3.0
print 2.5 < 3.0
print 2.5 <= 3.0
print 2.5 >= 3.0
print 2.5 == 3.0
print 2.5 <> 3.0


Const k0 As long = 10
Const k1 As long = k0 + 10

Dim i11 As Integer
Let i11 = 10

Do While i11 > 0
    print i11
    Let i11 = i11 - 1
Loop

Sub while1()
    Dim i As Integer
    Let i = 20
    Do While i > 10
        print i
        Let i = i - 1
    Loop
end Sub

Call while1()


Function times(m As integer, Optional n As long = 10 + 2) As Long
    Let times = n * m
end Function


Function times2(Optional n As long = 10 + 2) As Long
    Let times2 = n * 2
end Function

print times(2)
print times(3, 5)

Sub listParmArray(ParamArray a() As integer)
    print a(0)
    print a(1)
    print a(2)
end Sub

call listParmArray(2, 5, 9)

sub addParmArray(ParamArray a() As integer)
    print a(0) + a(1) + a(2)
end Sub
call addParmArray(17, -19, 25)


function addDefaultValue(optional a As integer = 1, optional  b As integer = 2) As integer
    let addDefaultValue = a + b
end function

print addDefaultValue(10, 4)
print addDefaultValue(10)
print addDefaultValue()


sub addByRef(ByRef a As integer)
    Let a = a + 1
    print "in addByRef A: ", a
end sub

sub addByVal(ByVal a As integer)
    Let a = a + 1
    print "in addByVal A: ", a
end sub

dim addBy As Integer
Let addBy = 1
print addBy
call addByRef(addBy)
print addBy

print addBy
call addByVal(addBy)
print addBy


sub FloatAddByRef(ByRef a As Single)
    Let a = a + 0.99
    print "in FloatAddByRef A: ", a
end sub

sub FloatAddByVal(ByVal a As Single)
    Let a = a + 4.0
    print "in FloatAddByVal A: ", a
end sub

dim FloatAddBy As Single
Let FloatAddBy = 3.0
print FloatAddBy
call FloatAddByRef(FloatAddBy)
print FloatAddBy

print FloatAddBy
call FloatAddByVal(FloatAddBy)
print FloatAddBy



print 1.0/3
Print 3/2.0
Print 3/2
'print 1/0    ' generate a runtime error

Dim b8(3, 3) As integer
Let b8(0,0) = 2
' print b8(0,0)
Let b8(0,1) = 3
' print b8(0,1)
Let b8(0,2) = 4
' print b8(0,2)
Let b8(1,0) = 5
' print b8(1,0)
Let b8(1,1) = 6
' print b8(1,1)
Let b8(1,1) = 7
' print b8(1,1)
Let b8(1,2) = 8
' print b8(1,2)
Let b8(2,1) = 9
 print b8(1,2) + b8(2,1)
 print b8(0,0) - b8(1,1) + b8(2,1) * -b8(1,2)
 print b8(1,2)  - (b8(1,1) + b8(2,1) * b8(1,2)) 
print 8 * b8(1,2)  - (b8(1,1) + b8(2,1) * b8(1,2)) 
Dim a(2) As Long
Let a(1) = 2
Let a(0) = 23
print a(1) , ", ", a(0)
' let a(2) = 24
' print a(2)

dim curr As Currency
dim lon As Long
dim sng As Single
dim dbl As Double
dim str As String
dim int As Integer
dim bool As Boolean

Let curr = 1.0$
let lon = 10.99
let sng = 1.1
let dbl = 1.3
let str = "hello"
let int = 1
let bool = false

print curr
print lon
print sng
print dbl
print str
print int
print bool

dim curra(2) As Currency
dim lona(2) As Long
dim snga(2) As Single
dim dbla(2) As Double
dim stra(2) As String
dim inta(2) As Integer
dim boola(2) As Boolean

let curra(1) = 1.05$
let lona(1) = 123456789
let snga(1) = 1.1
let dbla(1) = 1.3
let stra(1) = "hello"
let inta(1) = 1
let boola(1) = false

print curra(1)
print lona(1)
print snga(1)
print dbla(1)
print stra(1)
print inta(1)
print boola(1)

const currcg As Currency = 10 + 8
const loncg as long = 10.9 - 0.07 + currcg
const sngcg  as single = 100 * 2 - loncg
const dblcg as double = 17 * 3 + currcg
const strcg as string = "hello" & " world"
const intcg as integer = 1.5 * 2.3 + dblcg
const boolcg as Boolean = False
const datecg as Date = #2024/01/01#

sub setLocalConstants() 
    const currcl As Currency = 10 + 8 + currcg
    const loncl as long = 10.9 - 0.07 + currcl
    const sngcl  as single = 100 * 2 - loncl
    const dblcl as double = 17 * 3 + currcl
    const strcl as string = "hello" & " world! " & strcg
    const intcl as integer = 1.5 * 2.3 + dblcl
    const boolcl as Boolean = False
    const datecl as Date = #2024/01/01#

    print currcl
    print loncl
    print sngcl
    print dblcl
    print strcl
    print intcl
    print boolcl
    print datecl
end sub

print currcg
print loncg
print sngcg
print dblcg
print strcg
print intcg
print boolcg
print datecg

call setLocalConstants()

dim arrayf(3,2) As integer
dim  arrayd() as Currency
sub declareLocalArrays()
    dim arrayf(2,3) As integer
    dim  arrayd() as Currency
    Let arrayf(0,1) = 4
    print arrayf(0,1) 
    print arrayf(1,1)

end sub
let arrayf(2,0) = 3
call declareLocalArrays()
print arrayf(2,0)
