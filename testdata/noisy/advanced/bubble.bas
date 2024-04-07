' // This program illustrates the bubblesort algorithm by sorting an
' // array of char and printing the intermediate states of the array.


Sub BubbleSort(byref MyArray() As Variant)
  'Sorts a one-dimensional VBA array from smallest to largest
  'using the bubble sort algorithm.
  Dim i As Long, j As Long
  Dim Temp As Variant
  
  For i = LBound(MyArray, 1) To UBound(MyArray, 1) - 1
      For j = i + 1 To UBound(MyArray, 1)
          If MyArray(i) > MyArray(j) Then
              Let Temp = MyArray(j)
              Let MyArray(j) = MyArray(i)
              Let MyArray(i) = Temp
          End If
      Next j
  Next i
End Sub

Sub Main()
  Dim MyArray(26) As Variant
  Dim i As Long
  
  'Fill the array with a permutation of the characters a-z
  For i = 0 To 25
      Let MyArray(i) = Chr(97 + i)
  Next i
  
  'Print the original array
  For i = 1 To 26
      Debug.Print MyArray(i)
  Next i

  
  'Sort the array
  call BubbleSort(MyArray)
  
  'Print the sorted array
  For i = 1 To 26
      Debug.Print MyArray(i)
  Next i
End Sub

call Main()
