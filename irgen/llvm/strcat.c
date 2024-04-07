
#include <stdlib.h>
#include <stdio.h>
#include <string.h>

char *p;
const char *hello = "Hello, ";
const char *world = "world!";
int main(void)
{
	p = malloc(sizeof *p * 100 + 1);
	strcat(p, hello);
	strcat(p, world);
	printf("Final concatenated string : %s\n", p);
	free(p);
	return (0);
}
