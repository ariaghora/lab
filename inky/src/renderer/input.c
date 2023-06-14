#include "input.h"

#include "../cmd.h"
#include "../filebuf.h"

void input_handle_key_b(Renderer *r) {
    if (r->e->input_mode == INPUT_MODE_NORMAL)
        cmd_prev_word_start(r->e->active_buf, -1);
}

void input_handle_key_h(Renderer *r) {
    if (r->e->input_mode == INPUT_MODE_NORMAL)
        cmd_inc_cx(r->e->active_buf, -1);
}

void input_handle_key_i(Renderer *r) {
    if (r->e->input_mode == INPUT_MODE_NORMAL)
        r->e->input_mode = INPUT_MODE_INSERT;
}

void input_handle_key_j(Renderer *r) {
    if (r->e->input_mode == INPUT_MODE_NORMAL)
        cmd_inc_cy(r->e->active_buf, 1);
}

void input_handle_key_k(Renderer *r) {
    if (r->e->input_mode == INPUT_MODE_NORMAL)
        cmd_inc_cy(r->e->active_buf, -1);
}

void input_handle_key_l(Renderer *r) {
    if (r->e->input_mode == INPUT_MODE_NORMAL)
        cmd_inc_cx(r->e->active_buf, 1);
}

void input_handle_key_w(Renderer *r) {
    if (r->e->input_mode == INPUT_MODE_NORMAL)
        cmd_next_word_start(r->e->active_buf, 1);
}

void input_handle_key_esc(Renderer *r) {
    r->e->input_mode = INPUT_MODE_NORMAL;
}
