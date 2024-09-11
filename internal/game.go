package internal

import (
	"fmt"
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/stollenaar/gamblingjam/util"
)

type Place string

const (
	HOME  Place = "HOME"
	STORE Place = "STORE"
)

type Effect struct {
	name     string
	duration float64
	active   bool
}

type Animation struct {
	image       *ebiten.Image
	drawOptions *ebiten.DrawImageOptions
	startTime   time.Time
	effect      *Effect
	isBlocking  bool
}

type Circle struct {
	x, y, r float32
}
type Game struct {
	currentPlace       Place
	isGameOver         bool
	isDrawingNewsPaper bool

	stats *Stats
	input util.InputHandler

	inventorySlotLocations []Circle
	activeAnimations       []*Animation
	activeButtons          []*util.Button

	startX, startY, startWidth, startHeight, startRot float64

	kbHandler *util.KeyboardHandler
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
		&Member{
			alive:  true,
			sick:   false,
			name:   "Son",
			health: 10,
			hunger: 10,
			warmth: 10,
		},
		&Member{
			alive:  true,
			sick:   false,
			name:   "Daugther",
			health: 10,
			hunger: 10,
			warmth: 10,
		},
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
		currentPlace:       HOME,
		isDrawingNewsPaper: true,
		startX:             73.2,
		startY:             7.7,
		startWidth:         251.9,
		startHeight:        283.1,
		startRot:           45.3,
		stats: &Stats{
			money:     util.ConfigFile.StartingMoney,
			day:       0,
			inventory: make([]Item, 8),
			home: &Home{
				family: members,
			},
			store: &Store{},
		},
		activeAnimations:       []*Animation{},
		inventorySlotLocations: circles,
		input:                  input,
		kbHandler:              util.NewKBHandler(),
	}, nil
}

// Draw implements ebiten.Game.
func (g *Game) Draw(screen *ebiten.Image) {
	if g.isGameOver {
		util.DrawText(screen, 200, 200, util.SelectActiveColor(true), "GAME OVER", nil)
		return
	}
	if g.isDrawingNewsPaper {
		g.drawNewsPaper(screen)
		return
	}
	// Draw end of day information
	vector.DrawFilledRect(screen, 335, 15, 160, 30, color.RGBA{50, 50, 50, 255}, false)
	util.DrawText(screen, 340, 20, color.White, fmt.Sprintf("DAY %d", g.stats.day+1), nil)
	util.DrawText(screen, 410, 20, color.White, g.stats.time.Format(layoutTime), nil)

	g.drawInventory(screen)
	switch g.currentPlace {
	case HOME:
		g.activeButtons = g.stats.home.Draw(screen, g.stats.money)
	case STORE:
		util.DrawCenteredTextInRect(screen, 340, 489, color.RGBA{50, 50, 50, 255}, color.White, fmt.Sprintf("SAVINGS $%d", g.stats.money))
		g.activeButtons = g.stats.store.Draw(screen)
	}
	for _, animation := range g.activeAnimations {
		screen.DrawImage(animation.image, animation.drawOptions)
	}
}

// Layout implements ebiten.Game.
func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth int, screenHeight int) {
	return util.ConfigFile.ScreenWidth, util.ConfigFile.ScreenHeight
}

// Update implements ebiten.Game.
func (g *Game) Update() error {
	if g.isGameOver {
		return nil
	}
	if g.isDrawingNewsPaper && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		g.isDrawingNewsPaper = false
		g.stats.time = time.Date(0, 0, 0, 7, 0, 0, 0, time.Local)
		return nil
	}

	switch key := g.kbHandler.Read(); key {
	case util.KeyA:
		g.startWidth -= 0.1
	case util.KeyD:
		g.startWidth += 0.1
	case util.KeyW:
		g.startHeight += 0.1
	case util.KeyS:
		g.startHeight -= 0.1
	case util.KeyUp:
		g.startY += 0.1
	case util.KeyDown:
		g.startY -= 0.1
	case util.KeyLeft:
		g.startX -= 0.1
	case util.KeyRight:
		g.startX += 0.1
	case util.KeyF:
		g.startRot -= 0.1
	case util.KeyG:
		g.startRot += 0.1
	}

	// fmt.Printf("startWidth: %v, startHeight: %v, startX: %v, startY: %v, startRot: %v\n", g.startWidth, g.startHeight, g.startX, g.startY, g.startRot)
	var isBlocking bool
	var afterLoop []*Animation
	for _, animation := range g.activeAnimations {
		if time.Since(animation.startTime).Seconds() >= animation.effect.duration && !animation.effect.active {
			animation.effect.active = true
		}
		if animation.isBlocking {
			isBlocking = true
		}
		if animation.effect.active {
			switch animation.effect.name {
			case "fadeout":
				animation.drawOptions.ColorScale.Scale(1, 1, 1, animation.drawOptions.ColorScale.A()-0.02)
				if animation.drawOptions.ColorScale.A() == 0 {
					continue
				}
			}
		}
		afterLoop = append(afterLoop, animation)
	}
	g.activeAnimations = afterLoop
	if !isBlocking {
		g.HandleButtons()
	}
	if g.stats.money < 3 && g.hasEmptyInventory() {
		// Game over
		g.isGameOver = true
	}

	return nil
}

