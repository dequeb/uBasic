'/* Test file for semantic errors. Contains exactly one error. */

Dim a As Long = 42

' Attempt to redefine  'a' as Sub
Sub a()     
  MsgBox "Hello"
End Sub
