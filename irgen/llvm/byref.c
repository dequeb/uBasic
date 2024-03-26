void foo(int* var1) { (*var1) = (*var1) + 1; }
int main() {int a=1; foo(&a ); }