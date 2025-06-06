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
}

func (ft *FloatingText) Update() {
	// rise
	ft.Y -= 0.5
	ft.Lifetime++
	ft.Alpha = 1 - float64(ft.Lifetime)/float64(ft.MaxLifetime)
}

func (ft *FloatingText) Draw(screen *ebiten.Image, cam *Camera) {
	x := int(ft.X - cam.X)
	y := int(ft.Y - cam.Y)

	col := color.RGBA{255, 0, 0, uint8(ft.Alpha * 255)}

	text.Draw(screen, ft.Text, basicfont.Face7x13, x, y, col)
}
