package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

// Cuts a sprite sheet into evenly sized frames
func sliceSpriteSheet(sheet *ebiten.Image, frameCount, frameWidth, frameHeight int) []*ebiten.Image {
	frames := []*ebiten.Image{}
	for i := 0; i < frameCount; i++ {
		rect := image.Rect(i*frameWidth, 0, (i+1)*frameWidth, frameHeight)
		frame := sheet.SubImage(rect).(*ebiten.Image)
		frames = append(frames, frame)
	}
	return frames
}
