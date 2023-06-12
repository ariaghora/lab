#include <stdio.h>

#include "arr.h"
#include "inky.h"
#include "renderer/renderer.h"

int main(int argc, char** argv) {
    Editor e = editor_init();

    FileBuf buf;
    editor_open_file(&e, "Makefile", &buf);

    Renderer r = (Renderer){
        .e = &e,
        .w = 800,
        .h = 600,
    };
    renderer_render(&r);

    editor_close_all(&e);

    return 0;
}