func (g *Game) hasEmptyInventory() bool {
	for _, i := range g.stats.inventory {
		if i != nil {
			return false
		}
	}
	return true
}

func (g *Game) HandleButtons() bool {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if button := g.isMouseWithinButtons(); button != nil {
			switch g.currentPlace {
			case HOME:
				switch button.Name {
				case "STORE":
					g.currentPlace = STORE
					g.stats.advanceTime(time.Hour)
				case "SLEEP":
					if g.stats.time.Hour() >= 19 {
						g.stats.day++ // Advance to the next day
						g.stats.home.UpdateMembers()
						g.stats.money -= g.stats.home.totalBills
						g.stats.time = time.Date(0, 0, 0, 7, 0, 0, 0, time.Local)
					}
				case "food":
					g.stats.home.isBuyingFood = !g.stats.home.isBuyingFood
				case "heating":
					g.stats.home.isBuyingHeating = !g.stats.home.isBuyingHeating
				case "Wife":
					g.stats.home.family[0].isBuyingMedicine = !g.stats.home.family[0].isBuyingMedicine
				case "Son":
					g.stats.home.family[1].isBuyingMedicine = !g.stats.home.family[1].isBuyingMedicine
				case "Daugther":
					g.stats.home.family[2].isBuyingMedicine = !g.stats.home.family[2].isBuyingMedicine
				default:
					return false
				}
				return true
			case STORE:
				if button.Name == "HOME" {
					g.currentPlace = HOME
					g.stats.advanceTime(time.Hour)
					return true
				}
			}
		}

		switch {
		case g.stats.HandleButtons(g.currentPlace):
			return true
		case g.findInventorySlot():
			return true
		}
	}
	return false
}

func (g *Game) isMouseWithinButtons() *util.Button {
	for _, button := range g.activeButtons {
		if button.Width != 0 && button.Height != 0 && g.isMouseWithinSquareButton(button) {
			return button
		} else if button.Radius != 0 && g.isMouseWithinCircleButton(button) {
			return button
		}
	}
	return nil
}

// Check if the mouse click is within the button's bounds
func (g *Game) isMouseWithinCircleButton(button *util.Button) bool {
	mouseX, mouseY := ebiten.CursorPosition()
	return distance(float32(mouseX), float32(mouseY), button.X, button.Y) < button.Radius
}

// Check if the mouse click is within the button's bounds
func (g *Game) isMouseWithinSquareButton(button *util.Button) bool {
	mouseX, mouseY := ebiten.CursorPosition()
	return float32(mouseX) >= button.X && float32(mouseX) <= button.X+button.Width &&
		float32(mouseY) >= button.Y && float32(mouseY) <= button.Y+button.Height
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
		if ticket := g.stats.CheckTicket(slot); ticket != nil {
			asset := ebiten.NewImage(winnerTicketAsset.Bounds().Size().X, winnerTicketAsset.Bounds().Size().Y)
			ticketAsset := TicketAssets[ticket.Name]

			doptions := &ebiten.DrawImageOptions{}
			doptions.GeoM.Scale(251.9/float64(ticketAsset.Bounds().Dx()), 283.1/float64(ticketAsset.Bounds().Dy()))
			// Move the image's center to the origin for rotation
			centerX, centerY := float64(asset.Bounds().Dx())/2, float64(asset.Bounds().Dy())/2
			doptions.GeoM.Translate(-centerX, -centerY)
			// Rotate the image anti-clockwise by 90 degrees (π/2 radians)
			doptions.GeoM.Rotate(-(45.3 * math.Pi / 180))
			// Move the image back from the origin to its original position plus any desired offset
			doptions.GeoM.Translate(73.2+centerX, 7.7+centerY)

			asset.DrawImage(ticketAsset, doptions)
			asset.DrawImage(winnerTicketAsset, &ebiten.DrawImageOptions{})

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(200, 100)
			g.activeAnimations = append(g.activeAnimations, &Animation{
				image:       asset,
				drawOptions: op,
				startTime:   time.Now(),
				isBlocking:  true,
				effect: &Effect{
					name:     "fadeout",
					duration: 4,
				},
			})
		}
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

func (g *Game) DoGameLoop() {
	for {
		time.Sleep(time.Second * 5)
		if !g.isDrawingNewsPaper {
			g.stats.advanceTime(time.Minute * 10)
		}
	}
}
