#include "libtinyplot.h"

#include <stdio.h>
#include <stdlib.h>
#include <string.h>

typedef struct arr_header_t arr_header_t;
struct arr_header_t {
        size_t size;
        size_t cap;
};

/* Helper dynamic array data structure */

#define arr_header(a) ((arr_header_t *)((char *)(a) - sizeof(arr_header_t)))
#define arr_size(a) ((a) ? arr_header(a)->size : 0)
#define arr_full(a) ((a) ? (arr_size(a) == arr_header(a)->cap) : 1)
#define arr_push(a, item)                              \
        arr_full(a) ? a = arr_grow(a, sizeof(*a)) : 0, \
                      a[arr_header(a)->size++] = item

void *arr_grow(void *arr, size_t item_size) {
        size_t               new_cap  = arr ? 2 * arr_header(arr)->cap : 1;
        size_t               new_size = sizeof(struct arr_header_t) + new_cap * item_size;
        struct arr_header_t *new_arr  = realloc(arr ? (char *)(arr) - sizeof(struct arr_header_t) : 0, new_size);

        if (new_arr == NULL) {
                return NULL;  // allocation failed
        }

        if (!arr) {
                new_arr->size = 0;
        }

        new_arr->cap = new_cap;
        return (char *)(new_arr) + sizeof(struct arr_header_t);
}

#define arr_push_n(arr, items, type, n)                      \
        do {                                                 \
                for (size_t i = 0; i < n; ++i) {             \
                        arr_push(arr, ((type *)(items))[i]); \
                }                                            \
        } while (0)

void arr_free(void *arr) {
        if (arr) {
                free((char *)(arr) - sizeof(struct arr_header_t));
        }
}

tp_Figure tp_new_figure(const char *title, int w, int h) {
        return (tp_Figure){
            .title  = title,
            .n_plot = 0,
            .w      = w,
            .h      = h,
        };
}

tp_LinePlotConfig tp_line_plot_default_cfg(void) {
        return (tp_LinePlotConfig){
            .line_type = TP_LINE_SOLID,
            .width     = 1,
        };
}

void tp_line_plot(tp_Figure *fig, float *x, float *y, int n, tp_LinePlotConfig *cfg) {
        tp_Plot lp = {
            .type          = TP_LINE_PLOT,
            .x             = x,
            .y             = y,
            .n             = n,
            .line_plot_cfg = cfg,
        };
        fig->plots[fig->n_plot++] = lp;
}

void tp_scatter_plot(tp_Figure *fig, float *x, float *y, int n) {
        tp_Plot lp = {
            .type          = TP_SCATTER_PLOT,
            .x             = x,
            .y             = y,
            .n             = n,
            .line_plot_cfg = NULL,
        };
        fig->plots[fig->n_plot++] = lp;
}
