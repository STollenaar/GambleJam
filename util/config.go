package util

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Button struct {
	Name   string
	X      float32 `json:"x"`
	Y      float32 `json:"y"`
	Width  float32 `json:"width"`
	Height float32 `json:"height"`
}

type Config struct {
	StartingMoney int     `json:"startingMoney"`
	ScreenWidth   int     `json:"screenWidth"`
	ScreenHeight  int     `json:"screenHeight"`
	Margin        float32 `json:"margin"`

	Buttons map[string]*Button `json:"buttons"`
}

var (
	ConfigFile *Config

	DefaultFont *text.GoTextFace
	TitleFont *text.GoTextFace
)

func init() {
	ConfigFile = new(Config)
	data, err := os.ReadFile("./configs/config.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(data, ConfigFile)
	if err != nil {
		log.Fatal(err)
	}

			// Define font
			s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
			if err != nil {
				log.Fatal(err)
			}
	
			DefaultFont = &text.GoTextFace{
				Source: s,
				Size:   15,
			}
			TitleFont = &text.GoTextFace{
				Source: s,
				Size:   30,
			}
	
}
