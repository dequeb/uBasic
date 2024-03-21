Sub f(x() As Long)
  let x(3) = x(5) + 7
End Sub

Dim a(10) As Long
call f(a)
