' Fibonacci Sequence
' uBASIC Math Project
' ------------------------ 
' The array Fibo holds the Fibonacci numbers
Dim Fibo(52) As Long
Let Fibo(0) = 0
Let Fibo(1) = 1
Print "Fibonacci Sequence"
Print "-------------------"
Print "0"

Dim N As Integer
For N = 1 To 50
    Let Fibo(N+1) = Fibo(N) + Fibo(N-1)
    Print Fibo(N), ", ";
Next N
