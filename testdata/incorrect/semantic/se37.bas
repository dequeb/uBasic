'/* Test file for semantic errors. Contains exactly one error. */


Enum values 
    a
    b
End Enum

Sub Quicksort(list() As Long, min As Long, max As values)
    Print "Quicksort"
End Sub

Dim list(100) As Long
' not an error
Quicksort(list, 0, 99)  


  ' error: Enuum is not a valid Long
Quicksort(list, a, b)
