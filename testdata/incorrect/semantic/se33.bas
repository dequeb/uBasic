'/* Test file for semantic errors. Contains exactly one error. */

' valid date-time constant
Dim a As DateTime = #1999/12/31 00:01:02#
' invalid date-time constant
Dim c As DateTime = #1999/01/01 12:99:99#
