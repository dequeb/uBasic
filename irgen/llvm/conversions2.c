
#include <stdio.h>
#include <string.h>

int  i;
long l;
double d;
float f;
int main(void)
{
	i = 10;
	f = i;
	l = 1234567890;
	d = l;

	printf("i: %i, l: %li, d: %lf, f: %f\n", i, l, d, f);
	return (0);
}
