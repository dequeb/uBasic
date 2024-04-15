dim bzzz(2,2) as integer
let bzzz(0,0) = 4

Sub AddArrayNumbers ()
    Dim a(3, 2) As Integer    
    Let a(0, 1) = 1
    Let a(1, 0) = 2
    Let a(1, 1) = 3
    Let bzzz(0, 0) = 4
    print a(0,0), a(0,1), a(1, 0), a(1, 1), bzzz(0, 0)
    print a(0,0)+ a(0,1)+ a(1, 0)+ a(1, 1) + bzzz(0, 0)
end Sub

call AddArrayNumbers()