#include "input.h"

#include "../cmd.h"
#include "../filebuf.h"

void input_handle_key_l(Renderer *r) {
    FileBuf *active_buffer = r->e->active_buf;
    cmd_inc_cx(active_buffer, 1);
}

void input_handle_key_h(Renderer *r) {
    FileBuf *active_buffer = r->e->active_buf;
    cmd_inc_cx(active_buffer, -1);
}

void input_handle_key_j(Renderer *r) {
    FileBuf *active_buffer = r->e->active_buf;
    cmd_inc_cy(active_buffer, 1);
}

void input_handle_key_k(Renderer *r) {
    FileBuf *active_buffer = r->e->active_buf;
    cmd_inc_cy(active_buffer, -1);
}
