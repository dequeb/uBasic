' /* Let's try a bunch of function calls. Should output
'    01234567890123456789 */


Sub foo(v1 As Long, v2 As Long, v3 As Long, v4 As Long, v5 As Long)
  Debug.Print v1, v2, v3, v4, v5
End Sub

Function f(x As Long) As Long
  let f = x + 1
End Function

Function g(x As Long) As Long
  If (x = 0) Then
    let g = 1
  Else
    let g = 2 * g(x-1)
  End If
End Function

Dim x As Long
call foo(0,1,2,3,4)
let x = 5
call foo(x+0, x+1, x+2, x+x-2, x*2-1)

call foo(0, f(0), f(f(0)), f(f(f(0))), g(2))
call foo(g(2)+g(0), g(2)+g(1), g(0)+g(1)+g(2), g(3), g(4)-7)

