#ifndef _RENDERER_H_
#define _RENDERER_H_

#include "../inky.h"

typedef struct Renderer Renderer;

struct Renderer {
    int     w, h;
    Editor *e;
};

void renderer_render(Renderer *r);

#endif
