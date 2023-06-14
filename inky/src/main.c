#include <stdio.h>

#include "arr.h"
#include "inky.h"
#include "renderer/renderer.h"

int main(int argc, char** argv) {
    Editor  e   = editor_init();
    FileBuf buf = {
        .cur_x         = 0,
        .cur_y         = 0,
        .offset_x      = 0,
        .offset_y      = 0,
        .rows          = NULL,
        .screen_row_no = 20,
        .screen_col_no = 80,
    };
    editor_open_file(&e, "Makefile", &buf);

    Renderer r = renderer_init(&e, 1000, 600);
    renderer_render(&r);

    editor_close_all(&e);
    return 0;
}
