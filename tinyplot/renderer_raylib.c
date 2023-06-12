#include <ctype.h>
#include <raylib.h>

#include "libtinyplot.h"

float maxf(float *arr, int n) {
        float cur_max = arr[0];
        for (size_t i = 1; i < n; i++) {
                if (arr[i] > cur_max) cur_max = arr[i];
        }
        return cur_max;
}

static void draw_plot(tp_Plot *plot, int w, int h, int pad) {
        int   wpad = w - 2 * pad;
        int   hpad = h - 2 * pad;
        float xmax = maxf(plot->x, plot->n);
        float ymax = maxf(plot->y, plot->n);

        // prepare the scaling coefficient to convert position unit to pixel.
        // `pad` is already in pixel.
        float xscal = wpad / xmax;
        float yscal = hpad / ymax;

        for (size_t i = 0; i < plot->n; i++) {
                switch (plot->type) {
                        case TP_LINE_PLOT: {
                                if (i > 0) {
                                        int x0_norm = plot->x[i - 1] * xscal + pad;
                                        int x1_norm = plot->x[i] * xscal + pad;
                                        int y0_flip = GetScreenHeight() - (plot->y[i - 1] * yscal) - pad;
                                        int y1_flip = GetScreenHeight() - (plot->y[i] * yscal) - pad;

                                        DrawLineEx((Vector2){x0_norm, y0_flip},
                                                   (Vector2){x1_norm, y1_flip},
                                                   plot->line_plot_cfg->width,
                                                   RED);
                                }
                                break;
                        }
                        case TP_SCATTER_PLOT: {
                                int xnorm = plot->x[i] * xscal + pad;
                                int ynorm = GetScreenHeight() - (plot->y[i] * yscal) - pad;
                                DrawCircle(xnorm, ynorm, 2, GREEN);
                                break;
                        }
                        default:
                                break;
                }
        }
}

void tp_render_figure(tp_Figure *figure, int w, int h) {
        int pad = 40;

        SetConfigFlags(FLAG_MSAA_4X_HINT);
        SetTraceLogLevel(LOG_NONE);
        InitWindow(w, h, figure->title);
        SetTargetFPS(60);
        while (!WindowShouldClose()) {
                BeginDrawing();
                ClearBackground(RAYWHITE);

                for (size_t i = 0; i < figure->n_plot; i++)
                        draw_plot(&figure->plots[i], w, h, pad);

                EndDrawing();
        }
        CloseWindow();
}
