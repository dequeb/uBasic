
#include <stdio.h>
#include <string.h>

int  i;
long l;
double d;
float f;
int main(void)
{
	f = 12.34;
	d = 12345678.901;
	i = f;
	l = d;
	printf("i: %i, l: %li, d: %lf, f: %f\n", i, l, d, f);
	return (0);
}
