#ifndef _FILE_BUF_
#define _FILE_BUF_

#include <ctype.h>

typedef struct FileBuf FileBuf;
typedef struct Row     Row;

struct Row {
    char *chars;
};

struct FileBuf {
    size_t cur_x, cur_y;
    size_t offset_x, offset_y;
    Row   *rows;
};

void   filebuf_from_file(FileBuf *buf, const char *file_name);
size_t filebuf_row_no(FileBuf *buf);

#endif
