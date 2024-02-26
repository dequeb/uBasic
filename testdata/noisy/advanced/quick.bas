Sub Quicksort(list() As Long, min As Long, max As Long)
Dim med_value As Long
Dim hi As Long
Dim lo As Long
Dim i As Long

    ' If min >= max, the list contains 0 or 1 items so it
    ' is sorted.
    If min >= max Then 
        Exit Sub
    End If

    ' Pick the dividing value.
    let i = Int((max - min + 1) * Rnd() + min)
    let med_value = list(i)

    ' Swap it to the front.
    let list(i) = list(min)

    let lo = min
    let hi = max
    Do While True
        ' Look down from hi for a value < med_value.
        Do While list(hi) >= med_value
            let hi = hi - 1
            If hi <= lo Then 
                Exit Do
            End If
        Loop
        If hi <= lo Then
            let list(lo) = med_value
            Exit Do
        End If

        ' Swap the lo and hi values.
        let list(lo) = list(hi)
        
        ' Look up from lo for a value >= med_value.
        let lo = lo + 1
        Do While list(lo) < med_value
            let lo = lo + 1
            If lo >= hi Then 
                Exit Do
            End If
        Loop
        If lo >= hi Then
            let lo = hi
            let list(hi) = med_value
            Exit Do
        End If
        
        ' Swap the lo and hi values.
        let list(hi) = list(lo)
    Loop
    
    ' Sort the two sublists.
    call Quicksort(list, min, lo - 1)
    call Quicksort(list, lo + 1, max)
End Sub
 
Dim list(100) As Long
Dim i As Long
For i = 0 To 99
    let list(i) = Int(1000 * Rnd())
Next i

call Quicksort(list, 0, 99)
