#include "arr.h"

void *arr_grow(void *arr, size_t item_size) {
    size_t               new_cap  = arr ? 2 * arr_header(arr)->cap : 1;
    size_t               new_size = sizeof(struct arr_header_t) + new_cap * item_size;
    struct arr_header_t *new_arr  = (arr_header_t *)realloc(arr ? (char *)(arr) - sizeof(struct arr_header_t) : 0, new_size);

    if (new_arr == NULL) return NULL;  // allocation failed

    if (!arr) new_arr->size = 0;

    new_arr->cap = new_cap;
    return (char *)(new_arr) + sizeof(struct arr_header_t);
}

void *arr_copy(void *arr, size_t start, size_t end, size_t item_size) {
    if (!arr || start >= arr_size(arr) || end > arr_size(arr) || start > end) {
        return NULL;  // Array is NULL or indexes are out of bounds.
    }

    // Calculate number of elements to copy and allocate new array.
    size_t n = end - start;
    if (n == 0) return NULL;

    void *new_arr = malloc(sizeof(struct arr_header_t) + n * item_size);
    if (!new_arr) {
        return NULL;  // Allocation failed.
    }

    // Set up the header for the new array.
    struct arr_header_t *header = (struct arr_header_t *)new_arr;
    header->size                = n;
    header->cap                 = n;
    new_arr                     = (char *)new_arr + sizeof(struct arr_header_t);

    // Copy items from the source array to the new array.
    memcpy(new_arr, (char *)arr + start * item_size, n * item_size);

    return new_arr;
}

void *arr_del(void *arr, size_t idx, size_t item_size) {
    if (!arr || idx >= arr_size(arr)) {
        return NULL;  // Array is NULL or index is out of bounds.
    }

    // Calculate pointers to the item to delete and the item after it.
    char *item_to_del = (char *)(arr) + idx * item_size;
    char *next_item   = item_to_del + item_size;

    // Shift all items after the one to delete one place to the left.
    memmove(item_to_del, next_item, (arr_size(arr) - idx - 1) * item_size);

    // Decrement the size of the array.
    --arr_header(arr)->size;

    return arr;
}

void arr_free(void *arr) {
    if (arr) free((char *)(arr) - sizeof(struct arr_header_t));
}
