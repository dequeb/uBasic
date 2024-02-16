Dim Primes() As Long

Function IsPrime(Number As Long) As Boolean
  Dim I As Long
  For I = LBound(Primes) To UBound(Primes)
      If (Number Mod Primes(I)== 0) Then 
        Let IsPrime = False
        Exit Function
      End If
      If (Primes(I) >= Sqr(Number)) Then 
        Exit For
      End If
  Next
  Let IsPrime = True
End Function


 Sub BuildPrimes(Max As Long)
  If (Max < 3) Then 
    Exit Sub
  End If

  Dim I As Long
  For I = 2 To Max
    If (IsPrime(I)) Then
      ReDim Preserve Primes(UBound(Primes) + 1)
      Let Primes(UBound(Primes)) = I
    End If
  Next
End Sub

Sub Initialize()
  ReDim Primes(1)
  Let Primes(0) = 2
End Sub
call Initialize()
call BuildPrimes(100)

