' /* Test file for semantic errors. Contains exactly one error. */

' redeclare a function
Function a() As Long
  Let a = 2
  Let a(2) = 2    
End Function

