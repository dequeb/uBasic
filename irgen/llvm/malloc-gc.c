
#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include "gc.h"

char *p;
int main(int argc, char *argv[])
{
	gc_start(&gc, &argc);
	p = gc_malloc(&gc, sizeof *p * 100 + 1);
	strcpy(p, "This is tutorialspoint.com");
	printf("Final copied string : %s\n", p);
	gc_stop(&gc);
	return (0);
}
