#include "renderer.h"

#include <stdlib.h>

#include "../arr.h"

Renderer renderer_init(Editor *e, int w, int h) {
    SetTraceLogLevel(LOG_NONE);
    SetConfigFlags(FLAG_WINDOW_RESIZABLE);
    SetConfigFlags(FLAG_WINDOW_HIGHDPI);

    InitWindow(w, h, "Inky");
    SetTargetFPS(30);

    // Set up editor font
    Font f = LoadFont("IosevkaNerdFont-Regular.ttf");
    SetTextureFilter(f.texture, TEXTURE_FILTER_TRILINEAR);

    Renderer r = (Renderer){
        .e           = e,
        .w           = w,
        .h           = h,
        .editor_font = f,
    };
    return r;
}

void renderer_render_active_buffer(Renderer *r, FileBuf *buf) {
    size_t file_line_no = filebuf_row_no(buf);
    for (size_t i = buf->offset_y; i < buf->offset_y + buf->screen_row_no; i++) {
        // Draw only rows within screen range
        if (i < file_line_no) {
            char *line = NULL;
            for (size_t j = 0; j < arr_size(buf->rows[i].chars); j++) {
                char c = buf->rows[i].chars[j];
                switch (c) {
                    case '\t':
                        arr_push_n(line, "     ", char, 4);
                        break;
                    default:
                        arr_push(line, c);
                        break;
                }
            }

            DrawTextEx(r->editor_font, line, (Vector2){buf->offset_x, i * r->e->cfg.font_height}, r->e->cfg.font_height, 0, BLACK);
            arr_free(line);
        }
    }
}

void renderer_render(Renderer *r) {
    while (!WindowShouldClose()) {
        BeginDrawing();

        //---- Drawing section
        ClearBackground(WHITE);
        // Render active buffer
        FileBuf *buf = r->e->active_buf;
        renderer_render_active_buffer(r, buf);
        // TODO: render all buffers with a certain strategy
        //----

        EndDrawing();
    }

    CloseWindow();
}
