'/* Test file for semantic errors. Contains exactly one error. */

Enum a
    b
    c
End Enum

Dim d As Variant
Dim e As Variant
' valid
Print d - e
' invalid
Print e - c
