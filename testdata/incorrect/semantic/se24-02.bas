' /* Test file for semantic errors. Contains exactly one error. */

Enum values
    a
End Enum

'this is legal
Dim b As values
Let b = a       

Function foo() As values
    let foo = a
End Function

Dim c(5) As values

Const w As values = a

' a cannot be assigned to, because it is a constant
Let a = 20      
