' /* Test file for semantic errors. Contains exactly one error. */

Dim b As Long

' multiply scalars with array
Function q(a As Long, ParamArray c() As Long) As Long   
  Let q = a*a + b*b + c*c
End Function

