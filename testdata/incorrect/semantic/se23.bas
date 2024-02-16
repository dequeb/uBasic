' /* Test file for semantic errors. Contains exactly one error. */

Function first(a As String) As String
  Dim b As String
'  //not an array!
  Let a = b(0) 
End Function

