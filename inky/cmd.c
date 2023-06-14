#include "cmd.h"

#include <stdio.h>

#include "arr.h"

typedef struct Word Word;
struct Word {
    char text[1024];
    int  start, end;
};

char punctuations[] = " `~!@#$%^&*()-_=+[{]}\\|;:\'\",<.>/?]";

//------------------------------------------------------------------------------
// Tokenization helper functions. The functions are used to split a string into
// a collection of words (consecutive alphanumerics or consecutive punctuations).
//------------------------------------------------------------------------------

int is_delimiter(char c, char* delimiters) {
    for (int i = 0; i < strlen(delimiters); i++) {
        if (c == delimiters[i]) return 1;
    }
    return 0;
}

void tokenize_string(char* str, int str_len, char* delimiters, Word* out_words, int* count) {
    if (str_len == 0) return;

    int start        = 0;
    int token_length = 0;

    int idx = 0;
    for (int i = 0; i <= str_len; i++) {
        if (is_delimiter(str[i], delimiters) == is_delimiter(str[start], delimiters) && str[i] != '\0') {
            token_length++;
        } else {
            out_words[idx] = (Word){
                .end   = start + token_length - 1,
                .start = start,
            };
            strncpy(out_words[idx].text, &str[start], token_length);

            // adjust the word start when there's leading white space
            if (!isspace(out_words[idx].text[0]) || (token_length != 1))
                for (int i = 0; i < token_length; i++) {
                    if (isspace(out_words[idx].text[i])) {
                        out_words[idx].start++;
                    } else {
                        break;
                    }
                }

            idx++;
            *count       = idx;
            start        = i;
            token_length = 1;
        }
    }
}

//------------------------------------------------------------------------------
// Editor command implementations
//------------------------------------------------------------------------------

// moving n characters horizontally
void cmd_inc_cx(FileBuf* buf, int amount) {
    buf->cur_x += amount;

    char* line = filebuf_current_row(buf).chars;
    if (buf->cur_x + buf->offset_x > arr_size(line) - 1) buf->cur_x = arr_size(line) - 1;
    if (buf->cur_x + buf->offset_x < 0) buf->cur_x = 0;
}

// moving n lines vertically
void cmd_inc_cy(FileBuf* buf, int amount) {
    buf->cur_y += amount;

    if ((buf->cur_y + buf->offset_y) < 0) buf->cur_y = 0;
    if ((buf->cur_y + buf->offset_y) > arr_size(buf->rows) - 1) buf->cur_y = arr_size(buf->rows) - 1;

    // Prevent cursor exceeding allowed linea area after moving lines
    // above/below. It does not move along x-axis, but it performs cur_x
    // checking and adjust accordingly.
    cmd_inc_cx(buf, 0);
}

// movie to the next word's start
void cmd_next_word_start(FileBuf* buf, int amount) {
    Word words[1024] = {0};
    int  count       = 0;

    char* line = filebuf_current_row(buf).chars;
    tokenize_string(line, arr_size(line), punctuations, words, &count);

    for (int i = 0; i < count; i++) {
        if (buf->cur_x + buf->offset_x < words[i].start) {
            buf->cur_x = words[i].start;
            break;
        }
    }
}

// movie to the previous word's start
void cmd_prev_word_start(FileBuf* buf, int amount) {
    Word words[1024] = {0};
    int  count       = 0;

    char* line = filebuf_current_row(buf).chars;
    tokenize_string(line, arr_size(line), punctuations, words, &count);

    for (int i = count - 1; i >= 0; i--) {
        if (buf->cur_x + buf->offset_x > words[i].start) {
            buf->cur_x = words[i].start;
            break;
        }
    }
}
