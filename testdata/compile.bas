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








Function times2(n As long) as Long
    let times2 = n * 2
end Function

print times2(10)