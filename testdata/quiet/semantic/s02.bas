Dim a As Long

Function foo(x As Long) As Long
  foo = x + 1
End Function

Sub main()
  a = 42
  a = foo(a)
End Sub
main()
