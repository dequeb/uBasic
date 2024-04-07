' /* A simple program that tests output of strings. One string is a
'    global array, the other is local. Of course, neither array has a
'    length that is a multiple of 4. Outputs

'    Hello
'    Good bye
' */

Dim x(6) As String

Sub main()
  Dim y(9) As String
  let x(3) = "Hello World"
  let y(5) = "Good bye"
  Print x(3), y(5)
End Sub
call main()
