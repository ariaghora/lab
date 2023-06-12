CFLAGS=-Wall -Werror -pedantic -std=c99 -lraylib

all:
	@gcc -o tinyplot.out main.c libtinyplot.c renderer_raylib.c $(CFLAGS)

clean:
	@rm -f *.out

.PHONY: install
install: tinyplot.out
	@cp tinyplot.out /usr/local/bin/tinyplot
