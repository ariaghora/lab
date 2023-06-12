#include "arr.h"

void *arr_grow(void *arr, size_t item_size) {
    size_t               new_cap  = arr ? 2 * arr_header(arr)->cap : 1;
    size_t               new_size = sizeof(struct arr_header_t) + new_cap * item_size;
    struct arr_header_t *new_arr  = (arr_header_t *)realloc(arr ? (char *)(arr) - sizeof(struct arr_header_t) : 0, new_size);

    if (new_arr == NULL) {
        return NULL;  // allocation failed
    }

    if (!arr) {
        new_arr->size = 0;
    }

    new_arr->cap = new_cap;
    return (char *)(new_arr) + sizeof(struct arr_header_t);
}

void arr_free(void *arr) {
    if (arr) {
        free((char *)(arr) - sizeof(struct arr_header_t));
    }
}
