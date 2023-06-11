package main

import (
	"strings"
	"unicode"

	"github.com/veandco/go-sdl2/sdl"
)

func removeAtIndex(s []Row, index int) []Row {
	return append(s[:index], s[index+1:]...)
}

func insertAtIndex(s []Row, index int, value Row) []Row {
	s = append(s[:index], append([]Row{value}, s[index:]...)...)
	return s
}

func (r *RaptorCfg) IMInsertChar(c rune) {
	old := r.Rows[r.CurrentRowOffset+r.CY].Chars
	new := string(old[:r.CX+r.CurrentColOffset]) + string(c) + string(old[r.CX+r.CurrentColOffset:])
	r.Rows[r.CurrentRowOffset+r.CY].Chars = new
	r.CX += 1
}

func (r *RaptorCfg) IMBackspace() {
	x := r.CX + r.CurrentColOffset
	y := r.CurrentRowOffset + r.CY
	old := r.Rows[y].Chars

	// if cursor at the letfmost, then merge current line with previous line
	if x == 0 {
		// if we're at the first line of the file, do nothing
		if y == 0 {
			return
		} else {
			r.CY -= 1
			r.CX = len(r.Rows[y-1].Chars)
			r.Rows[y-1].Chars += old
			r.Rows = removeAtIndex(r.Rows, y)
			r.CurrentFileNumRows -= 1
		}
	} else {
		if r.CX > 0 {
			r.CX -= 1
		}
		new := old[:x-1] + old[x:]
		r.Rows[y].Chars = new
	}

}

func (r *RaptorCfg) GetCurrentLeadingIndent() string {
	line := r.Rows[r.CY+r.CurrentRowOffset].Chars
	if len(strings.TrimSpace(line)) == 0 {
		return ""
	}

	leadingIndent := ""

	i := 0
	c := rune(line[i])
	for unicode.IsSpace(c) && i < len(line) {
		if r.ConvertTabToSpace && c == '\t' {
			leadingIndent += strings.Repeat(" ", r.TabWidth)
		} else {
			leadingIndent += string(c)
		}

		i += 1
		c = rune(line[i])
	}

	if unicode.IsSpace(c) {
		return ""
	}

	return leadingIndent
}

func (r *RaptorCfg) IMLineBreak() {
	x := r.CX + r.CurrentColOffset
	y := r.CurrentRowOffset + r.CY
	newCurrent := r.Rows[y].Chars[:x]
	newLineContent := r.Rows[y].Chars[x:]

	leadingIdent := r.GetCurrentLeadingIndent()
	newLineContent = leadingIdent + newLineContent

	r.Rows[y].Chars = newCurrent
	r.Rows = insertAtIndex(r.Rows, y+1, Row{newLineContent, []int{}})
	r.CY += 1
	r.CX = len(leadingIdent)
	r.CurrentFileNumRows += 1
}

func (r *RaptorCfg) HandleKeyPressInsertMode(ev *sdl.KeyboardEvent) {
	shift := false
	caps := false

	if (ev.Keysym.Mod & sdl.KMOD_SHIFT) != 0 {
		shift = true
	}

	if (ev.Keysym.Mod & sdl.KMOD_CAPS) != 0 {
		caps = true
	}

	var c rune
	if ev.Keysym.Sym >= 'a' && ev.Keysym.Sym <= 'z' {
		if shift || caps {
			c = rune(strings.ToUpper(string(ev.Keysym.Sym))[0])
		} else {
			c = rune(string(ev.Keysym.Sym)[0])
		}
		r.IMInsertChar(c)
	} else if ev.Keysym.Sym >= '0' && ev.Keysym.Sym <= '9' {
		if shift {
			shifted := ")!@#$%^&*("
			c = rune(shifted[ev.Keysym.Sym-'0'])
		} else {
			c = rune(ev.Keysym.Sym)
		}
		r.IMInsertChar(c)
	} else if ev.Keysym.Sym == ' ' {
		r.IMInsertChar(' ')
	} else if ev.Keysym.Sym == '\t' {
		r.IMInsertChar('\t')
	} else if ev.Keysym.Scancode == sdl.SCANCODE_BACKSPACE {
		r.IMBackspace()
	} else if ev.Keysym.Scancode == sdl.SCANCODE_RETURN {
		r.IMLineBreak()
	} else {
		// other input
		if ev.Keysym.Sym >= 32 && ev.Keysym.Sym <= 126 {
			if shift {
				c = rune(r.ShiftMapper[string(ev.Keysym.Sym)][0])
			} else {
				c = rune(string(ev.Keysym.Sym)[0])
			}
			r.IMInsertChar(c)
		}
	}
}
