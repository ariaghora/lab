package main

import "github.com/veandco/go-sdl2/sdl"

func (r *RaptorCfg) DrawSBarInsert() {
	origX, origY := int32(0), int32(EditorHeight-r.SbarHeight)

	r.renderer.SetDrawColor(145, 36, 80, 255)
	r.renderer.FillRect(&sdl.Rect{X: origX, Y: origY, W: EditorWidth, H: int32(r.SbarHeight)})
	s, _ := r.sdlFont.RenderUTF8Blended("Insert", sdl.Color{R: 255, G: 255, B: 255, A: 255})
	defer s.Free()

	t, _ := r.renderer.CreateTextureFromSurface(s)
	defer t.Destroy()

	r.renderer.Copy(t, nil, &sdl.Rect{
		X: EditorWidth - s.W - 10,
		Y: origY + (int32(r.SbarHeight) / 2) - s.H/2,
		W: s.W,
		H: s.H},
	)
}
