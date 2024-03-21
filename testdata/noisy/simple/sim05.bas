' will make sure ByRef is able to return value to caller

Sub foo(ByRef a As Long, b As Long)
  let a = 1
  let b = 2
End Sub


Sub main()
  Dim a As Long
  Dim b As Long
  let a = 0
  let b = 0
  call foo(a, b)
  Debug.Print a, b
End Sub

call main()