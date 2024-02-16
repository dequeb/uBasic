'/* Test file for semantic errors. Contains exactly one error. */


Dim a As String = "Hello"
Dim b As String = "World"
' valid
Print a & " " &  b 
' Cannot add String (& is used for concatenation instead of +)
Print a + b
