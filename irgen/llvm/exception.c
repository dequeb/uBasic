#include <stdio.h>
#include <setjmp.h>

jmp_buf jump_buffer;

void throw_exception() {
    longjmp(jump_buffer, 1);
}

void function_that_might_throw() {
    // Some code...
    throw_exception();
    // More code...
}

int main() {
    if (setjmp(jump_buffer) == 0) {
        function_that_might_throw();
    } else {
        printf("An exception was thrown!\n");
    }
    return 0;
}