#include <stdio.h>
#include <stdlib.h>

int *a1;
int _a1_length = 0;

int main()
{
    // dynanmicaly allocate memory for a1 and initialize it
    a1 = (int *)calloc(2, sizeof(int));
    _a1_length = 2;

    int index;
    index = 10;
    if (index < _a1_length)
    {
        a1[index] = 2;
    }
    else
    {
        printf("Index out of bounds\n");
        return 1;
    }
    return 0;
}