Dim Primes() As Long

Function IsPrime(Number As Long) As Boolean
  Dim I As Long
  dim SqrNumber as Long
  Let SqrNumber = Sqr(Number)
  For I = LBound(Primes) To UBound(Primes)
      If (Number Mod Primes(I)== 0) Then 
        Let IsPrime = False
        Exit Function
      End If
      If (Primes(I) >= SqrNumber) Then 
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
  For I = 3 To Max
    If (IsPrime(I)) Then
      ReDim Preserve Primes(UBound(Primes)  + 2)	'option base 0 -> ubound == 0
      Let Primes(UBound(Primes)) = I
    End If
  Next
End Sub

Sub Initialize()
  ReDim Primes(1)
  Let Primes(0) = 2
End Sub
call Initialize()
dim before as Date, after as Date
let before = Time()
call BuildPrimes(20000)
let after = Time()

dim i as Long
print before, " - ", after

for i = LBound(Primes) to UBound(Primes)
  Debug.Print Primes(i), ", ";
next i

