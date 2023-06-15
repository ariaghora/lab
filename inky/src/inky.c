#include "inky.h"

#include <stdlib.h>

#include "arr.h"
#include "cmd.h"

EditorConfig editor_config_default(void) {
    return (EditorConfig){
        .font_height = 20,
        .line_height = 20,
    };
}

Editor editor_init(void) {
    return (Editor){
        .file_buff_no = 0,
        .cfg          = editor_config_default(),
        .active_buf   = NULL,
        .input_mode   = INPUT_MODE_NORMAL,
        .cmd_buf      = NULL,
        .mul_buf      = NULL,
    };
}

void editor_open_file(Editor* e, const char* file_name, FileBuf* buf) {
    // load into buffer and append buffer to the editor
    filebuf_from_file(buf, file_name);
    e->file_buf[e->file_buff_no++] = buf;
    e->active_buf                  = buf;
}

void editor_close_all(Editor* e) {
    for (size_t i = 0; i < e->file_buff_no; i++) {
        filebuf_close(e->file_buf[i]);
    }
}

void editor_clear_cmd_buffer(Editor* e) {
    arr_free(e->cmd_buf);
    arr_free(e->mul_buf);
    e->cmd_buf = NULL;
    e->mul_buf = NULL;
}
