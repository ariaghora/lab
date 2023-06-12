#ifndef _INKY_H_
#define _INKY_H_

#define MAX_FILE_BUF_NO 10

#include "filebuf.h"

typedef struct Editor Editor;

//------------------------------------------------------------------------------
// The API
//------------------------------------------------------------------------------

// Main structure that holds all editor's state. It enables loading multiple
// files (i.e., each of which into a buffer).
struct Editor {
    int      file_buff_no;
    FileBuf* file_buf[MAX_FILE_BUF_NO];

    int font_height;
    int line_height;
};

// Initialize an editor instance. There should be only one instance of editor
// in a project.
Editor editor_init(void);

// Open a file into a buffer, then append that buffer into buffer list in the
// editor.
void editor_open_file(Editor* e, const char* file_name, FileBuf* buf);
void editor_close_all(Editor* e);

#endif
