'/* Test file for semantic errors. Contains exactly one error. */

Sub foo(n As Long )
  Print ""
End Sub 

';	// 'foo' does not return a value
Dim A As Long 
Let A = foo(0)
