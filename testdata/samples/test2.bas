
Dim array(2) As String
Dim str(2) As String

Let array(0) = "Hello "
Let array(1) = "World!"
ReDim Preserve array(4)
Let array(2) = "Bonjour "
Let array(3) = "le monde!"
Let str(0) = array(0) & array(1)
Let str(1) = array(2) & array(3)
Print str(0), ", ", str(1); ' a comment...