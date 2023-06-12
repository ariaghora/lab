#include <stdio.h>

#include "arr.h"
#include "inky.h"

int main(int argc, char** argv) {
    Editor  e = editor_init();
    FileBuf buf;

    editor_open_file(&e, "Makefile", &buf);

    size_t file_row_no = filebuf_row_no(&buf);
    printf("%ld\n", file_row_no);

    for (size_t i = 0; i < file_row_no; i++) {
        printf("%s\n", buf.rows[i].chars);
    }

    return 0;
}
