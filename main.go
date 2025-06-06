package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Spawner struct {
	X, Y          float64
	SpawnInterval int
	Timer         int
	MaxEnemies    int
	Spawned       int
}

func (s *Spawner) Update(game *Game) {
	s.Timer++
	if s.Timer >= s.SpawnInterval && s.Spawned < s.MaxEnemies {
		s.Timer = 0
		s.Spawned++

		enemy := &Enemy{
			X:          s.X,
			Y:          s.Y,
			Speed:      0.5,
			Frames:     LoadEnemySpriteSheet("assets/slime/green.png", 11, 16, 32),
			FrameDelay: 7,
			FrameIndex: 0,
			FrameTimer: 0,
			HP:         1,
		}

		game.Enemies = append(game.Enemies, enemy)
	}
}

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

	// greenSlimeFrames := LoadEnemySpriteSheet("assets/slime/green.png", 11, 16, 32)
	// redSlimeFrames := LoadEnemySpriteSheet("assets/slime/red.png", 11, 16, 32)
	spawner := &Spawner{
		X:             300,
		Y:             300,
		SpawnInterval: 180, // every 3 seconds at 60 FPS
		MaxEnemies:    10,
	}
	enemies := []*Enemy{
		// {
		// 	X:          10,
		// 	Y:          100,
		// 	Speed:      1,
		// 	Frames:     greenSlimeFrames,
		// 	FrameDelay: 7,
		// 	FrameIndex: 0,
		// 	FrameTimer: 0,
		// 	HP:         1,
		// },
		// {
		// 	X:          50,
		// 	Y:          50,
		// 	Speed:      1,
		// 	Frames:     redSlimeFrames,
		// 	FrameDelay: 7,
		// 	FrameIndex: 0,
		// 	FrameTimer: 0,
		// 	HP:         5,
		// },
	}

	camera := &Camera{
		X:      0,
		Y:      0,
		Zoom:   1.0, // Or 1.0 for no zoom
		Width:  screenWidth,
		Height: screenHeight,
	}

	g := &Game{
		Spawners: []*Spawner{spawner},
		Player:   player,
		Enemies:  enemies,
		Camera:   camera,
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
