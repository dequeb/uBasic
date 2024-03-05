'/* Test file for semantic errors. Contains exactly one error. */


Dim Primes(1) As Long 
' Cannot redim a fixed array
ReDim Primes(5)      
