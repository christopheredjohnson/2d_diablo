package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

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
	Attacking
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
	if p.FrameTimer < p.FrameDelay {
		return
	}
	p.FrameTimer = 0

	frames := p.Animations[p.State][p.Dir]

	if p.State == Attacking {
		// Play attack animation once
		if p.FrameIndex < len(frames)-1 {
			p.FrameIndex++
		} else {
			// Animation finished — return to idle
			p.State = Idle
			p.FrameIndex = 0
			p.IsAttacking = false
		}
		return
	}

	// Looping animation (Idle/Running)
	p.FrameIndex = (p.FrameIndex + 1) % len(frames)
}

// Called each frame to move and animate the player
func (p *Player) Update() {

	if p.AttackTimer > 0 {
		p.AttackTimer--
	}

	if p.State == Attacking {
		p.AdvanceFrame()

		// If attack animation is done, go back to Idle
		if p.FrameIndex >= len(p.Animations[Attacking][p.Dir]) {
			p.State = Idle
			p.FrameIndex = 0
			p.IsAttacking = false
		}
		return
	}

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
	screen.Fill(color.RGBA{30, 30, 30, 255})
	frames := p.Animations[p.State][p.Dir]
	if len(frames) == 0 {
		log.Printf("Missing animation for state %v dir %v", p.State, p.Dir)
		return
	}

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

	states := []PlayerState{Idle, Running, Attacking}
	directions := []Direction{Down, Up, Left, Right}
	stateNames := map[PlayerState]string{
		Idle:      "idle",
		Running:   "run",
		Attacking: "attack2",
	}
	dirNames := map[Direction]string{
		Down:  "down",
		Up:    "up",
		Left:  "left",
		Right: "right",
	}

	frameCounts := map[PlayerState]int{
		Idle:      8,
		Running:   8,
		Attacking: 8,
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

func (p *Player) Attack(enemies []*Enemy, game *Game) {
	if p.AttackTimer > 0 || p.State == Attacking {
		return
	}

	p.IsAttacking = true
	p.State = Attacking
	p.FrameIndex = 0
	p.AttackTimer = p.AttackCooldown

	for _, e := range enemies {
		dx := e.X - p.X
		dy := e.Y - p.Y
		dist := math.Hypot(dx, dy)

		if dist <= p.AttackRange {
			e.TakeDamage(1)
			// Spawn floating text
			fx := &FloatingText{
				X:           e.X,
				Y:           e.Y - 20, // slight offset above the enemy
				Text:        "-1",
				Alpha:       1.0,
				Lifetime:    0,
				MaxLifetime: 60, // about 1 second at 60fps
			}
			game.FloatingTexts = append(game.FloatingTexts, fx)
		}
	}
}
