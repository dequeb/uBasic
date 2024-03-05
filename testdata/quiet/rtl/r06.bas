
Function fac(n As Long) As Long
  Dim i As Long
  Dim p As Long
  If n < 0 Then
    Let fac = 0
    Exit Function
  End If
  Let i = 0
  Let p = 1
  Do While i < n
    Let i = i + 1
    Let p = p * i
  Loop
  Let fac = p
End Function  ' fac

Dim a As Long
Let a = fac(5)
