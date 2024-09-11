package internal

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/stollenaar/gamblingjam/util"
)

type Home struct {
	family     []*Member
	totalBills int

	isBuyingFood    bool
	isBuyingHeating bool
}

func (h *Home) UpdateMembers() {
	for _, member := range h.family {
		member.doUpdate(h.isBuyingFood, h.isBuyingHeating)
	}
}

func (h *Home) Draw(screen *ebiten.Image, savings int, time time.Time) (activeButtons []*util.Button) {
	h.drawMemberStats(screen)
	activeButtons = append(activeButtons, h.drawUtilities(screen, savings)...)
	activeButtons = append(activeButtons, h.drawSleepButton(screen), h.drawStoreButton(screen, time))
	return activeButtons
}

func (h *Home) drawMemberStats(screen *ebiten.Image) {
	for i, member := range h.family {
		conditions := []string{}
		if member.hunger <= 5 {
			conditions = append(conditions, "Hungry")
		}
		if member.warmth <= 2 {
			conditions = append(conditions, "Froze")
		} else if member.warmth <= 5 {
			conditions = append(conditions, "Cold")
		} else if member.warmth <= 7 {
			conditions = append(conditions, "Chilly")
		}
		if member.sick {
			conditions = append(conditions, "Sick")
		}
		if !member.alive {
			conditions = []string{"Dead"}
		}

		// Get the text size
		textWidth, _ := text.Measure(member.name, util.DefaultFont, 0)

		// Draw the member name text at (x, y)
		x := 100 + (i * 200)
		y := 140
		circleRadius := 29.0

		totalCircleWidth := float32(len(conditions)-1) * float32(circleRadius*2)

		textX := float32(x) - float32(textWidth)/2

		if len(conditions) > 1 {
			textX += (totalCircleWidth / 2)
		}

		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(textX), float64(y))
		op.ColorScale.ScaleWithColor(color.White)
		text.Draw(screen, member.name, util.DefaultFont, op)

		for i, condition := range conditions {
			// Measure the size of the condition text
			conditionWidth, conditionHeight := text.Measure(condition, util.DefaultFont, 0)

			circleX := x
			circleY := float32(y) - float32(circleRadius) - float32(10) // Position circles above the text

			circleColor := color.RGBA{50, 50, 50, 255}
			if condition == "Dead" {
				circleColor = color.RGBA{204, 46, 46, 255}
			}
			vector.DrawFilledCircle(screen, float32(circleX)+float32(i)*(float32(circleRadius)*2), circleY, float32(circleRadius), circleColor, false)

			// Calculate position to center the text inside the circle
			textX := float64(circleX) + float64(i)*(circleRadius*2) - conditionWidth/2
			textY := float64(circleY) - conditionHeight/2

			// Draw the condition text in the center of the circle
			op := &text.DrawOptions{}
			op.GeoM.Translate(textX, textY)
			op.ColorScale.ScaleWithColor(color.White)
			text.Draw(screen, condition, util.DefaultFont, op)
		}
		i++
	}
}

