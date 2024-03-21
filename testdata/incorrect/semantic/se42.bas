'/* Test file for semantic errors. Contains exactly one error. */

Enum Seasons
    Spring
    Summer
    Autumn
    Winter
End Enum

Dim b As Boolean
' type mismatch, cannot assign Seasons to Boolean
Let b = Spring 
