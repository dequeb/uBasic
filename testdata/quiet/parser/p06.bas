




sub  f() 
  
end sub 

function g() as integer
  let g = 42
end function

sub  main() 
  call f()
  print g()
end sub

call main()