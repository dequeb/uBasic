function one () as boolean
    on error goto help
    let one = true
end function

sub two ()
help:   ' invalid jump target for an error 
    print "error"
    resume next
end sub