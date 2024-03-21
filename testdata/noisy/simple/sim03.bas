' /* Test the implementation of global and local arrays. Note that
'    output will appear as "123456".*/

Dim a(10) As Long

Sub main()
  Dim b(10) As Long
  Let a(7) = 123: let  b(5) = 456
  Debug.Print a(7), b(5)
End Sub
call main()
