package util

import (
	"encoding/json"
	"log"
	"os"
)

type Button struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

type Config struct {
	StartingMoney int `json:"startingMoney"`
	ScreenWidth   int `json:"screenWidth"`
	ScreenHeight  int `json:"screenHeight"`

	Buttons map[string]*Button `json:"buttons"`
}

var (
	ConfigFile *Config
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
}
