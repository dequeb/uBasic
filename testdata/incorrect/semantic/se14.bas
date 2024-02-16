' /* Test file for semantic errors. Contains exactly one error. */

Function f(n As Long ) As Long
  Let f = n / 2
End Function


Sub p(n As Long)
  Dim f As Long, g As Boolean
  '  'f' refers only to the local variable
  Let f = n * 2 * f(n)  
End Sub


Dim n As Long
call p(n)

