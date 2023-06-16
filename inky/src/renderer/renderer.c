#include "renderer.h"

#include <math.h>
#include <stdio.h>
#include <stdlib.h>

#include "../arr.h"
#include "../cmd.h"

//------------------------------------------------------------------------------
// A list mapping from character sequence (in normal mode) to the corresponding
// command.
//------------------------------------------------------------------------------

EditorCommand normal_editor_commands[] = {
    // Navigations
    {.cseq = "h", .cmd_func = cmd_dec_cx},
    {.cseq = "j", .cmd_func = cmd_inc_cy},
    {.cseq = "k", .cmd_func = cmd_dec_cy},
    {.cseq = "l", .cmd_func = cmd_inc_cx},
    {.cseq = "w", .cmd_func = cmd_next_word_start},
    {.cseq = "b", .cmd_func = cmd_prev_word_start},

    // Altering
    {.cseq = "dd", .cmd_func = cmd_del_current_line},
    {.cseq = "dw", .cmd_func = cmd_del_current_line},

    {.cseq = "i", .cmd_func = cmd_enter_input_mode},
};

//------------------------------------------------------------------------------
//
//

Renderer renderer_init(Editor *e, int w, int h) {
    SetTraceLogLevel(LOG_NONE);
    SetConfigFlags(FLAG_WINDOW_RESIZABLE);

    InitWindow(w, h, "Inky");
    SetTargetFPS(30);
    SetExitKey(0);

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

void handle_key_press(Renderer *r) {
    if (r->e->input_mode == INPUT_MODE_NORMAL) {
        char *cmdbuf = r->e->cmd_buf;

        char c = GetCharPressed();
        if (c >= 'a' && c <= 'z')
            arr_push(r->e->cmd_buf, c);

        // find command candidates to execture
        EditorCommand candidate_cmd[sizeof(normal_editor_commands) / sizeof(normal_editor_commands[0])];
        int           cnt = 0;
        for (size_t i = 0; i < sizeof(normal_editor_commands) / sizeof(normal_editor_commands[0]); i++) {
            if (arr_size(cmdbuf) > 0) {
                if (strnstr(normal_editor_commands[i].cseq, cmdbuf, arr_size(cmdbuf))) {
                    candidate_cmd[cnt++] = normal_editor_commands[i];
                }
            }
        }
        // We only execute the command if only match found. Then clear the command
        // buffer. Otherwise, if we have something on command buffer (and mul buffer),
        // but no match found, then the command is unknown, so clear the cmd buffer.
        if (cnt == 1) {
            candidate_cmd[0].cmd_func(r->e, 1);
            editor_clear_cmd_buffer(r->e);
        } else if (cnt == 0 && (arr_size(cmdbuf) + arr_size(r->e->mul_buf) > 0)) {
            editor_clear_cmd_buffer(r->e);
        }
    } else if (r->e->input_mode == INPUT_MODE_INSERT) {
        char c        = GetCharPressed();
        int  line_idx = r->e->active_buf->cur_y + r->e->active_buf->offset_y;
        int  col_idx  = r->e->active_buf->cur_x + r->e->active_buf->offset_x;

        if (isprint(c)) {
            arr_insert(r->e->active_buf->rows[line_idx].chars, col_idx, c);
            ++r->e->active_buf->cur_x;
        } else {
            // Handle non-ascii keys
            int o = GetKeyPressed();
            if (o == KEY_ENTER) {
                cmd_line_break(r->e, 1);
            }
        }
    }

    // From whatever mode, when user presses escape key, we will go back to the
    // normal mode
    if (IsKeyPressed(KEY_ESCAPE)) {
        r->e->input_mode = INPUT_MODE_NORMAL;
        editor_clear_cmd_buffer(r->e);
    }
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
                       r->e->cfg.font_height, 0, LIGHTGRAY);
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

void renderer_render_statusbar(Renderer *r) {
    int   bar_h    = 32;
    int   bar_x    = 0;
    int   bar_rpad = 10;
    int   bar_y    = GetScreenHeight() - bar_h;
    Color bar_color;
    char  mode_text[30];
    char  cmd_text[30];

    switch (r->e->input_mode) {
        case INPUT_MODE_NORMAL: {
            memcpy(mode_text, "Normal", 6);
            bar_color = (Color){255, 0, 0, 255};

            memcpy(cmd_text, r->e->cmd_buf, arr_size(r->e->cmd_buf));
            cmd_text[arr_size(r->e->cmd_buf)] = 0;
            break;
        }
        case INPUT_MODE_INSERT: {
            memcpy(mode_text, "Insert", 6);
            bar_color = (Color){0, 0, 255, 255};
            break;
        }
    }

    DrawRectangle(bar_x, bar_y, GetScreenWidth(), bar_h, bar_color);

    int w = MeasureTextEx(r->editor_font, mode_text, r->e->cfg.font_height, 0).x;
    DrawTextEx(r->editor_font, mode_text,
               (Vector2){GetScreenWidth() - w - bar_rpad, bar_y + (bar_h / 2) - (r->e->cfg.font_height / 2)},
               r->e->cfg.font_height, 0, WHITE);

    // rendering command buffer's command
    DrawTextEx(r->editor_font, cmd_text,
               (Vector2){0, bar_y + (bar_h / 2) - (r->e->cfg.font_height / 2)},
               r->e->cfg.font_height, 0, GREEN);
}

void renderer_render(Renderer *r) {
    while (!WindowShouldClose()) {
        //---- Update section
        handle_key_press(r);

        //---- Drawing section
        BeginDrawing();
        ClearBackground(DARKGRAY);

        // Render active buffer
        FileBuf *buf = r->e->active_buf;
        renderer_render_active_buffer(r, buf);
        renderer_render_statusbar(r);
        // TODO: render all buffers with a certain strategy
        //----

        EndDrawing();
    }

    CloseWindow();
    UnloadFont(r->editor_font);
}
