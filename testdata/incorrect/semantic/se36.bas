'/* Test file for semantic errors. Contains exactly one error. */


  ' Cannot redim preserve to a negative size
Dim Primes() As Long 
ReDim Preserve Primes(-1) 
