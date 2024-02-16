'/* Test file for semantic errors. Contains exactly one error. */


Dim a As DateTime = #1999/12/31#
Dim b As DateTime = #2000/12/31#
' valid
Print a <= b
' invalid
Print a > 4
