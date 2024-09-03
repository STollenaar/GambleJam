package util

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func DrawClearRectangle(screen *ebiten.Image, color color.RGBA, x, y, width, height, borderWidth float32) {
	// Draw top border
	vector.DrawFilledRect(screen, x, y, width, borderWidth, color, false)

	// Draw bottom border
	vector.DrawFilledRect(screen, x, y+height-borderWidth, width, borderWidth, color, false)

	// Draw left border
	vector.DrawFilledRect(screen, x, y, borderWidth, height, color, false)

	// Draw right border
	vector.DrawFilledRect(screen, x+width-borderWidth, y, borderWidth, height, color, false)
}
