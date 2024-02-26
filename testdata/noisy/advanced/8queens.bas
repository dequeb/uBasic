
' /* Test.
' ** 8-queen.

' Expected output:

' 04752613

' This corresponds to the following board:

' --x-----
' -----x--
' ---x----
' -x------
' -------x
' ----x---
' ------x-
' x-------

' */

Dim n As Long 
Dim board(8) As Long

Sub printboard(board() As Long)
  Dim i As Long
  Let i = 0
  Do While i < n
      Print board(i)
      Let i = i+1
  Loop
End Sub 

Function check(col As Long, row As Long) As Boolean
  Dim i As Long
  Dim j As Long
  Let i = col-1
  Do While i >= 0
    Let j = board(i)
    If j = row Then
      Let check = False
      Exit Function
    End If
    If j > row And col-i == j-row Then
      Let check = False
      Exit Function
    End If
    If col-i == row-j Then
      Let check = False
      Exit Function
    End If
    Let i = i-1
  Loop
  Let check = True
End Function

Function queen(col As Long, row As Long) As Boolean
  If col >= n Then 
    ' Returning false will generate all solutions...
    Let queen = True
    Exit Function
  End If
  Do While row < n
    Let board(col) = row
    If check(col,row) And queen(col+1,0) Then
      Let queen = True
      Exit Function
    End If
    Let row = row + 1
  Loop
  Let queen = False
End Function

Sub main()
  Let n = 8
  Dim null As Boolean
  let null = queen(0,0)
  call printboard(board)
End Sub

call main()

