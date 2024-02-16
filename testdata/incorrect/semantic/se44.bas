'/* Test file for semantic errors. Contains exactly one error. */

' ParamArray must be last
Sub foo(ParamArray b() As Long, a As Long)

End Sub
