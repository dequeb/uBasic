Dim n As Long
Dim f As Long
Dim i As Long

Function Factorial(n As Long) As Long
  let f = 1
  For i = 1 To n
      let f = f * i
  Next
  let Factorial = f
End Function

For n = 0 To 10
    ' Print result.
  Debug.Print "The factorial of ", n, " is ", Factorial(n)
Next n

