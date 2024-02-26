' option base 0 
Const dimension As Long = 10
Dim array(dimension) As Long
Dim a(10+8) As Double
Dim b(dimension*33) As String

Let array(0) = 0
' in fact, it's 0 to 9 but 0 has been initialized before loop
Dim i As Long
For i = 1 To 9	
	Let array(i) = i + array(i - 1)
	Print i, array(i)
	If i = 5 Then
		Exit For
	End If
Next

Do While True
	Print "In Do While"
	Exit Do
Loop

Do Until False
	Print "In Do Until"
	Exit Do
Loop

Do
	Print "Do..."
	Print "Loop While"
Loop While False

Do
	Print "Do..."
	Print "Loop Until"
Loop Until True
