package internal

import (
	"fmt"
	"image/color"
	"math"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/stollenaar/gamblingjam/util"
)

type Place string

const (
	HOME  Place = "HOME"
	STORE Place = "STORE"
)

type Circle struct {
	x, y, r float32
}
type Game struct {
	currentPlace Place

	stats *Stats
	input util.InputHandler

	inventorySlotLocations []Circle
}

func NewGame(input util.InputHandler) (*Game, error) {
	op := ebiten.DrawImageOptions{}
	op.ColorM.ChangeHSV(0, 1, .8)
	op.GeoM.Scale(.25, .25)

	members := []*Member{
		&Member{
			alive:  true,
			sick:   false,
			name:   "Wife",
			health: 10,
			hunger: 10,
			warmth: 10,
		},
	}

	r := rand.IntN(4) + 1

	for ; r > 0; r-- {
		members = append(members, &Member{
			alive:  true,
			sick:   false,
			name:   "Son",
			health: 10,
			hunger: 10,
			warmth: 10,
		})
	}
	var circles []Circle
	radius := 30

	for i := 0; i < 8; i++ {
		circles = append(circles, Circle{
			float32(50+i*90) + float32(radius),
			float32(util.ConfigFile.ScreenHeight) - 40,
			float32(radius),
		})
	}

	return &Game{
		currentPlace: STORE,
		stats: &Stats{
			money:     util.ConfigFile.StartingMoney,
			day:       0,
			inventory: make([]Item, 8),
			home: &Home{
				hasMedicine: false,
				family:      members,
			},
			store: &Store{},
		},
		inventorySlotLocations: circles,
		input:                  input,
	}, nil
}

// Draw implements ebiten.Game.
func (g *Game) Draw(screen *ebiten.Image) {
	// Draw end of day information
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("DAY %d", g.stats.day+1), 340, 20)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("SAVINGS $%d", g.stats.money), 340, 500)
	g.drawInventory(screen)
	switch g.currentPlace {
	case HOME:
		g.stats.home.Draw(screen)
	case STORE:
		g.stats.store.Draw(screen)
	}
}

// Layout implements ebiten.Game.
func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return util.ConfigFile.ScreenWidth, util.ConfigFile.ScreenHeight
}

// Update implements ebiten.Game.
func (g *Game) Update() error {
	g.HandleButtons()
	return nil
}

func (g *Game) HandleButtons() bool {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if g.isMouseWithinButton(util.ConfigFile.Buttons["sleep"]) {
			switch g.currentPlace {
			case HOME:
				g.stats.day++ // Advance to the next day
				g.stats.home.UpdateMembers()
			case STORE:
			}
			return true
		} else if g.isMouseWithinButton(util.ConfigFile.Buttons["store"]) {
			switch g.currentPlace {
			case HOME:
				g.currentPlace = STORE
			case STORE:
				g.currentPlace = HOME
			}
			return true
		} else {
			switch {
			case g.stats.HandleButtons(g.currentPlace):
			case g.findInventorySlot():
			}
		}
	}
	return false
}

// Check if the mouse click is within the button's bounds
func (g *Game) isMouseWithinButton(button *util.Button) bool {
	mouseX, mouseY := ebiten.CursorPosition()
	return mouseX >= button.X && mouseX <= button.X+button.Width &&
		mouseY >= button.Y && mouseY <= button.Y+button.Height
}

// Function to draw the inventory section at the bottom of the screen
func (g *Game) drawInventory(screen *ebiten.Image) {
	inventoryColor := color.RGBA{128, 128, 128, 255} // Grey color for inventory background
	slotColor := color.White                         // White color forthe slots

	// Draw the inventory background
	vector.DrawFilledRect(screen, 0, float32(util.ConfigFile.ScreenHeight-80), float32(util.ConfigFile.ScreenWidth), 80, inventoryColor, false)
	// Draw inventory slots (8 slots)
	for i, slot := range g.inventorySlotLocations {
		vector.DrawFilledCircle(screen, slot.x, slot.y, slot.r, slotColor, false)
		if item := g.stats.inventory[i]; item != nil {
			image := TicketAssets[item.(*Ticket).Name]
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(float64(slot.r/800.0), float64((slot.r*2)/900.0))
			op.GeoM.Translate(float64(slot.x-slot.r/2), float64(slot.y-slot.r))
			screen.DrawImage(image, &ebiten.DrawImageOptions{
				GeoM: op.GeoM,
			})
		}
	}
}

func (g *Game) findInventorySlot() bool {
	mouseX, mouseY := ebiten.CursorPosition()

	slot := g.pointInWhichCircle(float32(mouseX), float32(mouseY))
	if slot != -1 {
		g.stats.CheckTicket(slot)
		return true
	}
	return false
}

func (g *Game) pointInWhichCircle(x, y float32) int {
	for i, circle := range g.inventorySlotLocations {
		if distance(x, y, circle.x, circle.y) < circle.r {
			return i
		}
	}
	return -1
}

func distance(x1, y1, x2, y2 float32) float32 {
	return float32(math.Sqrt(float64((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1))))
}
