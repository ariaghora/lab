#ifndef _INKY_H_
#define _INKY_H_

#define MAX_FILE_BUF_NO 10

#include "filebuf.h"

typedef struct Editor Editor;

struct Editor {
    int      file_buff_no;
    FileBuf* file_buf[MAX_FILE_BUF_NO];
};

Editor editor_init(void);
void   editor_open_file(Editor* e, const char* file_name, FileBuf* buf);

#endif
