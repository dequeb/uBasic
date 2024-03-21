' /* Test the implementation of global and local variables. Note that
'    output will appear as "7654321".*/

Dim foo As Long
Sub main() 
  Dim bar As Long
  let foo = 76: let  bar = 54321
  Debug.Print foo, bar
End Sub
call main()

