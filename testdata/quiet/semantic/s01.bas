' /* Some strange but legal expressions and statements. 
'    For more examples of semantically correct code, see noisy/simple 
' */

Dim x As Long
Dim y As Double

Sub main()
  Dim z As Long
  Dim w As Double

  x = x+y+z+w
  
  x = z = 42

  Dim a As Boolean
  a = z == 42

  x = (z = 99)  

  Do While (a) 
    x = 0
  Loop

  If (a)  Then
    y = 4
  Else
    y = 7
  End If
  
  a = x > y;

  a = 0 < x < 10;

End Sub
main()










