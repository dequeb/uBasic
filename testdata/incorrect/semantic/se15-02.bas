' /* Test file for semantic errors. Contains exactly one error. */

Function q(a As Long, Optional b As Long = 9) As Long   
  Let q = a*a + b*b
End Function

Print q(1,2)
Print q(1)
Print q() ' this should fail

