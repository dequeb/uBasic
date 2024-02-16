' call to different functions and subroutines with different parameters
Function a1 (a As Long, b As Long) As Long
    Let a1 = a + b
End Function

Function b_1 (d As Variant) As Variant
    Let b_1 = d
End Function

Function c_1 (e As String, f As String) As String
    Let c_1 = e & f
End Function

Function d__ (g As Boolean, h As Boolean) As Boolean
    Let d__ = g And h
End Function

' optional parameters
Function e_1(j As Long, Optional i As Long = 9.0) As Long
    Let e_1 = i - j
End Function

Sub f_1 (ByRef k As Long, Optional ByRef l As Long = 9.0)
    Print k - l
    Let k = k + 1
    Let l = l - 1
End Sub

