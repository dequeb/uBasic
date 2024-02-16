' /* Passing arrays as arguments of functions poses some particular
'    difficulties. This program should output
' */

Sub f(a() As Long)
  Print a(3)
End Sub

Sub g(b() As Long)
  call f(b)
End Sub

Sub fc(a() As Double) 
  Print a
End Sub

Sub gc(a() As Double) 
  call fc(a)
End Sub 

Dim x(10) As Long
Dim xc(10) As Double

Sub main ()
  Dim y(10) As Long
  Dim yc(10) As Double

  let x(3) = 12
  let y(3) = 34

  call f(x)
  call f(y)
  
  let x(3) = 56
  let y(3) = 78

  call g(x)
  call g(y)

  let xc(0) = 1.0
  let xc(1) = 1.1
  let xc(2) = 1.2

  let yc(0) = 1.2
  let yc(1) = 2.3

  call fc(xc)
  call fc(yc)

  let xc(0) = 99.34
  let xc(1) = 0

  let yc(0) = 1008983.42345
  let yc(1) = -1234.1234
  let yc(2) = 0
 
  call gc(xc)
  call gc(yc)
End Sub
