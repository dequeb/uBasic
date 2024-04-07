' /* Test file for semantic errors. Contains exactly one error. */

' Default paramater must be last
Function q(a As Long, Optional b As Long = 9, ByRef c As Long) As Long   
  Let q = a*a + b*b
End Function

