' /* Fibbonacci, in the simple and naive form */
' /* Prints fibbonacci numbers for n=1..12 */

Function fib(n As Long) As Long
  If (n == 0) Then
    Let fib = 1
  Else
    If (n == 1) Then
      let fib = 1
    Else
      let fib = fib(n-1) + fib(n-2)
    End If
  End If
End Function

Dim i As Long
For i = 1 To 12
  Debug.Print i, " ", fib(i)  
Next i
