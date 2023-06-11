package main

import "github.com/veandco/go-sdl2/sdl"

func (r *RaptorCfg) DrawSBarVisual() {
	origX, origY := int32(0), int32(EditorHeight-r.SbarHeight)

	r.renderer.SetDrawColor(40, 80, 145, 255)
	r.renderer.FillRect(&sdl.Rect{X: origX, Y: origY, W: EditorWidth, H: int32(r.SbarHeight)})
	s, _ := r.sdlFont.RenderUTF8Blended("Normal", sdl.Color{R: 255, G: 255, B: 255, A: 255})
	t, _ := r.renderer.CreateTextureFromSurface(s)

	r.renderer.Copy(t, nil, &sdl.Rect{
		X: EditorWidth - s.W - 10,
		Y: origY + (int32(r.SbarHeight) / 2) - s.H/2,
		W: s.W,
		H: s.H},
	)
	s.Free()
	t.Destroy()

	if r.CommandBuffer+r.MultiplierBuffer != "" {
		col := r.ColorScheme.Foreground()
		s, _ = r.sdlFont.RenderUTF8Blended(r.MultiplierBuffer+r.CommandBuffer, col)
		t, _ = r.renderer.CreateTextureFromSurface(s)
		r.renderer.Copy(t, nil, &sdl.Rect{
			X: 10,
			Y: origY + (int32(r.SbarHeight) / 2) - s.H/2,
			W: s.W,
			H: s.H},
		)
		s.Free()
		t.Destroy()
	}
}
