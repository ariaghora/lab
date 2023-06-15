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

void cmd_enter_input_mode(Editor* e, int amount) {
    e->input_mode = INPUT_MODE_INSERT;
}

// moving n characters horizontally
void cmd_inc_cx(Editor* e, int amount) {
    e->active_buf->cur_x += amount;

    char* line = filebuf_current_row(e->active_buf).chars;
    if (e->active_buf->cur_x + e->active_buf->offset_x > arr_size(line) - 1) e->active_buf->cur_x = arr_size(line) - 1;
    if (e->active_buf->cur_x + e->active_buf->offset_x < 0) e->active_buf->cur_x = 0;
}

void cmd_dec_cx(Editor* e, int amount) {
    cmd_inc_cx(e, -amount);
}

// moving n lines vertically
void cmd_inc_cy(Editor* e, int amount) {
    e->active_buf->cur_y += amount;

    if ((e->active_buf->cur_y + e->active_buf->offset_y) < 0) e->active_buf->cur_y = 0;
    if ((e->active_buf->cur_y + e->active_buf->offset_y) > arr_size(e->active_buf->rows) - 1)
        e->active_buf->cur_y = arr_size(e->active_buf->rows) - 1;

    // Prevent cursor exceeding allowed linea area after moving lines
    // above/below. It does not move along x-axis, but it performs cur_x
    // checking and adjust accordingly.
    cmd_inc_cx(e, 0);
}
void cmd_dec_cy(Editor* e, int amount) {
    cmd_inc_cy(e, -amount);
}

void cmd_del_current_line(Editor* e, int amount) {
    int cur_line = e->active_buf->cur_y + e->active_buf->offset_y;
    arr_free(e->active_buf->rows[cur_line].chars);
    arr_del(e->active_buf->rows, cur_line, sizeof(Row));
}

// movie to the next word's start
void cmd_next_word_start(Editor* e, int amount) {
    Word words[1024] = {0};
    int  count       = 0;

    char* line = filebuf_current_row(e->active_buf).chars;
    tokenize_string(line, arr_size(line), punctuations, words, &count);

    for (int i = 0; i < count; i++) {
        if (e->active_buf->cur_x + e->active_buf->offset_x < words[i].start) {
            e->active_buf->cur_x = words[i].start;
            break;
        }
    }

    // Try to move to the next row at following cases as long as cursor is not at the
    // last line of buffer:
    // Case 1: reaching last word of the row,
    // Case 2: empty line
    if (count > 0) {  // Case 1
        if (e->active_buf->cur_x + e->active_buf->offset_x >= words[count - 1].start) {
            if (e->active_buf->cur_y + e->active_buf->offset_y < arr_size(e->active_buf->rows) - 1) {
                e->active_buf->cur_x    = 0;
                e->active_buf->offset_x = 0;
                e->active_buf->cur_y += 1;
                // TODO: We want to skip empty space
            }
        }
    } else if (count == 0) {  // Case 2
        e->active_buf->cur_x    = 0;
        e->active_buf->offset_x = 0;
        e->active_buf->cur_y += 1;
    }
}

// movie to the previous word's start
void cmd_prev_word_start(Editor* e, int amount) {
    Word  words[1024] = {0};
    int   count       = 0;
    char* line        = filebuf_current_row(e->active_buf).chars;
    tokenize_string(line, arr_size(line), punctuations, words, &count);

    for (int i = count - 1; i >= 0; i--) {
        if (e->active_buf->cur_x + e->active_buf->offset_x > words[i].start) {
            e->active_buf->cur_x = words[i].start;
            break;
        }
    }

    // Try to move one line up when finding previous word's start
    int move_up = 0;
    if (count > 0 && (e->active_buf->cur_y + e->active_buf->offset_y > 0)) {
        if (e->active_buf->cur_x + e->active_buf->offset_x <= words[0].end) {
            move_up = 1;
        }
    } else if (count == 0 && (e->active_buf->cur_y + e->active_buf->offset_y > 0)) {
        move_up = 1;
    }

    if (move_up) {
        e->active_buf->cur_y -= 1;
        Word  words[1024] = {0};
        int   count       = 0;
        char* line        = filebuf_current_row(e->active_buf).chars;
        tokenize_string(line, arr_size(line), punctuations, words, &count);
        if (count > 0) {
            tokenize_string(line, arr_size(line), punctuations, words, &count);
            e->active_buf->cur_x = words[count - 1].start;
            // TODO: handle offset x when the last word is beyond screen column
        } else {
            e->active_buf->cur_x    = 0;
            e->active_buf->offset_x = 0;
        }
    }
}
