package main

import (
	"regexp"

	"github.com/veandco/go-sdl2/sdl"
)

type Token struct {
	Text  string
	Type  string
	Color sdl.Color
	Start int
	End   int
}

func (r *RaptorCfg) tokenize(code string) []Token {
	var rules = []struct {
		Pattern *regexp.Regexp
		Type    string
		Color   sdl.Color
	}{
		{Pattern: regexp.MustCompile(`\b(def|return|if|else|for|while|break|continue|pass)\b`), Type: "keyword", Color: r.ColorScheme.Keyword()},
		{Pattern: regexp.MustCompile(`#.*`), Type: "comment", Color: r.ColorScheme.Comment()},
		{Pattern: regexp.MustCompile(`"[^"]*"`), Type: "string", Color: r.ColorScheme.StringLiteral()},         // String literal enclosed in double quotes
		{Pattern: regexp.MustCompile(`\b[0-9]+\b`), Type: "number", Color: r.ColorScheme.NumberLiteral()},      // Number literal
		{Pattern: regexp.MustCompile(`\b[\+\-\*\/\=\%]\b`), Type: "operator", Color: r.ColorScheme.Operator()}, // Operators: +, -, *, /, =, %
	}

	tokens := []Token{}
	workingCode := code // Copy of the code that we can modify without affecting the original

	// Iterate over the rules.
	for _, rule := range rules {
		matches := rule.Pattern.FindAllStringIndex(code, -1)

		// For each match, create a new token.
		for _, match := range matches {
			tokens = append(tokens, Token{
				Text:  code[match[0]:match[1]],
				Type:  rule.Type,
				Color: rule.Color,
				Start: match[0],
				End:   match[1],
			})
		}

		// Remove matches from the code to avoid duplicate matches.
		workingCode = rule.Pattern.ReplaceAllString(workingCode, "")

	}

	return tokens
}

func (r *RaptorCfg) DrawHighlight(oscWidthCumsums [][]int) {
	w := EstimateCharWidth(r.sdlFont)
	for i, row := range r.Rows {
		tokens := r.tokenize(row.Chars)
		for _, token := range tokens {
			sur, _ := r.sdlFont.RenderUTF8Blended(token.Text, token.Color)
			tex, _ := r.renderer.CreateTextureFromSurface(sur)
			tex.SetBlendMode(sdl.BLENDMODE_ADD)
			rect := &sdl.Rect{
				X: int32(1 + r.CurrentLineNoColWidth + oscWidthCumsums[i][token.Start] - w),
				Y: int32(i * r.LineHeight),
				W: int32(sur.W),
				H: int32(sur.H),
			}

			bg := r.ColorScheme.Background()
			r.renderer.SetDrawColor(bg.R, bg.G, bg.B, 255)
			r.renderer.FillRect(rect)

			r.renderer.Copy(tex, nil, rect)

			sur.Free()
			tex.Destroy()
		}
	}
}
