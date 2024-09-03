package internal

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/stollenaar/gamblingjam/util"
)

type Home struct {
	family      []*Member
	hasMedicine bool
}

func (h *Home) UpdateMembers() {
	for _, member := range h.family {
		member.doUpdate()
	}
}

func (h *Home) Draw(screen *ebiten.Image) {
	h.drawMemberStats(screen)
	h.drawSleepButton(screen)
	h.drawStoreButton(screen)
}

func (h *Home) drawMemberStats(screen *ebiten.Image) {
	for i, member := range h.family {
		text := fmt.Sprintf("%s", member.name)
		conditions := []string{}
		if member.hunger <= 5 {
			conditions = append(conditions, "Hungry")
		}
		if member.warmth <= 2 {
			conditions = append(conditions, "Froze da det")
		} else if member.warmth <= 5 {
			conditions = append(conditions, "Cold")
		} else if member.warmth <= 7 {
			conditions = append(conditions, "Shivering")
		}
		if member.sick {
			conditions = append(conditions, "Sick")
		}
		if !member.alive {
			conditions = []string{"Deceased"}
		}
		if len(conditions) > 0 {
			text += ": " + strings.Join(conditions, ", ")
		}
		ebitenutil.DebugPrintAt(screen, text, 100+(i*200), 80)
	}
}

// Draw the sleep button
func (h *Home) drawSleepButton(screen *ebiten.Image) {
	buttonColor := color.RGBA{50, 50, 50, 255} // Dark grey color for the button
	sleepButton := util.ConfigFile.Buttons["sleep"]
	vector.DrawFilledRect(screen, float32(sleepButton.X), float32(sleepButton.Y), float32(sleepButton.Width), float32(sleepButton.Height), buttonColor, false)
	ebitenutil.DebugPrintAt(screen, "SLEEP", sleepButton.X+10, sleepButton.Y+10) // Text centered within the button
}

// Draw the store button
func (h *Home) drawStoreButton(screen *ebiten.Image) {
	buttonColor := color.RGBA{50, 50, 50, 255} // Dark grey color for the button
	storeButton := util.ConfigFile.Buttons["store"]
	vector.DrawFilledRect(screen, float32(storeButton.X), float32(storeButton.Y), float32(storeButton.Width), float32(storeButton.Height), buttonColor, false)
	ebitenutil.DebugPrintAt(screen, "STORE", storeButton.X+10, storeButton.Y+10) // Text centered within the button
}
