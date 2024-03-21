' /* Test file for semantic errors. Contains exactly one error. */

Function a(n As Long) As Long
  Dim bv(10) As String
  ' Return from function with erroneous type
  Let a = bv  
End Function
