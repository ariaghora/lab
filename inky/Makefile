CFLAGS=-Wall -Werror -pedantic -std=c99 -O3 -lraylib

all:
	@gcc -o inky.out src/*.c src/renderer/*.c $(CFLAGS)

install:
	@cp inky.out /usr/local/bin/inky

.PHONY: clean
clean:
	@rm -f *.out
