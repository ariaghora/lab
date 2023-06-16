/* STB-style helper dynamic array data structure */

#ifndef _ARR_H_
#define _ARR_H_

#include <ctype.h>
#include <stdlib.h>
#include <string.h>

typedef struct arr_header_t arr_header_t;
struct arr_header_t {
    long size;
    long cap;
};

#define arr_header(a) ((arr_header_t *)((char *)(a) - sizeof(arr_header_t)))
#define arr_size(a) ((a) ? arr_header(a)->size : 0)
#define arr_full(a) ((a) ? (arr_size(a) == arr_header(a)->cap) : 1)
#define arr_push(a, item)                          \
    arr_full(a) ? a = arr_grow(a, sizeof(*a)) : 0, \
                  a[arr_header(a)->size++] = (item)

#define arr_push_n(arr, items, type, n)          \
    do {                                         \
        for (size_t i = 0; i < n; ++i) {         \
            arr_push(arr, ((type *)(items))[i]); \
        }                                        \
    } while (0)

#define arr_insert(arr, idx, item)                           \
    do {                                                     \
        if (arr_full(arr)) {                                 \
            (arr) = arr_grow((arr), sizeof(*(arr)));         \
        }                                                    \
        if ((arr) && (idx) <= arr_size(arr)) {               \
            for (size_t i = arr_size(arr); i > (idx); --i) { \
                (arr)[i] = (arr)[i - 1];                     \
            }                                                \
            (arr)[idx] = (item);                             \
            ++arr_header(arr)->size;                         \
        }                                                    \
    } while (0)

void *arr_del(void *arr, size_t idx, size_t item_size);
void *arr_grow(void *arr, size_t item_size);
void *arr_copy(void *arr, size_t start, size_t end, size_t item_size);
void  arr_free(void *arr);

#endif
