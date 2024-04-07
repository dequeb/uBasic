' /* Test file for semantic errors. Contains exactly one error. */

Function q(a As Long,  b As Long , c As Long) As Long
  Let q = a*a + b*b + c*c 
End Function
' // Too few arguments to function 'q'
Print  1 + q(1, 3)  
