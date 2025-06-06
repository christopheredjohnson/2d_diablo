package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	Player        *Player
	Enemies       []*Enemy
	Camera        *Camera
	FloatingTexts []*FloatingText
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

	aliveTexts := []*FloatingText{}
	for _, ft := range g.FloatingTexts {
		ft.Update()
		if ft.Lifetime < ft.MaxLifetime {
			aliveTexts = append(aliveTexts, ft)
		}
	}
	g.FloatingTexts = aliveTexts

	g.Player.Update()

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.Player.Attack(g.Enemies, g)
	}

	g.Camera.CenterOn(g.Player.X, g.Player.Y)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Player.Draw(screen, g.Camera)

	for _, enemy := range g.Enemies {
		enemy.Draw(screen, g.Camera)
	}

	for _, ft := range g.FloatingTexts {
		ft.Draw(screen, g.Camera)
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("X: %.2f Y: %.2f Zoom: %.2f", g.Player.X, g.Player.Y, g.Camera.Zoom))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
