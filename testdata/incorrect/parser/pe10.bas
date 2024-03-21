'/* Test file for syntactic errors. Contains exactly one error. */

Sub foo(a As Long, b As Long, c As Long)
  Print a, b, c;
End Sub

foo(1, 2, ) ' // Unexpected token ')'

