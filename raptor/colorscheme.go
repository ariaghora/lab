package main

import (
	"strconv"
	"strings"

	"github.com/veandco/go-sdl2/sdl"
)

type ColorScheme struct {
	background    string
	foreground    string
	comment       string
	keyword       string
	operator      string
	numberLiteral string
	stringLiteral string
}

func parseColor(hex string) (uint8, uint8, uint8, error) {
	// Remove the "#" character from the beginning of the hex string
	hex = strings.TrimPrefix(hex, "#")

	// Convert the hex string to RGB values
	rgb, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return 0, 0, 0, err
	}

	// Extract the individual RGB components
	red := (rgb >> 16) & 0xFF
	green := (rgb >> 8) & 0xFF
	blue := rgb & 0xFF

	return uint8(red), uint8(green), uint8(blue), nil
}

func (c *ColorScheme) Background() sdl.Color {
	r, g, b, _ := parseColor(c.background)
	return sdl.Color{R: r, G: g, B: b, A: 255}
}

func (c *ColorScheme) Foreground() sdl.Color {
	r, g, b, _ := parseColor(c.foreground)
	return sdl.Color{R: r, G: g, B: b, A: 255}
}

func (c *ColorScheme) Comment() sdl.Color {
	r, g, b, _ := parseColor(c.comment)
	return sdl.Color{R: r, G: g, B: b, A: 255}
}

func (c *ColorScheme) Keyword() sdl.Color {
	r, g, b, _ := parseColor(c.keyword)
	return sdl.Color{R: r, G: g, B: b, A: 255}
}

func (c *ColorScheme) NumberLiteral() sdl.Color {
	r, g, b, _ := parseColor(c.numberLiteral)
	return sdl.Color{R: r, G: g, B: b, A: 255}
}

func (c *ColorScheme) StringLiteral() sdl.Color {
	r, g, b, _ := parseColor(c.stringLiteral)
	return sdl.Color{R: r, G: g, B: b, A: 255}
}

func (c *ColorScheme) Operator() sdl.Color {
	r, g, b, _ := parseColor(c.operator)
	return sdl.Color{R: r, G: g, B: b, A: 255}
}
