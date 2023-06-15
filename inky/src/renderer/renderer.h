#ifndef _RENDERER_H_
#define _RENDERER_H_

#include <raylib.h>

#include "../inky.h"

typedef struct Renderer   Renderer;
typedef struct LayoutTree LayoutTree;
typedef struct Viewport   Viewport;

struct Viewport {
    int x, y, w, h;
};

struct Renderer {
    int     w, h;
    Editor *e;
    Font    editor_font;
};

Renderer renderer_init(Editor *e, int w, int h);
void     renderer_render(Renderer *r);

#endif
