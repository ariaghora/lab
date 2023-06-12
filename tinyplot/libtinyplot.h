#ifndef _TINYPLOT_H_
#define _TINYPLOT_H_

#define MAX_PLOT 10

typedef enum { TP_LINE_PLOT,
               TP_SCATTER_PLOT } tp_PlotType;
typedef enum { TP_LINE_SOLID,
               TP_LINE_DASHED } tp_LinePlotType;

typedef struct tp_Figure         tp_Figure;
typedef struct tp_Plot           tp_Plot;
typedef struct tp_PixelData      tp_PixelData;
typedef struct tp_LinePlotConfig tp_LinePlotConfig;

struct tp_LinePlotConfig {
        int             width;
        tp_LinePlotType line_type;
};

struct tp_PixelData {
        int r, g, b, a;
};

struct tp_Plot {
        tp_PlotType    type;
        float         *x;
        float         *y;
        int            n;
        tp_PixelData **px;

        tp_LinePlotConfig *line_plot_cfg;
};

struct tp_Figure {
        const char *title;
        int         n_plot;
        tp_Plot     plots[MAX_PLOT];
        int         w;
        int         h;
};

// API
tp_Figure tp_new_figure(const char *, int, int);

void              tp_line_plot(tp_Figure *fig, float *x, float *y, int n, tp_LinePlotConfig *cfg);
tp_LinePlotConfig tp_line_plot_default_cfg(void);

void tp_scatter_plot(tp_Figure *fig, float *x, float *y, int n);

// should be linked with one of rendering backend
void tp_render_figure(tp_Figure *figure, int w, int h);

#endif
