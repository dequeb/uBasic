Function foo(x() As Long) As Long
  foo = x(0)
End Function  

Sub main()
  Dim a(10) As Long
  a(0) = foo(a)
End Sub
