' /* Test file for semantic errors. Contains exactly one error. */


Function d( a As Long , b As Long ) As Long
	Let d = b - a
End Function

 ' ;	// Too many arguments to function 'd'
Print d(1, 2, 3)   
