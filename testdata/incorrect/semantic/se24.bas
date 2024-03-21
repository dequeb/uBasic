' /* Test file for semantic errors. Contains exactly one error. */


Sub first (a() As String)
  Dim b(10) As String
  	' b cannot be assigned
  let b = a	
End Sub
