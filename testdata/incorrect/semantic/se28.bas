' /* Test file for semantic errors. Contains exactly one error. */

Sub a (n As Long)
  If 1<2 Then
  ' Attempt to return value from procedure
    Let a = 2 * n   
  End If
End Sub

