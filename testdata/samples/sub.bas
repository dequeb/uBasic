Sub Message(Optional onServer As Boolean = True)
    If onServer Then
        Debug.Print "Hello World!"
    Else
        MsgBox "Hello World!"
    End If
End Sub

Sub Message2(onServer As Boolean)
    If onServer Then
        Debug.Print "Hello World2!"
    Else
        MsgBox "Hello World2!"
    End If
End Sub
Call Message(True)
Call Message2(True)
