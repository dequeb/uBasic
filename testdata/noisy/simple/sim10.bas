' control structures

Dim a As Long
Dim b As Long 
let b = 5
Do While b > 0
    let a = 48+b-1+a
    let b = b - 1
Loop
Debug.Print a

Dim c As Boolean
Let c = False
Do Until c
    Debug.Print "in loop 1"
    let c = True
Loop

Do
    let a = 48+b-1+a
    let b = b + 1
Loop While b < 10

let c= false
Do
    Debug.Print "in loop 2"
    let c = true
Loop Until c
