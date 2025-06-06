package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	// clr := color.RGBA{255, 0, 0, uint8(ft.Alpha * 255)}
	// textOpt := &ebiten.DrawImageOptions{}
	textX := ft.X - cam.X
	textY := ft.Y - cam.Y

	ebitenutil.DebugPrintAt(screen, ft.Text, int(textX), int(textY))
}
