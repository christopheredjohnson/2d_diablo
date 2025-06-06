package main

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Enemy struct {
	X, Y       float64
	Speed      float64
	Sprite     *ebiten.Image
	FrameIndex int
	FrameTimer int
	FrameDelay int
	Frames     []*ebiten.Image // if animated
	HP         int
	Dead       bool
}

func LoadEnemySpriteSheet(path string, frameCount int, frameWidth, frameHeight int) []*ebiten.Image {
	sheet, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatalf("failed to load enemy sprite: %v", err)
	}
	return sliceSpriteSheet(sheet, frameCount, frameWidth, frameHeight)
}

func (e *Enemy) Update(targetX, targetY float64) {
	// Move toward the player
	dx := targetX - e.X
	dy := targetY - e.Y
	dist := math.Hypot(dx, dy)
	if dist > 1 {
		e.X += (dx / dist) * e.Speed
		e.Y += (dy / dist) * e.Speed
	}

	// Advance animation
	e.FrameTimer++
	if e.FrameTimer >= e.FrameDelay {
		e.FrameTimer = 0
		e.FrameIndex = (e.FrameIndex + 1) % len(e.Frames)
	}
}

func (e *Enemy) Draw(screen *ebiten.Image, camera *Camera) {
	img := e.Frames[e.FrameIndex]
	op := &ebiten.DrawImageOptions{}

	frameWidth := float64(img.Bounds().Dx())
	frameHeight := float64(img.Bounds().Dy())

	// ✅ 1. Center origin (world space)
	op.GeoM.Translate(-frameWidth/2, -frameHeight/2)

	// ✅ 2. Move to world position relative to camera
	op.GeoM.Translate(e.X-camera.X, e.Y-camera.Y)

	// ✅ 3. Apply zoom
	op.GeoM.Scale(camera.Zoom, camera.Zoom)

	screen.DrawImage(img, op)
}

func (e *Enemy) TakeDamage(amount int) {
	e.HP -= amount
	if e.HP <= 0 {
		e.Dead = true
	}
}
