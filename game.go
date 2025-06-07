package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	Spawners      []*Spawner
	Player        *Player
	Enemies       []*Enemy
	Camera        *Camera
	FloatingTexts []*FloatingText
	Inventory     *Inventory
}

func (g *Game) Update() error {
	for _, spawner := range g.Spawners {
		spawner.Update(g)
	}

	aliveEnemies := []*Enemy{}
	for _, enemy := range g.Enemies {
		if !enemy.Dead {
			enemy.Update(g.Player.X, g.Player.Y)
			aliveEnemies = append(aliveEnemies, enemy)
		}
	}
	g.Enemies = aliveEnemies

	for _, enemy := range g.Enemies {
		if enemy.Dead {
			continue
		}

		// Distance to player
		dx := enemy.X - g.Player.X
		dy := enemy.Y - g.Player.Y
		dist := math.Hypot(dx, dy)

		if dist < 20 && g.Player.DamageCooldown == 0 {
			g.Player.TakeDamage(1)
			g.Player.DamageCooldown = 30 // half second cooldown

			// Floating text
			g.FloatingTexts = append(g.FloatingTexts, &FloatingText{
				X:           g.Player.X,
				Y:           g.Player.Y - 10,
				Text:        "-1",
				Color:       color.RGBA{255, 0, 0, 255},
				Lifetime:    0,
				MaxLifetime: 60,
				Alpha:       1.0,
			})
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.Camera.Zoom += 0.01
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		g.Camera.Zoom -= 0.01
		if g.Camera.Zoom < 0.1 {
			g.Camera.Zoom = 0.1
		}
	}

	for i := 0; i < len(g.FloatingTexts); {
		g.FloatingTexts[i].Update()
		if g.FloatingTexts[i].Lifetime >= g.FloatingTexts[i].MaxLifetime {
			// Remove expired text
			g.FloatingTexts = append(g.FloatingTexts[:i], g.FloatingTexts[i+1:]...)
		} else {
			i++
		}
	}

	g.Player.Update()

	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		g.Player.Attack(g.Enemies, g)
	}

	g.Camera.CenterOn(g.Player.X, g.Player.Y)

	g.Inventory.Update()

	// For testing, press `P` to add a random item
	if ebiten.IsKeyPressed(ebiten.KeyP) {
		g.Inventory.AddItem(CreateTestItem())
	}

	if ebiten.IsKeyPressed(ebiten.KeyH) {
		g.Player.TakeDamage(1)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

	screen.Fill(color.White)
	g.Player.Draw(screen, g.Camera)

	for _, enemy := range g.Enemies {
		enemy.Draw(screen, g.Camera)
	}

	for _, ft := range g.FloatingTexts {
		ft.Draw(screen, g.Camera)
	}

	g.Inventory.Draw(screen)
	g.Inventory.DrawTooltip(screen)

	drawPlayerHPBar(screen, g.Player)
	ebitenutil.DebugPrint(screen, fmt.Sprintf("X: %.2f Y: %.2f Zoom: %.2f", g.Player.X, g.Player.Y, g.Camera.Zoom))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
