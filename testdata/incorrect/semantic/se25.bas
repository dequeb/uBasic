' /* Test file for semantic errors. Contains exactly one error. */

' erase a dynamic array
Dim b() As Long
ReDim b(10) 
Erase b

' erase a non dynamic array
Dim a(10) As Long : Erase a
