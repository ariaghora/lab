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
		file, err := os.OpenFile(r.FileName, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		writer := bufio.NewWriter(file)
		for _, row := range r.Rows {
			fmt.Fprintln(writer, row.Chars)
		}

		// Use the Flush method to make sure all buffered operations have been applied to the underlying writer.
		if err = writer.Flush(); err != nil {
			fmt.Println("Error flushing the buffer:", err)
		}
		fmt.Println("Saved")

	}
}

func (r *RaptorCfg) HandleEnteringInsertMode(ev *sdl.KeyboardEvent) {
	if ev.Keysym.Scancode == sdl.SCANCODE_A {
		r.LastInsertMethod = InsertMethodAppend
		r.CX += 1
	} else if ev.Keysym.Scancode == sdl.SCANCODE_O {
		r.LastInsertMethod = InsertMethodBreakLine
		r.IMLineBreak()
	}
	r.EditorMode = EditorModeInsert
	r.DrawScreen()
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
	r.DrawScreen()
}

func (r *RaptorCfg) HandleKeyPressH() {
	r.CX -= 1
	if r.CX < 0 {
		r.CX += 1
	}
	r.DrawScreen()
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
	r.DrawScreen()
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
	r.DrawScreen()
}
