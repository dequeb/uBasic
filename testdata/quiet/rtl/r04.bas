Sub f(ByRef x As String)
  Print x
End Sub

Sub g(x As String)
  Print x
End Sub

Sub main() 
  call f("hello")
  call g("world")
End Sub

call main()