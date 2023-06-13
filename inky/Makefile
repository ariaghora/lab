CFLAGS=-Wall -Werror -pedantic -std=c99 -O3 -lraylib

all:
	@gcc -o inky.out *.c renderer/*.c $(CFLAGS)

.PHONY: clean
clean:
	@rm -f *.out
