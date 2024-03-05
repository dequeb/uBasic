' /* Test file for semantic errors. Contains exactly one error. */


Sub a(n As Long)
  '	// Attempt to return value from procedure
  Let a =  2 * n      
End Sub