func (h *Home) drawUtilities(screen *ebiten.Image, money int) (activeButtons []*util.Button) {
	var totalBills int

	moneyText := util.DrawText(screen, 300, 230, color.White, "SAVINGS:", nil)
	util.DrawText(screen, 500, float64(moneyText.Y)-22, color.White, fmt.Sprintf("$%d", money), nil)

	mortgageText := util.DrawText(screen, 300, float64(moneyText.Y+moneyText.Height)+2, color.RGBA{204, 46, 46, 255}, "MORTGAGE:", nil)
	util.DrawText(screen, 500, float64(mortgageText.Y)-22, color.RGBA{204, 46, 46, 255}, fmt.Sprintf("- $%d", 80), nil)
	totalBills += 80

	foodText := util.DrawText(screen, 300, float64(mortgageText.Y+mortgageText.Height)+2, util.SelectActiveColor(h.isBuyingFood), "FOOD:", nil)
	foodAmount := util.DrawText(screen, 500, float64(foodText.Y)-22, util.SelectActiveColor(h.isBuyingFood), fmt.Sprintf("- $%d", 30), nil)
	vector.DrawFilledCircle(screen, foodAmount.X+foodAmount.Width+50, foodAmount.Y-20+foodAmount.Height/2, float32(10), util.SelectActiveColor(h.isBuyingFood), false)
	activeButtons = append(activeButtons, &util.Button{
		Name:   "food",
		X:      foodAmount.X + foodAmount.Width + 50,
		Y:      foodAmount.Y - 20 + foodAmount.Height/2,
		Radius: 10,
	})

	if h.isBuyingFood {
		totalBills += 30
	}

	heatingText := util.DrawText(screen, 300, float64(foodText.Y+foodText.Height)+2, util.SelectActiveColor(h.isBuyingHeating), "HEATING:", nil)
	heatingAmount := util.DrawText(screen, 500, float64(heatingText.Y)-22, util.SelectActiveColor(h.isBuyingHeating), fmt.Sprintf("- $%d", 30), nil)
	vector.DrawFilledCircle(screen, heatingAmount.X+heatingAmount.Width+50, heatingAmount.Y-20+heatingAmount.Height/2, float32(10), util.SelectActiveColor(h.isBuyingHeating), false)
	activeButtons = append(activeButtons, &util.Button{
		Name:   "heating",
		X:      heatingAmount.X + heatingAmount.Width + 50,
		Y:      heatingAmount.Y - 20 + heatingAmount.Height/2,
		Radius: 10,
	})

	if h.isBuyingHeating {
		totalBills += 30
	}

	var offset int
	previousText := heatingText
	for _, member := range h.family {
		if member.alive && member.sick {

			previousText = util.DrawText(screen, 300, float64(previousText.Y+previousText.Height)+2, util.SelectActiveColor(member.isBuyingMedicine), fmt.Sprintf("MEDICINE %s:", member.name), nil)
			previousAmount := util.DrawText(screen, 500, float64(previousText.Y)-22, util.SelectActiveColor(member.isBuyingMedicine), fmt.Sprintf("- $%d", 10), nil)
			vector.DrawFilledCircle(screen, previousAmount.X+previousAmount.Width+50, previousAmount.Y-20+previousAmount.Height/2, float32(10), util.SelectActiveColor(member.isBuyingMedicine), false)
			activeButtons = append(activeButtons, &util.Button{
				Name:   member.name,
				X:      previousAmount.X + previousAmount.Width + 50,
				Y:      previousAmount.Y - 20 + previousAmount.Height/2,
				Radius: 10,
			})

			offset += 1
			if member.isBuyingMedicine {
				totalBills += 10
			}
		}
	}
	remainingText := util.DrawText(screen, 300, float64(previousText.Y+previousText.Height+25), color.White, "REMAINING SAVINGS:", nil)
	util.DrawText(screen, 500, float64(remainingText.Y)-22, color.White, fmt.Sprintf("$%d", money-totalBills), nil)
	h.totalBills = totalBills
	return activeButtons
}

// Draw the sleep button
func (h *Home) drawSleepButton(screen *ebiten.Image) *util.Button {
	buttonColor := color.RGBA{50, 50, 50, 255} // Dark grey color for the button
	sleepButton := util.ConfigFile.Buttons["sleep"]
	return util.DrawCenteredTextInRect(screen, float32(sleepButton.X), float32(sleepButton.Y), buttonColor, color.White, "SLEEP")
}

// Draw the store button
func (h *Home) drawStoreButton(screen *ebiten.Image, time time.Time) *util.Button {
	buttonColor := color.RGBA{50, 50, 50, 255} // Dark grey color for the button
	if time.Hour() >= 21 {
		buttonColor = util.SelectActiveColor(true)
	}
	storeButton := util.ConfigFile.Buttons["store"]
	button := util.DrawCenteredTextInRect(screen, float32(storeButton.X), float32(storeButton.Y), buttonColor, color.White, "STORE")
	if time.Hour() >= 21 {
		vector.StrokeLine(screen, float32(storeButton.X), float32(storeButton.Y), float32(storeButton.X)+button.Width, float32(storeButton.Y)+button.Height, 2.5, color.Black, false)
		vector.StrokeLine(screen, float32(storeButton.X), float32(storeButton.Y)+button.Height, float32(storeButton.X)+button.Width, float32(storeButton.Y), 2.5, color.Black, false)
	}
	return button
}
