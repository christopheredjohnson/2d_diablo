package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("2D Diablo-Like")

	animations := LoadPlayerAnimations(96, 80, 5)
	inventory := NewInventory(4, 6)
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
	// spawner := &Spawner{
	// 	X:             300,
	// 	Y:             300,
	// 	SpawnInterval: 180, // every 3 seconds at 60 FPS
	// 	MaxEnemies:    10,
	// }
	enemies := []*Enemy{}

	camera := &Camera{
		X:      0,
		Y:      0,
		Zoom:   2.0, // Or 1.0 for no zoom
		Width:  screenWidth,
		Height: screenHeight,
	}

	g := &Game{
		// Spawners:  []*Spawner{spawner},
		Player:    player,
		Enemies:   enemies,
		Camera:    camera,
		Inventory: inventory,
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
