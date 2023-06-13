#ifndef _FILE_BUF_
#define _FILE_BUF_

#include <ctype.h>

typedef struct FileBuf FileBuf;
typedef struct Row     Row;

struct Row {
    char *chars;
};

// Managing content and editing states for a single file
struct FileBuf {
    int  cur_x, cur_y;
    int  offset_x, offset_y;
    int  screen_row_no;
    int  screen_col_no;
    Row *rows;
};

void filebuf_from_file(FileBuf *buf, const char *file_name);
void filebuf_close(FileBuf *buf);
Row  filebuf_current_row(FileBuf *buf);
int  filebuf_row_no(FileBuf *buf);

#endif
