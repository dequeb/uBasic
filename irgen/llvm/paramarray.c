#include <stdio.h>

int add (int count, int *arr) {
    int sum = 0;
    for (int i = 0; i < count; i++) {
        sum += arr[i];
    }
    return sum;
}

int main() {
    int arr[2];
    arr[0] = 10;
    arr[1] = 20;
    printf( "%d\n", add(2, arr));
    return 0;
}