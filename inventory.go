package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

type Item struct {
	ID          string
	Name        string
	Icon        *ebiten.Image
	Quantity    int
	MaxStack    int
	Description string
}

type Slot struct {
	Item *Item
}

type Inventory struct {
	Slots            [][]*Slot
	Rows, Cols       int
	IsOpen           bool
	slotImage        *ebiten.Image
	ToggleCooldown   int
	DraggingItem     *Item
	DragOffsetX      int
	DragOffsetY      int
	prevMousePressed bool
}

func NewInventory(rows, cols int) *Inventory {

	slots := make([][]*Slot, rows)
	for y := range slots {
		slots[y] = make([]*Slot, cols)

		for x := range slots[y] {
			slots[y][x] = &Slot{}
		}
	}

	img := ebiten.NewImage(32, 32)
	img.Fill(color.RGBA{40, 40, 40, 255})

	return &Inventory{
		Slots:     slots,
		Rows:      rows,
		Cols:      cols,
		IsOpen:    false,
		slotImage: img,
	}
}

func (inv *Inventory) Draw(screen *ebiten.Image) {
	if !inv.IsOpen {
		return
	}

	startX := 50
	startY := 50
	spacing := 36

	for y := 0; y < inv.Rows; y++ {
		for x := 0; x < inv.Cols; x++ {
			sx := startX + x*spacing
			sy := startY + y*spacing

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(sx), float64(sy))
			screen.DrawImage(inv.slotImage, op)

			slot := inv.Slots[y][x]
			if slot.Item != nil && slot.Item.Icon != nil {
				iconOp := &ebiten.DrawImageOptions{}
				iconOp.GeoM.Translate(float64(sx), float64(sy))
				screen.DrawImage(slot.Item.Icon, iconOp)

				if slot.Item.Quantity > 1 {
					text.Draw(screen,
						fmt.Sprintf("%d", slot.Item.Quantity),
						basicfont.Face7x13,
						sx+16, sy+30,
						color.White,
					)
				}
			}
		}
	}

	if inv.DraggingItem != nil && inv.DraggingItem.Icon != nil {
		cx, cy := ebiten.CursorPosition()
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(cx-inv.DragOffsetX), float64(cy-inv.DragOffsetY))
		screen.DrawImage(inv.DraggingItem.Icon, op)
	}
}

func (inv *Inventory) Update() {
	if inv.ToggleCooldown > 0 {
		inv.ToggleCooldown--
	}

	if ebiten.IsKeyPressed(ebiten.KeyTab) && inv.ToggleCooldown == 0 {
		inv.IsOpen = !inv.IsOpen
		inv.ToggleCooldown = 15
	}

	if !inv.IsOpen {
		return
	}

	mousePressed := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	if mousePressed && !inv.prevMousePressed {
		// <-- trigger logic only on new press
		cx, cy := ebiten.CursorPosition()
		startX, startY := 50, 50
		spacing := 36

		for y := 0; y < inv.Rows; y++ {
			for x := 0; x < inv.Cols; x++ {
				sx := startX + x*spacing
				sy := startY + y*spacing

				if cx >= sx && cx < sx+32 && cy >= sy && cy < sy+32 {
					slot := inv.Slots[y][x]

					if inv.DraggingItem == nil && slot.Item != nil {
						// pick up item
						inv.DraggingItem = slot.Item
						inv.DragOffsetX = cx - sx
						inv.DragOffsetY = cy - sy
						slot.Item = nil
					} else if inv.DraggingItem != nil {
						// drop or swap
						inv.DraggingItem, slot.Item = slot.Item, inv.DraggingItem
					}
					break
				}
			}
		}
	}

	inv.prevMousePressed = mousePressed // store state for next frame
}

func (inv *Inventory) AddItem(newItem *Item) bool {
	// Try to stack
	for y := 0; y < inv.Rows; y++ {
		for x := 0; x < inv.Cols; x++ {
			slot := inv.Slots[y][x]
			if slot.Item != nil &&
				slot.Item.ID == newItem.ID &&
				slot.Item.Quantity < slot.Item.MaxStack {
				slot.Item.Quantity += newItem.Quantity
				if slot.Item.Quantity > slot.Item.MaxStack {
					slot.Item.Quantity = slot.Item.MaxStack
				}
				return true
			}
		}
	}

	// Try to find empty
	for y := 0; y < inv.Rows; y++ {
		for x := 0; x < inv.Cols; x++ {
			slot := inv.Slots[y][x]
			if slot.Item == nil {
				slot.Item = newItem
				return true
			}
		}
	}

	return false // no space
}

func (inv *Inventory) DrawTooltip(screen *ebiten.Image) {
	if !inv.IsOpen {
		return
	}

	cx, cy := ebiten.CursorPosition()
	startX, startY := 50, 50
	spacing := 36
	face := basicfont.Face7x13

	for y := 0; y < inv.Rows; y++ {
		for x := 0; x < inv.Cols; x++ {
			sx := startX + x*spacing
			sy := startY + y*spacing

			if cx >= sx && cx < sx+32 && cy >= sy && cy < sy+32 {
				slot := inv.Slots[y][x]
				if slot.Item == nil {
					return
				}

				name := slot.Item.Name
				descLines := wrapText(slot.Item.Description, face, 140)

				// Measure tooltip size
				width := 150
				height := 20 + len(descLines)*14

				tooltip := ebiten.NewImage(width, height)
				tooltip.Fill(color.RGBA{0, 0, 0, 200})

				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(cx+10), float64(cy+10))
				screen.DrawImage(tooltip, op)

				// Draw item name
				text.Draw(screen, name, face, cx+14, cy+24, color.White)

				// Draw description
				for i, line := range descLines {
					text.Draw(screen, line, face, cx+14, cy+38+i*14, color.Gray{Y: 180})
				}
				return
			}
		}
	}
}

// Utility: Random test item generator
func CreateTestItem() *Item {
	rand.Seed(time.Now().UnixNano())
	id := fmt.Sprintf("potion-%d", rand.Intn(3))
	icon := ebiten.NewImage(32, 32)

	switch id {
	case "potion-0":
		icon.Fill(color.RGBA{200, 0, 0, 255}) // red
	case "potion-1":
		icon.Fill(color.RGBA{0, 200, 0, 255}) // green
	case "potion-2":
		icon.Fill(color.RGBA{0, 0, 200, 255}) // blue
	}

	return &Item{
		ID:          id,
		Name:        "Health Potion",
		Icon:        icon,
		Quantity:    1,
		MaxStack:    5,
		Description: "Restores 50 HP",
	}
}

func wrapText(txt string, face font.Face, maxWidth int) []string {
	words := strings.Fields(txt)
	if len(words) == 0 {
		return []string{}
	}

	var lines []string
	var line string

	for _, word := range words {
		testLine := line
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		w := text.BoundString(face, testLine).Dx()
		if w > maxWidth {
			lines = append(lines, line)
			line = word
		} else {
			line = testLine
		}
	}
	if line != "" {
		lines = append(lines, line)
	}

	return lines
}
