' /* Test file for semantic errors. Contains exactly one error. */

Function a(n As Long) As Long
    Let a = 2 * n 
End Function

' // Redeclaration of 'a'
Function a(n As Long) As Long 
    Let a = n / 2
End Function
