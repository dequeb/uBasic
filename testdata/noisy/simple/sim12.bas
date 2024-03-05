' numeric conversions

Dim i As Long
Dim j As Double
Dim k As Double 
For i = 1 To 2
    For j = 1 To -0.5 Step -0.5
        let k = i+j*i-j
    Next
' this should result incorrect iterator
Next i
