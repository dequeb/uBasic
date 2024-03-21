' /* Test file for semantic errors. Contains exactly one error. */

Function q(a As Long, Optional b As Long = 9, ParamArray c() As Long) As Long   
  Let q = a*a + b*b + c*c ' invalid use of c
End Function

