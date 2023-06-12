#include "inky.h"

Editor editor_init(void) {
    return (Editor){
        .file_buff_no = 0,
    };
}

void editor_open_file(Editor* e, const char* file_name, FileBuf* buf) {
    filebuf_from_file(buf, file_name);
    e->file_buf[e->file_buff_no++] = buf;
}
