package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 800
	screenHeight = 600
	playerSpeed  = 2.5
)

type Game struct {
	Player  *Player
	Enemies []*Enemy
	Camera  *Camera
}

func (g *Game) Update() error {
	g.Player.Update()

	// Camera centers on the player
	g.Camera.X = g.Player.X - float64(g.Camera.Width)/2
	g.Camera.Y = g.Player.Y - float64(g.Camera.Height)/2

	for _, enemy := range g.Enemies {
		enemy.Update(g.Player.X, g.Player.Y)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Player.Draw(screen, g.Camera)

	for _, enemy := range g.Enemies {
		enemy.Draw(screen, g.Camera)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("2D Diablo-Like")

	animations := LoadPlayerAnimations(96, 80, 5)

	player := &Player{
		Animations:  animations,
		X:           screenWidth / 2,
		Y:           screenHeight / 2,
		Speed:       playerSpeed,
		Dir:         Down,
		State:       Idle,
		FrameWidth:  96,
		FrameHeight: 80,
		FrameDelay:  10,
	}

	batFrames := LoadEnemySpriteSheet("assets/bat/default.png", 4, 32, 32)

	enemies := []*Enemy{
		{
			X:          100,
			Y:          100,
			Speed:      1.2,
			Frames:     batFrames,
			FrameDelay: 10,
			FrameIndex: 0,
			FrameTimer: 0,
		},
	}

	cam := &Camera{
		X:      0,
		Y:      0,
		Zoom:   1.5,
		Width:  screenWidth,
		Height: screenHeight,
	}

	g := &Game{
		Player:  player,
		Enemies: enemies,
		Camera:  cam,
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
