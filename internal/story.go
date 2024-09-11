package internal

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/stollenaar/gamblingjam/util"
)

const (
	dateX, dateY         = float64(260), float64(82)
	headlineX, headlineY = float64(155), float64(110)
	article1X, article1Y = float64(25), float64(160)
	article2X, article2Y = float64(395), float64(280)
	article3X, article3Y = float64(532), float64(115)
)

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Newspaper struct {
	Template *ebiten.Image

	Headline Point `json:"headline"`

	Article1 Point `json:"article1"`
	Article2 Point `json:"article2"`
	Article3 Point `json:"article3"`
}

type Story struct {
	NewspaperTemplate string `json:"newspaperTemplate"`
	Headline string `json:"headline"`
	Article1 string `json:"article1"`
	Article2 string `json:"article2"`
	Article3 string `json:"article3"`
}

var (
	fadedColor = color.RGBA{96, 95, 88, 255}

	newspapers map[string]Newspaper
	story map[string]Story
)

func init() {
	newspapers = make(map[string]Newspaper)
	story = make(map[string]Story)

	data, err := os.ReadFile("./configs/newspapers.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, &newspapers)
	if err != nil {
		log.Fatal(err)
	}

	data, err = os.ReadFile("./configs/story.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, &story)
	if err != nil {
		log.Fatal(err)
	}

	assets, err := os.ReadDir("./assets/newspapers")
	if err != nil {
		log.Fatal(err)
	}
	for _, asset := range assets {
		igbm, _, err := ebitenutil.NewImageFromFile(fmt.Sprintf("./assets/newspapers/%s", asset.Name()))
		if err != nil {
			log.Fatal(err)
		}
		nameNorm := strings.Title(strings.ReplaceAll(strings.ReplaceAll(asset.Name(), "_", " "), ".png", ""))
		tpl := newspapers[nameNorm]
		tpl.Template = igbm
		newspapers[nameNorm] = tpl
	}
}

func (g *Game) drawNewsPaper(screen *ebiten.Image) {
	newspaper := newspapers["1"]
	today := story["1"]

	paper := ebiten.NewImageFromImage(newspaper.Template)

	util.DrawText(paper, dateX, dateY, fadedColor, "September", nil)

	util.DrawText(paper, newspaper.Headline.X, newspaper.Headline.Y, color.Black, today.Headline, util.TitleFont)
	util.DrawText(paper, newspaper.Article1.X, newspaper.Article1.Y, color.Black, today.Article1, nil)
	util.DrawText(paper, newspaper.Article2.X, newspaper.Article2.Y, color.Black, today.Article2, nil)
	util.DrawText(paper, newspaper.Article3.X, newspaper.Article3.Y, color.Black, today.Article3, nil)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(util.ConfigFile.ScreenWidth)/2-float64(paper.Bounds().Dx())/2, float64(util.ConfigFile.ScreenHeight)/2-float64(paper.Bounds().Dy())/2)
	screen.DrawImage(paper, op)
}
