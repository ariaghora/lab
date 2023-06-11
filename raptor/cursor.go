package main

import "github.com/veandco/go-sdl2/sdl"

func (r *RaptorCfg) DrawCursor() {
	charWidth := EstimateCharWidth(r.sdlFont)
	rowCharWidths := []int{}
	rowCharWidthCumsums := []int{}
	for i, c := range r.Rows[r.RowOffset+r.CY].Chars {
		switch c {
		case '\t':
			rowCharWidths = append(rowCharWidths, r.TabWidth*charWidth)
		default:
			rowCharWidths = append(rowCharWidths, 1*charWidth)
		}
		rowCharWidthCumsums = append(rowCharWidthCumsums, rowCharWidths[i])

		// cumsum of char widths
		if i > 0 {
			rowCharWidthCumsums[i] += rowCharWidthCumsums[i-1]
		}
	}
	rowCharWidths = append(rowCharWidths, 1*charWidth)
	// extra 1 last width tracker
	if len(r.Rows[r.RowOffset+r.CY].Chars) > 0 {
		rowCharWidthCumsums = append(
			rowCharWidthCumsums,
			rowCharWidths[len(rowCharWidths)-1]+rowCharWidthCumsums[len(rowCharWidthCumsums)-1],
		)
	} else {
		rowCharWidthCumsums = append(rowCharWidthCumsums, 1*charWidth)
	}

	r.renderer.SetDrawColor(255, 255, 255, 255)
	posX := int32(rowCharWidthCumsums[r.CX+r.ColOffset] + int(charWidth+20+20) + 8 - (rowCharWidths[r.CX+r.ColOffset]))
	posY := int32(r.CY * r.LineHeight)
	r.renderer.FillRect(&sdl.Rect{
		X: posX, Y: posY, W: int32(rowCharWidths[r.CX+r.ColOffset]), H: int32(r.LineHeight)},
	)

	curRowChars := r.CurrentRowChars()
	char := " "
	if len(curRowChars) > 0 && r.CX+r.ColOffset < len(curRowChars) {
		char = string(curRowChars[r.CX+r.ColOffset])
	}

	// Print char with inverted color when it is printable
	// TODO: check more printables other than tab
	if char != "\t" {
		sur, _ := r.sdlFont.RenderUTF8Blended(char, sdl.Color{R: 0, G: 0, B: 0, A: 255})
		tex, _ := r.renderer.CreateTextureFromSurface(sur)
		r.renderer.Copy(tex, nil, &sdl.Rect{X: posX, Y: posY, W: sur.W, H: sur.H})
	}
}
