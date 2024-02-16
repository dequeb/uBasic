' Array ByRef and by value

Sub foo(ByRef a() As Long, b() As Long)
  let a(LBound(a))  = 1
  let b(LBound(b)) = 2
End Sub

Sub main()
  Dim a(10) As Long
  Dim b(10) As Long
  let a(LBound(a)) = 0
  let b(LBound(b)) = 0
  call foo(a, b)
  Debug.Print a(LBound(a)), b(LBound(b))
End Sub
