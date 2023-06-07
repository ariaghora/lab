#include <ctype.h>
#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/ioctl.h>
#include <termios.h>
#include <unistd.h>

#define CTRL_KEY(c) ((c)&0x1f)

struct termios orig;

void die(const char *s) {
        perror(s);
        exit(1);
}

void disable_input_raw_mode(void) {
        tcsetattr(STDIN_FILENO, TCSAFLUSH, &orig);
}

void enable_input_raw_mode(void) {
        tcgetattr(STDIN_FILENO, &orig);
        atexit(disable_input_raw_mode);

        struct termios raw = orig;
        raw.c_iflag &= ~(BRKINT | ICRNL | INPCK | ISTRIP | IXON);
        raw.c_oflag &= ~(OPOST);
        raw.c_cflag &= ~(CS8);
        raw.c_lflag &= ~(ECHO | ICANON | IEXTEN | ISIG);  // disable echo and canonical mode
        raw.c_cc[VMIN]  = 0;
        raw.c_cc[VTIME] = 1;

        tcsetattr(STDIN_FILENO, TCSAFLUSH, &raw);
}

void raptor_refresh_screen(void) {
        write(STDOUT_FILENO, "\x1b[2J", 4);  // clear screen
        write(STDOUT_FILENO, "\x1b[H", 3);   // reposition cursor to top-left
}

char raptor_read_key(void) {
        char c = '\0';
        int  nread;
        while ((nread = read(STDIN_FILENO, &c, 1)) != 1) {
                if (nread == -1 && errno != EAGAIN) die("read");
        }
        return c;
}

void raptor_process_input(void) {
        char c = raptor_read_key();
        switch (c) {
                case CTRL_KEY('x'):
                        exit(0);
                        break;
        }
}

int main(void) {
        enable_input_raw_mode();

        while (1) {
                raptor_refresh_screen();
                raptor_process_input();
        }

        return 0;
}
