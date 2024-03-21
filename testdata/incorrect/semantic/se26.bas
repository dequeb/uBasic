'/* Test file for semantic errors. Contains exactly one error. */


Sub f(a() As Double)
' valid numeric implicit conversion
  Let a(0) = 3
End Sub

Sub main() 
  Dim b(10) As Boolean
  '  //passing array of boolean to array of Double
  call f(b)  
End Sub

