package main

import (
	"image"
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
}

func (g *Game) Update() error {
	g.Player.Update()

	for _, enemy := range g.Enemies {
		enemy.Update(g.Player.X, g.Player.Y)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Player.Draw(screen)

	for _, enemy := range g.Enemies {
		enemy.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

// Cuts a sprite sheet into evenly sized frames
func sliceSpriteSheet(sheet *ebiten.Image, frameCount, frameWidth, frameHeight int) []*ebiten.Image {
	frames := []*ebiten.Image{}
	for i := range frameCount {
		rect := image.Rect(i*frameWidth, 0, (i+1)*frameWidth, frameHeight)
		frame := sheet.SubImage(rect).(*ebiten.Image)
		frames = append(frames, frame)
	}
	return frames
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
		FrameDelay:     10,
		AttackCooldown: 20, // frames between attacks
		AttackRange:    40.0,
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

	g := &Game{
		Player:  player,
		Enemies: enemies,
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
