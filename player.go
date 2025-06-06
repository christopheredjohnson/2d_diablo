package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Direction int

const (
	Down Direction = iota
	Up
	Left
	Right
)

type PlayerState int

const (
	Idle PlayerState = iota
	Running
)

type Player struct {
	Animations     map[PlayerState]map[Direction][]*ebiten.Image
	X, Y           float64
	Speed          float64
	Dir            Direction
	State          PlayerState
	FrameIndex     int
	FrameTimer     int
	FrameDelay     int
	FrameWidth     int
	FrameHeight    int
	IsAttacking    bool
	AttackTimer    int
	AttackCooldown int
	AttackRange    float64
}

func (p *Player) AdvanceFrame() {
	p.FrameTimer++
	if p.FrameTimer >= p.FrameDelay {
		p.FrameTimer = 0
		p.FrameIndex = (p.FrameIndex + 1) % len(p.Animations[p.State][p.Dir])
	}
}

// Called each frame to move and animate the player
func (p *Player) Update() {
	moved := false

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		p.Y -= p.Speed
		p.Dir = Up
		moved = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		p.Y += p.Speed
		p.Dir = Down
		moved = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		p.X -= p.Speed
		p.Dir = Left
		moved = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.X += p.Speed
		p.Dir = Right
		moved = true
	}

	if moved {
		p.State = Running
		p.AdvanceFrame()
	} else {
		p.State = Idle
		p.AdvanceFrame()
	}
}

// Draws the correct frame of the player based on state/direction
func (p *Player) Draw(screen *ebiten.Image, camera *Camera) {
	frames := p.Animations[p.State][p.Dir]
	img := frames[p.FrameIndex%len(frames)]

	op := &ebiten.DrawImageOptions{}

	// 1. Move the origin to the center of the sprite (unscaled)
	op.GeoM.Translate(-float64(p.FrameWidth)/2, -float64(p.FrameHeight)/2)

	// 2. Move to world position relative to the camera
	op.GeoM.Translate(p.X-camera.X, p.Y-camera.Y)

	// 3. Scale the whole thing
	op.GeoM.Scale(camera.Zoom, camera.Zoom)

	screen.DrawImage(img, op)
}

// Loads and slices all player animations from sprite sheets
func LoadPlayerAnimations(frameWidth, frameHeight, frameDelay int) map[PlayerState]map[Direction][]*ebiten.Image {
	animations := make(map[PlayerState]map[Direction][]*ebiten.Image)

	states := []PlayerState{Idle, Running}
	directions := []Direction{Down, Up, Left, Right}
	stateNames := map[PlayerState]string{
		Idle:    "idle",
		Running: "run",
	}
	dirNames := map[Direction]string{
		Down:  "down",
		Up:    "up",
		Left:  "left",
		Right: "right",
	}

	frameCounts := map[PlayerState]int{
		Idle:    8,
		Running: 8,
	}

	for _, state := range states {
		animations[state] = make(map[Direction][]*ebiten.Image)
		for _, dir := range directions {
			path := fmt.Sprintf("assets/player/%s_%s.png", stateNames[state], dirNames[dir])
			sheet, _, err := ebitenutil.NewImageFromFile(path)
			if err != nil {
				log.Fatalf("failed to load %s: %v", path, err)
			}
			frames := sliceSpriteSheet(sheet, frameCounts[state], frameWidth, frameHeight)
			animations[state][dir] = frames
		}
	}

	return animations
}
