' /* Test file for syntactic errors. Contains exactly one error. */

Funct bar() As Long 
  Funct foo() as Boolean
  End Funct   '	// Local procedure definitions are not allowed
End Funct


