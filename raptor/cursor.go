package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

func (r *RaptorCfg) EstimateOnScreenCharWidthsAndCumsum() ([][]int, [][]int) {
	charWidth := EstimateCharWidth(r.sdlFont)
	rowCharWidths := [][]int{}
	rowCharWidthCumsums := [][]int{}
	for i := r.CurrentRowOffset; i < r.CurrentRowOffset+r.CurrentScreenRows; i += 1 {
		if i < r.CurrentFuleNumRows {

			widths := []int{}
			widthCumsums := []int{}
			for j, c := range r.Rows[i].Chars {
				switch c {
				case '\t':
					widths = append(widths, r.TabWidth*charWidth)
				default:
					widths = append(widths, 1*charWidth)
				}
				widthCumsums = append(widthCumsums, widths[j])

				// cumsum of char widths
				if j > 0 {
					widthCumsums[j] += widthCumsums[j-1]
				}
			}
			widths = append(widths, 1*charWidth)
			// extra 1 last width tracker
			if len(r.Rows[i].Chars) > 0 {
				widthCumsums = append(
					widthCumsums,
					widths[len(widths)-1]+widthCumsums[len(widthCumsums)-1],
				)
			} else {
				widthCumsums = append(widthCumsums, 1*charWidth)
			}
			rowCharWidths = append(rowCharWidths, widths)
			rowCharWidthCumsums = append(rowCharWidthCumsums, widthCumsums)
		}
	}
	return rowCharWidths, rowCharWidthCumsums
}

func (r *RaptorCfg) DrawCursor(rowCharWidths []int, rowCharWidthCumsums []int) {
	w := EstimateCharWidth(r.sdlFont)
	r.renderer.SetDrawColor(255, 255, 255, 255)
	posX := int32(rowCharWidthCumsums[r.CX+r.CurrentColOffset] + int(w+20+20) + 8 - (rowCharWidths[r.CX+r.CurrentColOffset]))
	posY := int32(r.CY * r.LineHeight)
	r.renderer.FillRect(&sdl.Rect{
		X: posX, Y: posY, W: int32(rowCharWidths[r.CX+r.CurrentColOffset]), H: int32(r.LineHeight)},
	)

	curRowChars := r.CurrentRowChars()
	char := " "
	if len(curRowChars) > 0 && r.CX+r.CurrentColOffset < len(curRowChars) {
		char = string(curRowChars[r.CX+r.CurrentColOffset])
	}

	// Print char with inverted color when it is printable
	// TODO: check more printables other than tab
	if char != "\t" {
		sur, _ := r.sdlFont.RenderUTF8Blended(char, sdl.Color{R: 0, G: 0, B: 0, A: 255})
		tex, _ := r.renderer.CreateTextureFromSurface(sur)
		r.renderer.Copy(tex, nil, &sdl.Rect{X: posX, Y: posY, W: sur.W, H: sur.H})
	}
}
