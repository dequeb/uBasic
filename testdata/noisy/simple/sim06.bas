' /* Our first control structure! Should print 9876543210 */


Dim t(2) As Long

Sub main() 
  Dim b As Long
  let b = 10
  let t(1) = 0
  Do While b > 0
    let t(0) = 48+b-1
    let b = b - 1
  Loop
  let t(0) = 10
  Print t(0), t(1)
End Sub

call main()
