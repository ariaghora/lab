#include <stdio.h>

#include "libtinyplot.h"
int main(void) {
        int   n = 5;
        float x[n];
        float y[n];
        for (size_t i = 0; i < n; i++) {
                x[i] = i;
                y[i] = i * i;
        }

        tp_Figure         fig    = tp_new_figure("plot for y=x^2", 800, 600);
        tp_LinePlotConfig lp_cfg = tp_line_plot_default_cfg();
        lp_cfg.width             = 10;
        tp_line_plot(&fig, x, y, 5, &lp_cfg);

        tp_scatter_plot(&fig, x, y, 5);
        tp_render_figure(&fig, 800, 600);

        return 0;
}
