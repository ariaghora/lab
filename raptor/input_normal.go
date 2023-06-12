package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
)

func (r *RaptorCfg) HandleBufferWriteToFile(ev *sdl.KeyboardEvent) {
	ctrl := false
	if (ev.Keysym.Mod & sdl.KMOD_CTRL) != 0 {
		ctrl = true
	}

	if ctrl {
		file, err := os.OpenFile(r.CurrentFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		writer := bufio.NewWriter(file)
		for i, row := range r.Rows {
			line := row.Chars
			if len(line) == 0 {
				fmt.Fprint(writer, "\n")
				continue
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

		r.Toasts.PushFront(NewToast("Saved "+r.CurrentFileName, 2, r.renderer, r.sdlFont))
	}
}

func (r *RaptorCfg) CurrentRowChars() string {
	return r.Rows[r.CY+r.CurrentRowOffset].Chars
}

func (r *RaptorCfg) HandleEnteringInsertMode(ev *sdl.KeyboardEvent) {
	if ev.Keysym.Scancode == sdl.SCANCODE_A {
		if len(r.Rows[r.CurrentRowOffset+r.CY].Chars) > 0 {
			r.CX += 1
		}
	} else if ev.Keysym.Scancode == sdl.SCANCODE_O {
		shifted := false
		if ev.Keysym.Mod&sdl.KMOD_SHIFT != 0 {
			shifted = true
		}

		leadingIndent := r.GetCurrentLeadingIndent()
		if shifted {
			r.Rows = insertAtIndex(r.Rows, r.CY+r.CurrentRowOffset, Row{leadingIndent + "", []int{}})
			r.CX = len(leadingIndent)
		} else {
			r.Rows = insertAtIndex(r.Rows, r.CY+r.CurrentRowOffset+1, Row{leadingIndent + "", []int{}})
			r.CX = len(leadingIndent)
			r.CY += 1
		}

		r.CurrentFileNumRows += 1
	}
	r.CurrentEditorMode = EditorModeInsert
}

func (r *RaptorCfg) HandleKeyPressL() {
	// prevent moving to the right on a single-char line
	if len(r.Rows[r.CY+r.CurrentRowOffset].Chars) == 0 {
		return
	}

	r.CX += 1
	if r.CX > len(r.Rows[r.CurrentRowOffset+r.CY].Chars)-1 {
		r.CX = len(r.Rows[r.CurrentRowOffset+r.CY].Chars) - 1
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

	if r.CurrentRowOffset+r.CY == r.CurrentFileNumRows {
		r.CY -= 1
		return
	}

	if r.CY > r.CurrentScreenRows-1 {
		r.CurrentRowOffset += 1
		r.CY -= 1
	}

	if r.CurrentRowOffset-r.CurrentScreenRows >= r.CurrentFileNumRows {
		return
	}

	// fix CX position after succesfully moving vertically
	if r.CX > len(r.Rows[r.CY+r.CurrentRowOffset].Chars)-1 {
		r.CX = len(r.Rows[r.CY+r.CurrentRowOffset].Chars) - 1
		if r.CX < 0 {
			r.CX = 0
		}
	}
}

func (r *RaptorCfg) HandleKeyPressK() {
	r.CY -= 1
	if r.CY < 0 {
		r.CY += 1
		if r.CurrentRowOffset > 0 {
			r.CurrentRowOffset -= 1
		}
	}

	// fix CX position after succesfully moving vertically
	if r.CX > len(r.Rows[r.CY+r.CurrentRowOffset].Chars)-1 {
		r.CX = len(r.Rows[r.CY+r.CurrentRowOffset].Chars) - 1
		if r.CX < 0 {
			r.CX = 0
		}
	}
}

func (r *RaptorCfg) SpawnToast(message string, durationSeconds int) {
	r.Toasts.PushFront(NewToast(message, durationSeconds, r.renderer, r.sdlFont))
}

func (r *RaptorCfg) ClearCommandBuffer() {
	r.CommandBuffer = ""
	r.MultiplierBuffer = ""
}

func (r *RaptorCfg) RaiseNormalCommandError(message string) {
	r.ClearCommandBuffer()
	r.SpawnToast(message, 1)
}

func (r *RaptorCfg) HandleNormalCommand(ev *sdl.KeyboardEvent) {
	var editorCommands = map[string]func(){
		"dd": func() {
			r.DeleteRow(r.CY + r.CurrentRowOffset)
		},
		"diw": func() {
			r.SpawnToast("TODO: delete current word", 1)
		},
	}

	if ev.Keysym.Scancode == sdl.SCANCODE_RETURN {
		if r.CommandBuffer != "" {
			r.SpawnToast("command: "+r.MultiplierBuffer+r.CommandBuffer, 1)
		}
		r.ClearCommandBuffer()
		return
	} else {
		c := ev.Keysym.Sym
		if c >= '0' && c <= '9' {
			if r.CommandBuffer != "" {
				r.RaiseNormalCommandError("Command's multiplier number must be entered first")
				return
			}
			r.MultiplierBuffer += string(c)
		} else if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') {
			r.CommandBuffer += string(c)
		}
	}

	keys := make([]string, 0, len(editorCommands))
	for k := range editorCommands {
		if strings.HasPrefix(k, r.CommandBuffer) {
			keys = append(keys, k)
		}
	}
	fmt.Println(keys)
	if len(keys) == 1 {
		r.SpawnToast("Executing "+keys[0], 1)
		editorCommands[keys[0]]()
		r.ClearCommandBuffer()
	} else if len(keys) == 0 {
		r.RaiseNormalCommandError("Unknown command " + r.CommandBuffer)
	}
}
