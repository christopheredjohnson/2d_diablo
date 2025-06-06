package main

import (
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	aliveEnemies := []*Enemy{}
	for _, enemy := range g.Enemies {
		if !enemy.Dead {
			enemy.Update(g.Player.X, g.Player.Y)
			aliveEnemies = append(aliveEnemies, enemy)
		}
	}
	g.Enemies = aliveEnemies

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.Camera.Zoom += 0.01
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		g.Camera.Zoom -= 0.01
		if g.Camera.Zoom < 0.1 {
			g.Camera.Zoom = 0.1
		}
	}

	g.Player.Update()

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.Player.Attack(g.Enemies)
	}

	if g.Player.AttackTimer > 0 {
		g.Player.AttackTimer--
	}

	g.Camera.CenterOn(g.Player.X, g.Player.Y)

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

	if g.Player.AttackTimer == g.Player.AttackCooldown-1 {
		screenX := (g.Player.X - g.Camera.X) * g.Camera.Zoom
		screenY := (g.Player.Y - g.Camera.Y) * g.Camera.Zoom

		radius := g.Player.AttackRange * g.Camera.Zoom

		ebitenutil.DrawCircle(screen, screenX, screenY, radius, color.RGBA{255, 0, 0, 100})
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

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
		FrameDelay:     1,
		AttackCooldown: 20,
		AttackRange:    40,
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
			HP:         3,
		},
	}

	camera := &Camera{
		X:      0,
		Y:      0,
		Zoom:   2.0, // Or 1.0 for no zoom
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
