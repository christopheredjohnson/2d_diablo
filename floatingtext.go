package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type FloatingText struct {
	X, Y        float64
	Text        string
	Alpha       float64
	Lifetime    int
	MaxLifetime int
	Color       color.Color
}

func (ft *FloatingText) Update() {
	ft.Y -= 0.5 // rising effect
	ft.Lifetime++
	ft.Alpha = 1 - float64(ft.Lifetime)/float64(ft.MaxLifetime)
}

func (ft *FloatingText) Draw(screen *ebiten.Image, cam *Camera) {
	// Apply camera transformation to world position
	x := (ft.X - cam.X) * cam.Zoom
	y := (ft.Y - cam.Y) * cam.Zoom

	// Fade alpha
	r, g, b, _ := ft.Color.RGBA()
	col := color.RGBA{
		R: uint8(r >> 8),
		G: uint8(g >> 8),
		B: uint8(b >> 8),
		A: uint8(ft.Alpha * 255),
	}

	text.Draw(screen, ft.Text, basicfont.Face7x13, int(x), int(y), col)
}
