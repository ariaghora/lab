package main

import (
	"bufio"
	"container/list"
	"os"
	"strconv"

	"github.com/rivo/uniseg"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const EditorWidth = 1000
const EditorHeight = 600

type Row struct {
	Chars string
}

type EditorMode int64

const (
	EditorModeNormal EditorMode = iota
	EditorModeInsert
	EditorModeSelect
)

type InsertMethod int64

const (
	InsertMethodInsert InsertMethod = iota
	InsertMethodAppend
	InsertMethodBreakLine
)

type Renderer interface {
	Draw()
	Update()
}

type RaptorCfg struct {
	CX               int
	CY               int
	EditorMode       EditorMode
	EditorWidth      int
	EditorHeight     int
	FileName         string
	FontSize         int
	LastInsertMethod InsertMethod
	LineHeight       int
	ScreenCols       int
	ScreenRows       int
	NumRows          int
	RowOffset        int
	ColOffset        int
	Rows             []Row
	SbarHeight       int
	ShiftMapper      map[string]string

	Toasts list.List

	sdlFont   *ttf.Font
	sdlWindow *sdl.Window
	renderer  *sdl.Renderer
}

func EstimateCharWidth(font *ttf.Font) int {
	s, _ := font.RenderUTF8Blended(" ", sdl.Color{})
	return int(s.W)
}

func NewEditor() *RaptorCfg {
	window, err := sdl.CreateWindow(
		"raptor - untitled", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		EditorWidth, EditorHeight, sdl.WINDOW_SHOWN,
	)

	if err != nil {
		panic(err)
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	// defer rendererRef.Destroy()
	renderer.SetViewport(&sdl.Rect{X: 0, Y: 0, W: EditorWidth, H: EditorHeight})

	fontSize := 14
	lineHeight := 22
	sbarHeight := 24
	font, err := ttf.OpenFont("IosevkaNerdFontMono-Medium.ttf", fontSize)
	if err != nil {
		panic(err)
	}
	w, h := window.GetSize()

	cfg := &RaptorCfg{
		CX:               0,
		CY:               0,
		RowOffset:        0,
		ColOffset:        0,
		EditorWidth:      int(w),
		EditorHeight:     int(h),
		FontSize:         fontSize,
		LastInsertMethod: InsertMethodInsert,
		LineHeight:       lineHeight,
		ScreenCols:       80,
		ScreenRows:       (EditorHeight - sbarHeight) / lineHeight,
		NumRows:          0,
		Rows:             []Row{},
		SbarHeight:       sbarHeight,
		ShiftMapper: map[string]string{
			"`":  "~",
			"-":  "_",
			"=":  "+",
			"[":  "{",
			"]":  "}",
			"\\": "|",
			";":  ":",
			"'":  "\"",
			",":  "<",
			".":  ">",
			"/":  "?",
		},
		Toasts:    *list.New(),
		sdlFont:   font,
		sdlWindow: window,
		renderer:  renderer,
	}
	return cfg
}

func (r *RaptorCfg) OpenFile(fileName string) error {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewScanner(file)
	for reader.Scan() {
		line := reader.Text()
		r.Rows = append(r.Rows, Row{
			Chars: line,
		})
		r.NumRows += 1
	}

	if r.NumRows == 0 {
		r.Rows = append(r.Rows, Row{""})
		r.NumRows += 1
	}

	r.sdlWindow.SetTitle(fileName)
	r.FileName = fileName
	return nil
}

func (r *RaptorCfg) HandleKeyPress(ev *sdl.KeyboardEvent) {
	if r.EditorMode == EditorModeNormal {
		switch ev.Keysym.Scancode {
		// entering insert mode
		case sdl.SCANCODE_I, sdl.SCANCODE_A, sdl.SCANCODE_O:
			r.HandleEnteringInsertMode(ev)

		// saving file
		case sdl.SCANCODE_S:
			r.HandleBufferWriteToFile(ev)

		// Navigation
		case sdl.SCANCODE_H:
			r.HandleKeyPressH()
		case sdl.SCANCODE_J:
			r.HandleKeyPressJ()
		case sdl.SCANCODE_K:
			r.HandleKeyPressK()
		case sdl.SCANCODE_L:
			r.HandleKeyPressL()
		case sdl.SCANCODE_W:
			r.JumpToNextWordBeginning()
		case sdl.SCANCODE_B:
			r.JumpToPrevWordBeginning()
		}
	} else if r.EditorMode == EditorModeInsert {
		r.HandleKeyPressInsertMode(ev)
	}
}

func (r *RaptorCfg) SwitchToNormalMode() {
}

func (r *RaptorCfg) Run() {
	running := true
	for running {
		// events
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				if t.State == sdl.PRESSED {
					if t.Keysym.Scancode == sdl.SCANCODE_ESCAPE {
						if r.CX > 0 && r.EditorMode == EditorModeInsert {
							r.CX -= 1
						}
						r.EditorMode = EditorModeNormal
						r.LastInsertMethod = InsertMethodInsert
					}

					r.HandleKeyPress(t)
				}
			}
		}

		// draw
		r.DrawScreen()
		sdl.Delay(33)
	}
}

