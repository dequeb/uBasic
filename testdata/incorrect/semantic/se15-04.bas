' /* Test file for semantic errors. Contains exactly one error. */

' one optional parameter
Sub w(Optional b As Long = 9)
  Print ""
End Sub

' two optional parameters
Sub x(Optional a As Long = 9, Optional b As Double = 9)
  Print ""
End Sub

' one required parameter and one optional
Sub y(a As Long, Optional b As Double = 9)
  Print ""
End Sub

' one required parameter and one parameter array
Sub z(a As Long, ParamArray b() As String)
  Print ""
End Sub

'test call to all subtoutines
call w()
call w(2)
call x( )
call x(2)
call x(2, 2.2)
call y(2, 2.2  )
call y(2)
call z(2)
call z(2, "2")
call z(2, "2", "2")

' error in the last parameter type
call z(2, "2", "2", False)
