' /* Test file for semantic errors. Contains exactly one error. */

Dim a (10) As String 
'  Attempt to apply arithmetic to array reference
Dim b As Long
Let b = a + 1     
