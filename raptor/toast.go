package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Toast struct {
	posY  int
	textW int
	textH int

	targetY        int
	durationSecond int
	Message        string

	counter       int
	elapsedSecond int
	shouldDestroy bool

	renderer *sdl.Renderer
	tex      *sdl.Texture
}

func NewToast(message string, durationSecond int, renderer *sdl.Renderer, font *ttf.Font) *Toast {
	sur, _ := font.RenderUTF8Blended(message, sdl.Color{R: 0, G: 0, B: 0, A: 255})
	defer sur.Free()
	tex, _ := renderer.CreateTextureFromSurface(sur)
	posYInint := 10

	return &Toast{
		Message:        message,
		posY:           posYInint,
		durationSecond: durationSecond,
		counter:        0,
		elapsedSecond:  0,
		targetY:        posYInint,
		tex:            tex,
		textW:          int(sur.W),
		textH:          int(sur.H),
		renderer:       renderer,
	}
}

func (t *Toast) MoveY(amount int) {
	t.targetY += amount
}

func (t *Toast) Draw() {
	pad := 10
	tW := t.textW + pad*2
	tH := t.textH + pad*2
	posX := EditorWidth - tW - pad

	t.renderer.SetDrawColor(170, 220, 230, 255)
	t.renderer.FillRect(&sdl.Rect{X: int32(posX), Y: int32(t.posY), W: int32(tW), H: int32(tH)})
	t.renderer.Copy(t.tex,
		nil,
		&sdl.Rect{
			X: int32(posX + (tW / 2) - (t.textW / 2)),
			Y: int32(t.posY + (tH / 2) - (t.textH / 2)),
			W: int32(t.textW),
			H: int32(t.textH),
		},
	)
}

func (t *Toast) Destroy() {
	t.tex.Destroy()
}

func (t *Toast) Update() {
	t.counter += 1
	if t.counter%30 == 0 {
		t.counter = 0
		t.elapsedSecond += 1
	}

	if t.posY < t.targetY {
		t.posY += 10
	}

	if t.elapsedSecond >= t.durationSecond {
		t.shouldDestroy = true
	}
}
