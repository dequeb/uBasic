' /* Test file for semantic errors. Contains exactly one error. */

Function a(n As Long) As Long
  Dim v As Variant
  ' a variant may contain Nothing
  Let v = Nothing
  ' incompatible types
  Let a = Nothing 
End Function
