#include "cmd.h"

#include "arr.h"

void cmd_inc_cx(FileBuf *buf, int amount) {
    buf->cur_x += amount;

    char *line = filebuf_current_row(buf).chars;
    if (buf->cur_x + buf->offset_x > arr_size(line) - 1) buf->cur_x = arr_size(line) - 1;
    if (buf->cur_x + buf->offset_x < 0) buf->cur_x = 0;
}

void cmd_inc_cy(FileBuf *buf, int amount) {
    buf->cur_y += amount;

    if ((buf->cur_y + buf->offset_y) < 0) buf->cur_y = 0;
    if ((buf->cur_y + buf->offset_y) > arr_size(buf->rows) - 1) buf->cur_y = arr_size(buf->rows) - 1;

    // Prevent cursor exceeding allowed linea area after moving lines
    // above/below. It does not move along x-axis, but it performs cur_x
    // checking and adjust accordingly.
    cmd_inc_cx(buf, 0);
}
