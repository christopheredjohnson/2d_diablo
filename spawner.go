package main

import (
	"math/rand"

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

		var frames []*ebiten.Image
		var hp int

		switch rand.Intn(3) {
		case 0:
			frames = LoadEnemySpriteSheet("assets/slime/green.png", 11, 16, 32)
			hp = 1
		case 1:
			frames = LoadEnemySpriteSheet("assets/slime/red.png", 11, 16, 32)
			hp = 2
		case 2:
			frames = LoadEnemySpriteSheet("assets/bat/default.png", 4, 32, 32)
			hp = 5
		}

		enemy := &Enemy{
			X:          s.X + float64(rand.Intn(40)-20), // slight spawn offset
			Y:          s.Y + float64(rand.Intn(40)-20),
			Speed:      1.0,
			Frames:     frames,
			FrameDelay: 7,
			FrameIndex: 0,
			FrameTimer: 0,
			HP:         hp,
		}

		game.Enemies = append(game.Enemies, enemy)
	}
}
