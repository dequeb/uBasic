Dim a As Boolean
Dim b As Boolean 
Let b = True

If a Or b Then
	MsgBox "Or"
End If

If a And b Then
	MsgBox "And Robin Hood"
ElseIf a Then
	MsgBox "ElseIf"
Else
	Print "Else, A:", a, ", B:", b
End If 

