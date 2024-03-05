'/* Test file for semantic errors. Contains exactly one error. */


Dim a(1) As Long
' index out of bounds
Let a(-1) = 1  
' will be catch at run-time
 