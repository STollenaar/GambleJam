package internal

import (
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"os"
	"strconv"
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

type Rect struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type Newspaper struct {
	Template *ebiten.Image

	Headline Rect `json:"headline"`

	Article1 Rect `json:"article1"`
	Article2 Rect `json:"article2"`
	Article3 Rect `json:"article3"`
	Article4 Rect `json:"article4"`
}

type Story struct {
	NewspaperTemplate string `json:"newspaperTemplate"`
	Headline          string `json:"headline"`
	Article1          string `json:"article1"`
	Article2          string `json:"article2"`
	Article3          string `json:"article3"`
	Article4          string `json:"article4"`
}

var (
	fadedColor = color.RGBA{96, 95, 88, 255}

	newspapers map[string]Newspaper
	story      map[string]Story
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

	today := story[strconv.Itoa(g.stats.day+1)]
	newspaper := newspapers[today.NewspaperTemplate]

	paper := ebiten.NewImageFromImage(newspaper.Template)

	util.DrawText(paper, dateX, dateY, fadedColor, fmt.Sprintf("%d September 2008", g.stats.day+1), nil)

	util.DrawTextInRect(paper, today.Headline, newspaper.Headline.X, newspaper.Headline.Y, newspaper.Headline.Width, newspaper.Headline.Height, color.Black, util.TitleFont)
	util.DrawTextInRect(paper, today.Article1, newspaper.Article1.X, newspaper.Article1.Y, newspaper.Article1.Width, newspaper.Article1.Height, color.Black, nil)
	util.DrawTextInRect(paper, today.Article2, newspaper.Article2.X, newspaper.Article2.Y, newspaper.Article2.Width, newspaper.Article2.Height, color.Black, nil)
	util.DrawTextInRect(paper, today.Article3, newspaper.Article3.X, newspaper.Article3.Y, newspaper.Article3.Width, newspaper.Article3.Height, color.Black, nil)
	util.DrawTextInRect(paper, today.Article4, newspaper.Article4.X, newspaper.Article4.Y, newspaper.Article4.Width, newspaper.Article4.Height, color.Black, nil)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(util.ConfigFile.ScreenWidth)/2-float64(paper.Bounds().Dx())/2, float64(util.ConfigFile.ScreenHeight)/2-float64(paper.Bounds().Dy())/2)
	screen.DrawImage(paper, op)
}
