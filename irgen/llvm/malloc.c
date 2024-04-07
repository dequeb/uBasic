
#include <stdlib.h>
#include <stdio.h>
#include <string.h>

char *p;
int main(void)
{
	p = malloc(sizeof *p * 100 + 1);
	strcpy(p, "This is tutorialspoint.com");
	printf("Final copied string : %s\n", p);
	free(p);
	return (0);
}
