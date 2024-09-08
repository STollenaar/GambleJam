package internal

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	dateX, dateY         = float64(260), float64(90)
	headlineX, headlineY = float64(155), float64(120)
	article1X, article1Y = float64(25), float64(160)
	article2X, article2Y = float64(470), float64(285)
	article3X, article3Y = float64(585), float64(120)
)

var (
	newsPaper *ebiten.Image

	isDrawingNewsPaper bool
)

func init() {
	igbm, _, err := ebitenutil.NewImageFromFile("./assets/newspaper_blank.png")
	if err != nil {
		log.Fatal(err)
	}
	newsPaper = igbm
}

func drawNewsPaper(screen *ebiten.Image) {

}