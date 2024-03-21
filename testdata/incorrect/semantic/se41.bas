'/* Test file for semantic errors. Contains exactly one error. */

Function a() As Long
    a = 1
End Function

' function value not used
call a() 
