'/* Test file for semantic errors. Contains exactly one error. */

' valid date constant
Dim a As DateTime = #1999/12/31#
' invalid date constant
a = #1999/99/99#