func (r *RaptorCfg) Destroy() {
	r.sdlWindow.Destroy()
	r.renderer.Destroy()
}

func (r *RaptorCfg) DrawSBar() {
	if r.EditorMode == EditorModeNormal {
		r.DrawSBarVisual()
	} else if r.EditorMode == EditorModeInsert {
		r.DrawSBarInsert()
	}
}

func (r *RaptorCfg) DrawScreen() {
	r.renderer.Clear()

	// Text background
	r.renderer.SetDrawColor(40, 40, 40, 255)
	r.renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: int32(r.EditorWidth), H: int32(r.EditorHeight)})

	// Line number area
	screenMaxLineNo := strconv.Itoa(r.RowOffset + r.ScreenRows)
	surf, _ := r.sdlFont.RenderUTF8Blended(screenMaxLineNo, sdl.Color{})
	lineNoColWidth := surf.W + 20 + 20
	r.renderer.SetDrawColor(20, 20, 20, 255)
	r.renderer.FillRect(&sdl.Rect{
		X: 0, Y: 0, W: int32(lineNoColWidth), H: int32(r.EditorHeight),
	})
	for y := 0; y < r.ScreenRows; y += 1 {
		if y < r.NumRows {
			numStr := strconv.Itoa(y + r.RowOffset + 1)
			s, _ := r.sdlFont.RenderUTF8Blended(numStr, sdl.Color{R: 255, G: 255, B: 255, A: 255})
			t, _ := r.renderer.CreateTextureFromSurface(s)

			r.renderer.Copy(t, nil, &sdl.Rect{X: 20, Y: int32(y * r.LineHeight), W: s.W, H: s.H})
			s.Free()
			t.Destroy()
		}
	}

	// Texts
	for y := r.RowOffset; y < r.RowOffset+r.ScreenRows; y += 1 {
		if y < r.NumRows {
			r.renderBufferText(
				int(lineNoColWidth+8),
				(y-r.RowOffset)*r.LineHeight,
				r.Rows[y].Chars,
				r.LineHeight,
				r.renderer)
		}
	}

	// Statusbar
	r.DrawSBar()

	// Toast
	for e := r.Toasts.Front(); e != nil; e = e.Next() {
		if e.Value == nil {
			continue
		}
		next := e.Next()
		t := e.Value.(*Toast)
		t.Update()
		t.Draw()

		if t.shouldDestroy {
			t.Destroy()
			r.Toasts.Remove(e)
		}
		e = next
	}

	r.renderer.Present()

}

func (r *RaptorCfg) renderBufferText(x int, y int, text string, lineHeight int, renderer *sdl.Renderer) {
	charOffsetX := x
	offsetY := y

	relY := y / r.LineHeight
	relX := 0

	gr := uniseg.NewGraphemes(text)
	for gr.Next() {
		c := gr.Str()
		if c == "\t" {
			c = "    "
		}
		var s *sdl.Surface
		if r.CX == relX && r.CY == relY {
			s, _ = r.sdlFont.RenderUTF8Blended(c, sdl.Color{R: 0, G: 0, B: 0, A: 255})
			renderer.SetDrawColor(255, 255, 255, 255)
			renderer.FillRect(&sdl.Rect{X: int32(charOffsetX + r.ColOffset), Y: int32(y), W: s.W, H: int32(r.LineHeight)})
		} else {
			s, _ = r.sdlFont.RenderUTF8Blended(c, sdl.Color{R: 255, G: 255, B: 255, A: 255})
		}
		t, _ := renderer.CreateTextureFromSurface(s)
		t.SetBlendMode(sdl.BLENDMODE_BLEND)
		renderer.Copy(t, nil, &sdl.Rect{X: int32(charOffsetX + r.ColOffset), Y: int32(offsetY), W: s.W, H: s.H})
		charOffsetX += int(s.W)
		relX += 1

		s.Free()
		t.Destroy()
	}
	// Hack: when line is an empty string, no text to highlight. Then, draw a dummy cursor
	if r.CX == relX && r.CY == relY && len(text) == 0 {
		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.FillRect(&sdl.Rect{
			X: int32(charOffsetX + r.ColOffset),
			Y: int32(y),
			W: int32(EstimateCharWidth(r.sdlFont)),
			H: int32(r.LineHeight),
		})
	}
}

func main() {
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	if err := ttf.Init(); err != nil {
		panic(err)
	}

	cfg := NewEditor()
	defer cfg.Destroy()

	cfg.Toasts.PushFront(NewToast("Editing "+"README.md", 2, cfg.renderer, cfg.sdlFont))

	cfg.OpenFile("README.md")
	cfg.Run()

}
