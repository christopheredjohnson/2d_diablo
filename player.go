package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
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
	HP             int
	MaxHP          int
	DamageCooldown int // frames remaining until next hit
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

	if p.DamageCooldown > 0 {
		p.DamageCooldown--
	}

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
			game.FloatingTexts = append(game.FloatingTexts, &FloatingText{
				X:           e.X,
				Y:           e.Y - 10,
				Text:        "-1",
				Color:       color.RGBA{255, 0, 0, 255},
				Lifetime:    0,
				MaxLifetime: 60,
				Alpha:       1.0,
			})
		}
	}
}

func (p *Player) TakeDamage(dmg int) {
	p.HP -= dmg
	if p.HP < 0 {
		p.HP = 0
	}
}

func drawPlayerHPBar(screen *ebiten.Image, player *Player) {
	x := 20
	y := 20
	width := 200
	height := 16

	// Avoid divide-by-zero
	maxHP := math.Max(1, float64(player.MaxHP))
	hpRatio := float64(player.HP) / maxHP

	// Clamp between 0 and 1
	hpRatio = math.Max(0, math.Min(1, hpRatio))
	barWidth := int(hpRatio * float64(width))

	// Background
	bg := ebiten.NewImage(width, height)
	bg.Fill(color.RGBA{40, 40, 40, 255})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(bg, op)

	// Foreground
	if barWidth > 0 {
		fg := ebiten.NewImage(barWidth, height)
		fg.Fill(color.RGBA{255, 0, 0, 255})
		op2 := &ebiten.DrawImageOptions{}
		op2.GeoM.Translate(float64(x), float64(y))
		screen.DrawImage(fg, op2)
	}

	// Text
	text.Draw(screen,
		fmt.Sprintf("HP: %d / %d", player.HP, player.MaxHP),
		basicfont.Face7x13,
		x+4, y+height+14,
		color.White)
}
