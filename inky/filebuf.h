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
    size_t cur_x, cur_y;
    size_t offset_x, offset_y;
    size_t screen_row_no;
    Row   *rows;
};

// Read file content line by line and store them as buffer's rows
void filebuf_from_file(FileBuf *buf, const char *file_name);
void filebuf_close(FileBuf *buf);

// A helper to get number of line of the file loaded in a buffer
size_t filebuf_row_no(FileBuf *buf);

#endif
