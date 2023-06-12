package main

import (
	"bufio"
	"container/list"
	"fmt"
	"os"
	"strconv"
	"strings"

	_ "embed"

	"github.com/rivo/uniseg"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

//go:embed IosevkaNerdFontMono-Medium.ttf
var FontBytes []byte

const EditorWidth = 1000
const EditorHeight = 600

type Row struct {
	Chars       string
	CharWidthts []int
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
	CX           int
	CY           int
	EditorWidth  int
	EditorHeight int
	FontSize     int
	LineHeight   int
	Rows         []Row
	SbarHeight   int
	ShiftMapper  map[string]string
	TabWidth     int
	ColorScheme  ColorScheme

	// Buffer-specific
	CurrentColOffset      int
	CurrentRowOffset      int
	CurrentEditorMode     EditorMode
	CurrentFileName       string
	CurrentFileNumRows    int
	CurrentLineNoColWidth int
	CurrentScreenCols     int
	CurrentScreenRows     int
	ConvertTabToSpace     bool
	CommandBuffer         string
	MultiplierBuffer      string

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

	renderer.SetViewport(&sdl.Rect{X: 0, Y: 0, W: EditorWidth, H: EditorHeight})

	fontSize := 14
	lineHeight := 22
	sbarHeight := 24

	rawFont, err := sdl.RWFromMem(FontBytes)
	if err != nil {
		panic(err)
	}
	font, err := ttf.OpenFontRW(rawFont, 1, fontSize)
	if err != nil {
		panic(err)
	}

	w, h := window.GetSize()

	cfg := &RaptorCfg{
		CX:           0,
		CY:           0,
		EditorWidth:  int(w),
		EditorHeight: int(h),
		FontSize:     fontSize,
		LineHeight:   lineHeight,
		Rows:         []Row{},
		SbarHeight:   sbarHeight,
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
		TabWidth: 4,
		Toasts:   *list.New(),
		ColorScheme: ColorScheme{
			background:    "#272822",
			foreground:    "#ffffff",
			comment:       "#75715e",
			keyword:       "#f92672",
			operator:      "#f92672",
			numberLiteral: "#ae81ff",
			stringLiteral: "#e6db74",
		},

		CurrentRowOffset:   0,
		CurrentColOffset:   0,
		CurrentScreenCols:  80,
		CurrentScreenRows:  (EditorHeight - sbarHeight) / lineHeight,
		CurrentFileNumRows: 0,
		ConvertTabToSpace:  false,

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
		r.CurrentFileNumRows += 1
	}

	if r.CurrentFileNumRows == 0 {
		r.Rows = append(r.Rows, Row{"", []int{}})
		r.CurrentFileNumRows += 1
	}

	r.sdlWindow.SetTitle(fileName)
	r.sdlWindow.SetResizable(true)
	r.CurrentFileName = fileName
	return nil
}

func (r *RaptorCfg) HandleKeyPress(ev *sdl.KeyboardEvent) {
	if r.CurrentEditorMode == EditorModeNormal {
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
		default:
			r.HandleNormalCommand(ev)
		}
	} else if r.CurrentEditorMode == EditorModeInsert {
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
						if r.CX > 0 && r.CurrentEditorMode == EditorModeInsert {
							r.CX -= 1
						}
						r.CurrentEditorMode = EditorModeNormal
						r.ClearCommandBuffer()
					}

					r.HandleKeyPress(t)
				}
			}
		}

		r.DrawScreen()
		sdl.Delay(33)
	}
}

func (r *RaptorCfg) Destroy() {
	r.sdlWindow.Destroy()
	r.renderer.Destroy()
}

func (r *RaptorCfg) DrawSBar() {
	if r.CurrentEditorMode == EditorModeNormal {
		r.DrawSBarVisual()
	} else if r.CurrentEditorMode == EditorModeInsert {
		r.DrawSBarInsert()
	}
}

func (r *RaptorCfg) DrawScreen() {
	r.renderer.Clear()

	// Text background
	bgColor := r.ColorScheme.Background()
	r.renderer.SetDrawColor(bgColor.R, bgColor.G, bgColor.B, 255)
	winW, winH := r.sdlWindow.GetSize()
	r.renderer.FillRect(&sdl.Rect{X: 0, Y: 0, W: winW, H: winH})

	//=== Line number area
	screenMaxLineNo := strconv.Itoa(r.CurrentRowOffset + r.CurrentScreenRows)
	surf, _ := r.sdlFont.RenderUTF8Blended(screenMaxLineNo, sdl.Color{})
	lineNoColWidth := surf.W + 20 + 20
	r.renderer.SetDrawColor(20, 20, 20, 255)
	r.renderer.FillRect(&sdl.Rect{
		X: 0, Y: 0, W: int32(lineNoColWidth), H: winH,
	})
	for y := 0; y < r.CurrentScreenRows; y += 1 {
		if y < r.CurrentFileNumRows {
			numStr := strconv.Itoa(y + r.CurrentRowOffset + 1)
			s, _ := r.sdlFont.RenderUTF8Blended(numStr, r.ColorScheme.Foreground())
			t, _ := r.renderer.CreateTextureFromSurface(s)

			r.renderer.Copy(t, nil, &sdl.Rect{X: 20, Y: int32(y * r.LineHeight), W: s.W, H: s.H})
			s.Free()
			t.Destroy()
		}
	}
	r.CurrentLineNoColWidth = int(lineNoColWidth)

	//=== Render the actual texts
	w := EstimateCharWidth(r.sdlFont)
	for y := r.CurrentRowOffset; y < r.CurrentRowOffset+r.CurrentScreenRows; y += 1 {
		if y < r.CurrentFileNumRows {
			r.renderBufferText(
				int(r.CurrentLineNoColWidth+8)-w,
				(y-r.CurrentRowOffset)*r.LineHeight,
				r.Rows[y].Chars,
				r.LineHeight,
				r.renderer)
		}
	}

	oscWidths, oscWidthCumsums := r.EstimateOnScreenCharWidthsAndCumsum()
	//=== Syntax highlight
	r.DrawHighlight(oscWidthCumsums)

	//=== Cursor
	r.DrawCursor(oscWidths[r.CurrentRowOffset+r.CY], oscWidthCumsums[r.CurrentRowOffset+r.CY])

	//=== Statusbar
	r.DrawSBar()

	//=== Update and draw toast list
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

	gr := uniseg.NewGraphemes(text)
	for gr.Next() {
		c := gr.Str()
		if c == "\t" {
			c = strings.Repeat(" ", r.TabWidth)
		}
		s, _ := r.sdlFont.RenderUTF8Blended(c, r.ColorScheme.Foreground())
		t, _ := renderer.CreateTextureFromSurface(s)
		t.SetBlendMode(sdl.BLENDMODE_BLEND)
		renderer.Copy(t, nil, &sdl.Rect{X: int32(charOffsetX + r.CurrentColOffset), Y: int32(offsetY), W: s.W, H: s.H})
		charOffsetX += int(s.W)

		s.Free()
		t.Destroy()
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

	if len(os.Args) != 2 {
		fmt.Println("Specify filename")
		os.Exit(1)
	}

	fileName := os.Args[1]
	if err := cfg.OpenFile(fileName); err != nil {
		fmt.Println("Error opening", fileName)
		os.Exit(1)
	}

	cfg.Toasts.PushFront(NewToast("Editing "+fileName, 1, cfg.renderer, cfg.sdlFont))
	cfg.Run()

}
