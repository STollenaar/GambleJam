package internal

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/stollenaar/gamblingjam/util"
)

type Store struct {
	cart []*Ticket
}

func (s *Store) Draw(screen *ebiten.Image) {
	s.drawDisplayCase(screen)
	s.drawHomeButton(screen)
	s.drawPlayButton(screen)
}

// Draw tickets in store
func (s *Store) drawDisplayCase(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, 10.0, 50.0, 660.0, 420.0, color.RGBA{255, 255, 255, 255}, false)
	util.DrawClearRectangle(screen, color.RGBA{213, 43, 30, 255}, 10, 50, 660, 420, 10)
	var xoff, yoff int
	for _, name := range TicketNames {
		ticket := TicketAssets[name]
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(128.0/800.0, 200.0/900.0)
		op.GeoM.Translate(float64(20+(128*xoff)), float64(60+(200*yoff)))
		screen.DrawImage(ticket, &ebiten.DrawImageOptions{
			GeoM: op.GeoM,
		})
		xoff++
		if xoff == 5 {
			xoff = 0
			yoff++
		}
	}
}

// Draw the play button
func (s *Store) drawPlayButton(screen *ebiten.Image) {
	buttonColor := color.RGBA{50, 50, 50, 255} // Dark grey color for the button
	sleepButton := util.ConfigFile.Buttons["sleep"]
	vector.DrawFilledRect(screen, float32(sleepButton.X), float32(sleepButton.Y), float32(sleepButton.Width), float32(sleepButton.Height), buttonColor, false)
	ebitenutil.DebugPrintAt(screen, "PLAY", sleepButton.X+10, sleepButton.Y+10) // Text centered within the button
}

// Draw the home button
func (s *Store) drawHomeButton(screen *ebiten.Image) {
	buttonColor := color.RGBA{50, 50, 50, 255} // Dark grey color for the button
	storeButton := util.ConfigFile.Buttons["store"]
	vector.DrawFilledRect(screen, float32(storeButton.X), float32(storeButton.Y), float32(storeButton.Width), float32(storeButton.Height), buttonColor, false)
	ebitenutil.DebugPrintAt(screen, "HOME", storeButton.X+10, storeButton.Y+10) // Text centered within the button
}

func (s *Store) FindTicket() *Ticket {
	mouseX, mouseY := ebiten.CursorPosition()

	mouseX -= 20
	mouseY -= 60
	if mouseX < 0 || mouseY < 0 || mouseX >= 640 || mouseY >= 400 {
		return nil
	}
	xOff := mouseX / 128
	yOff := mouseY / 200

	index := xOff + (yOff * 5)

	return Tickets[index]
}
