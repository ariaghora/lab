#include "renderer.h"

#include <raylib.h>

void renderer_render(Renderer *r) {
    SetTraceLogLevel(LOG_NONE);
    SetConfigFlags(FLAG_WINDOW_RESIZABLE);

    InitWindow(r->w, r->h, "Inky");
    SetTargetFPS(30);

    while (!WindowShouldClose()) {
        BeginDrawing();

        //---- Drawing section
        ClearBackground(WHITE);
        // TODO
        //----

        EndDrawing();
    }

    CloseWindow();
}
