'/* Test file for semantic errors. Contains exactly one error. */


Const a As Long = 1
Dim b As Double
' valid
Print a + b: Print a > b And a < b
' invalid
Print a Or b
