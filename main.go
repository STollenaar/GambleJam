package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/stollenaar/gamblingjam/internal"
	"github.com/stollenaar/gamblingjam/util"
)

func main() {
	input := util.NewKBHandler()

	game, err := internal.NewGame(input)
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(util.ConfigFile.ScreenWidth, util.ConfigFile.ScreenHeight)
	ebiten.SetWindowTitle("Scratch Tickets, Please!")
	go game.DoGameLoop()
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
