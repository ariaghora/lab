#ifndef _INKY_H_
#define _INKY_H_

#define MAX_FILE_BUF_NO 10

#include "filebuf.h"

typedef struct EditorConfig EditorConfig;
typedef struct Editor       Editor;

//------------------------------------------------------------------------------
// The API
//------------------------------------------------------------------------------

typedef enum { INPUT_MODE_NORMAL,
               INPUT_MODE_INSERT } EditorInputMode;

struct EditorConfig {
    int font_height;
    int line_height;
};

// Main structure that holds all editor's state. It enables loading multiple
// files (i.e., each of which into a buffer).
struct Editor {
    int             file_buff_no;
    FileBuf*        file_buf[MAX_FILE_BUF_NO];
    FileBuf*        active_buf;
    EditorConfig    cfg;
    EditorInputMode input_mode;

    char* cmd_buf;  // command buffer
    char* mul_buf;  // multiplier buffer
};

// Initialize default editor config
EditorConfig editor_config_default(void);

Editor editor_init(void);
void   editor_open_file(Editor* e, const char* file_name, FileBuf* buf);
void   editor_close_all(Editor* e);
void   editor_clear_cmd_buffer(Editor* e);

#endif
