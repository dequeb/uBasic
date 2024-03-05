'/* Test file for semantic errors. Contains exactly one error. */


Sub a (n As Long)
  If 1<2 Then
  ' Attempt Exit For without For
    Exit For  
  End If
End Sub
