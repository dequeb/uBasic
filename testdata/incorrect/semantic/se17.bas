' /* Test file for semantic errors. Contains exactly one error. */


Dim hello(5) As String 
  ' ; //  Attempt to use string array in arithmetic
Print hello(0) + 1 

