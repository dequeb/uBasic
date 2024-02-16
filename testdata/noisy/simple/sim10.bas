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
    Debug.Print "in loop"
    let c = True
Loop

Do
    let a = 48+b-1+a
    let b = b + 1
Loop While b < 10

let c= True
Do
    Debug.Print "in loop"
    let c = False
Loop Until c
