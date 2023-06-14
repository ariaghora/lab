#include "renderer.h"

#include <math.h>
#include <stdlib.h>

#include "../arr.h"
#include "input.h"

Renderer renderer_init(Editor *e, int w, int h) {
    SetTraceLogLevel(LOG_NONE);
    SetConfigFlags(FLAG_WINDOW_RESIZABLE);

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
    int widths[buf->screen_row_no][buf->screen_col_no];
    int offsets[buf->screen_row_no][buf->screen_col_no];
    memset(widths, 0, sizeof(widths));
    memset(offsets, 0, sizeof(widths));

    size_t file_line_no = filebuf_row_no(buf);
    for (size_t i = buf->offset_y; i < buf->offset_y + buf->screen_row_no; i++) {
        // Draw only rows within screen range
        if (i < file_line_no) {
            char *line = NULL;
            int   high = (int)fmin(fmax(arr_size(buf->rows[i].chars), 0), buf->screen_col_no);

            for (int j = buf->offset_x; j < buf->offset_x + high; j++) {
                // Keep track offset for each character. This is useful to determine
                // the position of the cursor later.
                // The line should be terminated with '\0' before measuring offsets
                char line_term[arr_size(line) + 1];
                memcpy(line_term, line, arr_size(line));
                line_term[arr_size(line)] = '\0';
                offsets[i][j]             = (int)MeasureTextEx(r->editor_font, line_term, r->e->cfg.font_height, 0).x;

                char c        = buf->rows[i].chars[j];
                char c_str[2] = {c, '\0'};
                int  c_w      = MeasureTextEx(r->editor_font, c_str, r->e->cfg.font_height, 0).x;
                switch (c) {
                    case '\t':
                        arr_push_n(line, "     ", char, 4);
                        widths[i][j] = 4 * (c_w);
                        break;
                    default:
                        arr_push(line, c);
                        widths[i][j] = 1 * (c_w);
                        break;
                }
            }

            // We add terminating character just for proper rendering
            arr_push(line, '\0');

            DrawTextEx(r->editor_font, line,
                       (Vector2){buf->offset_x, i * r->e->cfg.line_height},
                       r->e->cfg.font_height, 0, BLACK);
            arr_free(line);
        }
    }

    //---- Drawing cursors
    int cur_w = widths[buf->cur_y + buf->offset_y][buf->cur_x + buf->offset_x];
    if (arr_size(buf->rows[buf->cur_y + buf->offset_y].chars) == 0)
        cur_w = (int)MeasureTextEx(r->editor_font, " ", r->e->cfg.font_height, 0).x;

    int cur_h   = r->e->cfg.line_height;
    int cur_y   = buf->cur_y;
    int cur_off = offsets[buf->cur_y + buf->offset_y][buf->cur_x + buf->offset_x] - offsets[buf->cur_y + buf->offset_y][buf->offset_x];
    DrawRectangleLines(cur_off, cur_y * cur_h, cur_w + r->editor_font.glyphPadding / 2, cur_h, RED);
}

void handle_key_press(Renderer *r) {
    if (IsKeyPressed(KEY_B)) input_handle_key_b(r);
    if (IsKeyPressed(KEY_H)) input_handle_key_h(r);
    if (IsKeyPressed(KEY_J)) input_handle_key_j(r);
    if (IsKeyPressed(KEY_K)) input_handle_key_k(r);
    if (IsKeyPressed(KEY_L)) input_handle_key_l(r);
    if (IsKeyPressed(KEY_W)) input_handle_key_w(r);
}

void renderer_render(Renderer *r) {
    while (!WindowShouldClose()) {
        //---- Update section
        handle_key_press(r);

        //---- Drawing section
        BeginDrawing();
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
