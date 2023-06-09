package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/veandco/go-sdl2/sdl"
)

func (r *RaptorCfg) HandleBufferWriteToFile(ev *sdl.KeyboardEvent) {
	ctrl := false
	if (ev.Keysym.Mod & sdl.KMOD_CTRL) != 0 {
		ctrl = true
	}

	if ctrl {
		file, err := os.OpenFile(r.FileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		writer := bufio.NewWriter(file)
		for i, row := range r.Rows {
			line := row.Chars
			if len(line) == 0 {
				line = "\n"
			}

			if i > 0 {
				fmt.Fprint(writer, "\n")
			}
			fmt.Fprint(writer, line)
		}

		// Use the Flush method to make sure all buffered operations have been applied to the underlying writer.
		if err = writer.Flush(); err != nil {
			fmt.Println("Error flushing the buffer:", err)
		}

		r.Toasts.PushFront(NewToast("Saved "+r.FileName, 2, r.renderer, r.sdlFont))
	}
}

func (r *RaptorCfg) HandleEnteringInsertMode(ev *sdl.KeyboardEvent) {
	if ev.Keysym.Scancode == sdl.SCANCODE_A {
		r.LastInsertMethod = InsertMethodAppend
		if len(r.Rows[r.RowOffset+r.CY].Chars) > 0 {
			r.CX += 1
		}
	} else if ev.Keysym.Scancode == sdl.SCANCODE_O {
		shifted := false
		if ev.Keysym.Mod&sdl.KMOD_SHIFT != 0 {
			shifted = true
		}

		if shifted {
			r.Rows = insertAtIndex(r.Rows, r.CY+r.RowOffset, Row{""})
		} else {
			r.Rows = insertAtIndex(r.Rows, r.CY+r.RowOffset+1, Row{""})
			r.CX = 0
			r.CY += 1
		}

		r.LastInsertMethod = InsertMethodBreakLine
		r.NumRows += 1
	}
	r.EditorMode = EditorModeInsert
}

func (r *RaptorCfg) HandleKeyPressL() {
	// prevent moving to the right on a single-char line
	if len(r.Rows[r.CY+r.RowOffset].Chars) == 0 {
		return
	}

	r.CX += 1
	if r.CX > len(r.Rows[r.RowOffset+r.CY].Chars)-1 {
		r.CX = len(r.Rows[r.RowOffset+r.CY].Chars) - 1
	}
}

func (r *RaptorCfg) HandleKeyPressH() {
	r.CX -= 1
	if r.CX < 0 {
		r.CX += 1
	}
}

func (r *RaptorCfg) HandleKeyPressJ() {
	r.CY += 1

	if r.RowOffset+r.CY == r.NumRows {
		r.CY -= 1
		return
	}

	if r.CY > r.ScreenRows-1 {
		r.RowOffset += 1
		r.CY -= 1
	}

	if r.RowOffset-r.ScreenRows >= r.NumRows {
		return
	}

	// fix CX position after succesfully moving vertically
	if r.CX > len(r.Rows[r.CY+r.RowOffset].Chars)-1 {
		r.CX = len(r.Rows[r.CY+r.RowOffset].Chars) - 1
		if r.CX < 0 {
			r.CX = 0
		}
	}
}

func (r *RaptorCfg) HandleKeyPressK() {
	r.CY -= 1
	if r.CY < 0 {
		r.CY += 1
		if r.RowOffset > 0 {
			r.RowOffset -= 1
		}
	}

	// fix CX position after succesfully moving vertically
	if r.CX > len(r.Rows[r.CY+r.RowOffset].Chars)-1 {
		r.CX = len(r.Rows[r.CY+r.RowOffset].Chars) - 1
		if r.CX < 0 {
			r.CX = 0
		}
	}
}
