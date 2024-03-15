
#include <stdio.h>
#include <string.h>

int  i;
long l;
double d;
float f;
int main(void)
{
	i = 10;
	l = i;
	i = l;
	f = 35.0008777;
	d = f;
	f = d;

	printf("i: %i, l: %li, d: %lf, f: %f\n", i, l, d, f);
	return (0);
}
