' /* Test file for semantic errors. Contains exactly one error. */

Sub foo(n As Long)
  Print n
End Sub

' 'foo' does not return a value
Print 1 + foo(0)	
