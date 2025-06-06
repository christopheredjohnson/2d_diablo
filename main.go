package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("2D Diablo-Like")

	animations := LoadPlayerAnimations(96, 80, 5)

	player := &Player{
		Animations:     animations,
		X:              screenWidth / 2,
		Y:              screenHeight / 2,
		Speed:          playerSpeed,
		Dir:            Down,
		State:          Idle,
		FrameWidth:     96,
		FrameHeight:    80,
		FrameDelay:     5,
		AttackCooldown: 20,
		AttackRange:    40,
	}

	greenSlimeFrames := LoadEnemySpriteSheet("assets/slime/green.png", 11, 16, 32)
	redSlimeFrames := LoadEnemySpriteSheet("assets/slime/red.png", 11, 16, 32)

	enemies := []*Enemy{
		{
			X:          10,
			Y:          100,
			Speed:      0.5,
			Frames:     greenSlimeFrames,
			FrameDelay: 6,
			FrameIndex: 0,
			FrameTimer: 0,
			HP:         1,
		},
		{
			X:          50,
			Y:          50,
			Speed:      0.5,
			Frames:     redSlimeFrames,
			FrameDelay: 6,
			FrameIndex: 0,
			FrameTimer: 0,
			HP:         5,
		},
	}

	camera := &Camera{
		X:      0,
		Y:      0,
		Zoom:   1.0, // Or 1.0 for no zoom
		Width:  screenWidth,
		Height: screenHeight,
	}

	g := &Game{
		Player:  player,
		Enemies: enemies,
		Camera:  camera,
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
