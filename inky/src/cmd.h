#ifndef _CMD_H_
#define _CMD_H_

#include "inky.h"

typedef void(CmdFunc)(Editor *, int);
typedef struct EditorCommand EditorCommand;
struct EditorCommand {
    char     cseq[5];
    CmdFunc *cmd_func;
};

void cmd_enter_input_mode(Editor *e, int amount);

void cmd_inc_cx(Editor *e, int amount);
void cmd_dec_cx(Editor *e, int amount);
void cmd_inc_cy(Editor *e, int amount);
void cmd_dec_cy(Editor *e, int amount);
void cmd_del_current_line(Editor *e, int amount);
void cmd_next_word_start(Editor *e, int amount);
void cmd_prev_word_start(Editor *e, int amount);

#endif
