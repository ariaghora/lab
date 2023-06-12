#include "filebuf.h"

#include <stdio.h>

#include "arr.h"

void filebuf_from_file(FileBuf *buf, const char *file_name) {
    FILE *fp = fopen(file_name, "rw");
    if (fp == NULL)
        exit(EXIT_FAILURE);

    char   *line = NULL;
    size_t  len  = 0;
    ssize_t read;

    buf->rows = NULL;

    while ((read = getline(&line, &len, fp)) != -1) {
        // Set current row's text to `line`, along with terminating char
        Row row = (Row){NULL};
        arr_push_n(row.chars, line, char, read - 1);
        arr_push(row.chars, '\0');

        // Add this row to the buffer
        arr_push(buf->rows, row);
    }

    fclose(fp);

    if (line) free(line);
}

void filebuf_close(FileBuf *buf) {
    for (size_t i = 0; i < arr_size(buf->rows); i++) {
        arr_free(buf->rows[i].chars);
    }
    arr_free(buf->rows);
}

size_t filebuf_row_no(FileBuf *buf) {
    return arr_size(buf->rows);
}
