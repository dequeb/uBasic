'/* Test file for semantic errors. Contains exactly one error. */


Const a As Boolean = True
Const b As Boolean = False
' valid
Print a Or b

' invalid
Print a Mod b
