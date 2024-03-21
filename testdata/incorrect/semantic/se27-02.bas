'/* Test file for semantic errors. Contains exactly one error. */


Function a (n As Long) As Long 
  If 1<2 Then
  ' Attempt Exit Sub in function
    Exit Sub   
  End If
  Let a = 2
End Function 
