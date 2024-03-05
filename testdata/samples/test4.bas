Dim index As Double
Debug.Print """ Quotes """
For index = 10 To 0 Step -0.5
	Print index
	If index <= 2 Then
		Print "will exit for loop"
		Exit For
	End If
Next index

Dim var As Variant
Let var = Nothing
