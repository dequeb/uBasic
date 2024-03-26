

sub addByRef(byref a As Long)
    Let a = a + 1
    ' print "in addByRef A: ", a
end sub

' sub addByVal(ByVal a As integer)
'     Let a = a + 1
'     print "in addByVal A: ", a
' end sub

'call addByRef(0)

dim addBy As Long
Let addBy = 1
print addBy
call addByRef(addBy)
print addBy

' print addBy
' call addByVal(addBy)
' print addBy




