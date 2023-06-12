#include "inky.h"

Editor editor_init(void) {
    return (Editor){
        .file_buff_no = 0,
    };
}

void editor_open_file(Editor* e, const char* file_name, FileBuf* buf) {
    // load into buffer and append buffer to the editor
    filebuf_from_file(buf, file_name);
    e->file_buf[e->file_buff_no++] = buf;
}

void editor_close_all(Editor* e) {
    for (size_t i = 0; i < e->file_buff_no; i++) {
        filebuf_close(e->file_buf[i]);
    }
}
